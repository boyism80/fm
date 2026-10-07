package server

import (
	"context"
	"log"
	"strings"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/core/async"
	"github.com/boyism80/fm/protocol/dto"
	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
	g_actor "github.com/boyism80/fm/services/game/actor"
	"github.com/boyism80/fm/services/game/constant"
	"github.com/boyism80/fm/services/game/entity"
)

func (s guildSystem) DisbandAsync(ctx actor.Context, ch *entity.Character, result *int) *async.Promise[*internal.DisbandGuildReply] {
	fail := int(constant.GuildDisbandResultFailed)
	if result == nil {
		result = &fail
	}
	if s.gs == nil || s.gs.internalClient == nil || ch == nil {
		*result = fail
		return nil
	}
	*result = fail
	worldID := s.gs.config.WorldId
	charID := ch.GetID()
	promise := async.NewTask(ctx, core.InternalRPCPerStepTimeout)
	return promise.ThenRPC(func(c context.Context) (*internal.DisbandGuildReply, error) {
		return s.gs.internalClient.DisbandGuild(c, &internal.DisbandGuildRequest{
			WorldId:     worldID,
			CharacterId: charID,
		})
	}, func(reply *internal.DisbandGuildReply) error {
		if reply == nil {
			*result = int(constant.GuildDisbandResultFailed)
			return nil
		}
		if reply.GetOk() {
			*result = int(constant.GuildDisbandResultOK)
			return nil
		}
		switch reply.GetErrorCode() {
		case internal.GuildErrorCode_GUILD_ERROR_NOT_IN_GUILD:
			*result = int(constant.GuildDisbandResultNotInGuild)
		case internal.GuildErrorCode_GUILD_ERROR_NOT_AUTHORIZED:
			*result = int(constant.GuildDisbandResultNotMaster)
		default:
			*result = int(constant.GuildDisbandResultFailed)
		}
		return nil
	}).OnError(func(err error) {
		log.Printf("guild disband: character=%d: %v", charID, err)
		*result = int(constant.GuildDisbandResultFailed)
	})
}

func (s guildSystem) IncCapacityAsync(ctx actor.Context, ch *entity.Character, extendedCap bool, result *int) *async.Promise[*internal.IncreaseGuildCapacityReply] {
	fail := int(constant.GuildIncreaseCapacityResultFailed)
	if result == nil {
		result = &fail
	}
	if s.gs == nil || s.gs.internalClient == nil || ch == nil {
		*result = fail
		return nil
	}
	mesoCost := int32(0)
	if !extendedCap {
		mesoCost = constant.GuildCapacityIncreaseMesoCost
		if ch.Inventory.RemoveMeso(mesoCost) == false {
			*result = int(constant.GuildIncreaseCapacityResultInsufficientMeso)
			return nil
		}
	}
	*result = fail
	worldID := s.gs.config.WorldId
	charID := ch.GetID()
	promise := async.NewTask(ctx, core.InternalRPCPerStepTimeout)
	return promise.ThenRPC(func(c context.Context) (*internal.IncreaseGuildCapacityReply, error) {
		return s.gs.internalClient.IncreaseGuildCapacity(c, &internal.IncreaseGuildCapacityRequest{
			WorldId:     worldID,
			CharacterId: charID,
			ExtendedCap: extendedCap,
		})
	}, func(reply *internal.IncreaseGuildCapacityReply) error {
		if reply.GetOk() {
			*result = int(constant.GuildIncreaseCapacityResultOK)
			return nil
		}
		ch.Inventory.AddMeso(mesoCost)
		switch reply.GetErrorCode() {
		case internal.GuildErrorCode_GUILD_ERROR_NOT_IN_GUILD:
			*result = int(constant.GuildIncreaseCapacityResultNotInGuild)
		case internal.GuildErrorCode_GUILD_ERROR_NOT_AUTHORIZED:
			*result = int(constant.GuildIncreaseCapacityResultNotMaster)
		case internal.GuildErrorCode_GUILD_ERROR_CAPACITY_REACHED:
			*result = int(constant.GuildIncreaseCapacityResultCapacityReached)
		case internal.GuildErrorCode_GUILD_ERROR_INSUFFICIENT_GUILD_GP:
			*result = int(constant.GuildIncreaseCapacityResultInsufficientGP)
		default:
			*result = int(constant.GuildIncreaseCapacityResultFailed)
		}
		return nil
	}).OnError(func(err error) {
		log.Printf("guild inc_capacity: character=%d: %v", charID, err)
		if mesoCost > 0 {
			s.gs.WriteOperationLogAsync(ctx, charID, constant.OperationLogGuildIncCapacityUnsettled, mesoCost, err.Error())
		}
		*result = int(constant.GuildIncreaseCapacityResultFailed)
	})
}

