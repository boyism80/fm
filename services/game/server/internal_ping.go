package server

import (
	"context"
	"fmt"
	"time"

	fminternalpb "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// PingInternalService dials the internal gRPC server and calls Ping. address is host:port.
// Empty address skips the check.
func PingInternalService(ctx context.Context, address string) error {
	if address == "" {
		return nil
	}

	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return fmt.Errorf("internal grpc dial %q: %w", address, err)
	}
	defer func() { _ = conn.Close() }()

	client := fminternalpb.NewInternalClient(conn)
	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	reply, err := client.Ping(pingCtx, &fminternalpb.PingRequest{})
	if err != nil {
		return fmt.Errorf("internal Ping: %w", err)
	}
	if reply.GetMessage() != "pong" {
		return fmt.Errorf("internal Ping: want message %q, got %q", "pong", reply.GetMessage())
	}
	return nil
}
