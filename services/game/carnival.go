package game

import (
	"sync"
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/maple/constant"
)

type CarnivalState int

const (
	CarnivalStateEmpty   CarnivalState = 0
	CarnivalStateWaiting CarnivalState = 1
	CarnivalStateReady   CarnivalState = 2
	CarnivalStateBattle  CarnivalState = 3
	CarnivalStateReward  CarnivalState = 4
)

type CarnivalResult int

const (
	CarnivalResultUndecided CarnivalResult = 0
	CarnivalResultRedWin    CarnivalResult = 1
	CarnivalResultBlueWin   CarnivalResult = 2
	CarnivalResultDraw      CarnivalResult = 3
)

type CarnivalMember struct {
	id    uint32
	name  string
	level uint8
	class constant.ClassType
}

type CarnivalChallenge struct {
	leaderID    uint32
	partyID     uint32
	onMap       *Map
	members     []*CarnivalMember
	requestedAt time.Time
}

func (m *CarnivalMatch) PendingChallengeMembers() []*CarnivalMember {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.pending == nil {
		return nil
	}
	return m.pending.members
}

type CarnivalMatch struct {
	mu sync.Mutex

	slotIndex    int
	waitingMapID uint32
	fieldMapID   uint32
	reviveMapID  uint32
	winMapID     uint32
	loseMapID    uint32
	maxMembers   int
	state        CarnivalState
	teams        map[constant.CarnivalTeam]*CarnivalTeam
	pending      *CarnivalChallenge
	queue        []*CarnivalChallenge
	smID         string
	registry     *CarnivalSystem
}

func NewCarnivalMatch(slotIndex int, waitingMapID, fieldMapID, reviveMapID, winMapID, loseMapID uint32, maxMembers int) *CarnivalMatch {
	return &CarnivalMatch{
		slotIndex:    slotIndex,
		waitingMapID: waitingMapID,
		fieldMapID:   fieldMapID,
		reviveMapID:  reviveMapID,
		winMapID:     winMapID,
		loseMapID:    loseMapID,
		maxMembers:   maxMembers,
		state:        CarnivalStateEmpty,
		teams:        make(map[constant.CarnivalTeam]*CarnivalTeam),
	}
}

func (m *CarnivalMatch) StateMachine() *StateMachine {
	if m.smID == "" {
		return nil
	}
	group := m.registry.Group()
	if group == nil {
		return nil
	}
	return group.Get(m.smID)
}

func (m *CarnivalMatch) gameWorld() GameWorld {
	group := m.registry.Group()
	if group == nil {
		return nil
	}
	return group.gameWorld
}

func (m *CarnivalMatch) Team(teamID constant.CarnivalTeam) *CarnivalTeam {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.teams[teamID]
}

func (m *CarnivalMatch) FindTeam(ch *Character) *CarnivalTeam {
	m.mu.Lock()
	defer m.mu.Unlock()
	for _, team := range m.teams {
		if team != nil && team.HasMember(ch.GetID()) {
			return team
		}
	}
	return nil
}

func (m *CarnivalMatch) HasPendingChallenge() bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.pending != nil
}

func (m *CarnivalMatch) Result() CarnivalResult {
	m.mu.Lock()
	defer m.mu.Unlock()
	red := m.teams[constant.CarnivalTeamRed]
	blue := m.teams[constant.CarnivalTeamBlue]
	if red == nil || blue == nil {
		return CarnivalResultUndecided
	}
	if m.state != CarnivalStateBattle && m.state != CarnivalStateReward {
		return CarnivalResultUndecided
	}
	if red.totalCP > blue.totalCP {
		return CarnivalResultRedWin
	}
	if blue.totalCP > red.totalCP {
		return CarnivalResultBlueWin
	}
	return CarnivalResultDraw
}

func (m *CarnivalMatch) SetState(state CarnivalState) {
	m.mu.Lock()
	m.state = state
	m.mu.Unlock()
}

func (m *CarnivalMatch) GetState() CarnivalState {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.state
}

func (m *CarnivalMatch) ShowResult() {
	for _, teamID := range []constant.CarnivalTeam{constant.CarnivalTeamRed, constant.CarnivalTeamBlue} {
		if team := m.Team(teamID); team != nil {
			team.ShowResult()
		}
	}
}

func (m *CarnivalMatch) OnMobSpawn(mob *Mob) {
	if mob.carnivalTeam == constant.CarnivalTeamNone {
		return
	}
	team := m.Team(mob.carnivalTeam)
	if team == nil {
		return
	}
	team.buffGuardians(mob)
}

func (m *CarnivalMatch) DestroyGuardian(reactor *Reactor) {
	for _, teamID := range []constant.CarnivalTeam{constant.CarnivalTeamRed, constant.CarnivalTeamBlue} {
		if team := m.Team(teamID); team != nil && team.destroyGuardian(reactor) {
			return
		}
	}
}

func (m *CarnivalMatch) Finish(gw GameWorld) bool {
	m.mu.Lock()
	for _, team := range m.teams {
		if team != nil {
			team.Clear()
		}
	}
	m.mu.Unlock()
	return m.Conclude(gw)
}

