package entity

import (
	"log"

	"github.com/boyism80/fm/services/game/constant"
)

func (m *Mob) MarkSponge() {
	if m == nil {
		return
	}
	m.SpongeMob = true
}

func (m *Mob) IsSpongeMob() bool {
	return m != nil && m.SpongeMob
}

func (m *Mob) SetSponge(sponge *Mob) {
	if m == nil {
		return
	}
	if sponge == nil {
		m.SpongeOID = 0
		return
	}
	m.SpongeOID = sponge.OID
	if m.SpawnLink == 0 {
		m.SpawnLink = sponge.OID
	}
	sponge.MarkSponge()
}

func (m *Mob) GetSponge() *Mob {
	if m == nil || m.SpongeOID == 0 {
		return nil
	}
	mapInstance := m.GetMap()
	if mapInstance == nil {
		return nil
	}
	return mapInstance.GetMob(m.SpongeOID)
}

func (m *Mob) damageSponge(attacker *Character, damage uint32) {
	sponge := m.GetSponge()
	if sponge == nil || sponge.GetHp() == 0 {
		return
	}

	spongeDamage := damage
	if spongeDamage > sponge.GetHp() {
		spongeDamage = sponge.GetHp()
	}
	sponge.AddHp(-int(spongeDamage))

	if sponge.GetHp() == 0 {
		if sponge.Listener != nil {
			sponge.Listener.OnShowBossHp(sponge, true)
		}
		sponge.onKill(attacker)
	} else if sponge.Listener != nil {
		sponge.Listener.OnShowBossHp(sponge, false)
	}
}

func (m *Mob) onKill(attacker *Character) bool {
	mapInstance := m.GetMap()
	if mapInstance != nil {
		m.grantKillExp()
	}

	if attacker != nil {
		if m.Wz != nil {
			log.Printf("Mob %d (ID: %d) killed by character %d", m.OID, m.Wz.ID, attacker.GetID())
		}
		m.dropItems(attacker)
	} else if m.Wz != nil {
		log.Printf("Mob %d (ID: %d) killed with no attacker", m.OID, m.Wz.ID)
	}

	if mapInstance != nil {
		pos := m.Position
		linkOID := m.OID
		revives := []uint32(nil)
		if !m.IsFake() && m.Wz != nil && len(m.Wz.Revives) > 0 {
			revives = m.Wz.Revives
		}

		mapInstance.runDieScript(m, attacker)
		if len(revives) > 0 {
			m.handleRevives(mapInstance, pos, linkOID, revives)
		}

		oldSponge := m.GetSponge()
		m.SetSponge(nil)
		if oldSponge != nil && oldSponge.GetHp() > 0 && !mapInstance.spongePartsAlive(oldSponge, m.OID) {
			if remaining := oldSponge.GetHp(); remaining > 0 {
				oldSponge.ApplyDamage(attacker, remaining)
			}
		}

		mapInstance.RemoveMob(m.OID, constant.MobDieAnimationTypeFadeOut)
	}

	if attacker != nil {
		attacker.Listener.OnShowMobHp(attacker, m, 0)
	}

	return true
}
