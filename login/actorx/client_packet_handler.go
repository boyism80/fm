package actorx

import (
	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/common/handler"
	common_req "github.com/boyism80/fm/common/protocol/req"
	"github.com/boyism80/fm/common/types"
	"github.com/boyism80/fm/game/entity"
	"github.com/boyism80/fm/login/protocol/req"
	"github.com/boyism80/fm/login/protocol/resp"
)

func RegisterLoginClientPacketHandler(ctx actor.Context, client *LoginClientActor, h *handler.PacketHandler) {

	handler.RegisterPacketHandler(0x0A, ctx, client, h, onLoginClientPong)
	handler.RegisterPacketHandler(0x01, ctx, client, h, onLoginClientLogin)
	handler.RegisterPacketHandler(0x04, ctx, client, h, onLoginClientCharacterList)
	handler.RegisterPacketHandler(0x07, ctx, client, h, onLoginClientCheckName)
	handler.RegisterPacketHandler(0x08, ctx, client, h, onLoginClientCreateCharacter)
	handler.RegisterPacketHandler(0x09, ctx, client, h, onLoginClientDeleteCharacter)
	handler.RegisterPacketHandler(0x05, ctx, client, h, onLoginClientSelectCharacter)
}

func onLoginClientPong(ctx actor.Context, client *LoginClientActor, request *common_req.Pong) {
}

func onLoginClientLogin(ctx actor.Context, client *LoginClientActor, request *req.Login) {

	if request.ID == "cshyeon" {
		client.Send(&resp.Authenticate{
			AccountId:     2390,
			Gender:        0,
			Admin:         true,
			AccountName:   request.ID,
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

func onLoginClientCharacterList(ctx actor.Context, client *LoginClientActor, request *req.CharacterList) {
	client.Send(&resp.CharacterList{
		Characters: []entity.Character{
			entity.NewDummyCharacter(client, nil, 1, "채승현", nil),
			entity.NewDummyCharacter(client, nil, 2, "채진영", nil),
		},
		SlotCount: 6,
	}, types.SEND_POLICY_ENCRYPT)
}

func onLoginClientCheckName(ctx actor.Context, client *LoginClientActor, request *req.CheckName) {
	exists := request.Name == "채승현"
	client.Send(&resp.CheckName{Name: request.Name, Exists: exists}, types.SEND_POLICY_ENCRYPT)
}

func onLoginClientCreateCharacter(ctx actor.Context, client *LoginClientActor, request *req.CreateCharacter) {

	success := request.Name != "채진영"
	client.Send(&resp.CreateCharacter{
		Success: success,
		Character: &entity.Character{
			Life: entity.Life{
				Hp:    50,
				MaxHp: 50,
				Mp:    5,
				MaxMp: 5,
			},
			ID:         1,
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
			SpawnPoint: 3,
		},
	}, types.SEND_POLICY_ENCRYPT)
}

func onLoginClientDeleteCharacter(ctx actor.Context, client *LoginClientActor, request *req.DeleteCharacter) {
	client.Send(&resp.DeleteCharacter{
		ID:      request.ID,
		Success: true,
	}, types.SEND_POLICY_ENCRYPT)
}

func onLoginClientSelectCharacter(ctx actor.Context, client *LoginClientActor, request *req.SelectCharacter) {
	client.Send(&resp.Transfer{
		IP:          "127.0.0.1",
		Port:        7111,
		CharacterId: request.CharacterId,
	}, types.SEND_POLICY_ENCRYPT)
}
