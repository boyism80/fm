package entity

import (
	pconst "github.com/boyism80/fm/protocol/constant"
	"github.com/boyism80/fm/services/game/constant"
	"github.com/boyism80/fm/types"
)

func (m *Map) AddEntrustedShop(es *EntrustedShop) {
	es.GameWorld = m.GameWorld
	es.Map = m
	es.OID = m.allocateOID()
	es.initTimers()
	if m.objects[constant.ObjectTypeEntrustedShop] == nil {
		m.objects[constant.ObjectTypeEntrustedShop] = make(map[uint32]Object)
	}
	m.objects[constant.ObjectTypeEntrustedShop][es.OID] = es
	m.sections.add(es)
}

func (m *Map) RemoveEntrustedShop(es *EntrustedShop) {
	if m.objects[constant.ObjectTypeEntrustedShop][es.OID] != es {
		return
	}

	es.BroadcastCall(func(obj Object) {
		viewer, ok := obj.(*Character)
		if ok == false {
			return
		}
		es.SendDestroySyncToViewer(viewer)
	}, nil)
	es.ClearTimers()
	m.sections.remove(es)
	delete(m.objects[constant.ObjectTypeEntrustedShop], es.OID)
	m.releaseOID(es.OID)
	es.Map = nil
}

func (m *Map) checkEntrustedShopSpot(position types.Vector2[int16]) error {
	near := func(x, y int16) bool {
		dx := int(position.X) - int(x)
		dy := int(position.Y) - int(y)
		return dx*dx+dy*dy < EntrustedShopSpacingSq
	}
	for _, portal := range m.Wz.Portals {
		if portal.Type == entrustedShopPortalType && near(portal.Position.X, portal.Position.Y) {
			return &MiniRoomEnterError{Code: pconst.MiniRoomEnterNearPortal}
		}
	}
	for _, obj := range m.objects[constant.ObjectTypeEntrustedShop] {
		other := obj.GetPosition()
		if near(other.X, other.Y) {
			return &MiniRoomEnterError{Code: pconst.MiniRoomEnterCannotOpen}
		}
	}
	return nil
}
