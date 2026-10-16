package game

import "github.com/boyism80/fm/protocol/dto"

func (qp *Quest) BuildS2C() (dto.QuestStatus, bool) {
	if qp == nil || qp.questID == 0 || qp.questID > 0xFFFF {
		return dto.QuestStatus{}, false
	}
	return dto.QuestStatus{
		QuestID:        uint16(qp.questID),
		Status:         uint8(qp.status),
		MobKills:       qp.StartedMobKills(),
		Deadline:       qp.deadline,
		StatusRecord:   qp.statusRecord.AsString(),
		CompletionTime: qp.completionTime,
	}, true
}

func (qp *Quest) TimeLimit() int {
	if qp == nil || qp.wz == nil {
		return 0
	}
	return qp.wz.Meta.TimeLimit
}
