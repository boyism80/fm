package server

import (
	"log"

	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/protocol/dto"
	"github.com/boyism80/fm/protocol/request"
	"github.com/boyism80/fm/protocol/response"
	"github.com/boyism80/fm/types"
)

// CreateCharacter handles create character packet requests
type CreateCharacter struct {
	ls     *LoginServer
	opcode byte
}

func (CreateCharacter) New(ls *LoginServer) *CreateCharacter {
	return &CreateCharacter{
		ls:     ls,
		opcode: 0x08,
	}
}

func (h *CreateCharacter) GetOpcode() byte {
	return h.opcode
}

func (h *CreateCharacter) Handle(ctx *core.ClientContext, req *request.CreateCharacter) error {
	log.Printf("Create character packet received from %s - Name: %s, Face: %d, Hair: %d, Top: %d, Bottom: %d, Shoes: %d, Weapon: %d",
		ctx.Client.GetConnection().RemoteAddr(), req.Name, req.Face, req.Hair, req.Top, req.Bottom, req.Shoes, req.Weapon)

	// Check if creation is successful (hardcoded for demo)
	success := req.Name != "채진영"

	createResp := &response.CreateCharacter{
		Success: success,
		Character: &dto.Character{
			ID:         1,
			Name:       req.Name,
			Gender:     0,
			SkinColor:  0,
			Face:       req.Face,
			Hair:       req.Hair,
			Level:      1,
			Class:      0,
			Str:        12,
			Dex:        5,
			Int:        4,
			Luk:        4,
			Hp:         50,
			MaxHp:      50,
			Mp:         5,
			MaxMp:      5,
			SpawnPoint: 3,
			BaseLooks:  make(map[int8]uint32),
			Overlays:   make(map[int8]uint32),
		},
	}
	if err := ctx.Client.Send(createResp, types.SEND_POLICY_ENCRYPT); err != nil {
		log.Printf("Failed to send create character response: %v", err)
		return err
	}

	return nil
}
