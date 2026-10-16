package game

import (
	"context"
	"encoding/json"
	"log"
	"strings"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/core/async"
	"github.com/boyism80/fm/maple/constant"
	"github.com/boyism80/fm/protocol/dto"
	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
)

type GuildCreateDetail struct {
	Name string
	Meso int32
}

type guildCapacityDetail struct {
	Extended bool
	Meso     int32
}

type GuildEmblemDetail struct {
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

func (s *GuildSystem) ApplyGuildFromProto(ch *Character, guildPb *internal.Guild) {
	guildID := guildPb.GetGuildId()
	ch.Guild.SetID(&guildID)
	s.Update(guildPb)
	ch.Listener.OnShowGuildInfo(ch)
	ch.Listener.OnBroadcastGuildAppearance(ch)
}

func (s *GuildSystem) CreateAsync(ctx actor.Context, ch *Character, requestID uint64, name string) *async.Promise[*internal.CreateGuildReply] {
	worldID := s.gw.GetWorldID()
	req := &internal.CreateGuildRequest{
		WorldId:   worldID,
		GuildName: name,
		Leader:    GuildMemberFromCharacter(ch, worldID, s.gw.GetChannelID(), internal.GuildMemberRank_GUILD_MEMBER_RANK_MASTER).BuildRPC(),
		RequestId: requestID,
	}
	return async.NewTask(ctx, core.InternalRPCPerStepTimeout).Bind(ch).ThenRPC(func(c context.Context) (*internal.CreateGuildReply, error) {
		return s.gw.Systems().Get[*InternalSystem]().client.CreateGuild(c, req)
	}, nil)
}

func (s *GuildSystem) SettleCreate(ctx actor.Context, ch *Character, request *UnsettledRequest) {
	var detail GuildCreateDetail
	if err := json.Unmarshal([]byte(request.Detail), &detail); err != nil {
		log.Printf("guild settle create: character=%d id=%d: %v", ch.GetID(), request.ID, err)
		return
	}
	s.CreateAsync(ctx, ch, request.ID, detail.Name).Bind(ch).Do(func(reply *internal.CreateGuildReply) error {
		ch.Unsettled.Erase(request.ID)
		if !reply.GetOk() {
			ch.Inventory.AddMeso(detail.Meso)
			return nil
		}
		if guildPb := reply.GetGuild(); guildPb != nil {
			s.ApplyGuildFromProto(ch, guildPb)
		}
		return nil
	}).OnError(func(err error) {
		request.InFlight = false

		log.Printf("guild settle create: character=%d id=%d: %v", ch.GetID(), request.ID, err)
	})
}

func (s *GuildSystem) ChangeEmblemAsync(ctx actor.Context, ch *Character, requestID uint64, logo *internal.GuildLogo) *async.Promise[*internal.ChangeGuildEmblemReply] {
	req := &internal.ChangeGuildEmblemRequest{
		WorldId:     s.gw.GetWorldID(),
		CharacterId: ch.GetID(),
		Logo:        logo,
		RequestId:   requestID,
	}
	return async.NewTask(ctx, core.InternalRPCPerStepTimeout).Bind(ch).ThenRPC(func(c context.Context) (*internal.ChangeGuildEmblemReply, error) {
		return s.gw.Systems().Get[*InternalSystem]().client.ChangeGuildEmblem(c, req)
	}, nil)
}

func (s *GuildSystem) RefundEmblem(ch *Character, detail GuildEmblemDetail) {
	if !detail.CashItem {
		ch.Inventory.AddMeso(detail.Meso)
		return
	}
	item, err := NewItem(constant.GuildEmblemChangeCashItemID, 1, s.gw)
	if err != nil {
		log.Printf("guild emblem refund: character=%d: %v", ch.GetID(), err)
		return
	}
	if _, err := ch.Inventory.AddItem(item, false); err != nil {
		log.Printf("guild emblem refund: character=%d: %v", ch.GetID(), err)
	}
}

func (s *GuildSystem) SettleEmblem(ctx actor.Context, ch *Character, request *UnsettledRequest) {
	var detail GuildEmblemDetail
	if err := json.Unmarshal([]byte(request.Detail), &detail); err != nil {
		log.Printf("guild settle emblem: character=%d id=%d: %v", ch.GetID(), request.ID, err)
		return
	}
	s.ChangeEmblemAsync(ctx, ch, request.ID, detail.Logo).Bind(ch).Do(func(reply *internal.ChangeGuildEmblemReply) error {
		ch.Unsettled.Erase(request.ID)
		if !reply.GetOk() {
			s.RefundEmblem(ch, detail)
		}
		return nil
	}).OnError(func(err error) {
		request.InFlight = false

		log.Printf("guild settle emblem: character=%d id=%d: %v", ch.GetID(), request.ID, err)
	})
}

func (s *GuildSystem) DisbandAsync(ctx actor.Context, ch *Character, result *int) *async.Promise[*internal.DisbandGuildReply] {
	fail := int(constant.GuildDisbandResultFailed)
	if result == nil {
		result = &fail
	}
	if ch == nil {
		*result = fail
		return nil
	}
	*result = fail
	worldID := s.gw.GetWorldID()
	charID := ch.GetID()
	promise := async.NewTask(ctx, core.InternalRPCPerStepTimeout)
	return promise.ThenRPC(func(c context.Context) (*internal.DisbandGuildReply, error) {
		return s.gw.Systems().Get[*InternalSystem]().client.DisbandGuild(c, &internal.DisbandGuildRequest{
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
	}).Bind(ch).OnError(func(err error) {
		log.Printf("guild disband: character=%d: %v", charID, err)
		*result = int(constant.GuildDisbandResultFailed)
	})
}

func (s *GuildSystem) IncCapacityAsync(ctx actor.Context, ch *Character, extendedCap bool, result *int) *async.Promise[*internal.IncreaseGuildCapacityReply] {
	fail := int(constant.GuildIncreaseCapacityResultFailed)
	if result == nil {
		result = &fail
	}
	if ch == nil {
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
	return s.incCapacityAsync(ctx, ch, id, extendedCap).Bind(ch).Do(func(reply *internal.IncreaseGuildCapacityReply) error {
		ch.Unsettled.Erase(id)
		if !reply.GetOk() {
			ch.Inventory.AddMeso(detail.Meso)
		}

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
		ch.Unsettled.Find(id).InFlight = false

		log.Printf("guild inc_capacity: character=%d id=%d: %v", ch.GetID(), id, err)
		*result = int(constant.GuildIncreaseCapacityResultFailed)
	})
}

func (s *GuildSystem) incCapacityAsync(ctx actor.Context, ch *Character, requestID uint64, extendedCap bool) *async.Promise[*internal.IncreaseGuildCapacityReply] {
	req := &internal.IncreaseGuildCapacityRequest{
		WorldId:     s.gw.GetWorldID(),
		CharacterId: ch.GetID(),
		ExtendedCap: extendedCap,
		RequestId:   requestID,
	}
	return async.NewTask(ctx, core.InternalRPCPerStepTimeout).Bind(ch).ThenRPC(func(c context.Context) (*internal.IncreaseGuildCapacityReply, error) {
		return s.gw.Systems().Get[*InternalSystem]().client.IncreaseGuildCapacity(c, req)
	}, nil)
}

func (s *GuildSystem) SettleCapacity(ctx actor.Context, ch *Character, request *UnsettledRequest) {
	var detail guildCapacityDetail
	if err := json.Unmarshal([]byte(request.Detail), &detail); err != nil {
		log.Printf("guild settle capacity: character=%d id=%d: %v", ch.GetID(), request.ID, err)
		return
	}
	s.incCapacityAsync(ctx, ch, request.ID, detail.Extended).Bind(ch).Do(func(reply *internal.IncreaseGuildCapacityReply) error {
		ch.Unsettled.Erase(request.ID)
		if !reply.GetOk() {
			ch.Inventory.AddMeso(detail.Meso)
		}
		return nil
	}).OnError(func(err error) {
		request.InFlight = false

		log.Printf("guild settle capacity: character=%d id=%d: %v", ch.GetID(), request.ID, err)
	})
}

func (s *GuildSystem) GainGPAsync(ctx actor.Context, ch *Character, amount int32) *async.Promise[*internal.GainGuildGPReply] {
	guildID, inGuild := ch.Guild.ID()
	if !inGuild || amount == 0 {
		return nil
	}

	id := ch.Unsettled.Add(internal.UnsettledType_UNSETTLED_TYPE_GUILD_GP, uint64(guildID), guildGPDetail{Amount: amount})
	return s.gainGPAsync(ctx, id, guildID, amount).Bind(ch).Do(func(reply *internal.GainGuildGPReply) error {
		ch.Unsettled.Erase(id)
		if !reply.GetOk() {
			log.Printf("guild gain_gp: guild=%d amount=%d: %s", guildID, amount, reply.GetErrorCode())
		}
		return nil
	}).OnError(func(err error) {
		ch.Unsettled.Find(id).InFlight = false

		log.Printf("guild gain_gp: guild=%d amount=%d id=%d: %v", guildID, amount, id, err)
	})
}

func (s *GuildSystem) gainGPAsync(ctx actor.Context, requestID uint64, guildID uint32, amount int32) *async.Promise[*internal.GainGuildGPReply] {
	req := &internal.GainGuildGPRequest{
		WorldId:   s.gw.GetWorldID(),
		GuildId:   guildID,
		Amount:    amount,
		RequestId: requestID,
	}
	return async.NewTask(ctx, core.InternalRPCPerStepTimeout).ThenRPC(func(c context.Context) (*internal.GainGuildGPReply, error) {
		return s.gw.Systems().Get[*InternalSystem]().client.GainGuildGP(c, req)
	}, nil)
}

func (s *GuildSystem) SettleGP(ctx actor.Context, ch *Character, request *UnsettledRequest) {
	var detail guildGPDetail
	if err := json.Unmarshal([]byte(request.Detail), &detail); err != nil {
		log.Printf("guild settle gp: character=%d id=%d: %v", ch.GetID(), request.ID, err)
		return
	}
	s.gainGPAsync(ctx, request.ID, uint32(request.Target), detail.Amount).Bind(ch).Do(func(reply *internal.GainGuildGPReply) error {
		ch.Unsettled.Erase(request.ID)
		return nil
	}).OnError(func(err error) {
		request.InFlight = false

		log.Printf("guild settle gp: character=%d id=%d: %v", ch.GetID(), request.ID, err)
	})
}

func (s *GuildSystem) SendMessageAsync(ctx actor.Context, guildID uint32, messageType constant.ServerMessageType, message string) *async.Promise[*internal.SendGuildMessageReply] {
	if guildID == 0 || message == "" {
		return nil
	}
	worldID := s.gw.GetWorldID()
	promise := async.NewTask(ctx, core.InternalRPCPerStepTimeout)
	return promise.ThenRPC(func(c context.Context) (*internal.SendGuildMessageReply, error) {
		return s.gw.Systems().Get[*InternalSystem]().client.SendGuildMessage(c, &internal.SendGuildMessageRequest{
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

func (s *GuildSystem) ShowRankingAsync(ctx actor.Context, ch *Character, npcID uint32) *async.Promise[*internal.GetGuildRankingReply] {
	worldID := s.gw.GetWorldID()
	charID := ch.GetID()
	promise := async.NewTask(ctx, core.InternalRPCPerStepTimeout)
	return promise.ThenRPC(func(c context.Context) (*internal.GetGuildRankingReply, error) {
		return s.gw.Systems().Get[*InternalSystem]().client.GetGuildRanking(c, &internal.GetGuildRankingRequest{
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
		s.gw.Systems().Get[*DispatchSystem]().Visit(charID, func(_ actor.Context, ch *Character) {
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

func (s *GuildSystem) CreateAllianceAsync(ctx actor.Context, ch *Character, allianceName string, result *int) *async.Promise[*internal.CreateAllianceReply] {
	fail := int(constant.AllianceCreateResultFailed)
	if result == nil {
		result = &fail
	}
	if ch == nil {
		*result = fail
		return nil
	}
	partnerID, ok := s.gw.Systems().Get[*AllianceSystem]().ValidateCreateAlliance(ch)
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
	return s.createAllianceAsync(ctx, ch, id, detail).Bind(ch).Do(func(reply *internal.CreateAllianceReply) error {
		ch.Unsettled.Erase(id)
		if !reply.GetOk() {
			ch.Inventory.AddMeso(detail.Meso)
		}

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
		s.gw.Systems().Get[*AllianceSystem]().Update(alliancePb)
		s.gw.Systems().Get[*AllianceSystem]().BroadcastCreate(alliancePb)
		*result = int(constant.AllianceCreateResultOK)
		return nil
	}).OnError(func(err error) {
		ch.Unsettled.Find(id).InFlight = false

		log.Printf("alliance create: character=%d id=%d: %v", ch.GetID(), id, err)
		*result = int(constant.AllianceCreateResultFailed)
	})
}

func (s *GuildSystem) createAllianceAsync(ctx actor.Context, ch *Character, requestID uint64, detail allianceCreateDetail) *async.Promise[*internal.CreateAllianceReply] {
	req := &internal.CreateAllianceRequest{
		WorldId:            s.gw.GetWorldID(),
		AllianceName:       detail.Name,
		LeaderCharacterId:  ch.GetID(),
		PartnerCharacterId: detail.Partner,
		RequestId:          requestID,
	}
	return async.NewTask(ctx, core.InternalRPCPerStepTimeout).Bind(ch).ThenRPC(func(c context.Context) (*internal.CreateAllianceReply, error) {
		return s.gw.Systems().Get[*InternalSystem]().client.CreateAlliance(c, req)
	}, nil)
}

func (s *GuildSystem) SettleAllianceCreate(ctx actor.Context, ch *Character, request *UnsettledRequest) {
	var detail allianceCreateDetail
	if err := json.Unmarshal([]byte(request.Detail), &detail); err != nil {
		log.Printf("alliance settle create: character=%d id=%d: %v", ch.GetID(), request.ID, err)
		return
	}
	s.createAllianceAsync(ctx, ch, request.ID, detail).Bind(ch).Do(func(reply *internal.CreateAllianceReply) error {
		ch.Unsettled.Erase(request.ID)
		if !reply.GetOk() {
			ch.Inventory.AddMeso(detail.Meso)
			return nil
		}
		if alliancePb := reply.GetAlliance(); alliancePb != nil {
			s.gw.Systems().Get[*AllianceSystem]().Update(alliancePb)
			s.gw.Systems().Get[*AllianceSystem]().BroadcastCreate(alliancePb)
		}
		return nil
	}).OnError(func(err error) {
		request.InFlight = false

		log.Printf("alliance settle create: character=%d id=%d: %v", ch.GetID(), request.ID, err)
	})
}

func (s *GuildSystem) DisbandAllianceAsync(ctx actor.Context, ch *Character, result *int) *async.Promise[*internal.DisbandAllianceReply] {
	fail := int(constant.AllianceDisbandResultFailed)
	if result == nil {
		result = &fail
	}
	if ch == nil {
		*result = fail
		return nil
	}
	*result = fail
	guildID, ok := ch.Guild.ID()
	if !ok {
		*result = int(constant.AllianceDisbandResultNotInAlliance)
		return nil
	}
	g := s.Get(guildID)
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
	worldID := s.gw.GetWorldID()
	promise := async.NewTask(ctx, core.InternalRPCPerStepTimeout)
	return promise.ThenRPC(func(c context.Context) (*internal.DisbandAllianceReply, error) {
		return s.gw.Systems().Get[*InternalSystem]().client.DisbandAlliance(c, &internal.DisbandAllianceRequest{
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
		guildIDs := s.gw.Systems().Get[*AllianceSystem]().GuildIDs(allianceID)
		if len(guildIDs) == 0 {
			guildIDs = []uint32{guildID}
		}
		s.gw.Systems().Get[*AllianceSystem]().DisbandAsync(ctx, allianceID, guildIDs)
		*result = int(constant.AllianceDisbandResultOK)
		return nil
	}).Bind(ch).OnError(func(err error) {
		log.Printf("alliance disband: character=%d: %v", charID, err)
		*result = int(constant.AllianceDisbandResultFailed)
	})
}
