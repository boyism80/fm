package game

import (
	"time"

	"github.com/boyism80/fm/core/clock"

	"github.com/boyism80/fm/maple/constant"
	"github.com/boyism80/fm/maple/wz"
)

type MobBuff struct {
	startTime  time.Time
	duration   time.Duration
	values     map[constant.MobBuffFlag]int32
	stacks     map[constant.MobBuffFlag]uint8
	wz         *wz.Skill
	skillLevel uint8
	causer     uint32
}

func (e *MobBuff) RemainingDuration(now time.Time) time.Duration {
	if e == nil || e.duration <= 0 {
		return 0
	}
	end := e.startTime.Add(e.duration)
	if !now.Before(end) {
		return 0
	}
	return end.Sub(now)
}

// Reason is read by the client as one int: a mob skill packs its level above the id, a character skill sends the whole id.
func (e *MobBuff) Reason() uint32 {
	if e.wz == nil {
		return 0
	}
	if e.wz.MobSkill {
		return e.wz.ID | uint32(e.skillLevel)<<16
	}
	return e.wz.ID
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
	if e.wz != nil {
		entSid = e.wz.ID
	}
	return sid == entSid && e.skillLevel == skillLevel && e.causer == causer
}

func (e *MobBuff) Equal(values map[constant.MobBuffFlag]int32, stacks map[constant.MobBuffFlag]uint8) bool {
	if e == nil || len(e.values) != len(values) {
		return false
	}
	for f, v := range values {
		if ev, ok := e.values[f]; !ok || ev != v {
			return false
		}
		wantStack := uint8(1)
		if stacks != nil {
			if sv, ok := stacks[f]; ok && sv >= 1 {
				wantStack = sv
			}
		}
		got := e.stacks[f]
		if got < 1 {
			got = 1
		}
		if got != wantStack {
			return false
		}
	}
	return true
}

type MobBuffs struct {
	owner       *Mob
	byFlag      map[constant.MobBuffFlag]*MobBuff
	entities    map[*MobBuff]struct{}
	reflections []int32
}

func NewMobBuffs(owner *Mob) *MobBuffs {
	if owner == nil {
		panic("MobBuffs owner is nil")
	}
	if owner.listener == nil {
		panic("MobBuffs: mob listener must not be nil")
	}
	return &MobBuffs{
		owner:    owner,
		byFlag:   make(map[constant.MobBuffFlag]*MobBuff),
		entities: make(map[*MobBuff]struct{}),
	}
}

func (bc *MobBuffs) refreshDuration(now time.Time, duration time.Duration, skillWz *wz.Skill, skillLevel uint8, causer uint32, values map[constant.MobBuffFlag]int32, stacks map[constant.MobBuffFlag]uint8) bool {
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
	if !ent.Equal(values, stacks) || !ent.SameSource(skillWz, skillLevel, causer) {
		return false
	}
	ent.startTime = now
	ent.duration = duration
	return true
}

