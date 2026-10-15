package server

import (
	"encoding/json"
	"log"

	"github.com/boyism80/fm/services/game"
)

func decodeAllianceEventEnvelope(raw json.RawMessage) (game.AllianceEventEnvelope, bool) {
	var evt game.AllianceEventEnvelope
	if err := json.Unmarshal(raw, &evt); err != nil {
		log.Printf("alliance consumer: invalid envelope: %v", err)
		return evt, false
	}
	return evt, true
}
