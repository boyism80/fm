package entity

import (
	"encoding/json"
	"log"
	"slices"
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/core/clock"
	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
	"github.com/boyism80/fm/services/game/constant"
)

type familyMoveDetail struct {
	Complete bool `json:"complete,omitempty"`
}

// familyCall is an invite or a summon waiting for the owner's answer.
type familyCall struct {
	from      uint32
	name      string
	expiresAt time.Time
}

func (c familyCall) waiting(from uint32, name string) bool {
	return c.from != 0 && c.from == from && c.name == name && clock.Now().Before(c.expiresAt)
}

type Family struct {
	owner      *Character
	Tree       FamilyTree
	Progress   uint64
	Reputation map[uint32]uint64
	Pending    bool
	invite     familyCall
	summoned   familyCall
}

func (f *Family) Set(family *internal.Family) {
	if family == nil {
		return
	}

	f.Tree = FamilyTree{family: family}
	f.owner.Listener.OnFamilyChanged(f.owner)
}

func (f *Family) Self() *internal.FamilyMember {
	return f.Tree.Member(f.owner.GetID())
}

func (f *Family) fail(result internal.FamilyResult) {
	f.owner.Listener.OnFamilyResult(f.owner, result)
	f.owner.Listener.OnUnlockAction(f.owner)
}

func (f *Family) finishRequest(reply *internal.FamilyReply) error {
	f.Pending = false
	if reply.GetStatus() != internal.FamilyResult_FAMILY_RESULT_OK {
		f.fail(reply.GetStatus())
		return nil
	}

	f.Set(reply.GetFamily())
	f.owner.Listener.OnUnlockAction(f.owner)
	return nil
}

func (f *Family) failRequest(error) {
	f.Pending = false
	f.fail(internal.FamilyResult_FAMILY_RESULT_FAILED)
}

func (f *Family) Open() {
	f.owner.Listener.OnFamilyChanged(f.owner)
	f.owner.Listener.OnUnlockAction(f.owner)
}

func (f *Family) ShowPedigree(actx actor.Context, name string) {
	f.owner.Listener.GetFamilyAsync(actx, f.owner, name).Do(func(reply *internal.GetFamilyReply) error {
		if !reply.GetFound() {
			f.fail(internal.FamilyResult_FAMILY_RESULT_WRONG_NAME)
			return nil
		}

		f.owner.Listener.OnFamilyChart(f.owner, FamilyTree{family: reply.GetFamily()}, reply.GetCharacterId())
		f.owner.Listener.OnUnlockAction(f.owner)
		return nil
	}).OnError(func(error) {
		f.fail(internal.FamilyResult_FAMILY_RESULT_FAILED)
	})
}

// checkJoin is internal's join check, run early on the live characters; the owner is the senior.
func (f *Family) checkJoin(junior *Character) internal.FamilyResult {
	senior := f.owner
	switch {
	case junior == senior || junior.GetLevel() < constant.FamilyMinJuniorLevel:
		return internal.FamilyResult_FAMILY_RESULT_CANNOT_BE_JUNIOR
	case senior.GetLevel() <= junior.GetLevel():
		return internal.FamilyResult_FAMILY_RESULT_JUNIOR_LEVEL_TOO_HIGH
	case senior.GetLevel()-junior.GetLevel() > constant.FamilyMaxLevelGap:
		return internal.FamilyResult_FAMILY_RESULT_LEVEL_GAP_TOO_LARGE
	case junior.Family.Self().GetSeniorId() != 0:
		return internal.FamilyResult_FAMILY_RESULT_ALREADY_JUNIOR
	case len(f.Tree.Juniors(senior.GetID())) >= constant.FamilyMaxJuniors:
		return internal.FamilyResult_FAMILY_RESULT_CANNOT_BE_JUNIOR
	case f.Tree.Leader().GetCharacterId() == junior.GetID():
		return internal.FamilyResult_FAMILY_RESULT_SAME_FAMILY
	}
	return internal.FamilyResult_FAMILY_RESULT_OK
}

