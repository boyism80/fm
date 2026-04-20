package actor

import (
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
	Client core.Client
	Server core.Server
	timer  *scheduler.TimerScheduler
}

type pingTick struct{}

func (a *LoginLogicActor) Receive(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case *actor.Started:
		a.onStarted(ctx)
	case *c_actor.HandlePacket:
		a.handlePacket(ctx, msg)
	case *c_actor.ScheduleTimer:
		a.scheduleTimer(ctx, msg)
	case *c_actor.ExecuteTimer:
		a.executeTimer(ctx, msg)
	case *pingTick:
		a.onPingTick()
	case *actor.Stopped:
		a.onStopped(ctx)
	}
}

func (a *LoginLogicActor) onStarted(ctx actor.Context) {
	a.timer = scheduler.NewTimerScheduler(ctx)
	a.timer.SendRepeatedly(
		5*time.Second,
		5*time.Second,
		ctx.Self(),
		&pingTick{},
	)
}

func (a *LoginLogicActor) handlePacket(ctx actor.Context, msg *c_actor.HandlePacket) {
	err := core.ExecutePacketHandler(ctx, a.Server, msg.Client, msg.Opcode, msg.Data, msg.LogicActorPID)
	if err != nil {
		log.Printf("Error handling packet 0x%02X: %v", msg.Opcode, err)
	}
}

func (a *LoginLogicActor) scheduleTimer(ctx actor.Context, msg *c_actor.ScheduleTimer) {
	if msg.Logic == nil {
		return
	}

	go func() {
		time.Sleep(msg.Interval)
		ctx.Send(ctx.Self(), &c_actor.ExecuteTimer{
			Logic: msg.Logic,
		})
	}()
}

func (a *LoginLogicActor) executeTimer(ctx actor.Context, msg *c_actor.ExecuteTimer) {
	if msg.Logic == nil {
		return
	}

	if err := msg.Logic(); err != nil {
		log.Printf("Error executing timer logic: %v", err)
	}
}

func (a *LoginLogicActor) onStopped(ctx actor.Context) {
	if loginClient, ok := a.Client.(*client.LoginClient); ok {
		log.Printf("LoginLogicActor stopped for client %d", loginClient.GetFd())
	} else {
		log.Printf("LoginLogicActor stopped")
	}
}

func (a *LoginLogicActor) onPingTick() {
	loginClient, ok := a.Client.(*client.LoginClient)
	if !ok {
		return
	}

	sendPing, disconnect := loginClient.NextPingAction(time.Now(), 5*time.Second)
	if disconnect {
		_ = loginClient.GetConnection().Close()
		return
	}
	if !sendPing {
		return
	}
	if err := loginClient.Send(&response.Ping{}, types.SEND_POLICY_ENCRYPT); err != nil {
		_ = loginClient.GetConnection().Close()
	}
}
