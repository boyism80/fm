package core

import "github.com/boyism80/fm/core/delivery"

type Server interface {
	delivery.Coordinator
	GetPacketHandler() *PacketHandler
}
