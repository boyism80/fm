package entity

import (
	"github.com/boyism80/fm/protocol/dto"
)

func (n *Npc) ToDTO() *dto.Npc {
	if n == nil {
		return nil
	}

	npcId := uint32(0)
	foothold := int16(0)
	rx0 := int16(0)
	rx1 := int16(0)
	cy := int16(0)

	if n.Wz != nil && n.Wz.BaseSpawn != nil {
		npcId = n.Wz.ID
		foothold = n.Wz.Foothold
		rx0 = n.Wz.RenderX0
		rx1 = n.Wz.RenderX1
		cy = n.Wz.CollisionY
	}

	return &dto.Npc{
		OID:      n.OID,
		NpcId:    npcId,
		Position: n.Position,
		Foothold: foothold,
		Rx0:      rx0,
		Rx1:      rx1,
		Cy:       cy,
	}
}
