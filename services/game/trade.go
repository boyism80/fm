package game

import (
	"fmt"
	"sync/atomic"

	"github.com/boyism80/fm/maple/constant"
	pconst "github.com/boyism80/fm/protocol/constant"
)

const TradeSlots = 9

var nextTradeInviteSN atomic.Uint32

type tradeSide struct {
	ch        *Character
	items     [TradeSlots]Item
	confirmed bool
}

type Trade struct {
	sides   [2]*tradeSide
	invited *Character
	sn      uint32
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
	_, ok := ch.miniRoom.(*Trade)
	return ok
}

func NewTrade(owner *Character) (*Trade, error) {
	if owner.GetMap() == nil || owner.miniRoom != nil {
		return nil, ErrMiniRoomInvalid
	}

	t := &Trade{}
	t.sides[0] = &tradeSide{ch: owner}
	return t, nil
}

func (t *Trade) Open() {
	t.sides[0].ch.miniRoom = t
	t.sides[0].ch.listener.OnTradeEntered(t.sides[0].ch, t)
}

func (t *Trade) Invite(ch *Character, targetID uint32) error {
	if t.sides[0].ch != ch || ch.miniRoom != t || t.invited != nil {
		return ErrMiniRoomInvalid
	}

	m := ch.GetMap()
	if m == nil {
		return ErrMiniRoomInvalid
	}
	target := m.GetPlayer(targetID)
	switch {
	case target == nil || target == ch:
		ch.listener.OnTradeInviteResult(ch, pconst.MiniRoomInviteNotFound, "")
		t.Leave(ch)
	case target.miniRoom != nil:
		ch.listener.OnTradeInviteResult(ch, pconst.MiniRoomInviteBusy, target.GetName())
		t.Leave(ch)
	default:
		t.invited = target
		t.sn = nextTradeInviteSN.Add(1)
		target.miniRoom = t
		target.listener.OnTradeInvited(target, ch.GetName(), t.sn)
	}
	return nil
}

func (t *Trade) Decline(ch *Character, sn uint32, reason uint8) error {
	if ch.miniRoom != t || t.sn != sn || t.invited != ch || t.sides[1] != nil {
		return ErrMiniRoomInvalid
	}

	result := pconst.MiniRoomInviteDeclined
	if pconst.MiniRoomInviteResult(reason) == pconst.MiniRoomInviteBlocked {
		result = pconst.MiniRoomInviteBlocked
	}
	inviter := t.sides[0].ch
	inviter.listener.OnTradeInviteResult(inviter, result, ch.GetName())
	t.Leave(inviter)
	return nil
}

func (t *Trade) Visit(ch *Character) error {
	if ch.miniRoom != t {
		return ErrMiniRoomInvalid
	}
	if t.invited != ch || t.sides[1] != nil || t.sides[0].ch.GetMap() == nil || t.sides[0].ch.GetMap() != ch.GetMap() {
		return &MiniRoomEnterError{Code: pconst.MiniRoomEnterClosed}
	}

	t.invited = nil
	t.sides[1] = &tradeSide{ch: ch}
	inviter := t.sides[0].ch
	inviter.listener.OnMiniRoomVisited(inviter, 1, ch)
	ch.listener.OnTradeEntered(ch, t)
	return nil
}

