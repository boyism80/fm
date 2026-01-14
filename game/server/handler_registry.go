package server

type PacketHandlerRegistry struct {
	gs *GameServer
}

func NewPacketHandlerRegistry(gs *GameServer) *PacketHandlerRegistry {
	return &PacketHandlerRegistry{
		gs: gs,
	}
}
