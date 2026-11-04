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
	Logo     *GuildLogo
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

// One update per character is in flight at a time; changes made meanwhile
// are sent as one update with the latest values when it ends.
func (s *GuildSystem) UpdateMemberAsync(ctx actor.Context, ch *Character) *async.Task {
	p := async.NewTask(ctx, core.InternalRPCPerStepTimeout)
	guildID, inGuild := ch.guild.ID()
	if !inGuild {
		return p
	}

	req := &internal.UpdateGuildMemberRequest{
		WorldId:     s.gw.GetWorldID(),
		GuildId:     guildID,
		CharacterId: ch.GetID(),
		Level:       uint32(ch.GetLevel()),
		ClassId:     uint32(ch.class),
	}
	s.mu.Lock()
	if _, sending := s.memberUpdates[req.CharacterId]; sending {
		s.memberUpdates[req.CharacterId] = req
		s.mu.Unlock()
		return p
	}
	s.memberUpdates[req.CharacterId] = nil
	s.mu.Unlock()
	return s.sendMemberUpdate(p, req)
}

func (s *GuildSystem) sendMemberUpdate(p *async.Task, req *internal.UpdateGuildMemberRequest) *async.Task {
	p.OnError(func(err error) {
		log.Printf("guild update member: character=%d guild=%d: %v", req.CharacterId, req.GuildId, err)
	})
	p.Finally(func() {
		s.mu.Lock()
		next := s.memberUpdates[req.CharacterId]
		if next == nil {
			delete(s.memberUpdates, req.CharacterId)
			s.mu.Unlock()
			return
		}
		s.memberUpdates[req.CharacterId] = nil
		s.mu.Unlock()
		s.sendMemberUpdate(async.NewTask(nil, core.InternalRPCPerStepTimeout), next)
	})
	p.ThenRPC(func(c context.Context) (*internal.UpdateGuildMemberReply, error) {
		return s.gw.Systems().Get[*InternalSystem]().client.UpdateGuildMember(c, req)
	}, nil)
	return p
}

func (s *GuildSystem) StoreGuild(ch *Character, guildPb *internal.Guild) {
	guildID := guildPb.GetGuildId()
	ch.guild.SetID(&guildID)
	s.Update(guildPb)
	ch.listener.OnShowGuildInfo(ch)
	ch.listener.OnBroadcastGuildAppearance(ch)
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
	if err := json.Unmarshal([]byte(request.detail), &detail); err != nil {
		log.Printf("guild settle create: character=%d id=%d: %v", ch.GetID(), request.id, err)
		return
	}
	s.CreateAsync(ctx, ch, request.id, detail.Name).Bind(ch).Do(func(reply *internal.CreateGuildReply) error {
		ch.unsettled.Erase(request.id)
		if !reply.GetOk() {
			ch.inventory.AddMeso(detail.Meso)
			return nil
		}
		if guildPb := reply.GetGuild(); guildPb != nil {
			s.StoreGuild(ch, guildPb)
		}
		return nil
	}).OnError(func(err error) {
		request.inFlight = false

		log.Printf("guild settle create: character=%d id=%d: %v", ch.GetID(), request.id, err)
	})
}

func (s *GuildSystem) ChangeEmblemAsync(ctx actor.Context, ch *Character, requestID uint64, guildID uint32, logo *internal.GuildLogo) *async.Promise[*internal.ChangeGuildEmblemReply] {
	req := &internal.ChangeGuildEmblemRequest{
		WorldId:     s.gw.GetWorldID(),
		CharacterId: ch.GetID(),
		Logo:        logo,
		RequestId:   requestID,
		GuildId:     guildID,
	}
	return async.NewTask(ctx, core.InternalRPCPerStepTimeout).Bind(ch).ThenRPC(func(c context.Context) (*internal.ChangeGuildEmblemReply, error) {
		return s.gw.Systems().Get[*InternalSystem]().client.ChangeGuildEmblem(c, req)
	}, nil)
}

func (s *GuildSystem) RefundEmblem(ch *Character, detail GuildEmblemDetail) {
	if !detail.CashItem {
		ch.inventory.AddMeso(detail.Meso)
		return
	}
	item, err := NewItem(constant.GuildEmblemChangeCashItemID, 1, s.gw)
	if err != nil {
		log.Printf("guild emblem refund: character=%d: %v", ch.GetID(), err)
		return
	}
	if _, err := ch.inventory.AddItem(item, false); err != nil {
		log.Printf("guild emblem refund: character=%d: %v", ch.GetID(), err)
	}
}

