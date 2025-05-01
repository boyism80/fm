package actor

import (
	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/entity"
	"github.com/boyism80/fm/handler"
	"github.com/boyism80/fm/msg"
	"github.com/boyism80/fm/packet/req"
	"github.com/boyism80/fm/packet/resp"
	"github.com/boyism80/fm/types"
)

func RegisterPacketHandler(ctx actor.Context, client *ClientActor, h *handler.PacketHandler) {
	handler.RegisterPacketHandler(0x0A, ctx, client, h, onPong)
	handler.RegisterPacketHandler(0x01, ctx, client, h, onLogin)
	handler.RegisterPacketHandler(0x04, ctx, client, h, onCharacterList)
	handler.RegisterPacketHandler(0x07, ctx, client, h, onCheckName)
	handler.RegisterPacketHandler(0x08, ctx, client, h, onCreateCharacter)
	handler.RegisterPacketHandler(0x09, ctx, client, h, onDeleteCharacter)
	handler.RegisterPacketHandler(0x05, ctx, client, h, onSelectCharacter)

	// game
	handler.RegisterPacketHandler(0x06, ctx, client, h, onLoginGame)
}

func onPong(ctx actor.Context, client *ClientActor, request *req.Pong) {
}

func onLogin(ctx actor.Context, client *ClientActor, request *req.Login) {

	if request.Id == "cshyeon" {
		client.Send(&resp.Authenticate{
			AccountId:     2390,
			Gender:        0,
			Admin:         true,
			AccountName:   request.Id,
			IsChatBlocked: false,
			ChatBlockTime: 116445060000000000,
		}, types.SEND_POLICY_ENCRYPT)

		for i := 0; i < 10; i++ {
			client.Send(&resp.ServerList{
				ServerId:     uint8(i),
				ChannelSize:  5,
				WorldName:    "cshyeon",
				Flag:         0,
				EventMessage: "채승현이 간다.",
			}, types.SEND_POLICY_ENCRYPT)
		}

		client.Send(&resp.EndOfServerList{}, types.SEND_POLICY_ENCRYPT)
	} else {
		client.Send(&resp.LoginFailed{Reason: resp.LoginFailedReasonNoPopup}, types.SEND_POLICY_ENCRYPT)
		client.Send(&resp.Notice{
			Type:    resp.NoticeTypePopup,
			Channel: 0,
			Message: "Hello",
			MegaEar: false}, types.SEND_POLICY_ENCRYPT)
	}
}

func onCharacterList(ctx actor.Context, client *ClientActor, request *req.CharacterList) {
	client.Send(&resp.CharacterList{
		Characters: []entity.Character{
			entity.NewDummyCharacter(1, "채승현"),
			entity.NewDummyCharacter(2, "채진영"),
		},
		SlotCount: 6,
	}, types.SEND_POLICY_ENCRYPT)
}

func onCheckName(ctx actor.Context, client *ClientActor, request *req.CheckName) {
	exists := request.Name == "채승현"
	client.Send(&resp.CheckName{Name: request.Name, Exists: exists}, types.SEND_POLICY_ENCRYPT)
}

func onCreateCharacter(ctx actor.Context, client *ClientActor, request *req.CreateCharacter) {

	success := request.Name != "채진영"
	client.Send(&resp.CreateCharacter{
		Success: success,
		Character: &entity.Character{
			Id:         1,
			Name:       request.Name,
			Gender:     0,
			SkinColor:  0,
			Face:       request.Face,
			Hair:       request.Hair,
			Level:      1,
			Class:      0,
			Str:        12,
			Dex:        5,
			Int:        4,
			Luk:        4,
			Hp:         50,
			MaxHp:      50,
			Mp:         5,
			MaxMp:      5,
			SpawnPoint: 3,
		},
	}, types.SEND_POLICY_ENCRYPT)
}

func onDeleteCharacter(ctx actor.Context, client *ClientActor, request *req.DeleteCharacter) {
	client.Send(&resp.DeleteCharacter{
		Id:      request.Id,
		Success: true,
	}, types.SEND_POLICY_ENCRYPT)
}

func onSelectCharacter(ctx actor.Context, client *ClientActor, request *req.SelectCharacter) {
	client.Send(&resp.Transfer{
		IP:          "127.0.0.1",
		Port:        7100,
		CharacterId: request.CharacterId,
	}, types.SEND_POLICY_ENCRYPT)
}

func onLoginGame(ctx actor.Context, client *ClientActor, request *req.LoginGame) {
	ch := entity.NewDummyCharacter(request.PlayerId, "강원기")
	client.BindCharacter(ctx, &ch)

	mapActor := client.Context.MapActor(client.Character.Map)
	ctx.Send(mapActor, &msg.EnterMap{
		Id:  ch.Id,
		PID: ctx.Self(),
	})

	// TODO: 맵에 EnterMap 메시지가 전달된 이후에
	// 맵의 모든 오브젝트에게 Warped 메시지가 전달된다.
	// client.Send(&resp.Warp{Character: &ch}, types.SEND_POLICY_ENCRYPT)
}
