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
			{
				Id:         1,
				Name:       "채승현",
				Gender:     0,
				SkinColor:  0,
				Face:       20401,
				Hair:       30027,
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
			{
				Id:         2,
				Name:       "채진영",
				Gender:     0,
				SkinColor:  0,
				Face:       20401,
				Hair:       30027,
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
		},
		SlotCount: 6,
	}, SEND_POLICY_ENCRYPT)
}
