package core

import "github.com/boyism80/fm/core/route"

type Server interface {
	route.Coordinator
	GetPacketHandler() *PacketHandler
}
