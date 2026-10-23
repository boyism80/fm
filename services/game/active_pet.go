package game

import (
	"errors"
	"math/rand"
	"slices"
	"time"
	"unicode"

	"github.com/boyism80/fm/maple/constant"
	"github.com/boyism80/fm/maple/wz"
	"github.com/boyism80/fm/protocol/dto"
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
	item     *Pet
	position types.Vector2[int16]
	foothold int16
	stance   uint8
}

func (p *ActivePet) BuildS2C() *dto.ActivePet {
	return &dto.ActivePet{
		ItemID:   p.item.GetModel().GetID(),
		Name:     p.item.name,
		SN:       *p.item.uniqueId,
		Position: p.position,
		Stance:   p.stance,
		Foothold: p.foothold,
	}
}

func (p *ActivePet) Move(start types.Vector2[int16], fragments []dto.MoveFragment) {
	for _, m := range fragments {
		if move, ok := m.(*dto.AbsoluteLifeMovement); ok {
			p.position = move.Position
			p.foothold = move.Foothold
		}
		p.stance = m.GetStance()
	}
	p.owner.listener.OnPetMove(p.owner, start, fragments)
}

func (p *ActivePet) Chat(typ uint8, action uint8, text string) {
	if text == "" {
		return
	}
	p.owner.listener.OnPetChat(p.owner, typ, action, text)
}

func (p *ActivePet) Command(index uint8, calledByName bool) {
	commands := p.item.GetModel().(*wz.Pet).Commands
	if int(index) >= len(commands) {
		p.owner.listener.OnPetCommand(p.owner, index, false)
		return
	}
	command := commands[index]
	level := int(p.item.level)
	if level < command.MinLevel || level > command.MaxLevel {
		p.owner.listener.OnPetCommand(p.owner, index, false)
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
	p.owner.listener.OnPetCommand(p.owner, index, success)
}

func (p *ActivePet) Feed(slot int16, itemID uint32) error {
	consume, ok := p.owner.inventory.GetItem(constant.InventoryTypeConsume, slot).(*Consume)
	if !ok || consume.GetModel().GetID() != itemID {
		return ErrPetFoodInvalid
	}
	food := consume.GetModel().(*wz.Consume).PetFood
	if food == nil || !food.Feeds(p.item.GetModel().GetID()) {
		return ErrPetFoodInvalid
	}
	if p.item.fullness >= constant.PetMaxFullness {
		p.owner.listener.OnPetFood(p.owner, false)
		return ErrPetFull
	}

	p.owner.inventory.RemoveItem(constant.InventoryTypeConsume, slot, 1)
	p.item.fullness = uint8(min(int(p.item.fullness)+food.Fullness, constant.PetMaxFullness))
	if rand.Intn(100) < constant.PetFoodClosenessRate {
		p.AddCloseness(1)
	} else {
		p.owner.listener.OnPetUpdated(p.owner, p.item)
	}
	p.owner.listener.OnPetFood(p.owner, true)
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
		if p.item.skills&constant.PetSkillPickupItem == 0 || slices.Contains(p.item.exceptions, item.GetModel().GetID()) {
			return ErrPetLootRefused
		}
	}
	if position.DistanceSq(obj.GetPosition()) > constant.PetLootClientRange*constant.PetLootClientRange {
		return ErrPetLootRefused
	}
	if p.position.DistanceSq(obj.GetPosition()) > constant.PetLootRange*constant.PetLootRange {
		return ErrPetLootRefused
	}

	switch m.LootItem(obj, p.owner, position) {
	case constant.LootSuccess:
		return m.RemoveItem(oid, constant.RemoveItemTypeLootByPet, p.owner.GetID())
	case constant.LootPartial:
		p.owner.listener.OnItemGainFailed(p.owner, constant.ItemGainFailedTypeFull)
		return nil
	case constant.LootFailedInventoryFull, constant.LootFailedMesoFull:
		p.owner.listener.OnItemGainFailed(p.owner, constant.ItemGainFailedTypeFull)
		return ErrPetLootRefused
	default:
		return ErrPetLootRefused
	}
}

