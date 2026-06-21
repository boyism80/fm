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
	var damage int32
	if isBlock {
		h.callOnBlocked(ctx, character, req)
		damage = 0
	} else {
		damage = h.resolveDamageByScript(ctx, character, req, req.Damage)
	}

	if !character.Invincible {
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
	} else {
		character.Listener.OnUpdateStats(character, map[constant.Stat]int32{}, true)
	}

	return nil
}

func (h *Damaged) resolveDamageByScript(ctx *core.ClientContext, character *entity.Character, req *request.Damaged, damage int32) int32 {
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

	attackerArg := h.resolveDamageAttackerArg(character, req)
	skillArg := h.resolveDamageSkillArg(character, req)

	params := root.NewTable()
	params.RawSetString("hit_type", lua.LNumber(int32(req.Type)))
	params.RawSetString("reflect_ratio", lua.LNumber(int32(req.Reflect)))

	thread, err := luax.NewThread(root, constant.CharacterHookScriptPath)
	if err != nil {
		log.Printf("Failed to call script on_damaged: %v", err)
		return damage
	}
	result, err := luax.Call(thread, "on_damaged", character, attackerArg, skillArg, damage, params)
	if err != nil {
		log.Printf("Failed to call script on_damaged: %v", err)
		return damage
	}
	if result == nil || result.Type() != lua.LTNumber {
		return damage
	}

	adjustedDamage := int32(lua.LVAsNumber(result))
	if adjustedDamage < 0 {
		return 0
	}
	return adjustedDamage
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
	attackerArg := h.resolveDamageAttackerArg(character, req)
	thread, err := luax.NewThread(root, constant.CharacterHookScriptPath)
	if err != nil {
		log.Printf("Failed to call script on_blocked: %v", err)
		return
	}
	_, err = luax.Call(thread, "on_blocked", character, attackerArg)
	if err != nil {
		log.Printf("Failed to call script on_blocked: %v", err)
	}
}

func (h *Damaged) resolveDamageAttackerArg(character *entity.Character, req *request.Damaged) interface{} {
	if req.OID == 0 {
		return lua.LNil
	}
	mapInstance := character.GetMap()
	if mapInstance == nil {
		return lua.LNil
	}
	mob := mapInstance.GetMob(req.OID)
	if mob == nil {
		return lua.LNil
	}
	return mob
}

func (h *Damaged) resolveDamageSkillArg(character *entity.Character, req *request.Damaged) interface{} {
	if character.GameWorld == nil || req.SkillID == 0 {
		return lua.LNil
	}
	resources := character.GameWorld.GetResources()
	if resources == nil {
		return lua.LNil
	}
	skillID := uint32(req.SkillID)
	skillModel := resources.GetSkill(skillID)
	if skillModel == nil {
		return lua.LNil
	}
	return entity.NewSkillEntry(character, skillModel, int(req.Level), 0)
}
