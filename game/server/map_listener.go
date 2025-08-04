package server

import (
	"github.com/boyism80/fm/core/types"
	"github.com/boyism80/fm/game/action"
	"github.com/boyism80/fm/game/constant"
	"github.com/boyism80/fm/game/entity"
	"github.com/boyism80/fm/game/protocol/resp"
)

// MapListenerImpl implements MapListener for game server
type MapListenerImpl struct {
	gameServer *GameServer
}

// NewGameMapListener creates a new GameMapListener instance
func NewGameMapListener(gameServer *GameServer) *MapListenerImpl {
	return &MapListenerImpl{
		gameServer: gameServer,
	}
}

// OnPlayerAdded sends spawn player packet to other players on the map
func (l *MapListenerImpl) OnPlayerAdded(mapID uint32, playerID uint32, character *entity.Character, init bool) {
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

	// 2. Send existing players' info to the new player
	if mapInstance.GetPlayerCount() > 1 { // More than just the new player
		for cid, p := range mapInstance.GetAllPlayers() {
			if cid == playerID {
				continue // Skip the new player
			}

			if ch, ok := p.(*entity.Character); ok {
				// Send existing player's spawn info to the new player
				character.Send(&resp.SpawnPlayer{
					Character:       ch,
					BuffStates:      [4]uint32{},
					Diseases:        [4]uint32{},
					CrushRings:      []*entity.Ring{},
					FriendshipRings: []*entity.Ring{},
					MarriageRings:   []*entity.Ring{},
				}, types.SEND_POLICY_ENCRYPT)
			}
		}
	}

	// 3. Send SpawnPlayer packet to all other players on the map
	spawnPacket := &resp.SpawnPlayer{
		Character:       character,
		BuffStates:      [4]uint32{},
		Diseases:        [4]uint32{},
		CrushRings:      []*entity.Ring{},
		FriendshipRings: []*entity.Ring{},
		MarriageRings:   []*entity.Ring{},
	}
	mapInstance.BroadcastToPlayers(spawnPacket, types.SEND_POLICY_ENCRYPT, playerID)

	// 4. Send NPC spawn packets to the new player (following old server pattern)
	for _, npc := range mapInstance.GetNpcs() {
		if npc, ok := npc.(*entity.Npc); ok {
			// Send SpawnNpc packet
			character.Send(&resp.SpawnNpc{
				NPC:     npc,
				Visible: true,
			}, types.SEND_POLICY_ENCRYPT)

			// Send NpcControl packet
			character.Send(&resp.NpcControl{
				NPC:     npc,
				MiniMap: true,
			}, types.SEND_POLICY_ENCRYPT)
		}
	}

	// 5. Send existing items and meso on the map to the new player (following old server pattern)
	for _, item := range mapInstance.GetItems() {
		// Handle regular items
		if item, ok := item.(entity.Item); ok {
			drop := item.GetDrop()
			if drop != nil {
				// Send SpawnItem packet for existing items
				character.Send(&resp.SpawnItem{
					ID:           drop.Object.OID,
					Animation:    constant.DROP_ITEM_ANIMATION_TYPE_NONE,
					DropType:     drop.DropType,
					Item:         item,
					OwnerID:      drop.Owner,
					SpawnedPoint: drop.SpawnedPoint,
					IsPlayerDrop: true,
				}, types.SEND_POLICY_ENCRYPT)
			}
		}

		// Handle meso separately
		if meso, ok := item.(*entity.Meso); ok {
			drop := meso.GetDrop()
			if drop != nil {
				// Send SpawnMeso packet for existing meso
				character.Send(&resp.SpawnMeso{
					ID:           drop.Object.OID,
					Animation:    constant.DROP_ITEM_ANIMATION_TYPE_NONE,
					DropType:     drop.DropType,
					Count:        meso.Count,
					OwnerID:      drop.Owner,
					Position:     drop.Position,
					SpawnedPoint: drop.SpawnedPoint,
					IsPlayerDrop: true,
				}, types.SEND_POLICY_ENCRYPT)
			}
		}
	}

	// 6. Send existing mobs on the map to the new player (following old server pattern)
	for _, mob := range mapInstance.GetMobs() {
		if mob, ok := mob.(*entity.Mob); ok {
			// Send SpawnMob packet for existing mobs (following old server pattern)
			character.Send(&resp.SpawnMob{
				Mob:       mob,
				SpawnType: constant.MOB_SPAWN_TYPE_NONE, // Use NONE for existing mobs (following old server pattern)
			}, types.SEND_POLICY_ENCRYPT)
		}
	}
}

