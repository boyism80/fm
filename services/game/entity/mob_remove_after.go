package entity

import (
	"time"

	"github.com/boyism80/fm/services/game/constant"
)

const mobTimerRemoveAfterKey = "mob:removeAfter"

func (m *Mob) scheduleRemoveAfter() {
	if m == nil || m.Wz == nil {
		return
	}
	if m.IsFake() {
		return
	}
	removeAfter := m.Wz.RemoveAfter
	if removeAfter <= 0 {
		return
	}

	delay := time.Duration(removeAfter) * time.Second
	m.RemoveTimer(mobTimerRemoveAfterKey)
	m.AddTimer(mobTimerRemoveAfterKey, delay, false, func() {
		if m == nil || m.GetHp() == 0 {
			return
		}
		m.RemoveTimer(mobTimerRemoveAfterKey)
		m.Kill(nil, constant.MobDieAnimationTypeFadeOut)
	})
}
