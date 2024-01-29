package storage

import (
	"errors"
	"sync"
)

type StateRepository interface {
	GetLastSuccessState(clusterName, appName string) (State, error)
	Save(clusterName, appName string, state State) error
}

var NotFoundError = errors.New("item not found")

type State struct {
	Image string
}

type ApplicationsStates struct {
	applications map[string]map[string]*Stack // map[clusterName][appName][]State
	sync.Mutex
}

func NewApplicationsStates() StateRepository {
	return &ApplicationsStates{
		applications: make(map[string]map[string]*Stack),
	}
}

func (s *ApplicationsStates) GetLastSuccessState(clusterName, appName string) (State, error) {
	s.Lock()
	defer s.Unlock()

	if _, ok := s.applications[clusterName]; !ok {
		return State{}, NotFoundError
	}

	if _, ok := s.applications[clusterName][appName]; !ok {
		return State{}, NotFoundError
	}

	return s.applications[clusterName][appName].Pop()
}

func (s *ApplicationsStates) Save(clusterName, appName string, state State) error {
	s.Lock()
	defer s.Unlock()

	if _, ok := s.applications[clusterName]; !ok {
		s.applications[clusterName] = make(map[string]*Stack)
	}

	if _, ok := s.applications[appName]; !ok {
		stack := &Stack{
			Sates: make([]State, 0),
		}
		stack.Push(state)
		s.applications[clusterName][appName] = stack
	}

	stateStack := s.applications[clusterName][appName]
	stateStack.Push(state)

	return nil
}

type Stack struct {
	Sates []State
	sync.Mutex
}

func (st *Stack) Push(v State) {
	st.Lock()
	defer st.Unlock()

	// do not store more then 10 artifacts
	if len(st.Sates) > 10 {
		_, _ = st.tail()
	}

	st.Sates = append(st.Sates, v)
}

func (st *Stack) Pop() (State, error) {
	st.Lock()
	defer st.Unlock()

	if len(st.Sates) == 0 {
		return State{}, NotFoundError
	}

	ret := (st.Sates)[len(st.Sates)-1]
	st.Sates = (st.Sates)[0 : len(st.Sates)-1]

	return ret, nil
}

func (st *Stack) tail() (State, error) {
	st.Lock()
	defer st.Unlock()

	if len(st.Sates) == 0 {
		return State{}, NotFoundError
	}

	ret := st.Sates[0]
	st.Sates = (st.Sates)[1:len(st.Sates)]

	return ret, nil
}
