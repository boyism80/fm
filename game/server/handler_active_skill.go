package server

import (
	"fmt"
	"log"

	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/core/luax"
	"github.com/boyism80/fm/game/client"
	"github.com/boyism80/fm/game/wz"
	"github.com/boyism80/fm/protocol/request"
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

	character := client.GetCharacter()
	if character == nil {
		return nil
	}

	// Check blocked inventory
	// TODO: Implement hasBlockedInventory check

	mapID := character.GetMap()
	mapInstance := h.gs.GetMap(mapID)
	if mapInstance == nil {
		if character.Listener != nil {
			character.Listener.OnUpdateStats(nil, true)
		}
		return nil
	}

	// Get skill
	var wzSkill *wz.Skill
	if character.Context != nil {
		resources := character.Context.GetResources()
		if resources != nil {
			wzSkill = resources.GetSkill(req.SkillID)
		}
	}

	if wzSkill == nil {
		if character.Listener != nil {
			character.Listener.OnUpdateStats(nil, true)
		}
		return nil
	}

	// Validate skill level
	skillLevel := character.GetTotalSkillLevel(req.SkillID)
	if skillLevel <= 0 || skillLevel != int(req.SkillLevel) {
		// TODO: Check for Mu Lung Dojo and Pyramid skills
		// For now, reject if skill level doesn't match
		if character.Listener != nil {
			character.Listener.OnUpdateStats(nil, true)
		}
		return nil
	}

	// Get skill level data
	levelData := wzSkill.GetLevelData(int(skillLevel))
	if levelData == nil {
		if character.Listener != nil {
			character.Listener.OnUpdateStats(nil, true)
		}
		return nil
	}

	// Check MP recovery skill HP requirement
	// TODO: Check if effect.isMPRecovery() and HP < 10%

	// Check cooldown
	if levelData.Cooldown > 0 && !character.Admin {
		if character.IsSkillCooling(req.SkillID) {
			if character.Listener != nil {
				character.Listener.OnUpdateStats(nil, true)
			}
			return nil
		}
		character.AddCooldown(req.SkillID, levelData.Cooldown)
	}

	skillEntry := character.SkillsMap[req.SkillID]
	if skillEntry == nil {
		if character.Listener != nil {
			character.Listener.OnUpdateStats(nil, true)
		}
		return nil
	}

	pid := ctx.LogicActorID
	root := luax.GetRootLuaState(pid)
	if root == nil {
		log.Printf("No lua root state for actor %s", pid)
		if character.Listener != nil {
			character.Listener.OnUpdateStats(nil, true)
		}
		return nil
	}

	scriptPath := fmt.Sprintf("script/skill/%d.lua", req.SkillID)
	thread, err := luax.NewThread(root, scriptPath)
	if err != nil {
		log.Printf("Skill script not found or failed to load %s: %v", scriptPath, err)
		if character.Listener != nil {
			character.Listener.OnUpdateStats(nil, true)
		}
		return nil
	}

	if err := h.gs.ExecuteScript(root, thread, "on_active", character, skillEntry); err != nil {
		log.Printf("Failed to execute skill script %s: %v", scriptPath, err)
		if character.Listener != nil {
			character.Listener.OnUpdateStats(nil, true)
		}
		return err
	}

	if character.Listener != nil {
		character.Listener.OnUpdateStats(nil, true)
	}
	return nil
}
