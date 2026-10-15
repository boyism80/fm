package game

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/core/async"
	"github.com/boyism80/fm/protocol/dto"
	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
)

type allianceStored struct {
	alliance  *Alliance
	info      *dto.AllianceInfo
	guilds    []*dto.GuildInfo
	guildIDs  []uint32
	updatedAt time.Time
}

type AllianceEventEnvelope struct {
	EventID    string `json:"event_id"`
	EventType  string `json:"event_type"`
	WorldID    uint32 `json:"world_id"`
	AllianceID uint32 `json:"alliance_id"`
	OccurredAt string `json:"occurred_at"`
}

type AllianceSystem struct {
	gw        GameWorld
	worldID   uint32
	mu        sync.Mutex
	alliances map[uint32]*allianceStored
}

func NewAllianceSystem(gw GameWorld) *AllianceSystem {
	return &AllianceSystem{
		gw:        gw,
		worldID:   gw.GetWorldID(),
		alliances: make(map[uint32]*allianceStored),
	}
}

func (s *AllianceSystem) Start() error { return nil }
func (s *AllianceSystem) Stop()        {}

func (s *AllianceSystem) UpdateAsync(ctx actor.Context, evt AllianceEventEnvelope) *async.Task {
	p := async.NewTask(ctx, core.InternalRPCPerStepTimeout)
	allianceID := evt.AllianceID
	p.OnError(func(err error) {
		log.Printf("alliance consumer: apply type=%s alliance_id=%d: %v", evt.EventType, allianceID, err)
	})
	p.ThenRPC(
		func(c context.Context) (*internal.GetAllianceReply, error) {
			return s.gw.Systems().Get[*InternalSystem]().client.GetAlliance(c, &internal.GetAllianceRequest{
				WorldId:    s.worldID,
				AllianceId: allianceID,
			})
		},
		func(reply *internal.GetAllianceReply) error {
			if reply == nil || !reply.GetFound() || reply.GetAlliance() == nil {
				s.Remove(allianceID)
				return nil
			}
			s.Update(reply.GetAlliance())
			return nil
		},
	)
	return p
}

func (s *AllianceSystem) Update(alliancePb *internal.Alliance) {
	if alliancePb == nil {
		return
	}
	allianceID := alliancePb.GetAllianceId()
	for _, g := range alliancePb.GetGuilds() {
		if g == nil {
			continue
		}
		s.gw.Systems().Get[*GuildSystem]().Update(g)
	}
	for _, guildID := range alliancePb.GetGuildIds() {
		ent := s.gw.Systems().Get[*GuildSystem]().Get(guildID)
		if ent == nil {
			continue
		}
		id, ok := ent.GetAllianceID()
		if ok && id == allianceID {
			continue
		}
		cloned := ent.Clone()
		if cloned == nil {
			continue
		}
		aid := allianceID
		cloned.AllianceID = &aid
		s.gw.Systems().Get[*GuildSystem]().Update(cloned.BuildRPC())
	}
	info := AllianceInfoFromProto(alliancePb)
	if info == nil {
		return
	}
	guilds := AllianceCreateGuildsFromProto(alliancePb)
	guildIDs := append([]uint32(nil), alliancePb.GetGuildIds()...)
	ent := AllianceFromInfo(s.gw, info)
	if ent == nil {
		return
	}
	ent.LeaderCharacterID = alliancePb.GetLeaderCharacterId()
	stored := &allianceStored{
		alliance:  ent,
		info:      info,
		guilds:    guilds,
		guildIDs:  guildIDs,
		updatedAt: time.UnixMilli(alliancePb.GetUpdatedAtUnixMs()),
	}
	s.mu.Lock()
	if prev := s.alliances[allianceID]; prev != nil && stored.updatedAt.Before(prev.updatedAt) {
		s.mu.Unlock()
		return
	}
	s.alliances[allianceID] = stored
	s.mu.Unlock()
}

func (s *AllianceSystem) Remove(allianceID uint32) []uint32 {
	s.mu.Lock()
	defer s.mu.Unlock()
	entry := s.alliances[allianceID]
	delete(s.alliances, allianceID)
	if entry == nil {
		return nil
	}
	return append([]uint32(nil), entry.guildIDs...)
}

func (s *AllianceSystem) Get(allianceID uint32) *Alliance {
	s.mu.Lock()
	stored := s.alliances[allianceID]
	s.mu.Unlock()
	if stored == nil || stored.info == nil {
		return nil
	}
	if stored.alliance == nil {
		return nil
	}
	return stored.alliance.Clone()
}

func (s *AllianceSystem) ShowData(allianceID uint32) (*dto.AllianceInfo, []*dto.GuildInfo) {
	s.mu.Lock()
	stored := s.alliances[allianceID]
	s.mu.Unlock()
	if stored == nil || stored.info == nil {
		return nil, nil
	}
	var guilds []*dto.GuildInfo
	if len(stored.guilds) > 0 {
		guilds = append([]*dto.GuildInfo(nil), stored.guilds...)
	}
	return stored.info, guilds
}

