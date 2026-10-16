package game

import (
	"time"

	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
)

func (qc *Quests) Load(pbs []*internal.Quest) {
	for _, pb := range pbs {
		if pb == nil {
			continue
		}
		questID := pb.GetQuestId()
		if questID == 0 {
			continue
		}
		status := QuestStatusType(pb.GetStatus())
		hasRecordEx := len(pb.GetRecordEx()) > 0
		if status != QuestStatusStarted && status != QuestStatusCompleted && pb.GetForfeited() == 0 && !hasRecordEx {
			continue
		}
		q := qc.Get(questID)
		if q == nil {
			q = qc.Create(questID, status)
		}
		if q == nil {
			continue
		}
		q.status = status
		q.mobKills = make(map[uint32]int, len(pb.GetMobKills()))
		for mobID, kills := range pb.GetMobKills() {
			q.mobKills[mobID] = int(kills)
		}
		q.statusRecord.WriteString(pb.GetStatusRecord())
		if ms := pb.GetDeadlineUnixMs(); ms > 0 {
			q.deadline = time.UnixMilli(ms)
		} else {
			q.deadline = time.Time{}
		}
		if ms := pb.GetStartTimeUnixMs(); ms > 0 {
			q.startTime = time.UnixMilli(ms)
		} else {
			q.startTime = time.Time{}
		}
		q.recordEx = make(map[string]string, len(pb.GetRecordEx()))
		for key, value := range pb.GetRecordEx() {
			q.recordEx[key] = value
		}
		if ms := pb.GetCompletionTimeUnixMs(); ms > 0 {
			q.completionTime = time.UnixMilli(ms)
		} else {
			q.completionTime = time.Time{}
		}
		q.forfeited = int(pb.GetForfeited())
	}
}

func (qc *Quests) BuildRPC() []*internal.Quest {
	out := make([]*internal.Quest, 0)
	qc.ForEach(func(questID uint32, q *Quest) {
		if q == nil {
			return
		}
		if q.status != QuestStatusStarted && q.status != QuestStatusCompleted && q.forfeited == 0 && len(q.recordEx) == 0 {
			return
		}
		mobKills := make(map[uint32]uint32, len(q.mobKills))
		for mobID, kills := range q.mobKills {
			if kills < 0 {
				continue
			}
			mobKills[mobID] = uint32(kills)
		}
		recordEx := make(map[string]string, len(q.recordEx))
		for key, value := range q.recordEx {
			recordEx[key] = value
		}
		var completionTimeUnixMs int64
		if !q.completionTime.IsZero() {
			completionTimeUnixMs = q.completionTime.UnixMilli()
		}
		var deadlineUnixMs int64
		if !q.deadline.IsZero() {
			deadlineUnixMs = q.deadline.UnixMilli()
		}
		var startTimeUnixMs int64
		if !q.startTime.IsZero() {
			startTimeUnixMs = q.startTime.UnixMilli()
		}
		out = append(out, &internal.Quest{
			CharacterId:          qc.owner.GetID(),
			QuestId:              questID,
			Status:               uint32(q.status),
			MobKills:             mobKills,
			StatusRecord:         q.statusRecord.AsString(),
			RecordEx:             recordEx,
			CompletionTimeUnixMs: completionTimeUnixMs,
			DeadlineUnixMs:       deadlineUnixMs,
			StartTimeUnixMs:      startTimeUnixMs,
			Forfeited:            uint32(q.forfeited),
		})
	})
	return out
}
