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

var kst = time.FixedZone("KST", 9*60*60)

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

// MedalSystem keeps the world's title medal boards; quest scripts on any actor read them and the MQ consumer
// replaces them.
type MedalSystem struct {
	gw     GameWorld
	client internal.InternalClient
	mu     sync.Mutex
	boards map[string]*MedalRanking
}

func NewMedalSystem(gw GameWorld, conn *rpc.InternalConnection) *MedalSystem {
	return &MedalSystem{gw: gw, client: conn.Client, boards: make(map[string]*MedalRanking)}
}

func (s *MedalSystem) Start() error {
	// A channel can boot before internal listens; the boards would stay empty until the next submit.
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	reply, err := s.client.ListMedalRankings(ctx, &internal.ListMedalRankingsRequest{WorldId: s.gw.GetWorldID()}, grpc.WaitForReady(true))
	cancel()
	if err != nil {
		log.Printf("list medal rankings: %v", err)
		return nil
	}
	for _, pb := range reply.GetRankings() {
		s.Set(NewMedalRankingFromInternalProto(pb))
	}
	return nil
}

func (s *MedalSystem) Stop() {}

func (s *MedalSystem) month() uint32 {
	now := clock.Now().In(kst)
	return uint32(now.Year()*100 + int(now.Month()))
}

// The boards reset when the KST month changes, so a board from an earlier month is empty.
func (s *MedalSystem) Entries(kind string) []MedalRankingEntry {
	s.mu.Lock()
	defer s.mu.Unlock()
	board := s.boards[kind]
	if board == nil || board.Month != s.month() {
		return nil
	}
	return board.Entries
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

// Set is idempotent: the channel that submitted the score also receives its MQ notice.
func (s *MedalSystem) Set(board *MedalRanking) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if current := s.boards[board.Kind]; current != nil && current.Month > board.Month {
		return
	}
	s.boards[board.Kind] = board
}

func (s *MedalSystem) Revoke(characterID uint32, kind string) {
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
			log.Printf("on_medal_lost: %v", err)
			return
		}
		luax.CallAsync(actorCtx, root, thread, "on_medal_lost", ch, kind).OnError(func(err error) {
			log.Printf("on_medal_lost: %v", err)
		})
	})
}
