package entity

import (
	"fmt"
	"log"
	"time"

	"github.com/boyism80/fm/core/clock"
	"github.com/boyism80/fm/core/luax"
	"github.com/boyism80/fm/services/game/constant"
	"github.com/boyism80/fm/services/game/wz"
)

type Buff interface {
	RemainingDuration(now time.Time) time.Duration
	GetFlags() []constant.BuffFlag
	GetValues() map[constant.BuffFlag]int32
	hookScript(ch *Character) (string, interface{}, bool)
	GetBuffID() int32
	expiresAt() (time.Time, bool)
}

type BaseBuff struct {
	StartTime time.Time
	Duration  time.Duration
	Flags     []constant.BuffFlag
	Values    map[constant.BuffFlag]int32
}

type ItemBuff struct {
	*BaseBuff
	Wz *wz.Consume
}

type SkillBuff struct {
	*BaseBuff
	Wz         *wz.Skill
	SkillLevel uint8
	CauserID   uint32
}

func (e *BaseBuff) RemainingDuration(now time.Time) time.Duration {
	if e == nil || e.Duration <= 0 {
		return 0
	}
	end := e.StartTime.Add(e.Duration)
	rem := end.Sub(now)
	if rem < 0 {
		return 0
	}
	return rem
}

func (e *BaseBuff) expiresAt() (time.Time, bool) {
	if e.Duration <= 0 {
		return time.Time{}, false
	}
	return e.StartTime.Add(e.Duration), true
}

func (e *BaseBuff) GetFlags() []constant.BuffFlag {
	if e == nil {
		return nil
	}
	return e.Flags
}

func (e *BaseBuff) GetValues() map[constant.BuffFlag]int32 {
	if e == nil {
		return nil
	}
	return e.Values
}

func (e *SkillBuff) GetBuffID() int32 {
	if e == nil {
		return 0
	}
	return int32(e.Wz.ID)
}

func (e *SkillBuff) hookScript(ch *Character) (string, interface{}, bool) {
	return fmt.Sprintf("script/skill/%d.lua", e.Wz.ID), e, true
}

func (e *ItemBuff) GetBuffID() int32 {
	if e == nil {
		return 0
	}
	return -int32(e.Wz.ID)
}

func (e *ItemBuff) hookScript(ch *Character) (string, interface{}, bool) {
	item, err := NewItem(e.Wz.ID, 1, ch.GameWorld)
	if err != nil {
		log.Printf("Failed to create temp item %d: %v", e.Wz.ID, err)
		return "", nil, false
	}
	consume, ok := item.(*Consume)
	if ok == false {
		return "", nil, false
	}
	return fmt.Sprintf("script/item/%d.lua", e.Wz.ID), consume, true
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
	if owner.Listener == nil {
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
	if ok == false {
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
	ch.Listener.OnBuffAdded(ch, entity.GetBuffID(), entity.RemainingDuration(now), entityValues)
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

	for existing := range conflicts {
		removed = append(removed, existing)
		bc.owner.RemoveTimer(bc.timerKey(existing))
		delete(bc.entities, existing)
		for _, flag := range existing.GetFlags() {
			if current := bc.byFlag[flag]; current == existing {
				delete(bc.byFlag, flag)
			}
		}
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
	if ok == false {
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
	if _, ok := bc.entities[entity]; ok == false {
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
		BaseBuff: &BaseBuff{
			StartTime: now,
			Duration:  duration,
			Flags:     flags,
			Values:    valCopy,
		},
		Wz:         wz,
		SkillLevel: skillLevel,
		CauserID:   causerID,
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
		mul := bc.owner.Stats.PotionDurationMultiplierPercent()
		scaledDuration = time.Duration(int64(duration) * int64(mul) / 100)
	}
	entity := &ItemBuff{
		BaseBuff: &BaseBuff{
			StartTime: now,
			Duration:  scaledDuration,
			Flags:     flags,
			Values:    valCopy,
		},
		Wz: consumeWz,
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
		bc.owner.Listener.OnBuffRemoved(bc.owner, removedFlags)
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
	bc.owner.Listener.OnBuffAdded(bc.owner, entity.GetBuffID(), entity.RemainingDuration(clock.Now()), map[constant.BuffFlag]int32{flag: value})
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
