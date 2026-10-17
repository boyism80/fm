package game

import (
	"log"
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/core/clock"
	"github.com/boyism80/fm/core/luax"
	"github.com/boyism80/fm/maple/constant"
	"github.com/boyism80/fm/maple/wz"
	"github.com/boyism80/fm/protocol/response"
	"github.com/boyism80/fm/types"
)

const (
	mistPoisonDelay    = 2 * time.Second
	mistPoisonInterval = 2500 * time.Millisecond
)

type Mist struct {
	objectCore
	causer               uint32
	skillWz              *wz.Skill
	mobLevelData         *wz.MobSkillLevelData
	skillLevel           uint8
	mistType             constant.MistType
	mobMist              bool
	carnivalTeam         constant.CarnivalTeam
	skillDelay           uint16
	bounds               types.Rect[int32]
	expiresAt            time.Time
	nextPoisonTickAt     time.Time
	poisonTickMultiplier float64
}

type mistSpec struct {
	mistType             constant.MistType
	bounds               types.Rect[int32]
	duration             time.Duration
	poisonDelay          time.Duration
	poisonTickMultiplier float64
}

func (mist *Mist) GetObjectType() constant.ObjectType {
	return constant.ObjectTypeMist
}

func (mist *Mist) Is(typ constant.ObjectType) bool {
	return mist.GetObjectType().Has(typ)
}

func (mist *Mist) IsExpired(now time.Time) bool {
	return !mist.expiresAt.IsZero() && !now.Before(mist.expiresAt)
}

func (mist *Mist) spread(m *Map, spec mistSpec) {
	now := clock.Now()
	mist.objectCore.self = mist
	mist.mistType = spec.mistType
	mist.bounds = spec.bounds
	mist.poisonTickMultiplier = spec.poisonTickMultiplier
	mist.nextPoisonTickAt = now.Add(spec.poisonDelay)
	if spec.duration > 0 {
		mist.expiresAt = now.Add(spec.duration)
	}
	m.AddMist(mist)
}

func (mist *Mist) Tick(actx actor.Context, now time.Time) {
	m := mist.onMap
	if mist.IsExpired(now) {
		m.RemoveMist(mist.oid)
		return
	}
	if mist.mistType != constant.MistTypePoison || m.GetPlayerCount() == 0 || now.Before(mist.nextPoisonTickAt) {
		return
	}
	mist.nextPoisonTickAt = now.Add(mistPoisonInterval)

	targets := mist.poisonTargets()
	if len(targets) == 0 {
		return
	}
	root := m.GetLuaRoot()
	if root == nil {
		return
	}
	thread, err := luax.NewThread(root, constant.CharacterHookScriptPath)
	if err != nil {
		log.Printf("on_poison: %v", err)
		return
	}
	luax.CallAsync(actx, root, thread, "on_poison", mist, targets).OnError(func(err error) {
		log.Printf("on_poison: %v", err)
	})
}

func (mist *Mist) poisonTargets() []luax.Luable {
	targets := make([]luax.Luable, 0)
	if mist.mobMist {
		for _, obj := range mist.onMap.GetObjectsIn(constant.ObjectTypeCharacter, mist.bounds) {
			ch, ok := obj.(*Character)
			if !ok || !ch.IsAlive() || ch.debuffs.Has(constant.DebuffFlagPoison) {
				continue
			}
			targets = append(targets, ch)
		}
		return targets
	}
	for _, obj := range mist.onMap.GetMobs() {
		mob, ok := obj.(*Mob)
		if !ok || !mob.IsAlive() || mob.buffs.Has(constant.MobBuffPoison) {
			continue
		}
		if mob.carnivalTeam != constant.CarnivalTeamNone && mob.carnivalTeam == mist.carnivalTeam {
			continue
		}
		pos := mob.GetPosition()
		if !mist.bounds.ContainsPoint(types.Point[int32]{X: int32(pos.X), Y: int32(pos.Y)}) {
			continue
		}
		targets = append(targets, mob)
	}
	return targets
}

func (mist *Mist) SendSpawnSyncToViewer(viewer *Character) {
	if mist == nil || viewer == nil {
		return
	}
	skillID := uint32(0)
	if mist.skillWz != nil {
		skillID = mist.skillWz.ID
	}
	viewer.Send(&response.SpawnMist{
		OID:        mist.oid,
		Type:       mist.mistType,
		MobMist:    mist.mobMist,
		CauserID:   mist.causer,
		SkillID:    skillID,
		SkillLevel: mist.skillLevel,
		SkillDelay: mist.skillDelay,
		Bounds:     mist.bounds,
		MobSkill:   mist.mobMist,
	}, types.SEND_POLICY_ENCRYPT)
}

func (mist *Mist) SendDestroySyncToViewer(viewer *Character) {
	if mist == nil || viewer == nil {
		return
	}
	viewer.Send(&response.RemoveMist{
		OID:      mist.oid,
		Eruption: false,
	}, types.SEND_POLICY_ENCRYPT)
}