func (s *GuildSystem) SettleEmblem(ctx actor.Context, ch *Character, request *UnsettledRequest) {
	var detail GuildEmblemDetail
	if err := json.Unmarshal([]byte(request.detail), &detail); err != nil {
		log.Printf("guild settle emblem: character=%d id=%d: %v", ch.GetID(), request.id, err)
		return
	}
	s.ChangeEmblemAsync(ctx, ch, request.id, uint32(request.target), detail.Logo.BuildRPC()).Bind(ch).Do(func(reply *internal.ChangeGuildEmblemReply) error {
		ch.unsettled.Erase(request.id)
		if !reply.GetOk() {
			s.RefundEmblem(ch, detail)
		}
		return nil
	}).OnError(func(err error) {
		request.inFlight = false

		log.Printf("guild settle emblem: character=%d id=%d: %v", ch.GetID(), request.id, err)
	})
}

func (s *GuildSystem) DisbandAsync(ctx actor.Context, ch *Character) (*async.Promise[int], int) {
	worldID := s.gw.GetWorldID()
	charID := ch.GetID()
	promise := async.NewTask(ctx, core.InternalRPCPerStepTimeout)
	return promise.ThenRPC(func(c context.Context) (*internal.DisbandGuildReply, error) {
		return s.gw.Systems().Get[*InternalSystem]().client.DisbandGuild(c, &internal.DisbandGuildRequest{
			WorldId:     worldID,
			CharacterId: charID,
		})
	}, nil).Then(func(reply *internal.DisbandGuildReply) (int, error) {
		if reply == nil {
			return int(constant.GuildDisbandResultFailed), nil
		}
		if reply.GetOk() {
			return int(constant.GuildDisbandResultOK), nil
		}
		switch reply.GetErrorCode() {
		case internal.GuildErrorCode_GUILD_ERROR_NOT_IN_GUILD:
			return int(constant.GuildDisbandResultNotInGuild), nil
		case internal.GuildErrorCode_GUILD_ERROR_NOT_AUTHORIZED:
			return int(constant.GuildDisbandResultNotMaster), nil
		default:
			return int(constant.GuildDisbandResultFailed), nil
		}
	}).Bind(ch).OnError(func(err error) {
		log.Printf("guild disband: character=%d: %v", charID, err)
	}), int(constant.GuildDisbandResultFailed)
}

func (s *GuildSystem) IncCapacityAsync(ctx actor.Context, ch *Character, extendedCap bool) (*async.Promise[int], int) {
	detail := guildCapacityDetail{Extended: extendedCap}
	if !extendedCap {
		detail.Meso = constant.GuildCapacityIncreaseMesoCost
		if !ch.inventory.RemoveMeso(detail.Meso) {
			return nil, int(constant.GuildIncreaseCapacityResultInsufficientMeso)
		}
	}
	guildID, _ := ch.guild.ID()
	id := ch.unsettled.Add(internal.UnsettledType_UNSETTLED_TYPE_GUILD_CAPACITY, uint64(guildID), detail)
	return s.incCapacityAsync(ctx, ch, id, guildID, extendedCap).Bind(ch).Then(func(reply *internal.IncreaseGuildCapacityReply) (int, error) {
		ch.unsettled.Erase(id)
		if !reply.GetOk() {
			ch.inventory.AddMeso(detail.Meso)
		}

		if reply.GetOk() {
			return int(constant.GuildIncreaseCapacityResultOK), nil
		}
		switch reply.GetErrorCode() {
		case internal.GuildErrorCode_GUILD_ERROR_NOT_IN_GUILD:
			return int(constant.GuildIncreaseCapacityResultNotInGuild), nil
		case internal.GuildErrorCode_GUILD_ERROR_NOT_AUTHORIZED:
			return int(constant.GuildIncreaseCapacityResultNotMaster), nil
		case internal.GuildErrorCode_GUILD_ERROR_CAPACITY_REACHED:
			return int(constant.GuildIncreaseCapacityResultCapacityReached), nil
		case internal.GuildErrorCode_GUILD_ERROR_INSUFFICIENT_GUILD_GP:
			return int(constant.GuildIncreaseCapacityResultInsufficientGP), nil
		default:
			return int(constant.GuildIncreaseCapacityResultFailed), nil
		}
	}).OnError(func(err error) {
		ch.unsettled.Find(id).inFlight = false

		log.Printf("guild inc_capacity: character=%d id=%d: %v", ch.GetID(), id, err)
	}), int(constant.GuildIncreaseCapacityResultFailed)
}

func (s *GuildSystem) incCapacityAsync(ctx actor.Context, ch *Character, requestID uint64, guildID uint32, extendedCap bool) *async.Promise[*internal.IncreaseGuildCapacityReply] {
	req := &internal.IncreaseGuildCapacityRequest{
		WorldId:     s.gw.GetWorldID(),
		CharacterId: ch.GetID(),
		ExtendedCap: extendedCap,
		RequestId:   requestID,
		GuildId:     guildID,
	}
	return async.NewTask(ctx, core.InternalRPCPerStepTimeout).Bind(ch).ThenRPC(func(c context.Context) (*internal.IncreaseGuildCapacityReply, error) {
		return s.gw.Systems().Get[*InternalSystem]().client.IncreaseGuildCapacity(c, req)
	}, nil)
}

