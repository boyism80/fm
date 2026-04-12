package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/boyism80/fm/game/server"
)

func main() {

	var (
		host       = flag.String("host", "0.0.0.0", "Game server host address")
		port       = flag.Int("port", 8485, "Game server port number")
		wzPath     = flag.String("wz-path", "resources/wz", "Path to WZ files directory")
		worldName  = flag.String("world", "Scania", "World/Channel name")
		maxPlayers = flag.Int("max-players", 1000, "Maximum number of players per world")
		expRate    = flag.Int("exp-rate", 100, "Experience rate multiplier")
		dropRate   = flag.Int("drop-rate", 10, "Drop rate multiplier")
		mesoRate   = flag.Int("meso-rate", 1, "Meso rate multiplier")
		withStats  = flag.Bool("stats", false, "Enable statistics monitoring")
		highRate   = flag.Bool("high-rate", false, "Enable high-rate server configuration")
		help       = flag.Bool("help", false, "Show help information")
	)

	flag.Parse()

	if *help {
		flag.Usage()
		log.Println("\nMapleStory Game Server")
		log.Println("Handles game mechanics, character movement, combat, inventory")
		log.Println("\nExample usage:")
		log.Println("  ./game-server")
		log.Println("  ./game-server -host=localhost -port=8485 -stats")
		log.Println("  ./game-server -world=Scania -exp-rate=2 -drop-rate=2 -meso-rate=2")
		log.Println("  ./game-server -high-rate -max-players=2000")
		return
	}

	if *highRate {

		*worldName = "HighRate"
		*maxPlayers = 2000
		*expRate = 10
		*dropRate = 5
		*mesoRate = 5
		log.Println("High-rate server configuration enabled")
	}

	config := &server.GameConfig{

		Host:       *host,
		Port:       *port,
		WzPath:     *wzPath,
		WorldName:  *worldName,
		MaxPlayers: *maxPlayers,
		ExpRate:    *expRate,
		DropRate:   *dropRate,
		MesoRate:   *mesoRate,
	}

	gs, err := server.NewGameServer(config)
	if err != nil {
		log.Fatalf("Failed to create game server: %v", err)
	}

	if err := gs.Start(); err != nil {
		log.Fatalf("Failed to start game server: %v", err)
	}

	if *withStats {
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

				time.Sleep(10 * time.Second)
			}
		}()
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	<-sigChan
	log.Println("Received shutdown signal, stopping game server...")

	if err := gs.Stop(); err != nil {
		log.Printf("Error stopping game server: %v", err)
	}

	log.Println("Game server stopped successfully")
}
