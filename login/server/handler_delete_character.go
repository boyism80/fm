package server

import (
	"log"

	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/protocol/request"
	"github.com/boyism80/fm/protocol/response"
	"github.com/boyism80/fm/types"
)

// DeleteCharacter handles delete character packet requests
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

	deleteResp := &response.DeleteCharacter{
		ID:      req.ID,
		Success: true,
	}
	if err := ctx.Client.Send(deleteResp, types.SEND_POLICY_ENCRYPT); err != nil {
		log.Printf("Failed to send delete character response: %v", err)
		return err
	}

	return nil
}
