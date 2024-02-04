package clusters

import (
	"fmt"
	"github.com/viktordevopscourse/codersincontrol/app/internal/services/clusters/controller"
	"github.com/viktordevopscourse/codersincontrol/app/internal/storage"
	"github.com/viktordevopscourse/codersincontrol/app/pkg/logger"
	v12 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"strings"
	"sync"
	"time"
)

type Cluster struct {
	client                 Client
	Applications           map[string]Application // map[appName]
	Namespaces             map[string]Namespace   // map[nsName]
	Controller             controller.Controller
	ClusterName            string
	lastAppResourceVersion map[string]string // map[appName]revision

	appsStatesStorage storage.StateRepository
	appsEventsStorage storage.EventsRepository
	sync.RWMutex
}

func NewCluster(clusterName string,
	client Client,
	appsStatesStorage storage.StateRepository,
	appsEventsStorage storage.EventsRepository) *Cluster {
	log := logger.FromDefaultContext()

	c := &Cluster{
		client:                 client,
		Applications:           make(map[string]Application),
		Namespaces:             make(map[string]Namespace),
		ClusterName:            clusterName,
		appsStatesStorage:      appsStatesStorage,
		appsEventsStorage:      appsEventsStorage,
		lastAppResourceVersion: make(map[string]string),
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

	log := logger.FromDefaultContext()
	deployment := obj.(*v12.Deployment)

	if !c.IsWatchedNamespace(deployment.Namespace) {
		return
	}

	c.Lock()
	defer c.Unlock()

	c.Namespaces[deployment.Namespace] = Namespace{
		Name: deployment.Namespace,
	}

	status, err := c.getDeploymentStatus(deployment)
	if err != nil {
		log.Error(err)
		return
	}

	c.updateAppCurrentState(deployment, status)
}

func (c *Cluster) updateEventHandler(oldObj, newObj interface{}) {
	log := logger.FromDefaultContext()

	newApp := newObj.(*v12.Deployment)

	if !c.IsWatchedNamespace(newApp.Namespace) {
		return
	}

	status, err := c.getDeploymentStatus(newApp)
	if err != nil {
		log.Error(err)
		return
	}

	if status == InProgressStatus {
		return
	}

	if prevVersion, ok := c.lastAppResourceVersion[newApp.GetName()]; ok {
		if prevVersion == newApp.ResourceVersion {
			log.Infof("updateEventHandler ResourceVersion for status %s same. Skip", newApp.ResourceVersion)
			return
		}
	}

	c.lastAppResourceVersion[newApp.GetName()] = newApp.ResourceVersion

	err = c.appsEventsStorage.Save(c.ClusterName, newApp.GetName(), storage.ApplicationEvent{
		AppName:         newApp.GetName(),
		Image:           c.getImageVersion(newApp),
		EventTime:       time.Now(),
		Status:          string(status),
		ResourceVersion: newApp.ResourceVersion,
	})
	if err != nil {
		log.Error(err)
	}

	c.updateAppCurrentState(newApp, status)

}

func (c *Cluster) deleteEventHandler(obj interface{}) {
	deployment := obj.(*v12.Deployment)
	c.Lock()
	defer c.Unlock()

	delete(c.Applications, deployment.GetName())
}

func (c *Cluster) updateAppCurrentState(app *v12.Deployment, status Status) {

	c.Applications[app.GetName()] = Application{
		Name:                 app.GetName(),
		Namespace:            app.GetNamespace(),
		AppliedConfiguration: app.Annotations["kubectl.kubernetes.io/last-applied-configuration"],
		CreatedAt:            app.CreationTimestamp.Time,
		Labels:               app.GetLabels(),
		Replicas:             app.Spec.Replicas,
		SelectorMatchLabels:  app.Spec.Selector.MatchLabels,
		Image:                c.getImageVersion(app),
		Status: Conditions{
			AvailableReplicas: app.Status.AvailableReplicas,
			ServiceStatus:     status,
		},
	}
}

func (c *Cluster) IsWatchedNamespace(namespace string) bool {
	if _, ok := excludeNamespaces[namespace]; ok {
		return false
	}
	return true
}
func (c *Cluster) getImageVersion(app *v12.Deployment) string {
	if len(app.Spec.Template.Spec.Containers) != 0 {
		imageParts := strings.Split(app.Spec.Template.Spec.Containers[0].Image, ":")
		if len(imageParts) == 2 {
			return imageParts[1]
		}
	}
	return ""
}

func (c *Cluster) getDeploymentStatus(deployment *v12.Deployment) (Status, error) {
	status, err := c.checkGenericProperties(deployment)
	if err != nil {

		return "", fmt.Errorf("deployemnt controller, error update event handle `%s`", err)
	}

	if status != "" {
		return status, nil
	}

	status, err = c.deploymentConditions(deployment)
	if err != nil {
		return "", fmt.Errorf("deployemnt controller, error update event handle `%s`", err)
	}

	return status, nil
}

func (c *Cluster) checkGenericProperties(deployment *v12.Deployment) (Status, error) {

	if !deployment.ObjectMeta.DeletionTimestamp.IsZero() {
		return TerminatingStatus, nil
	}

	res, err := c.checkGeneration(deployment)
	if res != "" || err != nil {
		return res, err
	}

	for _, cond := range deployment.Status.Conditions {
		if string(cond.Type) == ConditionReconciling && cond.Status == corev1.ConditionTrue {
			return InProgressStatus, nil
		}
		if string(cond.Type) == ConditionStalled && cond.Status == corev1.ConditionTrue {
			return FailedStatus, nil
		}
	}

	return "", nil
}

func (c *Cluster) checkGeneration(deployment *v12.Deployment) (Status, error) {
	if deployment.Status.ObservedGeneration != deployment.ObjectMeta.Generation {
		return InProgressStatus, nil
	}

	return "", nil
}

// deploymentConditions return standardized Conditions for Deployment.
//
// For Deployments, we look at .status.conditions as well as the other properties
// under .status. Status will be Failed if the progress deadline has been exceeded.
func (c *Cluster) deploymentConditions(deployment *v12.Deployment) (Status, error) {

	progressing := false
	available := false

	for _, c := range deployment.Status.Conditions {
		switch c.Type {
		case "Progressing": // appsv1.DeploymentProgressing:
			if c.Reason == "ProgressDeadlineExceeded" {
				return FailedStatus, nil
			}
			if c.Status == corev1.ConditionTrue && c.Reason == "NewReplicaSetAvailable" {
				progressing = true
			}
		case "Available": // appsv1.DeploymentAvailable:
			if c.Status == corev1.ConditionTrue {
				available = true
			}
		}
	}

	var specReplicas int32
	if deployment.Spec.Replicas != nil {
		specReplicas = *deployment.Spec.Replicas
	}
	statusReplicas := deployment.Status.Replicas
	updatedReplicas := deployment.Status.UpdatedReplicas
	readyReplicas := deployment.Status.ReadyReplicas
	availableReplicas := deployment.Status.AvailableReplicas

	if specReplicas > statusReplicas {
		return InProgressStatus, nil
	}

	if specReplicas > updatedReplicas {
		return InProgressStatus, nil
	}

	if statusReplicas > specReplicas {
		return InProgressStatus, nil
	}

	if updatedReplicas > availableReplicas {
		return InProgressStatus, nil
	}

	if specReplicas > readyReplicas {
		return InProgressStatus, nil
	}

	// check conditions
	if !progressing {
		return InProgressStatus, nil
	}

	if !available {
		return InProgressStatus, nil
	}

	// All ok
	return RunningStatus, nil
}
