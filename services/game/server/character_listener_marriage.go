package server

import (
	"github.com/boyism80/fm/maple/constant"
	pconst "github.com/boyism80/fm/protocol/constant"
	"github.com/boyism80/fm/protocol/dto"
	"github.com/boyism80/fm/protocol/response"
	"github.com/boyism80/fm/services/game"
	"github.com/boyism80/fm/types"
)

func (l *CharacterListenerImpl) OnEngageRequest(ch *game.Character, name string, characterID uint32) {
	ch.Send(&response.EngageRequest{
		Mode:        pconst.EngageRequestPropose,
		Name:        name,
		CharacterID: characterID,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnWeddingWishlistInput(ch *game.Character) {
	ch.Send(&response.EngageRequest{
		Mode: pconst.EngageRequestWishlist,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnEngageResult(ch *game.Character, result pconst.EngageResult) {
	packet := &response.EngageResult{Result: result}
	if result == pconst.EngageResultEngaged || result == pconst.EngageResultMarried {
		packet.Marriage = ch.Wedding.Marriage.BuildS2C()
	}
	ch.Send(packet, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnWeddingInvitation(ch *game.Character, groomName string, brideName string, weddingType uint16) {
	ch.Send(&response.EngageResult{
		Result:      pconst.EngageResultInvitation,
		GroomName:   groomName,
		BrideName:   brideName,
		WeddingType: weddingType,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnWeddingGift(ch *game.Character, mode pconst.WeddingGiftMode, wishes []string, tabs map[constant.InventoryType][]game.Item) {
	packet := &response.WeddingGift{Mode: mode, Wishes: wishes}
	if tabs != nil {
		packet.Tabs = response.WeddingGiftTabs{}
		for typ, items := range tabs {
			dtos := make([]dto.Item, 0, len(items))
			for _, item := range items {
				dtos = append(dtos, item.BuildS2C())
			}
			packet.Tabs[typ] = dtos
		}
	}
	ch.Send(packet, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnSpouseMap(ch *game.Character, mapID uint32, spouseID uint32) {
	ch.Send(&response.SpouseMap{
		MapID:       mapID,
		CharacterID: spouseID,
	}, types.SEND_POLICY_ENCRYPT)
}
