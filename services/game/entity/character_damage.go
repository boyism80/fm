package entity

import (
	"fmt"
	"log"

	"github.com/boyism80/fm/core/luax"
	"github.com/boyism80/fm/protocol/dto"
	"github.com/boyism80/fm/services/game/constant"
	"github.com/boyism80/fm/types"
)

func (ch *Character) DamageTo(damages []dto.AttackPair) {
	if ch == nil || len(damages) == 0 {
		return
	}
	mapInstance := ch.GetMap()
	if mapInstance == nil {
		return
	}

	var killed []*Mob
	seen := make(map[uint32]struct{})
	for _, damage := range damages {
		mob := mapInstance.GetMob(damage.OID)
		if mob == nil {
			log.Printf("Mob not found for OID: %d", damage.OID)
			continue
		}
		total := uint64(0)
		for _, damagePair := range damage.DamagePairs {
			total += uint64(damagePair.Damage)
		}
		if total > 0 {
			ch.chargeDojoEnergy(mob, total)
		}
		for _, damagePair := range damage.DamagePairs {
			if damagePair.Damage == 0 {
				continue
			}
			if !mob.TakeDamage(ch, damagePair.Damage) {
				continue
			}
			if _, ok := seen[mob.OID]; ok {
				continue
			}
			seen[mob.OID] = struct{}{}
			killed = append(killed, mob)
		}
	}
	if len(killed) > 0 {
		ch.OnKill(killed)
	}
}

const mesoExplosionMargin = 50

func (ch *Character) ExplodeMesos(oids []uint32) {
	if len(oids) == 0 {
		return
	}
	mapInstance := ch.GetMap()
	if mapInstance == nil {
		return
	}
	skillID := uint32(constant.SkillMesoExplosion)
	skill := ch.GameWorld.GetResources().GetSkill(skillID)
	if skill == nil {
		return
	}
	levelData := skill.GetLevelData(ch.GetTotalSkillLevel(skillID))
	if levelData == nil {
		return
	}

	reach := max(levelData.LT.X, -levelData.LT.X, levelData.RB.X, -levelData.RB.X) + mesoExplosionMargin
	origin := types.Point[int32]{X: int32(ch.Position.X), Y: int32(ch.Position.Y)}
	area := types.Rect[int32]{
		Left:   -reach,
		Top:    min(levelData.LT.Y, levelData.RB.Y) - mesoExplosionMargin,
		Right:  reach,
		Bottom: max(levelData.LT.Y, levelData.RB.Y) + mesoExplosionMargin,
	}.AtOrigin(origin)

	items := mapInstance.GetItems()
	for _, oid := range oids {
		meso, ok := items[oid].(*Meso)
		if ok == false {
			continue
		}
		fp := meso.GetFieldPlacement()
		if fp == nil || fp.Owner != ch.GetID() {
			continue
		}
		pos := meso.GetPosition()
		if area.ContainsPoint(types.Point[int32]{X: int32(pos.X), Y: int32(pos.Y)}) == false {
			continue
		}
		if err := mapInstance.RemoveItem(oid, constant.RemoveItemTypeExplosion, ch.GetID()); err != nil {
			log.Printf("failed to remove meso oid %d for explosion: %v", oid, err)
		}
	}
}

func (ch *Character) TakeDamage(damage int32) {
	if damage > 0 {
		ch.AddDojoEnergy(2)
	}

	hp := int(ch.GetHp()) - int(damage)
	if hp < 0 {
		hp = 0
	}
	if maxHp := int(ch.GetMaxHp()); hp > maxHp {
		hp = maxHp
	}
	ch.SetHp(uint32(hp), false)
	ch.Listener.OnUpdateStats(ch, map[constant.Stat]int32{
		constant.StatHP: int32(ch.GetHp()),
	}, true)
}

func (ch *Character) OnKill(mobs []*Mob) {
	if ch == nil || len(mobs) == 0 {
		return
	}
	mapInstance := ch.GetMap()
	if mapInstance == nil {
		return
	}
	root := mapInstance.GetLuaRoot()
	if root == nil {
		return
	}

	byID := make(map[uint32][]luax.Luable)
	order := make([]uint32, 0)
	all := make([]luax.Luable, 0, len(mobs))
	for _, mob := range mobs {
		if mob == nil || mob.Wz == nil {
			continue
		}
		id := mob.Wz.ID
		if _, ok := byID[id]; !ok {
			order = append(order, id)
		}
		byID[id] = append(byID[id], mob)
		all = append(all, mob)
	}
	if len(all) == 0 {
		return
	}

	for _, id := range order {
		group := byID[id]
		scriptPath := fmt.Sprintf("script/mob/%d.lua", id)
		mapInstance.runMobLuaHook(root, scriptPath, "on_mob_kill", ch, group)
	}
	mapInstance.runMobLuaHook(root, constant.CharacterHookScriptPath, "on_mob_kill", ch, all)
	if sm := ch.StateMachine(); sm != nil {
		sm.CallHook("on_mob_kill", ch, all)
	}
}
