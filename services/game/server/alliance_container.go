package server

import (
	"context"
	"log"
	"sync"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/core/async"
	"github.com/boyism80/fm/protocol/dto"
	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
	g_actor "github.com/boyism80/fm/services/game/actor"
	"github.com/boyism80/fm/services/game/entity"
)

type allianceStored struct {
	alliance *entity.Alliance
	info     *dto.AllianceInfo
	guilds   []*dto.GuildInfo
	guildIDs []uint32
}

type AllianceContainer struct {
	gs             *GameServer
	worldID        uint32
	internalClient internal.InternalClient
	mu             sync.Mutex
	revisions      map[uint32]uint64
	alliances      map[uint32]*allianceStored
}

func NewAllianceContainer(gs *GameServer, worldID uint32, ic internal.InternalClient) *AllianceContainer {
	return &AllianceContainer{
		gs:             gs,
		worldID:        worldID,
		internalClient: ic,
		revisions:      make(map[uint32]uint64),
		alliances:      make(map[uint32]*allianceStored),
	}
}

func (ac *AllianceContainer) UpdateAsync(ctx actor.Context, evt AllianceEventEnvelope) *async.Promise {
	p := async.NewPromise(ctx, core.InternalRPCPerStepTimeout)
	allianceID := evt.AllianceID
	p.OnError(func(err error) {
		log.Printf("alliance consumer: apply type=%s alliance_id=%d: %v", evt.EventType, allianceID, err)
	})
	if ac.internalClient == nil {
		p.Then(func() (interface{}, error) {
			return nil, nil
		}, func(interface{}) error {
			return nil
		})
		return p
	}
	async.ThenRPC(p,
		func(c context.Context) (*internal.GetAllianceReply, error) {
			return ac.internalClient.GetAlliance(c, &internal.GetAllianceRequest{
				WorldId:    ac.worldID,
				AllianceId: allianceID,
			})
		},
		func(reply *internal.GetAllianceReply) error {
			if reply == nil || !reply.GetFound() || reply.GetAlliance() == nil {
				ac.Remove(allianceID)
				log.Printf("alliance consumer: removed alliance_id=%d type=%s", allianceID, evt.EventType)
				return nil
			}
			ac.Update(reply.GetAlliance())
			log.Printf("alliance consumer: applied type=%s alliance_id=%d revision=%d", evt.EventType, allianceID, reply.GetAlliance().GetRevision())
			return nil
		},
	)
	return p
}

func (ac *AllianceContainer) Update(alliancePb *internal.Alliance) {
	if alliancePb == nil {
		return
	}
	allianceID := alliancePb.GetAllianceId()
	for _, g := range alliancePb.GetGuilds() {
		if g == nil {
			continue
		}
		ac.gs.guild.Update(g)
	}
	for _, guildID := range alliancePb.GetGuildIds() {
		ent := ac.gs.guild.Get(guildID)
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
		ac.gs.guild.Update(cloned.ToProto())
	}
	info := entity.AllianceInfoFromProto(alliancePb)
	if info == nil {
		return
	}
	guilds := entity.AllianceCreateGuildsFromProto(alliancePb)
	guildIDs := append([]uint32(nil), alliancePb.GetGuildIds()...)
	ent := entity.AllianceFromInfo(ac.gs, info)
	if ent == nil {
		return
	}
	ent.LeaderCharacterID = alliancePb.GetLeaderCharacterId()
	stored := &allianceStored{
		alliance: ent,
		info:     info,
		guilds:   guilds,
		guildIDs: guildIDs,
	}
	ac.mu.Lock()
	ac.revisions[allianceID] = alliancePb.GetRevision()
	ac.alliances[allianceID] = stored
	ac.mu.Unlock()
	log.Printf("alliance: hydrated alliance_id=%d revision=%d", allianceID, alliancePb.GetRevision())
}

func (ac *AllianceContainer) Remove(allianceID uint32) []uint32 {
	ac.mu.Lock()
	defer ac.mu.Unlock()
	entry := ac.alliances[allianceID]
	delete(ac.revisions, allianceID)
	delete(ac.alliances, allianceID)
	if entry == nil {
		return nil
	}
	return append([]uint32(nil), entry.guildIDs...)
}

func (ac *AllianceContainer) Get(allianceID uint32) *entity.Alliance {
	ac.mu.Lock()
	stored := ac.alliances[allianceID]
	ac.mu.Unlock()
	if stored == nil || stored.info == nil {
		return nil
	}
	if stored.alliance == nil {
		return nil
	}
	return stored.alliance.Clone()
}

