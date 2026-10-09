package entity

import (
	"github.com/boyism80/fm/services/game/constant"
)

func (m *Map) AddTrade(t *Trade) {
	t.GameWorld = m.GameWorld
	t.Map = m
	t.OID = m.allocateOID()
	if m.objects[constant.ObjectTypeTrade] == nil {
		m.objects[constant.ObjectTypeTrade] = make(map[uint32]Object)
	}
	m.objects[constant.ObjectTypeTrade][t.OID] = t
}

func (m *Map) RemoveTrade(t *Trade) {
	if m.objects[constant.ObjectTypeTrade][t.OID] != t {
		return
	}

	delete(m.objects[constant.ObjectTypeTrade], t.OID)
	m.releaseOID(t.OID)
	t.Map = nil
}
