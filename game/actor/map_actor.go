package actor

import (
	"log"
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/core"
	coreactor "github.com/boyism80/fm/core/actor"
	"github.com/boyism80/fm/game/entity"
)

type MapActor struct {
	MapData *entity.Map
	Context core.ServerContext
}

func (a *MapActor) Receive(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case *coreactor.HandlePacket:
		a.handlePacket(ctx, msg)
	case *coreactor.ScheduleTimer:
		a.scheduleTimer(ctx, msg)
	case *coreactor.ExecuteTimer:
		a.executeTimer(ctx, msg)
	case *AddCharacter:
		a.addCharacter(ctx, msg)
	case *RemoveCharacter:
		a.removeCharacter(ctx, msg)
	case *WarpCharacter:
		a.warpCharacter(ctx, msg)
	}
}

func (a *MapActor) handlePacket(ctx actor.Context, msg *coreactor.HandlePacket) {
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

func (a *MapActor) scheduleTimer(ctx actor.Context, msg *coreactor.ScheduleTimer) {
	if msg.Logic == nil {
		return
	}

	// Schedule timer using goroutine (protoactor-go doesn't have built-in ScheduleOnce)
	// Send ExecuteTimer message to self after the interval
	go func() {
		time.Sleep(msg.Interval)
		ctx.Send(ctx.Self(), &coreactor.ExecuteTimer{
			Logic: msg.Logic,
		})
	}()
}

func (a *MapActor) executeTimer(ctx actor.Context, msg *coreactor.ExecuteTimer) {
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
