package server

import (
	"errors"
	"math"

	"github.com/boyism80/fm/maple/constant"
	pconst "github.com/boyism80/fm/protocol/constant"
	"github.com/boyism80/fm/protocol/dto"
	"github.com/boyism80/fm/protocol/response"
	"github.com/boyism80/fm/services/game"
	"github.com/boyism80/fm/types"
)

func (l *CharacterListenerImpl) OnEntrustedShopCheck(ch *game.Character, result pconst.EntrustedShopCheck, mapID uint32, channel uint8) {
	ch.Send(&response.EntrustedShopCheckResult{
		Result:  result,
		MapID:   mapID,
		Channel: channel,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnEntrustedShopRemoteVisit(ch *game.Character, sn uint32) {
	ch.Send(&response.EntrustedShopCheckResult{
		Result: pconst.EntrustedShopRemoteVisit,
		SN:     sn,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) miniRoomItems(items []*game.ShopItem) response.MiniRoomItemList {
	out := make(response.MiniRoomItemList, 0, len(items))
	for _, listed := range items {
		out = append(out, response.MiniRoomItem{
			Bundles:   listed.Bundles(),
			PerBundle: listed.PerBundle(),
			Price:     listed.Price(),
			Item:      listed.Item().BuildS2C(),
		})
	}
	return out
}

func (l *CharacterListenerImpl) OnMiniRoomEntered(ch *game.Character, slot uint8, view game.EntrustedShopView, firstTime bool) {
	enter := &response.MiniRoomEnter{
		MySlot:       slot,
		PermitItemID: view.ItemID(),
		OwnerName:    view.OwnerName(),
		Elapsed:      uint32(view.Elapsed().Milliseconds()),
		FirstTime:    firstTime,
		Title:        view.Title(),
		MaxItems:     game.EntrustedShopMaxItems,
		MiniRoomItems: response.MiniRoomItems{
			Meso:  view.Meso(),
			Items: l.miniRoomItems(view.Items()),
		},
	}
	for i, visitor := range view.Visitors() {
		if visitor == nil {
			continue
		}
		enter.Visitors = append(enter.Visitors, response.MiniRoomVisitor{
			Slot:      uint8(i + 1),
			Character: visitor.BuildS2C(),
		})
	}
	for _, sale := range view.Sold() {
		enter.Sold = append(enter.Sold, response.MiniRoomSale{
			ItemID:  sale.ItemID,
			Bundles: sale.Bundles,
			Total:   sale.Total,
			Buyer:   sale.Buyer,
		})
	}
	ch.Send(enter, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnMiniRoomEnterFailed(ch *game.Character, code pconst.MiniRoomEnterError) {
	ch.Send(&response.MiniRoomEnterFailed{Error: code}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnMiniRoomFailed(ch *game.Character, err error) {
	var enterErr *game.MiniRoomEnterError
	var buyErr *game.MiniRoomBuyError
	switch {
	case errors.As(err, &enterErr):
		l.OnMiniRoomEnterFailed(ch, enterErr.Code)
	case errors.As(err, &buyErr):
		l.OnMiniRoomBuyFailed(ch, buyErr.Result)
	case err != nil:
		l.OnUnlockAction(ch)
	}
}

func (l *CharacterListenerImpl) OnMiniRoomVisited(ch *game.Character, slot uint8, visitor *game.Character) {
	ch.Send(&response.MiniRoomVisited{
		MiniRoomVisitor: response.MiniRoomVisitor{
			Slot:      slot,
			Character: visitor.BuildS2C(),
		},
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnMiniRoomChat(ch *game.Character, slot uint8, text string) {
	ch.Send(&response.MiniRoomChat{
		Slot:    slot,
		Message: text,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnMiniRoomLeft(ch *game.Character, slot uint8, reason pconst.MiniRoomLeaveReason) {
	ch.Send(&response.MiniRoomLeave{
		Slot:   slot,
		Reason: reason,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnMiniRoomItems(ch *game.Character, view game.EntrustedShopView) {
	ch.Send(&response.MiniRoomItems{
		Meso:  view.Meso(),
		Items: l.miniRoomItems(view.Items()),
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnMiniRoomBuyFailed(ch *game.Character, result pconst.MiniRoomBuyResult) {
	ch.Send(&response.MiniRoomBuyFailed{Result: result}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnMiniRoomArranged(ch *game.Character, view game.EntrustedShopView) {
	ch.Send(&response.MiniRoomArranged{Meso: view.Meso()}, types.SEND_POLICY_ENCRYPT)
	l.OnMiniRoomItems(ch, view)
}

func (l *CharacterListenerImpl) OnMiniRoomClosed(ch *game.Character, result pconst.MiniRoomCloseResult) {
	ch.Send(&response.MiniRoomClosed{Result: result}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnMiniRoomMesoWithdrawn(ch *game.Character) {
	ch.Send(&response.MiniRoomMesoWithdrawn{}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnOpenStoreBank(ch *game.Character, npcID uint32) {
	meso := int32(min(ch.StoreBank().Meso(), math.MaxInt32))
	tabs := map[constant.InventoryType][]dto.Item{}
	for invType := constant.InventoryTypeEquipment; invType <= constant.InventoryTypeCash; invType++ {
		tabs[invType] = []dto.Item{}
	}
	for _, listed := range ch.StoreBank().Items() {
		invType := listed.Item().GetInventoryType()
		tabs[invType] = append(tabs[invType], listed.Item().Clone(listed.Bundles()*listed.PerBundle()).BuildS2C())
	}
	ch.Send(&response.StoreBankOpen{
		NpcID: npcID,
		Storage: response.Storage{
			Slots: game.StoreBankSlots,
			Meso:  &meso,
			Tabs:  tabs,
		},
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnStoreBankFee(ch *game.Character, days uint32, fee int32) {
	ch.Send(&response.StoreBankFee{
		Days: days,
		Fee:  fee,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnStoreBankLocation(ch *game.Character, npcID uint32, mapID uint32, channel uint8) {
	ch.Send(&response.StoreBankLocation{
		NpcID:   npcID,
		MapID:   mapID,
		Channel: channel,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnStoreBankResult(ch *game.Character, result pconst.StoreBankResult) {
	ch.Send(&response.StoreBankResult{Result: result}, types.SEND_POLICY_ENCRYPT)
}
