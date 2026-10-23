package game

import (
	"fmt"
	"time"

	"github.com/boyism80/fm/core/clock"
	"github.com/boyism80/fm/maple/constant"
)

type Debuff struct {
	flag       constant.DebuffFlag
	startTime  time.Time
	duration   time.Duration
	x          int16
	skillID    uint16
	skillLevel uint16
}

type Debuffs struct {
	owner   *Character
	entries map[constant.DebuffFlag]*Debuff
}

func (d *Debuffs) timerKey(flag constant.DebuffFlag) string {
	return fmt.Sprintf("debuff_%d_%d", flag.Position, flag.Mask)
}

func (d *Debuffs) Has(flag constant.DebuffFlag) bool {
	_, ok := d.entries[flag]
	return ok
}

func (d *Debuffs) Add(holder *Debuff) {
	if holder == nil {
		return
	}
	d.owner.RemoveTimer(d.timerKey(holder.flag))
	d.entries[holder.flag] = holder
	if holder.duration > 0 {
		flag := holder.flag
		d.owner.AddTimer(d.timerKey(flag), holder.duration, false, func() {
			d.owner.RemoveTimer(d.timerKey(flag))
			if _, ok := d.entries[flag]; ok {
				delete(d.entries, flag)
				d.owner.listener.OnDebuffRemoved(d.owner, []constant.DebuffFlag{flag})
			}
		})
	}
}

func (d *Debuffs) Give(flag constant.DebuffFlag, duration time.Duration, x int16, skillID uint16, skillLevel uint16) {
	if skillID == 0 {
		skillID = flag.DiseaseSkillID
	}
	if skillLevel == 0 {
		skillLevel = 1
	}
	holder := &Debuff{
		flag:       flag,
		startTime:  clock.Now(),
		duration:   duration,
		x:          x,
		skillID:    skillID,
		skillLevel: skillLevel,
	}
	d.Add(holder)
	d.owner.listener.OnDebuffAdded(d.owner, flag, x, skillID, skillLevel, int32(duration.Milliseconds()))
}

func (d *Debuffs) restore() {
	restored := make([]*Debuff, 0, len(d.entries))
	for _, holder := range d.entries {
		restored = append(restored, holder)
	}
	for _, holder := range restored {
		d.Add(holder)
		d.owner.listener.OnDebuffAdded(d.owner, holder.flag, holder.x, holder.skillID, holder.skillLevel, int32(holder.duration.Milliseconds()))
	}
}

func (d *Debuffs) Remove(flags ...constant.DebuffFlag) {
	var removed []constant.DebuffFlag
	for _, flag := range flags {
		d.owner.RemoveTimer(d.timerKey(flag))
		if _, ok := d.entries[flag]; ok {
			delete(d.entries, flag)
			removed = append(removed, flag)
		}
	}
	if len(removed) > 0 {
		d.owner.listener.OnDebuffRemoved(d.owner, removed)
	}
}

func (d *Debuffs) DiseaseMask() [4]uint32 {
	var mask [4]uint32
	for flag := range d.entries {
		idx := flag.Position - 1
		if idx >= 0 && idx < constant.MaxBuffFlag {
			mask[idx] |= flag.Mask
		}
	}
	return mask
}

// SpawnReasons gives each disease's reason (mob skill id with its level above) and the poison damage.
func (d *Debuffs) SpawnReasons() (map[constant.DebuffFlag]uint32, int16) {
	reasons := make(map[constant.DebuffFlag]uint32, len(d.entries))
	var poison int16
	for flag, holder := range d.entries {
		reasons[flag] = uint32(holder.skillID) | uint32(holder.skillLevel)<<16
		if flag == constant.DebuffFlagPoison {
			poison = holder.x
		}
	}
	return reasons, poison
}

// Returns false when the mob skill is not a disease, so the caller can fall back to a dispel.
func (d *Debuffs) GiveMobSkill(skillID uint32, skillLevel uint8) bool {
	mobSkill := d.owner.gameWorld.GetResources().GetMobSkill(skillID, skillLevel)
	if mobSkill == nil {
		return false
	}
	for _, flag := range constant.AllDebuffFlags() {
		if flag.DiseaseSkillID != uint16(skillID) {
			continue
		}
		duration := time.Duration(mobSkill.DurationMs) * time.Millisecond
		if duration <= 0 {
			duration = 5 * time.Second
		}
		d.Give(flag, duration, int16(mobSkill.X), uint16(skillID), uint16(skillLevel))
		return true
	}
	return false
}
