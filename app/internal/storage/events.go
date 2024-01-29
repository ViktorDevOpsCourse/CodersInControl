package storage

import (
	"sync"
	"time"
)

type EventsRepository interface {
	GetAndRemove(envName, appName string) (ApplicationEvent, error)
	Save(envName, appName string, event ApplicationEvent) error
}

type ApplicationEvent struct {
	AppName         string
	Image           string
	ResourceVersion string
	EventTime       time.Time
	Status          string
}

type ApplicationsEvents struct {
	events map[string]map[string]ApplicationEvent // map[envName]map[appName]ApplicationData
	sync.Mutex
}

func NewApplicationsEvents() EventsRepository {
	return &ApplicationsEvents{
		events: make(map[string]map[string]ApplicationEvent),
	}
}

func (e *ApplicationsEvents) GetAndRemove(envName, appName string) (ApplicationEvent, error) {
	e.Lock()
	defer e.Unlock()

	if _, ok := e.events[envName]; !ok {
		return ApplicationEvent{}, NotFoundError
	}

	envApps := e.events[envName]

	if _, ok := envApps[appName]; !ok {
		return ApplicationEvent{}, NotFoundError
	}

	appEvent := envApps[appName]

	delete(envApps, appName)

	return appEvent, nil

}

func (e *ApplicationsEvents) Save(envName, appName string, event ApplicationEvent) error {
	if _, ok := e.events[envName]; !ok {
		e.events[envName] = make(map[string]ApplicationEvent)
	}

	e.events[envName][appName] = event
	return nil
}
