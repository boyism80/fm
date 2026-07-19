package entity

import (
	"fmt"

	"github.com/boyism80/fm/protocol/response"
	"github.com/boyism80/fm/services/game/constant"
)

func (m *Map) RemoveNpc(oid uint32) error {
	if m == nil {
		return fmt.Errorf("map is nil")
	}
	if m.objects[constant.ObjectTypeNpc] == nil {
		return fmt.Errorf("no npcs on map")
	}
	obj, ok := m.objects[constant.ObjectTypeNpc][oid]
	if !ok {
		return fmt.Errorf("npc %d not found on map", oid)
	}
	npc, ok := obj.(*Npc)
	if !ok || npc == nil {
		return fmt.Errorf("npc %d is not an npc", oid)
	}
	npc.ClearTimers()
	m.sections.remove(npc)
	delete(m.objects[constant.ObjectTypeNpc], oid)
	m.releaseOID(oid)
	m.Broadcast(&response.NpcRemoveControl{OID: oid}, nil)
	m.Broadcast(&response.RemoveNpc{OID: oid}, nil)
	return nil
}

func (m *Map) RemoveRuntimeNpcs() {
	if m == nil || m.objects[constant.ObjectTypeNpc] == nil {
		return
	}
	oids := make([]uint32, 0)
	for oid, obj := range m.objects[constant.ObjectTypeNpc] {
		npc, ok := obj.(*Npc)
		if !ok || npc == nil || !npc.RuntimeSpawn {
			continue
		}
		oids = append(oids, oid)
	}
	for _, oid := range oids {
		_ = m.RemoveNpc(oid)
	}
}

func (m *Map) NpcByTemplate(npcID uint32) *Npc {
	if m == nil {
		return nil
	}
	for _, obj := range m.GetNpcs() {
		npc, ok := obj.(*Npc)
		if !ok || npc == nil || npc.Wz == nil || npc.Wz.BaseSpawn == nil {
			continue
		}
		if npc.Wz.BaseSpawn.ID == npcID {
			return npc
		}
	}
	return nil
}
