package actor

import (
	"github.com/boyism80/fm/game/entity"
)

type AddCharacter struct {
	Character *entity.Character
	Init      bool
}

type RemoveCharacter struct {
	CharacterID uint32
}

type WarpCharacter struct {
	Character *entity.Character
	TargetMap uint32
	Portal    uint8
}
