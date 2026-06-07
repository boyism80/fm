package actor

import (
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/core/client"
	"github.com/boyism80/fm/services/game/constant"
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

type RunObjectTimer struct {
	ObjectType constant.ObjectType
	ID         uint32
	Key        string
}

type RunReactorRespawn struct {
	SpawnID uint32
}

type PacketResponse struct {
	Packet types.Packet
	Policy types.SendPolicy
	Error  error
}
