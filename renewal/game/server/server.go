package server

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/game/data"
	"github.com/boyism80/fm/game/entity"
	"github.com/boyism80/fm/renewal/game/client"
)

// logicThreadWrapper wraps core.LogicThread to match entity.LogicThread interface
type logicThreadWrapper struct {
	logicThread *core.LogicThread
}

func (w *logicThreadWrapper) Schedule(delay time.Duration, task func()) interface{} {
	return w.logicThread.Schedule(delay, func() error {
		task()
		return nil
	}, nil)
}

func (w *logicThreadWrapper) ScheduleAtFixedRate(initialDelay, period time.Duration, task func()) interface{} {
	// Not implemented in core.LogicThread, return nil
	return nil
}

func (w *logicThreadWrapper) ScheduleWithFixedDelay(initialDelay, delay time.Duration, task func()) interface{} {
	// Not implemented in core.LogicThread, return nil
	return nil
}

// GameServer represents the game server for MapleStory private server
type GameServer struct {
	server         *core.Server
	config         *GameConfig
	resources      *data.Resources        // Game data resources
	maps           map[uint32]*entity.Map // Map instances by map ID
	mapsMutex      sync.RWMutex
	commandHandler *CommandHandler
}

// GameConfig holds game server specific configuration
type GameConfig struct {
	LogicThreadCount int    // Number of logic threads for game processing
	Host             string // Game server host address
	Port             int    // Game server port number
	WorldName        string // World/Channel name
	MaxPlayers       int    // Maximum number of players per world
	ExpRate          int    // Experience rate multiplier
	DropRate         int    // Drop rate multiplier
	MesoRate         int    // Meso rate multiplier
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
	}
	server, err := core.NewServer(serverConfig)
	if err != nil {
		return nil, err
	}

	// Load game resources
	log.Println("Loading game resources...")
	resources := data.NewResources()
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

	// Pre-create all maps
	gameServer.preCreateMaps()

	// Register packet handlers
	gameServer.registerPacketHandlers()

	// Set client disconnect handler
	server.SetOnClientDisconnect(func(client interface{}) {
		if client, ok := client.(core.Client); ok {
			gameServer.handleClientDisconnect(client)
		}
	})

	return gameServer, nil
}

// GetResources returns the game resources
func (gs *GameServer) GetResources() *data.Resources {
	return gs.resources
}

// GetThreadHash returns a hash for thread assignment
func (gs *GameServer) GetThreadHash() int {
	return 0 // GameServer uses a single thread hash
}

// GetLogicThread returns the logic thread for scheduling tasks
func (gs *GameServer) GetLogicThread() entity.LogicThread {
	logicThread, err := gs.server.GetLogicThread(gs)
	if err != nil {
		return nil
	}
	return &logicThreadWrapper{logicThread: logicThread}
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
		log.Printf("Resources loaded: %d maps, %d monsters, %d items, %d drops, %d strings",
			len(gs.resources.Maps), len(gs.resources.Monsters),
			len(gs.resources.Items), len(gs.resources.Drops), len(gs.resources.Strings))
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
		stats["strings_loaded"] = len(gs.resources.Strings)
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
