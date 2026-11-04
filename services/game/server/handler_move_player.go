package server

import (
	"fmt"
	"log"

	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/protocol/dto"
	"github.com/boyism80/fm/protocol/request"
)

type MovePlayer struct {
	gs *GameServer
}

func (MovePlayer) New(gs *GameServer) *MovePlayer {
	return &MovePlayer{
		gs: gs,
	}
}

func (h *MovePlayer) Handle(ctx *core.ClientContext, req *request.MovePlayer) error {
	character := h.gs.character(ctx)
	if character == nil {
		return nil
	}
	// A move the client made on the field it left still carries that field's key.
	if req.FieldKey != character.FieldKey() {
		return nil
	}

	beforePosition := character.GetPosition()

	for _, frag := range req.Fragments {
		if move, ok := frag.(*dto.AbsoluteLifeMovement); ok {
			character.Place(move.Position)
		}
		character.SetStance(frag.GetStance())
	}

	mapInstance := character.GetMap()
	if mapInstance == nil {
		log.Printf("Map not found for character movement")
		return fmt.Errorf("map not found")
	}

	mapInstance.OnMoved(character, beforePosition)
	character.Listener().OnPlayerMove(character, beforePosition, req.Fragments)

	return nil
}
