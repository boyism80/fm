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
func (l *MapListenerImpl) OnPlayerAdded(mapInstance *entity.Map, character *entity.Character, init bool) {
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
	playerID := character.GetID()
	if mapInstance.GetPlayerCount() > 1 { // More than just the new player
		for cid, p := range mapInstance.GetAllPlayers() {
			if cid == playerID {
				continue // Skip the new player
			}

			if ch, ok := p.(*entity.Character); ok {
				if ch.IsHidden() && !character.HasRoleAtLeast(ch.Role) {
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
	mapInstance.Broadcast(spawnPacket, &entity.BroadcastOption{
		ExceptPlayerIDs:    []uint32{playerID},
		ReferenceCharacter: character,
		RecipientFilter:    entity.BroadcastVisibleByReference,
	})

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
func (l *MapListenerImpl) OnPlayerRemoved(mapInstance *entity.Map, character *entity.Character) {
	if mapInstance == nil {
		return
	}

	// Send LeavePlayer packet to all other players on the map
	leavePacket := &response.LeavePlayer{
		ID: character.GetID(),
	}
	mapInstance.Broadcast(leavePacket, nil)
}

// OnPlayerMoved sends move player packet to other players on the map
func (l *MapListenerImpl) OnPlayerMoved(mapInstance *entity.Map, character *entity.Character) {
	if mapInstance == nil {
		return
	}

	// TODO: Create and send move player packet
}

// OnPlayerMove sends move packet with fragments to all other players on the map
func (l *MapListenerImpl) OnPlayerMove(mapInstance *entity.Map, character *entity.Character, startPoint types.Vector2[int16], fragments []dto.MoveFragment) {
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

	mapInstance.Broadcast(movePacket, &entity.BroadcastOption{
		ExceptPlayerIDs:    []uint32{character.GetID()},
		ReferenceCharacter: character,
		RecipientFilter:    entity.BroadcastVisibleByReference,
	})
}

// OnPlayerChat sends chat packet to other players on the map
func (l *MapListenerImpl) OnPlayerChat(mapInstance *entity.Map, character *entity.Character, message string) {
	if mapInstance == nil {
		return
	}

	// TODO: Create and send chat packet
}

// OnItemSpawned sends spawn item packet to all players on the map
func (l *MapListenerImpl) OnItemSpawned(mapInstance *entity.Map, item entity.Item, drop *entity.Drop) {
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
	mapInstance.Broadcast(spawnPacket, nil)
}

// OnMesoSpawned sends spawn meso packet to all players on the map
func (l *MapListenerImpl) OnMesoSpawned(mapInstance *entity.Map, meso *entity.Meso) {
	if mapInstance == nil {
		return
	}

	drop := meso.GetDrop()
	// Create spawn meso packet
	spawnPacket := &response.SpawnMeso{
		ID:           drop.OID,
		Animation:    constant.DROP_ITEM_ANIMATION_TYPE_LOOTING,
		DropType:     drop.DropType,
		Count:        meso.Count,
		OwnerID:      drop.Owner,
		Position:     drop.Position,
		SpawnedPoint: drop.SpawnedPoint,
		IsPlayerDrop: true,
	}

	// Broadcast to all players on the map
	mapInstance.Broadcast(spawnPacket, nil)
}

// OnItemRemoved sends remove item packet to all players on the map
func (l *MapListenerImpl) OnItemRemoved(mapInstance *entity.Map, itemID uint32, looterID uint32, mode constant.RemoveItemType) {
	if mapInstance == nil {
		return
	}

	// Create remove item packet
	removePacket := &response.RemoveItem{
		Mode:        mode,
		OID:         itemID,
		CharacterId: looterID,
	}

	// Broadcast to all players on the map
	mapInstance.Broadcast(removePacket, nil)
}

// OnMobSpawned sends spawn mob packet to all players on the map
func (l *MapListenerImpl) OnMobSpawned(mapInstance *entity.Map, mob *entity.Mob) {
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
	mapInstance.Broadcast(spawnPacket, nil)
}

// OnMobRemoved sends remove mob packet to all players on the map
func (l *MapListenerImpl) OnMobRemoved(mapInstance *entity.Map, mob *entity.Mob, animationType constant.MobDieAnimationType) {
	if mapInstance == nil {
		return
	}

	// Create remove mob packet
	removePacket := &response.DieMob{
		OID:           mob.OID,
		AnimationType: animationType,
	}

	// Broadcast to all players on the map
	mapInstance.Broadcast(removePacket, nil)
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
func (l *MapListenerImpl) OnMobMoved(mapInstance *entity.Map, mob *entity.Mob, isAggroed bool, centerSplit int8, skill1 uint8, skill2 uint8, skill3 uint8, skill4 uint8, startPoint types.Vector2[int16], movements []dto.MoveFragment) {
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
		OID:         mob.OID,
		StartPoint:  startPoint,
		Movements:   movements,
	}

	// Broadcast to all players on the map
	mapInstance.Broadcast(movePacket, nil)
}

// OnAttack broadcasts attack to all players on the map
func (l *MapListenerImpl) OnAttack(mapInstance *entity.Map, character *entity.Character, attackInfo dto.AttackInfo, skillLevel uint8) {
	if mapInstance == nil {
		return
	}

	// Create attack packet
	attackPacket := &response.Attack{
		CharacterId: character.GetID(),
		AttackInfo:  attackInfo,
		SkillLevel:  skillLevel,
	}

	mapInstance.Broadcast(attackPacket, &entity.BroadcastOption{
		ExceptPlayerIDs:    []uint32{character.GetID()},
		ReferenceCharacter: character,
		RecipientFilter:    entity.BroadcastVisibleByReference,
	})
}

// OnMobDebuffApplied broadcasts APPLY_DEBUFF (0xAF) to all players on the map (including attacker).
func (l *MapListenerImpl) OnMobDebuffApplied(mapInstance *entity.Map, mob *entity.Mob, debuff constant.Debuff, value int32, skillID uint32, durationMs int64) {
	if mapInstance == nil {
		return
	}
	buffTime := int16(32767)
	if durationMs > 0 && durationMs/1000 < 32767 {
		buffTime = int16(durationMs / 1000)
	}
	pkt := &response.ApplyDebuff{
		OID:        mob.OID,
		Status:     int32(debuff.Mask),
		X:          int16(value),
		SkillID:    skillID,
		BuffTime:   buffTime,
		Delay:      0,
		StatusSize: 1,
	}
	mapInstance.Broadcast(pkt, nil)
}

// OnMobDebuffCancelled broadcasts CANCEL_DEBUFF (0xB0) to all players on the map.
func (l *MapListenerImpl) OnMobDebuffCancelled(mapInstance *entity.Map, mob *entity.Mob, debuff constant.Debuff) {
	if mapInstance == nil {
		return
	}
	pkt := &response.CancelDebuff{
		OID:    mob.OID,
		Status: int32(debuff.Mask),
		Size:   1,
	}
	mapInstance.Broadcast(pkt, nil)
}
