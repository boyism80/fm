package server

import (
	"context"
	"fmt"
	"log"
	"time"

	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
	"github.com/boyism80/fm/services/game/entity"
)

func saveCharactersRPC(ctx context.Context, client internal.InternalClient, worldId uint32, chars []*entity.Character) error {
	if len(chars) == 0 {
		return nil
	}

	entries := make([]*internal.CharacterSaveEntry, 0, len(chars))
	for _, ch := range chars {
		if ch == nil {
			continue
		}
		if entry := ch.ToPersisted(worldId); entry != nil {
			entries = append(entries, entry)
		}
	}

	if len(entries) == 0 {
		return nil
	}

	req := &internal.SaveCharactersRequest{
		Entries: entries,
	}
	_, err := client.SaveCharacters(ctx, req)
	if err != nil {
		return fmt.Errorf("SaveCharacters rpc: %w", err)
	}
	return nil
}

func (gs *GameServer) saveCharactersAsync(chars []*entity.Character) error {
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		if err := saveCharactersRPC(ctx, gs.internalClient, uint32(gs.config.WorldId), chars); err != nil {
			log.Printf("saveCharacters failed: %v", err)
		}
	}()
	return nil
}

func (gs *GameServer) saveCharacterAsync(ch *entity.Character) {
	if ch == nil || gs.internalClient == nil {
		return
	}
	gs.saveCharactersAsync([]*entity.Character{ch})
}

func (gs *GameServer) saveCharactersChunked(chars []*entity.Character) error {
	if gs.internalClient == nil {
		return nil
	}
	const chunkSize = 100
	for i := 0; i < len(chars); i += chunkSize {
		end := i + chunkSize
		if end > len(chars) {
			end = len(chars)
		}
		chunk := chars[i:end]
		if err := gs.saveCharactersAsync(chunk); err != nil {
			log.Printf("saveCharactersChunked: chunk %d-%d: %v", i, end, err)
		}
	}
	return nil
}
