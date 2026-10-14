package server

import (
	"encoding/json"
	"log"

	"github.com/boyism80/fm/services/game/domain"
)

func decodeAllianceEventEnvelope(raw json.RawMessage) (domain.AllianceEventEnvelope, bool) {
	var evt domain.AllianceEventEnvelope
	if err := json.Unmarshal(raw, &evt); err != nil {
		log.Printf("alliance consumer: invalid envelope: %v", err)
		return evt, false
	}
	return evt, true
}
