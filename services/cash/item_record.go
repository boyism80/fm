package cash

import (
	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
)

type equipBonus struct {
	str, dex, intStat, luk int32
	maxHP, maxMP           int32
	pad, mad, pdd, mdd     int32
	acc, avoid, hands      int32
	speed, jump            int32
}

func equipBonusFromProto(pb *internal.EquipmentBonusStats) *equipBonus {
	if pb == nil {
		return nil
	}
	return &equipBonus{
		str: pb.GetStr(), dex: pb.GetDex(), intStat: pb.GetIntStat(), luk: pb.GetLuk(),
		maxHP: pb.GetMaxHp(), maxMP: pb.GetMaxMp(),
		pad: pb.GetPad(), mad: pb.GetMad(), pdd: pb.GetPdd(), mdd: pb.GetMdd(),
		acc: pb.GetAcc(), avoid: pb.GetAvoid(), hands: pb.GetHands(),
		speed: pb.GetSpeed(), jump: pb.GetJump(),
	}
}

func (b *equipBonus) BuildRPC() *internal.EquipmentBonusStats {
	if b == nil {
		return nil
	}
	return &internal.EquipmentBonusStats{
		Str: b.str, Dex: b.dex, IntStat: b.intStat, Luk: b.luk,
		MaxHp: b.maxHP, MaxMp: b.maxMP,
		Pad: b.pad, Mad: b.mad, Pdd: b.pdd, Mdd: b.mdd,
		Acc: b.acc, Avoid: b.avoid, Hands: b.hands,
		Speed: b.speed, Jump: b.jump,
	}
}

type petRecord struct {
	name                              string
	level, closeness, fullness, speed uint32
	skills, secondsLeft               uint32
	exceptions                        []uint32
}

func petFromProto(pb *internal.Pet) *petRecord {
	if pb == nil {
		return nil
	}
	return &petRecord{
		name: pb.GetName(), level: pb.GetLevel(), closeness: pb.GetCloseness(), fullness: pb.GetFullness(),
		speed: pb.GetSpeed(), skills: pb.GetSkills(), secondsLeft: pb.GetSecondsLeft(),
		exceptions: append([]uint32(nil), pb.GetExceptions()...),
	}
}

func (p *petRecord) BuildRPC() *internal.Pet {
	if p == nil {
		return nil
	}
	return &internal.Pet{
		Name: p.name, Level: p.level, Closeness: p.closeness, Fullness: p.fullness,
		Speed: p.speed, Skills: p.skills, SecondsLeft: p.secondsLeft,
		Exceptions: append([]uint32(nil), p.exceptions...),
	}
}

// ItemRecord is one inventory item held by the cash shop character.
type ItemRecord struct {
	uniqueID      *uint64
	ownerID       uint32
	itemID        uint32
	slot          int32
	count         uint32
	expiration    int64
	enhanceChance uint32
	flag          uint32
	skillBonus    uint32
	ownerName     string
	inventoryType uint32
	bonus         *equipBonus
	enhanceCount  uint32
	pet           *petRecord
	marriageID    uint32
}

func ItemFromProto(pb *internal.Inventory) ItemRecord {
	if pb == nil {
		return ItemRecord{}
	}
	item := ItemRecord{
		ownerID: pb.GetOwnerId(), itemID: pb.GetItemId(), slot: pb.GetSlot(), count: pb.GetCount(),
		expiration: pb.GetExpirationUnixMs(), enhanceChance: pb.GetEnhanceChance(), flag: pb.GetFlag(),
		skillBonus: pb.GetSkillBonus(), ownerName: pb.GetOwnerName(), inventoryType: pb.GetInventoryType(),
		bonus: equipBonusFromProto(pb.GetEquipBonusStats()), enhanceCount: pb.GetEnhanceCount(),
		pet: petFromProto(pb.GetPet()), marriageID: pb.GetMarriageId(),
	}
	if pb.UniqueId != nil {
		id := pb.GetUniqueId()
		item.uniqueID = &id
	}
	return item
}

func (it ItemRecord) clone() ItemRecord {
	out := it
	if it.uniqueID != nil {
		id := *it.uniqueID
		out.uniqueID = &id
	}
	if it.bonus != nil {
		bonus := *it.bonus
		out.bonus = &bonus
	}
	if it.pet != nil {
		pet := *it.pet
		pet.exceptions = append([]uint32(nil), it.pet.exceptions...)
		out.pet = &pet
	}
	return out
}

func (it ItemRecord) Serial() uint64 {
	if it.uniqueID == nil {
		return 0
	}
	return *it.uniqueID
}

func (it ItemRecord) ItemID() uint32 { return it.itemID }

func (it ItemRecord) Slot() int32 { return it.slot }

