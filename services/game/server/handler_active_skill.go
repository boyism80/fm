package server

import (
	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/core/luax"
	"github.com/boyism80/fm/protocol/request"
	"github.com/boyism80/fm/services/game"
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
	ch, err := h.gs.character(ctx)
	if err != nil {
		return nil
	}

	var direction *uint8
	if req.MagnetMobData != nil {
		d := req.MagnetMobData.Direction
		direction = &d
	}
	used := ch.Skills().UseActive(ctx.ActorContext, req.SkillID, req.SkillLevel, direction, func(L *lua.LState) lua.LValue {
		return h.skillParams(L, ch.GetMap(), req)
	})
	if !used {
		ch.Listener().OnUpdateStats(ch, nil, true)
	}
	return nil
}

func (h *ActiveSkill) skillParams(L *lua.LState, mapInstance *game.Map, req *request.ActiveSkill) lua.LValue {
	params := L.NewTable()
	if req.MagnetMobData != nil {
		magnet := L.NewTable()
		mobs := L.NewTable()
		for _, entry := range req.MagnetMobData.Mobs {
			mob := mapInstance.GetMob(entry.OID)
			if mob == nil {
				continue
			}
			entryTbl := L.NewTable()
			entryTbl.RawSetString("mob", luax.NewLuable(L, mob))
			entryTbl.RawSetString("oid", lua.LNumber(entry.OID))
			entryTbl.RawSetString("success", lua.LBool(entry.Success))
			mobs.Append(entryTbl)
		}
		magnet.RawSetString("mobs", mobs)
		magnet.RawSetString("direction", lua.LNumber(req.MagnetMobData.Direction))
		params.RawSetString("magnet", magnet)
	}
	return params
}
