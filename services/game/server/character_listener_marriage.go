package server

import (
	"context"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/core/async"
	"github.com/boyism80/fm/core/clock"
	pconst "github.com/boyism80/fm/protocol/constant"
	"github.com/boyism80/fm/protocol/dto"
	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
	"github.com/boyism80/fm/protocol/response"
	"github.com/boyism80/fm/services/game/constant"
	"github.com/boyism80/fm/services/game/entity"
	"github.com/boyism80/fm/types"
)

func (l *CharacterListenerImpl) CreateMarriageAsync(ctx actor.Context, groom *entity.Character, bride *entity.Character, ringItemID uint32) *async.Promise[*internal.MarriageReply] {
	req := &internal.CreateMarriageRequest{
		WorldId:    l.gs.config.WorldId,
		GroomId:    groom.GetID(),
		GroomName:  groom.GetName(),
		BrideId:    bride.GetID(),
		BrideName:  bride.GetName(),
		RingItemId: ringItemID,
	}
	return async.NewTask(ctx, core.InternalRPCPerStepTimeout).ThenRPC(
		func(c context.Context) (*internal.MarriageReply, error) {
			return l.gs.internalClient.CreateMarriage(c, req)
		},
		nil,
	)
}

func (l *CharacterListenerImpl) LoadMarriageAsync(ctx actor.Context, ch *entity.Character) *async.Promise[*internal.MarriageReply] {
	req := &internal.LoadMarriageRequest{
		WorldId:     l.gs.config.WorldId,
		CharacterId: ch.GetID(),
	}
	return async.NewTask(ctx, core.InternalRPCPerStepTimeout).ThenRPC(
		func(c context.Context) (*internal.MarriageReply, error) {
			return l.gs.internalClient.LoadMarriage(c, req)
		},
		nil,
	)
}

func (l *CharacterListenerImpl) GetMarriageAsync(ctx actor.Context, ch *entity.Character, marriageID uint32) *async.Promise[*internal.MarriageReply] {
	req := &internal.GetMarriageRequest{
		WorldId:    l.gs.config.WorldId,
		MarriageId: marriageID,
	}
	return async.NewTask(ctx, core.InternalRPCPerStepTimeout).ThenRPC(
		func(c context.Context) (*internal.MarriageReply, error) {
			return l.gs.internalClient.GetMarriage(c, req)
		},
		nil,
	)
}

func (l *CharacterListenerImpl) BreakEngagementAsync(ctx actor.Context, ch *entity.Character, marriageID uint32) *async.Promise[*internal.MarriageReply] {
	req := &internal.BreakEngagementRequest{
		WorldId:     l.gs.config.WorldId,
		MarriageId:  marriageID,
		CharacterId: ch.GetID(),
	}
	return async.NewTask(ctx, core.InternalRPCPerStepTimeout).ThenRPC(
		func(c context.Context) (*internal.MarriageReply, error) {
			return l.gs.internalClient.BreakEngagement(c, req)
		},
		nil,
	)
}

func (l *CharacterListenerImpl) ReserveWeddingAsync(ctx actor.Context, ch *entity.Character, marriageID uint32, ticketItemID uint32) *async.Promise[*internal.MarriageReply] {
	req := &internal.ReserveWeddingRequest{
		WorldId:      l.gs.config.WorldId,
		MarriageId:   marriageID,
		TicketItemId: ticketItemID,
	}
	return async.NewTask(ctx, core.InternalRPCPerStepTimeout).ThenRPC(
		func(c context.Context) (*internal.MarriageReply, error) {
			return l.gs.internalClient.ReserveWedding(c, req)
		},
		nil,
	)
}

func (l *CharacterListenerImpl) SetWeddingWishlistAsync(ctx actor.Context, ch *entity.Character, marriageID uint32, wishes []string) *async.Promise[*internal.MarriageReply] {
	req := &internal.SetWeddingWishlistRequest{
		WorldId:     l.gs.config.WorldId,
		MarriageId:  marriageID,
		CharacterId: ch.GetID(),
		Wishes:      wishes,
	}
	return async.NewTask(ctx, core.InternalRPCPerStepTimeout).ThenRPC(
		func(c context.Context) (*internal.MarriageReply, error) {
			return l.gs.internalClient.SetWeddingWishlist(c, req)
		},
		nil,
	)
}

func (l *CharacterListenerImpl) InviteWeddingGuestAsync(ctx actor.Context, ch *entity.Character, marriageID uint32, guestName string) *async.Promise[*internal.InviteWeddingGuestReply] {
	req := &internal.InviteWeddingGuestRequest{
		WorldId:    l.gs.config.WorldId,
		MarriageId: marriageID,
		GuestName:  guestName,
	}
	return async.NewTask(ctx, core.InternalRPCPerStepTimeout).ThenRPC(
		func(c context.Context) (*internal.InviteWeddingGuestReply, error) {
			return l.gs.internalClient.InviteWeddingGuest(c, req)
		},
		nil,
	)
}

func (l *CharacterListenerImpl) RequestDivorceAsync(ctx actor.Context, ch *entity.Character, marriageID uint32) *async.Promise[*internal.MarriageReply] {
	req := &internal.RequestDivorceRequest{
		WorldId:     l.gs.config.WorldId,
		MarriageId:  marriageID,
		CharacterId: ch.GetID(),
		NowUnixMs:   uint64(clock.Now().UnixMilli()),
	}
	return async.NewTask(ctx, core.InternalRPCPerStepTimeout).ThenRPC(
		func(c context.Context) (*internal.MarriageReply, error) {
			return l.gs.internalClient.RequestDivorce(c, req)
		},
		nil,
	)
}

