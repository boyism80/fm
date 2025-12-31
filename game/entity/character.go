package entity

import (
	"fmt"
	"sync"
	"time"

	"github.com/boyism80/fm/game/constant"
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
	RegRocks         []uint32 // Basic warp rock slots (5 slots)?
	Rocks            []uint32 // VIP warp rock slots (10 slots)?
	MonsterBookCover uint32
	MonsterBook      *MonsterBook
	QuestInfo        map[uint16]string

	// Dialog state management
	currentDialog *lua.LState // Current dialog coroutine
	dialogMutex   sync.Mutex  // Mutex for dialog state access
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

// GetInventory returns the inventory for the given slot
// Slot format: 0xXXYY where XX is inventory type and YY is slot number
func (ch *Character) GetInventory(slot uint16) *Inventory {
	inventoryType := constant.InventoryType(slot >> 8)
	return ch.Inventory[inventoryType]
}

// GetID returns the character's ID
func (ch *Character) GetID() uint32 {
	return ch.ID
}

// Message sends a notice message to the character through the listener
func (ch *Character) Message(message string) {
	if ch.Listener != nil {
		ch.Listener.OnMessage(constant.MSG_LIGHT_BLUE_TEXT, message)
	}
}

// GetThreadHash returns the character's map ID for thread assignment
// Characters on the same map will be assigned to the same logic thread
func (ch *Character) GetThreadHash() int {
	return int(ch.Map)
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

// AddExp adds experience points to the character and notifies the listener
func (ch *Character) AddExp(exp uint32) {
	ch.Exp += exp

	// Notify listener about exp gain (following old server pattern)
	if ch.Listener != nil {
		ch.Listener.OnExpGain(exp)
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
