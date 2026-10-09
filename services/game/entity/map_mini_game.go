package entity

import (
	"github.com/boyism80/fm/services/game/constant"
)

func (m *Map) AddMiniGame(g *MiniGame) {
	g.GameWorld = m.GameWorld
	g.Map = m
	g.OID = m.allocateOID()
	if m.objects[constant.ObjectTypeMiniGame] == nil {
		m.objects[constant.ObjectTypeMiniGame] = make(map[uint32]Object)
	}
	m.objects[constant.ObjectTypeMiniGame][g.OID] = g
}

func (m *Map) RemoveMiniGame(g *MiniGame) {
	if m.objects[constant.ObjectTypeMiniGame][g.OID] != g {
		return
	}

	delete(m.objects[constant.ObjectTypeMiniGame], g.OID)
	m.releaseOID(g.OID)
	g.Map = nil
}
