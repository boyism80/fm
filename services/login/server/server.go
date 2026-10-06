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
	"github.com/boyism80/fm/common/config"
	"github.com/boyism80/fm/core"
	c_actor "github.com/boyism80/fm/core/actor"
	"github.com/boyism80/fm/core/async"
	"github.com/boyism80/fm/core/ensure"
	"github.com/boyism80/fm/core/mq"
	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
	loginactor "github.com/boyism80/fm/services/login/actor"
	"github.com/boyism80/fm/services/login/client"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type LoginServer struct {
	*core.ServerCore
	config           *LoginConfig
	packetHandler    *core.PacketHandler
	internalClient   internal.InternalClient
	internalConn     *grpc.ClientConn
	internalHBCancel context.CancelFunc
	packetHandlers   *PacketHandlerRegistry
	actorSystem      *c_actor.ActorSystem
	actorRegistry    *c_actor.ActorRegistry
	worldCatalog     []*internal.WorldCatalog
	channelRoutes    map[uint32]map[uint32]*internal.ChannelCatalog
	rabbitGlobalPID  *actor.PID
}

func (ls *LoginServer) GetPacketHandler() *core.PacketHandler {
	return ls.packetHandler
}

func (ls *LoginServer) EnsureRedispatch(_ *ensure.EnsureDeliver) {}

func (ls *LoginServer) EnsureComplete(_ uint64) {}

func (ls *LoginServer) handleClient(c core.Client) {
	loginClient := c.(*client.LoginClient)

	props := actor.PropsFromProducer(func() actor.Actor {
		return &loginactor.LoginLogicActor{
			Client: loginClient,
			Server: ls,
		}
	})

	pid := ls.actorRegistry.GetOrCreateActor(
		fmt.Sprintf("login_session_%d", loginClient.GetClientID()),
		props,
	)

	loginClient.SetLogicActorPID(pid)
}

type LoginConfig struct {
	Host                             string
	Port                             int
	InitialRole                      uint32
	LoginInstanceID                  string
	InternalHeartbeatIntervalSeconds int
	InternalHost                     string
	InternalPort                     int
	RabbitMQ                         config.RabbitMQEndpoint
	CatalogRetryIntervalSeconds      int
	CatalogRetryMaxAttempts          int
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

	addr := config.internalAddr()
	if addr == "" {
		return nil, fmt.Errorf("login server requires internal gRPC endpoint")
	}
	internalConn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("internal gRPC dial %s: %w", addr, err)
	}
	log.Printf("Internal gRPC client connected to %s", addr)

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
		internalClient: internal.NewInternalClient(internalConn),
		internalConn:   internalConn,
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

	if config.RabbitMQ.Enabled() && len(ls.worldCatalog) > 0 {
		routingKeys := make([]string, 0, len(ls.worldCatalog))
		for _, world := range ls.worldCatalog {
			routingKeys = append(routingKeys, fmt.Sprintf("fm.%d.all.global", world.GetWorldId()))
		}
		queueName := fmt.Sprintf("fm.login.%s.global.events", config.LoginInstanceID)
		if queueName == "fm.login..global.events" {
			queueName = "fm.login.global.events"
		}
		consumerTag := fmt.Sprintf("fm-login-%s-global", config.LoginInstanceID)
		if consumerTag == "fm-login--global" {
			consumerTag = "fm-login-global"
		}

		globalDisp := mq.NewDispatcher()
		mq.Bind[*LoginServer, loginGlobalMqServerDatetime](ls, globalDisp)

		globalRabbitCfg := mq.RabbitActorConfig{
			Root:        ls.GetRootContext(),
			Broker:      config.RabbitMQ,
			Exchange:    mq.DirectExchange,
			QueueName:   queueName,
			ConsumerTag: consumerTag,
			RoutingKeys: routingKeys,
			Dispatcher:  globalDisp,
		}
		globalRabbitProps := actor.PropsFromProducer(func() actor.Actor {
			return mq.NewRabbitActor(globalRabbitCfg)
		})
		ls.rabbitGlobalPID = ls.actorRegistry.GetOrCreateActor(
			fmt.Sprintf("rabbitmq_login_global_%s", config.LoginInstanceID),
			globalRabbitProps,
		)
	}

	return ls, nil
}

