package game

import (
	"fmt"
	"time"

	"github.com/boyism80/fm/maple/constant"
	"github.com/boyism80/fm/protocol/dto"
	"github.com/boyism80/fm/types"
)

type Summons struct {
	owner   *Character
	entries map[constant.SkillID]*Summon
}

func NewSummons(owner *Character) *Summons {
	return &Summons{
		owner:   owner,
		entries: make(map[constant.SkillID]*Summon),
	}
}

func (sc *Summons) Spawn(skillID constant.SkillID, skillLevel uint8, movementType constant.SummonMovementType, summonType constant.SummonType, position types.Point[int16], duration time.Duration) *Summon {
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
			sc.Remove(sc.Get(skillID), true)
		})
	}
	return s
}

func (sc *Summons) Remove(s *Summon, animated bool) {
	if s == nil {
		return
	}
	_ = sc.owner.RemoveTimer(sc.timerKey(s.SkillID))
	if m := s.GetMap(); m != nil && s.OID != 0 {
		m.RemoveSummon(s.OID, animated)
	}
	delete(sc.entries, s.SkillID)
}

func (sc *Summons) Attack(oid uint32, damages []dto.AttackPair) (*Summon, error) {
	if sc.owner.GetMap() == nil {
		return nil, nil
	}
	if sc.owner.Spectating() || !sc.owner.IsAlive() {
		return nil, nil
	}
	var summon *Summon
	for _, s := range sc.entries {
		if s.OID == oid {
			summon = s
			break
		}
	}
	if summon == nil {
		return nil, nil
	}
	if err := NewSummonValidator(sc.owner, summon).Validate(damages); err != nil {
		return nil, err
	}

	sc.owner.damageTo(damages)
	return summon, nil
}

func (sc *Summons) Get(skillID constant.SkillID) *Summon {
	return sc.entries[skillID]
}

func (sc *Summons) All() []*Summon {
	out := make([]*Summon, 0, len(sc.entries))
	for _, s := range sc.entries {
		out = append(out, s)
	}
	return out
}

func (sc *Summons) Clear() {
	for _, s := range sc.All() {
		sc.Remove(s, true)
	}
}

func (sc *Summons) timerKey(skillID constant.SkillID) string {
	return fmt.Sprintf("summon:%d", skillID)
}
