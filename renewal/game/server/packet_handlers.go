package server

import (
	"fmt"
	"log"

	common_req "github.com/boyism80/fm/common/protocol/req"
	"github.com/boyism80/fm/common/stream"
	"github.com/boyism80/fm/common/types"
	"github.com/boyism80/fm/core"
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
	character := entity.NewDummyCharacter(ctx.Client, nil, request.PlayerId, name, nil)

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

	log.Printf("Normal chat packet received from %s - Message: %s, DontRecord: %t",
		ctx.Client.GetConnection().RemoteAddr(), request.Message, request.DontRecordHistory)
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

	log.Printf("Attack packet received from %s - Attack Info: %+v",
		ctx.Client.GetConnection().RemoteAddr(), request.AttackInfo)
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

	log.Printf("Move item packet received from %s - Source: %d, Dest: %d, Count: %d",
		ctx.Client.GetConnection().RemoteAddr(), request.Source, request.Dest, request.Count)
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

	log.Printf("Sort inventory packet received from %s - Inventory Type: %d",
		ctx.Client.GetConnection().RemoteAddr(), request.InventoryType)
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

	log.Printf("Item loot packet received from %s - OID: %d",
		ctx.Client.GetConnection().RemoteAddr(), request.OID)
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

	log.Printf("Drop meso packet received from %s - Count: %d",
		ctx.Client.GetConnection().RemoteAddr(), request.Count)
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

	log.Printf("Warp packet received from %s - Target: %d",
		ctx.Client.GetConnection().RemoteAddr(), request.Target)
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

	log.Printf("NPC control packet received from %s - Bytes: %d bytes",
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

	log.Printf("Move mob packet received from %s - OID: %d",
		ctx.Client.GetConnection().RemoteAddr(), request.OID)
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
	return nil
}
