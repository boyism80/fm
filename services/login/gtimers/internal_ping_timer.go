package gtimers

import (
	"context"
	"time"

	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
)

const defaultPingTimeout = 10 * time.Second

var internalPingWire struct {
	client          internal.InternalClient
	interval        time.Duration
	loginInstanceID string
}

func WireInternalPing(client internal.InternalClient, interval time.Duration, loginInstanceID string) {
	internalPingWire.client = client
	internalPingWire.interval = interval
	internalPingWire.loginInstanceID = loginInstanceID
}

type InternalPingTimer struct{}

func (*InternalPingTimer) New() *InternalPingTimer {
	return &InternalPingTimer{}
}

func (*InternalPingTimer) GetName() string {
	return "login_internal_ping"
}

func (*InternalPingTimer) GetInterval() time.Duration {
	return internalPingWire.interval
}

func (*InternalPingTimer) GetInitialDelay() time.Duration {
	return internalPingWire.interval
}

func (*InternalPingTimer) Handle(ctx context.Context) error {
	client := internalPingWire.client
	if client == nil || internalPingWire.interval <= 0 {
		return nil
	}
	req := &internal.PingRequest{
		Role:            internal.ServerRole_SERVER_ROLE_LOGIN,
		LoginInstanceId: internalPingWire.loginInstanceID,
	}
	pingCtx, cancel := context.WithTimeout(ctx, defaultPingTimeout)
	defer cancel()
	_, err := client.Ping(pingCtx, req)
	return err
}
