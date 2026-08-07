package server

import (
	"fmt"
	"log"

	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/protocol/request"
	"github.com/boyism80/fm/protocol/response"
	"github.com/boyism80/fm/services/game/client"
	"github.com/boyism80/fm/types"
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

	pkt, ok := h.resolveShipPacket(req.MapID)
	if !ok {
		return nil
	}
	ch.Send(pkt, types.SEND_POLICY_ENCRYPT)
	return nil
}

func (h *ShipObject) resolveShipPacket(mapID uint32) (types.Packet, bool) {
	groupName := ""
	wantBalrog := false
	switch mapID {
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
		return nil, false
	}

	prop := h.groupProp(groupName)
	if wantBalrog {
		if prop("haveBalrog") != "true" {
			return nil, false
		}
		return &response.ShipSpecialEffect{Effect: response.ShipSpecialBalrog}, true
	}
	state := response.ShipStateLeaving
	if prop("docked") == "true" {
		state = response.ShipStateDocked
	}
	return &response.ShipState{State: state}, true
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
