package clusters

import (
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

type Client struct {
	http *kubernetes.Clientset
}

func NewClient(kubeConfig string) (Client, error) {
	config, err := clientcmd.BuildConfigFromFlags("", kubeConfig)
	if err != nil {
		return Client{}, err
	}

	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		return Client{}, err
	}
	return Client{
		http: clientSet,
	}, nil
}
