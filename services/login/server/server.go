package server

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/core"
	c_actor "github.com/boyism80/fm/core/actor"
	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
	loginactor "github.com/boyism80/fm/services/login/actor"
	"github.com/boyism80/fm/services/login/client"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type LoginServer struct {
	server         *core.Server
	config         *LoginConfig
	packetHandlers *PacketHandlerRegistry
	context        *LoginServerContext
	actorSystem    *c_actor.ActorSystem
	actorRegistry  *c_actor.ActorRegistry
	worldCatalog   []*internal.WorldCatalog
	channelRoutes  map[uint32]map[uint32]*internal.ChannelCatalog
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
	Host         string
	Port         int
	WorldId      uint32
	InitialRole  uint32
	InternalHost string
	InternalPort int
}

func (c *LoginConfig) internalAddr() string {
	if c.InternalHost == "" || c.InternalPort == 0 {
		return ""
	}
	return fmt.Sprintf("%s:%d", c.InternalHost, c.InternalPort)
}

func NewLoginServer(config *LoginConfig) (*LoginServer, error) {
	if config.internalAddr() == "" {
		return nil, fmt.Errorf("login server requires internal gRPC endpoint")
	}
	var internalClient internal.InternalClient
	if addr := config.internalAddr(); addr != "" {
		conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			return nil, fmt.Errorf("internal gRPC dial %s: %w", addr, err)
		}
		internalClient = internal.NewInternalClient(conn)
		log.Printf("Internal gRPC client connected to %s", addr)
	}

	context := NewLoginServerContext(internalClient)

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
		channelRoutes:  make(map[uint32]map[uint32]*internal.ChannelCatalog),
	}
	ls.packetHandlers.ls = ls

	serverConfig.OnClientConnect = ls.handleClient
	server.SetOnClientDisconnect(ls.handleClientDisconnect)

	ls.registerPacketHandlers()
	if err := ls.loadServerCatalog(); err != nil {
		return nil, err
	}

	return ls, nil
}

func (ls *LoginServer) loadServerCatalog() error {
	if ls.context.InternalClient == nil {
		return fmt.Errorf("internal gRPC client is required")
	}
	ctx, cancel := context.WithTimeout(context.Background(), core.InternalRPCPerStepTimeout)
	defer cancel()
	reply, err := ls.context.InternalClient.GetServerCatalog(ctx, &internal.GetServerCatalogRequest{})
	if err != nil {
		return fmt.Errorf("GetServerCatalog failed: %w", err)
	}
	if len(reply.GetWorlds()) == 0 {
		return fmt.Errorf("GetServerCatalog returned no world definitions")
	}
	ls.worldCatalog = reply.GetWorlds()
	for _, world := range ls.worldCatalog {
		wid := world.GetWorldId()
		if ls.channelRoutes[wid] == nil {
			ls.channelRoutes[wid] = make(map[uint32]*internal.ChannelCatalog)
		}
		for _, ch := range world.GetChannels() {
			ls.channelRoutes[wid][ch.GetChannelId()] = ch
		}
	}
	return nil
}

func (ls *LoginServer) GetWorldCatalog() []*internal.WorldCatalog {
	return ls.worldCatalog
}

func (ls *LoginServer) ResolveChannelRoute(worldId uint32, channelId uint32) (*internal.ChannelCatalog, bool) {
	worldRoutes, ok := ls.channelRoutes[worldId]
	if !ok {
		return nil, false
	}
	route, ok := worldRoutes[channelId]
	return route, ok
}

func (ls *LoginServer) handleClientDisconnect(c core.Client) {
	ic := ls.context.InternalClient
	if ic == nil {
		return
	}
	loginClient, ok := c.(*client.LoginClient)
	if !ok {
		return
	}
	accountId := loginClient.GetAccountId()
	if accountId == 0 {
		return
	}
	worldId := loginClient.GetWorldId()
	ctx, cancel := context.WithTimeout(context.Background(), core.InternalRPCPerStepTimeout)
	defer cancel()
	_, err := ic.LogoutSession(ctx, &internal.LogoutSessionRequest{
		WorldId:   worldId,
		AccountId: accountId,
	})
	if err != nil {
		log.Printf("LogoutSession (login disconnect) failed for world=%d account=%d: %v", worldId, accountId, err)
	}
}

func (ls *LoginServer) Start() error {
	log.Println("Starting MapleStory Login Server...")

	if err := ls.server.Start(ls.config.Host, ls.config.Port); err != nil {
		return err
	}

	log.Printf("Login server started on %s:%d", ls.config.Host, ls.config.Port)
	log.Printf("Loaded server catalog: worlds=%d", len(ls.worldCatalog))

	return nil
}

func (ls *LoginServer) Stop() error {
	log.Println("Shutting down login server...")
	return ls.server.Stop()
}

func (ls *LoginServer) GetStats() map[string]interface{} {
	stats := ls.server.GetStats()

	stats["server_type"] = "login"
	stats["catalog_world_count"] = len(ls.worldCatalog)

	return stats
}

func RunLoginServer() {
	config := &LoginConfig{
		Host:    "0.0.0.0",
		Port:    8484,
		WorldId: 0,
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
