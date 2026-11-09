package rpc

import (
	"errors"

	"github.com/asynkron/protoactor-go/actor"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// OutcomeUnknown reports whether an RPC that failed with err may still have been applied by internal: its reply never
// came back, or the promise's future gave up before it did (a future timeout is not a gRPC status).
func OutcomeUnknown(err error) bool {
	if errors.Is(err, actor.ErrTimeout) {
		return true
	}
	switch status.Code(err) {
	case codes.DeadlineExceeded, codes.Unavailable, codes.Canceled:
		return true
	}
	return false
}
