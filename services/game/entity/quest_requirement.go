package entity

import (
	"strconv"
	"time"

	"github.com/boyism80/fm/core/clock"
	"github.com/boyism80/fm/services/game/wz"
)

func phaseRequirementsMet(phase wz.QuestPhase, qc *QuestContainer, qp *Quest, opts QuestPhaseOpts) bool {
	if qc == nil || qc.owner == nil {
		return false
	}
	return requirementsMet(phase.Requirements, qc, qp, opts)
}

func requirementsMet(req wz.QuestRequirements, qc *QuestContainer, qp *Quest, opts QuestPhaseOpts) bool {
	if qc == nil || qc.owner == nil {
		return false
	}
	ch := qc.owner

	if req.NPC != 0 && opts.NpcID != nil {
		wireNPC := *opts.NpcID
		if wireNPC != 0 && wireNPC != req.NPC {
			return false
		}
		mapInst := ch.GetMap()
		if mapInst == nil {
			return false
		}
		found := false
		for _, obj := range mapInst.GetNpcs() {
			npc, ok := obj.(*Npc)
			if !ok || npc.Wz == nil || npc.Wz.BaseSpawn == nil {
				continue
			}
			if npc.Wz.ID == req.NPC {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}

	if req.LevelMin > 0 && int(ch.GetLevel()) < req.LevelMin {
		return false
	}
	if req.LevelMax > 0 && int(ch.GetLevel()) > req.LevelMax {
		return false
	}
	if req.Level > 0 && int(ch.GetLevel()) < req.Level {
		return false
	}
	if len(req.Job) > 0 && !matchesQuestClass(ch.Class, req.Job) {
		return false
	}
	for itemID, count := range req.Item {
		if !ch.HasItemCount(itemID, uint16(count)) {
			return false
		}
	}
	if len(req.Mob) > 0 && !qp.MeetsMobCounts(req.Mob) {
		return false
	}
	for questID, state := range req.Quest {
		if !qc.Get(questID).MatchesState(state) {
			return false
		}
	}
	if req.Pop > 0 && int(ch.population) < req.Pop {
		return false
	}
	if req.FieldEnter > 0 {
		mapInst := ch.GetMap()
		if mapInst == nil || int(mapInst.GetMapID()) != req.FieldEnter {
			return false
		}
	}
	if req.HasInterval {
		if qp != nil && qp.Status == QuestStatusCompleted && !qp.CompletionTime.IsZero() {
			if req.Interval > 0 && clock.Now().Before(qp.CompletionTime.Add(time.Duration(req.Interval)*time.Minute)) {
				return false
			}
		}
	}
	if req.DayByDay {
		if qp != nil && qp.Status == QuestStatusCompleted && !qp.CompletionTime.IsZero() {
			if sameCalendarDay(qp.CompletionTime, clock.Now()) {
				return false
			}
		}
	}
	if req.Start != "" && !meetsEventTimeBound(req.Start, true) {
		return false
	}
	if req.End != "" && !meetsEventTimeBound(req.End, false) {
		return false
	}
	if req.InfoNumber > 0 {
		if qp == nil || qp.container == nil {
			return false
		}
		refQP := qp.container.Get(uint32(req.InfoNumber))
		if refQP == nil || !refQP.IsStarted() {
			return false
		}
		expected := req.Info
		if len(expected) == 0 {
			if refQP.StatusRecord.IsEmpty() {
				return false
			}
		} else {
			record := refQP.StatusRecord.AsString()
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
		if !ch.meetsQuestSkillRequirement(skillID, acquire) {
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
		partyQuestIDs := []uint32{1200, 1201, 1202, 1203, 1204, 1205, 1206, 1300, 1301, 1302}
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
	if len(req.Pet) > 0 || req.PetTamenessMin > 0 || req.MBMin > 0 || len(req.MBCard) > 0 || req.SubJobFlags != 0 {
		return false
	}
	if req.EndMeso > 0 || req.EquipAllNeed > 0 || req.EquipSelectNeed > 0 || req.TamingMobLevelMin > 0 {
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

func parseQuestEventTime(raw string) (time.Time, bool) {
	if len(raw) != 10 {
		return time.Time{}, false
	}
	year, err := strconv.Atoi(raw[0:4])
	if err != nil {
		return time.Time{}, false
	}
	month, err := strconv.Atoi(raw[4:6])
	if err != nil {
		return time.Time{}, false
	}
	day, err := strconv.Atoi(raw[6:8])
	if err != nil {
		return time.Time{}, false
	}
	hour, err := strconv.Atoi(raw[8:10])
	if err != nil {
		return time.Time{}, false
	}
	if month < 1 || month > 12 || day < 1 || day > 31 || hour < 0 || hour > 23 {
		return time.Time{}, false
	}
	return time.Date(year, time.Month(month), day, hour, 0, 0, 0, time.Local), true
}

func sameCalendarDay(a, b time.Time) bool {
	ay, am, ad := a.In(time.Local).Date()
	by, bm, bd := b.In(time.Local).Date()
	return ay == by && am == bm && ad == bd
}

func meetsEventTimeBound(raw string, isStart bool) bool {
	eventTime, ok := parseQuestEventTime(raw)
	if !ok {
		return false
	}
	now := clock.Now()
	if isStart {
		return !now.Before(eventTime)
	}
	return now.Before(eventTime)
}
