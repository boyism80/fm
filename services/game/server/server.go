package server

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/common/config"
	"github.com/boyism80/fm/core"
	c_actor "github.com/boyism80/fm/core/actor"
	"github.com/boyism80/fm/core/async"
	"github.com/boyism80/fm/core/ensure"
	"github.com/boyism80/fm/core/fault"
	"github.com/boyism80/fm/core/luax"
	"github.com/boyism80/fm/core/mq"
	"github.com/boyism80/fm/core/uniqueid"
	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
	g_actor "github.com/boyism80/fm/services/game/actor"
	"github.com/boyism80/fm/services/game/client"
	"github.com/boyism80/fm/services/game/entity"
	"github.com/boyism80/fm/services/game/wz"
	lua "github.com/yuin/gopher-lua"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	mapActorCallTimeout = 30 * time.Second
	logoutWaitTimeout   = ensureWallTimeout + time.Second
)

type GameServer struct {
	*core.ServerCore
	config            *GameConfig
	resources         *wz.Resources
	maps              map[uint32]*entity.Map
	instanceMaps      map[uint32]*entity.Map
	nextInstanceID    atomic.Uint32
	nextMachineID     atomic.Uint32
	mapsMutex         sync.RWMutex
	slotInstances     map[instanceSlot]*entity.Map
	slotMutex         sync.Mutex
	packetHandler     *core.PacketHandler
	packetHandlers    *PacketHandlerRegistry
	actorSystem       *c_actor.ActorSystem
	actorRegistry     *c_actor.ActorRegistry
	nilActorPID       *actor.PID
	characterListener entity.CharacterListener
	mobListener       entity.MobListener
	mapListener       entity.MapListener
	internalClient    internal.InternalClient
	internalConn      *grpc.ClientConn
	rpcFaults         *fault.Injector
	party             *PartyCache
	guild             *GuildCache
	alliance          *AllianceCache
	stateMachines     *StateMachineRegistry
	carnivalRegistry  *entity.CarnivalRegistry
	expeditions       *entity.ExpeditionRegistry
	rabbitPartyPID    *actor.PID
	rabbitBuddyPID    *actor.PID
	rabbitGuildPID    *actor.PID
	rabbitAlliancePID *actor.PID
	rabbitGlobalPID   *actor.PID
	rabbitParcelPID   *actor.PID
	characterRuntime  *ServerCharacterRuntime
	mapSystem         mapSystem
	schedulerSystem   schedulerSystem
	partySystem       partySystem
	guildSystem       guildSystem
	allianceSystem    allianceSystem
	dispatchSystem    dispatchSystem
	ensureMu          sync.Mutex
	ensurePending     map[uint64]*ensure.EnsureDeliver
	ensureNext        atomic.Uint64
	internalHBCancel  context.CancelFunc
	megaphoneMuted    atomic.Bool
	uniqueIDs         *uniqueid.Generator
}

func (gs *GameServer) GetRootContext() *actor.RootContext {
	return gs.ServerCore.GetRootContext()
}

func (gs *GameServer) StartStateMachineActor(sm *entity.StateMachine) *actor.PID {
	if gs == nil || sm == nil || sm.Group == nil {
		return nil
	}
	name := fmt.Sprintf("state_machine_%s_%s_%d", sm.Group.Name, sm.ID, gs.nextMachineID.Add(1))
	smActor := g_actor.NewStateMachineActor(sm, gs)
	props := actor.PropsFromProducer(func() actor.Actor {
		return smActor
	})
	pid := gs.actorRegistry.GetOrCreateActor(name, props)
	if pid == nil {
		return nil
	}
	gs.stateMachines.registerActor(pid, sm)
	return pid
}

func (gs *GameServer) SendStateMachineMessage(pid *actor.PID, msg interface{}) {
	if gs == nil || pid == nil || msg == nil {
		return
	}
	gs.GetRootContext().Send(pid, msg)
}

