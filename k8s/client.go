package k8s

import (
	"fmt"

	"github.com/spf13/viper"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

// KubeAPI performs all the Kubernetes API operations
type KubeAPI struct {
	Client    kubernetes.Interface
	Namespace string
}

// NewClient returns a new KubeAPI client
func NewClient() (*KubeAPI, error) {
	cluster := viper.GetString("kubernetes_cluster")
	namespace := viper.GetString("kubernetes_namespace")
	if cluster == "" {
		return nil, fmt.Errorf("must supply kubernetes_cluster in duncan.yml")
	}
	if namespace == "" {
		return nil, fmt.Errorf("must supply kubernetes_namespace in duncan.yml")
	}
	loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
	configOverrides := &clientcmd.ConfigOverrides{
		CurrentContext: cluster,
	}
	kubeConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, configOverrides)
	config, err := kubeConfig.ClientConfig()
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
