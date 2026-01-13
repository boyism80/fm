package actor

import (
	"github.com/boyism80/fm/game/entity"
)

type AddCharacter struct {
	Character  *entity.Character
	SpawnPoint uint8
	Init       bool
}

type RemoveCharacter struct {
	CharacterID uint32
}

type WarpCharacter struct {
	Character *entity.Character
	Portal    uint8
}
