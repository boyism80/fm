"use strict";

const PARTY_EVENTS_EXCHANGE = "fm.party";

class PartyEventPublisher {
    constructor(rabbitmqService) {
        this.rabbitmqService = rabbitmqService;
    }

    async initialize() {
        await this.rabbitmqService.assertDirectExchange(PARTY_EVENTS_EXCHANGE);
    }

    async publish(eventType, worldId, partyId, revision, extraPayload = {}) {
        const routingKey = `fm.${worldId}.party.all`;
        const payload = {
            event_id: extraPayload.event_id ?? `${Date.now()}-${Math.random().toString(36).slice(2, 8)}`,
            event_type: eventType,
            party_id: partyId,
            revision: Number(revision),
            occurred_at: new Date().toISOString(),
            ...extraPayload,
        };
        return this.rabbitmqService.publish(PARTY_EVENTS_EXCHANGE, routingKey, payload);
    }

    async publishToGameChannel(worldId, channelId, payload) {
        const routingKey = `fm.${worldId}.party.game.${channelId}`;
        const body = {
            occurred_at: new Date().toISOString(),
            ...payload,
        };
        return this.rabbitmqService.publish(PARTY_EVENTS_EXCHANGE, routingKey, body);
    }
}

module.exports = {
    PartyEventPublisher,
    PARTY_EVENTS_EXCHANGE,
};
