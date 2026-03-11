package actor

import (
	"fmt"
	"sync"

	"github.com/asynkron/protoactor-go/actor"
)

type ActorRegistry struct {
	system *ActorSystem
	actors map[string]*actor.PID
	mutex  sync.RWMutex
}

func NewActorRegistry(system *ActorSystem) *ActorRegistry {
	return &ActorRegistry{
		system: system,
		actors: make(map[string]*actor.PID),
	}
}

func (r *ActorRegistry) GetOrCreateActor(name string, props *actor.Props) *actor.PID {
	r.mutex.RLock()
	if pid, exists := r.actors[name]; exists {
		r.mutex.RUnlock()
		return pid
	}
	r.mutex.RUnlock()

	r.mutex.Lock()
	defer r.mutex.Unlock()

	pid := r.system.root.Spawn(props)
	r.actors[name] = pid
	return pid
}

func (r *ActorRegistry) GetActor(name string) (*actor.PID, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	pid, exists := r.actors[name]
	if !exists {
		return nil, fmt.Errorf("actor %s not found", name)
	}
	return pid, nil
}

func (r *ActorRegistry) RemoveActor(name string) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	delete(r.actors, name)
}
