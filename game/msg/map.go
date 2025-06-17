// Package msg defines actor messages for the MapleStory private server.
// This file contains map-related messages that flow between actors.
// These messages handle map operations, player movement, and object management.
package msg

import (
	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/common/types"
	"github.com/boyism80/fm/game/constant"
	"github.com/boyism80/fm/game/entity"
	"github.com/boyism80/fm/game/protocol"
	"github.com/boyism80/fm/game/protocol/req"
	"github.com/boyism80/fm/game/protocol/resp"
)

// EnterMap adds a character to the map
// Flow: GameClientActor/MapActor -> MapActor
// Trigger: Player login, warp, portal use
// Side Effects:
//   - Adds player to map's character list
//   - Spawns player for other players in map
//   - Sends map data to entering player
//   - May send CharacterMapChanged to player
//
// Related Messages: LeaveMap, MapChange, CharacterMapChanged
type EnterMap struct {
	ID         uint32     // Character ID entering map
	PID        *actor.PID // GameClientActor PID
	SpawnPoint uint8      // Spawn point to use
	Init       bool       // Whether this is initial map entry
}

// LeaveMap removes a character from the map
// Flow: GameClientActor -> MapActor
// Trigger: Player logout, warp, disconnect
// Side Effects:
//   - Removes player from map's character list
//   - Despawns player for other players in map
//   - Cleans up player-related map state
//
// Related Messages: EnterMap, MapChange
type LeaveMap struct {
	ID uint32 // Character ID leaving map
}

// QueryAllPIDs requests all PIDs currently in the map
// Flow: Any Actor -> MapActor -> Reply with MapPIDList
// Trigger: Need to broadcast to all players, admin commands
// Side Effects: None (read-only query)
// Related Messages: MapPIDList (response)
type QueryAllPIDs struct {
	ReplyTo *actor.PID // Actor to send response to
}

// MapPIDList response containing all PIDs in the map
// Flow: MapActor -> Requesting Actor
// Trigger: Response to QueryAllPIDs
// Side Effects: None (response message)
// Related Messages: QueryAllPIDs (request)
type MapPIDList struct {
	Targets []*actor.PID // List of all actor PIDs in map
}

// MapNotifyCharacterWarped notifies map that a character has warped
// Flow: GameClientActor -> MapActor
// Trigger: Character completes warp transition
// Side Effects: Updates character state, may trigger map events
// Related Messages: MapChange, EnterMap
type MapNotifyCharacterWarped struct {
	Sender *actor.PID // Character actor that warped
}

// MapBroadcast sends a message to multiple actors in the map
// Flow: Any Actor -> MapActor -> broadcasts to target actors
// Trigger: Movement, chat, skill use, any action needing broadcast
// Side Effects: Sends message to all/filtered actors in map
// Related Messages: Used with almost all other messages
type MapBroadcast struct {
	Sender     *actor.PID              // Original sender
	Pivot      types.Vector2[int16]    // Center point for range-based broadcast
	Message    any                     // Message to broadcast
	ExceptSelf bool                    // Exclude sender from broadcast
	Excepts    map[*actor.PID]struct{} // Additional actors to exclude
}

// MapSpawnItem creates a new item drop in the map
// Flow: Any Actor -> MapActor -> creates ItemActor -> broadcasts to players
// Trigger: Monster death, player drop, quest reward
// Side Effects:
//   - Creates new ItemActor instance
//   - Broadcasts spawn packet to nearby players
//   - Starts item expiration timer
//
// Related Messages: MapRemoveItem, ItemLooting
type MapSpawnItem struct {
	Item    entity.Item // The item to spawn
	Owner   *actor.PID  // Player who can pick it up first (nil for FFA)
	OwnerID uint32      // Character ID of owner
}

// MapSpawnItems creates multiple item drops simultaneously
// Flow: Any Actor -> MapActor -> creates multiple ItemActors
// Trigger: Monster death with multiple drops, bag opening
// Side Effects: Creates multiple ItemActors, broadcasts multiple spawn packets
// Related Messages: MapSpawnItem, MapRemoveItem
type MapSpawnItems struct {
	Items    []entity.Dropable    // List of items to spawn
	Position types.Vector2[int16] // Drop position
	Owner    *actor.PID           // Owner actor (nil for FFA)
}

// MapSpawnMeso creates a meso drop in the map
// Flow: Any Actor -> MapActor -> creates MesoActor -> broadcasts to players
// Trigger: Monster death, player drop, quest reward
// Side Effects:
//   - Creates new MesoActor instance
//   - Broadcasts meso spawn packet to nearby players
//   - Starts meso expiration timer
//
// Related Messages: MapItemLoot (for meso pickup)
type MapSpawnMeso struct {
	Count        int32                // Amount of meso to drop
	SpawnedPoint types.Vector2[int16] // Initial spawn position
	DestPoint    types.Vector2[int16] // Final position after drop animation
	Owner        *actor.PID           // Owner actor (nil for FFA)
	OwnerID      uint32               // Character ID of owner
	DropType     constant.DropType    // FFA, Party, or Owned drop
}

// MapItemLoot handles item pickup attempts
// Flow: GameClientActor -> MapActor -> ItemActor -> back to GameClientActor
// Trigger: Player clicks on dropped item
// Side Effects:
//   - Validates pickup eligibility
//   - Adds item to player inventory
//   - Removes item from map if successful
//
// Related Messages: MapRemoveItem, CharacterItemLooting
type MapItemLoot struct {
	Actor       *actor.PID           // Player actor attempting pickup
	OID         uint32               // Object ID of item to pick up
	CharacterId uint32               // Character ID attempting pickup
	Position    types.Vector2[int16] // Player position (for range check)
}

