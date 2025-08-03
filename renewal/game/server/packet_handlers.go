package server

import (
	"fmt"
	"log"
	"strings"

	common_req "github.com/boyism80/fm/common/protocol/req"
	"github.com/boyism80/fm/common/stream"
	"github.com/boyism80/fm/common/types"
	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/game/constant"
	"github.com/boyism80/fm/game/data"
	"github.com/boyism80/fm/game/entity"
	"github.com/boyism80/fm/game/protocol"
	"github.com/boyism80/fm/game/protocol/req"
	"github.com/boyism80/fm/game/protocol/resp"
	"github.com/boyism80/fm/renewal/game/client"
)

// registerPacketHandlers registers all game server packet handlers
func (gs *GameServer) registerPacketHandlers() {
	// Register the 15 game server packet handlers
	gs.server.RegisterPacketHandler(0x0A, gs.handlePong)
	gs.server.RegisterPacketHandler(0x06, gs.handleLoginGame)
	gs.server.RegisterPacketHandler(0x18, gs.handleMovePlayer)
	gs.server.RegisterPacketHandler(0x20, gs.handleNormalChat)
	gs.server.RegisterPacketHandler(0x1B, gs.handleAttack)
	gs.server.RegisterPacketHandler(0x36, gs.handleMoveItem)
	gs.server.RegisterPacketHandler(0x34, gs.handleSortInventory)
	gs.server.RegisterPacketHandler(0xA3, gs.handleItemLoot)
	gs.server.RegisterPacketHandler(0x4D, gs.handleDropMeso)
	gs.server.RegisterPacketHandler(0x15, gs.handleWarp)
	gs.server.RegisterPacketHandler(0x9E, gs.handleNpcControl)
	gs.server.RegisterPacketHandler(0x2B, gs.handleDialog)
	gs.server.RegisterPacketHandler(0x29, gs.handleNpcClick)
	gs.server.RegisterPacketHandler(0x95, gs.handleMoveMob)
	gs.server.RegisterPacketHandler(0x1F, gs.handleDamaged)

	log.Printf("Registered %d game server packet handlers", gs.server.GetPacketHandler().GetHandlerCount())
}

// handlePong processes pong responses
func (gs *GameServer) handlePong(ctx *core.ClientContext, data []byte) error {
	reader := stream.NewStreamReader(&data, stream.LittleEndian)
	request := &common_req.Pong{}
	if err := request.Deserialize(reader); err != nil {
		log.Printf("Failed to deserialize pong packet from %s: %v", ctx.Client.GetConnection().RemoteAddr(), err)
		return err
	}

	log.Printf("Pong packet received from %s",
		ctx.Client.GetConnection().RemoteAddr())
	return nil
}

// handleLoginGame processes game login requests
func (gs *GameServer) handleLoginGame(ctx *core.ClientContext, data []byte) error {
	reader := stream.NewStreamReader(&data, stream.LittleEndian)
	request := &req.LoginGame{}
	if err := request.Deserialize(reader); err != nil {
		log.Printf("Failed to deserialize login game packet from %s: %v", ctx.Client.GetConnection().RemoteAddr(), err)
		return err
	}

	log.Printf("Login game packet received from %s - Player ID: %d",
		ctx.Client.GetConnection().RemoteAddr(), request.PlayerId)

	// Create character name based on player ID (following old server pattern)
	name := "채승현"
	if request.PlayerId != 1 {
		name = "채진영"
	}

	// Create character using NewDummyCharacter (following old server pattern)
	character := entity.NewDummyCharacter(ctx.Client, nil, request.PlayerId, name, gs)

	// Set character listener for packet sending
	character.Listener = NewGameCharacterListener(gs, &character)

	// Set character in game client
	client, ok := ctx.Client.(*client.GameClient)
	if !ok {
		log.Printf("Client is not a GameClient")
		return fmt.Errorf("client is not a GameClient")
	}
	client.SetCharacter(&character)

	// Get map instance and set character position based on SpawnPoint
	mapInstance := gs.GetMap(character.Map)
	if mapInstance == nil {
		log.Printf("Map %d not found (should have been pre-created)", character.Map)
		return fmt.Errorf("map %d not found", character.Map)
	}

	// Set character position based on SpawnPoint (following old server pattern)
	mapSpec := mapInstance.GetSpec()
	if mapSpec == nil {
		log.Printf("MapSpec not found for map %d", character.Map)
		return fmt.Errorf("mapSpec not found for map %d", character.Map)
	}

	portal, ok := mapSpec.Portals[character.SpawnPoint]
	if !ok {
		log.Printf("Portal %d not found in map %d", character.SpawnPoint, character.Map)
		return fmt.Errorf("portal %d not found in map %d", character.SpawnPoint, character.Map)
	}

	// Set character position to portal position (following old server pattern)
	character.Position = portal.Position
	character.Stance = 0

	if err := mapInstance.AddPlayer(character.ID, &character, true); err != nil {
		log.Printf("Failed to add player to map: %v", err)
		return err
	}

	log.Printf("Character created for player %d: %s (Level %d, Class %d, Map %d)",
		character.ID, character.Name, character.Level, character.Class, character.Map)

	return nil
}

