package entity

import (
	"fmt"
	"sync"
	"time"

	"github.com/asynkron/protoactor-go/actor"
	c_actor "github.com/boyism80/fm/core/actor"
	"github.com/boyism80/fm/core/luax"
	"github.com/boyism80/fm/protocol/dto"
	"github.com/boyism80/fm/protocol/response"
	"github.com/boyism80/fm/services/game/constant"
	"github.com/boyism80/fm/stream"
	"github.com/boyism80/fm/types"
	lua "github.com/yuin/gopher-lua"
)

type Character struct {
	LifeCore
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
	Listener         CharacterListener
	Class            uint16
	Role             constant.CharacterRole
	AccountID        uint32
	AbilityPoint     uint16
	SkillPoint       uint16
	HpApUsed         uint16
	Meso             int32
	Inventory        map[constant.InventoryType]*Inventory
	Equipments       map[constant.EquipmentPartsType]Equipment
	Rings            RingContainer
	Skills           *SkillContainer
	CurrentShopID    uint32
	Chair            uint32
	LastHealHPTime   time.Time
	LastHealMPTime   time.Time
	BaseStats        BaseStats
	BonusStats       BonusStats
	Buffs            *BuffContainer
	diseases         map[constant.DebuffFlag]*DiseaseValueHolder
	timers           map[string]*CharacterTimer
	summons          map[constant.SkillID]*Summon
	doors            map[constant.SkillID]*Door
	HomingTargetOID  *uint32
}

