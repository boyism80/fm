package entity

import (
	"errors"
	"log"
	"math"
	"sort"

	"github.com/asynkron/protoactor-go/actor"
	pconst "github.com/boyism80/fm/protocol/constant"
	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
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

func (ch *Character) OpenStorage(actx actor.Context, npcID uint32, storeFee int32, takeOutFee int32) {
	if ch.Storage != nil {
		ch.Storage.Open(npcID, storeFee, takeOutFee)
		return
	}

	ch.Listener.LoadStorageAsync(actx, ch).Then(func(v interface{}) (interface{}, error) {
		if ch.Storage != nil {
			return nil, nil
		}
		ch.Storage = NewStorageFromInternalProto(v.(*internal.LoadStorageReply).GetStorage(), ch)
		ch.Storage.Open(npcID, storeFee, takeOutFee)
		return nil, nil
	}).OnError(func(err error) {
		log.Printf("OpenStorage character=%d: %v", ch.GetID(), err)
	})
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

func (s *Storage) Store(slot int16, itemID uint32, count uint16) error {
	if s.opened == false {
		return ErrStorageClosed
	}
	if s.count() >= int(s.Slots) {
		return ErrStorageFull
	}

	invType := constant.GetInventoryTypeByItemID(itemID)
	item := s.owner.Inventory.GetItem(invType, slot)
	if item == nil || item.GetModel().GetID() != itemID {
		return ErrStorageItemNotFound
	}
	if constant.ItemCategoryOf(itemID) == constant.ItemCategoryPet {
		return ErrStorageNotStorable
	}
	if constant.IsRechargeable(itemID) {
		count = item.GetCount()
	}
	if count == 0 || item.GetCount() < count {
		return ErrStorageItemNotFound
	}
	if (ExchangeSpec{Cost: ExchangeSide{Meso: s.storeFee}}).Valid(s.owner) != ExchangeOK {
		return ErrStorageNotEnoughMeso
	}

	stored := item.Clone(count)
	s.owner.Inventory.removeMesoUnchecked(s.storeFee)
	s.owner.Inventory.RemoveItem(invType, slot, count)
	s.Tabs[invType] = append(s.Tabs[invType], stored)
	s.owner.Listener.OnStorageTabChanged(s.owner, pconst.StorageResultStore, invType)
	return nil
}

func (s *Storage) TakeOut(invType constant.InventoryType, index uint8) error {
	if s.opened == false {
		return ErrStorageClosed
	}
	items := s.Tabs[invType]
	if int(index) >= len(items) {
		return ErrStorageItemNotFound
	}
	item := items[index]

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
	if s.opened == false {
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
	if s.opened == false {
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
	if s.opened == false {
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
