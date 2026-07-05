package entity

import (
	"fmt"
	"github.com/boyism80/fm/core/clock"
	"time"

	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
)

func (e *SkillEntry) ToProto(characterID uint32, skillID uint32) *internal.SkillPersisted {
	if e == nil {
		return nil
	}
	cooldownEnd := int64(0)
	if e.CooldownEnd != nil {
		cooldownEnd = e.CooldownEnd.UnixMilli()
	}
	return &internal.SkillPersisted{
		CharacterId:       characterID,
		SkillId:           skillID,
		Level:             int32(e.Level()),
		MasterLevel:       int32(e.MasterLevel),
		CooldownEndUnixMs: cooldownEnd,
	}
}

func NewSkillEntryFromInternalProto(owner *Character, pb *internal.SkillPersisted, gw GameWorld) (*SkillEntry, error) {
	if pb == nil {
		return nil, fmt.Errorf("nil SkillPersisted")
	}
	if gw == nil {
		return nil, fmt.Errorf("nil GameWorld")
	}
	resources := gw.GetResources()
	if resources == nil {
		return nil, fmt.Errorf("nil resources")
	}
	skillID := pb.GetSkillId()
	w, ok := resources.Skills[skillID]
	if !ok {
		return nil, fmt.Errorf("skill %d not found in resources", skillID)
	}
	entry := NewSkillEntry(owner, w, int(pb.GetLevel()), int(pb.GetMasterLevel()))
	if cd := pb.GetCooldownEndUnixMs(); cd > 0 {
		t := time.UnixMilli(cd)
		if clock.Now().Before(t) {
			entry.CooldownEnd = &t
		}
	}
	return entry, nil
}
