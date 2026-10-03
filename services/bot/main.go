package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/boyism80/fm/common/config"
)

func main() {
	var (
		help    = flag.Bool("help", false, "Show help and example bot YAML layout")
		mode    = flag.String("mode", "integration", "Run mode: integration or load")
		cfgPath = flag.String("config", "config/bot.yaml", "Path to bot YAML (default: config/bot.yaml)")
		filter  = flag.String("filter", "", "Run suites whose name contains this substring")
		seats   = flag.Int("seats", -1, "Override config seats when >= 0")
		junit   = flag.String("junit", "", "Write JUnit XML to this path")
	)
	flag.Parse()

	if *help {
		flag.Usage()
		log.Println("\nMapleStory Bot")
		log.Println("Runtime options are read from a bot YAML file. CLI flags: -help, -mode, -config, -filter, -seats, -junit.")
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
		log.Println("\nExample usage:")
		log.Println("  go run ./services/bot -mode integration -config config/bot.yaml")
		log.Println("  go run ./services/bot -mode integration -filter npc_ -seats 3 -junit out/bot.xml")
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
	case "load":
		// TODO(load): per-bot actors and a spawn rate.
		log.Println("load mode: 미구현")
		os.Exit(2)
	default:
		log.Printf("unknown mode %q (want integration or load)", *mode)
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
	if filter != "" {
		log.Printf("filter: %s", filter)
	}
	if junit != "" {
		log.Printf("junit: %s", junit)
	}
	log.Println("integration mode: 미구현")
	return 2
}
