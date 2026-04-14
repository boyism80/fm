package server

import (
	"context"
	"log"

	"github.com/boyism80/fm/core"
	fminternalpb "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
	"github.com/boyism80/fm/protocol/request"
	"github.com/boyism80/fm/protocol/response"
	"github.com/boyism80/fm/services/login/client"
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

	ic := h.ls.context.InternalClient
	if ic == nil {
		log.Printf("Internal gRPC client not configured, rejecting login")
		failedResp := &response.LoginFailed{Reason: response.LoginFailedReasonSystemError6}
		return ctx.Client.Send(failedResp, types.SEND_POLICY_ENCRYPT)
	}

	remoteAddr := ctx.Client.GetConnection().RemoteAddr().String()
	reply, err := ic.LoginAccount(context.Background(), &fminternalpb.LoginAccountRequest{
		LoginId:     req.ID,
		Password:    req.Pw,
		MacAddress:  req.Mac,
		IpAddress:   remoteAddr,
		InitialRole: h.ls.config.InitialRole,
	})
	if err != nil {
		log.Printf("LoginAccount RPC error: %v", err)
		failedResp := &response.LoginFailed{Reason: response.LoginFailedReasonSystemError6}
		return ctx.Client.Send(failedResp, types.SEND_POLICY_ENCRYPT)
	}

	switch reply.Status {
	case fminternalpb.LoginAccountReply_WRONG_PASSWORD:
		failedResp := &response.LoginFailed{Reason: response.LoginFailedReasonIncorrectPassword}
		return ctx.Client.Send(failedResp, types.SEND_POLICY_ENCRYPT)
	case fminternalpb.LoginAccountReply_BANNED:
		failedResp := &response.LoginFailed{Reason: response.LoginFailedReasonIDDeletedOrBlocked}
		return ctx.Client.Send(failedResp, types.SEND_POLICY_ENCRYPT)
	}

	loginClient, ok := ctx.Client.(*client.LoginClient)
	if ok {
		loginClient.SetAccountId(reply.AccountId)
	}

	authResp := &response.Authenticate{
		AccountId:     reply.AccountId,
		Gender:        uint8(reply.Gender),
		Role:          uint8(reply.Role),
		AccountName:   req.ID,
		IsChatBlocked: reply.IsChatBlocked,
	}
	if err := ctx.Client.Send(authResp, types.SEND_POLICY_ENCRYPT); err != nil {
		return err
	}

	for i := 0; i < 10; i++ {
		serverResp := &response.ServerList{
			ServerId:     uint8(i),
			ChannelSize:  5,
			WorldName:    "Scania",
			Flag:         0,
			EventMessage: "",
		}
		if err := ctx.Client.Send(serverResp, types.SEND_POLICY_ENCRYPT); err != nil {
			return err
		}
	}

	return ctx.Client.Send(&response.EndOfServerList{}, types.SEND_POLICY_ENCRYPT)
}
