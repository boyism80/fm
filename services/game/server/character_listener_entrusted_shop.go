package server

import (
	"math"

	pconst "github.com/boyism80/fm/protocol/constant"
	"github.com/boyism80/fm/protocol/dto"
	"github.com/boyism80/fm/protocol/response"
	"github.com/boyism80/fm/services/game/constant"
	"github.com/boyism80/fm/services/game/domain"
	"github.com/boyism80/fm/types"
)

func (l *CharacterListenerImpl) OnEntrustedShopCheck(ch *domain.Character, result pconst.EntrustedShopCheck, mapID uint32, channel uint8) {
	ch.Send(&response.EntrustedShopCheckResult{
		Result:  result,
		MapID:   mapID,
		Channel: channel,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnEntrustedShopRemoteVisit(ch *domain.Character, sn uint32) {
	ch.Send(&response.EntrustedShopCheckResult{
		Result: pconst.EntrustedShopRemoteVisit,
		SN:     sn,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) miniRoomItems(items []*domain.ShopItem) response.MiniRoomItemList {
	out := make(response.MiniRoomItemList, 0, len(items))
	for _, listed := range items {
		out = append(out, response.MiniRoomItem{
			Bundles:   listed.Bundles,
			PerBundle: listed.PerBundle,
			Price:     listed.Price,
			Item:      listed.Item.BuildS2C(),
		})
	}
	return out
}

func (l *CharacterListenerImpl) OnMiniRoomEntered(ch *domain.Character, slot uint8, view domain.EntrustedShopView, firstTime bool) {
	enter := &response.MiniRoomEnter{
		MySlot:       slot,
		PermitItemID: view.ItemID,
		OwnerName:    view.OwnerName,
		Elapsed:      uint32(view.Elapsed.Milliseconds()),
		FirstTime:    firstTime,
		Title:        view.Title,
		MaxItems:     domain.EntrustedShopMaxItems,
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

func (l *CharacterListenerImpl) OnMiniRoomEnterFailed(ch *domain.Character, code pconst.MiniRoomEnterError) {
	ch.Send(&response.MiniRoomEnterFailed{Error: code}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnMiniRoomVisited(ch *domain.Character, slot uint8, visitor *domain.Character) {
	ch.Send(&response.MiniRoomVisited{
		MiniRoomVisitor: response.MiniRoomVisitor{
			Slot:      slot,
			Character: visitor.BuildS2C(),
		},
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnMiniRoomChat(ch *domain.Character, slot uint8, text string) {
	ch.Send(&response.MiniRoomChat{
		Slot:    slot,
		Message: text,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnMiniRoomLeft(ch *domain.Character, slot uint8, reason pconst.MiniRoomLeaveReason) {
	ch.Send(&response.MiniRoomLeave{
		Slot:   slot,
		Reason: reason,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnMiniRoomItems(ch *domain.Character, view domain.EntrustedShopView) {
	ch.Send(&response.MiniRoomItems{
		Meso:  view.Meso,
		Items: l.miniRoomItems(view.Items),
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnMiniRoomBuyFailed(ch *domain.Character, result pconst.MiniRoomBuyResult) {
	ch.Send(&response.MiniRoomBuyFailed{Result: result}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnMiniRoomArranged(ch *domain.Character, view domain.EntrustedShopView) {
	ch.Send(&response.MiniRoomArranged{Meso: view.Meso}, types.SEND_POLICY_ENCRYPT)
	l.OnMiniRoomItems(ch, view)
}

func (l *CharacterListenerImpl) OnMiniRoomClosed(ch *domain.Character, result pconst.MiniRoomCloseResult) {
	ch.Send(&response.MiniRoomClosed{Result: result}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnMiniRoomMesoWithdrawn(ch *domain.Character) {
	ch.Send(&response.MiniRoomMesoWithdrawn{}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnOpenStoreBank(ch *domain.Character, npcID uint32) {
	meso := int32(min(ch.StoreBank.Meso, math.MaxInt32))
	tabs := map[constant.InventoryType][]dto.Item{}
	for invType := constant.InventoryTypeEquipment; invType <= constant.InventoryTypeCash; invType++ {
		tabs[invType] = []dto.Item{}
	}
	for _, listed := range ch.StoreBank.Items {
		invType := listed.Item.GetInventoryType()
		tabs[invType] = append(tabs[invType], listed.Item.Clone(listed.Bundles*listed.PerBundle).BuildS2C())
	}
	ch.Send(&response.StoreBankOpen{
		NpcID: npcID,
		Storage: response.Storage{
			Slots: domain.StoreBankSlots,
			Meso:  &meso,
			Tabs:  tabs,
		},
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnStoreBankFee(ch *domain.Character, days uint32, fee int32) {
	ch.Send(&response.StoreBankFee{
		Days: days,
		Fee:  fee,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnStoreBankLocation(ch *domain.Character, npcID uint32, mapID uint32, channel uint8) {
	ch.Send(&response.StoreBankLocation{
		NpcID:   npcID,
		MapID:   mapID,
		Channel: channel,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnStoreBankResult(ch *domain.Character, result pconst.StoreBankResult) {
	ch.Send(&response.StoreBankResult{Result: result}, types.SEND_POLICY_ENCRYPT)
}