func (s guildSystem) GainGPAsync(ctx actor.Context, guildID uint32, amount int32) *async.Promise[*internal.GainGuildGPReply] {
	if s.gs == nil || s.gs.internalClient == nil || guildID == 0 || amount == 0 {
		return nil
	}
	worldID := s.gs.config.WorldId
	promise := async.NewTask(ctx, core.InternalRPCPerStepTimeout)
	return promise.ThenRPC(func(c context.Context) (*internal.GainGuildGPReply, error) {
		return s.gs.internalClient.GainGuildGP(c, &internal.GainGuildGPRequest{
			WorldId: worldID,
			GuildId: guildID,
			Amount:  amount,
		})
	}, func(reply *internal.GainGuildGPReply) error {
		if reply.GetOk() == false {
			log.Printf("guild gain_gp: guild=%d amount=%d: %s", guildID, amount, reply.GetErrorCode())
		}
		return nil
	}).OnError(func(err error) {
		log.Printf("guild gain_gp: guild=%d amount=%d: %v", guildID, amount, err)
	})
}

func (s guildSystem) SendMessageAsync(ctx actor.Context, guildID uint32, messageType constant.ServerMessageType, message string) *async.Promise[*internal.SendGuildMessageReply] {
	if s.gs == nil || s.gs.internalClient == nil || guildID == 0 || message == "" {
		return nil
	}
	worldID := s.gs.config.WorldId
	promise := async.NewTask(ctx, core.InternalRPCPerStepTimeout)
	return promise.ThenRPC(func(c context.Context) (*internal.SendGuildMessageReply, error) {
		return s.gs.internalClient.SendGuildMessage(c, &internal.SendGuildMessageRequest{
			WorldId:     worldID,
			GuildId:     guildID,
			MessageType: uint32(messageType),
			Message:     message,
		})
	}, func(reply *internal.SendGuildMessageReply) error {
		if reply.GetOk() == false {
			log.Printf("guild message: guild=%d rejected", guildID)
		}
		return nil
	}).OnError(func(err error) {
		log.Printf("guild message: guild=%d: %v", guildID, err)
	})
}

func (s guildSystem) ShowRankingAsync(ctx actor.Context, ch *entity.Character, npcID uint32) *async.Promise[*internal.GetGuildRankingReply] {
	if s.gs == nil || s.gs.internalClient == nil {
		return nil
	}
	worldID := s.gs.config.WorldId
	charID := ch.GetID()
	promise := async.NewTask(ctx, core.InternalRPCPerStepTimeout)
	return promise.ThenRPC(func(c context.Context) (*internal.GetGuildRankingReply, error) {
		return s.gs.internalClient.GetGuildRanking(c, &internal.GetGuildRankingRequest{
			WorldId: worldID,
		})
	}, func(reply *internal.GetGuildRankingReply) error {
		entries := make([]dto.GuildRankingEntry, 0, len(reply.GetEntries()))
		for _, e := range reply.GetEntries() {
			logo := e.GetLogo()
			entries = append(entries, dto.GuildRankingEntry{
				Name:        e.GetName(),
				GP:          e.GetGp(),
				Logo:        logo.GetLogo(),
				LogoColor:   logo.GetLogoColor(),
				LogoBG:      logo.GetLogoBg(),
				LogoBGColor: logo.GetLogoBgColor(),
			})
		}
		// The character may have changed maps (and actors) while the RPC was in flight.
		s.gs.EnsureSend(nil, charID, &g_actor.DeliverGuildRanking{
			CharacterID: charID,
			NPCID:       npcID,
			Entries:     entries,
		})
		return nil
	}).OnError(func(err error) {
		log.Printf("guild ranking: character=%d: %v", charID, err)
	})
}

