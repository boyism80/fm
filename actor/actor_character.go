package actor

import (
	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/handler"
	"github.com/boyism80/fm/model"
)

type CharacterActor struct {
	Character *model.Character
	handler   *handler.MessageHandler
}

func NewCharacterActor(character *model.Character) actor.Actor {
	act := &CharacterActor{
		Character: character,
		handler:   handler.NewMessageHandler(),
	}
	handler.RegisterCharacterHandlers(act.Character, act.handler)
	return act
}

func (state *CharacterActor) Receive(context actor.Context) {
	state.handler.Handle(context)
}
