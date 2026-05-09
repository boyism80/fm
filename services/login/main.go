package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/boyism80/fm/common/config"
	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/services/login/server"
)

func main() {
	var (
		help    = flag.Bool("help", false, "Show help and example login YAML layout")
		cfgPath = flag.String("config", "config/login.yaml", "Path to login-only YAML (default: config/login.yaml)")
	)
	flag.Parse()

	if *help {
		flag.Usage()
		log.Println("\nMapleStory Login Server")
		log.Println("Runtime options are read from a login-only YAML file. CLI flags: -help, -config (default config/login.yaml).")
		log.Println("If the file is not in cwd, the binary directory and parents of cwd are searched.")
		log.Println("\nExample config/login.yaml (root document = login fields only):")
		log.Println(strings.TrimSpace(`
host: "0.0.0.0"
port: 8484
initial_role: 1
login_instance_id: "login-1"
catalog_retry_interval_seconds: 2
catalog_retry_max_attempts: 30
internal:
  host: "127.0.0.1"
  port: 50051
  timeout: 10
  heartbeat_interval_seconds: 15
`))
		log.Println("Notes:")
		log.Println("  - internal.heartbeat_interval_seconds: periodic Ping to internal (seconds); omit or 0 to disable heartbeat loop.")
		log.Println("  - login_instance_id: unique id per login process (internal periodic Ping / Redis alive).")
		log.Println("  - internal is required; login server fetches world/channel routing from internal at startup.")
		log.Println("\nExample usage:")
		log.Println("  ./login-server")
		log.Println("  ./login-server -config=config/login.yaml")
		return
	}

	cfgFile, err := config.ResolvePath(*cfgPath)
	if err != nil {
		log.Fatalf("Config: %v", err)
	}
	log.Printf("Using config file: %s", cfgFile)

	l, err := config.LoadLogin(cfgFile)
	if err != nil {
		log.Fatalf("Config: %v", err)
	}
	if l.Internal.TimeoutSeconds > 0 {
		core.InternalRPCPerStepTimeout = time.Duration(l.Internal.TimeoutSeconds) * time.Second
	} else {
		core.InternalRPCPerStepTimeout = 10 * time.Second
	}

	srvCfg := &server.LoginConfig{
		Host:                             l.Host,
		Port:                             l.Port,
		InitialRole:                      uint32(l.InitialRole),
		LoginInstanceID:                  l.LoginInstanceID,
		InternalHeartbeatIntervalSeconds: l.Internal.HeartbeatIntervalSeconds,
		InternalHost:                     l.Internal.Host,
		InternalPort:                     l.Internal.Port,
		CatalogRetryIntervalSeconds:      l.CatalogRetryIntervalSeconds,
		CatalogRetryMaxAttempts:          l.CatalogRetryMaxAttempts,
	}

	ls, err := server.NewLoginServer(srvCfg)
	if err != nil {
		log.Fatalf("Failed to create login server: %v", err)
	}

	if err := ls.Start(); err != nil {
		log.Fatalf("Failed to start login server: %v", err)
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
