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

	if commonThread, commonErr := luax.NewThread(root, constant.SkillHookScriptPath); commonErr == nil {
		commonResult, commonErr := luax.Call(commonThread, "on_activating", ch, skillEntry, params)
		if commonErr != nil {
			log.Printf("Skill common on_activating: %v", commonErr)
		} else if commonResult != nil && commonResult.Type() == lua.LTBool && !lua.LVAsBool(commonResult) {
			ch.Listener.OnUpdateStats(ch, nil, true)
			return nil
		}
	}

	scriptPath := fmt.Sprintf("script/skill/%d.lua", req.SkillID)
	thread, err := luax.NewThread(root, scriptPath)
	if err != nil {
		thread = nil
	} else {
		result, err := luax.Call(thread, fmt.Sprintf("on_activating_%d", req.SkillID), ch, skillEntry, params)
		if err != nil {
			log.Printf("Skill script not found or failed %s: %v", scriptPath, err)
			ch.Listener.OnUpdateStats(ch, nil, true)
			return nil
		}
		if result != nil && result.Type() == lua.LTBool && !lua.LVAsBool(result) {
			ch.Listener.OnUpdateStats(ch, nil, true)
			return nil
		}
	}

	if levelData.ItemCon != 0 && levelData.ItemConNo > 0 {
		itemID := uint32(levelData.ItemCon)
		count := uint16(levelData.ItemConNo)
		if !ch.RemoveByItemIDCount(itemID, count) {
			ch.Listener.OnUpdateStats(ch, nil, true)
			return nil
		}
	}

	commonActivatedThread, commonActivatedErr := luax.NewThread(root, constant.SkillHookScriptPath)
	if commonActivatedErr != nil {
		log.Printf("common on_activated: %v", commonActivatedErr)
		ch.Listener.OnUpdateStats(ch, nil, true)
		return nil
	}
	commonActivatedResult, commonActivatedErr := luax.Call(commonActivatedThread, "on_activated", ch, skillEntry, params)
	if commonActivatedErr != nil {
		log.Printf("common on_activated: %v", commonActivatedErr)
		ch.Listener.OnUpdateStats(ch, nil, true)
		return nil
	}
	if commonActivatedResult != nil && commonActivatedResult.Type() == lua.LTBool && !lua.LVAsBool(commonActivatedResult) {
		ch.Listener.OnUpdateStats(ch, nil, true)
		return nil
	}

	if thread != nil {
		luax.SetConfiguration(thread, luax.Configuration{
			ActorContext: ctx.ActorContext,
		})
		if _, err := luax.Call(thread, fmt.Sprintf("on_activated_%d", req.SkillID), ch, skillEntry, params); err != nil {
			log.Printf("Failed to execute skill script %s: %v", scriptPath, err)
			ch.Listener.OnUpdateStats(ch, nil, true)
			return err
		}
	}

	var dir *uint8
	if req.MagnetMobData != nil {
		d := req.MagnetMobData.Direction
		dir = &d
	}
	ch.Listener.OnShowSkillEffect(ch, pconst.SkillEffectTypeCast, req.SkillID, req.SkillLevel, dir)
	ch.Listener.OnUpdateStats(ch, nil, true)
	return nil
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
