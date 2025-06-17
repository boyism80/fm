// Package msg defines actor messages for the MapleStory private server.
// This file contains character-related messages that flow between actors.
// These messages handle character interactions, scripting, and game mechanics.
package msg

import (
	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/game/entity"
	lua "github.com/yuin/gopher-lua"
)

// CharacterLootFailed notifies character that item looting failed
// Flow: ItemActor/MesoActor -> GameClientActor
// Trigger: Item pickup attempt failed (out of range, no space, ownership, etc.)
// Side Effects: Sends failure packet to client, unlocks player action
// Related Messages: CharacterItemLooting, MapItemLoot, ItemLooting
type CharacterLootFailed struct {
	OID uint32 // Object ID of item that failed to loot
}

// CharacterItemLooting notifies character of available item for pickup
// Flow: ItemActor -> GameClientActor
// Trigger: Item pickup validation passed, item ready to be added to inventory
// Side Effects:
//   - Validates inventory space and adds item if possible
//   - Sends ItemLooted back to ItemActor with success/failure
//   - Updates inventory UI if successful
//   - Shows item gain animation
//
// Related Messages: CharacterLootFailed, ItemLooted, ItemLooting
type CharacterItemLooting struct {
	PID  *actor.PID  // Item actor that is being looted
	OID  uint32      // Object ID of item being looted
	Item entity.Item // Item data to be picked up
}

// CharacterMesoLooting notifies character of available meso for pickup
// Flow: MesoActor -> GameClientActor
// Trigger: Meso pickup validation passed, meso ready to be added to character
// Side Effects:
//   - Validates meso cap and adds meso if possible
//   - Sends ItemLooted back to MesoActor with success/failure
//   - Updates meso display if successful
//   - Shows meso gain animation
//
// Related Messages: ItemLooted, ItemLooting
type CharacterMesoLooting struct {
	PID  *actor.PID // Meso actor that is being looted
	OID  uint32     // Object ID of meso being looted
	Meso int32      // Amount of meso to be picked up
}

// CharacterMapChanged notifies character of map transition completion
// Flow: MapActor -> GameClientActor
// Trigger: Character successfully entered new map
// Side Effects:
//   - Updates character's current map reference
//   - Initializes character state in new map
//   - May trigger map-specific events
//
// Related Messages: MapChange, EnterMap, LeaveMap
type CharacterMapChanged struct {
	MID        uint32     // Map ID of new map
	Map        *actor.PID // Map actor reference
	Init       bool       // Whether this is initial map entry
	SpawnPoint uint8      // Spawn point used for entry
}

// RunScript initiates execution of a Lua script
// Flow: NPCActor/GameClientActor -> GameClientActor -> LuaActor
// Trigger: NPC interaction, quest trigger, command execution
// Side Effects:
//   - Creates new Lua execution context
//   - Runs script in isolated environment
//   - May yield for player input
//
// Related Messages: ResumeScript, LuaRun, LuaYield
type RunScript struct {
	Script string // Lua script code to execute
}

// ResumeScript continues execution of a yielded Lua script
// Flow: GameClientActor -> LuaActor
// Trigger: Player responds to dialog, input, or other script interaction
// Side Effects:
//   - Resumes Lua execution from yield point
//   - Processes player response
//   - May yield again or complete
//
// Related Messages: RunScript, LuaResume, LuaYield
type ResumeScript struct {
	L    *lua.LState  // Lua state to resume
	Args []lua.LValue // Arguments from player response
}

// CharacterKillMob notifies character of monster kill
// Flow: MobActor -> GameClientActor
// Trigger: Monster dies and character gets credit
// Side Effects:
//   - Awards experience points
//   - Updates kill statistics
//   - May trigger quest progress
//
// Related Messages: MapDieMob, CharacterBuiltinDialog (for quest updates)
type CharacterKillMob struct {
	OID   uint32 // Object ID of killed monster
	MobID uint32 // Monster template ID
}

// CharacterBuiltinDialog shows a simple dialog to player
// Flow: LuaActor -> GameClientActor
// Trigger: Lua script calls dialog function
// Side Effects:
//   - Displays dialog window to player
//   - May pause script execution
//   - Waits for player acknowledgment
//
// Related Messages: CharacterBuiltinDialogList, CharacterBuiltinDialogAccept
type CharacterBuiltinDialog struct {
	NPC     uint32 // NPC ID showing dialog
	Message string // Dialog message text
	Prev    bool   // Show previous button
	Next    bool   // Show next button
}

