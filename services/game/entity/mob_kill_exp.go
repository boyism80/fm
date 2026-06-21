package entity

import (
	"fmt"
	"log"
	"math"
)

const soloDamageBucketID int64 = -1

// grantKillExp splits mob kill EXP among accDamage contributors on this mob's current map (alive).
// Party groups use Maple-style pool + bonus + per-damager splits; solos use ExpForDamage per row.
func (m *Mob) grantKillExp() {
	if m == nil || m.Wz == nil {
		return
	}
	mapInst := m.GetMap()
	if mapInst == nil {
		return
	}

	rawByChar := make(map[uint32]uint64)
	mobExpRate := m.ExpRate
	if mobExpRate <= 0 {
		mobExpRate = 100
	}

	for bucketID, damagers := range m.accDamage {
		if len(damagers) == 0 {
			continue
		}
		if bucketID == soloDamageBucketID {
			for cid, dmg := range damagers {
				if dmg == 0 {
					continue
				}
				raw := m.ExpForDamage(dmg)
				raw = raw * uint32(mobExpRate) / 100
				if raw == 0 {
					continue
				}
				rawByChar[cid] += uint64(raw)
			}
		} else {
			partyRaw, err := m.grantPartyKillExp(mapInst, uint32(bucketID), damagers)
			if err != nil {
				log.Printf("grantKillExp: mob_oid=%d party_id=%d err=%v", m.GetOID(), uint32(bucketID), err)
				continue
			}
			for cid, raw := range partyRaw {
				if raw == 0 {
					continue
				}
				rawByChar[cid] += raw
			}
		}
	}

	for cid, rawTotal := range rawByChar {
		if rawTotal == 0 {
			continue
		}
		ch := mapInst.GetPlayer(cid)
		if ch == nil || ch.GetMap() != mapInst || !ch.IsAlive() {
			continue
		}
		raw := rawTotal
		if raw > uint64(^uint32(0)) {
			raw = uint64(^uint32(0))
		}
		ch.AddExp(ch.ComputeMobKillExp(uint32(raw)))
	}
}

func (m *Mob) grantPartyKillExp(mapInst *Map, partyID uint32, damagers map[uint32]uint64) (map[uint32]uint64, error) {
	if len(damagers) == 0 {
		return map[uint32]uint64{}, nil
	}

	gw := mapInst.GameWorld
	if gw == nil {
		return nil, fmt.Errorf("nil GameWorld")
	}
	party := gw.GetPartySystem().Get(partyID)

	var totDamage uint64
	for _, dmg := range damagers {
		totDamage += dmg
	}
	if totDamage == 0 {
		return map[uint32]uint64{}, nil
	}
	poolRaw := m.ExpForDamage(totDamage)
	if poolRaw == 0 {
		return map[uint32]uint64{}, nil
	}
	mobExpRate := m.ExpRate
	if mobExpRate <= 0 {
		mobExpRate = 100
	}
	poolRaw = poolRaw * uint32(mobExpRate) / 100

	mobLv := int(m.Wz.Level)
	expApplicable := make([]*Character, 0)
	if party == nil {
		// Party disbanded or snapshot unavailable: degrade to current damagers only.
		for damagerID := range damagers {
			damager := mapInst.GetPlayer(damagerID)
			if damager == nil || damager.GetMap() != mapInst || !damager.IsAlive() {
				continue
			}
			expApplicable = append(expApplicable, damager)
		}
	} else {
		expApplicable = make([]*Character, 0, len(party.Members))
		for _, pm := range party.Members {
			if pm == nil {
				continue
			}
			pch := mapInst.GetPlayer(pm.CharacterID)
			if pch == nil || pch.GetMap() != mapInst || !pch.IsAlive() {
				continue
			}
			pl := int(pch.GetLevel())
			eligible := false
			for damagerID := range damagers {
				damager := mapInst.GetPlayer(damagerID)
				if damager == nil || damager.GetMap() != mapInst || !damager.IsAlive() {
					continue
				}
				dl := int(damager.GetLevel())
				if math.Abs(float64(dl-pl)) <= 5 || math.Abs(float64(mobLv-pl)) <= 5 {
					eligible = true
					break
				}
			}
			if eligible {
				expApplicable = append(expApplicable, pch)
			}
		}
	}

	if len(expApplicable) == 0 {
		return nil, fmt.Errorf("no exp applicable members")
	}

	n := len(expApplicable)
	var avgLv float64
	for _, p := range expApplicable {
		avgLv += float64(p.GetLevel())
	}
	if n > 1 {
		avgLv /= float64(n)
	} else {
		avgLv = float64(expApplicable[0].GetLevel())
	}
	if avgLv <= 0 {
		avgLv = 1
	}

	expBonus := 1.0
	if n > 1 {
		expBonus = 1.0 + 0.05*float64(n)
	}

	rawAcc := make(map[uint32]float64)
	for damagerID, dmg := range damagers {
		damager := mapInst.GetPlayer(damagerID)
		if damager == nil || damager.GetMap() != mapInst || !damager.IsAlive() {
			continue
		}
		innerBase := float64(poolRaw) * float64(dmg) / float64(totDamage)
		expFraction := innerBase * expBonus / float64(n+1)
		for _, recv := range expApplicable {
			w := 1.0
			if recv.GetID() == damagerID {
				w = 2.0
			}
			levelMod := float64(recv.GetLevel()) / avgLv
			if levelMod > 1.0 {
				levelMod = 1.0
			}
			if _, hit := damagers[recv.GetID()]; hit {
				levelMod = 1.0
			}
			rawAcc[recv.GetID()] += expFraction * w * levelMod
		}
	}

	out := make(map[uint32]uint64, len(rawAcc))
	for rid, v := range rawAcc {
		raw := uint32(math.Round(v))
		if raw == 0 {
			continue
		}
		out[rid] += uint64(raw)
	}
	return out, nil
}
