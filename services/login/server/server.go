package server

import (
	"context"
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
	"github.com/boyism80/fm/core/async"
	"github.com/boyism80/fm/core/ensure"
	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
	loginactor "github.com/boyism80/fm/services/login/actor"
	"github.com/boyism80/fm/services/login/client"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type LoginServer struct {
	*core.ServerCore
	config         *LoginConfig
	packetHandler  *core.PacketHandler
	internalClient internal.InternalClient
	packetHandlers *PacketHandlerRegistry
	actorSystem    *c_actor.ActorSystem
	actorRegistry  *c_actor.ActorRegistry
	worldCatalog   []*internal.WorldCatalog
	channelRoutes  map[uint32]map[uint32]*internal.ChannelCatalog
}

func (ls *LoginServer) GetPacketHandler() *core.PacketHandler {
	if ls == nil {
		return nil
	}
	return ls.packetHandler
}

func (ls *LoginServer) EnsureRedispatch(_ *ensure.EnsureDeliver) {}

func (ls *LoginServer) EnsureComplete(_ uint64) {}

func (ls *LoginServer) handleClient(c core.Client) {
	loginClient, ok := c.(*client.LoginClient)
	if !ok {
		log.Printf("Client is not a LoginClient")
		return
	}

	props := actor.PropsFromProducer(func() actor.Actor {
		return &loginactor.LoginLogicActor{
			Client: loginClient,
			Server: ls,
		}
	})

	pid := ls.actorRegistry.GetOrCreateActor(
		fmt.Sprintf("login_session_%d", loginClient.GetFd()),
		props,
	)

	loginClient.SetLogicActorPID(pid)
}

type LoginConfig struct {
	Host                        string
	Port                        int
	WorldId                     uint32
	InitialRole                 uint32
	InternalHost                string
	InternalPort                int
	CatalogRetryIntervalSeconds int
	CatalogRetryMaxAttempts     int
}

func (c *LoginConfig) internalAddr() string {
	if c.InternalHost == "" || c.InternalPort == 0 {
		return ""
	}
	return fmt.Sprintf("%s:%d", c.InternalHost, c.InternalPort)
}

func NewLoginServer(config *LoginConfig) (*LoginServer, error) {
	if config.CatalogRetryIntervalSeconds <= 0 {
		config.CatalogRetryIntervalSeconds = 2
	}
	if config.CatalogRetryMaxAttempts <= 0 {
		config.CatalogRetryMaxAttempts = 30
	}

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
		ServerCore:     server,
		config:         config,
		packetHandler:  core.NewPacketHandler(),
		internalClient: internalClient,
		packetHandlers: NewPacketHandlerRegistry(nil),
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
	if ls.internalClient == nil {
		return fmt.Errorf("internal gRPC client is required")
	}
	interval := time.Duration(ls.config.CatalogRetryIntervalSeconds) * time.Second
	maxAttempts := ls.config.CatalogRetryMaxAttempts

	var lastErr error
	for attempt := 1; attempt <= maxAttempts; attempt++ {
		ctx, cancel := context.WithTimeout(context.Background(), core.InternalRPCPerStepTimeout)
		reply, err := ls.internalClient.GetServerCatalog(ctx, &internal.GetServerCatalogRequest{})
		cancel()
		if err == nil {
			worlds := reply.GetWorlds()
			if len(worlds) == 0 {
				lastErr = fmt.Errorf("GetServerCatalog returned no world definitions")
			} else {
				ls.worldCatalog = worlds
				ls.channelRoutes = make(map[uint32]map[uint32]*internal.ChannelCatalog)
				for _, world := range ls.worldCatalog {
					wid := world.GetWorldId()
					if ls.channelRoutes[wid] == nil {
						ls.channelRoutes[wid] = make(map[uint32]*internal.ChannelCatalog)
					}
					for _, ch := range world.GetChannels() {
						ls.channelRoutes[wid][ch.GetChannelId()] = ch
					}
				}
				if attempt > 1 {
					log.Printf("GetServerCatalog succeeded after retry (%d/%d)", attempt, maxAttempts)
				}
				return nil
			}
		} else {
			lastErr = fmt.Errorf("GetServerCatalog failed: %w", err)
		}

		if attempt < maxAttempts {
			log.Printf(
				"GetServerCatalog attempt %d/%d failed: %v (retrying in %s)",
				attempt,
				maxAttempts,
				lastErr,
				interval,
			)
			time.Sleep(interval)
		}
	}

	return fmt.Errorf("GetServerCatalog failed after %d attempts: %w", maxAttempts, lastErr)
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
	ic := ls.internalClient
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
	transfer := loginClient.TakeTransferDisconnect()

	p := async.NewPromise(nil, core.InternalRPCPerStepTimeout)
	p.OnError(func(err error) {
		log.Printf("LogoutSession (login disconnect) failed for world=%d account=%d: %v", worldId, accountId, err)
	})
	async.ThenRPC(p, func(ctx context.Context) (*internal.LogoutSessionReply, error) {
		return ic.LogoutSession(ctx, &internal.LogoutSessionRequest{
			WorldId:            worldId,
			AccountId:          accountId,
			DisconnectSource:   internal.SessionDisconnectSource_SESSION_DISCONNECT_SOURCE_LOGIN_SERVER,
			TransferDisconnect: transfer,
		})
	}, func(*internal.LogoutSessionReply) error {
		return nil
	})
	p.Run()
}

func (ls *LoginServer) Start() error {
	log.Println("Starting MapleStory Login Server...")

	if err := ls.ServerCore.Start(ls.config.Host, ls.config.Port); err != nil {
		return err
	}

	log.Printf("Login server started on %s:%d", ls.config.Host, ls.config.Port)
	log.Printf("Loaded server catalog: worlds=%d", len(ls.worldCatalog))

	return nil
}

func (ls *LoginServer) Stop() error {
	log.Println("Shutting down login server...")
	return ls.ServerCore.Stop()
}

func (ls *LoginServer) GetStats() map[string]interface{} {
	stats := ls.ServerCore.GetStats()

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
