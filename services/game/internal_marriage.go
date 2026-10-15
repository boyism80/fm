package game

import (
	"context"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/core/async"
	"github.com/boyism80/fm/core/clock"
	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
)

func (s *InternalSystem) CreateMarriageAsync(ctx actor.Context, groom *Character, bride *Character, requestID uint64, ringItemID uint32) *async.Promise[*internal.MarriageReply] {
	req := &internal.CreateMarriageRequest{
		WorldId:    s.gw.GetWorldID(),
		GroomId:    groom.GetID(),
		GroomName:  groom.GetName(),
		BrideId:    bride.GetID(),
		BrideName:  bride.GetName(),
		RingItemId: ringItemID,
		RequestId:  requestID,
	}
	return async.NewTask(ctx, core.InternalRPCPerStepTimeout).ThenRPC(
		func(c context.Context) (*internal.MarriageReply, error) {
			return s.client.CreateMarriage(c, req)
		},
		nil,
	)
}

func (s *InternalSystem) LoadMarriageAsync(ctx actor.Context, ch *Character) *async.Promise[*internal.MarriageReply] {
	req := &internal.LoadMarriageRequest{
		WorldId:     s.gw.GetWorldID(),
		CharacterId: ch.GetID(),
	}
	return async.NewTask(ctx, core.InternalRPCPerStepTimeout).ThenRPC(
		func(c context.Context) (*internal.MarriageReply, error) {
			return s.client.LoadMarriage(c, req)
		},
		nil,
	)
}

func (s *InternalSystem) GetMarriageAsync(ctx actor.Context, ch *Character, marriageID uint32) *async.Promise[*internal.MarriageReply] {
	req := &internal.GetMarriageRequest{
		WorldId:    s.gw.GetWorldID(),
		MarriageId: marriageID,
	}
	return async.NewTask(ctx, core.InternalRPCPerStepTimeout).ThenRPC(
		func(c context.Context) (*internal.MarriageReply, error) {
			return s.client.GetMarriage(c, req)
		},
		nil,
	)
}

func (s *InternalSystem) BreakEngagementAsync(ctx actor.Context, ch *Character, marriageID uint32) *async.Promise[*internal.MarriageReply] {
	req := &internal.BreakEngagementRequest{
		WorldId:     s.gw.GetWorldID(),
		MarriageId:  marriageID,
		CharacterId: ch.GetID(),
	}
	return async.NewTask(ctx, core.InternalRPCPerStepTimeout).ThenRPC(
		func(c context.Context) (*internal.MarriageReply, error) {
			return s.client.BreakEngagement(c, req)
		},
		nil,
	)
}

func (s *InternalSystem) ReserveWeddingAsync(ctx actor.Context, ch *Character, requestID uint64, marriageID uint32, ticketItemID uint32) *async.Promise[*internal.MarriageReply] {
	req := &internal.ReserveWeddingRequest{
		WorldId:      s.gw.GetWorldID(),
		MarriageId:   marriageID,
		TicketItemId: ticketItemID,
		RequestId:    requestID,
		CharacterId:  ch.GetID(),
	}
	return async.NewTask(ctx, core.InternalRPCPerStepTimeout).ThenRPC(
		func(c context.Context) (*internal.MarriageReply, error) {
			return s.client.ReserveWedding(c, req)
		},
		nil,
	)
}

func (s *InternalSystem) AbortMarriageChangeAsync(ctx actor.Context, ch *Character, requestID uint64) *async.Promise[*internal.AbortMarriageChangeReply] {
	req := &internal.AbortMarriageChangeRequest{
		WorldId:   s.gw.GetWorldID(),
		RequestId: requestID,
	}
	return async.NewTask(ctx, core.InternalRPCPerStepTimeout).ThenRPC(
		func(c context.Context) (*internal.AbortMarriageChangeReply, error) {
			return s.client.AbortMarriageChange(c, req)
		},
		nil,
	)
}

func (s *InternalSystem) SetWeddingWishlistAsync(ctx actor.Context, ch *Character, requestID uint64, marriageID uint32, wishes []string) *async.Promise[*internal.MarriageReply] {
	req := &internal.SetWeddingWishlistRequest{
		WorldId:     s.gw.GetWorldID(),
		MarriageId:  marriageID,
		CharacterId: ch.GetID(),
		Wishes:      wishes,
		RequestId:   requestID,
	}
	return async.NewTask(ctx, core.InternalRPCPerStepTimeout).ThenRPC(
		func(c context.Context) (*internal.MarriageReply, error) {
			return s.client.SetWeddingWishlist(c, req)
		},
		nil,
	)
}

func (s *InternalSystem) InviteWeddingGuestAsync(ctx actor.Context, ch *Character, requestID uint64, marriageID uint32, guestName string, parcel *internal.Parcel) *async.Promise[*internal.InviteWeddingGuestReply] {
	req := &internal.InviteWeddingGuestRequest{
		WorldId:    s.gw.GetWorldID(),
		MarriageId: marriageID,
		GuestName:  guestName,
		Parcel:     parcel,
		RequestId:  requestID,
	}
	return async.NewTask(ctx, core.InternalRPCPerStepTimeout).ThenRPC(
		func(c context.Context) (*internal.InviteWeddingGuestReply, error) {
			return s.client.InviteWeddingGuest(c, req)
		},
		nil,
	)
}