type DiseaseValueHolder struct {
	Disease   constant.DebuffFlag
	StartTime time.Time
	Duration  time.Duration
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

func (ch *Character) SendSpawnSyncToViewer(viewer *Character) {
	if ch == nil || viewer == nil {
		return
	}
	if ch.GetID() == viewer.GetID() {
		return
	}
	if ch.IsHidden() && !viewer.HasRoleAtLeast(ch.Role) {
		return
	}
	spawnBuffData := ch.GetSpawnPlayerBuffData()
	viewer.Send(&response.SpawnPlayer{
		Character:         ch.ToDTO(),
		BuffStates:        spawnBuffData.BuffStates,
		Diseases:          ch.GetDiseaseMask(),
		SpeedBuff:         spawnBuffData.SpeedBuff,
		ComboCount:        spawnBuffData.ComboCount,
		WKChargeSkillId:   spawnBuffData.WKChargeSkillID,
		MorphId:           spawnBuffData.MorphID,
		SpiritClawSkillId: spawnBuffData.SpiritClawSkillID,
		MountLevel:        spawnBuffData.MountLevel,
		MountExp:          spawnBuffData.MountExp,
		MountFatigue:      spawnBuffData.MountFatigue,
		CrushRings:        RingsToDTO(ch.Rings.Left),
		FriendshipRings:   RingsToDTO(ch.Rings.Mid),
		MarriageRings:     RingsToDTO(ch.Rings.Right),
	}, types.SEND_POLICY_ENCRYPT)

	if mountID, active := ch.GetRiddingInfo(); active {
		viewer.Send(&response.UpdateRemoteRidding{
			CharacterID: int32(ch.GetID()),
			MountID:     mountID,
			Buffs:       []dto.BuffEntry{{Buff: constant.BuffFlagMonsterRiding, Value: mountID}},
		}, types.SEND_POLICY_ENCRYPT)
	}

	if buff, _, active := ch.Buffs.GetBuffValue(constant.BuffFlagEnergyCharge); active {
		viewer.Send(&response.UpdateRemoteBuff{
			CharacterID: int32(ch.GetID()),
			BuffID:      buff.GetBuffID(),
			Duration:    50 * time.Second,
			Buffs:       []dto.BuffEntry{{Buff: constant.BuffFlagEnergyCharge, Value: 10000}},
		}, types.SEND_POLICY_ENCRYPT)
	}

	if dashBuff, dashValue, active := ch.Buffs.GetBuffValue(constant.BuffFlagDashSpeed); active {
		buffs := []dto.BuffEntry{{Buff: constant.BuffFlagDashSpeed, Value: dashValue}}
		if _, dashJumpValue, hasDashJump := ch.Buffs.GetBuffValue(constant.BuffFlagDashJump); hasDashJump {
			buffs = append(buffs, dto.BuffEntry{Buff: constant.BuffFlagDashJump, Value: dashJumpValue})
		}
		viewer.Send(&response.UpdateRemoteBuff{
			CharacterID: int32(ch.GetID()),
			BuffID:      dashBuff.GetBuffID(),
			Duration:    dashBuff.RemainingDuration(time.Now()),
			Buffs:       buffs,
		}, types.SEND_POLICY_ENCRYPT)
	}
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
		ch.Context.DispatchRunCharacterTimer(pid, &c_actor.RunCharacterTimer{CharacterID: characterID, Key: key})
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

func summonTimerKey(skillID constant.SkillID) string {
	return fmt.Sprintf("summon:%d", skillID)
}

func (ch *Character) SpawnSummon(skillID constant.SkillID, skillLevel uint8, movementType constant.SummonMovementType, summonType constant.SummonType, position types.Point[int16], duration time.Duration) *Summon {
	if ch == nil {
		return nil
	}
	m := ch.GetMap()
	if m == nil {
		return nil
	}
	if ch.summons != nil {
		if current, ok := ch.summons[skillID]; ok && current != nil {
			ch.RemoveSummon(current, true)
		}
	}

	s := &Summon{
		LifeCore: LifeCore{
			ObjectCore: ObjectCore{
				Position: position,
				Context:  m.context,
				Map:      nil,
			},
			Hp:     1,
			BaseHp: 1,
			BaseMp: 1,
		},
		Owner:        ch,
		OwnerID:      ch.GetID(),
		SkillID:      skillID,
		SkillLevel:   skillLevel,
		MovementType: movementType,
		SummonType:   summonType,
	}
	s.LifeCore.ObjectCore.self = s
	if ch.summons == nil {
		ch.summons = make(map[constant.SkillID]*Summon)
	}
	ch.summons[skillID] = s
	m.AddSummon(s)
	if duration > 0 {
		_ = ch.AddTimerWithCallback(summonTimerKey(s.SkillID), duration, false, func() {
			ch.handleSummonExpireBySkill(s.SkillID)
		})
	}
	return s
}

func (ch *Character) SpawnMist(skill *SkillEntry, position types.Point[int16], mistType constant.MistType, bounds types.Rect[int32], duration time.Duration, initialDelay time.Duration, poisonTickMultiplier float64) *Mist {
	if ch == nil || ch.Context == nil {
		return nil
	}
	if skill == nil {
		return nil
	}
	m := ch.GetMap()
	if m == nil {
		return nil
	}
	wzSkill := ch.Context.GetResources().GetSkill(skill.Wz.ID)
	if wzSkill == nil {
		return nil
	}
	skillLevel := uint8(skill.Level())
	ld := wzSkill.GetLevelData(int(skillLevel))
	b := bounds
	if b.Left == 0 && b.Right == 0 && b.Top == 0 && b.Bottom == 0 {
		b = MistWorldBounds(position, ld)
	}
	if poisonTickMultiplier <= 0 {
		poisonTickMultiplier = 1.0
	}
	mist := &Mist{
		ObjectCore: ObjectCore{
			Position: position,
			Context:  m.context,
			Map:      nil,
		},
		Causer:               ch.GetID(),
		SkillWz:              wzSkill,
		SkillLevel:           skillLevel,
		MistType:             mistType,
		MobMist:              false,
		MobSkill:             false,
		SkillDelay:           8,
		Bounds:               b,
		ExpiresAt:            time.Time{},
		PoisonTickMultiplier: poisonTickMultiplier,
	}
	mist.ObjectCore.self = mist
	if initialDelay > 0 {
		mist.NextPoisonTickAt = time.Now().Add(initialDelay)
	}
	if duration > 0 {
		mist.ExpiresAt = time.Now().Add(duration)
	}
	m.AddMist(mist)
	return mist
}

func (ch *Character) SpawnDoor(skillID constant.SkillID) *Door {
	if ch == nil {
		return nil
	}
	m := ch.GetMap()
	if m == nil || m.Wz == nil {
		return nil
	}
	if ch.doors != nil {
		if current, ok := ch.doors[skillID]; ok && current != nil {
			ch.RemoveDoor(current, true)
		}
	}

	destMapID := uint32(m.Wz.ReturnMapId)
	fieldAnchor := types.Point[int16]{X: ch.Position.X, Y: ch.Position.Y}
	var destPortalID uint8
	var closestPortalID uint8
	if ch.Context != nil {
		if destMapWz, ok := ch.Context.GetResources().Maps[destMapID]; ok {
			if id, ok := destMapWz.DoorReturnPortalSpawnID(0); ok {
				destPortalID = id
			} else if _, ok := destMapWz.GetSpawnPosition(0); ok {
				destPortalID = 0
			}
		}
		if id, ok := m.Wz.FindClosestDoorReturnPortalSpawnID(fieldAnchor); ok {
			closestPortalID = id
		} else {
			closestPortalID = m.Wz.FindClosestPortalSpawnID(fieldAnchor)
		}
	}

	door := &Door{
		ObjectCore: ObjectCore{
			Position: ch.Position,
			Context:  m.context,
			Map:      nil,
		},
		OwnerID:        ch.GetID(),
		SkillID:        skillID,
		ReturnMapID:    destMapID,
		FieldMapID:     uint32(m.Wz.ID),
		ReturnPortalID: destPortalID,
		FieldPortalID:  closestPortalID,
	}
	door.ObjectCore.self = door
	if ch.doors == nil {
		ch.doors = make(map[constant.SkillID]*Door)
	}
	ch.doors[skillID] = door
	m.AddDoor(door)
	if skillID == constant.SkillMysticDoor && destMapID != 0 && destMapID != uint32(m.Wz.ID) && ch.Context != nil {
		ch.Context.NotifyDoorSpawn(destMapID, DoorSpawn{
			OwnerID:        ch.GetID(),
			SkillID:        skillID,
			FieldMapID:     uint32(m.Wz.ID),
			ReturnPortalID: destPortalID,
			FieldPortalID:  closestPortalID,
		})
	}
	return door
}

func (ch *Character) forgetDoorRegistrationIfSame(door *Door) {
	if ch == nil || door == nil {
		return
	}
	if ch.doors != nil {
		key := constant.SkillID(door.SkillID)
		if d := ch.doors[key]; d == door {
			delete(ch.doors, key)
		}
	}
}

func (ch *Character) RemoveMist(mist *Mist) {
	if mist == nil {
		return
	}
	if m := mist.GetMap(); m != nil && mist.OID != 0 {
		m.RemoveMist(mist.OID)
	}
}

func (ch *Character) RemoveSummon(target *Summon, animated bool) {
	if target == nil {
		return
	}
	_ = ch.RemoveTimer(summonTimerKey(target.SkillID))
	if m := target.GetMap(); m != nil && target.OID != 0 {
		m.RemoveSummon(target.OID, animated)
	}
	if ch.summons != nil {
		delete(ch.summons, constant.SkillID(target.SkillID))
	}
}

func (ch *Character) RemoveDoor(target *Door, animated bool) {
	if target == nil {
		return
	}
	if m := target.GetMap(); m != nil && target.OID != 0 {
		m.RemoveDoor(target.OID, animated)
	}
	if ch.doors != nil {
		delete(ch.doors, constant.SkillID(target.SkillID))
	}
}

func (ch *Character) RemoveDoorBySkill(skillID constant.SkillID, animated bool) {
	if ch.doors == nil {
		return
	}
	door := ch.doors[skillID]
	if door == nil {
		return
	}
	ch.RemoveDoor(door, animated)
}

func (ch *Character) GetSummons() []*Summon {
	if len(ch.summons) == 0 {
		return nil
	}
	out := make([]*Summon, 0, len(ch.summons))
	for _, s := range ch.summons {
		if s != nil {
			out = append(out, s)
		}
	}
	return out
}

func (ch *Character) GetSummon(skillId constant.SkillID) *Summon {
	return ch.summons[skillId]
}

func (ch *Character) GetSummonsSize() int {
	return len(ch.summons)
}

func (ch *Character) ClearSummons() {
	if len(ch.summons) == 0 {
		return
	}
	for _, s := range ch.GetSummons() {
		ch.RemoveSummon(s, true)
	}
}

func (ch *Character) GetDoors() []*Door {
	if len(ch.doors) == 0 {
		return nil
	}
	out := make([]*Door, 0, len(ch.doors))
	for _, d := range ch.doors {
		if d != nil {
			out = append(out, d)
		}
	}
	return out
}

func (ch *Character) ClearDoors() {
	if len(ch.doors) == 0 {
		return
	}
	for _, d := range ch.GetDoors() {
		ch.RemoveDoor(d, true)
	}
}

func (ch *Character) handleSummonExpireBySkill(skillID constant.SkillID) {
	if ch.summons == nil {
		return
	}
	s := ch.summons[skillID]
	if s == nil {
		return
	}
	ch.RemoveSummon(s, true)
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
			ch.Context.DispatchRunCharacterTimer(pid, &c_actor.RunCharacterTimer{CharacterID: characterID, Key: k})
		})
	}
}

