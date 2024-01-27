package clusters

import (
	"github.com/viktordevopscourse/codersincontrol/app/internal/services/clusters/controller"
	"time"
)

type Conditions struct {
	AvailableReplicas int32
	ServiceStatus     controller.Status
}

type Application struct {
	Namespace            string
	Name                 string
	CreatedAt            time.Time
	Labels               map[string]string
	Replicas             *int32
	SelectorMatchLabels  map[string]string
	Image                string
	Status               Conditions
	AppliedConfiguration string
}

func (a *Application) GetName() string {
	return a.Name
}

func (a *Application) GetNamespaceName() string {
	return a.Namespace
}
