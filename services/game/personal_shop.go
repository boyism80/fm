package game

import (
	"log"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/core/clock"
	"github.com/boyism80/fm/maple/constant"
	pconst "github.com/boyism80/fm/protocol/constant"
	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
	"github.com/boyism80/fm/protocol/response"
)

const (
	PersonalShopMaxItems      = 16
	PersonalShopLargeMaxItems = 24
	personalShopBasicPermit   = 5140000
)

type PersonalShop struct {
	shopRoom
	maxItems uint8
	banned   map[string]struct{}
}

func (ps *PersonalShop) GetObjectType() constant.ObjectType {
	return constant.ObjectTypePersonalShop
}

func (ps *PersonalShop) Is(typ constant.ObjectType) bool {
	return ps.GetObjectType().Has(typ)
}

func (ps *PersonalShop) balloon() response.MiniRoomBalloon {
	return response.MiniRoomBalloon{
		Type:     pconst.MiniRoomTypePersonalShop,
		SN:       ps.oid,
		Title:    ps.title,
		Spec:     uint8(ps.itemID % 10),
		Users:    ps.users(),
		MaxUsers: pconst.MiniRoomShopUsers,
	}
}

func (ps *PersonalShop) updateBalloon() {
	if !ps.published {
		return
	}
	balloon := ps.balloon()
	ps.owner.listener.OnMiniRoomBalloon(ps.owner, &balloon)
}

func (ps *PersonalShop) Visit(ch *Character) error {
	if ch.miniRoom != nil {
		return ErrMiniRoomInvalid
	}
	if !ps.published {
		return &MiniRoomEnterError{Code: pconst.MiniRoomEnterClosed}
	}
	if _, ok := ps.banned[ch.GetName()]; ok {
		return &MiniRoomEnterError{Code: pconst.MiniRoomEnterBlocked}
	}
	index, ok := ps.freeSlot()
	if !ok {
		return &MiniRoomEnterError{Code: pconst.MiniRoomEnterFull}
	}

	slot := uint8(index + 1)
	for _, member := range ps.Members() {
		member.listener.OnMiniRoomVisited(member, slot, ch)
	}
	ps.visitors[index] = ch
	ch.miniRoom = ps
	ps.updateBalloon()
	ch.listener.OnPersonalShopEntered(ch, ps)
	return nil
}

func (ps *PersonalShop) Leave(ch *Character) {
	slot, ok := ps.SlotOf(ch)
	if !ok {
		return
	}
	if slot == 0 {
		ps.close(pconst.MiniRoomLeaveShopClosed)
		return
	}

	ch.miniRoom = nil
	ps.visitors[slot-1] = nil
	for _, member := range ps.Members() {
		member.listener.OnMiniRoomLeft(member, slot, pconst.MiniRoomLeaveExit)
	}
	ps.updateBalloon()
}

func (ps *PersonalShop) AddItem(actx actor.Context, ch *Character, invType constant.InventoryType, slot int16, bundles uint16, perBundle uint16, price int32) error {
	if ps.owner != ch {
		return ErrMiniRoomNotOwner
	}
	if ps.published {
		return ErrMiniRoomInvalid
	}
	if len(ps.items) >= int(ps.maxItems) {
		return ErrMiniRoomFull
	}

	if err := ps.list(ch, invType, slot, bundles, perBundle, price); err != nil {
		return err
	}
	ch.listener.OnPersonalShopItems(ch, ps)
	ps.save(actx, false, ch)
	return nil
}

func (ps *PersonalShop) RemoveItem(actx actor.Context, ch *Character, index uint16) error {
	if ps.owner != ch {
		return ErrMiniRoomNotOwner
	}

	if err := ps.unlist(ch, index); err != nil {
		return err
	}
	ch.listener.OnMiniRoomItemRemoved(ch, uint8(len(ps.items)), index)
	for _, visitor := range ps.visitors {
		if visitor != nil {
			visitor.listener.OnPersonalShopItems(visitor, ps)
		}
	}
	ps.save(actx, false, ch)
	return nil
}

