package actor

import (
	"github.com/boyism80/fm/core/clock"
	"log"
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/asynkron/protoactor-go/scheduler"
	"github.com/boyism80/fm/core"
	c_actor "github.com/boyism80/fm/core/actor"
	"github.com/boyism80/fm/protocol/response"
	"github.com/boyism80/fm/services/login/client"
	"github.com/boyism80/fm/types"
)

type LoginLogicActor struct {
	Client     *client.LoginClient
	Server     core.Server
	timer      *scheduler.TimerScheduler
	cancelPing scheduler.CancelFunc
}

type pingTick struct{}

func (a *LoginLogicActor) Receive(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case *actor.Started:
		a.onStarted(ctx)
	case *c_actor.HandlePacket:
		a.handlePacket(ctx, msg)
	case *pingTick:
		a.onPingTick()
	case *actor.Stopped:
		a.onStopped(ctx)
	}
}

func (a *LoginLogicActor) onStarted(ctx actor.Context) {
	a.timer = scheduler.NewTimerScheduler(ctx)
	a.cancelPing = a.timer.SendRepeatedly(
		5*time.Second,
		5*time.Second,
		ctx.Self(),
		&pingTick{},
	)
}

func (a *LoginLogicActor) handlePacket(ctx actor.Context, msg *c_actor.HandlePacket) {
	if a.Client.Transferred() {
		log.Printf("Dropping packet 0x%02X from transferred client %d", msg.Opcode, a.Client.GetClientID())
		return
	}
	err := core.ExecutePacketHandler(ctx, a.Server, msg.Client, msg.Opcode, msg.Data, msg.LogicActorPID)
	if err != nil {
		log.Printf("Error handling packet 0x%02X: %v", msg.Opcode, err)
	}
}

func (a *LoginLogicActor) onStopped(ctx actor.Context) {
	if a.cancelPing != nil {
		a.cancelPing()
	}
	log.Printf("LoginLogicActor stopped for client %d", a.Client.GetClientID())
}

func (a *LoginLogicActor) onPingTick() {
	if a.Client.Transferred() {
		_ = a.Client.GetConnection().Close()
		return
	}

	sendPing, disconnect := a.Client.NextPingAction(clock.Now(), 5*time.Second)
	if disconnect {
		_ = a.Client.GetConnection().Close()
		return
	}
	if !sendPing {
		return
	}
	if err := a.Client.Send(&response.Ping{}, types.SEND_POLICY_ENCRYPT); err != nil {
		_ = a.Client.GetConnection().Close()
	}
}
