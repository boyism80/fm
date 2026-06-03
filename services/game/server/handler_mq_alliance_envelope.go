package server

import (
	"encoding/json"
	"log"
)

type AllianceEventEnvelope struct {
	EventID    string `json:"event_id"`
	EventType  string `json:"event_type"`
	WorldID    uint32 `json:"world_id"`
	AllianceID uint32 `json:"alliance_id"`
	Revision   uint64 `json:"revision"`
	OccurredAt string `json:"occurred_at"`
}

func decodeAllianceEventEnvelope(raw json.RawMessage) (AllianceEventEnvelope, bool) {
	var evt AllianceEventEnvelope
	if err := json.Unmarshal(raw, &evt); err != nil {
		log.Printf("alliance consumer: invalid envelope: %v", err)
		return evt, false
	}
	return evt, true
}
