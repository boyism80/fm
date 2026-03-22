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
	for _, obj := range mapData.GetMobs() {
		mob, ok := obj.(*entity.Mob)
		if !ok || mob == nil || !mob.IsAlive() {
			continue
		}
		if !mob.HasMobStatus(constant.MobStatusPoison) {
			continue
		}
		tick := mob.GetMobStatusValue(constant.MobStatusPoison)
		if tick <= 0 {
			continue
		}
		hp := mob.GetHp()
		if hp <= 1 {
			continue
		}
		damage := uint32(tick)
		if damage >= hp {
			damage = hp - 1
		}
		if damage == 0 {
			continue
		}
		causerID := mob.GetCauserCharacterID(constant.MobStatusPoison)
		var attacker *entity.Character
		if causerID != 0 {
			attacker = mapData.GetPlayer(causerID)
		}
		mob.ApplyDamage(attacker, damage)
	}
	return nil
}
