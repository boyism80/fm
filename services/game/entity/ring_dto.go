package entity

import (
	"github.com/boyism80/fm/protocol/dto"
)

func RingToDTO(ring *Ring) *dto.Ring {
	if ring == nil {
		return nil
	}

	return &dto.Ring{
		RingId:       ring.RingId,
		PartnerId:    ring.PartnerId,
		RingUniqueId: ring.RingUniqueId,
		PartnerChrId: ring.PartnerChrId,
		ItemId:       ring.ItemId,
		PartnerName:  ring.PartnerName,
		Equipped:     ring.Equipped,
	}
}

func RingsToDTO(rings []*Ring) []*dto.Ring {
	if rings == nil {
		return nil
	}

	result := make([]*dto.Ring, 0, len(rings))
	for _, ring := range rings {
		if dtoRing := RingToDTO(ring); dtoRing != nil {
			result = append(result, dtoRing)
		}
	}
	return result
}
