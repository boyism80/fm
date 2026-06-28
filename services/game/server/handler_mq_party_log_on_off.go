package server

import (
	"encoding/base64"
	"encoding/json"
	"log"

	"github.com/asynkron/protoactor-go/actor"
	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
	amqp "github.com/rabbitmq/amqp091-go"
	"google.golang.org/protobuf/proto"
)

type partyMqLogOnOff struct {
	partyMqHandler
}

func (partyMqLogOnOff) New(gs *GameServer) *partyMqLogOnOff {
	return &partyMqLogOnOff{partyMqHandler: partyMqHandler{gs: gs}}
}

func (*partyMqLogOnOff) EventType() string {
	return "log_onoff"
}

func (h *partyMqLogOnOff) Handle(ctx actor.Context, _ amqp.Delivery, _ string, raw json.RawMessage) error {
	if h.gs == nil {
		return nil
	}
	pc := h.gs.party
	evt, ok := decodePartyEventEnvelope(raw)
	if !ok {
		return nil
	}
	if raw != nil {
		var payload struct {
			PartyPB     string `json:"party_pb"`
			CharacterID uint32 `json:"character_id"`
		}
		if err := json.Unmarshal(raw, &payload); err == nil && payload.PartyPB != "" {
			wire, err := base64.StdEncoding.DecodeString(payload.PartyPB)
			if err == nil {
				var partyPb internal.Party
				if err := proto.Unmarshal(wire, &partyPb); err == nil {
					applied, err := pc.applyEmbeddedParty(evt, &partyPb, false)
					if err != nil {
						log.Printf("party consumer: log_onoff embedded party: %v", err)
						return nil
					}
					if applied && payload.CharacterID != 0 {
						if party := pc.Get(evt.PartyID); party != nil {
							h.sendLogOnOff(party)
						}
					}
					return nil
				}
			}
		}
	}
	pc.UpdateAsync(ctx, evt).Then(func(interface{}) (interface{}, error) {
		if raw == nil {
			return nil, nil
		}
		var extra struct {
			CharacterID uint32 `json:"character_id"`
		}
		if err := json.Unmarshal(raw, &extra); err != nil || extra.CharacterID == 0 {
			return nil, nil
		}
		party := pc.Get(evt.PartyID)
		if party != nil {
			h.sendLogOnOff(party)
		}
		return nil, nil
	})
	return nil
}
