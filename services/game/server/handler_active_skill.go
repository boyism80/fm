package server

import (
	"fmt"
	"log"

	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/core/luax"
	"github.com/boyism80/fm/protocol/request"
	"github.com/boyism80/fm/services/game/client"
	"github.com/boyism80/fm/services/game/entity"
	"github.com/boyism80/fm/services/game/wz"
	lua "github.com/yuin/gopher-lua"
)

type ActiveSkill struct {
	gs     *GameServer
	opcode byte
}

func (ActiveSkill) New(gs *GameServer) *ActiveSkill {
	return &ActiveSkill{
		gs:     gs,
		opcode: 0x4A,
	}
}

func (h *ActiveSkill) GetOpcode() byte {
	return h.opcode
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

	if ctx.LogicActorPID == nil {
		ch.Listener.OnUpdateStats(ch, nil, true)
		return nil
	}
	root := luax.GetRootLuaState(ctx.LogicActorPID.String())
	if root == nil {
		log.Printf("No lua root state for actor %s", ctx.LogicActorPID.String())
		ch.Listener.OnUpdateStats(ch, nil, true)
		return nil
	}

	params := buildActiveSkillParams(root, mapInstance, req)

	if commonResult, commonThread, commonErr := luax.Call(root, commonSkillScriptPath, "on_activating", ch, skillEntry, params); commonErr == nil {
		defer commonThread.Close()
		if commonResult != nil && commonResult.Type() == lua.LTBool && !lua.LVAsBool(commonResult) {
			ch.Listener.OnUpdateStats(ch, nil, true)
			return nil
		}
	}

	scriptPath := fmt.Sprintf("script/skill/%d.lua", req.SkillID)
	result, thread, err := luax.Call(root, scriptPath, luax.SkillScriptHookName("on_activating", req.SkillID), ch, skillEntry, params)
	if err != nil {
		log.Printf("Skill script not found or failed %s: %v", scriptPath, err)
		ch.Listener.OnUpdateStats(ch, nil, true)
		return nil
	}
	if result != nil && result.Type() == lua.LTBool && !lua.LVAsBool(result) {
		thread.Close()
		ch.Listener.OnUpdateStats(ch, nil, true)
		return nil
	}

	if levelData.ItemCon != 0 && levelData.ItemConNo > 0 {
		itemID := uint32(levelData.ItemCon)
		count := uint16(levelData.ItemConNo)
		if !ch.RemoveByItemIDCount(itemID, count) {
			ch.Listener.OnUpdateStats(ch, nil, true)
			return nil
		}
	}

	commonActivatedResult, commonActivatedThread, commonActivatedErr := luax.Call(root, commonSkillScriptPath, "on_activated", ch, skillEntry, params)
	if commonActivatedErr != nil {
		log.Printf("common on_activated: %v", commonActivatedErr)
		thread.Close()
		ch.Listener.OnUpdateStats(ch, nil, true)
		return nil
	}
	if commonActivatedThread != nil {
		defer commonActivatedThread.Close()
	}
	if commonActivatedResult != nil && commonActivatedResult.Type() == lua.LTBool && !lua.LVAsBool(commonActivatedResult) {
		thread.Close()
		ch.Listener.OnUpdateStats(ch, nil, true)
		return nil
	}

	luax.SetConfiguration(thread, luax.Configuration{
		ActorContext: ctx.ActorContext,
	})
	if _, err := luax.CallThread(thread, luax.SkillScriptHookName("on_activated", req.SkillID), ch, skillEntry, params); err != nil {
		log.Printf("Failed to execute skill script %s: %v", scriptPath, err)
		ch.Listener.OnUpdateStats(ch, nil, true)
		return err
	}

	if !ch.IsHidden() {
		var dir *uint8
		if req.MagnetMobData != nil {
			d := req.MagnetMobData.Direction
			dir = &d
		}
		ch.Listener.OnShowBuffEffect(ch, 1, req.SkillID, req.SkillLevel, dir)
	}
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
