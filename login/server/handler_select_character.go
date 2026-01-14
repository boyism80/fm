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
	ls     *LoginServer
	opcode byte
}

func (SelectCharacter) New(ls *LoginServer) *SelectCharacter {
	return &SelectCharacter{
		ls:     ls,
		opcode: 0x05,
	}
}

func (h *SelectCharacter) GetOpcode() byte {
	return h.opcode
}

func (h *SelectCharacter) Handle(ctx *core.ClientContext, req *request.SelectCharacter) error {
	log.Printf("Select character packet received from %s - Character ID: %d",
		ctx.Client.GetConnection().RemoteAddr(), req.CharacterId)

	transferResp := &response.Transfer{
		IP:          h.ls.config.GameServerHost,
		Port:        uint16(h.ls.config.GameServerPort),
		CharacterId: req.CharacterId,
	}
	if err := ctx.Client.Send(transferResp, types.SEND_POLICY_ENCRYPT); err != nil {
		log.Printf("Failed to send transfer response: %v", err)
		return err
	}

	return nil
}
