package entity

import (
	"fmt"

	pconst "github.com/boyism80/fm/protocol/constant"
	"github.com/boyism80/fm/services/game/constant"
)

const TradeSlots = 9

type tradeSide struct {
	ch        *Character
	items     [TradeSlots]Item
	confirmed bool
}

type Trade struct {
	ObjectCore
	sides   [2]*tradeSide
	invited uint32
}

func (t *Trade) GetObjectType() constant.ObjectType {
	return constant.ObjectTypeTrade
}

func (t *Trade) Is(typ constant.ObjectType) bool {
	return t.GetObjectType().Has(typ)
}

func (t *Trade) SlotOf(ch *Character) (uint8, bool) {
	for i, side := range t.sides {
		if side != nil && side.ch == ch {
			return uint8(i), true
		}
	}
	return 0, false
}

func (t *Trade) Members() []*Character {
	members := make([]*Character, 0, len(t.sides))
	for _, side := range t.sides {
		if side != nil {
			members = append(members, side.ch)
		}
	}
	return members
}

func (t *Trade) tax(meso int32) int32 {
	var permille int64
	switch {
	case meso >= 100000000:
		permille = 60
	case meso >= 25000000:
		permille = 50
	case meso >= 10000000:
		permille = 40
	case meso >= 5000000:
		permille = 30
	case meso >= 1000000:
		permille = 18
	case meso >= 100000:
		permille = 8
	}
	return int32(int64(meso) * permille / 1000)
}

func (ch *Character) Trading() bool {
	_, ok := ch.MiniRoom.(*Trade)
	return ok
}

func (ch *Character) CreateTrade() error {
	m := ch.GetMap()
	if m == nil || ch.MiniRoom != nil || ch.miniRoomPending {
		return ErrMiniRoomInvalid
	}

	t := &Trade{}
	t.ObjectCore.self = t
	t.Position = ch.Position
	t.sides[0] = &tradeSide{ch: ch}
	m.AddTrade(t)
	ch.MiniRoom = t
	ch.Listener.OnTradeEntered(ch, t)
	return nil
}

func (t *Trade) Invite(ch *Character, targetID uint32) error {
	if t.sides[0].ch != ch || t.invited != 0 {
		return ErrMiniRoomInvalid
	}

	target := t.Map.GetPlayer(targetID)
	switch {
	case target == nil || target == ch:
		ch.Listener.OnTradeInviteResult(ch, pconst.MiniRoomInviteNotFound, "")
		t.Leave(ch)
	case target.MiniRoom != nil || target.miniRoomPending:
		ch.Listener.OnTradeInviteResult(ch, pconst.MiniRoomInviteBusy, target.GetName())
		t.Leave(ch)
	default:
		t.invited = targetID
		target.Listener.OnTradeInvited(target, ch.GetName(), t.OID)
	}
	return nil
}

func (ch *Character) DeclineTrade(sn uint32, reason uint8) error {
	t, ok := ch.GetMap().GetObject(constant.ObjectTypeTrade, sn).(*Trade)
	if !ok || t.invited != ch.GetID() || t.sides[1] != nil {
		return ErrMiniRoomInvalid
	}

	result := pconst.MiniRoomInviteDeclined
	if pconst.MiniRoomInviteResult(reason) == pconst.MiniRoomInviteBlocked {
		result = pconst.MiniRoomInviteBlocked
	}
	inviter := t.sides[0].ch
	inviter.Listener.OnTradeInviteResult(inviter, result, ch.GetName())
	t.Leave(inviter)
	return nil
}

func (t *Trade) Visit(ch *Character) error {
	if ch.MiniRoom != nil || ch.miniRoomPending {
		return ErrMiniRoomInvalid
	}
	if t.invited != ch.GetID() || t.sides[1] != nil {
		return &MiniRoomEnterError{Code: pconst.MiniRoomEnterClosed}
	}

	t.sides[1] = &tradeSide{ch: ch}
	ch.MiniRoom = t
	inviter := t.sides[0].ch
	inviter.Listener.OnMiniRoomVisited(inviter, 1, ch)
	ch.Listener.OnTradeEntered(ch, t)
	return nil
}

func (t *Trade) Chat(ch *Character, message string) error {
	slot, ok := t.SlotOf(ch)
	if !ok {
		return ErrMiniRoomInvalid
	}

	text := fmt.Sprintf("%s : %s", ch.GetName(), message)
	for _, member := range t.Members() {
		member.Listener.OnMiniRoomChat(member, slot, text)
	}
	return nil
}

