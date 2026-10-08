package entity

import (
	"fmt"
	"time"

	"github.com/boyism80/fm/core/clock"
	"github.com/boyism80/fm/services/game/constant"
)

type Debuff struct {
	Flag       constant.DebuffFlag
	StartTime  time.Time
	Duration   time.Duration
	X          int16
	SkillID    uint16
	SkillLevel uint16
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
	d.owner.RemoveTimer(d.timerKey(holder.Flag))
	d.entries[holder.Flag] = holder
	if holder.Duration > 0 {
		flag := holder.Flag
		d.owner.AddTimer(d.timerKey(flag), holder.Duration, false, func() {
			d.owner.RemoveTimer(d.timerKey(flag))
			if _, ok := d.entries[flag]; ok {
				delete(d.entries, flag)
				d.owner.Listener.OnDebuffRemoved(d.owner, []constant.DebuffFlag{flag})
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
		Flag:       flag,
		StartTime:  clock.Now(),
		Duration:   duration,
		X:          x,
		SkillID:    skillID,
		SkillLevel: skillLevel,
	}
	d.Add(holder)
	d.owner.Listener.OnDebuffAdded(d.owner, flag, x, skillID, skillLevel, int32(duration.Milliseconds()))
}

func (d *Debuffs) restore() {
	restored := make([]*Debuff, 0, len(d.entries))
	for _, holder := range d.entries {
		restored = append(restored, holder)
	}
	for _, holder := range restored {
		d.Add(holder)
		d.owner.Listener.OnDebuffAdded(d.owner, holder.Flag, holder.X, holder.SkillID, holder.SkillLevel, int32(holder.Duration.Milliseconds()))
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
		d.owner.Listener.OnDebuffRemoved(d.owner, removed)
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
