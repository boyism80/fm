package server

import (
	"fmt"
	"log"

	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/core/luax"
	"github.com/boyism80/fm/game/client"
	"github.com/boyism80/fm/game/constant"
	"github.com/boyism80/fm/game/entity"
	"github.com/boyism80/fm/protocol/request"
	lua "github.com/yuin/gopher-lua"
)

// Damaged handles damage packet requests
type Damaged struct {
	gs     *GameServer
	opcode byte
}

func (Damaged) New(gs *GameServer) *Damaged {
	return &Damaged{
		gs:     gs,
		opcode: 0x1F,
	}
}

func (h *Damaged) GetOpcode() byte {
	return h.opcode
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

	stats := map[constant.Stat]int32{}
	isBlock := req.Damage == -1 && req.Type == constant.IncomingHitCollide
	var damage int32
	if isBlock {
		h.callOnBlocked(ctx, character, req)
		damage = 0
	} else {
		damage = h.resolveDamageByScript(ctx, character, req, req.Damage)
	}

	if !character.Invincible {
		n := int(character.Hp) - int(damage)
		if n < 0 {
			n = 0
		}
		maxHp := int(character.GetMaxHp())
		if n > maxHp {
			n = maxHp
		}

		character.Hp = uint32(n)
		stats[constant.STAT_HP] = int32(character.Hp)
	}

	character.Listener.OnUpdateStats(character, stats, true)

	return nil
}

func (h *Damaged) resolveDamageByScript(ctx *core.ClientContext, character *entity.Character, req *request.Damaged, damage int32) int32 {
	if ctx.LogicActorPID == nil {
		return damage
	}
	root := luax.GetRootLuaState(ctx.LogicActorPID.String())
	if root == nil {
		return damage
	}

	attackerArg := h.resolveDamageAttackerArg(character, req)
	skillArg := h.resolveDamageSkillArg(character, req)

	params := root.NewTable()
	params.RawSetString("hit_type", lua.LNumber(int32(req.Type)))
	params.RawSetString("reflect_ratio", lua.LNumber(int32(req.Reflect)))

	result, thread, err := luax.Call(root, "script/script.lua", "on_damaged", character, attackerArg, skillArg, damage, params)
	if thread != nil {
		thread.Close()
	}
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
	if ctx.LogicActorPID == nil {
		return
	}
	root := luax.GetRootLuaState(ctx.LogicActorPID.String())
	if root == nil {
		return
	}
	attackerArg := h.resolveDamageAttackerArg(character, req)
	_, thread, err := luax.Call(root, "script/script.lua", "on_blocked", character, attackerArg)
	if thread != nil {
		thread.Close()
	}
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
	if character.Context == nil || req.SkillID == 0 {
		return lua.LNil
	}
	resources := character.Context.GetResources()
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