// handleMovePlayer processes player movement requests
func (gs *GameServer) handleMovePlayer(ctx *core.ClientContext, data []byte) error {
	reader := stream.NewStreamReader(&data, stream.LittleEndian)
	request := &req.MovePlayer{}
	if err := request.Deserialize(reader); err != nil {
		log.Printf("Failed to deserialize move player packet from %s: %v", ctx.Client.GetConnection().RemoteAddr(), err)
		return err
	}

	// Get the game client and character
	client, ok := ctx.Client.(*client.GameClient)
	if !ok {
		log.Printf("Client is not a GameClient")
		return fmt.Errorf("client is not a GameClient")
	}

	character := client.GetCharacter()
	if character == nil {
		log.Printf("Character is nil for client")
		return fmt.Errorf("character is nil")
	}

	// Store the position before movement
	beforePosition := character.Position

	// Process movement fragments
	for _, frag := range request.Fragments {
		if move, ok := frag.(*protocol.AbsoluteLifeMovement); ok {
			character.Position = move.Position
		}
		character.Stance = frag.GetStance()
	}

	// Get the map instance
	mapInstance := gs.GetMap(character.GetMap())
	if mapInstance == nil {
		log.Printf("Map %d not found for character movement", character.GetMap())
		return fmt.Errorf("map %d not found", character.GetMap())
	}

	// Broadcast movement to other players on the map
	movePacket := &resp.Move{
		Character:  character,
		Fragments:  request.Fragments,
		StartPoint: beforePosition,
	}
	mapInstance.BroadcastToPlayers(movePacket, types.SEND_POLICY_ENCRYPT, character.GetID())

	log.Printf("Move player packet processed for character %d - Position: %v, Fragments: %d",
		character.GetID(), character.Position, len(request.Fragments))
	return nil
}

// handleNormalChat processes normal chat messages
func (gs *GameServer) handleNormalChat(ctx *core.ClientContext, data []byte) error {
	reader := stream.NewStreamReader(&data, stream.LittleEndian)
	request := &req.NormalChat{}
	if err := request.Deserialize(reader); err != nil {
		log.Printf("Failed to deserialize normal chat packet from %s: %v", ctx.Client.GetConnection().RemoteAddr(), err)
		return err
	}

	// Get the game client and character
	client, ok := ctx.Client.(*client.GameClient)
	if !ok {
		log.Printf("Client is not a GameClient")
		return fmt.Errorf("client is not a GameClient")
	}

	character := client.GetCharacter()
	if character == nil {
		log.Printf("Character is nil for client")
		return fmt.Errorf("character is nil")
	}

	// Check for command prefix (following old server pattern)
	if strings.HasPrefix(request.Message, "/") {
		params := strings.Split(strings.TrimPrefix(request.Message, "/"), " ")
		err := gs.commandHandler.Handle(client, params...)
		if err != nil {
			log.Printf("Command error: %v", err)
		}
		return nil
	}

	// Get the map instance
	mapInstance := gs.GetMap(character.GetMap())
	if mapInstance == nil {
		log.Printf("Map %d not found for character chat", character.GetMap())
		return fmt.Errorf("map %d not found", character.GetMap())
	}

	// Use character listener to handle chat (following old server pattern)
	character.Listener.OnChat(request.Message, false, request.DontRecordHistory)

	log.Printf("Normal chat packet processed for character %d - Message: %s, DontRecord: %t",
		character.GetID(), request.Message, request.DontRecordHistory)
	return nil
}

// handleAttack processes attack requests
func (gs *GameServer) handleAttack(ctx *core.ClientContext, data []byte) error {
	reader := stream.NewStreamReader(&data, stream.LittleEndian)
	request := &req.Attack{}
	if err := request.Deserialize(reader); err != nil {
		log.Printf("Failed to deserialize attack packet from %s: %v", ctx.Client.GetConnection().RemoteAddr(), err)
		return err
	}

	// Get the game client and character
	client, ok := ctx.Client.(*client.GameClient)
	if !ok {
		log.Printf("Client is not a GameClient")
		return fmt.Errorf("client is not a GameClient")
	}

	character := client.GetCharacter()
	if character == nil {
		log.Printf("Character is nil for client")
		return fmt.Errorf("character is nil")
	}

	// Get map instance
	mapID := character.GetMap()
	mapInstance := gs.GetMap(mapID)
	if mapInstance == nil {
		log.Printf("Character is not in a map")
		return fmt.Errorf("character is not in a map")
	}

	// Process damage to mobs (following old server pattern)
	for _, damage := range request.AttackInfo.Damages {
		mob := mapInstance.GetMob(damage.OID)
		if mob == nil {
			log.Printf("Mob not found for OID: %d", damage.OID)
			continue
		}

		// Apply damage to mob using Mob.Damage method
		for _, damagePair := range damage.DamagePairs {
			mob.Damage(uint16(damagePair.Damage), character)

			// Send mob HP update to attacker
			character.Send(&resp.ShowMobHp{
				OID:        mob.OID,
				Percentage: uint8(mob.Hp * 100 / mob.MaxHp),
			}, types.SEND_POLICY_ENCRYPT)
		}
	}

	// Broadcast attack to all players in map (following old server pattern)
	mapInstance.BroadcastToPlayers(&resp.Attack{
		CharacterId: character.GetID(),
		AttackInfo:  request.AttackInfo,
		SkillLevel:  0,
	}, types.SEND_POLICY_ENCRYPT, character.GetID())

	log.Printf("Attack packet processed for character %d - Targets: %d, Hits: %d, Skill: %d",
		character.GetID(), request.AttackInfo.Targets, request.AttackInfo.Hits, request.AttackInfo.Skill)
	return nil
}

