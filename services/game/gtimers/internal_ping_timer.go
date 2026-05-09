package gtimers

import (
	"context"
	"time"

	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
)

const defaultPingTimeout = 10 * time.Second

var internalPingWire struct {
	client    internal.InternalClient
	interval  time.Duration
	worldID   uint32
	channelID uint32
}

func WireInternalPing(client internal.InternalClient, interval time.Duration, worldID uint32, channelID uint32) {
	internalPingWire.client = client
	internalPingWire.interval = interval
	internalPingWire.worldID = worldID
	internalPingWire.channelID = channelID
}

type InternalPingTimer struct{}

func (*InternalPingTimer) New() *InternalPingTimer {
	return &InternalPingTimer{}
}

func (*InternalPingTimer) GetName() string {
	return "game_internal_ping"
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
		Role:      internal.ServerRole_SERVER_ROLE_GAME,
		WorldId:   internalPingWire.worldID,
		ChannelId: internalPingWire.channelID,
	}
	pingCtx, cancel := context.WithTimeout(ctx, defaultPingTimeout)
	defer cancel()
	_, err := client.Ping(pingCtx, req)
	return err
}
