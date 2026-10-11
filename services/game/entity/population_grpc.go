package entity

import (
	"time"

	"github.com/boyism80/fm/core/clock"
	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
)

func NewPopulationHistoryFromInternalProto(pb *internal.PopulationHistory) PopulationHistory {
	history := PopulationHistory{Given: map[uint32]time.Time{}}
	if pb.GetLastGivenUnixMs() != 0 {
		history.LastGivenAt = time.UnixMilli(pb.GetLastGivenUnixMs())
	}

	now := clock.Now()
	for id, unixMs := range pb.GetGivenUnixMs() {
		givenAt := time.UnixMilli(unixMs)
		if now.Sub(givenAt) < populationTargetCooldown {
			history.Given[id] = givenAt
		}
	}
	return history
}

func (h *PopulationHistory) ToProto() *internal.PopulationHistory {
	pb := &internal.PopulationHistory{GivenUnixMs: make(map[uint32]int64, len(h.Given))}
	if !h.LastGivenAt.IsZero() {
		pb.LastGivenUnixMs = h.LastGivenAt.UnixMilli()
	}
	for id, givenAt := range h.Given {
		pb.GivenUnixMs[id] = givenAt.UnixMilli()
	}
	return pb
}
