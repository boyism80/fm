package server

import (
	"context"
	"encoding/json"
	"log"
	"strings"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/core/async"
	"github.com/boyism80/fm/protocol/dto"
	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
	"github.com/boyism80/fm/services/game/constant"
	"github.com/boyism80/fm/services/game/entity"
)

type guildCreateDetail struct {
	Name string
	Meso int32
}

type guildCapacityDetail struct {
	Extended bool
	Meso     int32
}

type guildEmblemDetail struct {
	Logo     *internal.GuildLogo
	CashItem bool
	Meso     int32
}

type guildGPDetail struct {
	Amount int32
}

type allianceCreateDetail struct {
	Name    string
	Partner uint32
	Meso    int32
}

func (s guildSystem) applyGuildFromProto(ch *entity.Character, guildPb *internal.Guild) {
	guildID := guildPb.GetGuildId()
	ch.Guild.SetID(&guildID)
	s.gs.guild.Update(guildPb)
	ch.Listener.OnShowGuildInfo(ch)
	ch.Listener.OnBroadcastGuildAppearance(ch)
}

func (s guildSystem) createAsync(ctx actor.Context, ch *entity.Character, requestID uint64, name string) *async.Promise[*internal.CreateGuildReply] {
	worldID := s.gs.config.WorldId
	req := &internal.CreateGuildRequest{
		WorldId:   worldID,
		GuildName: name,
		Leader:    ch.ToProtoGuildMember(worldID, int32(s.gs.config.ChannelId), internal.GuildMemberRank_GUILD_MEMBER_RANK_MASTER),
		RequestId: requestID,
	}
	return async.NewTask(ctx, core.InternalRPCPerStepTimeout).ThenRPC(func(c context.Context) (*internal.CreateGuildReply, error) {
		return s.gs.internalClient.CreateGuild(c, req)
	}, nil)
}

func (s guildSystem) SettleCreate(ctx actor.Context, ch *entity.Character, request *entity.UnsettledRequest) {
	var detail guildCreateDetail
	if err := json.Unmarshal([]byte(request.Detail), &detail); err != nil {
		log.Printf("guild settle create: character=%d id=%d: %v", ch.GetID(), request.ID, err)
		return
	}
	s.createAsync(ctx, ch, request.ID, detail.Name).Do(func(reply *internal.CreateGuildReply) error {
		ch.Call(func(actor.Context) {
			ch.Unsettled.Erase(request.ID)
			if !reply.GetOk() {
				ch.Inventory.AddMeso(detail.Meso)
				return
			}
			if guildPb := reply.GetGuild(); guildPb != nil {
				s.applyGuildFromProto(ch, guildPb)
			}
		})
		return nil
	}).OnError(func(err error) {
		ch.Call(func(actor.Context) {
			request.InFlight = false
		})
		log.Printf("guild settle create: character=%d id=%d: %v", ch.GetID(), request.ID, err)
	})
}

func (s guildSystem) changeEmblemAsync(ctx actor.Context, ch *entity.Character, requestID uint64, logo *internal.GuildLogo) *async.Promise[*internal.ChangeGuildEmblemReply] {
	req := &internal.ChangeGuildEmblemRequest{
		WorldId:     s.gs.config.WorldId,
		CharacterId: ch.GetID(),
		Logo:        logo,
		RequestId:   requestID,
	}
	return async.NewTask(ctx, core.InternalRPCPerStepTimeout).ThenRPC(func(c context.Context) (*internal.ChangeGuildEmblemReply, error) {
		return s.gs.internalClient.ChangeGuildEmblem(c, req)
	}, nil)
}

func (s guildSystem) refundEmblem(ch *entity.Character, detail guildEmblemDetail) {
	if !detail.CashItem {
		ch.Inventory.AddMeso(detail.Meso)
		return
	}
	item, err := entity.NewItem(constant.GuildEmblemChangeCashItemID, 1, s.gs)
	if err != nil {
		log.Printf("guild emblem refund: character=%d: %v", ch.GetID(), err)
		return
	}
	if _, err := ch.Inventory.AddItem(item, false); err != nil {
		log.Printf("guild emblem refund: character=%d: %v", ch.GetID(), err)
	}
}