// handleMoveItem processes item movement requests
func (gs *GameServer) handleMoveItem(ctx *core.ClientContext, data []byte) error {
	reader := stream.NewStreamReader(&data, stream.LittleEndian)
	request := &req.MoveItem{}
	if err := request.Deserialize(reader); err != nil {
		log.Printf("Failed to deserialize move item packet from %s: %v", ctx.Client.GetConnection().RemoteAddr(), err)
		return err
	}

	// Get the game client and character
	client, ok := ctx.Client.(*client.GameClient)
	if !ok {
		log.Printf("Client is not a GameClient")
		return fmt.Errorf("client is not a GameClient")
	}

	character := client.GetCharacter()
	if character == nil {
		log.Printf("Character is nil for client")
		return fmt.Errorf("character is nil")
	}

	// Following old server pattern
	if request.Source < 0 {
		// Unequip item
		gs.handleUnequip(client, character, constant.EquipmentPartsType(request.Source), request.Dest)
	} else if request.Dest < 0 {
		// Equip item
		gs.handleEquip(client, character, constant.EquipmentPartsType(request.Dest), request.Source)
	} else if request.Dest == 0 {
		// Drop item
		gs.handleDrop(client, character, request.InventoryType, request.Source, request.Count)
	} else {
		// Move item within inventory
		gs.handleMoveItemInternal(client, character, request.InventoryType, request.Source, request.Dest)
	}

	log.Printf("Move item packet processed for character %d - Source: %d, Dest: %d, Count: %d",
		character.GetID(), request.Source, request.Dest, request.Count)
	return nil
}

// handleSortInventory processes inventory sorting requests
func (gs *GameServer) handleSortInventory(ctx *core.ClientContext, data []byte) error {
	reader := stream.NewStreamReader(&data, stream.LittleEndian)
	request := &req.SortInventory{}
	if err := request.Deserialize(reader); err != nil {
		log.Printf("Failed to deserialize sort inventory packet from %s: %v", ctx.Client.GetConnection().RemoteAddr(), err)
		return err
	}

	// Get the game client and character
	client, ok := ctx.Client.(*client.GameClient)
	if !ok {
		log.Printf("Client is not a GameClient")
		return fmt.Errorf("client is not a GameClient")
	}

	character := client.GetCharacter()
	if character == nil {
		log.Printf("Character is nil for client")
		return fmt.Errorf("character is nil")
	}

	// Following old server pattern
	gs.handleMergeItems(client, character, request.InventoryType)
	gs.handleSortInventoryInternal(client, character, request.InventoryType)

	// Send end sort inventory response
	client.Send(&resp.EndSortInventory{
		InventoryType: request.InventoryType,
	}, types.SEND_POLICY_ENCRYPT)

	// Send unlock action response
	client.Send(&resp.UpdateStats{
		UnlockAction: true,
	}, types.SEND_POLICY_ENCRYPT)

	log.Printf("Sort inventory packet processed for character %d - Inventory Type: %d",
		character.GetID(), request.InventoryType)
	return nil
}

// handleItemLoot processes item looting requests
func (gs *GameServer) handleItemLoot(ctx *core.ClientContext, data []byte) error {
	reader := stream.NewStreamReader(&data, stream.LittleEndian)
	request := &req.ItemLoot{}
	if err := request.Deserialize(reader); err != nil {
		log.Printf("Failed to deserialize item loot packet from %s: %v", ctx.Client.GetConnection().RemoteAddr(), err)
		return err
	}

	// Get the game client and character
	client, ok := ctx.Client.(*client.GameClient)
	if !ok {
		log.Printf("Client is not a GameClient")
		return fmt.Errorf("client is not a GameClient")
	}

	character := client.GetCharacter()
	if character == nil {
		log.Printf("Character is nil for client")
		return fmt.Errorf("character is nil")
	}

	// Get map instance
	mapInstance := gs.GetMap(character.GetMap())
	if mapInstance == nil {
		log.Printf("Map %d not found for character %d", character.GetMap(), character.GetID())
		client.Send(&resp.UpdateStats{
			UnlockAction: true,
		}, types.SEND_POLICY_ENCRYPT)
		return fmt.Errorf("map not found")
	}

	// Attempt to loot the item/meso
	lootedObject, reason := mapInstance.LootItem(request.OID, character, request.Position)
	if reason != constant.LOOT_SUCCESS {
		log.Printf("Failed to loot item %d for character %d, reason: %d", request.OID, character.GetID(), reason)

		// Only send ItemGainFailed for inventory/meso capacity issues
		if reason == constant.LOOT_FAILED_INVENTORY_FULL || reason == constant.LOOT_FAILED_MESO_FULL {
			client.Send(&resp.ItemGainFailed{
				Mode: constant.ITEM_GAIN_FAILED_TYPE_FULL,
			}, types.SEND_POLICY_ENCRYPT)
		}

		client.Send(&resp.UpdateStats{
			UnlockAction: true,
		}, types.SEND_POLICY_ENCRYPT)
		return nil
	}

	// Handle different types of looted objects
	switch obj := lootedObject.(type) {
	case entity.Item:
		// Handle item looting
		item := obj
		invenType := item.GetInventoryType()
		inven := character.Inventory[invenType]
		spec := item.GetSpec()
		gain := uint16(0)

		// Add item to inventory (capacity check already done in LootItem)
		for item.GetCount() > 0 {
			slot, ok := inven.FindSlot(spec)
			if !ok {
				break
			}

			exists, ok := inven.Items[int16(slot)]
			cap := uint16(0)
			if ok {
				// Increase existing item count
				cap = min(spec.GetCapacity()-exists.GetCount(), item.GetCount())
				exists.Increase(cap)
				client.Send(&resp.UpdateInventorySlot{
					InventoryType: invenType,
					Slot:          int16(slot),
					Item:          exists,
				}, types.SEND_POLICY_ENCRYPT)
			} else {
				// Add new item to slot
				cap = min(spec.GetCapacity(), item.GetCount())
				inven.Items[int16(slot)] = item.Clone(cap)
				client.Send(&resp.AddInventorySlot{
					InventoryType: invenType,
					Slot:          int16(slot),
					Item:          inven.Items[int16(slot)],
				}, types.SEND_POLICY_ENCRYPT)
			}
			if item.Reduce(cap) == 0 {
				break
			}
			gain += cap
		}

		// Show item gain message
		client.Send(&resp.ShowItemGain{
			ItemId: spec.GetID(),
			Count:  uint32(gain),
			Mode:   resp.ShowItemGainTypeStatus,
		}, types.SEND_POLICY_ENCRYPT)

	case *entity.Meso:
		// Handle meso looting
		meso := obj
		mesoCount := meso.GetCount32()

		// Add meso to character (capacity check already done in LootItem)
		character.Meso += int32(mesoCount)

		// Show meso gain message
		client.Send(&resp.ShowMesoGain{
			Count: int32(mesoCount),
			Mode:  resp.ShowMesoGainTypeStatus,
		}, types.SEND_POLICY_ENCRYPT)

		client.Send(&resp.UpdateStats{
			Stats: map[constant.Stat]int32{
				constant.STAT_MESO: character.Meso,
			},
		}, types.SEND_POLICY_ENCRYPT)

		log.Printf("Meso loot packet processed for character %d - OID: %d, Count: %d",
			character.GetID(), request.OID, mesoCount)

	default:
		log.Printf("Unknown looted object type for OID %d", request.OID)
		client.Send(&resp.UpdateStats{
			UnlockAction: true,
		}, types.SEND_POLICY_ENCRYPT)
		return nil
	}

	// Send unlock action response
	client.Send(&resp.UpdateStats{
		UnlockAction: true,
	}, types.SEND_POLICY_ENCRYPT)

	log.Printf("Item loot packet processed for character %d - OID: %d",
		character.GetID(), request.OID)
	return nil
}