func (ch *Character) GetHp() uint32        { return ch.Hp }
func (ch *Character) GetMp() uint32        { return ch.Mp }
func (ch *Character) GetBonusHp() int32    { return ch.BonusHp }
func (ch *Character) GetBonusMp() int32    { return ch.BonusMp }
func (ch *Character) GetInvincible() bool  { return ch.Invincible }
func (ch *Character) IsAlive() bool        { return ch.Hp > 0 }
func (ch *Character) GetGender() uint8     { return ch.gender }
func (ch *Character) GetSkinColor() uint8  { return ch.skinColor }
func (ch *Character) GetFace() uint32      { return ch.face }
func (ch *Character) GetHair() uint32      { return ch.hair }
func (ch *Character) GetLevel() uint8      { return ch.level }
func (ch *Character) GetExp() uint32       { return ch.exp }
func (ch *Character) GetSpawnPoint() uint8 { return ch.spawnPoint }

func (ch *Character) SetHp(v uint32, notify bool) {
	maxHp := ch.GetMaxHp()
	if v > maxHp {
		v = maxHp
	}
	ch.Hp = v
	if notify {
		ch.Listener.OnUpdateStats(ch, map[constant.Stat]int32{constant.STAT_HP: int32(ch.Hp)}, false)
	}
}

