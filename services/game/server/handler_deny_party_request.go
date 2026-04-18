package server

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/core/async"
	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
	"github.com/boyism80/fm/protocol/request"
	"github.com/boyism80/fm/services/game/client"
)

type DenyPartyRequest struct {
	gs     *GameServer
	opcode byte
}

func (DenyPartyRequest) New(gs *GameServer) *DenyPartyRequest {
	return &DenyPartyRequest{
		gs:     gs,
		opcode: 0x67,
	}
}

func (h *DenyPartyRequest) GetOpcode() byte {
	return h.opcode
}

func (h *DenyPartyRequest) Handle(ctx *core.ClientContext, req *request.DenyPartyRequest) error {
	if h.gs == nil {
		return fmt.Errorf("deny party request: game server is nil")
	}
	if h.gs.internalClient == nil {
		return fmt.Errorf("deny party request: internal client not configured")
	}
	gameClient, ok := ctx.Client.(*client.GameClient)
	if !ok {
		return fmt.Errorf("deny party request: client is not a GameClient")
	}
	character := gameClient.GetCharacter()
	if character == nil {
		return fmt.Errorf("deny party request: character not found")
	}
	if partyID, inParty := character.GetPartyID(); inParty && partyID != 0 {
		log.Printf("DenyPartyRequest: ignored, character=%d already in party=%d", character.GetID(), partyID)
		return nil
	}

	inviterName := strings.TrimSpace(req.InviterName)
	if inviterName == "" {
		log.Printf("DenyPartyRequest: ignored empty inviter, character=%d", character.GetID())
		return nil
	}
	if ctx.ActorContext == nil {
		return fmt.Errorf("deny party request: actor context required")
	}
	worldID := h.gs.config.WorldId
	deniedCharacterID := character.GetID()
	async.ThenRPC(async.NewPromise(ctx.ActorContext, core.InternalRPCPerStepTimeout),
		func(c context.Context) (*internal.DenyPartyReply, error) {
			return h.gs.internalClient.DenyParty(c, &internal.DenyPartyRequest{
				WorldId:           worldID,
				DeniedCharacterId: deniedCharacterID,
				InviterName:       inviterName,
				Action:            uint32(req.Action),
			})
		},
		func(reply *internal.DenyPartyReply) error {
			if !reply.GetOk() {
				log.Printf("DenyPartyRequest: failed denied=%d inviter=%q code=%v", deniedCharacterID, inviterName, reply.GetErrorCode())
			}
			return nil
		},
	).OnError(func(err error) {
		log.Printf("DenyPartyRequest async error: %v", err)
	}).Run()

	log.Printf("DenyPartyRequest: action=%d inviter=%s denied_by=%s forwarded_internal=true", req.Action, inviterName, character.GetName())
	return nil
}
