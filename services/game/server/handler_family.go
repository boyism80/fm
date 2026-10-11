package server

import (
	"context"
	"sync/atomic"
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/core/async"
	"github.com/boyism80/fm/core/clock"
	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
	"github.com/boyism80/fm/protocol/request"
	"github.com/boyism80/fm/services/game/client"
	"github.com/boyism80/fm/services/game/entity"
	"github.com/google/uuid"
)

type familyHandler struct {
	gs *GameServer
}

func (h *familyHandler) handle(ctx *core.ClientContext, req *request.Family) error {
	if h.gs.internalClient == nil || ctx.ActorContext == nil {
		return nil
	}

	gameClient, ok := ctx.Client.(*client.GameClient)
	if !ok {
		return nil
	}
	ch := gameClient.GetCharacter()
	if ch == nil || ch.Family.Pending {
		return nil
	}

	if req.Action == internal.FamilyAction_FAMILY_ACTION_PRECEPT {
		if req.Name != "" && !entity.ClientText(req.Name).Valid(200) {
			return nil
		}
	} else if req.Name != "" && !entity.ClientText(req.Name).Valid(13) {
		return nil
	}

	ch.Family.Pending = true
	fail := func(result internal.FamilyResult) {
		ch.Family.Pending = false
		ch.Listener.OnFamilyResult(ch, result)
		ch.Listener.OnUnlockAction(ch)
	}
	rpc := &internal.FamilyRequest{
		OperationId: uuid.NewString(),
		WorldId:     h.gs.config.WorldId,
		CharacterId: ch.GetID(),
		ChannelId:   h.gs.config.ChannelId,
		Action:      req.Action,
		TargetId:    req.TargetID,
		Name:        req.Name,
		Feature:     req.Feature,
		Accepted:    req.Accepted,
		Level:       uint32(ch.GetLevel()),
		ClassId:     uint32(ch.Class),
	}

	recruit := req.Action == internal.FamilyAction_FAMILY_ACTION_INVITE || req.Action == internal.FamilyAction_FAMILY_ACTION_ACCEPT
	move := req.Action == internal.FamilyAction_FAMILY_ACTION_SUMMON ||
		(req.Action == internal.FamilyAction_FAMILY_ACTION_USE &&
			(req.Feature == internal.FamilyFeature_FAMILY_FEATURE_TELEPORT || req.Feature == internal.FamilyFeature_FAMILY_FEATURE_SUMMON))
	if !recruit && !move {
		h.call(ctx.ActorContext, ch, req, rpc, nil, nil)
		return nil
	}

	targetID := req.TargetID
	if targetID == 0 {
		var found bool
		targetID, found = h.gs.GetDispatchSystem().FindCharacterID(req.Name)
		if !found {
			fail(internal.FamilyResult_FAMILY_RESULT_WRONG_NAME)
			return nil
		}
	}
	if targetID == ch.GetID() {
		if recruit {
			fail(internal.FamilyResult_FAMILY_RESULT_CANNOT_BE_JUNIOR)
		} else {
			fail(internal.FamilyResult_FAMILY_RESULT_FAILED)
		}
		return nil
	}

	origin := ch.GetMap()
	if origin == nil {
		fail(internal.FamilyResult_FAMILY_RESULT_FAILED)
		return nil
	}

	// The target lives on another actor, so the timeout and the target's
	// answer race; whichever settles first wins.
	sourceRole := ch.Role
	var settled atomic.Bool
	settle := func(run func(actor.Context, *entity.Character)) {
		h.gs.GetDispatchSystem().CallCharacter(ch.GetID(), func(c actor.Context, current *entity.Character) {
			if settled.CompareAndSwap(false, true) {
				run(c, current)
			}
		})
	}
	time.AfterFunc(3*time.Second, func() {
		settle(func(_ actor.Context, current *entity.Character) {
			current.Family.Pending = false
			current.Listener.OnFamilyResult(current, internal.FamilyResult_FAMILY_RESULT_FAILED)
			current.Listener.OnUnlockAction(current)
		})
	})

	h.gs.GetDispatchSystem().CallCharacter(targetID, func(_ actor.Context, target *entity.Character) {
		targetMap := target.GetMap()
		targetLevel := target.GetLevel()
		name := target.GetName()
		result := internal.FamilyResult_FAMILY_RESULT_OK
		switch {
		case targetMap == nil || target.LoggedOut() || (target.IsHidden() && sourceRole < target.Role):
			result = internal.FamilyResult_FAMILY_RESULT_WRONG_NAME
		case recruit && targetMap != origin:
			result = internal.FamilyResult_FAMILY_RESULT_NOT_SAME_MAP
		case move && (targetMap.BlocksTeleportStone() || origin.BlocksTeleportStone()):
			result = internal.FamilyResult_FAMILY_RESULT_FAILED
		}

		settle(func(c actor.Context, current *entity.Character) {
			if result != internal.FamilyResult_FAMILY_RESULT_OK {
				fail(result)
				return
			}
			if current.GetMap() != origin {
				fail(internal.FamilyResult_FAMILY_RESULT_FAILED)
				return
			}
			if req.Action == internal.FamilyAction_FAMILY_ACTION_ACCEPT && req.Name != name {
				fail(internal.FamilyResult_FAMILY_RESULT_WRONG_NAME)
				return
			}

			rpc.TargetId = targetID
			rpc.TargetLevel = uint32(targetLevel)
			if origin.Wz != nil {
				rpc.MapName, _ = h.gs.GetResources().GetMapName(origin.TemplateID())
			}
			h.call(c, current, req, rpc, targetMap, origin)
		})
	})
	return nil
}

