package game

import (
	"github.com/boyism80/fm/maple/constant"
	pconst "github.com/boyism80/fm/protocol/constant"
	"github.com/boyism80/fm/types"
)

func (m *Map) AddEntrustedShop(es *EntrustedShop) {
	es.gameWorld = m.gameWorld
	es.onMap = m
	es.home = m
	es.oid = m.allocateOID()
	es.initTimers()
	if m.objects[constant.ObjectTypeEntrustedShop] == nil {
		m.objects[constant.ObjectTypeEntrustedShop] = make(map[uint32]Object)
	}
	m.objects[constant.ObjectTypeEntrustedShop][es.oid] = es
	m.sections.add(es)
}

func (m *Map) RemoveEntrustedShop(es *EntrustedShop) {
	if m.objects[constant.ObjectTypeEntrustedShop][es.oid] != es {
		return
	}

	es.BroadcastCall(func(obj Object) {
		viewer, ok := obj.(*Character)
		if !ok {
			return
		}
		es.SendDestroySyncToViewer(viewer)
	}, nil)
	es.ClearTimers()
	m.sections.remove(es)
	delete(m.objects[constant.ObjectTypeEntrustedShop], es.oid)
	m.releaseOID(es.oid)
	es.onMap = nil
}

func (m *Map) FindEntrustedShopByOwner(ownerID uint32) (uint32, bool) {
	for _, obj := range m.objects[constant.ObjectTypeEntrustedShop] {
		es, ok := obj.(*EntrustedShop)
		if ok && es.ownerID == ownerID && es.published {
			return es.oid, true
		}
	}
	return 0, false
}

func (m *Map) AddPersonalShop(ps *PersonalShop) {
	ps.gameWorld = m.gameWorld
	ps.onMap = m
	ps.home = m
	ps.oid = m.allocateOID()
	if m.objects[constant.ObjectTypePersonalShop] == nil {
		m.objects[constant.ObjectTypePersonalShop] = make(map[uint32]Object)
	}
	m.objects[constant.ObjectTypePersonalShop][ps.oid] = ps
}

func (m *Map) RemovePersonalShop(ps *PersonalShop) {
	if m.objects[constant.ObjectTypePersonalShop][ps.oid] != ps {
		return
	}

	delete(m.objects[constant.ObjectTypePersonalShop], ps.oid)
	m.releaseOID(ps.oid)
	ps.onMap = nil
}

func (m *Map) checkShopSpot(position types.Vector2[int16]) error {
	near := func(x, y int16) bool {
		dx := int(position.X) - int(x)
		dy := int(position.Y) - int(y)
		return dx*dx+dy*dy < ShopSpacingSq
	}
	for _, portal := range m.wz.Portals {
		if portal.Type == shopPortalType && near(portal.Position.X, portal.Position.Y) {
			return &MiniRoomEnterError{Code: pconst.MiniRoomEnterNearPortal}
		}
	}
	for _, typ := range []constant.ObjectType{constant.ObjectTypeEntrustedShop, constant.ObjectTypePersonalShop, constant.ObjectTypeMiniGame} {
		for _, obj := range m.objects[typ] {
			other := obj.GetPosition()
			if near(other.X, other.Y) {
				return &MiniRoomEnterError{Code: pconst.MiniRoomEnterCannotOpen}
			}
		}
	}
	return nil
}