func (gs *GameServer) ResumeLua(pid *actor.PID, root *lua.LState, thread *lua.LState, args []lua.LValue) {
	if gs == nil || pid == nil || root == nil || thread == nil {
		return
	}
	gs.GetRootContext().Send(pid, &g_actor.ResumeLua{
		Root:   root,
		Thread: thread,
		Args:   args,
	})
}

func (gs *GameServer) StopStateMachineActor(sm *entity.StateMachine) {
	if gs == nil || sm == nil || sm.ActorPID == nil {
		return
	}
	gs.stateMachines.unregisterActor(sm.ActorPID)
	gs.actorRegistry.PoisonActor(sm.ActorPID.Id, sm.ActorPID)
}

func (gs *GameServer) GetPacketHandler() *core.PacketHandler {
	if gs == nil {
		return nil
	}
	return gs.packetHandler
}

type GameConfig struct {
	Host                             string
	Port                             int
	ChannelId                        uint32
	WzPath                           string
	WorldName                        string
	WorldId                          uint32
	MaxPlayers                       int
	ExpRate                          int
	DropRate                         int
	MesoRate                         int
	InternalAddr                     string
	InternalHeartbeatIntervalSeconds int
	RabbitMQ                         config.RabbitMQEndpoint
	LuaAlwaysReload                  bool
	DueyIdentityPrompt               bool
}

