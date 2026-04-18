package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/boyism80/fm/common/config"
	"github.com/boyism80/fm/services/game/server"
)

func main() {
	var (
		help    = flag.Bool("help", false, "Show help and example game YAML layout")
		cfgPath = flag.String("config", "config/game.yaml", "Path to game-only YAML (default: config/game.yaml)")
	)
	flag.Parse()

	if *help {
		flag.Usage()
		log.Println("\nMapleStory Game Server")
		log.Println("Runtime options are read from a game-only YAML file. CLI flags: -help, -config (default config/game.yaml).")
		log.Println("If the file is not in cwd, the binary directory and parents of cwd are searched (see common/config.ResolvePath).")
		log.Println("\nExample config/game.yaml (root document = game fields only):")
		log.Println(strings.TrimSpace(`
host: "0.0.0.0"
port: 8485
channel_id: 0
wz_path: "resources/wz"
world: "Scania"
max_players: 1000
rate:
  exp: 100
  drop: 10
  meso: 1
internal:
  host: "127.0.0.1"
  port: 50051
rabbitmq:
  ip: "127.0.0.1"
  port: 5672
  uid: "guest"
  pwd: "guest"
  vhost: "/fm"
high_rate: false
`))
		log.Println("Notes:")
		log.Println("  - internal: omit host or set port: 0 to skip the pre-start internal gRPC ping.")
		log.Println("  - high_rate: when true, overrides world, max_players, and rate (exp/drop/meso).")
		log.Println("\nExample usage:")
		log.Println("  ./game-server")
		log.Println("  ./game-server -config=config/game.yaml")
		return
	}

	cfgFile, err := config.ResolvePath(*cfgPath)
	if err != nil {
		log.Fatalf("Config: %v", err)
	}
	log.Printf("Using config file: %s", cfgFile)

	g, err := config.LoadGame(cfgFile)
	if err != nil {
		log.Fatalf("Config: %v", err)
	}

	if g.HighRate {
		g.World = "HighRate"
		g.MaxPlayers = 2000
		g.Rate.Exp = 10
		g.Rate.Drop = 5
		g.Rate.Meso = 5
		log.Println("High-rate server configuration enabled (from config)")
	}

	internalAddr := g.Internal.GRPCAddr()
	if internalAddr != "" {
		pingCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		err := server.PingInternalService(pingCtx, internalAddr)
		cancel()
		if err != nil {
			log.Fatalf("Internal server not reachable (%s): %v", internalAddr, err)
		}
		log.Printf("Internal server ping ok (%s)", internalAddr)
	}

	srvCfg := &server.GameConfig{
		Host:         g.Host,
		Port:         g.Port,
		ChannelId:    uint32(g.ChannelId),
		WzPath:       g.WzPath,
		WorldName:    g.World,
		WorldId:      uint32(g.WorldId),
		MaxPlayers:   g.MaxPlayers,
		ExpRate:      g.Rate.Exp,
		DropRate:     g.Rate.Drop,
		MesoRate:     g.Rate.Meso,
		InternalAddr: g.Internal.GRPCAddr(),
		RabbitMQ:     g.RabbitMQ,
	}

	gs, err := server.NewGameServer(srvCfg)
	if err != nil {
		log.Fatalf("Failed to create game server: %v", err)
	}

	if err := gs.Start(); err != nil {
		log.Fatalf("Failed to start game server: %v", err)
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