func (s *InternalSystem) RequestDivorceAsync(ctx actor.Context, ch *Character, marriageID uint32) *async.Promise[*internal.MarriageReply] {
	req := &internal.RequestDivorceRequest{
		WorldId:     s.gw.GetWorldID(),
		MarriageId:  marriageID,
		CharacterId: ch.GetID(),
		NowUnixMs:   uint64(clock.Now().UnixMilli()),
	}
	return async.NewTask(ctx, core.InternalRPCPerStepTimeout).ThenRPC(
		func(c context.Context) (*internal.MarriageReply, error) {
			return s.client.RequestDivorce(c, req)
		},
		nil,
	)
}

func (s *InternalSystem) CancelDivorceAsync(ctx actor.Context, ch *Character, marriageID uint32) *async.Promise[*internal.MarriageReply] {
	req := &internal.CancelDivorceRequest{
		WorldId:     s.gw.GetWorldID(),
		MarriageId:  marriageID,
		CharacterId: ch.GetID(),
	}
	return async.NewTask(ctx, core.InternalRPCPerStepTimeout).ThenRPC(
		func(c context.Context) (*internal.MarriageReply, error) {
			return s.client.CancelDivorce(c, req)
		},
		nil,
	)
}

func (s *InternalSystem) GiveWeddingGiftAsync(ctx actor.Context, ch *Character, requestID uint64, receiverID uint32, gift *internal.WeddingGift) *async.Promise[*internal.GiveWeddingGiftReply] {
	req := &internal.GiveWeddingGiftRequest{
		WorldId:    s.gw.GetWorldID(),
		ReceiverId: receiverID,
		Gift:       gift,
		RequestId:  requestID,
	}
	return async.NewTask(ctx, core.InternalRPCPerStepTimeout).ThenRPC(
		func(c context.Context) (*internal.GiveWeddingGiftReply, error) {
			return s.client.GiveWeddingGift(c, req)
		},
		nil,
	)
}

func (s *InternalSystem) AbortWeddingGiftAsync(ctx actor.Context, ch *Character, requestID uint64, receiverID uint32) *async.Promise[*internal.AbortWeddingGiftReply] {
	req := &internal.AbortWeddingGiftRequest{
		WorldId:    s.gw.GetWorldID(),
		ReceiverId: receiverID,
		RequestId:  requestID,
	}
	return async.NewTask(ctx, core.InternalRPCPerStepTimeout).ThenRPC(
		func(c context.Context) (*internal.AbortWeddingGiftReply, error) {
			return s.client.AbortWeddingGift(c, req)
		},
		nil,
	)
}

func (s *InternalSystem) LoadWeddingGiftsAsync(ctx actor.Context, ch *Character) *async.Promise[*internal.LoadWeddingGiftsReply] {
	req := &internal.LoadWeddingGiftsRequest{
		WorldId:     s.gw.GetWorldID(),
		CharacterId: ch.GetID(),
	}
	return async.NewTask(ctx, core.InternalRPCPerStepTimeout).ThenRPC(
		func(c context.Context) (*internal.LoadWeddingGiftsReply, error) {
			return s.client.LoadWeddingGifts(c, req)
		},
		nil,
	)
}

func (s *InternalSystem) ClaimWeddingGiftAsync(ctx actor.Context, ch *Character, requestID uint64, giftID uint32) *async.Promise[*internal.ClaimWeddingGiftReply] {
	req := &internal.ClaimWeddingGiftRequest{
		WorldId:     s.gw.GetWorldID(),
		CharacterId: ch.GetID(),
		GiftId:      giftID,
		RequestId:   requestID,
	}
	return async.NewTask(ctx, core.InternalRPCPerStepTimeout).ThenRPC(
		func(c context.Context) (*internal.ClaimWeddingGiftReply, error) {
			return s.client.ClaimWeddingGift(c, req)
		},
		nil,
	)
}

func (s *InternalSystem) RemoveNotesAsync(ctx actor.Context, ch *Character, requestID uint64, action internal.NotesAction, noteIDs []uint32) *async.Promise[*internal.NotesReply] {
	req := &internal.NotesRequest{
		WorldId:     s.gw.GetWorldID(),
		CharacterId: ch.GetID(),
		Action:      action,
		NoteIds:     noteIDs,
		RequestId:   requestID,
	}
	return async.NewTask(ctx, core.InternalRPCPerStepTimeout).ThenRPC(
		func(c context.Context) (*internal.NotesReply, error) {
			return s.client.Notes(c, req)
		},
		nil,
	)
}

func (s *InternalSystem) NotifySpouseMapAsync(ctx actor.Context, ch *Character, spouseID uint32, mapID uint32, reply bool) *async.Promise[*internal.NotifySpouseMapReply] {
	req := &internal.NotifySpouseMapRequest{
		WorldId:     s.gw.GetWorldID(),
		CharacterId: ch.GetID(),
		SpouseId:    spouseID,
		MapId:       mapID,
		Reply:       reply,
	}
	return async.NewTask(ctx, core.InternalRPCPerStepTimeout).ThenRPC(
		func(c context.Context) (*internal.NotifySpouseMapReply, error) {
			return s.client.NotifySpouseMap(c, req)
		},
		nil,
	)
}
