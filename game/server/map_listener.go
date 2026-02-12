package server

import (
	"github.com/boyism80/fm/game/constant"
	"github.com/boyism80/fm/game/entity"
	"github.com/boyism80/fm/protocol/dto"
	"github.com/boyism80/fm/protocol/response"
	"github.com/boyism80/fm/types"
)

// MapListenerImpl implements MapListener for game server
type MapListenerImpl struct {
	gs *GameServer
}

// NewGameMapListener creates a new GameMapListener instance
func NewGameMapListener(gs *GameServer) *MapListenerImpl {
	return &MapListenerImpl{
		gs: gs,
	}
}

// OnPlayerAdded sends spawn player packet to other players on the map
func (l *MapListenerImpl) OnPlayerAdded(mapID uint32, playerID uint32, character *entity.Character, init bool) {
	// Get the map instance
	mapInstance := l.gs.GetMap(mapID)
	if mapInstance == nil {
		return
	}

	// 1. Send Login/Warp packet to the new player based on init status
	if init {
		// First time entering the game
		// Convert entity to DTO (with full data for Login)
		characterDTO := character.ToFullDTO()

		loginPacket := &response.Login{
			Character: characterDTO,
		}
		character.Send(loginPacket, types.SEND_POLICY_ENCRYPT)
	} else {
		// Warping to a new map
		// Convert entity to DTO
		characterDTO := character.ToDTO()
		warpPacket := &response.Warp{
			Character: characterDTO,
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
				if ch.IsHidden() && !character.HasRoleAtLeast(ch.GetRole()) {
					continue
				}
				// Convert entity to DTO
				chDTO := ch.ToDTO()
				// Send existing player's spawn info to the new player
				character.Send(&response.SpawnPlayer{
					Character:       chDTO,
					BuffStates:      [4]uint32{},
					Diseases:        [4]uint32{},
					CrushRings:      entity.RingsToDTO(ch.Rings.Left),
					FriendshipRings: entity.RingsToDTO(ch.Rings.Mid),
					MarriageRings:   entity.RingsToDTO(ch.Rings.Right),
				}, types.SEND_POLICY_ENCRYPT)
			}
		}
	}

	// 3. Send SpawnPlayer packet to all other players on the map
	// Convert entity to DTO
	characterDTO := character.ToDTO()
	spawnPacket := &response.SpawnPlayer{
		Character:       characterDTO,
		BuffStates:      [4]uint32{},
		Diseases:        [4]uint32{},
		CrushRings:      entity.RingsToDTO(character.Rings.Left),
		FriendshipRings: entity.RingsToDTO(character.Rings.Mid),
		MarriageRings:   entity.RingsToDTO(character.Rings.Right),
	}
	mapInstance.Broadcast(character, spawnPacket, types.SEND_POLICY_ENCRYPT, playerID)

	for _, npc := range mapInstance.GetNpcs() {
		if npc, ok := npc.(*entity.Npc); ok {
			npcDTO := npc.ToDTO()
			character.Send(&response.SpawnNpc{
				NPC:     npcDTO,
				Visible: true,
			}, types.SEND_POLICY_ENCRYPT)

			character.Send(&response.NpcControl{
				NPC:     npcDTO,
				MiniMap: true,
			}, types.SEND_POLICY_ENCRYPT)
		}
	}

	for _, item := range mapInstance.GetItems() {
		if item, ok := item.(entity.Item); ok {
			drop := item.GetDrop()
			if drop != nil {
				character.Send(&response.SpawnItem{
					ID:           drop.Object.OID,
					Animation:    constant.DROP_ITEM_ANIMATION_TYPE_NONE,
					DropType:     drop.DropType,
					ItemModel:    item.GetModel(),
					Expiration:   item.GetExpiration(),
					Position:     drop.Object.Position,
					OwnerID:      drop.Owner,
					SpawnedPoint: drop.SpawnedPoint,
					IsPlayerDrop: true,
				}, types.SEND_POLICY_ENCRYPT)
			}
		}

		if meso, ok := item.(*entity.Meso); ok {
			drop := meso.GetDrop()
			if drop != nil {
				character.Send(&response.SpawnMeso{
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

	for _, mob := range mapInstance.GetMobs() {
		if mob, ok := mob.(*entity.Mob); ok {
			mobDTO := mob.ToDTO()
			character.Send(&response.SpawnMob{
				Mob:       mobDTO,
				SpawnType: constant.MOB_SPAWN_TYPE_NONE,
			}, types.SEND_POLICY_ENCRYPT)
		}
	}
}

// OnPlayerRemoved sends leave player packet to other players on the map
func (l *MapListenerImpl) OnPlayerRemoved(mapID uint32, playerID uint32) {
	// Get the map instance
	mapInstance := l.gs.GetMap(mapID)
	if mapInstance == nil {
		return
	}

	// Send LeavePlayer packet to all other players on the map
	leavePacket := &response.LeavePlayer{
		ID: playerID,
	}
	mapInstance.BroadcastToPlayers(leavePacket, types.SEND_POLICY_ENCRYPT, 0)
}

// OnPlayerMoved sends move player packet to other players on the map
func (l *MapListenerImpl) OnPlayerMoved(mapID uint32, playerID uint32, character *entity.Character) {
	// Get the map instance
	mapInstance := l.gs.GetMap(mapID)
	if mapInstance == nil {
		return
	}

	// TODO: Create and send move player packet
	// movePacket := &response.MovePlayer{
	//     Character: character,
	// }
	// mapInstance.BroadcastToPlayers(movePacket, types.SEND_POLICY_ENCRYPT, playerID)
}

// OnPlayerMove sends move packet with fragments to all other players on the map
func (l *MapListenerImpl) OnPlayerMove(mapID uint32, playerID uint32, character *entity.Character, startPoint types.Vector2[int16], fragments []dto.MoveFragment) {
	// Get the map instance
	mapInstance := l.gs.GetMap(mapID)
	if mapInstance == nil {
		return
	}

	// Convert entity to DTO
	characterDTO := character.ToDTO()
	// Create move packet
	movePacket := &response.Move{
		Character:  characterDTO,
		Fragments:  fragments,
		StartPoint: startPoint,
	}

	mapInstance.Broadcast(character, movePacket, types.SEND_POLICY_ENCRYPT, playerID)
}

// OnPlayerChat sends chat packet to other players on the map
func (l *MapListenerImpl) OnPlayerChat(mapID uint32, playerID uint32, message string) {
	// Get the map instance
	mapInstance := l.gs.GetMap(mapID)
	if mapInstance == nil {
		return
	}

	// TODO: Create and send chat packet
	// chatPacket := &response.NormalChat{
	//     Message: message,
	// }
	// mapInstance.BroadcastToPlayers(chatPacket, types.SEND_POLICY_ENCRYPT, playerID)
}

// OnItemSpawned sends spawn item packet to all players on the map
func (l *MapListenerImpl) OnItemSpawned(mapID uint32, itemID uint32, item entity.Item, drop *entity.Drop) {
	// Get the map instance
	mapInstance := l.gs.GetMap(mapID)
	if mapInstance == nil {
		return
	}

	dropObj := item.GetDrop()
	// Create spawn item packet
	spawnPacket := &response.SpawnItem{
		ID:           drop.OID,
		Animation:    constant.DROP_ITEM_ANIMATION_TYPE_LOOTING,
		DropType:     drop.DropType,
		ItemModel:    item.GetModel(),
		Expiration:   item.GetExpiration(),
		Position:     dropObj.Object.Position,
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
	mapInstance := l.gs.GetMap(mapID)
	if mapInstance == nil {
		return
	}

	// Create spawn meso packet
	spawnPacket := &response.SpawnMeso{
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
	mapInstance := l.gs.GetMap(mapID)
	if mapInstance == nil {
		return
	}

	// Create remove item packet
	removePacket := &response.RemoveItem{
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
	mapInstance := l.gs.GetMap(mapID)
	if mapInstance == nil {
		return
	}

	// Convert entity to DTO
	mobDTO := mob.ToDTO()

	// Create spawn mob packet
	spawnPacket := &response.SpawnMob{
		Mob:       mobDTO,
		SpawnType: constant.MOB_SPAWN_TYPE_ANIMATE,
	}

	// Broadcast to all players on the map
	mapInstance.BroadcastToAllPlayers(spawnPacket, types.SEND_POLICY_ENCRYPT)
}

// OnMobRemoved sends remove mob packet to all players on the map
func (l *MapListenerImpl) OnMobRemoved(mapID uint32, mobID uint32, animationType constant.MobDieAnimationType) {
	// Get the map instance
	mapInstance := l.gs.GetMap(mapID)
	if mapInstance == nil {
		return
	}

	// Create remove mob packet
	removePacket := &response.DieMob{
		OID:           mobID,
		AnimationType: animationType,
	}

	// Broadcast to all players on the map
	mapInstance.BroadcastToAllPlayers(removePacket, types.SEND_POLICY_ENCRYPT)
}

func (l *MapListenerImpl) OnMobControllerChange(mob *entity.Mob, before *entity.Character, after *entity.Character) {
	if after != nil {
		// Convert entity to DTO
		mobDTO := mob.ToDTO()
		after.Send(&response.StartControlMob{
			Mob:   mobDTO,
			Aggro: false,
		}, types.SEND_POLICY_ENCRYPT)
	}
	// Note: StopControlMob is commented out in old server, so we skip it here too
}

// OnMobMoved broadcasts mob movement to all players on the map
func (l *MapListenerImpl) OnMobMoved(mapID uint32, mobID uint32, isAggroed bool, centerSplit int8, skill1 uint8, skill2 uint8, skill3 uint8, skill4 uint8, startPoint types.Vector2[int16], movements []dto.MoveFragment) {
	// Get the map instance
	mapInstance := l.gs.GetMap(mapID)
	if mapInstance == nil {
		return
	}

	// Create move mob packet
	movePacket := &response.MoveMob{
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
func (l *MapListenerImpl) OnAttack(mapID uint32, characterID uint32, attackInfo dto.AttackInfo, skillLevel uint8) {
	// Get the map instance
	mapInstance := l.gs.GetMap(mapID)
	if mapInstance == nil {
		return
	}

	// Create attack packet
	attackPacket := &response.Attack{
		CharacterId: characterID,
		AttackInfo:  attackInfo,
		SkillLevel:  skillLevel,
	}

	source := mapInstance.GetPlayer(characterID)
	if source == nil {
		return
	}

	mapInstance.Broadcast(source, attackPacket, types.SEND_POLICY_ENCRYPT, characterID)
}
