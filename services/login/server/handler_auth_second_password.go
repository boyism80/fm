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

type AuthSecondPassword struct {
	ls *LoginServer
}

func (AuthSecondPassword) New(ls *LoginServer) *AuthSecondPassword {
	return &AuthSecondPassword{ls: ls}
}

func (h *AuthSecondPassword) Handle(ctx *core.ClientContext, req *request.AuthSecondPassword) error {
	loginClient := ctx.Client.(*client.LoginClient)
	accountID := loginClient.GetAccountId()
	if accountID == 0 {
		return fmt.Errorf("second password: account id is missing")
	}

	if req.Mode != request.SecondPasswordModeSet && req.Mode != request.SecondPasswordModeRemove {
		return ctx.Client.Send(&response.SecondPasswordResult{
			Mode:   request.SecondPasswordModeSet,
			Result: response.SecondPasswordResultInvalid,
		}, types.SEND_POLICY_ENCRYPT)
	}

	async.NewTask(ctx.ActorContext, core.InternalRPCPerStepTimeout).ThenRPC(
		func(c context.Context) (*internal.SecondPasswordReply, error) {
			if req.Mode == request.SecondPasswordModeSet {
				return h.ls.internalClient.SetSecondPassword(c, &internal.SetSecondPasswordRequest{
					WorldId:   loginClient.GetWorldId(),
					AccountId: accountID,
					Password:  req.Password,
				})
			}
			return h.ls.internalClient.RemoveSecondPassword(c, &internal.RemoveSecondPasswordRequest{
				WorldId:   loginClient.GetWorldId(),
				AccountId: accountID,
				Password:  req.Password,
			})
		}, func(reply *internal.SecondPasswordReply) error {
			result := uint8(response.SecondPasswordResultInvalid)
			mode := uint8(request.SecondPasswordModeSet)
			if reply.GetStatus() == internal.SecondPasswordStatus_SECOND_PASSWORD_SUCCESS {
				result = response.SecondPasswordResultSuccess
				mode = req.Mode
			}
			return ctx.Client.Send(&response.SecondPasswordResult{Mode: mode, Result: result}, types.SEND_POLICY_ENCRYPT)
		}).OnError(func(err error) {
		log.Printf("AuthSecondPassword (async): account=%d: %v", accountID, err)
		_ = ctx.Client.Send(&response.SecondPasswordResult{
			Mode:   request.SecondPasswordModeSet,
			Result: response.SecondPasswordResultInvalid,
		}, types.SEND_POLICY_ENCRYPT)
	})
	return nil
}