// handleDropMeso processes meso dropping requests
func (gs *GameServer) handleDropMeso(ctx *core.ClientContext, data []byte) error {
	reader := stream.NewStreamReader(&data, stream.LittleEndian)
	request := &req.DropMeso{}
	if err := request.Deserialize(reader); err != nil {
		log.Printf("Failed to deserialize drop meso packet from %s: %v", ctx.Client.GetConnection().RemoteAddr(), err)
		return err
	}

	// Get the game client and character
	client, ok := ctx.Client.(*client.GameClient)
	if !ok {
		log.Printf("Client is not a GameClient")
		return fmt.Errorf("client is not a GameClient")
	}

	character := client.GetCharacter()
	if character == nil {
		log.Printf("Character is nil for client")
		return fmt.Errorf("character is nil")
	}

	// Validate meso count (following old server pattern)
	if request.Count < 10 || request.Count > 50000 {
		log.Printf("Invalid meso count: %d", request.Count)
		client.Send(&resp.UpdateStats{
			UnlockAction: true,
		}, types.SEND_POLICY_ENCRYPT)
		return nil
	}

	// Check if character has enough meso
	if request.Count > character.Meso {
		log.Printf("Character doesn't have enough meso: %d < %d", character.Meso, request.Count)
		client.Send(&resp.UpdateStats{
			UnlockAction: true,
		}, types.SEND_POLICY_ENCRYPT)
		return nil
	}

	// Deduct meso from character
	character.Meso -= request.Count
	client.Send(&resp.UpdateStats{
		Stats: map[constant.Stat]int32{
			constant.STAT_MESO: character.Meso,
		},
		UnlockAction: true,
	}, types.SEND_POLICY_ENCRYPT)

	// Spawn meso on map (following old server pattern)
	mapInstance := gs.GetMap(character.GetMap())
	if mapInstance != nil {
		if err := mapInstance.SpawnMeso(request.Count, character.Position, character.ID, constant.DROP_TYPE_FFA); err != nil {
			log.Printf("Failed to spawn meso on map: %v", err)
		} else {
			log.Printf("Drop meso spawned: %d at position %v", request.Count, character.Position)
		}
	}

	log.Printf("Drop meso packet processed for character %d - Count: %d",
		character.GetID(), request.Count)
	return nil
}

