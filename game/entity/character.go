package entity

import (
	"fmt"
	"time"

	"github.com/boyism80/fm/common/context"
	"github.com/boyism80/fm/common/stream"
	"github.com/boyism80/fm/common/types"
	"github.com/boyism80/fm/common/util"
	"github.com/boyism80/fm/game/constant"
	lua "github.com/yuin/gopher-lua"
)

type Character struct {
	Life
	Sendable
	Listener      CharacterListener
	Dialog        *lua.LState
	ID            uint32
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
	AbilityPoint  uint16
	SkillPoint    []uint16
	Exp           uint32
	FamePoint     uint16
	Map           uint32
	SpawnPoint    uint8
	Mega          bool
	Meso          int32

	Random1          stream.RandomStream
	Random2          stream.RandomStream
	Random3          stream.RandomStream
	Inventory        map[constant.InventoryType]*Inventory
	Equipments       map[constant.EquipmentPartsType]*Equipment
	Rings            RingContainer
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

type CooldownEntry struct {
	SkillId   uint32
	StartTime int64 // milliseconds
	Length    int64 // milliseconds
}

type Ring struct {
	RingId       uint64
	PartnerId    uint64
	RingUniqueId uint64
	PartnerChrId uint32
	ItemId       uint32
	PartnerName  string
	Equipped     bool
}

type RingContainer struct {
	Left  []*Ring
	Mid   []*Ring
	Right []*Ring
}

type MonsterBook struct {
	Cards map[uint32]uint32
}

func (ch *Character) Send(p types.Packet, policy types.SendPolicy) {
	if ch.Sendable == nil {
		return
	}
	ch.Sendable.Send(p, policy)
}

func (mb *MonsterBook) Serialize(writer *stream.StreamWriter) {
	writer.WriteU16(uint16(len(mb.Cards)))

	for cardId := range mb.Cards {
		writer.WriteU16(uint16(cardId))
	}
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

func NewDummyCharacter(sender Sendable, listener CharacterListener, id uint32, name string, ctx *context.ServerContext) Character {
	ch := Character{
		Sendable: sender,
		Listener: listener,
		Life: Life{
			Hp:    50,
			MaxHp: 50,
			Mp:    5,
			MaxMp: 5,
		},
		ID:         id,
		Name:       name,
		Gender:     0,
		SkinColor:  0,
		Face:       20100,
		Hair:       30000,
		Level:      255,
		Class:      0,
		Str:        12,
		Dex:        5,
		Int:        4,
		Luk:        4,
		SpawnPoint: 1,
		Map:        200000301,
		Meso:       2135983647,

		Random1: stream.NewRandomStream(),
		Random2: stream.NewRandomStream(),
		Random3: stream.NewRandomStream(),

		Inventory: map[constant.InventoryType]*Inventory{
			constant.InventoryTypeEquipment:    NewInventory(constant.InventoryTypeEquipment),
			constant.InventoryTypeConsume:      NewInventory(constant.InventoryTypeConsume),
			constant.InventoryTypeInstallation: NewInventory(constant.InventoryTypeInstallation),
			constant.InventoryTypeEtc:          NewInventory(constant.InventoryTypeEtc),
			constant.InventoryTypeCash:         NewInventory(constant.InventoryTypeCash),
		},
		Rings: RingContainer{
			Left:  []*Ring{},
			Right: []*Ring{},
			Mid:   []*Ring{},
		},
		Equipments: map[constant.EquipmentPartsType]*Equipment{
			constant.EquipmentPartsWeapon: nil,
			constant.EquipmentPartsShield: nil,
		},

		RegRocks: []uint32{999999999, 999999999, 999999999, 999999999, 999999999},
		Rocks:    []uint32{999999999, 999999999, 999999999, 999999999, 999999999, 999999999, 999999999, 999999999, 999999999, 999999999},
	}

	if ctx != nil {
		ch.Equipments[constant.EquipmentPartsWeapon] = &Equipment{
			ItemCore: &ItemCore{
				Spec:       ctx.Resources.Items[1302000],
				Count:      1,
				UniqueId:   0,
				Expiration: util.TimeMax,
			},
			EnchantChance: 7,
		}

		petExpiration, err := time.ParseInLocation("2006-01-02 15:04:05", "2025-05-30 09:30:00", util.KST)
		if err != nil {
			fmt.Println(err)
		}
		ch.Inventory[constant.InventoryTypeCash].Items[1] = &Pet{
			ItemCore: &ItemCore{
				Spec:       ctx.Resources.Items[5000007],
				Count:      1,
				UniqueId:   1,
				Expiration: util.TimeMax,
			},
			Level:       1,
			Closeness:   0,
			Fullness:    0,
			Speed:       1,
			Flags:       0,
			SecondsLeft: 0,
			Expiration:  petExpiration,
		}

		ch.Inventory[constant.InventoryTypeEtc].Items[1] = &GeneralItem{
			ItemCore: &ItemCore{
				Spec:       ctx.Resources.Items[4000001],
				Count:      100,
				Expiration: util.TimeMax,
			},
		}

		ch.Inventory[constant.InventoryTypeEquipment].Items[1], err = NewItem(ctx, 1060002, 1)
		ch.Inventory[constant.InventoryTypeEquipment].Items[2], err = NewItem(ctx, 1060006, 1)
		ch.Inventory[constant.InventoryTypeEquipment].Items[3], err = NewItem(ctx, 1040002, 1)
		ch.Inventory[constant.InventoryTypeEquipment].Items[4], err = NewItem(ctx, 1040010, 1)
	}

	return ch
}

func (m *Character) LuaTypeName() string {
	return "Character"
}
func (m *Character) LuaBuiltinFuncs() map[string]lua.LGFunction {
	return characterBuiltinFuncs
}