func (ps *PersonalShop) Open(actx actor.Context, ch *Character) error {
	if ps.owner != ch || ps.published {
		return ErrMiniRoomNotOwner
	}
	if len(ps.items) == 0 {
		return ErrMiniRoomInvalid
	}

	ps.published = true
	ps.updateBalloon()
	ps.save(actx, false)
	return nil
}

func (ps *PersonalShop) Buy(actx actor.Context, ch *Character, index uint16, bundles uint16) error {
	slot, ok := ps.SlotOf(ch)
	if !ok || slot == 0 || !ps.published {
		return ErrMiniRoomInvalid
	}

	owner := ps.owner
	_, _, income, err := ps.sell(ch, index, bundles, owner.inventory.meso)
	if err != nil {
		return err
	}
	owner.inventory.addMesoUnchecked(income)
	for _, member := range ps.Members() {
		member.listener.OnPersonalShopItems(member, ps)
	}
	owner.listener.OnMiniRoomSold(owner, uint8(index), bundles, ch.GetName())
	ps.save(actx, false, ch, owner)

	for _, listed := range ps.items {
		if listed.bundles > 0 {
			return nil
		}
	}
	owner.listener.OnMiniRoomLeft(owner, 0, pconst.MiniRoomLeaveSoldOut)
	ps.close(pconst.MiniRoomLeaveSoldOut)
	return nil
}

func (ps *PersonalShop) Kick(ch *Character, slot uint8, name string, reason pconst.MiniRoomLeaveReason) error {
	if ps.owner != ch {
		return ErrMiniRoomNotOwner
	}
	if slot == 0 || int(slot) > ShopVisitors {
		return ErrMiniRoomInvalid
	}
	visitor := ps.visitors[slot-1]
	if visitor == nil || visitor.GetName() != name {
		return ErrMiniRoomInvalid
	}

	ps.visitors[slot-1] = nil
	visitor.miniRoom = nil
	ps.banned[name] = struct{}{}
	visitor.listener.OnMiniRoomLeft(visitor, slot, reason)
	for _, member := range ps.Members() {
		member.listener.OnMiniRoomLeft(member, slot, pconst.MiniRoomLeaveExit)
	}
	ps.updateBalloon()
	return nil
}

func (ps *PersonalShop) Ban(ch *Character, names []string) error {
	if ps.owner != ch {
		return ErrMiniRoomNotOwner
	}

	for _, name := range names {
		ps.banned[name] = struct{}{}
	}
	return nil
}

func (ps *PersonalShop) returnItems(ch *Character) {
	items := ps.items[:0]
	for _, listed := range ps.items {
		if listed.bundles == 0 {
			continue
		}
		var spec ExchangeSpec
		spec.Reward.include(listed.item, listed.count())
		if ch.inventory.hasUniqueItem(listed.item.GetModel()) || spec.Valid(ch, ch.escrow) != ExchangeOK {
			items = append(items, listed)
			continue
		}
		if _, err := ch.inventory.addItemUnchecked(listed.item.Clone(listed.count()), true); err != nil {
			items = append(items, listed)
		}
	}
	ps.items = items
}

func (ps *PersonalShop) close(reason pconst.MiniRoomLeaveReason) {
	m := ps.onMap
	if m == nil {
		return
	}

	for i, visitor := range ps.visitors {
		if visitor == nil {
			continue
		}
		ps.visitors[i] = nil
		visitor.miniRoom = nil
		visitor.listener.OnMiniRoomLeft(visitor, uint8(i+1), reason)
	}
	owner := ps.owner
	ps.returnItems(owner)
	ps.owner = nil
	owner.miniRoom = nil
	if ps.published {
		owner.listener.OnMiniRoomBalloon(owner, nil)
	}
	ps.entries[owner.GetID()] = owner.BuildRPC(ps.gameWorld.GetWorldID())
	m.RemovePersonalShop(ps)
	ps.gameWorld.Systems().Get[*MapSystem]().Visit(m, func(ctx actor.Context) {
		ps.save(ctx, true)
	})
}

