package entity

import (
	"github.com/boyism80/fm/protocol/response"
	"github.com/boyism80/fm/services/game/constant"
	"github.com/boyism80/fm/services/game/wz"
	"github.com/boyism80/fm/types"
)

type Npc struct {
	ObjectCore
	Wz *wz.NpcSpawn
}

func (n *Npc) GetObjectType() constant.ObjectType {
	return constant.ObjectTypeNpc
}

func (n *Npc) Is(typ constant.ObjectType) bool {
	return n.GetObjectType().Has(typ)
}

func (n *Npc) SendSpawnSyncToViewer(viewer *Character) {
	if n == nil || viewer == nil {
		return
	}
	npcDTO := n.ToDTO()
	viewer.Send(&response.SpawnNpc{
		NPC:     npcDTO,
		Visible: true,
	}, types.SEND_POLICY_ENCRYPT)
	viewer.Send(&response.NpcControl{
		NPC:     npcDTO,
		MiniMap: true,
	}, types.SEND_POLICY_ENCRYPT)
}
