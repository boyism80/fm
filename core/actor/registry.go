package actor

import (
	"fmt"
	"sync"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/core/fault"
)

type ActorRegistry struct {
	system *ActorSystem
	actors map[string]*actor.PID
	mutex  sync.RWMutex
	Faults *fault.Injector
}

func NewActorRegistry(system *ActorSystem) *ActorRegistry {
	return &ActorRegistry{
		system: system,
		actors: make(map[string]*actor.PID),
		Faults: fault.NewInjector(),
	}
}

func (r *ActorRegistry) Start() error { return nil }
func (r *ActorRegistry) Stop()        { r.system.GetSystem().Shutdown() }

func (r *ActorRegistry) System() *actor.ActorSystem { return r.system.GetSystem() }
func (r *ActorRegistry) Root() *actor.RootContext   { return r.system.GetRoot() }

func (r *ActorRegistry) Running(pid *actor.PID) bool {
	if pid == nil {
		return false
	}
	_, ok := r.System().ProcessRegistry.Get(pid)
	return ok
}

func (r *ActorRegistry) PredictPID(name string) *actor.PID {
	return actor.NewPID(r.system.GetSystem().Address(), name)
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

	if pid, exists := r.actors[name]; exists {
		return pid
	}

	props.Configure(actor.WithReceiverMiddleware(r.Faults.ReceiverMiddleware(name)))
	pid, err := r.system.root.SpawnNamed(props, name)
	if err != nil {
		panic(fmt.Sprintf("spawn actor %s: %v", name, err))
	}
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

func (r *ActorRegistry) StopActor(name string, pid *actor.PID) {
	r.mutex.Lock()
	delete(r.actors, name)
	r.mutex.Unlock()
	if pid != nil {
		r.system.root.Stop(pid)
	}
}

func (r *ActorRegistry) PoisonActor(name string, pid *actor.PID) {
	r.mutex.Lock()
	delete(r.actors, name)
	r.mutex.Unlock()
	if pid != nil {
		r.system.root.Poison(pid)
	}
}
