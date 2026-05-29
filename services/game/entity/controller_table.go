package entity

type ControllerTable struct {
	controllers        map[uint32]*Character
	mobs               map[uint32]*Mob
	controller2mob     map[uint32]map[uint32]struct{}
	mob2controller     map[uint32]uint32
	onControllerChange func(mob *Mob, before *Character, after *Character, aggro bool)
}

func (t *ControllerTable) choiceControllablePlayer(excludeID uint32) *Character {
	for id, controller := range t.controllers {
		if id == excludeID {
			continue
		}
		if controller != nil && !controller.IsHidden() {
			return controller
		}
	}
	return nil
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
	playerID := character.GetID()
	t.controllers[playerID] = character
	t.controller2mob[playerID] = make(map[uint32]struct{})

	if character == nil || character.IsHidden() {
		return
	}

	for mobOID, mob := range t.mobs {
		if t.mob2controller[mobOID] == 0 {
			t.assign(mob, nil, character, false)
		}
	}
}

func (t *ControllerTable) LeavePlayer(character *Character) {
	playerID := character.GetID()
	next := t.choiceControllablePlayer(playerID)

	for mobOID := range t.controller2mob[playerID] {
		if mob, exists := t.mobs[mobOID]; exists {
			t.assign(mob, character, next, false)
		}
	}

	delete(t.controllers, playerID)
	delete(t.controller2mob, playerID)
}

func (t *ControllerTable) EnterMob(mob *Mob) {
	mobOID := mob.OID
	t.mobs[mobOID] = mob
	t.mob2controller[mobOID] = 0

	for _, controller := range t.controllers {
		if controller == nil || controller.IsHidden() {
			continue
		}
		t.assign(mob, nil, controller, false)
		break
	}
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

func (t *ControllerTable) Update(character *Character) {
	if character == nil {
		return
	}
	playerID := character.GetID()
	if _, exists := t.controllers[playerID]; !exists {
		return
	}
	if character.IsHidden() {
		next := t.choiceControllablePlayer(playerID)
		for mobOID := range t.controller2mob[playerID] {
			if mob, exists := t.mobs[mobOID]; exists {
				t.assign(mob, character, next, false)
			}
		}
		return
	}
	for mobOID, mob := range t.mobs {
		if t.mob2controller[mobOID] == 0 {
			t.assign(mob, nil, character, false)
		}
	}
}
