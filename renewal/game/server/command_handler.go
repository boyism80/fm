package server

import (
	"fmt"
	"log"
	"math"
	"strconv"

	"github.com/boyism80/fm/common/types"
	"github.com/boyism80/fm/game/constant"
	"github.com/boyism80/fm/game/entity"
	"github.com/boyism80/fm/game/protocol/resp"
	"github.com/boyism80/fm/renewal/game/client"
)

// CommandHandler handles chat commands for the game server
type CommandHandler struct {
	gameServer *GameServer
}

// NewCommandHandler creates a new command handler instance
func NewCommandHandler(gameServer *GameServer) *CommandHandler {
	return &CommandHandler{
		gameServer: gameServer,
	}
}

// Handle processes a command with parameters
func (ch *CommandHandler) Handle(gameClient *client.GameClient, params ...string) error {
	if len(params) == 0 {
		return fmt.Errorf("no command specified")
	}

	command := params[0]
	args := params[1:]

	switch command {
	case "아이템생성":
		return ch.handleCreateItem(gameClient, args...)
	case "메소초기화":
		return ch.handleClearMeso(gameClient, args...)
	case "메소얻기":
		return ch.handleGainMeso(gameClient, args...)
	case "풀메소":
		return ch.handleFullMeso(gameClient, args...)
	case "맵이동":
		return ch.handleChangeMap(gameClient, args...)
	case "좌표":
		return ch.handleGetPosition(gameClient, args...)
	case "체력바꾸기":
		return ch.handleChangeHp(gameClient, args...)
	case "마력바꾸기":
		return ch.handleChangeMp(gameClient, args...)
	case "힘바꾸기":
		return ch.handleChangeStr(gameClient, args...)
	case "민첩바꾸기":
		return ch.handleChangeDex(gameClient, args...)
	case "지능바꾸기":
		return ch.handleChangeInt(gameClient, args...)
	case "행운바꾸기":
		return ch.handleChangeLuk(gameClient, args...)
	case "모든스탯바꾸기":
		return ch.handleChangeAllStats(gameClient, args...)
	case "레벨바꾸기":
		return ch.handleChangeLevel(gameClient, args...)
	case "무적":
		return ch.handleInvincible(gameClient, args...)
	case "직업바꾸기":
		return ch.handleChangeClass(gameClient, args...)
	default:
		return fmt.Errorf("unknown command: %s", command)
	}
}

// handleCreateItem creates an item for the character
func (ch *CommandHandler) handleCreateItem(gameClient *client.GameClient, args ...string) error {
	if len(args) < 1 {
		return fmt.Errorf("missing itemId or item name")
	}

	var itemId int
	var err error

	// Try to parse as item ID first
	itemId, err = strconv.Atoi(args[0])
	if err != nil {
		// If not a number, try to find by item name
		itemIdUint, ok := ch.gameServer.resources.NameToItem(args[0])
		if !ok {
			return fmt.Errorf("invalid itemId or item name: %s", args[0])
		}
		itemId = int(itemIdUint)
	}

	count := 1
	if len(args) >= 2 {
		if parsedCount, err := strconv.Atoi(args[1]); err == nil {
			count = parsedCount
		} else {
			log.Printf("Command: Invalid count, using default 1: %s", args[1])
		}
	}

	character := gameClient.GetCharacter()
	if character == nil {
		return fmt.Errorf("character not found")
	}

	// Get item spec from resources
	itemSpec, ok := ch.gameServer.resources.Items[uint32(itemId)]
	if !ok {
		return fmt.Errorf("item spec not found for id: %d", itemId)
	}

	// Create the item
	item, err := entity.NewItem(itemSpec, uint16(count))
	if err != nil {
		return fmt.Errorf("failed to create item: %v", err)
	}

	// Get inventory type for the item
	inventoryType := item.GetInventoryType()
	inventory := character.Inventory[inventoryType]
	if inventory == nil {
		return fmt.Errorf("inventory not found for type: %v", inventoryType)
	}

	// Find next available slot
	nextSlot, ok := inventory.NextSlot()
	if !ok {
		return fmt.Errorf("inventory is full for type: %v", inventoryType)
	}

	// Add item to inventory
	inventory.Items[int16(nextSlot)] = item

	// Send AddItem packet to client
	gameClient.Send(&resp.AddItem{
		IsDrop:        false,
		Slot:          nextSlot,
		InventoryType: inventoryType,
		Item:          item,
	}, types.SEND_POLICY_ENCRYPT)

	log.Printf("Command: Created item %d, count %d in slot %d for character %d", itemId, count, nextSlot, character.GetID())
	return nil
}