func (ch *Character) SetMp(v uint32, notify bool) {
	maxMp := ch.GetMaxMp()
	if v > maxMp {
		v = maxMp
	}
	ch.Mp = v
	if notify {
		ch.Listener.OnUpdateStats(ch, map[constant.Stat]int32{constant.STAT_MP: int32(ch.Mp)}, false)
	}
}

func (ch *Character) SetBonusHp(v int32) {
	ch.BonusHp = v
	if ch.Hp > ch.GetMaxHp() {
		ch.Hp = ch.GetMaxHp()
	}
	ch.notifyStatChange(constant.STAT_MAX_HP)
}

func (ch *Character) SetBonusMp(v int32) {
	ch.BonusMp = v
	if ch.Mp > ch.GetMaxMp() {
		ch.Mp = ch.GetMaxMp()
	}
	ch.notifyStatChange(constant.STAT_MAX_MP)
}

func (ch *Character) SetInvincible(b bool) { ch.Invincible = b }

func (ch *Character) AddHp(amount int) {
	n := int(ch.Hp) + amount
	if n < 0 {
		n = 0
	}
	maxHp := int(ch.GetMaxHp())
	if n > maxHp {
		n = maxHp
	}
	ch.Hp = uint32(n)
	ch.Listener.OnUpdateStats(ch, map[constant.Stat]int32{constant.STAT_HP: int32(ch.Hp)}, false)
}

