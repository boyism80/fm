// Package msg defines actor messages for the MapleStory private server.
// This file contains Lua script execution messages that flow between actors.
// These messages handle NPC scripts, quest logic, and game event scripting.
package msg

import (
	"github.com/asynkron/protoactor-go/actor"
	lua "github.com/yuin/gopher-lua"
)

// LuaRun initiates execution of a Lua script function
// Flow: GameClientActor/NPCActor -> LuaActor
// Trigger: NPC interaction, quest trigger, command execution, event script
// Side Effects:
//   - Creates new Lua execution context
//   - Loads and executes specified script file and function
//   - May yield for player input (dialogs, etc.)
//   - Returns results to calling actor
//
// Related Messages: LuaResume, LuaYield, RunScript
type LuaRun struct {
	PID      *actor.PID   // Actor that requested script execution
	FileName string       // Lua script file to load
	FuncName string       // Function name to execute
	Params   []lua.LValue // Parameters to pass to function
}

// LuaResume continues execution of a yielded Lua script
// Flow: GameClientActor -> LuaActor
// Trigger: Player responds to dialog, input, or other script interaction
// Side Effects:
//   - Resumes Lua execution from yield point
//   - Processes player response as script parameters
//   - May yield again or complete execution
//   - Sends results back to requesting actor
//
// Related Messages: LuaRun, LuaYield, ResumeScript
type LuaResume struct {
	PID    *actor.PID   // Actor that is resuming the script
	Lua    *lua.LState  // Lua state to resume from yield
	Params []lua.LValue // Parameters from player response
}

// LuaYield pauses script execution and waits for player input
// Flow: LuaActor -> GameClientActor
// Trigger: Lua script calls yield function (dialog, input, etc.)
// Side Effects:
//   - Pauses script execution
//   - Sends dialog/input request to player
//   - Stores Lua state for later resumption
//   - Waits for player response
//
// Related Messages: LuaRun, LuaResume, CharacterBuiltinDialog*
type LuaYield struct {
	Lua *lua.LState // Lua state that yielded execution
}
