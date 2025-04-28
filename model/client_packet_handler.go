package model

import (
	"github.com/boyism80/fm/dao"
	"github.com/boyism80/fm/handler"
	"github.com/boyism80/fm/packet/req"
	"github.com/boyism80/fm/packet/resp"
)

func RegisterPacketHandler(client *ClientActor, h *handler.PacketHandler) {
	handler.RegisterPacketHandler(0x0A, client, h, onPong)
	handler.RegisterPacketHandler(0x01, client, h, onLogin)
	handler.RegisterPacketHandler(0x04, client, h, onCharacterList)
	handler.RegisterPacketHandler(0x07, client, h, onCheckName)
	handler.RegisterPacketHandler(0x08, client, h, onCreateCharacter)
	handler.RegisterPacketHandler(0x09, client, h, onDeleteCharacter)
	handler.RegisterPacketHandler(0x05, client, h, onSelectCharacter)
}

func onPong(client *ClientActor, request *req.Pong) {
}

func onLogin(client *ClientActor, request *req.Login) {

	if request.Id == "cshyeon" {
		client.Send(&resp.Authenticate{
			AccountId:     2390,
			Gender:        0,
			Admin:         true,
			AccountName:   request.Id,
			IsChatBlocked: false,
			ChatBlockTime: 116445060000000000,
		}, SEND_POLICY_ENCRYPT)

		for i := 0; i < 10; i++ {
			client.Send(&resp.ServerList{
				ServerId:     uint8(i),
				ChannelSize:  5,
				WorldName:    "cshyeon",
				Flag:         0,
				EventMessage: "채승현이 간다.",
			}, SEND_POLICY_ENCRYPT)
		}

		client.Send(&resp.EndOfServerList{}, SEND_POLICY_ENCRYPT)
	} else {
		client.Send(&resp.LoginFailed{Reason: resp.LoginFailedReasonNoPopup}, SEND_POLICY_ENCRYPT)
		client.Send(&resp.Notice{
			Type:    resp.NoticeTypePopup,
			Channel: 0,
			Message: "Hello",
			MegaEar: false}, SEND_POLICY_ENCRYPT)
	}
}

func onCharacterList(client *ClientActor, request *req.CharacterList) {
	client.Send(&resp.CharacterList{
		Characters: []dao.Character{
			dao.NewDummyCharacter(1, "채승현"),
			dao.NewDummyCharacter(2, "채진영"),
		},
		SlotCount: 6,
	}, SEND_POLICY_ENCRYPT)
}

func onCheckName(client *ClientActor, request *req.CheckName) {
	exists := request.Name == "채승현"
	client.Send(&resp.CheckName{Name: request.Name, Exists: exists}, SEND_POLICY_ENCRYPT)
}

func onCreateCharacter(client *ClientActor, request *req.CreateCharacter) {

	success := request.Name != "채진영"
	client.Send(&resp.CreateCharacter{
		Success: success,
		Character: &dao.Character{
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
	}, SEND_POLICY_ENCRYPT)
}

func onDeleteCharacter(client *ClientActor, request *req.DeleteCharacter) {
	client.Send(&resp.DeleteCharacter{
		Id:      request.Id,
		Success: true,
	}, SEND_POLICY_ENCRYPT)
}

func onSelectCharacter(client *ClientActor, request *req.SelectCharacter) {
	client.Send(&resp.Transfer{
		IP:          "127.0.0.1",
		Port:        7100,
		CharacterId: request.CharacterId,
	}, SEND_POLICY_ENCRYPT)
}
