package entity

import (
	"fmt"
	"strconv"
	"sync"

	"github.com/boyism80/fm/core/luax"
)

type StateMachineGroup struct {
	mu             sync.Mutex
	Name           string
	ScriptPath     string
	GameWorld      GameWorld
	props          map[string]string
	machines       map[string]*StateMachine
	bootMapID      uint32
	declaredMapIDs []uint32
	minPlayers     int
	exitMapID      uint32
}

func NewStateMachineGroup(name, scriptPath string, gw GameWorld) *StateMachineGroup {
	return &StateMachineGroup{
		Name:       name,
		ScriptPath: scriptPath,
		GameWorld:  gw,
		props:      make(map[string]string),
		machines:   make(map[string]*StateMachine),
		bootMapID:  180000000,
	}
}

func (g *StateMachineGroup) SetProperty(key, value string) {
	if g == nil {
		return
	}
	g.mu.Lock()
	defer g.mu.Unlock()
	if g.props == nil {
		g.props = make(map[string]string)
	}
	g.props[key] = value
}

func (g *StateMachineGroup) GetProperty(key string) string {
	if g == nil {
		return ""
	}
	g.mu.Lock()
	defer g.mu.Unlock()
	if g.props == nil {
		return ""
	}
	return g.props[key]
}

func (g *StateMachineGroup) GetMap(mapID uint32) *Map {
	if g == nil || g.GameWorld == nil {
		return nil
	}
	return g.GameWorld.GetMapSystem().Get(mapID)
}

func (g *StateMachineGroup) DeclareMaps(ids []uint32) {
	if g == nil {
		return
	}
	g.mu.Lock()
	defer g.mu.Unlock()
	g.declaredMapIDs = append(g.declaredMapIDs[:0], ids...)
}

func (g *StateMachineGroup) DeclaredMaps() []uint32 {
	if g == nil {
		return nil
	}
	g.mu.Lock()
	defer g.mu.Unlock()
	ids := make([]uint32, len(g.declaredMapIDs))
	copy(ids, g.declaredMapIDs)
	return ids
}

func (g *StateMachineGroup) DeclareMinPlayers(n int) {
	if g == nil {
		return
	}
	g.mu.Lock()
	defer g.mu.Unlock()
	if n < 0 {
		n = 0
	}
	g.minPlayers = n
}

func (g *StateMachineGroup) MinPlayers() int {
	if g == nil {
		return 1
	}
	g.mu.Lock()
	defer g.mu.Unlock()
	if g.minPlayers < 1 {
		return 1
	}
	return g.minPlayers
}

func (g *StateMachineGroup) Machines() []*StateMachine {
	if g == nil {
		return nil
	}
	g.mu.Lock()
	defer g.mu.Unlock()
	out := make([]*StateMachine, 0, len(g.machines))
	for _, sm := range g.machines {
		if sm != nil && !sm.Disposed() {
			out = append(out, sm)
		}
	}
	return out
}

func (g *StateMachineGroup) DeclareExitMap(mapID uint32) {
	if g == nil {
		return
	}
	g.mu.Lock()
	defer g.mu.Unlock()
	g.exitMapID = mapID
}

func (g *StateMachineGroup) ExitMapID() uint32 {
	if g == nil {
		return 0
	}
	g.mu.Lock()
	defer g.mu.Unlock()
	return g.exitMapID
}

func (g *StateMachineGroup) Get(id string) *StateMachine {
	if g == nil {
		return nil
	}
	g.mu.Lock()
	defer g.mu.Unlock()
	return g.machines[id]
}

func (g *StateMachineGroup) RemoveMachine(id string) {
	if g == nil {
		return
	}
	g.mu.Lock()
	defer g.mu.Unlock()
	delete(g.machines, id)
}

func (g *StateMachineGroup) bootMap() *Map {
	if g == nil || g.GameWorld == nil {
		return nil
	}
	ms := g.GameWorld.GetMapSystem()
	if ms == nil {
		return nil
	}
	if m := ms.Get(g.bootMapID); m != nil {
		return m
	}
	return ms.Get(0)
}

func (g *StateMachineGroup) CallGroupHook(hook string, args ...interface{}) {
	if g == nil || hook == "" {
		return
	}
	m := g.bootMap()
	if m == nil {
		return
	}
	root := m.EnsureLuaRoot(nil)
	if root == nil {
		return
	}
	thread, err := luax.NewThread(root, g.ScriptPath)
	if err != nil {
		return
	}
	if !luax.HasFunc(thread, hook) {
		luax.Close(thread)
		return
	}
	luax.SetConfiguration(thread, luax.Configuration{
		ActorPID: m.LogicActorPID(),
	})
	luax.CallAsync(root, thread, hook, args...).OnError(func(err error) {
		fmt.Printf("state machine group %s hook %s: %v\n", g.Name, hook, err)
	})
}

func (g *StateMachineGroup) Init() {
	g.CallGroupHook("on_init", g)
}

