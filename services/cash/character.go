package cash

import (
	"slices"
	"time"

	"github.com/boyism80/fm/maple/constant"
	"github.com/boyism80/fm/maple/item"
	"github.com/boyism80/fm/maple/wz"
	"github.com/boyism80/fm/protocol/dto"
	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
	"github.com/boyism80/fm/util"
)

const defaultSlotLimit = 32

type ItemWorld interface {
	GetResources() *wz.Resources
	NewUniqueID() uint64
}

type Character struct {
	id             uint32
	accountID      uint32
	name           string
	gender         uint8
	skin           uint8
	face           uint32
	hair           uint32
	summonedPet    uint64
	level          uint8
	class          uint16
	str            uint16
	dex            uint16
	intStat        uint16
	luk            uint16
	hp             uint16
	maxHp          uint16
	mp             uint16
	maxMp          uint16
	ap             uint16
	sp             uint16
	exp            uint32
	population     int16
	mapID          uint32
	spawnPoint     uint8
	meso           int32
	buddyCapacity  uint8
	teleportStones []uint32
	vipStones      []uint32
	slotLimits     []uint8
	inventory      []ItemRecord
	skills         []skillRecord
	quests         []questRecord
	debuffs        []debuffRecord

	ReturnChannel  uint32
	NXCash         uint32
	MaplePoint     uint32
	locker         []*LockerItem
	Wishlist       []uint32
	CharacterSlots uint16
	StorageSlots   uint16
	gifts          []giftRecord
	Expired        []uint64
	Busy           bool
	resources      *wz.Resources
}

func NewCharacter(reply *internal.EnterCashShopReply, resources *wz.Resources) *Character {
	game := reply.GetGame()
	p := game.GetCharacter()
	ch := &Character{
		id:             p.GetCharacterId(),
		accountID:      p.GetAccountId(),
		name:           p.GetName(),
		gender:         uint8(p.GetGender()),
		skin:           uint8(p.GetSkinColor()),
		face:           p.GetFace(),
		hair:           p.GetHair(),
		summonedPet:    p.GetSummonedPet(),
		level:          uint8(p.GetLevel()),
		class:          uint16(p.GetClassId()),
		str:            uint16(p.GetStr()),
		dex:            uint16(p.GetDex()),
		intStat:        uint16(p.GetIntStat()),
		luk:            uint16(p.GetLuk()),
		hp:             uint16(p.GetHp()),
		maxHp:          uint16(p.GetMaxHp()),
		mp:             uint16(p.GetMp()),
		maxMp:          uint16(p.GetMaxMp()),
		ap:             uint16(p.GetAbilityPoint()),
		sp:             uint16(p.GetSkillPoint()),
		exp:            p.GetExp(),
		population:     int16(p.GetPopulation()),
		mapID:          p.GetMapId(),
		spawnPoint:     uint8(p.GetSpawnPoint()),
		meso:           p.GetMeso(),
		buddyCapacity:  uint8(min(game.GetBuddyCapacity(), 255)),
		teleportStones: append([]uint32(nil), p.GetTeleportStones()...),
		vipStones:      append([]uint32(nil), p.GetVipTeleportStones()...),
		ReturnChannel:  reply.GetReturnChannelId(),
		NXCash:         reply.GetNxCash(),
		MaplePoint:     reply.GetMaplePoint(),
		Wishlist:       reply.GetWishlist(),
		CharacterSlots: uint16(reply.GetCharacterSlotCount()),
		StorageSlots:   uint16(reply.GetStorageSlotCount()),
		Expired:        reply.GetExpired(),
		resources:      resources,
	}
	for _, limit := range p.GetSlotLimits() {
		ch.slotLimits = append(ch.slotLimits, uint8(limit))
	}
	for _, pb := range game.GetInventory() {
		ch.inventory = append(ch.inventory, ItemFromProto(pb))
	}
	for _, skill := range game.GetSkills() {
		ch.skills = append(ch.skills, skillRecord{id: skill.GetSkillId(), level: skill.GetLevel(), masterLevel: skill.GetMasterLevel()})
	}
	for _, quest := range game.GetQuests() {
		ch.quests = append(ch.quests, questRecord{
			id: quest.GetQuestId(), status: quest.GetStatus(), mobKills: quest.GetMobKills(),
			statusRecord: quest.GetStatusRecord(), completion: quest.GetCompletionTimeUnixMs(), deadline: quest.GetDeadlineUnixMs(),
		})
	}
	for _, debuff := range game.GetDebuffs() {
		ch.debuffs = append(ch.debuffs, debuffRecord{
			mask: debuff.GetMask(), position: debuff.GetPosition(), x: debuff.GetX(),
			skillID: debuff.GetSkillId(), skillLevel: debuff.GetSkillLevel(), end: debuff.GetEndUnixMs(),
		})
	}
	for _, pb := range reply.GetLocker() {
		ch.locker = append(ch.locker, LockerItemFromProto(pb))
	}
	for _, pb := range reply.GetGifts() {
		ch.gifts = append(ch.gifts, giftFromProto(pb))
	}
	return ch
}

