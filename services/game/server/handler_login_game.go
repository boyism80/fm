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
)

func optionalUint32(value uint32) *uint32 {
	if value == 0 {
		return nil
	}
	v := value
	return &v
}

type LoginGame struct {
	gs     *GameServer
	opcode byte
}

func (LoginGame) New(gs *GameServer) *LoginGame {
	return &LoginGame{
		gs:     gs,
		opcode: 0x06,
	}
}

func (h *LoginGame) GetOpcode() byte {
	return h.opcode
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
		if enterReply == nil || enterReply.GetPartyId() == 0 {
			return &internal.GetPartyReply{Found: false}, nil
		}
		return ic.GetParty(c, &internal.GetPartyRequest{
			WorldId: h.gs.config.WorldId,
			PartyId: enterReply.GetPartyId(),
		})
	}, func(reply *internal.GetPartyReply) error {
		partyReply = reply
		return h.finishLoginGame(ctx, req, enterReply, partyReply)
	})
	promise.
		OnError(func(err error) {
			log.Printf("LoginGame (async): %v", err)
		}).
		Run()
	return nil
}

func (h *LoginGame) finishLoginGame(ctx *core.ClientContext, req *request.LoginGame, reply *internal.EnterGameReply, partyReply *internal.GetPartyReply) error {
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
		SpawnPoint:   uint8(p.GetSpawnPoint()),
		PositionX:    int16(p.GetPositionX()),
		PositionY:    int16(p.GetPositionY()),
		Stance:       uint8(p.GetStance()),
		PartyID:      optionalUint32(reply.GetPartyId()),
		GuildID:      optionalUint32(reply.GetGuildId()),
	}

	character := entity.NewCharacter(ctx.Client, h.gs.characterListener, initData, h.gs)
	character.KeyLayout().LoadKeyLayoutProto(reply.GetKeyLayout())

	character.LoadInventory(reply.GetInventory())
	character.LoadSkills(reply.GetSkills())

	gameClient, ok := ctx.Client.(*client.GameClient)
	if !ok {
		return fmt.Errorf("client is not a GameClient")
	}
	gameClient.SetCharacter(character)

	if h.gs.partyEventConsumer != nil {
		partyID := reply.GetPartyId()
		if partyID != 0 && partyReply != nil && partyReply.GetFound() && partyReply.GetParty() != nil {
			h.gs.partyEventConsumer.ApplyPartySnapshot(partyReply.GetParty())
		}
	}

	mapID := p.GetMapId()
	spawnPoint := uint8(p.GetSpawnPoint())
	mapInstance := h.gs.GetMap(mapID)
	if mapInstance == nil {
		log.Printf("saved map %d not found, falling back to default", mapID)
		defaultMapID, ok := h.gs.resources.NameToMap("\xed\x97\xa4\xeb\x84\xa4\xec\x8b\x9c\xec\x8a\xa4")
		if !ok {
			return fmt.Errorf("default map not found")
		}
		mapID = defaultMapID
		spawnPoint = 0
		mapInstance = h.gs.GetMap(mapID)
		if mapInstance == nil {
			return fmt.Errorf("default map %d not found", mapID)
		}
	}

	if mapInstance.Wz != nil {
		if _, ok := mapInstance.Wz.Portals[spawnPoint]; !ok {
			spawnPoint = 0
		}
	}

	character.Stance = constant.StanceDefaultValue

	targetMapPID := mapInstance.GetActorPID()
	if targetMapPID == nil {
		return fmt.Errorf("MapActor PID not found for map %d", mapID)
	}

	if err := h.gs.characterRuntime.RegisterCharacter(character.GetID(), character.GetName()); err != nil {
		return fmt.Errorf("runtime register: %w", err)
	}

	rootContext := h.gs.GetServer().GetRootContext()
	if rootContext == nil {
		return fmt.Errorf("rootContext not set")
	}

	rootContext.Send(targetMapPID, &g_actor.AddCharacter{
		Character:  character,
		SpawnPoint: spawnPoint,
		Init:       true,
	})

	return nil
}
