package game

import (
	"fmt"
	"log"
	"slices"
	"time"

	"github.com/boyism80/fm/core/clock"
	"github.com/boyism80/fm/core/luax"
	"github.com/boyism80/fm/maple/constant"
	"github.com/boyism80/fm/maple/wz"
)

type Buff interface {
	RemainingDuration(now time.Time) time.Duration
	GetFlags() []constant.BuffFlag
	GetValues() map[constant.BuffFlag]int32
	hookScript(ch *Character) (string, interface{}, bool)
	GetBuffID() int32
	expiresAt() (time.Time, bool)
	dropFlag(flag constant.BuffFlag)
}

type BaseBuff struct {
	startTime time.Time
	duration  time.Duration
	flags     []constant.BuffFlag
	values    map[constant.BuffFlag]int32
}

type baseBuff = BaseBuff

type ItemBuff struct {
	*baseBuff
	wz *wz.Consume
}

type SkillBuff struct {
	*baseBuff
	wz         *wz.Skill
	skillLevel uint8
	causerID   uint32
}

func (e *BaseBuff) RemainingDuration(now time.Time) time.Duration {
	if e == nil || e.duration <= 0 {
		return 0
	}
	end := e.startTime.Add(e.duration)
	rem := end.Sub(now)
	if rem < 0 {
		return 0
	}
	return rem
}

func (e *BaseBuff) expiresAt() (time.Time, bool) {
	if e.duration <= 0 {
		return time.Time{}, false
	}
	return e.startTime.Add(e.duration), true
}

func (e *BaseBuff) dropFlag(flag constant.BuffFlag) {
	e.flags = slices.DeleteFunc(e.flags, func(f constant.BuffFlag) bool {
		return f == flag
	})
	delete(e.values, flag)
}

func (e *BaseBuff) GetFlags() []constant.BuffFlag {
	if e == nil {
		return nil
	}
	return e.flags
}

func (e *BaseBuff) GetValues() map[constant.BuffFlag]int32 {
	if e == nil {
		return nil
	}
	return e.values
}

func (e *SkillBuff) GetBuffID() int32 {
	if e == nil {
		return 0
	}
	return int32(e.wz.ID)
}

func (e *SkillBuff) hookScript(ch *Character) (string, interface{}, bool) {
	return fmt.Sprintf("script/skill/%d.lua", e.wz.ID), e, true
}

func (e *ItemBuff) GetBuffID() int32 {
	if e == nil {
		return 0
	}
	return -int32(e.wz.ID)
}

func (e *ItemBuff) hookScript(ch *Character) (string, interface{}, bool) {
	item, err := NewItem(e.wz.ID, 1, ch.gameWorld)
	if err != nil {
		log.Printf("Failed to create temp item %d: %v", e.wz.ID, err)
		return "", nil, false
	}
	consume, ok := item.(*Consume)
	if !ok {
		return "", nil, false
	}
	return fmt.Sprintf("script/item/%d.lua", e.wz.ID), consume, true
}

type Buffs struct {
	owner    *Character
	byFlag   map[constant.BuffFlag]Buff
	entities map[Buff]struct{}
}

func NewBuffs(owner *Character) *Buffs {
	if owner == nil {
		panic("Buffs owner is nil")
	}
	if owner.listener == nil {
		panic("Buffs: character listener must not be nil")
	}
	return &Buffs{
		owner:    owner,
		byFlag:   make(map[constant.BuffFlag]Buff),
		entities: make(map[Buff]struct{}),
	}
}

func copyBuffValues(values map[constant.BuffFlag]int32) ([]constant.BuffFlag, map[constant.BuffFlag]int32) {
	flags := make([]constant.BuffFlag, 0, len(values))
	valCopy := make(map[constant.BuffFlag]int32, len(values))
	for flag, value := range values {
		flags = append(flags, flag)
		valCopy[flag] = value
	}
	return flags, valCopy
}

