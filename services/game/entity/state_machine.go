package entity

import (
	"fmt"
	"sync"

	"github.com/asynkron/protoactor-go/actor"
)

type StateMachine struct {
	mu         sync.Mutex
	ID         string
	Group      *StateMachineGroup
	Party      *Party
	Leader     *Character
	ScaleLevel int
	players    []*Character
	playerSet  map[uint32]*Character
	props      map[string]string
	kills      map[uint32]int
	disposed   bool
	ActorPID   *actor.PID
	Maps       map[uint32]*Map
}

func NewStateMachine(id string, group *StateMachineGroup) *StateMachine {
	return &StateMachine{
		ID:        id,
		Group:     group,
		players:   make([]*Character, 0),
		playerSet: make(map[uint32]*Character),
		props:     make(map[string]string),
		kills:     make(map[uint32]int),
		Maps:      make(map[uint32]*Map),
	}
}

func (sm *StateMachine) Disposed() bool {
	if sm == nil {
		return true
	}
	sm.mu.Lock()
	defer sm.mu.Unlock()
	return sm.disposed
}

func (sm *StateMachine) Register(ch *Character) {
	if sm == nil || ch == nil {
		return
	}
	sm.mu.Lock()
	defer sm.mu.Unlock()
	if sm.disposed {
		return
	}
	id := ch.GetID()
	if _, ok := sm.playerSet[id]; ok {
		return
	}
	sm.players = append(sm.players, ch)
	sm.playerSet[id] = ch
	ch.stateMachine = sm
}

func (sm *StateMachine) Unregister(ch *Character) {
	if sm == nil || ch == nil {
		return
	}
	sm.mu.Lock()
	defer sm.mu.Unlock()
	sm.unregisterLocked(ch)
}

func (sm *StateMachine) unregisterLocked(ch *Character) {
	if ch == nil {
		return
	}
	id := ch.GetID()
	if _, ok := sm.playerSet[id]; !ok {
		return
	}
	delete(sm.playerSet, id)
	out := sm.players[:0]
	for _, p := range sm.players {
		if p != nil && p.GetID() != id {
			out = append(out, p)
		}
	}
	sm.players = out
	if ch.stateMachine == sm {
		ch.stateMachine = nil
	}
}

func (sm *StateMachine) Players() []*Character {
	if sm == nil {
		return nil
	}
	sm.mu.Lock()
	defer sm.mu.Unlock()
	out := make([]*Character, len(sm.players))
	copy(out, sm.players)
	return out
}

func (sm *StateMachine) SetProperty(key, value string) {
	if sm == nil {
		return
	}
	sm.mu.Lock()
	defer sm.mu.Unlock()
	if sm.props == nil {
		sm.props = make(map[string]string)
	}
	sm.props[key] = value
}

func (sm *StateMachine) GetProperty(key string) string {
	if sm == nil {
		return ""
	}
	sm.mu.Lock()
	defer sm.mu.Unlock()
	if sm.props == nil {
		return ""
	}
	return sm.props[key]
}

func (sm *StateMachine) AddKill(ch *Character, n int) {
	if sm == nil || ch == nil || n == 0 {
		return
	}
	sm.mu.Lock()
	defer sm.mu.Unlock()
	if sm.disposed {
		return
	}
	if sm.kills == nil {
		sm.kills = make(map[uint32]int)
	}
	sm.kills[ch.GetID()] += n
}

func (sm *StateMachine) StartTimer(ms int64) {
	if sm == nil || sm.Group == nil || sm.Group.GameWorld == nil || sm.ActorPID == nil || ms <= 0 {
		return
	}
	sm.Group.GameWorld.SendStateMachineMessage(sm.ActorPID, &ScheduleStateMachineTimeout{Milliseconds: ms})
}

func (sm *StateMachine) RecordMap(mapID uint32, m *Map) {
	if sm == nil || m == nil {
		return
	}
	sm.mu.Lock()
	defer sm.mu.Unlock()
	if sm.Maps == nil {
		sm.Maps = make(map[uint32]*Map)
	}
	sm.Maps[mapID] = m
}

func (sm *StateMachine) MapList() []*Map {
	if sm == nil {
		return nil
	}
	sm.mu.Lock()
	defer sm.mu.Unlock()
	maps := make([]*Map, 0, len(sm.Maps))
	for _, m := range sm.Maps {
		maps = append(maps, m)
	}
	return maps
}

func (sm *StateMachine) ClearMaps() {
	if sm == nil {
		return
	}
	sm.mu.Lock()
	defer sm.mu.Unlock()
	sm.Maps = make(map[uint32]*Map)
}

func (sm *StateMachine) AbortStart() {
	if sm == nil {
		return
	}
	sm.mu.Lock()
	if sm.disposed {
		sm.mu.Unlock()
		return
	}
	sm.disposed = true
	players := make([]*Character, len(sm.players))
	copy(players, sm.players)
	sm.players = nil
	sm.playerSet = make(map[uint32]*Character)
	sm.mu.Unlock()
	for _, ch := range players {
		if ch != nil && ch.stateMachine == sm {
			ch.stateMachine = nil
		}
	}
}

func (sm *StateMachine) Finish(ctx actor.Context, exitMapID uint32) {
	if sm == nil {
		return
	}
	sm.mu.Lock()
	if sm.disposed {
		sm.mu.Unlock()
		return
	}
	sm.disposed = true
	players := make([]*Character, len(sm.players))
	copy(players, sm.players)
	group := sm.Group
	pid := sm.ActorPID
	sm.players = nil
	sm.playerSet = make(map[uint32]*Character)
	sm.mu.Unlock()

	var exitMap *Map
	if exitMapID > 0 && group != nil && group.GameWorld != nil {
		exitMap = group.GameWorld.GetMapSystem().Get(exitMapID)
	}
	for _, ch := range players {
		if ch == nil {
			continue
		}
		if ch.stateMachine == sm {
			ch.stateMachine = nil
		}
		if exitMap != nil {
			_ = ch.Warp(ctx, exitMap, 0)
		}
	}
	if group != nil && group.GameWorld != nil && pid != nil {
		group.GameWorld.SendStateMachineMessage(pid, &StopStateMachine{})
	}
}

func (sm *StateMachine) CallHook(hook string, args ...interface{}) {
	if sm == nil || sm.Disposed() || sm.Group == nil || sm.Group.GameWorld == nil || sm.ActorPID == nil || hook == "" {
		return
	}
	sm.Group.GameWorld.SendStateMachineMessage(sm.ActorPID, &CallStateMachineHook{
		Hook: hook,
		Args: args,
	})
}

type BootstrapStateMachine struct{}

type StopStateMachine struct{}

type EnterStateMachinePlayers struct{}

type CallStateMachineHook struct {
	Hook string
	Args []interface{}
}

type ScheduleStateMachineTimeout struct {
	Milliseconds int64
}

type StateMachineTimeout struct {
	Version uint64
}

func (ch *Character) StateMachine() *StateMachine {
	if ch == nil {
		return nil
	}
	return ch.stateMachine
}

func (ch *Character) TryPartyQuest(questID uint32) {
	if ch == nil || questID == 0 {
		return
	}
	ch.RunQuestHook(questID, fmt.Sprintf("on_quest_try_%d", questID))
}

func (ch *Character) EndPartyQuest(questID uint32) {
	if ch == nil || questID == 0 {
		return
	}
	ch.RunQuestHook(questID, fmt.Sprintf("on_quest_end_%d", questID))
}
