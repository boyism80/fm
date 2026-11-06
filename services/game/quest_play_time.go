package game

import (
	"fmt"
	"time"
)

const (
	// playTimeDays is the days of play a daily play time quest asks for; WZ has no value for it.
	playTimeDays = 27
	// playTimeShownAfter is the seconds of today's play before the original server sent the record.
	playTimeShownAfter = 257
)

// CountPlayTime adds play time to every started quest that counts daily play time. A day counts once today's play
// reaches the quest's dailyPlayTime, and the day's tally starts over at midnight.
func (qc *Quests) CountPlayTime(d time.Duration) {
	daily := RecordReset{kind: RecordResetDaily}
	for id, qp := range qc.progress {
		if !qp.IsStarted() || qp.wz == nil || qp.wz.Meta.DailyPlayTime <= 0 {
			continue
		}
		records := qc.owner.records
		doneKey := fmt.Sprintf("quest.%d.play_time.done", id)
		if records.Get(doneKey) != 0 {
			continue
		}
		today := records.Add(fmt.Sprintf("quest.%d.play_time.today", id), int64(d/time.Second), daily)
		daysKey := fmt.Sprintf("quest.%d.play_time.days", id)
		days := records.Get(daysKey)
		if today >= int64(qp.wz.Meta.DailyPlayTime) {
			records.Set(doneKey, 1, daily)
			days = records.Add(daysKey, 1, RecordReset{})
			if days == playTimeDays {
				qc.owner.listener.OnShowQuestCompletion(qc.owner, id)
			}
		}
		if today >= playTimeShownAfter {
			qc.owner.listener.OnQuestPlayTime(qc.owner, id, uint16(min(today, 0xFFFF)), uint8(min(days, 0xFF)))
		}
	}
}
