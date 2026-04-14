package server

import (
	"context"
	"log"

	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/protocol/dto"
	fminternalpb "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
	"github.com/boyism80/fm/protocol/request"
	"github.com/boyism80/fm/protocol/response"
	"github.com/boyism80/fm/services/login/client"
	"github.com/boyism80/fm/types"
)

type CreateCharacter struct {
	ls     *LoginServer
	opcode byte
}

func (CreateCharacter) New(ls *LoginServer) *CreateCharacter {
	return &CreateCharacter{
		ls:     ls,
		opcode: 0x08,
	}
}

func (h *CreateCharacter) GetOpcode() byte {
	return h.opcode
}

func (h *CreateCharacter) Handle(ctx *core.ClientContext, req *request.CreateCharacter) error {
	log.Printf("Create character packet received from %s - Name: %s",
		ctx.Client.GetConnection().RemoteAddr(), req.Name)

	loginClient, ok := ctx.Client.(*client.LoginClient)
	if !ok {
		return nil
	}

	ic := h.ls.context.InternalClient
	if ic == nil {
		createResp := &response.CreateCharacter{Success: false}
		return ctx.Client.Send(createResp, types.SEND_POLICY_ENCRYPT)
	}

	reply, err := ic.CreateCharacter(context.Background(), &fminternalpb.CreateCharacterRequest{
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
	})
	if err != nil {
		log.Printf("CreateCharacter RPC error: %v", err)
		createResp := &response.CreateCharacter{Success: false}
		return ctx.Client.Send(createResp, types.SEND_POLICY_ENCRYPT)
	}

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
