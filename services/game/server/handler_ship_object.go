package server

import (
	"fmt"
	"log"

	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/protocol/request"
	"github.com/boyism80/fm/protocol/response"
	"github.com/boyism80/fm/services/game/client"
)

type ShipObject struct {
	gs *GameServer
}

func (ShipObject) New(gs *GameServer) *ShipObject {
	return &ShipObject{gs: gs}
}

func (h *ShipObject) Handle(ctx *core.ClientContext, req *request.ShipObject) error {
	cl, ok := ctx.Client.(*client.GameClient)
	if !ok {
		log.Printf("Client is not a GameClient")
		return fmt.Errorf("client is not a GameClient")
	}
	ch := cl.GetCharacter()
	if ch == nil {
		return nil
	}

	groupName := ""
	wantBalrog := false
	switch req.MapID {
	case 101000300, 200000111:
		groupName = "Boats"
	case 200000121, 220000110:
		groupName = "Trains"
	case 200000151, 260000100:
		groupName = "Geenie"
	case 240000110, 200000131:
		groupName = "Flight"
	case 200090010, 200090000:
		groupName = "Boats"
		wantBalrog = true
	default:
		return nil
	}

	prop := h.groupProp(groupName)
	state := response.ShipStateLeaving
	switch {
	case wantBalrog:
		if prop("haveBalrog") != "true" {
			return nil
		}
		state = response.ShipSpecialBalrog
	case prop("docked") == "true":
		state = response.ShipStateDocked
	}
	ch.Listener.OnShipState(ch, state)
	return nil
}

func (h *ShipObject) groupProp(groupName string) func(string) string {
	return func(key string) string {
		if h.gs == nil {
			return ""
		}
		reg := h.gs.GetStateMachineRegistry()
		if reg == nil {
			return ""
		}
		group := reg.Get(groupName)
		if group == nil {
			return ""
		}
		return group.GetProperty(key)
	}
}
