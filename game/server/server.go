package server

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"sort"

	"sync"
	"syscall"
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/core"
	c_actor "github.com/boyism80/fm/core/actor"

	"github.com/boyism80/fm/core/luax"
	g_actor "github.com/boyism80/fm/game/actor"
	"github.com/boyism80/fm/game/client"
	"github.com/boyism80/fm/game/constant"
	"github.com/boyism80/fm/game/entity"
	"github.com/boyism80/fm/game/wz"
	lua "github.com/yuin/gopher-lua"
)

// GameServer represents the game server for MapleStory private server
type GameServer struct {
	server         *core.Server
	config         *GameConfig
	resources      *wz.Resources          // Game data resources
	maps           map[uint32]*entity.Map // Map instances by map ID
	mapsMutex      sync.RWMutex
	commandHandler *CommandHandler
	packetHandlers *PacketHandlerRegistry
	context        *GameServerContext
	actorSystem    *c_actor.ActorSystem
	actorRegistry  *c_actor.ActorRegistry
	nilMapActorPID *actor.PID
}

func registerSkillConstants(luaState *lua.LState) {
	skillTable := luaState.NewTable()
	for key, skillID := range constant.AllSkillConstants() {
		skillTable.RawSetString(key, lua.LNumber(skillID))
	}
	luaState.SetGlobal("SKILL", skillTable)
}

func skillToLuaWzTable(luaState *lua.LState, skill *wz.Skill) *lua.LTable {
	skillTable := luaState.NewTable()
	skillTable.RawSetString("id", lua.LNumber(skill.ID))
	skillTable.RawSetString("skill_id", lua.LNumber(skill.ID))
	skillTable.RawSetString("max_level", lua.LNumber(skill.MaxLevel))
	skillTable.RawSetString("master_level", lua.LNumber(skill.MasterLevel))
	skillTable.RawSetString("invisible", lua.LBool(skill.Invisible))
	return skillTable
}

func (gs *GameServer) GetServer() *core.Server {
	return gs.server
}

func (gs *GameServer) GetServerContext() core.ServerContext {
	return gs.context
}

func (gs *GameServer) GetRootContext() *actor.RootContext {
	return gs.server.GetRootContext()
}

// GameContext provides access to game resources and services
type GameContext interface {
	GetResources() *wz.Resources
	GetMap(mapId uint32) *entity.Map
	GetExpRate() int  // Returns experience rate multiplier
	GetDropRate() int // Returns drop rate multiplier
	GetMesoRate() int // Returns meso rate multiplier
	RequestWarp(character *entity.Character, targetMap *entity.Map, spawnPoint uint8) error
}

// GameConfig holds game server specific configuration
type GameConfig struct {
	// LogicThreadCount int    // Commented out: LogicThread removed
	Host       string // Game server host address
	Port       int    // Game server port number
	WzPath     string // Path to WZ files directory
	WorldName  string // World/Channel name
	MaxPlayers int    // Maximum number of players per world
	ExpRate    int    // Experience rate multiplier
	DropRate   int    // Drop rate multiplier
	MesoRate   int    // Meso rate multiplier
}

