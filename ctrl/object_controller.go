package ctrl

import (
	"github.com/boyism80/fm/types"
)

type ObjectController interface {
	ID() int64
	Position() types.Vec2
	Type() types.ObjectType
	Send(p types.Packet, policy types.SendPolicy)
}
