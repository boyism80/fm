package game

import (
	"fmt"
	"math/rand"
	"slices"
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/maple/constant"
	"github.com/boyism80/fm/maple/wz"
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
	availableCP int
	totalCP     int
}

type CarnivalTeam struct {
	match       *CarnivalMatch
	teamID      constant.CarnivalTeam
	leaderID    uint32
	partyID     uint32
	roster      []*CarnivalMember
	availableCP int
	totalCP     int
	personal    map[uint32]*CarnivalPersonalCP
	winner      bool
	guardians   map[uint32]*Reactor
}

func (t *CarnivalTeam) HasMember(characterID uint32) bool {
	for _, member := range t.roster {
		if member.id == characterID {
			return true
		}
	}
	return false
}

func (t *CarnivalTeam) Members() []*Character {
	if t.match == nil {
		return nil
	}
	out := make([]*Character, 0, len(t.roster))
	for _, ch := range t.match.StateMachine().Players() {
		if t.HasMember(ch.GetID()) {
			out = append(out, ch)
		}
	}
	return out
}

func (t *CarnivalTeam) Leader() *Character {
	for _, ch := range t.Members() {
		if ch.GetID() == t.leaderID {
			return ch
		}
	}
	return nil
}

func (t *CarnivalTeam) PersonalCP(characterID uint32) *CarnivalPersonalCP {
	if t.personal == nil {
		return nil
	}
	return t.personal[characterID]
}

func (t *CarnivalTeam) AddCP(ch *Character, amount int) bool {
	if amount <= 0 {
		return false
	}
	if t.match == nil || t.match.GetState() != CarnivalStateBattle {
		return false
	}
	personal := t.PersonalCP(ch.GetID())
	if personal == nil {
		return false
	}
	personal.availableCP += amount
	personal.totalCP += amount
	t.availableCP += amount
	t.totalCP += amount
	ch.listener.OnCarnivalObtainedCP(ch, personal.availableCP, personal.totalCP)
	ch.listener.OnCarnivalPartyCP(ch, t.teamID, t.availableCP, t.totalCP)
	t.broadcastPartyCP(ch)
	return true
}

func (t *CarnivalTeam) payableCP(ch *Character, amount int) *CarnivalPersonalCP {
	personal := t.PersonalCP(ch.GetID())
	if personal == nil || personal.availableCP < amount || t.availableCP < amount {
		return nil
	}
	return personal
}

func (t *CarnivalTeam) payCP(ch *Character, personal *CarnivalPersonalCP, amount int) {
	personal.availableCP -= amount
	t.availableCP -= amount
	ch.listener.OnCarnivalObtainedCP(ch, personal.availableCP, personal.totalCP)
	ch.listener.OnCarnivalPartyCP(ch, t.teamID, t.availableCP, t.totalCP)
	t.broadcastPartyCP(ch)
}

func (t *CarnivalTeam) SummonMob(ch *Character, field *Map, num int) CarnivalSummonResult {
	carnival := field.wz.Carnival
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
	if field.SummonMob(entry.ID, gen.Pos, t.teamID) != nil {
		return CarnivalSummonNoSlot
	}

	t.payCP(ch, personal, entry.SpendCP)
	return CarnivalSummonSuccess
}

