package entity

import (
	"github.com/boyism80/fm/services/game/constant"
	"github.com/boyism80/fm/services/game/wz"
)

type FlowResult int

const (
	FlowOK FlowResult = iota
	FlowLackCost
	FlowLackCapacity
)

type FlowSide struct {
	Items      map[uint32]uint16
	Meso       int32
	Exp        uint32
	Population int32
}

type FlowSpec struct {
	Cost   FlowSide
	Reward FlowSide
}

func flowItemsEmpty(items map[uint32]uint16) bool {
	if len(items) == 0 {
		return true
	}
	for _, count := range items {
		if count > 0 {
			return false
		}
	}
	return true
}

func (side FlowSide) isEmpty() bool {
	return side.Meso <= 0 && side.Exp == 0 && side.Population <= 0 && flowItemsEmpty(side.Items)
}

func groupFlowItemsByInventory(
	items map[uint32]uint16,
	modelOf func(uint32) wz.Item,
) map[constant.InventoryType]map[uint32]uint16 {
	if len(items) == 0 {
		return nil
	}
	out := make(map[constant.InventoryType]map[uint32]uint16)
	for id, count := range items {
		if count == 0 {
			continue
		}
		if modelOf != nil && modelOf(id) == nil {
			continue
		}
		invType := constant.GetInventoryTypeByItemID(id)
		if out[invType] == nil {
			out[invType] = make(map[uint32]uint16)
		}
		out[invType][id] += count
	}
	return out
}
