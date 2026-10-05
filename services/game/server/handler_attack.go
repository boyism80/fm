package server

import (
	"fmt"
	"log"

	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/core/luax"
	"github.com/boyism80/fm/protocol/dto"
	"github.com/boyism80/fm/protocol/request"
	"github.com/boyism80/fm/services/game/client"
	"github.com/boyism80/fm/services/game/entity"
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

	if character.IsAlive() == false {
		return nil
	}

	if character.GetMap() == nil {
		log.Printf("Character is not in a map")
		return fmt.Errorf("character is not in a map")
	}

	var skillLevel uint8 = 0
	var skill *entity.SkillEntry
	if req.Skill != 0 {
		skill = character.Skills.Get(req.Skill)
		activated := character.UseAttackSkill(req.Skill, func() bool {
			return skill != nil && character.CallSkillHook(ctx.ActorContext, skill, "on_activating")
		})
		if activated == false {
			character.Listener.OnUpdateStats(character, nil, true)
			return nil
		}
		skillLevel = uint8(character.GetTotalSkillLevel(req.Skill))
	}

	CallOnAttackHooks(ctx.ActorContext, character, req.Damages, skill, false, false, 0)
	character.DamageTo(req.Damages)
	character.ExplodeMesos(req.MesoOIDs)
	character.Listener.OnAttack(character, req, skillLevel)
	if skill != nil {
		character.CallSkillHook(ctx.ActorContext, skill, "on_activated")
	}
	return nil
}

func buildDamagesTable(L *lua.LState, character *entity.Character, damages []dto.AttackPair) *lua.LTable {
	tbl := L.NewTable()
	if character == nil {
		return tbl
	}
	mapInstance := character.GetMap()
	if mapInstance == nil {
		return tbl
	}
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
