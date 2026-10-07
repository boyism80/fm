package entity

import (
	"slices"
	"time"

	"github.com/boyism80/fm/protocol/dto"
	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
	"github.com/boyism80/fm/services/game/constant"
	gentity "github.com/boyism80/fm/services/game/entity"
	"github.com/boyism80/fm/services/game/wz"
	"github.com/boyism80/fm/util"
)

const defaultSlotLimit = 32

type Character struct {
	Game           *internal.EnterGameReply
	ReturnChannel  uint32
	NXCash         uint32
	MaplePoint     uint32
	Locker         []*internal.CashItem
	Wishlist       []uint32
	CharacterSlots uint16
	StorageSlots   uint16
	Gifts          []*internal.CashGift
	Expired        []uint64
	Busy           bool
}

func NewCharacter(reply *internal.EnterCashShopReply) *Character {
	return &Character{
		Game:           reply.GetGame(),
		ReturnChannel:  reply.GetReturnChannelId(),
		NXCash:         reply.GetNxCash(),
		MaplePoint:     reply.GetMaplePoint(),
		Locker:         reply.GetLocker(),
		Wishlist:       reply.GetWishlist(),
		CharacterSlots: uint16(reply.GetCharacterSlotCount()),
		StorageSlots:   uint16(reply.GetStorageSlotCount()),
		Gifts:          reply.GetGifts(),
		Expired:        reply.GetExpired(),
	}
}

func (ch *Character) ID() uint32 {
	return ch.Game.GetCharacter().GetCharacterId()
}

func (ch *Character) AccountID() uint32 {
	return ch.Game.GetCharacter().GetAccountId()
}

func (ch *Character) Name() string {
	return ch.Game.GetCharacter().GetName()
}

func (ch *Character) Gender() uint8 {
	return uint8(ch.Game.GetCharacter().GetGender())
}

func (ch *Character) Meso() int32 {
	return ch.Game.GetCharacter().GetMeso()
}

func (ch *Character) SlotLimit(invType constant.InventoryType) uint8 {
	limits := ch.Game.GetCharacter().GetSlotLimits()
	index := int(invType) - 1
	if index < 0 || index >= len(limits) {
		return defaultSlotLimit
	}
	return uint8(limits[index])
}

func (ch *Character) SetSlotLimit(invType constant.InventoryType, limit uint8) {
	p := ch.Game.GetCharacter()
	for len(p.SlotLimits) < int(constant.InventoryTypeCash) {
		p.SlotLimits = append(p.SlotLimits, defaultSlotLimit)
	}
	p.SlotLimits[invType-1] = uint32(limit)
}

func (ch *Character) RemoveLocker(serial uint64) {
	ch.Locker = slices.DeleteFunc(ch.Locker, func(item *internal.CashItem) bool {
		return item.GetItem().GetUniqueId() == serial
	})
}

func (ch *Character) FindLocker(serial uint64) *internal.CashItem {
	for _, item := range ch.Locker {
		if item.GetItem().GetUniqueId() == serial {
			return item
		}
	}
	return nil
}

func (ch *Character) Balance(currency uint8) (internal.CashCurrency, uint32) {
	if currency == 1 {
		return internal.CashCurrency_CASH_CURRENCY_MAPLE_POINT, ch.MaplePoint
	}
	return internal.CashCurrency_CASH_CURRENCY_NX_CASH, ch.NXCash
}

func (ch *Character) FindInventory(serial uint64, invType constant.InventoryType) (int, *internal.Inventory) {
	for i, item := range ch.Game.Inventory {
		if item.GetSlot() > 0 && item.GetUniqueId() == serial && ch.inventoryType(item) == invType {
			return i, item
		}
	}
	return -1, nil
}

func (ch *Character) FreeSlot(invType constant.InventoryType) (int16, bool) {
	used := make(map[int16]bool)
	for _, item := range ch.Game.Inventory {
		if item.GetSlot() > 0 && ch.inventoryType(item) == invType {
			used[int16(item.GetSlot())] = true
		}
	}
	for slot := int16(1); slot <= int16(ch.SlotLimit(invType)); slot++ {
		if used[slot] == false {
			return slot, true
		}
	}
	return 0, false
}

func (ch *Character) inventoryType(item *internal.Inventory) constant.InventoryType {
	if item.GetInventoryType() != 0 {
		return constant.InventoryType(item.GetInventoryType())
	}
	return constant.GetInventoryTypeByItemID(item.GetItemId())
}

func (ch *Character) NewCashItem(commodity *wz.Commodity, world gentity.ItemWorld) (*internal.CashItem, error) {
	item, err := gentity.NewItem(commodity.ItemID, max(commodity.Count, 1), world)
	if err != nil {
		return nil, err
	}

	pb := item.ToProto(ch.ID(), 0)
	if pb.UniqueId == nil {
		serial := world.NewUniqueID()
		pb.UniqueId = &serial
	}
	if commodity.Period > 0 {
		pb.ExpirationUnixMs = time.Now().AddDate(0, 0, int(commodity.Period)).UnixMilli()
	}
	return &internal.CashItem{CommoditySn: commodity.SN, BuyerName: ch.Name(), Item: pb}, nil
}

