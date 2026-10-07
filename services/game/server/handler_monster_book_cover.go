package server

import (
	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/protocol/request"
	"github.com/boyism80/fm/services/game/client"
)

type MonsterBookCover struct {
	gs *GameServer
}

func (MonsterBookCover) New(gs *GameServer) *MonsterBookCover {
	return &MonsterBookCover{
		gs: gs,
	}
}

func (h *MonsterBookCover) Handle(ctx *core.ClientContext, req *request.MonsterBookCover) error {
	gameClient, ok := ctx.Client.(*client.GameClient)
	if !ok {
		return nil
	}
	character := gameClient.GetCharacter()
	if character == nil {
		return nil
	}

	character.SetMonsterBookCover(req.CardID)
	return nil
}