// NewGameServer creates a new game server with specified configuration
func NewGameServer(config *GameConfig) (*GameServer, error) {
	serverConfig := &core.ServerConfig{
		Host: config.Host,
		Port: config.Port,
		ClientFactory: func(conn net.Conn, clientID int) (core.Client, error) {
			return client.NewGameClient(conn, clientID)
		},
	}
	server, err := core.NewServer(serverConfig)
	if err != nil {
		return nil, err
	}

	// Load game resources
	// FindWzPath is called inside NewResources, so we just pass the config path
	resources := wz.NewResources(config.WzPath)
	if resources == nil {
		return nil, fmt.Errorf("failed to load game resources")
	}

	// Create ServerContext (temporary, will be set after gameServer is created)
	// MapListener will be set after gameServer is created in preCreateMaps
	context := NewGameServerContext(config.WzPath, nil)

	// Create Actor system
	actorSystem := c_actor.NewActorSystem()
	actorRegistry := c_actor.NewActorRegistry(actorSystem)

	// Set RootContext in server for actor message sending
	server.SetRootContext(actorSystem.GetRoot())

	gs := &GameServer{
		server:        server,
		config:        config,
		resources:     resources,
		maps:          make(map[uint32]*entity.Map),
		context:       context,
		actorSystem:   actorSystem,
		actorRegistry: actorRegistry,
	}

	// Set gameServer reference in context
	context.gs = gs

	// Initialize command handler
	gs.commandHandler = NewCommandHandler(gs)

	// Initialize packet handler registry
	gs.packetHandlers = NewPacketHandlerRegistry(gs)

	// Register hook so each root LuaState (e.g. per map actor) gets game-specific types
	luax.RegisterOnCreateHook(func(luaState *lua.LState) {
		// Register Lua types with inheritance
		luax.RegisterLuaType[*entity.Object](luaState)
		luax.RegisterLuaDerivedType[*entity.Drop, *entity.Object](luaState)
		luax.RegisterLuaDerivedType[*entity.Meso, *entity.Drop](luaState)
		luax.RegisterLuaDerivedType[*entity.Life, *entity.Object](luaState)
		luax.RegisterLuaDerivedType[*entity.Character, *entity.Life](luaState)
		luax.RegisterLuaDerivedType[*entity.Mob, *entity.Life](luaState)
		luax.RegisterLuaDerivedType[*entity.Npc, *entity.Object](luaState)
		luax.RegisterLuaType[*entity.Map](luaState)
		luax.RegisterLuaType[*entity.SkillEntry](luaState)
		luax.RegisterLuaType[*entity.ItemCore](luaState)
		luax.RegisterLuaDerivedType[*entity.EquipmentCore, *entity.ItemCore](luaState)
		luax.RegisterLuaDerivedType[*entity.Weapon, *entity.EquipmentCore](luaState)
		luax.RegisterLuaDerivedType[*entity.Shield, *entity.EquipmentCore](luaState)
		luax.RegisterLuaDerivedType[*entity.Cap, *entity.EquipmentCore](luaState)
		luax.RegisterLuaDerivedType[*entity.Face, *entity.EquipmentCore](luaState)
		luax.RegisterLuaDerivedType[*entity.Accessory, *entity.EquipmentCore](luaState)
		luax.RegisterLuaDerivedType[*entity.Top, *entity.EquipmentCore](luaState)
		luax.RegisterLuaDerivedType[*entity.Pants, *entity.EquipmentCore](luaState)
		luax.RegisterLuaDerivedType[*entity.Shoes, *entity.EquipmentCore](luaState)
		luax.RegisterLuaDerivedType[*entity.Glove, *entity.EquipmentCore](luaState)
		luax.RegisterLuaDerivedType[*entity.Cape, *entity.EquipmentCore](luaState)
		luax.RegisterLuaDerivedType[*entity.RingEquip, *entity.EquipmentCore](luaState)
		luax.RegisterLuaType[*wz.ItemWzCore](luaState)
		luax.RegisterLuaDerivedType[*wz.ItemWzEquipment, *wz.ItemWzCore](luaState)
		luax.RegisterLuaDerivedType[*wz.ItemWzWeapon, *wz.ItemWzEquipment](luaState)
		luax.RegisterLuaDerivedType[*wz.ItemWzArmor, *wz.ItemWzEquipment](luaState)
		luax.RegisterLuaDerivedType[*wz.ItemWzConsume, *wz.ItemWzCore](luaState)
		luax.RegisterLuaDerivedType[*wz.ItemWzPet, *wz.ItemWzCore](luaState)
		luax.RegisterLuaDerivedType[*wz.ItemWzGeneralItem, *wz.ItemWzCore](luaState)
		luax.RegisterLuaDerivedType[*wz.ItemWzCashItem, *wz.ItemWzCore](luaState)
		luax.RegisterLuaDerivedType[*wz.ItemWzInstallation, *wz.ItemWzCore](luaState)
		luax.RegisterLuaDerivedType[*entity.Consume, *entity.ItemCore](luaState)
		luax.RegisterLuaDerivedType[*entity.CashItem, *entity.ItemCore](luaState)
		luax.RegisterLuaDerivedType[*entity.GeneralItem, *entity.ItemCore](luaState)
		luax.RegisterLuaDerivedType[*entity.Installation, *entity.ItemCore](luaState)
		luax.RegisterLuaDerivedType[*entity.Pet, *entity.ItemCore](luaState)

		// BuffFlag table for scripts
		buffFlagTable := luaState.NewTable()
		for name, bf := range constant.AllBuffFlags() {
			entry := luaState.NewTable()
			entry.RawSetString("mask", lua.LNumber(bf.Mask))
			entry.RawSetString("position", lua.LNumber(bf.Position))
			buffFlagTable.RawSetString(name, entry)
		}
		luaState.SetGlobal("BuffFlag", buffFlagTable)
		debuffTable := luaState.NewTable()
		for name, d := range constant.AllDebuffs() {
			entry := luaState.NewTable()
			entry.RawSetString("mask", lua.LNumber(d.Mask))
			debuffTable.RawSetString(name, entry)
		}
		luaState.SetGlobal("Debuff", debuffTable)
		weaponTypeTable := luaState.NewTable()
		for name, wt := range constant.AllWeaponTypes() {
			weaponTypeTable.RawSetString(name, lua.LNumber(wt))
		}
		luaState.SetGlobal("WeaponType", weaponTypeTable)
		consumeTypeTable := luaState.NewTable()
		for name, ct := range constant.AllConsumeTypes() {
			consumeTypeTable.RawSetString(name, lua.LNumber(ct))
		}
		luaState.SetGlobal("ConsumeType", consumeTypeTable)
		registerSkillConstants(luaState)
		equipmentPartTable := luaState.NewTable()
		equipmentPartTable.RawSetString("Cap", lua.LNumber(constant.EQUIPMENT_PARTS_CAP))
		equipmentPartTable.RawSetString("Face", lua.LNumber(constant.EQUIPMENT_PARTS_FACE))
		equipmentPartTable.RawSetString("Eye", lua.LNumber(constant.EQUIPMENT_PARTS_EYE))
		equipmentPartTable.RawSetString("Ear", lua.LNumber(constant.EQUIPMENT_PARTS_EAR))
		equipmentPartTable.RawSetString("Top", lua.LNumber(constant.EQUIPMENT_PARTS_TOP))
		equipmentPartTable.RawSetString("Pants", lua.LNumber(constant.EQUIPMENT_PARTS_PANTS))
		equipmentPartTable.RawSetString("Shoes", lua.LNumber(constant.EQUIPMENT_PARTS_SHOES))
		equipmentPartTable.RawSetString("Glove", lua.LNumber(constant.EQUIPMENT_PARTS_GLOVE))
		equipmentPartTable.RawSetString("Cape", lua.LNumber(constant.EQUIPMENT_PARTS_CAPE))
		equipmentPartTable.RawSetString("Shield", lua.LNumber(constant.EQUIPMENT_PARTS_SHIELD))
		equipmentPartTable.RawSetString("Weapon", lua.LNumber(constant.EQUIPMENT_PARTS_WEAPON))
		equipmentPartTable.RawSetString("Ring", lua.LNumber(constant.EQUIPMENT_PARTS_RING))
		luaState.SetGlobal("EquipmentPart", equipmentPartTable)
		inventoryTypeTable := luaState.NewTable()
		inventoryTypeTable.RawSetString("Equipment", lua.LNumber(constant.INVENTORY_TYPE_EQUIPMENT))
		inventoryTypeTable.RawSetString("Use", lua.LNumber(constant.INVENTORY_TYPE_CONSUME))
		inventoryTypeTable.RawSetString("Consume", lua.LNumber(constant.INVENTORY_TYPE_CONSUME))
		inventoryTypeTable.RawSetString("Installation", lua.LNumber(constant.INVENTORY_TYPE_INSTALLATION))
		inventoryTypeTable.RawSetString("Etc", lua.LNumber(constant.INVENTORY_TYPE_ETC))
		inventoryTypeTable.RawSetString("Cash", lua.LNumber(constant.INVENTORY_TYPE_CASH))
		luaState.SetGlobal("InventoryType", inventoryTypeTable)
		luax.RegisterFunc(luaState, "class_learnable_skill_wzs", func(L *lua.LState) int {
			class := uint16(L.CheckInt(1))
			result := L.NewTable()
			if gs.resources == nil {
				L.Push(result)
				return 1
			}

			classIDs := getClassAdvancementClasses(class)
			skillByID := make(map[uint32]*wz.Skill)
			for _, classID := range classIDs {
				skillIDStart := uint32(classID) * 10000
				skillIDEnd := skillIDStart + 9999
				for skillID := skillIDStart; skillID <= skillIDEnd; skillID++ {
					wzSkill := gs.resources.GetSkill(skillID)
					if wzSkill == nil {
						continue
					}
					skillByID[skillID] = wzSkill
				}
			}

			orderedSkillIDs := make([]int, 0, len(skillByID))
			for skillID := range skillByID {
				orderedSkillIDs = append(orderedSkillIDs, int(skillID))
			}
			sort.Ints(orderedSkillIDs)

			index := 1
			for _, skillID := range orderedSkillIDs {
				skillTable := skillToLuaWzTable(L, skillByID[uint32(skillID)])
				result.RawSetInt(index, skillTable)
				index++
			}

			L.Push(result)
			return 1
		})

		// Register utility functions (yield + AfterFunc send ResumeLua; pid from luax.GetThreadPID)
		luax.RegisterFunc(luaState, "sleep", func(L *lua.LState) int {
			duration := L.CheckNumber(1)
			pid := luax.GetThreadPID(L)
			if pid == nil {
				return 0
			}
			root := luax.GetRootLuaState(pid.String())
			if root == nil {
				return 0
			}
			d := time.Duration(float64(duration) * float64(time.Millisecond))
			time.AfterFunc(d, func() {
				gs.GetRootContext().Send(pid, &g_actor.ResumeLua{Root: root, Thread: L})
			})
			return L.Yield(lua.LNil)
		})
	})

	// Pre-create all maps
	gs.preCreateMaps()

	// Register packet handlers
	gs.registerPacketHandlers()

	// Register command handlers
	gs.registerCommandHandlers()

	// Set client disconnect handler
	server.SetOnClientDisconnect(func(client core.Client) {
		gs.handleClientDisconnect(client)
	})

	return gs, nil
}

