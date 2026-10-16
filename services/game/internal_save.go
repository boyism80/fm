package game

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/core"
	c_actor "github.com/boyism80/fm/core/actor"
	"github.com/boyism80/fm/core/async"
	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
	"google.golang.org/grpc"
)

const (
	SaveCharactersPromiseTimeout = 30 * time.Second
	MapActorCallTimeout          = 30 * time.Second
)

const saveCharactersChunkSize = 100

func (s *InternalSystem) SaveAsync(ctx actor.Context, entries []*internal.CharacterSaveEntry) *async.Task {
	p := async.NewTask(ctx, SaveCharactersPromiseTimeout)
	if len(entries) == 0 {
		return p
	}
	p.OnError(func(err error) {
		log.Printf("saveCharactersChunked: %v", err)
	})
	p.DoAsync(func() error {
		var wg sync.WaitGroup
		var mu sync.Mutex
		var firstErr error

		for i := 0; i < len(entries); i += saveCharactersChunkSize {
			end := i + saveCharactersChunkSize
			if end > len(entries) {
				end = len(entries)
			}
			chunk := append([]*internal.CharacterSaveEntry(nil), entries[i:end]...)
			wg.Add(1)
			cp := async.NewTask(nil, SaveCharactersPromiseTimeout)
			cp.ThenRPC(func(c context.Context) (*internal.SaveCharactersReply, error) {
				return s.client.SaveCharacters(c, &internal.SaveCharactersRequest{Entries: chunk}, grpc.WaitForReady(true))
			}, nil)
			cp.OnError(func(err error) {
				mu.Lock()
				if firstErr == nil {
					firstErr = err
				}
				mu.Unlock()
			})
			cp.Finally(func() {
				wg.Done()
			})
		}

		wg.Wait()
		mu.Lock()
		err := firstErr
		mu.Unlock()
		return err
	})
	return p
}

type SaveAllSummary struct {
	Maps  int
	Saved int
}

// SaveAllCharactersAsync asks every map actor, fixed and instance, to persist its online characters in parallel.
func (s *InternalSystem) SaveAllCharactersAsync(ctx actor.Context) *async.Promise[*SaveAllSummary] {
	return async.NewTask(ctx, SaveCharactersPromiseTimeout).ThenAsync(func() (*SaveAllSummary, error) {
		root := s.gw.Systems().Get[*c_actor.ActorRegistry]().Root()
		targets := make([]*actor.PID, 0)
		seen := map[string]bool{}
		s.gw.Systems().Get[*MapSystem]().Each(func(m *Map) {
			pid := m.LogicActorPID()
			if pid == nil || seen[pid.Id] {
				return
			}
			seen[pid.Id] = true
			targets = append(targets, pid)
		})
		if len(targets) == 0 {
			return &SaveAllSummary{}, nil
		}

		var wg sync.WaitGroup
		var mu sync.Mutex
		var errs []error
		summary := &SaveAllSummary{Maps: len(targets)}
		for _, t := range targets {
			wg.Add(1)
			go func(pid *actor.PID) {
				defer wg.Done()
				res, err := root.RequestFuture(pid, &SaveMapCharacters{}, MapActorCallTimeout).Result()
				if err != nil {
					mu.Lock()
					errs = append(errs, fmt.Errorf("map actor %s save request: %w", pid.Id, err))
					mu.Unlock()
					return
				}
				ack, ok := res.(*SaveMapCharactersAck)
				if !ok {
					mu.Lock()
					errs = append(errs, fmt.Errorf("map actor %s unexpected ack type %T", pid.Id, res))
					mu.Unlock()
					return
				}
				mu.Lock()
				summary.Saved += ack.Saved
				if ack.Err != "" {
					errs = append(errs, fmt.Errorf("map actor %s save failed: %s", pid.Id, ack.Err))
				}
				mu.Unlock()
			}(t)
		}
		wg.Wait()
		if len(errs) > 0 {
			return summary, errors.Join(errs...)
		}
		return summary, nil
	})
}

func (s *InternalSystem) SaveShopAsync(ctx actor.Context, shop *internal.Shop, entries []*internal.CharacterSaveEntry, storeBank []*internal.Shop, close bool) *async.Promise[*internal.SaveShopReply] {
	shop.WorldId = s.gw.GetWorldID()
	shop.ChannelId = s.gw.GetChannelID()
	req := &internal.SaveShopRequest{
		Shop:       shop,
		Characters: entries,
		Close:      close,
		StoreBank:  storeBank,
	}
	return async.NewTask(ctx, core.InternalRPCPerStepTimeout).ThenRPC(
		func(c context.Context) (*internal.SaveShopReply, error) {
			return s.client.SaveShop(c, req)
		},
		nil,
	)
}
