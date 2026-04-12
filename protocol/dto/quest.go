package dto

import (
	"time"
)

type QuestStatus struct {
	QuestID        uint16
	Status         uint8
	MobKills       []uint16
	CustomData     string
	CompletionTime time.Time
}
