package server

import (
	"fmt"
	"log"
	"time"

	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/core/luax"
	"github.com/boyism80/fm/protocol/dto"
	"github.com/boyism80/fm/protocol/request"
	"github.com/boyism80/fm/services/game/client"
	"github.com/boyism80/fm/services/game/constant"
	"github.com/boyism80/fm/services/game/entity"
)

const mobSkillScriptTimerKey = "mobSkillScript"

type MoveMob struct {
	gs *GameServer
}

func (MoveMob) New(gs *GameServer) *MoveMob {
	return &MoveMob{
		gs: gs,
	}
}

func (h *MoveMob) Handle(ctx *core.ClientContext, req *request.MoveMob) error {
	client, ok := ctx.Client.(*client.GameClient)
	if !ok {
		log.Printf("Client is not a GameClient")
		return fmt.Errorf("client is not a GameClient")
	}

	ch := client.GetCharacter()
	if ch == nil {
		log.Printf("Character not found for client")
		return fmt.Errorf("character not found")
	}

	mapInstance := ch.GetMap()
	if mapInstance == nil {
		log.Printf("Character not on a map")
		return fmt.Errorf("map not found")
	}

	mob := mapInstance.GetMob(req.OID)
	if mob == nil {
		return nil
	}

	controllerTable := mapInstance.GetControllerTable()
	controller, exists := controllerTable.GetController(mob)
	if !exists {
		log.Printf("No controller found for mob %d", req.OID)
		return fmt.Errorf("no controller found for this mob")
	}

	if controller.GetID() != ch.GetID() {
		return nil
	}

	startPoint := mob.Position

	for _, mnt := range req.Movements {
		if move, ok := mnt.(*dto.AbsoluteLifeMovement); ok {
			mob.Position = move.Position
		}

		mob.Stance = mnt.GetStance()
	}

	selectedSkillID := uint32(0)
	selectedSkillLevel := uint8(0)
	needChoice := !constant.MobSkillCastActionInRange(req.Action)
	if req.ActiveSkill || !needChoice {
		if needChoice {
			if skill := mob.Skills.Choice(ch); skill != nil {
				selectedSkillID = skill.Slot.SkillID
				selectedSkillLevel = skill.Slot.Level
			}
		} else {
			delay := time.Duration(int(req.SkillDelay)*100+450) * time.Millisecond
			mob.RemoveTimer(mobSkillScriptTimerKey)
			mob.AddTimer(mobSkillScriptTimerKey, delay, false, func() {
				if mob == nil || !mob.IsAlive() {
					return
				}
				h.runScript(ctx, mapInstance, mob, ch, req)
			})
		}
	}

	ch.Listener.OnControlMoveMob(ch, mob, req.MovementId, req.ActiveSkill, uint16(min(mob.GetMp(), 65535)), selectedSkillID, selectedSkillLevel)
	ch.Listener.OnMobMoved(
		ch,
		mob,
		req.ActiveSkill,
		req.CenterSplit,
		req.SkillId,
		req.SkillLevel,
		req.Unknown,
		req.SkillDelay,
		startPoint,
		req.Movements,
	)

	return nil
}

func (h *MoveMob) runScript(ctx *core.ClientContext, mapInstance *entity.Map, mob *entity.Mob, controller *entity.Character, req *request.MoveMob) {
	if req.SkillId == 0 || req.SkillLevel == 0 {
		return
	}
	if mob.Skills == nil {
		return
	}
	skill := mob.Skills.Get(uint32(req.SkillId), req.SkillLevel)
	if skill == nil {
		return
	}
	root := mapInstance.GetLuaRoot()
	if root == nil {
		return
	}
	scriptPath := fmt.Sprintf("script/mob/skill/%d.lua", req.SkillId)
	thread, err := luax.NewThread(root, scriptPath)
	if err != nil {
		log.Printf("mob skill script %s: %v", scriptPath, err)
		return
	}
	luax.SetConfiguration(thread, luax.Configuration{
		ActorContext: ctx.ActorContext,
	})
	hook := fmt.Sprintf("on_mob_skill_%d", req.SkillId)
	if _, err := luax.Call(thread, hook, mob, controller, skill); err != nil {
		log.Printf("mob skill %s: %v", hook, err)
	}
	h.consumeMobSkillMp(mob, skill)
}

func (h *MoveMob) consumeMobSkillMp(mob *entity.Mob, skill *entity.MobSkill) {
	if mob == nil || skill == nil || skill.LevelData == nil {
		return
	}
	mpCon := skill.LevelData.MpCon
	if mpCon <= 0 {
		return
	}
	floor := mob.GetMaxMp()
	if floor > 500 {
		floor = 500
	}
	current := mob.GetMp()
	var next uint32
	if current > uint32(mpCon) {
		next = current - uint32(mpCon)
	} else {
		next = 0
	}
	if next < floor {
		next = floor
	}
	mob.SetMp(next, false)
}
