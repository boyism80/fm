package entity

type SkillContainer struct {
	owner   *Character
	entries map[uint32]*SkillEntry
}

func NewSkillContainer(owner *Character) *SkillContainer {
	return &SkillContainer{
		owner:   owner,
		entries: make(map[uint32]*SkillEntry),
	}
}

func (sc *SkillContainer) Bind(skillID uint32, entry *SkillEntry) {
	if entry == nil || entry.Level() < 1 {
		return
	}
	entry.Owner = sc.owner
	sc.entries[skillID] = entry
}

func (sc *SkillContainer) Register(skillID uint32, entry *SkillEntry) {
	if entry == nil || entry.Level() < 1 {
		return
	}
	sc.Bind(skillID, entry)
	ch := sc.owner
	if entry.Level() > 0 {
		ch.Listener.OnSkillPassiveHook(ch, skillID, "on_passive")
	}
	ch.Listener.OnUpdateSkill(ch, skillID, int32(entry.Level()), int32(entry.MasterLevel))
}

func (sc *SkillContainer) Remove(skillID uint32) bool {
	entry := sc.entries[skillID]
	if entry == nil {
		return false
	}
	ch := sc.owner
	if entry.Level() > 0 {
		ch.Listener.OnSkillPassiveHook(ch, skillID, "on_unpassive")
	}
	delete(sc.entries, skillID)
	return true
}

func (sc *SkillContainer) Get(skillID uint32) *SkillEntry {
	return sc.entries[skillID]
}

func (sc *SkillContainer) ForEach(fn func(skillID uint32, entry *SkillEntry)) {
	for id, e := range sc.entries {
		fn(id, e)
	}
}
