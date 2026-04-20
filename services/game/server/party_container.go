package server

import (
	"context"
	"errors"
	"log"
	"sync"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/core/async"
	"github.com/boyism80/fm/protocol/constant"
	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
	"github.com/boyism80/fm/protocol/response"
	g_actor "github.com/boyism80/fm/services/game/actor"
	"github.com/boyism80/fm/services/game/entity"
	"google.golang.org/protobuf/proto"
)

type PartyEventEnvelope struct {
	EventID    string `json:"event_id"`
	EventType  string `json:"event_type"`
	WorldID    uint32 `json:"world_id"`
	PartyID    uint32 `json:"party_id"`
	Revision   uint64 `json:"revision"`
	OccurredAt string `json:"occurred_at"`
}

type PartyContainer struct {
	gs             *GameServer
	worldID        uint32
	internalClient internal.InternalClient
	mu             sync.Mutex
	revisions      map[uint32]uint64
	snapshots      map[uint32]*internal.PartySnapshot
}

func NewPartyContainer(gs *GameServer, worldID uint32, ic internal.InternalClient) *PartyContainer {
	return &PartyContainer{
		gs:             gs,
		worldID:        worldID,
		internalClient: ic,
		revisions:      make(map[uint32]uint64),
		snapshots:      make(map[uint32]*internal.PartySnapshot),
	}
}

// UpdateAsync builds a Promise for party-event reconciliation. Chain .Then(...) for follow-up work on the
// same actor, then .Run(). GetParty runs off the actor mailbox; merge and state updates run in ThenRPC's use on ctx's actor.
func (pc *PartyContainer) UpdateAsync(ctx actor.Context, evt PartyEventEnvelope) *async.Promise {
	p := async.NewPromise(ctx, core.InternalRPCPerStepTimeout)
	if pc == nil {
		return p
	}
	partyID := evt.PartyID
	p.OnError(func(err error) {
		log.Printf("party consumer: apply type=%s party_id=%d: %v", evt.EventType, partyID, err)
	})

	if evt.EventType == "disbanded" {
		pc.mu.Lock()
		var memberIDs []uint32
		if snapshot, ok := pc.snapshots[evt.PartyID]; ok && snapshot != nil {
			for _, m := range snapshot.GetMembers() {
				memberIDs = append(memberIDs, m.GetCharacterId())
			}
		}
		if len(memberIDs) > 0 {
			delete(pc.revisions, evt.PartyID)
			delete(pc.snapshots, evt.PartyID)
			pc.ClearPartyMembers(memberIDs)
			log.Printf("party consumer: disbanded applied party_id=%d revision=%d", evt.PartyID, evt.Revision)
			pc.mu.Unlock()
			return p
		}
		pc.mu.Unlock()
	}

	if pc.internalClient == nil {
		p.Then(func() (interface{}, error) {
			return nil, errors.New("party apply: internal client unavailable")
		}, func(interface{}) error { return nil })
		return p
	}

	async.ThenRPC(p,
		func(c context.Context) (*internal.GetPartyReply, error) {
			return pc.internalClient.GetParty(c, &internal.GetPartyRequest{
				WorldId: pc.worldID,
				PartyId: partyID,
			})
		},
		func(reply *internal.GetPartyReply) error {
			pc.mu.Lock()
			defer pc.mu.Unlock()
			pc.mergeGetPartyReplyLocked(partyID, reply)
			return pc.applyStateAfterRehydrateLocked(evt)
		},
	)
	return p
}

func (pc *PartyContainer) mergeGetPartyReplyLocked(partyID uint32, reply *internal.GetPartyReply) {
	if reply == nil || !reply.GetFound() || reply.GetParty() == nil {
		delete(pc.revisions, partyID)
		delete(pc.snapshots, partyID)
		return
	}
	pc.revisions[partyID] = reply.GetParty().GetRevision()
	pc.snapshots[partyID] = reply.GetParty()
	log.Printf("party consumer: rehydrated party_id=%d revision=%d", partyID, pc.revisions[partyID])
}

