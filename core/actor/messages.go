package actor

import (
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/core/client"
	"github.com/boyism80/fm/types"
)

type HandlePacket struct {
	Opcode        int
	Data          []byte
	Client        client.Client
	LogicActorPID *actor.PID
}

type ScheduleTimer struct {
	Interval time.Duration
	Logic    func() error
}

type ExecuteTimer struct {
	Logic func() error
}

type RunCharacterTimer struct {
	CharacterID uint32
	Key         string
}

type PacketResponse struct {
	Packet types.Packet
	Policy types.SendPolicy
	Error  error
}
