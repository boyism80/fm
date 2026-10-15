package game

import (
	"time"

	"github.com/boyism80/fm/core/clock"
	"github.com/boyism80/fm/maple/constant"
	"github.com/boyism80/fm/maple/wz"
	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
)

func (bc *Buffs) Load(pbs []*internal.Buff) {
	res := bc.owner.GameWorld.GetResources()
	if res == nil {
		return
	}
	const permanentDur = 100 * 365 * 24 * time.Hour
	for _, pb := range pbs {
		if pb == nil {
			continue
		}
		values := make(map[constant.BuffFlag]int32, len(pb.GetFlagValues()))
		for _, fv := range pb.GetFlagValues() {
			if fv == nil {
				continue
			}
			values[constant.BuffFlag{Mask: fv.GetMask(), Position: int(fv.GetPosition())}] = fv.GetValue()
		}
		if len(values) == 0 {
			continue
		}
		var dur time.Duration
		if pb.RemainingDurationMs != nil {
			if *pb.RemainingDurationMs == 0 {
				continue
			}
			dur = time.Duration(*pb.RemainingDurationMs) * time.Millisecond
		} else {
			dur = permanentDur
		}
		switch pb.GetKind() {
		case internal.BuffKind_BUFF_KIND_SKILL:
			sid := uint32(pb.GetBuffSourceId())
			wzSkill := res.GetSkill(sid)
			if wzSkill == nil {
				continue
			}
			bc.AddBuff(wzSkill, dur, uint8(pb.GetSkillLevel()), pb.GetCauserId(), values, false)
		case internal.BuffKind_BUFF_KIND_ITEM:
			itemID := uint32(-pb.GetBuffSourceId())
			model, ok := res.Items[itemID]
			if !ok {
				continue
			}
			cw, ok := model.(*wz.Consume)
			if !ok || cw == nil {
				continue
			}
			bc.AddItemBuff(cw, dur, values, false, false)
		default:
			continue
		}
	}
}

func (bc *Buffs) BuildRPC() []*internal.Buff {
	now := clock.Now()
	out := make([]*internal.Buff, 0)
	for _, ent := range bc.Entities() {
		if ent == nil {
			continue
		}
		flags := ent.GetFlags()
		values := ent.GetValues()
		if len(flags) == 0 || len(values) == 0 {
			continue
		}
		flagPB := make([]*internal.BuffFlagValue, 0, len(flags))
		for _, f := range flags {
			v, ok := values[f]
			if !ok {
				continue
			}
			flagPB = append(flagPB, &internal.BuffFlagValue{
				Mask:     f.Mask,
				Position: int32(f.Position),
				Value:    v,
			})
		}
		if len(flagPB) == 0 {
			continue
		}
		var kind internal.BuffKind
		var skillLevel uint32
		var causerID uint32
		var rem *uint64
		switch b := ent.(type) {
		case *SkillBuff:
			kind = internal.BuffKind_BUFF_KIND_SKILL
			skillLevel = uint32(b.SkillLevel)
			causerID = b.CauserID
			if b.BaseBuff != nil && b.Duration > 0 {
				ms := uint64(ent.RemainingDuration(now) / time.Millisecond)
				if ms == 0 {
					continue
				}
				rem = &ms
			}
		case *ItemBuff:
			kind = internal.BuffKind_BUFF_KIND_ITEM
			if b.BaseBuff != nil && b.Duration > 0 {
				ms := uint64(ent.RemainingDuration(now) / time.Millisecond)
				if ms == 0 {
					continue
				}
				rem = &ms
			}
		default:
			continue
		}
		out = append(out, &internal.Buff{
			CharacterId:         bc.owner.GetID(),
			BuffSourceId:        ent.GetBuffID(),
			Kind:                kind,
			FlagValues:          flagPB,
			RemainingDurationMs: rem,
			SkillLevel:          skillLevel,
			CauserId:            causerID,
		})
	}
	return out
}
