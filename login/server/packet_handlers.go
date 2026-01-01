package server

import (
	"github.com/boyism80/fm/core"
)

// registerPacketHandlers registers all login server packet handlers
func (ls *LoginServer) registerPacketHandlers() {
	core.Bind[*LoginServer, Pong](ls)
	core.Bind[*LoginServer, Login](ls)
	core.Bind[*LoginServer, CharacterList](ls)
	core.Bind[*LoginServer, CheckName](ls)
	core.Bind[*LoginServer, CreateCharacter](ls)
	core.Bind[*LoginServer, DeleteCharacter](ls)
	core.Bind[*LoginServer, SelectCharacter](ls)
}
