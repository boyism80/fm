package server

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/core/async"
	"github.com/boyism80/fm/core/luax"
	pconst "github.com/boyism80/fm/protocol/constant"
	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
	"github.com/boyism80/fm/protocol/request"
	"github.com/boyism80/fm/protocol/response"
	g_actor "github.com/boyism80/fm/services/game/actor"
	"github.com/boyism80/fm/services/game/client"
	gameconst "github.com/boyism80/fm/services/game/constant"
	"github.com/boyism80/fm/services/game/entity"
	"github.com/boyism80/fm/types"
	lua "github.com/yuin/gopher-lua"
)

type GuildOperation struct {
	gs *GameServer
}

func (GuildOperation) New(gs *GameServer) *GuildOperation {
	return &GuildOperation{
		gs: gs,
	}
}

func (h *GuildOperation) resumeGuildCreate(ch *entity.Character, ok bool) {
	if ch == nil {
		return
	}
	thread := ch.GetDialog()
	if thread == nil {
		return
	}
	mapInstance := ch.GetMap()
	if mapInstance == nil {
		return
	}
	root := mapInstance.GetLuaRoot()
	if root == nil {
		return
	}
	if _, err := luax.Resume(root, thread, lua.LBool(ok)); err != nil {
		log.Printf("GuildOperation(create): failed to resume npc script character=%d: %v", ch.GetID(), err)
	}
}

func (h *GuildOperation) inviterCanInvite(characterID, guildID uint32) bool {
	if characterID == 0 || guildID == 0 {
		return false
	}
	g := h.gs.GetGuildByID(guildID)
	if g == nil {
		return false
	}
	for _, m := range g.GetMembers() {
		if m == nil || m.GetCharacterId() != characterID {
			continue
		}
		rank := m.GetRank()
		if rank == internal.GuildMemberRank_GUILD_MEMBER_RANK_MASTER {
			return true
		}
		if rank == internal.GuildMemberRank_GUILD_MEMBER_RANK_JUNIOR {
			return true
		}
		return false
	}
	return false
}

func (h *GuildOperation) Handle(ctx *core.ClientContext, req *request.GuildOperation) error {
	if h.gs.internalClient == nil {
		return fmt.Errorf("internal client not configured")
	}
	if ctx.ActorContext == nil {
		return fmt.Errorf("guild operation: actor context required")
	}

	gameClient, ok := ctx.Client.(*client.GameClient)
	if !ok {
		return fmt.Errorf("client is not a GameClient")
	}
	ch := gameClient.GetCharacter()
	if ch == nil {
		return fmt.Errorf("character not found")
	}

	worldID := h.gs.config.WorldId
	charID := ch.GetID()

	switch req.Operation {
	case pconst.GuildC2SCreate:
		if _, inGuild := ch.GetGuildID(); inGuild {
			_ = ch.Send(&response.GuildMessage{
				Code: pconst.GuildResponseAlreadyInGuild,
			}, types.SEND_POLICY_ENCRYPT)
			h.resumeGuildCreate(ch, false)
			return nil
		}

		promise := async.NewPromise(ctx.ActorContext, core.InternalRPCPerStepTimeout)
		promise = async.ThenRPC(promise, func(c context.Context) (*internal.CreateGuildReply, error) {
			leader := ch.ToProtoGuildMember(worldID, int32(h.gs.config.ChannelId), internal.GuildMemberRank_GUILD_MEMBER_RANK_MASTER)
			return h.gs.internalClient.CreateGuild(c, &internal.CreateGuildRequest{
				WorldId:   worldID,
				GuildName: req.GuildName,
				Leader:    leader,
			})
		}, func(reply *internal.CreateGuildReply) error {
			if !reply.GetOk() {
				if reply.GetErrorCode() == internal.GuildErrorCode_GUILD_ERROR_ALREADY_IN_GUILD {
					_ = ch.Send(&response.GuildMessage{
						Code: pconst.GuildResponseAlreadyInGuild,
					}, types.SEND_POLICY_ENCRYPT)
				} else {
					log.Printf("GuildOperation(create): failed character=%d code=%v", charID, reply.GetErrorCode())
				}
				h.resumeGuildCreate(ch, false)
				return nil
			}
			guildPb := reply.GetGuild()
			if guildPb == nil || guildPb.GetGuildId() == 0 {
				log.Printf("GuildOperation(create): ok but missing guild character=%d", charID)
				h.resumeGuildCreate(ch, false)
				return nil
			}
			if h.gs.guild != nil {
				h.gs.guild.OnGuildCreated(ch, guildPb)
			}
			h.resumeGuildCreate(ch, true)
			return nil
		})
		promise.OnError(func(err error) {
			log.Printf("GuildOperation(create) async error: %v", err)
			h.resumeGuildCreate(ch, false)
		}).Run()
		return nil
	case pconst.GuildC2SInvite:
		guildID, inGuild := ch.GetGuildID()
		if !inGuild || !h.inviterCanInvite(charID, guildID) {
			return nil
		}
		targetName := req.TargetName
		if targetName == "" {
			return nil
		}
		if h.gs.characterRuntime == nil {
			return nil
		}
		targetID, ok := h.gs.characterRuntime.GetCharacterIDByName(targetName)
		if !ok || targetID == 0 {
			_ = ch.Send(&response.GuildMessage{
				Code: pconst.GuildResponseNotInChannel,
			}, types.SEND_POLICY_ENCRYPT)
			return nil
		}
		if mapInstance := ch.GetMap(); mapInstance != nil {
			if target := mapInstance.GetPlayer(targetID); target != nil {
				if _, targetInGuild := target.GetGuildID(); targetInGuild {
					_ = ch.Send(&response.GuildMessage{
						Code: pconst.GuildResponseAlreadyInGuild,
					}, types.SEND_POLICY_ENCRYPT)
					return nil
				}
				now := time.Now()
				for id, expiresAt := range target.GuildInvites {
					if !now.Before(expiresAt) {
						delete(target.GuildInvites, id)
					}
				}
				if len(target.GuildInvites) > 0 {
					ch.Listener.OnMessage(ch, gameconst.MSG_PINK_TEXT, gameconst.GuildInviteTargetBusyMessage)
					return nil
				}
				target.GuildInvites[guildID] = now.Add(gameconst.GuildInviteDuration)
				target.Listener.OnGuildInvite(target, guildID, ch.GetName())
				return nil
			}
		}
		h.gs.EnsureSend(nil, targetID, &g_actor.DeliverGuildInvite{
			CharacterID:        targetID,
			InviterCharacterID: charID,
			GuildID:            guildID,
			InviterName:        ch.GetName(),
		})
		return nil
	default:
		return nil
	}
}
