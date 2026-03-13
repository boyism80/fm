package entity

import (
	"fmt"
	"sync"
	"time"

	"github.com/asynkron/protoactor-go/actor"
	c_actor "github.com/boyism80/fm/core/actor"
	"github.com/boyism80/fm/core/luax"
	"github.com/boyism80/fm/game/constant"
	"github.com/boyism80/fm/protocol/response"
	"github.com/boyism80/fm/stream"
	"github.com/boyism80/fm/types"
	"github.com/boyism80/fm/util"
	lua "github.com/yuin/gopher-lua"
)

type Character struct {
	Life
	Sendable

	dialog           *lua.LState
	id               uint32
	name             string
	gender           uint8
	skinColor        uint8
	face             uint32
	hair             uint32
	level            uint8
	rank             uint32
	rankDiff         int32
	classRank        uint32
	classRankDiff    int32
	exp              uint32
	famePoint        uint16
	spawnPoint       uint8
	mega             bool
	random1          stream.RandomStream
	random2          stream.RandomStream
	random3          stream.RandomStream
	questStatuses    map[int]*QuestStatus
	marriageId       uint32
	regRocks         []uint32
	rocks            []uint32
	monsterBookCover uint32
	monsterBook      *MonsterBook
	quests           map[uint16]string
	currentDialog    *lua.LState
	dialogMutex      sync.Mutex
	hidden           bool

	Listener       CharacterListener
	Class          uint16
	Role           constant.CharacterRole
	AbilityPoint   uint16
	SkillPoint     uint16
	HpApUsed       uint16
	Meso           int32
	Inventory      map[constant.InventoryType]*Inventory
	Equipments     map[constant.EquipmentPartsType]Equipment
	Rings          RingContainer
	Skills         map[uint32]*SkillEntry
	CurrentShopID  uint32
	Chair          uint32
	LastHealHPTime time.Time // used for heal-over-time rate limit
	LastHealMPTime time.Time // used for heal-over-time rate limit
	BaseStats      BaseStats
	BonusStats     BonusStats
	Buffs          *BuffContainer
	timers         map[string]*CharacterTimer
}

type CharacterTimer struct {
	Timer      *time.Timer
	Interval   time.Duration
	Repeat     bool
	Callback   func()
	NextFireAt time.Time
	Remaining  time.Duration
}

func (ch *Character) GetObjectType() constant.ObjectType {
	return constant.ObjectTypeCharacter
}

func (ch *Character) Is(typ constant.ObjectType) bool {
	return ch.GetObjectType().Has(typ)
}

func (ch *Character) AddTimerWithCallback(key string, interval time.Duration, repeat bool, callback func()) bool {
	return ch.addTimer(key, interval, repeat, callback)
}

func (ch *Character) addTimer(key string, interval time.Duration, repeat bool, callback func()) bool {
	if ch.timers == nil {
		ch.timers = make(map[string]*CharacterTimer)
	}
	if _, exists := ch.timers[key]; exists {
		ch.RemoveTimer(key)
	}
	if ch.Context == nil {
		return false
	}
	m := ch.GetMap()
	if m == nil {
		return false
	}
	pid := m.GetActorPID()
	if pid == nil {
		return false
	}
	characterID := ch.GetID()
	entry := &CharacterTimer{
		Interval:   interval,
		Repeat:     repeat,
		Callback:   callback,
		NextFireAt: time.Now().Add(interval),
	}
	entry.Timer = time.AfterFunc(interval, func() {
		ch.Context.SendToActor(pid, &c_actor.RunCharacterTimer{CharacterID: characterID, Key: key})
	})
	ch.timers[key] = entry
	return true
}

func (ch *Character) RemoveTimer(key string) bool {
	if ch.timers == nil {
		return false
	}
	entry := ch.timers[key]
	if entry == nil {
		return false
	}
	if entry.Timer != nil {
		entry.Timer.Stop()
	}
	delete(ch.timers, key)
	return true
}

func (ch *Character) GetTimerEntry(key string) *CharacterTimer {
	if ch.timers == nil {
		return nil
	}
	return ch.timers[key]
}

func (ch *Character) ClearTimers() {
	if ch.timers == nil {
		return
	}
	for key, entry := range ch.timers {
		if entry != nil && entry.Timer != nil {
			entry.Timer.Stop()
		}
		delete(ch.timers, key)
	}
}

func (ch *Character) SuspendTimers() {
	if ch.timers == nil {
		return
	}
	now := time.Now()
	for _, entry := range ch.timers {
		if entry == nil || entry.Timer == nil {
			continue
		}
		entry.Remaining = entry.NextFireAt.Sub(now)
		if entry.Remaining < 0 {
			entry.Remaining = 0
		}
		entry.Timer.Stop()
		entry.Timer = nil
	}
}

