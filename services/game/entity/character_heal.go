package entity

import (
	"time"

	"github.com/boyism80/fm/core/luax"
	"github.com/boyism80/fm/services/game/constant"
	lua "github.com/yuin/gopher-lua"
)

const (
	healOverTimeHPInterval = 10 * time.Second
	healOverTimeMPInterval = 8 * time.Second
)

type lastHeal struct {
	hp time.Time
	mp time.Time
}

func (ch *Character) HealOverTime(hp, mp uint16, endure bool, now time.Time) {
	if ch.GetHp() == 0 {
		return
	}
	m := ch.GetMap()
	if m == nil || m.GetLuaRoot() == nil {
		return
	}
	thread, err := luax.NewThread(m.GetLuaRoot(), constant.CharacterQueryScriptPath)
	if err != nil {
		return
	}
	ret, err := luax.Call(thread, "get_heal_over_time_cap", ch)
	if err != nil {
		return
	}
	limit, ok := ret.(*lua.LTable)
	if ok == false {
		return
	}

	hpInterval := healOverTimeHPInterval
	if endure && hp > 0 {
		seconds := float64(lua.LVAsNumber(limit.RawGetString("endure_hp_interval")))
		if seconds <= 0 {
			return
		}
		hpInterval = time.Duration(seconds * float64(time.Second))
	}
	if hp > 0 && now.Sub(ch.lastHeal.hp) < hpInterval {
		return
	}
	if mp > 0 && now.Sub(ch.lastHeal.mp) < healOverTimeMPInterval {
		return
	}
	if hp > 0 {
		ch.lastHeal.hp = now
	}
	if mp > 0 {
		ch.lastHeal.mp = now
	}

	hp = min(hp, uint16(lua.LVAsNumber(limit.RawGetString("max_hp"))))
	mp = min(mp, uint16(lua.LVAsNumber(limit.RawGetString("max_mp"))))
	if hp == 0 && mp == 0 {
		return
	}
	ch.AddHpMp(int(hp), int(mp))
}
