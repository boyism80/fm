package server

import (
	"log"

	"github.com/asynkron/protoactor-go/actor"
	g_actor "github.com/boyism80/fm/services/game/actor"
)

// Warps that reach a stopped actor would strand the character between maps and leak the instance ticket.
func (gs *GameServer) subscribeDeadLetters() {
	gs.actorSystem.GetSystem().EventStream.Subscribe(func(evt interface{}) {
		dead, ok := evt.(*actor.DeadLetterEvent)
		if !ok {
			return
		}
		switch msg := dead.Message.(type) {
		case *g_actor.WarpCharacter:
			gs.recoverWarp(dead.PID, msg)
		case *g_actor.HandoffCharacter:
			gs.recoverHandoff(dead.PID, msg)
		}
	})
}

func (gs *GameServer) recoverWarp(deadPID *actor.PID, msg *g_actor.WarpCharacter) {
	log.Printf("invariant: warp of character %d to map %d reached stopped actor %v", msg.Character.GetID(), msg.TargetMap.GetMapID(), deadPID)

	msg.Ticket.Release()
	msg.Character.FinishMove(msg.TargetMap)

	exitMap := mapSystem{gs}.returnMap(msg.TargetMap)
	if exitMap == nil {
		log.Printf("invariant: character %d stranded; map %d has no return map", msg.Character.GetID(), msg.TargetMap.GetMapID())
		return
	}
	if !msg.Character.BeginMove(exitMap) {
		log.Printf("invariant: character %d stranded; it is already moving elsewhere", msg.Character.GetID())
		return
	}
	gs.GetRootContext().Send(exitMap.LogicActorPID(), &g_actor.WarpCharacter{
		Character: msg.Character,
		TargetMap: exitMap,
	})
}

func (gs *GameServer) recoverHandoff(deadPID *actor.PID, msg *g_actor.HandoffCharacter) {
	log.Printf("invariant: handoff of character %d to map %d reached stopped actor %v", msg.Character.GetID(), msg.TargetMap.GetMapID(), deadPID)

	source := msg.Character.GetMap()
	if source == nil {
		msg.Ticket.Release()
		return
	}
	pid := source.LogicActorPID()
	if pid == nil || pid.Equal(deadPID) {
		msg.Ticket.Release()
		return
	}
	gs.GetRootContext().Send(pid, msg)
}
