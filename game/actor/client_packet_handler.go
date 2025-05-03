package actor

import (
	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/common/handler"
	"github.com/boyism80/fm/common/types"
	"github.com/boyism80/fm/game/entity"
	"github.com/boyism80/fm/game/entity/movement"
	"github.com/boyism80/fm/game/msg"
	"github.com/boyism80/fm/game/protocol/req"
	"github.com/boyism80/fm/game/protocol/resp"

	common_msg "github.com/boyism80/fm/common/msg"
	common_req "github.com/boyism80/fm/common/protocol/req"
)

func RegisterGameClientPacketHandler(ctx actor.Context, client *GameClientActor, h *handler.PacketHandler) {
	handler.RegisterPacketHandler(0x0A, ctx, client, h, onGameClientPong)
	handler.RegisterPacketHandler(0x06, ctx, client, h, onLoginGame)
	handler.RegisterPacketHandler(0x18, ctx, client, h, onGameClientMovePlayer)
}

func onGameClientPong(ctx actor.Context, client *GameClientActor, request *common_req.Pong) {
}

func onLoginGame(ctx actor.Context, client *GameClientActor, request *req.LoginGame) {
	name := "채승현"
	if request.PlayerId != 35177 {
		name = "채진영"
	}
	ch := entity.NewDummyCharacter(request.PlayerId, name)
	spawnPoint := client.Context.GameData.Maps[ch.Map].Portals[ch.SpawnPoint].Position
	ch.Position = spawnPoint
	client.BindCharacter(ctx, &ch)

	// TODO: 맵에 EnterMap 메시지가 전달된 이후에
	// 맵의 모든 오브젝트에게 Warped 메시지가 전달된다.
	client.Send(&resp.Warp{Character: &ch}, types.SEND_POLICY_ENCRYPT)

	mapActor := client.Context.MapActors[client.Character.Map]
	if mapActor != nil {
		ctx.Send(mapActor, &msg.EnterMap{
			Id:  ch.Id,
			PID: ctx.Self(),
		})

		// 기존 오브젝트들에게 날 보여줌
		spawnResp := resp.SpawnPlayer{
			Character:       &ch,
			BuffStates:      [4]uint32{},
			Diseases:        [4]uint32{},
			CrushRings:      []*entity.Ring{},
			FriendshipRings: []*entity.Ring{},
			MarriageRings:   []*entity.Ring{},
		}
		ctx.Send(mapActor, &msg.MapBroadcastRange{
			Sender: ctx.Self(),
			Pivot:  ch.Position,
			Message: &common_msg.SendProtocol{
				Protocol: &spawnResp,
				Policy:   types.SEND_POLICY_ENCRYPT,
			},
		})
	}
}

func onGameClientMovePlayer(ctx actor.Context, client *GameClientActor, req *req.MovePlayer) {

	beforePosition := client.Character.Position

	for _, frag := range req.Fragments {
		if move, ok := frag.(*movement.AbsoluteLifeMovement); ok {
			client.Character.Position = move.Position
		}

		client.Character.Stance = frag.GetNewState()
	}

	mapActor := client.Context.MapActors[client.Character.Map]
	if mapActor == nil {
		return
	}

	ctx.Send(mapActor, &msg.MapBroadcastRange{
		Sender: ctx.Self(),
		Pivot:  beforePosition,
		Message: &common_msg.SendProtocol{
			Protocol: &resp.Move{
				Character:     client.Character,
				MoveFragments: req.Fragments,
				StartPoint:    beforePosition,
			},
			Policy: types.SEND_POLICY_ENCRYPT,
		},
	})
}
