package server

import (
	"encoding/json"
)

func decodePartyEventEnvelope(raw json.RawMessage) (PartyEventEnvelope, bool) {
	var evt PartyEventEnvelope
	if err := json.Unmarshal(raw, &evt); err != nil {
		return evt, false
	}
	return evt, true
}
