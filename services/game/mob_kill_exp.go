package game

import (
	"fmt"
	"log"
	"math"

	"github.com/boyism80/fm/core/luax"
	"github.com/boyism80/fm/maple/constant"
)

const soloDamageBucketID int64 = -1

// mobCredit holds resolved recipients and values for mob kill rewards,
// computed in a single pass over m.accDamage.
type mobCredit struct {
	dropOwner  *Character            // highest damager for item drops (falls back to attacker if none)
	expRewards map[*Character]uint64 // EXP share per character; receivers with 0 exp still qualify for kill count
}

// credit resolves mob kill reward recipients (drop owner and exp/quest recipients)
// in a single pass over m.accDamage.
func (m *Mob) credit(attacker *Character) mobCredit {
	credit := mobCredit{
		expRewards: make(map[*Character]uint64),
	}
	if m == nil || m.wz == nil {
		return credit
	}
	mapInst := m.GetMap()
	if mapInst == nil {
		return credit
	}

	mobExpRate := m.expRate
	if mobExpRate <= 0 {
		mobExpRate = 100
	}
	mobExpRate += m.buffs.GetValue(constant.MobBuffShowdown)

	var winBucketTotal uint64
	var winTopCID uint32

	for bucketID, damagers := range m.accDamage {
		if len(damagers) == 0 {
			continue
		}

		var bucketTotal uint64
		var bucketTopCID uint32
		var bucketTopDmg uint64

		for cid, dmg := range damagers {
			bucketTotal += dmg
			if dmg > bucketTopDmg {
				bucketTopDmg = dmg
				bucketTopCID = cid
			}
		}

		if bucketTotal > winBucketTotal {
			winBucketTotal = bucketTotal
			winTopCID = bucketTopCID
		}

		if bucketID == soloDamageBucketID {
			for cid, dmg := range damagers {
				if dmg == 0 {
					continue
				}
				ch := mapInst.GetPlayer(cid)
				if ch == nil || ch.GetMap() != mapInst || !ch.IsAlive() {
					continue
				}
				raw := m.ExpForDamage(dmg)
				raw = raw * uint32(mobExpRate) / 100
				credit.expRewards[ch] += uint64(raw)
			}
		} else {
			partyRewards, err := m.grantKillExp(mapInst, uint32(bucketID), damagers, mobExpRate)
			if err != nil {
				log.Printf("credit: mob_oid=%d party_id=%d err=%v", m.GetOID(), uint32(bucketID), err)
				continue
			}
			for ch, raw := range partyRewards {
				credit.expRewards[ch] += raw
			}
		}
	}

	if winTopCID != 0 {
		if topChar := mapInst.GetPlayer(winTopCID); topChar != nil && topChar.GetMap() == mapInst && topChar.IsAlive() {
			credit.dropOwner = topChar
		}
	}
	if credit.dropOwner == nil {
		credit.dropOwner = attacker
	}

	return credit
}

// applyCredit distributes EXP, records quest kill progress, and triggers character kill hooks
// for all qualifying reward receivers.
func (m *Mob) applyCredit(credit mobCredit) {
	if m == nil || m.wz == nil {
		return
	}
	mapInst := m.GetMap()
	if mapInst == nil {
		return
	}

	for ch, rawTotal := range credit.expRewards {
		if ch == nil || ch.GetMap() != mapInst || !ch.IsAlive() {
			continue
		}
		if rawTotal > 0 {
			raw := rawTotal
			if raw > uint64(^uint32(0)) {
				raw = uint64(^uint32(0))
			}
			ch.AddExp(ch.ComputeMobKillExp(uint32(raw)))
		}
		ch.quests.RecordMobKill(m.wz.ID)
		mapInst.runMobLuaHook(mapInst.GetLuaRoot(), constant.CharacterHookScriptPath, "on_mob_kill", ch, []luax.Luable{m})
	}
}

func (m *Mob) grantKillExp(mapInst *Map, partyID uint32, damagers map[uint32]uint64, mobExpRate int32) (map[*Character]uint64, error) {
	if len(damagers) == 0 {
		return map[*Character]uint64{}, nil
	}

	gw := mapInst.gameWorld
	if gw == nil {
		return nil, fmt.Errorf("nil GameWorld")
	}
	party := gw.Systems().Get[*PartySystem]().Get(partyID)

	var totDamage uint64
	for _, dmg := range damagers {
		totDamage += dmg
	}
	if totDamage == 0 {
		return map[*Character]uint64{}, nil
	}
	poolRaw := m.ExpForDamage(totDamage)
	poolRaw = poolRaw * uint32(mobExpRate) / 100

	mobLv := int(m.wz.Level)

	// Collect active, alive damagers on map once to avoid repeated map lookups
	activeDamagers := make([]*Character, 0, len(damagers))
	for damagerID := range damagers {
		d := mapInst.GetPlayer(damagerID)
		if d == nil || d.GetMap() != mapInst || !d.IsAlive() {
			continue
		}
		activeDamagers = append(activeDamagers, d)
	}

	expApplicable := make([]*Character, 0)
	if party == nil {
		// Party disbanded or snapshot unavailable: degrade to current damagers only.
		expApplicable = append(expApplicable, activeDamagers...)
	} else {
		expApplicable = make([]*Character, 0, len(party.members))
		for _, pm := range party.members {
			if pm == nil {
				continue
			}
			pch := mapInst.GetPlayer(pm.characterID)
			if pch == nil || pch.GetMap() != mapInst || !pch.IsAlive() {
				continue
			}
			pl := int(pch.GetLevel())
			eligible := false
			for _, damager := range activeDamagers {
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

	rawAcc := make(map[*Character]float64, len(expApplicable))
	for _, damager := range activeDamagers {
		dmg := damagers[damager.GetID()]
		innerBase := float64(poolRaw) * float64(dmg) / float64(totDamage)
		expFraction := innerBase * expBonus / float64(n+1)
		for _, recv := range expApplicable {
			w := 1.0
			if recv.GetID() == damager.GetID() {
				w = 2.0
			}
			levelMod := float64(recv.GetLevel()) / avgLv
			if levelMod > 1.0 {
				levelMod = 1.0
			}
			if _, hit := damagers[recv.GetID()]; hit {
				levelMod = 1.0
			}
			rawAcc[recv] += expFraction * w * levelMod
		}
	}

	out := make(map[*Character]uint64, len(expApplicable))
	for _, recv := range expApplicable {
		out[recv] = uint64(math.Round(rawAcc[recv]))
	}
	return out, nil
}
