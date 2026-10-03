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

	damage := int32(0)
	if req.Damage == -1 && req.Type == constant.IncomingHitCollide {
		h.callOnBlocked(ctx, character, req)
	} else {
		damage = h.adjustDamage(character, req)
	}

	if character.Invincible {
		character.Listener.OnUpdateStats(character, map[constant.Stat]int32{}, true)
		return nil
	}
	character.TakeDamage(damage)
	return nil
}

func (h *Damaged) adjustDamage(character *entity.Character, req *request.Damaged) int32 {
	damage := req.Damage
	mapInstance := character.GetMap()
	if mapInstance == nil {
		return damage
	}
	if req.OID != 0 {
		if mob := mapInstance.GetMob(req.OID); mob != nil && mob.IsFake() {
			return 0
		}
	}
	root := mapInstance.GetLuaRoot()
	if root == nil {
		return damage
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
		return damage
	}
	ret, err := luax.Call(thread, "on_damaged", character, attackerArg, skillArg, damage, params)
	if err != nil {
		log.Printf("Failed to call script on_damaged: %v", err)
		return damage
	}
	if ret == nil || ret.Type() != lua.LTNumber {
		return damage
	}
	adjusted := int32(lua.LVAsNumber(ret))
	if adjusted < 0 {
		return 0
	}
	return adjusted
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
	luax.CallAsync(ctx.ActorContext, root, thread, "on_blocked", character, attackerArg).OnError(func(err error) {
		log.Printf("Failed to call script on_blocked: %v", err)
	})
}
