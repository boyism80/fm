package entity

import (
	"errors"

	"github.com/boyism80/fm/services/game/constant"
)

var ErrMessageBoxSpot = errors.New("message box: too close to another object")

func (m *Map) AddMessageBox(box *MessageBox) error {
	for _, typ := range []constant.ObjectType{constant.ObjectTypeNpc, constant.ObjectTypeMessageBox, constant.ObjectTypeEntrustedShop, constant.ObjectTypePersonalShop} {
		for _, obj := range m.objects[typ] {
			other := obj.GetPosition()
			dx := int(box.Position.X) - int(other.X)
			dy := int(box.Position.Y) - int(other.Y)
			if dx*dx+dy*dy < messageBoxSpacingSq {
				return ErrMessageBoxSpot
			}
		}
	}

	box.ObjectCore.self = box
	box.GameWorld = m.GameWorld
	box.Map = m
	box.OID = m.allocateOID()
	box.initTimers()
	if m.objects[constant.ObjectTypeMessageBox] == nil {
		m.objects[constant.ObjectTypeMessageBox] = make(map[uint32]Object)
	}
	m.objects[constant.ObjectTypeMessageBox][box.OID] = box
	m.sections.add(box)
	m.listener.OnMessageBoxSpawned(m, box)
	box.AddTimer("expire", MessageBoxLifetime, false, func() {
		m.RemoveMessageBox(box)
	})
	return nil
}

func (m *Map) RemoveMessageBox(box *MessageBox) {
	if m.objects[constant.ObjectTypeMessageBox][box.OID] != box {
		return
	}

	m.listener.OnMessageBoxRemoved(m, box)
	box.ClearTimers()
	m.sections.remove(box)
	delete(m.objects[constant.ObjectTypeMessageBox], box.OID)
	m.releaseOID(box.OID)
	box.Map = nil
}
