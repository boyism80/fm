package server

import (
	"fmt"
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
	gs     *GameServer
	opcode byte
}

func (HealOverTime) New(gs *GameServer) *HealOverTime {
	return &HealOverTime{
		gs:     gs,
		opcode: 0x48,
	}
}

func (h *HealOverTime) GetOpcode() byte {
	return h.opcode
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

	now := time.Now()
	if healHP > 0 {
		hpInterval := healOverTimeHPInterval
		isHanging := req.PRate&1 == 1
		if isHanging {
			intervalSec := h.getEndureHPInterval(ctx, character)
			if intervalSec <= 0 {
				return nil
			}
			hpInterval = time.Duration(intervalSec * float64(time.Second))
		}
		if !allowCooldown(&character.LastHealHPTime, hpInterval, now) {
			return nil
		}
	}
	if healMP > 0 && !allowCooldown(&character.LastHealMPTime, healOverTimeMPInterval, now) {
		return nil
	}

	maxHP, maxMP := h.getHealCap(ctx, character)
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
			constant.STAT_HP: int32(character.GetHp()),
			constant.STAT_MP: int32(character.GetMp()),
		}
		character.Listener.OnUpdateStats(character, stats, false)
	}

	return nil
}

func (h *HealOverTime) getHealCap(ctx *core.ClientContext, character *entity.Character) (uint16, uint16) {
	mapInstance := character.GetMap()
	if mapInstance == nil {
		return 0, 0
	}
	root := mapInstance.GetLuaRoot()
	if root == nil {
		return 0, 0
	}
	result, thread, err := luax.Call(root, "script/script.lua", "get_heal_over_time_cap", character)
	if thread != nil {
		thread.Close()
	}
	if err != nil || result == nil || result.Type() != lua.LTTable {
		return 0, 0
	}
	tbl := result.(*lua.LTable)
	maxHP := uint16(0)
	maxMP := uint16(0)
	if v := tbl.RawGetString("max_hp"); v != nil && v.Type() == lua.LTNumber {
		maxHP = uint16(lua.LVAsNumber(v))
	}
	if v := tbl.RawGetString("max_mp"); v != nil && v.Type() == lua.LTNumber {
		maxMP = uint16(lua.LVAsNumber(v))
	}
	return maxHP, maxMP
}

func (h *HealOverTime) getEndureHPInterval(ctx *core.ClientContext, character *entity.Character) float64 {
	mapInstance := character.GetMap()
	if mapInstance == nil {
		return 0
	}
	root := mapInstance.GetLuaRoot()
	if root == nil {
		return 0
	}
	result, thread, err := luax.Call(root, "script/script.lua", "get_endure_hp_interval", character)
	if thread != nil {
		thread.Close()
	}
	if err != nil || result == nil || result.Type() != lua.LTNumber {
		return 0
	}
	return float64(lua.LVAsNumber(result))
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
