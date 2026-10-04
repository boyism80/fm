package server

import (
	"fmt"
	"log"

	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/protocol/request"
	"github.com/boyism80/fm/services/game/client"
	"github.com/boyism80/fm/types"
)

type UseInnerPortal struct {
	gs *GameServer
}

func (UseInnerPortal) New(gs *GameServer) *UseInnerPortal {
	return &UseInnerPortal{
		gs: gs,
	}
}

func (h *UseInnerPortal) Handle(ctx *core.ClientContext, req *request.UseInnerPortal) error {
	gameClient, ok := ctx.Client.(*client.GameClient)
	if !ok {
		return fmt.Errorf("client is not a GameClient")
	}
	ch := gameClient.GetCharacter()
	if ch == nil {
		return fmt.Errorf("character not found")
	}
	if err := ch.EnterInnerPortal(req.PortalName, types.Vector2[int16]{X: req.ToX, Y: req.ToY}); err != nil {
		log.Printf("UseInnerPortal: character=%d: %v", ch.GetID(), err)
	}
	return nil
}