func (p *ActivePet) UsePotion(slot int16, itemID uint32) error {
	var skill constant.PetSkill
	if itemID == p.owner.pets.hpItem {
		skill |= constant.PetSkillConsumeHP
	}
	if itemID == p.owner.pets.mpItem {
		skill |= constant.PetSkillConsumeMP
	}
	if p.item.skills&skill == 0 {
		return ErrPetPotionRefused
	}
	if !p.owner.IsAlive() {
		return ErrPetPotionRefused
	}
	if m := p.owner.GetMap(); m == nil || m.wz.Limits(constant.FieldLimitPotion) {
		return ErrPetPotionRefused
	}
	consume, ok := p.owner.inventory.GetItem(constant.InventoryTypeConsume, slot).(*Consume)
	if !ok || consume.GetModel().GetID() != itemID {
		return ErrPetPotionRefused
	}

	if p.owner.UseConsumeInSlot(slot, consume) != nil {
		return ErrPetPotionRefused
	}
	return nil
}

func (p *ActivePet) FeedCash(itemID uint32) error {
	model, ok := p.owner.gameWorld.GetResources().Items[itemID].(*wz.CashItem)
	if !ok || model.PetFood == nil || !model.PetFood.Feeds(p.item.GetModel().GetID()) {
		return ErrPetFoodInvalid
	}

	p.item.fullness = constant.PetMaxFullness
	p.AddCloseness(model.PetFood.Fullness)
	p.owner.listener.OnPetFood(p.owner, true)
	return nil
}

func (p *ActivePet) Rename(name string) error {
	encoded, err := korean.EUCKR.NewEncoder().String(name)
	if err != nil || len(encoded) < constant.PetNameMinBytes || len(encoded) > constant.PetNameMaxBytes {
		return ErrPetNameInvalid
	}
	for _, r := range name {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) {
			return ErrPetNameInvalid
		}
	}

	p.item.name = name
	p.owner.listener.OnPetNameChanged(p.owner, name)
	p.owner.listener.OnPetUpdated(p.owner, p.item)
	return nil
}

func (p *ActivePet) SetExceptions(itemIDs []uint32) {
	p.item.exceptions = itemIDs
	p.owner.listener.OnPetExceptions(p.owner, *p.item.uniqueId, itemIDs)
}

func (p *ActivePet) AddCloseness(n int) {
	p.SetCloseness(int(p.item.closeness) + n)
}

func (p *ActivePet) SetCloseness(n int) {
	levelUp := p.item.SetCloseness(n)
	p.owner.listener.OnPetUpdated(p.owner, p.item)
	if levelUp {
		p.owner.listener.OnPetLevelUp(p.owner)
	}
}

func (p *ActivePet) SetFullness(n int) {
	p.item.fullness = uint8(min(max(n, 0), constant.PetMaxFullness))
	p.owner.listener.OnPetUpdated(p.owner, p.item)
}

func (p *ActivePet) scheduleTimers() {
	p.owner.AddTimer(petHungerTimer, constant.PetHungerInterval, true, p.hunger)
	if p.item.GetModel().(*wz.Pet).LimitedLife > 0 {
		p.owner.AddTimer(petLifeTimer, time.Second, true, p.age)
	}
}

func (p *ActivePet) stopTimers() {
	p.owner.RemoveTimer(petHungerTimer)
	p.owner.RemoveTimer(petLifeTimer)
}

func (p *ActivePet) hunger() {
	if !p.item.Alive(time.Now()) {
		p.owner.pets.Dismiss(constant.PetRemoveReasonExpired)
		return
	}

	fullness := int(p.item.fullness) - p.item.GetModel().(*wz.Pet).Hungry
	if fullness <= constant.PetStarveFullness {
		p.item.fullness = constant.PetStarvedFullness
		p.owner.pets.Dismiss(constant.PetRemoveReasonHungry)
		return
	}
	p.item.fullness = uint8(fullness)
	p.owner.listener.OnPetUpdated(p.owner, p.item)
}

func (p *ActivePet) age() {
	if p.item.secondsLeft > 0 {
		p.item.secondsLeft--
	}
	if p.item.secondsLeft > 0 {
		return
	}

	p.owner.pets.Dismiss(constant.PetRemoveReasonExpired)
	if slot, ok := p.owner.inventory.FindSlot(constant.InventoryTypeCash, p.item); ok {
		p.owner.inventory.RemoveItem(constant.InventoryTypeCash, slot, 1)
	}
}
