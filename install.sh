#!/usr/bin/env bash

# Copyright 2018 The Kubeflow Authors
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#       http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

set -e

SCRIPT_DIR="$(cd "$(dirname "$(readlink "$0" || echo "$0")")"; pwd)"

function help() {
    echo -e "
Usage:

    install.sh [OPTION1] [OPTION2] ...

Options:
    --kubeconfig string              Specify the kubeconfig file
    --namespace string               Specify the namespace that operators will be installed in 
    --only-binary                    Only install arena binary
    --region-id string               Specify the region id(it is available in Alibaba Cloud)
    --host-network                   Enable host network
    --docker-registry string         Specify the docker registry
    --registry-repo-namespace string Specify the docker registry repo namespace
    --loadbalancer                   Specify k8s service type with loadbalancer
    --prometheus                     Install prometheus
    --platform string                Specify the platform(eg: ack)
    --rdma                           Enable rdma feature
"

}


function logger() {
    timestr=$(date +"%Y-%m-%d/%H:%M:%S")
    level=$(echo $1 | tr 'a-z' 'A-Z')
    echo ${timestr}"  "${level}"  "$2
}

# if pull images by aliyun vpc,change the image
function support_image_regionalization(){
	for file in $(ls $1);do
		local path=$1"/"$file
		if [ -d $path ];then
			support_image_regionalization $path
		else
        	if [[ $PULL_IMAGE_BY_VPC_NETWORK == "true" ]];then
				sed -i  "s@registry\..*aliyuncs.com@registry-vpc.${REGION}.aliyuncs.com@g" $path
			else
				sed -i  "s@registry\..*aliyuncs.com@registry.${REGION}.aliyuncs.com@g" $path
			fi
		fi
	done
}

# if execute the install.sh needs sudo,add sudo command to all command
function set_sudo() {
    export sudo_prefix=""
    if [ `id -u` -ne 0 ]; then
        export sudo_prefix="sudo"
    fi  
}

# install kubectl and rename arena-kubectl
# install helm and rename arena-helm 
function install_kubectl_and_helm() {
    logger "debug" "start to install arena-kubectl and arena-helm"
    ${sudo_prefix} rm -rf /usr/local/bin/arena-kubectl
    ${sudo_prefix} rm -rf /usr/local/bin/arena-helm

    ${sudo_prefix} cp $SCRIPT_DIR/bin/kubectl /usr/local/bin/arena-kubectl
    ${sudo_prefix} cp $SCRIPT_DIR/bin/helm /usr/local/bin/arena-helm

    if ! ${sudo_prefix} arena-kubectl cluster-info >/dev/null 2>&1; then
        logger "error" "failed to execute 'arena-kubectl cluster-info'"
        logger "error" "Please setup kubeconfig correctly before installing arena"
        exit 1
    fi
    logger "debug" "succeed to install arena-kubectl and arena-helm"
}

function custom_charts() {
    if [[ $REGION != "" ]];then
        logger "debug" "enable image regionalization"
	    support_image_regionalization $SCRIPT_DIR/charts
    fi

    if [ "$USE_HOSTNETWORK" == "true" ]; then
        logger "debug" "enable host network mode"
        find $SCRIPT_DIR/charts/ -name values.yaml | xargs sed -i "/useHostNetwork/s/false/true/g"
    fi

    if [[ ${DOCKER_REGISTRY} != "" ]]; then
        logger "debug" "custom the docker registry with ${DOCKER_REGISTRY}"
        find $SCRIPT_DIR/charts/ -name *.yaml | xargs sed -i "s/registry.cn-zhangjiakou.aliyuncs.com/${DOCKER_REGISTRY}/g"
        find $SCRIPT_DIR/charts/ -name *.yaml | xargs sed -i "s/registry.cn-hangzhou.aliyuncs.com/${DOCKER_REGISTRY}/g"
    fi 

    if [[ ${REGISTRY_REPO_NAMESPACE} != "" ]]; then
        logger "debug" "custom the docker registry repo namespace with ${REGISTRY_REPO_NAMESPACE}"
        find $SCRIPT_DIR/charts/ -name *.yaml | xargs sed -i "s/tensorflow-samples/${REGISTRY_REPO_NAMESPACE}/g"
    fi

    if [ "$USE_LOADBALANCER" == "true" ]; then
        logger "debug" "specify service with loadbalancer type"
        find $SCRIPT_DIR/charts/ -name *.yaml | xargs sed -i "s/NodePort/LoadBalancer/g"
    fi

    if [[ $USE_RDMA == "true" ]];then
        find $SCRIPT_DIR/charts/ -name *.yaml | xargs sed -i "/enableRDMA/s/false/true/g" 
    fi
}

