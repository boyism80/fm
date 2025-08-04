package server

import (
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/login/client"
)

// LoginServer represents the login server for MapleStory private server
type LoginServer struct {
	server *core.Server
	config *LoginConfig
}

// LoginConfig holds login server specific configuration
type LoginConfig struct {
	LogicThreadCount int    // Number of logic threads for login processing
	Host             string // Login server host address
	Port             int    // Login server port number
	GameServerHost   string // Game server host for redirection
	GameServerPort   int    // Game server port for redirection
}

// NewLoginServer creates a new login server with specified configuration
func NewLoginServer(config *LoginConfig) (*LoginServer, error) {
	// Create core server configuration
	serverConfig := &core.ServerConfig{
		LogicThreadCount: config.LogicThreadCount,
		Host:             config.Host,
		Port:             config.Port,
		ClientFactory: func(conn net.Conn, clientID int) (core.Client, error) {
			return client.NewLoginClient(conn, clientID)
		},
		LogicThreadInit: func(thread *core.LogicThread) {
			// Login server doesn't need Lua initialization
			log.Printf("Login server logic thread %d initialized (no Lua required)", thread.GetID())
		},
	}

	// Create core server
	server, err := core.NewServer(serverConfig)
	if err != nil {
		return nil, err
	}

	loginServer := &LoginServer{
		server: server,
		config: config,
	}

	// Register packet handlers
	loginServer.registerPacketHandlers()

	return loginServer, nil
}

// Start initializes and starts the login server
func (ls *LoginServer) Start() error {
	log.Println("Starting MapleStory Login Server...")

	// Register login server packet handlers
	ls.registerPacketHandlers()

	// Start the core server
	if err := ls.server.Start(ls.config.Host, ls.config.Port); err != nil {
		return err
	}

	log.Printf("Login server started on %s:%d", ls.config.Host, ls.config.Port)
	log.Printf("Game server redirect: %s:%d", ls.config.GameServerHost, ls.config.GameServerPort)

	return nil
}

// Stop gracefully shuts down the login server
func (ls *LoginServer) Stop() error {
	log.Println("Shutting down login server...")
	return ls.server.Stop()
}

// GetStats returns login server statistics for monitoring
func (ls *LoginServer) GetStats() map[string]interface{} {
	stats := ls.server.GetStats()

	// Add login-specific stats
	stats["server_type"] = "login"
	stats["game_server_host"] = ls.config.GameServerHost
	stats["game_server_port"] = ls.config.GameServerPort

	return stats
}

// RunLoginServer demonstrates how to run the login server with signal handling
func RunLoginServer() {
	// Create login server configuration
	config := &LoginConfig{
		LogicThreadCount: 4, // 4 logic threads for login processing
		Host:             "0.0.0.0",
		Port:             8484, // MapleStory login port
		GameServerHost:   "localhost",
		GameServerPort:   8485, // Game server port
	}

	// Create login server
	loginServer, err := NewLoginServer(config)
	if err != nil {
		log.Fatalf("Failed to create login server: %v", err)
	}

	// Start the login server
	if err := loginServer.Start(); err != nil {
		log.Fatalf("Failed to start login server: %v", err)
	}

	// Set up signal handling for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Wait for shutdown signal
	<-sigChan
	log.Println("Received shutdown signal, stopping login server...")

	// Stop the login server gracefully
	if err := loginServer.Stop(); err != nil {
		log.Printf("Error stopping login server: %v", err)
	}
}

// RunLoginServerWithStats demonstrates login server with statistics monitoring
func RunLoginServerWithStats() {
	config := &LoginConfig{
		LogicThreadCount: 4,
		Host:             "localhost",
		Port:             8484,
		GameServerHost:   "localhost",
		GameServerPort:   8485,
	}

	loginServer, err := NewLoginServer(config)
	if err != nil {
		log.Fatalf("Failed to create login server: %v", err)
	}

	if err := loginServer.Start(); err != nil {
		log.Fatalf("Failed to start login server: %v", err)
	}

	// Monitor server stats periodically
	go func() {
		for {
			stats := loginServer.GetStats()
			log.Printf("Login Server Stats: IO Threads=%d, Logic Threads=%d, Clients=%d, Listening=%v, Game Server=%s:%d",
				stats["io_thread_count"],
				stats["logic_thread_count"],
				stats["client_count"],
				stats["listening"],
				stats["game_server_host"],
				stats["game_server_port"])

			// Sleep for 10 seconds between stats
			time.Sleep(10 * time.Second)
		}
	}()

	// Wait for shutdown signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	log.Println("Shutting down login server...")
	loginServer.Stop()
}
