package entity

import (
	"time"

	"github.com/boyism80/fm/game/data"
)

// LogicThread interface to avoid import cycle
type LogicThread interface {
	Schedule(delay time.Duration, task func()) interface{}                             // Returns timer object
	ScheduleAtFixedRate(initialDelay, period time.Duration, task func()) interface{}   // Returns timer object
	ScheduleWithFixedDelay(initialDelay, delay time.Duration, task func()) interface{} // Returns timer object
}

// GameContext provides access to game resources and services
type GameContext interface {
	GetResources() *data.Resources
	GetMap(mapId uint32) *Map
	GetLogicThread() LogicThread
}
