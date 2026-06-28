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
