package dto

import (
	"time"
)

// QuestStatus represents quest status data for protocol
type QuestStatus struct {
	QuestID   uint16
	Status    uint8 // 1 = Started, 2 = Completed
	MobKills  []uint16
	CustomData string
	CompletionTime time.Time
}

