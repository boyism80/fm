package server

import (
	"log"

	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/protocol/request"
)

type Pong struct {
	ls     *LoginServer
	opcode byte
}

func (Pong) New(ls *LoginServer) *Pong {
	return &Pong{
		ls:     ls,
		opcode: 0x0A,
	}
}

func (h *Pong) GetOpcode() byte {
	return h.opcode
}

func (h *Pong) Handle(ctx *core.ClientContext, req *request.Pong) error {
	log.Printf("Pong packet received from %s", ctx.Client.GetConnection().RemoteAddr())
	return nil
}
