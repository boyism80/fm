package core

import "github.com/boyism80/fm/core/ensure"

type Server interface {
	ensure.EnsureCoordinator
	GetPacketHandler() *PacketHandler
}
