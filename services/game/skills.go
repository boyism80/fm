package game

import (
	"log"

	"github.com/boyism80/fm/maple/constant"
	"github.com/boyism80/fm/maple/wz"
)

type Skills struct {
	owner           *Character
	entries         map[uint32]*SkillEntry
	Macros          *SkillMacros
	HomingTargetOID *uint32
}

func NewSkills(owner *Character) *Skills {
	return &Skills{
		owner:   owner,
		entries: make(map[uint32]*SkillEntry),
		Macros:  NewSkillMacros(owner),
	}
}

func (sc *Skills) Bind(skillID uint32, entry *SkillEntry) {
	if entry == nil || entry.Level() < 1 && entry.MasterLevel < 1 {
		return
	}
	entry.Owner = sc.owner
	sc.entries[skillID] = entry
}

func (sc *Skills) Register(skillID uint32, entry *SkillEntry) {
	if entry == nil || entry.Level() < 1 && entry.MasterLevel < 1 {
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
		if entry.Level() > 0 {
			sc.owner.CallSkillHook(nil, entry, "on_passive")
		}
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
	if entry := sc.entries[skillID]; entry != nil && entry.Level() > 0 && !entry.Wz.Dojo {
		return entry
	}
	return sc.dojoSkill(skillID)
}

func (sc *Skills) Learned(skillID uint32) *SkillEntry {
	return sc.entries[skillID]
}

func (sc *Skills) dojoSkill(skillID uint32) *SkillEntry {
	if sc.owner.GameWorld == nil {
		return nil
	}
	w := sc.owner.GameWorld.GetResources().GetSkill(skillID)
	if w == nil || !w.Dojo {
		return nil
	}
	if sc.owner.Dojo.Energy < DojoEnergyFull || !sc.owner.Dojo.onField() {
		return nil
	}
	return NewSkillEntry(sc.owner, w, 1, 1)
}

func (sc *Skills) ForEach(fn func(skillID uint32, entry *SkillEntry)) {
	for id, e := range sc.entries {
		fn(id, e)
	}
}

func (sc *Skills) Level(skillID uint32) int {
	entry := sc.Get(skillID)
	if entry == nil {
		return 0
	}
	return entry.Level()
}

func (sc *Skills) UseAttack(skillID uint32, activate func() bool) bool {
	wzSkill := sc.owner.GameWorld.GetResources().GetSkill(skillID)
	if wzSkill == nil {
		log.Printf("Skill not found: %d", skillID)
		return false
	}
	skillLevel := sc.Level(skillID)
	if skillLevel <= 0 {
		log.Printf("Character does not have skill %d or skill level is 0", skillID)
		return false
	}
	levelData := wzSkill.GetLevelData(skillLevel)
	if levelData == nil {
		return false
	}
	skillEntry := sc.Get(skillID)
	if levelData.Cooldown > 0 && (skillEntry == nil || skillEntry.IsCooling()) {
		return false
	}
	if !activate() {
		return false
	}
	if levelData.Cooldown > 0 {
		skillEntry.StartCooldown(levelData.Cooldown)
	}
	return true
}

func (sc *Skills) LevelData(id constant.SkillID) *wz.SkillLevelData {
	entry := sc.Get(uint32(id))
	if entry == nil || entry.Level() <= 0 {
		return nil
	}
	return entry.Wz.GetLevelData(entry.Level())
}
