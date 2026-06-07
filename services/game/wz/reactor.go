package wz

import (
	"github.com/boyism80/fm/services/game/constant"
	"github.com/boyism80/fm/types"
)

type Reactor struct {
	ID     uint32
	Info   ReactorInfo
	Action string
	States map[byte]*ReactorEvent
}

type ReactorInfo struct {
	Link            uint32
	ActivateByTouch int
}

type ReactorEvent struct {
	Type         constant.ReactorEventType
	NextState    byte
	TimeOut      int
	ItemID       int
	ItemQuantity int
	LT           types.Vector2[int32]
	RB           types.Vector2[int32]
	HasLT        bool
	HasRB        bool
	TouchFlag    int
	HasClickArea bool
}