func (s guildSystem) SettleEmblem(ctx actor.Context, ch *entity.Character, request *entity.UnsettledRequest) {
	var detail guildEmblemDetail
	if err := json.Unmarshal([]byte(request.Detail), &detail); err != nil {
		log.Printf("guild settle emblem: character=%d id=%d: %v", ch.GetID(), request.ID, err)
		return
	}
	s.changeEmblemAsync(ctx, ch, request.ID, detail.Logo).Do(func(reply *internal.ChangeGuildEmblemReply) error {
		ch.Call(func(actor.Context) {
			ch.Unsettled.Erase(request.ID)
			if !reply.GetOk() {
				s.refundEmblem(ch, detail)
			}
		})
		return nil
	}).OnError(func(err error) {
		ch.Call(func(actor.Context) {
			request.InFlight = false
		})
		log.Printf("guild settle emblem: character=%d id=%d: %v", ch.GetID(), request.ID, err)
	})
}

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
	detail := guildCapacityDetail{Extended: extendedCap}
	if !extendedCap {
		detail.Meso = constant.GuildCapacityIncreaseMesoCost
		if !ch.Inventory.RemoveMeso(detail.Meso) {
			*result = int(constant.GuildIncreaseCapacityResultInsufficientMeso)
			return nil
		}
	}
	*result = fail
	guildID, _ := ch.Guild.ID()
	id := ch.Unsettled.Add(internal.UnsettledType_UNSETTLED_TYPE_GUILD_CAPACITY, uint64(guildID), detail)
	return s.incCapacityAsync(ctx, ch, id, extendedCap).Do(func(reply *internal.IncreaseGuildCapacityReply) error {
		ch.Call(func(actor.Context) {
			ch.Unsettled.Erase(id)
			if !reply.GetOk() {
				ch.Inventory.AddMeso(detail.Meso)
			}
		})
		if reply.GetOk() {
			*result = int(constant.GuildIncreaseCapacityResultOK)
			return nil
		}
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
		ch.Call(func(actor.Context) {
			ch.Unsettled.Find(id).InFlight = false
		})
		log.Printf("guild inc_capacity: character=%d id=%d: %v", ch.GetID(), id, err)
		*result = int(constant.GuildIncreaseCapacityResultFailed)
	})
}

func (s guildSystem) incCapacityAsync(ctx actor.Context, ch *entity.Character, requestID uint64, extendedCap bool) *async.Promise[*internal.IncreaseGuildCapacityReply] {
	req := &internal.IncreaseGuildCapacityRequest{
		WorldId:     s.gs.config.WorldId,
		CharacterId: ch.GetID(),
		ExtendedCap: extendedCap,
		RequestId:   requestID,
	}
	return async.NewTask(ctx, core.InternalRPCPerStepTimeout).ThenRPC(func(c context.Context) (*internal.IncreaseGuildCapacityReply, error) {
		return s.gs.internalClient.IncreaseGuildCapacity(c, req)
	}, nil)
}

func (s guildSystem) SettleCapacity(ctx actor.Context, ch *entity.Character, request *entity.UnsettledRequest) {
	var detail guildCapacityDetail
	if err := json.Unmarshal([]byte(request.Detail), &detail); err != nil {
		log.Printf("guild settle capacity: character=%d id=%d: %v", ch.GetID(), request.ID, err)
		return
	}
	s.incCapacityAsync(ctx, ch, request.ID, detail.Extended).Do(func(reply *internal.IncreaseGuildCapacityReply) error {
		ch.Call(func(actor.Context) {
			ch.Unsettled.Erase(request.ID)
			if !reply.GetOk() {
				ch.Inventory.AddMeso(detail.Meso)
			}
		})
		return nil
	}).OnError(func(err error) {
		ch.Call(func(actor.Context) {
			request.InFlight = false
		})
		log.Printf("guild settle capacity: character=%d id=%d: %v", ch.GetID(), request.ID, err)
	})
}

func (s guildSystem) GainGPAsync(ctx actor.Context, ch *entity.Character, amount int32) *async.Promise[*internal.GainGuildGPReply] {
	guildID, inGuild := ch.Guild.ID()
	if !inGuild || amount == 0 {
		return nil
	}

	id := ch.Unsettled.Add(internal.UnsettledType_UNSETTLED_TYPE_GUILD_GP, uint64(guildID), guildGPDetail{Amount: amount})
	return s.gainGPAsync(ctx, id, guildID, amount).Do(func(reply *internal.GainGuildGPReply) error {
		ch.Call(func(actor.Context) {
			ch.Unsettled.Erase(id)
			if !reply.GetOk() {
				log.Printf("guild gain_gp: guild=%d amount=%d: %s", guildID, amount, reply.GetErrorCode())
			}
		})
		return nil
	}).OnError(func(err error) {
		ch.Call(func(actor.Context) {
			ch.Unsettled.Find(id).InFlight = false
		})
		log.Printf("guild gain_gp: guild=%d amount=%d id=%d: %v", guildID, amount, id, err)
	})
}

func (s guildSystem) gainGPAsync(ctx actor.Context, requestID uint64, guildID uint32, amount int32) *async.Promise[*internal.GainGuildGPReply] {
	req := &internal.GainGuildGPRequest{
		WorldId:   s.gs.config.WorldId,
		GuildId:   guildID,
		Amount:    amount,
		RequestId: requestID,
	}
	return async.NewTask(ctx, core.InternalRPCPerStepTimeout).ThenRPC(func(c context.Context) (*internal.GainGuildGPReply, error) {
		return s.gs.internalClient.GainGuildGP(c, req)
	}, nil)
}

