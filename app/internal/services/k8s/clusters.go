package k8s

import (
	"context"
	"fmt"
	"github.com/viktordevopscourse/codersincontrol/app/pkg/logger"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Cluster struct {
	client       Client
	Applications []Application
	Namespaces   []Namespace
}

func NewCluster(client Client) Cluster {
	log := logger.FromDefaultContext()

	c := Cluster{
		client:       client,
		Applications: make([]Application, 0),
		Namespaces:   make([]Namespace, 0),
	}

	err := c.initReadCluster()
	if err != nil {
		log.Fatalf("Failed init read cluster. Err `%s`", err)
		return Cluster{}
	}

	return c
}

func (c *Cluster) initReadCluster() error {
	ctx := context.TODO()

	err := c.readAllNamespaces(ctx)
	if err != nil {
		return err
	}

	for _, namespace := range c.Namespaces {
		err := c.readAllApps(ctx, namespace.GetName())
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *Cluster) readAllNamespaces(ctx context.Context) error {
	namespaces, err := c.client.http.CoreV1().Namespaces().List(ctx, v1.ListOptions{})
	if err != nil {
		return err
	}

	for _, n := range namespaces.Items {
		if _, ok := excludeNamespaces[n.GetName()]; ok {
			continue
		}

		c.Namespaces = append(c.Namespaces, Namespace{
			Name: n.GetName(),
		})
	}

	return nil
}

func (c *Cluster) readAllApps(ctx context.Context, namespace string) error {
	deployments, err := c.client.http.AppsV1().Deployments(namespace).List(ctx, v1.ListOptions{})
	if err != nil {
		return fmt.Errorf("error get list apps in `%s` namespace. Error `%s", namespace, err)
	}

	for _, deploy := range deployments.Items {
		c.Applications = append(c.Applications, Application{
			Name:                 deploy.GetName(),
			Namespace:            deploy.GetNamespace(),
			AppliedConfiguration: deploy.Annotations["kubectl.kubernetes.io/last-applied-configuration"],
			CreatedAt:            deploy.CreationTimestamp.Time,
			Labels:               deploy.GetLabels(),
			Replicas:             deploy.Spec.Replicas,
			SelectorMatchLabels:  deploy.Spec.Selector.MatchLabels,
			Image:                deploy.Spec.Template.Spec.Containers[0].Image, // todo understand how update this field
			Status: ApplicationStatus{
				AvailableReplicas: deploy.Status.AvailableReplicas,
			},
		})
	}

	return nil
}
