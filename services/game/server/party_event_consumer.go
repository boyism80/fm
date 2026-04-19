package server

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/boyism80/fm/common/config"
	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
	amqp "github.com/rabbitmq/amqp091-go"
	"google.golang.org/protobuf/proto"
)

const partyEventsExchange = "fm.party"

type PartyEventEnvelope struct {
	EventID    string `json:"event_id"`
	EventType  string `json:"event_type"`
	WorldID    uint32 `json:"world_id"`
	PartyID    uint32 `json:"party_id"`
	Revision   uint64 `json:"revision"`
	OccurredAt string `json:"occurred_at"`
}

type PartyEventHandler interface {
	EventType() string
	Handle(evt PartyEventEnvelope, raw json.RawMessage) error
}

type PartyEventDispatcher struct {
	mu       sync.RWMutex
	handlers map[string]PartyEventHandler
	state    *partyEventState
}

func NewPartyEventDispatcher(state *partyEventState) *PartyEventDispatcher {
	return &PartyEventDispatcher{
		handlers: make(map[string]PartyEventHandler),
		state:    state,
	}
}

func (d *PartyEventDispatcher) Register(handler PartyEventHandler) {
	if handler == nil {
		return
	}
	eventType := handler.EventType()
	if eventType == "" {
		return
	}
	d.mu.Lock()
	d.handlers[eventType] = handler
	d.mu.Unlock()
}

func (d *PartyEventDispatcher) Dispatch(raw json.RawMessage, evt PartyEventEnvelope) error {
	if d.state == nil {
		return errors.New("party event state is nil")
	}
	d.mu.RLock()
	handler, ok := d.handlers[evt.EventType]
	d.mu.RUnlock()
	if !ok {
		log.Printf("party consumer: ignored unknown event type=%s party_id=%d revision=%d", evt.EventType, evt.PartyID, evt.Revision)
		return nil
	}
	return handler.Handle(evt, raw)
}

type simplePartyEventHandler struct {
	eventType string
	state     *partyEventState
	fn        func(state *partyEventState, evt PartyEventEnvelope, raw json.RawMessage) error
}

func (h *simplePartyEventHandler) EventType() string {
	return h.eventType
}

func (h *simplePartyEventHandler) Handle(evt PartyEventEnvelope, raw json.RawMessage) error {
	return h.fn(h.state, evt, raw)
}

type partyEventState struct {
	worldID               uint32
	internalClient        internal.InternalClient
	mu                    sync.Mutex
	revisions             map[uint32]uint64
	snapshots             map[uint32]*internal.PartySnapshot
	onSnapshot            func(*internal.PartySnapshot)
	onDisbanded           func([]uint32)
	onMemberJoined        func(*internal.PartySnapshot, uint32)
	onMemberLeft          func(*internal.PartySnapshot, *internal.PartySnapshot, uint32, bool)
	onDisbandedPkt        func(*internal.PartySnapshot, uint32)
	onPartyLeaderChange   func(snapshot *internal.PartySnapshot, newLeaderCharacterID uint32, byDisconnect bool)
	onPartyLogOnOff       func(snapshot *internal.PartySnapshot, characterID uint32)
	onPartySnapshotSilent func(snapshot *internal.PartySnapshot)
	onPartyInvite         func(targetCharacterID uint32, partyID uint32, inviterName string, partySearch bool)
	onPartyDeny           func(targetCharacterID uint32, action uint8, deniedCharacterName string)
}

type PartyEventConsumer struct {
	consumer *RabbitMQConsumer
	state    *partyEventState
}