func NewPersonalShop(owner *Character, title string, slot int16, itemID uint32) (*PersonalShop, error) {
	m := owner.GetMap()
	if m == nil || owner.miniRoom != nil {
		return nil, &MiniRoomEnterError{Code: pconst.MiniRoomEnterCannotOpen}
	}
	if !m.wz.PersonalShop {
		return nil, &MiniRoomEnterError{Code: pconst.MiniRoomEnterFreeMarket}
	}
	if owner.unsettled.FindByTarget(internal.UnsettledType_UNSETTLED_TYPE_SHOP_OPEN, uint64(owner.accountID)) != nil {
		return nil, ErrMiniRoomBusy
	}
	permit := owner.inventory.GetItem(constant.InventoryTypeCash, slot)
	if permit == nil || permit.GetModel().GetID() != itemID || permit.GetModel().GetTraits()&constant.ItemTraitPersonalShopPermit == 0 {
		return nil, ErrMiniRoomInvalid
	}
	if expiration := permit.GetExpiration(); !expiration.IsZero() && clock.Now().After(expiration) {
		return nil, ErrMiniRoomInvalid
	}
	if title == "" {
		return nil, ErrMiniRoomInvalid
	}
	if err := m.checkShopSpot(owner.position); err != nil {
		return nil, err
	}

	maxItems := uint8(PersonalShopLargeMaxItems)
	if itemID == personalShopBasicPermit {
		maxItems = PersonalShopMaxItems
	}
	ps := &PersonalShop{
		shopRoom: shopRoom{
			kind:      pconst.MiniRoomTypePersonalShop,
			accountID: owner.accountID,
			ownerID:   owner.GetID(),
			ownerName: owner.GetName(),
			mapID:     m.TemplateID(),
			itemID:    itemID,
			title:     title,
			openedAt:  clock.Now(),
			entries:   make(map[uint32]*internal.CharacterSaveEntry),
		},
		maxItems: maxItems,
		banned:   make(map[string]struct{}),
	}
	ps.objectCore.self = ps
	ps.position = owner.position
	ps.gameWorld = owner.gameWorld
	ps.home = m
	return ps, nil
}

func (ps *PersonalShop) Place(actx actor.Context, owner *Character) {
	id := owner.unsettled.Add(internal.UnsettledType_UNSETTLED_TYPE_SHOP_OPEN, uint64(owner.accountID), nil)
	owner.gameWorld.Systems().Get[*InternalSystem]().OpenShopAsync(actx, owner, id, ps.BuildRPC()).Bind(owner).Do(func(v *internal.OpenShopReply) error {
		owner.unsettled.Erase(id)
		if v.GetExisting() != nil || v.GetAborted() || v.GetShopId() == 0 {
			owner.listener.OnMiniRoomEnterFailed(owner, pconst.MiniRoomEnterCannotOpen)
			return nil
		}

		ps.id = v.GetShopId()
		if owner.GetMap() != ps.home || owner.miniRoom != nil || ps.home.checkShopSpot(ps.position) != nil {
			ps.save(actx, true)
			owner.listener.OnMiniRoomEnterFailed(owner, pconst.MiniRoomEnterCannotOpen)
			return nil
		}
		ps.home.AddPersonalShop(ps)
		ps.owner = owner
		owner.miniRoom = ps
		owner.listener.OnPersonalShopEntered(owner, ps)
		return nil
	}).OnError(func(err error) {
		log.Printf("PersonalShop.Place character=%d: %v", owner.GetID(), err)
		owner.listener.OnMiniRoomEnterFailed(owner, pconst.MiniRoomEnterCannotOpen)
		if request := owner.unsettled.Find(id); request != nil {
			owner.storeBank.settleShopOpen(actx, request)
		}

	})
}