function custom_manifests() {
    if [[ $REGION != "" ]];then
        support_image_regionalization $SCRIPT_DIR/kubernetes-artifacts
    fi

    if [[ ${DOCKER_REGISTRY} != "" ]]; then
        find $SCRIPT_DIR/kubernetes-artifacts/ -name *.yaml | xargs sed -i "s/registry.cn-zhangjiakou.aliyuncs.com/${DOCKER_REGISTRY}/g"
        find $SCRIPT_DIR/kubernetes-artifacts/ -name *.yaml | xargs sed -i "s/registry.cn-hangzhou.aliyuncs.com/${DOCKER_REGISTRY}/g"
    fi

    if [[ "${NAMESPACE}" != "" ]]; then
        logger "debug" "custom the namespace(${NAMESPACE}) which operators will be installed in"
        find $SCRIPT_DIR/kubernetes-artifacts/ -name *.yaml | xargs sed -i "s/arena-system/${NAMESPACE}/g"
    fi

    if [[ ${REGISTRY_REPO_NAMESPACE} != "" ]]; then
        find $SCRIPT_DIR/kubernetes-artifacts/ -name *.yaml | xargs sed -i "s/tensorflow-samples/${REGISTRY_REPO_NAMESPACE}/g"
    fi

    if [ "$USE_LOADBALANCER" == "true" ]; then
        find $SCRIPT_DIR/kubernetes-artifacts/ -name *.yaml | xargs sed -i "s/NodePort/LoadBalancer/g"
    fi

    if [ "$PLATFORM" == "ack" ]; then
        sed -i 's|accelerator/nvidia_gpu|aliyun.accelerator/nvidia_count|g' $SCRIPT_DIR/kubernetes-artifacts/prometheus/gpu-exporter.yaml
    fi
}

# install arena command tool and charts
function install_arena_and_charts() {
    now=$(date "+%Y%m%d%H%M%S")
    if [ -f "/usr/local/bin/arena" ]; then
        ${sudo_prefix} cp /usr/local/bin/arena /usr/local/bin/arena-$now
    fi
    ${sudo_prefix} cp $SCRIPT_DIR/bin/arena /usr/local/bin/arena
    if [ -f /usr/local/bin/arena-uninstall ];then
        ${sudo_prefix} rm -rf /usr/local/bin/arena-uninstall
    fi
    ${sudo_prefix} cp $SCRIPT_DIR/bin/arena-uninstall /usr/local/bin
    if [ ! -d $SCRIPT_DIR/charts/kubernetes-artifacts ];then 
        cp -a $SCRIPT_DIR/kubernetes-artifacts $SCRIPT_DIR/charts
    fi 
    # For non-root user, put the charts dir to the home directory
    if [ `id -u` -eq 0 ];then  
        if [ -d "/charts" ]; then
           mv /charts /charts-$now
        fi
        cp -r $SCRIPT_DIR/charts / 
    else  
        if [ -d "~/charts" ]; then
          mv ~/charts ~/charts-$now
        fi
        cp -r $SCRIPT_DIR/charts ~/  
    fi   
}

function install_arena_gen_kubeconfig() {
    ${sudo_prefix} rm -rf /usr/local/bin/arena-gen-kubeconfig.sh
    ${sudo_prefix} cp $SCRIPT_DIR/bin/arena-gen-kubeconfig.sh /usr/local/bin/arena-gen-kubeconfig.sh
}


function apply_jobmon() {
    if ! arena-kubectl get serviceaccount --all-namespaces | grep jobmon; then
        arena-kubectl apply -f $SCRIPT_DIR/kubernetes-artifacts/jobmon/jobmon-role.yaml
    fi
}


function apply_tf() {
    # if KubeDL has installed, will skip install some CRDs of framework operator
    if arena-kubectl get serviceaccount --all-namespaces | grep kubedl; then
        logger "warning" "KubeDL has been detected, will skip install tf-operator"
        return 
    fi
    if ! arena-kubectl get serviceaccount --all-namespaces | grep tf-job-operator; then
        arena-kubectl apply -f ${SCRIPT_DIR}/kubernetes-artifacts/tf-operator/tf-crd.yaml
        arena-kubectl apply -f ${SCRIPT_DIR}/kubernetes-artifacts/tf-operator/tf-operator.yaml
        return 
    fi
    if arena-kubectl get crd tfjobs.kubeflow.org -oyaml |grep -i 'version: v1alpha2'; then
        arena-kubectl delete -f ${SCRIPT_DIR}/kubernetes-artifacts/tf-operator/tf-operator-v1alpha2.yaml
        arena-kubectl apply -f ${SCRIPT_DIR}/kubernetes-artifacts/tf-operator/tf-crd.yaml
        arena-kubectl apply -f ${SCRIPT_DIR}/kubernetes-artifacts/tf-operator/tf-operator.yaml
        return 
    fi 

}

# TODO: the pytorch-operator update
function apply_pytorch() {
    if arena-kubectl get serviceaccount --all-namespaces | grep kubedl; then
        logger "warning" "KubeDL has been detected, will skip install pytorch-operator"
        return
    fi
    if ! arena-kubectl get serviceaccount --all-namespaces | grep pytorch-operator; then
        arena-kubectl apply -f $SCRIPT_DIR/kubernetes-artifacts/pytorch-operator/pytorch-operator.yaml
        return 
    fi    
}

