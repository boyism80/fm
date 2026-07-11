package entity

import (
	"strconv"
	"time"

	"github.com/boyism80/fm/core/clock"
	"github.com/boyism80/fm/services/game/wz"
)

func requirementMet(req wz.QuestRequirement, qc *QuestContainer, qp *Quest, phase wz.QuestPhase, opts QuestPhaseOpts) bool {
	if qc == nil || qc.owner == nil {
		return false
	}
	ch := qc.owner
	switch req.Kind {
	case wz.QuestReqNPC:
		if opts.NpcID == nil {
			return true
		}
		required := uint32(req.IntValue)
		if required == 0 {
			return true
		}
		wireNPC := *opts.NpcID
		if wireNPC != 0 && wireNPC != required {
			return false
		}
		mapInst := ch.GetMap()
		if mapInst == nil {
			return false
		}
		for _, obj := range mapInst.GetNpcs() {
			npc, ok := obj.(*Npc)
			if !ok || npc.Wz == nil || npc.Wz.BaseSpawn == nil {
				continue
			}
			if npc.Wz.ID == required {
				return true
			}
		}
		return false
	case wz.QuestReqLvMin:
		return int(ch.GetLevel()) >= req.IntValue
	case wz.QuestReqLvMax:
		return int(ch.GetLevel()) <= req.IntValue
	case wz.QuestReqClass:
		return matchesQuestClass(ch.Class, req.Classes)
	case wz.QuestReqItem:
		for itemID, count := range req.Items {
			if !ch.HasItemCount(itemID, uint16(count)) {
				return false
			}
		}
		return true
	case wz.QuestReqMob:
		return qp.MeetsMobCounts(req.Mobs)
	case wz.QuestReqQuest:
		for questID, state := range req.Quests {
			if !qc.Get(questID).MatchesState(state) {
				return false
			}
		}
		return true
	case wz.QuestReqStartScript, wz.QuestReqEndScript:
		return true
	case wz.QuestReqPop:
		return int(ch.population) >= req.IntValue
	case wz.QuestReqFieldEnter:
		mapID := req.IntValue
		if mapID <= 0 {
			return true
		}
		mapInst := ch.GetMap()
		if mapInst == nil {
			return false
		}
		return int(mapInst.GetMapID()) == mapID
	case wz.QuestReqNormalAutoStart:
		return true
	case wz.QuestReqInterval:
		if qp == nil || qp.Status != QuestStatusCompleted {
			return true
		}
		if qp.CompletionTime.IsZero() {
			return true
		}
		minutes := req.IntValue
		if minutes <= 0 {
			return true
		}
		return !clock.Now().Before(qp.CompletionTime.Add(time.Duration(minutes) * time.Minute))
	case wz.QuestReqDayByDay:
		if qp == nil || qp.Status != QuestStatusCompleted {
			return true
		}
		if qp.CompletionTime.IsZero() {
			return true
		}
		return !sameCalendarDay(qp.CompletionTime, clock.Now())
	case wz.QuestReqTimeStart, wz.QuestReqTimeEnd:
		return meetsEventTimeRequirement(req.Kind, req.StrValue)
	case wz.QuestReqInfoNumber:
		if qp == nil {
			return false
		}
		refID := req.IntValue
		if refID <= 0 {
			return true
		}
		if qp.container == nil {
			return false
		}
		refQP := qp.container.Get(uint32(refID))
		if refQP == nil || !refQP.IsStarted() {
			return false
		}
		expected := infoStringsFromPhase(phase)
		if len(expected) == 0 {
			return !refQP.StatusRecord.IsEmpty()
		}
		record := refQP.StatusRecord.AsString()
		for _, want := range expected {
			if want == "" {
				continue
			}
			if record == want {
				return true
			}
		}
		return false
	case wz.QuestReqInfo:
		return true
	case wz.QuestReqSkill:
		for skillID, acquire := range req.Skills {
			if !ch.meetsQuestSkillRequirement(skillID, acquire) {
				return false
			}
		}
		return true
	case wz.QuestReqQuestComplete:
		return qc.completedQuestCount() >= req.IntValue
	case wz.QuestReqPartyQuestS:
		need := req.IntValue
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
		return sRankings >= need
	case wz.QuestReqPet,
		wz.QuestReqPetTamenessMin, wz.QuestReqMBMin, wz.QuestReqMBCard,
		wz.QuestReqSubClassFlags:
		return false
	default:
		if req.IntValue != 0 || req.StrValue != "" || len(req.Items) > 0 {
			return false
		}
		return true
	}
}

func infoStringsFromPhase(phase wz.QuestPhase) []string {
	for _, req := range phase.Requirements {
		if req.Kind != wz.QuestReqInfo {
			continue
		}
		if len(req.InfoStrings) > 0 {
			return req.InfoStrings
		}
	}
	return nil
}

func phaseRequirementsMet(phase wz.QuestPhase, qc *QuestContainer, qp *Quest, opts QuestPhaseOpts) bool {
	if qc == nil || qc.owner == nil {
		return false
	}
	for _, req := range phase.Requirements {
		if !requirementMet(req, qc, qp, phase, opts) {
			return false
		}
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

func meetsEventTimeRequirement(kind wz.QuestRequirementKind, raw string) bool {
	eventTime, ok := parseQuestEventTime(raw)
	if !ok {
		return false
	}
	now := clock.Now()
	if kind == wz.QuestReqTimeStart {
		return !now.Before(eventTime)
	}
	if kind == wz.QuestReqTimeEnd {
		return now.Before(eventTime)
	}
	return false
}
