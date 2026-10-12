package delivery

import "github.com/asynkron/protoactor-go/actor"

// Envelope wraps an inner actor message with correlation and deadline for Deliver delivery.
type Envelope struct {
	CorrelationID    uint64
	CharacterID      uint32
	DeadlineUnixNano int64
	Caller           *actor.PID
	Inner            interface{}
}

// Result is sent back to Caller when a delivery finishes or fails.
type Result struct {
	CorrelationID uint64
	OK            bool
	Reason        string
}

// Coordinator is implemented by the game server (or login server) to complete or retry deliveries.
type Coordinator interface {
	Redeliver(d *Envelope)
	CompleteDelivery(correlationID uint64)
}
