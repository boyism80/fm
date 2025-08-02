package server

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/game/data"
	"github.com/boyism80/fm/game/entity"
)

// GameClientData represents game-specific client data
type GameClientData struct {
	Character *entity.Character // Character data for the connected player
}

// GameServer represents the game server for MapleStory private server
type GameServer struct {
	server    *core.Server[GameClientData]
	config    *GameConfig
	resources *data.Resources        // Game data resources
	maps      map[uint32]*entity.Map // Map instances by map ID
	mapsMutex sync.RWMutex
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
	// Create core server configuration
	serverConfig := &core.ServerConfig{
		LogicThreadCount: config.LogicThreadCount,
		Host:             config.Host,
		Port:             config.Port,
	}

	// Create core server
	server, err := core.NewServer[GameClientData](serverConfig)
	if err != nil {
		return nil, err
	}

	gameServer := &GameServer{
		server: server,
		config: config,
		maps:   make(map[uint32]*entity.Map),
	}

	// Load game resources
	log.Println("Loading game resources...")
	resources := data.NewResources()
	if resources == nil {
		return nil, fmt.Errorf("failed to load game resources")
	}
	gameServer.resources = resources
	log.Println("Game resources loaded successfully")

	// Set up map listeners for packet broadcasting
	log.Println("Setting up map listeners...")
	gameMapListener := NewGameMapListener(gameServer)

	// Pre-create all map instances from loaded resources
	log.Println("Pre-creating map instances...")
	for mapID, _ := range resources.Maps {
		mapInstance := entity.NewMap(mapID, gameMapListener)
		// TODO: Initialize map with MapSpec data (mob spawns, etc.)
		gameServer.maps[mapID] = mapInstance
	}
	log.Printf("Pre-created %d map instances", len(gameServer.maps))
	log.Println("Map listeners configured")

	// Register game server packet handlers
	gameServer.registerPacketHandlers()

	// Set up client disconnect callback
	gameServer.server.SetOnClientDisconnect(func(client interface{}) {
		if gameClient, ok := client.(*core.Client[GameClientData]); ok {
			gameServer.handleClientDisconnect(gameClient)
		}
	})

	return gameServer, nil
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
func (gs *GameServer) handleClientDisconnect(client *core.Client[GameClientData]) {
	clientData := client.GetData()
	if clientData.Character == nil {
		return
	}
	character := clientData.Character
	mapInstance := gs.GetMap(character.Map)
	if mapInstance == nil {
		return
	}
	mapInstance.RemovePlayer(character.ID)
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