func (ch *Character) ResumeTimers(pid *actor.PID) {
	if ch.timers == nil || pid == nil || ch.Context == nil {
		return
	}
	characterID := ch.GetID()
	for key, entry := range ch.timers {
		if entry == nil || entry.Timer != nil {
			continue
		}
		duration := entry.Remaining
		if duration <= 0 {
			duration = entry.Interval
		}
		entry.Remaining = 0
		entry.NextFireAt = time.Now().Add(duration)
		k := key
		entry.Timer = time.AfterFunc(duration, func() {
			ch.Context.SendToActor(pid, &c_actor.RunCharacterTimer{CharacterID: characterID, Key: k})
		})
	}
}

func (ch *Character) GetObject() *Object {
	return &ch.Life.Object
}

func (ch *Character) GetHp() uint16       { return ch.Life.Hp }
func (ch *Character) GetMp() uint16       { return ch.Life.Mp }
func (ch *Character) GetBonusHp() int16   { return ch.Life.BonusHp }
func (ch *Character) GetBonusMp() int16   { return ch.Life.BonusMp }
func (ch *Character) GetInvincible() bool { return ch.Life.Invincible }
func (ch *Character) IsAlive() bool       { return ch.Life.Hp > 0 }

const characterMaxHpMpCap = 32767

func (ch *Character) SetHp(v uint16, notify bool) {
	ch.Life.SetHp(v, false)
	if notify && ch.Listener != nil {
		ch.Listener.OnUpdateStats(map[constant.Stat]int32{constant.STAT_HP: int32(ch.Hp)}, false)
	}
}

func (ch *Character) SetMp(v uint16, notify bool) {
	ch.Life.SetMp(v, false)
	if notify && ch.Listener != nil {
		ch.Listener.OnUpdateStats(map[constant.Stat]int32{constant.STAT_MP: int32(ch.Mp)}, false)
	}
}

func (ch *Character) SetMaxHp(v uint16, notify bool) {
	if v > characterMaxHpMpCap {
		v = characterMaxHpMpCap
	}
	ch.Life.SetMaxHp(v, false)
	if ch.Life.Hp > ch.GetMaxHp() {
		ch.Life.Hp = ch.GetMaxHp()
	}
	if notify && ch.Listener != nil {
		ch.Listener.OnUpdateStats(map[constant.Stat]int32{
			constant.STAT_HP:     int32(ch.Hp),
			constant.STAT_MAX_HP: int32(ch.GetMaxHp()),
		}, false)
	}
}

func (ch *Character) SetMaxMp(v uint16, notify bool) {
	if v > characterMaxHpMpCap {
		v = characterMaxHpMpCap
	}
	ch.Life.SetMaxMp(v, false)
	if ch.Life.Mp > ch.GetMaxMp() {
		ch.Life.Mp = ch.GetMaxMp()
	}
	if notify && ch.Listener != nil {
		ch.Listener.OnUpdateStats(map[constant.Stat]int32{
			constant.STAT_MP:     int32(ch.Mp),
			constant.STAT_MAX_MP: int32(ch.GetMaxMp()),
		}, false)
	}
}

func (ch *Character) SetBonusHp(v int16) {
	ch.Life.SetBonusHp(v)
	if ch.Listener != nil {
		ch.notifyStatChange(constant.STAT_MAX_HP)
	}
}

func (ch *Character) SetBonusMp(v int16) {
	ch.Life.SetBonusMp(v)
	if ch.Listener != nil {
		ch.notifyStatChange(constant.STAT_MAX_MP)
	}
}

func (ch *Character) SetInvincible(b bool) { ch.Life.Invincible = b }

func (ch *Character) AddHp(amount int) {
	ch.Life.AddHp(amount)
	if ch.Listener != nil {
		ch.Listener.OnUpdateStats(map[constant.Stat]int32{constant.STAT_HP: int32(ch.Hp)}, false)
	}
}

func (ch *Character) AddMp(amount int) {
	ch.Life.AddMp(amount)
	if ch.Listener != nil {
		ch.Listener.OnUpdateStats(map[constant.Stat]int32{constant.STAT_MP: int32(ch.Mp)}, false)
	}
}

func (ch *Character) AddHpMp(hpDelta, mpDelta int) {
	ch.Life.AddHpMp(hpDelta, mpDelta)
	if ch.Listener != nil {
		ch.Listener.OnUpdateStats(map[constant.Stat]int32{
			constant.STAT_HP: int32(ch.Hp),
			constant.STAT_MP: int32(ch.Mp),
		}, false)
	}
}

func (ch *Character) GetBaseHp() uint16 {
	return ch.Life.BaseHp
}

