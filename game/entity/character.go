package entity

import (
	"fmt"
	"sync"
	"time"

	"github.com/boyism80/fm/core/luax"
	"github.com/boyism80/fm/game/constant"
	"github.com/boyism80/fm/game/wz"
	"github.com/boyism80/fm/protocol/response"
	"github.com/boyism80/fm/stream"
	"github.com/boyism80/fm/types"
	"github.com/boyism80/fm/util"
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
	SkillPoint    uint16
	HpApUsed      uint16
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
	SkillsMap        map[uint32]*SkillEntry
	CoolDowns        map[uint32]*CooldownEntry
	Quests           map[int]*QuestStatus
	MarriageId       uint32
	RegRocks         []uint32 // Basic warp rock slots (5 slots)?
	Rocks            []uint32 // VIP warp rock slots (10 slots)?
	MonsterBookCover uint32
	MonsterBook      *MonsterBook
	QuestInfo        map[uint16]string

	// Dialog state management
	currentDialog *lua.LState // Current dialog coroutine
	dialogMutex   sync.Mutex  // Mutex for dialog state access

	// Shop state management
	CurrentShopID uint32 // Current shop NPC ID (0 if no shop is open)

	// Chair state
	Chair uint32 // Chair item ID (0 if not sitting)

	// Base stats (permanent) - only for STR, DEX, INT, LUK
	BaseStats BaseStats

	// Bonus stats (temporary from buffs/equipment) - only for STR, DEX, INT, LUK
	BonusStats BonusStats
}

// BaseStats represents permanent character statistics
type BaseStats struct {
	Str uint16 // Base strength
	Dex uint16 // Base dexterity
	Int uint16 // Base intelligence
	Luk uint16 // Base luck
}

// BonusStats represents temporary stat increases from buffs/equipment
type BonusStats struct {
	Str   int16 // Bonus strength (can be negative)
	Dex   int16 // Bonus dexterity (can be negative)
	Int   int16 // Bonus intelligence (can be negative)
	Luk   int16 // Bonus luck (can be negative)
	Watk  int16 // Bonus weapon attack
	Matk  int16 // Bonus magic attack
	Wdef  int16 // Bonus weapon defense
	Mdef  int16 // Bonus magic defense
	Acc   int16 // Bonus accuracy
	Avoid int16 // Bonus evasion
	Speed int16 // Bonus speed
	Jump  int16 // Bonus jump
}

