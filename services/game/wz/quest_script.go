package wz

func (q *Quest) HasStartScript() bool {
	if q == nil {
		return false
	}
	return phaseHasRequirementKind(q.Start.Requirements, QuestReqStartScript)
}

func (q *Quest) HasEndScript() bool {
	if q == nil {
		return false
	}
	return phaseHasRequirementKind(q.Complete.Requirements, QuestReqEndScript)
}

func (q *Quest) NextQuestID() uint32 {
	if q == nil {
		return 0
	}
	for _, act := range q.Complete.Actions {
		if act.Kind == QuestActNextQuest && act.IntValue > 0 {
			return uint32(act.IntValue)
		}
	}
	return 0
}

func phaseHasRequirementKind(reqs []QuestRequirement, kind QuestRequirementKind) bool {
	for _, req := range reqs {
		if req.Kind == kind {
			return true
		}
	}
	return false
}
