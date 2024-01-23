package k8s

import (
	"context"
	"github.com/viktordevopscourse/codersincontrol/app/pkg/logger"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Cluster struct {
	client     Client
	Namespaces []Namespace
}

func NewCluster() Cluster {
	log := logger.FromDefaultContext()
	client, err := NewClient()
	if err != nil {
		log.Fatalf("Failed connection to cluster. Err `%s`", err)
		return Cluster{}
	}
	c := Cluster{
		client: client,
	}

	err = c.initReadCluster()
	if err != nil {
		log.Fatalf("Failed connection to cluster. Err `%s`", err)
		return Cluster{}
	}

	return c
}

func (c *Cluster) initReadCluster() error {
	ctx := context.TODO()
	log := logger.FromDefaultContext()
	namespaces, err := c.client.http.CoreV1().Namespaces().List(ctx, v1.ListOptions{})
	if err != nil {
		return err
	}

	for _, n := range namespaces.Items {
		deployments, err := c.client.http.AppsV1().Deployments(n.GetName()).List(ctx, v1.ListOptions{})
		if err != nil {
			log.WithField("Namespace", n.Namespace).Errorf("Error get list deploymrnts. Error `%s", err)
			continue
		}

		namespace := Namespace{
			Name:        n.GetName(),
			Deployments: make([]Deployment, 0),
		}

		for _, deploy := range deployments.Items {
			namespace.Deployments = append(namespace.Deployments, Deployment{
				Name: deploy.GetName(),
			})
		}

		c.Namespaces = append(c.Namespaces, namespace)

	}

	return nil
}
