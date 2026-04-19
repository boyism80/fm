package server

import (
	"encoding/json"
	"log"
)

func decodePartyEventEnvelope(raw json.RawMessage) (PartyEventEnvelope, bool) {
	var evt PartyEventEnvelope
	if err := json.Unmarshal(raw, &evt); err != nil {
		log.Printf("party consumer: invalid envelope: %v", err)
		return evt, false
	}
	return evt, true
}