func (ch *Character) SetBaseHp(v uint16, notify bool) {
	if v > constant.STAT_MAX_HP_MP {
		v = constant.STAT_MAX_HP_MP
	}
	ch.Life.BaseHp = v
	if ch.Life.Hp > ch.GetMaxHp() {
		ch.Life.Hp = ch.GetMaxHp()
	}
	if notify && ch.Listener != nil {
		ch.Listener.OnUpdateStats(map[constant.Stat]int32{
			constant.STAT_HP:     int32(ch.Hp),
			constant.STAT_MAX_HP: int32(ch.GetMaxHp()),
		}, false)
	}
}

func (ch *Character) GetBaseMp() uint16 {
	return ch.Life.BaseMp
}

func (ch *Character) SetBaseMp(v uint16, notify bool) {
	if v > constant.STAT_MAX_HP_MP {
		v = constant.STAT_MAX_HP_MP
	}
	ch.Life.BaseMp = v
	if ch.Life.Mp > ch.GetMaxMp() {
		ch.Life.Mp = ch.GetMaxMp()
	}
	if notify && ch.Listener != nil {
		ch.Listener.OnUpdateStats(map[constant.Stat]int32{
			constant.STAT_MP:     int32(ch.Mp),
			constant.STAT_MAX_MP: int32(ch.GetMaxMp()),
		}, false)
	}
}

func (ch *Character) SetAbilityPoint(v uint16, notify bool) {
	ch.AbilityPoint = v
	if notify && ch.Listener != nil {
		ch.Listener.OnUpdateStats(map[constant.Stat]int32{
			constant.STAT_AVAILABLE_AP: int32(ch.AbilityPoint),
		}, false)
	}
}

func (ch *Character) SetSkillPoint(v uint16, notify bool) {
	ch.SkillPoint = v
	if notify && ch.Listener != nil {
		ch.Listener.OnUpdateStats(map[constant.Stat]int32{
			constant.STAT_AVAILABLE_SP: int32(ch.SkillPoint),
		}, false)
	}
}

func (ch *Character) GetMap() *Map {
	return ch.GetObject().GetMap()
}

func (ch *Character) Warp(targetMap *Map, spawnPoint uint8) error {
	if ch.Context == nil {
		return fmt.Errorf("no game context")
	}
	return ch.Context.RequestWarp(ch, targetMap, spawnPoint)
}

func (ch *Character) Send(p types.Packet, policy types.SendPolicy) error {
	if ch.Sendable == nil {
		return nil
	}
	return ch.Sendable.Send(p, policy)
}

func (ch *Character) IsHidden() bool {
	return ch.hidden
}

func (ch *Character) SetHidden(hidden bool) {
	if ch.hidden == hidden {
		return
	}
	ch.hidden = hidden
	if ch.Listener != nil {
		ch.Listener.OnHiddenChanged(hidden)
	}
}

func (ch *Character) GetID() uint32 {
	return ch.id
}

func (ch *Character) Message(message string) {
	if ch.Listener != nil {
		ch.Listener.OnMessage(constant.MSG_LIGHT_BLUE_TEXT, message)
	}
}

func (ch *Character) HasRoleAtLeast(role constant.CharacterRole) bool {
	return ch.Role >= role
}

func (ch *Character) SetMeso(meso int32) {
	if meso < 0 {
		ch.Meso = 0
		return
	}
	ch.Meso = meso
}

func (ch *Character) IsRanked() bool {
	if ch.HasRoleAtLeast(constant.RoleAdmin) {
		return false
	}

	if ch.level < 30 {
		return false
	}

	return true
}

func (ch *Character) AddExp(exp uint32) {

	if ch.Context != nil {
		expRate := ch.Context.GetExpRate()
		if expRate > 0 {
			exp = exp * uint32(expRate)
		}
	}

	ch.exp += exp
	ch.Listener.OnExpGain(exp)

	if !ch.tryLevelUp() {
		ch.Listener.OnUpdateStats(map[constant.Stat]int32{
			constant.STAT_EXP: int32(ch.exp),
		}, false)
	}
}