type CooldownEntry struct {
	SkillId   uint32
	StartTime time.Time
	Duration  time.Duration
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

// GetTotalStr returns the total strength (BaseStr + BonusStr)
func (ch *Character) GetTotalStr() uint16 {
	// Migration: if BaseStats.Str is 0 and Str is set, copy Str to BaseStats.Str
	if ch.BaseStats.Str == 0 && ch.Str > 0 {
		ch.BaseStats.Str = ch.Str
		ch.Str = 0
	}

	total := int32(ch.BaseStats.Str) + int32(ch.BonusStats.Str)
	if total < 0 {
		return 0
	}
	if total > int32(constant.STAT_MAX_STR_DEX_INT_LUK) {
		return constant.STAT_MAX_STR_DEX_INT_LUK
	}
	return uint16(total)
}

// GetTotalDex returns the total dexterity (BaseDex + BonusDex)
func (ch *Character) GetTotalDex() uint16 {
	// Migration: if BaseStats.Dex is 0 and Dex is set, copy Dex to BaseStats.Dex
	if ch.BaseStats.Dex == 0 && ch.Dex > 0 {
		ch.BaseStats.Dex = ch.Dex
		ch.Dex = 0
	}

	total := int32(ch.BaseStats.Dex) + int32(ch.BonusStats.Dex)
	if total < 0 {
		return 0
	}
	if total > int32(constant.STAT_MAX_STR_DEX_INT_LUK) {
		return constant.STAT_MAX_STR_DEX_INT_LUK
	}
	return uint16(total)
}

// GetTotalInt returns the total intelligence (BaseInt + BonusInt)
func (ch *Character) GetTotalInt() uint16 {
	// Migration: if BaseStats.Int is 0 and Int is set, copy Int to BaseStats.Int
	if ch.BaseStats.Int == 0 && ch.Int > 0 {
		ch.BaseStats.Int = ch.Int
		ch.Int = 0
	}

	total := int32(ch.BaseStats.Int) + int32(ch.BonusStats.Int)
	if total < 0 {
		return 0
	}
	if total > int32(constant.STAT_MAX_STR_DEX_INT_LUK) {
		return constant.STAT_MAX_STR_DEX_INT_LUK
	}
	return uint16(total)
}

// GetTotalLuk returns the total luck (BaseLuk + BonusLuk)
func (ch *Character) GetTotalLuk() uint16 {
	// Migration: if BaseStats.Luk is 0 and Luk is set, copy Luk to BaseStats.Luk
	if ch.BaseStats.Luk == 0 && ch.Luk > 0 {
		ch.BaseStats.Luk = ch.Luk
		ch.Luk = 0
	}

	total := int32(ch.BaseStats.Luk) + int32(ch.BonusStats.Luk)
	if total < 0 {
		return 0
	}
	if total > int32(constant.STAT_MAX_STR_DEX_INT_LUK) {
		return constant.STAT_MAX_STR_DEX_INT_LUK
	}
	return uint16(total)
}

// notifyStatChange notifies the listener of a stat change
func (ch *Character) notifyStatChange(stat constant.Stat) {
	if ch.Listener == nil {
		return
	}

	stats := make(map[constant.Stat]int32)
	switch stat {
	case constant.STAT_STR:
		stats[constant.STAT_STR] = int32(ch.GetTotalStr())
	case constant.STAT_DEX:
		stats[constant.STAT_DEX] = int32(ch.GetTotalDex())
	case constant.STAT_INT:
		stats[constant.STAT_INT] = int32(ch.GetTotalInt())
	case constant.STAT_LUK:
		stats[constant.STAT_LUK] = int32(ch.GetTotalLuk())
	case constant.STAT_MAX_HP:
		stats[constant.STAT_MAX_HP] = int32(ch.Life.GetMaxHp())
	case constant.STAT_MAX_MP:
		stats[constant.STAT_MAX_MP] = int32(ch.Life.GetMaxMp())
	}

	ch.Listener.OnUpdateStats(stats, false)
}

func (ch *Character) Send(p types.Packet, policy types.SendPolicy) error {
	if ch.Sendable == nil {
		return nil
	}
	return ch.Sendable.Send(p, policy)
}

// GetMap returns the character's current map ID
func (ch *Character) GetMap() uint32 {
	return ch.Map
}

func (ch *Character) GetID() uint32 {
	return ch.ID
}

func (ch *Character) Message(message string) {
	if ch.Listener != nil {
		ch.Listener.OnMessage(constant.MSG_LIGHT_BLUE_TEXT, message)
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

func (ch *Character) IsAdventurer() bool {
	return ch.Class < 1000
}

// IsBeginner returns true if the character's class is a beginner class
func (ch *Character) IsBeginner() bool {
	return ch.Class == 0
}

func (ch *Character) getJobAdvancementLevel() int {
	class := ch.Class

	if class == 0 {
		return 0
	}

	if class >= 100 {
		secondDigit := (class / 10) % 10
		thirdDigit := class % 10

		if secondDigit == 0 {
			return 1
		} else if thirdDigit == 0 {
			return 2
		} else if thirdDigit == 1 {
			return 3
		} else {
			return 4
		}
	}

	return 0
}

func (ch *Character) IsCannon() bool {
	return ch.Class == 1 || ch.Class == 501 || (ch.Class >= 530 && ch.Class <= 532)
}

func (ch *Character) IsMagician() bool {
	return ch.Class >= 200 && ch.Class < 300
}

func (ch *Character) getSkillBookIndexByLevel(level uint8) int {
	if ch.IsBeginner() {
		return -1
	}

	isMagician := ch.IsMagician()
	minLevel := uint8(10)
	if isMagician {
		minLevel = 8
	}

	if level < minLevel {
		return -1
	}

	if level <= 30 {
		return 0
	} else if level <= 70 {
		return 1
	} else if level <= 120 {
		return 2
	} else {
		return 3
	}
}

func (ch *Character) GetSkillBookIndex() int {
	class := ch.Class

	if class >= 100 {
		secondDigit := (class / 10) % 10
		thirdDigit := class % 10

		if secondDigit == 0 {
			return 0
		} else if thirdDigit == 0 {
			return 1
		} else if thirdDigit == 1 {
			return 2
		} else if thirdDigit == 2 {
			return 3
		}
	}

	return 0
}

func (ch *Character) GetSkillBookIndexForSkill(skillID uint32) int {
	classID := skillID / 10000

	if classID >= 100 {
		secondDigit := (classID / 10) % 10
		thirdDigit := classID % 10

		if secondDigit == 0 {
			return 0
		} else if thirdDigit == 0 {
			return 1
		} else if thirdDigit == 1 {
			return 2
		} else if thirdDigit == 2 {
			return 3
		}
	}

	return 0
}

func (ch *Character) RemainingSkillPoints() uint16 {
	return ch.SkillPoint
}

func (ch *Character) GetTotalSkillLevel(skillID uint32) int {
	if ch.SkillsMap == nil {
		return 0
	}
	skillEntry, exists := ch.SkillsMap[skillID]
	if !exists || skillEntry == nil {
		return 0
	}
	return skillEntry.SkillLevel
}

func (ch *Character) IsSkillCooling(skillID uint32) bool {
	if ch.CoolDowns == nil {
		return false
	}
	cooldown, exists := ch.CoolDowns[skillID]
	if !exists || cooldown == nil {
		return false
	}
	return time.Since(cooldown.StartTime) < cooldown.Duration
}

func (ch *Character) AddCooldown(skillID uint32, cooldownSeconds int) {
	if ch.CoolDowns == nil {
		ch.CoolDowns = make(map[uint32]*CooldownEntry)
	}
	ch.CoolDowns[skillID] = &CooldownEntry{
		SkillId:   skillID,
		StartTime: time.Now(),
		Duration:  time.Duration(cooldownSeconds) * time.Second,
	}
}

func (ch *Character) ConsumeMP(amount uint16) bool {
	if ch.Mp < amount {
		return false
	}
	ch.Mp -= amount
	if ch.Listener != nil {
		ch.Listener.OnUpdateStats(map[constant.Stat]int32{
			constant.STAT_MP: int32(ch.Mp),
		}, false)
	}
	return true
}

func (ch *Character) ChangeClass(newClass uint16) {
	oldClass := ch.Class
	ch.Class = newClass

	ch.grantClassChangeSP(newClass)
	ch.initializeBaseSkills(newClass)

	if ch.Listener != nil {
		ch.Listener.OnClassChange(oldClass, newClass)
	}
}

func (ch *Character) grantClassChangeSP(newClass uint16) {
	if ch.IsBeginner() {
		return
	}

	ch.SkillPoint++

	if newClass >= 100 {
		thirdDigit := newClass % 10
		if thirdDigit >= 2 {
			ch.SkillPoint += 2
		}
	}

	if newClass%100 == 0 {
		minLevel := uint8(10)
		if newClass == 200 {
			minLevel = 8
		}

		if ch.Level > minLevel {
			spToGrant := uint16(3 * (int(ch.Level) - int(minLevel)))
			ch.SkillPoint += spToGrant
		}
	}
}

func (ch *Character) initializeBaseSkills(newClass uint16) {
	if ch.Context == nil {
		return
	}

	resources := ch.Context.GetResources()
	if resources == nil {
		return
	}

	advancementLevel := ch.getJobAdvancementLevel()
	if advancementLevel < 3 {
		return
	}

	classID := uint32(newClass)
	skillIDStart := classID * 10000
	skillIDEnd := skillIDStart + 9999

	if ch.SkillsMap == nil {
		ch.SkillsMap = make(map[uint32]*SkillEntry)
	}

	for skillID := skillIDStart; skillID <= skillIDEnd; skillID++ {
		wzSkill := resources.GetSkill(skillID)
		if wzSkill == nil {
			continue
		}

		if wzSkill.Invisible {
			continue
		}

		if !ch.isFourthClassSkill(skillID, wzSkill) {
			continue
		}

		masterLevel := 0
		if wzSkill.MasterLevel > 0 {
			masterLevel = wzSkill.MasterLevel
		} else if wzSkill.MaxLevel > 0 {
			masterLevel = wzSkill.MaxLevel
		} else {
			continue
		}

		existingEntry, exists := ch.SkillsMap[skillID]
		if exists && existingEntry != nil {
			if existingEntry.SkillLevel > 0 || existingEntry.MasterLevel > 0 {
				continue
			}
		}

		skillEntry := &SkillEntry{
			Skill:       wzSkill,
			SkillLevel:  0,
			MasterLevel: masterLevel,
			Expiration:  time.Time{},
		}
		ch.SkillsMap[skillID] = skillEntry

		if ch.Listener != nil {
			ch.Send(&response.UpdateSkills{
				SkillID:     skillID,
				Level:       0,
				MasterLevel: int32(skillEntry.MasterLevel),
			}, types.SEND_POLICY_ENCRYPT)
		}
	}
}

func (ch *Character) isFourthClassSkill(skillID uint32, wzSkill *wz.Skill) bool {
	classID := skillID / 10000

	if classID == 2312 {
		return true
	}

	if (wzSkill.MaxLevel <= 15 && !wzSkill.Invisible && wzSkill.MasterLevel <= 0) ||
		skillID == 3220010 || skillID == 3120011 || skillID == 33120010 || skillID == 32120009 ||
		skillID == 5321006 || skillID == 21120011 || skillID == 22181004 || skillID == 4340010 {
		return false
	}

	if classID >= 2212 && classID < 3000 {
		return (classID % 10) >= 7
	}

	if classID >= 430 && classID <= 434 {
		return (classID%10) == 4 || wzSkill.MasterLevel > 0
	}

	return (classID%10) == 2 && skillID < 90000000
}

// AddExp adds experience points to the character and notifies the listener
func (ch *Character) AddExp(exp uint32) {
	// Apply experience rate multiplier if context is available
	if ch.Context != nil {
		expRate := ch.Context.GetExpRate()
		if expRate > 0 {
			exp = exp * uint32(expRate)
		}
	}

	ch.Exp += exp
	ch.Listener.OnExpGain(exp)

	// Check for level up
	if !ch.tryLevelUp() {
		ch.Listener.OnUpdateStats(map[constant.Stat]int32{
			constant.STAT_EXP: int32(ch.Exp),
		}, false)
	}
}

// AddMeso adds meso to character
func (ch *Character) AddMeso(amount int32) {
	if amount <= 0 {
		return
	}

	// Check for overflow
	if ch.Meso > 0 && amount > 0 && ch.Meso+amount < ch.Meso {
		ch.Meso = int32(^uint32(0) >> 1) // int32.MaxValue
	} else {
		ch.Meso += amount
	}

	// Notify listener about meso change
	if ch.Listener != nil {
		ch.Listener.OnMesoChanged(ch.Meso)
	}
}

// RemoveMeso removes meso from character
func (ch *Character) RemoveMeso(amount int32) {
	if amount <= 0 {
		return
	}

	if ch.Meso < amount {
		ch.Meso = 0
	} else {
		ch.Meso -= amount
	}

	// Notify listener about meso change
	if ch.Listener != nil {
		ch.Listener.OnMesoChanged(ch.Meso)
	}
}

// AddItem adds an item to character's inventory
// If allOrNothing is true, adds all items or returns an error if cannot add all.
// If allOrNothing is false, adds as many items as possible and returns the count added.
// Returns (addedCount, error). If allOrNothing is true and not all items can be added, returns (0, error).
func (ch *Character) AddItem(item Item, allOrNothing bool) (uint16, error) {
	if item == nil {
		return 0, fmt.Errorf("item is nil")
	}

	invenType := item.GetInventoryType()
	inven := ch.Inventory[invenType]
	if inven == nil {
		return 0, fmt.Errorf("inventory type %d not found", invenType)
	}

	model := item.GetModel()
	requestedCount := item.GetCount()

	if allOrNothing {
		if !inven.IsFree(model, requestedCount) {
			return 0, fmt.Errorf("not enough inventory space for %d items", requestedCount)
		}
	}

	remainingCount := requestedCount
	addedCount := uint16(0)

	for remainingCount > 0 {
		slot, ok := inven.FindSlot(model)
		if !ok {
			// No more slots available
			if allOrNothing && addedCount == 0 {
				return 0, fmt.Errorf("no available slot found")
			}
			break
		}

		exists, ok := inven.Items[int16(slot)]
		cap := uint16(0)
		if ok {
			cap = min(model.GetCapacity()-exists.GetCount(), remainingCount)
			exists.Increase(cap)
			if ch.Listener != nil {
				ch.Listener.OnInventorySlotUpdated(invenType, int16(slot), exists)
			}
		} else {
			cap = min(model.GetCapacity(), remainingCount)
			inven.Items[int16(slot)] = item.Clone(cap)
			if ch.Listener != nil {
				ch.Listener.OnInventorySlotAdded(invenType, int16(slot), inven.Items[int16(slot)])
			}
		}
		remainingCount -= cap
		addedCount += cap
	}

	// Update original item's count to reflect remaining items
	if !allOrNothing && addedCount < requestedCount {
		item.SetCount(remainingCount)
	}

	if ch.Listener != nil && addedCount > 0 {
		ch.Listener.OnShowItemGain(model.GetID(), uint32(addedCount), constant.ShowItemGainTypeStatus)
	}

	return addedCount, nil
}

// GainMeso adds meso to character and shows gain notification
func (ch *Character) GainMeso(amount int32) {
	if amount <= 0 {
		return
	}

	// Check for overflow
	if ch.Meso > 0 && amount > 0 && ch.Meso+amount < ch.Meso {
		ch.Meso = int32(^uint32(0) >> 1) // int32.MaxValue
	} else {
		ch.Meso += amount
	}

	// Notify listener about meso change
	if ch.Listener != nil {
		ch.Listener.OnMesoChanged(ch.Meso)
		ch.Listener.OnShowMesoGain(amount, constant.ShowMesoGainTypeStatus)
		ch.Listener.OnUpdateStats(map[constant.Stat]int32{
			constant.STAT_MESO: ch.Meso,
		}, false)
	}
}

// Luable interface implementation
func (ch *Character) LuaTypeName() string {
	return "LuaCharacter"
}

func (ch *Character) LuaBuiltinFuncs() map[string]lua.LGFunction {
	return map[string]lua.LGFunction{
		"id": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			ch, ok := ud.Value.(*Character)
			if !ok {
				L.ArgError(1, "Character expected")
				return 0
			}

			argc := L.GetTop()
			if argc == 1 {
				// Getter: return id
				L.Push(lua.LNumber(ch.ID))
				return 1
			} else {
				L.ArgError(2, "id() is read-only")
				return 0
			}
		},
		"name": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			ch, ok := ud.Value.(*Character)
			if !ok {
				L.ArgError(1, "Character expected")
				return 0
			}

			argc := L.GetTop()
			if argc == 1 {
				// Getter: return name
				L.Push(lua.LString(ch.Name))
				return 1
			} else if argc == 2 {
				// Setter: name(value)
				name := L.CheckString(2)
				ch.Name = name
				return 0
			} else {
				L.ArgError(2, "name() requires 0 or 1 arguments")
				return 0
			}
		},
		"level": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			ch, ok := ud.Value.(*Character)
			if !ok {
				L.ArgError(1, "Character expected")
				return 0
			}

			argc := L.GetTop()
			if argc == 1 {
				// Getter: return level
				L.Push(lua.LNumber(ch.Level))
				return 1
			} else if argc == 2 {
				// Setter: level(value)
				level := L.CheckInt(2)
				if level < 1 {
					level = 1
				}
				if level > 200 {
					level = 200
				}
				ch.Level = uint8(level)
				return 0
			} else {
				L.ArgError(2, "level() requires 0 or 1 arguments")
				return 0
			}
		},
		"exp": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			ch, ok := ud.Value.(*Character)
			if !ok {
				L.ArgError(1, "Character expected")
				return 0
			}

			argc := L.GetTop()
			if argc == 1 {
				// Getter: return exp
				L.Push(lua.LNumber(ch.Exp))
				return 1
			} else if argc == 2 {
				// Setter: exp(value) or exp(+value) or exp(-value)
				value := L.CheckNumber(2)
				if value >= 0 {
					// Direct set or add operation
					if value < 0 {
						value = 0
					}
					ch.Exp = uint32(value)
				} else {
					// Add operation (negative value)
					amount := uint32(-value)
					ch.AddExp(amount)
					return 0
				}
				return 0
			} else {
				L.ArgError(2, "exp() requires 0 or 1 arguments")
				return 0
			}
		},
		"meso": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			ch, ok := ud.Value.(*Character)
			if !ok {
				L.ArgError(1, "Character expected")
				return 0
			}

			argc := L.GetTop()
			if argc == 1 {
				// Getter: return meso
				L.Push(lua.LNumber(ch.Meso))
				return 1
			} else if argc == 2 {
				// Setter: meso(value) or meso(+value) or meso(-value)
				value := L.CheckNumber(2)
				if value >= 0 {
					// Direct set or add operation
					if value <= 2147483647 { // int32.MaxValue
						ch.Meso = int32(value)
					} else {
						ch.Meso = 2147483647
					}
				} else {
					// Subtract operation (negative value)
					amount := int32(-value)
					if ch.Meso < amount {
						ch.Meso = 0
					} else {
						ch.Meso -= amount
					}
				}
				if ch.Listener != nil {
					ch.Listener.OnMesoChanged(ch.Meso)
				}
				return 0
			} else {
				L.ArgError(2, "meso() requires 0 or 1 arguments")
				return 0
			}
		},
		"chat": func(L *lua.LState) int {
			argc := L.GetTop()
			ud := L.CheckUserData(1)
			ch, ok := ud.Value.(*Character)
			if !ok {
				L.ArgError(1, "Character expected")
				return 0
			}

			message := L.CheckString(2)
			highlight := false
			if argc > 2 {
				highlight = L.CheckBool(3)
			}
			dontRecordHistory := false
			if argc > 3 {
				dontRecordHistory = L.CheckBool(4)
			}

			if ch.Listener != nil {
				ch.Listener.OnChat(message, highlight, dontRecordHistory)
			}
			return 0
		},
		"dialog": func(L *lua.LState) int {
			argc := L.GetTop()
			ud := L.CheckUserData(1)
			ch, ok := ud.Value.(*Character)
			if !ok {
				L.ArgError(1, "Character expected")
				return 0
			}

			npc := 0
			if argc > 1 {
				npc = L.CheckInt(2)
			}

			message := ""
			if argc > 2 {
				message = L.CheckString(3)
			}

			prev := false
			if argc > 3 {
				prev = L.CheckBool(4)
			}
			next := false
			if argc > 4 {
				next = L.CheckBool(5)
			}

			if ch.Listener != nil {
				ch.Listener.OnDialog(uint32(npc), message, prev, next)
			}
			return L.Yield(lua.LNumber(0))
		},
		"dialog_yes_no": func(L *lua.LState) int {
			argc := L.GetTop()
			ud := L.CheckUserData(1)
			ch, ok := ud.Value.(*Character)
			if !ok {
				L.ArgError(1, "Character expected")
				return 0
			}

			npc := 0
			if argc > 1 {
				npc = L.CheckInt(2)
			}

			message := ""
			if argc > 2 {
				message = L.CheckString(3)
			}

			prev := false
			if argc > 3 {
				prev = L.CheckBool(4)
			}
			next := false
			if argc > 4 {
				next = L.CheckBool(5)
			}

			if ch.Listener != nil {
				ch.Listener.OnDialogYesNo(uint32(npc), message, prev, next)
			}
			return L.Yield(lua.LNumber(0))
		},
		"dialog_list": func(L *lua.LState) int {
			argc := L.GetTop()
			ud := L.CheckUserData(1)
			ch, ok := ud.Value.(*Character)
			if !ok {
				L.ArgError(1, "Character expected")
				return 0
			}

			npc := 0
			if argc > 1 {
				npc = L.CheckInt(2)
			}

			message := ""
			if argc > 2 {
				message = L.CheckString(3)
			}

			selections := []string{}
			if argc > 3 {
				tbl := L.CheckTable(4)
				tbl.ForEach(func(_, value lua.LValue) {
					if str, ok := value.(lua.LString); ok {
						selections = append(selections, string(str))
					}
				})
			}

			if ch.Listener != nil {
				ch.Listener.OnDialogList(uint32(npc), message, selections)
			}
			return L.Yield(lua.LNumber(0))
		},
		"dialog_accept": func(L *lua.LState) int {
			argc := L.GetTop()
			ud := L.CheckUserData(1)
			ch, ok := ud.Value.(*Character)
			if !ok {
				L.ArgError(1, "Character expected")
				return 0
			}

			npc := 0
			if argc > 1 {
				npc = L.CheckInt(2)
			}

			message := ""
			if argc > 2 {
				message = L.CheckString(3)
			}

			enableEscape := false
			if argc > 3 {
				enableEscape = L.CheckBool(4)
			}

			if ch.Listener != nil {
				ch.Listener.OnDialogAccept(uint32(npc), message, enableEscape)
			}
			return L.Yield(lua.LNumber(0))
		},
		"dialog_input": func(L *lua.LState) int {
			argc := L.GetTop()
			ud := L.CheckUserData(1)
			ch, ok := ud.Value.(*Character)
			if !ok {
				L.ArgError(1, "Character expected")
				return 0
			}

			npc := 0
			if argc > 1 {
				npc = L.CheckInt(2)
			}

			message := ""
			if argc > 2 {
				message = L.CheckString(3)
			}

			if ch.Listener != nil {
				ch.Listener.OnDialogInput(uint32(npc), message)
			}
			return L.Yield(lua.LNumber(0))
		},
		"notice": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			ch, ok := ud.Value.(*Character)
			if !ok {
				L.ArgError(1, "Character expected")
				return 0
			}
			message := L.CheckString(2)
			if ch.Listener != nil {
				ch.Listener.OnMessage(constant.MSG_LIGHT_BLUE_TEXT, message)
			}
			return 0
		},
		"skill": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			ch, ok := ud.Value.(*Character)
			if !ok {
				L.ArgError(1, "Character expected")
				return 0
			}

			skillID := uint32(L.CheckInt(2))

			if ch.SkillsMap == nil {
				L.Push(lua.LNil)
				return 1
			}

			skillEntry, exists := ch.SkillsMap[skillID]
			if !exists || skillEntry == nil {
				L.Push(lua.LNil)
				return 1
			}

			skillUD := luax.NewLuable(L, skillEntry)
			L.Push(skillUD)
			return 1
		},
		"get_hp": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			ch, ok := ud.Value.(*Character)
			if !ok {
				L.ArgError(1, "Character expected")
				return 0
			}
			L.Push(lua.LNumber(ch.Hp))
			return 1
		},
		"map": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			ch, ok := ud.Value.(*Character)
			if !ok {
				L.ArgError(1, "Character expected")
				return 0
			}
			if ch.Context == nil {
				L.Push(lua.LNil)
				return 1
			}
			mapInstance := ch.Context.GetMap(ch.Map)
			if mapInstance == nil {
				L.Push(lua.LNil)
				return 1
			}
			L.Push(luax.NewLuable(L, mapInstance))
			return 1
		},
		"get_max_hp": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			ch, ok := ud.Value.(*Character)
			if !ok {
				L.ArgError(1, "Character expected")
				return 0
			}
			L.Push(lua.LNumber(ch.Life.GetMaxHp()))
			return 1
		},
		"set_hp": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			ch, ok := ud.Value.(*Character)
			if !ok {
				L.ArgError(1, "Character expected")
				return 0
			}
			hp := int(L.CheckNumber(2))
			if hp < 0 {
				hp = 0
			}
			maxHp := ch.Life.GetMaxHp()
			if hp > int(maxHp) {
				hp = int(maxHp)
			}
			ch.Hp = uint16(hp)
			if ch.Listener != nil {
				ch.Listener.OnUpdateStats(map[constant.Stat]int32{
					constant.STAT_HP: int32(ch.Hp),
				}, false)
			}
			return 0
		},
		"add_hp": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			ch, ok := ud.Value.(*Character)
			if !ok {
				L.ArgError(1, "Character expected")
				return 0
			}
			amount := int(L.CheckNumber(2))
			newHp := int(ch.Hp) + amount
			if newHp < 0 {
				newHp = 0
			}
			maxHp := ch.Life.GetMaxHp()
			if newHp > int(maxHp) {
				newHp = int(maxHp)
			}
			ch.Hp = uint16(newHp)
			if ch.Listener != nil {
				ch.Listener.OnUpdateStats(map[constant.Stat]int32{
					constant.STAT_HP: int32(ch.Hp),
				}, false)
			}
			return 0
		},
		"get_mp": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			ch, ok := ud.Value.(*Character)
			if !ok {
				L.ArgError(1, "Character expected")
				return 0
			}
			L.Push(lua.LNumber(ch.Mp))
			return 1
		},
		"get_max_mp": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			ch, ok := ud.Value.(*Character)
			if !ok {
				L.ArgError(1, "Character expected")
				return 0
			}
			L.Push(lua.LNumber(ch.Life.GetMaxMp()))
			return 1
		},
		"set_mp": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			ch, ok := ud.Value.(*Character)
			if !ok {
				L.ArgError(1, "Character expected")
				return 0
			}
			mp := int(L.CheckNumber(2))
			if mp < 0 {
				mp = 0
			}
			maxMp := ch.Life.GetMaxMp()
			if mp > int(maxMp) {
				mp = int(maxMp)
			}
			ch.Mp = uint16(mp)
			if ch.Listener != nil {
				ch.Listener.OnUpdateStats(map[constant.Stat]int32{
					constant.STAT_MP: int32(ch.Mp),
				}, false)
			}
			return 0
		},
		"add_mp": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			ch, ok := ud.Value.(*Character)
			if !ok {
				L.ArgError(1, "Character expected")
				return 0
			}
			amount := int(L.CheckNumber(2))
			newMp := int(ch.Mp) + amount
			if newMp < 0 {
				newMp = 0
			}
			maxMp := ch.Life.GetMaxMp()
			if newMp > int(maxMp) {
				newMp = int(maxMp)
			}
			ch.Mp = uint16(newMp)
			if ch.Listener != nil {
				ch.Listener.OnUpdateStats(map[constant.Stat]int32{
					constant.STAT_MP: int32(ch.Mp),
				}, false)
			}
			return 0
		},
		"add_mp_hp": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			ch, ok := ud.Value.(*Character)
			if !ok {
				L.ArgError(1, "Character expected")
				return 0
			}
			hpChange := int(L.CheckNumber(2))
			mpChange := int(L.CheckNumber(3))

			newHp := int(ch.Hp) + hpChange
			if newHp < 0 {
				newHp = 0
			}
			maxHp := ch.Life.GetMaxHp()
			if newHp > int(maxHp) {
				newHp = int(maxHp)
			}
			ch.Hp = uint16(newHp)

			newMp := int(ch.Mp) + mpChange
			if newMp < 0 {
				newMp = 0
			}
			maxMp := ch.Life.GetMaxMp()
			if newMp > int(maxMp) {
				newMp = int(maxMp)
			}
			ch.Mp = uint16(newMp)

			if ch.Listener != nil {
				stats := map[constant.Stat]int32{
					constant.STAT_HP: int32(ch.Hp),
					constant.STAT_MP: int32(ch.Mp),
				}
				ch.Listener.OnUpdateStats(stats, false)
			}
			return 0
		},
		"is_alive": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			ch, ok := ud.Value.(*Character)
			if !ok {
				L.ArgError(1, "Character expected")
				return 0
			}
			L.Push(lua.LBool(ch.Hp > 0))
			return 1
		},
		"get_bonus_str": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			ch, ok := ud.Value.(*Character)
			if !ok {
				L.ArgError(1, "Character expected")
				return 0
			}
			L.Push(lua.LNumber(ch.BonusStats.Str))
			return 1
		},
		"set_bonus_str": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			ch, ok := ud.Value.(*Character)
			if !ok {
				L.ArgError(1, "Character expected")
				return 0
			}
			bonusStr := L.CheckInt(2)
			ch.BonusStats.Str = int16(bonusStr)
			ch.notifyStatChange(constant.STAT_STR)
			return 0
		},
		"get_bonus_dex": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			ch, ok := ud.Value.(*Character)
			if !ok {
				L.ArgError(1, "Character expected")
				return 0
			}
			L.Push(lua.LNumber(ch.BonusStats.Dex))
			return 1
		},
		"set_bonus_dex": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			ch, ok := ud.Value.(*Character)
			if !ok {
				L.ArgError(1, "Character expected")
				return 0
			}
			bonusDex := L.CheckInt(2)
			ch.BonusStats.Dex = int16(bonusDex)
			ch.notifyStatChange(constant.STAT_DEX)
			return 0
		},
		"get_bonus_int": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			ch, ok := ud.Value.(*Character)
			if !ok {
				L.ArgError(1, "Character expected")
				return 0
			}
			L.Push(lua.LNumber(ch.BonusStats.Int))
			return 1
		},
		"set_bonus_int": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			ch, ok := ud.Value.(*Character)
			if !ok {
				L.ArgError(1, "Character expected")
				return 0
			}
			bonusInt := L.CheckInt(2)
			ch.BonusStats.Int = int16(bonusInt)
			ch.notifyStatChange(constant.STAT_INT)
			return 0
		},
		"get_bonus_luk": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			ch, ok := ud.Value.(*Character)
			if !ok {
				L.ArgError(1, "Character expected")
				return 0
			}
			L.Push(lua.LNumber(ch.BonusStats.Luk))
			return 1
		},
		"set_bonus_luk": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			ch, ok := ud.Value.(*Character)
			if !ok {
				L.ArgError(1, "Character expected")
				return 0
			}
			bonusLuk := L.CheckInt(2)
			ch.BonusStats.Luk = int16(bonusLuk)
			ch.notifyStatChange(constant.STAT_LUK)
			return 0
		},
		"get_bonus_hp": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			ch, ok := ud.Value.(*Character)
			if !ok {
				L.ArgError(1, "Character expected")
				return 0
			}
			L.Push(lua.LNumber(ch.Life.BonusHp))
			return 1
		},
		"set_bonus_hp": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			ch, ok := ud.Value.(*Character)
			if !ok {
				L.ArgError(1, "Character expected")
				return 0
			}
			bonusHp := L.CheckInt(2)
			ch.Life.BonusHp = int16(bonusHp)
			// Adjust HP if it exceeds new max
			if ch.Hp > ch.Life.GetMaxHp() {
				ch.Hp = ch.Life.GetMaxHp()
			}
			ch.notifyStatChange(constant.STAT_MAX_HP)
			return 0
		},
		"get_bonus_mp": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			ch, ok := ud.Value.(*Character)
			if !ok {
				L.ArgError(1, "Character expected")
				return 0
			}
			L.Push(lua.LNumber(ch.Life.BonusMp))
			return 1
		},
		"set_bonus_mp": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			ch, ok := ud.Value.(*Character)
			if !ok {
				L.ArgError(1, "Character expected")
				return 0
			}
			bonusMp := L.CheckInt(2)
			ch.Life.BonusMp = int16(bonusMp)
			// Adjust MP if it exceeds new max
			if ch.Mp > ch.Life.GetMaxMp() {
				ch.Mp = ch.Life.GetMaxMp()
			}
			ch.notifyStatChange(constant.STAT_MAX_MP)
			return 0
		},
	}
}

