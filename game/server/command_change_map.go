package server

import (
	"fmt"
	"strconv"

	"github.com/boyism80/fm/game/client"
)

type ChangeMap struct {
	gameServer *GameServer
}

func (*ChangeMap) New(gameServer *GameServer) *ChangeMap {
	return &ChangeMap{
		gameServer: gameServer,
	}
}

func (h *ChangeMap) GetCommandName() string {
	return "맵이동"
}

func (h *ChangeMap) GetUsage() string {
	return "<맵ID/이름> - 맵 이동"
}

func (h *ChangeMap) Handle(gameClient *client.GameClient, args ...string) error {
	if len(args) < 1 {
		return fmt.Errorf("missing mapId or map name")
	}

	var mapId int
	var err error

	mapId, err = strconv.Atoi(args[0])
	if err != nil {
		mapIdUint, ok := h.gameServer.resources.NameToMap(args[0])
		if !ok {
			return fmt.Errorf("invalid mapId or map name: %s", args[0])
		}
		mapId = int(mapIdUint)
	}

	character := gameClient.GetCharacter()
	if character == nil {
		return fmt.Errorf("character not found")
	}

	currentMap := h.gameServer.GetMap(character.GetMap())
	if currentMap == nil {
		return fmt.Errorf("current map not found")
	}

	targetMap := h.gameServer.GetMap(uint32(mapId))
	if targetMap == nil {
		return fmt.Errorf("target map %d not found", mapId)
	}

	currentMap.RemovePlayer(character.GetID())

	character.Map = uint32(mapId)

	mapSpec := targetMap.GetSpec()
	if mapSpec != nil && len(mapSpec.Portals) > 0 {
		for portalId, portal := range mapSpec.Portals {
			character.Position = portal.Position
			character.SpawnPoint = portalId
			break
		}
	}

	if err := targetMap.AddPlayer(character.GetID(), character, true); err != nil {
		return fmt.Errorf("failed to add character to new map: %v", err)
	}

	return nil
}