func NewGameServer(config *GameConfig) (*GameServer, error) {
	serverConfig := &core.ServerConfig{
		Host: config.Host,
		Port: config.Port,
		ClientFactory: func(conn net.Conn, clientID int) (core.Client, error) {
			return client.NewGameClient(conn, clientID)
		},
	}
	server, err := core.NewServer(serverConfig)
	if err != nil {
		return nil, err
	}

	resources := wz.NewResources(config.WzPath)
	if resources == nil {
		return nil, fmt.Errorf("failed to load game resources")
	}

	actorSystem := c_actor.NewActorSystem()
	actorRegistry := c_actor.NewActorRegistry(actorSystem)
	luax.SetAlwaysReload(config.LuaAlwaysReload)

	server.SetRootContext(actorSystem.GetRoot())

	gs := &GameServer{
		ServerCore:       server,
		config:           config,
		resources:        resources,
		maps:             make(map[uint32]*entity.Map),
		instanceMaps:     make(map[uint32]*entity.Map),
		slotInstances:    make(map[instanceSlot]*entity.Map),
		packetHandler:    core.NewPacketHandler(),
		actorSystem:      actorSystem,
		actorRegistry:    actorRegistry,
		rpcFaults:        fault.NewInjector(),
		characterRuntime: nil,
		uniqueIDs:        uniqueid.NewForChannel(uint32(config.WorldId), uint32(config.ChannelId)),
	}
	gs.characterRuntime = NewServerCharacterRuntime(gs)
	gs.mapSystem = mapSystem{gs}
	gs.schedulerSystem = schedulerSystem{gs}
	gs.partySystem = partySystem{gs}
	gs.guildSystem = guildSystem{gs}
	gs.allianceSystem = allianceSystem{gs}
	gs.dispatchSystem = dispatchSystem{gs}
	gs.subscribeDeadLetters()

	if config.InternalAddr != "" {
		conn, err := grpc.NewClient(config.InternalAddr, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithUnaryInterceptor(gs.rpcFaults.UnaryClientInterceptor()))
		if err != nil {
			log.Printf("internal grpc dial %q failed: %v", config.InternalAddr, err)
		} else {
			gs.internalConn = conn
			gs.internalClient = internal.NewInternalClient(conn)
		}
	}

	gs.party = NewPartyCache(gs, config.WorldId, gs.internalClient)
	gs.guild = NewGuildCache(gs, config.WorldId, gs.internalClient)
	gs.alliance = NewAllianceCache(gs, config.WorldId, gs.internalClient)
	gs.stateMachines = NewStateMachineRegistry(gs)
	gs.carnivalRegistry = entity.NewCarnivalRegistry()
	gs.expeditions = entity.NewExpeditionRegistry(gs)

	if config.RabbitMQ.Enabled() {
		queueName := fmt.Sprintf("fm.game.w%d.c%d.party.events", config.WorldId, config.ChannelId)
		consumerTag := fmt.Sprintf("fm-game-w%d-c%d-party", config.WorldId, config.ChannelId)
		routeAll := fmt.Sprintf("fm.%d.all.party", config.WorldId)
		routeGame := fmt.Sprintf("fm.%d.%d.party", config.WorldId, config.ChannelId)

		partyDisp := mq.NewDispatcher()
		partyDisp.Bind[partyMqMemberJoined](gs)
		partyDisp.Bind[partyMqMemberLeft](gs)
		partyDisp.Bind[partyMqLeaderChanged](gs)
		partyDisp.Bind[partyMqLogOnOff](gs)
		partyDisp.Bind[partyMqDisbanded](gs)
		partyDisp.Bind[partyMqPartySync](gs)
		partyDisp.Bind[partyMqPartyInvite](gs)
		partyDisp.Bind[partyMqChat](gs)
		partyDisp.Bind[partyMqPartyInviteDenied](gs)

		rabbitCfg := mq.RabbitActorConfig{
			Root:        gs.GetRootContext(),
			Broker:      config.RabbitMQ,
			Exchange:    mq.DirectExchange,
			QueueName:   queueName,
			ConsumerTag: consumerTag,
			RoutingKeys: []string{routeAll, routeGame},
			Dispatcher:  partyDisp,
		}
		rabbitProps := actor.PropsFromProducer(func() actor.Actor {
			return mq.NewRabbitActor(rabbitCfg)
		})
		gs.rabbitPartyPID = gs.actorRegistry.GetOrCreateActor(
			fmt.Sprintf("rabbitmq_party_w%d_c%d", config.WorldId, config.ChannelId),
			rabbitProps,
		)

		buddyQueueName := fmt.Sprintf("fm.game.w%d.c%d.buddy.events", config.WorldId, config.ChannelId)
		buddyConsumerTag := fmt.Sprintf("fm-game-w%d-c%d-buddy", config.WorldId, config.ChannelId)
		buddyRouteAll := fmt.Sprintf("fm.%d.all.buddy", config.WorldId)

		buddyDisp := mq.NewDispatcher()
		buddyDisp.Bind[buddyMqChannelUpdate](gs)
		buddyDisp.Bind[buddyMqUpdate](gs)
		buddyDisp.Bind[buddyMqAddRequest](gs)
		buddyDisp.Bind[buddyMqChat](gs)

		buddyRabbitCfg := mq.RabbitActorConfig{
			Root:        gs.GetRootContext(),
			Broker:      config.RabbitMQ,
			Exchange:    mq.DirectExchange,
			QueueName:   buddyQueueName,
			ConsumerTag: buddyConsumerTag,
			RoutingKeys: []string{buddyRouteAll},
			Dispatcher:  buddyDisp,
		}
		buddyRabbitProps := actor.PropsFromProducer(func() actor.Actor {
			return mq.NewRabbitActor(buddyRabbitCfg)
		})
		gs.rabbitBuddyPID = gs.actorRegistry.GetOrCreateActor(
			fmt.Sprintf("rabbitmq_buddy_w%d_c%d", config.WorldId, config.ChannelId),
			buddyRabbitProps,
		)

		guildQueueName := fmt.Sprintf("fm.game.w%d.c%d.guild.events", config.WorldId, config.ChannelId)
		guildConsumerTag := fmt.Sprintf("fm-game-w%d-c%d-guild", config.WorldId, config.ChannelId)
		guildRouteAll := fmt.Sprintf("fm.%d.all.guild", config.WorldId)

		guildDisp := mq.NewDispatcher()
		guildDisp.Bind[guildMqCreated](gs)
		guildDisp.Bind[guildMqMemberJoined](gs)
		guildDisp.Bind[guildMqMemberLeft](gs)
		guildDisp.Bind[guildMqRankTitlesChanged](gs)
		guildDisp.Bind[guildMqMemberRankChanged](gs)
		guildDisp.Bind[guildMqEmblemChanged](gs)
		guildDisp.Bind[guildMqNoticeChanged](gs)
		guildDisp.Bind[guildMqCapacityChanged](gs)
		guildDisp.Bind[guildMqGPChanged](gs)
		guildDisp.Bind[guildMqMemberOnlineChanged](gs)
		guildDisp.Bind[guildMqDisbanded](gs)
		guildDisp.Bind[guildMqChat](gs)
		guildDisp.Bind[guildMqMessage](gs)

		guildRabbitCfg := mq.RabbitActorConfig{
			Root:        gs.GetRootContext(),
			Broker:      config.RabbitMQ,
			Exchange:    mq.DirectExchange,
			QueueName:   guildQueueName,
			ConsumerTag: guildConsumerTag,
			RoutingKeys: []string{guildRouteAll},
			Dispatcher:  guildDisp,
		}
		guildRabbitProps := actor.PropsFromProducer(func() actor.Actor {
			return mq.NewRabbitActor(guildRabbitCfg)
		})
		gs.rabbitGuildPID = gs.actorRegistry.GetOrCreateActor(
			fmt.Sprintf("rabbitmq_guild_w%d_c%d", config.WorldId, config.ChannelId),
			guildRabbitProps,
		)

		allianceQueueName := fmt.Sprintf("fm.game.w%d.c%d.alliance.events", config.WorldId, config.ChannelId)
		allianceConsumerTag := fmt.Sprintf("fm-game-w%d-c%d-alliance", config.WorldId, config.ChannelId)
		allianceRouteAll := fmt.Sprintf("fm.%d.all.alliance", config.WorldId)

		allianceDisp := mq.NewDispatcher()
		allianceDisp.Bind[allianceMqCreated](gs)
		allianceDisp.Bind[allianceMqDisbanded](gs)
		allianceDisp.Bind[allianceMqGuildLeft](gs)
		allianceDisp.Bind[allianceMqGuildAdded](gs)
		allianceDisp.Bind[allianceMqCapacityChanged](gs)
		allianceDisp.Bind[allianceMqRankTitlesChanged](gs)
		allianceDisp.Bind[allianceMqMemberRankChanged](gs)
		allianceDisp.Bind[allianceMqLeaderChanged](gs)
		allianceDisp.Bind[allianceMqNoticeChanged](gs)
		allianceDisp.Bind[allianceMqChat](gs)

		allianceRabbitCfg := mq.RabbitActorConfig{
			Root:        gs.GetRootContext(),
			Broker:      config.RabbitMQ,
			Exchange:    mq.DirectExchange,
			QueueName:   allianceQueueName,
			ConsumerTag: allianceConsumerTag,
			RoutingKeys: []string{allianceRouteAll},
			Dispatcher:  allianceDisp,
		}
		allianceRabbitProps := actor.PropsFromProducer(func() actor.Actor {
			return mq.NewRabbitActor(allianceRabbitCfg)
		})
		gs.rabbitAlliancePID = gs.actorRegistry.GetOrCreateActor(
			fmt.Sprintf("rabbitmq_alliance_w%d_c%d", config.WorldId, config.ChannelId),
			allianceRabbitProps,
		)

		globalQueueName := fmt.Sprintf("fm.game.w%d.c%d.global.events", config.WorldId, config.ChannelId)
		globalConsumerTag := fmt.Sprintf("fm-game-w%d-c%d-global", config.WorldId, config.ChannelId)
		globalRouteAll := fmt.Sprintf("fm.%d.all.global", config.WorldId)

		globalDisp := mq.NewDispatcher()
		globalDisp.Bind[globalMqServerDatetime](gs)
		globalDisp.Bind[globalMqNotice](gs)

		globalRabbitCfg := mq.RabbitActorConfig{
			Root:        gs.GetRootContext(),
			Broker:      config.RabbitMQ,
			Exchange:    mq.DirectExchange,
			QueueName:   globalQueueName,
			ConsumerTag: globalConsumerTag,
			RoutingKeys: []string{globalRouteAll},
			Dispatcher:  globalDisp,
		}
		globalRabbitProps := actor.PropsFromProducer(func() actor.Actor {
			return mq.NewRabbitActor(globalRabbitCfg)
		})
		gs.rabbitGlobalPID = gs.actorRegistry.GetOrCreateActor(
			fmt.Sprintf("rabbitmq_global_w%d_c%d", config.WorldId, config.ChannelId),
			globalRabbitProps,
		)

		parcelDisp := mq.NewDispatcher()
		parcelDisp.Bind[parcelMqArrived](gs)
		parcelDisp.Bind[cashGiftMqArrived](gs)
		parcelDisp.Bind[marriageMqChanged](gs)
		parcelDisp.Bind[marriageMqSpouseMoved](gs)

		parcelRabbitCfg := mq.RabbitActorConfig{
			Root:        gs.GetRootContext(),
			Broker:      config.RabbitMQ,
			Exchange:    mq.DirectExchange,
			QueueName:   fmt.Sprintf("fm.game.w%d.c%d.parcel.events", config.WorldId, config.ChannelId),
			ConsumerTag: fmt.Sprintf("fm-game-w%d-c%d-parcel", config.WorldId, config.ChannelId),
			RoutingKeys: []string{fmt.Sprintf("fm.%d.%d.parcel", config.WorldId, config.ChannelId)},
			Dispatcher:  parcelDisp,
		}
		parcelRabbitProps := actor.PropsFromProducer(func() actor.Actor {
			return mq.NewRabbitActor(parcelRabbitCfg)
		})
		gs.rabbitParcelPID = gs.actorRegistry.GetOrCreateActor(
			fmt.Sprintf("rabbitmq_parcel_w%d_c%d", config.WorldId, config.ChannelId),
			parcelRabbitProps,
		)
	}

	gs.characterListener = &CharacterListenerImpl{gs: gs}
	gs.mobListener = &MobListenerImpl{}
	gs.packetHandlers = NewPacketHandlerRegistry(gs)
	luax.RegisterOnCreateHook(func(luaState *lua.LState) {
		gs.registerGameLuaState(luaState)
	})
	gs.runStartupScript()

	gs.preCreateMaps()
	if err := gs.stateMachines.LoadFromScripts("script/state_machine"); err != nil {
		log.Printf("state machine load: %v", err)
	}

	gs.registerPacketHandlers()

	server.SetOnClientDisconnect(func(client core.Client) {
		gs.handleClientDisconnect(client)
	})

	return gs, nil
}