// handleClearMeso clears the character's meso
func (ch *CommandHandler) handleClearMeso(gameClient *client.GameClient, args ...string) error {
	character := gameClient.GetCharacter()
	if character == nil {
		return fmt.Errorf("character not found")
	}

	character.Meso = 0
	gameClient.Send(&resp.UpdateStats{
		Stats: map[constant.Stat]int32{
			constant.STAT_MESO: character.Meso,
		},
	}, types.SEND_POLICY_ENCRYPT)

	log.Printf("Command: Cleared meso for character %d", character.GetID())
	return nil
}

// handleGainMeso adds meso to the character
func (ch *CommandHandler) handleGainMeso(gameClient *client.GameClient, args ...string) error {
	if len(args) < 1 {
		return fmt.Errorf("missing meso amount")
	}

	amount, err := strconv.Atoi(args[0])
	if err != nil {
		return fmt.Errorf("invalid meso amount: %s", args[0])
	}

	character := gameClient.GetCharacter()
	if character == nil {
		return fmt.Errorf("character not found")
	}

	character.Meso += int32(amount)
	gameClient.Send(&resp.UpdateStats{
		Stats: map[constant.Stat]int32{
			constant.STAT_MESO: character.Meso,
		},
	}, types.SEND_POLICY_ENCRYPT)

	log.Printf("Command: Added %d meso to character %d", amount, character.GetID())
	return nil
}

// handleFullMeso sets the character's meso to maximum
func (ch *CommandHandler) handleFullMeso(gameClient *client.GameClient, args ...string) error {
	character := gameClient.GetCharacter()
	if character == nil {
		return fmt.Errorf("character not found")
	}

	character.Meso = math.MaxInt32
	gameClient.Send(&resp.UpdateStats{
		Stats: map[constant.Stat]int32{
			constant.STAT_MESO: character.Meso,
		},
	}, types.SEND_POLICY_ENCRYPT)

	log.Printf("Command: Set full meso for character %d", character.GetID())
	return nil
}

// handleChangeMap changes the character's map
func (ch *CommandHandler) handleChangeMap(gameClient *client.GameClient, args ...string) error {
	if len(args) < 1 {
		return fmt.Errorf("missing mapId or map name")
	}

	var mapId int
	var err error

	// Try to parse as map ID first
	mapId, err = strconv.Atoi(args[0])
	if err != nil {
		// If not a number, try to find by map name
		mapIdUint, ok := ch.gameServer.resources.NameToMap(args[0])
		if !ok {
			return fmt.Errorf("invalid mapId or map name: %s", args[0])
		}
		mapId = int(mapIdUint)
	}

	character := gameClient.GetCharacter()
	if character == nil {
		return fmt.Errorf("character not found")
	}

	// Get current map
	currentMap := ch.gameServer.GetMap(character.GetMap())
	if currentMap == nil {
		return fmt.Errorf("current map not found")
	}

	// Get target map
	targetMap := ch.gameServer.GetMap(uint32(mapId))
	if targetMap == nil {
		return fmt.Errorf("target map %d not found", mapId)
	}

	// Remove character from current map
	currentMap.RemovePlayer(character.GetID())

	// Update character's map
	character.Map = uint32(mapId)

	// Set character position to spawn point of new map
	mapSpec := targetMap.GetSpec()
	if mapSpec != nil && len(mapSpec.Portals) > 0 {
		// Use first portal as spawn point
		for portalId, portal := range mapSpec.Portals {
			character.Position = portal.Position
			character.SpawnPoint = portalId
			break
		}
	}

	// Add character to new map
	if err := targetMap.AddPlayer(character.GetID(), character, true); err != nil {
		return fmt.Errorf("failed to add character to new map: %v", err)
	}

	log.Printf("Command: Changed map from %d to %d for character %d", character.GetMap(), mapId, character.GetID())
	return nil
}

// handleGetPosition shows the character's current position
func (ch *CommandHandler) handleGetPosition(gameClient *client.GameClient, args ...string) error {
	character := gameClient.GetCharacter()
	if character == nil {
		return fmt.Errorf("character not found")
	}

	log.Printf("Command: Character %d position - Map: %d, Position: %v",
		character.GetID(), character.GetMap(), character.Position)
	return nil
}

