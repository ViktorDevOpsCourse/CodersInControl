package k8s

import (
	"context"
	"fmt"
	"github.com/viktordevopscourse/codersincontrol/app/pkg/logger"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Cluster struct {
	client       Client
	Applications map[string]Application
	Namespaces   map[string]Namespace
}

func NewCluster(client Client) Cluster {
	log := logger.FromDefaultContext()

	c := Cluster{
		client:       client,
		Applications: make(map[string]Application),
		Namespaces:   make(map[string]Namespace),
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

		c.Namespaces[n.GetName()] = Namespace{
			Name: n.GetName(),
		}
	}

	return nil
}

func (c *Cluster) readAllApps(ctx context.Context, namespace string) error {
	log := logger.FromContext(ctx)

	deployments, err := c.client.http.AppsV1().Deployments(namespace).List(ctx, v1.ListOptions{})
	if err != nil {
		return fmt.Errorf("error get list apps in `%s` namespace. Error `%s", namespace, err)
	}

	for _, deploy := range deployments.Items {
		if app, ok := c.Applications[deploy.GetName()]; ok {
			log.Errorf("Application `%s` already present in namespace `%s`", deploy.GetName(), app.Namespace)
			continue
		}

		c.Applications[deploy.GetName()] = Application{
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
		}
	}

	return nil
}
