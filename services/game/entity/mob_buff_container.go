package entity

import (
	"fmt"
	"log"
	"sort"
	"time"

	"github.com/boyism80/fm/core/luax"
	"github.com/boyism80/fm/services/game/constant"
	"github.com/boyism80/fm/services/game/wz"
)

type mobBuffForPacket struct {
	Status  constant.MobBuffFlag
	Value   int32
	SkillID uint32
}

type MobSkillBuff struct {
	StartTime  time.Time
	Duration   time.Duration
	Values     map[constant.MobBuffFlag]int32
	Stacks     map[constant.MobBuffFlag]uint8
	Wz         *wz.Skill
	SkillLevel uint8
	Causer     uint32
}

func (e *MobSkillBuff) CallOnMobBuffScript(mob *Mob) {
	e.callMobSkillHook(mob, "on_mob_buff")
}

func (e *MobSkillBuff) CallOnMobUnbuffScript(mob *Mob) {
	e.callMobSkillHook(mob, "on_mob_unbuff")
}

func (e *MobSkillBuff) callMobSkillHook(mob *Mob, hookPrefix string) {
	if e == nil || mob == nil {
		return
	}
	mapInstance := mob.GetMap()
	if mapInstance == nil {
		return
	}
	root := mapInstance.GetLuaRoot()
	if root == nil {
		return
	}
	skillID := e.Wz.ID
	var causer *Character
	if e.Causer != 0 {
		causer = mapInstance.GetPlayer(e.Causer)
	}
	if causer == nil {
		return
	}
	scriptPath := fmt.Sprintf("script/skill/%d.lua", skillID)
	hookName := luax.SkillScriptHookName(hookPrefix, skillID)
	thread, err := luax.NewThread(root, scriptPath)
	if err != nil {
		log.Printf("mob skill hook %s %d: %v", hookPrefix, skillID, err)
		return
	}
	if _, err := luax.Call(thread, hookName, mob, e, causer); err != nil {
		log.Printf("mob skill hook %s %d: %v", hookPrefix, skillID, err)
	}
}

type MobBuffContainer struct {
	owner    *Mob
	byFlag   map[constant.MobBuffFlag]*MobSkillBuff
	entities map[*MobSkillBuff]struct{}
}

func NewMobBuffContainer(owner *Mob) *MobBuffContainer {
	if owner == nil {
		panic("MobBuffContainer owner is nil")
	}
	return &MobBuffContainer{
		owner:    owner,
		byFlag:   make(map[constant.MobBuffFlag]*MobSkillBuff),
		entities: make(map[*MobSkillBuff]struct{}),
	}
}

func (bc *MobBuffContainer) add(entity *MobSkillBuff) (removed []*MobSkillBuff) {
	if bc == nil || entity == nil || len(entity.Values) == 0 {
		return nil
	}
	conflicts := make(map[*MobSkillBuff]struct{})
	for flag := range entity.Values {
		if existing := bc.byFlag[flag]; existing != nil && existing != entity {
			conflicts[existing] = struct{}{}
		}
	}
	for existing := range conflicts {
		removed = append(removed, existing)
		bc.removeEntity(existing)
	}
	bc.addEntity(entity)
	return removed
}

func (bc *MobBuffContainer) addEntity(entity *MobSkillBuff) {
	bc.entities[entity] = struct{}{}
	for flag := range entity.Values {
		bc.byFlag[flag] = entity
	}
}

func (bc *MobBuffContainer) removeEntity(entity *MobSkillBuff) {
	if bc == nil || entity == nil {
		return
	}
	delete(bc.entities, entity)
	for flag := range entity.Values {
		if current := bc.byFlag[flag]; current == entity {
			delete(bc.byFlag, flag)
		}
	}
}

func mobSkillBuffSourceEqual(ent *MobSkillBuff, skillWz *wz.Skill, skillLevel uint8, causer uint32) bool {
	if ent == nil {
		return false
	}
	var sid uint32
	if skillWz != nil {
		sid = skillWz.ID
	}
	var entSid uint32
	if ent.Wz != nil {
		entSid = ent.Wz.ID
	}
	return sid == entSid && ent.SkillLevel == skillLevel && ent.Causer == causer
}