func (it ItemRecord) Count() uint32 { return it.count }

func (it ItemRecord) Expiration() int64 { return it.expiration }

func (it ItemRecord) InventoryType() uint32 { return it.inventoryType }

func (it *ItemRecord) SetUniqueID(id uint64) { it.uniqueID = &id }

func (it *ItemRecord) SetOwnerID(id uint32) { it.ownerID = id }

func (it *ItemRecord) SetSlot(slot int32) { it.slot = slot }

func (it *ItemRecord) SetCount(count uint32) { it.count = count }

func (it *ItemRecord) SetFlag(flag uint32) { it.flag = flag }

func (it *ItemRecord) SetOwnerName(name string) { it.ownerName = name }

func (it *ItemRecord) SetInventoryType(invType uint32) { it.inventoryType = invType }

func (it *ItemRecord) SetExpiration(ms int64) { it.expiration = ms }

func (it *ItemRecord) ensurePet(name string, fullness, skills, secondsLeft uint32) {
	if it.pet != nil {
		return
	}
	it.pet = &petRecord{name: name, level: 1, fullness: fullness, skills: skills, secondsLeft: secondsLeft}
}

func (it ItemRecord) BuildRPC() *internal.Inventory {
	pb := &internal.Inventory{
		OwnerId: it.ownerID, ItemId: it.itemID, Slot: it.slot, Count: it.count,
		ExpirationUnixMs: it.expiration, EnhanceChance: it.enhanceChance, Flag: it.flag,
		SkillBonus: it.skillBonus, OwnerName: it.ownerName, InventoryType: it.inventoryType,
		EquipBonusStats: it.bonus.BuildRPC(), EnhanceCount: it.enhanceCount,
		Pet: it.pet.BuildRPC(), MarriageId: it.marriageID,
	}
	if it.uniqueID != nil {
		id := *it.uniqueID
		pb.UniqueId = &id
	}
	return pb
}

// LockerItem is one cash locker entry.
type LockerItem struct {
	commoditySN uint32
	buyerName   string
	item        ItemRecord
}

func LockerItemFromProto(pb *internal.CashItem) *LockerItem {
	if pb == nil {
		return nil
	}
	return &LockerItem{commoditySN: pb.GetCommoditySn(), buyerName: pb.GetBuyerName(), item: ItemFromProto(pb.GetItem())}
}

func (it *LockerItem) Serial() uint64 {
	if it == nil {
		return 0
	}
	return it.item.Serial()
}

func (it *LockerItem) ItemID() uint32 {
	if it == nil {
		return 0
	}
	return it.item.ItemID()
}

func (it *LockerItem) CommoditySN() uint32 {
	if it == nil {
		return 0
	}
	return it.commoditySN
}

func (it *LockerItem) BuyerName() string {
	if it == nil {
		return ""
	}
	return it.buyerName
}

func (it *LockerItem) Count() uint32 {
	if it == nil {
		return 0
	}
	return it.item.Count()
}

func (it *LockerItem) Expiration() int64 {
	if it == nil {
		return 0
	}
	return it.item.Expiration()
}

func (it *LockerItem) Item() ItemRecord {
	if it == nil {
		return ItemRecord{}
	}
	return it.item.clone()
}

func (it *LockerItem) BuildRPC() *internal.CashItem {
	if it == nil {
		return nil
	}
	return &internal.CashItem{CommoditySn: it.commoditySN, BuyerName: it.buyerName, Item: it.item.BuildRPC()}
}

type giftRecord struct {
	serial     uint64
	itemID     uint32
	senderName string
	message    string
}

func giftFromProto(pb *internal.CashGift) giftRecord {
	return giftRecord{serial: pb.GetSerial(), itemID: pb.GetItemId(), senderName: pb.GetSenderName(), message: pb.GetMessage()}
}

func (g giftRecord) Serial() uint64     { return g.serial }
func (g giftRecord) ItemID() uint32     { return g.itemID }
func (g giftRecord) SenderName() string { return g.senderName }
func (g giftRecord) Message() string    { return g.message }

type skillRecord struct {
	id          uint32
	level       int32
	masterLevel int32
}

type questRecord struct {
	id           uint32
	status       uint32
	mobKills     map[uint32]uint32
	statusRecord string
	completion   int64
	deadline     int64
}

type debuffRecord struct {
	mask       uint32
	position   int32
	x          int32
	skillID    uint32
	skillLevel uint32
	end        int64
}

func (d debuffRecord) BuildRPC() *internal.Debuff {
	return &internal.Debuff{Mask: d.mask, Position: d.position, X: d.x, SkillId: d.skillID, SkillLevel: d.skillLevel, EndUnixMs: d.end}
}