func (bc *Buffs) callHook(entity Buff, hook string) {
	mapInstance := bc.owner.GetMap()
	if mapInstance == nil {
		return
	}
	root := mapInstance.GetLuaRoot()
	if root == nil {
		return
	}
	path, arg, ok := entity.hookScript(bc.owner)
	if !ok {
		return
	}
	thread, err := luax.NewThread(root, path)
	if err != nil {
		log.Printf("Failed to call %s for %s: %v", hook, path, err)
		return
	}
	luax.CallAsync(nil, root, thread, hook, bc.owner, arg).OnError(func(err error) {
		log.Printf("Failed to call %s for %s: %v", hook, path, err)
	})
}

func (bc *Buffs) callUnbuffScripts(removed []Buff) {
	for _, entity := range removed {
		if entity == nil {
			continue
		}
		bc.callHook(entity, "on_unbuff")
	}
}

func (bc *Buffs) restoreEffects() {
	for entity := range bc.entities {
		bc.callHook(entity, "on_buff")
	}
}

func (bc *Buffs) updateMaxHpMp(flags []constant.BuffFlag) {
	for _, flag := range flags {
		if flag == constant.BuffFlagMaxHp || flag == constant.BuffFlagMaxMp {
			bc.owner.updateMaxHpMp()
			return
		}
	}
}

func (bc *Buffs) notifyAdded(entity Buff, now time.Time) {
	ch := bc.owner
	values := entity.GetValues()
	entityValues := make(map[constant.BuffFlag]int32, len(values))
	for flag, value := range values {
		entityValues[flag] = value
	}
	ch.listener.OnBuffAdded(ch, entity.GetBuffID(), entity.RemainingDuration(now), entityValues)
}

func (bc *Buffs) add(entity Buff) (removed []Buff) {
	if entity == nil || len(entity.GetFlags()) == 0 {
		return nil
	}

	conflicts := make(map[Buff]struct{})
	for _, flag := range entity.GetFlags() {
		if existing := bc.byFlag[flag]; existing != nil && existing != entity {
			conflicts[existing] = struct{}{}
		}
	}

	// A buff that keeps some flags keeps its timer; its expiry cancels only the flags it still holds.
	for existing := range conflicts {
		for _, flag := range entity.GetFlags() {
			existing.dropFlag(flag)
		}
		if len(existing.GetFlags()) > 0 {
			continue
		}
		removed = append(removed, existing)
		bc.owner.RemoveTimer(bc.timerKey(existing))
		delete(bc.entities, existing)
	}

	bc.entities[entity] = struct{}{}
	for _, flag := range entity.GetFlags() {
		bc.byFlag[flag] = entity
	}
	return removed
}

func (bc *Buffs) remove(flags []constant.BuffFlag) (removed []Buff, removedFlags []constant.BuffFlag) {
	if len(flags) == 0 {
		return nil, nil
	}

	targets := make(map[Buff]struct{})
	for _, flag := range flags {
		if existing := bc.byFlag[flag]; existing != nil {
			targets[existing] = struct{}{}
		}
	}

	removedFlagSet := make(map[constant.BuffFlag]struct{})
	for entity := range targets {
		removed = append(removed, entity)
		for _, flag := range entity.GetFlags() {
			removedFlagSet[flag] = struct{}{}
		}
		bc.owner.RemoveTimer(bc.timerKey(entity))
		delete(bc.entities, entity)
		for _, flag := range entity.GetFlags() {
			if current := bc.byFlag[flag]; current == entity {
				delete(bc.byFlag, flag)
			}
		}
	}

	removedFlags = make([]constant.BuffFlag, 0, len(removedFlagSet))
	for flag := range removedFlagSet {
		removedFlags = append(removedFlags, flag)
	}

	return removed, removedFlags
}

func (bc *Buffs) commitBuff(entity Buff, now time.Time, notify bool) {
	removed := bc.add(entity)
	bc.scheduleExpire(entity, now)
	bc.callUnbuffScripts(removed)
	bc.callHook(entity, "on_buff")
	if notify {
		bc.notifyAdded(entity, now)
		flags := entity.GetFlags()
		for _, existing := range removed {
			flags = append(flags, existing.GetFlags()...)
		}
		bc.updateMaxHpMp(flags)
	}
}

