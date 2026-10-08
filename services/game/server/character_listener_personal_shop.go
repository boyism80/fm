package server

import (
	"github.com/boyism80/fm/protocol/response"
	"github.com/boyism80/fm/services/game/entity"
	"github.com/boyism80/fm/types"
)

func (l *CharacterListenerImpl) OnPersonalShopEntered(ch *entity.Character, ps *entity.PersonalShop) {
	slot, _ := ps.SlotOf(ch)
	enter := &response.PersonalShopEnter{
		MySlot:   slot,
		Title:    ps.Title,
		MaxItems: ps.MaxItems,
		Items:    l.miniRoomItems(ps.Items),
	}
	for _, member := range ps.Members() {
		memberSlot, _ := ps.SlotOf(member)
		enter.Members = append(enter.Members, response.MiniRoomVisitor{
			Slot:      memberSlot,
			Character: member.ToDTO(),
		})
	}
	ch.Send(enter, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnPersonalShopItems(ch *entity.Character, ps *entity.PersonalShop) {
	ch.Send(&response.PersonalShopItems{Items: l.miniRoomItems(ps.Items)}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnMiniRoomSold(ch *entity.Character, index uint8, bundles uint16, buyer string) {
	ch.Send(&response.MiniRoomSold{
		Index:   index,
		Bundles: bundles,
		Buyer:   buyer,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnMiniRoomItemRemoved(ch *entity.Character, count uint8, index uint16) {
	ch.Send(&response.MiniRoomItemRemoved{
		Count: count,
		Index: index,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnMiniRoomBalloon(ch *entity.Character, balloon *response.MiniRoomBalloon) {
	ch.Broadcast(&response.UserMiniRoomBalloon{
		CharacterID: ch.GetID(),
		Balloon:     balloon,
	}, &entity.ObjectBroadcastOption{WithMe: true})
}
