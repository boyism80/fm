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

type partyMqLogOnOff struct{ gs *GameServer }

func (partyMqLogOnOff) New(gs *GameServer) *partyMqLogOnOff { return &partyMqLogOnOff{gs: gs} }
func (*partyMqLogOnOff) EventType() string                  { return "log_onoff" }
func (h *partyMqLogOnOff) Handle(ctx actor.Context, _ amqp.Delivery, _ string, raw json.RawMessage) error {
	gs := h.gs
	if gs == nil || gs.party == nil {
		return nil
	}
	pc := gs.party
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
							pc.DeliverPartyLogOnOff(party, payload.CharacterID)
						}
					}
					return nil
				}
			}
		}
	}
	pc.UpdateAsync(ctx, evt).
		Then(func() (interface{}, error) { return nil, nil }, func(interface{}) error {
			if raw == nil {
				return nil
			}
			var extra struct {
				CharacterID uint32 `json:"character_id"`
			}
			if err := json.Unmarshal(raw, &extra); err != nil || extra.CharacterID == 0 {
				return nil
			}
			party := pc.Get(evt.PartyID)
			if party != nil {
				pc.DeliverPartyLogOnOff(party, extra.CharacterID)
			}
			return nil
		}).
		Run()
	return nil
}
