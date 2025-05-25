package actor

import (
	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/common/context"
	"github.com/boyism80/fm/common/handler"
	"github.com/boyism80/fm/common/types"
	"github.com/boyism80/fm/game/data"
	"github.com/boyism80/fm/game/msg"
	"github.com/boyism80/fm/game/protocol/resp"

	common_msg "github.com/boyism80/fm/common/msg"
)

type MapActor struct {
	objectPIDs map[uint32]*actor.PID // 오브젝트 ID → PID
	sequence   uint32
	handler    *handler.MessageHandler
	Spec       *data.MapSpec
	ctx        *context.ServerContext
}

func NewMapActorProps(ctx actor.Context, serverCtx *context.ServerContext, spec *data.MapSpec) *actor.Props {
	return actor.PropsFromProducer(func() actor.Actor {

		act := &MapActor{
			objectPIDs: make(map[uint32]*actor.PID),
			handler:    handler.NewMessageHandler(),
			Spec:       spec,
			ctx:        serverCtx,
		}

		handler.RegisterHandler(ctx, act, act.handler, onMapEnter)
		handler.RegisterHandler(ctx, act, act.handler, onMapLeave)
		handler.RegisterHandler(ctx, act, act.handler, onMapPidList)
		handler.RegisterHandler(ctx, act, act.handler, onMapBroadcastRange)
		handler.RegisterHandler(ctx, act, act.handler, onMapSpawnItem)

		return act
	})
}

func (state *MapActor) Receive(ctx actor.Context) {
	state.handler.Handle(ctx)
}

func onMapEnter(ctx actor.Context, state *MapActor, m *msg.EnterMap) {

	// 기존에 있던 오브젝트들에게 새로 추가된 오브젝트 알림
	for _, pid := range state.objectPIDs {
		ctx.Send(pid, &msg.Warped{
			Sender: m.PID,
		})
	}

	// 맵에 플레이어를 추가
	state.objectPIDs[m.Id] = m.PID
}

func onMapLeave(ctx actor.Context, state *MapActor, m *msg.LeaveMap) {
	delete(state.objectPIDs, m.Id)

	for _, pid := range state.objectPIDs {
		ctx.Send(pid, &common_msg.SendProtocol{
			Protocol: &resp.LeavePlayer{
				Id: m.Id,
			},
			Policy: types.SEND_POLICY_ENCRYPT,
		})
	}
}

func onMapPidList(ctx actor.Context, state *MapActor, m *msg.MapPidList) {
	var pids []*actor.PID
	for _, pid := range state.objectPIDs {
		pids = append(pids, pid)
	}
	ctx.Send(ctx.Sender(), &msg.MapPidList{Targets: pids})
}

func onMapBroadcastRange(ctx actor.Context, state *MapActor, m *msg.MapBroadcastRange) {
	// 나중에 섹터 추가하고 섹터 찾아서 섹터 액터한테 던짐

	for _, pid := range state.objectPIDs {
		if m.ExceptSelf && pid == m.Sender {
			continue
		}

		ctx.Send(pid, m.Message)
	}
}

func onMapSpawnItem(ctx actor.Context, state *MapActor, m *msg.MapSpawnItem) {
	obj := m.Item.GetObject()
	dropPoint, ok := state.Spec.DropPoint(obj.Position)
	if ok {
		dropPoint = obj.Position
	}
	obj.Position = dropPoint

	props := actor.PropsFromProducer(func() actor.Actor {
		state.sequence++
		return NewItemActor(ctx,
			state.ctx,
			m.Item,
			state.sequence,
			ctx.Self())
	})
	pid := ctx.Spawn(props)

	ctx.Send(pid, &msg.ItemSpawn{
		OwnerId:      m.OwnerId,
		SpawnedPoint: m.SpawnedPoint,
	})
}