// Conclude frees the slot but leaves each member bound to its detached team for the reward NPC.
func (m *CarnivalMatch) Conclude(gw GameWorld) bool {
	m.mu.Lock()
	for id, team := range m.teams {
		if team != nil {
			team.match = nil
		}
		delete(m.teams, id)
	}
	m.pending = nil
	m.queue = nil
	m.state = CarnivalStateEmpty
	smID := m.smID
	m.smID = ""
	m.mu.Unlock()

	if smID != "" {
		group := m.registry.Group()
		if group != nil {
			if sm := group.Get(smID); sm != nil && !sm.Disposed() {
				sm.AbortStart()
				group.gameWorld.Systems().Get[*StateMachineSystem]().Send(sm.actorPID, &StopStateMachine{})
			}
		}
	}
	return true
}

func (m *CarnivalMatch) AcceptPendingChallenge(gw GameWorld) bool {
	m.mu.Lock()
	challenge := m.pending
	if challenge == nil || m.teams[constant.CarnivalTeamRed] == nil || m.state != CarnivalStateWaiting {
		m.mu.Unlock()
		return false
	}
	m.mu.Unlock()

	gw.Systems().Get[*MapSystem]().Visit(challenge.onMap, func(actor.Context) {
		m.admitChallengers(challenge)
	})
	return true
}

func (m *CarnivalMatch) admitChallengers(challenge *CarnivalChallenge) {
	members := make([]*Character, 0, len(challenge.members))
	for _, member := range challenge.members {
		ch := challenge.onMap.GetPlayer(member.id)
		if ch == nil {
			m.dropChallenge(challenge)
			return
		}
		members = append(members, ch)
	}

	personal := make(map[uint32]*CarnivalPersonalCP, len(challenge.members))
	for _, member := range challenge.members {
		personal[member.id] = &CarnivalPersonalCP{}
	}
	blue := &CarnivalTeam{
		match:    m,
		teamID:   constant.CarnivalTeamBlue,
		leaderID: challenge.leaderID,
		partyID:  challenge.partyID,
		roster:   challenge.members,
		personal: personal,
	}

	m.mu.Lock()
	if m.pending != challenge || m.state != CarnivalStateWaiting {
		m.mu.Unlock()
		return
	}
	m.pending = nil
	m.queue = nil
	m.teams[constant.CarnivalTeamBlue] = blue
	m.state = CarnivalStateReady
	m.mu.Unlock()

	sm := m.StateMachine()
	for _, ch := range members {
		ch.carnival.team = blue
		sm.EnterPlayer(ch)
	}
	sm.CallHook("on_challenge_accepted")
}

func (m *CarnivalMatch) dropChallenge(challenge *CarnivalChallenge) {
	m.mu.Lock()
	if m.pending != challenge {
		m.mu.Unlock()
		return
	}
	m.queue = m.queue[1:]
	m.pending = nil
	if len(m.queue) > 0 {
		m.pending = m.queue[0]
	}
	m.mu.Unlock()

	m.StateMachine().CallHook("on_challenge_failed")
}

func (m *CarnivalMatch) RejectPendingChallenge() {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.pending == nil || len(m.queue) == 0 {
		return
	}
	m.queue = m.queue[1:]
	m.pending = nil
	if len(m.queue) > 0 {
		m.pending = m.queue[0]
	}
}

func (m *CarnivalMatch) enemyTeam(teamID constant.CarnivalTeam) *CarnivalTeam {
	if teamID == constant.CarnivalTeamRed {
		return m.Team(constant.CarnivalTeamBlue)
	}
	return m.Team(constant.CarnivalTeamRed)
}

func (m *CarnivalMatch) NotifyCarnivalStart(ch *Character) {
	team := m.FindTeam(ch)
	if team == nil {
		return
	}
	personal := team.PersonalCP(ch.GetID())
	if personal == nil {
		return
	}
	enemy := m.enemyTeam(team.teamID)
	enemyAvail, enemyTotal := 0, 0
	if enemy != nil {
		enemyAvail = enemy.availableCP
		enemyTotal = enemy.totalCP
	}
	ch.listener.OnCarnivalStart(ch, team.teamID,
		personal.availableCP, personal.totalCP,
		team.availableCP, team.totalCP,
		enemyAvail, enemyTotal)
}

func (m *CarnivalMatch) OnMobKilled(killer *Character, cpAmount int) {
	if cpAmount <= 0 {
		return
	}
	team := m.FindTeam(killer)
	if team == nil {
		return
	}
	team.AddCP(killer, cpAmount)
}

func (m *CarnivalMatch) OnPlayerDied(ch *Character) {
	team := m.FindTeam(ch)
	if team == nil {
		return
	}
	ch.listener.OnCarnivalDied(ch, team.teamID, ch.GetName(), 0)
}

func NewCarnivalMember(ch *Character) *CarnivalMember {
	return &CarnivalMember{
		id:    ch.GetID(),
		name:  ch.GetName(),
		level: ch.GetLevel(),
		class: ch.class,
	}
}
