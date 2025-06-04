package entity

import "github.com/boyism80/fm/common/types"

type Sendable interface {
	Send(p types.Packet, policy types.SendPolicy)
}
