package server

import (
	"encoding/json"
	"log"

	"github.com/asynkron/protoactor-go/actor"
	amqp "github.com/rabbitmq/amqp091-go"
)

type allianceMqNoticeChanged struct{ gs *GameServer }

func (allianceMqNoticeChanged) New(gs *GameServer) *allianceMqNoticeChanged {
	return &allianceMqNoticeChanged{gs: gs}
}

func (*allianceMqNoticeChanged) EventType() string {
	return "notice_changed"
}

func (h *allianceMqNoticeChanged) Handle(_ actor.Context, _ amqp.Delivery, _ string, raw json.RawMessage) error {
	gs := h.gs
	if gs == nil {
		return nil
	}
	evt, ok := decodeAllianceEventEnvelope(raw)
	if !ok {
		return nil
	}
	alliancePb, ok := allianceFromMQPayload(raw)
	if !ok {
		log.Printf("alliance consumer: notice_changed alliance_id=%d missing alliance_pb", evt.AllianceID)
		return nil
	}
	var extra struct {
		Notice string `json:"notice"`
	}
	_ = json.Unmarshal(raw, &extra)
	gs.alliance.Update(alliancePb)
	gs.alliance.BroadcastNoticeChanged(alliancePb, extra.Notice)
	log.Printf("alliance consumer: applied notice_changed alliance_id=%d revision=%d", evt.AllianceID, evt.Revision)
	return nil
}
