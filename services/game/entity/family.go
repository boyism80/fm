package entity

import (
	"encoding/json"
	"log"
	"time"

	"github.com/asynkron/protoactor-go/actor"
	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
)

type familyMoveDetail struct {
	Complete bool `json:"complete,omitempty"`
}

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

// A paid move ends by refund until the character has arrived; after that it ends by completion.
func (f *Family) FinishMove(actx actor.Context, id uint64, complete bool) {
	request := f.owner.Unsettled.Find(id)
	if request == nil {
		return
	}
	if complete {
		data, _ := json.Marshal(familyMoveDetail{Complete: true})
		request.Detail = string(data)
	}
	f.finishMove(actx, request, complete)
}

func (f *Family) settleMove(actx actor.Context, request *UnsettledRequest) {
	var detail familyMoveDetail
	if err := json.Unmarshal([]byte(request.Detail), &detail); err != nil {
		log.Printf("Family.settleMove character=%d id=%d: %v", f.owner.GetID(), request.ID, err)
		return
	}
	f.finishMove(actx, request, detail.Complete)
}

func (f *Family) finishMove(actx actor.Context, request *UnsettledRequest, complete bool) {
	f.owner.Listener.FinishFamilyMoveAsync(actx, f.owner, request.ID, complete).Do(func(reply *internal.FamilyReply) error {
		f.owner.Unsettled.Erase(request.ID)
		f.Set(reply.Family)
		return nil
	}).OnError(func(err error) {
		request.InFlight = false
		log.Printf("Family.finishMove character=%d id=%d: %v", f.owner.GetID(), request.ID, err)
	})
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
