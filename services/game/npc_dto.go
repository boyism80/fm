package game

import (
	"github.com/boyism80/fm/protocol/dto"
)

func (n *Npc) BuildS2C() *dto.Npc {
	if n == nil {
		return nil
	}

	npcId := uint32(0)
	foothold := int16(0)
	rx0 := int16(0)
	rx1 := int16(0)
	cy := int16(0)

	if n.wz != nil && n.wz.BaseSpawn != nil {
		npcId = n.wz.ID
		foothold = n.wz.Foothold
		rx0 = n.wz.RenderX0
		rx1 = n.wz.RenderX1
		cy = n.wz.CollisionY
	}

	return &dto.Npc{
		OID:      n.oid,
		NpcId:    npcId,
		Position: n.position,
		Foothold: foothold,
		Rx0:      rx0,
		Rx1:      rx1,
		Cy:       cy,
	}
}
