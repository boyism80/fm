package game

import (
	"github.com/boyism80/fm/maple/constant"
	"github.com/boyism80/fm/maple/wz"
	"github.com/boyism80/fm/protocol/response"
	"github.com/boyism80/fm/types"
)

type Npc struct {
	objectCore
	wz           *wz.NpcSpawn
	runtimeSpawn bool
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
	npcDTO := n.BuildS2C()
	viewer.Send(&response.SpawnNpc{
		NPC:     npcDTO,
		Visible: true,
	}, types.SEND_POLICY_ENCRYPT)
	viewer.Send(&response.NpcControl{
		NPC:     npcDTO,
		MiniMap: true,
	}, types.SEND_POLICY_ENCRYPT)
}

func (n *Npc) SendDestroySyncToViewer(viewer *Character) {
	if n == nil || viewer == nil {
		return
	}
	viewer.Send(&response.NpcRemoveControl{OID: n.oid}, types.SEND_POLICY_ENCRYPT)
	viewer.Send(&response.RemoveNpc{OID: n.oid}, types.SEND_POLICY_ENCRYPT)
}

func (n *Npc) ShowEffect(action string) {
	if n == nil || action == "" {
		return
	}
	n.Broadcast(&response.ShowNpcEffect{
		OID:    n.oid,
		Action: action,
	}, nil)
}
