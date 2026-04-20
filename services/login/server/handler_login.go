package server

import (
	"context"
	"fmt"
	"log"

	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/core/async"
	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
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

	ic := h.ls.internalClient
	if ic == nil {
		log.Printf("Internal gRPC client not configured, rejecting login")
		failedResp := &response.LoginFailed{Reason: response.LoginFailedReasonSystemError6}
		return ctx.Client.Send(failedResp, types.SEND_POLICY_ENCRYPT)
	}

	remoteAddr := ctx.Client.GetConnection().RemoteAddr().String()
	reqMsg := &internal.LoginAccountRequest{
		LoginId:     req.ID,
		Password:    req.Pw,
		MacAddress:  req.Mac,
		IpAddress:   remoteAddr,
		InitialRole: h.ls.config.InitialRole,
	}

	if ctx.ActorContext == nil {
		log.Printf("Login: no actor context, cannot run internal RPC")
		_ = ctx.Client.Send(&response.LoginFailed{Reason: response.LoginFailedReasonSystemError6}, types.SEND_POLICY_ENCRYPT)
		return fmt.Errorf("login: actor context required for internal RPC")
	}

	async.ThenRPC(async.NewPromise(ctx.ActorContext, core.InternalRPCPerStepTimeout),
		func(c context.Context) (*internal.LoginAccountReply, error) {
			return ic.LoginAccount(c, reqMsg)
		}, func(reply *internal.LoginAccountReply) error {
			return h.sendLoginAccountResult(ctx, req, reply)
		}).
		OnError(func(err error) {
			log.Printf("Login (async): %v", err)
			_ = ctx.Client.Send(&response.LoginFailed{Reason: response.LoginFailedReasonSystemError6}, types.SEND_POLICY_ENCRYPT)
		}).
		Run()
	return nil
}

func (h *Login) sendLoginAccountResult(ctx *core.ClientContext, req *request.Login, reply *internal.LoginAccountReply) error {
	switch reply.Status {
	case internal.LoginAccountReply_WRONG_PASSWORD:
		failedResp := &response.LoginFailed{Reason: response.LoginFailedReasonIncorrectPassword}
		return ctx.Client.Send(failedResp, types.SEND_POLICY_ENCRYPT)
	case internal.LoginAccountReply_BANNED:
		failedResp := &response.LoginFailed{Reason: response.LoginFailedReasonIDDeletedOrBlocked}
		return ctx.Client.Send(failedResp, types.SEND_POLICY_ENCRYPT)
	case internal.LoginAccountReply_ALREADY_LOGGED_IN:
		failedResp := &response.LoginFailed{Reason: response.LoginFailedReasonAlreadyLoggedIn}
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
		ChatBlockTime: uint64(reply.ChatBlockedUntilUnixMs),
	}
	if err := ctx.Client.Send(authResp, types.SEND_POLICY_ENCRYPT); err != nil {
		return err
	}

	for _, world := range h.ls.GetWorldCatalog() {
		if len(world.GetChannels()) == 0 {
			continue
		}
		channels := make([]response.ServerChannel, 0, len(world.GetChannels()))
		for _, ch := range world.GetChannels() {
			channelName := ch.GetName()
			if channelName == "" {
				channelName = fmt.Sprintf("%s-%d", world.GetWorldName(), ch.GetChannelId()+1)
			}
			channels = append(channels, response.ServerChannel{
				ChannelID: uint16(ch.GetChannelId()),
				Name:      channelName,
				Load:      1200,
			})
		}
		serverResp := &response.ServerList{
			ServerId:     uint8(world.GetWorldId()),
			Channels:     channels,
			WorldName:    world.GetWorldName(),
			Flag:         byte(world.GetFlag()),
			EventMessage: world.GetEventMessage(),
		}
		if err := ctx.Client.Send(serverResp, types.SEND_POLICY_ENCRYPT); err != nil {
			return err
		}
	}

	return ctx.Client.Send(&response.EndOfServerList{}, types.SEND_POLICY_ENCRYPT)
}
