package actor

import (
	"log"
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/asynkron/protoactor-go/scheduler"
	"github.com/boyism80/fm/core"
	c_actor "github.com/boyism80/fm/core/actor"
	"github.com/boyism80/fm/core/luax"
	"github.com/boyism80/fm/game/actor/timers"
	"github.com/boyism80/fm/game/entity"
)

type MapActor struct {
	MapData   *entity.Map
	Context   core.ServerContext
	scheduler *scheduler.TimerScheduler
	timerReg  *TimerRegistry
}

func (a *MapActor) Receive(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case *actor.Started:
		a.onStarted(ctx)
	case *actor.Stopped:
		a.onStopped(ctx)
	case *c_actor.HandlePacket:
		a.handlePacket(ctx, msg)
	case *c_actor.ScheduleTimer:
		a.scheduleTimer(ctx, msg)
	case *c_actor.ExecuteTimer:
		a.executeTimer(ctx, msg)
	case *AddCharacter:
		a.addCharacter(ctx, msg)
	case *RemoveCharacter:
		a.removeCharacter(ctx, msg)
	case *WarpCharacter:
		a.warpCharacter(ctx, msg)
	case *TimerTick:
		a.onTimerTick(ctx, msg)
	}
}

func (a *MapActor) handlePacket(ctx actor.Context, msg *c_actor.HandlePacket) {
	client, ok := msg.Client.(core.Client)
	if !ok {
		log.Printf("Invalid client type in HandlePacket")
		return
	}
	err := core.ExecutePacketHandler(a.Context, client, msg.Opcode, msg.Data, ctx.Self().String())
	if err != nil {
		log.Printf("Error handling packet 0x%02X: %v", msg.Opcode, err)
	}
}

func (a *MapActor) scheduleTimer(ctx actor.Context, msg *c_actor.ScheduleTimer) {
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

func (a *MapActor) executeTimer(ctx actor.Context, msg *c_actor.ExecuteTimer) {
	if msg.Logic == nil {
		return
	}

	if err := msg.Logic(); err != nil {
		log.Printf("Error executing timer logic: %v", err)
	}
}

func (a *MapActor) addCharacter(ctx actor.Context, msg *AddCharacter) {
	if a.MapData == nil {
		return
	}
	a.MapData.AddPlayer(msg.Character.ID, msg.Character, msg.SpawnPoint, msg.Init)
}

func (a *MapActor) removeCharacter(ctx actor.Context, msg *RemoveCharacter) {
	if a.MapData == nil {
		return
	}
	a.MapData.RemovePlayer(msg.CharacterID)
}

func (a *MapActor) warpCharacter(ctx actor.Context, msg *WarpCharacter) {
	if a.MapData == nil {
		return
	}
	a.MapData.AddPlayer(msg.Character.ID, msg.Character, msg.Portal, false)
}

func (a *MapActor) onStarted(ctx actor.Context) {
	L := luax.NewState()
	luax.RegisterRootLuaState(ctx.Self().String(), L)

	a.scheduler = scheduler.NewTimerScheduler(ctx)
	a.timerReg = NewTimerRegistry()
	a.registerTimers()

	for _, handler := range a.timerReg.GetAllHandlers() {
		a.scheduler.SendRepeatedly(
			handler.GetInitialDelay(),
			handler.GetInterval(),
			ctx.Self(),
			&TimerTick{
				HandlerName: handler.GetName(),
			},
		)
	}
}

func (a *MapActor) onStopped(ctx actor.Context) {
	luax.UnregisterRootLuaState(ctx.Self().String())
}

func (a *MapActor) registerTimers() {
	RegisterTimer[*timers.MobSpawnTimer, *timers.MobSpawnTimer](a.timerReg)
	RegisterTimer[*timers.ItemCleanupTimer, *timers.ItemCleanupTimer](a.timerReg)
}

func (a *MapActor) onTimerTick(ctx actor.Context, msg *TimerTick) {
	if a.MapData == nil {
		return
	}

	for _, handler := range a.timerReg.GetAllHandlers() {
		if handler.GetName() == msg.HandlerName {
			if a.MapData.GetPlayerCount() > 0 {
				if err := handler.Handle(ctx, a.MapData); err != nil {
					log.Printf("Timer handler %s error: %v", handler.GetName(), err)
				}
			}
			return
		}
	}
}
