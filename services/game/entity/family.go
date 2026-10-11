package entity

import (
	"time"

	"github.com/asynkron/protoactor-go/actor"
	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
)

type Family struct {
	owner      *Character
	State      *internal.FamilyState
	Progress   uint64
	Reputation map[uint32]uint64
	Pending    bool
}

func (f *Family) Set(state *internal.FamilyState) {
	if state == nil {
		return
	}

	f.State = state
	f.owner.Listener.OnFamilyChanged(f.owner)
}

func (f *Family) Refresh(ctx actor.Context) {
	f.owner.Listener.LoadFamilyAsync(ctx, f.owner).Do(func(reply *internal.FamilyReply) error {
		if f.owner.LoggedOut() {
			return nil
		}

		f.Set(reply.Family)
		return nil
	}).OnError(func(error) {})
}

func (f *Family) BuffRemaining() time.Duration {
	if f.State == nil {
		return 0
	}
	return time.Until(time.UnixMilli(f.State.BuffUntil))
}

func (f *Family) ExpRate() int {
	if f.BuffRemaining() <= 0 {
		return 100
	}

	switch f.State.BuffFeature {
	case internal.FamilyFeature_FAMILY_FEATURE_EXP_15_MINUTES, internal.FamilyFeature_FAMILY_FEATURE_EXP_30_MINUTES:
		return 120
	case internal.FamilyFeature_FAMILY_FEATURE_UNITY:
		return 150
	}
	return 100
}

func (f *Family) DropRate() int {
	if f.BuffRemaining() <= 0 {
		return 100
	}

	switch f.State.BuffFeature {
	case internal.FamilyFeature_FAMILY_FEATURE_DROP_15_MINUTES, internal.FamilyFeature_FAMILY_FEATURE_DROP_30_MINUTES:
		return 120
	case internal.FamilyFeature_FAMILY_FEATURE_UNITY:
		return 150
	}
	return 100
}

func (f *Family) gainProgress(exp uint32) {
	if f.State == nil || f.State.SeniorId == 0 || f.owner.GameWorld == nil {
		return
	}

	resources := f.owner.GameWorld.GetResources()
	if resources == nil {
		return
	}

	level := f.owner.level
	old := uint64(f.owner.exp)
	remaining := uint64(exp)
	for remaining > 0 && level < 200 {
		needed := uint64(resources.GetExpNeededForLevel(level))
		if needed == 0 || old >= needed {
			return
		}

		step := min(remaining, needed-old)
		// Whole percentage crossings survive small gains and level changes.
		units := (old+step)*100/needed - old*100/needed
		if old+step == needed {
			units += uint64(level+1) / 2
		}
		f.Progress += units

		seniorID := f.State.SeniorId
		for _, multiplier := range []uint64{10, 5} {
			for _, member := range f.State.Members {
				if member.CharacterId != seniorID {
					continue
				}

				amount := units * multiplier
				if uint32(level) > member.Level {
					amount /= 2
				}
				f.Reputation[seniorID] += amount
				seniorID = member.SeniorId
				break
			}
		}
		remaining -= step
		old += step
		if old == needed {
			level++
			old = 0
		}
	}
}
