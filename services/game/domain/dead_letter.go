package domain

import (
	"log"

	"github.com/asynkron/protoactor-go/actor"
	c_actor "github.com/boyism80/fm/core/actor"
)

// Warps that reach a stopped actor would strand the character between maps and leak the instance ticket.
func (s *MapSystem) subscribeDeadLetters() {
	s.gw.Systems().Get[*c_actor.ActorRegistry]().System().EventStream.Subscribe(func(evt interface{}) {
		dead, ok := evt.(*actor.DeadLetterEvent)
		if !ok {
			return
		}
		switch msg := dead.Message.(type) {
		case *WarpCharacter:
			s.recoverWarp(dead.PID, msg)
		case *HandoffCharacter:
			s.recoverHandoff(dead.PID, msg)
		}
	})
}

func (s *MapSystem) recoverWarp(deadPID *actor.PID, msg *WarpCharacter) {
	log.Printf("invariant: warp of character %d to map %d reached stopped actor %v", msg.Character.GetID(), msg.TargetMap.GetMapID(), deadPID)

	msg.Ticket.Release()
	msg.Character.FinishMove(msg.TargetMap)

	exitMap := s.returnMap(msg.TargetMap)
	if exitMap == nil {
		log.Printf("invariant: character %d stranded; map %d has no return map", msg.Character.GetID(), msg.TargetMap.GetMapID())
		return
	}
	if !msg.Character.BeginMove(exitMap) {
		log.Printf("invariant: character %d stranded; it is already moving elsewhere", msg.Character.GetID())
		return
	}
	s.gw.Systems().Get[*c_actor.ActorRegistry]().Root().Send(exitMap.LogicActorPID(), &WarpCharacter{
		Character: msg.Character,
		TargetMap: exitMap,
	})
}

func (s *MapSystem) recoverHandoff(deadPID *actor.PID, msg *HandoffCharacter) {
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
	s.gw.Systems().Get[*c_actor.ActorRegistry]().Root().Send(pid, msg)
}
