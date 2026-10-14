package domain

import (
	"github.com/boyism80/fm/protocol/dto"
)

func (r *Ring) BuildS2C() *dto.Ring {
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

func (r Rings) BuildS2C() dto.Rings {
	return dto.Rings{
		Left: r.buildS2C(r.Left),
		Mid:  r.buildS2C(r.Mid),
	}
}

func (r Rings) buildS2C(rings []*Ring) []*dto.Ring {
	if rings == nil {
		return nil
	}

	result := make([]*dto.Ring, 0, len(rings))
	for _, ring := range rings {
		if ring == nil {
			continue
		}
		result = append(result, ring.BuildS2C())
	}
	return result
}
