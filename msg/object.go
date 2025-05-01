package msg

import (
	"github.com/asynkron/protoactor-go/actor"
)

type ObjectMove struct {
	X int
	Y int
}

type Warped struct {
	Sender *actor.PID
}
