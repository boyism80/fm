// Package msg defines actor messages for the MapleStory private server.
// This file contains NPC-related messages that flow between actors.
// These messages handle NPC interactions and behavior.
package msg

import "github.com/asynkron/protoactor-go/actor"

// NpcClick handles player interaction with an NPC
// Flow: GameClientActor -> MapActor -> NPCActor -> GameClientActor (via script)
// Trigger: Player clicks on an NPC or uses NPC interaction key
// Side Effects:
//   - Initiates NPC script execution
//   - May open shop, dialog, or quest interface
//   - Creates Lua execution context for NPC script
//   - May trigger character state changes
// Related Messages: RunScript, CharacterBuiltinDialog, LuaRun
type NpcClick struct {
	Sender *actor.PID // Player actor that clicked the NPC
}
