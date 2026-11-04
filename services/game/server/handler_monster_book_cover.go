package server

import (
	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/protocol/request"
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
	character, err := h.gs.character(ctx)
	if err != nil {
		return nil
	}

	character.MonsterBook().SetCover(req.CardID)
	return nil
}
