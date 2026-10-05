package entity

import (
	"sync"
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/services/game/constant"
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
	ID    uint32
	Name  string
	Level uint8
	Class uint16
}

type CarnivalChallenge struct {
	LeaderID    uint32
	PartyID     uint32
	Map         *Map
	Members     []*CarnivalMember
	RequestedAt time.Time
}

func (m *CarnivalMatch) PendingChallengeMembers() []*CarnivalMember {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.Pending == nil {
		return nil
	}
	return m.Pending.Members
}

type CarnivalMatch struct {
	mu sync.Mutex

	SlotIndex    int
	WaitingMapID uint32
	FieldMapID   uint32
	ReviveMapID  uint32
	WinMapID     uint32
	LoseMapID    uint32
	MaxMembers   int
	State        CarnivalState
	Teams        map[constant.CarnivalTeam]*CarnivalTeam
	Pending      *CarnivalChallenge
	queue        []*CarnivalChallenge
	smID         string
	registry     *CarnivalRegistry
}

func NewCarnivalMatch(slotIndex int, waitingMapID, fieldMapID, reviveMapID, winMapID, loseMapID uint32, maxMembers int) *CarnivalMatch {
	return &CarnivalMatch{
		SlotIndex:    slotIndex,
		WaitingMapID: waitingMapID,
		FieldMapID:   fieldMapID,
		ReviveMapID:  reviveMapID,
		WinMapID:     winMapID,
		LoseMapID:    loseMapID,
		MaxMembers:   maxMembers,
		State:        CarnivalStateEmpty,
		Teams:        make(map[constant.CarnivalTeam]*CarnivalTeam),
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
	return group.GameWorld
}

func (m *CarnivalMatch) Team(teamID constant.CarnivalTeam) *CarnivalTeam {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.Teams[teamID]
}

func (m *CarnivalMatch) FindTeam(ch *Character) *CarnivalTeam {
	m.mu.Lock()
	defer m.mu.Unlock()
	for _, team := range m.Teams {
		if team != nil && team.HasMember(ch.GetID()) {
			return team
		}
	}
	return nil
}

func (m *CarnivalMatch) HasPendingChallenge() bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.Pending != nil
}

func (m *CarnivalMatch) Result() CarnivalResult {
	m.mu.Lock()
	defer m.mu.Unlock()
	red := m.Teams[constant.CarnivalTeamRed]
	blue := m.Teams[constant.CarnivalTeamBlue]
	if red == nil || blue == nil {
		return CarnivalResultUndecided
	}
	if m.State != CarnivalStateBattle && m.State != CarnivalStateReward {
		return CarnivalResultUndecided
	}
	if red.TotalCP > blue.TotalCP {
		return CarnivalResultRedWin
	}
	if blue.TotalCP > red.TotalCP {
		return CarnivalResultBlueWin
	}
	return CarnivalResultDraw
}

func (m *CarnivalMatch) SetState(state CarnivalState) {
	m.mu.Lock()
	m.State = state
	m.mu.Unlock()
}

func (m *CarnivalMatch) GetState() CarnivalState {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.State
}

func (m *CarnivalMatch) WarpAll(ctx actor.Context, mapID uint32, portalName string) bool {
	ok := false
	for _, team := range m.Teams {
		if team != nil {
			ok = team.Warp(ctx, mapID, portalName) || ok
		}
	}
	return ok
}

func (m *CarnivalMatch) Finish(gw GameWorld) bool {
	m.mu.Lock()
	for _, team := range m.Teams {
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
	for id, team := range m.Teams {
		if team != nil {
			team.Match = nil
		}
		delete(m.Teams, id)
	}
	m.Pending = nil
	m.queue = nil
	m.State = CarnivalStateEmpty
	smID := m.smID
	m.smID = ""
	m.mu.Unlock()

	if smID != "" {
		group := m.registry.Group()
		if group != nil {
			if sm := group.Get(smID); sm != nil {
				if !sm.Disposed() {
					sm.AbortStart()
				}
				group.RemoveMachine(sm)
			}
		}
	}
	return true
}

func (m *CarnivalMatch) AcceptPendingChallenge(gw GameWorld) bool {
	m.mu.Lock()
	challenge := m.Pending
	if challenge == nil || m.Teams[constant.CarnivalTeamRed] == nil || m.State != CarnivalStateWaiting {
		m.mu.Unlock()
		return false
	}
	m.mu.Unlock()

	gw.GetMapSystem().Call(challenge.Map, func(actor.Context) {
		m.admitChallengers(challenge)
	})
	return true
}

func (m *CarnivalMatch) admitChallengers(challenge *CarnivalChallenge) {
	members := make([]*Character, 0, len(challenge.Members))
	for _, member := range challenge.Members {
		ch := challenge.Map.GetPlayer(member.ID)
		if ch == nil {
			m.dropChallenge(challenge)
			return
		}
		members = append(members, ch)
	}

	personal := make(map[uint32]*CarnivalPersonalCP, len(challenge.Members))
	for _, member := range challenge.Members {
		personal[member.ID] = &CarnivalPersonalCP{}
	}
	blue := &CarnivalTeam{
		Match:    m,
		TeamID:   constant.CarnivalTeamBlue,
		LeaderID: challenge.LeaderID,
		PartyID:  challenge.PartyID,
		Roster:   challenge.Members,
		Personal: personal,
	}

	m.mu.Lock()
	if m.Pending != challenge || m.State != CarnivalStateWaiting {
		m.mu.Unlock()
		return
	}
	m.Pending = nil
	m.queue = nil
	m.Teams[constant.CarnivalTeamBlue] = blue
	m.State = CarnivalStateReady
	m.mu.Unlock()

	sm := m.StateMachine()
	for _, ch := range members {
		ch.BindCarnival(blue)
		sm.EnterPlayer(ch)
	}
	sm.CallHook("on_challenge_accepted")
}

func (m *CarnivalMatch) dropChallenge(challenge *CarnivalChallenge) {
	m.mu.Lock()
	if m.Pending != challenge {
		m.mu.Unlock()
		return
	}
	m.queue = m.queue[1:]
	m.Pending = nil
	if len(m.queue) > 0 {
		m.Pending = m.queue[0]
	}
	m.mu.Unlock()

	m.StateMachine().CallHook("on_challenge_failed")
}

func (m *CarnivalMatch) RejectPendingChallenge() {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.Pending == nil || len(m.queue) == 0 {
		return
	}
	m.queue = m.queue[1:]
	m.Pending = nil
	if len(m.queue) > 0 {
		m.Pending = m.queue[0]
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
	enemy := m.enemyTeam(team.TeamID)
	enemyAvail, enemyTotal := 0, 0
	if enemy != nil {
		enemyAvail = enemy.AvailableCP
		enemyTotal = enemy.TotalCP
	}
	ch.Listener.OnCarnivalStart(ch, team.TeamID,
		personal.AvailableCP, personal.TotalCP,
		team.AvailableCP, team.TotalCP,
		enemyAvail, enemyTotal)
}

func (m *CarnivalMatch) OnMobKilled(killer *Character, cpAmount int) {
	if cpAmount <= 0 {
		return
	}
	if m.GetState() != CarnivalStateBattle {
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
	ch.Listener.OnCarnivalDied(ch, team.TeamID, ch.GetName(), 0)
}
