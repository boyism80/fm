package entity

import (
	"fmt"
	"log"
	"time"

	"github.com/boyism80/fm/core/luax"
	"github.com/boyism80/fm/services/game/constant"
	"github.com/boyism80/fm/services/game/wz"
)

type MobBuff struct {
	StartTime  time.Time
	Duration   time.Duration
	Values     map[constant.MobBuffFlag]int32
	Stacks     map[constant.MobBuffFlag]uint8
	Wz         *wz.Skill
	SkillLevel uint8
	Causer     uint32
}

func (e *MobBuff) RemainingDuration(now time.Time) time.Duration {
	if e == nil || e.Duration <= 0 {
		return 0
	}
	end := e.StartTime.Add(e.Duration)
	if !now.Before(end) {
		return 0
	}
	return end.Sub(now)
}

func (e *MobBuff) SameSource(skillWz *wz.Skill, skillLevel uint8, causer uint32) bool {
	if e == nil {
		return false
	}
	var sid uint32
	if skillWz != nil {
		sid = skillWz.ID
	}
	var entSid uint32
	if e.Wz != nil {
		entSid = e.Wz.ID
	}
	return sid == entSid && e.SkillLevel == skillLevel && e.Causer == causer
}

func (e *MobBuff) SameValuesAndStacks(values map[constant.MobBuffFlag]int32, stacks map[constant.MobBuffFlag]uint8) bool {
	if e == nil || len(e.Values) != len(values) {
		return false
	}
	for f, v := range values {
		if ev, ok := e.Values[f]; !ok || ev != v {
			return false
		}
		wantStack := uint8(1)
		if stacks != nil {
			if sv, ok := stacks[f]; ok && sv >= 1 {
				wantStack = sv
			}
		}
		got := e.Stacks[f]
		if got < 1 {
			got = 1
		}
		if got != wantStack {
			return false
		}
	}
	return true
}

func (e *MobBuff) callMobSkillHook(mob *Mob, hookPrefix string) {
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
	scriptPath := fmt.Sprintf("script/mob/skill/%d.lua", skillID)
	hookName := fmt.Sprintf("%s_%d", hookPrefix, skillID)
	thread, err := luax.NewThread(root, scriptPath)
	if err != nil {
		log.Printf("mob skill hook %s %d: %v", hookPrefix, skillID, err)
		return
	}
	luax.CallAsync(root, thread, hookName, mob, e, causer).OnError(func(err error) {
		log.Printf("mob skill hook %s %d: %v", hookPrefix, skillID, err)
	})
}

type MobBuffContainer struct {
	owner       *Mob
	byFlag      map[constant.MobBuffFlag]*MobBuff
	entities    map[*MobBuff]struct{}
	reflections []int32
}

func NewMobBuffContainer(owner *Mob) *MobBuffContainer {
	if owner == nil {
		panic("MobBuffContainer owner is nil")
	}
	if owner.Listener == nil {
		panic("MobBuffContainer: mob listener must not be nil")
	}
	return &MobBuffContainer{
		owner:    owner,
		byFlag:   make(map[constant.MobBuffFlag]*MobBuff),
		entities: make(map[*MobBuff]struct{}),
	}
}

func (bc *MobBuffContainer) refreshDuration(now time.Time, duration time.Duration, skillWz *wz.Skill, skillLevel uint8, causer uint32, values map[constant.MobBuffFlag]int32, stacks map[constant.MobBuffFlag]uint8) bool {
	if bc == nil || len(values) == 0 {
		return false
	}
	var ent *MobBuff
	for flag := range values {
		current := bc.byFlag[flag]
		if current == nil {
			return false
		}
		if ent == nil {
			ent = current
		} else if ent != current {
			return false
		}
	}
	if !ent.SameValuesAndStacks(values, stacks) || !ent.SameSource(skillWz, skillLevel, causer) {
		return false
	}
	ent.StartTime = now
	ent.Duration = duration
	return true
}

