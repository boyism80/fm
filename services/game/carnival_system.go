package game

import (
	"strconv"
	"sync"
	"time"

	"github.com/boyism80/fm/maple/constant"
)

type CarnivalSystem struct {
	mu             sync.RWMutex
	matches        map[int]*CarnivalMatch
	mapMatches     map[uint32]*CarnivalMatch
	group          *StateMachineGroup
	skillHitChance map[uint32]int
}

func (r *CarnivalSystem) Start() error { return nil }
func (r *CarnivalSystem) Stop()        {}

func NewCarnivalSystem() *CarnivalSystem {
	return &CarnivalSystem{
		matches:    make(map[int]*CarnivalMatch),
		mapMatches: make(map[uint32]*CarnivalMatch),
	}
}

func (r *CarnivalSystem) SetSkillHitChance(skillID uint32, chance int) {
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

func (r *CarnivalSystem) SkillHitChance(skillID uint32, fallback int) int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if v, ok := r.skillHitChance[skillID]; ok {
		return v
	}
	return fallback
}

func (r *CarnivalSystem) registerMaps(match *CarnivalMatch) {
	ids := []uint32{
		match.WaitingMapID,
		match.FieldMapID,
		match.ReviveMapID,
		match.WinMapID,
		match.LoseMapID,
	}
	for _, id := range ids {
		r.mapMatches[id] = match
	}
}

func (r *CarnivalSystem) BindGroup(group *StateMachineGroup) {
	r.mu.Lock()
	r.group = group
	r.mu.Unlock()
}

func (r *CarnivalSystem) Group() *StateMachineGroup {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.group
}

func (r *CarnivalSystem) Register(slotIndex int, waitingMapID, fieldMapID, reviveMapID, winMapID, loseMapID uint32, maxMembers int) *CarnivalMatch {
	match := NewCarnivalMatch(slotIndex, waitingMapID, fieldMapID, reviveMapID, winMapID, loseMapID, maxMembers)
	match.registry = r
	r.mu.Lock()
	r.matches[slotIndex] = match
	r.registerMaps(match)
	r.mu.Unlock()
	return match
}

func (r *CarnivalSystem) Match(slotIndex int) *CarnivalMatch {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.matches[slotIndex]
}

func (r *CarnivalSystem) Matches() []*CarnivalMatch {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]*CarnivalMatch, 0, len(r.matches))
	for _, m := range r.matches {
		out = append(out, m)
	}
	return out
}

func (r *CarnivalSystem) MapMatch(mapID uint32) *CarnivalMatch {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.mapMatches[mapID]
}

func (r *CarnivalSystem) Enter(slotIndex int, leader *Character) bool {
	match := r.Match(slotIndex)
	if match == nil {
		return false
	}
	members, party := leader.PartyOnMap()
	if party == nil || len(members) > match.MaxMembers {
		return false
	}

	roster := make([]*CarnivalMember, len(members))
	personal := make(map[uint32]*CarnivalPersonalCP, len(members))
	for i, ch := range members {
		roster[i] = ch.CarnivalMember()
		personal[ch.GetID()] = &CarnivalPersonalCP{}
	}
	red := &CarnivalTeam{
		Match:    match,
		TeamID:   constant.CarnivalTeamRed,
		LeaderID: leader.GetID(),
		PartyID:  party.GetPartyId(),
		Roster:   roster,
		Personal: personal,
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
		ch.carnivalTeam = red
	}
	return true
}

func (r *CarnivalSystem) Challenge(slotIndex int, leader *Character) bool {
	match := r.Match(slotIndex)
	if match == nil {
		return false
	}
	members, party := leader.PartyOnMap()
	if party == nil {
		return false
	}

	roster := make([]*CarnivalMember, len(members))
	for i, ch := range members {
		roster[i] = ch.CarnivalMember()
	}
	challenge := &CarnivalChallenge{
		LeaderID:    leader.GetID(),
		PartyID:     party.GetPartyId(),
		Map:         leader.GetMap(),
		Members:     roster,
		RequestedAt: time.Now(),
	}

	match.mu.Lock()
	red := match.Teams[constant.CarnivalTeamRed]
	if match.State != CarnivalStateWaiting || red == nil {
		match.mu.Unlock()
		return false
	}
	if len(members) != len(red.Roster) {
		match.mu.Unlock()
		return false
	}
	wasEmpty := len(match.queue) == 0
	match.queue = append(match.queue, challenge)
	if wasEmpty {
		match.Pending = challenge
	}
	match.mu.Unlock()

	if wasEmpty {
		match.StateMachine().CallHook("on_challenge")
	}
	return true
}
