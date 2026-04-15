package server

import (
	"context"
	"fmt"
	"log"
	"time"

	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
	"github.com/boyism80/fm/services/game/entity"
	"github.com/boyism80/fm/services/game/wz"
)

func equipmentLooks(ch *entity.Character) (baseLooks, overlays map[int32]uint32) {
	baseLooks = make(map[int32]uint32)
	overlays = make(map[int32]uint32)

	for parts, equipment := range ch.Equipments {
		if equipment == nil || parts < -127 {
			continue
		}
		model, ok := equipment.GetModel().(wz.Equipment)
		if !ok {
			continue
		}
		absoluteParts := int32(parts * -1)
		if absoluteParts < 100 {
			if _, exists := baseLooks[absoluteParts]; !exists {
				baseLooks[absoluteParts] = model.GetID()
			}
		} else if absoluteParts > 100 && absoluteParts != 111 {
			adjustedParts := absoluteParts - 100
			if existing, exists := baseLooks[adjustedParts]; exists {
				overlays[adjustedParts] = existing
			}
			baseLooks[adjustedParts] = model.GetID()
		} else if _, exists := baseLooks[absoluteParts]; exists {
			overlays[absoluteParts] = model.GetID()
		}
	}
	return
}

func characterToSaveEntry(ch *entity.Character, worldId uint32) *internal.CharacterSaveEntry {
	mapID := uint32(0)
	if m := ch.GetMap(); m != nil {
		mapID = m.GetMapID()
	}

	baseLooks, overlays := equipmentLooks(ch)

	persisted := &internal.CharacterPersisted{
		CharacterId:     ch.GetID(),
		AccountId:       ch.AccountID,
		WorldId:         worldId,
		Name:            ch.GetName(),
		Gender:          uint32(ch.GetGender()),
		SkinColor:       uint32(ch.GetSkinColor()),
		Face:            ch.GetFace(),
		Hair:            ch.GetHair(),
		Level:           uint32(ch.GetLevel()),
		ClassId:         uint32(ch.Class),
		Role:            uint32(ch.Role),
		Str:             uint32(ch.BaseStats.Str),
		Dex:             uint32(ch.BaseStats.Dex),
		IntStat:         uint32(ch.BaseStats.Int),
		Luk:             uint32(ch.BaseStats.Luk),
		Hp:              ch.Hp,
		MaxHp:           ch.BaseHp,
		Mp:              ch.Mp,
		MaxMp:           ch.BaseMp,
		AbilityPoint:    uint32(ch.AbilityPoint),
		Exp:             ch.GetExp(),
		MapId:           mapID,
		SpawnPoint:      uint32(ch.GetSpawnPoint()),
		PositionX:       int32(ch.Position.X),
		PositionY:       int32(ch.Position.Y),
		Stance:          uint32(ch.Stance),
		Meso:            ch.Meso,
		SkillPoint:      uint32(ch.SkillPoint),
		UpdatedAtUnixMs: time.Now().UnixMilli(),
	}

	return &internal.CharacterSaveEntry{
		Character: persisted,
		BaseLooks: baseLooks,
		Overlays:  overlays,
	}
}

func saveCharactersRPC(ctx context.Context, client internal.InternalClient, worldId uint32, chars []*entity.Character) error {
	if len(chars) == 0 {
		return nil
	}

	entries := make([]*internal.CharacterSaveEntry, 0, len(chars))
	for _, ch := range chars {
		if ch == nil {
			continue
		}
		entries = append(entries, characterToSaveEntry(ch, worldId))
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