// handleChangeHp changes the character's HP
func (ch *CommandHandler) handleChangeHp(gameClient *client.GameClient, args ...string) error {
	if len(args) < 1 {
		return fmt.Errorf("missing HP value")
	}

	hp, err := strconv.Atoi(args[0])
	if err != nil {
		return fmt.Errorf("invalid HP value: %s", args[0])
	}

	// Limit to maximum value
	if hp > 32767 {
		hp = 32767
	}

	character := gameClient.GetCharacter()
	if character == nil {
		return fmt.Errorf("character not found")
	}

	character.Hp = uint16(hp)
	character.MaxHp = uint16(hp)
	gameClient.Send(&resp.UpdateStats{
		Stats: map[constant.Stat]int32{
			constant.STAT_HP:     int32(character.Hp),
			constant.STAT_MAX_HP: int32(character.MaxHp),
		},
	}, types.SEND_POLICY_ENCRYPT)

	log.Printf("Command: Changed HP to %d for character %d", hp, character.GetID())
	return nil
}

// handleChangeMp changes the character's MP
func (ch *CommandHandler) handleChangeMp(gameClient *client.GameClient, args ...string) error {
	if len(args) < 1 {
		return fmt.Errorf("missing MP value")
	}

	mp, err := strconv.Atoi(args[0])
	if err != nil {
		return fmt.Errorf("invalid MP value: %s", args[0])
	}

	// Limit to maximum value
	if mp > 32767 {
		mp = 32767
	}

	character := gameClient.GetCharacter()
	if character == nil {
		return fmt.Errorf("character not found")
	}

	character.Mp = uint16(mp)
	character.MaxMp = uint16(mp)
	gameClient.Send(&resp.UpdateStats{
		Stats: map[constant.Stat]int32{
			constant.STAT_MP:     int32(character.Mp),
			constant.STAT_MAX_MP: int32(character.MaxMp),
		},
	}, types.SEND_POLICY_ENCRYPT)

	log.Printf("Command: Changed MP to %d for character %d", mp, character.GetID())
	return nil
}

// handleChangeStr changes the character's STR
func (ch *CommandHandler) handleChangeStr(gameClient *client.GameClient, args ...string) error {
	if len(args) < 1 {
		return fmt.Errorf("missing STR value")
	}

	str, err := strconv.Atoi(args[0])
	if err != nil {
		return fmt.Errorf("invalid STR value: %s", args[0])
	}

	// Limit to maximum value
	if str > 32767 {
		str = 32767
	}

	character := gameClient.GetCharacter()
	if character == nil {
		return fmt.Errorf("character not found")
	}

	character.Str = uint16(str)
	gameClient.Send(&resp.UpdateStats{
		Stats: map[constant.Stat]int32{
			constant.STAT_STR: int32(character.Str),
		},
	}, types.SEND_POLICY_ENCRYPT)

	log.Printf("Command: Changed STR to %d for character %d", str, character.GetID())
	return nil
}

// handleChangeDex changes the character's DEX
func (ch *CommandHandler) handleChangeDex(gameClient *client.GameClient, args ...string) error {
	if len(args) < 1 {
		return fmt.Errorf("missing DEX value")
	}

	dex, err := strconv.Atoi(args[0])
	if err != nil {
		return fmt.Errorf("invalid DEX value: %s", args[0])
	}

	// Limit to maximum value
	if dex > 32767 {
		dex = 32767
	}

	character := gameClient.GetCharacter()
	if character == nil {
		return fmt.Errorf("character not found")
	}

	character.Dex = uint16(dex)
	gameClient.Send(&resp.UpdateStats{
		Stats: map[constant.Stat]int32{
			constant.STAT_DEX: int32(character.Dex),
		},
	}, types.SEND_POLICY_ENCRYPT)

	log.Printf("Command: Changed DEX to %d for character %d", dex, character.GetID())
	return nil
}

// handleChangeInt changes the character's INT
func (ch *CommandHandler) handleChangeInt(gameClient *client.GameClient, args ...string) error {
	if len(args) < 1 {
		return fmt.Errorf("missing INT value")
	}

	intVal, err := strconv.Atoi(args[0])
	if err != nil {
		return fmt.Errorf("invalid INT value: %s", args[0])
	}

	// Limit to maximum value
	if intVal > 32767 {
		intVal = 32767
	}

	character := gameClient.GetCharacter()
	if character == nil {
		return fmt.Errorf("character not found")
	}

	character.Int = uint16(intVal)
	gameClient.Send(&resp.UpdateStats{
		Stats: map[constant.Stat]int32{
			constant.STAT_INT: int32(character.Int),
		},
	}, types.SEND_POLICY_ENCRYPT)

	log.Printf("Command: Changed INT to %d for character %d", intVal, character.GetID())
	return nil
}

