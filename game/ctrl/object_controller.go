package ctrl

import (
	"github.com/boyism80/fm/common/types"
)

type ObjectController interface {
	ID() int64
	Position() types.Vec2[int32]
	Type() types.ObjectType
	Send(p types.Packet, policy types.SendPolicy)
}
