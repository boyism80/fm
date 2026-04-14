package server

import (
	"context"
	"log"

	"github.com/boyism80/fm/core"
	fminternalpb "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
	"github.com/boyism80/fm/protocol/request"
	"github.com/boyism80/fm/protocol/response"
	"github.com/boyism80/fm/types"
)

type CheckName struct {
	ls     *LoginServer
	opcode byte
}

func (CheckName) New(ls *LoginServer) *CheckName {
	return &CheckName{
		ls:     ls,
		opcode: 0x07,
	}
}

func (h *CheckName) GetOpcode() byte {
	return h.opcode
}

func (h *CheckName) Handle(ctx *core.ClientContext, req *request.CheckName) error {
	log.Printf("Check name packet received from %s - Name: %s",
		ctx.Client.GetConnection().RemoteAddr(), req.Name)

	ic := h.ls.context.InternalClient
	exists := false
	if ic != nil {
		reply, err := ic.CheckCharacterName(context.Background(), &fminternalpb.CheckCharacterNameRequest{
			Name: req.Name,
		})
		if err != nil {
			log.Printf("CheckCharacterName RPC error: %v", err)
		} else {
			exists = reply.Exists
		}
	}

	checkResp := &response.CheckName{
		Name:   req.Name,
		Exists: exists,
	}
	return ctx.Client.Send(checkResp, types.SEND_POLICY_ENCRYPT)
}
