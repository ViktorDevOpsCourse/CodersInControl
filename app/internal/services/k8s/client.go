package k8s

import (
	"flag"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"path/filepath"
)

type Client struct {
	http *kubernetes.Clientset
}

func NewClient() (Client, error) {
	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "path to kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "path to kubeconfig file")
	}
	flag.Parse()

	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
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
