package game

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/asynkron/protoactor-go/actor"
	c_actor "github.com/boyism80/fm/core/actor"
	lua "github.com/yuin/gopher-lua"
)

type StateMachineSystem struct {
	gw            GameWorld
	mu            sync.RWMutex
	groups        map[string]*StateMachineGroup
	actors        map[string]*StateMachine
	nextMachineID atomic.Uint32
}

func NewStateMachineSystem(gw GameWorld) *StateMachineSystem {
	return &StateMachineSystem{
		gw:     gw,
		groups: make(map[string]*StateMachineGroup),
		actors: make(map[string]*StateMachine),
	}
}

func (s *StateMachineSystem) Start() error { return nil }
func (s *StateMachineSystem) Stop()        {}

func (s *StateMachineSystem) Get(name string) *StateMachineGroup {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.groups[name]
}

func (s *StateMachineSystem) Groups() []*StateMachineGroup {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]*StateMachineGroup, 0, len(s.groups))
	for _, group := range s.groups {
		out = append(out, group)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].name < out[j].name })
	return out
}

func (s *StateMachineSystem) FindByActor(pid *actor.PID) *StateMachine {
	if pid == nil {
		return nil
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.actors[pid.Id]
}

func (s *StateMachineSystem) registerActor(pid *actor.PID, sm *StateMachine) {
	s.mu.Lock()
	s.actors[pid.Id] = sm
	s.mu.Unlock()
}

func (s *StateMachineSystem) unregisterActor(pid *actor.PID) {
	s.mu.Lock()
	delete(s.actors, pid.Id)
	s.mu.Unlock()
}

func (s *StateMachineSystem) Register(group *StateMachineGroup) {
	if group == nil || group.name == "" {
		return
	}
	s.mu.Lock()
	s.groups[group.name] = group
	s.mu.Unlock()
}

func (s *StateMachineSystem) LoadFromScripts(dir string) error {
	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			log.Printf("state machine: script dir %s not found, skipping", dir)
			return nil
		}
		return err
	}
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		name := e.Name()
		if !strings.HasSuffix(name, ".lua") {
			continue
		}
		stem := strings.TrimSuffix(name, ".lua")
		scriptPath := filepath.ToSlash(filepath.Join(dir, name))
		group := NewStateMachineGroup(stem, scriptPath, s.gw)
		s.Register(group)
		group.Init()
		log.Printf("state machine: loaded group %s (%s)", stem, scriptPath)
	}
	return nil
}

func (s *StateMachineSystem) StartActor(sm *StateMachine) *actor.PID {
	if sm == nil || sm.group == nil {
		return nil
	}
	name := fmt.Sprintf("state_machine_%s_%s_%d", sm.group.name, sm.id, s.nextMachineID.Add(1))
	smActor := s.gw.Systems().Get[ActorFactory]().NewStateMachineActor(sm)
	props := actor.PropsFromProducer(func() actor.Actor {
		return smActor
	})
	pid := s.gw.Systems().Get[*c_actor.ActorRegistry]().GetOrCreateActor(name, props)
	if pid == nil {
		return nil
	}
	s.registerActor(pid, sm)
	return pid
}

func (s *StateMachineSystem) StopActor(sm *StateMachine) {
	if sm == nil || sm.actorPID == nil {
		return
	}
	s.unregisterActor(sm.actorPID)
	s.gw.Systems().Get[*c_actor.ActorRegistry]().PoisonActor(sm.actorPID.Id, sm.actorPID)
}

func (s *StateMachineSystem) Send(pid *actor.PID, msg interface{}) {
	if pid == nil || msg == nil {
		return
	}
	s.gw.Systems().Get[*c_actor.ActorRegistry]().Root().Send(pid, msg)
}

func (s *StateMachineSystem) ResumeLua(pid *actor.PID, root *lua.LState, thread *lua.LState, args []lua.LValue) {
	if pid == nil || root == nil || thread == nil {
		return
	}
	s.gw.Systems().Get[*c_actor.ActorRegistry]().Root().Send(pid, &ResumeLua{
		Root:   root,
		Thread: thread,
		Args:   args,
	})
}
