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
	MaxItems uint8
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
		SN:       ps.OID,
		Title:    ps.Title,
		Spec:     uint8(ps.ItemID % 10),
		Users:    ps.users(),
		MaxUsers: pconst.MiniRoomShopUsers,
	}
}

func (ps *PersonalShop) updateBalloon() {
	if !ps.published {
		return
	}
	balloon := ps.balloon()
	ps.owner.Listener.OnMiniRoomBalloon(ps.owner, &balloon)
}

func (ps *PersonalShop) Visit(ch *Character) error {
	if ch.MiniRoom != nil {
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
		member.Listener.OnMiniRoomVisited(member, slot, ch)
	}
	ps.Visitors[index] = ch
	ch.MiniRoom = ps
	ps.updateBalloon()
	ch.Listener.OnPersonalShopEntered(ch, ps)
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

	ch.MiniRoom = nil
	ps.Visitors[slot-1] = nil
	for _, member := range ps.Members() {
		member.Listener.OnMiniRoomLeft(member, slot, pconst.MiniRoomLeaveExit)
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
	if len(ps.Items) >= int(ps.MaxItems) {
		return ErrMiniRoomFull
	}

	if err := ps.list(ch, invType, slot, bundles, perBundle, price); err != nil {
		return err
	}
	ch.Listener.OnPersonalShopItems(ch, ps)
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
	ch.Listener.OnMiniRoomItemRemoved(ch, uint8(len(ps.Items)), index)
	for _, visitor := range ps.Visitors {
		if visitor != nil {
			visitor.Listener.OnPersonalShopItems(visitor, ps)
		}
	}
	ps.save(actx, false, ch)
	return nil
}

func (ps *PersonalShop) Open(actx actor.Context, ch *Character) error {
	if ps.owner != ch || ps.published {
		return ErrMiniRoomNotOwner
	}
	if len(ps.Items) == 0 {
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
	_, _, income, err := ps.sell(ch, index, bundles, owner.Inventory.Meso)
	if err != nil {
		return err
	}
	owner.Inventory.addMesoUnchecked(income)
	for _, member := range ps.Members() {
		member.Listener.OnPersonalShopItems(member, ps)
	}
	owner.Listener.OnMiniRoomSold(owner, uint8(index), bundles, ch.GetName())
	ps.save(actx, false, ch, owner)

	for _, listed := range ps.Items {
		if listed.Bundles > 0 {
			return nil
		}
	}
	owner.Listener.OnMiniRoomLeft(owner, 0, pconst.MiniRoomLeaveSoldOut)
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
	visitor := ps.Visitors[slot-1]
	if visitor == nil || visitor.GetName() != name {
		return ErrMiniRoomInvalid
	}

	ps.Visitors[slot-1] = nil
	visitor.MiniRoom = nil
	ps.banned[name] = struct{}{}
	visitor.Listener.OnMiniRoomLeft(visitor, slot, reason)
	for _, member := range ps.Members() {
		member.Listener.OnMiniRoomLeft(member, slot, pconst.MiniRoomLeaveExit)
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
	items := ps.Items[:0]
	for _, listed := range ps.Items {
		if listed.Bundles == 0 {
			continue
		}
		model := listed.Item.GetModel()
		spec := ExchangeSpec{Reward: ExchangeSide{Items: map[uint32]uint16{model.GetID(): listed.count()}}}
		if (model.GetTraits()&constant.ItemTraitOnly != 0 && ch.Inventory.HasItem(model.GetID())) || spec.Valid(ch) != ExchangeOK {
			items = append(items, listed)
			continue
		}
		ch.Inventory.addItemUnchecked(listed.Item.Clone(listed.count()), true)
	}
	ps.Items = items
}

func (ps *PersonalShop) close(reason pconst.MiniRoomLeaveReason) {
	m := ps.Map
	if m == nil {
		return
	}

	for i, visitor := range ps.Visitors {
		if visitor == nil {
			continue
		}
		ps.Visitors[i] = nil
		visitor.MiniRoom = nil
		visitor.Listener.OnMiniRoomLeft(visitor, uint8(i+1), reason)
	}
	owner := ps.owner
	ps.returnItems(owner)
	ps.owner = nil
	owner.MiniRoom = nil
	if ps.published {
		owner.Listener.OnMiniRoomBalloon(owner, nil)
	}
	ps.entries[owner.GetID()] = owner.BuildRPC(ps.GameWorld.GetWorldID())
	m.RemovePersonalShop(ps)
	ps.GameWorld.Systems().Get[*MapSystem]().Call(m, func(ctx actor.Context) {
		ps.save(ctx, true)
	})
}

func NewPersonalShop(owner *Character, title string, slot int16, itemID uint32) (*PersonalShop, error) {
	m := owner.GetMap()
	if m == nil || owner.MiniRoom != nil {
		return nil, &MiniRoomEnterError{Code: pconst.MiniRoomEnterCannotOpen}
	}
	if !m.Wz.PersonalShop {
		return nil, &MiniRoomEnterError{Code: pconst.MiniRoomEnterFreeMarket}
	}
	if owner.Unsettled.FindByTarget(internal.UnsettledType_UNSETTLED_TYPE_SHOP_OPEN, uint64(owner.AccountID)) != nil {
		return nil, ErrMiniRoomBusy
	}
	permit := owner.Inventory.GetItem(constant.InventoryTypeCash, slot)
	if permit == nil || permit.GetModel().GetID() != itemID || permit.GetModel().GetTraits()&constant.ItemTraitPersonalShopPermit == 0 {
		return nil, ErrMiniRoomInvalid
	}
	if expiration := permit.GetExpiration(); !expiration.IsZero() && clock.Now().After(expiration) {
		return nil, ErrMiniRoomInvalid
	}
	if title == "" {
		return nil, ErrMiniRoomInvalid
	}
	if err := m.checkShopSpot(owner.Position); err != nil {
		return nil, err
	}

	maxItems := uint8(PersonalShopLargeMaxItems)
	if itemID == personalShopBasicPermit {
		maxItems = PersonalShopMaxItems
	}
	ps := &PersonalShop{
		shopRoom: shopRoom{
			kind:      pconst.MiniRoomTypePersonalShop,
			AccountID: owner.AccountID,
			OwnerID:   owner.GetID(),
			OwnerName: owner.GetName(),
			MapID:     m.TemplateID(),
			ItemID:    itemID,
			Title:     title,
			OpenedAt:  clock.Now(),
			entries:   make(map[uint32]*internal.CharacterSaveEntry),
		},
		MaxItems: maxItems,
		banned:   make(map[string]struct{}),
	}
	ps.ObjectCore.self = ps
	ps.Position = owner.Position
	ps.GameWorld = owner.GameWorld
	ps.home = m
	return ps, nil
}

func (ps *PersonalShop) Place(actx actor.Context, owner *Character) {
	id := owner.Unsettled.Add(internal.UnsettledType_UNSETTLED_TYPE_SHOP_OPEN, uint64(owner.AccountID), nil)
	owner.GameWorld.Systems().Get[*InternalSystem]().OpenShopAsync(actx, owner, id, ps.BuildRPC()).Do(func(v *internal.OpenShopReply) error {
		owner.Call(func(actx actor.Context) {
			owner.Unsettled.Erase(id)
			if v.GetExisting() != nil || v.GetAborted() || v.GetShopId() == 0 {
				owner.Listener.OnMiniRoomEnterFailed(owner, pconst.MiniRoomEnterCannotOpen)
				return
			}

			ps.ID = v.GetShopId()
			if owner.GetMap() != ps.home || owner.MiniRoom != nil || ps.home.checkShopSpot(ps.Position) != nil {
				ps.save(actx, true)
				owner.Listener.OnMiniRoomEnterFailed(owner, pconst.MiniRoomEnterCannotOpen)
				return
			}
			ps.home.AddPersonalShop(ps)
			ps.owner = owner
			owner.MiniRoom = ps
			owner.Listener.OnPersonalShopEntered(owner, ps)
		})
		return nil
	}).OnError(func(err error) {
		log.Printf("PersonalShop.Place character=%d: %v", owner.GetID(), err)
		owner.Call(func(actx actor.Context) {
			owner.Listener.OnMiniRoomEnterFailed(owner, pconst.MiniRoomEnterCannotOpen)
			if request := owner.Unsettled.Find(id); request != nil {
				owner.StoreBank.settleShopOpen(actx, request)
			}
		})
	})
}
