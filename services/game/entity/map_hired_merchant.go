package entity

import (
	pconst "github.com/boyism80/fm/protocol/constant"
	"github.com/boyism80/fm/services/game/constant"
	"github.com/boyism80/fm/types"
)

func (m *Map) AddHiredMerchant(hm *HiredMerchant) {
	hm.GameWorld = m.GameWorld
	hm.Map = m
	hm.OID = m.allocateOID()
	hm.initTimers()
	if m.objects[constant.ObjectTypeHiredMerchant] == nil {
		m.objects[constant.ObjectTypeHiredMerchant] = make(map[uint32]Object)
	}
	m.objects[constant.ObjectTypeHiredMerchant][hm.OID] = hm
	m.sections.add(hm)
}

func (m *Map) RemoveHiredMerchant(hm *HiredMerchant) {
	if m.objects[constant.ObjectTypeHiredMerchant][hm.OID] != hm {
		return
	}

	hm.BroadcastCall(func(obj Object) {
		viewer, ok := obj.(*Character)
		if ok == false {
			return
		}
		hm.SendDestroySyncToViewer(viewer)
	}, nil)
	hm.ClearTimers()
	m.sections.remove(hm)
	delete(m.objects[constant.ObjectTypeHiredMerchant], hm.OID)
	m.releaseOID(hm.OID)
	hm.Map = nil
}

func (m *Map) checkHiredMerchantSpot(position types.Vector2[int16]) error {
	near := func(x, y int16) bool {
		dx := int(position.X) - int(x)
		dy := int(position.Y) - int(y)
		return dx*dx+dy*dy < HiredMerchantSpacingSq
	}
	for _, portal := range m.Wz.Portals {
		if portal.Type == hiredMerchantPortalType && near(portal.Position.X, portal.Position.Y) {
			return &MiniRoomEnterError{Code: pconst.MiniRoomEnterNearPortal}
		}
	}
	for _, obj := range m.objects[constant.ObjectTypeHiredMerchant] {
		other := obj.GetPosition()
		if near(other.X, other.Y) {
			return &MiniRoomEnterError{Code: pconst.MiniRoomEnterCannotOpen}
		}
	}
	return nil
}
