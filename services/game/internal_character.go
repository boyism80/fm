package game

import (
	"context"
	"fmt"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/core/async"
	"github.com/boyism80/fm/maple/constant"
	"github.com/boyism80/fm/protocol/dto"
	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
	"github.com/boyism80/fm/protocol/response"
	"github.com/boyism80/fm/stream"
)

func (s *InternalSystem) BroadcastNoticeAsync(ctx actor.Context, ch *Character, messageType constant.ServerMessageType, message string, ear bool) *async.Promise[*internal.BroadcastNoticeReply] {
	fail := func(err error) *async.Promise[*internal.BroadcastNoticeReply] {
		p := async.NewDeferred[*internal.BroadcastNoticeReply](ctx)
		p.SetError(err)
		return p
	}
	if ctx == nil {
		return fail(fmt.Errorf("broadcast notice unavailable"))
	}
	channelID := uint32(s.gw.GetChannelID())
	worldID := s.gw.GetWorldID()
	return async.NewTask(ctx, core.InternalRPCPerStepTimeout).ThenRPC(
		func(c context.Context) (*internal.BroadcastNoticeReply, error) {
			return s.client.BroadcastNotice(c, &internal.BroadcastNoticeRequest{
				WorldId:         worldID,
				SourceChannelId: channelID,
				MessageType:     uint32(messageType),
				Message:         message,
				MegaEar:         ear,
			})
		},
		func(reply *internal.BroadcastNoticeReply) error {
			if reply == nil || !reply.GetOk() {
				return fmt.Errorf("broadcast notice failed")
			}
			s.gw.Systems().Get[*NoticeSystem]().BroadcastNotice(messageType, message, int(channelID)+1, ear)
			return nil
		},
	)
}

func (s *InternalSystem) BroadcastMegaphoneAsync(ctx actor.Context, ch *Character, megaphone *response.Megaphone) *async.Promise[*internal.BroadcastNoticeReply] {
	if ctx == nil {
		p := async.NewDeferred[*internal.BroadcastNoticeReply](ctx)
		p.SetError(fmt.Errorf("broadcast megaphone unavailable"))
		return p
	}
	channelID := uint32(s.gw.GetChannelID())
	rpc := &internal.BroadcastNoticeRequest{
		WorldId:         s.gw.GetWorldID(),
		SourceChannelId: channelID,
		MessageType:     uint32(megaphone.Type),
		Message:         megaphone.Lines[0],
		MegaEar:         megaphone.Ear,
		NextLines:       megaphone.Lines[1:],
	}
	if megaphone.Item != nil {
		w := stream.NewStreamWriter(stream.LittleEndian)
		megaphone.Item.Serialize(w, dto.ItemSerializeOption{SlotMode: dto.SlotEncodeOmit})
		rpc.Item = w.Bytes()
	}
	return async.NewTask(ctx, core.InternalRPCPerStepTimeout).ThenRPC(
		func(c context.Context) (*internal.BroadcastNoticeReply, error) {
			return s.client.BroadcastNotice(c, rpc)
		},
		func(reply *internal.BroadcastNoticeReply) error {
			if reply == nil || !reply.GetOk() {
				return fmt.Errorf("broadcast megaphone failed")
			}
			megaphone.Channel = int(channelID) + 1
			s.gw.Systems().Get[*NoticeSystem]().BroadcastMegaphone(megaphone)
			return nil
		},
	)
}

func (s *InternalSystem) SetAccountRoleAsync(ctx actor.Context, ch *Character, role constant.CharacterRole) *async.Promise[*internal.SetAccountRoleReply] {
	if ctx == nil {
		p := async.NewDeferred[*internal.SetAccountRoleReply](ctx)
		p.SetError(fmt.Errorf("set account role unavailable"))
		return p
	}
	accountID := ch.accountID
	return async.NewTask(ctx, core.InternalRPCPerStepTimeout).ThenRPC(
		func(c context.Context) (*internal.SetAccountRoleReply, error) {
			return s.client.SetAccountRole(c, &internal.SetAccountRoleRequest{
				AccountId: accountID,
				Role:      uint32(role),
			})
		},
		func(reply *internal.SetAccountRoleReply) error {
			if !reply.GetOk() {
				return fmt.Errorf("set account role failed")
			}
			return nil
		},
	)
}

func (s *InternalSystem) LoadParcelsAsync(ctx actor.Context, ch *Character) *async.Promise[*internal.LoadParcelsReply] {
	req := &internal.LoadParcelsRequest{
		WorldId:     s.gw.GetWorldID(),
		CharacterId: ch.GetID(),
	}
	return async.NewTask(ctx, core.InternalRPCPerStepTimeout).ThenRPC(
		func(c context.Context) (*internal.LoadParcelsReply, error) {
			return s.client.LoadParcels(c, req)
		},
		nil,
	)
}

