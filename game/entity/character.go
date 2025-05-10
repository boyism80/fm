package entity

import (
	"github.com/boyism80/fm/common/context"
	"github.com/boyism80/fm/common/stream"
	"github.com/boyism80/fm/common/types"
)

type Character struct {
	Life
	Id            uint32
	Name          string
	Gender        uint8
	SkinColor     uint8
	Face          uint32
	Hair          uint32
	Level         uint8
	Class         uint16
	Rank          uint32
	RankDiff      int32
	ClassRank     uint32
	ClassRankDiff int32
	Admin         bool
	Str           uint16
	Dex           uint16
	Int           uint16
	Luk           uint16
	Hp            uint16
	MaxHp         uint16
	Mp            uint16
	MaxMp         uint16
	AbilityPoint  uint16
	SkillPoint    []uint16
	Exp           uint32
	FamePoint     uint16
	Map           uint32
	SpawnPoint    uint8
	Mega          bool
	Meso          uint32

	Random1          stream.RandomStream
	Random2          stream.RandomStream
	Random3          stream.RandomStream
	Inventory        map[InventoryType]*Inventory
	Equipments       map[EquipmentPartsType]*Equipment
	SkillsMap        map[*Skill]*SkillEntry
	CoolDowns        map[uint32]*CooldownEntry
	Quests           map[int]*QuestStatus
	MarriageId       uint32
	RegRocks         []uint32 // 기본 순간이동 장소 (5개)
	Rocks            []uint32 // VIP 순간이동 장소 (10개)
	MonsterBookCover uint32
	MonsterBook      *MonsterBook
	QuestInfo        map[uint16]string
}

func (ch *Character) IsRanked() bool {
	if ch.Admin {
		return false
	}

	if ch.Level < 30 {
		return false
	}

	return true
}

func (ch *Character) IsEvan() bool {
	return ch.Class == 2001 || (ch.Class >= 2200 && ch.Class <= 2218)
}

func (ch *Character) IsKOC() bool {
	return ch.Class >= 1000 && ch.Class < 2000
}

func (ch *Character) IsMercedes() bool {
	return ch.Class == 2002 || (ch.Class >= 2300 && ch.Class <= 2312)
}

func (ch *Character) IsDemon() bool {
	return ch.Class == 3001 || (ch.Class >= 3100 && ch.Class <= 3112)
}

func (ch *Character) IsAran() bool {
	return (ch.Class >= 2000 && ch.Class <= 2112) && ch.Class != 2001 && ch.Class != 2002
}

func (ch *Character) IsResist() bool {
	return ch.Class >= 3000 && ch.Class <= 3512
}

func (ch *Character) IsAdventurer() bool {
	return ch.Class < 1000
}

func (ch *Character) IsCannon() bool {
	return ch.Class == 1 || ch.Class == 501 || (ch.Class >= 530 && ch.Class <= 532)
}

func (ch *Character) RemainingSkillPoints() uint16 {
	ret := 0
	for _, sp := range ch.SkillPoint {
		if sp > 0 {
			ret++
		}
	}
	return uint16(ret)
}

func NewDummyCharacter(id uint32, name string, ctx *context.ServerContext) Character {
	ch := Character{
		Id:         id,
		Name:       name,
		Gender:     0,
		SkinColor:  0,
		Face:       20100,
		Hair:       30000,
		Level:      1,
		Class:      0,
		Str:        12,
		Dex:        5,
		Int:        4,
		Luk:        4,
		Hp:         50,
		MaxHp:      50,
		Mp:         5,
		MaxMp:      5,
		SpawnPoint: 1,
		Map:        200000301,
		Meso:       2135983647,

		Random1: stream.NewRandomStream(),
		Random2: stream.NewRandomStream(),
		Random3: stream.NewRandomStream(),

		Inventory: map[InventoryType]*Inventory{
			InventoryTypeEquip: NewMapleInventory(InventoryTypeEquip),
			InventoryTypeUse:   NewMapleInventory(InventoryTypeUse),
			InventoryTypeSetUp: NewMapleInventory(InventoryTypeSetUp),
			InventoryTypeEtc:   NewMapleInventory(InventoryTypeEtc),
			InventoryTypeCash:  NewMapleInventory(InventoryTypeCash),
		},
		Equipments: map[EquipmentPartsType]*Equipment{
			EquipmentPartsWeapon: nil,
			EquipmentPartsShield: nil,
		},

		RegRocks: []uint32{999999999, 999999999, 999999999, 999999999, 999999999},
		Rocks:    []uint32{999999999, 999999999, 999999999, 999999999, 999999999, 999999999, 999999999, 999999999, 999999999, 999999999},
	}

	if ctx != nil {
		ch.Equipments[EquipmentPartsWeapon] = &Equipment{
			BaseItem: &BaseItem{
				Object:     nil,
				Template:   ctx.GameData.Items[1302000],
				UniqueId:   0,
				Expiration: -1,
			},
			EnchantChance: 7,
		}

		// TODO: 펫 추가하고 이어서 작업
		// ch.Inventory[InventoryTypeCash].Items[1] = &Pet{
		// 	BaseItem: &BaseItem{
		// 		Object:     nil,
		// 		Template:   ctx.GameData.Items[5000007],
		// 		UniqueId:   1,
		// 		Expiration: -1,
		// 	},
		// 	Level:       1,
		// 	Closeness:   0,
		// 	Fullness:    0,
		// 	Speed:       1,
		// 	Flags:       0,
		// 	SecondsLeft: 0,
		// 	Expiration:  -1,
		// }
	}

	return ch
}

func (ch *Character) Send(p types.Packet, policy types.SendPolicy) {

}
