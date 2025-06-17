// Package msg defines actor messages for the MapleStory private server.
// This file contains life entity messages that flow between actors.
// These messages handle health and lifecycle management for living entities.
package msg

// LifeAddHp increases health points for a life entity
// Flow: Any Actor -> LifeActor (Character/Mob/NPC)
// Trigger: Healing potion use, skill effect, regeneration, admin command
// Side Effects:
//   - Increases entity's current HP
//   - Caps HP at maximum value
//   - May trigger visual healing effects
//   - Updates HP display for players
// Related Messages: MobDamaged (opposite effect), CharacterBuiltinMeso (for potion cost)
type LifeAddHp struct {
	Hp int // Amount of HP to add (positive value)
}