func (gs *GameServer) GetResources() *wz.Resources {
	return gs.resources
}

func (gs *GameServer) GetWorldID() uint32 {
	return gs.config.WorldId
}

func (gs *GameServer) GetExpRate() int {
	return gs.config.ExpRate
}

func (gs *GameServer) GetDropRate() int {
	return gs.config.DropRate
}

func (gs *GameServer) GetMesoRate() int {
	return gs.config.MesoRate
}

func (gs *GameServer) DueyIdentityPrompt() bool {
	return gs.config.DueyIdentityPrompt
}

func (gs *GameServer) NewUniqueID() uint64 {
	return gs.uniqueIDs.Next()
}

func (gs *GameServer) preCreateMaps() {
	log.Println("Setting up map listeners...")
	gameMapListener := NewGameMapListener(gs)
	gs.mapListener = gameMapListener

	nilMapProps := actor.PropsFromProducer(func() actor.Actor {
		return g_actor.NewMapActor(nil, gs)
	})

	nilMapPID := gs.actorRegistry.GetOrCreateActor(
		"map_nil",
		nilMapProps,
	)
	gs.nilActorPID = nilMapPID

	gs.ServerCore.SetNilActorPID(nilMapPID)

	log.Println("Pre-creating map instances...")
	for mapID := range gs.resources.Maps {
		name := fmt.Sprintf("map_%d", mapID)
		pid := gs.actorRegistry.PredictPID(name)
		mapInstance := entity.NewMap(mapID, gameMapListener, gs.mobListener, mapID, gs, pid)
		gs.maps[mapID] = mapInstance

		props := actor.PropsFromProducer(func() actor.Actor {
			return g_actor.NewMapActor(mapInstance, gs)
		})

		actual := gs.actorRegistry.GetOrCreateActor(name, props)
		if !actual.Equal(pid) {
			panic(fmt.Sprintf("map %d actor pid mismatch", mapID))
		}
	}
	log.Printf("Pre-created %d map instances", len(gs.maps))
	log.Println("Map listeners configured")
}

