package entity

import (
	"github.com/boyism80/fm/services/game/constant"
)

type Sponge struct {
	me       *Mob
	parent   *Mob
	children []*Mob
}

func (s *Sponge) SetParent(mob *Mob) bool {
	if s == nil || s.me == nil {
		return false
	}
	if mob == nil {
		oldParent := s.parent
		if oldParent != nil {
			children := make([]*Mob, 0, len(oldParent.sponge.children))
			for _, c := range oldParent.sponge.children {
				if c != s.me {
					children = append(children, c)
				}
			}
			oldParent.sponge.children = children
		}
		s.parent = nil
		return true
	}
	if mob == s.me {
		return false
	}
	if s.parent != nil {
		return false
	}
	s.parent = mob
	mob.sponge.children = append(mob.sponge.children, s.me)
	return true
}

func (s *Sponge) isFinish() bool {
	if s == nil || len(s.children) == 0 {
		return true
	}
	for _, child := range s.children {
		if child == nil || child.Wz == nil || child.GetHp() == 0 {
			continue
		}
		if child.Wz.Level > 1 {
			return false
		}
	}
	return true
}

func (s *Sponge) applyDamageFromHit(damage uint32) {
	if s == nil || s.me == nil || damage == 0 {
		return
	}
	parent := s.parent
	if parent == nil {
		return
	}

	spongeDamage := damage
	if spongeDamage > parent.GetHp() {
		spongeDamage = parent.GetHp()
	}
	if spongeDamage == 0 {
		return
	}

	parent.LifeCore.AddHp(-int(spongeDamage))
	if parent.Listener == nil {
		return
	}
	parent.Listener.OnShowBossHp(parent, parent.GetHp() == 0)
}

func (s *Sponge) removeAllChildren() {
	if s == nil || len(s.children) == 0 {
		return
	}
	children := append([]*Mob(nil), s.children...)
	for _, child := range children {
		if child == nil {
			continue
		}

		if !child.IsAlive() {
			continue
		}

		child.sponge.SetParent(nil)
		child.Kill(nil, constant.MobDieAnimationTypeFadeOut)
	}
}

func (s *Sponge) Disconnect() {
	if s == nil {
		return
	}
	if s.parent != nil {
		s.SetParent(nil)
	}
	for _, child := range s.children {
		if child != nil {
			child.sponge.SetParent(nil)
		}
	}
	s.children = make([]*Mob, 0)
}

func (s *Sponge) onDead(attacker *Character) {
	if s == nil || s.me == nil {
		return
	}

	parent := s.parent
	if parent != nil {
		s.SetParent(nil)

		if parent.sponge.isFinish() {
			parent.Kill(attacker, constant.MobDieAnimationTypeFadeOut)
		}
	}
	s.removeAllChildren()
}
