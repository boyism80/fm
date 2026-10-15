package wz

func (q *Quest) HasStartScript() bool {
	if q == nil {
		return false
	}
	return q.Start.Requirements.StartScript != ""
}

func (q *Quest) HasEndScript() bool {
	if q == nil {
		return false
	}
	return q.Complete.Requirements.EndScript != ""
}

func (q *Quest) NextQuestID() uint32 {
	if q == nil {
		return 0
	}
	return q.Complete.Actions.NextQuest
}
