package game

import (
	"sync"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/protocol/dto"
	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
)

type allianceStored struct {
	alliance *Alliance
	guildIDs []uint32
	revision uint64
}

type AllianceEventEnvelope struct {
	EventType  string `json:"event_type"`
	WorldID    uint32 `json:"world_id"`
	AllianceID uint32 `json:"alliance_id"`
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

func (s *AllianceSystem) Update(alliancePb *internal.Alliance) bool {
	if alliancePb == nil {
		return false
	}
	allianceID := alliancePb.GetAllianceId()
	if s.Get(allianceID) == nil {
		found := false
		for _, guildID := range alliancePb.GetGuildIds() {
			if s.gw.Systems().Get[*GuildSystem]().Get(guildID) != nil {
				found = true
				break
			}
		}
		for _, g := range alliancePb.GetGuilds() {
			if found {
				break
			}
			for _, m := range g.GetMembers() {
				if s.gw.Systems().Get[*DispatchSystem]().Exists(m.GetCharacterId()) {
					found = true
					break
				}
			}
		}
		if !found {
			return false
		}
	}

	for _, g := range alliancePb.GetGuilds() {
		s.gw.Systems().Get[*GuildSystem]().store(g)
	}

	ent := allianceFromProto(s.gw, alliancePb)
	if ent == nil {
		return false
	}
	stored := &allianceStored{
		alliance: ent,
		guildIDs: append([]uint32(nil), alliancePb.GetGuildIds()...),
		revision: alliancePb.GetRevision(),
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if prev := s.alliances[allianceID]; prev != nil && stored.revision < prev.revision {
		return false
	}
	s.alliances[allianceID] = stored
	return true
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
	if stored == nil || stored.alliance == nil {
		return nil
	}
	return stored.alliance.Clone()
}

func (s *AllianceSystem) ShowData(allianceID uint32) (*dto.AllianceInfo, []*dto.GuildInfo) {
	alliance := s.Get(allianceID)
	if alliance == nil {
		return nil, nil
	}
	var guilds []*dto.GuildInfo
	for _, guildID := range s.GuildIDs(allianceID) {
		guild := s.gw.Systems().Get[*GuildSystem]().Get(guildID)
		if guild == nil {
			continue
		}
		if info := guild.BuildS2C(); info != nil {
			guilds = append(guilds, info)
		}
	}
	return alliance.BuildS2C(), guilds
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
	guildID, inGuild := ch.guild.ID()
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
	classID := uint32(ch.class)
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
				ch.listener.OnAllianceMemberFieldsChange(ch, allianceID, subjectGuildID, subjectID, level, classID)
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
				ch.listener.OnAllianceMemberOnlineChange(ch, allianceID, subjectGuildID, subjectID, online)
			})
		}
	}
}

func (s *AllianceSystem) ValidateCreateAlliance(ch *Character) (uint32, bool) {
	if s.gw == nil || ch == nil {
		return 0, false
	}
	leaderID := ch.GetID()
	guildID, inGuild := ch.guild.ID()
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

	partyID := ch.party.ID()
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
		if pm == nil || pm.characterID == leaderID {
			continue
		}
		partnerID = pm.characterID
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