func (gs *GameServer) Start() error {
	log.Println("Starting MapleStory Game Server...")

	if gs.internalClient != nil {
		closeCtx, cancel := context.WithTimeout(context.Background(), core.InternalRPCPerStepTimeout)
		reply, err := gs.internalClient.CloseChannelMerchants(closeCtx, &internal.CloseChannelMerchantsRequest{
			WorldId:   gs.config.WorldId,
			ChannelId: int32(gs.config.ChannelId),
		})
		cancel()
		if err != nil {
			log.Printf("close channel hired merchants: %v", err)
		} else {
			log.Printf("Closed %d hired merchants into the store bank", reply.GetCount())
		}
	}

	if err := gs.ServerCore.Start(gs.config.Host, gs.config.Port); err != nil {
		return err
	}

	if gs.internalClient != nil && gs.config.InternalHeartbeatIntervalSeconds > 0 {
		hbCtx, cancel := context.WithCancel(context.Background())
		gs.internalHBCancel = cancel
		go gs.pingInternal(hbCtx, time.Duration(gs.config.InternalHeartbeatIntervalSeconds)*time.Second)
	}

	log.Printf("Game server started on %s:%d", gs.config.Host, gs.config.Port)
	log.Printf("World: %s, Max Players: %d", gs.config.WorldName, gs.config.MaxPlayers)
	log.Printf("Rates: Exp=%dx, Drop=%dx, Meso=%dx",
		gs.config.ExpRate, gs.config.DropRate, gs.config.MesoRate)
	if gs.rabbitPartyPID != nil {
		log.Printf("Party MQ RabbitActor running (%s)", fmt.Sprintf("fm.game.w%d.c%d.party.events", gs.config.WorldId, gs.config.ChannelId))
	}
	if gs.rabbitBuddyPID != nil {
		log.Printf("Buddy MQ RabbitActor running (%s)", fmt.Sprintf("fm.game.w%d.c%d.buddy.events", gs.config.WorldId, gs.config.ChannelId))
	}
	if gs.rabbitGuildPID != nil {
		log.Printf("Guild MQ RabbitActor running (%s)", fmt.Sprintf("fm.game.w%d.c%d.guild.events", gs.config.WorldId, gs.config.ChannelId))
	}
	if gs.rabbitAlliancePID != nil {
		log.Printf("Alliance MQ RabbitActor running (%s)", fmt.Sprintf("fm.game.w%d.c%d.alliance.events", gs.config.WorldId, gs.config.ChannelId))
	}
	if gs.rabbitGlobalPID != nil {
		log.Printf("Global MQ RabbitActor running (%s)", fmt.Sprintf("fm.game.w%d.c%d.global.events", gs.config.WorldId, gs.config.ChannelId))
	}

	if gs.resources != nil {
		stringCount := 0
		if gs.resources.Strings != nil {
			stringCount = gs.resources.Strings.CountStrings()
		}
		log.Printf("Resources loaded: %d maps, %d monsters, %d items, %d mob drops, %d reactor drops, %d strings",
			len(gs.resources.Maps), len(gs.resources.Monsters),
			len(gs.resources.Items), len(gs.resources.MobDrops), len(gs.resources.ReactorDrops), stringCount)
	}

	return nil
}

