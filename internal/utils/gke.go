package utils

import (
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type GKEService struct {
	kubeclient *kubernetes.Clientset
}

func NewGKEService(kubeconfig string, namespace string) (GKEService, error) {
	var config *rest.Config
	var err error
	config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return GKEService{}, err
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return GKEService{}, err
	}
	service := GKEService{
		kubeclient: clientset,
	}
	return service, nil
}
