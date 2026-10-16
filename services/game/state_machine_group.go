package game

import (
	"errors"
	"fmt"
	"strconv"
	"sync"

	"github.com/boyism80/fm/core/luax"
)

type StateMachineGroup struct {
	mu          sync.Mutex
	name        string
	scriptPath  string
	gameWorld   GameWorld
	props       map[string]string
	machines    map[string]*StateMachine
	bootMapID   uint32
	minPlayers  int
	exitMapID   uint32
	machineCap  MachineCap
	capOverride *MachineCap
}

type MachineCap struct {
	limited bool
	max     int
}

var ErrStateMachineLimit = errors.New("state machine limit reached")

func NewStateMachineGroup(name, scriptPath string, gw GameWorld) *StateMachineGroup {
	return &StateMachineGroup{
		name:       name,
		scriptPath: scriptPath,
		gameWorld:  gw,
		props:      make(map[string]string),
		machines:   make(map[string]*StateMachine),
		bootMapID:  180000000,
		machineCap: MachineCap{limited: true, max: 1},
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

func (g *StateMachineGroup) SetMaxMachines(n int) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.machineCap = MachineCap{limited: true, max: n}
}

func (g *StateMachineGroup) RemoveMachineLimit() {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.machineCap = MachineCap{}
}

func (g *StateMachineGroup) OverrideMachineCap(c MachineCap) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.capOverride = &c
}

func (g *StateMachineGroup) ClearMachineCapOverride() {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.capOverride = nil
}

func (g *StateMachineGroup) MachineCap() (MachineCap, bool) {
	g.mu.Lock()
	defer g.mu.Unlock()
	if g.capOverride != nil {
		return *g.capOverride, true
	}
	return g.machineCap, false
}

func (g *StateMachineGroup) SetMinPlayers(n int) {
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

func (g *StateMachineGroup) SetExitMap(mapID uint32) {
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

func (g *StateMachineGroup) RemoveMachine(sm *StateMachine) {
	if g == nil {
		return
	}
	g.mu.Lock()
	defer g.mu.Unlock()
	if g.machines[sm.id] != sm {
		return
	}
	delete(g.machines, sm.id)
}

func (g *StateMachineGroup) bootMap() *Map {
	if g == nil || g.gameWorld == nil {
		return nil
	}
	ms := g.gameWorld.Systems().Get[*MapSystem]()
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
	root := m.LuaRoot(nil)
	if root == nil {
		return
	}
	thread, err := luax.NewThread(root, g.scriptPath)
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
	luax.CallAsync(nil, root, thread, hook, args...).OnError(func(err error) {
		fmt.Printf("state machine group %s hook %s: %v\n", g.name, hook, err)
	})
}

func (g *StateMachineGroup) Init() {
	g.CallGroupHook("on_init", g)
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
	sm.minPlayers = 0
	return sm, nil
}

type StartPartyOpts struct {
	maxLevel int
}

type CreateOpts struct {
	leader     *Character
	party      *Party
	scaleLevel int
}

func (g *StateMachineGroup) Create(id string, opts CreateOpts) (*StateMachine, error) {
	if g == nil {
		return nil, fmt.Errorf("state machine group is nil")
	}
	if id == "" {
		return nil, fmt.Errorf("state machine id is empty")
	}
	if g.gameWorld == nil {
		return nil, fmt.Errorf("game world is nil")
	}

	sm := NewStateMachine(id, g)
	sm.party = opts.party
	sm.leader = opts.leader
	sm.scaleLevel = opts.scaleLevel
	sm.minPlayers = g.MinPlayers()
	sm.exitMapID = g.ExitMapID()
	sm.actorPID = g.gameWorld.Systems().Get[*StateMachineSystem]().StartActor(sm)
	if sm.actorPID == nil {
		return nil, fmt.Errorf("failed to start state machine actor")
	}

	g.mu.Lock()
	err := g.checkCreateLocked(id)
	if err == nil {
		g.machines[id] = sm
	}
	g.mu.Unlock()
	if err != nil {
		g.gameWorld.Systems().Get[*StateMachineSystem]().StopActor(sm)
		return nil, err
	}

	g.gameWorld.Systems().Get[*StateMachineSystem]().Send(sm.actorPID, &BootstrapStateMachine{})
	return sm, nil
}

func (g *StateMachineGroup) checkCreateLocked(id string) error {
	if existing := g.machines[id]; existing != nil && !existing.Disposed() {
		return fmt.Errorf("state machine %s is already running", id)
	}
	c := g.machineCap
	if g.capOverride != nil {
		c = *g.capOverride
	}
	if !c.limited {
		return nil
	}
	running := 0
	for _, sm := range g.machines {
		if !sm.Disposed() {
			running++
		}
	}
	if running >= c.max {
		return ErrStateMachineLimit
	}
	return nil
}

func (g *StateMachineGroup) Full() bool {
	if g == nil {
		return true
	}
	g.mu.Lock()
	defer g.mu.Unlock()
	c := g.machineCap
	if g.capOverride != nil {
		c = *g.capOverride
	}
	if !c.limited {
		return false
	}
	running := 0
	for _, sm := range g.machines {
		if sm != nil && !sm.Disposed() {
			running++
		}
	}
	return running >= c.max
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
	if opts.maxLevel > 0 && scale > opts.maxLevel {
		scale = opts.maxLevel
	}

	id := "p" + strconv.FormatUint(uint64(party.GetPartyId()), 10)
	sm, err := g.Create(id, CreateOpts{
		leader:     leader,
		party:      party,
		scaleLevel: scale,
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
	if opts.maxLevel > 0 && scale > opts.maxLevel {
		scale = opts.maxLevel
	}

	id := "c" + strconv.FormatUint(uint64(player.GetID()), 10)
	sm, err := g.Create(id, CreateOpts{
		leader:     player,
		scaleLevel: scale,
	})
	if err != nil {
		return nil, err
	}
	g.EnterPlayer(sm, player)
	g.Start(sm)
	return sm, nil
}
