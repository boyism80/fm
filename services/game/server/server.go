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
	"github.com/boyism80/fm/core/luax"
	"github.com/boyism80/fm/core/mq"
	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
	g_actor "github.com/boyism80/fm/services/game/actor"
	"github.com/boyism80/fm/services/game/client"
	gameconst "github.com/boyism80/fm/services/game/constant"
	"github.com/boyism80/fm/services/game/entity"
	"github.com/boyism80/fm/services/game/wz"
	"github.com/boyism80/fm/types"
	lua "github.com/yuin/gopher-lua"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const mapActorCallTimeout = 30 * time.Second

type GameServer struct {
	server            *core.Server
	config            *GameConfig
	resources         *wz.Resources
	maps              map[uint32]*entity.Map
	mapsMutex         sync.RWMutex
	packetHandlers    *PacketHandlerRegistry
	context           *GameServerContext
	actorSystem       *c_actor.ActorSystem
	actorRegistry     *c_actor.ActorRegistry
	nilMapActorPID    *actor.PID
	characterListener entity.CharacterListener
	internalClient    internal.InternalClient
	internalConn      *grpc.ClientConn
	party             *PartyContainer
	consumerParty     *mq.JSONConsumer[*GameServer]
	characterRuntime  *ServerCharacterRuntime
	ensureMu          sync.Mutex
	ensurePending     map[uint64]*g_actor.EnsureDeliver
	ensureNext        atomic.Uint64
}

func (gs *GameServer) GetServer() *core.Server {
	return gs.server
}

func (gs *GameServer) GetServerContext() core.ServerContext {
	return gs.context
}

func (gs *GameServer) GetRootContext() *actor.RootContext {
	return gs.server.GetRootContext()
}

type GameContext interface {
	GetResources() *wz.Resources
	GetMap(mapId uint32) *entity.Map
	GetExpRate() int
	GetDropRate() int
	GetMesoRate() int
	RequestWarp(character *entity.Character, targetMap *entity.Map, spawnPoint uint8) error
}

