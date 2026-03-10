package entity

import (
	// "github.com/boyism80/fm/core" // Commented out: LogicThread removed
	"github.com/boyism80/fm/game/wz"
)

// GameContext provides access to game resources and services
type GameContext interface {
	GetResources() *wz.Resources
	GetMap(mapId uint32) *Map
	GetExpRate() int  // Returns experience rate multiplier
	GetDropRate() int // Returns drop rate multiplier
	GetMesoRate() int // Returns meso rate multiplier
	RequestWarp(character *Character, targetMap *Map, spawnPoint uint8) error
}
