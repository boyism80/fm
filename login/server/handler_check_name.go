package server

import (
	"log"

	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/protocol/request"
	"github.com/boyism80/fm/protocol/response"
	"github.com/boyism80/fm/types"
)

type CheckName struct {
	loginServer *LoginServer
	opcode      byte
}

func (CheckName) New(loginServer *LoginServer) *CheckName {
	return &CheckName{
		loginServer: loginServer,
		opcode:      0x07,
	}
}

func (h *CheckName) GetOpcode() byte {
	return h.opcode
}

func (h *CheckName) Handle(ctx *core.ClientContext, req *request.CheckName) error {
	log.Printf("Check name packet received from %s - Name: %s",
		ctx.Client.GetConnection().RemoteAddr(), req.Name)

	// Check if name exists (hardcoded for demo)
	exists := req.Name == "채승현"

	checkResp := &response.CheckName{
		Name:   req.Name,
		Exists: exists,
	}
	if err := ctx.Client.Send(checkResp, types.SEND_POLICY_ENCRYPT); err != nil {
		log.Printf("Failed to send check name response: %v", err)
		return err
	}

	return nil
}
