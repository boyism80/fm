package entity

type QuestStatus struct {
	Quest          *Quest
	Status         uint8 // 1 = Started, 2 = Completed
	MobKills       map[int]int
	CustomData     string
	CompletionTime int64
}

type Quest struct {
	Id int
}

func (q *QuestStatus) HasMobKills() bool {
	return len(q.MobKills) > 0
}
