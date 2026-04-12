package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/boyism80/fm/login/server"
)

func main() {

	var (
		host = flag.String("host", "0.0.0.0", "Login server host address")
		port = flag.Int("port", 8484, "Login server port number")

		gameServerHost = flag.String("game-host", "localhost", "Game server host for redirection")
		gameServerPort = flag.Int("game-port", 8485, "Game server port for redirection")
		withStats      = flag.Bool("stats", false, "Enable statistics monitoring")
		help           = flag.Bool("help", false, "Show help information")
	)

	flag.Parse()

	if *help {
		flag.Usage()
		log.Println("\nMapleStory Login Server")
		log.Println("Handles user authentication and redirects to game server")
		log.Println("\nExample usage:")
		log.Println("  ./login-server")
		log.Println("  ./login-server -host=localhost -port=8484 -stats")

		return
	}

	config := &server.LoginConfig{

		Host:           *host,
		Port:           *port,
		GameServerHost: *gameServerHost,
		GameServerPort: *gameServerPort,
	}

	ls, err := server.NewLoginServer(config)
	if err != nil {
		log.Fatalf("Failed to create login server: %v", err)
	}

	if err := ls.Start(); err != nil {
		log.Fatalf("Failed to start login server: %v", err)
	}

	if *withStats {
		go func() {
			for {
				stats := ls.GetStats()
				log.Printf("Login Server Stats: Clients=%d, Listening=%v, Game Server=%s:%d",
					stats["client_count"],
					stats["listening"],
					stats["game_server_host"],
					stats["game_server_port"])

				time.Sleep(10 * time.Second)
			}
		}()
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	<-sigChan
	log.Println("Received shutdown signal, stopping login server...")

	if err := ls.Stop(); err != nil {
		log.Printf("Error stopping login server: %v", err)
	}

	log.Println("Login server stopped successfully")
}
