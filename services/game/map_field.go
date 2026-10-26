package game

import (
	"math/rand"
	"time"

	"github.com/boyism80/fm/core/clock"
	"github.com/boyism80/fm/maple/constant"
	"github.com/boyism80/fm/maple/wz"
)

const songDuration = 5 * time.Minute

// A started weather stays this long; the client keeps it until the map clears it.
const weatherDuration = 30 * time.Second

type mapSong struct {
	itemID uint32
	name   string
	until  time.Time
}

type mapWeather struct {
	itemID  uint32
	message string
	until   time.Time
}

const healerRiseInterval = time.Second

type mapHealer struct {
	y      int16
	rising bool
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

// StartWeather shows a weather message for weatherDuration, to the map and to whoever enters meanwhile.
func (m *Map) StartWeather(itemID uint32, message string) bool {
	now := clock.Now()
	if now.Before(m.weather.until) {
		return false
	}
	until := now.Add(weatherDuration)
	m.weather = mapWeather{itemID: itemID, message: message, until: until}
	m.listener.OnWeather(m, itemID, message)
	m.After(weatherDuration, func() {
		if m.weather.until != until {
			return
		}
		m.weather = mapWeather{}
		m.listener.OnWeather(m, 0, "")
	})
	return true
}

// GiveStateChange gives everyone on the map the consume item effect a weather item names.
func (m *Map) GiveStateChange(itemID uint32) {
	model, ok := m.gameWorld.GetResources().Items[itemID].(*wz.CashItem)
	if !ok || model.StateChangeItem == 0 {
		return
	}
	for _, obj := range m.objects[constant.ObjectTypeCharacter] {
		ch, ok := obj.(*Character)
		if !ok {
			continue
		}
		item, err := NewItem(model.StateChangeItem, 1, m.gameWorld)
		if err != nil {
			return
		}
		if consume, ok := item.(*Consume); ok {
			ch.UseConsume(consume)
		}
	}
}

func (m *Map) CurrentWeather() (uint32, string, bool) {
	if !clock.Now().Before(m.weather.until) {
		return 0, "", false
	}
	return m.weather.itemID, m.weather.message, true
}

// PullPulley lowers the healer; at the bottom it heals everyone in the field and goes back to the top.
func (m *Map) PullPulley() {
	healer := m.wz.Healer
	if healer == nil {
		return
	}

	m.healer.y += healer.Fall
	if m.healer.y >= healer.YMax {
		amount := healer.HealMin + rand.Intn(healer.HealMax-healer.HealMin+1)
		for _, obj := range m.objects[constant.ObjectTypeCharacter] {
			ch, ok := obj.(*Character)
			if !ok || ch.GetHp() == 0 {
				continue
			}
			ch.AddHp(amount)
		}
		m.healer.y = healer.YMin
	}
	m.listener.OnHealerMove(m, m.healer.y)
	m.scheduleHealerRise()
}

func (m *Map) scheduleHealerRise() {
	if m.healer.rising || m.healer.y <= m.wz.Healer.YMin {
		return
	}
	m.healer.rising = true
	m.After(healerRiseInterval, m.riseHealer)
}

func (m *Map) riseHealer() {
	m.healer.rising = false
	healer := m.wz.Healer
	if m.healer.y <= healer.YMin {
		return
	}
	m.healer.y = max(healer.YMin, m.healer.y-healer.Rise)
	m.listener.OnHealerMove(m, m.healer.y)
	m.scheduleHealerRise()
}

// HealerY reports where a lowered healer hangs, for players who enter meanwhile.
func (m *Map) HealerY() (int16, bool) {
	if m.wz.Healer == nil || m.healer.y <= m.wz.Healer.YMin {
		return 0, false
	}
	return m.healer.y, true
}

func (m *Map) ShowOXQuiz(ask bool, set uint8, id uint16) {
	m.listener.OnOXQuiz(m, ask, set, id)
}

func (m *Map) ShowCoconutHit(id int16, delay int16, state uint8) {
	m.listener.OnCoconutHit(m, id, delay, state)
}

func (m *Map) ShowCoconutScore(maple int16, story int16) {
	m.listener.OnCoconutScore(m, maple, story)
}

func (m *Map) ShowSnowballState(state uint8, snowmanHP [2]int32, position [2]int16, stage [2]uint8, damage [3]int16) {
	m.listener.OnSnowballState(m, state, snowmanHP, position, stage, damage)
}

func (m *Map) ShowSnowballHit(target uint8, damage int16, delay int16) {
	m.listener.OnSnowballHit(m, target, damage, delay)
}

func (m *Map) ShowSnowballMessage(team uint8, message uint8) {
	m.listener.OnSnowballMessage(m, team, message)
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
