package game

import (
	"errors"
	"math"
	"sort"

	"github.com/boyism80/fm/maple/constant"
	pconst "github.com/boyism80/fm/protocol/constant"
	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
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
	slots      uint8
	meso       int32
	tabs       map[constant.InventoryType][]Item
	kept       []*internal.Inventory
	opened     bool
	storeFee   int32
	takeOutFee int32
}

func (s *Storage) Open(npcID uint32, storeFee int32, takeOutFee int32) {
	s.opened = true
	s.storeFee = storeFee
	s.takeOutFee = takeOutFee
	s.owner.listener.OnOpenStorage(s.owner, npcID)
}

func (s *Storage) Close() {
	s.opened = false
}

func (s *Storage) count() int {
	count := len(s.kept)
	for _, items := range s.tabs {
		count += len(items)
	}
	return count
}

func (s *Storage) has(invType constant.InventoryType, itemID uint32) bool {
	for _, stored := range s.tabs[invType] {
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
	if s.count() >= int(s.slots) {
		return ErrStorageFull
	}

	model := s.owner.itemModel(itemID)
	if model == nil {
		return ErrStorageItemNotFound
	}
	invType := model.GetInventoryType()
	item := s.owner.inventory.GetItem(invType, slot)
	if item == nil || item.GetModel().GetID() != itemID {
		return ErrStorageItemNotFound
	}
	if !item.Storable() {
		return ErrStorageNotStorable
	}
	if model.GetTraits()&constant.ItemTraitUnique != 0 && s.has(invType, itemID) {
		return ErrStorageNotStorable
	}
	if model.GetTraits()&constant.ItemTraitRechargeable != 0 {
		count = item.GetCount()
	}
	if count == 0 || item.GetCount() < count {
		return ErrStorageItemNotFound
	}
	if (ExchangeSpec{Cost: ExchangeSide{Meso: s.storeFee}}).Valid(s.owner, s.owner.escrow) != ExchangeOK {
		return ErrStorageNotEnoughMeso
	}

	stored := item.Clone(count)
	stored.consumeKarma()
	s.owner.inventory.removeMesoUnchecked(s.storeFee)
	s.owner.inventory.RemoveItem(invType, slot, count)
	s.tabs[invType] = append(s.tabs[invType], stored)
	s.owner.listener.OnStorageTabChanged(s.owner, pconst.StorageResultStore, invType)
	return nil
}

func (s *Storage) TakeOut(invType constant.InventoryType, index uint8) error {
	if !s.opened {
		return ErrStorageClosed
	}
	items := s.tabs[invType]
	if int(index) >= len(items) {
		return ErrStorageItemNotFound
	}
	item := items[index]
	if s.owner.inventory.hasUniqueItem(item.GetModel()) {
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
	switch spec.Valid(s.owner, s.owner.escrow) {
	case ExchangeLackCost:
		return ErrStorageNotEnoughMeso
	case ExchangeLackCapacity:
		return ErrInventoryFull
	}

	s.tabs[invType] = append(items[:index:index], items[index+1:]...)
	s.owner.inventory.removeMesoUnchecked(s.takeOutFee)
	s.owner.inventory.addItemUnchecked(item, true)
	s.owner.listener.OnStorageTabChanged(s.owner, pconst.StorageResultTakeOut, invType)
	return nil
}

func (s *Storage) Arrange() error {
	if !s.opened {
		return ErrStorageClosed
	}

	for _, items := range s.tabs {
		sort.SliceStable(items, func(i, j int) bool {
			return items[i].GetModel().GetID() < items[j].GetModel().GetID()
		})
	}
	s.owner.listener.OnStorageArranged(s.owner)
	return nil
}

func (s *Storage) DepositMeso(amount int32) error {
	if !s.opened {
		return ErrStorageClosed
	}
	amount = min(amount, math.MaxInt32-s.meso)
	if amount <= 0 {
		return ErrStorageNotEnoughMeso
	}
	if (ExchangeSpec{Cost: ExchangeSide{Meso: amount}}).Valid(s.owner, s.owner.escrow) != ExchangeOK {
		return ErrStorageNotEnoughMeso
	}

	s.owner.inventory.removeMesoUnchecked(amount)
	s.meso += amount
	s.owner.listener.OnStorageMesoChanged(s.owner)
	return nil
}

func (s *Storage) WithdrawMeso(amount int32) error {
	if !s.opened {
		return ErrStorageClosed
	}
	amount = min(amount, math.MaxInt32-s.owner.inventory.meso)
	if amount <= 0 || amount > s.meso {
		return ErrStorageNotEnoughMeso
	}

	s.meso -= amount
	s.owner.inventory.addMesoUnchecked(amount)
	s.owner.listener.OnStorageMesoChanged(s.owner)
	return nil
}