func (ch *Character) String() string {
	return ch.LuaTypeName()
}

func (ch *Character) Type() lua.LValueType {
	return lua.LTUserData
}

func NewDummyCharacter(sender Sendable, listener CharacterListener, id uint32, name string, ctx GameContext) Character {
	ch := Character{
		Sendable: sender,
		Listener: listener,
		Life: Life{
			Object: Object{
				Context: ctx,
			},
			Hp:     50,
			BaseHp: 50,
			Mp:     5,
			BaseMp: 5,
		},
		ID:           id,
		Name:         name,
		Gender:       0,
		SkinColor:    0,
		Face:         20100,
		Hair:         30000,
		Level:        1,
		Class:        0,
		Str:          12,
		Dex:          5,
		Int:          4,
		Luk:          4,
		AbilityPoint: 0,
		SkillPoint:   0,
		HpApUsed:     0,
		SpawnPoint:   1,
		Map:          200000301,
		Meso:         2135983647,

		Random1: stream.NewRandomStream(),
		Random2: stream.NewRandomStream(),
		Random3: stream.NewRandomStream(),

		Inventory: map[constant.InventoryType]*Inventory{
			constant.INVENTORY_TYPE_EQUIPMENT:    NewInventory(constant.INVENTORY_TYPE_EQUIPMENT),
			constant.INVENTORY_TYPE_CONSUME:      NewInventory(constant.INVENTORY_TYPE_CONSUME),
			constant.INVENTORY_TYPE_INSTALLATION: NewInventory(constant.INVENTORY_TYPE_INSTALLATION),
			constant.INVENTORY_TYPE_ETC:          NewInventory(constant.INVENTORY_TYPE_ETC),
			constant.INVENTORY_TYPE_CASH:         NewInventory(constant.INVENTORY_TYPE_CASH),
		},
		Rings: RingContainer{
			Left:  []*Ring{},
			Right: []*Ring{},
			Mid:   []*Ring{},
		},
		Equipments: map[constant.EquipmentPartsType]*Equipment{
			constant.EQUIPMENT_PARTS_WEAPON: nil,
			constant.EQUIPMENT_PARTS_SHIELD: nil,
		},

		RegRocks: []uint32{999999999, 999999999, 999999999, 999999999, 999999999},
		Rocks:    []uint32{999999999, 999999999, 999999999, 999999999, 999999999, 999999999, 999999999, 999999999, 999999999, 999999999},
	}

	if ctx != nil {
		resources := ctx.GetResources()
		ch.Equipments[constant.EQUIPMENT_PARTS_WEAPON] = &Equipment{
			ItemCore: &ItemCore{
				Wz:         resources.Items[1302000],
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
		ch.Inventory[constant.INVENTORY_TYPE_CASH].Items[1] = &Pet{
			ItemCore: &ItemCore{
				Wz:         resources.Items[5000007],
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

		ch.Inventory[constant.INVENTORY_TYPE_ETC].Items[1] = &GeneralItem{
			ItemCore: &ItemCore{
				Wz:         resources.Items[4000001],
				Count:      100,
				Expiration: util.TimeMax,
			},
		}

		ch.Inventory[constant.INVENTORY_TYPE_EQUIPMENT].Items[1], err = NewItem(1060002, 1, ctx)
		ch.Inventory[constant.INVENTORY_TYPE_EQUIPMENT].Items[2], err = NewItem(1060006, 1, ctx)
		ch.Inventory[constant.INVENTORY_TYPE_EQUIPMENT].Items[3], err = NewItem(1040002, 1, ctx)
		ch.Inventory[constant.INVENTORY_TYPE_EQUIPMENT].Items[4], err = NewItem(1040010, 1, ctx)
	}

	return ch
}

// GetCurrentDialog returns the current dialog coroutine
func (ch *Character) GetCurrentDialog() *lua.LState {
	ch.dialogMutex.Lock()
	defer ch.dialogMutex.Unlock()
	return ch.currentDialog
}

// SetCurrentDialog sets the current dialog coroutine
func (ch *Character) SetCurrentDialog(dialog *lua.LState) {
	ch.dialogMutex.Lock()
	defer ch.dialogMutex.Unlock()
	ch.currentDialog = dialog
}

// ClearCurrentDialog clears the current dialog coroutine
func (ch *Character) ClearCurrentDialog() {
	ch.dialogMutex.Lock()
	defer ch.dialogMutex.Unlock()
	ch.currentDialog = nil
}

// tryLevelUp checks if the character can level up and performs level up if possible
// Returns true if level up occurred, false otherwise
func (ch *Character) tryLevelUp() bool {
	if ch.Context == nil {
		return false
	}

	resources := ch.Context.GetResources()
	if resources == nil {
		return false
	}

	if ch.Level >= 200 {
		return false
	}

	oldLevel := ch.Level
	remainingExp := ch.Exp
	targetLevel := ch.Level

	// Calculate the maximum level we can reach with current experience
	for targetLevel < 200 {
		expNeeded := resources.GetExpNeededForLevel(targetLevel)
		if expNeeded == 0 {
			break
		}

		if remainingExp < expNeeded {
			break
		}

		remainingExp -= expNeeded
		targetLevel++
	}

	if targetLevel == oldLevel {
		return false
	}

	// Update experience
	ch.Exp = remainingExp

	// Set level once with all stat increases
	ch.SetLevel(targetLevel)

	return true
}

// SetLevel sets the character's level and applies all stat changes
// This simulates leveling up from the current level to the target level
func (ch *Character) SetLevel(newLevel uint8) {
	if newLevel < 1 {
		newLevel = 1
	}
	if newLevel > 200 {
		newLevel = 200
	}

	if newLevel == ch.Level {
		return
	}

	oldLevel := ch.Level
	levelDiff := int(newLevel) - int(oldLevel)

	if levelDiff > 0 {
		totalAPIncrease := uint16(0)
		totalSPIncrease := uint16(0)
		totalHPIncrease := uint16(0)
		totalMPIncrease := uint16(0)

		for level := oldLevel + 1; level <= newLevel; level++ {
			totalAPIncrease += 5

			if !ch.IsBeginner() {
				totalSPIncrease += 3
			}

			hpIncrease := uint16(20 + int(level)*2)
			mpIncrease := uint16(10 + int(level))
			totalHPIncrease += hpIncrease
			totalMPIncrease += mpIncrease
		}

		ch.Level = newLevel
		ch.AbilityPoint += totalAPIncrease
		ch.SkillPoint += totalSPIncrease

		ch.Life.AddBaseHp(totalHPIncrease)
		ch.Life.AddBaseMp(totalMPIncrease)

		// Restore HP/MP to max
		ch.Hp = ch.Life.GetMaxHp()
		ch.Mp = ch.Life.GetMaxMp()

		// Notify listener
		if ch.Listener != nil {
			stats := map[constant.Stat]int32{
				constant.STAT_LEVEL:        int32(ch.Level),
				constant.STAT_EXP:          int32(ch.Exp),
				constant.STAT_MAX_HP:       int32(ch.Life.GetMaxHp()),
				constant.STAT_MAX_MP:       int32(ch.Life.GetMaxMp()),
				constant.STAT_HP:           int32(ch.Hp),
				constant.STAT_MP:           int32(ch.Mp),
				constant.STAT_AVAILABLE_AP: int32(ch.AbilityPoint),
				constant.STAT_AVAILABLE_SP: int32(ch.SkillPoint),
			}
			ch.Listener.OnUpdateStats(stats, false)

			// Broadcast level up effect for each level gained
			for level := oldLevel + 1; level <= newLevel; level++ {
				ch.broadcastLevelUpEffect()
			}
		}
	} else {
		// Leveling down: just set the level and exp, don't decrease stats
		ch.Level = newLevel

		if ch.Context != nil {
			resources := ch.Context.GetResources()
			if resources != nil {
				if newLevel > 1 {
					ch.Exp = resources.GetExpNeededForLevel(newLevel - 1)
				} else {
					ch.Exp = 0
				}
			}
		}

		if ch.Listener != nil {
			stats := map[constant.Stat]int32{
				constant.STAT_LEVEL: int32(ch.Level),
				constant.STAT_EXP:   int32(ch.Exp),
			}
			ch.Listener.OnUpdateStats(stats, false)
		}
	}
}

// broadcastLevelUpEffect broadcasts level up effect to other players on the map
func (ch *Character) broadcastLevelUpEffect() {
	if ch.Context == nil {
		return
	}

	mapInstance := ch.Context.GetMap(ch.Map)
	if mapInstance == nil {
		return
	}

	// Create SHOW_FOREIGN_EFFECT packet for level up (EffectID 0)
	levelUpPacket := &response.ShowForeignEffect{
		CharacterID: ch.ID,
		EffectID:    0, // Level up effect
	}

	// Broadcast to all players on the map except the character who leveled up
	mapInstance.BroadcastToPlayers(levelUpPacket, types.SEND_POLICY_ENCRYPT, ch.ID)
}
