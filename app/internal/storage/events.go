package storage

import (
	"sync"
	"time"
)

type EventsRepository interface {
	Get(AppName string) (ApplicationEvent, error)
	Save(AppName string, event ApplicationEvent) error
}

type ApplicationEvent struct {
	AppName   string
	Image     string
	EventTime time.Time
	Status    string
}

type ApplicationsEvents struct {
	events map[string]ApplicationEvent // map[appName]ApplicationData
	sync.RWMutex
}

func NewApplicationsEvents() EventsRepository {
	return &ApplicationsEvents{
		events: make(map[string]ApplicationEvent),
	}
}

func (e *ApplicationsEvents) Get(AppName string) (ApplicationEvent, error) {
	e.RLock()
	defer e.RUnlock()

	if appEvent, ok := e.events[AppName]; ok {
		delete(e.events, AppName)
		
		return appEvent, nil
	}

	return ApplicationEvent{}, NotFoundError
}

func (e *ApplicationsEvents) Save(AppName string, event ApplicationEvent) error {
	e.events[AppName] = event
	return nil
}
