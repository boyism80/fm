package game

import (
	"errors"
	"time"

	"github.com/boyism80/fm/maple/constant"
	"github.com/boyism80/fm/maple/wz"
)

var (
	ErrPetNotFound     = errors.New("pet not found")
	ErrPetDead         = errors.New("pet is dead")
	ErrPetSkillInvalid = errors.New("invalid pet skill change")
	ErrPetNotRevivable = errors.New("pet cannot be revived")
)

type Pets struct {
	owner    *Character
	summoned uint64
	active   *ActivePet
	hpItem   uint32
	mpItem   uint32
}

func (p *Pets) Summon(slot int16) error {
	pet, ok := p.owner.inventory.GetItem(constant.InventoryTypeCash, slot).(*Pet)
	if !ok {
		return ErrPetNotFound
	}
	if p.active != nil && p.active.item == pet {
		p.Dismiss(constant.PetRemoveReasonNone)
		return nil
	}
	if !pet.Alive(time.Now()) {
		return ErrPetDead
	}

	p.Dismiss(constant.PetRemoveReasonNone)
	p.active = &ActivePet{owner: p.owner, item: pet, position: p.owner.position, stance: p.owner.stance}
	p.active.scheduleTimers()
	p.summoned = *pet.uniqueId
	p.owner.listener.OnPetSpawn(p.owner, *pet.uniqueId)
	return nil
}

func (p *Pets) ShowExceptions() {
	if p.active == nil || p.active.item == nil || p.active.item.uniqueId == nil {
		return
	}
	p.owner.listener.OnPetExceptions(p.owner, *p.active.item.uniqueId, p.active.item.exceptions)
}

func (p *Pets) Dismiss(reason constant.PetRemoveReason) {
	if p.active == nil {
		return
	}
	p.active.stopTimers()
	p.active = nil
	p.summoned = 0
	p.owner.listener.OnPetRemove(p.owner, reason)
}

func (p *Pets) ChangeSkill(sn uint64, itemID uint32) error {
	model, ok := p.owner.gameWorld.GetResources().Items[itemID].(*wz.CashItem)
	if !ok || model.PetSkill == 0 {
		return ErrPetSkillInvalid
	}
	var pet *Pet
	for _, item := range p.owner.inventory.tabs[constant.InventoryTypeCash].items {
		if candidate, ok := item.(*Pet); ok && *candidate.uniqueId == sn {
			pet = candidate
			break
		}
	}
	if pet == nil {
		return ErrPetNotFound
	}

	skill := model.PetSkill
	if model.PetSkillAdd {
		if pet.skills&skill != 0 {
			return ErrPetSkillInvalid
		}
		if required, ok := constant.PetSkillRequires[skill]; ok && pet.skills&required == 0 {
			return ErrPetSkillInvalid
		}
		pet.skills |= skill
	} else {
		if pet.skills&skill == 0 {
			return ErrPetSkillInvalid
		}
		for dependent, required := range constant.PetSkillRequires {
			if required == skill && pet.skills&dependent != 0 {
				return ErrPetSkillInvalid
			}
		}
		pet.skills &^= skill
	}
	p.owner.listener.OnPetSkillChanged(p.owner, pet, *pet.uniqueId, skill, model.PetSkillAdd)
	return nil
}

func (p *Pets) Revive(slot int16) error {
	pet, ok := p.owner.inventory.GetItem(constant.InventoryTypeCash, slot).(*Pet)
	if !ok {
		return ErrPetNotFound
	}
	model := pet.GetModel().(*wz.Pet)
	if model.NoRevive || pet.Alive(time.Now()) {
		return ErrPetNotRevivable
	}

	pet.itemCore.expiration = time.Now().AddDate(0, 0, model.Life)
	p.owner.listener.OnPetUpdated(p.owner, pet)
	return nil
}

func (p *Pets) Expired() map[int16]*Pet {
	now := time.Now()
	pets := make(map[int16]*Pet)
	for slot, item := range p.owner.inventory.tabs[constant.InventoryTypeCash].items {
		if pet, ok := item.(*Pet); ok && !pet.Alive(now) && !pet.GetModel().(*wz.Pet).NoRevive {
			pets[slot] = pet
		}
	}
	return pets
}

func (p *Pets) restore() {
	if p.summoned == 0 {
		return
	}
	now := time.Now()
	for _, item := range p.owner.inventory.tabs[constant.InventoryTypeCash].items {
		pet, ok := item.(*Pet)
		if !ok || *pet.uniqueId != p.summoned {
			continue
		}
		if pet.Alive(now) {
			p.active = &ActivePet{owner: p.owner, item: pet, position: p.owner.position, stance: p.owner.stance}
			p.active.scheduleTimers()
			return
		}
		break
	}
	p.summoned = 0
}