func (s *AllianceSystem) GuildIDs(allianceID uint32) []uint32 {
	s.mu.Lock()
	stored := s.alliances[allianceID]
	s.mu.Unlock()
	if stored == nil || len(stored.guildIDs) == 0 {
		return nil
	}
	return append([]uint32(nil), stored.guildIDs...)
}

func (s *AllianceSystem) NotifyMemberFieldsChanged(ch *Character) {
	if s.gw == nil || ch == nil {
		return
	}
	s.gw.Systems().Get[*GuildSystem]().NotifyMemberFieldsChanged(ch)
	guildID, inGuild := ch.Guild.ID()
	if !inGuild {
		return
	}
	g := s.gw.Systems().Get[*GuildSystem]().Get(guildID)
	if g == nil {
		return
	}
	allianceID, inAlliance := g.GetAllianceID()
	if !inAlliance {
		return
	}
	subjectID := ch.GetID()
	level := uint32(ch.GetLevel())
	classID := uint32(ch.Class)
	s.broadcastMemberFieldsChanged(allianceID, guildID, subjectID, level, classID)
}

func (s *AllianceSystem) broadcastMemberFieldsChanged(allianceID uint32, subjectGuildID uint32, subjectID uint32, level uint32, classID uint32) {
	if s.gw == nil || subjectID == 0 {
		return
	}
	guildIDs := s.GuildIDs(allianceID)
	if len(guildIDs) == 0 {
		guildIDs = []uint32{subjectGuildID}
	}
	for _, gid := range guildIDs {
		g := s.gw.Systems().Get[*GuildSystem]().Get(gid)
		if g == nil {
			continue
		}
		for _, m := range g.GetMembers() {
			if m == nil {
				continue
			}
			memberID := m.GetCharacterId()
			if memberID == 0 {
				continue
			}
			if !s.gw.Systems().Get[*DispatchSystem]().Exists(memberID) {
				continue
			}
			s.gw.Systems().Get[*DispatchSystem]().Visit(memberID, func(_ actor.Context, ch *Character) {
				if ch == nil {
					return
				}
				ch.Listener.OnAllianceMemberFieldsChange(ch, allianceID, subjectGuildID, subjectID, level, classID)
			})
		}
	}
}

func (s *AllianceSystem) BroadcastMemberOnlineChanged(allianceID uint32, subjectGuildID uint32, subjectID uint32, online bool) {
	if s.gw == nil || subjectID == 0 {
		return
	}
	guildIDs := s.GuildIDs(allianceID)
	if len(guildIDs) == 0 {
		guildIDs = []uint32{subjectGuildID}
	}
	for _, gid := range guildIDs {
		if gid == subjectGuildID {
			continue
		}
		g := s.gw.Systems().Get[*GuildSystem]().Get(gid)
		if g == nil {
			continue
		}
		for _, m := range g.GetMembers() {
			if m == nil {
				continue
			}
			memberID := m.GetCharacterId()
			if memberID == 0 {
				continue
			}
			if !s.gw.Systems().Get[*DispatchSystem]().Exists(memberID) {
				continue
			}
			s.gw.Systems().Get[*DispatchSystem]().Visit(memberID, func(_ actor.Context, ch *Character) {
				if ch == nil {
					return
				}
				ch.Listener.OnAllianceMemberOnlineChange(ch, allianceID, subjectGuildID, subjectID, online)
			})
		}
	}
}

func (s *AllianceSystem) ValidateCreateAlliance(ch *Character) (uint32, bool) {
	if s.gw == nil || ch == nil {
		return 0, false
	}
	leaderID := ch.GetID()
	guildID, inGuild := ch.Guild.ID()
	if !inGuild {
		return 0, false
	}
	leaderGuild := s.gw.Systems().Get[*GuildSystem]().Get(guildID)
	if leaderGuild == nil {
		return 0, false
	}
	if !leaderGuild.IsGuildMaster(leaderID) {
		return 0, false
	}
	if _, inAlliance := leaderGuild.GetAllianceID(); inAlliance {
		return 0, false
	}

	partyID := ch.Party.ID()
	if partyID == nil {
		return 0, false
	}
	party := s.gw.Systems().Get[*PartySystem]().Get(*partyID)
	if party == nil {
		return 0, false
	}
	members := party.GetMembers()
	if len(members) != 2 {
		return 0, false
	}
	var partnerID uint32
	for _, pm := range members {
		if pm == nil || pm.CharacterID == 0 || pm.CharacterID == leaderID {
			continue
		}
		partnerID = pm.CharacterID
		break
	}
	if partnerID == 0 {
		return 0, false
	}

	partnerGuildID, partnerInGuild := s.gw.Systems().Get[*GuildSystem]().GuildIDForCharacter(partnerID)
	if !partnerInGuild || partnerGuildID == guildID {
		return 0, false
	}
	partnerGuild := s.gw.Systems().Get[*GuildSystem]().Get(partnerGuildID)
	if partnerGuild == nil {
		return 0, false
	}
	if !partnerGuild.IsGuildMaster(partnerID) {
		return 0, false
	}
	if _, inAlliance := partnerGuild.GetAllianceID(); inAlliance {
		return 0, false
	}
	return partnerID, true
}