// handleWarp processes warp requests
func (gs *GameServer) handleWarp(ctx *core.ClientContext, data []byte) error {
	reader := stream.NewStreamReader(&data, stream.LittleEndian)
	request := &req.Warp{}
	if err := request.Deserialize(reader); err != nil {
		log.Printf("Failed to deserialize warp packet from %s: %v", ctx.Client.GetConnection().RemoteAddr(), err)
		return err
	}

	// Get the game client
	client, ok := ctx.Client.(*client.GameClient)
	if !ok {
		log.Printf("Client is not a GameClient")
		return fmt.Errorf("client is not a GameClient")
	}

	character := client.GetCharacter()
	if character == nil {
		return fmt.Errorf("character not found")
	}

	var targetMapId uint32
	var spawnPoint uint8
	stats := map[constant.Stat]int32{}

	if request.Target != 0xFFFFFFFF {
		// Return map warp (when HP is 0)
		if character.Hp == 0 {
			character.Hp = 50
			character.Stance = 0

			// Get return map from current map spec
			currentMap := gs.GetMap(character.Map)
			if currentMap == nil {
				return fmt.Errorf("current map not found")
			}

			mapSpec := currentMap.GetSpec()
			if mapSpec == nil {
				return fmt.Errorf("map spec not found")
			}

			targetMapId = uint32(mapSpec.ReturnMapId)
			spawnPoint = 0
			stats[constant.STAT_HP] = int32(character.Hp)

			// Send stats update
			client.Send(&resp.UpdateStats{
				Stats:        stats,
				UnlockAction: true,
			}, types.SEND_POLICY_ENCRYPT)
		} else {
			// Invalid target for non-zero HP
			client.Send(&resp.UpdateStats{
				UnlockAction: true,
			}, types.SEND_POLICY_ENCRYPT)
			return nil
		}
	} else {
		// Normal portal warp
		currentMap := gs.GetMap(character.Map)
		if currentMap == nil {
			return fmt.Errorf("current map not found")
		}

		mapSpec := currentMap.GetSpec()
		if mapSpec == nil {
			return fmt.Errorf("map spec not found")
		}

		// Find portal by name
		portal, ok := mapSpec.FindPortal(request.PortalName)
		if !ok {
			client.Send(&resp.UpdateStats{
				UnlockAction: true,
			}, types.SEND_POLICY_ENCRYPT)
			return nil
		}

		// Get target map spec
		targetMapSpec, ok := gs.resources.Maps[uint32(portal.TargetMapId)]
		if !ok {
			client.Send(&resp.UpdateStats{
				UnlockAction: true,
			}, types.SEND_POLICY_ENCRYPT)
			return nil
		}

		// Find target portal in new map
		targetPortal, ok := targetMapSpec.FindPortal(portal.Target)
		if !ok {
			client.Send(&resp.UpdateStats{
				UnlockAction: true,
			}, types.SEND_POLICY_ENCRYPT)
			return nil
		}

		targetMapId = uint32(portal.TargetMapId)
		spawnPoint = targetPortal.ID
	}

	// Perform the warp
	if err := gs.performWarp(client, character, targetMapId, spawnPoint); err != nil {
		return fmt.Errorf("failed to perform warp: %v", err)
	}

	log.Printf("Warp packet processed from %s - Target Map: %d, Spawn Point: %d",
		ctx.Client.GetConnection().RemoteAddr(), targetMapId, spawnPoint)
	return nil
}

// handleNpcControl processes NPC control requests
func (gs *GameServer) handleNpcControl(ctx *core.ClientContext, data []byte) error {
	reader := stream.NewStreamReader(&data, stream.LittleEndian)
	request := &req.NpcAction{}
	if err := request.Deserialize(reader); err != nil {
		log.Printf("Failed to deserialize NPC control packet from %s: %v", ctx.Client.GetConnection().RemoteAddr(), err)
		return err
	}

	// Get the game client
	client, ok := ctx.Client.(*client.GameClient)
	if !ok {
		log.Printf("Client is not a GameClient")
		return fmt.Errorf("client is not a GameClient")
	}

	// Send NPC action response back to the client (following old server pattern)
	response := &resp.NpcAction{
		Bytes: request.Bytes,
	}
	client.Send(response, types.SEND_POLICY_ENCRYPT)

	log.Printf("NPC control packet processed from %s - Bytes: %d bytes",
		ctx.Client.GetConnection().RemoteAddr(), len(request.Bytes))
	return nil
}

// handleDialog processes dialog requests
func (gs *GameServer) handleDialog(ctx *core.ClientContext, data []byte) error {
	reader := stream.NewStreamReader(&data, stream.LittleEndian)
	request := &req.Dialog{}
	if err := request.Deserialize(reader); err != nil {
		log.Printf("Failed to deserialize dialog packet from %s: %v", ctx.Client.GetConnection().RemoteAddr(), err)
		return err
	}

	log.Printf("Dialog packet received from %s - Dialog Type: %d, Next: %t",
		ctx.Client.GetConnection().RemoteAddr(), request.DialogType, request.Next)
	return nil
}

// handleNpcClick processes NPC click requests
func (gs *GameServer) handleNpcClick(ctx *core.ClientContext, data []byte) error {
	reader := stream.NewStreamReader(&data, stream.LittleEndian)
	request := &req.NpcClick{}
	if err := request.Deserialize(reader); err != nil {
		log.Printf("Failed to deserialize NPC click packet from %s: %v", ctx.Client.GetConnection().RemoteAddr(), err)
		return err
	}

	log.Printf("NPC click packet received from %s - OID: %d",
		ctx.Client.GetConnection().RemoteAddr(), request.OID)
	return nil
}

