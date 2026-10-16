package game

import (
	"github.com/boyism80/fm/maple/constant"
	"github.com/boyism80/fm/protocol/dto"
	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
)

type Ring struct {
	ringId       uint64
	partnerId    uint64
	ringUniqueId uint64
	partnerChrId uint32
	itemId       uint32
	partnerName  string
	equipped     bool
}

type Rings struct {
	left []*Ring
	mid  []*Ring
}

func (inv *Inventory) LoadRings(rings []*internal.CashRing) {
	for _, pb := range rings {
		ring := &Ring{
			ringId:       pb.GetSerial(),
			partnerId:    pb.GetPartnerSerial(),
			ringUniqueId: pb.GetSerial(),
			partnerChrId: pb.GetPartnerId(),
			itemId:       pb.GetItemId(),
			partnerName:  pb.GetPartnerName(),
		}
		model := inv.owner.itemModel(ring.itemId)
		if model == nil {
			continue
		}
		switch {
		case model.GetTraits()&constant.ItemTraitCrushRing != 0:
			inv.rings.left = append(inv.rings.left, ring)
		case model.GetTraits()&constant.ItemTraitFriendshipRing != 0:
			inv.rings.mid = append(inv.rings.mid, ring)
		}
	}
}

func (inv *Inventory) WornRing(rings []*Ring) *dto.Ring {
	for _, ring := range rings {
		for _, equipment := range inv.equipped {
			if equipment == nil {
				continue
			}
			if uniqueID := equipment.GetEquipmentCore().uniqueId; uniqueID != nil && *uniqueID == ring.ringId {
				return ring.BuildS2C()
			}
		}
	}
	return nil
}