// Invite asks a character on the owner's map to become the owner's junior.
func (f *Family) Invite(name string) {
	id, found := f.owner.GameWorld.GetDispatchSystem().FindCharacterID(name)
	if !found {
		f.fail(internal.FamilyResult_FAMILY_RESULT_WRONG_NAME)
		return
	}
	junior := f.owner.GetMap().GetPlayer(id)
	if junior == nil || (junior.IsHidden() && f.owner.Role < junior.Role) {
		f.fail(internal.FamilyResult_FAMILY_RESULT_NOT_SAME_MAP)
		return
	}
	if result := f.checkJoin(junior); result != internal.FamilyResult_FAMILY_RESULT_OK {
		f.fail(result)
		return
	}

	junior.Family.invite = familyCall{from: f.owner.GetID(), name: f.owner.GetName(), expiresAt: clock.Now().Add(constant.FamilyAnswerDuration)}
	junior.Listener.OnFamilyInvite(junior, f.owner.GetID(), f.owner.GetName())
	f.owner.Listener.OnUnlockAction(f.owner)
}

func (f *Family) AnswerInvite(actx actor.Context, seniorID uint32, name string, accepted bool) {
	if f.Pending {
		return
	}
	invite := f.invite
	f.invite = familyCall{}
	if !invite.waiting(seniorID, name) {
		f.fail(internal.FamilyResult_FAMILY_RESULT_FAILED)
		return
	}

	senior := f.owner.GetMap().GetPlayer(seniorID)
	if !accepted {
		if senior != nil {
			senior.Listener.OnFamilyInviteResult(senior, false, f.owner.GetName())
		}
		f.owner.Listener.OnUnlockAction(f.owner)
		return
	}
	if senior == nil {
		f.fail(internal.FamilyResult_FAMILY_RESULT_NOT_SAME_MAP)
		return
	}
	if result := senior.Family.checkJoin(f.owner); result != internal.FamilyResult_FAMILY_RESULT_OK {
		f.fail(result)
		return
	}

	f.Pending = true
	f.owner.Listener.JoinFamilyAsync(actx, senior, f.owner).Do(func(reply *internal.FamilyReply) error {
		f.Pending = false
		if reply.GetStatus() != internal.FamilyResult_FAMILY_RESULT_OK {
			f.fail(reply.GetStatus())
			return nil
		}

		f.Set(reply.GetFamily())
		f.owner.Listener.OnFamilyJoined(f.owner, name)
		if m := f.owner.GetMap(); m != nil {
			if senior := m.GetPlayer(seniorID); senior != nil {
				senior.Listener.OnFamilyInviteResult(senior, true, f.owner.GetName())
			}
		}
		f.owner.Listener.OnUnlockAction(f.owner)
		return nil
	}).OnError(f.failRequest)
}

// Separate takes juniorID away from its senior: the owner removes its junior, or leaves its own senior.
func (f *Family) Separate(actx actor.Context, juniorID uint32) {
	if f.Pending {
		return
	}
	seniorID := f.Tree.Member(juniorID).GetSeniorId()
	if seniorID == 0 || (juniorID != f.owner.GetID() && seniorID != f.owner.GetID()) {
		f.fail(internal.FamilyResult_FAMILY_RESULT_FAILED)
		return
	}

	f.Pending = true
	f.owner.Listener.SeparateFamilyAsync(actx, f.owner, juniorID).Do(f.finishRequest).OnError(f.failRequest)
}

func (f *Family) SetPrecept(actx actor.Context, precept string) {
	if f.Pending {
		return
	}
	if f.Tree.Leader().GetCharacterId() != f.owner.GetID() {
		f.fail(internal.FamilyResult_FAMILY_RESULT_FAILED)
		return
	}

	f.Pending = true
	f.owner.Listener.SetFamilyPreceptAsync(actx, f.owner, precept).Do(f.finishRequest).OnError(f.failRequest)
}

