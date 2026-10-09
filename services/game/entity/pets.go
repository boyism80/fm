package entity

import (
	"errors"
	"time"

	"github.com/boyism80/fm/services/game/constant"
	"github.com/boyism80/fm/services/game/wz"
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
	Active   *ActivePet
	HPItem   uint32
	MPItem   uint32
}

func (p *Pets) Summon(slot int16) error {
	pet, ok := p.owner.Inventory.GetItem(constant.InventoryTypeCash, slot).(*Pet)
	if !ok {
		return ErrPetNotFound
	}
	if p.Active != nil && p.Active.Item == pet {
		p.Dismiss(constant.PetRemoveReasonNone)
		return nil
	}
	if !pet.Alive(time.Now()) {
		return ErrPetDead
	}

	p.Dismiss(constant.PetRemoveReasonNone)
	p.Active = &ActivePet{owner: p.owner, Item: pet, Position: p.owner.Position, Stance: p.owner.Stance}
	p.Active.scheduleTimers()
	p.summoned = *pet.UniqueId
	p.owner.Listener.OnPetSpawn(p.owner)
	return nil
}

func (p *Pets) Dismiss(reason constant.PetRemoveReason) {
	if p.Active == nil {
		return
	}
	p.Active.stopTimers()
	p.Active = nil
	p.summoned = 0
	p.owner.Listener.OnPetRemove(p.owner, reason)
}

func (p *Pets) ChangeSkill(sn uint64, itemID uint32) error {
	model, ok := p.owner.GameWorld.GetResources().Items[itemID].(*wz.CashItem)
	if !ok || model.PetSkill == 0 {
		return ErrPetSkillInvalid
	}
	var pet *Pet
	for _, item := range p.owner.Inventory.Tabs[constant.InventoryTypeCash].Items {
		if candidate, ok := item.(*Pet); ok && *candidate.UniqueId == sn {
			pet = candidate
			break
		}
	}
	if pet == nil {
		return ErrPetNotFound
	}

	skill := model.PetSkill
	if model.PetSkillAdd {
		if pet.Skills&skill != 0 {
			return ErrPetSkillInvalid
		}
		if required, ok := constant.PetSkillRequires[skill]; ok && pet.Skills&required == 0 {
			return ErrPetSkillInvalid
		}
		pet.Skills |= skill
	} else {
		if pet.Skills&skill == 0 {
			return ErrPetSkillInvalid
		}
		for dependent, required := range constant.PetSkillRequires {
			if required == skill && pet.Skills&dependent != 0 {
				return ErrPetSkillInvalid
			}
		}
		pet.Skills &^= skill
	}
	p.owner.Listener.OnPetSkillChanged(p.owner, pet, skill, model.PetSkillAdd)
	return nil
}

func (p *Pets) Revive(slot int16) error {
	pet, ok := p.owner.Inventory.GetItem(constant.InventoryTypeCash, slot).(*Pet)
	if !ok {
		return ErrPetNotFound
	}
	model := pet.GetModel().(*wz.Pet)
	if model.NoRevive || pet.Alive(time.Now()) {
		return ErrPetNotRevivable
	}

	pet.ItemCore.Expiration = time.Now().AddDate(0, 0, model.Life)
	p.owner.Listener.OnPetUpdated(p.owner, pet)
	return nil
}

func (p *Pets) Expired() map[int16]*Pet {
	now := time.Now()
	pets := make(map[int16]*Pet)
	for slot, item := range p.owner.Inventory.Tabs[constant.InventoryTypeCash].Items {
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
	for _, item := range p.owner.Inventory.Tabs[constant.InventoryTypeCash].Items {
		pet, ok := item.(*Pet)
		if !ok || *pet.UniqueId != p.summoned {
			continue
		}
		if pet.Alive(now) {
			p.Active = &ActivePet{owner: p.owner, Item: pet, Position: p.owner.Position, Stance: p.owner.Stance}
			p.Active.scheduleTimers()
			return
		}
		break
	}
	p.summoned = 0
}
