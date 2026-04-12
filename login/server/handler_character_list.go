package server

import (
	"log"

	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/protocol/dto"
	"github.com/boyism80/fm/protocol/request"
	"github.com/boyism80/fm/protocol/response"
	"github.com/boyism80/fm/types"
)

type CharacterList struct {
	ls     *LoginServer
	opcode byte
}

func (CharacterList) New(ls *LoginServer) *CharacterList {
	return &CharacterList{
		ls:     ls,
		opcode: 0x04,
	}
}

func (h *CharacterList) GetOpcode() byte {
	return h.opcode
}

func (h *CharacterList) Handle(ctx *core.ClientContext, req *request.CharacterList) error {
	log.Printf("Character list packet received from %s - Server: %d, Channel: %d",
		ctx.Client.GetConnection().RemoteAddr(), req.Server, req.Channel)

	characters := []dto.Character{
		{
			ID:         1,
			Name:       "채승현",
			Gender:     0,
			SkinColor:  0,
			Face:       20100,
			Hair:       30000,
			Level:      255,
			Class:      0,
			Str:        12,
			Dex:        5,
			Int:        4,
			Luk:        4,
			Hp:         50,
			MaxHp:      50,
			Mp:         5,
			MaxMp:      5,
			SpawnPoint: 1,
			BaseLooks:  make(map[int8]uint32),
			Overlays:   make(map[int8]uint32),
		},
		{
			ID:         2,
			Name:       "채진영",
			Gender:     0,
			SkinColor:  0,
			Face:       20100,
			Hair:       30000,
			Level:      255,
			Class:      0,
			Str:        12,
			Dex:        5,
			Int:        4,
			Luk:        4,
			Hp:         50,
			MaxHp:      50,
			Mp:         5,
			MaxMp:      5,
			SpawnPoint: 1,
			BaseLooks:  make(map[int8]uint32),
			Overlays:   make(map[int8]uint32),
		},
	}

	charListResp := &response.CharacterList{
		Characters: characters,
		SlotCount:  6,
	}
	if err := ctx.Client.Send(charListResp, types.SEND_POLICY_ENCRYPT); err != nil {
		log.Printf("Failed to send character list response: %v", err)
		return err
	}

	return nil
}
