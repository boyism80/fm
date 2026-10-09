package entity

import (
	"time"

	"github.com/boyism80/fm/core/clock"
	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
	"github.com/boyism80/fm/services/game/constant"
)

func (d *Debuffs) Load(pbs []*internal.Debuff) {
	now := clock.Now()
	flags := constant.AllDebuffFlags()
	for _, pb := range pbs {
		var duration time.Duration
		if pb.GetEndUnixMs() != 0 {
			duration = time.UnixMilli(pb.GetEndUnixMs()).Sub(now)
			if duration <= 0 {
				continue
			}
		}
		for _, flag := range flags {
			if flag.Mask != pb.GetMask() || flag.Position != int(pb.GetPosition()) {
				continue
			}
			d.entries[flag] = &Debuff{
				Flag:       flag,
				StartTime:  now,
				Duration:   duration,
				X:          int16(pb.GetX()),
				SkillID:    uint16(pb.GetSkillId()),
				SkillLevel: uint16(pb.GetSkillLevel()),
			}
		}
	}
}

func (d *Debuffs) ToProto() []*internal.Debuff {
	now := clock.Now()
	out := make([]*internal.Debuff, 0, len(d.entries))
	for _, holder := range d.entries {
		var endUnixMs int64
		if holder.Duration > 0 {
			end := holder.StartTime.Add(holder.Duration)
			if !end.After(now) {
				continue
			}
			endUnixMs = end.UnixMilli()
		}
		out = append(out, &internal.Debuff{
			Mask:       holder.Flag.Mask,
			Position:   int32(holder.Flag.Position),
			X:          int32(holder.X),
			SkillId:    uint32(holder.SkillID),
			SkillLevel: uint32(holder.SkillLevel),
			EndUnixMs:  endUnixMs,
		})
	}
	return out
}