func (s *partyEventState) apply(evt PartyEventEnvelope) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	current := s.revisions[evt.PartyID]
	if evt.EventType == "log_onoff" {
		if err := s.rehydrateLocked(evt.PartyID); err != nil {
			return err
		}
		if snapshot := s.snapshots[evt.PartyID]; snapshot != nil && s.onSnapshot != nil {
			s.onSnapshot(snapshot)
		}
		log.Printf("party consumer: applied type=log_onoff party_id=%d revision=%d", evt.PartyID, s.revisions[evt.PartyID])
		return nil
	}
	if evt.Revision <= current {
		log.Printf("party consumer: stale/duplicate event ignored type=%s party_id=%d revision=%d current=%d", evt.EventType, evt.PartyID, evt.Revision, current)
		return nil
	}
	if current > 0 && evt.Revision != current+1 {
		log.Printf("party consumer: revision gap detected party_id=%d current=%d incoming=%d", evt.PartyID, current, evt.Revision)
		if err := s.rehydrateLocked(evt.PartyID); err != nil {
			return err
		}
		current = s.revisions[evt.PartyID]
		if evt.Revision <= current {
			return nil
		}
	}
	switch evt.EventType {
	case "disbanded":
		var memberIDs []uint32
		if snapshot, ok := s.snapshots[evt.PartyID]; ok && snapshot != nil {
			for _, m := range snapshot.GetMembers() {
				memberIDs = append(memberIDs, m.GetCharacterId())
			}
		}
		if len(memberIDs) == 0 {
			if err := s.rehydrateLocked(evt.PartyID); err == nil {
				if snapshot, ok := s.snapshots[evt.PartyID]; ok && snapshot != nil {
					for _, m := range snapshot.GetMembers() {
						memberIDs = append(memberIDs, m.GetCharacterId())
					}
				}
			}
		}
		delete(s.revisions, evt.PartyID)
		delete(s.snapshots, evt.PartyID)
		if s.onDisbanded != nil {
			s.onDisbanded(memberIDs)
		}
		log.Printf("party consumer: disbanded applied party_id=%d revision=%d", evt.PartyID, evt.Revision)
	default:
		if err := s.rehydrateLocked(evt.PartyID); err != nil {
			return err
		}
		if snapshot := s.snapshots[evt.PartyID]; snapshot != nil && s.onSnapshot != nil {
			s.onSnapshot(snapshot)
		}
		log.Printf("party consumer: applied type=%s party_id=%d revision=%d", evt.EventType, evt.PartyID, s.revisions[evt.PartyID])
	}
	return nil
}

