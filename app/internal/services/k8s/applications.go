package k8s

import "time"

type ApplicationStatus struct {
	AvailableReplicas int32
}

type Application struct {
	Namespace            string
	Name                 string
	CreatedAt            time.Time
	Labels               map[string]string
	Replicas             *int32
	SelectorMatchLabels  map[string]string
	Image                string
	Status               ApplicationStatus
	AppliedConfiguration string
}

func (a *Application) GetName() string {
	return a.Name
}

func (a *Application) GetNamespaceName() string {
	return a.Namespace
}
