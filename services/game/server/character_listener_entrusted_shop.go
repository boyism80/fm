package server

import (
	"context"
	"math"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/core/async"
	pconst "github.com/boyism80/fm/protocol/constant"
	"github.com/boyism80/fm/protocol/dto"
	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
	"github.com/boyism80/fm/protocol/response"
	"github.com/boyism80/fm/services/game/constant"
	"github.com/boyism80/fm/services/game/entity"
	"github.com/boyism80/fm/types"
)

func (l *CharacterListenerImpl) FindEntrustedShopAsync(ctx actor.Context, ch *entity.Character) *async.Promise[*internal.FindEntrustedShopReply] {
	req := &internal.FindEntrustedShopRequest{
		WorldId:   l.gs.config.WorldId,
		AccountId: ch.AccountID,
	}
	return async.NewTask(ctx, core.InternalRPCPerStepTimeout).ThenRPC(
		func(c context.Context) (*internal.FindEntrustedShopReply, error) {
			return l.gs.internalClient.FindEntrustedShop(c, req)
		},
		nil,
	)
}

func (l *CharacterListenerImpl) OpenShopAsync(ctx actor.Context, ch *entity.Character, shop *internal.Shop) *async.Promise[*internal.OpenShopReply] {
	shop.WorldId = l.gs.config.WorldId
	shop.ChannelId = int32(l.gs.config.ChannelId)
	req := &internal.OpenShopRequest{Shop: shop}
	return async.NewTask(ctx, core.InternalRPCPerStepTimeout).ThenRPC(
		func(c context.Context) (*internal.OpenShopReply, error) {
			return l.gs.internalClient.OpenShop(c, req)
		},
		nil,
	)
}

func (l *CharacterListenerImpl) OnEntrustedShopCheck(ch *entity.Character, result pconst.EntrustedShopCheck, mapID uint32, channel uint8) {
	ch.Send(&response.EntrustedShopCheckResult{
		Result:  result,
		MapID:   mapID,
		Channel: channel,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnEntrustedShopRemoteVisit(ch *entity.Character, sn uint32) {
	ch.Send(&response.EntrustedShopCheckResult{
		Result: pconst.EntrustedShopRemoteVisit,
		SN:     sn,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) miniRoomItems(items []*entity.ShopItem) response.MiniRoomItemList {
	out := make(response.MiniRoomItemList, 0, len(items))
	for _, listed := range items {
		out = append(out, response.MiniRoomItem{
			Bundles:   listed.Bundles,
			PerBundle: listed.PerBundle,
			Price:     listed.Price,
			Item:      listed.Item.ToDTO(),
		})
	}
	return out
}

func (l *CharacterListenerImpl) OnMiniRoomEntered(ch *entity.Character, slot uint8, view entity.EntrustedShopView, firstTime bool) {
	enter := &response.MiniRoomEnter{
		MySlot:       slot,
		PermitItemID: view.ItemID,
		OwnerName:    view.OwnerName,
		Elapsed:      uint32(view.Elapsed.Milliseconds()),
		FirstTime:    firstTime,
		Title:        view.Title,
		MaxItems:     entity.EntrustedShopMaxItems,
		MiniRoomItems: response.MiniRoomItems{
			Meso:  view.Meso,
			Items: l.miniRoomItems(view.Items),
		},
	}
	for i, visitor := range view.Visitors {
		if visitor == nil {
			continue
		}
		enter.Visitors = append(enter.Visitors, response.MiniRoomVisitor{
			Slot:      uint8(i + 1),
			Character: visitor,
		})
	}
	for _, sale := range view.Sold {
		enter.Sold = append(enter.Sold, response.MiniRoomSale{
			ItemID:  sale.ItemID,
			Bundles: sale.Bundles,
			Total:   sale.Total,
			Buyer:   sale.Buyer,
		})
	}
	ch.Send(enter, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnMiniRoomEnterFailed(ch *entity.Character, code pconst.MiniRoomEnterError) {
	ch.Send(&response.MiniRoomEnterFailed{Error: code}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnMiniRoomVisited(ch *entity.Character, slot uint8, visitor *entity.Character) {
	ch.Send(&response.MiniRoomVisited{
		MiniRoomVisitor: response.MiniRoomVisitor{
			Slot:      slot,
			Character: visitor.ToDTO(),
		},
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnMiniRoomChat(ch *entity.Character, slot uint8, text string) {
	ch.Send(&response.MiniRoomChat{
		Slot:    slot,
		Message: text,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnMiniRoomLeft(ch *entity.Character, slot uint8, reason pconst.MiniRoomLeaveReason) {
	ch.Send(&response.MiniRoomLeave{
		Slot:   slot,
		Reason: reason,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnMiniRoomItems(ch *entity.Character, view entity.EntrustedShopView) {
	ch.Send(&response.MiniRoomItems{
		Meso:  view.Meso,
		Items: l.miniRoomItems(view.Items),
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnMiniRoomBuyFailed(ch *entity.Character, result pconst.MiniRoomBuyResult) {
	ch.Send(&response.MiniRoomBuyFailed{Result: result}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnMiniRoomArranged(ch *entity.Character, view entity.EntrustedShopView) {
	ch.Send(&response.MiniRoomArranged{Meso: view.Meso}, types.SEND_POLICY_ENCRYPT)
	l.OnMiniRoomItems(ch, view)
}

func (l *CharacterListenerImpl) OnMiniRoomClosed(ch *entity.Character, result pconst.MiniRoomCloseResult) {
	ch.Send(&response.MiniRoomClosed{Result: result}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnMiniRoomMesoWithdrawn(ch *entity.Character) {
	ch.Send(&response.MiniRoomMesoWithdrawn{}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnOpenStoreBank(ch *entity.Character, npcID uint32) {
	meso := int32(min(ch.StoreBank.Meso, math.MaxInt32))
	tabs := map[constant.InventoryType][]dto.Item{}
	for invType := constant.InventoryTypeEquipment; invType <= constant.InventoryTypeCash; invType++ {
		tabs[invType] = []dto.Item{}
	}
	for _, listed := range ch.StoreBank.Items {
		invType := listed.Item.GetInventoryType()
		tabs[invType] = append(tabs[invType], listed.Item.Clone(listed.Bundles*listed.PerBundle).ToDTO())
	}
	ch.Send(&response.StoreBankOpen{
		NpcID: npcID,
		Storage: response.Storage{
			Slots: entity.StoreBankSlots,
			Meso:  &meso,
			Tabs:  tabs,
		},
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnStoreBankFee(ch *entity.Character, days uint32, fee int32) {
	ch.Send(&response.StoreBankFee{
		Days: days,
		Fee:  fee,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnStoreBankLocation(ch *entity.Character, npcID uint32, mapID uint32, channel uint8) {
	ch.Send(&response.StoreBankLocation{
		NpcID:   npcID,
		MapID:   mapID,
		Channel: channel,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnStoreBankResult(ch *entity.Character, result pconst.StoreBankResult) {
	ch.Send(&response.StoreBankResult{Result: result}, types.SEND_POLICY_ENCRYPT)
}