func (ch *Character) AddMp(amount int) {
	n := int(ch.Mp) + amount
	if n < 0 {
		n = 0
	}
	maxMp := int(ch.GetMaxMp())
	if n > maxMp {
		n = maxMp
	}
	ch.Mp = uint32(n)
	ch.Listener.OnUpdateStats(ch, map[constant.Stat]int32{constant.STAT_MP: int32(ch.Mp)}, false)
}

func (ch *Character) AddHpMp(hpDelta, mpDelta int) {
	nh := int(ch.Hp) + hpDelta
	if nh < 0 {
		nh = 0
	}
	maxHp := int(ch.GetMaxHp())
	if nh > maxHp {
		nh = maxHp
	}
	ch.Hp = uint32(nh)
	nm := int(ch.Mp) + mpDelta
	if nm < 0 {
		nm = 0
	}
	maxMp := int(ch.GetMaxMp())
	if nm > maxMp {
		nm = maxMp
	}
	ch.Mp = uint32(nm)
	ch.Listener.OnUpdateStats(ch, map[constant.Stat]int32{
		constant.STAT_HP: int32(ch.Hp),
		constant.STAT_MP: int32(ch.Mp),
	}, false)
}

func (ch *Character) SetBaseHp(v uint32, notify bool) {
	if v > constant.STAT_MAX_HP_MP {
		v = constant.STAT_MAX_HP_MP
	}
	ch.BaseHp = v
	if ch.Hp > ch.GetMaxHp() {
		ch.Hp = ch.GetMaxHp()
	}
	if notify {
		ch.Listener.OnUpdateStats(ch, map[constant.Stat]int32{
			constant.STAT_HP:     int32(ch.Hp),
			constant.STAT_MAX_HP: int32(ch.GetMaxHp()),
		}, false)
	}
}

func (ch *Character) SetBaseMp(v uint32, notify bool) {
	if v > constant.STAT_MAX_HP_MP {
		v = constant.STAT_MAX_HP_MP
	}
	ch.BaseMp = v
	if ch.Mp > ch.GetMaxMp() {
		ch.Mp = ch.GetMaxMp()
	}
	if notify {
		ch.Listener.OnUpdateStats(ch, map[constant.Stat]int32{
			constant.STAT_MP:     int32(ch.Mp),
			constant.STAT_MAX_MP: int32(ch.GetMaxMp()),
		}, false)
	}
}

func (ch *Character) SetAbilityPoint(v uint16, notify bool) {
	ch.AbilityPoint = v
	if notify {
		ch.Listener.OnUpdateStats(ch, map[constant.Stat]int32{
			constant.STAT_AVAILABLE_AP: int32(ch.AbilityPoint),
		}, false)
	}
}

func (ch *Character) SetSkillPoint(v uint16, notify bool) {
	ch.SkillPoint = v
	if notify {
		ch.Listener.OnUpdateStats(ch, map[constant.Stat]int32{
			constant.STAT_AVAILABLE_SP: int32(ch.SkillPoint),
		}, false)
	}
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
	ch.Listener.OnHiddenChanged(ch, hidden)
}

func (ch *Character) GetID() uint32 {
	return ch.id
}

func (ch *Character) GetRole() constant.CharacterRole {
	return ch.Role
}

func (ch *Character) GetName() string {
	return ch.name
}

func (ch *Character) Message(message string) {
	ch.Listener.OnMessage(ch, constant.MSG_LIGHT_BLUE_TEXT, message)
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
	ch.Listener.OnExpGain(ch, exp)

	if !ch.tryLevelUp() {
		ch.Listener.OnUpdateStats(ch, map[constant.Stat]int32{
			constant.STAT_EXP: int32(ch.exp),
		}, false)
	}
}