func (gs *GameServer) Stop() error {
	log.Println("Shutting down game server...")
	if gs.internalHBCancel != nil {
		gs.internalHBCancel()
		gs.internalHBCancel = nil
	}
	if gs.rabbitPartyPID != nil {
		root := gs.GetRootContext()
		if root != nil {
			root.Poison(gs.rabbitPartyPID)
		}
	}
	if gs.rabbitBuddyPID != nil {
		root := gs.GetRootContext()
		if root != nil {
			root.Poison(gs.rabbitBuddyPID)
		}
	}
	if gs.rabbitGuildPID != nil {
		root := gs.GetRootContext()
		if root != nil {
			root.Poison(gs.rabbitGuildPID)
		}
	}
	if gs.rabbitAlliancePID != nil {
		root := gs.GetRootContext()
		if root != nil {
			root.Poison(gs.rabbitAlliancePID)
		}
	}
	if gs.rabbitGlobalPID != nil {
		root := gs.GetRootContext()
		if root != nil {
			root.Poison(gs.rabbitGlobalPID)
		}
	}
	if gs.rabbitParcelPID != nil {
		root := gs.GetRootContext()
		if root != nil {
			root.Poison(gs.rabbitParcelPID)
		}
	}
	if gs.internalConn != nil {
		_ = gs.internalConn.Close()
	}
	return gs.ServerCore.Stop()
}