// OnPlayerRemoved sends leave player packet to other players on the map
func (l *MapListenerImpl) OnPlayerRemoved(mapID uint32, playerID uint32) {
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
func (l *MapListenerImpl) OnPlayerMoved(mapID uint32, playerID uint32, character *entity.Character) {
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

// OnPlayerMove sends move packet with fragments to all other players on the map
func (l *MapListenerImpl) OnPlayerMove(mapID uint32, playerID uint32, character *entity.Character, startPoint types.Vector2[int16], fragments []action.MoveFragment) {
	// Get the map instance
	mapInstance := l.gameServer.GetMap(mapID)
	if mapInstance == nil {
		return
	}

	// Create move packet
	movePacket := &resp.Move{
		Character:  character,
		Fragments:  fragments,
		StartPoint: startPoint,
	}

	// Broadcast to all players on the map except the moving player
	mapInstance.BroadcastToPlayers(movePacket, types.SEND_POLICY_ENCRYPT, playerID)
}

// OnPlayerChat sends chat packet to other players on the map
func (l *MapListenerImpl) OnPlayerChat(mapID uint32, playerID uint32, message string) {
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

// OnItemSpawned sends spawn item packet to all players on the map
func (l *MapListenerImpl) OnItemSpawned(mapID uint32, itemID uint32, item entity.Item, drop *entity.Drop) {
	// Get the map instance
	mapInstance := l.gameServer.GetMap(mapID)
	if mapInstance == nil {
		return
	}

	// Create spawn item packet
	spawnPacket := &resp.SpawnItem{
		ID:           drop.OID,
		Animation:    constant.DROP_ITEM_ANIMATION_TYPE_LOOTING,
		DropType:     drop.DropType,
		Item:         item,
		OwnerID:      drop.Owner,
		SpawnedPoint: drop.SpawnedPoint,
		IsPlayerDrop: true,
	}

	// Broadcast to all players on the map
	mapInstance.BroadcastToAllPlayers(spawnPacket, types.SEND_POLICY_ENCRYPT)
}

// OnMesoSpawned sends spawn meso packet to all players on the map
func (l *MapListenerImpl) OnMesoSpawned(mapID uint32, itemID uint32, meso *entity.Meso) {
	// Get the map instance
	mapInstance := l.gameServer.GetMap(mapID)
	if mapInstance == nil {
		return
	}

	// Create spawn meso packet
	spawnPacket := &resp.SpawnMeso{
		ID:           meso.GetDrop().OID,
		Animation:    constant.DROP_ITEM_ANIMATION_TYPE_LOOTING,
		DropType:     meso.GetDrop().DropType,
		Count:        meso.Count,
		OwnerID:      meso.GetDrop().Owner,
		Position:     meso.GetDrop().Position,
		SpawnedPoint: meso.GetDrop().SpawnedPoint,
		IsPlayerDrop: true,
	}

	// Broadcast to all players on the map
	mapInstance.BroadcastToAllPlayers(spawnPacket, types.SEND_POLICY_ENCRYPT)
}

// OnItemRemoved sends remove item packet to all players on the map
func (l *MapListenerImpl) OnItemRemoved(mapID uint32, itemID uint32, characterID uint32, mode constant.RemoveItemType) {
	// Get the map instance
	mapInstance := l.gameServer.GetMap(mapID)
	if mapInstance == nil {
		return
	}

	// Create remove item packet
	removePacket := &resp.RemoveItem{
		Mode:        mode,
		OID:         itemID,
		CharacterId: characterID,
	}

	// Broadcast to all players on the map
	mapInstance.BroadcastToAllPlayers(removePacket, types.SEND_POLICY_ENCRYPT)
}

// OnMobSpawned sends spawn mob packet to all players on the map
func (l *MapListenerImpl) OnMobSpawned(mapID uint32, mobID uint32, mob *entity.Mob) {
	// Get the map instance
	mapInstance := l.gameServer.GetMap(mapID)
	if mapInstance == nil {
		return
	}

	// Create spawn mob packet
	spawnPacket := &resp.SpawnMob{
		Mob:       mob,
		SpawnType: constant.MOB_SPAWN_TYPE_ANIMATE,
	}

	// Broadcast to all players on the map
	mapInstance.BroadcastToAllPlayers(spawnPacket, types.SEND_POLICY_ENCRYPT)
}

// OnMobRemoved sends remove mob packet to all players on the map
func (l *MapListenerImpl) OnMobRemoved(mapID uint32, mobID uint32, animationType constant.MobDieAnimationType) {
	// Get the map instance
	mapInstance := l.gameServer.GetMap(mapID)
	if mapInstance == nil {
		return
	}

	// Create remove mob packet
	removePacket := &resp.DieMob{
		OID:           mobID,
		AnimationType: animationType,
	}

	// Broadcast to all players on the map
	mapInstance.BroadcastToAllPlayers(removePacket, types.SEND_POLICY_ENCRYPT)
}

// OnMobControllerChange sends StartControlMob packet to new controller (following old server pattern)
func (l *MapListenerImpl) OnMobControllerChange(mob *entity.Mob, before *entity.Character, after *entity.Character) {
	// Send StartControlMob packet to new controller (following old server pattern)
	if after != nil {
		after.Send(&resp.StartControlMob{
			Mob:   mob,
			Aggro: false,
		}, types.SEND_POLICY_ENCRYPT)
	}
	// Note: StopControlMob is commented out in old server, so we skip it here too
}

// OnMobMoved broadcasts mob movement to all players on the map
func (l *MapListenerImpl) OnMobMoved(mapID uint32, mobID uint32, isAggroed bool, centerSplit int8, skill1 uint8, skill2 uint8, skill3 uint8, skill4 uint8, startPoint types.Vector2[int16], movements []action.MoveFragment) {
	// Get the map instance
	mapInstance := l.gameServer.GetMap(mapID)
	if mapInstance == nil {
		return
	}

	// Create move mob packet
	movePacket := &resp.MoveMob{
		IsAggroed:   isAggroed,
		CenterSplit: centerSplit,
		Skill1:      skill1,
		Skill2:      skill2,
		Skill3:      skill3,
		Skill4:      skill4,
		OID:         mobID,
		StartPoint:  startPoint,
		Movements:   movements,
	}

	// Broadcast to all players on the map
	mapInstance.BroadcastToAllPlayers(movePacket, types.SEND_POLICY_ENCRYPT)
}

// OnAttack broadcasts attack to all players on the map
func (l *MapListenerImpl) OnAttack(mapID uint32, characterID uint32, attackInfo action.AttackInfo, skillLevel uint8) {
	// Get the map instance
	mapInstance := l.gameServer.GetMap(mapID)
	if mapInstance == nil {
		return
	}

	// Create attack packet
	attackPacket := &resp.Attack{
		CharacterId: characterID,
		AttackInfo:  attackInfo,
		SkillLevel:  skillLevel,
	}

	// Broadcast to all players on the map except the attacker
	mapInstance.BroadcastToPlayers(attackPacket, types.SEND_POLICY_ENCRYPT, characterID)
}
