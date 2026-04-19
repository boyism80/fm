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

	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/common/config"
	"github.com/boyism80/fm/core"
	c_actor "github.com/boyism80/fm/core/actor"
	"github.com/boyism80/fm/core/luax"
	"github.com/boyism80/fm/protocol/constant"
	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
	"github.com/boyism80/fm/protocol/response"
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

type GameServer struct {
	server             *core.Server
	config             *GameConfig
	resources          *wz.Resources
	maps               map[uint32]*entity.Map
	mapsMutex          sync.RWMutex
	packetHandlers     *PacketHandlerRegistry
	context            *GameServerContext
	actorSystem        *c_actor.ActorSystem
	actorRegistry      *c_actor.ActorRegistry
	nilMapActorPID     *actor.PID
	characterListener  entity.CharacterListener
	internalClient     internal.InternalClient
	internalConn       *grpc.ClientConn
	partyEventConsumer *PartyEventConsumer

	characterRuntime *ServerCharacterRuntime

	ensureMu      sync.Mutex
	ensurePending map[uint64]*g_actor.EnsureDeliver
	ensureNext    atomic.Uint64
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

	if config.RabbitMQ.Enabled() {
		gs.partyEventConsumer = NewPartyEventConsumer(
			config.RabbitMQ,
			config.WorldId,
			config.ChannelId,
			gs.internalClient,
			gs.SyncPartySnapshot,
			gs.ClearPartyMembers,
			gs.DeliverPartyJoinUpdate,
			gs.DeliverPartyLeaveUpdate,
			gs.DeliverPartyDisbandUpdate,
			gs.DeliverPartyLeaderChange,
			gs.DeliverPartyLogOnOff,
			gs.DeliverPartySilentFromSnapshot,
			gs.DeliverPartyInviteToCharacter,
			gs.DeliverPartyDenyStatusToCharacter,
		)
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
	if gs.partyEventConsumer != nil {
		if err := gs.partyEventConsumer.Start(); err != nil {
			return fmt.Errorf("party event consumer start failed: %w", err)
		}
		log.Printf("Party event consumer started (%s)", gs.partyEventConsumer.QueueName())
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
	if gs.partyEventConsumer != nil {
		_ = gs.partyEventConsumer.Close()
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

func (gs *GameServer) DeliverPartyInviteToCharacter(targetCharacterID uint32, partyID uint32, inviterName string, partySearch bool) {
	if gs == nil {
		return
	}
	gs.EnsureSend(nil, targetCharacterID, &g_actor.DeliverPartyInvite{
		CharacterID: targetCharacterID,
		PartyID:     partyID,
		InviterName: inviterName,
		PartySearch: partySearch,
	})
}

func (gs *GameServer) DeliverPartyDenyStatusToCharacter(targetCharacterID uint32, action uint8, deniedCharacterName string) {
	if gs == nil || targetCharacterID == 0 {
		return
	}
	gs.EnsureSend(nil, targetCharacterID, &g_actor.DeliverPartyStatusMessage{
		CharacterID: targetCharacterID,
		Code:        constant.PartyStatusCode(action),
		Name:        deniedCharacterName,
	})
}

func (gs *GameServer) SyncCharacterPartyState(characterID uint32, partyID *uint32) {
	if gs == nil || characterID == 0 {
		return
	}

	gs.EnsureSend(nil, characterID, &g_actor.SyncCharacterPartyState{
		CharacterID: characterID,
		PartyID:     partyID,
	})
}

func (gs *GameServer) SyncPartySnapshot(snapshot *internal.PartySnapshot) {
	if gs == nil || snapshot == nil {
		return
	}
	partyID := snapshot.GetPartyId()
	for _, member := range snapshot.GetMembers() {
		cid := member.GetCharacterId()
		pid := partyID
		gs.SyncCharacterPartyState(cid, &pid)
	}
}

func (gs *GameServer) ClearCharacterPartyID(characterID uint32) {
	if gs == nil || characterID == 0 {
		return
	}
	gs.SyncCharacterPartyState(characterID, nil)
}

func (gs *GameServer) ClearPartyMembers(memberIDs []uint32) {
	if gs == nil || len(memberIDs) == 0 {
		return
	}
	for _, cid := range memberIDs {
		gs.ClearCharacterPartyID(cid)
	}
}

func partyMembersToResponse(members []*internal.PartyMemberSnapshot) []response.PartyMemberStatus {
	out := make([]response.PartyMemberStatus, 0, len(members))
	for _, m := range members {
		if m == nil {
			continue
		}
		ch := int32(-2)
		if m.ChannelIndex != nil {
			ch = *m.ChannelIndex
		}
		doorTown := uint32(999999999)
		doorTarget := uint32(999999999)
		doorX := int32(0)
		doorY := int32(0)
		if d := m.GetDoor(); d != nil {
			doorTown = d.GetTown()
			doorTarget = d.GetTarget()
			doorX = d.GetX()
			doorY = d.GetY()
		}
		out = append(out, response.PartyMemberStatus{
			CharacterID: m.GetCharacterId(),
			Name:        m.GetCharacterName(),
			Class:       m.GetClassId(),
			Level:       m.GetLevel(),
			Channel:     ch,
			MapID:       m.GetMapId(),
			DoorTown:    doorTown,
			DoorTarget:  doorTarget,
			DoorX:       doorX,
			DoorY:       doorY,
		})
	}
	return out
}

func (gs *GameServer) DeliverPartyJoinUpdate(snapshot *internal.PartySnapshot, joinedCharacterID uint32) {
	if gs == nil || snapshot == nil || joinedCharacterID == 0 {
		return
	}
	members := snapshot.GetMembers()
	if len(members) == 0 {
		return
	}
	joinName := ""
	for _, m := range members {
		if m != nil && m.GetCharacterId() == joinedCharacterID {
			joinName = m.GetCharacterName()
			break
		}
	}
	respMembers := partyMembersToResponse(members)
	if joinName == "" {
		return
	}
	for _, m := range members {
		if m == nil || m.GetCharacterId() == 0 {
			continue
		}
		gs.EnsureSend(nil, m.GetCharacterId(), &g_actor.DeliverPartyUpdateJoin{
			CharacterID: m.GetCharacterId(),
			ForChannel:  int32(gs.config.ChannelId),
			PartyID:     snapshot.GetPartyId(),
			JoinName:    joinName,
			LeaderID:    snapshot.GetLeaderCharacterId(),
			Members:     respMembers,
		})
	}
}

func (gs *GameServer) DeliverPartyLeaveUpdate(prev, current *internal.PartySnapshot, targetCharacterID uint32, expelled bool) {
	if gs == nil || prev == nil || targetCharacterID == 0 {
		return
	}
	targetName := ""
	oldMembers := prev.GetMembers()
	for _, m := range oldMembers {
		if m != nil && m.GetCharacterId() == targetCharacterID {
			targetName = m.GetCharacterName()
			break
		}
	}
	if targetName == "" {
		return
	}
	var (
		partyID  = prev.GetPartyId()
		leaderID = prev.GetLeaderCharacterId()
		members  []*internal.PartyMemberSnapshot
	)
	if current != nil {
		partyID = current.GetPartyId()
		leaderID = current.GetLeaderCharacterId()
		members = current.GetMembers()
	}
	respMembers := partyMembersToResponse(members)
	for _, m := range oldMembers {
		if m == nil || m.GetCharacterId() == 0 {
			continue
		}
		gs.EnsureSend(nil, m.GetCharacterId(), &g_actor.DeliverPartyUpdateLeave{
			CharacterID: m.GetCharacterId(),
			ForChannel:  int32(gs.config.ChannelId),
			PartyID:     partyID,
			TargetID:    targetCharacterID,
			TargetName:  targetName,
			LeaderID:    leaderID,
			Members:     respMembers,
			Expelled:    expelled,
		})
	}
}

func (gs *GameServer) DeliverPartyDisbandUpdate(prev *internal.PartySnapshot, leaderCharacterID uint32) {
	if gs == nil || prev == nil || leaderCharacterID == 0 {
		return
	}
	for _, m := range prev.GetMembers() {
		if m == nil || m.GetCharacterId() == 0 {
			continue
		}
		gs.EnsureSend(nil, m.GetCharacterId(), &g_actor.DeliverPartyUpdateDisband{
			CharacterID: m.GetCharacterId(),
			PartyID:     prev.GetPartyId(),
			LeaderID:    leaderCharacterID,
		})
	}
}

func (gs *GameServer) DeliverPartyLeaderChange(snapshot *internal.PartySnapshot, newLeaderCharacterID uint32, byDisconnect bool) {
	if gs == nil || snapshot == nil || newLeaderCharacterID == 0 {
		return
	}
	for _, m := range snapshot.GetMembers() {
		if m == nil || m.GetCharacterId() == 0 {
			continue
		}
		gs.EnsureSend(nil, m.GetCharacterId(), &g_actor.DeliverPartyUpdateLeaderChange{
			CharacterID:          m.GetCharacterId(),
			NewLeaderCharacterID: newLeaderCharacterID,
			ByDisconnect:         byDisconnect,
		})
	}
}

func (gs *GameServer) DeliverPartyLogOnOff(snapshot *internal.PartySnapshot, _ uint32) {
	if gs == nil || snapshot == nil {
		return
	}
	respMembers := partyMembersToResponse(snapshot.GetMembers())
	for _, m := range snapshot.GetMembers() {
		if m == nil || m.GetCharacterId() == 0 {
			continue
		}
		gs.EnsureSend(nil, m.GetCharacterId(), &g_actor.DeliverPartyUpdateLogOnOff{
			CharacterID: m.GetCharacterId(),
			ForChannel:  int32(gs.config.ChannelId),
			PartyID:     snapshot.GetPartyId(),
			LeaderID:    snapshot.GetLeaderCharacterId(),
			Members:     respMembers,
		})
	}
}

func (gs *GameServer) DeliverPartySilentFromSnapshot(snapshot *internal.PartySnapshot) {
	if gs == nil || snapshot == nil {
		return
	}
	respMembers := partyMembersToResponse(snapshot.GetMembers())
	for _, m := range snapshot.GetMembers() {
		if m == nil || m.GetCharacterId() == 0 {
			continue
		}
		gs.EnsureSend(nil, m.GetCharacterId(), &g_actor.DeliverPartyUpdateSilent{
			CharacterID: m.GetCharacterId(),
			ForChannel:  int32(gs.config.ChannelId),
			PartyID:     snapshot.GetPartyId(),
			LeaderID:    snapshot.GetLeaderCharacterId(),
			Members:     respMembers,
		})
	}
}

func (gs *GameServer) partySnapshotForMapEnter(partyID uint32) *internal.PartySnapshot {
	if gs.partyEventConsumer != nil {
		if s := gs.partyEventConsumer.CachedPartySnapshot(partyID); s != nil {
			return s
		}
	}
	if gs.internalClient == nil {
		return nil
	}
	ctx, cancel := context.WithTimeout(context.Background(), core.InternalRPCPerStepTimeout)
	defer cancel()
	reply, err := gs.internalClient.GetParty(ctx, &internal.GetPartyRequest{
		WorldId: gs.config.WorldId,
		PartyId: partyID,
	})
	if err != nil || !reply.GetFound() || reply.GetParty() == nil {
		return nil
	}
	return reply.GetParty()
}

func (gs *GameServer) SendPartySilentOnMapEnter(ch *entity.Character) {
	if gs == nil || ch == nil {
		return
	}
	partyIDPtr := ch.GetPartyID()
	if partyIDPtr == nil {
		return
	}
	snap := gs.partySnapshotForMapEnter(*partyIDPtr)
	if snap == nil {
		return
	}
	members := partyMembersToResponse(snap.GetMembers())
	gs.EnsureSend(nil, ch.GetID(), &g_actor.DeliverPartyUpdateSilent{
		CharacterID: ch.GetID(),
		ForChannel:  int32(gs.config.ChannelId),
		PartyID:     snap.GetPartyId(),
		LeaderID:    snap.GetLeaderCharacterId(),
		Members:     members,
	})
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
	slot := gs.PartyMemberIndex(ch.GetID(), ch.GetPartyID())
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

func (gs *GameServer) PartyMemberIndex(characterID uint32, partyID *uint32) int {
	if partyID == nil || *partyID == 0 || gs.partyEventConsumer == nil {
		return 0
	}
	snap := gs.partyEventConsumer.CachedPartySnapshot(*partyID)
	if snap == nil {
		return 0
	}
	for i, mem := range snap.GetMembers() {
		if mem.GetCharacterId() == characterID {
			return i
		}
	}
	return 0
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
		_, err := gs.internalClient.LogoutSession(ctx, &internal.LogoutSessionRequest{
			WorldId:   gs.config.WorldId,
			AccountId: character.AccountID,
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
		mapInstance.RemovePlayer(character.GetID())
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
