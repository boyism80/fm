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
	"github.com/boyism80/fm/game/constant"
	"github.com/boyism80/fm/game/entity"
	lua "github.com/yuin/gopher-lua"
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
		a.handlePacket(msg)
	case *c_actor.ScheduleTimer:
		a.scheduleTimer(ctx, msg)
	case *c_actor.ExecuteTimer:
		a.executeTimer(msg)
	case *AddCharacter:
		a.addCharacter(ctx, msg)
	case *RemoveCharacter:
		a.removeCharacter(msg)
	case *WarpCharacter:
		a.warpCharacter(ctx, msg)
	case *SpawnDoor:
		a.spawnTownMysticDoor(msg)
	case *RemoveDoor:
		a.removeMysticDoor(msg)
	case *ResumeLua:
		a.resumeLua(msg)
	case *c_actor.RunCharacterTimer:
		a.runCharacterTimer(ctx, msg)
	case *TimerTick:
		a.onTimerTick(ctx, msg)
	}
}

func (a *MapActor) resumeLua(msg *ResumeLua) {
	if msg.Root == nil || msg.Thread == nil {
		return
	}
	state, _, _ := msg.Root.Resume(msg.Thread, nil)
	if state == lua.ResumeOK {
		luax.ClearThreadPID(msg.Thread)
		msg.Thread.Close()
	}
}

func (a *MapActor) handlePacket(msg *c_actor.HandlePacket) {
	err := core.ExecutePacketHandler(a.Context, msg.Client, msg.Opcode, msg.Data, msg.LogicActorPID)
	if err != nil {
		log.Printf("Error handling packet 0x%02X: %v", msg.Opcode, err)
	}
}

func (a *MapActor) scheduleTimer(ctx actor.Context, msg *c_actor.ScheduleTimer) {
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

func (a *MapActor) executeTimer(msg *c_actor.ExecuteTimer) {
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
	a.MapData.AddPlayer(msg.Character.GetID(), msg.Character, msg.SpawnPoint, msg.Init)
	msg.Character.ResumeTimers(ctx.Self())
}

func (a *MapActor) removeCharacter(msg *RemoveCharacter) {
	if a.MapData == nil {
		return
	}
	a.MapData.RemovePlayer(msg.CharacterID)
}

func (a *MapActor) warpCharacter(ctx actor.Context, msg *WarpCharacter) {
	if a.MapData == nil {
		return
	}
	a.MapData.AddPlayer(msg.Character.GetID(), msg.Character, msg.Portal, false)
	msg.Character.ResumeTimers(ctx.Self())
}

func (a *MapActor) spawnTownMysticDoor(msg *SpawnDoor) {
	if a.MapData == nil || msg == nil || a.MapData.Wz == nil {
		return
	}
	wz := a.MapData.Wz
	position, ok := wz.GetSpawnPosition(msg.ReturnPortalID)
	if !ok {
		if p, ok2 := wz.Portals[msg.ReturnPortalID]; ok2 {
			position = p.Position
		} else {
			return
		}
	}
	door := entity.NewDoor(position, msg.OwnerID, msg.SkillID, wz.ID, msg.FieldMapID, msg.ReturnPortalID, msg.FieldPortalID)
	a.MapData.AddDoor(door)
}

func (a *MapActor) removeMysticDoor(msg *RemoveDoor) {
	if a.MapData == nil || msg == nil {
		return
	}
	a.MapData.RemoveMysticDoorByOwnerSkill(msg.OwnerID, constant.SkillID(msg.SkillID), true)
}

func (a *MapActor) runCharacterTimer(ctx actor.Context, msg *c_actor.RunCharacterTimer) {
	if a.MapData == nil {
		return
	}
	ch := a.MapData.GetPlayer(msg.CharacterID)
	if ch == nil {
		return
	}
	entry := ch.GetTimerEntry(msg.Key)
	if entry == nil {
		return
	}
	if entry.Callback != nil {
		entry.Callback()
	}
	if entry.Repeat && ch.GetTimerEntry(msg.Key) != nil {
		characterID := msg.CharacterID
		key := msg.Key
		entry.NextFireAt = time.Now().Add(entry.Interval)
		entry.Timer = time.AfterFunc(entry.Interval, func() {
			ctx.Send(ctx.Self(), &c_actor.RunCharacterTimer{CharacterID: characterID, Key: key})
		})
	} else if !entry.Repeat {
		ch.RemoveTimer(msg.Key)
	}
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
	RegisterTimer[*timers.MobSpawnTimer](a.timerReg)
	RegisterTimer[*timers.ItemCleanupTimer](a.timerReg)
	RegisterTimer[*timers.CooldownCheckTimer](a.timerReg)
	RegisterTimer[*timers.BuffExpireTimer](a.timerReg)
	RegisterTimer[*timers.MobBuffExpireTimer](a.timerReg)
	RegisterTimer[*timers.MobPoisonTickTimer](a.timerReg)
	RegisterTimer[*timers.MistExpireTimer](a.timerReg)
	RegisterTimer[*timers.MistPoisonTickTimer](a.timerReg)
}

func (a *MapActor) onTimerTick(ctx actor.Context, msg *TimerTick) {
	if a.MapData == nil {
		return
	}
	handler := a.timerReg.GetHandler(msg.HandlerName)
	if handler == nil {
		return
	}
	if err := handler.Handle(ctx, a.MapData); err != nil {
		log.Printf("Timer handler %s error: %v", handler.GetName(), err)
	}
}