func (t *Trade) Chat(ch *Character, message string) error {
	slot, ok := t.SlotOf(ch)
	if !ok {
		return ErrMiniRoomInvalid
	}

	text := fmt.Sprintf("%s : %s", ch.GetName(), message)
	for _, member := range t.Members() {
		member.listener.OnMiniRoomChat(member, slot, text)
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
	item := ch.inventory.GetItem(invType, slot)
	if item == nil {
		return ErrMiniRoomItemNotFound
	}
	if !item.Tradable() {
		return ErrMiniRoomInvalid
	}
	if invType == constant.InventoryTypeCash || item.GetModel().GetTraits()&constant.ItemTraitCash != 0 {
		return ErrMiniRoomInvalid
	}
	if invType == constant.InventoryTypeEquipment || item.GetModel().GetTraits()&constant.ItemTraitRechargeable != 0 {
		count = item.GetCount()
	}

	held, err := ch.escrow.hold(invType, slot, count)
	if err != nil {
		return err
	}
	side.items[tradeSlot-1] = held
	for _, member := range t.sides {
		who := uint8(1)
		if member == side {
			who = 0
		}
		member.ch.listener.OnTradeItem(member.ch, who, tradeSlot, held)
	}
	return nil
}

func (t *Trade) PutMeso(ch *Character, meso int32) error {
	side, err := t.unlockedSide(ch)
	if err != nil {
		return err
	}
	if err := ch.escrow.holdMeso(meso); err != nil {
		return err
	}

	for _, member := range t.sides {
		who := uint8(1)
		if member == side {
			who = 0
		}
		member.ch.listener.OnTradeMeso(member.ch, who, ch.escrow.meso)
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
		other.ch.listener.OnTradeConfirmed(other.ch)
		return nil
	}
	t.complete()
	return nil
}

func (t *Trade) Leave(ch *Character) {
	if t.invited == ch && ch.miniRoom == t {
		inviter := t.sides[0].ch
		t.close()
		inviter.listener.OnMiniRoomLeft(inviter, 0, pconst.MiniRoomLeaveTradeCancel)
		return
	}
	slot, ok := t.SlotOf(ch)
	if !ok || ch.miniRoom != t {
		return
	}

	t.close()
	for i, side := range t.sides {
		if side == nil {
			continue
		}
		side.ch.escrow.restore()
		reason := pconst.MiniRoomLeaveTradeCancel
		if uint8(i) == slot {
			reason = pconst.MiniRoomLeaveExit
		}
		side.ch.listener.OnMiniRoomLeft(side.ch, uint8(i), reason)
	}
}

func (t *Trade) close() {
	if t.invited != nil && t.invited.miniRoom == t {
		t.invited.miniRoom = nil
	}
	t.invited = nil
	for _, side := range t.sides {
		if side != nil {
			if side.ch.miniRoom == t {
				side.ch.miniRoom = nil
			}
		}
	}
}

func (t *Trade) check() pconst.MiniRoomLeaveReason {
	for i, side := range t.sides {
		if side.ch.LoggedOut() {
			return pconst.MiniRoomLeaveTradeFail
		}
		giver := t.sides[1-i]
		var reward ExchangeSide
		uniques := make(map[uint32]struct{})
		for _, item := range giver.items {
			if item == nil {
				continue
			}
			model := item.GetModel()
			// The receiver's own offer leaves with the trade, so only what the receiver carries counts.
			if side.ch.inventory.carriesUniqueItem(model) {
				return pconst.MiniRoomLeaveTradeOnly
			}
			if model.GetTraits()&constant.ItemTraitUnique != 0 {
				if _, twice := uniques[model.GetID()]; twice {
					return pconst.MiniRoomLeaveTradeOnly
				}
				uniques[model.GetID()] = struct{}{}
			}
			reward.include(item, item.GetCount())
		}
		meso := giver.ch.escrow.meso
		reward.Meso = meso - t.tax(meso)
		if (ExchangeSpec{Reward: reward}).Valid(side.ch, nil) != ExchangeOK {
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
			side.ch.escrow.restore()
			side.ch.listener.OnMiniRoomLeft(side.ch, uint8(i), reason)
		}
		return
	}

	var given [2][]Item
	var meso [2]int32
	for i, side := range t.sides {
		given[i], meso[i] = side.ch.escrow.release()
	}
	for i, side := range t.sides {
		for _, item := range given[1-i] {
			side.ch.inventory.receive(item)
		}
		side.ch.inventory.addMesoUnchecked(meso[1-i] - t.tax(meso[1-i]))
	}
	for i, side := range t.sides {
		side.ch.listener.OnMiniRoomLeft(side.ch, uint8(i), pconst.MiniRoomLeaveTradeDone)
	}
}
