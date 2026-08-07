package entity

import (
	"fmt"
	"sync"
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/core/async"
	"github.com/boyism80/fm/core/luax"
	lua "github.com/yuin/gopher-lua"
)

type StateMachine struct {
	mu              sync.Mutex
	ID              string
	Group           *StateMachineGroup
	Party           *Party
	Leader          *Character
	ScaleLevel      int
	MinPlayers      int
	ExitMapID       uint32
	players         []*Character
	playerSet       map[uint32]*Character
	props           map[string]string
	kills           map[uint32]int
	disposed        bool
	ActorPID        *actor.PID
	Maps            map[uint32]*Map
	timeoutDeadline time.Time
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

func (sm *StateMachine) LeavePlayer(ctx actor.Context, ch *Character, warpLeaver bool) bool {
	if sm == nil || ch == nil {
		return false
	}
	sm.mu.Lock()
	if sm.disposed {
		sm.mu.Unlock()
		return false
	}
	sm.unregisterLocked(ch)
	count := len(sm.players)
	minPlayers := sm.MinPlayers
	exitMapID := sm.ExitMapID
	group := sm.Group
	belowMin := minPlayers > 0 && count < minPlayers
	sm.mu.Unlock()

	if belowMin {
		sm.Finish(ctx, exitMapID)
		return true
	}
	if warpLeaver && exitMapID > 0 && group != nil && group.GameWorld != nil {
		if exitMap := group.GameWorld.GetMapSystem().Get(exitMapID); exitMap != nil {
			_ = ch.Warp(ctx, exitMap, 0)
		}
	}
	return false
}

func (sm *StateMachine) RequestLeave(ch *Character, warpLeaver bool) {
	if sm == nil || ch == nil || sm.Disposed() || sm.Group == nil || sm.Group.GameWorld == nil || sm.ActorPID == nil {
		return
	}
	sm.Group.GameWorld.SendStateMachineMessage(sm.ActorPID, &LeaveStateMachinePlayer{
		Character:  ch,
		WarpLeaver: warpLeaver,
	})
}

func (sm *StateMachine) EnterPlayer(ch *Character) {
	if sm == nil || ch == nil || sm.Disposed() || sm.Group == nil || sm.Group.GameWorld == nil || sm.ActorPID == nil {
		return
	}
	sm.Group.GameWorld.SendStateMachineMessage(sm.ActorPID, &EnterStateMachinePlayer{
		Character: ch,
	})
}

func (sm *StateMachine) Start() {
	if sm == nil || sm.Disposed() || sm.Group == nil || sm.Group.GameWorld == nil || sm.ActorPID == nil {
		return
	}
	sm.Group.GameWorld.SendStateMachineMessage(sm.ActorPID, &StartStateMachine{})
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

func (sm *StateMachine) StartTimerAsync(ctx actor.Context, ms int64) *async.Promise {
	if sm == nil || sm.Group == nil || sm.Group.GameWorld == nil || sm.ActorPID == nil || ms <= 0 || ctx == nil {
		return nil
	}
	gw := sm.Group.GameWorld
	pid := sm.ActorPID
	return async.Ask(ctx, pid, 10*time.Second, func(replyTo *actor.PID) {
		gw.SendStateMachineMessage(pid, &ScheduleStateMachineTimeout{
			Milliseconds: ms,
			ReplyTo:      replyTo,
		})
	})
}

func (sm *StateMachine) StopTimer() {
	if sm == nil || sm.Group == nil || sm.Group.GameWorld == nil || sm.ActorPID == nil {
		return
	}
	sm.Group.GameWorld.SendStateMachineMessage(sm.ActorPID, &CancelStateMachineTimeout{})
}

func (sm *StateMachine) StopTimerAsync(ctx actor.Context) *async.Promise {
	if sm == nil || sm.Group == nil || sm.Group.GameWorld == nil || sm.ActorPID == nil || ctx == nil {
		return nil
	}
	gw := sm.Group.GameWorld
	pid := sm.ActorPID
	return async.Ask(ctx, pid, 10*time.Second, func(replyTo *actor.PID) {
		gw.SendStateMachineMessage(pid, &CancelStateMachineTimeout{ReplyTo: replyTo})
	})
}

func (sm *StateMachine) TimeLeft() int64 {
	if sm == nil {
		return 0
	}
	sm.mu.Lock()
	defer sm.mu.Unlock()
	if sm.timeoutDeadline.IsZero() {
		return 0
	}
	left := time.Until(sm.timeoutDeadline).Milliseconds()
	if left < 0 {
		return 0
	}
	return left
}

func (sm *StateMachine) SetTimeoutDeadline(deadline time.Time) {
	if sm == nil {
		return
	}
	sm.mu.Lock()
	defer sm.mu.Unlock()
	sm.timeoutDeadline = deadline
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

func (sm *StateMachine) RegisterMap(m *Map) error {
	if sm == nil || m == nil {
		return fmt.Errorf("invalid map register")
	}
	if sm.Disposed() {
		return fmt.Errorf("state machine disposed")
	}
	if m.StateMachine() == sm {
		sm.RecordMap(m.GetMapID(), m)
		return nil
	}
	if err := m.AttachStateMachine(sm); err != nil {
		return err
	}
	if sm.ActorPID != nil {
		m.RebindObjectTimers(sm.ActorPID)
	}
	sm.RecordMap(m.GetMapID(), m)
	return nil
}

func (sm *StateMachine) OwnsMap(m *Map) bool {
	if sm == nil || m == nil {
		return false
	}
	if m.StateMachine() == sm {
		return true
	}
	sm.mu.Lock()
	defer sm.mu.Unlock()
	if sm.Maps == nil {
		return false
	}
	return sm.Maps[m.GetMapID()] == m
}

func (sm *StateMachine) HandlePlayerMapEnter(ch *Character, m *Map) {
	if sm == nil || ch == nil || m == nil {
		return
	}
	if sm.OwnsMap(m) {
		sm.CallHook("on_changed_map", sm, ch, m.TemplateID())
		return
	}
	sm.RequestLeave(ch, false)
}

func (sm *StateMachine) HasRegisteredMaps() bool {
	if sm == nil {
		return false
	}
	sm.mu.Lock()
	defer sm.mu.Unlock()
	return len(sm.Maps) > 0
}

func ResolveCreateMaps(result interface{}, ms MapSystem) ([]*Map, error) {
	if ms == nil {
		return nil, fmt.Errorf("map system is nil")
	}
	vals := luax.ResultValues(result)
	if len(vals) == 0 || vals[0] == nil || vals[0] == lua.LNil {
		return nil, fmt.Errorf("on_create must return a non-empty map array")
	}
	tbl, ok := vals[0].(*lua.LTable)
	if !ok {
		return nil, fmt.Errorf("on_create must return a table")
	}
	out := make([]*Map, 0)
	var resolveErr error
	tbl.ForEach(func(_ lua.LValue, value lua.LValue) {
		if resolveErr != nil {
			return
		}
		switch v := value.(type) {
		case lua.LNumber:
			id := uint32(v)
			m := ms.Get(id)
			if m == nil {
				m = ms.GetInstance(id)
			}
			if m == nil {
				resolveErr = fmt.Errorf("map %d not found", id)
				return
			}
			out = append(out, m)
		case *lua.LUserData:
			if v == nil || v.Value == nil {
				resolveErr = fmt.Errorf("nil map in on_create return")
				return
			}
			m, ok := v.Value.(*Map)
			if !ok || m == nil {
				resolveErr = fmt.Errorf("invalid map userdata in on_create return")
				return
			}
			out = append(out, m)
		default:
			if value == nil || value == lua.LNil {
				return
			}
			resolveErr = fmt.Errorf("on_create map entry must be map id or Map")
		}
	})
	if resolveErr != nil {
		return nil, resolveErr
	}
	if len(out) == 0 {
		return nil, fmt.Errorf("on_create must return a non-empty map array")
	}
	return out, nil
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
	scriptPath := ""
	if group != nil {
		scriptPath = group.ScriptPath
	}
	sm.players = nil
	sm.playerSet = make(map[uint32]*Character)
	sm.mu.Unlock()

	sm.callOnFinish(ctx, scriptPath)

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

func (sm *StateMachine) callOnFinish(ctx actor.Context, scriptPath string) {
	if sm == nil || scriptPath == "" {
		return
	}
	root := luax.NewState()
	defer root.Close()
	thread, err := luax.NewThread(root, scriptPath)
	if err != nil {
		return
	}
	if !luax.HasFunc(thread, "on_finish") {
		luax.Close(thread)
		return
	}
	luax.SetConfiguration(thread, luax.Configuration{
		ActorContext: ctx,
		ActorPID:     sm.ActorPID,
	})
	if _, err := luax.Call(thread, "on_finish", sm); err != nil {
		name := ""
		if sm.Group != nil {
			name = sm.Group.Name
		}
		fmt.Printf("state machine %s hook on_finish: %v\n", name, err)
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

type FinishStateMachineCreate struct{}

type ApplyStateMachineCreateMaps struct {
	Maps []*Map
	Err  string
}

type EnterStateMachinePlayer struct {
	Character *Character
}

type StartStateMachine struct{}

type LeaveStateMachinePlayer struct {
	Character  *Character
	WarpLeaver bool
}

type CallStateMachineHook struct {
	Hook string
	Args []interface{}
}

type ScheduleStateMachineTimeout struct {
	Milliseconds int64
	ReplyTo      *actor.PID
}

type CancelStateMachineTimeout struct {
	ReplyTo *actor.PID
}

type StateMachineTimerAck struct{}

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
	ch.RunQuestHook(questID, "on_quest_try")
}

func (ch *Character) EndPartyQuest(questID uint32) {
	if ch == nil || questID == 0 {
		return
	}
	ch.RunQuestHook(questID, "on_quest_end")
}
