package domain

import (
	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/core/system"
)

// ActorFactory builds the actors that own maps and state machines. The actor package implements it,
// since that package depends on this one.
type ActorFactory interface {
	system.System
	NewMapActor(m *Map) actor.Actor
	NewStateMachineActor(sm *StateMachine) actor.Actor
}

// MapCallRunner is implemented by actors that own maps, so a map call made from Lua on the target actor
// itself runs in place instead of waiting on its own mailbox.
type MapCallRunner interface {
	RunMapCall(ctx actor.Context, msg *MapCall)
	RunMapCallAsync(ctx actor.Context, msg *MapCallAsync)
}
