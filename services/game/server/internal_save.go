package server

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/core/async"
	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
	g_actor "github.com/boyism80/fm/services/game/actor"
	"github.com/boyism80/fm/services/game/entity"
)

const saveCharactersPromiseTimeout = 30 * time.Second

const saveCharactersChunkSize = 100

// saveCharactersAsync builds a Promise that saves the given characters in one RPC. Caller must Run().
func (gs *GameServer) saveCharactersAsync(ctx actor.Context, chars []*entity.Character) *async.Promise {
	p := async.NewPromise(ctx, saveCharactersPromiseTimeout)
	if gs == nil || len(chars) == 0 {
		return p
	}
	charsCopy := append([]*entity.Character(nil), chars...)
	p.OnError(func(err error) {
		log.Printf("saveCharacters failed: %v", err)
	})
	async.ThenRPC(p, func(c context.Context) (*internal.SaveCharactersReply, error) {
		return gs.grpcSaveCharacters(c, charsCopy)
	}, func(*internal.SaveCharactersReply) error {
		return nil
	})
	return p
}

// saveCharacterAsync builds a Promise that saves one character. Caller must Run().
func (gs *GameServer) saveCharacterAsync(ctx actor.Context, ch *entity.Character) *async.Promise {
	if ch == nil {
		return async.NewPromise(ctx, saveCharactersPromiseTimeout)
	}
	return gs.saveCharactersAsync(ctx, []*entity.Character{ch})
}

// SaveCharactersAsync builds a Promise that saves chars in parallel chunks of saveCharactersChunkSize. Caller must Run().
func (gs *GameServer) SaveCharactersAsync(ctx actor.Context, chars []*entity.Character) *async.Promise {
	p := async.NewPromise(ctx, saveCharactersPromiseTimeout)
	if gs == nil || gs.internalClient == nil || len(chars) == 0 {
		return p
	}
	p.OnError(func(err error) {
		log.Printf("saveCharactersChunked: %v", err)
	})
	p.Then(func() (interface{}, error) {
		done := make(chan struct{})
		var wg sync.WaitGroup
		var mu sync.Mutex
		var firstErr error

		for i := 0; i < len(chars); i += saveCharactersChunkSize {
			end := i + saveCharactersChunkSize
			if end > len(chars) {
				end = len(chars)
			}
			chunk := append([]*entity.Character(nil), chars[i:end]...)
			wg.Add(1)
			cp := async.NewPromise(nil, saveCharactersPromiseTimeout)
			async.ThenRPC(cp, func(c context.Context) (*internal.SaveCharactersReply, error) {
				return gs.grpcSaveCharacters(c, chunk)
			}, func(*internal.SaveCharactersReply) error {
				return nil
			})
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
			go cp.Run()
		}

		go func() {
			wg.Wait()
			mu.Lock()
			err := firstErr
			mu.Unlock()
			if err != nil {
				p.SetError(err)
			} else {
				p.SetResult()
			}
			close(done)
		}()

		<-done
		return nil, nil
	}, func(interface{}) error {
		return nil
	})
	return p
}

type saveAllSummary struct {
	Maps  int
	Saved int
}

// SaveAllCharactersAsync builds a Promise that asks all map actors to persist online characters in parallel.
// Caller must Run().
func (gs *GameServer) SaveAllCharactersAsync(ctx actor.Context) *async.Promise {
	p := async.NewPromise(ctx, saveCharactersPromiseTimeout)
	if gs == nil {
		return p
	}
	p.Then(func() (interface{}, error) {
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
			pid := m.GetActorPID()
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
	}, func(v interface{}) error {
		s, _ := v.(*saveAllSummary)
		if s != nil {
			log.Printf("save all complete: maps=%d saved=%d", s.Maps, s.Saved)
		}
		return nil
	})
	return p
}
