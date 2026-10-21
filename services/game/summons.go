package game

import (
	"fmt"
	"github.com/asynkron/protoactor-go/actor"
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
		lifeCore: LifeCore{
			objectCore: ObjectCore{
				position:  position,
				gameWorld: m.gameWorld,
				onMap:     nil,
			},
			hp:     1,
			baseHp: 1,
			baseMp: 1,
		},
		owner:        sc.owner,
		ownerID:      sc.owner.GetID(),
		skillID:      skillID,
		skillLevel:   skillLevel,
		movementType: movementType,
		summonType:   summonType,
	}
	s.lifeCore.objectCore.self = s
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
	_ = sc.owner.RemoveTimer(sc.timerKey(s.skillID))
	if m := s.GetMap(); m != nil && s.oid != 0 {
		m.RemoveSummon(s.oid, animated)
	}
	delete(sc.entries, s.skillID)
}

func (sc *Summons) Attack(actx actor.Context, oid uint32, damages []dto.AttackPair) (*Summon, error) {
	if sc.owner.Spectating() || !sc.owner.IsAlive() {
		return nil, nil
	}
	summon := sc.FindByOID(oid)
	if summon == nil {
		return nil, nil
	}
	sc.owner.cheat.PaceSummon(actx, summon)
	if err := NewSummonValidator(sc.owner, summon).Validate(damages); err != nil {
		sc.owner.cheat.Rejected(actx, dto.AttackInfo{AttackHeader: dto.AttackHeader{Skill: uint32(summon.skillID)}}, err)
		return nil, err
	}

	sc.owner.damageTo(damages)
	return summon, nil
}

func (sc *Summons) Get(skillID constant.SkillID) *Summon {
	return sc.entries[skillID]
}

func (sc *Summons) FindByOID(oid uint32) *Summon {
	m := sc.owner.GetMap()
	if m == nil {
		return nil
	}
	for _, s := range sc.entries {
		if s.oid == oid && s.onMap == m {
			return s
		}
	}
	return nil
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