func (t *CarnivalTeam) freeMobGenPos(field *Map) (wz.CarnivalGenPos, bool) {
	for _, gen := range field.wz.Carnival.MobGenPos {
		if gen.Team != int(t.teamID) && gen.Team != int(constant.CarnivalTeamNone) {
			continue
		}
		used := slices.ContainsFunc(field.summonedMobSpawns, func(spawn *MobSpawn) bool {
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
	carnival := field.wz.Carnival
	if carnival == nil || t.match == nil || num < 0 || num >= len(carnival.Skills) {
		return CarnivalSummonFailed
	}
	skillID := carnival.Skills[num]
	skillDef := field.gameWorld.GetResources().GetCarnivalSkill(skillID)
	if skillDef == nil {
		return CarnivalSummonFailed
	}
	personal := t.payableCP(ch, skillDef.SpendCP)
	if personal == nil {
		return CarnivalSummonLackCP
	}

	enemy := t.match.enemyTeam(t.teamID)
	if enemy == nil || !enemy.Debuff(field, skillID) {
		return CarnivalSummonFailed
	}

	t.payCP(ch, personal, skillDef.SpendCP)
	return CarnivalSummonSuccess
}

func (t *CarnivalTeam) SummonGuardian(ch *Character, field *Map, num uint32) CarnivalSummonResult {
	carnival := field.wz.Carnival
	if carnival == nil {
		return CarnivalSummonFailed
	}
	guardian := field.gameWorld.GetResources().GetCarnivalGuardian(num)
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
	if t.teamID == constant.CarnivalTeamRed {
		reactorID = carnival.ReactorRed
	}
	reactor, err := field.SpawnReactorTemplate(reactorID, gen.Pos, fmt.Sprintf("%d%d", t.teamID, num), carnivalGuardianSpawnState)
	if err != nil {
		return CarnivalSummonNoSlot
	}
	if t.guardians == nil {
		t.guardians = make(map[uint32]*Reactor)
	}
	t.guardians[num] = reactor

	for _, obj := range field.GetMobs() {
		mob, ok := obj.(*Mob)
		if ok && mob.carnivalTeam == t.teamID {
			t.buffGuardian(mob, guardian)
		}
	}

	t.payCP(ch, personal, guardian.SpendCP)
	return CarnivalSummonSuccess
}

func (t *CarnivalTeam) isGuardianAlive(reactor *Reactor) bool {
	if reactor == nil || reactor.onMap == nil {
		return false
	}
	return reactor.onMap.GetReactor(reactor.oid) == reactor && reactor.state < carnivalGuardianDestroyedState
}

func (t *CarnivalTeam) freeGuardianGenPos(field *Map) (wz.CarnivalGenPos, bool) {
	for _, gen := range field.wz.Carnival.GuardianGenPos {
		if gen.Team != int(t.teamID) && gen.Team != int(constant.CarnivalTeamNone) {
			continue
		}
		used := false
		for _, obj := range field.GetReactors() {
			reactor, ok := obj.(*Reactor)
			if ok && reactor.position.X == gen.Pos.X && reactor.position.Y == gen.Pos.Y && reactor.state < carnivalGuardianDestroyedState {
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
	levelData := mob.gameWorld.GetResources().GetMobSkill(guardian.MobSkillID, guardian.Level)
	if levelData == nil {
		return
	}
	mob.buffs.Add(time.Duration(levelData.DurationMs)*time.Millisecond, levelData.SkillWz, guardian.Level, 0,
		map[constant.MobBuffFlag]int32{flag: int32(levelData.X)}, nil)
}

func (t *CarnivalTeam) buffGuardians(mob *Mob) {
	for num, reactor := range t.guardians {
		if !t.isGuardianAlive(reactor) {
			continue
		}
		guardian := mob.gameWorld.GetResources().GetCarnivalGuardian(num)
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
		guardian := reactor.gameWorld.GetResources().GetCarnivalGuardian(num)
		if guardian == nil || reactor.onMap == nil {
			return true
		}
		for _, obj := range reactor.onMap.GetMobs() {
			mob, ok := obj.(*Mob)
			if ok && mob.carnivalTeam == t.teamID {
				mob.buffs.Dispel(guardian.MobSkillID)
			}
		}
		return true
	}
	return false
}

func (t *CarnivalTeam) ShowResult() {
	for _, ch := range t.Members() {
		ch.listener.OnCarnivalResult(ch, t.winner)
	}
}

func (t *CarnivalTeam) broadcastPartyCP(origin *Character) {
	if t.match == nil {
		return
	}
	sm := t.match.StateMachine()
	if sm == nil {
		return
	}
	for _, p := range sm.Players() {
		if p == nil || p.GetID() == origin.GetID() {
			continue
		}
		p.listener.OnCarnivalPartyCP(p, t.teamID, t.availableCP, t.totalCP)
	}
}

// Only members standing on field can be hit; a skill that targets all still misses each one by chance.
func (t *CarnivalTeam) Debuff(field *Map, skillID uint32) bool {
	if t.match == nil || t.match.GetState() != CarnivalStateBattle {
		return false
	}
	skillDef := field.gameWorld.GetResources().GetCarnivalSkill(skillID)
	if skillDef == nil {
		return false
	}

	targets := make([]*Character, 0, len(t.roster))
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

	chance := t.match.registry.SkillHitChance(skillID, skillDef.HitChance)
	for _, target := range targets {
		if skillDef.TargetsAll && rand.Intn(100) >= chance {
			continue
		}
		if !target.debuffs.GiveMobSkill(skillDef.MobSkillID, skillDef.Level) {
			target.buffs.Dispel()
		}
	}
	return true
}

func (t *CarnivalTeam) Warp(ctx actor.Context, mapID uint32, portalName string) bool {
	if t.match == nil {
		return false
	}
	gw := t.match.gameWorld()
	if gw == nil {
		return false
	}
	dest := gw.Systems().Get[*MapSystem]().Find(t.match.StateMachine(), mapID)
	if dest == nil {
		return false
	}
	spawn := uint8(0)
	if portalName != "" && dest.wz != nil {
		for _, p := range dest.wz.Portals {
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
	t.roster = slices.DeleteFunc(t.roster, func(member *CarnivalMember) bool {
		return member.id == ch.GetID()
	})
	ch.carnival.team = nil
}

func (t *CarnivalTeam) Clear() {
	for _, ch := range t.Members() {
		ch.carnival.team = nil
	}
	t.roster = nil
	t.leaderID = 0
	t.availableCP = 0
	t.totalCP = 0
	t.personal = nil
	t.winner = false
}
