// Package msg defines actor messages for the MapleStory private server.
// This file contains general object-related messages that flow between actors.
// These messages handle basic object operations like movement and warping.
package msg

import (
	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/common/types"
)

// CharacterMove notifies of character position change
// Flow: GameClientActor -> MapActor -> broadcasts to nearby players
// Trigger: Player moves using keyboard/mouse input
// Side Effects:
//   - Updates character position in map
//   - Broadcasts movement packet to nearby players
//   - May trigger foothold calculations
//   - Updates mob controller assignments if needed
//
// Related Messages: MapBroadcast, MobControllerChange
type CharacterMove struct {
	Position types.Vector2[int16] // New character position
}

// Warped notifies that an object has completed a warp/teleport
// Flow: GameClientActor/NPCActor/MobActor -> MapActor
// Trigger: Object completes teleportation or map transition
// Side Effects:
//   - Confirms warp completion
//   - Updates object state in new location
//   - May trigger post-warp events or scripts
//
// Related Messages: MapChange, MapNotifyCharacterWarped, EnterMap
type Warped struct {
	Sender *actor.PID // Actor that completed the warp
}