func (bc *Buffs) scheduleExpire(entity Buff, now time.Time) {
	end, ok := entity.expiresAt()
	if !ok {
		return
	}
	delay := end.Sub(now)
	if delay < 0 {
		delay = 0
	}
	bc.owner.AddTimer(bc.timerKey(entity), delay, false, func() {
		bc.expire(entity)
	})
}

func (bc *Buffs) scheduleExpires() {
	now := clock.Now()
	for entity := range bc.entities {
		if bc.owner.GetTimerEntry(bc.timerKey(entity)) != nil {
			continue
		}
		bc.scheduleExpire(entity, now)
	}
}

func (bc *Buffs) expire(entity Buff) {
	if _, ok := bc.entities[entity]; !ok {
		return
	}
	bc.RemoveBuff(entity.GetFlags())
}

func (bc *Buffs) timerKey(entity Buff) string {
	return fmt.Sprintf("buff:%p", entity)
}

func (bc *Buffs) Has(flag constant.BuffFlag) bool {
	_, exists := bc.byFlag[flag]
	return exists
}

func (bc *Buffs) GetEntity(flag constant.BuffFlag) Buff {
	return bc.byFlag[flag]
}

func (bc *Buffs) Entities() []Buff {
	entities := make([]Buff, 0, len(bc.entities))
	for entity := range bc.entities {
		entities = append(entities, entity)
	}
	return entities
}

func (bc *Buffs) AddBuff(wz *wz.Skill, duration time.Duration, skillLevel uint8, causerID uint32, values map[constant.BuffFlag]int32, notify bool) {
	if len(values) == 0 {
		return
	}
	now := clock.Now()
	flags, valCopy := copyBuffValues(values)
	entity := &SkillBuff{
		baseBuff: &BaseBuff{
			startTime: now,
			duration:  duration,
			flags:     flags,
			values:    valCopy,
		},
		wz:         wz,
		skillLevel: skillLevel,
		causerID:   causerID,
	}
	bc.commitBuff(entity, now, notify)
}

func (bc *Buffs) AddItemBuff(consumeWz *wz.Consume, duration time.Duration, values map[constant.BuffFlag]int32, applyPotionDurationScale bool, notify bool) {
	if len(values) == 0 || consumeWz == nil {
		return
	}
	now := clock.Now()
	flags, valCopy := copyBuffValues(values)
	scaledDuration := duration
	if applyPotionDurationScale && duration > 0 {
		mul := bc.owner.stats.PotionDurationMultiplierPercent()
		scaledDuration = time.Duration(int64(duration) * int64(mul) / 100)
	}
	entity := &ItemBuff{
		baseBuff: &BaseBuff{
			startTime: now,
			duration:  scaledDuration,
			flags:     flags,
			values:    valCopy,
		},
		wz: consumeWz,
	}
	bc.commitBuff(entity, now, notify)
}

func (bc *Buffs) Dispel() {
	entities := bc.Entities()
	if len(entities) == 0 {
		return
	}
	flags := make([]constant.BuffFlag, 0)
	for _, entity := range entities {
		if entity == nil {
			continue
		}
		flags = append(flags, entity.GetFlags()...)
	}
	bc.RemoveBuff(flags)
}

func (bc *Buffs) RemoveBuff(flags []constant.BuffFlag) {
	if len(flags) == 0 {
		return
	}
	removed, removedFlags := bc.remove(flags)
	if len(removed) == 0 {
		return
	}
	bc.callUnbuffScripts(removed)
	if len(removedFlags) > 0 {
		bc.owner.listener.OnBuffRemoved(bc.owner, removedFlags)
	}
	bc.updateMaxHpMp(removedFlags)
}

