package system

import (
	"fmt"
	"reflect"
	"slices"
)

type System interface {
	Start() error
	Stop()
}

// Systems are added while the server is built and only read after that, so the maps need no lock.
// A system shared by entity and server code is added once under each type it is fetched by.
type Systems struct {
	byType map[reflect.Type]System
	order  []System
}

func NewSystems() *Systems {
	return &Systems{byType: make(map[reflect.Type]System)}
}

func (s *Systems) Add[T System](v T) {
	key := reflect.TypeFor[T]()
	if _, ok := s.byType[key]; ok {
		panic(fmt.Sprintf("system %v added twice", key))
	}
	s.byType[key] = v
	if !slices.Contains(s.order, System(v)) {
		s.order = append(s.order, v)
	}
}

func (s *Systems) Get[T System]() T {
	v, ok := s.byType[reflect.TypeFor[T]()]
	if !ok {
		panic(fmt.Sprintf("system %v not added", reflect.TypeFor[T]()))
	}
	return v.(T)
}

// Start runs in the order the systems were added, so a system starts after the systems it uses.
func (s *Systems) Start() error {
	for _, v := range s.order {
		if err := v.Start(); err != nil {
			return err
		}
	}
	return nil
}

// Stop runs in reverse order, so a system stops before the systems it uses.
func (s *Systems) Stop() {
	for i := len(s.order) - 1; i >= 0; i-- {
		s.order[i].Stop()
	}
}

// Listeners are added while the server is built and only read after that, so the map needs no lock.
type Listeners struct {
	byType map[reflect.Type]any
}

func NewListeners() *Listeners {
	return &Listeners{byType: make(map[reflect.Type]any)}
}

func (l *Listeners) Add[T any](v T) {
	key := reflect.TypeFor[T]()
	if _, ok := l.byType[key]; ok {
		panic(fmt.Sprintf("listener %v added twice", key))
	}
	l.byType[key] = v
}

func (l *Listeners) Get[T any]() T {
	v, ok := l.byType[reflect.TypeFor[T]()]
	if !ok {
		panic(fmt.Sprintf("listener %v not added", reflect.TypeFor[T]()))
	}
	return v.(T)
}
