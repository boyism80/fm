package server

import (
	"log"

	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/protocol/request"
	"github.com/boyism80/fm/protocol/response"
	"github.com/boyism80/fm/types"
)

// SelectCharacter handles select character packet requests
type SelectCharacter struct {
	loginServer *LoginServer
	opcode      byte
}

func (SelectCharacter) New(loginServer *LoginServer) *SelectCharacter {
	return &SelectCharacter{
		loginServer: loginServer,
		opcode:      0x05,
	}
}

func (h *SelectCharacter) GetOpcode() byte {
	return h.opcode
}

func (h *SelectCharacter) Handle(ctx *core.ClientContext, req *request.SelectCharacter) error {
	log.Printf("Select character packet received from %s - Character ID: %d",
		ctx.Client.GetConnection().RemoteAddr(), req.CharacterId)

	transferResp := &response.Transfer{
		IP:          h.loginServer.config.GameServerHost,
		Port:        uint16(h.loginServer.config.GameServerPort),
		CharacterId: req.CharacterId,
	}
	if err := ctx.SendFunc(transferResp, types.SEND_POLICY_ENCRYPT); err != nil {
		log.Printf("Failed to send transfer response: %v", err)
		return err
	}

	return nil
}