func (l *CharacterListenerImpl) CancelDivorceAsync(ctx actor.Context, ch *entity.Character, marriageID uint32) *async.Promise[*internal.MarriageReply] {
	req := &internal.CancelDivorceRequest{
		WorldId:     l.gs.config.WorldId,
		MarriageId:  marriageID,
		CharacterId: ch.GetID(),
	}
	return async.NewTask(ctx, core.InternalRPCPerStepTimeout).ThenRPC(
		func(c context.Context) (*internal.MarriageReply, error) {
			return l.gs.internalClient.CancelDivorce(c, req)
		},
		nil,
	)
}

func (l *CharacterListenerImpl) GiveWeddingGiftAsync(ctx actor.Context, ch *entity.Character, receiverID uint32, gift *internal.WeddingGift, sender *internal.CharacterSaveEntry) *async.Promise[*internal.GiveWeddingGiftReply] {
	req := &internal.GiveWeddingGiftRequest{
		WorldId:    l.gs.config.WorldId,
		ReceiverId: receiverID,
		Gift:       gift,
		Sender:     sender,
	}
	return async.NewTask(ctx, core.InternalRPCPerStepTimeout).ThenRPC(
		func(c context.Context) (*internal.GiveWeddingGiftReply, error) {
			return l.gs.internalClient.GiveWeddingGift(c, req)
		},
		nil,
	)
}

func (l *CharacterListenerImpl) LoadWeddingGiftsAsync(ctx actor.Context, ch *entity.Character) *async.Promise[*internal.LoadWeddingGiftsReply] {
	req := &internal.LoadWeddingGiftsRequest{
		WorldId:     l.gs.config.WorldId,
		CharacterId: ch.GetID(),
	}
	return async.NewTask(ctx, core.InternalRPCPerStepTimeout).ThenRPC(
		func(c context.Context) (*internal.LoadWeddingGiftsReply, error) {
			return l.gs.internalClient.LoadWeddingGifts(c, req)
		},
		nil,
	)
}

func (l *CharacterListenerImpl) ClaimWeddingGiftAsync(ctx actor.Context, ch *entity.Character, giftID uint32) *async.Promise[*internal.ClaimWeddingGiftReply] {
	req := &internal.ClaimWeddingGiftRequest{
		WorldId:     l.gs.config.WorldId,
		CharacterId: ch.GetID(),
		GiftId:      giftID,
	}
	return async.NewTask(ctx, core.InternalRPCPerStepTimeout).ThenRPC(
		func(c context.Context) (*internal.ClaimWeddingGiftReply, error) {
			return l.gs.internalClient.ClaimWeddingGift(c, req)
		},
		nil,
	)
}

func (l *CharacterListenerImpl) NotifySpouseMapAsync(ctx actor.Context, ch *entity.Character, spouseID uint32, mapID uint32, reply bool) *async.Promise[*internal.NotifySpouseMapReply] {
	req := &internal.NotifySpouseMapRequest{
		WorldId:     l.gs.config.WorldId,
		CharacterId: ch.GetID(),
		SpouseId:    spouseID,
		MapId:       mapID,
		Reply:       reply,
	}
	return async.NewTask(ctx, core.InternalRPCPerStepTimeout).ThenRPC(
		func(c context.Context) (*internal.NotifySpouseMapReply, error) {
			return l.gs.internalClient.NotifySpouseMap(c, req)
		},
		nil,
	)
}

func (l *CharacterListenerImpl) OnEngageRequest(ch *entity.Character, name string, characterID uint32) {
	ch.Send(&response.EngageRequest{
		Mode:        pconst.EngageRequestPropose,
		Name:        name,
		CharacterID: characterID,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnWeddingWishlistInput(ch *entity.Character) {
	ch.Send(&response.EngageRequest{
		Mode: pconst.EngageRequestWishlist,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnEngageResult(ch *entity.Character, result pconst.EngageResult) {
	packet := &response.EngageResult{Result: result}
	if result == pconst.EngageResultEngaged || result == pconst.EngageResultMarried {
		packet.Marriage = ch.Wedding.Marriage.ToDTO()
	}
	ch.Send(packet, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnWeddingInvitation(ch *entity.Character, groomName string, brideName string, weddingType uint16) {
	ch.Send(&response.EngageResult{
		Result:      pconst.EngageResultInvitation,
		GroomName:   groomName,
		BrideName:   brideName,
		WeddingType: weddingType,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnWeddingGift(ch *entity.Character, mode pconst.WeddingGiftMode, wishes []string, tabs map[constant.InventoryType][]entity.Item) {
	packet := &response.WeddingGift{Mode: mode, Wishes: wishes}
	if tabs != nil {
		packet.Tabs = response.WeddingGiftTabs{}
		for typ, items := range tabs {
			dtos := make([]dto.Item, 0, len(items))
			for _, item := range items {
				dtos = append(dtos, entity.ItemToDTO(item))
			}
			packet.Tabs[typ] = dtos
		}
	}
	ch.Send(packet, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnSpouseMap(ch *entity.Character, mapID uint32, spouseID uint32) {
	ch.Send(&response.SpouseMap{
		MapID:       mapID,
		CharacterID: spouseID,
	}, types.SEND_POLICY_ENCRYPT)
}
