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
	"github.com/boyism80/fm/core/rpc"
	"github.com/boyism80/fm/core/system"
	"github.com/boyism80/fm/core/uniqueid"
	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
	"github.com/boyism80/fm/services/cashshop/client"
	"github.com/boyism80/fm/services/game/wz"
)

type CashShopServer struct {
	*core.ServerCore
	config  *config.CashShop
	systems *system.Systems
}

func NewCashShopServer(cfg *config.CashShop) (*CashShopServer, error) {
	serverConfig := &core.ServerConfig{
		Host: cfg.Host,
		Port: cfg.Port,
		ClientFactory: func(conn net.Conn, clientID int) (core.Client, error) {
			return client.NewCashShopClient(conn, clientID)
		},
	}
	server, err := core.NewServer(serverConfig)
	if err != nil {
		return nil, err
	}

	cs := &CashShopServer{
		ServerCore: server,
		config:     cfg,
		systems:    system.NewSystems(),
	}

	actors := c_actor.NewActorRegistry(c_actor.NewActorSystem())
	cs.systems.Add[*c_actor.ActorRegistry](actors)
	server.SetRootContext(actors.Root())

	heartbeat := cfg.Internal.HeartbeatIntervalSeconds
	if heartbeat <= 0 {
		heartbeat = 15
	}
	conn, err := rpc.NewInternalConnection(cfg.Internal.GRPCAddr(), time.Duration(heartbeat)*time.Second, cs.ping)
	if err != nil {
		return nil, err
	}
	cs.systems.Add[*rpc.InternalConnection](conn)
	cs.systems.Add[*wz.Resources](wz.NewResources(cfg.WzPath))
	cs.systems.Add[*uniqueid.Generator](uniqueid.NewForCashShop(uint32(cfg.WorldID), uint32(cfg.CashShopID)))

	serverConfig.OnClientConnect = cs.handleClient
	server.SetOnClientDisconnect(cs.handleClientDisconnect)
	cs.registerPacketHandlers()
	return cs, nil
}

func (cs *CashShopServer) GetResources() *wz.Resources {
	return cs.systems.Get[*wz.Resources]()
}

func (cs *CashShopServer) NewUniqueID() uint64 {
	return cs.systems.Get[*uniqueid.Generator]().Next()
}

func (cs *CashShopServer) handleClient(c core.Client) {
	csClient := c.(*client.CashShopClient)
	props := actor.PropsFromProducer(func() actor.Actor {
		return &sessionActor{client: csClient, server: cs}
	})
	pid := cs.systems.Get[*c_actor.ActorRegistry]().GetOrCreateActor(fmt.Sprintf("cashshop_session_%d", csClient.GetClientID()), props)
	csClient.SetLogicActorPID(pid)
}

func (cs *CashShopServer) handleClientDisconnect(c core.Client) {
	csClient := c.(*client.CashShopClient)
	cs.systems.Get[*c_actor.ActorRegistry]().StopActor(fmt.Sprintf("cashshop_session_%d", csClient.GetClientID()), csClient.GetLogicActorPID())

	character := csClient.Disconnect()
	if character == nil {
		return
	}

	characterID := character.ID()
	cashShopID := uint32(cs.config.CashShopID)
	p := async.NewTask(nil, core.InternalRPCPerStepTimeout)
	p.OnError(func(err error) {
		log.Printf("LogoutSession (cash shop disconnect) failed for character %d: %v", characterID, err)
	})
	p.ThenRPC(func(ctx context.Context) (*internal.LogoutSessionReply, error) {
		return cs.systems.Get[*rpc.InternalConnection]().Client.LogoutSession(ctx, &internal.LogoutSessionRequest{
			WorldId:          uint32(cs.config.WorldID),
			AccountId:        character.AccountID(),
			DisconnectSource: internal.SessionDisconnectSource_SESSION_DISCONNECT_SOURCE_CASH_SHOP_SERVER,
			CharacterId:      &characterID,
			CashShopId:       &cashShopID,
		})
	}, nil)
}

func (cs *CashShopServer) ping() *internal.PingRequest {
	online := uint32(0)
	for _, c := range cs.Clients() {
		if csClient, ok := c.(*client.CashShopClient); ok && csClient.GetCharacter() != nil {
			online++
		}
	}
	return &internal.PingRequest{
		Role:        internal.ServerRole_SERVER_ROLE_CASH_SHOP,
		WorldId:     uint32(cs.config.WorldID),
		CashShopId:  uint32(cs.config.CashShopID),
		OnlineUsers: online,
	}
}

func (cs *CashShopServer) Start() error {
	if err := cs.systems.Start(); err != nil {
		return err
	}
	if err := cs.ServerCore.Start(cs.config.Host, cs.config.Port); err != nil {
		return err
	}

	log.Printf("Cash shop server %d (world %d) started on %s:%d", cs.config.CashShopID, cs.config.WorldID, cs.config.Host, cs.config.Port)
	return nil
}

func (cs *CashShopServer) Stop() error {
	err := cs.ServerCore.Stop()
	cs.systems.Stop()
	return err
}
