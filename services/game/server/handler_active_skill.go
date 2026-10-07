package server

import (
	"fmt"
	"log"

	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/core/luax"
	pconst "github.com/boyism80/fm/protocol/constant"
	"github.com/boyism80/fm/protocol/request"
	"github.com/boyism80/fm/services/game/client"
	"github.com/boyism80/fm/services/game/constant"
	"github.com/boyism80/fm/services/game/entity"
	"github.com/boyism80/fm/services/game/wz"
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

	if ch.Spectating() {
		ch.Listener.OnUpdateStats(ch, nil, true)
		return nil
	}

	mapInstance := ch.GetMap()
	if mapInstance == nil {
		ch.Listener.OnUpdateStats(ch, nil, true)
		return nil
	}

	var wzSkill *wz.Skill
	if ch.GameWorld != nil {
		resources := ch.GameWorld.GetResources()
		if resources != nil {
			wzSkill = resources.GetSkill(req.SkillID)
		}
	}

	if wzSkill == nil {
		ch.Listener.OnUpdateStats(ch, nil, true)
		return nil
	}

	skillLevel := ch.GetTotalSkillLevel(req.SkillID)
	if skillLevel <= 0 || skillLevel != int(req.SkillLevel) {
		ch.Listener.OnUpdateStats(ch, nil, true)
		return nil
	}

	levelData := wzSkill.GetLevelData(int(skillLevel))
	if levelData == nil {
		ch.Listener.OnUpdateStats(ch, nil, true)
		return nil
	}

	skillEntry := ch.Skills.Get(req.SkillID)
	if skillEntry == nil {
		ch.Listener.OnUpdateStats(ch, nil, true)
		return nil
	}

	if levelData.Cooldown > 0 && skillEntry.IsCooling() {
		ch.Listener.OnUpdateStats(ch, nil, true)
		return nil
	}

	if constant.SkillID(req.SkillID) == constant.SkillMysticDoor && mapInstance.Wz.Limits(constant.FieldLimitMysticDoor) {
		ch.Listener.OnUpdateStats(ch, nil, true)
		return nil
	}

	root := mapInstance.GetLuaRoot()
	if root == nil {
		log.Printf("No lua root state for map %d", mapInstance.GetMapID())
		ch.Listener.OnUpdateStats(ch, nil, true)
		return nil
	}

	params := h.skillParams(root, mapInstance, req)

	if ch.CallSkillHook(ctx.ActorContext, skillEntry, "on_activating", params) == false {
		ch.Listener.OnUpdateStats(ch, nil, true)
		return nil
	}

	if levelData.ItemCon != 0 && levelData.ItemConNo > 0 {
		if !ch.Inventory.RemoveByItemIDCount(uint32(levelData.ItemCon), uint16(levelData.ItemConNo)) {
			ch.Listener.OnUpdateStats(ch, nil, true)
			return nil
		}
	}

	if ch.CallSkillHook(ctx.ActorContext, skillEntry, "on_activated", params) == false {
		ch.Listener.OnUpdateStats(ch, nil, true)
		return nil
	}
	if levelData.Cooldown > 0 {
		skillEntry.StartCooldown(levelData.Cooldown)
	}
	h.showActiveSkillEffect(ch, req)
	return nil
}

func (h *ActiveSkill) showActiveSkillEffect(ch *entity.Character, req *request.ActiveSkill) {
	var dir *uint8
	if req.MagnetMobData != nil {
		d := req.MagnetMobData.Direction
		dir = &d
	}
	ch.Listener.OnShowSkillEffect(ch, pconst.SkillEffectTypeCast, req.SkillID, req.SkillLevel, dir)
	ch.Listener.OnUpdateStats(ch, nil, true)
}

func (h *ActiveSkill) skillParams(L *lua.LState, mapInstance *entity.Map, req *request.ActiveSkill) lua.LValue {
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