func (s guildSystem) SettleGP(ctx actor.Context, ch *entity.Character, request *entity.UnsettledRequest) {
	var detail guildGPDetail
	if err := json.Unmarshal([]byte(request.Detail), &detail); err != nil {
		log.Printf("guild settle gp: character=%d id=%d: %v", ch.GetID(), request.ID, err)
		return
	}
	s.gainGPAsync(ctx, request.ID, uint32(request.Target), detail.Amount).Do(func(reply *internal.GainGuildGPReply) error {
		ch.Call(func(actor.Context) {
			ch.Unsettled.Erase(request.ID)
		})
		return nil
	}).OnError(func(err error) {
		ch.Call(func(actor.Context) {
			request.InFlight = false
		})
		log.Printf("guild settle gp: character=%d id=%d: %v", ch.GetID(), request.ID, err)
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
		if !reply.GetOk() {
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
		s.gs.GetDispatchSystem().CallCharacter(charID, func(_ actor.Context, ch *entity.Character) {
			if ch == nil {
				return
			}
			ch.Listener.OnGuildRanking(ch, npcID, entries)
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
	if !ch.Inventory.RemoveMeso(constant.AllianceCreateMesoCost) {
		*result = int(constant.AllianceCreateResultInsufficientMeso)
		return nil
	}
	*result = fail
	guildID, _ := ch.Guild.ID()
	detail := allianceCreateDetail{Name: trimmed, Partner: partnerID, Meso: constant.AllianceCreateMesoCost}
	id := ch.Unsettled.Add(internal.UnsettledType_UNSETTLED_TYPE_ALLIANCE_CREATE, uint64(guildID), detail)
	return s.createAllianceAsync(ctx, ch, id, detail).Do(func(reply *internal.CreateAllianceReply) error {
		ch.Call(func(actor.Context) {
			ch.Unsettled.Erase(id)
			if !reply.GetOk() {
				ch.Inventory.AddMeso(detail.Meso)
			}
		})
		if !reply.GetOk() {
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
		if alliancePb == nil {
			return nil
		}
		s.gs.alliance.Update(alliancePb)
		s.gs.alliance.BroadcastCreate(alliancePb)
		*result = int(constant.AllianceCreateResultOK)
		return nil
	}).OnError(func(err error) {
		ch.Call(func(actor.Context) {
			ch.Unsettled.Find(id).InFlight = false
		})
		log.Printf("alliance create: character=%d id=%d: %v", ch.GetID(), id, err)
		*result = int(constant.AllianceCreateResultFailed)
	})
}

func (s guildSystem) createAllianceAsync(ctx actor.Context, ch *entity.Character, requestID uint64, detail allianceCreateDetail) *async.Promise[*internal.CreateAllianceReply] {
	req := &internal.CreateAllianceRequest{
		WorldId:            s.gs.config.WorldId,
		AllianceName:       detail.Name,
		LeaderCharacterId:  ch.GetID(),
		PartnerCharacterId: detail.Partner,
		RequestId:          requestID,
	}
	return async.NewTask(ctx, core.InternalRPCPerStepTimeout).ThenRPC(func(c context.Context) (*internal.CreateAllianceReply, error) {
		return s.gs.internalClient.CreateAlliance(c, req)
	}, nil)
}

func (s guildSystem) SettleAllianceCreate(ctx actor.Context, ch *entity.Character, request *entity.UnsettledRequest) {
	var detail allianceCreateDetail
	if err := json.Unmarshal([]byte(request.Detail), &detail); err != nil {
		log.Printf("alliance settle create: character=%d id=%d: %v", ch.GetID(), request.ID, err)
		return
	}
	s.createAllianceAsync(ctx, ch, request.ID, detail).Do(func(reply *internal.CreateAllianceReply) error {
		ch.Call(func(actor.Context) {
			ch.Unsettled.Erase(request.ID)
			if !reply.GetOk() {
				ch.Inventory.AddMeso(detail.Meso)
				return
			}
			if alliancePb := reply.GetAlliance(); alliancePb != nil {
				s.gs.alliance.Update(alliancePb)
				s.gs.alliance.BroadcastCreate(alliancePb)
			}
		})
		return nil
	}).OnError(func(err error) {
		ch.Call(func(actor.Context) {
			request.InFlight = false
		})
		log.Printf("alliance settle create: character=%d id=%d: %v", ch.GetID(), request.ID, err)
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
	guildID, ok := ch.Guild.ID()
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
	if !g.IsAllianceLeader(charID) {
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
