package storage

import (
	"errors"
	"sync"
)

type StateRepository interface {
	GetLastApplied(appName string) (State, error)
	Save(appName string, state State) error
}

var NotFoundError = errors.New("item not found")

type State struct {
	Image string
}

type ApplicationsStates struct {
	applications map[string]*Stack // map[appName][]State
	sync.RWMutex
}

func NewApplicationsStates() StateRepository {
	return &ApplicationsStates{
		applications: make(map[string]*Stack),
	}
}

func (s *ApplicationsStates) GetLastApplied(appName string) (State, error) {
	s.RLock()
	defer s.RUnlock()

	if _, ok := s.applications[appName]; !ok {
		return State{}, NotFoundError
	}

	stateStack := s.applications[appName]

	return stateStack.Pop()
}

func (s *ApplicationsStates) Save(appName string, state State) error {
	s.Lock()
	defer s.Unlock()

	if _, ok := s.applications[appName]; !ok {
		stack := &Stack{
			Sates: make([]State, 0),
		}
		stack.Push(state)
		s.applications[appName] = stack
	}

	stateStack := s.applications[appName]
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
