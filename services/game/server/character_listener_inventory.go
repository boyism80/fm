package server

import (
	"github.com/boyism80/fm/maple/constant"
	"github.com/boyism80/fm/maple/wz"
	pconst "github.com/boyism80/fm/protocol/constant"
	"github.com/boyism80/fm/protocol/dto"
	"github.com/boyism80/fm/protocol/response"
	"github.com/boyism80/fm/services/game"
	"github.com/boyism80/fm/types"
)

func (l *CharacterListenerImpl) OnItemEffect(ch *game.Character) {
	ch.Broadcast(&response.ItemEffect{CharacterID: ch.GetID(), ItemID: ch.ItemEffectID()}, nil)
}

func (l *CharacterListenerImpl) OnOpenNpcShop(ch *game.Character, shopID uint32, shop *wz.Shop) {
	ch.Send(&response.OpenNpcShop{
		ShopID: int32(shopID),
		Shop:   shop,
		Items:  l.gs.GetResources().Items,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) storageTabs(storage *game.Storage, invTypes ...constant.InventoryType) map[constant.InventoryType][]dto.Item {
	tabs := make(map[constant.InventoryType][]dto.Item, len(invTypes))
	for _, invType := range invTypes {
		items := make([]dto.Item, 0, len(storage.Tabs()[invType]))
		for _, item := range storage.Tabs()[invType] {
			items = append(items, item.BuildS2C())
		}
		tabs[invType] = items
	}
	return tabs
}

func (l *CharacterListenerImpl) OnOpenStorage(ch *game.Character, npcID uint32) {
	meso := ch.Account().Storage().Meso()
	ch.Send(&response.Storage{
		Result: pconst.StorageResultOpen,
		NpcID:  npcID,
		Slots:  ch.Account().Storage().Slots(),
		Meso:   &meso,
		Tabs: l.storageTabs(ch.Account().Storage(),
			constant.InventoryTypeEquipment,
			constant.InventoryTypeConsume,
			constant.InventoryTypeInstallation,
			constant.InventoryTypeETC,
			constant.InventoryTypeCash,
		),
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnStorageTabChanged(ch *game.Character, result pconst.StorageResult, invType constant.InventoryType) {
	ch.Send(&response.Storage{
		Result: result,
		Slots:  ch.Account().Storage().Slots(),
		Tabs:   l.storageTabs(ch.Account().Storage(), invType),
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnStorageArranged(ch *game.Character) {
	ch.Send(&response.Storage{
		Result: pconst.StorageResultArrange,
		Slots:  ch.Account().Storage().Slots(),
		Tabs: l.storageTabs(ch.Account().Storage(),
			constant.InventoryTypeEquipment,
			constant.InventoryTypeConsume,
			constant.InventoryTypeInstallation,
			constant.InventoryTypeETC,
			constant.InventoryTypeCash,
		),
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnStorageMesoChanged(ch *game.Character) {
	meso := ch.Account().Storage().Meso()
	ch.Send(&response.Storage{
		Result: pconst.StorageResultMeso,
		Slots:  ch.Account().Storage().Slots(),
		Meso:   &meso,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnStorageError(ch *game.Character, result pconst.StorageResult) {
	ch.Send(&response.StorageError{
		Result: result,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnOpenDuey(ch *game.Character, fromArrival bool) {
	parcels, expired := ch.Duey().BuildS2C()
	ch.Send(&response.DueyOpen{
		FromArrival: fromArrival,
		Parcels:     parcels,
		Expired:     expired,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnDueyResult(ch *game.Character, result pconst.DueyResult) {
	ch.Send(&response.Duey{
		Result: result,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnDueyRemoved(ch *game.Character, parcelID uint32, reason uint8) {
	ch.Send(&response.DueyRemoved{
		ParcelID: parcelID,
		Reason:   reason,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnDueyArrival(ch *game.Character, sender string, quick bool, count int) {
	if count > 1 {
		ch.Send(&response.DueyArrivals{
			Quick: quick,
		}, types.SEND_POLICY_ENCRYPT)
		return
	}
	ch.Send(&response.DueyArrival{
		Sender: sender,
		Quick:  quick,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnItemGainFailed(ch *game.Character, mode constant.ItemGainFailedType) {
	ch.Send(&response.ItemGainFailed{
		Mode: mode,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnInventorySlotUpdated(ch *game.Character, inventoryType constant.InventoryType, slot int16, item game.Item) {
	var itemDTO dto.Item
	if item != nil {
		itemDTO = item.BuildS2C()
	}
	ch.Send(&response.UpdateInventorySlot{
		InventoryType: inventoryType,
		Slot:          slot,
		Item:          itemDTO,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnInventorySlotAdded(ch *game.Character, inventoryType constant.InventoryType, slot int16, item game.Item) {
	var itemDTO dto.Item
	if item != nil {
		itemDTO = item.BuildS2C()
	}
	ch.Send(&response.AddInventorySlot{
		InventoryType: inventoryType,
		Slot:          slot,
		Item:          itemDTO,
		FromDrop:      true,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnInventorySlotLimitChanged(ch *game.Character, inventoryType constant.InventoryType, limit uint8) {
	ch.Send(&response.InventorySlotLimit{InventoryType: inventoryType, Limit: limit}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnShowItemGain(ch *game.Character, itemId uint32, count uint32, mode constant.ShowItemGainType) {
	ch.Send(&response.ShowItemGain{
		ItemId: itemId,
		Count:  count,
		Mode:   mode,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnShowMesoGain(ch *game.Character, count int32, mode constant.ShowMesoGainType) {
	ch.Send(&response.ShowMesoGain{
		Count: count,
		Mode:  mode,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnItemMakerDestroyed(ch *game.Character) {
	ch.Send(&response.ShowSelfItemMakerSuccessEffect{}, types.SEND_POLICY_ENCRYPT)
	l.OnMessage(ch, constant.MsgPopup, "아이템이 파괴되었습니다.")
}

func (l *CharacterListenerImpl) OnUseSkillBook(ch *game.Character, skillID uint32, masterLevel int, canUse bool, success bool) {
	ch.Broadcast(&response.UseSkillBook{
		CharacterID: ch.GetID(),
		SkillID:     skillID,
		MasterLevel: uint32(masterLevel),
		CanUse:      canUse,
		Success:     success,
	}, &game.ObjectBroadcastOption{WithMe: true})
}

func (l *CharacterListenerImpl) OnScrolledItem(ch *game.Character, scrollInventoryType constant.InventoryType, scrollSlot int16, scrollCount uint16, upgradedSlot int16, destroyed bool, potential bool, upgradedItem game.Item) {
	var upgradedItemDTO dto.Item
	if upgradedItem != nil {
		upgradedItemDTO = upgradedItem.BuildS2C()
	}
	ch.Send(&response.ScrolledItem{
		ScrollInventoryType: scrollInventoryType,
		ScrollSlot:          scrollSlot,
		ScrollCount:         scrollCount,
		UpgradedSlot:        upgradedSlot,
		Destroyed:           destroyed,
		Potential:           potential,
		UpgradedItem:        upgradedItemDTO,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnEndSortInventory(ch *game.Character, inventoryType constant.InventoryType) {
	ch.Send(&response.EndSortInventory{
		InventoryType: inventoryType,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnSwapInventorySlot(ch *game.Character, inventoryType constant.InventoryType, source int16, dest int16, equipmentAction int8) {
	ch.Send(&response.SwapInventorySlot{
		InventoryType:   inventoryType,
		Source:          source,
		Dest:            dest,
		EquipmentAction: response.EquipmentActionType(equipmentAction),
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnRemoveInventorySlot(ch *game.Character, inventoryType constant.InventoryType, slot int16) {
	ch.Send(&response.RemoveInventorySlot{
		InventoryType: inventoryType,
		Slot:          slot,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnItemExpired(ch *game.Character, items []wz.Item) {
	var expired []uint32
	for _, model := range items {
		if model.GetTraits()&constant.ItemTraitCash != 0 {
			ch.Send(&response.CashItemExpired{ItemID: model.GetID()}, types.SEND_POLICY_ENCRYPT)
		} else {
			expired = append(expired, model.GetID())
		}
	}
	if len(expired) == 0 {
		return
	}
	ch.Send(&response.ItemExpired{ItemIDs: expired}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnFullMergeInventorySlot(ch *game.Character, inventoryType constant.InventoryType, source int16, dest int16, count uint16) {
	ch.Send(&response.FullMergeInventorySlot{
		InventoryType: inventoryType,
		Source:        source,
		Dest:          dest,
		Count:         count,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnPartialMergeInventorySlot(ch *game.Character, inventoryType constant.InventoryType, source int16, dest int16, sourceCount uint16, destCount uint16) {
	ch.Send(&response.PartialMergeInventorySlot{
		InventoryType: inventoryType,
		Source:        source,
		Dest:          dest,
		SourceCount:   sourceCount,
		DestCount:     destCount,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnConfirmShopTransaction(ch *game.Character, code uint8) {
	_ = ch.Send(&response.ConfirmShopTransaction{
		Code: code,
	}, types.SEND_POLICY_ENCRYPT)
}
