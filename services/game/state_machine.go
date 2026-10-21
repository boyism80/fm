package game

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
	id              string
	group           *StateMachineGroup
	party           *Party
	leader          *Character
	scaleLevel      int
	minPlayers      int
	exitMapID       uint32
	players         []*Character
	playerSet       map[uint32]*Character
	props           map[string]string
	kills           map[uint32]int
	disposed        bool
	actorPID        *actor.PID
	maps            map[uint32]*Map
	mapRefs         map[*Map]*MapRef
	timeoutDeadline time.Time
}

var ErrStateMachineEntryDenied = errors.New("state machine map entry denied")

type StateMachineMapSpec struct {
	templateID uint32
	opts       MapInitOpts
}

func NewStateMachine(id string, group *StateMachineGroup) *StateMachine {
	return &StateMachine{
		id:        id,
		group:     group,
		players:   make([]*Character, 0),
		playerSet: make(map[uint32]*Character),
		props:     make(map[string]string),
		kills:     make(map[uint32]int),
		maps:      make(map[uint32]*Map),
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
	sm.bind(ch)
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
	if sm.disposed || !sm.unregisterLocked(ch) {
		sm.mu.Unlock()
		return false
	}
	count := len(sm.players)
	minPlayers := sm.minPlayers
	exitMapID := sm.exitMapID
	group := sm.group
	belowMin := minPlayers > 0 && count < minPlayers
	sm.mu.Unlock()

	if warpLeaver && exitMapID > 0 && group != nil && group.gameWorld != nil {
		if exitMap := group.gameWorld.Systems().Get[*MapSystem]().Get(exitMapID); exitMap != nil {
			_ = ch.Warp(ctx, exitMap, 0)
		}
	}
	if belowMin {
		if group != nil {
			sm.callHookNow(ctx, group.scriptPath, "on_player_leave", ch, string(reason))
		}
		sm.Finish(ctx, exitMapID, 0)
		return true
	}
	sm.CallHook("on_player_leave", ch, string(reason))
	return false
}

func (sm *StateMachine) RequestLeave(ch *Character, warpLeaver bool, reason StateMachineLeaveReason) {
	if sm == nil || ch == nil || sm.Disposed() || sm.group == nil || sm.group.gameWorld == nil || sm.actorPID == nil {
		return
	}
	sm.group.gameWorld.Systems().Get[*StateMachineSystem]().Send(sm.actorPID, &LeaveStateMachinePlayer{
		Character:  ch,
		WarpLeaver: warpLeaver,
		Reason:     reason,
	})
}

func (sm *StateMachine) RequestFinish(exitMapID uint32, exitPortal uint8) {
	if sm == nil || sm.Disposed() || sm.group == nil || sm.group.gameWorld == nil || sm.actorPID == nil {
		return
	}
	sm.group.gameWorld.Systems().Get[*StateMachineSystem]().Send(sm.actorPID, &FinishStateMachine{
		ExitMapID:  exitMapID,
		ExitPortal: exitPortal,
	})
}

func (sm *StateMachine) EnterPlayer(ch *Character, args ...interface{}) {
	if sm == nil || ch == nil || sm.Disposed() || sm.group == nil || sm.group.gameWorld == nil || sm.actorPID == nil {
		return
	}
	sm.group.gameWorld.Systems().Get[*StateMachineSystem]().Send(sm.actorPID, &EnterStateMachinePlayer{
		Character: ch,
		Args:      args,
	})
}

func (sm *StateMachine) Start() {
	if sm == nil || sm.Disposed() || sm.group == nil || sm.group.gameWorld == nil || sm.actorPID == nil {
		return
	}
	sm.group.gameWorld.Systems().Get[*StateMachineSystem]().Send(sm.actorPID, &StartStateMachine{})
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
	sm.unbind(ch)
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
	for _, ch := range append(sm.Players(), sm.leader) {
		if ch == nil || ch.gm.timerLimit == 0 {
			continue
		}
		if limit := int64(ch.gm.timerLimit) * 1000; limit < ms {
			ms = limit
		}
	}
	return ms
}

func (sm *StateMachine) StartTimer(ms int64) {
	if sm == nil || sm.group == nil || sm.group.gameWorld == nil || sm.actorPID == nil || ms <= 0 {
		return
	}
	sm.group.gameWorld.Systems().Get[*StateMachineSystem]().Send(sm.actorPID, &ScheduleStateMachineTimeout{Milliseconds: ms})
}

func (sm *StateMachine) StartTimerAsync(ctx actor.Context, ms int64) *async.Promise[interface{}] {
	if sm == nil || sm.group == nil || sm.group.gameWorld == nil || sm.actorPID == nil || ms <= 0 || ctx == nil {
		return nil
	}
	gw := sm.group.gameWorld
	pid := sm.actorPID
	return async.Ask(ctx, pid, 10*time.Second, func(replyTo *actor.PID) {
		gw.Systems().Get[*StateMachineSystem]().Send(pid, &ScheduleStateMachineTimeout{
			Milliseconds: ms,
			ReplyTo:      replyTo,
		})
	})
}

func (sm *StateMachine) StopTimer() {
	if sm == nil || sm.group == nil || sm.group.gameWorld == nil || sm.actorPID == nil {
		return
	}
	sm.group.gameWorld.Systems().Get[*StateMachineSystem]().Send(sm.actorPID, &CancelStateMachineTimeout{})
}

func (sm *StateMachine) StopTimerAsync(ctx actor.Context) *async.Promise[interface{}] {
	if sm == nil || sm.group == nil || sm.group.gameWorld == nil || sm.actorPID == nil || ctx == nil {
		return nil
	}
	gw := sm.group.gameWorld
	pid := sm.actorPID
	return async.Ask(ctx, pid, 10*time.Second, func(replyTo *actor.PID) {
		gw.Systems().Get[*StateMachineSystem]().Send(pid, &CancelStateMachineTimeout{ReplyTo: replyTo})
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
		ch.listener.OnClock(ch, seconds)
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
	ch.listener.OnClock(ch, int32(left/1000))
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
	ms := sm.group.gameWorld.Systems().Get[*MapSystem]()
	for _, spec := range specs {
		if sm.Map(spec.templateID) != nil {
			continue
		}
		m, err := ms.CreateStateMachineMap(spec.templateID, sm, spec.opts)
		if err != nil {
			return err
		}
		ref, err := m.Reserve()
		if err != nil {
			return err
		}
		sm.mu.Lock()
		sm.maps[spec.templateID] = m
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
	return sm.maps[templateID]
}

func (sm *StateMachine) CloseMaps() {
	ms := sm.group.gameWorld.Systems().Get[*MapSystem]()
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
	if sm.maps[m.TemplateID()] == m {
		delete(sm.maps, m.TemplateID())
	}
	return len(sm.maps) == 0
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
			out = append(out, StateMachineMapSpec{templateID: uint32(v), opts: DefaultMapInitOpts()})
		case *lua.LTable:
			id, ok := v.RawGetString("id").(lua.LNumber)
			if !ok {
				parseErr = fmt.Errorf("on_create map entry table needs a numeric id")
				return
			}
			out = append(out, StateMachineMapSpec{templateID: uint32(id), opts: ParseMapInitOptsLua(v)})
		default:
			parseErr = fmt.Errorf("on_create map entry must be a map id")
		}
	})
	if parseErr != nil {
		return nil, parseErr
	}
	return out, nil
}

func (sm *StateMachine) MapList() []*Map {
	if sm == nil {
		return nil
	}
	sm.mu.Lock()
	defer sm.mu.Unlock()
	maps := make([]*Map, 0, len(sm.maps))
	for _, m := range sm.maps {
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
			sm.unbind(ch)
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
	group := sm.group
	pid := sm.actorPID
	scriptPath := ""
	if group != nil {
		scriptPath = group.scriptPath
	}
	sm.players = nil
	sm.playerSet = make(map[uint32]*Character)
	sm.mu.Unlock()

	sm.callHookNow(ctx, scriptPath, "on_finish")

	var exitMap *Map
	if exitMapID > 0 && group != nil && group.gameWorld != nil {
		exitMap = group.gameWorld.Systems().Get[*MapSystem]().Get(exitMapID)
	}
	for _, ch := range players {
		if ch == nil {
			continue
		}
		sm.unbind(ch)
		if exitMap != nil {
			_ = ch.Warp(ctx, exitMap, exitPortal)
		}
	}
	if group != nil && group.gameWorld != nil && pid != nil {
		group.gameWorld.Systems().Get[*StateMachineSystem]().Send(pid, &StopStateMachine{})
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
		ActorPID:     sm.actorPID,
	})
	if _, err := luax.Call(thread, hook, append([]interface{}{sm}, args...)...); err != nil {
		name := ""
		if sm.group != nil {
			name = sm.group.name
		}
		fmt.Printf("state machine %s hook %s: %v\n", name, hook, err)
	}
}

func (sm *StateMachine) CallHook(hook string, args ...interface{}) {
	if sm == nil || sm.Disposed() || sm.group == nil || sm.group.gameWorld == nil || sm.actorPID == nil || hook == "" {
		return
	}
	sm.group.gameWorld.Systems().Get[*StateMachineSystem]().Send(sm.actorPID, &CallStateMachineHook{
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

func (ch *Character) Spectating() bool {
	m := ch.GetMap()
	if m == nil {
		return true
	}
	sm := m.StateMachine()
	return sm != nil && !sm.HasPlayer(ch)
}

func (sm *StateMachine) bind(ch *Character) {
	ch.gameWorld.Systems().Get[*DispatchSystem]().Visit(ch.GetID(), func(actor.Context, *Character) {
		ch.stateMachine = sm
	})
}

func (sm *StateMachine) unbind(ch *Character) {
	ch.gameWorld.Systems().Get[*DispatchSystem]().Visit(ch.GetID(), func(actor.Context, *Character) {
		if ch.stateMachine == sm {
			ch.stateMachine = nil
		}
	})
}
