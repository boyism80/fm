package entity

import (
	"strconv"
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

type CarnivalPersonalCP struct {
	AvailableCP int
	TotalCP     int
}

type CarnivalChallenge struct {
	LeaderID    uint32
	PartyID     uint32
	MemberIDs   []uint32
	RequestedAt time.Time
}

func (m *CarnivalMatch) PendingChallengeMembers(gw GameWorld) ([]*Character, int) {
	if m == nil {
		return nil, 0
	}
	m.mu.Lock()
	challenge := m.Pending
	m.mu.Unlock()
	if challenge == nil {
		return nil, 0
	}
	size := len(challenge.MemberIDs)
	if gw == nil || gw.GetMapSystem() == nil {
		return nil, size
	}
	ms := gw.GetMapSystem()
	out := make([]*Character, 0, size)
	seen := make(map[uint32]bool)
	add := func(ch *Character) {
		if ch == nil || seen[ch.GetID()] {
			return
		}
		out = append(out, ch)
		seen[ch.GetID()] = true
	}
	tryMap := func(mapID uint32) {
		mp := ms.Get(mapID)
		if mp == nil {
			return
		}
		for _, id := range challenge.MemberIDs {
			if !seen[id] {
				add(mp.GetPlayer(id))
			}
		}
	}
	tryMap(m.WaitingMapID)
	if challenge.PartyID != 0 {
		if party := gw.GetPartySystem().Get(challenge.PartyID); party != nil {
			for _, mem := range party.GetMembers() {
				tryMap(mem.GetMapId())
			}
		}
	}
	return out, size
}

type CarnivalTeam struct {
	Match       *CarnivalMatch
	TeamID      constant.CarnivalTeam
	LeaderID    uint32
	PartyID     uint32
	MemberIDs   []uint32
	AvailableCP int
	TotalCP     int
	Personal    map[uint32]*CarnivalPersonalCP
	Winner      bool
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

func (m *CarnivalMatch) SetStateMachineID(id string) {
	if m == nil {
		return
	}
	m.mu.Lock()
	m.smID = id
	m.mu.Unlock()
}

func (m *CarnivalMatch) StateMachineID() string {
	if m == nil {
		return ""
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.smID
}

func (m *CarnivalMatch) StateMachine() *StateMachine {
	if m == nil || m.smID == "" {
		return nil
	}
	group := GlobalCarnivalRegistry().Group()
	if group == nil {
		return nil
	}
	return group.Get(m.smID)
}

func (m *CarnivalMatch) gameWorld() GameWorld {
	group := GlobalCarnivalRegistry().Group()
	if group == nil {
		return nil
	}
	return group.GameWorld
}

func (m *CarnivalMatch) Team(teamID constant.CarnivalTeam) *CarnivalTeam {
	if m == nil {
		return nil
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.Teams[teamID]
}

func (m *CarnivalMatch) TeamOf(ch *Character) *CarnivalTeam {
	if m == nil || ch == nil {
		return nil
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	for _, team := range m.Teams {
		if team != nil && team.HasMember(ch.GetID()) {
			return team
		}
	}
	return nil
}

func (m *CarnivalMatch) TeamIDOf(ch *Character) constant.CarnivalTeam {
	team := m.TeamOf(ch)
	if team == nil {
		return constant.CarnivalTeamNone
	}
	return team.TeamID
}

func (m *CarnivalMatch) HasPendingChallenge() bool {
	if m == nil {
		return false
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.Pending != nil
}

func (m *CarnivalMatch) Result() CarnivalResult {
	if m == nil {
		return CarnivalResultUndecided
	}
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
	if m == nil {
		return
	}
	m.mu.Lock()
	m.State = state
	m.mu.Unlock()
}

func (m *CarnivalMatch) GetState() CarnivalState {
	if m == nil {
		return CarnivalStateEmpty
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.State
}

func (t *CarnivalTeam) HasMember(characterID uint32) bool {
	if t == nil {
		return false
	}
	for _, id := range t.MemberIDs {
		if id == characterID {
			return true
		}
	}
	return false
}

func (t *CarnivalTeam) Members(gw GameWorld) []*Character {
	if t == nil || gw == nil {
		return nil
	}
	out := make([]*Character, 0, len(t.MemberIDs))
	seen := make(map[uint32]bool)
	add := func(ch *Character) {
		if ch == nil || seen[ch.GetID()] {
			return
		}
		out = append(out, ch)
		seen[ch.GetID()] = true
	}
	if sm := t.Match.StateMachine(); sm != nil {
		for _, ch := range sm.Players() {
			if ch != nil && t.HasMember(ch.GetID()) {
				add(ch)
			}
		}
		if len(out) == len(t.MemberIDs) {
			return out
		}
	}
	ms := gw.GetMapSystem()
	tryMap := func(mapID uint32) {
		m := ms.Get(mapID)
		if m == nil {
			return
		}
		for _, id := range t.MemberIDs {
			if !seen[id] {
				add(m.GetPlayer(id))
			}
		}
	}
	tryMap(t.Match.WaitingMapID)
	tryMap(t.Match.FieldMapID)
	tryMap(t.Match.ReviveMapID)
	tryMap(t.Match.WinMapID)
	tryMap(t.Match.LoseMapID)
	if t.PartyID != 0 {
		if party := gw.GetPartySystem().Get(t.PartyID); party != nil {
			for _, mem := range party.GetMembers() {
				tryMap(mem.GetMapId())
			}
		}
	}
	return out
}

func (t *CarnivalTeam) Leader(gw GameWorld) *Character {
	if t == nil || gw == nil {
		return nil
	}
	for _, ch := range t.Members(gw) {
		if ch.GetID() == t.LeaderID {
			return ch
		}
	}
	return nil
}

func (t *CarnivalTeam) PersonalCP(characterID uint32) *CarnivalPersonalCP {
	if t == nil || t.Personal == nil {
		return nil
	}
	return t.Personal[characterID]
}

func (t *CarnivalTeam) AddCP(ch *Character, amount int) bool {
	if t == nil || ch == nil || amount <= 0 {
		return false
	}
	personal := t.PersonalCP(ch.GetID())
	if personal == nil {
		return false
	}
	personal.AvailableCP += amount
	personal.TotalCP += amount
	t.AvailableCP += amount
	t.TotalCP += amount
	ch.Listener.OnCarnivalObtainedCP(ch, personal.AvailableCP, personal.TotalCP)
	ch.Listener.OnCarnivalPartyCP(ch, t.TeamID, t.AvailableCP, t.TotalCP)
	t.broadcastPartyCP(ch)
	return true
}

func (t *CarnivalTeam) UseCP(ch *Character, amount int) bool {
	if t == nil || ch == nil || amount <= 0 {
		return false
	}
	personal := t.PersonalCP(ch.GetID())
	if personal == nil {
		return false
	}
	if personal.AvailableCP < amount || t.AvailableCP < amount {
		return false
	}
	personal.AvailableCP -= amount
	t.AvailableCP -= amount
	ch.Listener.OnCarnivalObtainedCP(ch, personal.AvailableCP, personal.TotalCP)
	ch.Listener.OnCarnivalPartyCP(ch, t.TeamID, t.AvailableCP, t.TotalCP)
	t.broadcastPartyCP(ch)
	return true
}

func (t *CarnivalTeam) broadcastPartyCP(origin *Character) {
	if t == nil || t.Match == nil {
		return
	}
	sm := t.Match.StateMachine()
	if sm == nil {
		return
	}
	for _, p := range sm.Players() {
		if p == nil || p.GetID() == origin.GetID() {
			continue
		}
		p.Listener.OnCarnivalPartyCP(p, t.TeamID, t.AvailableCP, t.TotalCP)
	}
}

func (t *CarnivalTeam) SetWinner(flag bool) {
	if t == nil {
		return
	}
	t.Winner = flag
}

func (t *CarnivalTeam) IsWinner() bool {
	if t == nil {
		return false
	}
	return t.Winner
}

func (t *CarnivalTeam) Warp(ctx actor.Context, mapID uint32, portalName string) bool {
	if t == nil || t.Match == nil {
		return false
	}
	gw := t.Match.gameWorld()
	if gw == nil {
		return false
	}
	dest := gw.GetMapSystem().Get(mapID)
	if dest == nil {
		return false
	}
	spawn := uint8(0)
	if portalName != "" && dest.Wz != nil {
		for _, p := range dest.Wz.Portals {
			if p.Name == portalName {
				spawn = p.ID
				break
			}
		}
	}
	for _, ch := range t.Members(gw) {
		_ = ch.Warp(ctx, dest, spawn)
	}
	return true
}

func (t *CarnivalTeam) Clear(gw GameWorld) {
	if t == nil {
		return
	}
	for _, ch := range t.Members(gw) {
		ch.UnbindCarnival()
	}
	t.MemberIDs = nil
	t.LeaderID = 0
	t.AvailableCP = 0
	t.TotalCP = 0
	t.Personal = nil
	t.Winner = false
}

func (m *CarnivalMatch) WarpAll(ctx actor.Context, mapID uint32, portalName string) bool {
	if m == nil {
		return false
	}
	ok := false
	for _, team := range m.Teams {
		if team != nil {
			ok = team.Warp(ctx, mapID, portalName) || ok
		}
	}
	return ok
}

func (m *CarnivalMatch) Finish(gw GameWorld) bool {
	if m == nil {
		return false
	}
	m.mu.Lock()
	for id, team := range m.Teams {
		if team != nil {
			team.Clear(gw)
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
		group := GlobalCarnivalRegistry().Group()
		if sm := group.Get(smID); sm != nil && !sm.Disposed() {
			sm.AbortStart()
		}
		group.RemoveMachine(smID)
	}
	return true
}

type CarnivalRegistry struct {
	mu             sync.RWMutex
	matches        map[int]*CarnivalMatch
	byMapID        map[uint32]*CarnivalMatch
	group          *StateMachineGroup
	skillHitChance map[uint32]int
}

func NewCarnivalRegistry() *CarnivalRegistry {
	return &CarnivalRegistry{
		matches: make(map[int]*CarnivalMatch),
		byMapID: make(map[uint32]*CarnivalMatch),
	}
}

func (r *CarnivalRegistry) SetSkillHitChance(skillID uint32, chance int) {
	if r == nil {
		return
	}
	if chance < 0 {
		chance = 0
	}
	if chance > 100 {
		chance = 100
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.skillHitChance == nil {
		r.skillHitChance = make(map[uint32]int)
	}
	r.skillHitChance[skillID] = chance
}

func (r *CarnivalRegistry) SkillHitChance(skillID uint32, fallback int) int {
	if r == nil {
		return fallback
	}
	r.mu.RLock()
	defer r.mu.RUnlock()
	if v, ok := r.skillHitChance[skillID]; ok {
		return v
	}
	return fallback
}

func (r *CarnivalRegistry) registerMaps(match *CarnivalMatch) {
	ids := []uint32{
		match.WaitingMapID,
		match.FieldMapID,
		match.ReviveMapID,
		match.WinMapID,
		match.LoseMapID,
	}
	for _, id := range ids {
		r.byMapID[id] = match
	}
}

func (r *CarnivalRegistry) unregisterMaps(match *CarnivalMatch) {
	ids := []uint32{
		match.WaitingMapID,
		match.FieldMapID,
		match.ReviveMapID,
		match.WinMapID,
		match.LoseMapID,
	}
	for _, id := range ids {
		delete(r.byMapID, id)
	}
}

func (r *CarnivalRegistry) BindGroup(group *StateMachineGroup) {
	r.mu.Lock()
	r.group = group
	r.mu.Unlock()
}

func (r *CarnivalRegistry) Group() *StateMachineGroup {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.group
}

func (r *CarnivalRegistry) Register(slotIndex int, waitingMapID, fieldMapID, reviveMapID, winMapID, loseMapID uint32, maxMembers int) *CarnivalMatch {
	if r == nil {
		return nil
	}
	match := NewCarnivalMatch(slotIndex, waitingMapID, fieldMapID, reviveMapID, winMapID, loseMapID, maxMembers)
	r.mu.Lock()
	r.matches[slotIndex] = match
	r.registerMaps(match)
	r.mu.Unlock()
	return match
}

func (r *CarnivalRegistry) Match(slotIndex int) *CarnivalMatch {
	if r == nil {
		return nil
	}
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.matches[slotIndex]
}

func (r *CarnivalRegistry) Matches() []*CarnivalMatch {
	if r == nil {
		return nil
	}
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]*CarnivalMatch, 0, len(r.matches))
	for _, m := range r.matches {
		out = append(out, m)
	}
	return out
}

func (r *CarnivalRegistry) MatchByMap(mapID uint32) *CarnivalMatch {
	if r == nil {
		return nil
	}
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.byMapID[mapID]
}

func (r *CarnivalRegistry) Enter(slotIndex int, leader *Character) bool {
	if r == nil || leader == nil {
		return false
	}
	match := r.Match(slotIndex)
	if match == nil {
		return false
	}
	partyID := leader.GetPartyID()
	if partyID == nil {
		return false
	}
	gw := leader.GameWorld
	if gw == nil {
		return false
	}
	party := gw.GetPartySystem().Get(*partyID)
	if party == nil {
		return false
	}
	if party.GetLeaderCharacterId() != leader.GetID() {
		return false
	}
	leaderMap := leader.GetMap()
	if leaderMap == nil {
		return false
	}
	members := make([]*Character, 0)
	for _, mem := range party.GetMembers() {
		if mem == nil {
			continue
		}
		ch := leaderMap.GetPlayer(mem.GetCharacterId())
		if ch == nil {
			return false
		}
		members = append(members, ch)
	}
	if len(members) == 0 || len(members) > match.MaxMembers {
		return false
	}

	memberIDs := make([]uint32, len(members))
	for i, ch := range members {
		memberIDs[i] = ch.GetID()
	}
	personal := make(map[uint32]*CarnivalPersonalCP, len(memberIDs))
	for _, id := range memberIDs {
		personal[id] = &CarnivalPersonalCP{}
	}
	red := &CarnivalTeam{
		Match:     match,
		TeamID:    constant.CarnivalTeamRed,
		LeaderID:  leader.GetID(),
		PartyID:   party.GetPartyId(),
		MemberIDs: memberIDs,
		Personal:  personal,
	}

	smID := strconv.FormatUint(uint64(match.WaitingMapID), 10)
	match.mu.Lock()
	if match.State != CarnivalStateEmpty || match.Teams[constant.CarnivalTeamRed] != nil {
		match.mu.Unlock()
		return false
	}
	match.Teams[constant.CarnivalTeamRed] = red
	match.State = CarnivalStateWaiting
	match.smID = smID
	match.mu.Unlock()

	for _, ch := range members {
		ch.BindCarnival(match)
	}
	return true
}

func (r *CarnivalRegistry) Challenge(slotIndex int, leader *Character) (ok bool, shouldOpen bool) {
	if r == nil || leader == nil {
		return false, false
	}
	match := r.Match(slotIndex)
	if match == nil {
		return false, false
	}
	partyID := leader.GetPartyID()
	if partyID == nil {
		return false, false
	}
	gw := leader.GameWorld
	if gw == nil {
		return false, false
	}
	party := gw.GetPartySystem().Get(*partyID)
	if party == nil {
		return false, false
	}
	if party.GetLeaderCharacterId() != leader.GetID() {
		return false, false
	}
	leaderMap := leader.GetMap()
	if leaderMap == nil {
		return false, false
	}
	members := make([]*Character, 0)
	for _, mem := range party.GetMembers() {
		if mem == nil {
			continue
		}
		ch := leaderMap.GetPlayer(mem.GetCharacterId())
		if ch == nil {
			return false, false
		}
		members = append(members, ch)
	}
	if len(members) == 0 {
		return false, false
	}

	memberIDs := make([]uint32, len(members))
	for i, ch := range members {
		memberIDs[i] = ch.GetID()
	}
	challenge := &CarnivalChallenge{
		LeaderID:    leader.GetID(),
		PartyID:     party.GetPartyId(),
		MemberIDs:   memberIDs,
		RequestedAt: time.Now(),
	}

	match.mu.Lock()
	red := match.Teams[constant.CarnivalTeamRed]
	if match.State != CarnivalStateWaiting || red == nil {
		match.mu.Unlock()
		return false, false
	}
	if len(members) != len(red.MemberIDs) {
		match.mu.Unlock()
		return false, false
	}
	wasEmpty := len(match.queue) == 0
	match.queue = append(match.queue, challenge)
	if wasEmpty {
		match.Pending = challenge
	}
	match.mu.Unlock()

	if !wasEmpty {
		return true, false
	}
	redLeader := red.Leader(gw)
	if redLeader == nil || redLeader.GetDialog() != nil {
		return true, false
	}
	return true, true
}

func (m *CarnivalMatch) AcceptPendingChallenge(gw GameWorld) bool {
	if m == nil || gw == nil {
		return false
	}
	m.mu.Lock()
	challenge := m.Pending
	if challenge == nil || m.Teams[constant.CarnivalTeamRed] == nil || m.State != CarnivalStateWaiting {
		m.mu.Unlock()
		return false
	}
	m.mu.Unlock()

	blueMemberIDs := append([]uint32(nil), challenge.MemberIDs...)
	personal := make(map[uint32]*CarnivalPersonalCP, len(blueMemberIDs))
	for _, id := range blueMemberIDs {
		personal[id] = &CarnivalPersonalCP{}
	}
	blue := &CarnivalTeam{
		Match:     m,
		TeamID:    constant.CarnivalTeamBlue,
		LeaderID:  challenge.LeaderID,
		PartyID:   challenge.PartyID,
		MemberIDs: blueMemberIDs,
		Personal:  personal,
	}
	found := blue.Members(gw)
	if len(found) != len(blue.MemberIDs) {
		return false
	}

	m.mu.Lock()
	if m.Pending != challenge || m.State != CarnivalStateWaiting {
		m.mu.Unlock()
		return false
	}
	m.Pending = nil
	m.queue = nil
	m.Teams[constant.CarnivalTeamBlue] = blue
	m.State = CarnivalStateReady
	m.mu.Unlock()

	for _, ch := range found {
		ch.BindCarnival(m)
	}
	m.StateMachine().CallHook("on_challenge_accepted")
	return true
}

func (m *CarnivalMatch) RejectPendingChallenge() (ok bool, shouldOpen bool) {
	if m == nil {
		return false, false
	}
	m.mu.Lock()
	if m.Pending == nil || len(m.queue) == 0 {
		m.mu.Unlock()
		return false, false
	}
	m.queue = m.queue[1:]
	if len(m.queue) == 0 {
		m.Pending = nil
		m.mu.Unlock()
		return true, false
	}
	m.Pending = m.queue[0]
	m.mu.Unlock()
	return true, true
}

func (m *CarnivalMatch) GrantWinnerCP() {
	if m == nil {
		return
	}
	result := m.Result()
	m.mu.Lock()
	defer m.mu.Unlock()
	if red := m.Teams[constant.CarnivalTeamRed]; red != nil {
		red.Winner = result == CarnivalResultRedWin
	}
	if blue := m.Teams[constant.CarnivalTeamBlue]; blue != nil {
		blue.Winner = result == CarnivalResultBlueWin
	}
}

func (m *CarnivalMatch) NotifyCarnivalStart(ch *Character) {
	if m == nil || ch == nil {
		return
	}
	team := m.TeamOf(ch)
	if team == nil {
		return
	}
	personal := team.PersonalCP(ch.GetID())
	if personal == nil {
		return
	}
	friendTeam := m.Team(team.TeamID)
	var enemyTeam *CarnivalTeam
	if team.TeamID == constant.CarnivalTeamRed {
		enemyTeam = m.Team(constant.CarnivalTeamBlue)
	} else {
		enemyTeam = m.Team(constant.CarnivalTeamRed)
	}
	friendAvail, friendTotal := 0, 0
	enemyAvail, enemyTotal := 0, 0
	if friendTeam != nil {
		friendAvail = friendTeam.AvailableCP
		friendTotal = friendTeam.TotalCP
	}
	if enemyTeam != nil {
		enemyAvail = enemyTeam.AvailableCP
		enemyTotal = enemyTeam.TotalCP
	}
	ch.Listener.OnCarnivalStart(ch, team.TeamID,
		personal.AvailableCP, personal.TotalCP,
		friendAvail, friendTotal,
		enemyAvail, enemyTotal)
}

func (m *CarnivalMatch) OnMobKilled(killer *Character, cpAmount int) {
	if m == nil || killer == nil || cpAmount <= 0 {
		return
	}
	if m.GetState() != CarnivalStateBattle {
		return
	}
	team := m.TeamOf(killer)
	if team == nil {
		return
	}
	team.AddCP(killer, cpAmount)
}

func (m *CarnivalMatch) OnPlayerDied(ch *Character) {
	if m == nil || ch == nil {
		return
	}
	team := m.TeamOf(ch)
	if team == nil {
		return
	}
	ch.Listener.OnCarnivalDied(ch, team.TeamID, ch.GetName(), 0)
}

func (m *CarnivalMatch) DisposeAll(gw GameWorld, exitMapID uint32) {
	if m == nil {
		return
	}
	if sm := m.StateMachine(); sm != nil {
		for _, ch := range sm.Players() {
			ch.UnbindCarnival()
		}
	}
	m.Finish(gw)
}
