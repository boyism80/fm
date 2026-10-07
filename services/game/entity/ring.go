package entity

import (
	"github.com/boyism80/fm/protocol/dto"
	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
	"github.com/boyism80/fm/services/game/constant"
)

type Ring struct {
	RingId       uint64
	PartnerId    uint64
	RingUniqueId uint64
	PartnerChrId uint32
	ItemId       uint32
	PartnerName  string
	Equipped     bool
}

type RingContainer struct {
	Left []*Ring
	Mid  []*Ring
}

func (inv *Inventory) LoadRings(rings []*internal.CashRing) {
	for _, pb := range rings {
		ring := &Ring{
			RingId:       pb.GetSerial(),
			PartnerId:    pb.GetPartnerSerial(),
			RingUniqueId: pb.GetSerial(),
			PartnerChrId: pb.GetPartnerId(),
			ItemId:       pb.GetItemId(),
			PartnerName:  pb.GetPartnerName(),
		}
		switch {
		case constant.IsCrushRing(ring.ItemId):
			inv.Rings.Left = append(inv.Rings.Left, ring)
		case constant.IsFriendshipRing(ring.ItemId):
			inv.Rings.Mid = append(inv.Rings.Mid, ring)
		}
	}
}

func (inv *Inventory) WornRing(rings []*Ring) *dto.Ring {
	for _, ring := range rings {
		for _, equipment := range inv.Equipped {
			if equipment == nil {
				continue
			}
			if uniqueID := equipment.GetEquipmentCore().UniqueId; uniqueID != nil && *uniqueID == ring.RingId {
				return RingToDTO(ring)
			}
		}
	}
	return nil
}
