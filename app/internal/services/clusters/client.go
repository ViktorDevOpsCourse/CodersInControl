package clusters

import (
	"context"
	"fmt"
	v1 "k8s.io/api/core/v1"
	v12 "k8s.io/apimachinery/pkg/apis/meta/v1"
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

	c := Client{
		http: clientSet,
	}

	err = c.Ping()
	if err != nil {
		return Client{}, fmt.Errorf("failed connect to cluster `%s`. can't obtain pods", config.Host)
	}

	return c, nil
}

func (c Client) Ping() error {
	_, err := c.http.CoreV1().Pods(v1.NamespaceAll).List(context.TODO(), v12.ListOptions{})
	return err
}
