package server

import (
	"encoding/json"
	"log"

	"github.com/boyism80/fm/services/game"
)

func decodeGuildEventEnvelope(raw json.RawMessage) (game.GuildEventEnvelope, bool) {
	var evt game.GuildEventEnvelope
	if err := json.Unmarshal(raw, &evt); err != nil {
		log.Printf("guild consumer: invalid envelope: %v", err)
		return evt, false
	}
	return evt, true
}
