package entity

import "time"

type QuestStatus struct {
	Quest          *Quest
	Status         uint8 // 1 = Started, 2 = Completed
	MobKills       map[int]int
	CustomData     string
	CompletionTime time.Time
}

type Quest struct {
	ID int
}

func (q *QuestStatus) HasMobKills() bool {
	return len(q.MobKills) > 0
}
