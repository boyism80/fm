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
	pid := mapInstance.GetActorPID()
	if pid == nil {
		return
	}
	root := luax.GetRootLuaState(pid.String())
	if root == nil {
		return
	}

	skillID := e.Wz.ID
	scriptPath := fmt.Sprintf("script/skill/%d.lua", skillID)
	_, thread, err := luax.Call(root, scriptPath, luax.SkillScriptHookName("on_buff", skillID), ch, e)
	if err != nil {
		log.Printf("Failed to call on_buff for skill %d: %v", skillID, err)
	}
	if thread != nil {
		thread.Close()
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
	pid := mapInstance.GetActorPID()
	if pid == nil {
		return
	}
	root := luax.GetRootLuaState(pid.String())
	if root == nil {
		return
	}

	skillID := e.Wz.ID
	scriptPath := fmt.Sprintf("script/skill/%d.lua", skillID)
	_, thread, err := luax.Call(root, scriptPath, luax.SkillScriptHookName("on_unbuff", skillID), ch, e)
	if err != nil {
		log.Printf("Failed to call on_unbuff for skill %d: %v", skillID, err)
	}
	if thread != nil {
		thread.Close()
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
	pid := mapInstance.GetActorPID()
	if pid == nil {
		return
	}
	root := luax.GetRootLuaState(pid.String())
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

	_, thread, callErr := luax.Call(root, scriptPath, "on_buff", ch, consume)
	if callErr != nil {
		log.Printf("Failed to call on_buff for item %d: %v", itemWzID, callErr)
	}
	if thread != nil {
		thread.Close()
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
	pid := mapInstance.GetActorPID()
	if pid == nil {
		return
	}
	root := luax.GetRootLuaState(pid.String())
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

	_, thread, callErr := luax.Call(root, scriptPath, "on_unbuff", ch, consume)
	if callErr != nil {
		log.Printf("Failed to call on_unbuff for item %d: %v", itemWzID, callErr)
	}
	if thread != nil {
		thread.Close()
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
	return &BuffContainer{
		owner:    owner,
		byFlag:   make(map[constant.BuffFlag]Buff),
		entities: make(map[Buff]struct{}),
	}
}

func (bc *BuffContainer) add(entity Buff) (removed []Buff) {
	if bc == nil || entity == nil || len(entity.GetFlags()) == 0 {
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
		bc.removeEntity(existing)
	}

	bc.addEntity(entity)
	return removed
}

func (bc *BuffContainer) remove(flags []constant.BuffFlag) (removed []Buff, removedFlags []constant.BuffFlag) {
	if bc == nil || len(flags) == 0 {
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
		bc.removeEntity(entity)
	}

	removedFlags = make([]constant.BuffFlag, 0, len(removedFlagSet))
	for flag := range removedFlagSet {
		removedFlags = append(removedFlags, flag)
	}

	return removed, removedFlags
}

func (bc *BuffContainer) Has(flag constant.BuffFlag) bool {
	if bc == nil {
		return false
	}
	_, exists := bc.byFlag[flag]
	return exists
}

func (bc *BuffContainer) GetEntity(flag constant.BuffFlag) Buff {
	if bc == nil {
		return nil
	}
	return bc.byFlag[flag]
}

func (bc *BuffContainer) Entities() []Buff {
	if bc == nil {
		return nil
	}
	entities := make([]Buff, 0, len(bc.entities))
	for entity := range bc.entities {
		entities = append(entities, entity)
	}
	return entities
}

func (bc *BuffContainer) addEntity(entity Buff) {
	bc.entities[entity] = struct{}{}
	for _, flag := range entity.GetFlags() {
		bc.byFlag[flag] = entity
	}
}

func (bc *BuffContainer) removeEntity(entity Buff) {
	delete(bc.entities, entity)
	for _, flag := range entity.GetFlags() {
		if current := bc.byFlag[flag]; current == entity {
			delete(bc.byFlag, flag)
		}
	}
}

func (bc *BuffContainer) AddBuff(wz *wz.Skill, duration time.Duration, skillLevel uint8, causerID uint32, values map[constant.BuffFlag]int32) {
	if bc == nil {
		return
	}
	if len(values) == 0 {
		return
	}
	ch := bc.owner

	entityFlags := make([]constant.BuffFlag, 0, len(values))
	entityValues := make(map[constant.BuffFlag]int32, len(values))
	for flag, value := range values {
		entityFlags = append(entityFlags, flag)
		entityValues[flag] = value
	}

	now := time.Now()

	entity := &SkillBuff{
		BaseBuff: &BaseBuff{
			StartTime: now,
			Duration:  duration,
			Flags:     entityFlags,
			Values:    entityValues,
		},
		Wz:         wz,
		SkillLevel: skillLevel,
		CauserID:   causerID,
	}

	removed := bc.add(entity)
	if len(removed) > 0 {
		ch.handleRemovedBuffEntities(removed)
	}

	entity.CallOnBuffScript(ch)

	ch.Listener.OnBuffAdded(ch, entity.GetBuffID(), entity.RemainingDuration(now), entityValues)
}

func (bc *BuffContainer) AddItemBuff(consumeWz *wz.Consume, duration time.Duration, values map[constant.BuffFlag]int32) {
	if bc == nil {
		return
	}
	if len(values) == 0 {
		return
	}
	if consumeWz == nil {
		return
	}

	ch := bc.owner

	entityFlags := make([]constant.BuffFlag, 0, len(values))
	entityValues := make(map[constant.BuffFlag]int32, len(values))
	for flag, value := range values {
		entityFlags = append(entityFlags, flag)
		entityValues[flag] = value
	}

	now := time.Now()

	scaledDuration := duration
	if duration > 0 {
		mul := ch.PotionDurationMultiplierPercent()
		scaledDuration = time.Duration(int64(duration) * int64(mul) / 100)
	}

	entity := &ItemBuff{
		BaseBuff: &BaseBuff{
			StartTime: now,
			Duration:  scaledDuration,
			Flags:     entityFlags,
			Values:    entityValues,
		},
		Wz: consumeWz,
	}

	removed := bc.add(entity)
	if len(removed) > 0 {
		ch.handleRemovedBuffEntities(removed)
	}

	entity.CallOnBuffScript(ch)

	ch.Listener.OnBuffAdded(ch, entity.GetBuffID(), entity.RemainingDuration(now), entityValues)
}

func (bc *BuffContainer) RemoveBuff(flags []constant.BuffFlag) {
	if bc == nil {
		return
	}
	if len(flags) == 0 {
		return
	}
	ch := bc.owner
	removed, removedFlags := bc.remove(flags)
	if len(removed) == 0 {
		return
	}
	ch.handleRemovedBuffEntities(removed)
	if len(removedFlags) > 0 {
		ch.Listener.OnBuffRemoved(ch, removedFlags)
	}
}

func (bc *BuffContainer) RemoveSkillBuff(skillID uint32) {
	if bc == nil {
		return
	}
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
	if bc == nil {
		return nil, 0, false
	}
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
	if bc == nil {
		return nil, false
	}
	entity := bc.GetEntity(flag)
	if entity == nil {
		return nil, false
	}
	values := entity.GetValues()
	if values == nil {
		return nil, false
	}
	values[flag] = value
	if bc.owner != nil {
		bc.owner.Listener.OnBuffAdded(bc.owner, entity.GetBuffID(), entity.RemainingDuration(time.Now()), map[constant.BuffFlag]int32{flag: value})
	}
	return entity, true
}

func (ch *Character) handleRemovedBuffEntities(removed []Buff) {
	if len(removed) == 0 {
		return
	}
	for _, entity := range removed {
		if entity == nil {
			continue
		}
		entity.CallOnUnbuffScript(ch)
	}
}
