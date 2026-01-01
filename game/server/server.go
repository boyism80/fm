package server

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"path/filepath"
	"sync"
	"syscall"
	"time"

	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/core/luax"
	"github.com/boyism80/fm/game/client"
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
}

func (gs *GameServer) GetServer() *core.Server {
	return gs.server
}

// GameContext provides access to game resources and services
type GameContext interface {
	GetResources() *wz.Resources
	GetMap(mapId uint32) *entity.Map
	GetLogicThread() *core.LogicThread // Returns core.LogicThread
}

// GameConfig holds game server specific configuration
type GameConfig struct {
	LogicThreadCount int    // Number of logic threads for game processing
	Host             string // Game server host address
	Port             int    // Game server port number
	WzPath           string // Path to WZ files directory
	WorldName        string // World/Channel name
	MaxPlayers       int    // Maximum number of players per world
	ExpRate          int    // Experience rate multiplier
	DropRate         int    // Drop rate multiplier
	MesoRate         int    // Meso rate multiplier
}

// GetLuaState returns the Lua state for script execution

// ExecuteNpcScript executes the Lua script for an NPC
func (gs *GameServer) ExecuteNpcScript(character *entity.Character, npcInterface interface{}) error {
	// Get NPC model to determine script file
	npc, ok := npcInterface.(*entity.Npc)
	if !ok {
		return fmt.Errorf("invalid NPC type")
	}

	if npc.Wz == nil {
		return fmt.Errorf("NPC has no model")
	}

	// Get Lua state from LogicThread
	logicThread := gs.GetLogicThread()
	if logicThread == nil {
		return fmt.Errorf("logic thread not available")
	}

	luaState := logicThread.GetLuaState()
	if luaState == nil {
		return fmt.Errorf("lua state not available")
	}

	// Load NPC script
	path := filepath.Join("script", "npc", fmt.Sprintf("%d.lua", npc.Wz.ID))

	// Load script function
	fn, err := luaState.LoadFile(path)
	if err != nil {
		log.Printf("Failed to load NPC script %s: %v", path, err)
		return fmt.Errorf("failed to load NPC script: %w", err)
	}

	// Create new thread for script execution
	co, _ := luaState.NewThread()

	// Push script function to thread
	co.Push(fn)

	// Execute script (loads all functions)
	if err := co.PCall(0, lua.MultRet, nil); err != nil {
		return fmt.Errorf("failed to execute NPC script: %w", err)
	}

	// Get on_start function
	onStartFn := co.GetGlobal("on_start")
	if onStartFn.Type() != lua.LTFunction {
		return fmt.Errorf("on_start function not found in NPC script")
	}

	// Create character Lua object using luax.NewLuable
	characterLua := luax.NewLuable(co, character)

	// Call on_start(me) function
	resumeState, err, _ := luaState.Resume(co, onStartFn.(*lua.LFunction), characterLua)
	if err != nil {
		return fmt.Errorf("failed to call on_start: %w", err)
	}

	// If script yielded (waiting for dialog response), store the coroutine
	if resumeState == lua.ResumeYield {
		character.SetCurrentDialog(co)
	}

	return nil
}