func mobSkillBuffMapsEqualForRefresh(ent *MobSkillBuff, values map[constant.MobBuffFlag]int32, stacks map[constant.MobBuffFlag]uint8) bool {
	if ent == nil || len(ent.Values) != len(values) {
		return false
	}
	for f, v := range values {
		if ev, ok := ent.Values[f]; !ok || ev != v {
			return false
		}
		wantStack := uint8(1)
		if stacks != nil {
			if sv, ok := stacks[f]; ok && sv >= 1 {
				wantStack = sv
			}
		}
		got := ent.Stacks[f]
		if got < 1 {
			got = 1
		}
		if got != wantStack {
			return false
		}
	}
	return true
}

func (bc *MobBuffContainer) tryRefreshDurationOnly(now time.Time, durationMs int64, skillWz *wz.Skill, skillLevel uint8, causer uint32, values map[constant.MobBuffFlag]int32, stacks map[constant.MobBuffFlag]uint8) bool {
	if bc == nil || len(values) == 0 {
		return false
	}
	var sole *MobSkillBuff
	for flag := range values {
		ent := bc.byFlag[flag]
		if ent == nil {
			return false
		}
		if sole == nil {
			sole = ent
		} else if sole != ent {
			return false
		}
	}
	if sole == nil || !mobSkillBuffMapsEqualForRefresh(sole, values, stacks) || !mobSkillBuffSourceEqual(sole, skillWz, skillLevel, causer) {
		return false
	}
	var dur time.Duration
	if durationMs > 0 {
		dur = time.Duration(durationMs) * time.Millisecond
	}
	sole.StartTime = now
	sole.Duration = dur
	return true
}

func (bc *MobBuffContainer) AddSkillBuff(now time.Time, durationMs int64, skillWz *wz.Skill, skillLevel uint8, causer uint32, values map[constant.MobBuffFlag]int32, stacks map[constant.MobBuffFlag]uint8) {
	if bc == nil || len(values) == 0 {
		return
	}
	if bc.tryRefreshDurationOnly(now, durationMs, skillWz, skillLevel, causer, values, stacks) {
		return
	}
	valCopy := make(map[constant.MobBuffFlag]int32, len(values))
	stackCopy := make(map[constant.MobBuffFlag]uint8, len(values))
	for f, v := range values {
		valCopy[f] = v
		s := uint8(1)
		if stacks != nil {
			if sv, ok := stacks[f]; ok && sv >= 1 {
				s = sv
			}
		}
		stackCopy[f] = s
	}
	var dur time.Duration
	if durationMs > 0 {
		dur = time.Duration(durationMs) * time.Millisecond
	}
	ent := &MobSkillBuff{
		StartTime:  now,
		Duration:   dur,
		Values:     valCopy,
		Stacks:     stackCopy,
		Wz:         skillWz,
		SkillLevel: skillLevel,
		Causer:     causer,
	}
	removed := bc.add(ent)
	mob := bc.owner
	for _, old := range removed {
		mob.dispatchRemovedMobSkillBuff(old)
	}
	ent.CallOnMobBuffScript(mob)
	mob.dispatchMobSkillBuffAppliedPackets(ent)
}

func (m *Mob) dispatchRemovedMobSkillBuff(old *MobSkillBuff) {
	if old == nil {
		return
	}
	mapInstance := m.GetMap()
	if mapInstance != nil && mapInstance.listener != nil {
		for flag := range old.Values {
			mapInstance.listener.OnMobMobBuffCancelled(mapInstance, m, flag)
		}
	}
	old.CallOnMobUnbuffScript(m)
}

func (m *Mob) dispatchMobSkillBuffAppliedPackets(ent *MobSkillBuff) {
	if ent == nil {
		return
	}
	mapInstance := m.GetMap()
	if mapInstance == nil || mapInstance.listener == nil {
		return
	}
	durationMs := int64(ent.Duration / time.Millisecond)
	if durationMs < 0 {
		durationMs = 0
	}
	skillID := uint32(0)
	if ent.Wz != nil {
		skillID = ent.Wz.ID
	}
	for _, flag := range sortedMobBuffFlags(ent.Values) {
		mapInstance.listener.OnMobMobBuffApplied(mapInstance, m, flag, ent.Values[flag], skillID, durationMs)
	}
}