// MapRemoveItem removes an item from the map
// Flow: ItemActor/MapActor -> broadcasts to players
// Trigger: Item picked up, expired, or manually removed
// Side Effects:
//   - Destroys ItemActor
//   - Broadcasts item removal packet to nearby players
//   - Cleans up item state
//
// Related Messages: MapSpawnItem, MapItemLoot
type MapRemoveItem struct {
	Actor       *actor.PID           // Actor requesting removal
	OID         uint32               // Object ID of item to remove
	CharacterId uint32               // Character ID (if picked up by player)
	Mode        resp.RemoveItemType  // Removal animation type
	Position    types.Vector2[int16] // Item position
}

// MapChange handles character map transitions
// Flow: GameClientActor -> MapActor -> target MapActor
// Trigger: Portal use, warp command, teleport skill
// Side Effects:
//   - Removes character from current map
//   - Adds character to target map
//   - Updates character's map reference
//
// Related Messages: EnterMap, LeaveMap, MapNotifyCharacterWarped
type MapChange struct {
	CharacterId uint32     // Character ID changing maps
	Sender      *actor.PID // Character actor
	To          *actor.PID // Target map actor
	SpawnPoint  uint8      // Spawn point in target map
}

// SendMessage sends a message to a specific object in the map
// Flow: Any Actor -> MapActor -> target object actor
// Trigger: Direct communication with specific map object
// Side Effects: Forwards message to target object
// Related Messages: Depends on the forwarded message
type SendMessage struct {
	OID     uint32 // Target object ID
	Message any    // Message to send
}

// MapRepeatSpawnMobs triggers respawn cycle for all mobs in map
// Flow: Timer/MapActor -> MapActor -> spawns MobActors
// Trigger: Periodic respawn timer, map initialization
// Side Effects:
//   - Spawns missing mobs based on spawn points
//   - Creates new MobActor instances
//   - Broadcasts mob spawn packets
//
// Related Messages: MapSpawningMob, MapSpawnedMob
type MapRepeatSpawnMobs struct {
}

// MapDieMob handles monster death
// Flow: MobActor -> MapActor -> broadcasts to players
// Trigger: Monster HP reaches 0, killed by player
// Side Effects:
//   - Triggers drop calculation and spawning
//   - Broadcasts mob death packet
//   - Schedules mob respawn
//   - Awards EXP to players
//
// Related Messages: MapSpawnItem, MapSpawnMeso, MapRepeatSpawnMobs
type MapDieMob struct {
	Sender        *actor.PID                   // Mob actor that died
	OID           uint32                       // Mob object ID
	Position      types.Vector2[int16]         // Death position
	AnimationType constant.MobDieAnimationType // Death animation type
}

// MapClearMobs removes all mobs from the map
// Flow: Admin/MapActor -> MapActor -> all MobActors
// Trigger: Admin command, map reset, special event
// Side Effects:
//   - Kills all mobs with specified animation
//   - Clears mob spawn timers
//   - Broadcasts death packets for all mobs
//
// Related Messages: MapDieMob, MobKill
type MapClearMobs struct {
	AnimationType constant.MobDieAnimationType // Animation for all mob deaths
}

// MapMoveMob handles monster movement
// Flow: MobActor -> MapActor -> broadcasts to players
// Trigger: Monster AI movement, controller change
// Side Effects:
//   - Updates mob position
//   - Broadcasts movement packet to nearby players
//   - Updates mob controller if needed
//
// Related Messages: MobMove, MapBroadcast
type MapMoveMob struct {
	req.MoveMob            // Movement data from protocol
	Sender      *actor.PID // Mob actor that moved
}

// MapSpawningMob initiates mob spawn process
// Flow: MapActor -> MapActor (internal) -> creates MobActor
// Trigger: Respawn timer, map initialization
// Side Effects:
//   - Creates new MobActor instance
//   - Initializes mob state and position
//   - Prepares for spawn broadcast
//
// Related Messages: MapSpawnedMob, MapRepeatSpawnMobs
type MapSpawningMob struct {
	Position types.Vector2[int16] // Spawn position
	MobId    uint32               // Monster ID to spawn
}

// MapSpawnedMob confirms mob spawn completion
// Flow: MobActor -> MapActor -> broadcasts to players
// Trigger: MobActor successfully created and initialized
// Side Effects:
//   - Registers mob in map's object registry
//   - Broadcasts mob spawn packet to nearby players
//   - Starts mob AI and behavior
//
// Related Messages: MapSpawningMob, EnterMap
type MapSpawnedMob struct {
	PID *actor.PID // Newly created mob actor
	OID uint32     // Mob object ID
}

// MapCharacterAttack handles character attack actions
// Flow: GameClientActor -> MapActor -> target MobActors -> broadcasts
// Trigger: Player uses attack skill or normal attack
// Side Effects:
//   - Calculates damage to target mobs
//   - Broadcasts attack animation to nearby players
//   - May trigger mob death if HP reaches 0
//
// Related Messages: MobDamaged, MapDieMob, MapBroadcast
type MapCharacterAttack struct {
	Sender      *actor.PID           // Character actor attacking
	CharacterId uint32               // Character ID
	AttackInfo  protocol.AttackInfo  // Attack data (skill, targets, etc.)
	Position    types.Vector2[int16] // Attack position
}
