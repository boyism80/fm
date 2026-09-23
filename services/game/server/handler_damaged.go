package server

import (
	"fmt"
	"log"

	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/core/luax"
	"github.com/boyism80/fm/protocol/request"
	"github.com/boyism80/fm/services/game/client"
	"github.com/boyism80/fm/services/game/constant"
	"github.com/boyism80/fm/services/game/entity"
	lua "github.com/yuin/gopher-lua"
)

type Damaged struct {
	gs *GameServer
}

func (Damaged) New(gs *GameServer) *Damaged {
	return &Damaged{
		gs: gs,
	}
}

func (h *Damaged) Handle(ctx *core.ClientContext, req *request.Damaged) error {
	client, ok := ctx.Client.(*client.GameClient)
	if !ok {
		log.Printf("Client is not a GameClient")
		return fmt.Errorf("client is not a GameClient")
	}

	character := client.GetCharacter()
	if character == nil {
		log.Printf("Character not found for client")
		return fmt.Errorf("character not found")
	}

	isBlock := req.Damage == -1 && req.Type == constant.IncomingHitCollide
	if isBlock {
		h.callOnBlocked(ctx, character, req)
		h.takeDamage(character, 0)
		return nil
	}
	h.finalizeDamage(ctx, character, req, req.Damage, func(damage int32) {
		h.takeDamage(character, damage)
	})
	return nil
}

func (h *Damaged) takeDamage(character *entity.Character, damage int32) {
	if !character.Invincible {
		wasAlive := character.GetHp() > 0
		n := int(character.GetHp()) - int(damage)
		if n < 0 {
			n = 0
		}
		maxHp := int(character.GetMaxHp())
		if n > maxHp {
			n = maxHp
		}
		character.SetHp(uint32(n), false)
		character.Listener.OnUpdateStats(character, map[constant.Stat]int32{
			constant.StatHP: int32(character.GetHp()),
		}, true)
		if wasAlive && character.GetHp() == 0 {
			if sm := character.StateMachine(); sm != nil {
				sm.CallHook("on_player_dead", sm, character)
			}
		}
	} else {
		character.Listener.OnUpdateStats(character, map[constant.Stat]int32{}, true)
	}
}

func (h *Damaged) finalizeDamage(ctx *core.ClientContext, character *entity.Character, req *request.Damaged, damage int32, fn func(int32)) {
	mapInstance := character.GetMap()
	if mapInstance == nil {
		fn(damage)
		return
	}
	if req.OID != 0 {
		if mob := mapInstance.GetMob(req.OID); mob != nil && mob.IsFake() {
			fn(0)
			return
		}
	}
	root := mapInstance.GetLuaRoot()
	if root == nil {
		fn(damage)
		return
	}

	var attackerArg interface{} = lua.LNil
	if req.OID != 0 {
		if mob := mapInstance.GetMob(req.OID); mob != nil {
			attackerArg = mob
		}
	}
	var skillArg interface{} = lua.LNil
	if character.GameWorld != nil && req.SkillID != 0 {
		if resources := character.GameWorld.GetResources(); resources != nil {
			skillID := uint32(req.SkillID)
			if skillModel := resources.GetSkill(skillID); skillModel != nil {
				skillArg = entity.NewSkillEntry(character, skillModel, int(req.Level), 0)
			}
		}
	}

	params := root.NewTable()
	params.RawSetString("hit_type", lua.LNumber(int32(req.Type)))
	params.RawSetString("reflect_ratio", lua.LNumber(int32(req.Reflect)))

	thread, err := luax.NewThread(root, constant.CharacterHookScriptPath)
	if err != nil {
		log.Printf("Failed to call script on_damaged: %v", err)
		fn(damage)
		return
	}
	luax.CallAsync(root, thread, "on_damaged", character, attackerArg, skillArg, damage, params).Then(func(value interface{}) (interface{}, error) {
		vals := luax.ResultValues(value)
		if len(vals) == 0 || vals[0] == nil || vals[0].Type() != lua.LTNumber {
			fn(damage)
			return nil, nil
		}
		adjustedDamage := int32(lua.LVAsNumber(vals[0]))
		if adjustedDamage < 0 {
			fn(0)
			return nil, nil
		}
		fn(adjustedDamage)
		return nil, nil
	}).OnError(func(err error) {
		log.Printf("Failed to call script on_damaged: %v", err)
		fn(damage)
	})
}

func (h *Damaged) callOnBlocked(ctx *core.ClientContext, character *entity.Character, req *request.Damaged) {
	mapInstance := character.GetMap()
	if mapInstance == nil {
		return
	}
	root := mapInstance.GetLuaRoot()
	if root == nil {
		return
	}
	var attackerArg interface{} = lua.LNil
	if req.OID != 0 {
		if mob := mapInstance.GetMob(req.OID); mob != nil {
			attackerArg = mob
		}
	}
	thread, err := luax.NewThread(root, constant.CharacterHookScriptPath)
	if err != nil {
		log.Printf("Failed to call script on_blocked: %v", err)
		return
	}
	luax.CallAsync(root, thread, "on_blocked", character, attackerArg).OnError(func(err error) {
		log.Printf("Failed to call script on_blocked: %v", err)
	})
}
