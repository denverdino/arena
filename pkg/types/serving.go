package types

import (
	v1 "k8s.io/api/core/v1"

	"github.com/kubeflow/arena/pkg/util"
	log "github.com/sirupsen/logrus"
	app_v1 "k8s.io/api/apps/v1"
	"k8s.io/client-go/kubernetes"
)

type Serving struct {
	Name      string
	ServeType string
	Namespace string
	Version   string
	pods      []v1.Pod
	deploy    app_v1.Deployment
	client    *kubernetes.Clientset
}

var SERVING_CHARTS = map[string]string{
	"tensorflow-serving-0.2.0":        "Tensorflow",
	"tensorrt-inference-server-0.0.1": "TensorRT",
}
var SERVING_TYPE = map[string]string{
	"tf-serving":     "TENSORFLOW",
	"trt-serving":    "TENSORRT",
	"custom-serving": "CUSTOM",
}

func NewServingJob(client *kubernetes.Clientset, deploy app_v1.Deployment, allPods []v1.Pod) Serving {
	servingTypeLabel := deploy.Labels["servingType"]
	servingVersion := deploy.Labels["servingVersion"]
	servingName := deploy.Labels["servingName"]
	servingType := "Tensorflow"
	if serveType, ok := SERVING_TYPE[servingTypeLabel]; ok {
		servingType = serveType
	}
	serving := Serving{
		Name:      servingName,
		client:    client,
		ServeType: servingType,
		Namespace: deploy.Namespace,
		Version:   servingVersion,
		deploy:    deploy,
	}
	for _, pod := range allPods {
		if IsPodControllerByDeploment(pod, deploy) {
			serving.pods = append(serving.pods, pod)
		}
	}
	return serving
}

func (s Serving) GetName() string {
	// if s.Version != "" {
	// 	return fmt.Sprintf("%s-%s", s.Name, s.Version)
	// }
	return s.Name
}

func (s Serving) AllPods() []v1.Pod {
	return s.pods
}

func (s Serving) GetClusterIP() string {
	// serviceName := fmt.Sprintf("%s-%s", s.deploy.Labels["release"], s.deploy.Labels["app"])
	allServices, err := util.AcquireServingServices(s.Namespace, s.client)
	if err != nil {
		log.Errorf("failed to list serving services, err: %++v", err)
		return "N/A"
	}

	log.Debugf("try to get Endpoint IP for name %s and %s", s.Name, s.Version)
	for _, service := range allServices {
		if service.Labels["servingName"] == s.Name &&
			service.Labels["servingVersion"] == s.Version {
			return service.Spec.ClusterIP
		}

		// if service.Name == serviceName {
		// 	return service.Spec.ClusterIP
		// }
	}
	return "N/A"
}

// Available instances
func (s Serving) AvailableInstances() int32 {
	return s.deploy.Status.AvailableReplicas
}

// Desired Instances
func (s Serving) DesiredInstances() int32 {
	return s.deploy.Status.Replicas
}

func (s Serving) GetStatus() string {
	hasPendingPod := false
	for _, pod := range s.pods {
		if pod.Status.Phase == v1.PodPending {
			log.Debugf("pod %s is pending", pod.Name)
			hasPendingPod = true
			break
		}
		if hasPendingPod {
			return "PENDING"
		}
	}
	return "RUNNING"
}

func IsPodControllerByDeploment(pod v1.Pod, deploy app_v1.Deployment) bool {
	if len(pod.OwnerReferences) == 0 {
		return false
	}
	podLabel := pod.GetLabels()
	if len(podLabel) == 0 {
		return false
	}
	for key, value := range deploy.Spec.Selector.MatchLabels {
		if podLabel[key] != value {
			return false
		}
	}
	return true
}
