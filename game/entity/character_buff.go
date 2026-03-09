package entity

import (
	"fmt"
	"log"
	"time"

	"github.com/boyism80/fm/core/luax"
	"github.com/boyism80/fm/game/constant"
	"github.com/boyism80/fm/game/wz"
)

type BuffEntity struct {
	Wz        *wz.Skill
	StartTime time.Time
	Level     uint8
	Flags     []constant.BuffFlag
	Values    map[constant.BuffFlag]int32
}

type BuffContainer struct {
	owner    *Character
	byFlag   map[constant.BuffFlag]*BuffEntity
	entities map[*BuffEntity]struct{}
}

func NewBuffContainer(owner *Character) *BuffContainer {
	if owner == nil {
		panic("buff container owner is nil")
	}
	return &BuffContainer{
		owner:    owner,
		byFlag:   make(map[constant.BuffFlag]*BuffEntity),
		entities: make(map[*BuffEntity]struct{}),
	}
}

func (bc *BuffContainer) add(entity *BuffEntity) (removed []*BuffEntity) {
	if bc == nil || entity == nil || len(entity.Flags) == 0 {
		return nil
	}

	conflicts := make(map[*BuffEntity]struct{})
	for _, flag := range entity.Flags {
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

func (bc *BuffContainer) remove(flags []constant.BuffFlag) (removed []*BuffEntity, removedFlags []constant.BuffFlag) {
	if bc == nil || len(flags) == 0 {
		return nil, nil
	}

	targets := make(map[*BuffEntity]struct{})
	for _, flag := range flags {
		if existing := bc.byFlag[flag]; existing != nil {
			targets[existing] = struct{}{}
		}
	}

	removedFlagSet := make(map[constant.BuffFlag]struct{})
	for entity := range targets {
		removed = append(removed, entity)
		for _, flag := range entity.Flags {
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

func (bc *BuffContainer) GetEntity(flag constant.BuffFlag) *BuffEntity {
	if bc == nil {
		return nil
	}
	return bc.byFlag[flag]
}

func (bc *BuffContainer) Entities() []*BuffEntity {
	if bc == nil {
		return nil
	}
	entities := make([]*BuffEntity, 0, len(bc.entities))
	for entity := range bc.entities {
		entities = append(entities, entity)
	}
	return entities
}

func (bc *BuffContainer) addEntity(entity *BuffEntity) {
	bc.entities[entity] = struct{}{}
	for _, flag := range entity.Flags {
		bc.byFlag[flag] = entity
	}
}

func (bc *BuffContainer) removeEntity(entity *BuffEntity) {
	delete(bc.entities, entity)
	for _, flag := range entity.Flags {
		if current := bc.byFlag[flag]; current == entity {
			delete(bc.byFlag, flag)
		}
	}
}

func (bc *BuffContainer) AddBuff(wz *wz.Skill, level uint8, values map[constant.BuffFlag]int32) {
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

	entity := &BuffEntity{
		Wz:        wz,
		StartTime: time.Now(),
		Level:     level,
		Flags:     entityFlags,
		Values:    entityValues,
	}

	removed := bc.add(entity)
	if len(removed) > 0 {
		ch.handleRemovedBuffEntities(removed)
	}

	if ch.Listener != nil {
		ch.Listener.OnBuffAdded(ch, wz, level, entityValues)
	}
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
	if ch.Listener != nil && len(removedFlags) > 0 {
		ch.Listener.OnBuffRemoved(ch, removedFlags)
	}
}

func (bc *BuffContainer) GetBuffValue(flag constant.BuffFlag) (*BuffEntity, int32, bool) {
	if bc == nil {
		return nil, 0, false
	}
	entity := bc.GetEntity(flag)
	if entity == nil || entity.Values == nil {
		return nil, 0, false
	}
	value, exists := entity.Values[flag]
	if !exists {
		return nil, 0, false
	}
	return entity, value, true
}

func (bc *BuffContainer) SetBuffValue(flag constant.BuffFlag, value int32) (*BuffEntity, bool) {
	if bc == nil {
		return nil, false
	}
	entity := bc.GetEntity(flag)
	if entity == nil || entity.Values == nil {
		return nil, false
	}
	entity.Values[flag] = value
	if bc.owner != nil && bc.owner.Listener != nil {
		bc.owner.Listener.OnBuffAdded(bc.owner, entity.Wz, entity.Level, map[constant.BuffFlag]int32{flag: value})
	}
	return entity, true
}

func (ch *Character) handleRemovedBuffEntities(removed []*BuffEntity) {
	if len(removed) == 0 {
		return
	}
	ch.invokeOnDeactivateScripts(removed)
}

func (ch *Character) invokeOnDeactivateScripts(removed []*BuffEntity) {
	if len(removed) == 0 || ch.Context == nil {
		return
	}

	mapInstance := ch.Context.GetMap(ch.Map)
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

	for _, entity := range removed {
		if entity == nil || entity.Wz == nil {
			continue
		}
		skillID := entity.Wz.ID

		skillEntry := &SkillEntry{
			Skill:      entity.Wz,
			SkillLevel: int(entity.Level),
			Owner:      ch,
		}

		scriptPath := fmt.Sprintf("script/skill/%d.lua", skillID)
		_, thread, err := luax.Call(root, scriptPath, "on_deactivated", ch, skillEntry)
		if err != nil {
			log.Printf("Failed to call on_deactivated for skill %d: %v", skillID, err)
			continue
		}
		if thread != nil {
			thread.Close()
		}
	}
}