func sortedMobBuffFlags(values map[constant.MobBuffFlag]int32) []constant.MobBuffFlag {
	flags := make([]constant.MobBuffFlag, 0, len(values))
	for f := range values {
		flags = append(flags, f)
	}
	sort.Slice(flags, func(i, j int) bool { return flags[i] < flags[j] })
	return flags
}

func (bc *MobBuffContainer) RemoveBuffForFlag(flag constant.MobBuffFlag) {
	if bc == nil {
		return
	}
	ent := bc.byFlag[flag]
	if ent == nil {
		return
	}
	bc.removeEntity(ent)
	bc.owner.dispatchRemovedMobSkillBuff(ent)
}

func (bc *MobBuffContainer) collectExpired(now time.Time) []*MobSkillBuff {
	if bc == nil || len(bc.entities) == 0 {
		return nil
	}
	var out []*MobSkillBuff
	for ent := range bc.entities {
		if ent == nil {
			continue
		}
		if ent.Duration <= 0 {
			continue
		}
		if now.After(ent.StartTime.Add(ent.Duration)) {
			out = append(out, ent)
		}
	}
	return out
}

func (bc *MobBuffContainer) removeExpiredEntities(now time.Time) {
	expired := bc.collectExpired(now)
	for _, ent := range expired {
		bc.removeEntity(ent)
		bc.owner.dispatchRemovedMobSkillBuff(ent)
	}
}

func (bc *MobBuffContainer) flattenForSpawnPacket() (mask uint32, entries []mobBuffForPacket) {
	if bc == nil || len(bc.entities) == 0 {
		return 0, nil
	}
	type row struct {
		flag    constant.MobBuffFlag
		value   int32
		skillID uint32
	}
	var rows []row
	for ent := range bc.entities {
		if ent == nil {
			continue
		}
		skillID := uint32(0)
		if ent.Wz != nil {
			skillID = ent.Wz.ID
		}
		for f, v := range ent.Values {
			mask |= uint32(f)
			rows = append(rows, row{f, v, skillID})
		}
	}
	sort.Slice(rows, func(i, j int) bool { return rows[i].flag < rows[j].flag })
	entries = make([]mobBuffForPacket, 0, len(rows))
	for _, r := range rows {
		entries = append(entries, mobBuffForPacket{Status: r.flag, Value: r.value, SkillID: r.skillID})
	}
	return mask, entries
}

func (bc *MobBuffContainer) getValue(flag constant.MobBuffFlag) int32 {
	if bc == nil {
		return 0
	}
	ent := bc.byFlag[flag]
	if ent == nil {
		return 0
	}
	return ent.Values[flag]
}

func (bc *MobBuffContainer) getStack(flag constant.MobBuffFlag) uint8 {
	if bc == nil {
		return 0
	}
	ent := bc.byFlag[flag]
	if ent == nil {
		return 0
	}
	s := ent.Stacks[flag]
	if s < 1 {
		return 1
	}
	return s
}

func (bc *MobBuffContainer) setStack(flag constant.MobBuffFlag, stack uint8) bool {
	if bc == nil {
		return false
	}
	ent := bc.byFlag[flag]
	if ent == nil {
		return false
	}
	if stack < 1 {
		stack = 1
	}
	ent.Stacks[flag] = stack
	return true
}

func (bc *MobBuffContainer) hasFlag(flag constant.MobBuffFlag) bool {
	if bc == nil {
		return false
	}
	_, ok := bc.byFlag[flag]
	return ok
}

func (bc *MobBuffContainer) causerForFlag(flag constant.MobBuffFlag) uint32 {
	if bc == nil {
		return 0
	}
	ent := bc.byFlag[flag]
	if ent == nil {
		return 0
	}
	return ent.Causer
}
