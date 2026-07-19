package actor

import "github.com/asynkron/protoactor-go/actor"

type AttachStateMachineHandler struct{}

func (AttachStateMachineHandler) New() *AttachStateMachineHandler {
	return &AttachStateMachineHandler{}
}

func (h *AttachStateMachineHandler) Handle(ctx actor.Context, a *MapActor, msg *AttachStateMachine) {
	if msg == nil || msg.ReplyTo == nil || msg.StateMachine == nil {
		return
	}
	ack := &AttachStateMachineAck{MapID: msg.MapID}
	if a.Map != nil && a.Map.GetMapID() == msg.MapID {
		if err := a.Map.AttachStateMachine(msg.StateMachine); err == nil {
			ack.OK = true
		}
	}
	ctx.Send(msg.ReplyTo, ack)
}
