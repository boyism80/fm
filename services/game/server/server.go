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
	"github.com/boyism80/fm/core/luax"
	"github.com/boyism80/fm/core/mq"
	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
	"github.com/boyism80/fm/services/common/globaltimer"
	g_actor "github.com/boyism80/fm/services/game/actor"
	"github.com/boyism80/fm/services/game/client"
	gameconst "github.com/boyism80/fm/services/game/constant"
	"github.com/boyism80/fm/services/game/entity"
	gamegtimers "github.com/boyism80/fm/services/game/gtimers"
	"github.com/boyism80/fm/services/game/wz"
	"github.com/boyism80/fm/types"
	lua "github.com/yuin/gopher-lua"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const mapActorCallTimeout = 30 * time.Second

// grpcSaveCharacters performs one SaveCharacters RPC (used by persist promises).
func (gs *GameServer) grpcSaveCharacters(ctx context.Context, chars []*entity.Character) (*internal.SaveCharactersReply, error) {
	if gs == nil || gs.internalClient == nil || len(chars) == 0 {
		return &internal.SaveCharactersReply{}, nil
	}
	worldID := uint32(gs.config.WorldId)
	entries := make([]*internal.CharacterSaveEntry, 0, len(chars))
	for _, ch := range chars {
		if ch == nil {
			continue
		}
		if entry := ch.ToProto(worldID); entry != nil {
			entries = append(entries, entry)
		}
	}
	if len(entries) == 0 {
		return &internal.SaveCharactersReply{}, nil
	}
	reply, err := gs.internalClient.SaveCharacters(ctx, &internal.SaveCharactersRequest{Entries: entries})
	if err != nil {
		return nil, fmt.Errorf("SaveCharacters rpc: %w", err)
	}
	return reply, nil
}

type GameServer struct {
	*core.ServerCore
	config            *GameConfig
	resources         *wz.Resources
	maps              map[uint32]*entity.Map
	mapsMutex         sync.RWMutex
	packetHandler     *core.PacketHandler
	packetHandlers    *PacketHandlerRegistry
	actorSystem       *c_actor.ActorSystem
	actorRegistry     *c_actor.ActorRegistry
	nilMapActorPID    *actor.PID
	characterListener entity.CharacterListener
	internalClient    internal.InternalClient
	internalConn      *grpc.ClientConn
	party             *PartyContainer
	rabbitPartyPID    *actor.PID
	characterRuntime  *ServerCharacterRuntime
	ensureMu          sync.Mutex
	ensurePending     map[uint64]*ensure.EnsureDeliver
	ensureNext        atomic.Uint64
	internalHBCancel  context.CancelFunc
}

func (gs *GameServer) GetRootContext() *actor.RootContext {
	return gs.ServerCore.GetRootContext()
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
		packetHandler:    core.NewPacketHandler(),
		actorSystem:      actorSystem,
		actorRegistry:    actorRegistry,
		characterRuntime: nil,
	}
	gs.characterRuntime = NewServerCharacterRuntime(gs)

	if config.InternalAddr != "" {
		conn, err := grpc.NewClient(config.InternalAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Printf("internal grpc dial %q failed: %v", config.InternalAddr, err)
		} else {
			gs.internalConn = conn
			gs.internalClient = internal.NewInternalClient(conn)
		}
	}

	gs.party = NewPartyContainer(gs, config.WorldId, gs.internalClient)

	if config.RabbitMQ.Enabled() {
		queueName := fmt.Sprintf("fm.game.w%d.c%d.party.events", config.WorldId, config.ChannelId)
		consumerTag := fmt.Sprintf("fm-game-w%d-c%d-party", config.WorldId, config.ChannelId)
		routeAll := fmt.Sprintf("fm.%d.all.party", config.WorldId)
		routeGame := fmt.Sprintf("fm.%d.%d.party", config.WorldId, config.ChannelId)

		partyDisp := mq.NewDispatcher()
		mq.Bind[*GameServer, partyMqCreated](gs, partyDisp)
		mq.Bind[*GameServer, partyMqMemberJoined](gs, partyDisp)
		mq.Bind[*GameServer, partyMqMemberLeft](gs, partyDisp)
		mq.Bind[*GameServer, partyMqLeaderChanged](gs, partyDisp)
		mq.Bind[*GameServer, partyMqLogOnOff](gs, partyDisp)
		mq.Bind[*GameServer, partyMqDisbanded](gs, partyDisp)
		mq.Bind[*GameServer, partyMqPartySync](gs, partyDisp)
		mq.Bind[*GameServer, partyMqPartyInvite](gs, partyDisp)
		mq.Bind[*GameServer, partyMqMultiChat](gs, partyDisp)
		mq.Bind[*GameServer, partyMqPartyInviteDenied](gs, partyDisp)

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
	}

	gs.characterListener = &CharacterListenerImpl{gs: gs}
	gs.packetHandlers = NewPacketHandlerRegistry(gs)
	luax.RegisterOnCreateHook(func(luaState *lua.LState) {
		registerGameLuaState(gs, luaState)
	})

	gs.preCreateMaps()

	gs.registerPacketHandlers()

	server.SetOnClientDisconnect(func(client core.Client) {
		gs.handleClientDisconnect(client)
	})

	return gs, nil
}

