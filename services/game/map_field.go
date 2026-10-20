package game

import (
	"time"

	"github.com/boyism80/fm/core/clock"
	"github.com/boyism80/fm/maple/constant"
)

const songDuration = 5 * time.Minute

type mapSong struct {
	itemID uint32
	name   string
	until  time.Time
}

func (m *Map) ChangeMusic(song string) {
	if m == nil || song == "" || m.listener == nil {
		return
	}
	m.listener.OnMusicChanged(m, song)
}

func (m *Map) MapMessage(messageType constant.ServerMessageType, message string) {
	if message == "" {
		return
	}
	m.listener.OnMapMessage(m, messageType, message)
}

func (m *Map) Weather(itemID uint32, message string) {
	m.listener.OnWeather(m, itemID, message)
}

func (m *Map) ShowOXQuiz(ask bool, set uint8, id uint16) {
	m.listener.OnOXQuiz(m, ask, set, id)
}

// The client has no way to stop a song, so a map plays one at a time until it ends.
func (m *Map) PlaySong(itemID uint32, name string) bool {
	now := clock.Now()
	if now.Before(m.song.until) {
		return false
	}
	m.song = mapSong{itemID: itemID, name: name, until: now.Add(songDuration)}
	m.listener.OnSongPlayed(m, itemID, name)
	return true
}

func (m *Map) Song() (uint32, string, bool) {
	if !clock.Now().Before(m.song.until) {
		return 0, "", false
	}
	return m.song.itemID, m.song.name, true
}

func (m *Map) YellowChat(message string) {
	m.listener.OnYellowChat(m, message)
}

func (m *Map) WeddingEffect(groomID uint32, brideID uint32) {
	m.listener.OnWeddingEffect(m, groomID, brideID)
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

func (m *Map) Tremble(typ uint8, delay int32) {
	m.listener.OnTremble(m, typ, delay)
}

func (m *Map) PlayersInArea(index int) int {
	if m == nil || m.wz == nil || index < 0 || index >= len(m.wz.Areas) {
		return 0
	}
	area := m.wz.Areas[index]
	count := 0
	for _, obj := range m.GetAllPlayers() {
		ch, ok := obj.(*Character)
		if !ok || ch == nil {
			continue
		}
		if area.ContainsPoint(ch.position) {
			count++
		}
	}
	return count
}
