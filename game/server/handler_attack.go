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

	mapID := character.Map
	mapInstance := h.gs.GetMap(mapID)
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

	h.callOnAttackScript(ctx, character, mapInstance, req.AttackInfo.Damages, req.AttackInfo.Skill)
	h.applyDamageToMobs(character, mapInstance, req.AttackInfo.Damages)
	character.Listener.OnAttack(mapID, character.GetID(), req.AttackInfo, skillLevel)

	return nil
}

func (h *Attack) callOnAttackScript(ctx *core.ClientContext, character *entity.Character, mapInstance *entity.Map, damages []dto.AttackPair, skillID uint32) {
	if ctx.LogicActorPID == nil {
		return
	}
	root := luax.GetRootLuaState(ctx.LogicActorPID.String())
	if root == nil {
		return
	}

	targets := make([]luax.Luable, 0, len(damages))
	seen := make(map[uint32]bool)
	for _, damage := range damages {
		if seen[damage.OID] {
			continue
		}
		seen[damage.OID] = true
		mob := mapInstance.GetMob(damage.OID)
		if mob != nil {
			targets = append(targets, mob)
		}
	}

	var skillArg interface{} = lua.LNil
	if skillID != 0 {
		if skillEntry := character.Skills[skillID]; skillEntry != nil {
			skillArg = skillEntry
		}
	}

	_, thread, err := luax.Call(root, "script/script.lua", "on_attack", character, targets, skillArg)
	if err != nil {
		log.Printf("Failed to call script on_attack: %v", err)
		return
	}
	if thread != nil {
		thread.Close()
	}
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
