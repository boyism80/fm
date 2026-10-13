package entity

import (
	"github.com/boyism80/fm/protocol/dto"
)

func (r *Ring) ToDTO() *dto.Ring {
	return &dto.Ring{
		RingId:       r.RingId,
		PartnerId:    r.PartnerId,
		RingUniqueId: r.RingUniqueId,
		PartnerChrId: r.PartnerChrId,
		ItemId:       r.ItemId,
		PartnerName:  r.PartnerName,
		Equipped:     r.Equipped,
	}
}

func RingsToDTO(rings []*Ring) []*dto.Ring {
	if rings == nil {
		return nil
	}

	result := make([]*dto.Ring, 0, len(rings))
	for _, ring := range rings {
		if ring == nil {
			continue
		}
		result = append(result, ring.ToDTO())
	}
	return result
}