func (h *familyHandler) call(
	ctx actor.Context,
	ch *entity.Character,
	req *request.Family,
	rpc *internal.FamilyRequest,
	targetMap, origin *entity.Map,
) {
	teleport := (req.Action == internal.FamilyAction_FAMILY_ACTION_USE && req.Feature == internal.FamilyFeature_FAMILY_FEATURE_TELEPORT) ||
		(req.Action == internal.FamilyAction_FAMILY_ACTION_SUMMON && req.Accepted)
	deadline := clock.Now().Add(20 * time.Second)
	async.NewTask(ctx, core.InternalRPCPerStepTimeout).ThenRPC(
		func(c context.Context) (*internal.FamilyReply, error) {
			return h.gs.internalClient.Family(c, rpc)
		},
		func(reply *internal.FamilyReply) error {
			h.gs.GetDispatchSystem().CallCharacter(ch.GetID(), func(currentCtx actor.Context, current *entity.Character) {
				current.Family.Pending = false
				if reply.Status != internal.FamilyResult_FAMILY_RESULT_OK {
					current.Listener.OnFamilyResult(current, reply.Status)
					current.Listener.OnUnlockAction(current)
					return
				}

				if req.Action == internal.FamilyAction_FAMILY_ACTION_PEDIGREE {
					current.Listener.OnFamilyChart(current, reply.Family)
				} else {
					current.Family.Set(reply.Family)
				}

				if targetMap != nil && teleport {
					m := current.GetMap()
					if clock.Now().After(deadline) || m != origin || m.BlocksTeleportStone() || targetMap.BlocksTeleportStone() ||
						current.GetHp() == 0 || current.Dialog.Thread() != nil || current.Dialog.ShopID != 0 {
						h.refund(currentCtx, current, rpc)
						return
					}

					if err := h.gs.GetMapSystem().Warp(currentCtx, current, targetMap, 0, func(entered actor.Context) {
						h.settleMove(entered, rpc, true, 0)
					}); err != nil {
						h.refund(currentCtx, current, rpc)
					}
				}
				current.Listener.OnUnlockAction(current)
			})
			return nil
		},
	).OnError(func(error) {
		if teleport {
			h.settleMove(ctx, rpc, false, 0)
		}
		h.gs.GetDispatchSystem().CallCharacter(ch.GetID(), func(_ actor.Context, current *entity.Character) {
			current.Family.Pending = false
			current.Listener.OnFamilyResult(current, internal.FamilyResult_FAMILY_RESULT_FAILED)
			current.Listener.OnUnlockAction(current)
		})
	})
}

func (h *familyHandler) refund(ctx actor.Context, ch *entity.Character, original *internal.FamilyRequest) {
	h.settleMove(ctx, original, false, 0)
	ch.Listener.OnFamilyResult(ch, internal.FamilyResult_FAMILY_RESULT_FAILED)
	ch.Listener.OnUnlockAction(ch)
}

func (h *familyHandler) settleMove(ctx actor.Context, original *internal.FamilyRequest, complete bool, attempt int) {
	req := &internal.FamilyRequest{
		OperationId: original.OperationId,
		WorldId:     original.WorldId,
		CharacterId: original.CharacterId,
		ChannelId:   original.ChannelId,
		Action:      internal.FamilyAction_FAMILY_ACTION_REFUND,
	}
	if complete {
		req.Action = internal.FamilyAction_FAMILY_ACTION_MOVE_COMPLETE
	}

	async.NewTask(ctx, core.InternalRPCPerStepTimeout).ThenRPC(
		func(c context.Context) (*internal.FamilyReply, error) {
			return h.gs.internalClient.Family(c, req)
		},
		func(reply *internal.FamilyReply) error {
			h.gs.GetDispatchSystem().CallCharacter(original.CharacterId, func(_ actor.Context, current *entity.Character) {
				current.Family.Set(reply.Family)
			})
			return nil
		},
	).OnError(func(error) {
		if attempt >= 2 {
			return
		}

		h.gs.GetDispatchSystem().CallCharacter(original.CharacterId, func(currentCtx actor.Context, _ *entity.Character) {
			h.settleMove(currentCtx, original, complete, attempt+1)
		})
	})
}
