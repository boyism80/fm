package server

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/common/config"
	"github.com/boyism80/fm/core"
	c_actor "github.com/boyism80/fm/core/actor"
	"github.com/boyism80/fm/core/async"
	"github.com/boyism80/fm/core/mq"
	"github.com/boyism80/fm/core/rpc"
	"github.com/boyism80/fm/core/system"
	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
	loginactor "github.com/boyism80/fm/services/login/actor"
	"github.com/boyism80/fm/services/login/client"
)

type LoginServer struct {
	*core.ServerCore
	config  *LoginConfig
	systems *system.Systems
}

func (ls *LoginServer) handleClient(c core.Client) {
	loginClient := c.(*client.LoginClient)

	props := actor.PropsFromProducer(func() actor.Actor {
		return &loginactor.LoginLogicActor{
			Client: loginClient,
			Server: ls,
		}
	})

	pid := ls.systems.Get[*c_actor.ActorRegistry]().GetOrCreateActor(
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

	ls := &LoginServer{
		ServerCore: server,
		config:     config,
		systems:    system.NewSystems(),
	}

	actors := c_actor.NewActorRegistry(c_actor.NewActorSystem())
	ls.systems.Add[*c_actor.ActorRegistry](actors)
	server.SetRootContext(actors.Root())

	var ping func() *internal.PingRequest
	if config.LoginInstanceID != "" {
		ping = ls.ping
	} else {
		log.Printf("login_instance_id empty: internal heartbeat disabled")
	}
	conn, err := rpc.NewInternalConnection(config.internalAddr(), time.Duration(config.InternalHeartbeatIntervalSeconds)*time.Second, ping)
	if err != nil {
		return nil, err
	}
	ls.systems.Add[*rpc.InternalConnection](conn)

	catalog := newCatalogSystem(conn.Client, time.Duration(config.CatalogRetryIntervalSeconds)*time.Second, config.CatalogRetryMaxAttempts)
	if err := catalog.load(); err != nil {
		return nil, err
	}
	ls.systems.Add[*catalogSystem](catalog)
	ls.systems.Add[*mq.Consumers](mq.NewConsumers(actors, config.RabbitMQ, ls.mqQueues(catalog.worlds)))

	serverConfig.OnClientConnect = ls.handleClient
	server.SetOnClientDisconnect(ls.handleClientDisconnect)
	ls.registerPacketHandlers()

	return ls, nil
}

func (ls *LoginServer) mqQueues(worlds []*internal.WorldCatalog) []mq.Queue {
	routingKeys := make([]string, 0, len(worlds))
	for _, world := range worlds {
		routingKeys = append(routingKeys, fmt.Sprintf("fm.%d.all.global", world.GetWorldId()))
	}
	queueName := fmt.Sprintf("fm.login.%s.global.events", ls.config.LoginInstanceID)
	if queueName == "fm.login..global.events" {
		queueName = "fm.login.global.events"
	}
	consumerTag := fmt.Sprintf("fm-login-%s-global", ls.config.LoginInstanceID)
	if consumerTag == "fm-login--global" {
		consumerTag = "fm-login-global"
	}
	return []mq.Queue{{
		Name:        queueName,
		ConsumerTag: consumerTag,
		ActorName:   fmt.Sprintf("rabbitmq_login_global_%s", ls.config.LoginInstanceID),
		RoutingKeys: routingKeys,
		Bind: func(d *mq.Dispatcher) {
			d.Bind[loginGlobalMqServerDatetime](ls)
		},
	}}
}

func (ls *LoginServer) ping() *internal.PingRequest {
	return &internal.PingRequest{
		Role:            internal.ServerRole_SERVER_ROLE_LOGIN,
		LoginInstanceId: ls.config.LoginInstanceID,
	}
}

func (ls *LoginServer) handleClientDisconnect(c core.Client) {
	loginClient := c.(*client.LoginClient)
	ls.systems.Get[*c_actor.ActorRegistry]().StopActor(fmt.Sprintf("login_session_%d", loginClient.GetClientID()), loginClient.GetLogicActorPID())

	remoteAddr := c.GetConnection().RemoteAddr().String()
	accountId := loginClient.GetAccountId()
	if accountId == 0 {
		log.Printf("Login disconnect: remote=%s account=0, no logout", remoteAddr)
		return
	}
	worldId := loginClient.GetWorldId()
	transfer := loginClient.TakeTransferDisconnect()
	log.Printf("Login disconnect: remote=%s account=%d transfer=%v", remoteAddr, accountId, transfer)

	p := async.NewTask(nil, core.InternalRPCPerStepTimeout)
	p.OnError(func(err error) {
		log.Printf("LogoutSession (login disconnect) failed for world=%d account=%d: %v", worldId, accountId, err)
	})
	p.ThenRPC(func(ctx context.Context) (*internal.LogoutSessionReply, error) {
		return ls.systems.Get[*rpc.InternalConnection]().Client.LogoutSession(ctx, &internal.LogoutSessionRequest{
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

func (ls *LoginServer) Start() error {
	log.Println("Starting MapleStory Login Server...")

	if err := ls.systems.Start(); err != nil {
		return err
	}
	if err := ls.ServerCore.Start(ls.config.Host, ls.config.Port); err != nil {
		return err
	}

	log.Printf("Login server started on %s:%d", ls.config.Host, ls.config.Port)
	log.Printf("Loaded server catalog: worlds=%d", len(ls.systems.Get[*catalogSystem]().worlds))

	return nil
}

func (ls *LoginServer) Stop() error {
	log.Println("Shutting down login server...")
	err := ls.ServerCore.Stop()
	ls.systems.Stop()
	return err
}
