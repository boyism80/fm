package ensure

import "github.com/asynkron/protoactor-go/actor"

// EnsureDeliver wraps an inner actor message with correlation and deadline for EnsureSend delivery.
type EnsureDeliver struct {
	CorrelationID    uint64
	CharacterID      uint32
	DeadlineUnixNano int64
	Caller           *actor.PID
	Inner            interface{}
}

// EnsureResult is sent back to Caller when an ensured delivery finishes or fails.
type EnsureResult struct {
	CorrelationID uint64
	OK            bool
	Reason        string
}

// EnsureCoordinator is implemented by the game server (or login server) to complete or retry ensured deliveries.
type EnsureCoordinator interface {
	EnsureRedispatch(d *EnsureDeliver)
	EnsureComplete(correlationID uint64)
}
