package game

import (
	"encoding/json"
	"log"
	"slices"
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/core/clock"
	"github.com/boyism80/fm/maple/constant"
	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
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
	tree       FamilyTree
	progress   uint64
	reputation map[uint32]uint64
	pending    bool
	invite     familyCall
	summoned   familyCall
}

func (f *Family) Set(tree FamilyTree) {
	f.tree = tree
	f.owner.listener.OnFamilyChanged(f.owner)
}

func (f *Family) Self() *FamilyMember {
	return f.tree.Member(f.owner.GetID())
}

func (f *Family) fail(result internal.FamilyResult) {
	f.owner.listener.OnFamilyResult(f.owner, result)
	f.owner.listener.OnUnlockAction(f.owner)
}

func (f *Family) finishRequest(reply *internal.FamilyReply) error {
	f.pending = false
	if reply.GetStatus() != internal.FamilyResult_FAMILY_RESULT_OK {
		f.fail(reply.GetStatus())
		return nil
	}

	f.Set(FamilyTreeFromProto(reply.GetFamily()))
	f.owner.listener.OnUnlockAction(f.owner)
	return nil
}

func (f *Family) failRequest(error) {
	f.pending = false
	f.fail(internal.FamilyResult_FAMILY_RESULT_FAILED)

}

func (f *Family) Open() {
	f.owner.listener.OnFamilyChanged(f.owner)
	f.owner.listener.OnUnlockAction(f.owner)
}

