package server

import (
	"fmt"
	"strconv"

	g_actor "github.com/boyism80/fm/game/actor"
	"github.com/boyism80/fm/game/client"
)

type ChangeMap struct {
	gs *GameServer
}

func (*ChangeMap) New(gs *GameServer) *ChangeMap {
	return &ChangeMap{
		gs: gs,
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
		mapIdUint, ok := h.gs.resources.NameToMap(args[0])
		if !ok {
			return fmt.Errorf("invalid mapId or map name: %s", args[0])
		}
		mapId = int(mapIdUint)
	}

	character := gameClient.GetCharacter()
	if character == nil {
		return fmt.Errorf("character not found")
	}

	currentMap := h.gs.GetMap(character.GetMap())
	if currentMap == nil {
		return fmt.Errorf("current map not found")
	}

	targetMap := h.gs.GetMap(uint32(mapId))
	if targetMap == nil {
		return fmt.Errorf("target map %d not found", mapId)
	}

	targetMapPID := targetMap.GetActorPID()
	if targetMapPID == nil {
		return fmt.Errorf("target map actor PID not found")
	}

	currentMap.RemovePlayer(character.GetID())

	rootContext := h.gs.GetServer().GetRootContext()
	if rootContext == nil {
		return fmt.Errorf("rootContext not set")
	}

	rootContext.Send(targetMapPID, &g_actor.WarpCharacter{
		Character: character,
		Portal:    1,
	})

	return nil
}
