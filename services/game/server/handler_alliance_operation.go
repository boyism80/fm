package server

import (
	"context"
	"fmt"
	"log"

	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/core/async"
	pconst "github.com/boyism80/fm/protocol/constant"
	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
	"github.com/boyism80/fm/protocol/request"
	"github.com/boyism80/fm/services/game/client"
	"github.com/boyism80/fm/services/game/entity"
)

type AllianceOperation struct {
	gs *GameServer
}

func (AllianceOperation) New(gs *GameServer) *AllianceOperation {
	return &AllianceOperation{
		gs: gs,
	}
}

func (h *AllianceOperation) Handle(ctx *core.ClientContext, req *request.AllianceOperation) error {
	gameClient, ok := ctx.Client.(*client.GameClient)
	if !ok {
		return fmt.Errorf("client is not a GameClient")
	}
	ch := gameClient.GetCharacter()
	if ch == nil {
		return fmt.Errorf("character not found")
	}

	switch req.Operation {
	case pconst.AllianceC2SCreate:
		return h.handleCreate(ctx, ch, req)
	default:
		return nil
	}
}

func (h *AllianceOperation) handleCreate(ctx *core.ClientContext, ch *entity.Character, req *request.AllianceOperation) error {
	if h.gs == nil || h.gs.internalClient == nil {
		return nil
	}
	partnerCharacterID, ok := h.gs.ValidateCreateAlliance(ch)
	if !ok || partnerCharacterID == 0 {
		return nil
	}

	worldID := h.gs.config.WorldId
	charID := ch.GetID()
	allianceName := req.AllianceName
	if allianceName == "" {
		return nil
	}

	promise := async.NewPromise(ctx.ActorContext, core.InternalRPCPerStepTimeout)
	promise = async.ThenRPC(promise, func(c context.Context) (*internal.CreateAllianceReply, error) {
		return h.gs.internalClient.CreateAlliance(c, &internal.CreateAllianceRequest{
			WorldId:            worldID,
			AllianceName:       allianceName,
			LeaderCharacterId:  charID,
			PartnerCharacterId: partnerCharacterID,
		})
	}, func(reply *internal.CreateAllianceReply) error {
		if !reply.GetOk() {
			log.Printf("AllianceOperation(create): failed character=%d code=%v", charID, reply.GetErrorCode())
			return nil
		}
		alliancePb := reply.GetAlliance()
		if alliancePb == nil || alliancePb.GetAllianceId() == 0 {
			log.Printf("AllianceOperation(create): ok but missing alliance character=%d", charID)
			return nil
		}
		h.gs.applyAllianceFromProto(alliancePb)
		h.gs.broadcastAllianceCreate(alliancePb)
		return nil
	})
	promise.OnError(func(err error) {
		log.Printf("AllianceOperation(create) async error: %v", err)
	}).Run()
	return nil
}