func (s guildSystem) CreateAllianceAsync(ctx actor.Context, ch *entity.Character, allianceName string, result *int) *async.Promise[*internal.CreateAllianceReply] {
	fail := int(constant.AllianceCreateResultFailed)
	if result == nil {
		result = &fail
	}
	if s.gs == nil || s.gs.internalClient == nil || ch == nil {
		*result = fail
		return nil
	}
	partnerID, ok := s.gs.alliance.ValidateCreateAlliance(ch)
	if !ok || partnerID == 0 {
		*result = int(constant.AllianceCreateResultInvalidRequirements)
		return nil
	}
	trimmed := strings.TrimSpace(allianceName)
	if len(trimmed) < 3 || len(trimmed) > 12 {
		*result = int(constant.AllianceCreateResultInvalidName)
		return nil
	}
	if ch.Inventory.RemoveMeso(constant.AllianceCreateMesoCost) == false {
		*result = int(constant.AllianceCreateResultInsufficientMeso)
		return nil
	}
	*result = fail
	worldID := s.gs.config.WorldId
	charID := ch.GetID()
	promise := async.NewTask(ctx, core.InternalRPCPerStepTimeout)
	return promise.ThenRPC(func(c context.Context) (*internal.CreateAllianceReply, error) {
		return s.gs.internalClient.CreateAlliance(c, &internal.CreateAllianceRequest{
			WorldId:            worldID,
			AllianceName:       trimmed,
			LeaderCharacterId:  charID,
			PartnerCharacterId: partnerID,
		})
	}, func(reply *internal.CreateAllianceReply) error {
		if !reply.GetOk() || reply.GetAlliance() == nil {
			ch.Inventory.AddMeso(constant.AllianceCreateMesoCost)
			switch reply.GetErrorCode() {
			case internal.AllianceErrorCode_ALLIANCE_ERROR_ALLIANCE_NAME_INVALID:
				*result = int(constant.AllianceCreateResultInvalidName)
			case internal.AllianceErrorCode_ALLIANCE_ERROR_ALLIANCE_NAME_TAKEN:
				*result = int(constant.AllianceCreateResultNameTaken)
			default:
				*result = int(constant.AllianceCreateResultFailed)
			}
			return nil
		}
		alliancePb := reply.GetAlliance()
		s.gs.alliance.Update(alliancePb)
		s.gs.alliance.BroadcastCreate(alliancePb)
		*result = int(constant.AllianceCreateResultOK)
		return nil
	}).OnError(func(err error) {
		log.Printf("alliance create: character=%d: %v", charID, err)
		s.gs.WriteOperationLogAsync(ctx, charID, constant.OperationLogAllianceCreateUnsettled, constant.AllianceCreateMesoCost, err.Error())
		*result = int(constant.AllianceCreateResultFailed)
	})
}

func (s guildSystem) DisbandAllianceAsync(ctx actor.Context, ch *entity.Character, result *int) *async.Promise[*internal.DisbandAllianceReply] {
	fail := int(constant.AllianceDisbandResultFailed)
	if result == nil {
		result = &fail
	}
	if s.gs == nil || s.gs.internalClient == nil || ch == nil {
		*result = fail
		return nil
	}
	*result = fail
	guildID, ok := ch.GetGuildID()
	if !ok {
		*result = int(constant.AllianceDisbandResultNotInAlliance)
		return nil
	}
	g := s.gs.guild.Get(guildID)
	allianceID, inAlliance := g.GetAllianceID()
	if g == nil || !inAlliance {
		*result = int(constant.AllianceDisbandResultNotInAlliance)
		return nil
	}
	charID := ch.GetID()
	if g.LeaderCharacterID != charID {
		*result = int(constant.AllianceDisbandResultNotGuildMaster)
		return nil
	}
	m := g.FindMember(charID)
	if m == nil {
		*result = int(constant.AllianceDisbandResultNotLeader)
		return nil
	}
	ar, isAllianceLeader := m.GetAllianceRank()
	if !isAllianceLeader || ar != 1 {
		*result = int(constant.AllianceDisbandResultNotLeader)
		return nil
	}
	worldID := s.gs.config.WorldId
	promise := async.NewTask(ctx, core.InternalRPCPerStepTimeout)
	return promise.ThenRPC(func(c context.Context) (*internal.DisbandAllianceReply, error) {
		return s.gs.internalClient.DisbandAlliance(c, &internal.DisbandAllianceRequest{
			WorldId:     worldID,
			CharacterId: charID,
		})
	}, func(reply *internal.DisbandAllianceReply) error {
		if reply == nil {
			*result = int(constant.AllianceDisbandResultFailed)
			return nil
		}
		if !reply.GetOk() {
			switch reply.GetErrorCode() {
			case internal.AllianceErrorCode_ALLIANCE_ERROR_NOT_IN_ALLIANCE,
				internal.AllianceErrorCode_ALLIANCE_ERROR_ALLIANCE_NOT_FOUND:
				*result = int(constant.AllianceDisbandResultNotInAlliance)
			case internal.AllianceErrorCode_ALLIANCE_ERROR_NOT_ALLIANCE_LEADER,
				internal.AllianceErrorCode_ALLIANCE_ERROR_NOT_GUILD_MASTER:
				*result = int(constant.AllianceDisbandResultNotLeader)
			default:
				*result = int(constant.AllianceDisbandResultFailed)
			}
			return nil
		}
		guildIDs := s.gs.alliance.GuildIDs(allianceID)
		if len(guildIDs) == 0 {
			guildIDs = []uint32{guildID}
		}
		s.gs.alliance.DisbandAsync(ctx, allianceID, guildIDs)
		*result = int(constant.AllianceDisbandResultOK)
		return nil
	}).OnError(func(err error) {
		log.Printf("alliance disband: character=%d: %v", charID, err)
		*result = int(constant.AllianceDisbandResultFailed)
	})
}
