package server

import (
	"log"

	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/game/constant"
	"github.com/boyism80/fm/protocol/request"
	"github.com/boyism80/fm/protocol/response"
	"github.com/boyism80/fm/types"
)

type Login struct {
	ls     *LoginServer
	opcode byte
}

func (Login) New(ls *LoginServer) *Login {
	return &Login{
		ls:     ls,
		opcode: 0x01,
	}
}

func (h *Login) GetOpcode() byte {
	return h.opcode
}

func (h *Login) Handle(ctx *core.ClientContext, req *request.Login) error {
	log.Printf("Login packet received from %s - ID: %s, MAC: %s",
		ctx.Client.GetConnection().RemoteAddr(), req.ID, req.Mac)

	if req.ID == "cshyeon" {
		authResp := &response.Authenticate{
			AccountId:     2390,
			Gender:        0,
			Admin:         true,
			AccountName:   req.ID,
			IsChatBlocked: false,
			ChatBlockTime: 116445060000000000,
		}
		if err := ctx.Client.Send(authResp, types.SEND_POLICY_ENCRYPT); err != nil {
			log.Printf("Failed to send authenticate response: %v", err)
			return err
		}

		for i := 0; i < 10; i++ {
			serverResp := &response.ServerList{
				ServerId:     uint8(i),
				ChannelSize:  5,
				WorldName:    "cshyeon",
				Flag:         0,
				EventMessage: "채승현이 간다.",
			}
			if err := ctx.Client.Send(serverResp, types.SEND_POLICY_ENCRYPT); err != nil {
				log.Printf("Failed to send server list response: %v", err)
				return err
			}
		}

		endResp := &response.EndOfServerList{}
		if err := ctx.Client.Send(endResp, types.SEND_POLICY_ENCRYPT); err != nil {
			log.Printf("Failed to send end of server list: %v", err)
			return err
		}
	} else {
		failedResp := &response.LoginFailed{Reason: response.LoginFailedReasonNoPopup}
		if err := ctx.Client.Send(failedResp, types.SEND_POLICY_ENCRYPT); err != nil {
			log.Printf("Failed to send login failed response: %v", err)
			return err
		}

		noticeResp := &response.Notice{
			Type:    constant.MSG_POPUP,
			Channel: 0,
			Message: "Hello",
			MegaEar: false,
		}
		if err := ctx.Client.Send(noticeResp, types.SEND_POLICY_ENCRYPT); err != nil {
			log.Printf("Failed to send notice: %v", err)
			return err
		}
	}

	return nil
}
