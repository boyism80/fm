package entity

import (
	"fmt"
	"math/rand"
	"slices"
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/services/game/constant"
	"github.com/boyism80/fm/services/game/wz"
)

type CarnivalSummonResult int

const (
	CarnivalSummonSuccess CarnivalSummonResult = 0
	CarnivalSummonLackCP  CarnivalSummonResult = 1
	CarnivalSummonNoSlot  CarnivalSummonResult = 2
	CarnivalSummonFailed  CarnivalSummonResult = 3
)

const (
	carnivalGuardianSpawnState     byte = 1
	carnivalGuardianDestroyedState byte = 5
)

var carnivalGuardianBuffs = map[uint32]constant.MobBuffFlag{
	140: constant.MobBuffWeaponImmunity,
	141: constant.MobBuffMagicImmunity,
	150: constant.MobBuffWeaponAttackUp,
	151: constant.MobBuffMagicAttackUp,
	152: constant.MobBuffWeaponDefenseUp,
	153: constant.MobBuffMagicDefenseUp,
	154: constant.MobBuffAcc,
	155: constant.MobBuffAvoid,
	156: constant.MobBuffSpeed,
}

type CarnivalPersonalCP struct {
	AvailableCP int
	TotalCP     int
}

type CarnivalTeam struct {
	Match       *CarnivalMatch
	TeamID      constant.CarnivalTeam
	LeaderID    uint32
	PartyID     uint32
	Roster      []*CarnivalMember
	AvailableCP int
	TotalCP     int
	Personal    map[uint32]*CarnivalPersonalCP
	Winner      bool
	guardians   map[uint32]*Reactor
}

func (t *CarnivalTeam) HasMember(characterID uint32) bool {
	for _, member := range t.Roster {
		if member.ID == characterID {
			return true
		}
	}
	return false
}

func (t *CarnivalTeam) Members() []*Character {
	if t.Match == nil {
		return nil
	}
	out := make([]*Character, 0, len(t.Roster))
	for _, ch := range t.Match.StateMachine().Players() {
		if t.HasMember(ch.GetID()) {
			out = append(out, ch)
		}
	}
	return out
}

