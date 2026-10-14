package rpc

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/core/fault"
	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type InternalConnection struct {
	Client    internal.InternalClient
	Faults    *fault.Injector
	conn      *grpc.ClientConn
	heartbeat time.Duration
	ping      func() *internal.PingRequest
	cancel    context.CancelFunc
}

// NewInternalConnection dials the internal server. With a nil ping or a zero heartbeat the connection sends no heartbeat.
func NewInternalConnection(addr string, heartbeat time.Duration, ping func() *internal.PingRequest) (*InternalConnection, error) {
	if addr == "" {
		return nil, fmt.Errorf("internal gRPC endpoint is required")
	}
	faults := fault.NewInjector()
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithUnaryInterceptor(faults.UnaryClientInterceptor()))
	if err != nil {
		return nil, fmt.Errorf("internal gRPC dial %s: %w", addr, err)
	}
	return &InternalConnection{
		Client:    internal.NewInternalClient(conn),
		Faults:    faults,
		conn:      conn,
		heartbeat: heartbeat,
		ping:      ping,
	}, nil
}

func (s *InternalConnection) Start() error {
	if s.heartbeat <= 0 || s.ping == nil {
		return nil
	}
	ctx, cancel := context.WithCancel(context.Background())
	s.cancel = cancel
	go s.beat(ctx)
	return nil
}

func (s *InternalConnection) Stop() {
	if s.cancel != nil {
		s.cancel()
	}
	_ = s.conn.Close()
}

func (s *InternalConnection) beat(ctx context.Context) {
	ticker := time.NewTicker(s.heartbeat)
	defer ticker.Stop()
	for {
		pingCtx, cancel := context.WithTimeout(ctx, core.InternalRPCPerStepTimeout)
		_, err := s.Client.Ping(pingCtx, s.ping())
		cancel()
		if err != nil {
			log.Printf("internal ping: %v", err)
		}

		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
		}
	}
}