func (bc *MobBuffContainer) Add(duration time.Duration, skillWz *wz.Skill, skillLevel uint8, causer uint32, values map[constant.MobBuffFlag]int32, stacks map[constant.MobBuffFlag]uint8) {
	if bc == nil || len(values) == 0 {
		return
	}
	if mob := bc.owner; mob != nil && mob.IsFake() {
		filtered := make(map[constant.MobBuffFlag]int32, len(values))
		filteredStacks := make(map[constant.MobBuffFlag]uint8, len(values))
		for flag, value := range values {
			if !mob.canReceiveMobBuff(flag) {
				continue
			}
			filtered[flag] = value
			if stacks != nil {
				if stack, ok := stacks[flag]; ok {
					filteredStacks[flag] = stack
				}
			}
		}
		values = filtered
		stacks = filteredStacks
		if len(values) == 0 {
			return
		}
	}
	now := time.Now()
	if bc.refreshDuration(now, duration, skillWz, skillLevel, causer, values, stacks) {
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
	ent := &MobBuff{
		StartTime:  now,
		Duration:   duration,
		Values:     valCopy,
		Stacks:     stackCopy,
		Wz:         skillWz,
		SkillLevel: skillLevel,
		Causer:     causer,
	}
	for flag := range valCopy {
		bc.Remove(flag)
	}
	addedReflections := bc.appendReflectionAdds(valCopy)
	bc.entities[ent] = struct{}{}
	for flag := range ent.Values {
		bc.byFlag[flag] = ent
	}
	mob := bc.owner
	ent.callMobSkillHook(mob, "on_mob_buff")
	remaining := ent.RemainingDuration(now)
	mob.Listener.OnMobBuffApplied(mob, ent, addedReflections, remaining)
}

func (bc *MobBuffContainer) Remove(flag constant.MobBuffFlag) {
	if bc == nil {
		return
	}
	ent := bc.byFlag[flag]
	if ent == nil {
		return
	}
	delete(bc.entities, ent)
	for f := range ent.Values {
		if current := bc.byFlag[f]; current == ent {
			delete(bc.byFlag, f)
		}
	}
	mob := bc.owner
	for f := range ent.Values {
		if f == constant.MobBuffWeaponDamageReflect || f == constant.MobBuffMagicDamageReflect {
			bc.popReflection()
		}
	}
	for f := range ent.Values {
		mob.Listener.OnMobBuffCancelled(mob, f)
	}
	ent.callMobSkillHook(mob, "on_mob_unbuff")
}

func (bc *MobBuffContainer) Dispel(skillID uint32) {
	if bc == nil || skillID == 0 {
		return
	}
	var ents []*MobBuff
	for ent := range bc.entities {
		if ent == nil || ent.Wz == nil || ent.Wz.ID != skillID {
			continue
		}
		ents = append(ents, ent)
	}
	for _, ent := range ents {
		for flag := range ent.Values {
			bc.Remove(flag)
			break
		}
	}
}

func (bc *MobBuffContainer) getExpired(now time.Time) []*MobBuff {
	if bc == nil || len(bc.entities) == 0 {
		return nil
	}
	var out []*MobBuff
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

func (bc *MobBuffContainer) RemoveExpired() {
	expired := bc.getExpired(time.Now())
	for _, ent := range expired {
		for flag := range ent.Values {
			bc.Remove(flag)
			break
		}
	}
}

func (bc *MobBuffContainer) GetValue(flag constant.MobBuffFlag) int32 {
	if bc == nil {
		return 0
	}
	ent := bc.byFlag[flag]
	if ent == nil {
		return 0
	}
	return ent.Values[flag]
}

func (bc *MobBuffContainer) GetStack(flag constant.MobBuffFlag) uint8 {
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

func (bc *MobBuffContainer) SetStack(flag constant.MobBuffFlag, stack uint8) bool {
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

func (bc *MobBuffContainer) Has(flag constant.MobBuffFlag) bool {
	if bc == nil {
		return false
	}
	_, ok := bc.byFlag[flag]
	return ok
}

func (bc *MobBuffContainer) Causer(flag constant.MobBuffFlag) (uint32, bool) {
	if bc == nil {
		return 0, false
	}
	ent := bc.byFlag[flag]
	if ent == nil {
		return 0, false
	}
	return ent.Causer, true
}

func (bc *MobBuffContainer) Clear() {
	if bc == nil {
		return
	}
	flags := make([]constant.MobBuffFlag, 0, len(bc.byFlag))
	for flag := range bc.byFlag {
		flags = append(flags, flag)
	}
	for _, flag := range flags {
		bc.Remove(flag)
	}
	bc.reflections = nil
}

func (bc *MobBuffContainer) Reflections() []int32 {
	if bc == nil || len(bc.reflections) == 0 {
		return nil
	}
	out := make([]int32, len(bc.reflections))
	copy(out, bc.reflections)
	return out
}

func (bc *MobBuffContainer) appendReflectionAdds(values map[constant.MobBuffFlag]int32) []int32 {
	if bc == nil || len(values) == 0 {
		return nil
	}
	var added []int32
	if v, ok := values[constant.MobBuffWeaponDamageReflect]; ok {
		if bc.byFlag[constant.MobBuffWeaponDamageReflect] == nil {
			added = append(added, v)
		}
	}
	if v, ok := values[constant.MobBuffMagicDamageReflect]; ok {
		if bc.byFlag[constant.MobBuffMagicDamageReflect] == nil {
			added = append(added, v)
		}
	}
	if len(added) == 0 {
		return nil
	}
	bc.reflections = append(bc.reflections, added...)
	return added
}

func (bc *MobBuffContainer) popReflection() {
	if bc == nil || len(bc.reflections) == 0 {
		return
	}
	bc.reflections = bc.reflections[1:]
}
