package entity

import (
	"errors"
	"math/rand"
	"slices"
	"time"
	"unicode"

	"github.com/boyism80/fm/protocol/dto"
	"github.com/boyism80/fm/services/game/constant"
	"github.com/boyism80/fm/services/game/wz"
	"github.com/boyism80/fm/types"
	"golang.org/x/text/encoding/korean"
)

const (
	petHungerTimer = "pet_hunger"
	petLifeTimer   = "pet_life"
)

var (
	ErrPetFoodInvalid   = errors.New("invalid pet food")
	ErrPetFull          = errors.New("pet is full")
	ErrPetLootRefused   = errors.New("pet loot refused")
	ErrPetPotionRefused = errors.New("pet potion refused")
	ErrPetNameInvalid   = errors.New("invalid pet name")
)

type ActivePet struct {
	owner    *Character
	Item     *Pet
	Position types.Vector2[int16]
	Foothold int16
	Stance   uint8
}

func (p *ActivePet) ToDTO() *dto.ActivePet {
	return &dto.ActivePet{
		ItemID:   p.Item.GetModel().GetID(),
		Name:     p.Item.Name,
		SN:       *p.Item.UniqueId,
		Position: p.Position,
		Stance:   p.Stance,
		Foothold: p.Foothold,
	}
}

func (p *ActivePet) Move(start types.Vector2[int16], fragments []dto.MoveFragment) {
	for _, m := range fragments {
		if move, ok := m.(*dto.AbsoluteLifeMovement); ok {
			p.Position = move.Position
			p.Foothold = move.Foothold
		}
		p.Stance = m.GetStance()
	}
	p.owner.Listener.OnPetMove(p.owner, start, fragments)
}

func (p *ActivePet) Chat(typ uint8, action uint8, text string) {
	if text == "" {
		return
	}
	p.owner.Listener.OnPetChat(p.owner, typ, action, text)
}

func (p *ActivePet) Command(index uint8, calledByName bool) {
	commands := p.Item.GetModel().(*wz.Pet).Commands
	if int(index) >= len(commands) {
		p.owner.Listener.OnPetCommand(p.owner, index, false)
		return
	}
	command := commands[index]
	level := int(p.Item.Level)
	if level < command.MinLevel || level > command.MaxLevel {
		p.owner.Listener.OnPetCommand(p.owner, index, false)
		return
	}

	prob := command.Prob
	if calledByName {
		prob += constant.PetCalledByNameBonus
	}
	success := rand.Intn(100) < prob
	if success {
		p.AddCloseness(command.Inc)
	}
	p.owner.Listener.OnPetCommand(p.owner, index, success)
}

func (p *ActivePet) Feed(slot int16, itemID uint32) error {
	consume, ok := p.owner.Inventory.GetItem(constant.InventoryTypeConsume, slot).(*Consume)
	if ok == false || consume.GetModel().GetID() != itemID {
		return ErrPetFoodInvalid
	}
	food := consume.GetModel().(*wz.Consume).PetFood
	if food == nil || food.Feeds(p.Item.GetModel().GetID()) == false {
		return ErrPetFoodInvalid
	}
	if p.Item.Fullness >= constant.PetMaxFullness {
		p.owner.Listener.OnPetFood(p.owner, false)
		return ErrPetFull
	}

	p.owner.Inventory.RemoveItem(constant.InventoryTypeConsume, slot, 1)
	p.Item.Fullness = uint8(min(int(p.Item.Fullness)+food.Fullness, constant.PetMaxFullness))
	if rand.Intn(100) < constant.PetFoodClosenessRate {
		p.AddCloseness(1)
	} else {
		p.owner.Listener.OnPetUpdated(p.owner, p.Item)
	}
	p.owner.Listener.OnPetFood(p.owner, true)
	return nil
}

func (p *ActivePet) Loot(oid uint32, position types.Vector2[int16]) error {
	m := p.owner.GetMap()
	if m == nil {
		return ErrPetLootRefused
	}
	obj := m.GetItems()[oid]
	if obj == nil {
		return ErrPetLootRefused
	}
	if item, ok := obj.(Item); ok {
		if p.Item.Skills&constant.PetSkillPickupItem == 0 || slices.Contains(p.Item.Exceptions, item.GetModel().GetID()) {
			return ErrPetLootRefused
		}
	}
	if position.DistanceSq(obj.GetPosition()) > constant.PetLootClientRange*constant.PetLootClientRange {
		return ErrPetLootRefused
	}
	if p.Position.DistanceSq(obj.GetPosition()) > constant.PetLootRange*constant.PetLootRange {
		return ErrPetLootRefused
	}

	switch m.LootItem(obj, p.owner, position) {
	case constant.LootSuccess:
		return m.RemoveItem(oid, constant.RemoveItemTypeLootByPet, p.owner.GetID())
	case constant.LootPartial:
		p.owner.Listener.OnItemGainFailed(p.owner, constant.ItemGainFailedTypeFull)
		return nil
	case constant.LootFailedInventoryFull, constant.LootFailedMesoFull:
		p.owner.Listener.OnItemGainFailed(p.owner, constant.ItemGainFailedTypeFull)
		return ErrPetLootRefused
	default:
		return ErrPetLootRefused
	}
}