// GetResources returns the game resources
func (gs *GameServer) GetResources() *wz.Resources {
	return gs.resources
}

func (gs *GameServer) GetExpRate() int {
	return gs.config.ExpRate
}

func (gs *GameServer) GetDropRate() int {
	return gs.config.DropRate
}

func (gs *GameServer) GetMesoRate() int {
	return gs.config.MesoRate
}

// GetLogicThread returns the logic thread for scheduling tasks
// Commented out: LogicThread removed, will be replaced with Actor model
// func (gs *GameServer) GetLogicThread() *core.LogicThread {
// 	logicThread, err := gs.server.GetLogicThreadByHash(0)
// 	if err != nil {
// 		return nil
// 	}
// 	return logicThread
// }

// preCreateMaps pre-creates all map instances from loaded resources
func (gs *GameServer) preCreateMaps() {
	log.Println("Setting up map listeners...")
	gameMapListener := NewGameMapListener(gs)

	// Create nil MapActor (for characters before map assignment)
	nilMapProps := actor.PropsFromProducer(func() actor.Actor {
		return &g_actor.MapActor{
			MapData: nil,
			Context: gs.context,
		}
	})

	nilMapPID := gs.actorRegistry.GetOrCreateActor(
		"map_nil",
		nilMapProps,
	)
	gs.nilMapActorPID = nilMapPID

	// Set nil MapActor PID in server for handling nil PID clients
	gs.server.SetNilMapActorPID(nilMapPID)

	log.Println("Pre-creating map instances...")
	for mapID := range gs.resources.Maps {
		mapInstance := entity.NewMap(mapID, gameMapListener, mapID, gs)
		gs.maps[mapID] = mapInstance

		// Create MapActor for this map
		props := actor.PropsFromProducer(func() actor.Actor {
			return &g_actor.MapActor{
				MapData: mapInstance,
				Context: gs.context,
			}
		})

		pid := gs.actorRegistry.GetOrCreateActor(
			fmt.Sprintf("map_%d", mapID),
			props,
		)

		mapInstance.SetActorPID(pid)
	}
	log.Printf("Pre-created %d map instances", len(gs.maps))
	log.Println("Map listeners configured")
}

