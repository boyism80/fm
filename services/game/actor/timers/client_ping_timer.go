package timers

import (
	"github.com/boyism80/fm/core/clock"
	"log"
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/protocol/response"
	gameclient "github.com/boyism80/fm/services/game/client"
	"github.com/boyism80/fm/services/game/entity"
	"github.com/boyism80/fm/types"
)

type ClientPingTimer struct{}

func (*ClientPingTimer) New() *ClientPingTimer {
	return &ClientPingTimer{}
}

func (t *ClientPingTimer) GetName() string {
	return "ClientPing"
}

func (t *ClientPingTimer) GetInterval() time.Duration {
	return 1 * time.Second
}

func (t *ClientPingTimer) GetInitialDelay() time.Duration {
	return 1 * time.Second
}

func (t *ClientPingTimer) Handle(ctx actor.Context, mapData *entity.Map) error {
	_ = ctx
	if mapData.GetPlayerCount() == 0 {
		return nil
	}

	now := clock.Now()
	for _, obj := range mapData.GetAllPlayers() {
		ch, ok := obj.(*entity.Character)
		if !ok || ch == nil || ch.Sendable == nil {
			continue
		}

		client, ok := ch.Sendable.(*gameclient.GameClient)
		if !ok {
			continue
		}

		sendPing, disconnect := client.NextPingAction(now, 5*time.Second)
		if disconnect {
			log.Printf("disconnecting character %d (pong timeout)", ch.GetID())
			_ = client.GetConnection().Close()
			continue
		}
		if !sendPing {
			continue
		}

		if err := client.Send(&response.Ping{}, types.SEND_POLICY_ENCRYPT); err != nil {
			log.Printf("failed to send ping to character %d: %v", ch.GetID(), err)
			_ = client.GetConnection().Close()
		}
	}
	return nil
}
