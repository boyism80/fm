package server

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/core"
	c_actor "github.com/boyism80/fm/core/actor"
	loginactor "github.com/boyism80/fm/login/actor"
	"github.com/boyism80/fm/login/client"
)

type LoginServer struct {
	server         *core.Server
	config         *LoginConfig
	packetHandlers *PacketHandlerRegistry
	context        *LoginServerContext
	actorSystem    *c_actor.ActorSystem
	actorRegistry  *c_actor.ActorRegistry
}

func (ls *LoginServer) GetServer() *core.Server {
	return ls.server
}

func (ls *LoginServer) GetServerContext() core.ServerContext {
	return ls.context
}

func (ls *LoginServer) handleClient(c core.Client) {
	loginClient, ok := c.(*client.LoginClient)
	if !ok {
		log.Printf("Client is not a LoginClient")
		return
	}

	props := actor.PropsFromProducer(func() actor.Actor {
		return &loginactor.LoginLogicActor{
			Client:  loginClient,
			Context: ls.context,
		}
	})

	pid := ls.actorRegistry.GetOrCreateActor(
		fmt.Sprintf("login_session_%d", loginClient.GetFd()),
		props,
	)

	loginClient.SetLogicActorPID(pid)
}

type LoginConfig struct {
	Host           string
	Port           int
	GameServerHost string
	GameServerPort int
}

func NewLoginServer(config *LoginConfig) (*LoginServer, error) {

	context := NewLoginServerContext()

	actorSystem := c_actor.NewActorSystem()
	actorRegistry := c_actor.NewActorRegistry(actorSystem)

	serverConfig := &core.ServerConfig{
		Host: config.Host,
		Port: config.Port,
		ClientFactory: func(conn net.Conn, clientID int) (core.Client, error) {
			return client.NewLoginClient(conn, clientID)
		},
		OnClientConnect: func(c core.Client) {

		},
	}

	server, err := core.NewServer(serverConfig)
	if err != nil {
		return nil, err
	}

	server.SetRootContext(actorSystem.GetRoot())

	ls := &LoginServer{
		server:         server,
		config:         config,
		packetHandlers: NewPacketHandlerRegistry(nil),
		context:        context,
		actorSystem:    actorSystem,
		actorRegistry:  actorRegistry,
	}
	ls.packetHandlers.ls = ls

	serverConfig.OnClientConnect = ls.handleClient

	ls.registerPacketHandlers()

	return ls, nil
}

func (ls *LoginServer) Start() error {
	log.Println("Starting MapleStory Login Server...")

	if err := ls.server.Start(ls.config.Host, ls.config.Port); err != nil {
		return err
	}

	log.Printf("Login server started on %s:%d", ls.config.Host, ls.config.Port)
	log.Printf("Game server redirect: %s:%d", ls.config.GameServerHost, ls.config.GameServerPort)

	return nil
}

func (ls *LoginServer) Stop() error {
	log.Println("Shutting down login server...")
	return ls.server.Stop()
}

func (ls *LoginServer) GetStats() map[string]interface{} {
	stats := ls.server.GetStats()

	stats["server_type"] = "login"
	stats["game_server_host"] = ls.config.GameServerHost
	stats["game_server_port"] = ls.config.GameServerPort

	return stats
}

func RunLoginServer() {

	config := &LoginConfig{

		Host:           "0.0.0.0",
		Port:           8484,
		GameServerHost: "localhost",
		GameServerPort: 8485,
	}

	ls, err := NewLoginServer(config)
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
}

func RunLoginServerWithStats() {
	config := &LoginConfig{

		Host:           "localhost",
		Port:           8484,
		GameServerHost: "localhost",
		GameServerPort: 8485,
	}

	ls, err := NewLoginServer(config)
	if err != nil {
		log.Fatalf("Failed to create login server: %v", err)
	}

	if err := ls.Start(); err != nil {
		log.Fatalf("Failed to start login server: %v", err)
	}

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

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	log.Println("Shutting down login server...")
	ls.Stop()
}
