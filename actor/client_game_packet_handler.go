package actor

import (
	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/entity"
	"github.com/boyism80/fm/handler"
	"github.com/boyism80/fm/packet/req"
	"github.com/boyism80/fm/packet/resp"
	"github.com/boyism80/fm/types"
)

func RegisterGameClientPacketHandler(ctx actor.Context, client *GameClientActor, h *handler.PacketHandler) {
	handler.RegisterPacketHandler(0x0A, ctx, client, h, onGameClientPong)
	handler.RegisterPacketHandler(0x06, ctx, client, h, onLoginGame)
}

func onGameClientPong(ctx actor.Context, client *GameClientActor, request *req.Pong) {
}

func onLoginGame(ctx actor.Context, client *GameClientActor, request *req.LoginGame) {
	ch := entity.NewDummyCharacter(request.PlayerId, "채승현")
	client.BindCharacter(ctx, &ch)

	// mapActor := client.Context.MapActor(client.Character.Map)
	// ctx.Send(mapActor, &msg.EnterMap{
	// 	Id:  ch.Id,
	// 	PID: ctx.Self(),
	// })

	// TODO: 맵에 EnterMap 메시지가 전달된 이후에
	// 맵의 모든 오브젝트에게 Warped 메시지가 전달된다.
	client.Send(&resp.Warp{Character: &ch}, types.SEND_POLICY_ENCRYPT)
}