// Use spends the owner's reputation on a feature. Teleport and summon name a
// family member on the same channel.
func (f *Family) Use(actx actor.Context, feature internal.FamilyFeature, name string) {
	if f.Pending {
		return
	}
	self := f.Self()
	if feature < 0 || int(feature) >= len(constant.FamilyCosts) || !f.Tree.InFamily() ||
		self.GetReputation() < constant.FamilyCosts[feature] || slices.Contains(self.GetUsedFeatures(), feature) {
		f.fail(internal.FamilyResult_FAMILY_RESULT_FAILED)
		return
	}
	if feature != internal.FamilyFeature_FAMILY_FEATURE_TELEPORT && feature != internal.FamilyFeature_FAMILY_FEATURE_SUMMON {
		f.Pending = true
		f.owner.Listener.UseFamilyFeatureAsync(actx, f.owner, feature, 0, 0).Do(f.finishRequest).OnError(f.failRequest)
		return
	}

	targetID, found := f.owner.GameWorld.GetDispatchSystem().FindCharacterID(name)
	switch {
	case !found:
		f.fail(internal.FamilyResult_FAMILY_RESULT_WRONG_NAME)
	case targetID == f.owner.GetID():
		f.fail(internal.FamilyResult_FAMILY_RESULT_FAILED)
	case f.Tree.Member(targetID) == nil:
		f.fail(internal.FamilyResult_FAMILY_RESULT_NOT_SAME_FAMILY)
	case f.owner.GetMap().BlocksTeleportStone():
		f.fail(internal.FamilyResult_FAMILY_RESULT_FAILED)
	case feature == internal.FamilyFeature_FAMILY_FEATURE_TELEPORT:
		f.moveTo(internal.FamilyFeature_FAMILY_FEATURE_TELEPORT, targetID)
	default:
		f.summon(targetID)
	}
}

// summon asks the target first; the owner pays only when the target accepts, in AnswerSummon.
func (f *Family) summon(targetID uint32) {
	ownerID, ownerName, role := f.owner.GetID(), f.owner.GetName(), f.owner.Role
	mapName, _ := f.owner.GameWorld.GetResources().GetMapName(f.owner.GetMap().TemplateID())
	dispatch := f.owner.GameWorld.GetDispatchSystem()
	f.owner.Listener.OnUnlockAction(f.owner)

	dispatch.CallCharacter(targetID, func(_ actor.Context, target *Character) {
		result := internal.FamilyResult_FAMILY_RESULT_OK
		switch {
		case target == nil || target.LoggedOut() || target.GetMap() == nil || (target.IsHidden() && role < target.Role):
			result = internal.FamilyResult_FAMILY_RESULT_WRONG_NAME
		case target.GetMap().BlocksTeleportStone():
			result = internal.FamilyResult_FAMILY_RESULT_FAILED
		case target.Family.summoned.from != 0 && clock.Now().Before(target.Family.summoned.expiresAt):
			result = internal.FamilyResult_FAMILY_RESULT_SUMMON_PENDING
		}
		if result == internal.FamilyResult_FAMILY_RESULT_OK {
			target.Family.summoned = familyCall{from: ownerID, name: ownerName, expiresAt: clock.Now().Add(constant.FamilyAnswerDuration)}
			target.Listener.OnFamilySummon(target, ownerName, mapName)
			return
		}

		dispatch.CallCharacter(ownerID, func(_ actor.Context, owner *Character) {
			if owner != nil {
				owner.Listener.OnFamilyResult(owner, result)
			}
		})
	})
}

func (f *Family) AnswerSummon(name string, accepted bool) {
	summoned := f.summoned
	f.summoned = familyCall{}
	if !summoned.waiting(summoned.from, name) {
		f.fail(internal.FamilyResult_FAMILY_RESULT_SUMMON_FAILED)
		return
	}
	if !accepted {
		f.owner.Listener.OnUnlockAction(f.owner)
		return
	}
	if f.owner.GetMap().BlocksTeleportStone() {
		f.fail(internal.FamilyResult_FAMILY_RESULT_FAILED)
		return
	}

	f.moveTo(internal.FamilyFeature_FAMILY_FEATURE_SUMMON, summoned.from)
}

// moveTo reads the target's map on the target's actor, then starts the move
// back on the owner's actor.
func (f *Family) moveTo(feature internal.FamilyFeature, targetID uint32) {
	ownerID, role, origin := f.owner.GetID(), f.owner.Role, f.owner.GetMap()
	dispatch := f.owner.GameWorld.GetDispatchSystem()
	f.owner.Listener.OnUnlockAction(f.owner)

	dispatch.CallCharacter(targetID, func(_ actor.Context, target *Character) {
		if target == nil || target.LoggedOut() || target.GetMap() == nil || (target.IsHidden() && role < target.Role) || target.GetMap().BlocksTeleportStone() {
			dispatch.CallCharacter(ownerID, func(_ actor.Context, owner *Character) {
				if owner != nil {
					owner.Listener.OnFamilyResult(owner, internal.FamilyResult_FAMILY_RESULT_FAILED)
				}
			})
			return
		}

		destination := target.GetMap()
		dispatch.CallCharacter(ownerID, func(c actor.Context, owner *Character) {
			if owner != nil {
				owner.Family.move(c, feature, targetID, destination, origin)
			}
		})
	})
}