func (gs *GameServer) GetResources() *wz.Resources {
	return gs.resources
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

func (gs *GameServer) preCreateMaps() {
	log.Println("Setting up map listeners...")
	gameMapListener := NewGameMapListener(gs)

	nilMapProps := actor.PropsFromProducer(func() actor.Actor {
		return &g_actor.MapActor{
			Map:       nil,
			GameWorld: gs,
		}
	})

	nilMapPID := gs.actorRegistry.GetOrCreateActor(
		"map_nil",
		nilMapProps,
	)
	gs.nilMapActorPID = nilMapPID

	gs.ServerCore.SetNilMapActorPID(nilMapPID)

	log.Println("Pre-creating map instances...")
	for mapID := range gs.resources.Maps {
		mapInstance := entity.NewMap(mapID, gameMapListener, mapID, gs)
		gs.maps[mapID] = mapInstance

		props := actor.PropsFromProducer(func() actor.Actor {
			return &g_actor.MapActor{
				Map:       mapInstance,
				GameWorld: gs,
			}
		})

		pid := gs.actorRegistry.GetOrCreateActor(
			fmt.Sprintf("map_%d", mapID),
			props,
		)

		mapInstance.SetActorPID(pid)
	}
	log.Printf("Pre-created %d map instances", len(gs.maps))
	log.Println("Map listeners configured")
}

func (gs *GameServer) Start() error {
	log.Println("Starting MapleStory Game Server...")

	if err := gs.ServerCore.Start(gs.config.Host, gs.config.Port); err != nil {
		return err
	}

	if gs.internalClient != nil && gs.config.InternalHeartbeatIntervalSeconds > 0 {
		hbCtx, cancel := context.WithCancel(context.Background())
		gs.internalHBCancel = cancel
		wid := gs.config.WorldId
		ch := gs.config.ChannelId
		iv := gs.config.InternalHeartbeatIntervalSeconds
		gamegtimers.WireInternalPing(gs.internalClient, time.Duration(iv)*time.Second, wid, ch)
		reg := globaltimer.NewRegistry()
		globaltimer.RegisterTimer[*gamegtimers.InternalPingTimer](reg)
		reg.Start(hbCtx)
	}

	log.Printf("Game server started on %s:%d", gs.config.Host, gs.config.Port)
	log.Printf("World: %s, Max Players: %d", gs.config.WorldName, gs.config.MaxPlayers)
	log.Printf("Rates: Exp=%dx, Drop=%dx, Meso=%dx",
		gs.config.ExpRate, gs.config.DropRate, gs.config.MesoRate)
	if gs.rabbitPartyPID != nil {
		log.Printf("Party MQ RabbitActor running (%s)", fmt.Sprintf("fm.game.w%d.c%d.party.events", gs.config.WorldId, gs.config.ChannelId))
	}

	if gs.resources != nil {
		stringCount := 0
		if gs.resources.Strings != nil {
			stringCount = gs.resources.Strings.CountStrings()
		}
		log.Printf("Resources loaded: %d maps, %d monsters, %d items, %d drops, %d strings",
			len(gs.resources.Maps), len(gs.resources.Monsters),
			len(gs.resources.Items), len(gs.resources.Drops), stringCount)
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
	if gs.internalConn != nil {
		_ = gs.internalConn.Close()
	}
	return gs.ServerCore.Stop()
}

func (gs *GameServer) GetPartyByID(partyID uint32) *entity.Party {
	if gs == nil || gs.party == nil {
		return nil
	}
	return gs.party.Get(partyID)
}

func (gs *GameServer) GetMap(mapID uint32) *entity.Map {
	gs.mapsMutex.RLock()
	defer gs.mapsMutex.RUnlock()
	return gs.maps[mapID]
}

func (gs *GameServer) RequestWarp(character *entity.Character, targetMap *entity.Map, spawnPoint uint8) error {
	if targetMap == nil {
		return fmt.Errorf("target map is nil")
	}
	if character == nil {
		return fmt.Errorf("character is nil")
	}
	currentMap := character.GetMap()
	if currentMap != nil {
		currentMap.RemovePlayer(character.GetID())
	}
	targetPID := targetMap.GetActorPID()
	if targetPID == nil {
		return fmt.Errorf("target map actor not found")
	}
	gs.GetRootContext().Send(targetPID, &g_actor.WarpCharacter{
		Character: character,
		Portal:    spawnPoint,
	})
	return nil
}

func (gs *GameServer) DispatchRunCharacterTimer(pid *actor.PID, payload *c_actor.RunCharacterTimer) {
	if pid == nil || payload == nil {
		return
	}
	if root := gs.GetRootContext(); root != nil {
		root.Send(pid, payload)
	}
}

func (gs *GameServer) RequestSpawnReturnMapDoor(ch *entity.Character, skillID gameconst.SkillID) {
	if ch == nil {
		return
	}
	m := ch.GetMap()
	if m == nil || m.Wz == nil {
		return
	}
	destMapID := uint32(m.Wz.ReturnMapId)
	if destMapID == 0 || destMapID == uint32(m.Wz.ID) {
		return
	}
	destMap := gs.GetMap(destMapID)
	if destMap == nil {
		ch.Listener.OnMessage(ch, gameconst.MSG_PINK_TEXT, gameconst.DoorNoTownPortalMessage)
		return
	}
	destPID := destMap.GetActorPID()
	if destPID == nil {
		ch.Listener.OnMessage(ch, gameconst.MSG_PINK_TEXT, gameconst.DoorNoTownPortalMessage)
		return
	}
	srcPID := m.GetActorPID()
	if srcPID == nil {
		ch.Listener.OnMessage(ch, gameconst.MSG_PINK_TEXT, gameconst.DoorNoTownPortalMessage)
		return
	}
	root := gs.GetRootContext()
	if root == nil {
		return
	}
	fieldAnchorPt := types.Point[int16]{X: ch.Position.X, Y: ch.Position.Y}
	var closestPortalID uint8
	if id, ok := m.Wz.FindClosestDoorReturnPortalSpawnID(fieldAnchorPt); ok {
		closestPortalID = id
	} else {
		closestPortalID = m.Wz.FindClosestPortalSpawnID(fieldAnchorPt)
	}
	slot := 0
	if gs.party != nil {
		slot = gs.party.PartyMemberIndex(ch.GetID(), ch.GetPartyID())
	}
	root.Send(destPID, &g_actor.RequestSpawnDoor{
		ReplyTo:        srcPID,
		CharacterID:    ch.GetID(),
		OwnerID:        ch.GetID(),
		SkillID:        skillID,
		FieldMapID:     uint32(m.Wz.ID),
		FieldPortalID:  closestPortalID,
		PartyOwnerSlot: slot,
		PartyID:        ch.GetPartyID(),
		FieldAnchor:    ch.Position,
	})
}

func (gs *GameServer) NotifyDoorRemove(ownerID uint32, skillID uint32, counterpartMapWZID uint32) {
	mapInstance := gs.GetMap(counterpartMapWZID)
	if mapInstance == nil {
		return
	}
	pid := mapInstance.GetActorPID()
	if pid == nil {
		return
	}
	if root := gs.GetRootContext(); root != nil {
		root.Send(pid, &g_actor.RemoveDoor{
			OwnerID: ownerID,
			SkillID: skillID,
		})
	}
}

func (gs *GameServer) runCharacterLogoutScript(ch *entity.Character) {
	if ch == nil {
		return
	}
	mapInstance := ch.GetMap()
	if mapInstance == nil {
		return
	}
	root := mapInstance.GetLuaRoot()
	if root == nil {
		return
	}
	thread, err := luax.NewThread(root, "script/script.lua")
	if err != nil {
		log.Printf("on_logout: %v", err)
		return
	}
	if _, err := luax.Call(thread, "on_logout", ch); err != nil {
		log.Printf("on_logout: %v", err)
	}
}

func (gs *GameServer) getMapByActorPID(pid *actor.PID) *entity.Map {
	if gs == nil || pid == nil {
		return nil
	}
	key := pid.String()
	gs.mapsMutex.RLock()
	defer gs.mapsMutex.RUnlock()
	for _, m := range gs.maps {
		if m == nil {
			continue
		}
		actorPID := m.GetActorPID()
		if actorPID != nil && actorPID.String() == key {
			return m
		}
	}
	return nil
}

func (gs *GameServer) handleClientDisconnect(c core.Client) {
	client, ok := c.(*client.GameClient)
	if !ok {
		return
	}
	character := client.GetCharacter()
	if character == nil {
		return
	}

	p := async.NewPromise(nil, saveCharactersPromiseTimeout)
	p.OnError(func(err error) {
		log.Printf("disconnect async: %v", err)
	})

	if gs.internalClient != nil && character.AccountID != 0 {
		transfer := client.TakeTransferDisconnect()
		accID := character.AccountID
		wid := gs.config.WorldId
		p.Then(func() (interface{}, error) {
			ctx, cancel := context.WithTimeout(context.Background(), core.InternalRPCPerStepTimeout)
			defer cancel()
			req := &internal.LogoutSessionRequest{
				WorldId:            wid,
				AccountId:          accID,
				DisconnectSource:   internal.SessionDisconnectSource_SESSION_DISCONNECT_SOURCE_GAME_SERVER,
				TransferDisconnect: transfer,
			}
			if cid := character.GetID(); cid != 0 {
				v := cid
				req.CharacterId = &v
			}
			ch := gs.config.ChannelId
			req.ChannelId = &ch
			_, err := gs.internalClient.LogoutSession(ctx, req)
			return err, nil
		}, func(v interface{}) error {
			if err, _ := v.(error); err != nil {
				log.Printf("LogoutSession (game disconnect) failed for account %d: %v", accID, err)
			}
			return nil
		})
	}

	toSave := []*entity.Character{character}
	async.ThenRPC(p, func(c context.Context) (*internal.SaveCharactersReply, error) {
		return gs.grpcSaveCharacters(c, toSave)
	}, func(*internal.SaveCharactersReply) error {
		return nil
	})

	charID := character.GetID()
	p.Finally(func() {
		gs.runCharacterLogoutScript(character)
		mapInstance := character.GetMap()
		if mapInstance != nil {
			pid := mapInstance.GetActorPID()
			root := gs.GetRootContext()
			if pid != nil && root != nil {
				_, err := root.RequestFuture(pid, &g_actor.RemoveCharacter{CharacterID: charID}, mapActorCallTimeout).Result()
				if err != nil {
					log.Printf("RemoveCharacter on disconnect (char %d): %v", charID, err)
				}
			} else {
				log.Printf("disconnect: map actor missing for char %d; map remove skipped", charID)
			}
		}
		if gs.characterRuntime != nil {
			gs.ensureAbandonCharacter(charID)
			gs.characterRuntime.UnregisterCharacter(charID)
		}
		character.ClearTimers()
	})

	p.Run()
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
		stats["drops_loaded"] = len(gs.resources.Drops)
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
