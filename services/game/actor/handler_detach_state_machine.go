package actor

import "github.com/asynkron/protoactor-go/actor"

type DetachStateMachineHandler struct{}

func (DetachStateMachineHandler) New() *DetachStateMachineHandler {
	return &DetachStateMachineHandler{}
}

func (h *DetachStateMachineHandler) Handle(ctx actor.Context, a *MapActor, msg *DetachStateMachine) {
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
