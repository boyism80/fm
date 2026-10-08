package entity

import (
	"errors"
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
	mapRefs         map[*Map]*MapRef
	timeoutDeadline time.Time
}

var ErrStateMachineEntryDenied = errors.New("state machine map entry denied")

type StateMachineMapSpec struct {
	TemplateID uint32
	Opts       MapInitOpts
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
		mapRefs:   make(map[*Map]*MapRef),
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
	ch.bindStateMachine(sm)
}

func (sm *StateMachine) Unregister(ch *Character) {
	if sm == nil || ch == nil {
		return
	}
	sm.mu.Lock()
	defer sm.mu.Unlock()
	sm.unregisterLocked(ch)
}

type StateMachineLeaveReason string

const (
	StateMachineLeaveDisconnect StateMachineLeaveReason = "disconnect"
	StateMachineLeaveMap        StateMachineLeaveReason = "map"
	StateMachineLeaveParty      StateMachineLeaveReason = "party"
	StateMachineLeaveExit       StateMachineLeaveReason = "exit"
)

func (sm *StateMachine) LeavePlayer(ctx actor.Context, ch *Character, warpLeaver bool, reason StateMachineLeaveReason) bool {
	if sm == nil || ch == nil {
		return false
	}
	sm.mu.Lock()
	if sm.disposed || sm.unregisterLocked(ch) == false {
		sm.mu.Unlock()
		return false
	}
	count := len(sm.players)
	minPlayers := sm.MinPlayers
	exitMapID := sm.ExitMapID
	group := sm.Group
	belowMin := minPlayers > 0 && count < minPlayers
	sm.mu.Unlock()

	if belowMin {
		if group != nil {
			sm.callHookNow(ctx, group.ScriptPath, "on_player_leave", ch, string(reason))
		}
		sm.Finish(ctx, exitMapID, 0)
		return true
	}
	if warpLeaver && exitMapID > 0 && group != nil && group.GameWorld != nil {
		if exitMap := group.GameWorld.GetMapSystem().Get(exitMapID); exitMap != nil {
			_ = ch.Warp(ctx, exitMap, 0)
		}
	}
	sm.CallHook("on_player_leave", ch, string(reason))
	return false
}

func (sm *StateMachine) RequestLeave(ch *Character, warpLeaver bool, reason StateMachineLeaveReason) {
	if sm == nil || ch == nil || sm.Disposed() || sm.Group == nil || sm.Group.GameWorld == nil || sm.ActorPID == nil {
		return
	}
	sm.Group.GameWorld.SendStateMachineMessage(sm.ActorPID, &LeaveStateMachinePlayer{
		Character:  ch,
		WarpLeaver: warpLeaver,
		Reason:     reason,
	})
}

func (sm *StateMachine) RequestFinish(exitMapID uint32, exitPortal uint8) {
	if sm == nil || sm.Disposed() || sm.Group == nil || sm.Group.GameWorld == nil || sm.ActorPID == nil {
		return
	}
	sm.Group.GameWorld.SendStateMachineMessage(sm.ActorPID, &FinishStateMachine{
		ExitMapID:  exitMapID,
		ExitPortal: exitPortal,
	})
}

func (sm *StateMachine) EnterPlayer(ch *Character, args ...interface{}) {
	if sm == nil || ch == nil || sm.Disposed() || sm.Group == nil || sm.Group.GameWorld == nil || sm.ActorPID == nil {
		return
	}
	sm.Group.GameWorld.SendStateMachineMessage(sm.ActorPID, &EnterStateMachinePlayer{
		Character: ch,
		Args:      args,
	})
}

func (sm *StateMachine) Start() {
	if sm == nil || sm.Disposed() || sm.Group == nil || sm.Group.GameWorld == nil || sm.ActorPID == nil {
		return
	}
	sm.Group.GameWorld.SendStateMachineMessage(sm.ActorPID, &StartStateMachine{})
}

