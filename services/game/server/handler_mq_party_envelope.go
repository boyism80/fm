package server

import (
	"encoding/json"

	"github.com/boyism80/fm/services/game/domain"
)

func decodePartyEventEnvelope(raw json.RawMessage) (domain.PartyEventEnvelope, bool) {
	var evt domain.PartyEventEnvelope
	if err := json.Unmarshal(raw, &evt); err != nil {
		return evt, false
	}
	return evt, true
}
