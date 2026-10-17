package server

import (
	"encoding/base64"
	"encoding/json"

	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
	"github.com/boyism80/fm/services/game"
	"google.golang.org/protobuf/proto"
)

func decodePartyEventEnvelope(raw json.RawMessage) (game.PartyEventEnvelope, bool) {
	var evt game.PartyEventEnvelope
	if err := json.Unmarshal(raw, &evt); err != nil {
		return evt, false
	}
	return evt, true
}

func decodePartySnapshot(encoded string) *internal.Party {
	if encoded == "" {
		return nil
	}
	wire, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return nil
	}
	var partyPb internal.Party
	if err := proto.Unmarshal(wire, &partyPb); err != nil {
		return nil
	}
	return &partyPb
}
