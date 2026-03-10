package server

import (
	"fmt"
	"log"

	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/core/luax"
	"github.com/boyism80/fm/game/client"
	"github.com/boyism80/fm/game/wz"
	"github.com/boyism80/fm/protocol/request"
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

	character := client.GetCharacter()
	if character == nil {
		return nil
	}

	// Check blocked inventory
	// TODO: Implement hasBlockedInventory check

	mapInstance := character.GetMap()
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
		// TODO: Check for Mu Lung Dojo skills
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

	skillEntry := character.Skills[req.SkillID]
	if skillEntry == nil {
		if character.Listener != nil {
			character.Listener.OnUpdateStats(nil, true)
		}
		return nil
	}

	if levelData.Cooldown > 0 {
		if skillEntry.IsCooling() {
			if character.Listener != nil {
				character.Listener.OnUpdateStats(nil, true)
			}
			return nil
		}
		skillEntry.StartCooldown(levelData.Cooldown)
	}

	if ctx.LogicActorPID == nil {
		if character.Listener != nil {
			character.Listener.OnUpdateStats(nil, true)
		}
		return nil
	}
	root := luax.GetRootLuaState(ctx.LogicActorPID.String())
	if root == nil {
		log.Printf("No lua root state for actor %s", ctx.LogicActorPID.String())
		if character.Listener != nil {
			character.Listener.OnUpdateStats(nil, true)
		}
		return nil
	}

	// Common validation (e.g. event map check when event system exists). Optional; skip if script missing.
	const commonSkillScript = "script/skill/common.lua"
	if commonResult, commonThread, commonErr := luax.Call(root, commonSkillScript, "on_preactivated_common", character, skillEntry); commonErr == nil {
		defer commonThread.Close()
		if commonResult != nil && commonResult.Type() == lua.LTBool && !lua.LVAsBool(commonResult) {
			if character.Listener != nil {
				character.Listener.OnUpdateStats(nil, true)
			}
			return nil
		}
	}

	scriptPath := fmt.Sprintf("script/skill/%d.lua", req.SkillID)
	result, thread, err := luax.Call(root, scriptPath, "on_preactivated", character, skillEntry)
	if err != nil {
		log.Printf("Skill script not found or failed %s: %v", scriptPath, err)
		if character.Listener != nil {
			character.Listener.OnUpdateStats(nil, true)
		}
		return nil
	}
	if result != nil && result.Type() == lua.LTBool && !lua.LVAsBool(result) {
		thread.Close()
		if character.Listener != nil {
			character.Listener.OnUpdateStats(nil, true)
		}
		return nil
	}

	resumeState, err := luax.Execute(root, thread, ctx.LogicActorPID, "on_activated", character, skillEntry)
	if err != nil {
		thread.Close()
		log.Printf("Failed to execute skill script %s: %v", scriptPath, err)
		if character.Listener != nil {
			character.Listener.OnUpdateStats(nil, true)
		}
		return err
	}
	// Thread may have yielded (e.g. sleep); only close when it finished in this call (ResumeOK).
	// If it yielded, ResumeLua will resume it later and close it in MapActor.resumeLua when done.
	if resumeState == lua.ResumeOK {
		thread.Close()
	}

	if character.Listener != nil {
		character.Listener.OnUpdateStats(nil, true)
	}
	return nil
}