func (s *GuildSystem) SettleCapacity(ctx actor.Context, ch *Character, request *UnsettledRequest) {
	var detail guildCapacityDetail
	if err := json.Unmarshal([]byte(request.detail), &detail); err != nil {
		log.Printf("guild settle capacity: character=%d id=%d: %v", ch.GetID(), request.id, err)
		return
	}
	s.incCapacityAsync(ctx, ch, request.id, uint32(request.target), detail.Extended).Bind(ch).Do(func(reply *internal.IncreaseGuildCapacityReply) error {
		ch.unsettled.Erase(request.id)
		if !reply.GetOk() {
			ch.inventory.AddMeso(detail.Meso)
		}
		return nil
	}).OnError(func(err error) {
		request.inFlight = false

		log.Printf("guild settle capacity: character=%d id=%d: %v", ch.GetID(), request.id, err)
	})
}

func (s *GuildSystem) GainGPAsync(ctx actor.Context, ch *Character, amount int32) *async.Promise[*internal.GainGuildGPReply] {
	guildID, inGuild := ch.guild.ID()
	if !inGuild || amount == 0 {
		return nil
	}

	id := ch.unsettled.Add(internal.UnsettledType_UNSETTLED_TYPE_GUILD_GP, uint64(guildID), guildGPDetail{Amount: amount})
	return s.gainGPAsync(ctx, id, guildID, amount).Bind(ch).Do(func(reply *internal.GainGuildGPReply) error {
		ch.unsettled.Erase(id)
		if !reply.GetOk() {
			log.Printf("guild gain_gp: guild=%d amount=%d: %s", guildID, amount, reply.GetErrorCode())
		}
		return nil
	}).OnError(func(err error) {
		ch.unsettled.Find(id).inFlight = false

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
	if err := json.Unmarshal([]byte(request.detail), &detail); err != nil {
		log.Printf("guild settle gp: character=%d id=%d: %v", ch.GetID(), request.id, err)
		return
	}
	s.gainGPAsync(ctx, request.id, uint32(request.target), detail.Amount).Bind(ch).Do(func(reply *internal.GainGuildGPReply) error {
		ch.unsettled.Erase(request.id)
		return nil
	}).OnError(func(err error) {
		request.inFlight = false

		log.Printf("guild settle gp: character=%d id=%d: %v", ch.GetID(), request.id, err)
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
			ch.listener.OnGuildRanking(ch, npcID, entries)
		})
		return nil
	}).OnError(func(err error) {
		log.Printf("guild ranking: character=%d: %v", charID, err)
	})
}

func (s *GuildSystem) CreateAllianceAsync(ctx actor.Context, ch *Character, allianceName string) (*async.Promise[int], int) {
	partnerID, ok := s.gw.Systems().Get[*AllianceSystem]().ValidateCreateAlliance(ch)
	if !ok || partnerID == 0 {
		return nil, int(constant.AllianceCreateResultInvalidRequirements)
	}
	trimmed := strings.TrimSpace(allianceName)
	if !s.gw.GetResources().ValidName(trimmed) {
		return nil, int(constant.AllianceCreateResultInvalidName)
	}
	if !ch.inventory.RemoveMeso(constant.AllianceCreateMesoCost) {
		return nil, int(constant.AllianceCreateResultInsufficientMeso)
	}
	guildID, _ := ch.guild.ID()
	detail := allianceCreateDetail{Name: trimmed, Partner: partnerID, Meso: constant.AllianceCreateMesoCost}
	id := ch.unsettled.Add(internal.UnsettledType_UNSETTLED_TYPE_ALLIANCE_CREATE, uint64(guildID), detail)
	return s.createAllianceAsync(ctx, ch, id, guildID, detail).Bind(ch).Then(func(reply *internal.CreateAllianceReply) (int, error) {
		ch.unsettled.Erase(id)
		if !reply.GetOk() {
			ch.inventory.AddMeso(detail.Meso)
		}

		if !reply.GetOk() {
			switch reply.GetErrorCode() {
			case internal.AllianceErrorCode_ALLIANCE_ERROR_ALLIANCE_NAME_INVALID:
				return int(constant.AllianceCreateResultInvalidName), nil
			case internal.AllianceErrorCode_ALLIANCE_ERROR_ALLIANCE_NAME_TAKEN:
				return int(constant.AllianceCreateResultNameTaken), nil
			default:
				return int(constant.AllianceCreateResultFailed), nil
			}
		}
		alliancePb := reply.GetAlliance()
		if alliancePb == nil {
			return int(constant.AllianceCreateResultFailed), nil
		}
		s.gw.Systems().Get[*AllianceSystem]().Update(alliancePb)
		s.gw.Systems().Get[*AllianceSystem]().BroadcastCreate(alliancePb)
		return int(constant.AllianceCreateResultOK), nil
	}).OnError(func(err error) {
		ch.unsettled.Find(id).inFlight = false

		log.Printf("alliance create: character=%d id=%d: %v", ch.GetID(), id, err)
	}), int(constant.AllianceCreateResultFailed)
}