func (pc *PartyContainer) applyStateAfterRehydrateLocked(evt PartyEventEnvelope) error {
	partyID := evt.PartyID
	if evt.EventType == "log_onoff" {
		if snapshot := pc.snapshots[partyID]; snapshot != nil {
			pc.SyncPartySnapshot(snapshot)
		}
		log.Printf("party consumer: applied type=log_onoff party_id=%d revision=%d", partyID, pc.revisions[partyID])
		return nil
	}
	if evt.EventType == "disbanded" {
		var memberIDs []uint32
		if snapshot, ok := pc.snapshots[partyID]; ok && snapshot != nil {
			for _, m := range snapshot.GetMembers() {
				memberIDs = append(memberIDs, m.GetCharacterId())
			}
		}
		delete(pc.revisions, partyID)
		delete(pc.snapshots, partyID)
		pc.ClearPartyMembers(memberIDs)
		log.Printf("party consumer: disbanded applied party_id=%d revision=%d", partyID, evt.Revision)
		return nil
	}
	if snapshot := pc.snapshots[partyID]; snapshot != nil {
		pc.SyncPartySnapshot(snapshot)
	}
	log.Printf("party consumer: applied type=%s party_id=%d revision=%d", evt.EventType, partyID, pc.revisions[partyID])
	return nil
}

func (pc *PartyContainer) applyEmbeddedPartySnapshot(evt PartyEventEnvelope, snap *internal.PartySnapshot, notifySilent bool) (bool, error) {
	if pc == nil || snap == nil {
		return false, nil
	}
	partyID := evt.PartyID
	pc.mu.Lock()
	cloned, ok := proto.Clone(snap).(*internal.PartySnapshot)
	if !ok || cloned == nil {
		pc.mu.Unlock()
		return false, errors.New("party_snapshot: clone failed")
	}
	rev := cloned.GetRevision()
	pc.revisions[partyID] = rev
	pc.snapshots[partyID] = cloned
	pc.mu.Unlock()

	pc.SyncPartySnapshot(cloned)
	if notifySilent {
		pc.DeliverPartySilentFromSnapshot(cloned)
	}
	log.Printf("party consumer: applied embedded party snapshot party_id=%d revision=%d silent=%v", partyID, rev, notifySilent)
	return true, nil
}

func (pc *PartyContainer) ApplySnapshotFromLogin(snapshot *internal.PartySnapshot) {
	if pc == nil || snapshot == nil {
		return
	}
	partyID := snapshot.GetPartyId()
	pc.mu.Lock()
	defer pc.mu.Unlock()
	pc.revisions[partyID] = snapshot.GetRevision()
	pc.snapshots[partyID] = snapshot
	pc.SyncPartySnapshot(snapshot)
	log.Printf("party consumer: hydrated from login party_id=%d revision=%d", partyID, snapshot.GetRevision())
}

func (pc *PartyContainer) CachedSnapshot(partyID uint32) *internal.PartySnapshot {
	if pc == nil {
		return nil
	}
	pc.mu.Lock()
	s := pc.snapshots[partyID]
	pc.mu.Unlock()
	if s == nil {
		return nil
	}
	cloned, ok := proto.Clone(s).(*internal.PartySnapshot)
	if !ok || cloned == nil {
		return nil
	}
	return cloned
}

func partySnapshotMemberIDs(members []*internal.PartyMemberSnapshot) []uint32 {
	out := make([]uint32, 0, len(members))
	for _, m := range members {
		if m == nil || m.GetCharacterId() == 0 {
			continue
		}
		out = append(out, m.GetCharacterId())
	}
	return out
}

func partyMembersToResponse(members []*internal.PartyMemberSnapshot) []response.PartyMemberStatus {
	out := make([]response.PartyMemberStatus, 0, len(members))
	for _, m := range members {
		if m == nil {
			continue
		}
		ch := int32(-2)
		if m.ChannelIndex != nil {
			ch = *m.ChannelIndex
		}
		doorTown := uint32(999999999)
		doorTarget := uint32(999999999)
		doorX := int32(0)
		doorY := int32(0)
		if d := m.GetDoor(); d != nil {
			doorTown = d.GetTown()
			doorTarget = d.GetTarget()
			doorX = d.GetX()
			doorY = d.GetY()
		}
		out = append(out, response.PartyMemberStatus{
			CharacterID: m.GetCharacterId(),
			Name:        m.GetCharacterName(),
			Class:       m.GetClassId(),
			Level:       m.GetLevel(),
			Channel:     ch,
			MapID:       m.GetMapId(),
			DoorTown:    doorTown,
			DoorTarget:  doorTarget,
			DoorX:       doorX,
			DoorY:       doorY,
		})
	}
	return out
}