func (ch *Character) ID() uint32        { return ch.id }
func (ch *Character) AccountID() uint32 { return ch.accountID }
func (ch *Character) Name() string      { return ch.name }
func (ch *Character) Gender() uint8     { return ch.gender }
func (ch *Character) Meso() int32       { return ch.meso }

func (ch *Character) SlotLimit(invType constant.InventoryType) uint8 {
	index := int(invType) - 1
	if index < 0 || index >= len(ch.slotLimits) {
		return defaultSlotLimit
	}
	return ch.slotLimits[index]
}

func (ch *Character) SetSlotLimit(invType constant.InventoryType, limit uint8) {
	for len(ch.slotLimits) < int(constant.InventoryTypeCash) {
		ch.slotLimits = append(ch.slotLimits, defaultSlotLimit)
	}
	ch.slotLimits[invType-1] = limit
}

func (ch *Character) AddLocker(items ...*LockerItem) {
	ch.locker = append(ch.locker, items...)
}

func (ch *Character) RemoveLocker(serial uint64) {
	ch.locker = slices.DeleteFunc(ch.locker, func(item *LockerItem) bool {
		return item.Serial() == serial
	})
}

func (ch *Character) AddInventory(item ItemRecord) {
	ch.inventory = append(ch.inventory, item)
}

func (ch *Character) TakeOut(serial uint64, item ItemRecord) {
	ch.RemoveLocker(serial)
	ch.AddInventory(item)
}

func (ch *Character) PutIn(serial uint64, invType constant.InventoryType, item *LockerItem) {
	ch.inventory = slices.DeleteFunc(ch.inventory, func(pb ItemRecord) bool {
		return pb.slot > 0 && pb.Serial() == serial && ch.inventoryType(pb) == invType
	})
	ch.AddLocker(item)
}

func (ch *Character) SetBalance(nxCash, maplePoint uint32) {
	ch.NXCash = nxCash
	ch.MaplePoint = maplePoint
}

func (ch *Character) SetMeso(meso int32) {
	ch.meso = meso
}

