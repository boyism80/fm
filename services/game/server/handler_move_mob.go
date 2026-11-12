package server

import (
	"fmt"
	"log"
	"time"

	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/maple/constant"
	"github.com/boyism80/fm/protocol/dto"
	"github.com/boyism80/fm/protocol/request"
)

type MoveMob struct {
	gs *GameServer
}

func (MoveMob) New(gs *GameServer) *MoveMob {
	return &MoveMob{
		gs: gs,
	}
}

func (h *MoveMob) Handle(ctx *core.ClientContext, req *request.MoveMob) error {
	ch := h.gs.character(ctx)
	if ch == nil {
		return nil
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
		delay := time.Duration(int(req.SkillDelay)*100+450) * time.Millisecond
		mob.Skills().Cast(uint32(req.SkillId), req.SkillLevel, delay, ch.GetID())
	}

	ch.Listener().OnControlMoveMob(ch, mob, req.MovementId, req.ActiveSkill, uint16(min(mob.ControllerMP(), 65535)), selectedSkillID, selectedSkillLevel)
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
