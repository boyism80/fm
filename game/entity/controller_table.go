package entity

type ControllerTable struct {
	controllers        map[uint32]*Character
	mobs               map[uint32]*Mob
	controller2mob     map[uint32]map[uint32]struct{}
	mob2controller     map[uint32]uint32
	onControllerChange func(mob *Mob, before *Character, after *Character)
}

func NewControllerTable(onControllerChange func(mob *Mob, before *Character, after *Character)) *ControllerTable {
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

	for mobOID, mob := range t.mobs {
		if t.mob2controller[mobOID] == 0 {
			t.assign(mob, nil, character)
		}
	}
}

func (t *ControllerTable) LeavePlayer(character *Character) {
	playerID := character.GetID()

	var next *Character
	for id, controller := range t.controllers {
		if id != playerID {
			next = controller
			break
		}
	}

	for mobOID := range t.controller2mob[playerID] {
		if mob, exists := t.mobs[mobOID]; exists {
			t.assign(mob, character, next)
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
		t.assign(mob, nil, controller)
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
			t.onControllerChange(mob, controller, nil)
		}
	}

	delete(t.mobs, mobOID)
	delete(t.mob2controller, mobOID)
	return controller
}

func (t *ControllerTable) assign(mob *Mob, before *Character, after *Character) {
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
		t.onControllerChange(mob, before, after)
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

func (t *ControllerTable) SwitchController(mob *Mob, newController *Character) {
	before, _ := t.GetController(mob)
	if before == newController {
		return
	}
	t.assign(mob, before, newController)
}