func (ch *Character) FindLocker(serial uint64) *LockerItem {
	for _, item := range ch.locker {
		if item.Serial() == serial {
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

func (ch *Character) FindInventory(serial uint64, invType constant.InventoryType) (ItemRecord, bool) {
	for _, item := range ch.inventory {
		if item.slot > 0 && item.Serial() == serial && ch.inventoryType(item) == invType {
			return item.clone(), true
		}
	}
	return ItemRecord{}, false
}

func (ch *Character) FreeSlot(invType constant.InventoryType) (int16, bool) {
	used := make(map[int16]bool)
	for _, item := range ch.inventory {
		if item.slot > 0 && ch.inventoryType(item) == invType {
			used[int16(item.slot)] = true
		}
	}
	for slot := int16(1); slot <= int16(ch.SlotLimit(invType)); slot++ {
		if !used[slot] {
			return slot, true
		}
	}
	return 0, false
}

func (ch *Character) inventoryType(item ItemRecord) constant.InventoryType {
	if item.inventoryType != 0 {
		return constant.InventoryType(item.inventoryType)
	}
	model := ch.resources.Items[item.itemID]
	if model == nil {
		return 0
	}
	return model.GetInventoryType()
}

func (ch *Character) NewCashItem(commodity *wz.Commodity, world ItemWorld) (*LockerItem, error) {
	pb, err := item.NewInventory(commodity.ItemID, max(commodity.Count, 1), ch.ID(), 0, world.GetResources(), world.NewUniqueID)
	if err != nil {
		return nil, err
	}
	record := ItemFromProto(pb)
	if record.Serial() == 0 {
		record.SetUniqueID(world.NewUniqueID())
	}
	if commodity.Period > 0 {
		record.SetExpiration(time.Now().AddDate(0, 0, int(commodity.Period)).UnixMilli())
	}
	return &LockerItem{commoditySN: commodity.SN, buyerName: ch.Name(), item: record}, nil
}

func (ch *Character) PrepareTakeOut(serial uint64, slot int16, invType constant.InventoryType, world ItemWorld) (ItemRecord, bool) {
	found := ch.FindLocker(serial)
	if found == nil {
		return ItemRecord{}, false
	}
	model := ch.resources.Items[found.ItemID()]
	if model == nil {
		return ItemRecord{}, false
	}
	record := found.Item()
	record.SetOwnerID(ch.id)
	record.SetSlot(int32(slot))
	record.SetInventoryType(uint32(invType))
	switch m := model.(type) {
	case wz.Equipment, *wz.Installation:
		record.SetCount(1)
	case *wz.Pet:
		record.SetCount(1)
		record.SetFlag(0)
		record.SetOwnerName("")
		if record.Serial() == 0 {
			record.SetUniqueID(world.NewUniqueID())
		}
		record.ensurePet(world.GetResources().GetItemName(record.ItemID()), constant.PetMaxFullness, uint32(m.Skills), uint32(m.LimitedLife))
	}
	return record, true
}

func (ch *Character) PreparePutIn(serial uint64, invType constant.InventoryType) *LockerItem {
	record, ok := ch.FindInventory(serial, invType)
	if !ok {
		return nil
	}
	record.SetSlot(0)
	return &LockerItem{buyerName: ch.name, item: record}
}

func (ch *Character) DebuffsRPC() []*internal.Debuff {
	out := make([]*internal.Debuff, 0, len(ch.debuffs))
	for _, debuff := range ch.debuffs {
		out = append(out, debuff.BuildRPC())
	}
	return out
}

func (ch *Character) BuildLockerItemS2C(item *LockerItem) *dto.CashShopItem {
	expiration := util.TimeMax
	if exp := item.Expiration(); exp > 0 {
		expiration = time.UnixMilli(exp)
	}
	return &dto.CashShopItem{
		Serial:      item.Serial(),
		AccountID:   ch.AccountID(),
		ItemID:      item.ItemID(),
		CommoditySN: item.CommoditySN(),
		Count:       uint16(item.Count()),
		BuyerName:   item.BuyerName(),
		Expiration:  expiration,
	}
}

func (ch *Character) BuildLockerS2C() []*dto.CashShopItem {
	items := make([]*dto.CashShopItem, 0, len(ch.locker))
	for _, item := range ch.locker {
		items = append(items, ch.BuildLockerItemS2C(item))
	}
	return items
}

func (ch *Character) BuildGiftsS2C() []*dto.CashShopGift {
	gifts := make([]*dto.CashShopGift, 0, len(ch.gifts))
	for _, gift := range ch.gifts {
		gifts = append(gifts, &dto.CashShopGift{
			Serial:     gift.Serial(),
			ItemID:     gift.ItemID(),
			SenderName: gift.SenderName(),
			Message:    gift.Message(),
		})
	}
	return gifts
}

func (ch *Character) BuildS2C(world ItemWorld) *dto.Character {
	character := &dto.Character{
		ID:                ch.id,
		Name:              ch.name,
		Gender:            ch.gender,
		SkinColor:         ch.skin,
		Face:              ch.face,
		Hair:              ch.hair,
		Pet:               ch.summonedPet,
		Level:             ch.level,
		Class:             ch.class,
		Str:               ch.str,
		Dex:               ch.dex,
		Int:               ch.intStat,
		Luk:               ch.luk,
		Hp:                ch.hp,
		MaxHp:             ch.maxHp,
		Mp:                ch.mp,
		MaxMp:             ch.maxMp,
		AbilityPoint:      ch.ap,
		SkillPoint:        ch.sp,
		Exp:               ch.exp,
		Population:        ch.population,
		Map:               ch.mapID,
		SpawnPoint:        ch.spawnPoint,
		BuddyCapacity:     ch.buddyCapacity,
		TeleportStones:    constant.NewTeleportStoneSlots(ch.teleportStones, constant.TeleportStoneCount),
		VipTeleportStones: constant.NewTeleportStoneSlots(ch.vipStones, constant.VipTeleportStoneCount),
		Inventory: &dto.Inventory{
			Tabs:     make(map[constant.InventoryType]*dto.InventoryTab),
			Equipped: make(map[constant.EquipmentPartsType]*dto.Equipment),
			Meso:     ch.meso,
		},
	}

	for _, invType := range []constant.InventoryType{constant.InventoryTypeEquipment, constant.InventoryTypeConsume, constant.InventoryTypeInstallation, constant.InventoryTypeETC, constant.InventoryTypeCash} {
		character.Inventory.Tabs[invType] = &dto.InventoryTab{Type: invType, SlotLimit: ch.SlotLimit(invType), Items: make(map[int16]dto.Item)}
	}
	for _, record := range ch.inventory {
		shown, err := item.BuildS2C(record.BuildRPC(), world.GetResources(), world.NewUniqueID)
		if err != nil {
			continue
		}
		if record.Slot() < 0 {
			if equipment, ok := shown.(*dto.Equipment); ok {
				character.Inventory.Equipped[constant.EquipmentPartsType(record.Slot())] = equipment
			}
			continue
		}
		if tab := character.Inventory.Tabs[ch.inventoryType(record)]; tab != nil {
			tab.Items[int16(record.Slot())] = shown
		}
	}

	for _, skill := range ch.skills {
		character.Skills = append(character.Skills, &dto.Skill{
			ID:          skill.id,
			SkillLevel:  uint32(skill.level),
			MasterLevel: uint32(skill.masterLevel),
		})
	}

	resources := world.GetResources()
	for _, quest := range ch.quests {
		if quest.id > 0xFFFF {
			continue
		}
		status := &dto.QuestStatus{
			QuestID:        uint16(quest.id),
			Status:         uint8(quest.status),
			StatusRecord:   quest.statusRecord,
			CompletionTime: time.UnixMilli(quest.completion),
		}
		if quest.deadline > 0 {
			status.Deadline = time.UnixMilli(quest.deadline)
		}
		if wzQuest := resources.Quests[quest.id]; wzQuest != nil {
			for _, mobID := range wzQuest.MobIDs {
				status.MobKills = append(status.MobKills, uint16(quest.mobKills[mobID]))
			}
		}
		switch wz.QuestStatus(quest.status) {
		case wz.QuestStatusStarted:
			character.QuestsStarted = append(character.QuestsStarted, status)
		case wz.QuestStatusCompleted:
			character.QuestsCompleted = append(character.QuestsCompleted, status)
		}
	}
	return character
}
