package entity

import (
	"fmt"
	"sync"
	"time"

	"github.com/boyism80/fm/core/async"
	"github.com/boyism80/fm/core/luax"
	lua "github.com/yuin/gopher-lua"
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
	timer      *time.Timer
	disposed   bool
}

func NewStateMachine(id string, group *StateMachineGroup) *StateMachine {
	return &StateMachine{
		ID:        id,
		Group:     group,
		players:   make([]*Character, 0),
		playerSet: make(map[uint32]*Character),
		props:     make(map[string]string),
		kills:     make(map[uint32]int),
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
	if sm == nil || ms <= 0 {
		return
	}
	sm.mu.Lock()
	defer sm.mu.Unlock()
	if sm.disposed {
		return
	}
	if sm.timer != nil {
		sm.timer.Stop()
		sm.timer = nil
	}
	sm.timer = time.AfterFunc(time.Duration(ms)*time.Millisecond, func() {
		sm.CallHook("on_scheduled_timeout", sm)
	})
}

func (sm *StateMachine) stopTimerLocked() {
	if sm.timer != nil {
		sm.timer.Stop()
		sm.timer = nil
	}
}

func (sm *StateMachine) Finish(exitMapID uint32) {
	if sm == nil {
		return
	}
	sm.mu.Lock()
	if sm.disposed {
		sm.mu.Unlock()
		return
	}
	sm.disposed = true
	sm.stopTimerLocked()
	players := make([]*Character, len(sm.players))
	copy(players, sm.players)
	group := sm.Group
	id := sm.ID
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
			_ = ch.Warp(exitMap, 0)
		}
	}
	if group != nil {
		group.RemoveMachine(id)
	}
}

func (sm *StateMachine) luaMap() *Map {
	if sm == nil {
		return nil
	}
	sm.mu.Lock()
	defer sm.mu.Unlock()
	for _, ch := range sm.players {
		if ch != nil && ch.GetMap() != nil {
			return ch.GetMap()
		}
	}
	if sm.Leader != nil && sm.Leader.GetMap() != nil {
		return sm.Leader.GetMap()
	}
	if sm.Group != nil && sm.Group.GameWorld != nil {
		return sm.Group.bootMap()
	}
	return nil
}

func stateMachineHookDone() *async.Promise {
	p := async.NewDeferred()
	p.SetResult(nil)
	return p
}

func (sm *StateMachine) CallHook(hook string, args ...interface{}) *async.Promise {
	if sm == nil || sm.Group == nil || hook == "" {
		return stateMachineHookDone()
	}
	if sm.Disposed() && hook != "on_scheduled_timeout" {
		return stateMachineHookDone()
	}
	m := sm.luaMap()
	if m == nil {
		return stateMachineHookDone()
	}
	root := m.EnsureLuaRoot(nil)
	if root == nil {
		return stateMachineHookDone()
	}
	path := sm.Group.ScriptPath
	thread, err := luax.NewThread(root, path)
	if err != nil {
		return stateMachineHookDone()
	}
	fn := thread.GetGlobal(hook)
	if fn.Type() != lua.LTFunction {
		luax.Close(thread)
		return stateMachineHookDone()
	}
	luax.SetConfiguration(thread, luax.Configuration{
		MapActorPID: m.GetActorPID(),
	})
	return luax.CallAsync(root, thread, hook, args...).OnError(func(err error) {
		fmt.Printf("state machine %s hook %s: %v\n", sm.Group.Name, hook, err)
	})
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