// Start initializes and starts the game server
func (gs *GameServer) Start() error {
	log.Println("Starting MapleStory Game Server...")

	// Start the core server
	if err := gs.server.Start(gs.config.Host, gs.config.Port); err != nil {
		return err
	}

	log.Printf("Game server started on %s:%d", gs.config.Host, gs.config.Port)
	log.Printf("World: %s, Max Players: %d", gs.config.WorldName, gs.config.MaxPlayers)
	log.Printf("Rates: Exp=%dx, Drop=%dx, Meso=%dx",
		gs.config.ExpRate, gs.config.DropRate, gs.config.MesoRate)

	// Log resource information
	if gs.resources != nil {
		stringCount := 0
		if gs.resources.Strings != nil {
			stringCount = gs.resources.Strings.CountStrings()
		}
		log.Printf("Resources loaded: %d maps, %d monsters, %d items, %d drops, %d strings",
			len(gs.resources.Maps), len(gs.resources.Monsters),
			len(gs.resources.Items), len(gs.resources.Drops), stringCount)
	}

	return nil
}

// Stop gracefully shuts down the game server
func (gs *GameServer) Stop() error {
	log.Println("Shutting down game server...")

	// TODO: Save world state and character data before shutdown

	return gs.server.Stop()
}

// GetMap gets an existing map (returns nil if not found)
func (gs *GameServer) GetMap(mapID uint32) *entity.Map {
	gs.mapsMutex.RLock()
	defer gs.mapsMutex.RUnlock()
	return gs.maps[mapID]
}

