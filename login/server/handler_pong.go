package server

import (
	"log"

	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/protocol/request"
)

// Pong handles pong packet requests
type Pong struct {
	loginServer *LoginServer
	opcode      byte
}

func (Pong) New(loginServer *LoginServer) *Pong {
	return &Pong{
		loginServer: loginServer,
		opcode:      0x0A,
	}
}

func (h *Pong) GetOpcode() byte {
	return h.opcode
}

func (h *Pong) Handle(ctx *core.ClientContext, req *request.Pong) error {
	log.Printf("Pong packet received from %s", ctx.Client.GetConnection().RemoteAddr())
	return nil
}

