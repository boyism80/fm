package game

import (
	"time"

	"github.com/boyism80/fm/core/clock"
	"github.com/boyism80/fm/maple/constant"
	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
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
				flag:       flag,
				startTime:  now,
				duration:   duration,
				x:          int16(pb.GetX()),
				skillID:    uint16(pb.GetSkillId()),
				skillLevel: uint16(pb.GetSkillLevel()),
			}
		}
	}
}

func (d *Debuffs) BuildRPC() []*internal.Debuff {
	now := clock.Now()
	out := make([]*internal.Debuff, 0, len(d.entries))
	for _, holder := range d.entries {
		var endUnixMs int64
		if holder.duration > 0 {
			end := holder.startTime.Add(holder.duration)
			if !end.After(now) {
				continue
			}
			endUnixMs = end.UnixMilli()
		}
		out = append(out, &internal.Debuff{
			Mask:       holder.flag.Mask,
			Position:   int32(holder.flag.Position),
			X:          int32(holder.x),
			SkillId:    uint32(holder.skillID),
			SkillLevel: uint32(holder.skillLevel),
			EndUnixMs:  endUnixMs,
		})
	}
	return out
}
