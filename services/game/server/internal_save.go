package server

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/core/async"
	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
	g_actor "github.com/boyism80/fm/services/game/actor"
)

const saveCharactersPromiseTimeout = 30 * time.Second

const saveCharactersChunkSize = 100

// SaveAsync builds a Promise that saves entries in parallel chunks.
func (gs *GameServer) SaveAsync(ctx actor.Context, entries []*internal.CharacterSaveEntry) *async.Task {
	p := async.NewTask(ctx, saveCharactersPromiseTimeout)
	if gs == nil || gs.internalClient == nil || len(entries) == 0 {
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
			cp := async.NewTask(nil, saveCharactersPromiseTimeout)
			cp.ThenRPC(func(c context.Context) (*internal.SaveCharactersReply, error) {
				return gs.internalClient.SaveCharacters(c, &internal.SaveCharactersRequest{Entries: chunk})
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

type saveAllSummary struct {
	Maps  int
	Saved int
}

// SaveAllCharactersAsync builds a Promise that asks all map actors to persist online characters in parallel.
func (gs *GameServer) SaveAllCharactersAsync(ctx actor.Context) *async.Promise[*saveAllSummary] {
	return async.NewTask(ctx, saveCharactersPromiseTimeout).ThenAsync(func() (*saveAllSummary, error) {
		root := gs.GetRootContext()
		if root == nil {
			return nil, fmt.Errorf("save all: nil root context")
		}
		type mapTarget struct {
			mapID uint32
			pid   *actor.PID
		}
		targets := make([]mapTarget, 0)
		gs.mapsMutex.RLock()
		for mapID, m := range gs.maps {
			if m == nil {
				continue
			}
			pid := m.LogicActorPID()
			if pid == nil {
				continue
			}
			targets = append(targets, mapTarget{mapID: mapID, pid: pid})
		}
		gs.mapsMutex.RUnlock()
		if len(targets) == 0 {
			return &saveAllSummary{}, nil
		}

		var wg sync.WaitGroup
		var mu sync.Mutex
		var errs []error
		summary := &saveAllSummary{Maps: len(targets)}
		for _, t := range targets {
			wg.Add(1)
			go func(t mapTarget) {
				defer wg.Done()
				res, err := root.RequestFuture(t.pid, &g_actor.SaveMapCharacters{}, mapActorCallTimeout).Result()
				if err != nil {
					mu.Lock()
					errs = append(errs, fmt.Errorf("map %d save request: %w", t.mapID, err))
					mu.Unlock()
					return
				}
				ack, ok := res.(*g_actor.SaveMapCharactersAck)
				if !ok {
					mu.Lock()
					errs = append(errs, fmt.Errorf("map %d unexpected ack type %T", t.mapID, res))
					mu.Unlock()
					return
				}
				mu.Lock()
				summary.Saved += ack.Saved
				if ack.Err != "" {
					errs = append(errs, fmt.Errorf("map %d save failed: %s", t.mapID, ack.Err))
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

func (gs *GameServer) SaveShopAsync(ctx actor.Context, shop *internal.Shop, entries []*internal.CharacterSaveEntry, storeBank []*internal.Shop, close bool) *async.Promise[*internal.SaveShopReply] {
	shop.WorldId = gs.config.WorldId
	shop.ChannelId = int32(gs.config.ChannelId)
	req := &internal.SaveShopRequest{
		Shop:       shop,
		Characters: entries,
		Close:      close,
		StoreBank:  storeBank,
	}
	return async.NewTask(ctx, core.InternalRPCPerStepTimeout).ThenRPC(
		func(c context.Context) (*internal.SaveShopReply, error) {
			return gs.internalClient.SaveShop(c, req)
		},
		nil,
	)
}