func (f *Family) ShowPedigree(actx actor.Context, name string) {
	f.owner.gameWorld.Systems().Get[*InternalSystem]().GetFamilyAsync(actx, f.owner, name).Bind(f.owner).Do(func(reply *internal.GetFamilyReply) error {
		if !reply.GetFound() {
			f.fail(internal.FamilyResult_FAMILY_RESULT_WRONG_NAME)
			return nil
		}

		f.owner.listener.OnFamilyChart(f.owner, FamilyTreeFromProto(reply.GetFamily()), reply.GetCharacterId())
		f.owner.listener.OnUnlockAction(f.owner)
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
	case junior.family.Self().SeniorID() != 0:
		return internal.FamilyResult_FAMILY_RESULT_ALREADY_JUNIOR
	case len(f.tree.Juniors(senior.GetID())) >= constant.FamilyMaxJuniors:
		return internal.FamilyResult_FAMILY_RESULT_CANNOT_BE_JUNIOR
	case f.tree.Leader().CharacterID() == junior.GetID():
		return internal.FamilyResult_FAMILY_RESULT_SAME_FAMILY
	}
	return internal.FamilyResult_FAMILY_RESULT_OK
}

// Invite asks a character on the owner's map to become the owner's junior.
func (f *Family) Invite(name string) {
	id, found := f.owner.gameWorld.Systems().Get[*DispatchSystem]().FindCharacterID(name)
	if !found {
		f.fail(internal.FamilyResult_FAMILY_RESULT_WRONG_NAME)
		return
	}
	m := f.owner.GetMap()
	if m == nil {
		f.fail(internal.FamilyResult_FAMILY_RESULT_NOT_SAME_MAP)
		return
	}
	junior := m.GetPlayer(id)
	if junior == nil || (junior.IsHidden() && f.owner.role < junior.role) {
		f.fail(internal.FamilyResult_FAMILY_RESULT_NOT_SAME_MAP)
		return
	}
	if result := f.checkJoin(junior); result != internal.FamilyResult_FAMILY_RESULT_OK {
		f.fail(result)
		return
	}

	junior.family.invite = familyCall{from: f.owner.GetID(), name: f.owner.GetName(), expiresAt: clock.Now().Add(constant.FamilyAnswerDuration)}
	junior.listener.OnFamilyInvite(junior, f.owner.GetID(), f.owner.GetName())
	f.owner.listener.OnUnlockAction(f.owner)
}

func (f *Family) AnswerInvite(actx actor.Context, seniorID uint32, name string, accepted bool) {
	if f.pending {
		return
	}
	invite := f.invite
	f.invite = familyCall{}
	if !invite.waiting(seniorID, name) {
		f.fail(internal.FamilyResult_FAMILY_RESULT_FAILED)
		return
	}

	var senior *Character
	if m := f.owner.GetMap(); m != nil {
		senior = m.GetPlayer(seniorID)
	}
	if !accepted {
		if senior != nil {
			senior.listener.OnFamilyInviteResult(senior, false, f.owner.GetName())
		}
		f.owner.listener.OnUnlockAction(f.owner)
		return
	}
	if senior == nil {
		f.fail(internal.FamilyResult_FAMILY_RESULT_NOT_SAME_MAP)
		return
	}
	if result := senior.family.checkJoin(f.owner); result != internal.FamilyResult_FAMILY_RESULT_OK {
		f.fail(result)
		return
	}

	f.pending = true
	f.owner.gameWorld.Systems().Get[*InternalSystem]().JoinFamilyAsync(actx, senior, f.owner).Bind(f.owner).Do(func(reply *internal.FamilyReply) error {
		f.pending = false
		if reply.GetStatus() != internal.FamilyResult_FAMILY_RESULT_OK {
			f.fail(reply.GetStatus())
			return nil
		}

		f.Set(FamilyTreeFromProto(reply.GetFamily()))
		f.owner.listener.OnFamilyJoined(f.owner, name)
		if m := f.owner.GetMap(); m != nil {
			if senior := m.GetPlayer(seniorID); senior != nil {
				senior.listener.OnFamilyInviteResult(senior, true, f.owner.GetName())
			}
		}
		f.owner.listener.OnUnlockAction(f.owner)
		return nil
	}).OnError(f.failRequest)
}

// Separate takes juniorID away from its senior: the owner removes its junior, or leaves its own senior.
func (f *Family) Separate(actx actor.Context, juniorID uint32) {
	if f.pending {
		return
	}
	seniorID := f.tree.Member(juniorID).SeniorID()
	if seniorID == 0 || (juniorID != f.owner.GetID() && seniorID != f.owner.GetID()) {
		f.fail(internal.FamilyResult_FAMILY_RESULT_FAILED)
		return
	}

	f.pending = true
	f.owner.gameWorld.Systems().Get[*InternalSystem]().SeparateFamilyAsync(actx, f.owner, juniorID).Bind(f.owner).Do(f.finishRequest).OnError(f.failRequest)
}

func (f *Family) SetPrecept(actx actor.Context, precept string) {
	if f.pending {
		return
	}
	if f.tree.Leader().CharacterID() != f.owner.GetID() {
		f.fail(internal.FamilyResult_FAMILY_RESULT_FAILED)
		return
	}

	f.pending = true
	f.owner.gameWorld.Systems().Get[*InternalSystem]().SetFamilyPreceptAsync(actx, f.owner, precept).Bind(f.owner).Do(f.finishRequest).OnError(f.failRequest)
}

// Use spends the owner's reputation on a feature. Teleport and summon name a
// family member on the same channel.
func (f *Family) Use(actx actor.Context, feature internal.FamilyFeature, name string) {
	if f.pending {
		return
	}
	self := f.Self()
	if feature < 0 || int(feature) >= len(constant.FamilyCosts) || !f.tree.InFamily() ||
		self.Reputation() < constant.FamilyCosts[feature] || slices.Contains(self.UsedFeatures(), feature) {
		f.fail(internal.FamilyResult_FAMILY_RESULT_FAILED)
		return
	}
	if feature != internal.FamilyFeature_FAMILY_FEATURE_TELEPORT && feature != internal.FamilyFeature_FAMILY_FEATURE_SUMMON {
		f.pending = true
		f.owner.gameWorld.Systems().Get[*InternalSystem]().UseFamilyFeatureAsync(actx, f.owner, feature, 0, 0).Bind(f.owner).Do(f.finishRequest).OnError(f.failRequest)
		return
	}

	targetID, found := f.owner.gameWorld.Systems().Get[*DispatchSystem]().FindCharacterID(name)
	m := f.owner.GetMap()
	switch {
	case !found:
		f.fail(internal.FamilyResult_FAMILY_RESULT_WRONG_NAME)
	case targetID == f.owner.GetID():
		f.fail(internal.FamilyResult_FAMILY_RESULT_FAILED)
	case f.tree.Member(targetID) == nil:
		f.fail(internal.FamilyResult_FAMILY_RESULT_NOT_SAME_FAMILY)
	case m == nil || m.BlocksTeleportStone():
		f.fail(internal.FamilyResult_FAMILY_RESULT_FAILED)
	case feature == internal.FamilyFeature_FAMILY_FEATURE_TELEPORT:
		f.moveTo(internal.FamilyFeature_FAMILY_FEATURE_TELEPORT, targetID)
	default:
		f.summon(targetID)
	}
}

// summon asks the target first; the owner pays only when the target accepts, in AnswerSummon.
func (f *Family) summon(targetID uint32) {
	ownerID, ownerName, role := f.owner.GetID(), f.owner.GetName(), f.owner.role
	mapName, _ := f.owner.gameWorld.GetResources().GetMapName(f.owner.GetMap().TemplateID())
	dispatch := f.owner.gameWorld.Systems().Get[*DispatchSystem]()
	f.owner.listener.OnUnlockAction(f.owner)

	dispatch.Visit(targetID, func(_ actor.Context, target *Character) {
		result := internal.FamilyResult_FAMILY_RESULT_OK
		switch {
		case target == nil || target.LoggedOut() || target.GetMap() == nil || (target.IsHidden() && role < target.role):
			result = internal.FamilyResult_FAMILY_RESULT_WRONG_NAME
		case target.GetMap().BlocksTeleportStone():
			result = internal.FamilyResult_FAMILY_RESULT_FAILED
		case target.family.summoned.from != 0 && clock.Now().Before(target.family.summoned.expiresAt):
			result = internal.FamilyResult_FAMILY_RESULT_SUMMON_PENDING
		}
		if result == internal.FamilyResult_FAMILY_RESULT_OK {
			target.family.summoned = familyCall{from: ownerID, name: ownerName, expiresAt: clock.Now().Add(constant.FamilyAnswerDuration)}
			target.listener.OnFamilySummon(target, ownerName, mapName)
			return
		}

		dispatch.Visit(ownerID, func(_ actor.Context, owner *Character) {
			if owner != nil {
				owner.listener.OnFamilyResult(owner, result)
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
		f.owner.listener.OnUnlockAction(f.owner)
		return
	}
	if m := f.owner.GetMap(); m == nil || m.BlocksTeleportStone() {
		f.fail(internal.FamilyResult_FAMILY_RESULT_FAILED)
		return
	}

	f.moveTo(internal.FamilyFeature_FAMILY_FEATURE_SUMMON, summoned.from)
}

// moveTo reads the target's map on the target's actor, then starts the move
// back on the owner's actor.
func (f *Family) moveTo(feature internal.FamilyFeature, targetID uint32) {
	ownerID, role, origin := f.owner.GetID(), f.owner.role, f.owner.GetMap()
	dispatch := f.owner.gameWorld.Systems().Get[*DispatchSystem]()
	f.owner.listener.OnUnlockAction(f.owner)

	dispatch.Visit(targetID, func(_ actor.Context, target *Character) {
		if target == nil || target.LoggedOut() || target.GetMap() == nil || (target.IsHidden() && role < target.role) || target.GetMap().BlocksTeleportStone() {
			dispatch.Visit(ownerID, func(_ actor.Context, owner *Character) {
				if owner != nil {
					owner.listener.OnFamilyResult(owner, internal.FamilyResult_FAMILY_RESULT_FAILED)
				}
			})
			return
		}

		destination := target.GetMap()
		dispatch.Visit(ownerID, func(c actor.Context, owner *Character) {
			if owner != nil {
				owner.family.move(c, feature, targetID, destination, origin)
			}
		})
	})
}

// move pays for the move first. A reply that comes too late or finds the owner
// elsewhere refunds the reservation instead of moving.
func (f *Family) move(actx actor.Context, feature internal.FamilyFeature, targetID uint32, destination, origin *Map) {
	if f.pending || f.owner.GetMap() != origin {
		f.owner.listener.OnFamilyResult(f.owner, internal.FamilyResult_FAMILY_RESULT_FAILED)
		return
	}

	f.pending = true
	id := f.owner.unsettled.Add(internal.UnsettledType_UNSETTLED_TYPE_FAMILY_MOVE, 0, nil)
	deadline := clock.Now().Add(constant.FamilyMoveDuration)
	f.owner.gameWorld.Systems().Get[*InternalSystem]().UseFamilyFeatureAsync(actx, f.owner, feature, targetID, id).Bind(f.owner).Do(func(reply *internal.FamilyReply) error {
		f.pending = false
		if reply.GetStatus() != internal.FamilyResult_FAMILY_RESULT_OK {
			f.owner.unsettled.Erase(id)
			f.owner.listener.OnFamilyResult(f.owner, reply.GetStatus())
			return nil
		}

		f.Set(FamilyTreeFromProto(reply.GetFamily()))
		m := f.owner.GetMap()
		if clock.Now().After(deadline) || m != origin || m.BlocksTeleportStone() || destination.BlocksTeleportStone() ||
			f.owner.GetHp() == 0 || f.owner.dialog.Busy() {
			f.refund(actx, id)
			return nil
		}
		err := f.owner.gameWorld.Systems().Get[*MapSystem]().Warp(actx, f.owner, destination, 0, func(entered actor.Context) {
			f.FinishMove(entered, id, true)
		})
		if err != nil {
			f.refund(actx, id)
		}
		return nil
	}).OnError(func(error) {
		f.pending = false
		f.owner.unsettled.Find(id).inFlight = false
		f.owner.listener.OnFamilyResult(f.owner, internal.FamilyResult_FAMILY_RESULT_FAILED)

	})
}

func (f *Family) refund(actx actor.Context, id uint64) {
	f.FinishMove(actx, id, false)
	f.owner.listener.OnFamilyResult(f.owner, internal.FamilyResult_FAMILY_RESULT_FAILED)
}

// A paid move ends by refund until the character has arrived; after that it ends by completion.
func (f *Family) FinishMove(actx actor.Context, id uint64, complete bool) {
	request := f.owner.unsettled.Find(id)
	if request == nil {
		return
	}
	if complete {
		data, _ := json.Marshal(familyMoveDetail{Complete: true})
		request.detail = string(data)
	}
	f.finishMove(actx, request, complete)
}

func (f *Family) settleMove(actx actor.Context, request *UnsettledRequest) {
	var detail familyMoveDetail
	if err := json.Unmarshal([]byte(request.detail), &detail); err != nil {
		log.Printf("Family.settleMove character=%d id=%d: %v", f.owner.GetID(), request.id, err)
		return
	}
	f.finishMove(actx, request, detail.Complete)
}

func (f *Family) finishMove(actx actor.Context, request *UnsettledRequest, complete bool) {
	f.owner.gameWorld.Systems().Get[*InternalSystem]().FinishFamilyMoveAsync(actx, f.owner, request.id, complete).Bind(f.owner).Do(func(reply *internal.FamilyReply) error {
		f.owner.unsettled.Erase(request.id)
		if family := reply.GetFamily(); family != nil {
			f.Set(FamilyTreeFromProto(family))
		}
		return nil
	}).OnError(func(err error) {
		request.inFlight = false

		log.Printf("Family.finishMove character=%d id=%d: %v", f.owner.GetID(), request.id, err)
	})
}

func (f *Family) BuffRemaining() time.Duration {
	until := f.Self().BuffUntil()
	if until == 0 {
		return 0
	}
	return time.Until(time.UnixMilli(until))
}

func (f *Family) ExpRate() int {
	if f.BuffRemaining() <= 0 {
		return 100
	}

	switch f.Self().BuffFeature() {
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

	switch f.Self().BuffFeature() {
	case internal.FamilyFeature_FAMILY_FEATURE_DROP_15_MINUTES, internal.FamilyFeature_FAMILY_FEATURE_DROP_30_MINUTES:
		return 120
	case internal.FamilyFeature_FAMILY_FEATURE_UNITY:
		return 150
	}
	return 100
}

func (f *Family) gainProgress(exp uint32) {
	seniorID := f.Self().SeniorID()
	if seniorID == 0 || f.owner.gameWorld == nil {
		return
	}

	resources := f.owner.gameWorld.GetResources()
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
		f.progress += units

		receiverID := seniorID
		for _, multiplier := range []uint64{10, 5} {
			senior := f.tree.Member(receiverID)
			if senior == nil {
				break
			}

			amount := units * multiplier
			if uint32(level) > senior.Level() {
				amount /= 2
			}
			f.reputation[receiverID] += amount
			receiverID = senior.SeniorID()
		}
		remaining -= step
		old += step
		if old == needed {
			level++
			old = 0
		}
	}
}
