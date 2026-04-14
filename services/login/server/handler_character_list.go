package server

import (
	"context"
	"log"

	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/protocol/dto"
	fminternalpb "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
	"github.com/boyism80/fm/protocol/request"
	"github.com/boyism80/fm/protocol/response"
	"github.com/boyism80/fm/services/login/client"
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

	loginClient, ok := ctx.Client.(*client.LoginClient)
	if !ok {
		return nil
	}

	worldId := uint32(req.Server)
	loginClient.SetWorldId(worldId)
	accountId := loginClient.GetAccountId()

	ic := h.ls.context.InternalClient
	if ic == nil {
		charListResp := &response.CharacterList{Characters: nil, SlotCount: 6}
		return ctx.Client.Send(charListResp, types.SEND_POLICY_ENCRYPT)
	}

	reply, err := ic.GetCharacterList(context.Background(), &fminternalpb.GetCharacterListRequest{
		AccountId: accountId,
		WorldId:   worldId,
	})
	if err != nil {
		log.Printf("GetCharacterList RPC error: %v", err)
		charListResp := &response.CharacterList{Characters: nil, SlotCount: 6}
		return ctx.Client.Send(charListResp, types.SEND_POLICY_ENCRYPT)
	}

	characters := make([]dto.Character, 0, len(reply.Characters))
	for _, ov := range reply.Characters {
		characters = append(characters, overviewToDto(ov))
	}

	charListResp := &response.CharacterList{
		Characters: characters,
		SlotCount:  reply.SlotCount,
	}
	return ctx.Client.Send(charListResp, types.SEND_POLICY_ENCRYPT)
}

func overviewToDto(ov *fminternalpb.CharacterOverview) dto.Character {
	baseLooks := make(map[int8]uint32)
	for k, v := range ov.BaseLooks {
		baseLooks[int8(k)] = v
	}
	overlays := make(map[int8]uint32)
	for k, v := range ov.Overlays {
		overlays[int8(k)] = v
	}
	return dto.Character{
		ID:            ov.CharacterId,
		Name:          ov.Name,
		Gender:        uint8(ov.Gender),
		SkinColor:     uint8(ov.SkinColor),
		Face:          ov.Face,
		Hair:          ov.Hair,
		Level:         uint8(ov.Level),
		Class:         uint16(ov.ClassId),
		Map:           ov.MapId,
		SpawnPoint:    uint8(ov.SpawnPoint),
		Rank:          uint32(ov.Rank),
		RankDiff:      ov.RankDiff,
		ClassRank:     uint32(ov.ClassRank),
		ClassRankDiff: ov.ClassRankDiff,
		BaseLooks:     baseLooks,
		Overlays:      overlays,
	}
}
