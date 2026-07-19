package actor

import (
	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/services/game/entity"
)

type MapActor struct {
	GameLogicActor
	Map *entity.Map
}

func NewMapActor(m *entity.Map, gameWorld entity.GameWorld) *MapActor {
	a := &MapActor{Map: m}
	a.GameWorld = gameWorld
	a.maps = func() []*entity.Map {
		if a.Map == nil || a.Map.StateMachine() != nil {
			return nil
		}
		return []*entity.Map{a.Map}
	}
	return a
}

func (a *MapActor) Receive(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case *actor.Stopped:
		if a.Map != nil {
			a.Map.ClearLuaRoot()
		}
	case *AttachStateMachine:
		a.handleAttach(ctx, msg)
	case *DetachStateMachine:
		a.handleDetach(ctx, msg)
	default:
		a.GameLogicActor.Receive(ctx)
	}
}

func (a *MapActor) handleAttach(ctx actor.Context, msg *AttachStateMachine) {
	if msg == nil || msg.ReplyTo == nil || msg.StateMachine == nil {
		return
	}
	ack := &AttachStateMachineAck{MapID: msg.MapID}
	if a.Map != nil && a.Map.GetMapID() == msg.MapID {
		if err := a.Map.AttachStateMachine(msg.StateMachine); err == nil {
			ack.OK = true
			a.Map.RebindObjectTimers(msg.StateMachine.ActorPID)
		}
	}
	ctx.Send(msg.ReplyTo, ack)
}

func (a *MapActor) handleDetach(ctx actor.Context, msg *DetachStateMachine) {
	if msg == nil {
		return
	}
	if a.Map != nil && a.Map.GetMapID() == msg.MapID && msg.StateMachine != nil {
		a.Map.DetachStateMachine(msg.StateMachine)
		if a.Map.StateMachine() == nil {
			a.Map.RebindObjectTimers(ctx.Self())
		}
	}
	if msg.ReplyTo != nil {
		ctx.Send(msg.ReplyTo, &DetachStateMachineAck{MapID: msg.MapID})
	}
}