func (gs *GameServer) pingInternal(ctx context.Context, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
		}
		pingCtx, cancel := context.WithTimeout(ctx, core.InternalRPCPerStepTimeout)
		_, err := gs.internalClient.Ping(pingCtx, &internal.PingRequest{
			Role:      internal.ServerRole_SERVER_ROLE_GAME,
			WorldId:   gs.config.WorldId,
			ChannelId: gs.config.ChannelId,
		})
		cancel()
		if err != nil {
			log.Printf("internal ping: %v", err)
		}
	}
}

func (gs *GameServer) handleClientDisconnect(c core.Client) {
	client, ok := c.(*client.GameClient)
	if !ok {
		return
	}
	character := client.Logout()
	if character == nil {
		return
	}
	logout := character.MarkLoggedOut()

	p := async.NewTask(nil, saveCharactersPromiseTimeout)
	p.OnError(func(err error) {
		log.Printf("disconnect async: %v", err)
	})

	var entry *internal.CharacterSaveEntry
	p.DoAsync(func() error {
		entry = gs.removeCharacter(character, logout)
		return nil
	})

	if gs.internalClient == nil {
		return
	}

	if client.SessionLost() == false {
		p.ThenRPC(func(c context.Context) (*internal.SaveCharactersReply, error) {
			return gs.internalClient.SaveCharacters(c, &internal.SaveCharactersRequest{Entries: []*internal.CharacterSaveEntry{entry}}, grpc.WaitForReady(true))
		}, func(reply *internal.SaveCharactersReply) error {
			if reply.GetOk() == false {
				return fmt.Errorf("save character %d on disconnect failed; session kept", character.GetID())
			}
			return nil
		})
	}

	if character.AccountID == 0 {
		return
	}
	charID := character.GetID()
	channelID := gs.config.ChannelId
	p.ThenRPC(func(c context.Context) (*internal.LogoutSessionReply, error) {
		return gs.internalClient.LogoutSession(c, &internal.LogoutSessionRequest{
			WorldId:          gs.config.WorldId,
			AccountId:        character.AccountID,
			CharacterId:      &charID,
			ChannelId:        &channelID,
			DisconnectSource: internal.SessionDisconnectSource_SESSION_DISCONNECT_SOURCE_GAME_SERVER,
		}, grpc.WaitForReady(true))
	}, nil)
}