function apply_mpi() {
    if arena-kubectl get serviceaccount --all-namespaces | grep kubedl; then
        logger "warning" "KubeDL has been detected, will skip install mpi-operator"
        return 
    fi
    if ! arena-kubectl get serviceaccount --all-namespaces | grep mpi-operator; then
        arena-kubectl apply -f $SCRIPT_DIR/kubernetes-artifacts/mpi-operator/mpi-operator.yaml
        return 
    fi
}

function apply_et() { 
    if ! arena-kubectl get serviceaccount --all-namespaces | grep et-operator; then
        arena-kubectl apply -f $SCRIPT_DIR/kubernetes-artifacts/et-operator/et-operator.yaml
        return 
    fi    
}

function apply_prometheus() {
    if [[ "$USE_PROMETHEUS" != "true" ]];then
        return 
    fi  
    if ! arena-kubectl get serviceaccount --all-namespaces | grep prometheus; then
        arena-kubectl apply -f $SCRIPT_DIR/kubernetes-artifacts/prometheus/gpu-exporter.yaml
        arena-kubectl apply -f $SCRIPT_DIR/kubernetes-artifacts/prometheus/prometheus.yaml
        arena-kubectl apply -f $SCRIPT_DIR/kubernetes-artifacts/prometheus/grafana.yaml
        return 
    fi
}

function apply_rdma() {
    if [[ $USE_RDMA != "true" ]];then
        return 
    fi
    arena-kubectl apply -f $SCRIPT_DIR/kubernetes-artifacts/rdma/rdma-config.yaml   
    arena-kubectl apply -f $SCRIPT_DIR/kubernetes-artifacts/rdma/device-plugin.yaml   
}

function apply_kubedl() {
  if ! arena-kubectl get serviceaccount --all-namespaces | grep kubedl-operator; then
    arena-kubectl apply -f $SCRIPT_DIR/kubernetes-artifacts/kubedl/kubedl-crd.yaml
    arena-kubectl apply -f $SCRIPT_DIR/kubernetes-artifacts/kubedl/kubedl-operator.yaml
    return
  fi
}

function create_namespace() {
    namespace="arena-system"
    if [[ "${NAMESPACE}" != "" ]]; then
        namespace=${NAMESPACE}
    fi
    if arena-kubectl get ns | grep -E "\<$namespace\>" &> /dev/null;then
        logger "debug" "namespace $namespace has been existed,skip to create it"
        return
    fi
    arena-kubectl create ns $namespace
}

function binary() {
    install_kubectl_and_helm
    custom_charts
    install_arena_and_charts
    install_arena_gen_kubeconfig
}

function operators() {
    if [[ $ONLY_BINARY == "true" ]];then
        logger "debug" "skip to install operators,because --only-binary is enabled"
        return 
    fi
    custom_manifests
    create_namespace  
    apply_tf
    apply_pytorch
    apply_mpi
    apply_et
    apply_prometheus
    apply_jobmon
    apply_rdma
    apply_kubedl
}


function parse_args() {
    while
        [[ $# -gt 0 ]]
    do
        key="$1"
		case $key in
        --only-binary)
            export ONLY_BINARY="true"
			;;
        --host-network)
            export USE_HOSTNETWORK="true"
			;;
        --rdma)
            export USE_RDMA="true"
			;;
        --loadbalancer)
            export USE_LOADBALANCER="true"
			;;
        --prometheus)
            export USE_PROMETHEUS="true"
			;;
        --kubeconfig)
            check_option_value "--kubeconfig" $2
            export KUBECONFIG=$2
            shift
			;;    
        --platform)
            check_option_value "--platform" $2
            export PLATFORM=$2
            shift
			;;    
        --region-id)
            check_option_value "--region-id" $2
            export REGION=$2
            shift
			;;    
        --docker-registry)
            check_option_value "--docker-registry" $2
            export DOCKER_REGISTRY=$2
            shift
			;;    
        --registry-repo-namespace)
            check_option_value "--registry-repo-namespace" $2
            export REGISTRY_REPO_NAMESPACE=$2
            shift
			;;    
        --namespace)
            check_option_value "--namespace" $2
            export NAMESPACE=$2
            shift
			;;    
        --help|-h)
            help
            exit 0
            ;;    
        *)
            # unknown option
            logger error "unkonw option [$key]"
            help
            exit 3
            ;;
        esac
        shift
    done
}

function check_option_value() {
    option=$1
    value=$2
    if [[ $value == "" ]] || echo "$value" | grep -- "^--" &> /dev/null;then
        logger error "the option $option not set value,please set it"
        exit 3
    fi  
}


function main() {
    parse_args "$@"
    set_sudo
    binary
    operators
    logger "debug" "--------------------------------"
    logger "debug" "Arena has been installed successfully!"
}

main "$@"