type CharacterInitData struct {
	ID           uint32
	AccountID    uint32
	Name         string
	Gender       uint8
	SkinColor    uint8
	Face         uint32
	Hair         uint32
	Level        uint8
	Class        uint16
	Role         uint8
	Str          uint16
	Dex          uint16
	Int          uint16
	Luk          uint16
	Hp           uint32
	MaxHp        uint32
	Mp           uint32
	MaxMp        uint32
	AbilityPoint uint16
	SkillPoint   uint16
	Exp          uint32
	Meso         int32
	SpawnPoint   uint8
	PositionX    int16
	PositionY    int16
	Stance       uint8
}

func NewCharacter(sender Sendable, listener CharacterListener, data *CharacterInitData, ctx GameContext) *Character {
	if listener == nil {
		panic("NewCharacter: listener must not be nil")
	}
	ch := &Character{
		Sendable: sender,
		Listener: listener,
		LifeCore: LifeCore{
			ObjectCore: ObjectCore{
				Context: ctx,
				Position: types.Vector2[int16]{
					X: data.PositionX,
					Y: data.PositionY,
				},
			},
			Hp:     data.Hp,
			Mp:     data.Mp,
			BaseHp: data.MaxHp,
			BaseMp: data.MaxMp,
			Stance: data.Stance,
		},
		id:           data.ID,
		AccountID:    data.AccountID,
		name:         data.Name,
		gender:       data.Gender,
		skinColor:    data.SkinColor,
		face:         data.Face,
		hair:         data.Hair,
		level:        data.Level,
		Class:        data.Class,
		Role:         constant.CharacterRole(data.Role),
		BaseStats:    BaseStats{Str: data.Str, Dex: data.Dex, Int: data.Int, Luk: data.Luk},
		AbilityPoint: data.AbilityPoint,
		SkillPoint:   data.SkillPoint,
		exp:          data.Exp,
		spawnPoint:   data.SpawnPoint,
		Meso:         data.Meso,

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
		Equipments: map[constant.EquipmentPartsType]Equipment{},

		regRocks: []uint32{999999999, 999999999, 999999999, 999999999, 999999999},
		rocks:    []uint32{999999999, 999999999, 999999999, 999999999, 999999999, 999999999, 999999999, 999999999, 999999999, 999999999},
	}
	ch.Buffs = NewBuffContainer(ch)
	ch.Skills = NewSkillContainer(ch)
	ch.LifeCore.ObjectCore.self = ch
	return ch
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
		ch.Listener.OnUpdateStats(ch, stats, false)
		for i := 0; i < levelDiff; i++ {
			ch.broadcastLevelUpEffect()
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

	ch.Listener.OnUpdateStats(ch, map[constant.Stat]int32{
		constant.STAT_LEVEL: int32(ch.level),
		constant.STAT_EXP:   int32(ch.exp),
	}, false)
}

func (ch *Character) broadcastLevelUpEffect() {
	mapInstance := ch.GetMap()
	if mapInstance == nil {
		return
	}

	ch.Broadcast(&response.ShowForeignEffect{
		CharacterID: ch.id,
		EffectID:    0,
	}, nil)
}

func debuffTimerKey(flag constant.DebuffFlag) string {
	return fmt.Sprintf("debuff_%d_%d", flag.Position, flag.Mask)
}

func (ch *Character) HasDebuff(flag constant.DebuffFlag) bool {
	if ch.diseases == nil {
		return false
	}
	_, ok := ch.diseases[flag]
	return ok
}

func (ch *Character) AddDebuff(holder *DiseaseValueHolder) {
	if holder == nil {
		return
	}
	if ch.diseases == nil {
		ch.diseases = make(map[constant.DebuffFlag]*DiseaseValueHolder)
	}
	ch.RemoveTimer(debuffTimerKey(holder.Disease))
	ch.diseases[holder.Disease] = holder
	if holder.Duration > 0 {
		flag := holder.Disease
		ch.addTimer(debuffTimerKey(flag), holder.Duration, false, func() {
			ch.RemoveTimer(debuffTimerKey(flag))
			if _, ok := ch.diseases[flag]; ok {
				delete(ch.diseases, flag)
				ch.Listener.OnDebuffRemoved(ch, []constant.DebuffFlag{flag})
			}
		})
	}
}

func (ch *Character) GiveDebuff(flag constant.DebuffFlag, duration time.Duration, x int16, skillID uint16, skillLevel uint16) {
	if skillID == 0 {
		skillID = flag.DiseaseSkillID
	}
	if skillLevel == 0 {
		skillLevel = 1
	}
	holder := &DiseaseValueHolder{
		Disease:   flag,
		StartTime: time.Now(),
		Duration:  duration,
	}
	ch.AddDebuff(holder)
	ch.Listener.OnDebuffAdded(ch, flag, x, skillID, skillLevel, int32(duration.Milliseconds()))
}

func (ch *Character) RemoveDebuff(flags ...constant.DebuffFlag) {
	var removed []constant.DebuffFlag
	if ch.diseases != nil {
		for _, flag := range flags {
			ch.RemoveTimer(debuffTimerKey(flag))
			if _, ok := ch.diseases[flag]; ok {
				delete(ch.diseases, flag)
				removed = append(removed, flag)
			}
		}
	}
	if len(removed) > 0 {
		ch.Listener.OnDebuffRemoved(ch, removed)
	}
}

func (ch *Character) GetDiseaseMask() [4]uint32 {
	var mask [4]uint32
	if ch.diseases == nil {
		return mask
	}
	for flag := range ch.diseases {
		idx := flag.Position - 1
		if idx >= 0 && idx < constant.MaxBuffFlag {
			mask[idx] |= flag.Mask
		}
	}
	return mask
}

type SpawnPlayerBuffData struct {
	BuffStates        [4]uint32
	SpeedBuff         uint8
	ComboCount        uint8
	WKChargeSkillID   uint32
	MorphID           uint16
	SpiritClawSkillID uint32
	MountLevel        uint32
	MountExp          uint32
	MountFatigue      uint32
}

func (ch *Character) GetRiddingInfo() (mountID int32, active bool) {
	if ch == nil || ch.Buffs == nil {
		return 0, false
	}
	_, mountID, active = ch.Buffs.GetBuffValue(constant.BuffFlagMonsterRiding)
	return
}

func (ch *Character) GetSpawnPlayerBuffData() SpawnPlayerBuffData {
	data := SpawnPlayerBuffData{
		SpeedBuff:  1,
		ComboCount: 1,
		MountLevel: 1,
	}
	if ch == nil || ch.Buffs == nil {
		return data
	}

	for _, buff := range ch.Buffs.Entities() {
		if buff == nil {
			continue
		}
		values := buff.GetValues()
		if values == nil {
			continue
		}
		for flag, value := range values {
			idx := flag.Position - 1
			if idx < 0 || idx >= constant.MaxBuffFlag {
				continue
			}
			if constant.IsRemoteStatFlag(flag) {
				data.BuffStates[idx] |= flag.Mask
			}
			switch flag {
			case constant.BuffFlagSpeed:
				data.SpeedBuff = uint8(value)
			case constant.BuffFlagCombo:
				data.ComboCount = uint8(value)
			case constant.BuffFlagMorph:
				data.MorphID = uint16(value)
			}
		}
		switch typed := buff.(type) {
		case *SkillBuff:
			for _, flag := range typed.GetFlags() {
				switch flag {
				case constant.BuffFlagWkCharge:
					data.WKChargeSkillID = typed.Wz.ID
				case constant.BuffFlagSpiritClaw:
					data.SpiritClawSkillID = typed.Wz.ID
				}
			}
		}
	}

	return data
}
