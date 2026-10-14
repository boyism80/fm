package route

import "github.com/asynkron/protoactor-go/actor"

// Envelope wraps an inner actor message sent to the actor that owns a character.
type Envelope struct {
	CorrelationID    uint64
	CharacterID      uint32
	DeadlineUnixNano int64
	Caller           *actor.PID
	Inner            interface{}
}

// Result is sent back to Caller when an envelope is handled or fails.
type Result struct {
	CorrelationID uint64
	OK            bool
	Reason        string
}

// Coordinator is implemented by the game server (or login server) to complete or resend envelopes.
type Coordinator interface {
	Resend(d *Envelope)
	CompleteEnvelope(correlationID uint64)
}