func (pc *PartyContainer) SyncCharacterPartyState(characterID uint32, partyID *uint32) {
	if pc == nil || characterID == 0 || pc.gs == nil {
		return
	}
	pc.gs.EnsureSend(nil, characterID, &g_actor.SyncCharacterPartyState{
		CharacterID: characterID,
		PartyID:     partyID,
	})
}

func (pc *PartyContainer) SyncPartySnapshot(snapshot *internal.PartySnapshot) {
	if pc == nil || snapshot == nil || pc.gs == nil {
		return
	}
	partyID := snapshot.GetPartyId()
	for _, member := range snapshot.GetMembers() {
		cid := member.GetCharacterId()
		pid := partyID
		pc.SyncCharacterPartyState(cid, &pid)
	}
}

func (pc *PartyContainer) ClearCharacterPartyID(characterID uint32) {
	if pc == nil || characterID == 0 || pc.gs == nil {
		return
	}
	pc.SyncCharacterPartyState(characterID, nil)
}

func (pc *PartyContainer) ClearPartyMembers(memberIDs []uint32) {
	if pc == nil || len(memberIDs) == 0 || pc.gs == nil {
		return
	}
	for _, cid := range memberIDs {
		pc.ClearCharacterPartyID(cid)
	}
}

func (pc *PartyContainer) DeliverPartyInviteToCharacter(targetCharacterID uint32, partyID uint32, inviterName string, partySearch bool) {
	if pc == nil || pc.gs == nil {
		return
	}
	pc.gs.EnsureSend(nil, targetCharacterID, &g_actor.DeliverPartyInvite{
		CharacterID: targetCharacterID,
		PartyID:     partyID,
		InviterName: inviterName,
		PartySearch: partySearch,
	})
}

func (pc *PartyContainer) DeliverPartyDenyStatusToCharacter(targetCharacterID uint32, action uint8, deniedCharacterName string) {
	if pc == nil || pc.gs == nil || targetCharacterID == 0 {
		return
	}
	pc.gs.EnsureSend(nil, targetCharacterID, &g_actor.DeliverPartyStatusMessage{
		CharacterID: targetCharacterID,
		Code:        constant.PartyStatusCode(action),
		Name:        deniedCharacterName,
	})
}

func (pc *PartyContainer) deliverPartyMemberLeftToMaps(leaverID uint32, prev *internal.PartySnapshot) {
	if pc == nil || pc.gs == nil || pc.gs.characterRuntime == nil || leaverID == 0 || prev == nil {
		return
	}
	gs := pc.gs
	oldIDs := partySnapshotMemberIDs(prev.GetMembers())
	root := gs.GetRootContext()
	if root == nil {
		return
	}
	seen := make(map[string]struct{})
	for _, cid := range oldIDs {
		mapPID, ok := gs.characterRuntime.GetMapPID(cid)
		if !ok || mapPID == nil {
			continue
		}
		key := mapPID.String()
		if _, dup := seen[key]; dup {
			continue
		}
		seen[key] = struct{}{}
		root.Send(mapPID, &g_actor.PartyMemberLeft{
			LeaverID: leaverID,
		})
	}
}

func (pc *PartyContainer) deliverPartyDisbandToMaps(prev *internal.PartySnapshot) {
	if pc == nil || pc.gs == nil || pc.gs.characterRuntime == nil || prev == nil {
		return
	}
	gs := pc.gs
	formerIDs := partySnapshotMemberIDs(prev.GetMembers())
	if len(formerIDs) == 0 {
		return
	}
	root := gs.GetRootContext()
	if root == nil {
		return
	}
	seen := make(map[string]struct{})
	for _, cid := range formerIDs {
		mapPID, ok := gs.characterRuntime.GetMapPID(cid)
		if !ok || mapPID == nil {
			continue
		}
		key := mapPID.String()
		if _, dup := seen[key]; dup {
			continue
		}
		seen[key] = struct{}{}
		root.Send(mapPID, &g_actor.PartyDisband{
			FormerMemberIDs: formerIDs,
		})
	}
}

