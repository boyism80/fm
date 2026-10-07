package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/boyism80/fm/common/config"
	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/services/cashshop/server"
)

func main() {
	cfgPath := flag.String("config", "config/dev/cashshop-0.yaml", "Path to cash shop YAML")
	flag.Parse()

	cfgFile, err := config.FindConfigFilePath(*cfgPath)
	if err != nil {
		log.Fatalf("Config: %v", err)
	}
	cfg, err := config.LoadCashShop(cfgFile)
	if err != nil {
		log.Fatalf("Config: %v", err)
	}
	if cfg.Internal.TimeoutSeconds > 0 {
		core.InternalRPCPerStepTimeout = time.Duration(cfg.Internal.TimeoutSeconds) * time.Second
	} else {
		core.InternalRPCPerStepTimeout = 10 * time.Second
	}

	cs, err := server.NewCashShopServer(cfg)
	if err != nil {
		log.Fatalf("Failed to create cash shop server: %v", err)
	}
	if err := cs.Start(); err != nil {
		log.Fatalf("Failed to start cash shop server: %v", err)
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	if err := cs.Stop(); err != nil {
		log.Printf("Error stopping cash shop server: %v", err)
	}
}
