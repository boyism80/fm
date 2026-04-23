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

type CancelChair struct {
	gs *GameServer
}

func (CancelChair) New(gs *GameServer) *CancelChair {
	return &CancelChair{
		gs: gs,
	}
}

func (h *CancelChair) Handle(ctx *core.ClientContext, req *request.CancelChair) error {
	client, ok := ctx.Client.(*client.GameClient)
	if !ok {
		log.Printf("Client is not a GameClient")
		return fmt.Errorf("client is not a GameClient")
	}

	character := client.GetCharacter()
	if character == nil {
		return nil
	}

	mapInstance := character.GetMap()

	if req.ChairID == -1 {
		character.Chair = 0

		cancelChairPacket := &response.CancelChair{
			ChairID: -1,
		}
		character.Send(cancelChairPacket, types.SEND_POLICY_ENCRYPT)

		if mapInstance != nil {
			character.Broadcast(&response.ShowChair{
				CharacterID: character.GetID(),
				ItemID:      0,
			}, nil)
		}
	} else {
		character.Chair = uint32(req.ChairID)

		cancelChairPacket := &response.CancelChair{
			ChairID: req.ChairID,
		}
		character.Send(cancelChairPacket, types.SEND_POLICY_ENCRYPT)
	}

	character.Listener.OnUpdateStats(character, nil, true)

	return nil
}
