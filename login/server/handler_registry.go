package server

type PacketHandlerRegistry struct {
	loginServer *LoginServer
}

func NewPacketHandlerRegistry(loginServer *LoginServer) *PacketHandlerRegistry {
	return &PacketHandlerRegistry{
		loginServer: loginServer,
	}
}
