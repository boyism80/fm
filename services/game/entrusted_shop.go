package game

import (
	"fmt"
	"log"
	"math"
	"slices"
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/core/clock"
	"github.com/boyism80/fm/maple/constant"
	"github.com/boyism80/fm/maple/wz"
	pconst "github.com/boyism80/fm/protocol/constant"
	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
	"github.com/boyism80/fm/protocol/response"
	"github.com/boyism80/fm/types"
)

const (
	EntrustedShopMaxItems     = 10
	EntrustedShopDuration     = 24 * time.Hour
	RemoteEntrustedShopItemID = 5470000
	entrustedShopCloseTimer   = "entrusted_shop_close"
	entrustedShopKeptMessage  = "인벤토리에 공간이 부족해 받지 못한 물품은 프레드릭에게 보관되었습니다."
	entrustedShopEmptyMessage = "프레드릭에게서 아이템을 찾아가시기 바랍니다."
)

type EntrustedShop struct {
	shopRoom
	soldInform bool
	managing   bool
}

type EntrustedShopView struct {
	itemID    uint32
	ownerName string
	title     string
	elapsed   time.Duration
	meso      int32
	items     []*ShopItem
	sold      []ShopSale
	visitors  [ShopVisitors]*Character
}

func (es *EntrustedShop) GetObjectType() constant.ObjectType {
	return constant.ObjectTypeEntrustedShop
}

func (es *EntrustedShop) Is(typ constant.ObjectType) bool {
	return es.GetObjectType().Has(typ)
}

func (es *EntrustedShop) view() EntrustedShopView {
	view := EntrustedShopView{
		itemID:    es.itemID,
		ownerName: es.ownerName,
		title:     es.title,
		elapsed:   es.Elapsed(),
		meso:      es.meso,
		sold:      slices.Clone(es.sold),
	}
	for _, listed := range es.items {
		copied := *listed
		view.items = append(view.items, &copied)
	}
	for i, visitor := range es.visitors {
		if visitor != nil {
			view.visitors[i] = visitor
		}
	}
	return view
}

func (es *EntrustedShop) balloon() response.EntrustedShopBalloon {
	return response.EntrustedShopBalloon{
		SN:     es.oid,
		Title:  es.title,
		ItemID: es.itemID,
		Users:  es.users(),
	}
}

func (es *EntrustedShop) SendSpawnSyncToViewer(viewer *Character) {
	if !es.published {
		return
	}
	footholdID := uint16(0)
	if foothold, _, ok := es.onMap.wz.FootholdBelow(es.position); ok {
		footholdID = uint16(foothold.ID)
	}
	viewer.Send(&response.SpawnEntrustedShop{
		EmployerID:           es.ownerID,
		X:                    es.position.X,
		Y:                    es.position.Y,
		Foothold:             footholdID,
		OwnerName:            es.ownerName,
		EntrustedShopBalloon: es.balloon(),
	}, types.SEND_POLICY_ENCRYPT)
}

func (es *EntrustedShop) SendDestroySyncToViewer(viewer *Character) {
	if !es.published {
		return
	}
	viewer.Send(&response.DestroyEntrustedShop{EmployerID: es.ownerID}, types.SEND_POLICY_ENCRYPT)
}

func (es *EntrustedShop) updateBalloon() {
	if !es.published || es.onMap == nil {
		return
	}
	es.Broadcast(&response.UpdateEntrustedShop{
		EmployerID:           es.ownerID,
		EntrustedShopBalloon: es.balloon(),
	}, nil)
}

func (es *EntrustedShop) accepting() bool {
	return es.published && !es.managing && es.onMap != nil
}

func (es *EntrustedShop) release(member *Character, slot uint8, reason pconst.MiniRoomLeaveReason) {
	es.gameWorld.Systems().Get[*DispatchSystem]().Visit(member.GetID(), func(ctx actor.Context, c *Character) {
		if c == nil || c.miniRoom != es {
			return
		}
		c.miniRoom = nil
		c.listener.OnMiniRoomLeft(c, slot, reason)
	})
}