func (ac *AllianceContainer) ShowData(allianceID uint32) (*dto.AllianceInfo, []*dto.GuildInfo) {
	ac.mu.Lock()
	stored := ac.alliances[allianceID]
	ac.mu.Unlock()
	if stored == nil || stored.info == nil {
		return nil, nil
	}
	var guilds []*dto.GuildInfo
	if len(stored.guilds) > 0 {
		guilds = append([]*dto.GuildInfo(nil), stored.guilds...)
	}
	return stored.info, guilds
}

func (ac *AllianceContainer) GuildIDs(allianceID uint32) []uint32 {
	ac.mu.Lock()
	stored := ac.alliances[allianceID]
	ac.mu.Unlock()
	if stored == nil || len(stored.guildIDs) == 0 {
		return nil
	}
	return append([]uint32(nil), stored.guildIDs...)
}

func (ac *AllianceContainer) NotifyMemberFieldsChanged(ch *entity.Character) {
	if ac.gs == nil || ch == nil {
		return
	}
	ac.gs.guild.NotifyMemberFieldsChanged(ch)
	guildID, inGuild := ch.GetGuildID()
	if !inGuild {
		return
	}
	g := ac.gs.guild.Get(guildID)
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
	ac.deliverMemberFieldsChange(allianceID, guildID, subjectID, level, classID)
}

func (ac *AllianceContainer) deliverMemberFieldsChange(allianceID uint32, subjectGuildID uint32, subjectID uint32, level uint32, classID uint32) {
	if ac.gs == nil || subjectID == 0 {
		return
	}
	guildIDs := ac.GuildIDs(allianceID)
	if len(guildIDs) == 0 {
		guildIDs = []uint32{subjectGuildID}
	}
	for _, gid := range guildIDs {
		g := ac.gs.guild.Get(gid)
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
			if ac.gs.characterRuntime == nil || !ac.gs.characterRuntime.Exists(memberID) {
				continue
			}
			ac.gs.EnsureSend(nil, memberID, &g_actor.DeliverAllianceMemberFieldsChange{
				CharacterID: memberID,
				AllianceID:  allianceID,
				GuildID:     subjectGuildID,
				SubjectID:   subjectID,
				Level:       level,
				ClassID:     classID,
			})
		}
	}
}

func (ac *AllianceContainer) DeliverMemberOnlineChange(allianceID uint32, subjectGuildID uint32, subjectID uint32, online bool) {
	if ac.gs == nil || subjectID == 0 {
		return
	}
	guildIDs := ac.GuildIDs(allianceID)
	if len(guildIDs) == 0 {
		guildIDs = []uint32{subjectGuildID}
	}
	for _, gid := range guildIDs {
		if gid == subjectGuildID {
			continue
		}
		g := ac.gs.guild.Get(gid)
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
			if ac.gs.characterRuntime == nil || !ac.gs.characterRuntime.Exists(memberID) {
				continue
			}
			ac.gs.EnsureSend(nil, memberID, &g_actor.DeliverAllianceMemberOnlineChange{
				CharacterID: memberID,
				AllianceID:  allianceID,
				GuildID:     subjectGuildID,
				SubjectID:   subjectID,
				Online:      online,
			})
		}
	}
}

func (ac *AllianceContainer) ValidateCreateAlliance(ch *entity.Character) (uint32, bool) {
	if ac.gs == nil || ch == nil {
		return 0, false
	}
	leaderID := ch.GetID()
	guildID, inGuild := ch.GetGuildID()
	if !inGuild {
		return 0, false
	}
	leaderGuild := ac.gs.guild.Get(guildID)
	if leaderGuild == nil {
		return 0, false
	}
	if !leaderGuild.IsGuildMaster(leaderID) {
		return 0, false
	}
	if _, inAlliance := leaderGuild.GetAllianceID(); inAlliance {
		return 0, false
	}
	partnerID, ok := ac.partyPartnerCharacterID(ch)
	if !ok || partnerID == 0 {
		return 0, false
	}
	partnerGuildID, partnerInGuild := ac.gs.guild.GuildIDForCharacter(partnerID)
	if !partnerInGuild || partnerGuildID == guildID {
		return 0, false
	}
	partnerGuild := ac.gs.guild.Get(partnerGuildID)
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

func (ac *AllianceContainer) partyPartnerCharacterID(ch *entity.Character) (uint32, bool) {
	if ch == nil || ac.gs == nil {
		return 0, false
	}
	partyID := ch.GetPartyID()
	if partyID == nil {
		return 0, false
	}
	party := ac.gs.party.Get(*partyID)
	if party == nil {
		return 0, false
	}
	members := party.GetMembers()
	if len(members) != 2 {
		return 0, false
	}
	selfID := ch.GetID()
	for _, pm := range members {
		if pm == nil || pm.CharacterID == 0 || pm.CharacterID == selfID {
			continue
		}
		return pm.CharacterID, true
	}
	return 0, false
}
