package entity

import (
	"math/rand"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/services/game/constant"
)

type CarnivalPersonalCP struct {
	AvailableCP int
	TotalCP     int
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

func (t *CarnivalTeam) HasMember(characterID uint32) bool {
	for _, id := range t.MemberIDs {
		if id == characterID {
			return true
		}
	}
	return false
}

func (t *CarnivalTeam) Members(gw GameWorld) []*Character {
	if gw == nil || t.Match == nil {
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
	if ms == nil {
		return out
	}
	for _, ch := range t.Match.findCharacters(gw, t.MemberIDs, t.PartyID, []uint32{
		t.Match.WaitingMapID,
		t.Match.FieldMapID,
		t.Match.ReviveMapID,
		t.Match.WinMapID,
		t.Match.LoseMapID,
	}) {
		add(ch)
	}
	return out
}

func (t *CarnivalTeam) Leader(gw GameWorld) *Character {
	if gw == nil {
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
	if t.Personal == nil {
		return nil
	}
	return t.Personal[characterID]
}

func (t *CarnivalTeam) AddCP(ch *Character, amount int) bool {
	if amount <= 0 {
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
	if amount <= 0 {
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
	if t.Match == nil {
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

// Only members standing on field can be hit; a skill that targets all still misses each one by chance.
func (t *CarnivalTeam) Debuff(field *Map, skillID uint32) bool {
	if t.Match == nil {
		return false
	}
	skillDef := field.GameWorld.GetResources().GetCarnivalSkill(skillID)
	if skillDef == nil {
		return false
	}

	targets := make([]*Character, 0, len(t.MemberIDs))
	for _, member := range t.Members(field.GameWorld) {
		if member.GetMap() == field {
			targets = append(targets, member)
		}
	}
	if len(targets) == 0 {
		return false
	}
	rand.Shuffle(len(targets), func(i, j int) {
		targets[i], targets[j] = targets[j], targets[i]
	})
	if skillDef.TargetsAll == false {
		targets = targets[:1]
	}

	chance := t.Match.registry.SkillHitChance(skillID, skillDef.HitChance)
	for _, target := range targets {
		if skillDef.TargetsAll && rand.Intn(100) >= chance {
			continue
		}
		if target.GiveMobSkillDebuff(skillDef.MobSkillID, skillDef.Level) == false {
			target.Buffs.Dispel()
		}
	}
	return true
}

func (t *CarnivalTeam) Warp(ctx actor.Context, mapID uint32, portalName string) bool {
	if t.Match == nil {
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

func (t *CarnivalTeam) RemoveMember(ch *Character) {
	for i, id := range t.MemberIDs {
		if id == ch.GetID() {
			t.MemberIDs = append(t.MemberIDs[:i], t.MemberIDs[i+1:]...)
			break
		}
	}
	ch.UnbindCarnival()
}

func (t *CarnivalTeam) Clear(gw GameWorld) {
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
