package server

import (
	"fmt"
	"log"

	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/core/luax"
	"github.com/boyism80/fm/game/client"
	"github.com/boyism80/fm/game/constant"
	"github.com/boyism80/fm/game/entity"
	"github.com/boyism80/fm/game/wz"
	"github.com/boyism80/fm/protocol/dto"
	"github.com/boyism80/fm/protocol/request"
	lua "github.com/yuin/gopher-lua"
)

type Attack struct {
	gs     *GameServer
	opcode byte
}

func (Attack) New(gs *GameServer) *Attack {
	return &Attack{
		gs:     gs,
		opcode: 0x1B,
	}
}

func (h *Attack) GetOpcode() byte {
	return h.opcode
}

func (h *Attack) Handle(ctx *core.ClientContext, req *request.Attack) error {
	client, ok := ctx.Client.(*client.GameClient)
	if !ok {
		log.Printf("Client is not a GameClient")
		return fmt.Errorf("client is not a GameClient")
	}

	character := client.GetCharacter()
	if character == nil {
		log.Printf("Character is nil for client")
		return fmt.Errorf("character is nil")
	}

	mapInstance := character.GetMap()
	if mapInstance == nil {
		log.Printf("Character is not in a map")
		return fmt.Errorf("character is not in a map")
	}

	var skillLevel uint8 = 0
	if req.AttackInfo.Skill != 0 {
		if !h.validateAndConsumeSkill(character, req.AttackInfo.Skill) {
			return nil
		}
		skillLevel = uint8(character.GetTotalSkillLevel(req.AttackInfo.Skill))
	}

	h.callOnAttackScript(ctx, character, mapInstance, req.AttackInfo.Damages, req.AttackInfo.Skill, false, 0)
	h.applyDamageToMobs(character, mapInstance, req.AttackInfo.Damages)
	character.Listener.OnAttack(character, req.AttackInfo, skillLevel)

	return nil
}

func (h *Attack) callOnAttackScript(ctx *core.ClientContext, character *entity.Character, mapInstance *entity.Map, damages []dto.AttackPair, skillID uint32, ranged bool, consumeSlot uint16) {
	if ctx.LogicActorPID == nil {
		return
	}
	root := luax.GetRootLuaState(ctx.LogicActorPID.String())
	if root == nil {
		return
	}

	thread, err := luax.NewThread(root, "script/script.lua")
	if err != nil {
		log.Printf("Failed to load script: %v", err)
		return
	}
	defer thread.Close()

	f := thread.GetGlobal("on_attack")
	if f.Type() != lua.LTFunction {
		return
	}

	var skillLV lua.LValue = lua.LNil
	if skillID != 0 {
		if skillEntry := character.Skills[skillID]; skillEntry != nil {
			skillLV = luax.NewLuable(thread, skillEntry)
		}
	}

	damagesTable := buildDamagesTable(thread, mapInstance, damages)
	attackInfoTable := buildAttackInfoTable(thread, ranged, consumeSlot)
	thread.Push(f)
	thread.Push(luax.NewLuable(thread, character))
	thread.Push(skillLV)
	thread.Push(damagesTable)
	thread.Push(attackInfoTable)
	if err := thread.PCall(4, 1, nil); err != nil {
		log.Printf("Failed to call script on_attack: %v", err)
		return
	}
	thread.Pop(1)
}

func buildAttackInfoTable(L *lua.LState, ranged bool, consumeSlot uint16) *lua.LTable {
	tbl := L.NewTable()
	tbl.RawSetString("ranged", lua.LBool(ranged))
	tbl.RawSetString("consume_slot", lua.LNumber(consumeSlot))
	return tbl
}

// buildDamagesTable builds a Lua table: key = mob, value = array of damage amounts per hit.
// In Lua: for mob, hits in damages do ... for _, amount in ipairs(hits) do
func buildDamagesTable(L *lua.LState, mapInstance *entity.Map, damages []dto.AttackPair) *lua.LTable {
	tbl := L.NewTable()
	for _, ap := range damages {
		mob := mapInstance.GetMob(ap.OID)
		if mob == nil {
			continue
		}
		hits := L.NewTable()
		for i, dp := range ap.DamagePairs {
			hits.RawSetInt(i+1, lua.LNumber(dp.Damage))
		}
		tbl.RawSet(luax.NewLuable(L, mob), hits)
	}
	return tbl
}

func (h *Attack) validateAndConsumeSkill(character *entity.Character, skillID uint32) bool {
	var wzSkill *wz.Skill
	if character.Context != nil {
		resources := character.Context.GetResources()
		if resources != nil {
			wzSkill = resources.GetSkill(skillID)
		}
	}

	if wzSkill == nil {
		log.Printf("Skill not found: %d", skillID)
		return false
	}

	skillLevel := character.GetTotalSkillLevel(skillID)
	if skillLevel <= 0 {
		log.Printf("Character does not have skill %d or skill level is 0", skillID)
		return false
	}

	levelData := wzSkill.GetLevelData(skillLevel)
	if levelData == nil {
		return false
	}

	if levelData.Cooldown > 0 {
		skillEntry := character.Skills[skillID]
		if skillEntry == nil || skillEntry.IsCooling() {
			log.Printf("Skill %d is on cooldown", skillID)
			return false
		}
		skillEntry.StartCooldown(levelData.Cooldown)
	}

	if levelData.MPCon > 0 {
		mpCon := uint16(levelData.MPCon)
		if !character.ConsumeMP(mpCon) {
			log.Printf("Not enough MP for skill %d (required: %d, current: %d)", skillID, mpCon, character.Mp)
			return false
		}
	}

	return true
}

func (h *Attack) applyDamageToMobs(character *entity.Character, mapInstance *entity.Map, damages []dto.AttackPair) {
	for _, damage := range damages {
		mob := mapInstance.GetMob(damage.OID)
		if mob == nil {
			log.Printf("Mob not found for OID: %d", damage.OID)
			continue
		}

		if character.HasRoleAtLeast(constant.RoleAdmin) {
			mob.Damage(mob.Hp, character)
		} else {
			for _, damagePair := range damage.DamagePairs {
				mob.Damage(uint16(damagePair.Damage), character)
			}
		}
	}
}
