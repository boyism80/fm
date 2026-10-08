package constant

const (
	TeleportStoneCount    = 5
	VipTeleportStoneCount = 10
	TeleportStoneEmpty    = uint32(999999999)
)

type TeleportStoneSlots []uint32

func IsCashTeleportStone(itemID uint32) bool {
	return itemID/10000 == 504
}

func IsVipTeleportStone(itemID uint32) bool {
	return IsCashTeleportStone(itemID) && itemID/1000 != 5040
}

func NewTeleportStoneSlots(registered []uint32, size int) TeleportStoneSlots {
	stones := make(TeleportStoneSlots, size)
	for i := range stones {
		stones[i] = TeleportStoneEmpty
	}
	copy(stones, registered)
	return stones
}

func (r TeleportStoneSlots) Registered() []uint32 {
	out := make([]uint32, 0, len(r))
	for _, mapID := range r {
		if mapID != TeleportStoneEmpty {
			out = append(out, mapID)
		}
	}
	return out
}

func (r TeleportStoneSlots) Contains(mapID uint32) bool {
	for _, registered := range r {
		if registered == mapID {
			return true
		}
	}
	return false
}
