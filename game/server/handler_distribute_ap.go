package server

import (
	"fmt"
	"log"
	"math/rand"

	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/game/client"
	"github.com/boyism80/fm/game/constant"
	"github.com/boyism80/fm/protocol/request"
)

// DistributeAP handles DISTRIBUTE_AP packet requests (0x46)
type DistributeAP struct {
	gameServer *GameServer
	opcode     byte
}

func (DistributeAP) New(gameServer *GameServer) *DistributeAP {
	return &DistributeAP{
		gameServer: gameServer,
		opcode:     0x46,
	}
}

func (h *DistributeAP) GetOpcode() byte {
	return h.opcode
}

func (h *DistributeAP) Handle(ctx *core.ClientContext, req *request.DistributeAP) error {
	client, ok := ctx.Client.(*client.GameClient)
	if !ok {
		log.Printf("Client is not a GameClient")
		return fmt.Errorf("client is not a GameClient")
	}

	character := client.GetCharacter()
	if character == nil {
		log.Printf("Character is nil for client")
		return fmt.Errorf("character is nil")
	}

	// Send empty stat update packet first (for client synchronization)
	stats := map[constant.Stat]int32{}
	character.Listener.OnUpdateStats(stats, true)

	// Check if character has remaining AP
	if character.AbilityPoint == 0 {
		return nil
	}

	// Process stat distribution
	statUpdate := map[constant.Stat]int32{}
	success := false

	switch constant.StatType(req.StatType) {
	case constant.STAT_TYPE_STR:
		if character.Str >= constant.STAT_MAX_STR_DEX_INT_LUK {
			return nil
		}
		character.Str++
		statUpdate[constant.STAT_STR] = int32(character.Str)
		success = true

	case constant.STAT_TYPE_DEX:
		if character.Dex >= constant.STAT_MAX_STR_DEX_INT_LUK {
			return nil
		}
		character.Dex++
		statUpdate[constant.STAT_DEX] = int32(character.Dex)
		success = true

	case constant.STAT_TYPE_INT:
		if character.Int >= constant.STAT_MAX_STR_DEX_INT_LUK {
			return nil
		}
		character.Int++
		statUpdate[constant.STAT_INT] = int32(character.Int)
		success = true

	case constant.STAT_TYPE_LUK:
		if character.Luk >= constant.STAT_MAX_STR_DEX_INT_LUK {
			return nil
		}
		character.Luk++
		statUpdate[constant.STAT_LUK] = int32(character.Luk)
		success = true

	case constant.STAT_TYPE_HP:
		if character.HpApUsed >= constant.HP_AP_USED_MAX || character.MaxHp >= constant.STAT_MAX_HP_MP {
			return nil
		}
		// Calculate HP increase based on job
		hpIncrease := h.calculateHPIncrease(character.Class)
		newMaxHp := character.MaxHp + hpIncrease
		if newMaxHp > constant.STAT_MAX_HP_MP {
			newMaxHp = constant.STAT_MAX_HP_MP
		}
		character.MaxHp = newMaxHp
		character.HpApUsed++
		statUpdate[constant.STAT_MAX_HP] = int32(character.MaxHp)
		success = true

	case constant.STAT_TYPE_MP:
		if character.HpApUsed >= constant.HP_AP_USED_MAX || character.MaxMp >= constant.STAT_MAX_HP_MP {
			return nil
		}
		// Calculate MP increase based on job
		mpIncrease := h.calculateMPIncrease(character.Class)
		newMaxMp := character.MaxMp + mpIncrease
		if newMaxMp > constant.STAT_MAX_HP_MP {
			newMaxMp = constant.STAT_MAX_HP_MP
		}
		character.MaxMp = newMaxMp
		character.HpApUsed++
		statUpdate[constant.STAT_MAX_MP] = int32(character.MaxMp)
		success = true

	default:
		// Invalid stat type - send empty stat update
		character.Listener.OnUpdateStats(stats, true)
		return nil
	}

	if success {
		// Decrease AP
		character.AbilityPoint--
		statUpdate[constant.STAT_AVAILABLE_AP] = int32(character.AbilityPoint)

		// Send stat update packet
		character.Listener.OnUpdateStats(statUpdate, true)
	}

	return nil
}

// calculateHPIncrease calculates HP increase based on job
func (h *DistributeAP) calculateHPIncrease(job uint16) uint16 {
	// Beginner
	if job == constant.JOB_BEGINNER_MIN || job == constant.JOB_BEGINNER_1 || job == constant.JOB_BEGINNER_2 {
		return uint16(rand.Intn(5) + 8) // 8-12
	}
	// Warrior
	if job >= constant.JOB_WARRIOR_MIN && job <= constant.JOB_WARRIOR_MAX {
		return uint16(rand.Intn(9) + 12) // 12-20
	}
	// Magician
	if job >= constant.JOB_MAGICIAN_MIN && job <= constant.JOB_MAGICIAN_MAX {
		return uint16(rand.Intn(6) + 6) // 6-11
	}
	// Bowman/Thief
	if (job >= constant.JOB_BOWMAN_MIN && job <= constant.JOB_BOWMAN_MAX) ||
		(job >= constant.JOB_THIEF_MIN && job <= constant.JOB_THIEF_MAX) {
		return uint16(rand.Intn(5) + 14) // 14-18
	}
	// Default (GameMaster)
	return uint16(rand.Intn(51) + 50) // 50-100
}

// calculateMPIncrease calculates MP increase based on job
func (h *DistributeAP) calculateMPIncrease(job uint16) uint16 {
	// Beginner
	if job == constant.JOB_BEGINNER_MIN || job == constant.JOB_BEGINNER_1 || job == constant.JOB_BEGINNER_2 {
		return uint16(rand.Intn(3) + 6) // 6-8
	}
	// Magician
	if job >= constant.JOB_MAGICIAN_MIN && job <= constant.JOB_MAGICIAN_MAX {
		return uint16(rand.Intn(11) + 10) // 10-20
	}
	// Bowman/Thief
	if (job >= constant.JOB_BOWMAN_MIN && job <= constant.JOB_BOWMAN_MAX) ||
		(job >= constant.JOB_THIEF_MIN && job <= constant.JOB_THIEF_MAX) {
		return uint16(rand.Intn(5) + 8) // 8-12
	}
	// Warrior/Soul Master
	if job >= constant.JOB_WARRIOR_MIN && job <= constant.JOB_WARRIOR_MAX {
		return uint16(rand.Intn(4) + 4) // 4-7
	}
	// Default (GameMaster)
	return uint16(rand.Intn(51) + 50) // 50-100
}