// handleMoveMob processes mob movement requests
func (gs *GameServer) handleMoveMob(ctx *core.ClientContext, data []byte) error {
	reader := stream.NewStreamReader(&data, stream.LittleEndian)
	request := &req.MoveMob{}
	if err := request.Deserialize(reader); err != nil {
		log.Printf("Failed to deserialize move mob packet from %s: %v", ctx.Client.GetConnection().RemoteAddr(), err)
		return err
	}

	log.Printf("Move mob packet received from %s - OID: %d, MovementId: %d, IsAggroed: %t, Unknown2: %t",
		ctx.Client.GetConnection().RemoteAddr(), request.OID, request.MovementId, request.IsAggroed, request.Unknown2)

	// Get game client
	client, ok := ctx.Client.(*client.GameClient)
	if !ok {
		log.Printf("Client is not a GameClient")
		return fmt.Errorf("client is not a GameClient")
	}

	// Get character
	character := client.GetCharacter()
	if character == nil {
		log.Printf("Character not found for client")
		return fmt.Errorf("character not found")
	}

	// Get map instance
	mapInstance := gs.GetMap(character.GetMap())
	if mapInstance == nil {
		log.Printf("Map %d not found", character.GetMap())
		return fmt.Errorf("map %d not found", character.GetMap())
	}

	// Get mob from map
	mob := mapInstance.GetMob(request.OID)
	if mob == nil {
		log.Printf("Mob %d not found on map %d", request.OID, character.GetMap())
		return fmt.Errorf("mob %d not found", request.OID)
	}

	// Check controller (following old server pattern)
	controllerTable := mapInstance.GetControllerTable()
	controller, exists := controllerTable.GetController(mob)
	if !exists {
		log.Printf("No controller found for mob %d", request.OID)
		return fmt.Errorf("no controller found for this mob")
	}

	// If sender is not the controller, handle control switching (following old server pattern)
	if controller.GetID() != character.GetID() {
		if request.Unknown2 {
			// TODO: stopControl - Stop control logic
			log.Printf("Stop control requested for mob %d by character %d", request.OID, character.GetID())
			// Currently only logging, actual implementation to be added later
		} else {
			// TODO: switchControl - Switch control logic
			log.Printf("Switch control requested for mob %d by character %d", request.OID, character.GetID())
			// Currently only logging, actual implementation to be added later
		}
		return nil // Don't return error, handle normally
	}

	// Send ControlMoveMob packet to the controller (following old server pattern)
	controlPacket := &resp.ControlMoveMob{
		OID:          request.OID,
		MoveId:       request.MovementId,
		EnabledSkill: request.IsAggroed,
		MP:           mob.Mp,
		SkillId:      0,
		SkillLevel:   0,
	}
	client.Send(controlPacket, types.SEND_POLICY_ENCRYPT)

	// Update mob position based on movement data (following old server pattern)
	for _, movement := range request.Movements {
		if move, ok := movement.(*protocol.AbsoluteLifeMovement); ok {
			mob.Position = move.Position
		}

		mob.Stance = movement.GetStance()
	}

	// Broadcast the movement to all players on the map
	movePacket := &resp.MoveMob{
		IsAggroed:   request.IsAggroed,
		CenterSplit: request.CenterSplit,
		Skill1:      request.Skill1,
		Skill2:      request.Skill2,
		Skill3:      request.Skill3,
		Skill4:      request.Skill4,
		OID:         request.OID,
		StartPoint:  mob.Position, // Use current position as start point
		Movements:   request.Movements,
	}

	// Broadcast to all players on the map (excluding the sender)
	mapInstance.BroadcastToPlayers(movePacket, types.SEND_POLICY_ENCRYPT, character.GetID())

	log.Printf("Mob %d movement broadcasted - Position: %v, Stance: %d",
		request.OID, mob.Position, mob.Stance)

	return nil
}

// handleDamaged processes damage requests
func (gs *GameServer) handleDamaged(ctx *core.ClientContext, data []byte) error {
	reader := stream.NewStreamReader(&data, stream.LittleEndian)
	request := &req.Damaged{}
	if err := request.Deserialize(reader); err != nil {
		log.Printf("Failed to deserialize damaged packet from %s: %v", ctx.Client.GetConnection().RemoteAddr(), err)
		return err
	}

	log.Printf("Damaged packet received from %s - Damage: %d",
		ctx.Client.GetConnection().RemoteAddr(), request.Damage)

	// Get game client
	client, ok := ctx.Client.(*client.GameClient)
	if !ok {
		log.Printf("Client is not a GameClient")
		return fmt.Errorf("client is not a GameClient")
	}

	// Get character
	character := client.GetCharacter()
	if character == nil {
		log.Printf("Character not found for client")
		return fmt.Errorf("character not found")
	}

	// Process damage (following old server pattern)
	stats := map[constant.Stat]int32{}
	if !character.Invincible {
		// Calculate new HP (following old server pattern)
		newHp := int32(character.Hp) - int32(request.Damage)
		if newHp < 0 {
			newHp = 0
		}
		if newHp > int32(character.MaxHp) {
			newHp = int32(character.MaxHp)
		}

		character.Hp = uint16(newHp)
		stats[constant.STAT_HP] = int32(character.Hp)
	}

	// Send stats update to client (following old server pattern)
	client.Send(&resp.UpdateStats{
		Stats:        stats,
		UnlockAction: true,
	}, types.SEND_POLICY_ENCRYPT)

	log.Printf("Damage processed for character %d - Damage: %d, New HP: %d/%d",
		character.GetID(), request.Damage, character.Hp, character.MaxHp)

	return nil
}

// performWarp performs the actual warp operation
func (gs *GameServer) performWarp(client *client.GameClient, character *entity.Character, targetMapId uint32, spawnPoint uint8) error {
	// Get current map
	currentMap := gs.GetMap(character.Map)
	if currentMap == nil {
		return fmt.Errorf("current map not found")
	}

	// Get target map
	targetMap := gs.GetMap(targetMapId)
	if targetMap == nil {
		return fmt.Errorf("target map %d not found", targetMapId)
	}

	// Remove character from current map
	currentMap.RemovePlayer(character.GetID())

	// Update character's map and spawn point
	character.Map = targetMapId
	character.SpawnPoint = spawnPoint

	// Set character position based on spawn point
	mapSpec := targetMap.GetSpec()
	if mapSpec != nil && len(mapSpec.Portals) > 0 {
		// Find portal by spawn point
		for portalId, portal := range mapSpec.Portals {
			if portalId == spawnPoint {
				character.Position = portal.Position
				break
			}
		}
	}

	// Add character to target map in a new logic task
	// This ensures the task runs on the correct logic thread for the new map
	task := &core.LogicTask{
		Predicate: func() bool {
			// Check if character is still valid and connected
			return character != nil && client.GetConnection() != nil
		},
		Logic: func() error {
			// This will trigger GameMapListener.OnPlayerAdded which sends the Warp packet
			if err := targetMap.AddPlayer(character.GetID(), character, false); err != nil {
				return fmt.Errorf("failed to add character to target map: %v", err)
			}
			return nil
		},
		Callback: func(success bool, err error) {
			if err != nil {
				log.Printf("Failed to add character to target map: %v", err)
			} else {
				log.Printf("Character %d successfully warped to map %d", character.GetID(), targetMapId)
			}
		},
		Object:     client, // Use client for thread assignment (based on character's new map)
		MaxRetries: 3,
	}

	// Submit to appropriate logic thread based on character's new map
	if err := gs.server.SubmitLogicTaskForObject(client, task); err != nil {
		return fmt.Errorf("failed to submit warp task: %v", err)
	}

	return nil
}

