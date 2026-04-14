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

type DeleteCharacter struct {
	ls     *LoginServer
	opcode byte
}

func (DeleteCharacter) New(ls *LoginServer) *DeleteCharacter {
	return &DeleteCharacter{
		ls:     ls,
		opcode: 0x09,
	}
}

func (h *DeleteCharacter) GetOpcode() byte {
	return h.opcode
}

func (h *DeleteCharacter) Handle(ctx *core.ClientContext, req *request.DeleteCharacter) error {
	log.Printf("Delete character packet received from %s - Character ID: %d",
		ctx.Client.GetConnection().RemoteAddr(), req.ID)

	loginClient, ok := ctx.Client.(*client.LoginClient)
	if !ok {
		return nil
	}

	ic := h.ls.context.InternalClient
	if ic == nil {
		deleteResp := &response.DeleteCharacter{ID: req.ID, Success: false}
		return ctx.Client.Send(deleteResp, types.SEND_POLICY_ENCRYPT)
	}

	reply, err := ic.DeleteCharacter(context.Background(), &fminternalpb.DeleteCharacterRequest{
		AccountId:   loginClient.GetAccountId(),
		CharacterId: req.ID,
	})
	if err != nil {
		log.Printf("DeleteCharacter RPC error: %v", err)
		deleteResp := &response.DeleteCharacter{ID: req.ID, Success: false}
		return ctx.Client.Send(deleteResp, types.SEND_POLICY_ENCRYPT)
	}

	deleteResp := &response.DeleteCharacter{
		ID:      req.ID,
		Success: reply.Success,
	}
	return ctx.Client.Send(deleteResp, types.SEND_POLICY_ENCRYPT)
}
