package actor

import (
	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/common/context"
	"github.com/boyism80/fm/common/handler"
	common_msg "github.com/boyism80/fm/common/msg"
	"github.com/boyism80/fm/common/types"
	"github.com/boyism80/fm/game/constant"
	"github.com/boyism80/fm/game/entity"
	"github.com/boyism80/fm/game/msg"
	"github.com/boyism80/fm/game/protocol/resp"
)

type MesoActor struct {
	entity.Object
	oid     uint32
	ctx     *context.ServerContext
	handler *handler.MessageHandler
	mapPid  *actor.PID
	looting bool
	meso    int32
}

func NewMesoActor(ctx actor.Context,
	serverCtx *context.ServerContext,
	meso int32,
	oid uint32,
	mapPid *actor.PID,
	position types.Vector2[int16]) actor.Actor {

	act := &MesoActor{
		Object: entity.Object{
			Position: position,
		},
		oid:     oid,
		ctx:     serverCtx,
		handler: handler.NewMessageHandler(),
		meso:    meso,
		mapPid:  mapPid,
	}
	RegisterObjectHandlers(ctx, &act.Object, act.handler)
	handler.RegisterHandler(ctx, act, act.handler, onMesoSpawn)
	handler.RegisterHandler(ctx, act, act.handler, onMesoLooting)
	handler.RegisterHandler(ctx, act, act.handler, onMesoLooted)
	return act
}

func (state *MesoActor) Receive(context actor.Context) {
	state.handler.Handle(context)
}

func onMesoSpawn(ctx actor.Context, state *MesoActor, request *msg.MesoSpawn) {
	ctx.Send(state.mapPid, &msg.MapBroadcastRange{
		Sender:     ctx.Self(),
		Pivot:      state.Object.Position,
		ExceptSelf: true,
		Message: &common_msg.SendProtocol{
			Protocol: &resp.DropMeso{
				Id:           state.oid,
				Animation:    constant.DropItemAnimationTypeDefault,
				DropType:     2,
				Meso:         state.meso,
				OwnerId:      request.OwnerId,
				Position:     request.Position,
				SpawnedPoint: request.SpawnedPoint,
				IsPlayerDrop: true,
			},
			Policy: types.SEND_POLICY_ENCRYPT,
		},
	})
}

func onMesoLooting(ctx actor.Context, state *MesoActor, request *msg.ItemLooting) {
	if state.looting {
		ctx.Send(request.Actor, &msg.CharacterLootFailed{
			Oid: state.oid,
		})
	} else {
		ctx.Send(request.Actor, &msg.CharacterMesoLooting{
			Pid:  ctx.Self(),
			Oid:  state.oid,
			Meso: state.meso,
		})
		state.looting = true
	}
}

func onMesoLooted(ctx actor.Context, state *MesoActor, request *msg.ItemLooted) {
	state.looting = false
	if request.Success {
		state.meso -= request.Count
		if state.meso == 0 {
			ctx.Send(state.mapPid, &msg.MapItemLooted{
				Actor:       ctx.Self(),
				Oid:         state.oid,
				CharacterId: request.CharacterId,
				Mode:        resp.RemoveItemTypeAnimated,
				Position:    state.Object.Position,
			})
		}
	}
}
