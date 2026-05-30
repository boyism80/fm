package server

import (
	"encoding/json"
	"log"
)

func decodeGuildEventEnvelope(raw json.RawMessage) (GuildEventEnvelope, bool) {
	var evt GuildEventEnvelope
	if err := json.Unmarshal(raw, &evt); err != nil {
		log.Printf("guild consumer: invalid envelope: %v", err)
		return evt, false
	}
	return evt, true
}
