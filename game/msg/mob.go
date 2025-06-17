// Package msg defines actor messages for the MapleStory private server.
// This file contains monster-related messages that flow between actors.
// These messages handle monster behavior, combat, and lifecycle management.
package msg

import (
	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/game/constant"
	"github.com/boyism80/fm/game/protocol"
	"github.com/boyism80/fm/game/protocol/req"
)

// MobKill instructs a monster to die with specified animation
// Flow: MapActor/GameClientActor -> MobActor
// Trigger: Monster HP reaches 0, admin kill command, map clear
// Side Effects:
//   - Triggers death animation
//   - Calculates and spawns drops
//   - Removes mob from map
//   - Awards EXP to attackers
//
// Related Messages: MapDieMob, MapSpawnItem, MapSpawnMeso
type MobKill struct {
	AnimationType constant.MobDieAnimationType // Death animation type
}

// MobControllerChange notifies mob of controller change
// Flow: MapActor -> MobActor
// Trigger: Player moves out of range, disconnects, or new player comes closer
// Side Effects:
//   - Updates mob's controlling player
//   - Sends control packets to old/new controllers
//   - May affect mob AI behavior
//
// Related Messages: MapBroadcast (for control packets)
type MobControllerChange struct {
	Before *actor.PID // Previous controlling player (nil if none)
	After  *actor.PID // New controlling player (nil if none)
}

// MobMove handles monster movement from controller
// Flow: GameClientActor -> MapActor -> MobActor
// Trigger: Controlling player sends mob movement packet
// Side Effects:
//   - Updates mob position and state
//   - Validates movement legality
//   - Broadcasts movement to nearby players
//
// Related Messages: MapMoveMob, MapBroadcast
type MobMove struct {
	req.MoveMob            // Movement data from protocol
	Sender      *actor.PID // Player actor controlling the mob
}

// MobDamaged applies damage to a monster
// Flow: GameClientActor -> MapActor -> MobActor
// Trigger: Player attacks monster with skill or normal attack
// Side Effects:
//   - Reduces mob HP by damage amount
//   - Triggers death if HP reaches 0
//   - Updates damage statistics
//   - May trigger mob special abilities
//
// Related Messages: MobKill, MapCharacterAttack, MapDieMob
type MobDamaged struct {
	Sender      *actor.PID            // Player actor dealing damage
	CharacterId uint32                // Character ID of attacker
	DamagePairs []protocol.DamagePair // Damage amounts and types
}
