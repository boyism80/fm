package entity

import (
	"fmt"
	"time"

	"github.com/boyism80/fm/services/game/constant"
	"github.com/boyism80/fm/types"
)

type SummonContainer struct {
	owner   *Character
	entries map[constant.SkillID]*Summon
}

func NewSummonContainer(owner *Character) *SummonContainer {
	return &SummonContainer{
		owner:   owner,
		entries: make(map[constant.SkillID]*Summon),
	}
}

func (sc *SummonContainer) Spawn(skillID constant.SkillID, skillLevel uint8, movementType constant.SummonMovementType, summonType constant.SummonType, position types.Point[int16], duration time.Duration) *Summon {
	m := sc.owner.GetMap()
	if m == nil {
		return nil
	}
	if current := sc.Get(skillID); current != nil {
		sc.Remove(current, true)
	}

	s := &Summon{
		LifeCore: LifeCore{
			ObjectCore: ObjectCore{
				Position:  position,
				GameWorld: m.GameWorld,
				Map:       nil,
			},
			hp:     1,
			BaseHp: 1,
			BaseMp: 1,
		},
		Owner:        sc.owner,
		OwnerID:      sc.owner.GetID(),
		SkillID:      skillID,
		SkillLevel:   skillLevel,
		MovementType: movementType,
		SummonType:   summonType,
	}
	s.LifeCore.ObjectCore.self = s
	sc.entries[skillID] = s
	m.AddSummon(s)
	if duration > 0 {
		_ = sc.owner.AddTimer(sc.timerKey(skillID), duration, false, func() {
			sc.expire(skillID)
		})
	}
	return s
}

func (sc *SummonContainer) Remove(s *Summon, animated bool) {
	if s == nil {
		return
	}
	_ = sc.owner.RemoveTimer(sc.timerKey(s.SkillID))
	if m := s.GetMap(); m != nil && s.OID != 0 {
		m.RemoveSummon(s.OID, animated)
	}
	delete(sc.entries, s.SkillID)
}

func (sc *SummonContainer) Get(skillID constant.SkillID) *Summon {
	return sc.entries[skillID]
}

func (sc *SummonContainer) All() []*Summon {
	out := make([]*Summon, 0, len(sc.entries))
	for _, s := range sc.entries {
		out = append(out, s)
	}
	return out
}

func (sc *SummonContainer) Clear() {
	for _, s := range sc.All() {
		sc.Remove(s, true)
	}
}

func (sc *SummonContainer) expire(skillID constant.SkillID) {
	sc.Remove(sc.Get(skillID), true)
}

func (sc *SummonContainer) timerKey(skillID constant.SkillID) string {
	return fmt.Sprintf("summon:%d", skillID)
}