// removeCharacter takes a logged-out character out of this channel on whichever map actor owns it at that moment,
// and returns the save entry built there. It blocks, so it must not run on a map actor.
func (gs *GameServer) removeCharacter(character *entity.Character, logout <-chan *internal.CharacterSaveEntry) *internal.CharacterSaveEntry {
	charID := character.GetID()
	gs.expeditions.LeaveChannel(character)
	if sm := character.StateMachine(); sm != nil {
		sm.RequestLeave(character, false, entity.StateMachineLeaveDisconnect)
	}
	gs.GetDispatchSystem().Call(charID, func(ctx actor.Context) {
		_ = character.GetMap().LogoutPlayer(ctx, charID)
	})

	var entry *internal.CharacterSaveEntry
	select {
	case entry = <-logout:
	case <-time.After(logoutWaitTimeout):
		log.Printf("remove character %d: no map actor owns it; clearing timers here", charID)
		entry = character.ToProto(gs.config.WorldId)
		character.ClearTimers()
	}

	gs.ensureAbandonCharacter(charID)
	gs.characterRuntime.UnregisterCharacter(charID)
	return entry
}

func (gs *GameServer) GetStats() map[string]interface{} {
	stats := gs.ServerCore.GetStats()

	stats["server_type"] = "game"
	stats["world_name"] = gs.config.WorldName
	stats["max_players"] = gs.config.MaxPlayers
	stats["current_players"] = stats["client_count"]
	stats["exp_rate"] = gs.config.ExpRate
	stats["drop_rate"] = gs.config.DropRate
	stats["meso_rate"] = gs.config.MesoRate

	if gs.resources != nil {
		stats["maps_loaded"] = len(gs.resources.Maps)
		stats["monsters_loaded"] = len(gs.resources.Monsters)
		stats["items_loaded"] = len(gs.resources.Items)
		stats["mob_drops_loaded"] = len(gs.resources.MobDrops)
		stats["reactor_drops_loaded"] = len(gs.resources.ReactorDrops)
		if gs.resources.Strings != nil {
			stats["strings_loaded"] = gs.resources.Strings.CountStrings()
		} else {
			stats["strings_loaded"] = 0
		}
	}

	return stats
}

func RunGameServer() {

	config := &GameConfig{

		Host:            "0.0.0.0",
		Port:            8485,
		WorldName:       "Scania",
		MaxPlayers:      1000,
		ExpRate:         1,
		DropRate:        1,
		MesoRate:        1,
		LuaAlwaysReload: false,
	}

	gs, err := NewGameServer(config)
	if err != nil {
		log.Fatalf("Failed to create game server: %v", err)
	}

	if err := gs.Start(); err != nil {
		log.Fatalf("Failed to start game server: %v", err)
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	<-sigChan
	log.Println("Received shutdown signal, stopping game server...")

	if err := gs.Stop(); err != nil {
		log.Printf("Error stopping game server: %v", err)
	}
}

func RunHighRateGameServer() {
	config := &GameConfig{

		Host:            "0.0.0.0",
		Port:            8485,
		WorldName:       "HighRate",
		MaxPlayers:      2000,
		ExpRate:         10,
		DropRate:        5,
		MesoRate:        5,
		LuaAlwaysReload: false,
	}

	gs, err := NewGameServer(config)
	if err != nil {
		log.Fatalf("Failed to create high-rate game server: %v", err)
	}

	if err := gs.Start(); err != nil {
		log.Fatalf("Failed to start high-rate game server: %v", err)
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	log.Println("Shutting down high-rate game server...")
	gs.Stop()
}
