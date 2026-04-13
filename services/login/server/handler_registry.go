package server

type PacketHandlerRegistry struct {
	ls *LoginServer
}

func NewPacketHandlerRegistry(ls *LoginServer) *PacketHandlerRegistry {
	return &PacketHandlerRegistry{
		ls: ls,
	}
}