func (pc *PartyContainer) DeliverPartyJoinUpdate(snapshot *internal.PartySnapshot, joinedCharacterID uint32) {
	if pc == nil || pc.gs == nil || snapshot == nil || joinedCharacterID == 0 {
		return
	}
	gs := pc.gs
	members := snapshot.GetMembers()
	if len(members) == 0 {
		return
	}
	joinName := ""
	for _, m := range members {
		if m != nil && m.GetCharacterId() == joinedCharacterID {
			joinName = m.GetCharacterName()
			break
		}
	}
	respMembers := partyMembersToResponse(members)
	if joinName == "" {
		return
	}
	for _, m := range members {
		if m == nil || m.GetCharacterId() == 0 {
			continue
		}
		gs.EnsureSend(nil, m.GetCharacterId(), &g_actor.DeliverPartyUpdateJoin{
			CharacterID: m.GetCharacterId(),
			ForChannel:  int32(gs.config.ChannelId),
			PartyID:     snapshot.GetPartyId(),
			JoinName:    joinName,
			LeaderID:    snapshot.GetLeaderCharacterId(),
			Members:     respMembers,
		})
	}
}

func (pc *PartyContainer) DeliverPartyLeaveUpdate(prev, current *internal.PartySnapshot, targetCharacterID uint32, expelled bool) {
	if pc == nil || pc.gs == nil || prev == nil || targetCharacterID == 0 {
		return
	}
	gs := pc.gs
	targetName := ""
	oldMembers := prev.GetMembers()
	for _, m := range oldMembers {
		if m != nil && m.GetCharacterId() == targetCharacterID {
			targetName = m.GetCharacterName()
			break
		}
	}
	if targetName == "" {
		return
	}
	var (
		partyID  = prev.GetPartyId()
		leaderID = prev.GetLeaderCharacterId()
		members  []*internal.PartyMemberSnapshot
	)
	if current != nil {
		partyID = current.GetPartyId()
		leaderID = current.GetLeaderCharacterId()
		members = current.GetMembers()
	}
	respMembers := partyMembersToResponse(members)
	pc.deliverPartyMemberLeftToMaps(targetCharacterID, prev)
	for _, m := range oldMembers {
		if m == nil || m.GetCharacterId() == 0 {
			continue
		}
		gs.EnsureSend(nil, m.GetCharacterId(), &g_actor.DeliverPartyUpdateLeave{
			CharacterID: m.GetCharacterId(),
			ForChannel:  int32(gs.config.ChannelId),
			PartyID:     partyID,
			TargetID:    targetCharacterID,
			TargetName:  targetName,
			LeaderID:    leaderID,
			Members:     respMembers,
			Expelled:    expelled,
		})
	}
}

func (pc *PartyContainer) DeliverPartyDisbandUpdate(prev *internal.PartySnapshot, leaderCharacterID uint32) {
	if pc == nil || pc.gs == nil || prev == nil || leaderCharacterID == 0 {
		return
	}
	gs := pc.gs
	pc.deliverPartyDisbandToMaps(prev)
	for _, m := range prev.GetMembers() {
		if m == nil || m.GetCharacterId() == 0 {
			continue
		}
		gs.EnsureSend(nil, m.GetCharacterId(), &g_actor.DeliverPartyUpdateDisband{
			CharacterID: m.GetCharacterId(),
			PartyID:     prev.GetPartyId(),
			LeaderID:    leaderCharacterID,
		})
	}
}

func (pc *PartyContainer) DeliverPartyLeaderChange(snapshot *internal.PartySnapshot, newLeaderCharacterID uint32, byDisconnect bool) {
	if pc == nil || pc.gs == nil || snapshot == nil || newLeaderCharacterID == 0 {
		return
	}
	gs := pc.gs
	for _, m := range snapshot.GetMembers() {
		if m == nil || m.GetCharacterId() == 0 {
			continue
		}
		gs.EnsureSend(nil, m.GetCharacterId(), &g_actor.DeliverPartyUpdateLeaderChange{
			CharacterID:          m.GetCharacterId(),
			NewLeaderCharacterID: newLeaderCharacterID,
			ByDisconnect:         byDisconnect,
		})
	}
}

