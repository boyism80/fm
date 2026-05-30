package server

import (
	"context"
	"fmt"
	"log"

	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/core/async"
	"github.com/boyism80/fm/core/luax"
	"github.com/boyism80/fm/protocol/constant"
	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
	"github.com/boyism80/fm/protocol/request"
	"github.com/boyism80/fm/protocol/response"
	"github.com/boyism80/fm/services/game/client"
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
	case constant.GuildC2SCreate:
		if _, inGuild := ch.GetGuildID(); inGuild {
			_ = ch.Send(&response.GuildGenericMessage{
				Code: constant.GuildResponseAlreadyInGuild,
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
					_ = ch.Send(&response.GuildGenericMessage{
						Code: constant.GuildResponseAlreadyInGuild,
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
	default:
		return nil
	}
}
