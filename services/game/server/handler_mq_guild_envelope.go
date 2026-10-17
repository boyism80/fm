package server

import (
	"encoding/base64"
	"encoding/json"
	"log"

	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
	"github.com/boyism80/fm/services/game"
	"google.golang.org/protobuf/proto"
)

func decodeGuildEventEnvelope(raw json.RawMessage) (game.GuildEventEnvelope, bool) {
	var evt game.GuildEventEnvelope
	if err := json.Unmarshal(raw, &evt); err != nil {
		log.Printf("guild consumer: invalid envelope: %v", err)
		return evt, false
	}
	return evt, true
}

func decodeGuildSnapshot(encoded string) *internal.Guild {
	if encoded == "" {
		return nil
	}
	wire, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return nil
	}
	var guildPb internal.Guild
	if err := proto.Unmarshal(wire, &guildPb); err != nil {
		return nil
	}
	return &guildPb
}
