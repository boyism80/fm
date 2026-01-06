package server

import (
	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/game/entity"
	"github.com/boyism80/fm/game/wz"
)

type GameServerContext struct {
	packetHandler *core.PacketHandler
	resources     *wz.Resources
	gameServer    GameServerInterface
	mapListener   MapListenerInterface
}

type GameServerInterface interface {
	GetMap(mapID uint32) *entity.Map
}

type MapListenerInterface interface {
	OnWarpCharacter(character *entity.Character, targetMapID uint32, portal uint8)
}

func NewGameServerContext(wzPath string, gameServer GameServerInterface, mapListener MapListenerInterface) *GameServerContext {
	resources := wz.NewResources(wzPath)
	return &GameServerContext{
		packetHandler: core.NewPacketHandler(),
		resources:     resources,
		gameServer:    gameServer,
		mapListener:   mapListener,
	}
}

func (c *GameServerContext) GetMapListener() MapListenerInterface {
	return c.mapListener
}

func (c *GameServerContext) GetPacketHandler() *core.PacketHandler {
	return c.packetHandler
}

func (c *GameServerContext) GetResources() *wz.Resources {
	return c.resources
}

func (c *GameServerContext) GetMapActorPID(mapID uint32) *actor.PID {
	if c.gameServer == nil {
		return nil
	}
	mapInstance := c.gameServer.GetMap(mapID)
	if mapInstance == nil {
		return nil
	}
	return mapInstance.GetActorPID()
}
