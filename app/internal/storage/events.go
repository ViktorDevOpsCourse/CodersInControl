package storage

import (
	"sync"
	"time"
)

type EventsRepository interface {
	GetAndRemove(clusterName, appName string) (ApplicationEvent, error)
	Save(clusterName, appName string, event ApplicationEvent) error
}

type ApplicationEvent struct {
	AppName         string
	Image           string
	ResourceVersion string
	EventTime       time.Time
	Status          string
}

type ApplicationsEvents struct {
	events map[string]map[string]ApplicationEvent // map[clusterName]map[appName]ApplicationData
	sync.Mutex
}

func NewApplicationsEvents() EventsRepository {
	return &ApplicationsEvents{
		events: make(map[string]map[string]ApplicationEvent),
	}
}

func (e *ApplicationsEvents) GetAndRemove(clusterName, appName string) (ApplicationEvent, error) {
	e.Lock()
	defer e.Unlock()

	if _, ok := e.events[clusterName]; !ok {
		return ApplicationEvent{}, NotFoundError
	}

	if _, ok := e.events[clusterName][appName]; !ok {
		return ApplicationEvent{}, NotFoundError
	}

	appEvent := e.events[clusterName][appName]

	delete(e.events[clusterName], appName)

	return appEvent, nil

}

func (e *ApplicationsEvents) Save(clusterName, appName string, event ApplicationEvent) error {
	if _, ok := e.events[clusterName]; !ok {
		e.events[clusterName] = make(map[string]ApplicationEvent)
	}

	e.events[clusterName][appName] = event
	return nil
}