func (s *partyEventState) rehydrateLocked(partyID uint32) error {
	if s.internalClient == nil {
		return errors.New("internal client unavailable for rehydrate")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	reply, err := s.internalClient.GetParty(ctx, &internal.GetPartyRequest{
		WorldId: s.worldID,
		PartyId: partyID,
	})
	if err != nil {
		return fmt.Errorf("GetParty rehydrate failed: %w", err)
	}
	if !reply.GetFound() || reply.GetParty() == nil {
		delete(s.revisions, partyID)
		delete(s.snapshots, partyID)
		return nil
	}
	s.revisions[partyID] = reply.GetParty().GetRevision()
	s.snapshots[partyID] = reply.GetParty()
	log.Printf("party consumer: rehydrated party_id=%d revision=%d", partyID, s.revisions[partyID])
	return nil
}

func (s *partyEventState) applyEmbeddedPartySnapshot(evt PartyEventEnvelope, snap *internal.PartySnapshot, notifySilent bool) (bool, error) {
	if s == nil || snap == nil {
		return false, nil
	}
	partyID := evt.PartyID
	s.mu.Lock()
	current := s.revisions[partyID]
	if evt.Revision <= current {
		s.mu.Unlock()
		log.Printf("party consumer: stale party_snapshot ignored party_id=%d revision=%d current=%d", partyID, evt.Revision, current)
		return false, nil
	}
	if current > 0 && evt.Revision != current+1 {
		log.Printf("party consumer: party_snapshot revision gap party_id=%d current=%d incoming=%d (applying embedded)", partyID, current, evt.Revision)
	}
	cloned, ok := proto.Clone(snap).(*internal.PartySnapshot)
	if !ok || cloned == nil {
		s.mu.Unlock()
		return false, errors.New("party_snapshot: clone failed")
	}
	s.revisions[partyID] = evt.Revision
	s.snapshots[partyID] = cloned
	s.mu.Unlock()

	if s.onSnapshot != nil {
		s.onSnapshot(cloned)
	}
	if notifySilent && s.onPartySnapshotSilent != nil {
		s.onPartySnapshotSilent(cloned)
	}
	log.Printf("party consumer: applied embedded party snapshot party_id=%d revision=%d silent=%v", partyID, evt.Revision, notifySilent)
	return true, nil
}

func NewPartyEventConsumer(
	cfg config.RabbitMQEndpoint,
	worldID uint32,
	channelID uint32,
	internalClient internal.InternalClient,
	onSnapshot func(*internal.PartySnapshot),
	onDisbanded func([]uint32),
	onMemberJoined func(*internal.PartySnapshot, uint32),
	onMemberLeft func(*internal.PartySnapshot, *internal.PartySnapshot, uint32, bool),
	onDisbandedPkt func(*internal.PartySnapshot, uint32),
	onPartyLeaderChange func(snapshot *internal.PartySnapshot, newLeaderCharacterID uint32, byDisconnect bool),
	onPartyLogOnOff func(snapshot *internal.PartySnapshot, characterID uint32),
	onPartySnapshotSilent func(snapshot *internal.PartySnapshot),
	onPartyInvite func(targetCharacterID uint32, partyID uint32, inviterName string, partySearch bool),
	onPartyDeny func(targetCharacterID uint32, action uint8, deniedCharacterName string),
) *PartyEventConsumer {
	queueName := fmt.Sprintf("fm.game.w%d.c%d.party.events", worldID, channelID)
	consumerTag := fmt.Sprintf("fm-game-w%d-c%d-party", worldID, channelID)
	routeAll := fmt.Sprintf("fm.%d.party.all", worldID)
	routeGame := fmt.Sprintf("fm.%d.party.game.%d", worldID, channelID)
	state := &partyEventState{
		worldID:               worldID,
		internalClient:        internalClient,
		revisions:             make(map[uint32]uint64),
		snapshots:             make(map[uint32]*internal.PartySnapshot),
		onSnapshot:            onSnapshot,
		onDisbanded:           onDisbanded,
		onMemberJoined:        onMemberJoined,
		onMemberLeft:          onMemberLeft,
		onDisbandedPkt:        onDisbandedPkt,
		onPartyLeaderChange:   onPartyLeaderChange,
		onPartyLogOnOff:       onPartyLogOnOff,
		onPartySnapshotSilent: onPartySnapshotSilent,
		onPartyInvite:         onPartyInvite,
		onPartyDeny:           onPartyDeny,
	}
	dispatcher := NewPartyEventDispatcher(state)
	registerDefaultPartyEventHandlers(dispatcher, state)
	handler := func(msg amqp.Delivery) error {
		return handlePartyEventDelivery(msg, worldID, dispatcher)
	}
	return &PartyEventConsumer{
		consumer: NewRabbitMQConsumer(cfg, partyEventsExchange, queueName, consumerTag, handler, routeAll, routeGame),
		state:    state,
	}
}

func handlePartyEventDelivery(msg amqp.Delivery, worldID uint32, dispatcher *PartyEventDispatcher) error {
	if dispatcher == nil {
		return errors.New("party dispatcher is nil")
	}
	var evt PartyEventEnvelope
	if err := json.Unmarshal(msg.Body, &evt); err != nil {
		log.Printf("party consumer: invalid payload: %v", err)
		return err
	}
	if evt.WorldID != 0 && evt.WorldID != worldID {
		return nil
	}
	return dispatcher.Dispatch(msg.Body, evt)
}

func registerDefaultPartyEventHandlers(dispatcher *PartyEventDispatcher, state *partyEventState) {
	registerRevisionHandler(dispatcher, state, "created")
	registerRevisionHandler(dispatcher, state, "member_joined")
	registerRevisionHandler(dispatcher, state, "member_left")
	registerRevisionHandler(dispatcher, state, "leader_changed")
	registerRevisionHandler(dispatcher, state, "log_onoff")
	registerRevisionHandler(dispatcher, state, "disbanded")
	registerPartySnapshotHandler(dispatcher, state)
	registerPartyInviteHandler(dispatcher, state)
	registerPartyDenyHandler(dispatcher, state)
}

func registerPartySnapshotHandler(dispatcher *PartyEventDispatcher, state *partyEventState) {
	dispatcher.Register(&simplePartyEventHandler{
		eventType: "party_snapshot",
		state:     state,
		fn: func(s *partyEventState, evt PartyEventEnvelope, raw json.RawMessage) error {
			var payload struct {
				PartySnapshotPB string `json:"party_snapshot_pb"`
			}
			if err := json.Unmarshal(raw, &payload); err != nil {
				return err
			}
			if payload.PartySnapshotPB == "" {
				return errors.New("party_snapshot: missing party_snapshot_pb")
			}
			wire, err := base64.StdEncoding.DecodeString(payload.PartySnapshotPB)
			if err != nil {
				return fmt.Errorf("party_snapshot_pb: %w", err)
			}
			var snap internal.PartySnapshot
			if err := proto.Unmarshal(wire, &snap); err != nil {
				return fmt.Errorf("party_snapshot unmarshal: %w", err)
			}
			if snap.GetPartyId() != evt.PartyID {
				log.Printf("party_snapshot: party_id mismatch envelope=%d snapshot=%d", evt.PartyID, snap.GetPartyId())
			}
			_, err = s.applyEmbeddedPartySnapshot(evt, &snap, true)
			return err
		},
	})
}

func registerPartyInviteHandler(dispatcher *PartyEventDispatcher, state *partyEventState) {
	if state == nil || state.onPartyInvite == nil {
		return
	}
	dispatcher.Register(&simplePartyEventHandler{
		eventType: "party_invite",
		state:     state,
		fn: func(s *partyEventState, _ PartyEventEnvelope, raw json.RawMessage) error {
			var payload struct {
				TargetCharacterID uint32 `json:"target_character_id"`
				PartyID           uint32 `json:"party_id"`
				InviterName       string `json:"inviter_name"`
				PartySearch       bool   `json:"party_search"`
			}
			if err := json.Unmarshal(raw, &payload); err != nil {
				return err
			}
			if payload.TargetCharacterID == 0 {
				return nil
			}
			s.onPartyInvite(payload.TargetCharacterID, payload.PartyID, payload.InviterName, payload.PartySearch)
			return nil
		},
	})
}

func registerPartyDenyHandler(dispatcher *PartyEventDispatcher, state *partyEventState) {
	if state == nil || state.onPartyDeny == nil {
		return
	}
	dispatcher.Register(&simplePartyEventHandler{
		eventType: "party_invite_denied",
		state:     state,
		fn: func(s *partyEventState, _ PartyEventEnvelope, raw json.RawMessage) error {
			var payload struct {
				InviterCharacterID  uint32 `json:"inviter_character_id"`
				DeniedCharacterName string `json:"denied_character_name"`
				Action              uint8  `json:"action"`
			}
			if err := json.Unmarshal(raw, &payload); err != nil {
				return err
			}
			if payload.InviterCharacterID == 0 {
				return nil
			}
			s.onPartyDeny(payload.InviterCharacterID, payload.Action, payload.DeniedCharacterName)
			return nil
		},
	})
}

func registerRevisionHandler(dispatcher *PartyEventDispatcher, state *partyEventState, eventType string) {
	dispatcher.Register(&simplePartyEventHandler{
		eventType: eventType,
		state:     state,
		fn: func(s *partyEventState, evt PartyEventEnvelope, raw json.RawMessage) error {
			if eventType == "log_onoff" && raw != nil {
				var payload struct {
					PartySnapshotPB string `json:"party_snapshot_pb"`
					CharacterID     uint32 `json:"character_id"`
				}
				if err := json.Unmarshal(raw, &payload); err == nil && payload.PartySnapshotPB != "" {
					wire, err := base64.StdEncoding.DecodeString(payload.PartySnapshotPB)
					if err == nil {
						var snap internal.PartySnapshot
						if err := proto.Unmarshal(wire, &snap); err == nil {
							applied, err := s.applyEmbeddedPartySnapshot(evt, &snap, false)
							if err != nil {
								return err
							}
							if applied && s.onPartyLogOnOff != nil && payload.CharacterID != 0 {
								if snapshot := s.snapshots[evt.PartyID]; snapshot != nil {
									s.onPartyLogOnOff(snapshot, payload.CharacterID)
								}
							}
							return nil
						}
					}
				}
			}
			var prevSnapshot *internal.PartySnapshot
			if eventType == "member_left" || eventType == "disbanded" {
				s.mu.Lock()
				prevSnapshot = s.snapshots[evt.PartyID]
				s.mu.Unlock()
			}
			if err := s.apply(evt); err != nil {
				return err
			}

			if eventType == "member_left" && raw != nil {
				var extra struct {
					CharacterID          uint32 `json:"character_id"`
					ExpelledByCharacter  uint32 `json:"expelled_by_character_id"`
					LeaderChanged        bool   `json:"leader_changed"`
					NewLeaderCharacterID uint32 `json:"new_leader_character_id"`
				}
				if err := json.Unmarshal(raw, &extra); err == nil && extra.CharacterID != 0 {
					snapshot := s.snapshots[evt.PartyID]
					if s.onMemberLeft != nil {
						s.onMemberLeft(prevSnapshot, snapshot, extra.CharacterID, extra.ExpelledByCharacter != 0)
					}
					if s.onPartyLeaderChange != nil && extra.LeaderChanged && extra.NewLeaderCharacterID != 0 && snapshot != nil {
						s.onPartyLeaderChange(snapshot, extra.NewLeaderCharacterID, true)
					}
				}
			}
			if eventType == "member_joined" && s.onMemberJoined != nil && raw != nil {
				var extra struct {
					CharacterID uint32 `json:"character_id"`
				}
				if err := json.Unmarshal(raw, &extra); err == nil && extra.CharacterID != 0 {
					snapshot := s.snapshots[evt.PartyID]
					if snapshot != nil {
						s.onMemberJoined(snapshot, extra.CharacterID)
					}
				}
			}
			if eventType == "disbanded" && s.onDisbandedPkt != nil && raw != nil {
				var extra struct {
					CharacterID uint32 `json:"character_id"`
				}
				if err := json.Unmarshal(raw, &extra); err == nil && extra.CharacterID != 0 && prevSnapshot != nil {
					s.onDisbandedPkt(prevSnapshot, extra.CharacterID)
				}
			}
			if eventType == "leader_changed" && s.onPartyLeaderChange != nil && raw != nil {
				var extra struct {
					NewLeaderCharacterID uint32 `json:"new_leader_character_id"`
				}
				if err := json.Unmarshal(raw, &extra); err == nil && extra.NewLeaderCharacterID != 0 {
					snapshot := s.snapshots[evt.PartyID]
					if snapshot != nil {
						s.onPartyLeaderChange(snapshot, extra.NewLeaderCharacterID, false)
					}
				}
			}
			if eventType == "log_onoff" && s.onPartyLogOnOff != nil && raw != nil {
				var extra struct {
					CharacterID uint32 `json:"character_id"`
				}
				if err := json.Unmarshal(raw, &extra); err == nil && extra.CharacterID != 0 {
					snapshot := s.snapshots[evt.PartyID]
					if snapshot != nil {
						s.onPartyLogOnOff(snapshot, extra.CharacterID)
					}
				}
			}

			return nil
		},
	})
}

func (c *PartyEventConsumer) Start() error {
	if c == nil || c.consumer == nil {
		return nil
	}
	return c.consumer.Start()
}

func (c *PartyEventConsumer) Close() error {
	if c == nil || c.consumer == nil {
		return nil
	}
	return c.consumer.Close()
}

func (c *PartyEventConsumer) QueueName() string {
	if c == nil || c.consumer == nil {
		return ""
	}
	return c.consumer.QueueName()
}

func (c *PartyEventConsumer) ApplyPartySnapshot(snapshot *internal.PartySnapshot) {
	if c == nil || c.state == nil || snapshot == nil {
		return
	}
	partyID := snapshot.GetPartyId()
	state := c.state
	state.mu.Lock()
	defer state.mu.Unlock()
	state.revisions[partyID] = snapshot.GetRevision()
	state.snapshots[partyID] = snapshot
	if state.onSnapshot != nil {
		state.onSnapshot(snapshot)
	}
	log.Printf("party consumer: hydrated from login party_id=%d revision=%d", partyID, snapshot.GetRevision())
}

func (c *PartyEventConsumer) CachedPartySnapshot(partyID uint32) *internal.PartySnapshot {
	if c == nil || c.state == nil {
		return nil
	}
	c.state.mu.Lock()
	s := c.state.snapshots[partyID]
	c.state.mu.Unlock()
	if s == nil {
		return nil
	}
	cloned, ok := proto.Clone(s).(*internal.PartySnapshot)
	if !ok || cloned == nil {
		return nil
	}
	return cloned
}