func (ch *Character) LockerItemDTO(item *internal.CashItem) *dto.CashShopItem {
	expiration := util.TimeMax
	if exp := item.GetItem().GetExpirationUnixMs(); exp > 0 {
		expiration = time.UnixMilli(exp)
	}
	return &dto.CashShopItem{
		Serial:      item.GetItem().GetUniqueId(),
		AccountID:   ch.AccountID(),
		ItemID:      item.GetItem().GetItemId(),
		CommoditySN: item.GetCommoditySn(),
		Count:       uint16(item.GetItem().GetCount()),
		BuyerName:   item.GetBuyerName(),
		Expiration:  expiration,
	}
}

func (ch *Character) LockerDTO() []*dto.CashShopItem {
	items := make([]*dto.CashShopItem, 0, len(ch.Locker))
	for _, item := range ch.Locker {
		items = append(items, ch.LockerItemDTO(item))
	}
	return items
}

func (ch *Character) ToDTO(world gentity.ItemWorld) *dto.Character {
	p := ch.Game.GetCharacter()
	character := &dto.Character{
		ID:            p.GetCharacterId(),
		Name:          p.GetName(),
		Gender:        uint8(p.GetGender()),
		SkinColor:     uint8(p.GetSkinColor()),
		Face:          p.GetFace(),
		Hair:          p.GetHair(),
		Pet:           p.GetSummonedPet(),
		Level:         uint8(p.GetLevel()),
		Class:         uint16(p.GetClassId()),
		Str:           uint16(p.GetStr()),
		Dex:           uint16(p.GetDex()),
		Int:           uint16(p.GetIntStat()),
		Luk:           uint16(p.GetLuk()),
		Hp:            uint16(p.GetHp()),
		MaxHp:         uint16(p.GetMaxHp()),
		Mp:            uint16(p.GetMp()),
		MaxMp:         uint16(p.GetMaxMp()),
		AbilityPoint:  uint16(p.GetAbilityPoint()),
		SkillPoint:    uint16(p.GetSkillPoint()),
		Exp:           p.GetExp(),
		Population:    uint16(p.GetPopulation()),
		Map:           p.GetMapId(),
		SpawnPoint:    uint8(p.GetSpawnPoint()),
		BuddyCapacity: uint8(min(ch.Game.GetBuddyCapacity(), 255)),
		RegRocks:      []uint32{999999999, 999999999, 999999999, 999999999, 999999999},
		Rocks:         []uint32{999999999, 999999999, 999999999, 999999999, 999999999, 999999999, 999999999, 999999999, 999999999, 999999999},
		Inventory: &dto.Inventory{
			Tabs:     make(map[constant.InventoryType]*dto.ItemContainer),
			Equipped: make(map[constant.EquipmentPartsType]*dto.Equipment),
			Meso:     p.GetMeso(),
		},
	}

	for _, invType := range []constant.InventoryType{constant.InventoryTypeEquipment, constant.InventoryTypeConsume, constant.InventoryTypeInstallation, constant.InventoryTypeETC, constant.InventoryTypeCash} {
		character.Inventory.Tabs[invType] = &dto.ItemContainer{Type: invType, SlotLimit: ch.SlotLimit(invType), Items: make(map[int16]dto.Item)}
	}
	for _, pb := range ch.Game.Inventory {
		item, err := gentity.NewItemFromInternalProto(pb, world)
		if err != nil {
			continue
		}
		if pb.GetSlot() < 0 {
			if equipment, ok := item.(gentity.Equipment); ok {
				character.Inventory.Equipped[constant.EquipmentPartsType(pb.GetSlot())] = equipment.ToEquipmentDTO()
			}
			continue
		}
		if tab := character.Inventory.Tabs[ch.inventoryType(pb)]; tab != nil {
			tab.Items[int16(pb.GetSlot())] = item.ToDTO()
		}
	}

	for _, skill := range ch.Game.Skills {
		character.Skills = append(character.Skills, &dto.Skill{
			ID:          skill.GetSkillId(),
			SkillLevel:  uint32(skill.GetLevel()),
			MasterLevel: uint32(skill.GetMasterLevel()),
		})
	}

	resources := world.GetResources()
	for _, quest := range ch.Game.Quests {
		if quest.GetQuestId() > 0xFFFF {
			continue
		}
		status := &dto.QuestStatus{
			QuestID:        uint16(quest.GetQuestId()),
			Status:         uint8(quest.GetStatus()),
			StatusRecord:   quest.GetStatusRecord(),
			CompletionTime: time.UnixMilli(quest.GetCompletionTimeUnixMs()),
		}
		if deadline := quest.GetDeadlineUnixMs(); deadline > 0 {
			status.Deadline = time.UnixMilli(deadline)
		}
		if wzQuest := resources.Quests[quest.GetQuestId()]; wzQuest != nil {
			for _, mobID := range wzQuest.OrderedMobIDs() {
				status.MobKills = append(status.MobKills, uint16(quest.GetMobKills()[mobID]))
			}
		}
		switch wz.QuestStatus(quest.GetStatus()) {
		case wz.QuestStatusStarted:
			character.QuestsStarted = append(character.QuestsStarted, status)
		case wz.QuestStatusCompleted:
			character.QuestsCompleted = append(character.QuestsCompleted, status)
		}
	}
	return character
}
