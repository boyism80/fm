package server

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/boyism80/fm/core"
	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
)

type catalogSystem struct {
	client      internal.InternalClient
	interval    time.Duration
	maxAttempts int
	worlds      []*internal.WorldCatalog
	channels    map[uint32]map[uint32]*internal.ChannelCatalog
}

func newCatalogSystem(client internal.InternalClient, interval time.Duration, maxAttempts int) *catalogSystem {
	return &catalogSystem{
		client:      client,
		interval:    interval,
		maxAttempts: maxAttempts,
		channels:    make(map[uint32]map[uint32]*internal.ChannelCatalog),
	}
}

func (s *catalogSystem) Start() error { return nil }
func (s *catalogSystem) Stop()        {}

func (s *catalogSystem) load() error {
	var lastErr error
	for attempt := 1; attempt <= s.maxAttempts; attempt++ {
		ctx, cancel := context.WithTimeout(context.Background(), core.InternalRPCPerStepTimeout)
		reply, err := s.client.GetServerCatalog(ctx, &internal.GetServerCatalogRequest{})
		cancel()
		if err == nil {
			worlds := reply.GetWorlds()
			if len(worlds) == 0 {
				lastErr = fmt.Errorf("GetServerCatalog returned no world definitions")
			} else {
				s.worlds = worlds
				s.channels = make(map[uint32]map[uint32]*internal.ChannelCatalog)
				for _, world := range s.worlds {
					wid := world.GetWorldId()
					if s.channels[wid] == nil {
						s.channels[wid] = make(map[uint32]*internal.ChannelCatalog)
					}
					for _, ch := range world.GetChannels() {
						s.channels[wid][ch.GetChannelId()] = ch
					}
				}
				if attempt > 1 {
					log.Printf("GetServerCatalog succeeded after retry (%d/%d)", attempt, s.maxAttempts)
				}
				return nil
			}
		} else {
			lastErr = fmt.Errorf("GetServerCatalog failed: %w", err)
		}

		if attempt < s.maxAttempts {
			log.Printf(
				"GetServerCatalog attempt %d/%d failed: %v (retrying in %s)",
				attempt,
				s.maxAttempts,
				lastErr,
				s.interval,
			)
			time.Sleep(s.interval)
		}
	}

	return fmt.Errorf("GetServerCatalog failed after %d attempts: %w", s.maxAttempts, lastErr)
}
