package timers

import (
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/services/game/constant"
	"github.com/boyism80/fm/services/game/entity"
)

type PartySearchTimer struct{}

func (t *PartySearchTimer) GetName() string {
	return "PartySearch"
}

func (t *PartySearchTimer) GetInterval() time.Duration {
	return 5 * time.Second
}

func (t *PartySearchTimer) GetInitialDelay() time.Duration {
	return 5 * time.Second
}

func (t *PartySearchTimer) Handle(ctx actor.Context, mapData *entity.Map) error {
	_ = ctx
	if mapData == nil || mapData.GetPlayerCount() == 0 {
		return nil
	}

	players := mapData.GetAllPlayers()
	for _, src := range players {
		ch, ok := src.(*entity.Character)
		if !ok || ch == nil {
			continue
		}
		cfg := ch.Party.Search
		if cfg == nil {
			continue
		}
		pid := ch.Party.ID()
		if pid == nil {
			ch.Party.Search = nil
			continue
		}
		if reachedPartySearchTarget(mapData, *pid, cfg) {
			ch.Party.Search = nil
			continue
		}
		targetIDs := make([]uint32, 0)
		for _, dst := range players {
			target, ok := dst.(*entity.Character)
			if !ok || target == nil || target.GetID() == ch.GetID() {
				continue
			}
			if target.Party.ID() != nil {
				continue
			}
			lvl := int32(target.GetLevel())
			if lvl < cfg.MinLevel || lvl > cfg.MaxLevel {
				continue
			}
			if !ch.HasRoleAtLeast(constant.RoleAdmin) && target.HasRoleAtLeast(constant.RoleAdmin) {
				continue
			}
			if !entity.MatchesPartySearchClassMask(target, cfg.ClassMask) {
				continue
			}
			targetIDs = append(targetIDs, target.GetID())
		}
		if len(targetIDs) == 0 {
			continue
		}
		inviterName := ch.GetName()
		for _, targetID := range targetIDs {
			target := mapData.GetPlayer(targetID)
			if target == nil {
				continue
			}
			target.Listener.OnPartyInvite(target, *pid, inviterName, true)
		}
		if reachedPartySearchTarget(mapData, *pid, cfg) {
			ch.Party.Search = nil
		}
	}
	return nil
}

func reachedPartySearchTarget(mapData *entity.Map, partyID uint32, cfg *entity.PartySearchConfig) bool {
	if mapData == nil || cfg == nil {
		return true
	}
	members := len(mapData.GetPartyMembers(partyID))
	return members >= int(cfg.MembersNeeded) || members >= 6
}