func (es *EntrustedShop) notifyOwner(notify func(owner *Character)) {
	es.gameWorld.Systems().Get[*DispatchSystem]().Visit(es.ownerID, func(ctx actor.Context, c *Character) {
		if c == nil || c.miniRoom != es {
			return
		}
		notify(c)
	})
}

func (es *EntrustedShop) payOwner(meso int32, items []*ShopItem, closing bool, notify func(owner *Character)) {
	worldID := es.gameWorld.GetWorldID()
	es.gameWorld.Systems().Get[*DispatchSystem]().Visit(es.ownerID, func(ctx actor.Context, c *Character) {
		var entry *internal.CharacterSaveEntry
		if c != nil {
			meso, items = c.inventory.collectShopGoods(meso, items)
			entry = c.BuildRPC(worldID)
			if c.miniRoom == es {
				notify(c)
			}
			if meso > 0 || len(items) > 0 {
				c.listener.OnMessage(c, constant.MsgPopup, entrustedShopKeptMessage)
			}
		}

		es.gameWorld.Systems().Get[*MapSystem]().Visit(es.home, func(ctx actor.Context) {
			if entry != nil {
				es.entries[es.ownerID] = entry
			}
			if meso > 0 || len(items) > 0 {
				es.storeBank = append(es.storeBank, es.buildStoreBank(meso, items))
			}
			es.save(ctx, closing)
		})
	})
}

func (es *EntrustedShop) Elapsed() time.Duration {
	return clock.Now().Sub(es.openedAt)
}

func (es *EntrustedShop) Visit(ch *Character) error {
	if ch.miniRoom != nil {
		return ErrMiniRoomInvalid
	}
	if !es.published || es.onMap == nil {
		return &MiniRoomEnterError{Code: pconst.MiniRoomEnterClosed}
	}
	if es.managing {
		return &MiniRoomEnterError{Code: pconst.MiniRoomEnterOrganizing}
	}
	index, ok := es.freeSlot()
	if !ok {
		return &MiniRoomEnterError{Code: pconst.MiniRoomEnterFull}
	}

	slot := uint8(index + 1)
	for _, member := range es.Members() {
		member.listener.OnMiniRoomVisited(member, slot, ch)
	}
	es.visitors[index] = ch
	ch.miniRoom = es
	es.updateBalloon()
	ch.listener.OnMiniRoomEntered(ch, slot, es.view(), false)
	return nil
}

func (es *EntrustedShop) Manage(ownerID uint32) {
	if !es.accepting() || es.ownerID != ownerID {
		es.gameWorld.Systems().Get[*DispatchSystem]().Visit(ownerID, func(ctx actor.Context, c *Character) {
			if c == nil {
				return
			}
			c.listener.OnMiniRoomEnterFailed(c, pconst.MiniRoomEnterClosed)
		})
		return
	}

	for i, visitor := range es.visitors {
		if visitor == nil {
			continue
		}
		es.visitors[i] = nil
		es.release(visitor, uint8(i+1), pconst.MiniRoomLeaveOrganizing)
	}
	es.managing = true
	es.updateBalloon()

	view := es.view()
	es.gameWorld.Systems().Get[*DispatchSystem]().Visit(ownerID, func(ctx actor.Context, c *Character) {
		if c != nil && c.miniRoom == nil {
			c.miniRoom = es
			c.listener.OnMiniRoomEntered(c, 0, view, false)
			return
		}
		es.gameWorld.Systems().Get[*MapSystem]().Visit(es.home, es.stopManaging)
	})
}

func (es *EntrustedShop) stopManaging(ctx actor.Context) {
	if !es.managing {
		return
	}

	es.managing = false
	if es.published && len(es.items) > 0 {
		es.updateBalloon()
		return
	}
	emptied := es.published && es.onMap != nil
	es.close(ctx)
	if emptied {
		es.gameWorld.Systems().Get[*DispatchSystem]().Visit(es.ownerID, func(ctx actor.Context, c *Character) {
			if c != nil {
				c.listener.OnMessage(c, constant.MsgPopup, entrustedShopEmptyMessage)
			}
		})
	}
}

