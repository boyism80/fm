package entity

import (
	"time"

	"github.com/boyism80/fm/services/game/constant"
)

const mobTimerRemoveAfterKey = "mob:removeAfter"

func (m *Mob) armRemoveAfter() {
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
		m.expireRemoveAfter()
	})
}

func (m *Mob) expireRemoveAfter() {
	if m == nil || m.GetHp() == 0 {
		return
	}

	m.RemoveTimer(mobTimerRemoveAfterKey)
	// removeAfterDieAnimation()로 처리하면 핑크빈에서 석상 스폰 이펙트가 발생안함
	m.Kill(nil, constant.MobDieAnimationTypeFadeOut)
}
