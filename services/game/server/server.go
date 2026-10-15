package server

import (
	"context"
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/common/config"
	"github.com/boyism80/fm/core"
	c_actor "github.com/boyism80/fm/core/actor"
	"github.com/boyism80/fm/core/async"
	"github.com/boyism80/fm/core/luax"
	"github.com/boyism80/fm/core/mq"
	"github.com/boyism80/fm/core/rpc"
	"github.com/boyism80/fm/core/system"
	"github.com/boyism80/fm/core/uniqueid"
	"github.com/boyism80/fm/maple/wz"
	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
	"github.com/boyism80/fm/services/game"
	g_actor "github.com/boyism80/fm/services/game/actor"
	"github.com/boyism80/fm/services/game/client"
	lua "github.com/yuin/gopher-lua"
	"google.golang.org/grpc"
)

const (
	logoutWaitTimeout   = game.EnvelopeTimeout + time.Second
	shutdownWaitTimeout = time.Minute
)

type GameServer struct {
	*core.ServerCore
	config      *GameConfig
	systems     *system.Systems
	listeners   *system.Listeners
	disconnects sync.WaitGroup
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
	Dev                              bool
	LuaAlwaysReload                  bool
	DueyIdentityPrompt               bool
}

func NewGameServer(config *GameConfig) (*GameServer, error) {
	serverConfig := &core.ServerConfig{
		Host:         config.Host,
		Port:         config.Port,
		EnterTimeout: 30 * time.Second,
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
	luax.SetAlwaysReload(config.LuaAlwaysReload)

	gs := &GameServer{
		ServerCore: server,
		config:     config,
		systems:    system.NewSystems(),
		listeners:  system.NewListeners(),
	}
	gs.systems.Add[*wz.Resources](resources)
	gs.systems.Add[*uniqueid.Generator](uniqueid.NewForChannel(uint32(config.WorldId), uint32(config.ChannelId)))

	gs.listeners.Add[game.CharacterListener](&CharacterListenerImpl{gs: gs})
	gs.listeners.Add[game.MapListener](NewGameMapListener(gs))
	gs.listeners.Add[game.MobListener](&MobListenerImpl{})

	// Systems start in the order they are added and stop in reverse: the actor system stops after every
	// system that uses it, and the MQ consumers stop first so no event arrives after the systems it uses have stopped.
	actors := c_actor.NewActorRegistry(c_actor.NewActorSystem())
	gs.systems.Add[*c_actor.ActorRegistry](actors)
	server.SetRootContext(actors.Root())

	conn, err := rpc.NewInternalConnection(config.InternalAddr, time.Duration(config.InternalHeartbeatIntervalSeconds)*time.Second, gs.ping)
	if err != nil {
		return nil, err
	}
	gs.systems.Add[*rpc.InternalConnection](conn)
	gs.systems.Add[game.ActorFactory](g_actor.NewFactory(gs))
	gs.systems.Add[*game.InternalSystem](game.NewInternalSystem(gs, conn))
	gs.systems.Add[*game.DispatchSystem](game.NewDispatchSystem(gs))
	maps := game.NewMapSystem(gs, server, gs.listeners.Get[game.MapListener](), gs.listeners.Get[game.MobListener]())
	gs.systems.Add[*game.MapSystem](maps)
	gs.systems.Add[*game.SchedulerSystem](game.NewSchedulerSystem(gs))
	stateMachines := game.NewStateMachineSystem(gs)
	gs.systems.Add[*game.StateMachineSystem](stateMachines)
	gs.systems.Add[*game.PartySystem](game.NewPartySystem(gs))
	gs.systems.Add[*game.GuildSystem](game.NewGuildSystem(gs))
	gs.systems.Add[*game.AllianceSystem](game.NewAllianceSystem(gs))
	gs.systems.Add[*game.CarnivalSystem](game.NewCarnivalSystem())
	gs.systems.Add[*game.ExpeditionSystem](game.NewExpeditionSystem(gs))
	gs.systems.Add[*game.NoticeSystem](game.NewNoticeSystem(server))
	gs.systems.Add[*mq.Consumers](mq.NewConsumers(actors, config.RabbitMQ, gs.mqQueues()))

	luax.RegisterOnCreateHook(func(luaState *lua.LState) {
		gs.registerGameLuaState(luaState)
	})
	gs.runStartupScript()

	maps.PreCreate()
	if err := stateMachines.LoadFromScripts("script/state_machine"); err != nil {
		log.Printf("state machine load: %v", err)
	}

	gs.registerPacketHandlers()

	server.SetOnClientDisconnect(func(client core.Client) {
		gs.handleClientDisconnect(client)
	})

	return gs, nil
}

func (gs *GameServer) Systems() *system.Systems {
	return gs.systems
}

func (gs *GameServer) GetResources() *wz.Resources {
	return gs.systems.Get[*wz.Resources]()
}

func (gs *GameServer) GetWorldID() uint32 {
	return gs.config.WorldId
}

func (gs *GameServer) GetChannelID() int32 {
	return int32(gs.config.ChannelId)
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
	return gs.systems.Get[*uniqueid.Generator]().Next()
}

func (gs *GameServer) Start() error {
	log.Println("Starting MapleStory Game Server...")

	if err := gs.systems.Start(); err != nil {
		return err
	}
	if err := gs.ServerCore.Start(gs.config.Host, gs.config.Port); err != nil {
		return err
	}

	log.Printf("Game server started on %s:%d", gs.config.Host, gs.config.Port)
	log.Printf("World: %s, Max Players: %d", gs.config.WorldName, gs.config.MaxPlayers)
	log.Printf("Rates: Exp=%dx, Drop=%dx, Meso=%dx",
		gs.config.ExpRate, gs.config.DropRate, gs.config.MesoRate)
	resources := gs.systems.Get[*wz.Resources]()
	log.Printf("Resources loaded: %d maps, %d monsters, %d items, %d mob drops, %d global drops, %d reactor drops, %d strings",
		len(resources.Maps), len(resources.Monsters),
		len(resources.Items), len(resources.MobDrops), len(resources.GlobalDrops), len(resources.ReactorDrops), resources.Strings.CountStrings())
	return nil
}

func (gs *GameServer) Stop() error {
	log.Println("Shutting down game server...")

	// Closing the clients saves and logs out every character through internal, so the systems stay up
	// until those disconnects finish.
	err := gs.ServerCore.Stop()
	done := make(chan struct{})
	go func() {
		gs.disconnects.Wait()
		close(done)
	}()
	select {
	case <-done:
	case <-time.After(shutdownWaitTimeout):
		log.Printf("shutdown: disconnect saves did not finish within %v", shutdownWaitTimeout)
	}

	gs.systems.Stop()
	return err
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

	gs.disconnects.Add(1)
	p := async.NewTask(nil, game.SaveCharactersPromiseTimeout)
	p.OnError(func(err error) {
		log.Printf("disconnect async: %v", err)
	}).Finally(gs.disconnects.Done)

	var entry *internal.CharacterSaveEntry
	p.DoAsync(func() error {
		entry = gs.removeCharacter(character, logout)
		return nil
	})

	if !client.SessionLost() {
		p.ThenRPC(func(c context.Context) (*internal.SaveCharactersReply, error) {
			if entry == nil {
				return &internal.SaveCharactersReply{Ok: true}, nil
			}
			return gs.systems.Get[*rpc.InternalConnection]().Client.SaveCharacters(c, &internal.SaveCharactersRequest{Entries: []*internal.CharacterSaveEntry{entry}}, grpc.WaitForReady(true))
		}, func(reply *internal.SaveCharactersReply) error {
			if !reply.GetOk() {
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
		return gs.systems.Get[*rpc.InternalConnection]().Client.LogoutSession(c, &internal.LogoutSessionRequest{
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
func (gs *GameServer) removeCharacter(character *game.Character, logout <-chan *internal.CharacterSaveEntry) *internal.CharacterSaveEntry {
	charID := character.GetID()
	gs.systems.Get[*game.ExpeditionSystem]().LeaveChannel(character)
	if sm := character.StateMachine(); sm != nil {
		sm.RequestLeave(character, false, game.StateMachineLeaveDisconnect)
	}
	dispatch := gs.systems.Get[*game.DispatchSystem]()
	dispatch.Call(charID, func(ctx actor.Context) {
		if m := character.GetMap(); m != nil {
			_ = m.LogoutPlayer(ctx, charID)
		}
	})

	var entry *internal.CharacterSaveEntry
	select {
	case entry = <-logout:
	case <-time.After(logoutWaitTimeout):
		log.Printf("remove character %d: no map actor owns it; clearing timers here", charID)
		entry = character.BuildRPC(gs.config.WorldId)
		character.ClearTimers()
	}

	dispatch.AbandonEnvelopes(charID)
	dispatch.UnregisterCharacter(charID)
	return entry
}
