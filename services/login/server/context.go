package server

import (
	"github.com/boyism80/fm/core"
	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
)

type LoginServerContext struct {
	packetHandler  *core.PacketHandler
	InternalClient internal.InternalClient
}

func NewLoginServerContext(internalClient internal.InternalClient) *LoginServerContext {
	return &LoginServerContext{
		packetHandler:  core.NewPacketHandler(),
		InternalClient: internalClient,
	}
}

func (c *LoginServerContext) GetPacketHandler() *core.PacketHandler {
	return c.packetHandler
}
