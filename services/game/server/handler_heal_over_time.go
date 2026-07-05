package server

import (
	"fmt"
	"github.com/boyism80/fm/core/clock"
	"log"
	"time"

	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/core/luax"
	"github.com/boyism80/fm/protocol/request"
	"github.com/boyism80/fm/services/game/client"
	"github.com/boyism80/fm/services/game/constant"
	"github.com/boyism80/fm/services/game/entity"
	lua "github.com/yuin/gopher-lua"
)

const (
	healOverTimeHPInterval = 10 * time.Second
	healOverTimeMPInterval = 8 * time.Second
)

type HealOverTime struct {
	gs *GameServer
}

func (HealOverTime) New(gs *GameServer) *HealOverTime {
	return &HealOverTime{
		gs: gs,
	}
}

func (h *HealOverTime) Handle(ctx *core.ClientContext, req *request.HealOverTime) error {
	client, ok := ctx.Client.(*client.GameClient)
	if !ok {
		log.Printf("Client is not a GameClient")
		return fmt.Errorf("client is not a GameClient")
	}

	character := client.GetCharacter()
	if character == nil {
		return nil
	}

	if character.GetHp() <= 0 {
		return nil
	}

	healHP := req.HealHP
	healMP := req.HealMP

	now := clock.Now()
	if healHP > 0 {
		hpInterval := healOverTimeHPInterval
		isHanging := req.PRate&1 == 1
		if isHanging {
			h.getEndureHPInterval(ctx, character, func(intervalSec float64) {
				if intervalSec <= 0 {
					return
				}
				hpInterval = time.Duration(intervalSec * float64(time.Second))
				h.applyHealOverTime(ctx, character, req, healHP, healMP, hpInterval, now)
			})
			return nil
		}
		if !allowCooldown(&character.LastHealHPTime, hpInterval, now) {
			return nil
		}
	}
	if healMP > 0 && !allowCooldown(&character.LastHealMPTime, healOverTimeMPInterval, now) {
		return nil
	}

	h.getHealCap(ctx, character, func(maxHP, maxMP uint16) {
		h.applyHealAmounts(character, healHP, healMP, maxHP, maxMP)
	})
	return nil
}

func (h *HealOverTime) applyHealOverTime(ctx *core.ClientContext, character *entity.Character, req *request.HealOverTime, healHP, healMP uint16, hpInterval time.Duration, now time.Time) {
	if !allowCooldown(&character.LastHealHPTime, hpInterval, now) {
		return
	}
	if healMP > 0 && !allowCooldown(&character.LastHealMPTime, healOverTimeMPInterval, now) {
		return
	}
	h.getHealCap(ctx, character, func(maxHP, maxMP uint16) {
		h.applyHealAmounts(character, healHP, healMP, maxHP, maxMP)
	})
}

func (h *HealOverTime) applyHealAmounts(character *entity.Character, healHP, healMP, maxHP, maxMP uint16) {
	if healHP > maxHP {
		healHP = maxHP
	}
	if healMP > maxMP {
		healMP = maxMP
	}

	if healHP > 0 {
		newHP := character.GetHp() + uint32(healHP)
		maxHp := character.GetMaxHp()
		if newHP > maxHp {
			newHP = maxHp
		}
		character.SetHp(newHP, false)
	}

	if healMP > 0 {
		newMP := character.GetMp() + uint32(healMP)
		maxMp := character.GetMaxMp()
		if newMP > maxMp {
			newMP = maxMp
		}
		character.SetMp(newMP, false)
	}

	if healHP > 0 || healMP > 0 {
		stats := map[constant.Stat]int32{
			constant.StatHP: int32(character.GetHp()),
			constant.StatMP: int32(character.GetMp()),
		}
		character.Listener.OnUpdateStats(character, stats, false)
	}
}

func (h *HealOverTime) getHealCap(ctx *core.ClientContext, character *entity.Character, fn func(uint16, uint16)) {
	mapInstance := character.GetMap()
	if mapInstance == nil {
		fn(0, 0)
		return
	}
	root := mapInstance.GetLuaRoot()
	if root == nil {
		fn(0, 0)
		return
	}
	thread, err := luax.NewThread(root, constant.CharacterQueryScriptPath)
	if err != nil {
		fn(0, 0)
		return
	}
	luax.CallAsync(root, thread, "get_heal_over_time_cap", character).Then(func(value interface{}) (interface{}, error) {
		vals := luax.ResultValues(value)
		if len(vals) == 0 || vals[0] == nil || vals[0].Type() != lua.LTTable {
			fn(0, 0)
			return nil, nil
		}
		tbl := vals[0].(*lua.LTable)
		maxHP := uint16(0)
		maxMP := uint16(0)
		if v := tbl.RawGetString("max_hp"); v != nil && v.Type() == lua.LTNumber {
			maxHP = uint16(lua.LVAsNumber(v))
		}
		if v := tbl.RawGetString("max_mp"); v != nil && v.Type() == lua.LTNumber {
			maxMP = uint16(lua.LVAsNumber(v))
		}
		fn(maxHP, maxMP)
		return nil, nil
	}).OnError(func(err error) {
		fn(0, 0)
	})
}

func (h *HealOverTime) getEndureHPInterval(ctx *core.ClientContext, character *entity.Character, fn func(float64)) {
	mapInstance := character.GetMap()
	if mapInstance == nil {
		fn(0)
		return
	}
	root := mapInstance.GetLuaRoot()
	if root == nil {
		fn(0)
		return
	}
	thread, err := luax.NewThread(root, constant.CharacterQueryScriptPath)
	if err != nil {
		fn(0)
		return
	}
	luax.CallAsync(root, thread, "get_endure_hp_interval", character).Then(func(value interface{}) (interface{}, error) {
		vals := luax.ResultValues(value)
		if len(vals) == 0 || vals[0] == nil || vals[0].Type() != lua.LTNumber {
			fn(0)
			return nil, nil
		}
		fn(float64(lua.LVAsNumber(vals[0])))
		return nil, nil
	}).OnError(func(err error) {
		fn(0)
	})
}

func allowCooldown(last *time.Time, interval time.Duration, now time.Time) bool {
	if last == nil {
		return false
	}
	if last.IsZero() || now.Sub(*last) >= interval {
		*last = now
		return true
	}
	return false
}