// NewGameServer creates a new game server with specified configuration
func NewGameServer(config *GameConfig) (*GameServer, error) {
	serverConfig := &core.ServerConfig{
		LogicThreadCount: config.LogicThreadCount,
		Host:             config.Host,
		Port:             config.Port,
		ClientFactory: func(conn net.Conn, clientID int) (core.Client, error) {
			return client.NewGameClient(conn, clientID)
		},
		LogicThreadInit: func(thread *core.LogicThread) {
			// Initialize Lua state for game server logic thread
			luaState := thread.GetLuaState()
			if luaState == nil {
				log.Printf("Failed to get Lua state for logic thread %d", thread.GetID())
				return
			}

			// Register Lua types with inheritance
			luax.RegisterLuaType[*entity.Object](luaState)
			luax.RegisterLuaDerivedType[*entity.Life, *entity.Object](luaState)
			luax.RegisterLuaDerivedType[*entity.Character, *entity.Life](luaState)
			luax.RegisterLuaDerivedType[*entity.Mob, *entity.Life](luaState)

			// Register utility functions
			luax.RegisterFunc(luaState, "sleep", func(L *lua.LState) int {
				duration := L.CheckNumber(1)
				time.Sleep(time.Duration(float64(duration) * float64(time.Second)))
				return 0
			})

			log.Printf("Lua state initialized for logic thread %d", thread.GetID())
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

	gameServer := &GameServer{
		server:    server,
		config:    config,
		resources: resources,
		maps:      make(map[uint32]*entity.Map),
	}

	// Initialize command handler
	gameServer.commandHandler = NewCommandHandler(gameServer)

	// Initialize packet handler registry
	gameServer.packetHandlers = NewPacketHandlerRegistry(gameServer)

	// Pre-create all maps
	gameServer.preCreateMaps()

	// Register packet handlers
	gameServer.registerPacketHandlers()

	// Register command handlers
	gameServer.registerCommandHandlers()

	// Set client disconnect handler
	server.SetOnClientDisconnect(func(client interface{}) {
		if client, ok := client.(core.Client); ok {
			gameServer.handleClientDisconnect(client)
		}
	})

	return gameServer, nil
}

// GetResources returns the game resources
func (gs *GameServer) GetResources() *wz.Resources {
	return gs.resources
}

// GetThreadHash returns a hash for thread assignment
func (gs *GameServer) GetThreadHash() int {
	return 0 // GameServer uses a single thread hash
}

// GetLogicThread returns the logic thread for scheduling tasks
func (gs *GameServer) GetLogicThread() *core.LogicThread {
	logicThread, err := gs.server.GetLogicThread(gs)
	if err != nil {
		return nil
	}
	return logicThread
}

// preCreateMaps pre-creates all map instances from loaded resources
func (gs *GameServer) preCreateMaps() {
	log.Println("Setting up map listeners...")
	gameMapListener := NewGameMapListener(gs)

	log.Println("Pre-creating map instances...")
	for mapID := range gs.resources.Maps {
		mapInstance := entity.NewMap(mapID, gameMapListener, mapID, gs) // Pass mapID instead of mapSpec
		// TODO: Initialize map with MapSpec data (mob spawns, etc.)
		gs.maps[mapID] = mapInstance
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
	mapID := character.GetMap()
	mapInstance := gs.GetMap(mapID)
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
		LogicThreadCount: 8, // 8 logic threads for game processing
		Host:             "0.0.0.0",
		Port:             8485, // MapleStory game port
		WorldName:        "Scania",
		MaxPlayers:       1000,
		ExpRate:          1, // 1x experience rate
		DropRate:         1, // 1x drop rate
		MesoRate:         1, // 1x meso rate
	}

	// Create game server
	gameServer, err := NewGameServer(config)
	if err != nil {
		log.Fatalf("Failed to create game server: %v", err)
	}

	// Start the game server
	if err := gameServer.Start(); err != nil {
		log.Fatalf("Failed to start game server: %v", err)
	}

	// Set up signal handling for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Wait for shutdown signal
	<-sigChan
	log.Println("Received shutdown signal, stopping game server...")

	// Stop the game server gracefully
	if err := gameServer.Stop(); err != nil {
		log.Printf("Error stopping game server: %v", err)
	}
}

// RunGameServerWithStats runs game server with statistics monitoring
func RunGameServerWithStats() {
	config := &GameConfig{
		LogicThreadCount: 8,
		Host:             "localhost",
		Port:             8485,
		WorldName:        "Scania",
		MaxPlayers:       1000,
		ExpRate:          2, // 2x experience rate
		DropRate:         2, // 2x drop rate
		MesoRate:         2, // 2x meso rate
	}

	gameServer, err := NewGameServer(config)
	if err != nil {
		log.Fatalf("Failed to create game server: %v", err)
	}

	if err := gameServer.Start(); err != nil {
		log.Fatalf("Failed to start game server: %v", err)
	}

	// Monitor server stats periodically
	go func() {
		for {
			stats := gameServer.GetStats()
			log.Printf("Game Server Stats: IO Threads=%d, Logic Threads=%d, Players=%d/%d, World=%s, Rates: Exp=%dx, Drop=%dx, Meso=%dx",
				stats["io_thread_count"],
				stats["logic_thread_count"],
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
	gameServer.Stop()
}

// RunHighRateGameServer runs a high-rate game server configuration
func RunHighRateGameServer() {
	config := &GameConfig{
		LogicThreadCount: 12, // More logic threads for complex game mechanics
		Host:             "0.0.0.0",
		Port:             8485,
		WorldName:        "HighRate",
		MaxPlayers:       2000,
		ExpRate:          10, // 10x experience rate
		DropRate:         5,  // 5x drop rate
		MesoRate:         5,  // 5x meso rate
	}

	gameServer, err := NewGameServer(config)
	if err != nil {
		log.Fatalf("Failed to create high-rate game server: %v", err)
	}

	if err := gameServer.Start(); err != nil {
		log.Fatalf("Failed to start high-rate game server: %v", err)
	}

	// Monitor high-rate server stats
	go func() {
		for {
			stats := gameServer.GetStats()
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
	gameServer.Stop()
}
