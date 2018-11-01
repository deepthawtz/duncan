package k8s

import (
	"fmt"
	"path/filepath"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

// KubeAPI performs all the Kubernetes API operations
type KubeAPI struct {
	Client    kubernetes.Interface
	Namespace string
}

// NewClient returns a new KubeAPI client
func NewClient(namespace string) (*KubeAPI, error) {
	if namespace == "" {
		return nil, fmt.Errorf("must supply kubernetes_namespace in duncan.yml")
	}
	kubeconfig := filepath.Join(homedir.HomeDir(), ".kube", "config")
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return nil, err
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return &KubeAPI{
		Client:    clientset,
		Namespace: namespace,
	}, nil
}
