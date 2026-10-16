package game

import (
	"github.com/boyism80/fm/protocol/dto"
)

func (r *Ring) BuildS2C() *dto.Ring {
	return &dto.Ring{
		RingId:       r.ringId,
		PartnerId:    r.partnerId,
		RingUniqueId: r.ringUniqueId,
		PartnerChrId: r.partnerChrId,
		ItemId:       r.itemId,
		PartnerName:  r.partnerName,
		Equipped:     r.equipped,
	}
}

func (r Rings) BuildS2C() dto.Rings {
	return dto.Rings{
		Left: r.buildS2C(r.left),
		Mid:  r.buildS2C(r.mid),
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