func (bc *Buffs) CancelBySource(buffID int32) {
	var flags []constant.BuffFlag
	for entity := range bc.entities {
		if entity.GetBuffID() == buffID {
			flags = append(flags, entity.GetFlags()...)
		}
	}
	bc.RemoveBuff(flags)
}

func (bc *Buffs) GetBuffValue(flag constant.BuffFlag) (Buff, int32, bool) {
	entity := bc.GetEntity(flag)
	if entity == nil {
		return nil, 0, false
	}
	values := entity.GetValues()
	if values == nil {
		return nil, 0, false
	}
	value, exists := values[flag]
	if !exists {
		return nil, 0, false
	}
	return entity, value, true
}

func (bc *Buffs) SetBuffValue(flag constant.BuffFlag, value int32) (Buff, bool) {
	entity := bc.GetEntity(flag)
	if entity == nil {
		return nil, false
	}
	values := entity.GetValues()
	if values == nil {
		return nil, false
	}
	values[flag] = value
	bc.owner.listener.OnBuffAdded(bc.owner, entity.GetBuffID(), entity.RemainingDuration(clock.Now()), map[constant.BuffFlag]int32{flag: value})
	return entity, true
}

func (bc *Buffs) ShowAll() {
	now := clock.Now()
	for entity := range bc.entities {
		if entity == nil {
			continue
		}
		values := entity.GetValues()
		if len(values) == 0 {
			continue
		}
		bc.notifyAdded(entity, now)
	}
}

func (bc *Buffs) value(flag constant.BuffFlag) float64 {
	_, v, ok := bc.GetBuffValue(flag)
	if !ok {
		return 0
	}
	return float64(v)
}

type SpawnPlayerBuffData struct {
	BuffStates        [4]uint32
	SpeedBuff         uint8
	ComboCount        uint8
	WKChargeSkillID   uint32
	MorphID           uint16
	SpiritClawSkillID uint32
}

func (bc *Buffs) SpawnPlayerData() SpawnPlayerBuffData {
	data := SpawnPlayerBuffData{
		SpeedBuff:  1,
		ComboCount: 1,
	}

	for _, buff := range bc.Entities() {
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
			if constant.TargetStatFlags[flag] {
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
					data.WKChargeSkillID = typed.wz.ID
				case constant.BuffFlagSpiritClaw:
					data.SpiritClawSkillID = typed.wz.ID
				}
			}
		}
	}

	return data
}

func (bc *Buffs) Cancel(sourceID int32) {
	if sourceID > 0 {
		skill := bc.owner.skills.Learned(uint32(sourceID))
		if skill != nil && skill.wz.Charge {
			bc.owner.listener.OnChargeSkillCancelled(bc.owner, uint32(sourceID))
			return
		}
	}
	bc.CancelBySource(sourceID)
}

func (bc *Buffs) addConsumeBuff(consumeItem *wz.Consume) bool {
	if consumeItem == nil {
		return false
	}
	buffValues := consumeItem.BuffSpecValues()
	if len(buffValues) == 0 {
		return false
	}
	if m := bc.owner.GetMap(); m != nil && !consumeItem.IsAvailableAt(m.TemplateID()) {
		return false
	}
	bc.AddItemBuff(consumeItem, consumeItem.BuffDuration, buffValues, true, true)
	return true
}

func (bc *Buffs) HolySymbolExpRate() int32 {
	ent := bc.GetEntity(constant.BuffFlagHolySymbol)
	sb, typeOk := ent.(*SkillBuff)
	if !typeOk || sb == nil || sb.wz == nil {
		return 100
	}
	x, valOk := sb.values[constant.BuffFlagHolySymbol]
	if !valOk || x <= 0 {
		return 100
	}
	full := sb.wz.FullHolySymbol
	if !full {
		pid := bc.owner.party.ID()
		m := bc.owner.GetMap()
		full = pid != nil && m != nil && len(m.GetPartyMembers(*pid)) >= 2
	}
	if full {
		return 100 + x
	}
	return int32(100 * int64(150+x) / 150)
}