func (t *Trade) unlockedSide(ch *Character) (*tradeSide, error) {
	slot, ok := t.SlotOf(ch)
	if !ok || t.sides[1] == nil || t.sides[0].confirmed || t.sides[1].confirmed {
		return nil, ErrMiniRoomInvalid
	}
	return t.sides[slot], nil
}

func (t *Trade) PutItem(ch *Character, invType constant.InventoryType, slot int16, count uint16, tradeSlot uint8) error {
	side, err := t.unlockedSide(ch)
	if err != nil {
		return err
	}
	if tradeSlot < 1 || tradeSlot > TradeSlots || side.items[tradeSlot-1] != nil {
		return ErrMiniRoomInvalid
	}
	item, err := ch.Inventory.findTradable(invType, slot)
	if err != nil {
		return err
	}
	if invType == constant.InventoryTypeCash || item.GetModel().IsCash() {
		return ErrMiniRoomInvalid
	}
	if invType == constant.InventoryTypeEquipment || constant.IsRechargeable(item.GetModel().GetID()) {
		count = item.GetCount()
	}

	held, err := ch.Escrow.hold(invType, slot, count)
	if err != nil {
		return err
	}
	side.items[tradeSlot-1] = held
	for _, member := range t.sides {
		who := uint8(1)
		if member == side {
			who = 0
		}
		member.ch.Listener.OnTradeItem(member.ch, who, tradeSlot, held)
	}
	return nil
}

func (t *Trade) PutMeso(ch *Character, meso int32) error {
	side, err := t.unlockedSide(ch)
	if err != nil {
		return err
	}
	if err := ch.Escrow.holdMeso(meso); err != nil {
		return err
	}

	for _, member := range t.sides {
		who := uint8(1)
		if member == side {
			who = 0
		}
		member.ch.Listener.OnTradeMeso(member.ch, who, ch.Escrow.Meso)
	}
	return nil
}

func (t *Trade) Confirm(ch *Character) error {
	slot, ok := t.SlotOf(ch)
	if !ok || t.sides[1] == nil || t.sides[slot].confirmed {
		return ErrMiniRoomInvalid
	}

	t.sides[slot].confirmed = true
	other := t.sides[1-slot]
	if !other.confirmed {
		other.ch.Listener.OnTradeConfirmed(other.ch)
		return nil
	}
	t.complete()
	return nil
}

func (t *Trade) Leave(ch *Character) {
	slot, ok := t.SlotOf(ch)
	if !ok {
		return
	}

	t.close()
	for i, side := range t.sides {
		if side == nil {
			continue
		}
		side.ch.Escrow.restore()
		reason := pconst.MiniRoomLeaveTradeCancel
		if uint8(i) == slot {
			reason = pconst.MiniRoomLeaveExit
		}
		side.ch.Listener.OnMiniRoomLeft(side.ch, uint8(i), reason)
	}
}

func (t *Trade) close() {
	t.Map.RemoveTrade(t)
	for _, side := range t.sides {
		if side != nil {
			side.ch.MiniRoom = nil
		}
	}
}

func (t *Trade) check() pconst.MiniRoomLeaveReason {
	for i, side := range t.sides {
		giver := t.sides[1-i]
		rewards := make(map[uint32]uint16)
		for _, item := range giver.items {
			if item == nil {
				continue
			}
			model := item.GetModel()
			if model.IsOnly() && (side.ch.Inventory.HasItem(model.GetID()) || rewards[model.GetID()] > 0) {
				return pconst.MiniRoomLeaveTradeOnly
			}
			rewards[model.GetID()] += item.GetCount()
		}
		meso := giver.ch.Escrow.Meso
		spec := ExchangeSpec{Reward: ExchangeSide{Items: rewards, Meso: meso - t.tax(meso)}}
		if spec.Valid(side.ch) != ExchangeOK {
			return pconst.MiniRoomLeaveTradeFail
		}
	}
	return pconst.MiniRoomLeaveTradeDone
}

func (t *Trade) complete() {
	t.close()
	reason := t.check()
	if reason != pconst.MiniRoomLeaveTradeDone {
		for i, side := range t.sides {
			side.ch.Escrow.restore()
			side.ch.Listener.OnMiniRoomLeft(side.ch, uint8(i), reason)
		}
		return
	}

	var given [2][]Item
	var meso [2]int32
	for i, side := range t.sides {
		given[i], meso[i] = side.ch.Escrow.release()
	}
	for i, side := range t.sides {
		for _, item := range given[1-i] {
			side.ch.Inventory.receive(item)
		}
		side.ch.Inventory.addMesoUnchecked(meso[1-i] - t.tax(meso[1-i]))
	}
	for i, side := range t.sides {
		side.ch.Listener.OnMiniRoomLeft(side.ch, uint8(i), pconst.MiniRoomLeaveTradeDone)
	}
}
