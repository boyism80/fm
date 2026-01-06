package core

type ServerContext interface {
	GetPacketHandler() *PacketHandler
}