// RequestWarp implements GameContext. Removes character from current map (if any) and sends WarpCharacter to the target map's actor.
func (gs *GameServer) RequestWarp(character *entity.Character, targetMap *entity.Map, spawnPoint uint8) error {
	if targetMap == nil {
		return fmt.Errorf("target map is nil")
	}
	currentMap := character.GetMap()
	if currentMap != nil {
		currentMap.RemovePlayer(character.GetID())
	}
	targetPID := targetMap.GetActorPID()
	if targetPID == nil {
		return fmt.Errorf("target map actor not found")
	}
	gs.GetRootContext().Send(targetPID, &g_actor.WarpCharacter{
		Character: character,
		Portal:    spawnPoint,
	})
	return nil
}

// handleClientDisconnect handles client disconnection by removing character from map
func (gs *GameServer) handleClientDisconnect(c core.Client) {
	client, ok := c.(*client.GameClient)
	if !ok {
		return
	}
	character := client.GetCharacter()
	if character == nil {
		return
	}
	mapInstance := character.GetMap()
	if mapInstance == nil {
		return
	}
	mapInstance.RemovePlayer(character.GetID())
}

// GetStats returns game server statistics for monitoring
func (gs *GameServer) GetStats() map[string]interface{} {
	stats := gs.server.GetStats()

	// Add game-specific stats
	stats["server_type"] = "game"
	stats["world_name"] = gs.config.WorldName
	stats["max_players"] = gs.config.MaxPlayers
	stats["current_players"] = stats["client_count"] // For now, client count = player count
	stats["exp_rate"] = gs.config.ExpRate
	stats["drop_rate"] = gs.config.DropRate
	stats["meso_rate"] = gs.config.MesoRate

	// Add resource stats
	if gs.resources != nil {
		stats["maps_loaded"] = len(gs.resources.Maps)
		stats["monsters_loaded"] = len(gs.resources.Monsters)
		stats["items_loaded"] = len(gs.resources.Items)
		stats["drops_loaded"] = len(gs.resources.Drops)
		if gs.resources.Strings != nil {
			stats["strings_loaded"] = gs.resources.Strings.CountStrings()
		} else {
			stats["strings_loaded"] = 0
		}
	}

	return stats
}

