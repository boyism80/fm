package server

import (
	"log"

	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/protocol/request"
)

type SelectCharacterWithSecondPassword struct {
	selection *SelectCharacter
}

func (SelectCharacterWithSecondPassword) New(ls *LoginServer) *SelectCharacterWithSecondPassword {
	return &SelectCharacterWithSecondPassword{
		selection: &SelectCharacter{ls: ls},
	}
}

func (h *SelectCharacterWithSecondPassword) Handle(ctx *core.ClientContext, req *request.SelectCharacterWithSecondPassword) error {
	log.Printf("Select character with second password packet received from %s - Character ID: %d",
		ctx.Client.GetConnection().RemoteAddr(), req.CharacterID)
	return h.selection.selectCharacter(ctx, req.CharacterID, &req.Password)
}
