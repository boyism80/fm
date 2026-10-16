package game

import (
	"time"

	"github.com/boyism80/fm/maple/constant"
)

const mobTimerRemoveAfterKey = "mob:removeAfter"

func (m *Mob) scheduleRemoveAfter() {
	if m == nil || m.wz == nil {
		return
	}
	if m.fake {
		return
	}
	removeAfter := m.wz.RemoveAfter
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
