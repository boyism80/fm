package entity

import "github.com/boyism80/fm/services/game/constant"

type Skills struct {
	owner   *Character
	entries map[uint32]*SkillEntry
}

func NewSkills(owner *Character) *Skills {
	return &Skills{
		owner:   owner,
		entries: make(map[uint32]*SkillEntry),
	}
}

func (sc *Skills) Bind(skillID uint32, entry *SkillEntry) {
	if entry == nil || entry.Level() < 1 {
		return
	}
	entry.Owner = sc.owner
	sc.entries[skillID] = entry
}

func (sc *Skills) Register(skillID uint32, entry *SkillEntry) {
	if entry == nil || entry.Level() < 1 {
		return
	}
	sc.Bind(skillID, entry)
	ch := sc.owner
	if entry.Level() > 0 {
		ch.CallSkillHook(nil, entry, "on_passive")
	}
	ch.Listener.OnUpdateSkill(ch, skillID, int32(entry.Level()), int32(entry.MasterLevel))
}

func (sc *Skills) RestorePassives() {
	for _, entry := range sc.entries {
		sc.owner.CallSkillHook(nil, entry, "on_passive")
	}
}

func (sc *Skills) Remove(skillID uint32) bool {
	entry := sc.entries[skillID]
	if entry == nil {
		return false
	}
	ch := sc.owner
	if entry.Level() > 0 {
		ch.CallSkillHook(nil, entry, "on_unpassive")
	}
	delete(sc.entries, skillID)
	return true
}

func (sc *Skills) Get(skillID uint32) *SkillEntry {
	if entry := sc.entries[skillID]; entry != nil {
		return entry
	}
	return sc.dojoSkill(skillID)
}

func (sc *Skills) dojoSkill(skillID uint32) *SkillEntry {
	switch constant.SkillID(skillID) {
	case constant.SkillBambooRain, constant.SkillInvincibility, constant.SkillPowerExplosion,
		constant.SkillBambooRainCygnus, constant.SkillInvincibilityCygnus, constant.SkillPowerExplosionCygnus:
	default:
		return nil
	}
	if sc.owner.dojoEnergy < DojoEnergyFull || !sc.owner.OnDojoField() {
		return nil
	}
	w := sc.owner.GameWorld.GetResources().GetSkill(skillID)
	if w == nil {
		return nil
	}
	return NewSkillEntry(sc.owner, w, 1, 1)
}

func (sc *Skills) ForEach(fn func(skillID uint32, entry *SkillEntry)) {
	for id, e := range sc.entries {
		fn(id, e)
	}
}
