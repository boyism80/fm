package server

import (
	"fmt"
	"log"

	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/core/luax"
	"github.com/boyism80/fm/protocol/dto"
	"github.com/boyism80/fm/protocol/request"
	"github.com/boyism80/fm/services/game/client"
	"github.com/boyism80/fm/services/game/constant"
	"github.com/boyism80/fm/services/game/entity"
	"github.com/boyism80/fm/services/game/wz"
	lua "github.com/yuin/gopher-lua"
)

type Attack struct {
	gs *GameServer
}

func (Attack) New(gs *GameServer) *Attack {
	return &Attack{
		gs: gs,
	}
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
	skillID := req.Skill
	if skillID != 0 {
		if !h.validateSkillForAttack(character, skillID) {
			return nil
		}
		skillLevel = uint8(character.GetTotalSkillLevel(skillID))
		if !CallSkillHook(ctx, character, skillID, "on_activating") {
			character.Listener.OnUpdateStats(character, nil, true)
			return nil
		}
	}

	damages := req.Damages
	CallOnAttackHooks(ctx, character, mapInstance, damages, skillID, false, 0)
	ApplyDamageToMobs(character, mapInstance, damages)

	if len(req.MesoOIDs) > 0 {
		items := mapInstance.GetItems()
		for _, oid := range req.MesoOIDs {
			obj, exists := items[oid]
			if !exists {
				log.Printf("invalid meso explosion oid: %d", oid)
				return nil
			}
			if _, ok := obj.(*entity.Meso); !ok {
				log.Printf("non-meso object used in meso explosion: %d", oid)
				return nil
			}
			if err := mapInstance.RemoveItem(oid, constant.REMOVE_ITEM_TYPE_EXPLOSION, character.GetID()); err != nil {
				log.Printf("failed to remove meso oid %d for explosion: %v", oid, err)
				return nil
			}
		}
	}
	character.Listener.OnAttack(character, req, skillLevel)
	if skillID != 0 {
		CallSkillHook(ctx, character, skillID, "on_activated")
	}

	return nil
}

func buildAttackInfoTable(L *lua.LState, ranged bool, consumeSlot uint16) *lua.LTable {
	tbl := L.NewTable()
	tbl.RawSetString("ranged", lua.LBool(ranged))
	tbl.RawSetString("consume_slot", lua.LNumber(consumeSlot))
	return tbl
}

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

func readDamagesFromLuaTableInto(damagesTable *lua.LTable, damages []dto.AttackPair) {
	oidToPairs := make(map[uint32][]dto.DamagePair)
	damagesTable.ForEach(func(key lua.LValue, value lua.LValue) {
		ud, ok := key.(*lua.LUserData)
		if !ok || ud.Value == nil {
			return
		}
		mob, ok := ud.Value.(*entity.Mob)
		if !ok {
			return
		}
		oid := mob.GetOID()
		hitsTbl, ok := value.(*lua.LTable)
		if !ok {
			return
		}
		n := hitsTbl.Len()
		pairs := make([]dto.DamagePair, 0, n)
		for i := 1; i <= n; i++ {
			lv := hitsTbl.RawGetInt(i)
			if num, ok := lv.(lua.LNumber); ok {
				pairs = append(pairs, dto.DamagePair{Damage: uint32(num), Unknown: false})
			}
		}
		oidToPairs[oid] = pairs
	})
	for i := range damages {
		if pairs, ok := oidToPairs[damages[i].OID]; ok {
			damages[i].DamagePairs = pairs
		}
	}
}

func (h *Attack) validateSkillForAttack(character *entity.Character, skillID uint32) bool {
	var wzSkill *wz.Skill
	if character.GameWorld != nil {
		resources := character.GameWorld.GetResources()
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
		skillEntry := character.Skills.Get(skillID)
		if skillEntry == nil || skillEntry.IsCooling() {
			log.Printf("Skill %d is on cooldown", skillID)
			return false
		}
		skillEntry.StartCooldown(levelData.Cooldown)
	}

	return true
}
