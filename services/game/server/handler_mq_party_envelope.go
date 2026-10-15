package server

import (
	"encoding/json"

	"github.com/boyism80/fm/services/game"
)

func decodePartyEventEnvelope(raw json.RawMessage) (game.PartyEventEnvelope, bool) {
	var evt game.PartyEventEnvelope
	if err := json.Unmarshal(raw, &evt); err != nil {
		return evt, false
	}
	return evt, true
}
