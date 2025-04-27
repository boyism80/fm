package model

import (
	"log"

	"github.com/boyism80/fm/handler"
	"github.com/boyism80/fm/packet/req"
	"github.com/boyism80/fm/packet/resp"
)

func RegisterPacketHandler(client *ClientActor, h *handler.PacketHandler) {
	handler.RegisterPacketHandler(0x0A, client, h, onPong)
	handler.RegisterPacketHandler(0x01, client, h, onLogin)
}

func onPong(client *ClientActor, request *req.Pong) {
	log.Printf("pong")
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