// handleUnequip handles unequipping items
func (gs *GameServer) handleUnequip(client *client.GameClient, character *entity.Character, parts constant.EquipmentPartsType, slot int16) {
	// Check if equipment exists
	if character.Equipments[parts] == nil {
		return
	}

	// Check if inventory slot is empty
	inven := character.Inventory[constant.INVENTORY_TYPE_EQUIPMENT]
	if inven.Items[slot] != nil {
		return
	}

	// Move equipment to inventory
	inven.Items[slot] = character.Equipments[parts]
	delete(character.Equipments, parts)

	// Send swap inventory slot response
	client.Send(&resp.SwapInventorySlot{
		InventoryType:   constant.INVENTORY_TYPE_EQUIPMENT,
		Source:          int16(parts),
		Dest:            slot,
		EquipmentAction: resp.EQUIPMENT_ACTION_TYPE_OFF,
	}, types.SEND_POLICY_ENCRYPT)

	// Broadcast character look update
	mapInstance := gs.GetMap(character.GetMap())
	if mapInstance != nil {
		mapInstance.BroadcastToPlayers(&resp.UpdateCharacterLook{
			Character: character,
		}, types.SEND_POLICY_ENCRYPT, character.GetID())
	}
}

// handleEquip handles equipping items
func (gs *GameServer) handleEquip(client *client.GameClient, character *entity.Character, parts constant.EquipmentPartsType, slot int16) {
	inven := character.Inventory[constant.INVENTORY_TYPE_EQUIPMENT]
	if inven.Items[slot] == nil {
		return
	}

	// Check if item is equipment
	new, ok := inven.Items[slot].(*entity.Equipment)
	if !ok {
		return
	}

	old, swap := character.Equipments[parts]

	// Handle overall equipment logic
	switch parts {
	case constant.EQUIPMENT_PARTS_TOP:
		if new.IsOverall() {
			_, isWearPants := character.Equipments[constant.EQUIPMENT_PARTS_PANTS]
			if isWearPants {
				// Unequip pants when wearing overall
				storageSlot, isFree := inven.NextSlot()
				if !isFree {
					client.Send(&resp.ItemGainFailed{
						Mode: constant.ITEM_GAIN_FAILED_TYPE_FULL,
					}, types.SEND_POLICY_ENCRYPT)
					return
				}
				gs.handleUnequip(client, character, constant.EQUIPMENT_PARTS_PANTS, int16(storageSlot))
			}
		}

	case constant.EQUIPMENT_PARTS_PANTS:
		top, isWearTop := character.Equipments[constant.EQUIPMENT_PARTS_TOP]
		if isWearTop && top.IsOverall() {
			storageSlot, isFree := inven.NextSlot()
			if swap && !isFree {
				client.Send(&resp.ItemGainFailed{
					Mode: constant.ITEM_GAIN_FAILED_TYPE_FULL,
				}, types.SEND_POLICY_ENCRYPT)
				return
			}
			gs.handleUnequip(client, character, constant.EQUIPMENT_PARTS_TOP, int16(storageSlot))
		}
	}

	// Swap equipment
	character.Equipments[parts], inven.Items[slot] = new, old
	if !swap {
		delete(inven.Items, slot)
	}

	// Send swap inventory slot response
	client.Send(&resp.SwapInventorySlot{
		InventoryType:   constant.INVENTORY_TYPE_EQUIPMENT,
		Source:          slot,
		Dest:            int16(parts),
		EquipmentAction: resp.EQUIPMENT_ACTION_TYPE_ON,
	}, types.SEND_POLICY_ENCRYPT)

	// Broadcast character look update
	mapInstance := gs.GetMap(character.GetMap())
	if mapInstance != nil {
		mapInstance.BroadcastToPlayers(&resp.UpdateCharacterLook{
			Character: character,
		}, types.SEND_POLICY_ENCRYPT, character.GetID())
	}
}

// handleDrop handles dropping items
func (gs *GameServer) handleDrop(client *client.GameClient, character *entity.Character, invenType constant.InventoryType, slot int16, count uint16) {
	item, ok := character.Inventory[invenType].Items[slot]
	if !ok {
		return
	}

	// Remove item from inventory
	removed := (item.Reduce(count) == 0)
	if removed {
		client.Send(&resp.RemoveInventorySlot{
			InventoryType: invenType,
			Slot:          slot,
		}, types.SEND_POLICY_ENCRYPT)
		delete(character.Inventory[invenType].Items, slot)
	} else {
		client.Send(&resp.UpdateInventorySlot{
			InventoryType: invenType,
			Slot:          slot,
			Item:          item,
		}, types.SEND_POLICY_ENCRYPT)
	}

	// Create drop item
	spawned := item.Clone(count)
	spawned.BindDrop(&entity.Drop{
		Object: &entity.Object{
			Position: character.Position,
		},
		Owner:        character.ID,
		SpawnedPoint: character.Position,
		DropType:     constant.DROP_TYPE_FFA,
	})

	// Spawn item on map (following old server pattern)
	mapInstance := gs.GetMap(character.GetMap())
	if mapInstance != nil {
		if err := mapInstance.SpawnItem(spawned, character.ID, constant.DROP_TYPE_FFA); err != nil {
			log.Printf("Failed to spawn item on map: %v", err)
		} else {
			log.Printf("Drop item spawned: %d x %d at position %v", spawned.GetSpec().GetID(), count, character.Position)
		}
	}
}