// move pays for the move first. A reply that comes too late or finds the owner
// elsewhere refunds the reservation instead of moving.
func (f *Family) move(actx actor.Context, feature internal.FamilyFeature, targetID uint32, destination, origin *Map) {
	if f.Pending || f.owner.GetMap() != origin {
		f.owner.Listener.OnFamilyResult(f.owner, internal.FamilyResult_FAMILY_RESULT_FAILED)
		return
	}

	f.Pending = true
	id := f.owner.Unsettled.Add(internal.UnsettledType_UNSETTLED_TYPE_FAMILY_MOVE, 0, nil)
	deadline := clock.Now().Add(constant.FamilyMoveDuration)
	f.owner.Listener.UseFamilyFeatureAsync(actx, f.owner, feature, targetID, id).Do(func(reply *internal.FamilyReply) error {
		f.Pending = false
		if reply.GetStatus() != internal.FamilyResult_FAMILY_RESULT_OK {
			f.owner.Unsettled.Erase(id)
			f.owner.Listener.OnFamilyResult(f.owner, reply.GetStatus())
			return nil
		}

		f.Set(reply.GetFamily())
		m := f.owner.GetMap()
		if clock.Now().After(deadline) || m != origin || m.BlocksTeleportStone() || destination.BlocksTeleportStone() ||
			f.owner.GetHp() == 0 || f.owner.Dialog.Busy() {
			f.refund(actx, id)
			return nil
		}
		err := f.owner.GameWorld.GetMapSystem().Warp(actx, f.owner, destination, 0, func(entered actor.Context) {
			f.FinishMove(entered, id, true)
		})
		if err != nil {
			f.refund(actx, id)
		}
		return nil
	}).OnError(func(error) {
		f.Pending = false
		f.owner.Unsettled.Find(id).InFlight = false
		f.owner.Listener.OnFamilyResult(f.owner, internal.FamilyResult_FAMILY_RESULT_FAILED)
	})
}

func (f *Family) refund(actx actor.Context, id uint64) {
	f.FinishMove(actx, id, false)
	f.owner.Listener.OnFamilyResult(f.owner, internal.FamilyResult_FAMILY_RESULT_FAILED)
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
		f.Set(reply.GetFamily())
		return nil
	}).OnError(func(err error) {
		request.InFlight = false
		log.Printf("Family.finishMove character=%d id=%d: %v", f.owner.GetID(), request.ID, err)
	})
}

func (f *Family) BuffRemaining() time.Duration {
	until := f.Self().GetBuffUntil()
	if until == 0 {
		return 0
	}
	return time.Until(time.UnixMilli(until))
}

func (f *Family) ExpRate() int {
	if f.BuffRemaining() <= 0 {
		return 100
	}

	switch f.Self().GetBuffFeature() {
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

	switch f.Self().GetBuffFeature() {
	case internal.FamilyFeature_FAMILY_FEATURE_DROP_15_MINUTES, internal.FamilyFeature_FAMILY_FEATURE_DROP_30_MINUTES:
		return 120
	case internal.FamilyFeature_FAMILY_FEATURE_UNITY:
		return 150
	}
	return 100
}

func (f *Family) gainProgress(exp uint32) {
	seniorID := f.Self().GetSeniorId()
	if seniorID == 0 || f.owner.GameWorld == nil {
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

		receiverID := seniorID
		for _, multiplier := range []uint64{10, 5} {
			senior := f.Tree.Member(receiverID)
			if senior == nil {
				break
			}

			amount := units * multiplier
			if uint32(level) > senior.GetLevel() {
				amount /= 2
			}
			f.Reputation[receiverID] += amount
			receiverID = senior.GetSeniorId()
		}
		remaining -= step
		old += step
		if old == needed {
			level++
			old = 0
		}
	}
}
