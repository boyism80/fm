package entity

import (
	"github.com/boyism80/fm/game/data"
)

// GameContext provides access to game resources and services
type GameContext interface {
	GetResources() *data.Resources
	GetMap(mapId uint32) *Map
	GetLogicThread() interface{} // Returns core.LogicThread
}
