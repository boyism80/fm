package actor

import (
	"time"

	"github.com/boyism80/fm/core/client"
	"github.com/boyism80/fm/types"
)

type HandlePacket struct {
	Opcode int
	Data   []byte
	Client client.Client
}

type ScheduleTimer struct {
	Interval time.Duration
	Logic    func() error
}

type ExecuteTimer struct {
	Logic func() error
}

type PacketResponse struct {
	Packet types.Packet
	Policy types.SendPolicy
	Error  error
}