func (bc *MobBuffs) Add(duration time.Duration, skillWz *wz.Skill, skillLevel uint8, causer uint32, values map[constant.MobBuffFlag]int32, stacks map[constant.MobBuffFlag]uint8) {
	if bc == nil || len(values) == 0 {
		return
	}
	{
		filtered := make(map[constant.MobBuffFlag]int32, len(values))
		filteredStacks := make(map[constant.MobBuffFlag]uint8, len(values))
		for flag, value := range values {
			if !bc.owner.canReceiveMobBuff(flag, skillWz) {
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
	now := clock.Now()
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
		startTime:  now,
		duration:   duration,
		values:     valCopy,
		stacks:     stackCopy,
		wz:         skillWz,
		skillLevel: skillLevel,
		causer:     causer,
	}
	for flag := range valCopy {
		bc.Remove(flag)
	}
	addedReflections := bc.appendReflectionAdds(valCopy)
	bc.entities[ent] = struct{}{}
	for flag := range ent.values {
		bc.byFlag[flag] = ent
	}
	mob := bc.owner
	remaining := ent.RemainingDuration(now)
	mob.listener.OnMobBuffApplied(mob, ent, addedReflections, remaining)
}

// Remove cancels one status; the rest of the entry it came with stays.
func (bc *MobBuffs) Remove(flag constant.MobBuffFlag) {
	if bc == nil {
		return
	}
	ent := bc.byFlag[flag]
	if ent == nil {
		return
	}
	delete(bc.byFlag, flag)
	delete(ent.values, flag)
	delete(ent.stacks, flag)
	if len(ent.values) == 0 {
		delete(bc.entities, ent)
	}
	if (flag == constant.MobBuffWeaponDamageReflect || flag == constant.MobBuffMagicDamageReflect) && len(bc.reflections) > 0 {
		bc.reflections = bc.reflections[1:]
	}
	mob := bc.owner
	mob.listener.OnMobBuffCancelled(mob, flag)
}

func (bc *MobBuffs) removeEntry(ent *MobBuff) {
	flags := make([]constant.MobBuffFlag, 0, len(ent.values))
	for flag := range ent.values {
		flags = append(flags, flag)
	}
	for _, flag := range flags {
		bc.Remove(flag)
	}
}

func (bc *MobBuffs) Dispel(skillID uint32) {
	if bc == nil || skillID == 0 {
		return
	}
	var ents []*MobBuff
	for ent := range bc.entities {
		if ent == nil || ent.wz == nil || ent.wz.ID != skillID {
			continue
		}
		ents = append(ents, ent)
	}
	for _, ent := range ents {
		bc.removeEntry(ent)
	}
}

func (bc *MobBuffs) getExpired(now time.Time) []*MobBuff {
	if bc == nil || len(bc.entities) == 0 {
		return nil
	}
	var out []*MobBuff
	for ent := range bc.entities {
		if ent == nil {
			continue
		}
		if ent.duration <= 0 {
			continue
		}
		if now.After(ent.startTime.Add(ent.duration)) {
			out = append(out, ent)
		}
	}
	return out
}

func (bc *MobBuffs) RemoveExpired() {
	expired := bc.getExpired(clock.Now())
	for _, ent := range expired {
		bc.removeEntry(ent)
	}
}

func (bc *MobBuffs) GetValue(flag constant.MobBuffFlag) int32 {
	if bc == nil {
		return 0
	}
	ent := bc.byFlag[flag]
	if ent == nil {
		return 0
	}
	return ent.values[flag]
}

func (bc *MobBuffs) GetStack(flag constant.MobBuffFlag) uint8 {
	if bc == nil {
		return 0
	}
	ent := bc.byFlag[flag]
	if ent == nil {
		return 0
	}
	s := ent.stacks[flag]
	if s < 1 {
		return 1
	}
	return s
}

func (bc *MobBuffs) SetStack(flag constant.MobBuffFlag, stack uint8) bool {
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
	ent.stacks[flag] = stack
	return true
}

func (bc *MobBuffs) Has(flag constant.MobBuffFlag) bool {
	if bc == nil {
		return false
	}
	_, ok := bc.byFlag[flag]
	return ok
}

func (bc *MobBuffs) Causer(flag constant.MobBuffFlag) (uint32, bool) {
	if bc == nil {
		return 0, false
	}
	ent := bc.byFlag[flag]
	if ent == nil {
		return 0, false
	}
	return ent.causer, true
}

func (bc *MobBuffs) Clear() {
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

func (bc *MobBuffs) Reflections() []int32 {
	if bc == nil || len(bc.reflections) == 0 {
		return nil
	}
	out := make([]int32, len(bc.reflections))
	copy(out, bc.reflections)
	return out
}

func (bc *MobBuffs) appendReflectionAdds(values map[constant.MobBuffFlag]int32) []int32 {
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
