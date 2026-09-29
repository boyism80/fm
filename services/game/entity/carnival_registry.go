package entity

import (
	"strconv"
	"sync"
	"time"

	"github.com/boyism80/fm/services/game/constant"
)

type CarnivalRegistry struct {
	mu             sync.RWMutex
	matches        map[int]*CarnivalMatch
	mapMatches     map[uint32]*CarnivalMatch
	group          *StateMachineGroup
	skillHitChance map[uint32]int
}

func NewCarnivalRegistry() *CarnivalRegistry {
	return &CarnivalRegistry{
		matches:    make(map[int]*CarnivalMatch),
		mapMatches: make(map[uint32]*CarnivalMatch),
	}
}

func (r *CarnivalRegistry) SetSkillHitChance(skillID uint32, chance int) {
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
		r.mapMatches[id] = match
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
	match := NewCarnivalMatch(slotIndex, waitingMapID, fieldMapID, reviveMapID, winMapID, loseMapID, maxMembers)
	match.registry = r
	r.mu.Lock()
	r.matches[slotIndex] = match
	r.registerMaps(match)
	r.mu.Unlock()
	return match
}

func (r *CarnivalRegistry) Match(slotIndex int) *CarnivalMatch {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.matches[slotIndex]
}

func (r *CarnivalRegistry) Matches() []*CarnivalMatch {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]*CarnivalMatch, 0, len(r.matches))
	for _, m := range r.matches {
		out = append(out, m)
	}
	return out
}

func (r *CarnivalRegistry) MapMatch(mapID uint32) *CarnivalMatch {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.mapMatches[mapID]
}

func (r *CarnivalRegistry) Enter(slotIndex int, leader *Character) bool {
	match := r.Match(slotIndex)
	if match == nil {
		return false
	}
	members, party := leader.PartyOnMap()
	if party == nil || len(members) > match.MaxMembers {
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
	match := r.Match(slotIndex)
	if match == nil {
		return false, false
	}
	members, party := leader.PartyOnMap()
	if party == nil {
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
	redLeader := red.Leader(leader.GameWorld)
	if redLeader == nil || redLeader.GetDialog() != nil {
		return true, false
	}
	return true, true
}
