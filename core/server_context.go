package core

type Server interface {
	GetPacketHandler() *PacketHandler
}
