package game

import (
	"github.com/boyism80/fm/maple/constant"
)

func (m *Map) AddMiniGame(g *MiniGame) {
	g.gameWorld = m.gameWorld
	g.onMap = m
	g.oid = m.allocateOID()
	if m.objects[constant.ObjectTypeMiniGame] == nil {
		m.objects[constant.ObjectTypeMiniGame] = make(map[uint32]Object)
	}
	m.objects[constant.ObjectTypeMiniGame][g.oid] = g
}

func (m *Map) RemoveMiniGame(g *MiniGame) {
	if m.objects[constant.ObjectTypeMiniGame][g.oid] != g {
		return
	}

	delete(m.objects[constant.ObjectTypeMiniGame], g.oid)
	m.releaseOID(g.oid)
	g.onMap = nil
}
