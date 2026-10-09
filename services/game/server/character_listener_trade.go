package server

import (
	pconst "github.com/boyism80/fm/protocol/constant"
	"github.com/boyism80/fm/protocol/response"
	"github.com/boyism80/fm/services/game/entity"
	"github.com/boyism80/fm/types"
)

func (l *CharacterListenerImpl) OnTradeEntered(ch *entity.Character, t *entity.Trade) {
	slot, _ := t.SlotOf(ch)
	enter := &response.TradeEnter{MySlot: slot}
	for _, member := range t.Members() {
		memberSlot, _ := t.SlotOf(member)
		enter.Members = append(enter.Members, response.MiniRoomVisitor{
			Slot:      memberSlot,
			Character: member.ToDTO(),
		})
	}
	ch.Send(enter, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnTradeInvited(ch *entity.Character, inviter string, sn uint32) {
	ch.Send(&response.TradeInvite{
		Inviter: inviter,
		SN:      sn,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnTradeInviteResult(ch *entity.Character, result pconst.MiniRoomInviteResult, name string) {
	ch.Send(&response.TradeInviteResult{
		Result: result,
		Name:   name,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnTradeItem(ch *entity.Character, who uint8, slot uint8, item entity.Item) {
	ch.Send(&response.TradeItem{
		Who:  who,
		Slot: slot,
		Item: item.ToDTO(),
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnTradeMeso(ch *entity.Character, who uint8, meso int32) {
	ch.Send(&response.TradeMeso{
		Who:  who,
		Meso: meso,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnTradeConfirmed(ch *entity.Character) {
	ch.Send(&response.TradeConfirm{}, types.SEND_POLICY_ENCRYPT)
}
