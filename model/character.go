package model

import (
	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/handler"
)

type CharacterActor struct {
	Life
	PlayerID int64
	Class    string
	Exp      int
	handler  *handler.MessageHandler
}

func NewCharacterActor() actor.Actor {
	act := &CharacterActor{
		handler: handler.NewMessageHandler(),
	}
	RegisterCharacterHandlers(act, act.handler)
	return act
}

func (state *CharacterActor) Receive(context actor.Context) {
	state.handler.Handle(context)
}
