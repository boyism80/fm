package entity

import (
	"github.com/boyism80/fm/services/game/wz"
)

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

func (qp *Quest) meetsInfoNumberRequirement(refID int, phase wz.QuestPhase) bool {
	if refID <= 0 {
		return true
	}
	if qp == nil || qp.container == nil {
		return false
	}
	refQP := qp.container.Get(uint32(refID))
	if refQP == nil || !refQP.IsStarted() {
		return false
	}
	expected := infoStringsFromPhase(phase)
	if len(expected) == 0 {
		return refQP.StatusRecord != ""
	}
	record := refQP.StatusRecord
	for _, want := range expected {
		if want == "" {
			continue
		}
		if record == want {
			return true
		}
	}
	return false
}

func (qc *QuestContainer) applyInfoNumberAction(refID uint32) {
	if qc == nil || refID == 0 {
		return
	}
	refQP := qc.Get(refID)
	if refQP == nil || !refQP.IsStarted() {
		return
	}
	refQP.Status = QuestStatusCompleted
	refQP.CompletionTime = questRequirementNow()
}

func (qp *Quest) meetsDayByDayRequirement() bool {
	if qp == nil {
		return false
	}
	if qp.Status != QuestStatusCompleted {
		return true
	}
	if qp.CompletionTime.IsZero() {
		return true
	}
	return !questSameCalendarDay(qp.CompletionTime, questRequirementNow())
}

func (qp *Quest) meetsQuestEventTimeRequirement(kind wz.QuestRequirementKind, raw string) bool {
	eventTime, ok := parseQuestEventTime(raw)
	if !ok {
		return false
	}
	now := questRequirementNow()
	if kind == wz.QuestReqTimeStart {
		return !now.Before(eventTime)
	}
	if kind == wz.QuestReqTimeEnd {
		return now.Before(eventTime)
	}
	return false
}
