// +build !ignore_autogenerated

// Copyright 2020 The Kubeflow Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Code generated by openapi-gen. DO NOT EDIT.

// This file was autogenerated by openapi-gen. Do not edit it manually!

package v1

import (
	spec "github.com/go-openapi/spec"
	common "k8s.io/kube-openapi/pkg/common"
)

func GetOpenAPIDefinitions(ref common.ReferenceCallback) map[string]common.OpenAPIDefinition {
	return map[string]common.OpenAPIDefinition{
		"github.com/kubeflow/pytorch-operator/pkg/apis/pytorch/v1.PyTorchJob":     schema_pkg_apis_pytorch_v1_PyTorchJob(ref),
		"github.com/kubeflow/pytorch-operator/pkg/apis/pytorch/v1.PyTorchJobList": schema_pkg_apis_pytorch_v1_PyTorchJobList(ref),
		"github.com/kubeflow/pytorch-operator/pkg/apis/pytorch/v1.PyTorchJobSpec": schema_pkg_apis_pytorch_v1_PyTorchJobSpec(ref),
	}
}

func schema_pkg_apis_pytorch_v1_PyTorchJob(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "Represents a PyTorchJob resource.",
				Type:        []string{"object"},
				Properties: map[string]spec.Schema{
					"kind": {
						SchemaProps: spec.SchemaProps{
							Description: "Kind is a string value representing the REST resource this object represents. Servers may infer this from the endpoint the client submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#types-kinds",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"apiVersion": {
						SchemaProps: spec.SchemaProps{
							Description: "APIVersion defines the versioned schema of this representation of an object. Servers should convert recognized schemas to the latest internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#resources",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"metadata": {
						SchemaProps: spec.SchemaProps{
							Description: "Standard Kubernetes object's metadata.",
							Ref:         ref("k8s.io/apimachinery/pkg/apis/meta/v1.ObjectMeta"),
						},
					},
					"spec": {
						SchemaProps: spec.SchemaProps{
							Description: "Specification of the desired state of the PyTorchJob.",
							Ref:         ref("github.com/kubeflow/pytorch-operator/pkg/apis/pytorch/v1.PyTorchJobSpec"),
						},
					},
					"status": {
						SchemaProps: spec.SchemaProps{
							Description: "Most recently observed status of the PyTorchJob. Read-only (modified by the system).",
							Ref:         ref("github.com/kubeflow/common/pkg/apis/common/v1.JobStatus"),
						},
					},
				},
			},
		},
		Dependencies: []string{
			"github.com/kubeflow/common/pkg/apis/common/v1.JobStatus", "github.com/kubeflow/pytorch-operator/pkg/apis/pytorch/v1.PyTorchJobSpec", "k8s.io/apimachinery/pkg/apis/meta/v1.ObjectMeta"},
	}
}

func schema_pkg_apis_pytorch_v1_PyTorchJobList(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "PyTorchJobList is a list of PyTorchJobs.",
				Type:        []string{"object"},
				Properties: map[string]spec.Schema{
					"kind": {
						SchemaProps: spec.SchemaProps{
							Description: "Kind is a string value representing the REST resource this object represents. Servers may infer this from the endpoint the client submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#types-kinds",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"apiVersion": {
						SchemaProps: spec.SchemaProps{
							Description: "APIVersion defines the versioned schema of this representation of an object. Servers should convert recognized schemas to the latest internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#resources",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"metadata": {
						SchemaProps: spec.SchemaProps{
							Description: "Standard list metadata.",
							Ref:         ref("k8s.io/apimachinery/pkg/apis/meta/v1.ListMeta"),
						},
					},
					"items": {
						SchemaProps: spec.SchemaProps{
							Description: "List of PyTorchJobs.",
							Type:        []string{"array"},
							Items: &spec.SchemaOrArray{
								Schema: &spec.Schema{
									SchemaProps: spec.SchemaProps{
										Ref: ref("github.com/kubeflow/pytorch-operator/pkg/apis/pytorch/v1.PyTorchJob"),
									},
								},
							},
						},
					},
				},
				Required: []string{"items"},
			},
		},
		Dependencies: []string{
			"github.com/kubeflow/pytorch-operator/pkg/apis/pytorch/v1.PyTorchJob", "k8s.io/apimachinery/pkg/apis/meta/v1.ListMeta"},
	}
}

func schema_pkg_apis_pytorch_v1_PyTorchJobSpec(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "PyTorchJobSpec is a desired state description of the PyTorchJob.",
				Type:        []string{"object"},
				Properties: map[string]spec.Schema{
					"activeDeadlineSeconds": {
						SchemaProps: spec.SchemaProps{
							Description: "Specifies the duration (in seconds) since startTime during which the job can remain active before it is terminated. Must be a positive integer. This setting applies only to pods where restartPolicy is OnFailure or Always.",
							Type:        []string{"integer"},
							Format:      "int64",
						},
					},
					"backoffLimit": {
						SchemaProps: spec.SchemaProps{
							Description: "Number of retries before marking this job as failed.",
							Type:        []string{"integer"},
							Format:      "int32",
						},
					},
					"cleanPodPolicy": {
						SchemaProps: spec.SchemaProps{
							Description: "Defines the policy for cleaning up pods after the PyTorchJob completes. Defaults to None.",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"ttlSecondsAfterFinished": {
						SchemaProps: spec.SchemaProps{
							Description: "Defines the TTL for cleaning up finished PyTorchJobs (temporary before Kubernetes adds the cleanup controller). It may take extra ReconcilePeriod seconds for the cleanup, since reconcile gets called periodically. Defaults to infinite.",
							Type:        []string{"integer"},
							Format:      "int32",
						},
					},
					"pytorchReplicaSpecs": {
						SchemaProps: spec.SchemaProps{
							Description: "A map of PyTorchReplicaType (type) to ReplicaSpec (value). Specifies the PyTorch cluster configuration. For example,\n  {\n    \"Master\": PyTorchReplicaSpec,\n    \"Worker\": PyTorchReplicaSpec,\n  }",
							Type:        []string{"object"},
							AdditionalProperties: &spec.SchemaOrBool{
								Allows: true,
								Schema: &spec.Schema{
									SchemaProps: spec.SchemaProps{
										Ref: ref("github.com/kubeflow/common/pkg/apis/common/v1.ReplicaSpec"),
									},
								},
							},
						},
					},
				},
				Required: []string{"pytorchReplicaSpecs"},
			},
		},
		Dependencies: []string{
			"github.com/kubeflow/common/pkg/apis/common/v1.ReplicaSpec"},
	}
}
