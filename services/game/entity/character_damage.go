package entity

import (
	"fmt"
	"log"

	"github.com/boyism80/fm/core/luax"
	"github.com/boyism80/fm/protocol/dto"
	"github.com/boyism80/fm/services/game/constant"
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
		for _, damagePair := range damage.DamagePairs {
			if damagePair.Damage == 0 {
				continue
			}
			if !mob.ApplyDamage(ch, damagePair.Damage) {
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
		mapInstance.runMobLuaHook(root, scriptPath, fmt.Sprintf("on_mob_kill_%d", id), ch, group)
	}
	mapInstance.runMobLuaHook(root, constant.CharacterHookScriptPath, "on_mob_kill", ch, all)
}