func (es *EntrustedShop) Leave(ch *Character) {
	if ch.GetID() == es.ownerID {
		ch.miniRoom = nil
		es.gameWorld.Systems().Get[*MapSystem]().Visit(es.home, es.stopManaging)
		return
	}

	slot, ok := es.SlotOf(ch)
	if !ok {
		return
	}
	ch.miniRoom = nil
	es.visitors[slot-1] = nil
	for _, member := range es.Members() {
		member.listener.OnMiniRoomLeft(member, slot, pconst.MiniRoomLeaveExit)
	}
	es.updateBalloon()
}

func (es *EntrustedShop) Chat(ch *Character, message string) error {
	if ch.GetID() == es.ownerID {
		ch.listener.OnMiniRoomChat(ch, 0, fmt.Sprintf("%s : %s", ch.GetName(), message))
		return nil
	}
	return es.shopRoom.Chat(ch, message)
}

func (es *EntrustedShop) AddItem(actx actor.Context, ch *Character, invType constant.InventoryType, slot int16, bundles uint16, perBundle uint16, price int32) error {
	if ch.GetID() != es.ownerID {
		return ErrMiniRoomNotOwner
	}
	listed, err := ch.inventory.takeShopItem(invType, slot, bundles, perBundle, price)
	if err != nil {
		return err
	}

	entry := ch.BuildRPC(es.gameWorld.GetWorldID())
	es.gameWorld.Systems().Get[*MapSystem]().Visit(es.home, func(ctx actor.Context) {
		if !es.managing || len(es.items) >= EntrustedShopMaxItems || es.holds(listed.item.GetModel().GetID()) {
			es.payOwner(0, []*ShopItem{listed}, false, func(owner *Character) {
				owner.listener.OnUnlockAction(owner)
			})
			return
		}

		es.items = append(es.items, listed)
		es.entries[es.ownerID] = entry
		es.save(ctx, false)
		view := es.view()
		es.notifyOwner(func(owner *Character) {
			owner.listener.OnMiniRoomItems(owner, view)
		})
	})
	return nil
}

func (es *EntrustedShop) RemoveItem(actx actor.Context, ch *Character, index uint16) error {
	if ch.GetID() != es.ownerID {
		return ErrMiniRoomNotOwner
	}

	es.gameWorld.Systems().Get[*MapSystem]().Visit(es.home, func(ctx actor.Context) {
		if !es.managing || int(index) >= len(es.items) {
			es.notifyOwner(func(owner *Character) {
				owner.listener.OnUnlockAction(owner)
			})
			return
		}

		listed := es.items[index]
		es.items = append(es.items[:index:index], es.items[index+1:]...)
		var items []*ShopItem
		if listed.bundles > 0 {
			items = append(items, listed)
		}
		count := uint8(len(es.items))
		es.payOwner(0, items, false, func(owner *Character) {
			owner.listener.OnMiniRoomItemRemoved(owner, count, index)
		})
	})
	return nil
}

func (es *EntrustedShop) Open(actx actor.Context, ch *Character) error {
	if ch.GetID() != es.ownerID {
		return ErrMiniRoomNotOwner
	}

	es.gameWorld.Systems().Get[*MapSystem]().Visit(es.home, func(ctx actor.Context) {
		if !es.managing || es.published || len(es.items) == 0 || es.onMap == nil {
			es.notifyOwner(func(owner *Character) {
				owner.listener.OnUnlockAction(owner)
			})
			return
		}

		es.managing = false
		es.published = true
		es.BroadcastCall(func(obj Object) {
			viewer, ok := obj.(*Character)
			if !ok {
				return
			}
			es.SendSpawnSyncToViewer(viewer)
		}, nil)
		es.scheduleClose()
		es.save(ctx, false)
		es.notifyOwner(func(owner *Character) {
			owner.miniRoom = nil
		})
	})
	return nil
}