func (s *InternalSystem) SendParcelAsync(ctx actor.Context, ch *Character, requestID uint64, senderAccountID uint32, recipient string, parcel *internal.Parcel, oneOfAKind bool, system bool) *async.Promise[*internal.SendParcelReply] {
	req := &internal.SendParcelRequest{
		WorldId:         s.gw.GetWorldID(),
		RecipientName:   recipient,
		Parcel:          parcel,
		SenderAccountId: senderAccountID,
		OneOfAKind:      oneOfAKind,
		RequestId:       requestID,
		System:          system,
	}
	return async.NewTask(ctx, core.InternalRPCPerStepTimeout).ThenRPC(
		func(c context.Context) (*internal.SendParcelReply, error) {
			return s.client.SendParcel(c, req)
		},
		nil,
	)
}

func (s *InternalSystem) ReturnParcelAsync(ctx actor.Context, ch *Character, requestID uint64, parcels []*internal.Parcel) *async.Promise[*internal.ReturnParcelReply] {
	req := &internal.ReturnParcelRequest{
		WorldId:     s.gw.GetWorldID(),
		CharacterId: ch.GetID(),
		RequestId:   requestID,
		Parcels:     parcels,
	}
	return async.NewTask(ctx, core.InternalRPCPerStepTimeout).ThenRPC(
		func(c context.Context) (*internal.ReturnParcelReply, error) {
			return s.client.ReturnParcel(c, req)
		},
		nil,
	)
}

func (s *InternalSystem) AbortParcelSendAsync(ctx actor.Context, ch *Character, requestID uint64, recipient string) *async.Promise[*internal.AbortParcelSendReply] {
	req := &internal.AbortParcelSendRequest{
		WorldId:       s.gw.GetWorldID(),
		RecipientName: recipient,
		RequestId:     requestID,
	}
	return async.NewTask(ctx, core.InternalRPCPerStepTimeout).ThenRPC(
		func(c context.Context) (*internal.AbortParcelSendReply, error) {
			return s.client.AbortParcelSend(c, req)
		},
		nil,
	)
}

func (s *InternalSystem) AddCashAsync(ctx actor.Context, ch *Character, nxCash int32, maplePoint int32) *async.Promise[*internal.AddCashReply] {
	req := &internal.AddCashRequest{
		WorldId:    s.gw.GetWorldID(),
		AccountId:  ch.accountID,
		NxCash:     nxCash,
		MaplePoint: maplePoint,
	}
	return async.NewTask(ctx, core.InternalRPCPerStepTimeout).ThenRPC(
		func(c context.Context) (*internal.AddCashReply, error) {
			return s.client.AddCash(c, req)
		},
		nil,
	)
}

func (s *InternalSystem) CreateCashCouponsAsync(ctx actor.Context, ch *Character, kind internal.CashCouponKind, value uint32, count uint32) *async.Promise[*internal.CreateCashCouponsReply] {
	req := &internal.CreateCashCouponsRequest{
		WorldId: s.gw.GetWorldID(),
		Kind:    kind,
		Value:   value,
		Count:   count,
	}
	return async.NewTask(ctx, core.InternalRPCPerStepTimeout).ThenRPC(
		func(c context.Context) (*internal.CreateCashCouponsReply, error) {
			return s.client.CreateCashCoupons(c, req)
		},
		nil,
	)
}

func (s *InternalSystem) ClaimParcelAsync(ctx actor.Context, ch *Character, requestID uint64, parcelID uint32) *async.Promise[*internal.ClaimParcelReply] {
	req := &internal.ClaimParcelRequest{
		WorldId:     s.gw.GetWorldID(),
		CharacterId: ch.GetID(),
		ParcelId:    parcelID,
		RequestId:   requestID,
	}
	return async.NewTask(ctx, core.InternalRPCPerStepTimeout).ThenRPC(
		func(c context.Context) (*internal.ClaimParcelReply, error) {
			return s.client.ClaimParcel(c, req)
		},
		nil,
	)
}

func (s *InternalSystem) DeleteParcelAsync(ctx actor.Context, ch *Character, parcelID uint32) *async.Promise[*internal.DeleteParcelReply] {
	req := &internal.DeleteParcelRequest{
		WorldId:     s.gw.GetWorldID(),
		CharacterId: ch.GetID(),
		ParcelId:    parcelID,
	}
	return async.NewTask(ctx, core.InternalRPCPerStepTimeout).ThenRPC(
		func(c context.Context) (*internal.DeleteParcelReply, error) {
			return s.client.DeleteParcel(c, req)
		},
		nil,
	)
}

func (s *InternalSystem) CheckParcelArrivalsAsync(ctx actor.Context, ch *Character) *async.Promise[*internal.CheckParcelArrivalsReply] {
	req := &internal.CheckParcelArrivalsRequest{
		WorldId:     s.gw.GetWorldID(),
		CharacterId: ch.GetID(),
	}
	return async.NewTask(ctx, core.InternalRPCPerStepTimeout).ThenRPC(
		func(c context.Context) (*internal.CheckParcelArrivalsReply, error) {
			return s.client.CheckParcelArrivals(c, req)
		},
		nil,
	)
}
