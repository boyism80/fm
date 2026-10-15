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
		q.Status = status
		q.MobKills = make(map[uint32]int, len(pb.GetMobKills()))
		for mobID, kills := range pb.GetMobKills() {
			q.MobKills[mobID] = int(kills)
		}
		q.StatusRecord.WriteString(pb.GetStatusRecord())
		if ms := pb.GetDeadlineUnixMs(); ms > 0 {
			q.Deadline = time.UnixMilli(ms)
		} else {
			q.Deadline = time.Time{}
		}
		if ms := pb.GetStartTimeUnixMs(); ms > 0 {
			q.StartTime = time.UnixMilli(ms)
		} else {
			q.StartTime = time.Time{}
		}
		q.RecordEx = make(map[string]string, len(pb.GetRecordEx()))
		for key, value := range pb.GetRecordEx() {
			q.RecordEx[key] = value
		}
		if ms := pb.GetCompletionTimeUnixMs(); ms > 0 {
			q.CompletionTime = time.UnixMilli(ms)
		} else {
			q.CompletionTime = time.Time{}
		}
		q.Forfeited = int(pb.GetForfeited())
	}
}

func (qc *Quests) BuildRPC() []*internal.Quest {
	out := make([]*internal.Quest, 0)
	qc.ForEach(func(questID uint32, q *Quest) {
		if q == nil {
			return
		}
		if q.Status != QuestStatusStarted && q.Status != QuestStatusCompleted && q.Forfeited == 0 && len(q.RecordEx) == 0 {
			return
		}
		mobKills := make(map[uint32]uint32, len(q.MobKills))
		for mobID, kills := range q.MobKills {
			if kills < 0 {
				continue
			}
			mobKills[mobID] = uint32(kills)
		}
		recordEx := make(map[string]string, len(q.RecordEx))
		for key, value := range q.RecordEx {
			recordEx[key] = value
		}
		var completionTimeUnixMs int64
		if !q.CompletionTime.IsZero() {
			completionTimeUnixMs = q.CompletionTime.UnixMilli()
		}
		var deadlineUnixMs int64
		if !q.Deadline.IsZero() {
			deadlineUnixMs = q.Deadline.UnixMilli()
		}
		var startTimeUnixMs int64
		if !q.StartTime.IsZero() {
			startTimeUnixMs = q.StartTime.UnixMilli()
		}
		out = append(out, &internal.Quest{
			CharacterId:          qc.owner.GetID(),
			QuestId:              questID,
			Status:               uint32(q.Status),
			MobKills:             mobKills,
			StatusRecord:         q.StatusRecord.AsString(),
			RecordEx:             recordEx,
			CompletionTimeUnixMs: completionTimeUnixMs,
			DeadlineUnixMs:       deadlineUnixMs,
			StartTimeUnixMs:      startTimeUnixMs,
			Forfeited:            uint32(q.Forfeited),
		})
	})
	return out
}