func (s *GuildSystem) createAllianceAsync(ctx actor.Context, ch *Character, requestID uint64, guildID uint32, detail allianceCreateDetail) *async.Promise[*internal.CreateAllianceReply] {
	req := &internal.CreateAllianceRequest{
		WorldId:            s.gw.GetWorldID(),
		AllianceName:       detail.Name,
		LeaderCharacterId:  ch.GetID(),
		PartnerCharacterId: detail.Partner,
		RequestId:          requestID,
		GuildId:            guildID,
	}
	return async.NewTask(ctx, core.InternalRPCPerStepTimeout).Bind(ch).ThenRPC(func(c context.Context) (*internal.CreateAllianceReply, error) {
		return s.gw.Systems().Get[*InternalSystem]().client.CreateAlliance(c, req)
	}, nil)
}

func (s *GuildSystem) SettleAllianceCreate(ctx actor.Context, ch *Character, request *UnsettledRequest) {
	var detail allianceCreateDetail
	if err := json.Unmarshal([]byte(request.detail), &detail); err != nil {
		log.Printf("alliance settle create: character=%d id=%d: %v", ch.GetID(), request.id, err)
		return
	}
	s.createAllianceAsync(ctx, ch, request.id, uint32(request.target), detail).Bind(ch).Do(func(reply *internal.CreateAllianceReply) error {
		ch.unsettled.Erase(request.id)
		if !reply.GetOk() {
			ch.inventory.AddMeso(detail.Meso)
			return nil
		}
		if alliancePb := reply.GetAlliance(); alliancePb != nil {
			s.gw.Systems().Get[*AllianceSystem]().Update(alliancePb)
			s.gw.Systems().Get[*AllianceSystem]().BroadcastCreate(alliancePb)
		}
		return nil
	}).OnError(func(err error) {
		request.inFlight = false

		log.Printf("alliance settle create: character=%d id=%d: %v", ch.GetID(), request.id, err)
	})
}

func (s *GuildSystem) DisbandAllianceAsync(ctx actor.Context, ch *Character) (*async.Promise[int], int) {
	guildID, ok := ch.guild.ID()
	if !ok {
		return nil, int(constant.AllianceDisbandResultNotInAlliance)
	}
	g := s.Get(guildID)
	_, inAlliance := g.GetAllianceID()
	if g == nil || !inAlliance {
		return nil, int(constant.AllianceDisbandResultNotInAlliance)
	}
	charID := ch.GetID()
	if g.leaderCharacterID != charID {
		return nil, int(constant.AllianceDisbandResultNotGuildMaster)
	}
	if !g.IsAllianceLeader(charID) {
		return nil, int(constant.AllianceDisbandResultNotLeader)
	}
	worldID := s.gw.GetWorldID()
	promise := async.NewTask(ctx, core.InternalRPCPerStepTimeout)
	return promise.ThenRPC(func(c context.Context) (*internal.DisbandAllianceReply, error) {
		return s.gw.Systems().Get[*InternalSystem]().client.DisbandAlliance(c, &internal.DisbandAllianceRequest{
			WorldId:     worldID,
			CharacterId: charID,
		})
	}, nil).Then(func(reply *internal.DisbandAllianceReply) (int, error) {
		if reply == nil {
			return int(constant.AllianceDisbandResultFailed), nil
		}
		if !reply.GetOk() {
			switch reply.GetErrorCode() {
			case internal.AllianceErrorCode_ALLIANCE_ERROR_NOT_IN_ALLIANCE,
				internal.AllianceErrorCode_ALLIANCE_ERROR_ALLIANCE_NOT_FOUND:
				return int(constant.AllianceDisbandResultNotInAlliance), nil
			case internal.AllianceErrorCode_ALLIANCE_ERROR_NOT_ALLIANCE_LEADER,
				internal.AllianceErrorCode_ALLIANCE_ERROR_NOT_GUILD_MASTER:
				return int(constant.AllianceDisbandResultNotLeader), nil
			default:
				return int(constant.AllianceDisbandResultFailed), nil
			}
		}
		return int(constant.AllianceDisbandResultOK), nil
	}).Bind(ch).OnError(func(err error) {
		log.Printf("alliance disband: character=%d: %v", charID, err)
	}), int(constant.AllianceDisbandResultFailed)
}
