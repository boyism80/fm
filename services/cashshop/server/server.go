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
	"github.com/boyism80/fm/core/ensure"
	"github.com/boyism80/fm/core/uniqueid"
	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
	"github.com/boyism80/fm/services/cashshop/client"
	"github.com/boyism80/fm/services/game/wz"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type CashShopServer struct {
	*core.ServerCore
	config           *config.CashShop
	packetHandler    *core.PacketHandler
	internalClient   internal.InternalClient
	internalConn     *grpc.ClientConn
	internalHBCancel context.CancelFunc
	actorRegistry    *c_actor.ActorRegistry
	resources        *wz.Resources
	uniqueIDs        *uniqueid.Generator
}

func NewCashShopServer(cfg *config.CashShop) (*CashShopServer, error) {
	addr := cfg.Internal.GRPCAddr()
	if addr == "" {
		return nil, fmt.Errorf("cash shop server requires internal gRPC endpoint")
	}
	internalConn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("internal gRPC dial %s: %w", addr, err)
	}

	actorSystem := c_actor.NewActorSystem()
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
	server.SetRootContext(actorSystem.GetRoot())

	cs := &CashShopServer{
		ServerCore:     server,
		config:         cfg,
		packetHandler:  core.NewPacketHandler(),
		internalClient: internal.NewInternalClient(internalConn),
		internalConn:   internalConn,
		actorRegistry:  c_actor.NewActorRegistry(actorSystem),
		resources:      wz.NewResources(cfg.WzPath),
		uniqueIDs:      uniqueid.NewForCashShop(uint32(cfg.WorldID), uint32(cfg.CashShopID)),
	}
	serverConfig.OnClientConnect = cs.handleClient
	server.SetOnClientDisconnect(cs.handleClientDisconnect)
	cs.registerPacketHandlers()
	return cs, nil
}

func (cs *CashShopServer) GetPacketHandler() *core.PacketHandler {
	return cs.packetHandler
}

func (cs *CashShopServer) EnsureRedispatch(_ *ensure.EnsureDeliver) {}

func (cs *CashShopServer) EnsureComplete(_ uint64) {}

func (cs *CashShopServer) GetResources() *wz.Resources {
	return cs.resources
}

func (cs *CashShopServer) NewUniqueID() uint64 {
	return cs.uniqueIDs.Next()
}

func (cs *CashShopServer) worldID() uint32 {
	return uint32(cs.config.WorldID)
}

func (cs *CashShopServer) cashShopID() uint32 {
	return uint32(cs.config.CashShopID)
}

func (cs *CashShopServer) handleClient(c core.Client) {
	csClient := c.(*client.CashShopClient)
	props := actor.PropsFromProducer(func() actor.Actor {
		return &sessionActor{client: csClient, server: cs}
	})
	pid := cs.actorRegistry.GetOrCreateActor(fmt.Sprintf("cashshop_session_%d", csClient.GetClientID()), props)
	csClient.SetLogicActorPID(pid)
}

func (cs *CashShopServer) handleClientDisconnect(c core.Client) {
	csClient := c.(*client.CashShopClient)
	cs.actorRegistry.StopActor(fmt.Sprintf("cashshop_session_%d", csClient.GetClientID()), csClient.GetLogicActorPID())

	character := csClient.Disconnect()
	if character == nil {
		return
	}

	characterID := character.ID()
	cashShopID := cs.cashShopID()
	p := async.NewTask(nil, core.InternalRPCPerStepTimeout)
	p.OnError(func(err error) {
		log.Printf("LogoutSession (cash shop disconnect) failed for character %d: %v", characterID, err)
	})
	p.ThenRPC(func(ctx context.Context) (*internal.LogoutSessionReply, error) {
		return cs.internalClient.LogoutSession(ctx, &internal.LogoutSessionRequest{
			WorldId:          cs.worldID(),
			AccountId:        character.AccountID(),
			DisconnectSource: internal.SessionDisconnectSource_SESSION_DISCONNECT_SOURCE_CASH_SHOP_SERVER,
			CharacterId:      &characterID,
			CashShopId:       &cashShopID,
		})
	}, nil)
}

func (cs *CashShopServer) pingInternal(ctx context.Context, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		pingCtx, cancel := context.WithTimeout(ctx, core.InternalRPCPerStepTimeout)
		_, err := cs.internalClient.Ping(pingCtx, &internal.PingRequest{
			Role:       internal.ServerRole_SERVER_ROLE_CASH_SHOP,
			WorldId:    cs.worldID(),
			CashShopId: cs.cashShopID(),
		})
		cancel()
		if err != nil {
			log.Printf("internal ping: %v", err)
		}

		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
		}
	}
}

func (cs *CashShopServer) Start() error {
	if err := cs.ServerCore.Start(cs.config.Host, cs.config.Port); err != nil {
		return err
	}

	interval := cs.config.Internal.HeartbeatIntervalSeconds
	if interval <= 0 {
		interval = 15
	}
	hbCtx, cancel := context.WithCancel(context.Background())
	cs.internalHBCancel = cancel
	go cs.pingInternal(hbCtx, time.Duration(interval)*time.Second)

	log.Printf("Cash shop server %d (world %d) started on %s:%d", cs.config.CashShopID, cs.config.WorldID, cs.config.Host, cs.config.Port)
	return nil
}

func (cs *CashShopServer) Stop() error {
	if cs.internalHBCancel != nil {
		cs.internalHBCancel()
	}
	if cs.internalConn != nil {
		_ = cs.internalConn.Close()
	}
	return cs.ServerCore.Stop()
}