func (pc *PartyContainer) DeliverPartyLogOnOff(snapshot *internal.PartySnapshot, _ uint32) {
	if pc == nil || pc.gs == nil || snapshot == nil {
		return
	}
	gs := pc.gs
	respMembers := partyMembersToResponse(snapshot.GetMembers())
	for _, m := range snapshot.GetMembers() {
		if m == nil || m.GetCharacterId() == 0 {
			continue
		}
		gs.EnsureSend(nil, m.GetCharacterId(), &g_actor.DeliverPartyUpdateLogOnOff{
			CharacterID: m.GetCharacterId(),
			ForChannel:  int32(gs.config.ChannelId),
			PartyID:     snapshot.GetPartyId(),
			LeaderID:    snapshot.GetLeaderCharacterId(),
			Members:     respMembers,
		})
	}
}

func (pc *PartyContainer) DeliverPartySilentFromSnapshot(snapshot *internal.PartySnapshot) {
	if pc == nil || pc.gs == nil || snapshot == nil {
		return
	}
	gs := pc.gs
	respMembers := partyMembersToResponse(snapshot.GetMembers())
	for _, m := range snapshot.GetMembers() {
		if m == nil || m.GetCharacterId() == 0 {
			continue
		}
		gs.EnsureSend(nil, m.GetCharacterId(), &g_actor.DeliverPartyUpdateSilent{
			CharacterID: m.GetCharacterId(),
			ForChannel:  int32(gs.config.ChannelId),
			PartyID:     snapshot.GetPartyId(),
			LeaderID:    snapshot.GetLeaderCharacterId(),
			Members:     respMembers,
		})
	}
}

// SendPartySilentAsync builds a Promise that sends party silent UI state to the character. Uses cache
// when warm; otherwise schedules GetParty via ThenRPC. Caller must .Run() (from map actor Receive).
func (pc *PartyContainer) SendPartySilentAsync(ctx actor.Context, ch *entity.Character) *async.Promise {
	p := async.NewPromise(ctx, core.InternalRPCPerStepTimeout)
	if pc == nil || ch == nil || pc.gs == nil || ctx == nil {
		return p
	}
	partyIDPtr := ch.GetPartyID()
	if partyIDPtr == nil {
		return p
	}
	partyID := *partyIDPtr
	p.OnError(func(err error) {
		log.Printf("SendPartySilentAsync: GetParty char %d party %d: %v", ch.GetID(), partyID, err)
	})
	if snap := pc.CachedSnapshot(partyID); snap != nil {
		p.Then(func() (interface{}, error) { return nil, nil }, func(interface{}) error {
			pc.sendPartySilentSnapshotToCharacter(ch, snap)
			return nil
		})
		return p
	}
	if pc.internalClient == nil {
		return p
	}
	async.ThenRPC(p,
		func(c context.Context) (*internal.GetPartyReply, error) {
			return pc.internalClient.GetParty(c, &internal.GetPartyRequest{
				WorldId: pc.worldID,
				PartyId: partyID,
			})
		},
		func(reply *internal.GetPartyReply) error {
			if reply == nil || !reply.GetFound() || reply.GetParty() == nil {
				return nil
			}
			pc.sendPartySilentSnapshotToCharacter(ch, reply.GetParty())
			return nil
		},
	)
	return p
}

func (pc *PartyContainer) sendPartySilentSnapshotToCharacter(ch *entity.Character, snap *internal.PartySnapshot) {
	if pc == nil || ch == nil || snap == nil || pc.gs == nil {
		return
	}
	gs := pc.gs
	members := partyMembersToResponse(snap.GetMembers())
	gs.EnsureSend(nil, ch.GetID(), &g_actor.DeliverPartyUpdateSilent{
		CharacterID: ch.GetID(),
		ForChannel:  int32(gs.config.ChannelId),
		PartyID:     snap.GetPartyId(),
		LeaderID:    snap.GetLeaderCharacterId(),
		Members:     members,
	})
}

func (pc *PartyContainer) PartyMemberIndex(characterID uint32, partyID *uint32) int {
	if partyID == nil || *partyID == 0 || pc == nil {
		return 0
	}
	snap := pc.CachedSnapshot(*partyID)
	if snap == nil {
		return 0
	}
	for i, mem := range snap.GetMembers() {
		if mem.GetCharacterId() == characterID {
			return i
		}
	}
	return 0
}