func (g *StateMachineGroup) CancelSchedule() {
	g.CallGroupHook("on_cancel_schedule", g)
	for _, sm := range g.Machines() {
		sm.CancelAllSchedules()
	}
}

func (g *StateMachineGroup) StartPersistent(id string) (*StateMachine, error) {
	if g == nil {
		return nil, fmt.Errorf("state machine group is nil")
	}
	if id == "" {
		id = "persistent"
	}
	sm, err := g.Create(id, CreateOpts{})
	if err != nil {
		return nil, err
	}
	sm.MinPlayers = 0
	return sm, nil
}

type StartPartyOpts struct {
	MaxLevel int
}

type CreateOpts struct {
	Leader     *Character
	Party      *Party
	ScaleLevel int
}

func (g *StateMachineGroup) Create(id string, opts CreateOpts) (*StateMachine, error) {
	if g == nil {
		return nil, fmt.Errorf("state machine group is nil")
	}
	if id == "" {
		return nil, fmt.Errorf("state machine id is empty")
	}
	if g.GameWorld == nil {
		return nil, fmt.Errorf("game world is nil")
	}
	if len(g.DeclaredMaps()) == 0 {
		return nil, fmt.Errorf("state machine group %s declared no maps", g.Name)
	}

	g.mu.Lock()
	existing := g.machines[id]
	g.mu.Unlock()
	if existing != nil && !existing.Disposed() {
		return nil, fmt.Errorf("state machine %s is already running", id)
	}

	sm := NewStateMachine(id, g)
	sm.Party = opts.Party
	sm.Leader = opts.Leader
	sm.ScaleLevel = opts.ScaleLevel
	sm.MinPlayers = g.MinPlayers()
	sm.ExitMapID = g.ExitMapID()

	sm.ActorPID = g.GameWorld.StartStateMachineActor(sm)
	if sm.ActorPID == nil {
		return nil, fmt.Errorf("failed to start state machine actor")
	}

	g.mu.Lock()
	if g.machines == nil {
		g.machines = make(map[string]*StateMachine)
	}
	g.machines[id] = sm
	g.mu.Unlock()

	g.GameWorld.SendStateMachineMessage(sm.ActorPID, &BootstrapStateMachine{})
	return sm, nil
}

func (g *StateMachineGroup) EnterPlayer(sm *StateMachine, ch *Character) {
	if sm == nil {
		return
	}
	sm.EnterPlayer(ch)
}

func (g *StateMachineGroup) Start(sm *StateMachine) {
	if sm == nil {
		return
	}
	sm.Start()
}

func (g *StateMachineGroup) StartParty(leader *Character, party *Party, opts StartPartyOpts) (*StateMachine, error) {
	if g == nil {
		return nil, fmt.Errorf("state machine group is nil")
	}
	if leader == nil {
		return nil, fmt.Errorf("leader is nil")
	}
	if party == nil {
		return nil, fmt.Errorf("party is nil")
	}
	leaderMap := leader.GetMap()
	if leaderMap == nil {
		return nil, fmt.Errorf("leader map is nil")
	}

	avgLevel, count := 0, 0
	members := make([]*Character, 0)
	for _, mem := range party.GetMembers() {
		if mem == nil {
			continue
		}
		ch := leaderMap.GetPlayer(mem.GetCharacterId())
		if ch == nil {
			continue
		}
		avgLevel += int(ch.GetLevel())
		count++
		members = append(members, ch)
	}
	if count <= 0 {
		return nil, fmt.Errorf("no party members on map")
	}
	avgLevel /= count
	scale := avgLevel
	if opts.MaxLevel > 0 && scale > opts.MaxLevel {
		scale = opts.MaxLevel
	}

	id := strconv.FormatUint(uint64(party.GetPartyId()), 10)
	sm, err := g.Create(id, CreateOpts{
		Leader:     leader,
		Party:      party,
		ScaleLevel: scale,
	})
	if err != nil {
		return nil, err
	}
	for _, ch := range members {
		g.EnterPlayer(sm, ch)
	}
	g.Start(sm)
	return sm, nil
}

func (g *StateMachineGroup) StartSolo(player *Character, opts StartPartyOpts) (*StateMachine, error) {
	if g == nil {
		return nil, fmt.Errorf("state machine group is nil")
	}
	if player == nil {
		return nil, fmt.Errorf("player is nil")
	}
	if player.GetMap() == nil {
		return nil, fmt.Errorf("player map is nil")
	}

	scale := int(player.GetLevel())
	if opts.MaxLevel > 0 && scale > opts.MaxLevel {
		scale = opts.MaxLevel
	}

	id := strconv.FormatUint(uint64(player.GetID()), 10)
	sm, err := g.Create(id, CreateOpts{
		Leader:     player,
		ScaleLevel: scale,
	})
	if err != nil {
		return nil, err
	}
	g.EnterPlayer(sm, player)
	g.Start(sm)
	return sm, nil
}
