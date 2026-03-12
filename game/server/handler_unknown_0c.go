package server

import (
	"log"

	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/game/client"
	"github.com/boyism80/fm/protocol/request"
)

// Unknown0C handles client packet 0x0C and logs the raw payload for inspection.
type Unknown0C struct {
	gs     *GameServer
	opcode byte
}

func (Unknown0C) New(gs *GameServer) *Unknown0C {
	return &Unknown0C{
		gs:     gs,
		opcode: 0x0C,
	}
}

func (h *Unknown0C) GetOpcode() byte {
	return h.opcode
}

func (h *Unknown0C) Handle(ctx *core.ClientContext, req *request.Unknown0C) error {
	gameClient, ok := ctx.Client.(*client.GameClient)
	remoteAddr := ctx.Client.GetConnection().RemoteAddr().String()
	charID := uint32(0)
	if ok {
		if ch := gameClient.GetCharacter(); ch != nil {
			charID = ch.GetID()
		}
	}
	if len(req.Payload) > 0 {
		log.Printf("[0x0C] from %s char=%d type=%d paramA=%d paramB=%d extra=%x", remoteAddr, charID, req.Type, req.ParamA, req.ParamB, req.Payload)
	} else {
		log.Printf("[0x0C] from %s char=%d type=%d paramA=%d paramB=%d", remoteAddr, charID, req.Type, req.ParamA, req.ParamB)
	}
	return nil
}
