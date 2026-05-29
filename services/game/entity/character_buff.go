package entity

import (
	"fmt"
	"log"
	"time"

	"github.com/boyism80/fm/core/luax"
	"github.com/boyism80/fm/services/game/constant"
	"github.com/boyism80/fm/services/game/wz"
)

type Buff interface {
	RemainingDuration(now time.Time) time.Duration
	GetFlags() []constant.BuffFlag
	GetValues() map[constant.BuffFlag]int32
	CallOnBuffScript(ch *Character)
	CallOnUnbuffScript(ch *Character)
	GetBuffID() int32
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

func (e *SkillBuff) CallOnBuffScript(ch *Character) {
	if e == nil || ch == nil || ch.GameWorld == nil {
		return
	}
	mapInstance := ch.GetMap()
	if mapInstance == nil {
		return
	}
	root := mapInstance.GetLuaRoot()
	if root == nil {
		return
	}

	skillID := e.Wz.ID
	scriptPath := fmt.Sprintf("script/skill/character/%d.lua", skillID)
	thread, err := luax.NewThread(root, scriptPath)
	if err != nil {
		log.Printf("Failed to call on_buff for skill %d: %v", skillID, err)
		return
	}
	if _, err := luax.Call(thread, fmt.Sprintf("on_buff_%d", skillID), ch, e); err != nil {
		log.Printf("Failed to call on_buff for skill %d: %v", skillID, err)
	}
}

func (e *SkillBuff) CallOnUnbuffScript(ch *Character) {
	if e == nil || ch == nil || ch.GameWorld == nil {
		return
	}
	mapInstance := ch.GetMap()
	if mapInstance == nil {
		return
	}
	root := mapInstance.GetLuaRoot()
	if root == nil {
		return
	}

	skillID := e.Wz.ID
	scriptPath := fmt.Sprintf("script/skill/character/%d.lua", skillID)
	thread, err := luax.NewThread(root, scriptPath)
	if err != nil {
		log.Printf("Failed to call on_unbuff for skill %d: %v", skillID, err)
		return
	}
	if _, err := luax.Call(thread, fmt.Sprintf("on_unbuff_%d", skillID), ch, e); err != nil {
		log.Printf("Failed to call on_unbuff for skill %d: %v", skillID, err)
	}
}

func (e *ItemBuff) GetBuffID() int32 {
	if e == nil {
		return 0
	}
	return -int32(e.Wz.ID)
}

func (e *ItemBuff) CallOnBuffScript(ch *Character) {
	if e == nil || ch == nil || ch.GameWorld == nil {
		return
	}
	mapInstance := ch.GetMap()
	if mapInstance == nil {
		return
	}
	root := mapInstance.GetLuaRoot()
	if root == nil {
		return
	}

	itemWzID := e.Wz.ID
	scriptPath := fmt.Sprintf("script/item/%d.lua", itemWzID)

	tempItem, err := NewItem(itemWzID, 1, ch.GameWorld)
	if err != nil {
		log.Printf("Failed to create temp item %d: %v", itemWzID, err)
		return
	}
	consume, ok := tempItem.(*Consume)
	if !ok || consume == nil {
		return
	}

	thread, callErr := luax.NewThread(root, scriptPath)
	if callErr != nil {
		log.Printf("Failed to call on_buff for item %d: %v", itemWzID, callErr)
		return
	}
	if _, callErr := luax.Call(thread, "on_buff", ch, consume); callErr != nil {
		log.Printf("Failed to call on_buff for item %d: %v", itemWzID, callErr)
	}
}

func (e *ItemBuff) CallOnUnbuffScript(ch *Character) {
	if e == nil || ch == nil || ch.GameWorld == nil {
		return
	}
	mapInstance := ch.GetMap()
	if mapInstance == nil {
		return
	}
	root := mapInstance.GetLuaRoot()
	if root == nil {
		return
	}

	itemWzID := e.Wz.ID
	scriptPath := fmt.Sprintf("script/item/%d.lua", itemWzID)

	tempItem, err := NewItem(itemWzID, 1, ch.GameWorld)
	if err != nil {
		log.Printf("Failed to create temp item %d: %v", itemWzID, err)
		return
	}
	consume, ok := tempItem.(*Consume)
	if !ok || consume == nil {
		return
	}

	thread, callErr := luax.NewThread(root, scriptPath)
	if callErr != nil {
		log.Printf("Failed to call on_unbuff for item %d: %v", itemWzID, callErr)
		return
	}
	if _, callErr := luax.Call(thread, "on_unbuff", ch, consume); callErr != nil {
		log.Printf("Failed to call on_unbuff for item %d: %v", itemWzID, callErr)
	}
}

type BuffContainer struct {
	owner    *Character
	byFlag   map[constant.BuffFlag]Buff
	entities map[Buff]struct{}
}

func NewBuffContainer(owner *Character) *BuffContainer {
	if owner == nil {
		panic("buff container owner is nil")
	}
	if owner.Listener == nil {
		panic("BuffContainer: character listener must not be nil")
	}
	return &BuffContainer{
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

func (bc *BuffContainer) callUnbuffScripts(removed []Buff) {
	ch := bc.owner
	for _, entity := range removed {
		if entity == nil {
			continue
		}
		entity.CallOnUnbuffScript(ch)
	}
}

func (bc *BuffContainer) notifyAdded(entity Buff, now time.Time) {
	ch := bc.owner
	values := entity.GetValues()
	entityValues := make(map[constant.BuffFlag]int32, len(values))
	for flag, value := range values {
		entityValues[flag] = value
	}
	ch.Listener.OnBuffAdded(ch, entity.GetBuffID(), entity.RemainingDuration(now), entityValues)
}

func (bc *BuffContainer) add(entity Buff) (removed []Buff) {
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

func (bc *BuffContainer) remove(flags []constant.BuffFlag) (removed []Buff, removedFlags []constant.BuffFlag) {
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

func (bc *BuffContainer) applyAdded(entity Buff, now time.Time, notify bool) {
	removed := bc.add(entity)
	bc.callUnbuffScripts(removed)
	entity.CallOnBuffScript(bc.owner)
	if notify {
		bc.notifyAdded(entity, now)
	}
}

func (bc *BuffContainer) Has(flag constant.BuffFlag) bool {
	_, exists := bc.byFlag[flag]
	return exists
}

func (bc *BuffContainer) GetEntity(flag constant.BuffFlag) Buff {
	return bc.byFlag[flag]
}

func (bc *BuffContainer) Entities() []Buff {
	entities := make([]Buff, 0, len(bc.entities))
	for entity := range bc.entities {
		entities = append(entities, entity)
	}
	return entities
}

func (bc *BuffContainer) AddBuff(wz *wz.Skill, duration time.Duration, skillLevel uint8, causerID uint32, values map[constant.BuffFlag]int32, notify bool) {
	if len(values) == 0 {
		return
	}
	now := time.Now()
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
	bc.applyAdded(entity, now, notify)
}

func (bc *BuffContainer) AddItemBuff(consumeWz *wz.Consume, duration time.Duration, values map[constant.BuffFlag]int32, applyPotionDurationScale bool, notify bool) {
	if len(values) == 0 || consumeWz == nil {
		return
	}
	now := time.Now()
	flags, valCopy := copyBuffValues(values)
	scaledDuration := duration
	if applyPotionDurationScale && duration > 0 {
		mul := bc.owner.PotionDurationMultiplierPercent()
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
	bc.applyAdded(entity, now, notify)
}

func (bc *BuffContainer) RemoveBuff(flags []constant.BuffFlag) {
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
}

func (bc *BuffContainer) RemoveSkillBuff(skillID uint32) {
	var target Buff
	for entity := range bc.entities {
		skillBuff, ok := entity.(*SkillBuff)
		if !ok || skillBuff == nil {
			continue
		}
		if skillBuff.Wz.ID == skillID {
			target = entity
			break
		}
	}
	if target == nil {
		return
	}
	flags := target.GetFlags()
	if len(flags) == 0 {
		return
	}
	bc.RemoveBuff(flags)
}

func (bc *BuffContainer) GetBuffValue(flag constant.BuffFlag) (Buff, int32, bool) {
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

func (bc *BuffContainer) SetBuffValue(flag constant.BuffFlag, value int32) (Buff, bool) {
	entity := bc.GetEntity(flag)
	if entity == nil {
		return nil, false
	}
	values := entity.GetValues()
	if values == nil {
		return nil, false
	}
	values[flag] = value
	bc.owner.Listener.OnBuffAdded(bc.owner, entity.GetBuffID(), entity.RemainingDuration(time.Now()), map[constant.BuffFlag]int32{flag: value})
	return entity, true
}

func (bc *BuffContainer) EmitAllBuffAddedEvents() {
	now := time.Now()
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
