package actor

import (
	"github.com/asynkron/protoactor-go/actor"
)

type ActorSystem struct {
	system *actor.ActorSystem
	root   *actor.RootContext
}

func NewActorSystem() *ActorSystem {
	system := actor.NewActorSystem()
	root := system.Root
	return &ActorSystem{
		system: system,
		root:   root,
	}
}

func (as *ActorSystem) GetSystem() *actor.ActorSystem {
	return as.system
}

func (as *ActorSystem) GetRoot() *actor.RootContext {
	return as.root
}