// CharacterBuiltinDialogList shows a selection dialog to player
// Flow: LuaActor -> GameClientActor
// Trigger: Lua script calls selection dialog function
// Side Effects:
//   - Displays selection dialog window
//   - Pauses script execution until selection
//   - Returns selected index to script
//
// Related Messages: CharacterBuiltinDialog, ResumeScript
type CharacterBuiltinDialogList struct {
	NPC        uint32   // NPC ID showing dialog
	Message    string   // Dialog message text
	Selections []string // List of selectable options
}

// CharacterBuiltinDialogAccept shows an accept/cancel dialog
// Flow: LuaActor -> GameClientActor
// Trigger: Lua script calls accept dialog function
// Side Effects:
//   - Displays accept/cancel dialog
//   - Pauses script execution until response
//   - Returns boolean result to script
//
// Related Messages: CharacterBuiltinDialog, ResumeScript
type CharacterBuiltinDialogAccept struct {
	NPC          uint32 // NPC ID showing dialog
	Message      string // Dialog message text
	EnableEscape bool   // Allow escape/cancel option
}

// CharacterBuiltinDialogInput shows a text input dialog
// Flow: LuaActor -> GameClientActor
// Trigger: Lua script calls input dialog function
// Side Effects:
//   - Displays text input dialog
//   - Pauses script execution until input
//   - Returns input string to script
//
// Related Messages: CharacterBuiltinDialog, ResumeScript
type CharacterBuiltinDialogInput struct {
	NPC     uint32 // NPC ID showing dialog
	Message string // Dialog prompt message
}

// CharacterBuiltinChat sends a chat message from script
// Flow: LuaActor -> GameClientActor -> MapActor (broadcast)
// Trigger: Lua script calls chat function
// Side Effects:
//   - Sends chat message to nearby players
//   - May highlight message differently
//   - Can bypass chat history recording
//
// Related Messages: MapBroadcast (for chat distribution)
type CharacterBuiltinChat struct {
	Lua               *lua.LState // Lua context for response
	Message           string      // Chat message text
	Highlight         bool        // Highlight message in chat
	DontRecordHistory bool        // Skip chat history recording
}

// CharacterBuiltinName retrieves character name for script
// Flow: LuaActor -> GameClientActor -> LuaActor
// Trigger: Lua script requests character name
// Side Effects: Returns character name to Lua script
// Related Messages: Used within script execution context
type CharacterBuiltinName struct {
	Lua *lua.LState // Lua context for response
}

// CharacterBuiltinMeso retrieves character meso amount for script
// Flow: LuaActor -> GameClientActor -> LuaActor
// Trigger: Lua script requests current meso amount
// Side Effects: Returns current meso amount to Lua script
// Related Messages: CharacterBuiltinAddMeso, CharacterBuiltinRemoveMeso
type CharacterBuiltinMeso struct {
	Lua *lua.LState // Lua context for response
}

// CharacterBuiltinRemoveMeso removes meso from character via script
// Flow: LuaActor -> GameClientActor
// Trigger: Lua script removes meso (shop purchase, quest cost, etc.)
// Side Effects:
//   - Deducts meso from character
//   - Sends meso update packet to client
//   - Returns success/failure to script
//
// Related Messages: CharacterBuiltinMeso, CharacterBuiltinAddMeso
type CharacterBuiltinRemoveMeso struct {
	Lua   *lua.LState // Lua context for response
	Count int32       // Amount of meso to remove
}

// CharacterBuiltinAddMeso adds meso to character via script
// Flow: LuaActor -> GameClientActor
// Trigger: Lua script adds meso (quest reward, shop sale, etc.)
// Side Effects:
//   - Adds meso to character
//   - Sends meso update packet to client
//   - Returns success/failure to script
//
// Related Messages: CharacterBuiltinMeso, CharacterBuiltinRemoveMeso
type CharacterBuiltinAddMeso struct {
	Lua   *lua.LState // Lua context for response
	Count int32       // Amount of meso to add
}
