package server

import (
	"context"
	"fmt"
	"log"

	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/core/async"
	"github.com/boyism80/fm/protocol/dto"
	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
	"github.com/boyism80/fm/protocol/request"
	"github.com/boyism80/fm/protocol/response"
	"github.com/boyism80/fm/services/login/client"
	"github.com/boyism80/fm/types"
)

type CreateCharacter struct {
	ls *LoginServer
}

func (CreateCharacter) New(ls *LoginServer) *CreateCharacter {
	return &CreateCharacter{
		ls: ls,
	}
}

func (h *CreateCharacter) Handle(ctx *core.ClientContext, req *request.CreateCharacter) error {
	log.Printf("Create character packet received from %s - Name: %s",
		ctx.Client.GetConnection().RemoteAddr(), req.Name)

	loginClient, ok := ctx.Client.(*client.LoginClient)
	if !ok {
		return nil
	}

	ic := h.ls.internalClient
	if ic == nil {
		createResp := &response.CreateCharacter{Success: false}
		return ctx.Client.Send(createResp, types.SEND_POLICY_ENCRYPT)
	}

	reqProto := &internal.CreateCharacterRequest{
		AccountId:    loginClient.GetAccountId(),
		WorldId:      loginClient.GetWorldId(),
		Name:         req.Name,
		Face:         req.Face,
		Hair:         req.Hair,
		SkinColor:    0,
		TopItemId:    req.Top,
		BottomItemId: req.Bottom,
		ShoesItemId:  req.Shoes,
		WeaponItemId: req.Weapon,
	}

	if ctx.ActorContext == nil {
		log.Printf("CreateCharacter: no actor context, cannot run internal RPC")
		createResp := &response.CreateCharacter{Success: false}
		_ = ctx.Client.Send(createResp, types.SEND_POLICY_ENCRYPT)
		return fmt.Errorf("create character: actor context required for internal RPC")
	}

	async.ThenRPC(async.NewPromise(ctx.ActorContext, core.InternalRPCPerStepTimeout),
		func(c context.Context) (*internal.CreateCharacterReply, error) {
			return ic.CreateCharacter(c, reqProto)
		}, func(reply *internal.CreateCharacterReply) error {
			return h.sendCreateCharacterResult(ctx, reply)
		}).
		OnError(func(err error) {
			log.Printf("CreateCharacter (async): %v", err)
			createResp := &response.CreateCharacter{Success: false}
			_ = ctx.Client.Send(createResp, types.SEND_POLICY_ENCRYPT)
		}).
		Run()
	return nil
}

func (h *CreateCharacter) sendCreateCharacterResult(ctx *core.ClientContext, reply *internal.CreateCharacterReply) error {
	createResp := &response.CreateCharacter{
		Success: reply.Success,
	}
	if reply.Success && reply.Character != nil {
		ch := overviewToDto(reply.Character)
		createResp.Character = &ch
	} else {
		createResp.Character = &dto.Character{
			BaseLooks: make(map[int8]uint32),
			Overlays:  make(map[int8]uint32),
		}
	}
	return ctx.Client.Send(createResp, types.SEND_POLICY_ENCRYPT)
}
