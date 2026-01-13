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
}

type GameServerInterface interface {
	GetMap(mapID uint32) *entity.Map
}

func NewGameServerContext(wzPath string, gameServer GameServerInterface) *GameServerContext {
	resources := wz.NewResources(wzPath)
	return &GameServerContext{
		packetHandler: core.NewPacketHandler(),
		resources:     resources,
		gameServer:    gameServer,
	}
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
