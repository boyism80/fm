package entity

import (
	"github.com/boyism80/fm/services/game/constant"
)

type ControllerTable struct {
	controllers        map[uint32]*Character
	mobs               map[uint32]*Mob
	controller2mob     map[uint32]map[uint32]struct{}
	mob2controller     map[uint32]uint32
	onControllerChange func(mob *Mob, before *Character, after *Character, aggro bool)
}

func NewControllerTable(onControllerChange func(mob *Mob, before *Character, after *Character, aggro bool)) *ControllerTable {
	return &ControllerTable{
		onControllerChange: onControllerChange,
		controllers:        make(map[uint32]*Character),
		mobs:               make(map[uint32]*Mob),
		controller2mob:     make(map[uint32]map[uint32]struct{}),
		mob2controller:     make(map[uint32]uint32),
	}
}

func (t *ControllerTable) EnterPlayer(character *Character) {
	if character == nil {
		return
	}
	playerID := character.GetID()
	t.controllers[playerID] = character
	if _, ok := t.controller2mob[playerID]; !ok {
		t.controller2mob[playerID] = make(map[uint32]struct{})
	}
	if character.IsHidden() {
		return
	}
	for _, mob := range t.nearbyMobs(character) {
		t.reassign(mob)
	}
}

func (t *ControllerTable) LeavePlayer(character *Character) {
	if character == nil {
		return
	}
	playerID := character.GetID()
	released := make([]*Mob, 0, len(t.controller2mob[playerID]))
	for mobOID := range t.controller2mob[playerID] {
		if mob, exists := t.mobs[mobOID]; exists {
			released = append(released, mob)
			t.assign(mob, character, nil, false)
		}
	}

	delete(t.controllers, playerID)
	delete(t.controller2mob, playerID)

	for _, mob := range released {
		t.reassign(mob)
	}
}

func (t *ControllerTable) EnterMob(mob *Mob) {
	if mob == nil {
		return
	}
	mobOID := mob.OID
	t.mobs[mobOID] = mob
	if _, ok := t.mob2controller[mobOID]; !ok {
		t.mob2controller[mobOID] = 0
	}
	t.reassign(mob)
}

func (t *ControllerTable) LeaveMob(mob *Mob) *Character {
	mobOID := mob.OID
	controllerID := t.mob2controller[mobOID]

	var controller *Character
	if controllerID != 0 {
		controller = t.controllers[controllerID]
		delete(t.controller2mob[controllerID], mobOID)
		if t.onControllerChange != nil {
			t.onControllerChange(mob, controller, nil, false)
		}
	}

	delete(t.mobs, mobOID)
	delete(t.mob2controller, mobOID)
	return controller
}

func (t *ControllerTable) MovePlayer(character *Character) {
	if character == nil {
		return
	}
	playerID := character.GetID()
	if _, exists := t.controllers[playerID]; !exists {
		return
	}
	for _, mob := range t.MobsControlledBy(playerID) {
		t.reassign(mob)
	}
	if character.IsHidden() {
		return
	}
	for _, mob := range t.nearbyMobs(character) {
		t.reassign(mob)
	}
}

func (t *ControllerTable) MoveMob(mob *Mob) {
	if mob == nil {
		return
	}
	if _, exists := t.mobs[mob.OID]; !exists {
		return
	}
	t.reassign(mob)
}

func (t *ControllerTable) Update(character *Character) {
	if character == nil {
		return
	}
	playerID := character.GetID()
	if _, exists := t.controllers[playerID]; !exists {
		return
	}
	if character.IsHidden() {
		for _, mob := range t.MobsControlledBy(playerID) {
			t.reassign(mob)
		}
		return
	}
	for _, mob := range t.nearbyMobs(character) {
		t.reassign(mob)
	}
}

func (t *ControllerTable) GetController(mob *Mob) (*Character, bool) {
	mobOID := mob.OID
	controllerID := t.mob2controller[mobOID]
	if controllerID == 0 {
		return nil, false
	}
	controller, exists := t.controllers[controllerID]
	return controller, exists
}

func (t *ControllerTable) SwitchController(mob *Mob, newController *Character, aggro bool) {
	if newController == nil || newController.IsHidden() {
		return
	}
	before, _ := t.GetController(mob)
	if before == newController {
		return
	}
	t.assign(mob, before, newController, aggro)
}

func (t *ControllerTable) MobsControlledBy(playerID uint32) []*Mob {
	ids := t.controller2mob[playerID]
	out := make([]*Mob, 0, len(ids))
	for mobOID := range ids {
		if mob, ok := t.mobs[mobOID]; ok {
			out = append(out, mob)
		}
	}
	return out
}

func (t *ControllerTable) reassign(mob *Mob) {
	current, _ := t.GetController(mob)
	next := t.choiceNearbyPlayer(mob, 0)
	if current != nil && !current.IsHidden() && t.canControl(current, mob) {
		return
	}
	if next == current {
		return
	}
	t.assign(mob, current, next, false)
}

func (t *ControllerTable) choiceNearbyPlayer(mob *Mob, excludeID uint32) *Character {
	mapInstance := mob.GetMap()
	if mapInstance == nil {
		return nil
	}
	for _, candidate := range mapInstance.GetObjectsNear(mob.GetPosition(), constant.ObjectTypeCharacter) {
		viewer, ok := candidate.(*Character)
		if !ok || viewer.IsHidden() {
			continue
		}
		if viewer.GetID() == excludeID {
			continue
		}
		if _, registered := t.controllers[viewer.GetID()]; !registered {
			continue
		}
		return viewer
	}
	return nil
}

func (t *ControllerTable) canControl(character *Character, mob *Mob) bool {
	if character == nil || mob == nil || character.IsHidden() {
		return false
	}
	mapInstance := mob.GetMap()
	if mapInstance == nil || character.GetMap() != mapInstance {
		return false
	}
	for _, candidate := range mapInstance.GetObjectsNear(mob.GetPosition(), constant.ObjectTypeCharacter) {
		if candidate == character {
			return true
		}
	}
	return false
}

func (t *ControllerTable) nearbyMobs(character *Character) []*Mob {
	mapInstance := character.GetMap()
	if mapInstance == nil {
		return nil
	}
	out := make([]*Mob, 0)
	for _, candidate := range mapInstance.GetObjectsNear(character.GetPosition(), constant.ObjectTypeMob) {
		mob, ok := candidate.(*Mob)
		if !ok {
			continue
		}
		if _, registered := t.mobs[mob.OID]; !registered {
			continue
		}
		out = append(out, mob)
	}
	return out
}

func (t *ControllerTable) assign(mob *Mob, before *Character, after *Character, aggro bool) {
	mobOID := mob.OID

	prevControllerID := t.mob2controller[mobOID]
	if prevControllerID != 0 {
		delete(t.controller2mob[prevControllerID], mobOID)
	}

	if after != nil {
		afterID := after.GetID()
		t.mob2controller[mobOID] = afterID
		if _, ok := t.controller2mob[afterID]; !ok {
			t.controller2mob[afterID] = map[uint32]struct{}{}
		}
		t.controller2mob[afterID][mobOID] = struct{}{}
	} else {
		t.mob2controller[mobOID] = 0
	}

	if t.onControllerChange != nil {
		t.onControllerChange(mob, before, after, aggro)
	}
}