func (es *EntrustedShop) scheduleClose() {
	remaining := max(EntrustedShopDuration-es.Elapsed(), 0)
	es.AddTimer(entrustedShopCloseTimer, remaining, false, func() {
		es.gameWorld.Systems().Get[*MapSystem]().Visit(es.home, func(ctx actor.Context) {
			if es.onMap == nil {
				return
			}
			for i, visitor := range es.visitors {
				if visitor == nil {
					continue
				}
				es.release(visitor, uint8(i+1), pconst.MiniRoomLeaveTimeUp)
			}
			es.visitors = [ShopVisitors]*Character{}
			if es.managing {
				es.managing = false
				es.notifyOwner(func(owner *Character) {
					owner.miniRoom = nil
					owner.listener.OnMiniRoomLeft(owner, 0, pconst.MiniRoomLeaveTimeUp)
				})
			}
			es.close(ctx)
		})
	})
}

func (es *EntrustedShop) EndMaintenance(ch *Character) error {
	if ch.GetID() != es.ownerID {
		return ErrMiniRoomNotOwner
	}

	es.Leave(ch)
	return nil
}

func (es *EntrustedShop) Buy(actx actor.Context, ch *Character, index uint16, bundles uint16) error {
	if _, ok := es.SlotOf(ch); !ok || !es.accepting() {
		return ErrMiniRoomInvalid
	}

	listed, total, income, err := es.sell(ch, index, bundles, es.meso)
	if err != nil {
		return err
	}
	es.meso += income
	model := listed.item.GetModel()
	es.sold = append(es.sold, ShopSale{
		ItemID:  model.GetID(),
		Bundles: bundles,
		Total:   total,
		Buyer:   ch.GetName(),
	})
	if len(es.sold) > math.MaxUint8 {
		es.sold = es.sold[len(es.sold)-math.MaxUint8:]
	}
	view := es.view()
	for _, member := range es.Members() {
		member.listener.OnMiniRoomItems(member, view)
	}
	es.save(actx, false, ch)

	if !es.soldInform {
		return nil
	}
	message := fmt.Sprintf("고용상점에서 %s %d개가 판매되었습니다.", es.gameWorld.GetResources().GetItemName(model.GetID()), bundles*listed.perBundle)
	es.gameWorld.Systems().Get[*DispatchSystem]().Visit(es.ownerID, func(ctx actor.Context, owner *Character) {
		if owner == nil {
			return
		}
		owner.listener.OnMessage(owner, constant.MsgLightBlueText, message)
	})
	return nil
}

func (es *EntrustedShop) Arrange(actx actor.Context, ch *Character) error {
	if ch.GetID() != es.ownerID {
		return ErrMiniRoomNotOwner
	}

	es.gameWorld.Systems().Get[*MapSystem]().Visit(es.home, func(ctx actor.Context) {
		if !es.managing {
			es.notifyOwner(func(owner *Character) {
				owner.listener.OnUnlockAction(owner)
			})
			return
		}

		meso := es.meso
		es.meso = 0
		items := es.items[:0]
		for _, listed := range es.items {
			if listed.bundles > 0 {
				items = append(items, listed)
			}
		}
		es.items = items
		view := es.view()
		es.payOwner(meso, nil, false, func(owner *Character) {
			owner.listener.OnMiniRoomArranged(owner, view)
		})
	})
	return nil
}

func (es *EntrustedShop) WithdrawMeso(actx actor.Context, ch *Character) error {
	if ch.GetID() != es.ownerID {
		return ErrMiniRoomNotOwner
	}

	es.gameWorld.Systems().Get[*MapSystem]().Visit(es.home, func(ctx actor.Context) {
		if !es.managing || es.meso <= 0 {
			es.notifyOwner(func(owner *Character) {
				owner.listener.OnUnlockAction(owner)
			})
			return
		}

		meso := es.meso
		es.meso = 0
		es.payOwner(meso, nil, false, func(owner *Character) {
			owner.listener.OnMiniRoomMesoWithdrawn(owner)
		})
	})
	return nil
}

