package server

import (
	"context"
	"fmt"
	"time"

	"github.com/boyism80/fm/core/async"
	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const pingAsyncPerStepTimeout = 10 * time.Second

// InternalPingAsync builds a Promise that dials address and calls internal Ping (reply message must be "pong").
// Empty address returns an idle Promise. Caller must Run() (and typically OnError / Finally).
func InternalPingAsync(address string) *async.Promise {
	p := async.NewPromise(nil, pingAsyncPerStepTimeout)
	if address == "" {
		return p
	}
	addr := address
	async.ThenRPC(p, func(c context.Context) (*internal.PingReply, error) {
		conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			return nil, fmt.Errorf("internal grpc dial %q: %w", addr, err)
		}
		defer func() { _ = conn.Close() }()
		client := internal.NewInternalClient(conn)
		return client.Ping(c, &internal.PingRequest{})
	}, func(reply *internal.PingReply) error {
		if reply == nil || reply.GetMessage() != "pong" {
			return fmt.Errorf("internal Ping: want message %q, got %q", "pong", reply.GetMessage())
		}
		return nil
	})
	return p
}