func (ls *LoginServer) loadServerCatalog() error {
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

func (ls *LoginServer) handleClientDisconnect(c core.Client) {
	loginClient := c.(*client.LoginClient)
	ls.actorRegistry.StopActor(fmt.Sprintf("login_session_%d", loginClient.GetClientID()), loginClient.GetLogicActorPID())

	remoteAddr := c.GetConnection().RemoteAddr().String()
	accountId := loginClient.GetAccountId()
	if accountId == 0 {
		log.Printf("Login disconnect: remote=%s account=0, no logout", remoteAddr)
		return
	}
	worldId := loginClient.GetWorldId()
	transfer := loginClient.TakeTransferDisconnect()
	log.Printf("Login disconnect: remote=%s account=%d transfer=%v", remoteAddr, accountId, transfer)

	p := async.NewPromise(nil, core.InternalRPCPerStepTimeout)
	p.OnError(func(err error) {
		log.Printf("LogoutSession (login disconnect) failed for world=%d account=%d: %v", worldId, accountId, err)
	})
	async.ThenRPC(p, func(ctx context.Context) (*internal.LogoutSessionReply, error) {
		return ls.internalClient.LogoutSession(ctx, &internal.LogoutSessionRequest{
			WorldId:            worldId,
			AccountId:          accountId,
			DisconnectSource:   internal.SessionDisconnectSource_SESSION_DISCONNECT_SOURCE_LOGIN_SERVER,
			TransferDisconnect: transfer,
		})
	}, func(reply *internal.LogoutSessionReply) error {
		log.Printf("LogoutSession (login disconnect): account=%d ok=%v", accountId, reply.GetOk())
		return nil
	})

}

func (ls *LoginServer) pingInternal(ctx context.Context, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
		}
		pingCtx, cancel := context.WithTimeout(ctx, core.InternalRPCPerStepTimeout)
		_, err := ls.internalClient.Ping(pingCtx, &internal.PingRequest{
			Role:            internal.ServerRole_SERVER_ROLE_LOGIN,
			LoginInstanceId: ls.config.LoginInstanceID,
		})
		cancel()
		if err != nil {
			log.Printf("internal ping: %v", err)
		}
	}
}

func (ls *LoginServer) Start() error {
	log.Println("Starting MapleStory Login Server...")

	if err := ls.ServerCore.Start(ls.config.Host, ls.config.Port); err != nil {
		return err
	}

	if ls.config.InternalHeartbeatIntervalSeconds > 0 && ls.config.LoginInstanceID != "" {
		hbCtx, cancel := context.WithCancel(context.Background())
		ls.internalHBCancel = cancel
		go ls.pingInternal(hbCtx, time.Duration(ls.config.InternalHeartbeatIntervalSeconds)*time.Second)
	} else if ls.config.LoginInstanceID == "" {
		log.Printf("login_instance_id empty: internal heartbeat disabled")
	}

	log.Printf("Login server started on %s:%d", ls.config.Host, ls.config.Port)
	log.Printf("Loaded server catalog: worlds=%d", len(ls.worldCatalog))

	return nil
}

func (ls *LoginServer) Stop() error {
	log.Println("Shutting down login server...")
	if ls.internalHBCancel != nil {
		ls.internalHBCancel()
		ls.internalHBCancel = nil
	}
	if ls.rabbitGlobalPID != nil {
		root := ls.GetRootContext()
		if root != nil {
			root.Poison(ls.rabbitGlobalPID)
		}
	}
	if ls.internalConn != nil {
		_ = ls.internalConn.Close()
		ls.internalConn = nil
	}
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
		Host:            "0.0.0.0",
		Port:            8484,
		LoginInstanceID: "dev-login",
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
