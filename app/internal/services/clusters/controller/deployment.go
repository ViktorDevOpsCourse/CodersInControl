package controller

import (
	v12 "k8s.io/api/apps/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"time"
)

type DeploymentsController struct {
	watchList      cache.ListerWatcher
	k8sClient      *kubernetes.Clientset
	stopController chan struct{}
	config         ConfigController
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