func (sm *StateMachine) unregisterLocked(ch *Character) bool {
	if ch == nil {
		return false
	}
	id := ch.GetID()
	if _, ok := sm.playerSet[id]; !ok {
		return false
	}
	delete(sm.playerSet, id)
	out := sm.players[:0]
	for _, p := range sm.players {
		if p != nil && p.GetID() != id {
			out = append(out, p)
		}
	}
	sm.players = out
	ch.unbindStateMachine(sm)
	return true
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

func (sm *StateMachine) LimitTimer(ms int64) int64 {
	for _, ch := range append(sm.Players(), sm.Leader) {
		if ch == nil || ch.GM.TimerLimit == 0 {
			continue
		}
		if limit := int64(ch.GM.TimerLimit) * 1000; limit < ms {
			ms = limit
		}
	}
	return ms
}

func (sm *StateMachine) StartTimer(ms int64) {
	if sm == nil || sm.Group == nil || sm.Group.GameWorld == nil || sm.ActorPID == nil || ms <= 0 {
		return
	}
	sm.Group.GameWorld.SendStateMachineMessage(sm.ActorPID, &ScheduleStateMachineTimeout{Milliseconds: ms})
}

func (sm *StateMachine) StartTimerAsync(ctx actor.Context, ms int64) *async.Promise[interface{}] {
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

func (sm *StateMachine) StopTimerAsync(ctx actor.Context) *async.Promise[interface{}] {
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

func (sm *StateMachine) BroadcastClock(seconds int32) {
	if sm == nil {
		return
	}
	for _, ch := range sm.Players() {
		if ch == nil {
			continue
		}
		if !sm.OwnsMap(ch.GetMap()) {
			continue
		}
		ch.Listener.OnClock(ch, seconds)
	}
}

func (sm *StateMachine) SyncClock(ch *Character) {
	if sm == nil || ch == nil {
		return
	}
	if !sm.OwnsMap(ch.GetMap()) {
		return
	}
	left := sm.TimeLeft()
	if left <= 0 {
		return
	}
	ch.Listener.OnClock(ch, int32(left/1000))
}

func (sm *StateMachine) SetTimeoutDeadline(deadline time.Time) {
	if sm == nil {
		return
	}
	sm.mu.Lock()
	defer sm.mu.Unlock()
	sm.timeoutDeadline = deadline
}

func (sm *StateMachine) CreateMaps(specs []StateMachineMapSpec) error {
	ms := sm.Group.GameWorld.GetMapSystem()
	for _, spec := range specs {
		if sm.Map(spec.TemplateID) != nil {
			continue
		}
		m, err := ms.CreateStateMachineMap(spec.TemplateID, sm, spec.Opts)
		if err != nil {
			return err
		}
		ref, err := m.Reserve()
		if err != nil {
			return err
		}
		sm.mu.Lock()
		sm.Maps[spec.TemplateID] = m
		sm.mapRefs[m] = ref
		sm.mu.Unlock()
	}
	return nil
}

func (sm *StateMachine) Map(templateID uint32) *Map {
	if sm == nil {
		return nil
	}
	sm.mu.Lock()
	defer sm.mu.Unlock()
	return sm.Maps[templateID]
}

func (sm *StateMachine) CloseMaps() {
	ms := sm.Group.GameWorld.GetMapSystem()
	sm.mu.Lock()
	refs := sm.mapRefs
	sm.mapRefs = make(map[*Map]*MapRef)
	sm.mu.Unlock()

	for m, ref := range refs {
		ms.CloseInstance(m)
		ref.Release()
	}
}

func (sm *StateMachine) RemoveMap(m *Map) (empty bool) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	if sm.Maps[m.TemplateID()] == m {
		delete(sm.Maps, m.TemplateID())
	}
	return len(sm.Maps) == 0
}

func (sm *StateMachine) OwnsMap(m *Map) bool {
	if sm == nil || m == nil {
		return false
	}
	return m.StateMachine() == sm
}

func (sm *StateMachine) HasPlayer(ch *Character) bool {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	_, ok := sm.playerSet[ch.GetID()]
	return ok
}

func (sm *StateMachine) Admits(ch *Character) bool {
	if sm.HasPlayer(ch) {
		return true
	}
	return ch.ActsAsGM()
}

func (sm *StateMachine) HandlePlayerMapEnter(ch *Character, m *Map) {
	if sm == nil || ch == nil || m == nil {
		return
	}
	if sm.OwnsMap(m) {
		sm.SyncClock(ch)
		sm.CallHook("on_changed_map", ch, m.TemplateID())
		return
	}
	sm.RequestLeave(ch, false, StateMachineLeaveMap)
}

func ParseCreateMaps(tbl *lua.LTable) ([]StateMachineMapSpec, error) {
	if tbl == nil {
		return nil, fmt.Errorf("on_create must return a table")
	}
	out := make([]StateMachineMapSpec, 0)
	var parseErr error
	tbl.ForEach(func(_ lua.LValue, value lua.LValue) {
		if parseErr != nil {
			return
		}
		switch v := value.(type) {
		case lua.LNumber:
			out = append(out, StateMachineMapSpec{TemplateID: uint32(v), Opts: DefaultMapInitOpts()})
		case *lua.LTable:
			id, ok := v.RawGetString("id").(lua.LNumber)
			if ok == false {
				parseErr = fmt.Errorf("on_create map entry table needs a numeric id")
				return
			}
			out = append(out, StateMachineMapSpec{TemplateID: uint32(id), Opts: ParseMapInitOptsLua(v)})
		default:
			parseErr = fmt.Errorf("on_create map entry must be a map id")
		}
	})
	if parseErr != nil {
		return nil, parseErr
	}
	if len(out) == 0 {
		return nil, fmt.Errorf("on_create must return a non-empty map id array")
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
		if ch != nil {
			ch.unbindStateMachine(sm)
		}
	}
}

func (sm *StateMachine) Finish(ctx actor.Context, exitMapID uint32, exitPortal uint8) {
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

	sm.callHookNow(ctx, scriptPath, "on_finish")

	var exitMap *Map
	if exitMapID > 0 && group != nil && group.GameWorld != nil {
		exitMap = group.GameWorld.GetMapSystem().Get(exitMapID)
	}
	for _, ch := range players {
		if ch == nil {
			continue
		}
		ch.unbindStateMachine(sm)
		if exitMap != nil {
			_ = ch.Warp(ctx, exitMap, exitPortal)
		}
	}
	if group != nil && group.GameWorld != nil && pid != nil {
		group.GameWorld.SendStateMachineMessage(pid, &StopStateMachine{})
	}
}

func (sm *StateMachine) callHookNow(ctx actor.Context, scriptPath string, hook string, args ...interface{}) {
	if sm == nil || scriptPath == "" {
		return
	}
	root := luax.NewState()
	defer root.Close()
	thread, err := luax.NewThread(root, scriptPath)
	if err != nil {
		return
	}
	if !luax.HasFunc(thread, hook) {
		luax.Close(thread)
		return
	}
	luax.SetConfiguration(thread, luax.Configuration{
		ActorContext: ctx,
		ActorPID:     sm.ActorPID,
	})
	if _, err := luax.Call(thread, hook, append([]interface{}{sm}, args...)...); err != nil {
		name := ""
		if sm.Group != nil {
			name = sm.Group.Name
		}
		fmt.Printf("state machine %s hook %s: %v\n", name, hook, err)
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

type CreateStateMachineMaps struct {
	Specs []StateMachineMapSpec
	Err   string
}

type StateMachineMapRemoved struct {
	Map *Map
}

type EnterStateMachinePlayer struct {
	Character *Character
	Args      []interface{}
}

type StartStateMachine struct{}

type FinishStateMachine struct {
	ExitMapID  uint32
	ExitPortal uint8
}

type LeaveStateMachinePlayer struct {
	Character  *Character
	WarpLeaver bool
	Reason     StateMachineLeaveReason
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

func (ch *Character) Spectating() bool {
	sm := ch.GetMap().StateMachine()
	return sm != nil && sm.HasPlayer(ch) == false
}

func (ch *Character) bindStateMachine(sm *StateMachine) {
	ch.GameWorld.GetDispatchSystem().Call(ch.GetID(), func(actor.Context) {
		ch.stateMachine = sm
	})
}

func (ch *Character) unbindStateMachine(sm *StateMachine) {
	ch.GameWorld.GetDispatchSystem().Call(ch.GetID(), func(actor.Context) {
		if ch.stateMachine == sm {
			ch.stateMachine = nil
		}
	})
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
