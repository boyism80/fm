package game

import (
	"slices"
	"time"

	"github.com/boyism80/fm/core/clock"
	"github.com/boyism80/fm/maple/constant"
	"github.com/boyism80/fm/maple/wz"
)

func (qc *Quests) passesRequirements(req wz.QuestRequirements, qp *Quest, opts QuestPhaseOpts) bool {
	if req.NPC != 0 && opts.NpcID != nil {
		wireNPC := *opts.NpcID
		if wireNPC != 0 && wireNPC != req.NPC {
			return false
		}
		mapInst := qc.owner.GetMap()
		if mapInst == nil {
			return false
		}
		found := false
		for _, obj := range mapInst.GetNpcs() {
			npc, ok := obj.(*Npc)
			if !ok || npc.wz == nil || npc.wz.BaseSpawn == nil {
				continue
			}
			if npc.wz.ID == req.NPC {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}

	if req.LevelMin > 0 && int(qc.owner.GetLevel()) < req.LevelMin {
		return false
	}
	if req.LevelMax > 0 && int(qc.owner.GetLevel()) > req.LevelMax {
		return false
	}
	if req.Level > 0 && int(qc.owner.GetLevel()) < req.Level {
		return false
	}
	if len(req.Classes) > 0 && !matchesQuestClass(qc.owner.class, req.Classes) {
		return false
	}
	for itemID, count := range req.Item {
		if count <= 0 {
			if qc.owner.inventory.HasItem(itemID) {
				return false
			}
			continue
		}
		if !qc.owner.inventory.HasItemCount(itemID, uint16(count)) {
			return false
		}
	}
	if len(req.Mob) > 0 {
		for mobID, count := range req.Mob {
			kills := 0
			if qp != nil && qp.mobKills != nil {
				kills = qp.mobKills[mobID]
			}
			if kills < count {
				return false
			}
		}
	}
	for questID, state := range req.Quest {
		if !qc.Get(questID).MatchesState(state) {
			return false
		}
	}
	if req.Pop > 0 && int(qc.owner.stats.population) < req.Pop {
		return false
	}
	if req.FieldEnter > 0 {
		mapInst := qc.owner.GetMap()
		if mapInst == nil || int(mapInst.TemplateID()) != req.FieldEnter {
			return false
		}
	}
	if req.HasInterval {
		if qp != nil && qp.status == QuestStatusCompleted && !qp.completionTime.IsZero() {
			if req.Interval > 0 && clock.Now().Before(qp.completionTime.Add(time.Duration(req.Interval)*time.Minute)) {
				return false
			}
		}
	}
	if req.DayByDay {
		if qp != nil && qp.status == QuestStatusCompleted && !qp.completionTime.IsZero() {
			if sameCalendarDay(qp.completionTime, clock.Now()) {
				return false
			}
		}
	}
	if !req.Start.IsZero() && clock.Now().Before(req.Start) {
		return false
	}
	if !req.End.IsZero() && !clock.Now().Before(req.End) {
		return false
	}
	if req.InfoNumber > 0 {
		refQP := qc.Get(uint32(req.InfoNumber))
		if refQP == nil || !refQP.IsStarted() {
			return false
		}
		expected := req.Info
		if len(expected) == 0 {
			if refQP.statusRecord.IsEmpty() {
				return false
			}
		} else {
			record := refQP.statusRecord.AsString()
			matched := false
			for _, want := range expected {
				if want == "" {
					continue
				}
				if record == want {
					matched = true
					break
				}
			}
			if !matched {
				return false
			}
		}
	}
	for skillID, acquire := range req.Skill {
		if skillID == 0 {
			continue
		}
		mustAcquire := acquire > 0
		var wzSkill *wz.Skill
		if qc.owner.gameWorld != nil {
			resources := qc.owner.gameWorld.GetResources()
			if resources != nil {
				wzSkill = resources.GetSkill(skillID)
			}
		}
		entry := qc.owner.skills.Learned(skillID)
		skillLevel := 0
		masterLevel := 0
		if entry != nil {
			skillLevel = entry.Level()
			masterLevel = entry.masterLevel
		}
		if mustAcquire {
			ok := false
			if wzSkill != nil && wzSkill.FourthClass {
				ok = masterLevel > 0
			} else {
				ok = skillLevel > 0
			}
			if !ok {
				return false
			}
		} else if skillLevel != 0 || masterLevel != 0 {
			return false
		}
	}
	if req.QuestComplete > 0 && qc.completedQuestCount() < req.QuestComplete {
		return false
	}
	if req.PartyQuestS != 0 {
		need := req.PartyQuestS
		if need <= 0 {
			need = 5
		}
		partyQuestIDs := []uint32{1200, 1201, 1202, 1203, 1204, 1205, 1206, 1300, 1301}
		sRankings := 0
		for _, questID := range partyQuestIDs {
			qp := qc.Get(questID)
			if qp == nil {
				continue
			}
			rank, ok := qp.RecordExField("rank")
			if ok && rank == "S" {
				sRankings++
			}
		}
		if sRankings < need {
			return false
		}
	}
	if len(req.Pet) > 0 && (qc.owner.pets.active == nil || !slices.Contains(req.Pet, qc.owner.pets.active.item.GetModel().GetID())) {
		return false
	}
	if req.PetTamenessMin > 0 && (qc.owner.pets.active == nil || int(qc.owner.pets.active.item.closeness) < req.PetTamenessMin) {
		return false
	}
	if req.MBMin > 0 || len(req.MBCard) > 0 || req.SubClassFlags != 0 {
		return false
	}
	if req.TamingMobLevelMin > 0 && (qc.owner.inventory.equipped[constant.EquipmentPartsTamingMob] == nil || int(qc.owner.mount.level) < req.TamingMobLevelMin) {
		return false
	}
	if req.EndMeso > 0 || req.EquipAllNeed > 0 || req.EquipSelectNeed > 0 {
		return false
	}
	if req.Buff != "" || req.ExceptBuff != "" || req.WorldMin != "" || req.WorldMax != "" {
		return false
	}
	if req.Premium {
		return false
	}

	return true
}

func sameCalendarDay(a, b time.Time) bool {
	ay, am, ad := a.In(time.Local).Date()
	by, bm, bd := b.In(time.Local).Date()
	return ay == by && am == bm && ad == bd
}