// handleMoveItemInternal handles moving items within inventory
func (gs *GameServer) handleMoveItemInternal(client *client.GameClient, character *entity.Character, invenType constant.InventoryType, sourceSlot int16, destSlot int16) {
	inven := character.Inventory[invenType]
	src, ok := inven.Items[sourceSlot]
	if !ok {
		return
	}

	dst, ok := inven.Items[destSlot]

	if !ok {
		// Move to empty slot
		inven.Items[destSlot] = inven.Items[sourceSlot]
		delete(inven.Items, sourceSlot)
		client.Send(&resp.SwapInventorySlot{
			InventoryType: invenType,
			Source:        sourceSlot,
			Dest:          destSlot,
		}, types.SEND_POLICY_ENCRYPT)
		return
	}

	// Check if items are the same type
	specSrc := src.GetSpec()
	specDst := dst.GetSpec()
	if specSrc != specDst {
		// Swap different items
		inven.Items[sourceSlot], inven.Items[destSlot] = inven.Items[destSlot], inven.Items[sourceSlot]
		client.Send(&resp.SwapInventorySlot{
			InventoryType: invenType,
			Source:        sourceSlot,
			Dest:          destSlot,
		}, types.SEND_POLICY_ENCRYPT)
		return
	}

	// Merge same items
	limit := min(src.GetCount(), specSrc.GetCapacity()-dst.GetCount())
	dst.Increase(limit)
	if src.Reduce(limit) == 0 {
		client.Send(&resp.FullMergeInventorySlot{
			InventoryType: invenType,
			Source:        sourceSlot,
			Dest:          destSlot,
			Count:         dst.GetCount(),
		}, types.SEND_POLICY_ENCRYPT)
		delete(inven.Items, sourceSlot)
	} else {
		client.Send(&resp.PartialMergeInventorySlot{
			InventoryType: invenType,
			Source:        sourceSlot,
			Dest:          destSlot,
			SourceCount:   src.GetCount(),
			DestCount:     dst.GetCount(),
		}, types.SEND_POLICY_ENCRYPT)
	}
}

// handleMergeItems handles merging items in inventory
func (gs *GameServer) handleMergeItems(client *client.GameClient, character *entity.Character, inventoryType constant.InventoryType) {
	inven := character.Inventory[inventoryType]
	buckets := map[data.ItemSpec]map[int16]entity.Item{}

	for i := range inven.SlotLimit {
		item := inven.Items[int16(i+1)]
		if item == nil {
			continue
		}

		spec := item.GetSpec()
		if buckets[spec] == nil {
			buckets[spec] = map[int16]entity.Item{}
		}

		buckets[spec][int16(i+1)] = item
	}

	for spec, bucket := range buckets {
		count := uint16(0)
		for _, v := range bucket {
			count += v.GetCount()
		}

		capacity := spec.GetCapacity()
		for slot, item := range bucket {
			value := min(capacity, count)
			if item.GetCount() != value {
				item.SetCount(value)
				if value == 0 {
					client.Send(&resp.RemoveInventorySlot{
						InventoryType: inventoryType,
						Slot:          slot,
					}, types.SEND_POLICY_ENCRYPT)
					delete(inven.Items, slot)
				} else {
					client.Send(&resp.UpdateInventorySlot{
						InventoryType: inventoryType,
						Slot:          slot,
						Item:          item,
					}, types.SEND_POLICY_ENCRYPT)
				}
			}
			count -= value
		}
	}
}

// handleSortInventoryInternal handles sorting inventory
func (gs *GameServer) handleSortInventoryInternal(client *client.GameClient, character *entity.Character, inventoryType constant.InventoryType) {
	inven := character.Inventory[inventoryType]
	n := inven.SlotLimit
	buffer := make([]entity.Item, n)
	for i := range n {
		buffer[i] = inven.Items[int16(i+1)]
	}

	less := func(item1, item2 entity.Item) bool {
		if item1 == nil && item2 == nil {
			return false
		}
		if item1 == nil {
			return false
		}
		if item2 == nil {
			return true
		}
		id1, id2 := item1.GetSpec().GetID(), item2.GetSpec().GetID()
		if id1 != id2 {
			return id1 < id2
		}
		return item1.GetCount() > item2.GetCount()
	}

	partition := func(low, high int) int {
		pivot := buffer[(low+high)/2]
		i1, i2 := low, high
		for i1 <= i2 {
			for less(buffer[i1], pivot) {
				i1++
			}
			for less(pivot, buffer[i2]) {
				i2--
			}
			if i1 <= i2 {
				buffer[i1], buffer[i2] = buffer[i2], buffer[i1]

				client.Send(&resp.SwapInventorySlot{
					InventoryType: inventoryType,
					Source:        int16(i1 + 1),
					Dest:          int16(i2 + 1),
				}, types.SEND_POLICY_ENCRYPT)
				i1++
				i2--
			}
		}
		return i1
	}

	var qsort func(low, high int)
	qsort = func(low, high int) {
		if low < high {
			p := partition(low, high)
			qsort(low, p-1)
			qsort(p, high)
		}
	}

	qsort(0, int(n-1))

	inven.Items = map[int16]entity.Item{}
	for i := range buffer {
		if buffer[i] != nil {
			inven.Items[int16(i+1)] = buffer[i]
		}
	}
}
