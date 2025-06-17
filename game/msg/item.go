// Package msg defines actor messages for the MapleStory private server.
// This file contains item-related messages that flow between actors.
// These messages handle item drops, pickup, and lifecycle management.
package msg

import (
	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/common/types"
	"github.com/boyism80/fm/game/constant"
)

// Spawn instructs an item to spawn in the map
// Flow: MapActor -> ItemActor
// Trigger: Item created and ready to be visible
// Side Effects:
//   - Makes item visible to players
//   - Starts item expiration timer
//   - Broadcasts spawn packet to nearby players
//
// Related Messages: MapSpawnItem, ItemDestroy
type Spawn struct {
}

// ItemLooting handles item pickup attempt by a player
// Flow: GameClientActor -> MapActor -> ItemActor/MesoActor
// Trigger: Player clicks on dropped item to pick it up
// Side Effects:
//   - Validates pickup eligibility (range, ownership, inventory space)
//   - If valid: sends CharacterItemLooting/CharacterMesoLooting to player
//   - If invalid: sends CharacterLootFailed to player
//   - Sets looting flag to prevent duplicate attempts
//
// Related Messages: CharacterItemLooting, CharacterMesoLooting, CharacterLootFailed
type ItemLooting struct {
	Actor       *actor.PID           // Player actor attempting pickup
	Position    types.Vector2[int16] // Player position for range validation
	CharacterId uint32               // Character ID attempting pickup
}

// ItemLooted confirms the result of an item pickup attempt
// Flow: GameClientActor -> ItemActor/MesoActor
// Trigger: Player finishes processing item/meso pickup (success or failure)
// Side Effects:
//   - If successful: reduces item count, may destroy ItemActor if count reaches 0
//   - If failed: resets looting flag for retry
//   - Sends MapRemoveItem if item/meso is fully consumed
//
// Related Messages: ItemLooting, MapRemoveItem, CharacterItemLooting, CharacterMesoLooting
type ItemLooted struct {
	Success     bool   // Whether pickup was successful
	CharacterId uint32 // Character ID that attempted pickup
	Count       int32  // Amount picked up (for stackable items)
}

// ItemDropTypeChanged updates item ownership/drop type
// Flow: Timer/MapActor -> ItemActor/MesoActor
// Trigger: Drop ownership timer expires, party changes, admin command
// Side Effects:
//   - Changes who can pick up the item
//   - Updates item visual appearance if needed
//   - May reset pickup timers
//
// Related Messages: MapSpawnItem, ItemLooting
type ItemDropTypeChanged struct {
	Mode constant.DropType // New drop type (FFA, Party, Owned)
}

// ItemDestroy removes an item from the map with animation
// Flow: Timer/MapActor -> ItemActor/MesoActor -> MapActor
// Trigger: Item expires, picked up, or manually destroyed
// Side Effects:
//   - Plays destruction animation
//   - Removes item from map
//   - Broadcasts removal packet to nearby players
//   - Destroys ItemActor/MesoActor
//
// Related Messages: MapRemoveItem, Spawn
type ItemDestroy struct {
	Animation constant.DropItemAnimationType // Destruction animation type
}