func (p *ActivePet) UsePotion(slot int16, itemID uint32) error {
	var skill constant.PetSkill
	if itemID == p.owner.PetHPItem {
		skill |= constant.PetSkillConsumeHP
	}
	if itemID == p.owner.PetMPItem {
		skill |= constant.PetSkillConsumeMP
	}
	if p.Item.Skills&skill == 0 {
		return ErrPetPotionRefused
	}
	if p.owner.IsAlive() == false {
		return ErrPetPotionRefused
	}
	if m := p.owner.GetMap(); m == nil || m.Wz.Limits(constant.FieldLimitPotion) {
		return ErrPetPotionRefused
	}
	consume, ok := p.owner.Inventory.GetItem(constant.InventoryTypeConsume, slot).(*Consume)
	if ok == false || consume.GetModel().GetID() != itemID {
		return ErrPetPotionRefused
	}

	if p.owner.UseConsume(consume) == false {
		return ErrPetPotionRefused
	}
	p.owner.Inventory.RemoveItem(constant.InventoryTypeConsume, slot, 1)
	return nil
}

func (p *ActivePet) FeedCash(itemID uint32) error {
	model, ok := p.owner.GameWorld.GetResources().Items[itemID].(*wz.CashItem)
	if ok == false || model.PetFood == nil || model.PetFood.Feeds(p.Item.GetModel().GetID()) == false {
		return ErrPetFoodInvalid
	}

	p.Item.Fullness = constant.PetMaxFullness
	p.AddCloseness(model.PetFood.Fullness)
	p.owner.Listener.OnPetFood(p.owner, true)
	return nil
}

func (p *ActivePet) Rename(name string) error {
	encoded, err := korean.EUCKR.NewEncoder().String(name)
	if err != nil || len(encoded) < constant.PetNameMinBytes || len(encoded) > constant.PetNameMaxBytes {
		return ErrPetNameInvalid
	}
	for _, r := range name {
		if unicode.IsLetter(r) == false && unicode.IsDigit(r) == false {
			return ErrPetNameInvalid
		}
	}

	p.Item.Name = name
	p.owner.Listener.OnPetNameChanged(p.owner)
	p.owner.Listener.OnPetUpdated(p.owner, p.Item)
	return nil
}

func (p *ActivePet) SetExceptions(itemIDs []uint32) {
	p.Item.Exceptions = itemIDs
	p.owner.Listener.OnPetExceptions(p.owner)
}

func (p *ActivePet) AddCloseness(n int) {
	levelUp := p.Item.AddCloseness(n)
	p.owner.Listener.OnPetUpdated(p.owner, p.Item)
	if levelUp {
		p.owner.Listener.OnPetLevelUp(p.owner)
	}
}

func (p *ActivePet) scheduleTimers() {
	p.owner.AddTimer(petHungerTimer, constant.PetHungerInterval, true, p.hunger)
	if p.Item.GetModel().(*wz.Pet).LimitedLife > 0 {
		p.owner.AddTimer(petLifeTimer, time.Second, true, p.age)
	}
}

func (p *ActivePet) stopTimers() {
	p.owner.RemoveTimer(petHungerTimer)
	p.owner.RemoveTimer(petLifeTimer)
}

func (p *ActivePet) hunger() {
	if p.Item.Alive(time.Now()) == false {
		p.owner.DismissPet(constant.PetRemoveReasonExpired)
		return
	}

	fullness := int(p.Item.Fullness) - p.Item.GetModel().(*wz.Pet).Hungry
	if fullness <= constant.PetStarveFullness {
		p.Item.Fullness = constant.PetStarvedFullness
		p.owner.DismissPet(constant.PetRemoveReasonHungry)
		return
	}
	p.Item.Fullness = uint8(fullness)
	p.owner.Listener.OnPetUpdated(p.owner, p.Item)
}

func (p *ActivePet) age() {
	if p.Item.SecondsLeft > 0 {
		p.Item.SecondsLeft--
	}
	if p.Item.SecondsLeft > 0 {
		return
	}

	p.owner.DismissPet(constant.PetRemoveReasonExpired)
	if slot, ok := p.owner.Inventory.FindSlot(constant.InventoryTypeCash, p.Item); ok {
		p.owner.Inventory.RemoveItem(constant.InventoryTypeCash, slot, 1)
	}
}
