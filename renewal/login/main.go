package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/boyism80/fm/renewal/login/server"
)

// main function for MapleStory Login Server
func main() {
	// Define command line flags
	var (
		host           = flag.String("host", "0.0.0.0", "Login server host address")
		port           = flag.Int("port", 8484, "Login server port number")
		logicThreads   = flag.Int("logic-threads", 4, "Number of logic threads")
		gameServerHost = flag.String("game-host", "localhost", "Game server host for redirection")
		gameServerPort = flag.Int("game-port", 8485, "Game server port for redirection")
		withStats      = flag.Bool("stats", false, "Enable statistics monitoring")
		help           = flag.Bool("help", false, "Show help information")
	)

	// Parse command line arguments
	flag.Parse()

	// Show help if requested
	if *help {
		flag.Usage()
		log.Println("\nMapleStory Login Server")
		log.Println("Handles user authentication and redirects to game server")
		log.Println("\nExample usage:")
		log.Println("  ./login-server")
		log.Println("  ./login-server -host=localhost -port=8484 -stats")
		log.Println("  ./login-server -logic-threads=8 -game-host=192.168.1.100")
		return
	}

	// Create login server configuration
	config := &server.LoginConfig{
		LogicThreadCount: *logicThreads,
		Host:             *host,
		Port:             *port,
		GameServerHost:   *gameServerHost,
		GameServerPort:   *gameServerPort,
	}

	// Create login server
	loginServer, err := server.NewLoginServer(config)
	if err != nil {
		log.Fatalf("Failed to create login server: %v", err)
	}

	// Start the login server
	if err := loginServer.Start(); err != nil {
		log.Fatalf("Failed to start login server: %v", err)
	}

	// Start statistics monitoring if enabled
	if *withStats {
		go func() {
			for {
				stats := loginServer.GetStats()
				log.Printf("Login Server Stats: Logic Threads=%d, Clients=%d, Listening=%v, Game Server=%s:%d",
					stats["logic_thread_count"],
					stats["client_count"],
					stats["listening"],
					stats["game_server_host"],
					stats["game_server_port"])

				// Sleep for 10 seconds between stats
				time.Sleep(10 * time.Second)
			}
		}()
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

	log.Println("Login server stopped successfully")
}
