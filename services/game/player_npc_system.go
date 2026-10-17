package game

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/core/async"
	"github.com/boyism80/fm/core/rpc"
	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
	"google.golang.org/grpc"
)

// PlayerNpcSystem keeps every player npc of the world; Lua on map actors and the MQ consumer both reach it.
type PlayerNpcSystem struct {
	gw     GameWorld
	client internal.InternalClient
	mu     sync.Mutex
	npcs   map[uint32]*PlayerNpc
}

func NewPlayerNpcSystem(gw GameWorld, conn *rpc.InternalConnection) *PlayerNpcSystem {
	return &PlayerNpcSystem{gw: gw, client: conn.Client, npcs: make(map[uint32]*PlayerNpc)}
}

func (s *PlayerNpcSystem) Start() error {
	// A channel can boot before internal listens; the npcs would stay missing until the next restart.
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	reply, err := s.client.ListPlayerNpcs(ctx, &internal.ListPlayerNpcsRequest{WorldId: s.gw.GetWorldID()}, grpc.WaitForReady(true))
	cancel()
	if err != nil {
		log.Printf("list player npcs: %v", err)
		return nil
	}
	for _, pb := range reply.GetNpcs() {
		s.Place(NewPlayerNpcFromInternalProto(pb))
	}
	log.Printf("Placed %d player npcs", len(reply.GetNpcs()))
	return nil
}

func (s *PlayerNpcSystem) Stop() {}

func (s *PlayerNpcSystem) Find(npcID uint32) *PlayerNpc {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.npcs[npcID]
}

func (s *PlayerNpcSystem) CreateAsync(actorCtx actor.Context, npc *PlayerNpc) *async.Promise[bool] {
	return async.NewTask(actorCtx, core.InternalRPCPerStepTimeout).ThenRPC(
		func(c context.Context) (*internal.CreatePlayerNpcReply, error) {
			return s.client.CreatePlayerNpc(c, &internal.CreatePlayerNpcRequest{WorldId: s.gw.GetWorldID(), Npc: npc.BuildRPC()})
		},
		nil,
	).Then(func(reply *internal.CreatePlayerNpcReply) (bool, error) {
		if reply.GetCreated() {
			s.Place(npc)
		}
		return reply.GetCreated(), nil
	})
}

func (s *PlayerNpcSystem) RemoveAsync(actorCtx actor.Context, npcID uint32) *async.Promise[bool] {
	return async.NewTask(actorCtx, core.InternalRPCPerStepTimeout).ThenRPC(
		func(c context.Context) (*internal.RemovePlayerNpcReply, error) {
			return s.client.RemovePlayerNpc(c, &internal.RemovePlayerNpcRequest{WorldId: s.gw.GetWorldID(), NpcId: npcID})
		},
		nil,
	).Then(func(reply *internal.RemovePlayerNpcReply) (bool, error) {
		if reply.GetRemoved() {
			s.Remove(npcID)
		}
		return reply.GetRemoved(), nil
	})
}

// Place and Remove are idempotent: the channel that made the change also receives its MQ notice.
func (s *PlayerNpcSystem) Place(npc *PlayerNpc) {
	s.mu.Lock()
	if s.npcs[npc.NpcID] != nil {
		s.mu.Unlock()
		return
	}
	s.npcs[npc.NpcID] = npc
	s.mu.Unlock()

	maps := s.gw.Systems().Get[*MapSystem]()
	m := maps.Get(npc.MapID)
	if m == nil {
		return
	}
	maps.Visit(m, func(actor.Context) {
		m.PlacePlayerNpc(npc)
	})
}

func (s *PlayerNpcSystem) Remove(npcID uint32) {
	s.mu.Lock()
	npc := s.npcs[npcID]
	if npc == nil {
		s.mu.Unlock()
		return
	}
	delete(s.npcs, npcID)
	s.mu.Unlock()

	maps := s.gw.Systems().Get[*MapSystem]()
	m := maps.Get(npc.MapID)
	if m == nil {
		return
	}
	maps.Visit(m, func(actor.Context) {
		m.RemovePlayerNpc(npcID)
	})
}
