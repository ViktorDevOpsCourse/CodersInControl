package controller

import (
	"fmt"
	v12 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"time"
)

// Status defines the set of statuses a resource can have.
type Status string

const (
	InProgressStatus  Status = "InProgress"
	FailedStatus      Status = "Failed"
	RunningStatus     Status = "Running"
	TerminatingStatus Status = "Terminating"

	ConditionStalled     string = "Stalled"
	ConditionReconciling string = "Reconciling"
)

type DeploymentsControllerConfig struct {
	AddFunc    func(obj interface{})
	UpdateFunc func(oldObj, newObj interface{})
	DeleteFunc func(obj interface{})
}

type DeploymentsController struct {
	watchList      cache.ListerWatcher
	k8sClient      *kubernetes.Clientset
	stopController chan struct{}
	config         DeploymentsControllerConfig
}

func NewController(k8sClient *kubernetes.Clientset, cfg DeploymentsControllerConfig) (*DeploymentsController, error) {
	if cfg.AddFunc == nil || cfg.UpdateFunc == nil || cfg.DeleteFunc == nil {
		return nil, fmt.Errorf("deployment controller events funcs should be defined got: AddFunc %s, UpdateFunc %#v, DeleteFunc %#v", cfg.AddFunc, cfg.UpdateFunc, cfg.DeleteFunc)
	}

	return &DeploymentsController{
		k8sClient: k8sClient,
		watchList: cache.NewListWatchFromClient(
			k8sClient.AppsV1().RESTClient(),
			"deployments",
			v1.NamespaceAll,
			fields.Everything()),
		config: cfg,
	}, nil
}

func (d *DeploymentsController) Run() error {

	informer := cache.NewSharedIndexInformer(
		d.watchList,
		&v12.Deployment{},
		time.Second*5,
		cache.Indexers{},
	)

	_, err := informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    d.config.AddFunc,
		UpdateFunc: d.config.UpdateFunc,
		DeleteFunc: d.config.DeleteFunc,
	})

	if err != nil {
		return err
	}

	d.stopController = make(chan struct{})
	go informer.Run(d.stopController)

	return nil
}

func (d *DeploymentsController) Stop() {
	close(d.stopController)
}

func checkGenericProperties(deployment *v12.Deployment) (Status, error) {

	if !deployment.ObjectMeta.DeletionTimestamp.IsZero() {
		return TerminatingStatus, nil
	}

	res, err := checkGeneration(deployment)
	if res != "" || err != nil {
		return res, err
	}

	for _, cond := range deployment.Status.Conditions {
		if string(cond.Type) == string(ConditionReconciling) && cond.Status == corev1.ConditionTrue {
			return InProgressStatus, nil
		}
		if string(cond.Type) == string(ConditionStalled) && cond.Status == corev1.ConditionTrue {
			return FailedStatus, nil
		}
	}

	return "", nil
}

func checkGeneration(deployment *v12.Deployment) (Status, error) {
	if deployment.Status.ObservedGeneration != deployment.ObjectMeta.Generation {
		return InProgressStatus, nil
	}

	return "", nil
}

// deploymentConditions return standardized Conditions for Deployment.
//
// For Deployments, we look at .status.conditions as well as the other properties
// under .status. Status will be Failed if the progress deadline has been exceeded.
func deploymentConditions(deployment *v12.Deployment) (Status, error) {

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
