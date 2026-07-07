package entity

func (qc *QuestContainer) ForceStart(questID uint32, record string) (*Quest, error) {
	if qc == nil || questID == 0 {
		return nil, ErrQuestNotStartable
	}
	if qc.questDef(questID) != nil {
		qp, err := qc.Start(questID, QuestPrepareOpts{Force: true})
		if err != nil {
			return nil, err
		}
		if record != "" {
			qp.StatusRecord.WriteString(record)
			if qp.WiresToClient() && qc.owner != nil && qc.owner.Listener != nil {
				qc.owner.Listener.OnQuestProgress(qc.owner, qp)
			}
		}
		return qp, nil
	}

	qp := qc.Get(questID)
	if qp == nil {
		qp = qc.Create(questID, QuestStatusStarted)
		if qp == nil {
			qp = qc.Get(questID)
		}
	}
	if qp == nil {
		return nil, ErrQuestInvalidState
	}
	qp.Status = QuestStatusStarted
	qp.StatusRecord.WriteString(record)
	qp.ClearDeadline()
	if qp.MobKills == nil {
		qp.MobKills = make(map[uint32]int)
	}
	qc.RunAutoTriggers(nil, AutoQuestTriggerLogin, 0)
	qc.RunAutoTriggers(nil, AutoQuestTriggerLevelUp, 0)
	return qp, nil
}
