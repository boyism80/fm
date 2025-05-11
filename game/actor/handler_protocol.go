package actor

import (
	"strings"

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
	handler.RegisterPacketHandler(0x20, ctx, client, h, onGameClientNormalChat)
	handler.RegisterPacketHandler(0x1B, ctx, client, h, onGameClientAttack)
}

func onGameClientPong(ctx actor.Context, client *GameClientActor, request *common_req.Pong) {
}

func onLoginGame(ctx actor.Context, client *GameClientActor, request *req.LoginGame) {
	name := "채승현"
	if request.PlayerId != 1 {
		name = "채진영"
	}
	ch := entity.NewDummyCharacter(request.PlayerId, name, client.ctx)
	spawnPoint := client.ctx.Resources.Maps[ch.Map].Portals[ch.SpawnPoint].Position
	ch.Position = spawnPoint
	client.ch = &ch
	RegisterLifeHandlers(ctx, &ch.Life, client.messageHandler)

	// TODO: 맵에 EnterMap 메시지가 전달된 이후에
	// 맵의 모든 오브젝트에게 Warped 메시지가 전달된다.
	client.Send(&resp.Warp{Character: &ch}, types.SEND_POLICY_ENCRYPT)

	mapActor := client.ctx.MapActors[client.ch.Map]
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
			ExceptSelf: true,
		})
	}
}

func onGameClientMovePlayer(ctx actor.Context, client *GameClientActor, req *req.MovePlayer) {

	beforePosition := client.ch.Position

	for _, frag := range req.Fragments {
		if move, ok := frag.(*movement.AbsoluteLifeMovement); ok {
			client.ch.Position = move.Position
		}

		client.ch.Stance = frag.GetNewState()
	}

	mapActor := client.ctx.MapActors[client.ch.Map]
	if mapActor == nil {
		return
	}

	ctx.Send(mapActor, &msg.MapBroadcastRange{
		Sender: ctx.Self(),
		Pivot:  beforePosition,
		Message: &common_msg.SendProtocol{
			Protocol: &resp.Move{
				Character:     client.ch,
				MoveFragments: req.Fragments,
				StartPoint:    beforePosition,
			},
			Policy: types.SEND_POLICY_ENCRYPT,
		},
	})
}

func onGameClientNormalChat(ctx actor.Context, client *GameClientActor, req *req.NormalChat) {
	position := client.ch.Position
	mapActor := client.ctx.MapActors[client.ch.Map]
	if mapActor == nil {
		return
	}

	if strings.HasPrefix(req.Message, "/") {
		params := strings.Split(strings.TrimPrefix(req.Message, "/"), " ")
		err := client.commandHandler.Handle(ctx, params...)
		if err == nil {
			return
		}
	}

	ctx.Send(mapActor, &msg.MapBroadcastRange{
		Sender: ctx.Self(),
		Pivot:  position,
		Message: &common_msg.SendProtocol{
			Protocol: &resp.NormalChat{
				CharacterId: client.ch.Id,
				Highlight:   false,
				Message:     req.Message,
				Show:        req.Show,
			},
			Policy: types.SEND_POLICY_ENCRYPT,
		},
	})
}

func onGameClientAttack(ctx actor.Context, client *GameClientActor, req *req.Attack) {
	position := client.ch.Position
	mapActor := client.ctx.MapActors[client.ch.Map]
	if mapActor == nil {
		return
	}

	ctx.Send(mapActor, &msg.MapBroadcastRange{
		Sender: ctx.Self(),
		Pivot:  position,
		Message: &common_msg.SendProtocol{
			Protocol: &resp.Attack{
				CharacterId: client.ch.Id,
				AttackInfo:  req.AttackInfo,
				SkillLevel:  0,
			},
			Policy: types.SEND_POLICY_ENCRYPT,
		},
	})
}
