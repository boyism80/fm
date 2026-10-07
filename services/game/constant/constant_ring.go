package constant

import "slices"

var crushRings = []uint32{1112001, 1112002, 1112003, 1112005, 1112006, 1112007, 1112012, 1112015, 1048000, 1048001, 1048002}

var friendshipRings = []uint32{1112800, 1112801, 1112802, 1112810, 1112811, 1112812, 1112816, 1112817, 1049000}

func IsCrushRing(itemID uint32) bool {
	return slices.Contains(crushRings, itemID)
}

func IsFriendshipRing(itemID uint32) bool {
	return slices.Contains(friendshipRings, itemID)
}
