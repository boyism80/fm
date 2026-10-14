package domain

import (
	"errors"
	"math"
	"sort"

	pconst "github.com/boyism80/fm/protocol/constant"
	"github.com/boyism80/fm/services/game/constant"
)

var (
	ErrStorageClosed        = errors.New("storage is not open")
	ErrStorageFull          = errors.New("storage is full")
	ErrStorageNotEnoughMeso = errors.New("not enough meso")
	ErrStorageItemNotFound  = errors.New("storage item not found")
	ErrStorageNotStorable   = errors.New("item cannot be stored")
)

type Storage struct {
	owner      *Character
	Slots      uint8
	Meso       int32
	Tabs       map[constant.InventoryType][]Item
	opened     bool
	storeFee   int32
	takeOutFee int32
}

func (s *Storage) Open(npcID uint32, storeFee int32, takeOutFee int32) {
	s.opened = true
	s.storeFee = storeFee
	s.takeOutFee = takeOutFee
	s.owner.Listener.OnOpenStorage(s.owner, npcID)
}

func (s *Storage) Close() {
	s.opened = false
}

func (s *Storage) count() int {
	count := 0
	for _, items := range s.Tabs {
		count += len(items)
	}
	return count
}

func (s *Storage) has(invType constant.InventoryType, itemID uint32) bool {
	for _, stored := range s.Tabs[invType] {
		if stored.GetModel().GetID() == itemID {
			return true
		}
	}
	return false
}

func (s *Storage) Store(slot int16, itemID uint32, count uint16) error {
	if !s.opened {
		return ErrStorageClosed
	}
	if s.count() >= int(s.Slots) {
		return ErrStorageFull
	}

	model := s.owner.itemModel(itemID)
	if model == nil {
		return ErrStorageItemNotFound
	}
	invType := model.GetInventoryType()
	item := s.owner.Inventory.GetItem(invType, slot)
	if item == nil || item.GetModel().GetID() != itemID {
		return ErrStorageItemNotFound
	}
	if !item.Storable() {
		return ErrStorageNotStorable
	}
	if model.GetTraits()&constant.ItemTraitOnly != 0 && s.has(invType, itemID) {
		return ErrStorageNotStorable
	}
	if model.GetTraits()&constant.ItemTraitRechargeable != 0 {
		count = item.GetCount()
	}
	if count == 0 || item.GetCount() < count {
		return ErrStorageItemNotFound
	}
	if (ExchangeSpec{Cost: ExchangeSide{Meso: s.storeFee}}).Valid(s.owner) != ExchangeOK {
		return ErrStorageNotEnoughMeso
	}

	stored := item.Clone(count)
	stored.consumeKarma()
	s.owner.Inventory.removeMesoUnchecked(s.storeFee)
	s.owner.Inventory.RemoveItem(invType, slot, count)
	s.Tabs[invType] = append(s.Tabs[invType], stored)
	s.owner.Listener.OnStorageTabChanged(s.owner, pconst.StorageResultStore, invType)
	return nil
}

func (s *Storage) TakeOut(invType constant.InventoryType, index uint8) error {
	if !s.opened {
		return ErrStorageClosed
	}
	items := s.Tabs[invType]
	if int(index) >= len(items) {
		return ErrStorageItemNotFound
	}
	item := items[index]
	if item.GetModel().GetTraits()&constant.ItemTraitOnly != 0 && s.owner.Inventory.HasItem(item.GetModel().GetID()) {
		return ErrStorageNotStorable
	}

	spec := ExchangeSpec{
		Cost: ExchangeSide{
			Meso: s.takeOutFee,
		},
		Reward: ExchangeSide{
			Items: map[uint32]uint16{item.GetModel().GetID(): item.GetCount()},
		},
	}
	switch spec.Valid(s.owner) {
	case ExchangeLackCost:
		return ErrStorageNotEnoughMeso
	case ExchangeLackCapacity:
		return ErrInventoryFull
	}

	s.Tabs[invType] = append(items[:index:index], items[index+1:]...)
	s.owner.Inventory.removeMesoUnchecked(s.takeOutFee)
	s.owner.Inventory.addItemUnchecked(item, true)
	s.owner.Listener.OnStorageTabChanged(s.owner, pconst.StorageResultTakeOut, invType)
	return nil
}

func (s *Storage) Arrange() error {
	if !s.opened {
		return ErrStorageClosed
	}

	for _, items := range s.Tabs {
		sort.SliceStable(items, func(i, j int) bool {
			return items[i].GetModel().GetID() < items[j].GetModel().GetID()
		})
	}
	s.owner.Listener.OnStorageArranged(s.owner)
	return nil
}

func (s *Storage) DepositMeso(amount int32) error {
	if !s.opened {
		return ErrStorageClosed
	}
	amount = min(amount, math.MaxInt32-s.Meso)
	if amount <= 0 {
		return ErrStorageNotEnoughMeso
	}
	if (ExchangeSpec{Cost: ExchangeSide{Meso: amount}}).Valid(s.owner) != ExchangeOK {
		return ErrStorageNotEnoughMeso
	}

	s.owner.Inventory.removeMesoUnchecked(amount)
	s.Meso += amount
	s.owner.Listener.OnStorageMesoChanged(s.owner)
	return nil
}

func (s *Storage) WithdrawMeso(amount int32) error {
	if !s.opened {
		return ErrStorageClosed
	}
	amount = min(amount, math.MaxInt32-s.owner.Inventory.Meso)
	if amount <= 0 || amount > s.Meso {
		return ErrStorageNotEnoughMeso
	}

	s.Meso -= amount
	s.owner.Inventory.addMesoUnchecked(amount)
	s.owner.Listener.OnStorageMesoChanged(s.owner)
	return nil
}
