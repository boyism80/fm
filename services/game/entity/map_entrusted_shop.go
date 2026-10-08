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
		es.sendDestroy(viewer)
	}, nil)
	es.ClearTimers()
	m.sections.remove(es)
	delete(m.objects[constant.ObjectTypeEntrustedShop], es.OID)
	m.releaseOID(es.OID)
	es.Map = nil
}

func (m *Map) FindEntrustedShopByOwner(ownerID uint32) *EntrustedShop {
	for _, obj := range m.objects[constant.ObjectTypeEntrustedShop] {
		es, ok := obj.(*EntrustedShop)
		if ok == false || es.OwnerID != ownerID {
			continue
		}
		es.mu.Lock()
		published := es.published && es.ended == false
		es.mu.Unlock()
		if published {
			return es
		}
	}
	return nil
}

func (m *Map) AddPersonalShop(ps *PersonalShop) {
	ps.GameWorld = m.GameWorld
	ps.Map = m
	ps.OID = m.allocateOID()
	if m.objects[constant.ObjectTypePersonalShop] == nil {
		m.objects[constant.ObjectTypePersonalShop] = make(map[uint32]Object)
	}
	m.objects[constant.ObjectTypePersonalShop][ps.OID] = ps
}

func (m *Map) RemovePersonalShop(ps *PersonalShop) {
	if m.objects[constant.ObjectTypePersonalShop][ps.OID] != ps {
		return
	}

	delete(m.objects[constant.ObjectTypePersonalShop], ps.OID)
	m.releaseOID(ps.OID)
	ps.Map = nil
}

func (m *Map) checkShopSpot(position types.Vector2[int16]) error {
	near := func(x, y int16) bool {
		dx := int(position.X) - int(x)
		dy := int(position.Y) - int(y)
		return dx*dx+dy*dy < ShopSpacingSq
	}
	for _, portal := range m.Wz.Portals {
		if portal.Type == shopPortalType && near(portal.Position.X, portal.Position.Y) {
			return &MiniRoomEnterError{Code: pconst.MiniRoomEnterNearPortal}
		}
	}
	for _, typ := range []constant.ObjectType{constant.ObjectTypeEntrustedShop, constant.ObjectTypePersonalShop} {
		for _, obj := range m.objects[typ] {
			other := obj.GetPosition()
			if near(other.X, other.Y) {
				return &MiniRoomEnterError{Code: pconst.MiniRoomEnterCannotOpen}
			}
		}
	}
	return nil
}
