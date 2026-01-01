package server

import (
	"fmt"
	"log"

	"github.com/boyism80/fm/game/client"
	"github.com/boyism80/fm/game/constant"
	"github.com/boyism80/fm/game/entity"
)

type KillAllMobs struct {
	gameServer *GameServer
}

func (*KillAllMobs) New(gameServer *GameServer) *KillAllMobs {
	return &KillAllMobs{
		gameServer: gameServer,
	}
}

func (h *KillAllMobs) GetCommandName() string {
	return "몬스터죽이기"
}

func (h *KillAllMobs) Handle(gameClient *client.GameClient, args ...string) error {
	character := gameClient.GetCharacter()
	if character == nil {
		return fmt.Errorf("character not found")
	}

	mapInstance := h.gameServer.GetMap(character.Map)
	if mapInstance == nil {
		return fmt.Errorf("map %d not found", character.Map)
	}

	mobs := mapInstance.GetMobs()
	for mobID, mobInterface := range mobs {
		if _, ok := mobInterface.(*entity.Mob); ok {
			if err := mapInstance.RemoveMob(mobID, constant.MOB_DIE_ANIMATION_TYPE_FADE_OUT); err != nil {
				log.Printf("Failed to remove mob %d: %v", mobID, err)
			} else {
				log.Printf("Command: Killed mob OID: %d", mobID)
			}
		}
	}

	return nil
}