// RunGameServer runs the game server with signal handling
func RunGameServer() {
	// Create game server configuration
	config := &GameConfig{
		// LogicThreadCount: 8, // Commented out: LogicThread removed
		Host:       "0.0.0.0",
		Port:       8485, // MapleStory game port
		WorldName:  "Scania",
		MaxPlayers: 1000,
		ExpRate:    1, // 1x experience rate
		DropRate:   1, // 1x drop rate
		MesoRate:   1, // 1x meso rate
	}

	// Create game server
	gs, err := NewGameServer(config)
	if err != nil {
		log.Fatalf("Failed to create game server: %v", err)
	}

	// Start the game server
	if err := gs.Start(); err != nil {
		log.Fatalf("Failed to start game server: %v", err)
	}

	// Set up signal handling for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Wait for shutdown signal
	<-sigChan
	log.Println("Received shutdown signal, stopping game server...")

	// Stop the game server gracefully
	if err := gs.Stop(); err != nil {
		log.Printf("Error stopping game server: %v", err)
	}
}

// RunGameServerWithStats runs game server with statistics monitoring
func RunGameServerWithStats() {
	config := &GameConfig{
		// LogicThreadCount: 8, // Commented out: LogicThread removed
		Host:       "localhost",
		Port:       8485,
		WorldName:  "Scania",
		MaxPlayers: 1000,
		ExpRate:    2, // 2x experience rate
		DropRate:   2, // 2x drop rate
		MesoRate:   2, // 2x meso rate
	}

	gs, err := NewGameServer(config)
	if err != nil {
		log.Fatalf("Failed to create game server: %v", err)
	}

	if err := gs.Start(); err != nil {
		log.Fatalf("Failed to start game server: %v", err)
	}

	// Monitor server stats periodically
	go func() {
		for {
			stats := gs.GetStats()
			log.Printf("Game Server Stats: Players=%d/%d, World=%s, Rates: Exp=%dx, Drop=%dx, Meso=%dx",
				stats["current_players"],
				stats["max_players"],
				stats["world_name"],
				stats["exp_rate"],
				stats["drop_rate"],
				stats["meso_rate"])

			// Sleep for 10 seconds between stats
			time.Sleep(10 * time.Second)
		}
	}()

	// Wait for shutdown signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	log.Println("Shutting down game server...")
	gs.Stop()
}

// RunHighRateGameServer runs a high-rate game server configuration
func RunHighRateGameServer() {
	config := &GameConfig{
		// LogicThreadCount: 12, // Commented out: LogicThread removed
		Host:       "0.0.0.0",
		Port:       8485,
		WorldName:  "HighRate",
		MaxPlayers: 2000,
		ExpRate:    10, // 10x experience rate
		DropRate:   5,  // 5x drop rate
		MesoRate:   5,  // 5x meso rate
	}

	gs, err := NewGameServer(config)
	if err != nil {
		log.Fatalf("Failed to create high-rate game server: %v", err)
	}

	if err := gs.Start(); err != nil {
		log.Fatalf("Failed to start high-rate game server: %v", err)
	}

	// Monitor high-rate server stats
	go func() {
		for {
			stats := gs.GetStats()
			log.Printf("High-Rate Game Server Stats: Players=%d/%d, World=%s, Rates: Exp=%dx, Drop=%dx, Meso=%dx",
				stats["current_players"],
				stats["max_players"],
				stats["world_name"],
				stats["exp_rate"],
				stats["drop_rate"],
				stats["meso_rate"])

			time.Sleep(10 * time.Second)
		}
	}()

	// Wait for shutdown signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	log.Println("Shutting down high-rate game server...")
	gs.Stop()
}