func (es *EntrustedShop) Close(actx actor.Context, ch *Character) error {
	if ch.GetID() != es.ownerID {
		return ErrMiniRoomNotOwner
	}

	es.gameWorld.Systems().Get[*MapSystem]().Visit(es.home, func(ctx actor.Context) {
		if !es.managing || es.onMap == nil {
			es.notifyOwner(func(owner *Character) {
				owner.listener.OnUnlockAction(owner)
			})
			return
		}

		meso := es.meso
		var items []*ShopItem
		for _, listed := range es.items {
			if listed.bundles > 0 {
				items = append(items, listed)
			}
		}
		es.meso = 0
		es.items = nil
		es.managing = false
		es.onMap.RemoveEntrustedShop(es)
		es.payOwner(meso, items, true, func(owner *Character) {
			owner.miniRoom = nil
			owner.listener.OnMiniRoomClosed(owner, pconst.MiniRoomCloseAll)
		})
	})
	return nil
}

func (es *EntrustedShop) close(actx actor.Context) {
	m := es.onMap
	if m == nil {
		return
	}
	m.RemoveEntrustedShop(es)
	es.save(actx, true)
}

func NewEntrustedShop(owner *Character, title string, slot int16, itemID uint32) (*EntrustedShop, error) {
	m := owner.GetMap()
	if m == nil || !m.wz.EntrustedShop || owner.miniRoom != nil {
		return nil, &MiniRoomEnterError{Code: pconst.MiniRoomEnterCannotOpen}
	}
	if owner.unsettled.FindByTarget(internal.UnsettledType_UNSETTLED_TYPE_SHOP_OPEN, uint64(owner.accountID)) != nil {
		return nil, ErrMiniRoomBusy
	}
	permit := owner.inventory.GetItem(constant.InventoryTypeCash, slot)
	if permit == nil || permit.GetModel().GetID() != itemID || permit.GetModel().GetTraits()&constant.ItemTraitEntrustedShopPermit == 0 {
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

	soldInform := false
	if model, ok := permit.GetModel().(*wz.CashItem); ok {
		soldInform = model.SoldInform
	}
	es := &EntrustedShop{
		kind:       pconst.MiniRoomTypeEntrustedShop,
		accountID:  owner.accountID,
		ownerID:    owner.GetID(),
		ownerName:  owner.GetName(),
		mapID:      m.TemplateID(),
		itemID:     itemID,
		title:      title,
		openedAt:   clock.Now(),
		entries:    make(map[uint32]*internal.CharacterSaveEntry),
		soldInform: soldInform,
	}
	es.objectCore.self = es
	es.position = owner.position
	es.gameWorld = owner.gameWorld
	es.home = m
	return es, nil
}

func (es *EntrustedShop) Place(actx actor.Context, owner *Character) {
	id := owner.unsettled.Add(internal.UnsettledType_UNSETTLED_TYPE_SHOP_OPEN, uint64(owner.accountID), nil)
	owner.gameWorld.Systems().Get[*InternalSystem]().OpenShopAsync(actx, owner, id, es.BuildRPC()).Bind(owner).Do(func(v *internal.OpenShopReply) error {
		owner.unsettled.Erase(id)
		if v.GetExisting() != nil || v.GetAborted() || v.GetShopId() == 0 {
			owner.listener.OnMiniRoomEnterFailed(owner, pconst.MiniRoomEnterCannotOpen)
			return nil
		}

		es.id = v.GetShopId()
		if owner.GetMap() != es.home || owner.miniRoom != nil || es.home.checkShopSpot(es.position) != nil {
			es.save(actx, true)
			owner.listener.OnMiniRoomEnterFailed(owner, pconst.MiniRoomEnterCannotOpen)
			return nil
		}
		es.home.AddEntrustedShop(es)
		es.managing = true
		owner.miniRoom = es
		owner.listener.OnMiniRoomEntered(owner, 0, es.view(), true)
		return nil
	}).OnError(func(err error) {
		log.Printf("EntrustedShop.Place character=%d: %v", owner.GetID(), err)
		owner.listener.OnMiniRoomEnterFailed(owner, pconst.MiniRoomEnterCannotOpen)
		if request := owner.unsettled.Find(id); request != nil {
			owner.storeBank.settleShopOpen(actx, request)
		}

	})
}
