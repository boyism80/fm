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

	if levelData.Cooldown > 0 {
		if skillEntry.IsCooling() {
			ch.Listener.OnUpdateStats(ch, nil, true)
			return nil
		}
		skillEntry.StartCooldown(levelData.Cooldown)
	}

	root := mapInstance.GetLuaRoot()
	if root == nil {
		log.Printf("No lua root state for map %d", mapInstance.GetMapID())
		ch.Listener.OnUpdateStats(ch, nil, true)
		return nil
	}

	params := buildActiveSkillParams(root, mapInstance, req)

	if !h.runActivatingHooks(root, ch, req, skillEntry, params) {
		ch.Listener.OnUpdateStats(ch, nil, true)
		return nil
	}

	if levelData.ItemCon != 0 && levelData.ItemConNo > 0 {
		if !ch.Inventory.RemoveByItemIDCount(uint32(levelData.ItemCon), uint16(levelData.ItemConNo)) {
			ch.Listener.OnUpdateStats(ch, nil, true)
			return nil
		}
	}

	if !h.runActivatedHooks(ctx, root, ch, req, skillEntry, params) {
		ch.Listener.OnUpdateStats(ch, nil, true)
		return nil
	}
	h.showActiveSkillEffect(ch, req)
	return nil
}

func (h *ActiveSkill) runActivatingHooks(root *lua.LState, ch *entity.Character, req *request.ActiveSkill, skillEntry *entity.SkillEntry, params lua.LValue) bool {
	commonThread, err := luax.NewThread(root, constant.SkillHookScriptPath)
	if err != nil {
		log.Printf("Skill common on_activating: %v", err)
		return true
	}
	commonRet, err := luax.Call(commonThread, "on_activating", ch, skillEntry, params)
	if !skillHookAllowed(commonRet, err, "on_activating", true) {
		return false
	}
	scriptPath := fmt.Sprintf("script/skill/%d.lua", req.SkillID)
	skillThread, err := luax.NewThread(root, scriptPath)
	if err != nil {
		return true
	}
	skillHook := "on_activating"
	skillRet, err := luax.Call(skillThread, skillHook, ch, skillEntry, params)
	if err != nil {
		log.Printf("Skill script not found or failed %s: %v", scriptPath, err)
		return false
	}
	return skillHookAllowed(skillRet, err, skillHook, true)
}

func (h *ActiveSkill) runActivatedHooks(ctx *core.ClientContext, root *lua.LState, ch *entity.Character, req *request.ActiveSkill, skillEntry *entity.SkillEntry, params lua.LValue) bool {
	commonThread, err := luax.NewThread(root, constant.SkillHookScriptPath)
	if err != nil {
		log.Printf("common on_activated: %v", err)
		return false
	}
	commonRet, err := luax.Call(commonThread, "on_activated", ch, skillEntry, params)
	if err != nil {
		log.Printf("common on_activated: %v", err)
		return false
	}
	if !skillHookAllowed(commonRet, err, "on_activated", true) {
		return false
	}
	scriptPath := fmt.Sprintf("script/skill/%d.lua", req.SkillID)
	skillThread, err := luax.NewThread(root, scriptPath)
	if err != nil {
		return true
	}
	luax.SetConfiguration(skillThread, luax.Configuration{
		ActorContext: ctx.ActorContext,
	})
	skillHook := "on_activated"
	if _, err := luax.Call(skillThread, skillHook, ch, skillEntry, params); err != nil {
		log.Printf("Failed to execute skill script %s: %v", scriptPath, err)
		return false
	}
	return true
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

func buildActiveSkillParams(L *lua.LState, mapInstance *entity.Map, req *request.ActiveSkill) lua.LValue {
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
