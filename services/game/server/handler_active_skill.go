package server

import (
	"fmt"
	"log"

	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/core/luax"
	"github.com/boyism80/fm/protocol/request"
	"github.com/boyism80/fm/services/game"
	"github.com/boyism80/fm/services/game/client"
	lua "github.com/yuin/gopher-lua"
)

type ActiveSkill struct {
	gs *GameServer
}

func (ActiveSkill) New(gs *GameServer) *ActiveSkill {
	return &ActiveSkill{
		gs: gs,
	}
}

func (h *ActiveSkill) Handle(ctx *core.ClientContext, req *request.ActiveSkill) error {
	client, ok := ctx.Client.(*client.GameClient)
	if !ok {
		log.Printf("Client is not a GameClient")
		return fmt.Errorf("client is not a GameClient")
	}

	ch := client.GetCharacter()
	if ch == nil {
		return nil
	}

	var direction *uint8
	if req.MagnetMobData != nil {
		d := req.MagnetMobData.Direction
		direction = &d
	}
	used := ch.Skills.UseActive(ctx.ActorContext, req.SkillID, req.SkillLevel, direction, func(L *lua.LState) lua.LValue {
		return h.skillParams(L, ch.GetMap(), req)
	})
	if !used {
		ch.Listener.OnUpdateStats(ch, nil, true)
	}
	return nil
}

func (h *ActiveSkill) skillParams(L *lua.LState, mapInstance *game.Map, req *request.ActiveSkill) lua.LValue {
	params := L.NewTable()
	if req.MagnetMobData != nil {
		magnet := L.NewTable()
		mobs := L.NewTable()
		for i, entry := range req.MagnetMobData.Mobs {
			mob := mapInstance.GetMob(entry.OID)
			if mob != nil {
				entryTbl := L.NewTable()
				entryTbl.RawSetString("mob", luax.NewLuable(L, mob))
				entryTbl.RawSetString("oid", lua.LNumber(entry.OID))
				entryTbl.RawSetString("success", lua.LBool(entry.Success))
				mobs.RawSetInt(i+1, entryTbl)
			}
		}
		magnet.RawSetString("mobs", mobs)
		magnet.RawSetString("direction", lua.LNumber(req.MagnetMobData.Direction))
		params.RawSetString("magnet", magnet)
	}
	return params
}
