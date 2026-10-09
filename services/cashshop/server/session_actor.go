package server

import (
	"log"
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/asynkron/protoactor-go/scheduler"
	"github.com/boyism80/fm/core"
	c_actor "github.com/boyism80/fm/core/actor"
	"github.com/boyism80/fm/core/clock"
	"github.com/boyism80/fm/protocol/response"
	"github.com/boyism80/fm/services/cashshop/client"
	"github.com/boyism80/fm/types"
)

type sessionActor struct {
	client     *client.CashShopClient
	server     *CashShopServer
	cancelPing scheduler.CancelFunc
}

type pingTick struct{}

func (a *sessionActor) Receive(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case *actor.Started:
		a.cancelPing = scheduler.NewTimerScheduler(ctx).SendRepeatedly(5*time.Second, 5*time.Second, ctx.Self(), &pingTick{})
	case *c_actor.HandlePacket:
		if a.client.Leaving() {
			return
		}
		err := core.ExecutePacketHandler(ctx, a.server, msg.Client, msg.Opcode, msg.Data, msg.LogicActorPID)
		if err != nil {
			log.Printf("Error handling packet 0x%02X: %v", msg.Opcode, err)
		}
	case *pingTick:
		a.ping()
	case *actor.Stopped:
		if a.cancelPing != nil {
			a.cancelPing()
		}
	}
}

func (a *sessionActor) ping() {
	sendPing, disconnect := a.client.NextPingAction(clock.Now(), 5*time.Second)
	if disconnect {
		_ = a.client.GetConnection().Close()
		return
	}
	if !sendPing {
		return
	}
	if err := a.client.Send(&response.Ping{}, types.SEND_POLICY_ENCRYPT); err != nil {
		_ = a.client.GetConnection().Close()
	}
}
