package server

import (
	"github.com/boyism80/fm/common/types"
	"github.com/boyism80/fm/game/entity"
	"github.com/boyism80/fm/game/protocol/resp"
)

// GameMapListener implements MapListener for game server
type GameMapListener struct {
	gameServer *GameServer
}

// NewGameMapListener creates a new GameMapListener instance
func NewGameMapListener(gameServer *GameServer) *GameMapListener {
	return &GameMapListener{
		gameServer: gameServer,
	}
}

// OnPlayerAdded sends spawn player packet to other players on the map
func (l *GameMapListener) OnPlayerAdded(mapID uint32, playerID uint32, character *entity.Character, init bool) {
	// Get the map instance
	mapInstance := l.gameServer.GetMap(mapID)
	if mapInstance == nil {
		return
	}

	// 1. Send Login/Warp packet to the new player based on init status
	if init {
		// First time entering the game
		loginPacket := &resp.Login{
			Character: character,
		}
		character.Send(loginPacket, types.SEND_POLICY_ENCRYPT)
	} else {
		// Warping to a new map
		warpPacket := &resp.Warp{
			Character: character,
			Channel:   0,
		}
		character.Send(warpPacket, types.SEND_POLICY_ENCRYPT)
	}

	// 2. Send SpawnPlayer packet to all other players on the map
	spawnPacket := &resp.SpawnPlayer{
		Character:       character,
		BuffStates:      [4]uint32{},
		Diseases:        [4]uint32{},
		CrushRings:      []*entity.Ring{},
		FriendshipRings: []*entity.Ring{},
		MarriageRings:   []*entity.Ring{},
	}
	mapInstance.BroadcastToPlayers(spawnPacket, types.SEND_POLICY_ENCRYPT, playerID)
}

// OnPlayerRemoved sends leave player packet to other players on the map
func (l *GameMapListener) OnPlayerRemoved(mapID uint32, playerID uint32) {
	// Get the map instance
	mapInstance := l.gameServer.GetMap(mapID)
	if mapInstance == nil {
		return
	}

	// Send LeavePlayer packet to all other players on the map
	leavePacket := &resp.LeavePlayer{
		ID: playerID,
	}
	mapInstance.BroadcastToPlayers(leavePacket, types.SEND_POLICY_ENCRYPT, 0)
}

// OnPlayerMoved sends move player packet to other players on the map
func (l *GameMapListener) OnPlayerMoved(mapID uint32, playerID uint32, character *entity.Character) {
	// Get the map instance
	mapInstance := l.gameServer.GetMap(mapID)
	if mapInstance == nil {
		return
	}

	// TODO: Create and send move player packet
	// movePacket := &resp.MovePlayer{
	//     Character: character,
	// }
	// mapInstance.BroadcastToPlayers(movePacket, types.SEND_POLICY_ENCRYPT, playerID)
}

// OnPlayerChat sends chat packet to other players on the map
func (l *GameMapListener) OnPlayerChat(mapID uint32, playerID uint32, message string) {
	// Get the map instance
	mapInstance := l.gameServer.GetMap(mapID)
	if mapInstance == nil {
		return
	}

	// TODO: Create and send chat packet
	// chatPacket := &resp.NormalChat{
	//     Message: message,
	// }
	// mapInstance.BroadcastToPlayers(chatPacket, types.SEND_POLICY_ENCRYPT, playerID)
}
