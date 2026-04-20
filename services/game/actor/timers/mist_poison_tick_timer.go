package timers

import (
	"log"
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/core/luax"
	"github.com/boyism80/fm/services/game/constant"
	"github.com/boyism80/fm/services/game/entity"
	"github.com/boyism80/fm/types"
)

type MistPoisonTickTimer struct{}

func (*MistPoisonTickTimer) New() *MistPoisonTickTimer {
	return &MistPoisonTickTimer{}
}

func (t *MistPoisonTickTimer) GetName() string {
	return "MistPoisonTick"
}

func (t *MistPoisonTickTimer) GetInterval() time.Duration {
	return 500 * time.Millisecond
}

func (t *MistPoisonTickTimer) GetInitialDelay() time.Duration {
	return 500 * time.Millisecond
}

func (t *MistPoisonTickTimer) Handle(ctx actor.Context, mapData *entity.Map) error {
	if mapData.GetPlayerCount() == 0 {
		return nil
	}
	now := time.Now()
	objects := mapData.GetObjects(constant.ObjectTypeMist)
	for _, obj := range objects {
		mist, ok := obj.(*entity.Mist)
		if !ok || mist == nil || mist.OID == 0 {
			continue
		}
		if mist.MistType != constant.MistTypePoison {
			continue
		}
		if mist.NextPoisonTickAt.IsZero() {
			mist.NextPoisonTickAt = now.Add(2 * time.Second)
			continue
		}
		if now.Before(mist.NextPoisonTickAt) {
			continue
		}
		candidates := make([]luax.Luable, 0)
		for _, mobObj := range mapData.GetMobs() {
			mob, ok := mobObj.(*entity.Mob)
			if !ok || mob == nil || !mob.IsAlive() {
				continue
			}
			if mob.HasBuff(constant.MobBuffPoison) {
				continue
			}
			pos := mob.GetPosition()
			if !mist.Bounds.ContainsPoint(types.Point[int32]{X: int32(pos.X), Y: int32(pos.Y)}) {
				continue
			}
			candidates = append(candidates, mob)
		}
		if len(candidates) > 0 {
			root := mapData.GetLuaRoot()
			if root != nil {
				thread, err := luax.NewThread(root, "script/script.lua")
				if err == nil {
					_, err = luax.Call(thread, "on_poison", mist, candidates)
				}
				if err != nil {
					log.Printf("on_poison failed: %v", err)
				}
			} else {
				log.Printf("on_poison skipped: root Lua state not found")
			}
		}
		mist.NextPoisonTickAt = now.Add(2500 * time.Millisecond)
	}
	return nil
}
