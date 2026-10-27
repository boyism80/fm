package server

import (
	"fmt"
	"log"
	"time"

	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/core/clock"
	"github.com/boyism80/fm/core/luax"
	"github.com/boyism80/fm/maple/constant"
	"github.com/boyism80/fm/protocol/dto"
	"github.com/boyism80/fm/protocol/request"
	"github.com/boyism80/fm/services/game"
	"github.com/boyism80/fm/services/game/client"
	lua "github.com/yuin/gopher-lua"
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
	if controller, _ := controllerTable.GetController(mob); controller != ch {
		// The client flags a move it sends after it was told to stop; a move without the flag means it still thinks it controls the mob.
		if !req.ControlStopped {
			ch.Listener().OnStopControlMob(ch, mob)
		}
		return nil
	}

	startPoint := mob.GetPosition()

	for _, mnt := range req.Movements {
		if move, ok := mnt.(*dto.AbsoluteLifeMovement); ok {
			mob.Place(move.Position)
		}

		mob.SetStance(mnt.GetStance())
	}

	selectedSkillID := uint32(0)
	selectedSkillLevel := uint8(0)
	needChoice := !constant.MobSkillCastActionInRange(req.Action)
	if req.ActiveSkill && needChoice {
		if skill := mob.Skills().Choice(ch); skill != nil {
			selectedSkillID = skill.Slot().SkillID
			selectedSkillLevel = skill.Slot().Level
		}
	} else if !needChoice {
		if skill := mob.Skills().TakePending(uint32(req.SkillId), req.SkillLevel, clock.Now()); skill != nil {
			delay := time.Duration(int(req.SkillDelay)*100+450) * time.Millisecond
			controllerID := ch.GetID()
			mob.RemoveTimer(mobSkillScriptTimerKey)
			mob.AddTimer(mobSkillScriptTimerKey, delay, false, func() {
				if !mob.IsAlive() {
					return
				}
				h.runScript(ctx, mapInstance, mob, mapInstance.GetPlayer(controllerID), skill)
			})
		}
	}

	ch.Listener().OnControlMoveMob(ch, mob, req.MovementId, req.ActiveSkill, uint16(min(mob.GetMp(), 65535)), selectedSkillID, selectedSkillLevel)
	ch.Listener().OnMobMoved(
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
	mapInstance.OnMoved(mob, startPoint)

	return nil
}

func (h *MoveMob) runScript(ctx *core.ClientContext, mapInstance *game.Map, mob *game.Mob, controller *game.Character, skill *game.MobSkill) {
	root := mapInstance.GetLuaRoot()
	if root == nil {
		return
	}
	scriptPath := fmt.Sprintf("script/mob/skill/%d.lua", skill.Slot().SkillID)
	thread, err := luax.NewThread(root, scriptPath)
	if err != nil {
		log.Printf("mob skill script %s: %v", scriptPath, err)
		return
	}
	luax.SetConfiguration(thread, luax.Configuration{
		ActorContext: ctx.ActorContext,
	})
	var target interface{}
	if controller != nil {
		target = controller
	}
	hook := "on_mob_skill"
	luax.CallAsync(ctx.ActorContext, root, thread, hook, mob, target, skill).Do(func([]lua.LValue) error {
		mob.ConsumeSkillMP(skill)
		return nil
	}).OnError(func(err error) {
		log.Printf("mob skill %s: %v", hook, err)
	})
}
