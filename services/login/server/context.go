package server

import (
	"github.com/boyism80/fm/core"
)

type LoginServerContext struct {
	packetHandler *core.PacketHandler
}

func NewLoginServerContext() *LoginServerContext {
	return &LoginServerContext{
		packetHandler: core.NewPacketHandler(),
	}
}

func (c *LoginServerContext) GetPacketHandler() *core.PacketHandler {
	return c.packetHandler
}
