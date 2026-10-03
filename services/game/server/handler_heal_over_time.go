package server

import (
	"fmt"
	"log"
	"time"

	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/core/clock"
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

	if healHP > 0 && req.PRate&1 == 1 {
		intervalSec := h.endureHPInterval(character)
		if intervalSec <= 0 {
			return nil
		}
		hpInterval := time.Duration(intervalSec * float64(time.Second))
		if !allowCooldown(&character.LastHeal.HP, hpInterval, now) {
			return nil
		}
		if healMP > 0 && !allowCooldown(&character.LastHeal.MP, healOverTimeMPInterval, now) {
			return nil
		}
		h.heal(ctx, character, healHP, healMP)
		return nil
	}

	if healHP > 0 && !allowCooldown(&character.LastHeal.HP, healOverTimeHPInterval, now) {
		return nil
	}
	if healMP > 0 && !allowCooldown(&character.LastHeal.MP, healOverTimeMPInterval, now) {
		return nil
	}
	h.heal(ctx, character, healHP, healMP)
	return nil
}

func (h *HealOverTime) heal(ctx *core.ClientContext, character *entity.Character, healHP, healMP uint16) {
	mapInstance := character.GetMap()
	if mapInstance == nil {
		return
	}
	root := mapInstance.GetLuaRoot()
	if root == nil {
		return
	}
	thread, err := luax.NewThread(root, constant.CharacterQueryScriptPath)
	if err != nil {
		return
	}
	ret, err := luax.Call(thread, "get_heal_over_time_cap", character)
	if err != nil {
		return
	}

	maxHP := uint16(0)
	maxMP := uint16(0)
	if tbl, ok := ret.(*lua.LTable); ok {
		if v := tbl.RawGetString("max_hp"); v != nil && v.Type() == lua.LTNumber {
			maxHP = uint16(lua.LVAsNumber(v))
		}
		if v := tbl.RawGetString("max_mp"); v != nil && v.Type() == lua.LTNumber {
			maxMP = uint16(lua.LVAsNumber(v))
		}
	}
	if healHP > maxHP {
		healHP = maxHP
	}
	if healMP > maxMP {
		healMP = maxMP
	}
	if healHP > 0 {
		newHP := character.GetHp() + uint32(healHP)
		if maxHp := character.GetMaxHp(); newHP > maxHp {
			newHP = maxHp
		}
		character.SetHp(newHP, false)
	}
	if healMP > 0 {
		newMP := character.GetMp() + uint32(healMP)
		if maxMp := character.GetMaxMp(); newMP > maxMp {
			newMP = maxMp
		}
		character.SetMp(newMP, false)
	}
	if healHP > 0 || healMP > 0 {
		character.Listener.OnUpdateStats(character, map[constant.Stat]int32{
			constant.StatHP: int32(character.GetHp()),
			constant.StatMP: int32(character.GetMp()),
		}, false)
	}
}

func (h *HealOverTime) endureHPInterval(character *entity.Character) float64 {
	mapInstance := character.GetMap()
	if mapInstance == nil {
		return 0
	}
	root := mapInstance.GetLuaRoot()
	if root == nil {
		return 0
	}
	thread, err := luax.NewThread(root, constant.CharacterQueryScriptPath)
	if err != nil {
		return 0
	}
	ret, err := luax.Call(thread, "get_endure_hp_interval", character)
	if err != nil || ret == nil || ret.Type() != lua.LTNumber {
		return 0
	}
	return float64(lua.LVAsNumber(ret))
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
