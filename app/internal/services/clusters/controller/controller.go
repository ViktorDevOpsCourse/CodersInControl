package controller

import (
	"fmt"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
)

type Controller interface {
	Run() error
	Stop()
}

type ConfigController struct {
	AddFunc    func(obj interface{})
	UpdateFunc func(oldObj, newObj interface{})
	DeleteFunc func(obj interface{})
}

func NewController(k8sClient *kubernetes.Clientset, cfg ConfigController) (Controller, error) {
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
