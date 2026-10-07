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
	ls *LoginServer
}

func (Login) New(ls *LoginServer) *Login {
	return &Login{
		ls: ls,
	}
}

func (h *Login) Handle(ctx *core.ClientContext, req *request.Login) error {
	log.Printf("Login packet received from %s - ID: %s, MAC: %s",
		ctx.Client.GetConnection().RemoteAddr(), req.ID, req.Mac)

	remoteAddr := ctx.Client.GetConnection().RemoteAddr().String()
	loginClient := ctx.Client.(*client.LoginClient)
	if loginClient.GetAccountId() != 0 {
		log.Printf("Login: packet on authenticated connection remote=%s id=%s account=%d", remoteAddr, req.ID, loginClient.GetAccountId())
	}
	reqMsg := &internal.LoginAccountRequest{
		LoginId:     req.ID,
		Password:    req.Pw,
		MacAddress:  req.Mac,
		IpAddress:   remoteAddr,
		InitialRole: h.ls.config.InitialRole,
	}

	async.NewTask(ctx.ActorContext, core.InternalRPCPerStepTimeout).ThenRPC(
		func(c context.Context) (*internal.LoginAccountReply, error) {
			return h.ls.internalClient.LoginAccount(c, reqMsg)
		}, func(reply *internal.LoginAccountReply) error {
			log.Printf("Login: reply remote=%s id=%s account=%d status=%v", remoteAddr, req.ID, reply.GetAccountId(), reply.GetStatus())
			return h.sendLoginAccountResult(ctx, loginClient, req, reply)
		}).OnError(func(err error) {
		log.Printf("Login (async): remote=%s id=%s: %v", remoteAddr, req.ID, err)
		_ = ctx.Client.Send(&response.LoginFailed{Reason: response.LoginFailedReasonSystemError6}, types.SEND_POLICY_ENCRYPT)
	})
	return nil
}

func (h *Login) sendLoginAccountResult(ctx *core.ClientContext, loginClient *client.LoginClient, req *request.Login, reply *internal.LoginAccountReply) error {
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

	loginClient.SetAccountId(reply.AccountId)

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
