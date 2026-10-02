package server

import (
	"context"
	"fmt"
	"log"

	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/core/async"
	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
	"github.com/boyism80/fm/protocol/request"
	g_actor "github.com/boyism80/fm/services/game/actor"
	"github.com/boyism80/fm/services/game/client"
	"github.com/boyism80/fm/services/game/constant"
	"github.com/boyism80/fm/services/game/entity"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type LoginGame struct {
	gs *GameServer
}

func (LoginGame) New(gs *GameServer) *LoginGame {
	return &LoginGame{
		gs: gs,
	}
}

func (h *LoginGame) Handle(ctx *core.ClientContext, req *request.LoginGame) error {
	if h.gs.internalClient == nil {
		return fmt.Errorf("internal client not configured")
	}

	worldId := h.gs.config.WorldId
	ic := h.gs.internalClient

	reqMsg := &internal.EnterGameRequest{
		WorldId:     worldId,
		CharacterId: req.PlayerId,
		ChannelId:   h.gs.config.ChannelId,
	}

	if ctx.ActorContext == nil {
		log.Printf("LoginGame: no actor context, cannot run internal RPC")
		return fmt.Errorf("login game: actor context required for internal RPC")
	}

	var enterReply *internal.EnterGameReply
	var partyReply *internal.GetPartyReply
	var guildReply *internal.GetGuildReply
	var allianceReply *internal.GetAllianceReply

	promise := async.NewPromise(ctx.ActorContext, core.InternalRPCPerStepTimeout)
	promise = async.ThenRPC(promise, func(c context.Context) (*internal.EnterGameReply, error) {
		return ic.EnterGame(c, reqMsg)
	}, func(reply *internal.EnterGameReply) error {
		if !reply.GetFound() || reply.GetCharacter() == nil {
			return fmt.Errorf("character %d not found", req.PlayerId)
		}
		enterReply = reply
		return nil
	})
	promise = async.ThenRPC(promise, func(c context.Context) (*internal.GetPartyReply, error) {
		if enterReply == nil || enterReply.PartyId == nil {
			return &internal.GetPartyReply{Found: false}, nil
		}
		return ic.GetParty(c, &internal.GetPartyRequest{
			WorldId: h.gs.config.WorldId,
			PartyId: *enterReply.PartyId,
		})
	}, func(reply *internal.GetPartyReply) error {
		partyReply = reply
		return nil
	})
	promise = async.ThenRPC(promise, func(c context.Context) (*internal.GetGuildReply, error) {
		if enterReply == nil || enterReply.GuildId == nil {
			return &internal.GetGuildReply{Found: false}, nil
		}
		guildID := *enterReply.GuildId
		return ic.GetGuild(c, &internal.GetGuildRequest{
			WorldId: h.gs.config.WorldId,
			GuildId: guildID,
		})
	}, func(reply *internal.GetGuildReply) error {
		guildReply = reply
		return nil
	})
	promise = async.ThenRPC(promise, func(c context.Context) (*internal.GetAllianceReply, error) {
		if guildReply == nil || !guildReply.GetFound() || guildReply.GetGuild() == nil {
			return &internal.GetAllianceReply{Found: false}, nil
		}
		g := guildReply.GetGuild()
		if g.AllianceId == nil {
			return &internal.GetAllianceReply{Found: false}, nil
		}
		return ic.GetAlliance(c, &internal.GetAllianceRequest{
			WorldId:    h.gs.config.WorldId,
			AllianceId: *g.AllianceId,
		})
	}, func(reply *internal.GetAllianceReply) error {
		allianceReply = reply
		return h.finishLoginGame(ctx, req, enterReply, partyReply, guildReply, allianceReply)
	})
	promise.OnError(func(err error) {
		log.Printf("LoginGame (async): %v", err)

		// Only an EnterGame that internal answered with an error is known not to have attached; logging out then would end someone else's session.
		// A promise future timeout is not a gRPC status, and its RPC may still attach.
		if enterReply == nil {
			s, ok := status.FromError(err)
			if ok && s.Code() != codes.DeadlineExceeded && s.Code() != codes.Unavailable && s.Code() != codes.Canceled {
				return
			}
		}

		characterID := req.PlayerId
		channelID := h.gs.config.ChannelId
		logout := async.NewPromise(ctx.ActorContext, core.InternalRPCPerStepTimeout)
		logout = async.ThenRPC(logout, func(c context.Context) (*internal.LogoutSessionReply, error) {
			return ic.LogoutSession(c, &internal.LogoutSessionRequest{
				WorldId:          worldId,
				AccountId:        enterReply.GetCharacter().GetAccountId(),
				CharacterId:      &characterID,
				ChannelId:        &channelID,
				DisconnectSource: internal.SessionDisconnectSource_SESSION_DISCONNECT_SOURCE_GAME_SERVER,
			})
		}, func(*internal.LogoutSessionReply) error {
			return nil
		})
		logout.OnError(func(err error) {
			log.Printf("LoginGame: logout session of character %d after login failure: %v", characterID, err)
		})
	})
	return nil
}

func (h *LoginGame) finishLoginGame(ctx *core.ClientContext, req *request.LoginGame, reply *internal.EnterGameReply, partyReply *internal.GetPartyReply, guildReply *internal.GetGuildReply, allianceReply *internal.GetAllianceReply) error {
	if !reply.GetFound() {
		return fmt.Errorf("character %d not found", req.PlayerId)
	}

	p := reply.GetCharacter()
	if p == nil {
		return fmt.Errorf("character %d not found: empty character payload", req.PlayerId)
	}
	initData := &entity.CharacterInitData{
		ID:           p.GetCharacterId(),
		AccountID:    p.GetAccountId(),
		Name:         p.GetName(),
		Gender:       uint8(p.GetGender()),
		SkinColor:    uint8(p.GetSkinColor()),
		Face:         p.GetFace(),
		Hair:         p.GetHair(),
		Level:        uint8(p.GetLevel()),
		Class:        uint16(p.GetClassId()),
		Role:         uint8(p.GetRole()),
		Str:          uint16(p.GetStr()),
		Dex:          uint16(p.GetDex()),
		Int:          uint16(p.GetIntStat()),
		Luk:          uint16(p.GetLuk()),
		Hp:           p.GetHp(),
		MaxHp:        p.GetMaxHp(),
		Mp:           p.GetMp(),
		MaxMp:        p.GetMaxMp(),
		AbilityPoint: uint16(p.GetAbilityPoint()),
		SkillPoint:   uint16(p.GetSkillPoint()),
		Exp:          p.GetExp(),
		Meso:         p.GetMeso(),
		Population:   uint16(p.GetPopulation()),
		PositionX:    int16(p.GetPositionX()),
		PositionY:    int16(p.GetPositionY()),
		Stance:       uint8(p.GetStance()),
		Hidden:       p.GetHidden(),
	}
	if reply.PartyId != nil {
		v := *reply.PartyId
		initData.PartyID = &v
	}
	if reply.GuildId != nil {
		v := *reply.GuildId
		initData.GuildID = &v
	}

	character := entity.NewCharacter(ctx.Client, h.gs.characterListener, initData, h.gs)
	character.KeyLayout().LoadKeyLayoutProto(reply.GetKeyLayout())

	character.LoadInventory(reply.GetInventory())
	character.LoadSkills(reply.GetSkills())
	character.LoadBuffs(reply.GetBuffs())
	character.LoadQuests(reply.GetQuests())
	character.LoadSavedLocations(reply.GetSavedLocations())
	character.LoadBuddyList(reply.GetBuddies(), reply.GetBuddyCapacity())

	gameClient, ok := ctx.Client.(*client.GameClient)
	if !ok {
		return fmt.Errorf("client is not a GameClient")
	}

	mapID := p.GetMapId()
	spawnPoint := uint8(p.GetSpawnPoint())
	mapInstance := h.gs.GetMapSystem().Get(mapID)
	if mapInstance == nil {
		log.Printf("saved map %d not found, falling back to default", mapID)
		defaultMapID, ok := h.gs.resources.NameToMap("\xed\x97\xa4\xeb\x84\xa4\xec\x8b\x9c\xec\x8a\xa4")
		if !ok {
			return fmt.Errorf("default map not found")
		}
		mapID = defaultMapID
		spawnPoint = 0
		mapInstance = h.gs.GetMapSystem().Get(mapID)
		if mapInstance == nil {
			return fmt.Errorf("default map %d not found", mapID)
		}
	}
	if mapInstance.Wz != nil && mapInstance.Wz.HasForcedReturn() {
		forcedID := uint32(mapInstance.Wz.ForcedReturn)
		if alt := h.gs.GetMapSystem().Get(forcedID); alt != nil {
			mapID = forcedID
			spawnPoint = 0
			mapInstance = alt
		}
	}

	if mapInstance.FindPortal(spawnPoint) == nil {
		spawnPoint = 0
	}

	character.Stance = constant.StanceDefaultValue

	targetMapPID := mapInstance.LogicActorPID()
	if targetMapPID == nil {
		return fmt.Errorf("LogicActor PID not found for map %d", mapID)
	}

	rootContext := h.gs.GetRootContext()
	if rootContext == nil {
		return fmt.Errorf("rootContext not set")
	}

	// Fail before SetCharacter: the login error path ends the session only while no character is set.
	if err := h.gs.characterRuntime.RegisterCharacter(character.GetID(), character.GetName()); err != nil {
		return fmt.Errorf("runtime register: %w", err)
	}
	if gameClient.SetCharacter(character) == false {
		h.gs.characterRuntime.UnregisterCharacter(character.GetID())
		return fmt.Errorf("character %d: client disconnected during login", character.GetID())
	}

	if reply.PartyId != nil && partyReply != nil && partyReply.GetFound() && partyReply.GetParty() != nil {
		h.gs.party.Update(partyReply.GetParty())
	}

	if guildReply != nil && guildReply.GetFound() && guildReply.GetGuild() != nil {
		h.gs.guild.Update(guildReply.GetGuild())
	}

	if allianceReply != nil && allianceReply.GetFound() && allianceReply.GetAlliance() != nil {
		h.gs.alliance.Update(allianceReply.GetAlliance())
	}

	rootContext.Send(targetMapPID, &g_actor.AddCharacter{
		Character:  character,
		TargetMap:  mapInstance,
		SpawnPoint: spawnPoint,
		Init:       true,
	})

	return nil
}
