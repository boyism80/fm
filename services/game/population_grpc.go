package game

import (
	"time"

	"github.com/boyism80/fm/core/clock"
	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
)

func NewPopulationHistoryFromInternalProto(pb *internal.PopulationHistory) PopulationHistory {
	history := PopulationHistory{given: map[uint32]time.Time{}}
	if pb.GetLastGivenUnixMs() != 0 {
		history.lastGivenAt = time.UnixMilli(pb.GetLastGivenUnixMs())
	}

	now := clock.Now()
	for id, unixMs := range pb.GetGivenUnixMs() {
		givenAt := time.UnixMilli(unixMs)
		if now.Sub(givenAt) < populationTargetCooldown {
			history.given[id] = givenAt
		}
	}
	return history
}

func (h *PopulationHistory) BuildRPC() *internal.PopulationHistory {
	pb := &internal.PopulationHistory{GivenUnixMs: make(map[uint32]int64, len(h.given))}
	if !h.lastGivenAt.IsZero() {
		pb.LastGivenUnixMs = h.lastGivenAt.UnixMilli()
	}
	for id, givenAt := range h.given {
		pb.GivenUnixMs[id] = givenAt.UnixMilli()
	}
	return pb
}
