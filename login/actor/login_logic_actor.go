package actor

import (
	"log"
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/core"
	c_actor "github.com/boyism80/fm/core/actor"
	"github.com/boyism80/fm/login/client"
)

type LoginLogicActor struct {
	Client  core.Client
	Context core.ServerContext
}

func (a *LoginLogicActor) Receive(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case *c_actor.HandlePacket:
		a.handlePacket(ctx, msg)
	case *c_actor.ScheduleTimer:
		a.scheduleTimer(ctx, msg)
	case *c_actor.ExecuteTimer:
		a.executeTimer(ctx, msg)
	case *actor.Stopped:
		a.onStopped(ctx)
	}
}

func (a *LoginLogicActor) handlePacket(ctx actor.Context, msg *c_actor.HandlePacket) {
	client, ok := msg.Client.(core.Client)
	if !ok {
		log.Printf("Invalid client type in HandlePacket")
		return
	}
	err := core.ExecutePacketHandler(a.Context, client, msg.Opcode, msg.Data)
	if err != nil {
		log.Printf("Error handling packet 0x%02X: %v", msg.Opcode, err)
	}
}

func (a *LoginLogicActor) scheduleTimer(ctx actor.Context, msg *c_actor.ScheduleTimer) {
	if msg.Logic == nil {
		return
	}

	// Schedule timer using goroutine (protoactor-go doesn't have built-in ScheduleOnce)
	// Send ExecuteTimer message to self after the interval
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
