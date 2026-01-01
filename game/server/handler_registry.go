package server

type PacketHandlerRegistry struct {
	gameServer *GameServer
}

func NewPacketHandlerRegistry(gameServer *GameServer) *PacketHandlerRegistry {
	return &PacketHandlerRegistry{
		gameServer: gameServer,
	}
}
