package game

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/core/async"
	"github.com/boyism80/fm/core/clock"
	"github.com/boyism80/fm/core/luax"
	"github.com/boyism80/fm/core/rpc"
	"github.com/boyism80/fm/maple/constant"
	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
	"google.golang.org/grpc"
)

type MedalRankingEntry struct {
	CharacterID uint32
	Name        string
	Score       uint64
}

type MedalRanking struct {
	Kind    string
	Month   uint32
	Entries []MedalRankingEntry
}

// MedalSystem keeps a copy of the boards for display; a medal is given or taken back only on internal's answer.
type MedalSystem struct {
	gw     GameWorld
	client internal.InternalClient
	mu     sync.Mutex
	boards map[string]*MedalRanking
	// A kind missing from the last full list had no board in this season.
	listedMonth uint32
}

func NewMedalSystem(gw GameWorld, conn *rpc.InternalConnection) *MedalSystem {
	return &MedalSystem{gw: gw, client: conn.Client, boards: make(map[string]*MedalRanking)}
}

func (s *MedalSystem) Start() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	reply, err := s.client.ListMedalRankings(ctx, &internal.ListMedalRankingsRequest{WorldId: s.gw.GetWorldID()}, grpc.WaitForReady(true))
	cancel()
	if err != nil {
		log.Printf("list medal rankings: %v", err)
		return nil
	}
	s.list(reply.GetRankings())
	return nil
}

func (s *MedalSystem) Stop() {}

func (s *MedalSystem) month() uint32 {
	now := clock.Now().In(time.Local)
	return uint32(now.Year()*100 + int(now.Month()))
}

func (s *MedalSystem) Entries(kind string) ([]MedalRankingEntry, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	month := s.month()
	board := s.boards[kind]
	if board != nil && board.Month == month {
		return board.Entries, true
	}
	return nil, s.listedMonth == month
}

func (s *MedalSystem) SubmitAsync(actorCtx actor.Context, ch *Character, kind string, score uint64) *async.Promise[uint32] {
	req := &internal.SubmitMedalScoreRequest{
		WorldId:     s.gw.GetWorldID(),
		Kind:        kind,
		Month:       s.month(),
		CharacterId: ch.GetID(),
		Name:        ch.GetName(),
		Score:       score,
	}
	return async.NewTask(actorCtx, core.InternalRPCPerStepTimeout).ThenRPC(
		func(c context.Context) (*internal.SubmitMedalScoreReply, error) {
			return s.client.SubmitMedalScore(c, req)
		},
		nil,
	).Then(func(reply *internal.SubmitMedalScoreReply) (uint32, error) {
		s.Set(NewMedalRankingFromInternalProto(reply.GetRanking()))
		return reply.GetRank(), nil
	})
}

func (s *MedalSystem) RefreshAsync(actorCtx actor.Context) *async.Promise[bool] {
	req := &internal.ListMedalRankingsRequest{WorldId: s.gw.GetWorldID()}
	return async.NewTask(actorCtx, core.InternalRPCPerStepTimeout).ThenRPC(
		func(c context.Context) (*internal.ListMedalRankingsReply, error) {
			return s.client.ListMedalRankings(c, req)
		},
		nil,
	).Then(func(reply *internal.ListMedalRankingsReply) (bool, error) {
		s.list(reply.GetRankings())
		return true, nil
	})
}

// Set is idempotent: the channel that submitted the score also receives its MQ notice.
func (s *MedalSystem) Set(ranking *MedalRanking) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if current := s.boards[ranking.Kind]; current != nil && current.Month > ranking.Month {
		return
	}
	s.boards[ranking.Kind] = ranking
}

// Internal keeps earlier seasons until a new submit prunes them.
func (s *MedalSystem) list(rankings []*internal.MedalRanking) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.listedMonth = s.month()
	boards := make(map[string]*MedalRanking)
	for _, pb := range rankings {
		ranking := NewMedalRankingFromInternalProto(pb)
		if current := boards[ranking.Kind]; current == nil || current.Month < ranking.Month {
			boards[ranking.Kind] = ranking
		}
	}
	s.boards = boards
}

func (s *MedalSystem) Revoke(characterID uint32) {
	s.callHook(characterID, "on_medal_lost")
}

// The character may change maps while internal answers, so the medals are taken back on the actor that holds it by then.
func (s *MedalSystem) Recheck(actorCtx actor.Context, characterID uint32) {
	s.RefreshAsync(actorCtx).Do(func(bool) error {
		s.callHook(characterID, "on_medal_rechecked")
		return nil
	}).OnError(func(err error) {
		log.Printf("refresh medal rankings: %v", err)
	})
}

func (s *MedalSystem) callHook(characterID uint32, hook string) {
	s.gw.Systems().Get[*DispatchSystem]().Visit(characterID, func(actorCtx actor.Context, ch *Character) {
		if ch == nil || ch.LoggedOut() {
			return
		}
		m := ch.GetMap()
		if m == nil {
			return
		}
		root := m.GetLuaRoot()
		if root == nil {
			return
		}
		thread, err := luax.NewThread(root, constant.CharacterHookScriptPath)
		if err != nil {
			log.Printf("%s: %v", hook, err)
			return
		}
		luax.SetConfiguration(thread, luax.Configuration{ActorContext: actorCtx, ActorPID: m.LogicActorPID()})
		luax.CallAsync(actorCtx, root, thread, hook, ch).OnError(func(err error) {
			log.Printf("%s: %v", hook, err)
		})
	})
}
