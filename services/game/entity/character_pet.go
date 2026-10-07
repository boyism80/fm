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

func (ch *Character) SummonPet(slot int16) error {
	pet, ok := ch.Inventory.GetItem(constant.InventoryTypeCash, slot).(*Pet)
	if ok == false {
		return ErrPetNotFound
	}
	if ch.Pet != nil && ch.Pet.Item == pet {
		ch.DismissPet(constant.PetRemoveReasonNone)
		return nil
	}
	if pet.Alive(time.Now()) == false {
		return ErrPetDead
	}

	ch.DismissPet(constant.PetRemoveReasonNone)
	ch.Pet = &ActivePet{owner: ch, Item: pet, Position: ch.Position, Stance: ch.Stance}
	ch.Pet.scheduleTimers()
	ch.summonedPet = *pet.UniqueId
	ch.Listener.OnPetSpawn(ch)
	return nil
}

func (ch *Character) DismissPet(reason constant.PetRemoveReason) {
	if ch.Pet == nil {
		return
	}
	ch.Pet.stopTimers()
	ch.Pet = nil
	ch.summonedPet = 0
	ch.Listener.OnPetRemove(ch, reason)
}

func (ch *Character) ChangePetSkill(sn uint64, itemID uint32) error {
	model, ok := ch.GameWorld.GetResources().Items[itemID].(*wz.CashItem)
	if ok == false || model.PetSkill == 0 {
		return ErrPetSkillInvalid
	}
	var pet *Pet
	for _, item := range ch.Inventory.Tabs[constant.InventoryTypeCash].Items {
		if p, ok := item.(*Pet); ok && *p.UniqueId == sn {
			pet = p
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
	ch.Listener.OnPetSkillChanged(ch, pet, skill, model.PetSkillAdd)
	return nil
}

func (ch *Character) RevivePet(slot int16) error {
	pet, ok := ch.Inventory.GetItem(constant.InventoryTypeCash, slot).(*Pet)
	if ok == false {
		return ErrPetNotFound
	}
	model := pet.GetModel().(*wz.Pet)
	if model.NoRevive || pet.Alive(time.Now()) {
		return ErrPetNotRevivable
	}

	pet.ItemCore.Expiration = time.Now().AddDate(0, 0, model.Life)
	ch.Listener.OnPetUpdated(ch, pet)
	return nil
}

func (ch *Character) ExpiredPets() map[int16]*Pet {
	now := time.Now()
	pets := make(map[int16]*Pet)
	for slot, item := range ch.Inventory.Tabs[constant.InventoryTypeCash].Items {
		if pet, ok := item.(*Pet); ok && pet.Alive(now) == false && pet.GetModel().(*wz.Pet).NoRevive == false {
			pets[slot] = pet
		}
	}
	return pets
}

func (ch *Character) restorePet() {
	if ch.summonedPet == 0 {
		return
	}
	now := time.Now()
	for _, item := range ch.Inventory.Tabs[constant.InventoryTypeCash].Items {
		pet, ok := item.(*Pet)
		if ok == false || *pet.UniqueId != ch.summonedPet {
			continue
		}
		if pet.Alive(now) {
			ch.Pet = &ActivePet{owner: ch, Item: pet, Position: ch.Position, Stance: ch.Stance}
			ch.Pet.scheduleTimers()
			return
		}
		break
	}
	ch.summonedPet = 0
}
