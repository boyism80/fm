package entity

import (
	"github.com/boyism80/fm/services/game/constant"
)

func (m *Map) ChangeMusic(song string) {
	if m == nil || song == "" || m.listener == nil {
		return
	}
	m.listener.OnMusicChanged(m, song)
}

func (m *Map) MapMessage(message string) {
	if m == nil || message == "" || m.listener == nil {
		return
	}
	m.listener.OnMapMessage(m, constant.MsgPinkText, message)
}

func (m *Map) ClearEffect() {
	if m == nil || m.listener == nil {
		return
	}
	m.listener.OnClearEffect(m)
}

func (m *Map) ShowEffect(path string) {
	if m == nil || path == "" || m.listener == nil {
		return
	}
	m.listener.OnShowEffect(m, path)
}

func (m *Map) PlaySound(path string) {
	if m == nil || path == "" || m.listener == nil {
		return
	}
	m.listener.OnPlaySound(m, path)
}

func (m *Map) PlayersInArea(index int) int {
	if m == nil || m.Wz == nil || index < 0 || index >= len(m.Wz.Areas) {
		return 0
	}
	area := m.Wz.Areas[index]
	count := 0
	for _, obj := range m.GetAllPlayers() {
		ch, ok := obj.(*Character)
		if !ok || ch == nil {
			continue
		}
		if area.ContainsPoint(ch.Position) {
			count++
		}
	}
	return count
}
