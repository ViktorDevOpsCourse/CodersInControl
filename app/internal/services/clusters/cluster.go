package clusters

import (
	"github.com/viktordevopscourse/codersincontrol/app/internal/services/clusters/controller"
	"github.com/viktordevopscourse/codersincontrol/app/pkg/logger"
	v12 "k8s.io/api/apps/v1"
	"sync"
)

type Cluster struct {
	client       Client
	Applications map[string]Application
	Namespaces   map[string]Namespace
	Controller   controller.Controller
	sync.RWMutex
}

func NewCluster(client Client) *Cluster {
	log := logger.FromDefaultContext()

	c := &Cluster{
		client:       client,
		Applications: make(map[string]Application),
		Namespaces:   make(map[string]Namespace),
	}

	deploymentController, err := controller.NewController(client.http, controller.ConfigController{
		AddFunc:    c.addEventHandler,
		UpdateFunc: c.updateEventHandler,
		DeleteFunc: c.deleteEventHandler,
	})
	if err != nil {
		log.Errorf("failed setup new controller while creating cluster. err `%s`", err)
		return nil
	}

	c.Controller = deploymentController

	return c
}

func (c *Cluster) Run() {
	log := logger.FromDefaultContext()

	err := c.Controller.Run()
	if err != nil {
		log.Errorf("run cluster controller failed. err `%s`", err)
		return
	}
}

func (c *Cluster) Stop() {
	c.Controller.Stop()
}

func (c *Cluster) GetApplications() map[string]Application {
	return c.Applications
}

func (c *Cluster) GetApplicationByName(name string) Application {
	return c.Applications[name]
}

func (c *Cluster) addEventHandler(obj interface{}) {

	deployment := obj.(*v12.Deployment)

	if _, ok := excludeNamespaces[deployment.Namespace]; ok {
		return
	}

	c.Lock()
	defer c.Unlock()

	c.Namespaces[deployment.Namespace] = Namespace{
		Name: deployment.Namespace,
	}

	c.Applications[deployment.GetName()] = Application{
		Name:                 deployment.GetName(),
		Namespace:            deployment.GetNamespace(),
		AppliedConfiguration: deployment.Annotations["kubectl.kubernetes.io/last-applied-configuration"],
		CreatedAt:            deployment.CreationTimestamp.Time,
		Labels:               deployment.GetLabels(),
		Replicas:             deployment.Spec.Replicas,
		SelectorMatchLabels:  deployment.Spec.Selector.MatchLabels,
		Image:                deployment.Spec.Template.Spec.Containers[0].Image, // todo understand how update this field
		Status: Conditions{
			AvailableReplicas: deployment.Status.AvailableReplicas,
			ServiceStatus:     controller.RunningStatus,
		},
	}
}

func (c *Cluster) updateEventHandler(oldObj, newObj interface{}) {
	//log := logger.FromDefaultContext()
	//
	//newDeployment := newObj.(*v12.Deployment)
	//
	//status, err := checkGenericProperties(newDeployment)
	//if err != nil {
	//	log.Errorf("Deployemnt controller, error update event handle `%s`", err)
	//	return
	//}
	//
	//if status != "" {
	//	// TODO handel status
	//	return
	//}
	//
	//status, err = deploymentConditions(newDeployment)
	//if err != nil {
	//	log.Errorf("Deployemnt controller, error update event handle `%s`", err)
	//	return
	//}
	//
	//// TODO handel status
}

func (c *Cluster) deleteEventHandler(obj interface{}) {
	deployment := obj.(*v12.Deployment)
	c.Lock()
	defer c.Unlock()

	delete(c.Applications, deployment.GetName())
}
