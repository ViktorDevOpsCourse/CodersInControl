package clusters

import (
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

type Conditions struct {
	AvailableReplicas int32
	ServiceStatus     Status
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