func (t *CarnivalTeam) Leader() *Character {
	for _, ch := range t.Members() {
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
	if t.Match == nil || t.Match.GetState() != CarnivalStateBattle {
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

func (t *CarnivalTeam) payableCP(ch *Character, amount int) *CarnivalPersonalCP {
	personal := t.PersonalCP(ch.GetID())
	if personal == nil || personal.AvailableCP < amount || t.AvailableCP < amount {
		return nil
	}
	return personal
}

func (t *CarnivalTeam) payCP(ch *Character, personal *CarnivalPersonalCP, amount int) {
	personal.AvailableCP -= amount
	t.AvailableCP -= amount
	ch.Listener.OnCarnivalObtainedCP(ch, personal.AvailableCP, personal.TotalCP)
	ch.Listener.OnCarnivalPartyCP(ch, t.TeamID, t.AvailableCP, t.TotalCP)
	t.broadcastPartyCP(ch)
}

func (t *CarnivalTeam) SummonMob(ch *Character, field *Map, num int) CarnivalSummonResult {
	carnival := field.Wz.Carnival
	if carnival == nil || num < 0 || num >= len(carnival.Mobs) {
		return CarnivalSummonFailed
	}
	entry := carnival.Mobs[num]
	personal := t.payableCP(ch, entry.SpendCP)
	if personal == nil {
		return CarnivalSummonLackCP
	}

	gen, ok := t.freeMobGenPos(field)
	if !ok {
		return CarnivalSummonNoSlot
	}
	if field.SummonMob(entry.ID, gen.Pos, t.TeamID) != nil {
		return CarnivalSummonNoSlot
	}

	t.payCP(ch, personal, entry.SpendCP)
	return CarnivalSummonSuccess
}

func (t *CarnivalTeam) freeMobGenPos(field *Map) (wz.CarnivalGenPos, bool) {
	for _, gen := range field.Wz.Carnival.MobGenPos {
		if gen.Team != int(t.TeamID) && gen.Team != int(constant.CarnivalTeamNone) {
			continue
		}
		used := slices.ContainsFunc(field.SummonedMobSpawns, func(spawn *MobSpawn) bool {
			pos := spawn.Wz.Position
			return pos.X == gen.Pos.X && pos.Y == gen.Pos.Y && (gen.Team == int(constant.CarnivalTeamNone) || spawn.Wz.Team == gen.Team)
		})
		if !used {
			return gen, true
		}
	}
	return wz.CarnivalGenPos{}, false
}

func (t *CarnivalTeam) UseSkill(ch *Character, field *Map, num int) CarnivalSummonResult {
	carnival := field.Wz.Carnival
	if carnival == nil || t.Match == nil || num < 0 || num >= len(carnival.Skills) {
		return CarnivalSummonFailed
	}
	skillID := carnival.Skills[num]
	skillDef := field.GameWorld.GetResources().GetCarnivalSkill(skillID)
	if skillDef == nil {
		return CarnivalSummonFailed
	}
	personal := t.payableCP(ch, skillDef.SpendCP)
	if personal == nil {
		return CarnivalSummonLackCP
	}

	enemy := t.Match.enemyTeam(t.TeamID)
	if enemy == nil || !enemy.Debuff(field, skillID) {
		return CarnivalSummonFailed
	}

	t.payCP(ch, personal, skillDef.SpendCP)
	return CarnivalSummonSuccess
}

func (t *CarnivalTeam) SummonGuardian(ch *Character, field *Map, num uint32) CarnivalSummonResult {
	carnival := field.Wz.Carnival
	if carnival == nil {
		return CarnivalSummonFailed
	}
	guardian := field.GameWorld.GetResources().GetCarnivalGuardian(num)
	if guardian == nil {
		return CarnivalSummonFailed
	}
	personal := t.payableCP(ch, guardian.SpendCP)
	if personal == nil {
		return CarnivalSummonLackCP
	}
	if t.isGuardianAlive(t.guardians[num]) {
		return CarnivalSummonNoSlot
	}

	gen, ok := t.freeGuardianGenPos(field)
	if !ok {
		return CarnivalSummonNoSlot
	}
	reactorID := carnival.ReactorBlue
	if t.TeamID == constant.CarnivalTeamRed {
		reactorID = carnival.ReactorRed
	}
	reactor, err := field.SpawnReactorTemplate(reactorID, gen.Pos, fmt.Sprintf("%d%d", t.TeamID, num), carnivalGuardianSpawnState)
	if err != nil {
		return CarnivalSummonNoSlot
	}
	if t.guardians == nil {
		t.guardians = make(map[uint32]*Reactor)
	}
	t.guardians[num] = reactor

	for _, obj := range field.GetMobs() {
		mob, ok := obj.(*Mob)
		if ok && mob.CarnivalTeam == t.TeamID {
			t.buffGuardian(mob, guardian)
		}
	}

	t.payCP(ch, personal, guardian.SpendCP)
	return CarnivalSummonSuccess
}

func (t *CarnivalTeam) isGuardianAlive(reactor *Reactor) bool {
	if reactor == nil || reactor.Map == nil {
		return false
	}
	return reactor.Map.GetReactor(reactor.OID) == reactor && reactor.State < carnivalGuardianDestroyedState
}

func (t *CarnivalTeam) freeGuardianGenPos(field *Map) (wz.CarnivalGenPos, bool) {
	for _, gen := range field.Wz.Carnival.GuardianGenPos {
		if gen.Team != int(t.TeamID) && gen.Team != int(constant.CarnivalTeamNone) {
			continue
		}
		used := false
		for _, obj := range field.GetReactors() {
			reactor, ok := obj.(*Reactor)
			if ok && reactor.Position.X == gen.Pos.X && reactor.Position.Y == gen.Pos.Y && reactor.State < carnivalGuardianDestroyedState {
				used = true
				break
			}
		}
		if !used {
			return gen, true
		}
	}
	return wz.CarnivalGenPos{}, false
}

func (t *CarnivalTeam) buffGuardian(mob *Mob, guardian *wz.CarnivalGuardian) {
	flag, ok := carnivalGuardianBuffs[guardian.MobSkillID]
	if !ok {
		return
	}
	levelData := mob.GameWorld.GetResources().GetMobSkill(guardian.MobSkillID, guardian.Level)
	if levelData == nil {
		return
	}
	mob.Buffs.Add(time.Duration(levelData.DurationMs)*time.Millisecond, levelData.SkillWz, guardian.Level, 0,
		map[constant.MobBuffFlag]int32{flag: int32(levelData.X)}, nil)
}

func (t *CarnivalTeam) buffGuardians(mob *Mob) {
	for num, reactor := range t.guardians {
		if !t.isGuardianAlive(reactor) {
			continue
		}
		guardian := mob.GameWorld.GetResources().GetCarnivalGuardian(num)
		if guardian != nil {
			t.buffGuardian(mob, guardian)
		}
	}
}

func (t *CarnivalTeam) destroyGuardian(reactor *Reactor) bool {
	for num, guardianReactor := range t.guardians {
		if guardianReactor != reactor {
			continue
		}
		delete(t.guardians, num)
		guardian := reactor.GameWorld.GetResources().GetCarnivalGuardian(num)
		if guardian == nil || reactor.Map == nil {
			return true
		}
		for _, obj := range reactor.Map.GetMobs() {
			mob, ok := obj.(*Mob)
			if ok && mob.CarnivalTeam == t.TeamID {
				mob.Buffs.Dispel(guardian.MobSkillID)
			}
		}
		return true
	}
	return false
}

func (t *CarnivalTeam) ShowResult() {
	for _, ch := range t.Members() {
		ch.Listener.OnCarnivalResult(ch, t.Winner)
	}
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
	if t.Match == nil || t.Match.GetState() != CarnivalStateBattle {
		return false
	}
	skillDef := field.GameWorld.GetResources().GetCarnivalSkill(skillID)
	if skillDef == nil {
		return false
	}

	targets := make([]*Character, 0, len(t.Roster))
	for _, member := range t.Members() {
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
	if !skillDef.TargetsAll {
		targets = targets[:1]
	}

	chance := t.Match.registry.SkillHitChance(skillID, skillDef.HitChance)
	for _, target := range targets {
		if skillDef.TargetsAll && rand.Intn(100) >= chance {
			continue
		}
		if !target.GiveMobSkillDebuff(skillDef.MobSkillID, skillDef.Level) {
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
	dest := gw.GetMapSystem().Find(t.Match.StateMachine(), mapID)
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
	for _, ch := range t.Members() {
		_ = ch.Warp(ctx, dest, spawn)
	}
	return true
}

func (t *CarnivalTeam) RemoveMember(ch *Character) {
	t.Roster = slices.DeleteFunc(t.Roster, func(member *CarnivalMember) bool {
		return member.ID == ch.GetID()
	})
	ch.UnbindCarnival()
}

func (t *CarnivalTeam) Clear() {
	for _, ch := range t.Members() {
		ch.UnbindCarnival()
	}
	t.Roster = nil
	t.LeaderID = 0
	t.AvailableCP = 0
	t.TotalCP = 0
	t.Personal = nil
	t.Winner = false
}
