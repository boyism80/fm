package server

import (
	"github.com/boyism80/fm/core"
	fminternalpb "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
)

type LoginServerContext struct {
	packetHandler  *core.PacketHandler
	InternalClient fminternalpb.InternalClient
}

func NewLoginServerContext(internalClient fminternalpb.InternalClient) *LoginServerContext {
	return &LoginServerContext{
		packetHandler:  core.NewPacketHandler(),
		InternalClient: internalClient,
	}
}

func (c *LoginServerContext) GetPacketHandler() *core.PacketHandler {
	return c.packetHandler
}
