package timers

import (
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/game/constant"
	"github.com/boyism80/fm/game/entity"
)

type MobPoisonTickTimer struct{}

func (*MobPoisonTickTimer) New() *MobPoisonTickTimer {
	return &MobPoisonTickTimer{}
}

func (t *MobPoisonTickTimer) GetName() string {
	return "MobPoisonTick"
}

func (t *MobPoisonTickTimer) GetInterval() time.Duration {
	return 1 * time.Second
}

func (t *MobPoisonTickTimer) GetInitialDelay() time.Duration {
	return 1 * time.Second
}

func (t *MobPoisonTickTimer) Handle(ctx actor.Context, mapData *entity.Map) error {
	if mapData.GetPlayerCount() == 0 {
		return nil
	}
	for _, obj := range mapData.GetMobs() {
		mob, ok := obj.(*entity.Mob)
		if !ok || mob == nil || !mob.IsAlive() {
			continue
		}
		applyMobStatusDotDamage(mapData, mob, constant.MobStatusPoison)
		applyMobStatusDotDamage(mapData, mob, constant.MobStatusVenom)
	}
	return nil
}

func applyMobStatusDotDamage(mapData *entity.Map, mob *entity.Mob, status constant.MobStatus) {
	if !mob.HasMobStatus(status) {
		return
	}
	tick := mob.GetMobStatusValue(status)
	if tick <= 0 {
		return
	}
	hp := mob.GetHp()
	if hp <= 1 {
		return
	}
	damage := uint32(tick)
	if damage >= hp {
		damage = hp - 1
	}
	if damage == 0 {
		return
	}
	causerID := mob.GetCauserCharacterID(status)
	var attacker *entity.Character
	if causerID != 0 {
		attacker = mapData.GetPlayer(causerID)
	}
	mob.ApplyDamage(attacker, damage)
}
