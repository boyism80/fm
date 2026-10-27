package server

import (
	"fmt"
	"log"

	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/protocol/request"
	"github.com/boyism80/fm/services/game/client"
)

type Aggro struct {
	gs *GameServer
}

func (Aggro) New(gs *GameServer) *Aggro {
	return &Aggro{
		gs: gs,
	}
}

func (h *Aggro) Handle(ctx *core.ClientContext, req *request.Aggro) error {
	client, ok := ctx.Client.(*client.GameClient)
	if !ok {
		log.Printf("Client is not a GameClient")
		return fmt.Errorf("client is not a GameClient")
	}
	ch := client.GetCharacter()
	if ch == nil {
		return fmt.Errorf("character is nil")
	}

	m := ch.GetMap()
	if m == nil {
		return nil
	}
	mob := m.GetMob(req.OID)
	if mob == nil {
		return nil
	}
	m.GetControllerTable().AutoAggro(mob, ch)
	return nil
}
