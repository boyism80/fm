package entity

import (
	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/game/wz"
)

// GameContext provides access to game resources and services
type GameContext interface {
	GetResources() *wz.Resources
	GetMap(mapId uint32) *Map
	GetLogicThread() *core.LogicThread // Returns core.LogicThread
}
