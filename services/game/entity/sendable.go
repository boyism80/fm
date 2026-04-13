package entity

import "github.com/boyism80/fm/types"

type Sendable interface {
	Send(p types.Packet, policy types.SendPolicy) error
}