// handleChangeLuk changes the character's LUK
func (ch *CommandHandler) handleChangeLuk(gameClient *client.GameClient, args ...string) error {
	if len(args) < 1 {
		return fmt.Errorf("missing LUK value")
	}

	luk, err := strconv.Atoi(args[0])
	if err != nil {
		return fmt.Errorf("invalid LUK value: %s", args[0])
	}

	// Limit to maximum value
	if luk > 32767 {
		luk = 32767
	}

	character := gameClient.GetCharacter()
	if character == nil {
		return fmt.Errorf("character not found")
	}

	character.Luk = uint16(luk)
	gameClient.Send(&resp.UpdateStats{
		Stats: map[constant.Stat]int32{
			constant.STAT_LUK: int32(character.Luk),
		},
	}, types.SEND_POLICY_ENCRYPT)

	log.Printf("Command: Changed LUK to %d for character %d", luk, character.GetID())
	return nil
}

// handleChangeAllStats changes all stats of the character
func (ch *CommandHandler) handleChangeAllStats(gameClient *client.GameClient, args ...string) error {
	if len(args) < 1 {
		return fmt.Errorf("missing stat value")
	}

	statValue, err := strconv.Atoi(args[0])
	if err != nil {
		return fmt.Errorf("invalid stat value: %s", args[0])
	}

	// Limit to maximum value
	if statValue > 32767 {
		statValue = 32767
	}

	character := gameClient.GetCharacter()
	if character == nil {
		return fmt.Errorf("character not found")
	}

	character.Str = uint16(statValue)
	character.Dex = uint16(statValue)
	character.Int = uint16(statValue)
	character.Luk = uint16(statValue)

	gameClient.Send(&resp.UpdateStats{
		Stats: map[constant.Stat]int32{
			constant.STAT_STR: int32(character.Str),
			constant.STAT_DEX: int32(character.Dex),
			constant.STAT_INT: int32(character.Int),
			constant.STAT_LUK: int32(character.Luk),
		},
	}, types.SEND_POLICY_ENCRYPT)

	log.Printf("Command: Changed all stats to %d for character %d", statValue, character.GetID())
	return nil
}

// handleChangeLevel changes the character's level
func (ch *CommandHandler) handleChangeLevel(gameClient *client.GameClient, args ...string) error {
	if len(args) < 1 {
		return fmt.Errorf("missing level value")
	}

	level, err := strconv.Atoi(args[0])
	if err != nil {
		return fmt.Errorf("invalid level value: %s", args[0])
	}

	// Limit to valid range (1-255)
	if level < 1 {
		level = 1
	} else if level > 255 {
		level = 255
	}

	character := gameClient.GetCharacter()
	if character == nil {
		return fmt.Errorf("character not found")
	}

	character.Level = uint8(level)
	gameClient.Send(&resp.UpdateStats{
		Stats: map[constant.Stat]int32{
			constant.STAT_LEVEL: int32(character.Level),
		},
	}, types.SEND_POLICY_ENCRYPT)

	log.Printf("Command: Changed level to %d for character %d", level, character.GetID())
	return nil
}

// handleInvincible toggles the character's invincibility
func (ch *CommandHandler) handleInvincible(gameClient *client.GameClient, args ...string) error {
	character := gameClient.GetCharacter()
	if character == nil {
		return fmt.Errorf("character not found")
	}

	character.Invincible = !character.Invincible
	status := "enabled"
	if !character.Invincible {
		status = "disabled"
	}

	log.Printf("Command: Invincibility %s for character %d", status, character.GetID())
	return nil
}

// handleChangeClass changes the character's class
func (ch *CommandHandler) handleChangeClass(gameClient *client.GameClient, args ...string) error {
	if len(args) < 1 {
		return fmt.Errorf("missing class value")
	}

	class, err := strconv.Atoi(args[0])
	if err != nil {
		return fmt.Errorf("invalid class value: %s", args[0])
	}

	character := gameClient.GetCharacter()
	if character == nil {
		return fmt.Errorf("character not found")
	}

	character.Class = uint16(class)
	gameClient.Send(&resp.UpdateStats{
		Stats: map[constant.Stat]int32{
			constant.STAT_JOB: int32(character.Class),
		},
	}, types.SEND_POLICY_ENCRYPT)

	log.Printf("Command: Changed class to %d for character %d", class, character.GetID())
	return nil
}
