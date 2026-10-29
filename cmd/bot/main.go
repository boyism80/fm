package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/asynkron/protoactor-go/actor"
	botactor "github.com/boyism80/fm/bot/actor"
	"github.com/boyism80/fm/bot/cluster"
	"github.com/boyism80/fm/common/config"
	coreactor "github.com/boyism80/fm/core/actor"
	"github.com/boyism80/fm/core/luax"
	"github.com/boyism80/fm/maple/wz"
)

func main() {
	host, _ := os.Hostname()
	var (
		help    = flag.Bool("help", false, "Show help and example bot YAML layout")
		mode    = flag.String("mode", "integration", "Run mode: integration, master, agent or load")
		cfgPath = flag.String("config", "config/bot.yaml", "Path to bot YAML (default: config/bot.yaml)")
		filter  = flag.String("filter", "", "Run suites whose name contains this substring")
		seats   = flag.Int("seats", -1, "Override config seats when >= 0")
		junit   = flag.String("junit", "", "Write JUnit XML to this path")
		listen  = flag.String("listen", "127.0.0.1:8090", "Master: dashboard and agent address")
		runs    = flag.String("runs", "out/runs", "Master: run history directory")
		token   = flag.String("token", "", "Master and agent: agent token; without one the master takes only local agents")
		master  = flag.String("master", "ws://127.0.0.1:8090", "Agent: master address")
		name    = flag.String("name", host, "Agent: name shown on the dashboard")
	)
	flag.Parse()

	if *help {
		flag.Usage()
		log.Println("\nMapleStory Bot")
		log.Println("Runtime options are read from a bot YAML file.")
		log.Println("If the file is not in cwd, the binary directory and parents of cwd are searched.")
		log.Println("\nExample config/bot.yaml:")
		log.Println(strings.TrimSpace(`
login:
  host: "127.0.0.1"
  port: 8484
world_id: 0
channel: 0
seats: 3
password: "bot-pass"
timeout_ms: 10000
suite_timeout_ms: 180000
script_dir: "script/integration"
wz_path: "resources/wz"
`))
		log.Println("Notes:")
		log.Println("  - channel and world_id are 0-based. 0 is a valid channel.")
		log.Println("  - seats is the parallel suite limit and the instance slot count.")
		log.Println("  - -seats overrides seats when the flag is >= 0.")
		log.Println("  - run from bot: suites load script/integration relative to cwd.")
		log.Println("  - master and agent re-read a changed script on the next suite, so edits apply without a restart.")
		log.Println("\nExample usage (cwd bot):")
		log.Println("  go run ../cmd/bot -mode integration -config config/bot.yaml")
		log.Println("  go run ../cmd/bot -mode integration -filter npc_ -seats 3 -junit out/bot.xml")
		log.Println("  go run ../cmd/bot -mode master -listen 127.0.0.1:8090   (then open http://127.0.0.1:8090)")
		log.Println("  go run ../cmd/bot -mode agent -master ws://127.0.0.1:8090 -name agent-1 -seats 4")
		return
	}

	cfgFile, err := config.FindConfigFilePath(*cfgPath)
	if err != nil {
		log.Printf("Config: %v", err)
		os.Exit(2)
	}
	log.Printf("Using config file: %s", cfgFile)

	cfg, err := config.LoadBot(cfgFile)
	if err != nil {
		log.Printf("Config: %v", err)
		os.Exit(2)
	}
	if *seats >= 0 {
		cfg.Seats = *seats
	}
	if err := cfg.Validate(); err != nil {
		log.Printf("Config: %v", err)
		os.Exit(2)
	}

	switch *mode {
	case "integration":
		os.Exit(runIntegration(cfg, *filter, *junit))
	case "master":
		os.Exit(runMaster(cfg, *listen, *runs, *token))
	case "agent":
		os.Exit(runAgent(cfg, *master, *token, *name))
	case "load":
		// TODO(load): per-bot actors and a spawn rate.
		log.Println("load mode: 미구현")
		os.Exit(2)
	default:
		log.Printf("unknown mode %q (want integration, master, agent or load)", *mode)
		os.Exit(2)
	}
}

func runIntegration(cfg *config.Bot, filter, junit string) int {
	log.Printf("integration: login=%s world=%d channel=%d seats=%d timeout_ms=%d",
		fmt.Sprintf("%s:%d", cfg.Login.Host, cfg.Login.Port),
		cfg.WorldID,
		cfg.Channel,
		cfg.Seats,
		cfg.TimeoutMs,
	)

	resources := loadWz(cfg)
	if resources == nil {
		return 2
	}

	done := make(chan int, 1)
	system := coreactor.NewActorSystem()
	system.GetRoot().Spawn(actor.PropsFromProducer(func() actor.Actor {
		return botactor.NewRunnerActor(cfg, resources, filter, junit, done)
	}))
	return <-done
}

func runMaster(cfg *config.Bot, listen, runs, token string) int {
	luax.SetRecompileChanged(true)
	resources := loadWz(cfg)
	if resources == nil {
		return 2
	}

	m := cluster.NewMaster(cfg, resources, runs, token)
	log.Printf("master: dashboard http://%s, agents ws://%s/ws/agent", listen, listen)
	if err := m.ListenAndServe(listen); err != nil {
		log.Printf("master: %v", err)
		return 2
	}
	return 0
}

func runAgent(cfg *config.Bot, master, token, name string) int {
	luax.SetRecompileChanged(true)
	log.Printf("agent %s: login=%s:%d seats=%d", name, cfg.Login.Host, cfg.Login.Port, cfg.Seats)
	resources := loadWz(cfg)
	if resources == nil {
		return 2
	}

	system := coreactor.NewActorSystem()
	cluster.NewAgent(cfg, resources, system.GetRoot(), master, token, name).Run()
	return 0
}

func loadWz(cfg *config.Bot) *wz.Resources {
	log.Printf("loading wz %s", cfg.WzPath)
	resources := wz.NewResources(cfg.WzPath)
	if resources == nil {
		log.Printf("failed to load wz %s", cfg.WzPath)
	}
	return resources
}
