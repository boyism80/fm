package game

import (
	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
)

func NewMedalRankingFromInternalProto(pb *internal.MedalRanking) *MedalRanking {
	entries := make([]MedalRankingEntry, 0, len(pb.GetEntries()))
	for _, entry := range pb.GetEntries() {
		entries = append(entries, MedalRankingEntry{CharacterID: entry.GetCharacterId(), Name: entry.GetName(), Score: entry.GetScore()})
	}
	return &MedalRanking{Kind: pb.GetKind(), Month: pb.GetMonth(), Entries: entries}
}
