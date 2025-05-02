package actor

import (
	"github.com/asynkron/protoactor-go/actor"
	protoactor "github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/common/context"
	"github.com/boyism80/fm/common/handler"
	"github.com/boyism80/fm/game/msg"
)

type MapActor struct {
	objectPIDs map[uint32]*actor.PID // 오브젝트 ID → PID
	handler    *handler.MessageHandler
}

func NewMapActorProps(ctx protoactor.Context, serverCtx context.IServerContext) *actor.Props {
	return actor.PropsFromProducer(func() actor.Actor {

		act := &MapActor{
			objectPIDs: make(map[uint32]*actor.PID),
			handler:    handler.NewMessageHandler(),
		}

		handler.RegisterHandler(ctx, act, act.handler, onMapEnter)
		handler.RegisterHandler(ctx, act, act.handler, onMapLeave)
		handler.RegisterHandler(ctx, act, act.handler, onMapPidList)
		handler.RegisterHandler(ctx, act, act.handler, onMapBroadcastRange)

		return act
	})
}

func (state *MapActor) Receive(ctx protoactor.Context) {
	state.handler.Handle(ctx)
}

func onMapEnter(ctx protoactor.Context, state *MapActor, m *msg.EnterMap) {

	// 기존에 있던 오브젝트들에게 새로 추가된 오브젝트 알림
	for _, pid := range state.objectPIDs {
		ctx.Send(pid, &msg.Warped{
			Sender: m.PID,
		})
	}

	// 맵에 플레이어를 추가
	state.objectPIDs[m.Id] = m.PID
}

func onMapLeave(ctx protoactor.Context, state *MapActor, m *msg.LeaveMap) {
	delete(state.objectPIDs, m.Id)
}

func onMapPidList(ctx protoactor.Context, state *MapActor, m *msg.MapPidList) {
	var pids []*actor.PID
	for _, pid := range state.objectPIDs {
		pids = append(pids, pid)
	}
	ctx.Send(ctx.Sender(), &msg.MapPidList{Targets: pids})
}

func onMapBroadcastRange(ctx protoactor.Context, state *MapActor, m *msg.MapBroadcastRange) {
	// 나중에 섹터 추가하고 섹터 찾아서 섹터 액터한테 던짐

	for _, pid := range state.objectPIDs {
		if pid == m.Sender {
			continue
		}

		ctx.Send(pid, m.Message)
	}
}