func NewDummyCharacter(sender Sendable, listener CharacterListener, id uint32, name string, ctx GameContext) *Character {
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
		id:        id,
		name:      name,
		gender:    0,
		skinColor: 0,
		face:      20100,
		hair:      30000,
		level:     1,
		Class:     0,
		BaseStats: BaseStats{
			Str: 12,
			Dex: 5,
			Int: 4,
			Luk: 4,
		},
		AbilityPoint: 0,
		SkillPoint:   0,
		HpApUsed:     0,
		spawnPoint:   1,
		Meso:         2135983647,

		random1: stream.NewRandomStream(),
		random2: stream.NewRandomStream(),
		random3: stream.NewRandomStream(),

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
		Equipments: map[constant.EquipmentPartsType]Equipment{
			constant.EQUIPMENT_PARTS_WEAPON: nil,
			constant.EQUIPMENT_PARTS_SHIELD: nil,
		},
		Skills: make(map[uint32]*SkillEntry),

		regRocks: []uint32{999999999, 999999999, 999999999, 999999999, 999999999},
		rocks:    []uint32{999999999, 999999999, 999999999, 999999999, 999999999, 999999999, 999999999, 999999999, 999999999, 999999999},
	}
	ch.Buffs = NewBuffContainer(&ch)

	if ctx != nil {
		resources := ctx.GetResources()
		weaponItem, err := NewItem(1302000, 1, ctx)
		if err == nil {
			if eq, ok := weaponItem.(Equipment); ok {
				ch.Equipments[constant.EQUIPMENT_PARTS_WEAPON] = eq
			}
		}
		if ch.Equipments[constant.EQUIPMENT_PARTS_WEAPON] == nil && resources != nil {
			core := &EquipmentCore{
				ItemCore: &ItemCore{
					Wz:         resources.Items[1302000],
					Count:      1,
					UniqueId:   0,
					Expiration: util.TimeMax,
				},
				EnchantChance: 7,
			}
			ch.Equipments[constant.EQUIPMENT_PARTS_WEAPON] = &Weapon{EquipmentCore: core}
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

	return &ch
}

func (ch *Character) GetCurrentDialog() *lua.LState {
	ch.dialogMutex.Lock()
	defer ch.dialogMutex.Unlock()
	return ch.currentDialog
}

func (ch *Character) SetCurrentDialog(dialog *lua.LState) {
	ch.dialogMutex.Lock()
	defer ch.dialogMutex.Unlock()
	ch.currentDialog = dialog
}

func (ch *Character) ClearCurrentDialog() {
	ch.dialogMutex.Lock()
	defer ch.dialogMutex.Unlock()
	ch.currentDialog = nil
}

func (ch *Character) tryLevelUp() bool {
	if ch.Context == nil {
		return false
	}

	resources := ch.Context.GetResources()
	if resources == nil {
		return false
	}

	if ch.level >= 200 {
		return false
	}

	oldLevel := ch.level
	remainingExp := ch.exp
	targetLevel := ch.level

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

	mapInstance := ch.GetMap()
	if mapInstance == nil {
		return false
	}

	pid := mapInstance.GetActorPID()
	if pid == nil {
		return false
	}

	root := luax.GetRootLuaState(pid.String())
	if root == nil {
		return false
	}

	ch.exp = remainingExp
	ch.SetLevel(targetLevel)
	levelDiff := int(targetLevel) - int(oldLevel)
	if levelDiff >= 1 {
		_, thread, _ := luax.Call(root, "script/script.lua", "on_level_up", ch, int32(oldLevel), int32(targetLevel))
		if thread != nil {
			thread.Close()
		}

		if ch.Listener != nil {
			stats := map[constant.Stat]int32{
				constant.STAT_LEVEL:        int32(ch.level),
				constant.STAT_EXP:          int32(ch.exp),
				constant.STAT_MAX_HP:       int32(ch.GetMaxHp()),
				constant.STAT_MAX_MP:       int32(ch.GetMaxMp()),
				constant.STAT_HP:           int32(ch.Hp),
				constant.STAT_MP:           int32(ch.Mp),
				constant.STAT_AVAILABLE_AP: int32(ch.AbilityPoint),
				constant.STAT_AVAILABLE_SP: int32(ch.SkillPoint),
			}
			ch.Listener.OnUpdateStats(stats, false)
			for i := 0; i < levelDiff; i++ {
				ch.broadcastLevelUpEffect()
			}
		}
	}
	return true
}

func (ch *Character) SetLevel(newLevel uint8) {
	if newLevel < 1 {
		newLevel = 1
	}
	if newLevel > 200 {
		newLevel = 200
	}

	if newLevel == ch.level {
		return
	}

	oldLevel := ch.level
	ch.level = newLevel
	if newLevel < oldLevel {
		if ch.Context != nil {
			resources := ch.Context.GetResources()
			if resources != nil {
				if newLevel > 1 {
					ch.exp = resources.GetExpNeededForLevel(newLevel - 1)
				} else {
					ch.exp = 0
				}
			}
		}
	}

	if ch.Listener != nil {
		ch.Listener.OnUpdateStats(map[constant.Stat]int32{
			constant.STAT_LEVEL: int32(ch.level),
			constant.STAT_EXP:   int32(ch.exp),
		}, false)
	}
}

func (ch *Character) broadcastLevelUpEffect() {
	mapInstance := ch.GetMap()
	if mapInstance == nil {
		return
	}

	mapInstance.Broadcast(&response.ShowForeignEffect{
		CharacterID: ch.id,
		EffectID:    0,
	}, &BroadcastOption{
		ExceptPlayerIDs:    []uint32{ch.GetID()},
		ReferenceCharacter: ch,
		RecipientFilter:    BroadcastVisibleByReference,
	})
}