type GameConfig struct {
	Host         string
	Port         int
	ChannelId    uint32
	WzPath       string
	WorldName    string
	WorldId      uint32
	MaxPlayers   int
	ExpRate      int
	DropRate     int
	MesoRate     int
	InternalAddr string
	RabbitMQ     config.RabbitMQEndpoint
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

	context := NewGameServerContext(config.WzPath, nil)

	actorSystem := c_actor.NewActorSystem()
	actorRegistry := c_actor.NewActorRegistry(actorSystem)

	server.SetRootContext(actorSystem.GetRoot())

	gs := &GameServer{
		server:           server,
		config:           config,
		resources:        resources,
		maps:             make(map[uint32]*entity.Map),
		context:          context,
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
		consumerParty := mq.NewJSONConsumer(config.RabbitMQ, gs, queueName, consumerTag)
		consumerParty.Route(routeAll).Route(routeGame)

		mq.BindOn[*GameServer, partyMqCreated](consumerParty)
		mq.BindOn[*GameServer, partyMqMemberJoined](consumerParty)
		mq.BindOn[*GameServer, partyMqMemberLeft](consumerParty)
		mq.BindOn[*GameServer, partyMqLeaderChanged](consumerParty)
		mq.BindOn[*GameServer, partyMqLogOnOff](consumerParty)
		mq.BindOn[*GameServer, partyMqDisbanded](consumerParty)
		mq.BindOn[*GameServer, partyMqPartySnapshot](consumerParty)
		mq.BindOn[*GameServer, partyMqPartyInvite](consumerParty)
		mq.BindOn[*GameServer, partyMqPartyInviteDenied](consumerParty)
		gs.consumerParty = consumerParty
	}

	gs.characterListener = &CharacterListenerImpl{gs: gs}
	context.gs = gs
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
			MapData:        nil,
			Context:        gs.context,
			SaveCharacters: gs.saveCharactersChunked,
			Ensure:         gs,
		}
	})

	nilMapPID := gs.actorRegistry.GetOrCreateActor(
		"map_nil",
		nilMapProps,
	)
	gs.nilMapActorPID = nilMapPID

	gs.server.SetNilMapActorPID(nilMapPID)

	log.Println("Pre-creating map instances...")
	for mapID := range gs.resources.Maps {
		mapInstance := entity.NewMap(mapID, gameMapListener, mapID, gs)
		gs.maps[mapID] = mapInstance

		props := actor.PropsFromProducer(func() actor.Actor {
			return &g_actor.MapActor{
				MapData:        mapInstance,
				Context:        gs.context,
				SaveCharacters: gs.saveCharactersChunked,
				Ensure:         gs,
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

	if err := gs.server.Start(gs.config.Host, gs.config.Port); err != nil {
		return err
	}

	log.Printf("Game server started on %s:%d", gs.config.Host, gs.config.Port)
	log.Printf("World: %s, Max Players: %d", gs.config.WorldName, gs.config.MaxPlayers)
	log.Printf("Rates: Exp=%dx, Drop=%dx, Meso=%dx",
		gs.config.ExpRate, gs.config.DropRate, gs.config.MesoRate)
	if gs.consumerParty != nil {
		if err := gs.consumerParty.Start(); err != nil {
			return fmt.Errorf("party mq start failed: %w", err)
		}
		log.Printf("Party MQ started (%s)", gs.consumerParty.QueueName())
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
	if gs.consumerParty != nil {
		_ = gs.consumerParty.Close()
	}
	if gs.internalConn != nil {
		_ = gs.internalConn.Close()
	}
	return gs.server.Stop()
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
	var pid *actor.PID
	if m := ch.GetMap(); m != nil {
		pid = m.GetActorPID()
	}
	if pid == nil && gs.nilMapActorPID != nil {
		pid = gs.nilMapActorPID
	}
	if pid == nil {
		return
	}
	root := luax.GetRootLuaState(pid.String())
	if root == nil {
		return
	}
	_, thread, err := luax.Call(root, "script/script.lua", "on_logout", ch)
	if err != nil {
		log.Printf("on_logout: %v", err)
	}
	if thread != nil {
		thread.Close()
	}
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
	if gs.internalClient != nil && character.AccountID != 0 {
		ctx, cancel := context.WithTimeout(context.Background(), core.InternalRPCPerStepTimeout)
		transfer := client.TakeTransferDisconnect()
		_, err := gs.internalClient.LogoutSession(ctx, &internal.LogoutSessionRequest{
			WorldId:            gs.config.WorldId,
			AccountId:          character.AccountID,
			DisconnectSource:   internal.SessionDisconnectSource_SESSION_DISCONNECT_SOURCE_GAME_SERVER,
			TransferDisconnect: transfer,
		})
		cancel()
		if err != nil {
			log.Printf("LogoutSession (game disconnect) failed for account %d: %v", character.AccountID, err)
		}
	}
	gs.saveCharacterAsync(character)
	gs.runCharacterLogoutScript(character)
	mapInstance := character.GetMap()
	if mapInstance != nil {
		pid := mapInstance.GetActorPID()
		root := gs.GetRootContext()
		if pid != nil && root != nil {
			_, err := root.RequestFuture(pid, &g_actor.RemoveCharacter{CharacterID: character.GetID()}, mapActorCallTimeout).Result()
			if err != nil {
				log.Printf("RemoveCharacter on disconnect (char %d): %v", character.GetID(), err)
			}
		} else {
			log.Printf("disconnect: map actor missing for char %d; map remove skipped", character.GetID())
		}
	}
	if gs.characterRuntime != nil {
		gs.ensureAbandonCharacter(character.GetID())
		gs.characterRuntime.UnregisterCharacter(character.GetID())
	}
	character.ClearTimers()
}

func (gs *GameServer) GetStats() map[string]interface{} {
	stats := gs.server.GetStats()

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

		Host:       "0.0.0.0",
		Port:       8485,
		WorldName:  "Scania",
		MaxPlayers: 1000,
		ExpRate:    1,
		DropRate:   1,
		MesoRate:   1,
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

		Host:       "0.0.0.0",
		Port:       8485,
		WorldName:  "HighRate",
		MaxPlayers: 2000,
		ExpRate:    10,
		DropRate:   5,
		MesoRate:   5,
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
