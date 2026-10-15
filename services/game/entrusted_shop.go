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
	"github.com/boyism80/fm/protocol/dto"
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
)

type EntrustedShop struct {
	shopRoom
	soldInform bool
	managing   bool
}

type EntrustedShopView struct {
	ItemID    uint32
	OwnerName string
	Title     string
	Elapsed   time.Duration
	Meso      int32
	Items     []*ShopItem
	Sold      []ShopSale
	Visitors  [ShopVisitors]*dto.Character
}

func (es *EntrustedShop) GetObjectType() constant.ObjectType {
	return constant.ObjectTypeEntrustedShop
}

func (es *EntrustedShop) Is(typ constant.ObjectType) bool {
	return es.GetObjectType().Has(typ)
}

func (es *EntrustedShop) view() EntrustedShopView {
	view := EntrustedShopView{
		ItemID:    es.ItemID,
		OwnerName: es.OwnerName,
		Title:     es.Title,
		Elapsed:   es.Elapsed(),
		Meso:      es.Meso,
		Sold:      slices.Clone(es.Sold),
	}
	for _, listed := range es.Items {
		copied := *listed
		view.Items = append(view.Items, &copied)
	}
	for i, visitor := range es.Visitors {
		if visitor != nil {
			view.Visitors[i] = visitor.BuildS2C()
		}
	}
	return view
}

func (es *EntrustedShop) balloon() response.EntrustedShopBalloon {
	return response.EntrustedShopBalloon{
		SN:     es.OID,
		Title:  es.Title,
		ItemID: es.ItemID,
		Users:  es.users(),
	}
}

func (es *EntrustedShop) SendSpawnSyncToViewer(viewer *Character) {
	if !es.published {
		return
	}
	footholdID := uint16(0)
	if foothold, _, ok := es.Map.Wz.FootholdBelow(es.Position); ok {
		footholdID = uint16(foothold.ID)
	}
	viewer.Send(&response.SpawnEntrustedShop{
		EmployerID:           es.OwnerID,
		X:                    es.Position.X,
		Y:                    es.Position.Y,
		Foothold:             footholdID,
		OwnerName:            es.OwnerName,
		EntrustedShopBalloon: es.balloon(),
	}, types.SEND_POLICY_ENCRYPT)
}

func (es *EntrustedShop) SendDestroySyncToViewer(viewer *Character) {
	if !es.published {
		return
	}
	viewer.Send(&response.DestroyEntrustedShop{EmployerID: es.OwnerID}, types.SEND_POLICY_ENCRYPT)
}

func (es *EntrustedShop) updateBalloon() {
	if !es.published || es.Map == nil {
		return
	}
	es.Broadcast(&response.UpdateEntrustedShop{
		EmployerID:           es.OwnerID,
		EntrustedShopBalloon: es.balloon(),
	}, nil)
}

func (es *EntrustedShop) accepting() bool {
	return es.published && !es.managing && es.Map != nil
}

func (es *EntrustedShop) release(member *Character, slot uint8, reason pconst.MiniRoomLeaveReason) {
	es.GameWorld.Systems().Get[*DispatchSystem]().Visit(member.GetID(), func(ctx actor.Context, c *Character) {
		if c == nil || c.MiniRoom != es {
			return
		}
		c.MiniRoom = nil
		c.Listener.OnMiniRoomLeft(c, slot, reason)
	})
}

func (es *EntrustedShop) notifyOwner(notify func(owner *Character)) {
	es.GameWorld.Systems().Get[*DispatchSystem]().Visit(es.OwnerID, func(ctx actor.Context, c *Character) {
		if c == nil || c.MiniRoom != es {
			return
		}
		notify(c)
	})
}

func (es *EntrustedShop) payOwner(meso int32, items []*ShopItem, closing bool, notify func(owner *Character)) {
	worldID := es.GameWorld.GetWorldID()
	es.GameWorld.Systems().Get[*DispatchSystem]().Visit(es.OwnerID, func(ctx actor.Context, c *Character) {
		var entry *internal.CharacterSaveEntry
		if c != nil {
			meso, items = c.Inventory.collectShopGoods(meso, items)
			entry = c.BuildRPC(worldID)
			if c.MiniRoom == es {
				notify(c)
			}
			if meso > 0 || len(items) > 0 {
				c.Listener.OnMessage(c, constant.MsgPopup, entrustedShopKeptMessage)
			}
		}

		es.GameWorld.Systems().Get[*MapSystem]().Visit(es.home, func(ctx actor.Context) {
			if entry != nil {
				es.entries[es.OwnerID] = entry
			}
			if meso > 0 || len(items) > 0 {
				es.storeBank = append(es.storeBank, es.buildStoreBankRPC(meso, items))
			}
			es.save(ctx, closing)
		})
	})
}

func (es *EntrustedShop) Elapsed() time.Duration {
	return clock.Now().Sub(es.OpenedAt)
}

func (es *EntrustedShop) Visit(ch *Character) error {
	if ch.MiniRoom != nil {
		return ErrMiniRoomInvalid
	}
	if !es.published || es.Map == nil {
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
		member.Listener.OnMiniRoomVisited(member, slot, ch)
	}
	es.Visitors[index] = ch
	ch.MiniRoom = es
	es.updateBalloon()
	ch.Listener.OnMiniRoomEntered(ch, slot, es.view(), false)
	return nil
}

func (es *EntrustedShop) Manage(ownerID uint32) {
	if !es.accepting() || es.OwnerID != ownerID {
		es.GameWorld.Systems().Get[*DispatchSystem]().Visit(ownerID, func(ctx actor.Context, c *Character) {
			if c == nil {
				return
			}
			c.Listener.OnMiniRoomEnterFailed(c, pconst.MiniRoomEnterClosed)
		})
		return
	}

	for i, visitor := range es.Visitors {
		if visitor == nil {
			continue
		}
		es.Visitors[i] = nil
		es.release(visitor, uint8(i+1), pconst.MiniRoomLeaveOrganizing)
	}
	es.managing = true
	es.updateBalloon()

	view := es.view()
	es.GameWorld.Systems().Get[*DispatchSystem]().Visit(ownerID, func(ctx actor.Context, c *Character) {
		if c != nil && c.MiniRoom == nil {
			c.MiniRoom = es
			c.Listener.OnMiniRoomEntered(c, 0, view, false)
			return
		}
		es.GameWorld.Systems().Get[*MapSystem]().Visit(es.home, es.stopManaging)
	})
}

func (es *EntrustedShop) stopManaging(ctx actor.Context) {
	if !es.managing {
		return
	}

	es.managing = false
	if es.published {
		es.updateBalloon()
		return
	}
	es.close(ctx)
}

func (es *EntrustedShop) Leave(ch *Character) {
	if ch.GetID() == es.OwnerID {
		ch.MiniRoom = nil
		es.GameWorld.Systems().Get[*MapSystem]().Visit(es.home, es.stopManaging)
		return
	}

	slot, ok := es.SlotOf(ch)
	if !ok {
		return
	}
	ch.MiniRoom = nil
	es.Visitors[slot-1] = nil
	for _, member := range es.Members() {
		member.Listener.OnMiniRoomLeft(member, slot, pconst.MiniRoomLeaveExit)
	}
	es.updateBalloon()
}

func (es *EntrustedShop) Chat(ch *Character, message string) error {
	if ch.GetID() == es.OwnerID {
		ch.Listener.OnMiniRoomChat(ch, 0, fmt.Sprintf("%s : %s", ch.GetName(), message))
		return nil
	}
	return es.shopRoom.Chat(ch, message)
}

func (es *EntrustedShop) AddItem(actx actor.Context, ch *Character, invType constant.InventoryType, slot int16, bundles uint16, perBundle uint16, price int32) error {
	if ch.GetID() != es.OwnerID {
		return ErrMiniRoomNotOwner
	}
	listed, err := ch.Inventory.takeShopItem(invType, slot, bundles, perBundle, price)
	if err != nil {
		return err
	}

	entry := ch.BuildRPC(es.GameWorld.GetWorldID())
	es.GameWorld.Systems().Get[*MapSystem]().Visit(es.home, func(ctx actor.Context) {
		if !es.managing || len(es.Items) >= EntrustedShopMaxItems || es.holds(listed.Item.GetModel().GetID()) {
			es.payOwner(0, []*ShopItem{listed}, false, func(owner *Character) {
				owner.Listener.OnUnlockAction(owner)
			})
			return
		}

		es.Items = append(es.Items, listed)
		es.entries[es.OwnerID] = entry
		es.save(ctx, false)
		view := es.view()
		es.notifyOwner(func(owner *Character) {
			owner.Listener.OnMiniRoomItems(owner, view)
		})
	})
	return nil
}

func (es *EntrustedShop) RemoveItem(actx actor.Context, ch *Character, index uint16) error {
	if ch.GetID() != es.OwnerID {
		return ErrMiniRoomNotOwner
	}

	es.GameWorld.Systems().Get[*MapSystem]().Visit(es.home, func(ctx actor.Context) {
		if !es.managing || int(index) >= len(es.Items) {
			es.notifyOwner(func(owner *Character) {
				owner.Listener.OnUnlockAction(owner)
			})
			return
		}

		listed := es.Items[index]
		es.Items = append(es.Items[:index:index], es.Items[index+1:]...)
		var items []*ShopItem
		if listed.Bundles > 0 {
			items = append(items, listed)
		}
		count := uint8(len(es.Items))
		es.payOwner(0, items, false, func(owner *Character) {
			owner.Listener.OnMiniRoomItemRemoved(owner, count, index)
		})
	})
	return nil
}

func (es *EntrustedShop) Open(actx actor.Context, ch *Character) error {
	if ch.GetID() != es.OwnerID {
		return ErrMiniRoomNotOwner
	}

	es.GameWorld.Systems().Get[*MapSystem]().Visit(es.home, func(ctx actor.Context) {
		if !es.managing || es.published || len(es.Items) == 0 || es.Map == nil {
			es.notifyOwner(func(owner *Character) {
				owner.Listener.OnUnlockAction(owner)
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
			owner.MiniRoom = nil
		})
	})
	return nil
}

func (es *EntrustedShop) scheduleClose() {
	remaining := max(EntrustedShopDuration-es.Elapsed(), 0)
	es.AddTimer(entrustedShopCloseTimer, remaining, false, func() {
		es.GameWorld.Systems().Get[*MapSystem]().Visit(es.home, func(ctx actor.Context) {
			if es.Map == nil {
				return
			}
			for i, visitor := range es.Visitors {
				if visitor == nil {
					continue
				}
				es.release(visitor, uint8(i+1), pconst.MiniRoomLeaveTimeUp)
			}
			es.Visitors = [ShopVisitors]*Character{}
			if es.managing {
				es.managing = false
				es.notifyOwner(func(owner *Character) {
					owner.MiniRoom = nil
					owner.Listener.OnMiniRoomLeft(owner, 0, pconst.MiniRoomLeaveTimeUp)
				})
			}
			es.close(ctx)
		})
	})
}

func (es *EntrustedShop) EndMaintenance(ch *Character) error {
	if ch.GetID() != es.OwnerID {
		return ErrMiniRoomNotOwner
	}

	es.Leave(ch)
	return nil
}

func (es *EntrustedShop) Buy(actx actor.Context, ch *Character, index uint16, bundles uint16) error {
	if _, ok := es.SlotOf(ch); !ok || !es.accepting() {
		return ErrMiniRoomInvalid
	}

	listed, total, income, err := es.sell(ch, index, bundles, es.Meso)
	if err != nil {
		return err
	}
	es.Meso += income
	model := listed.Item.GetModel()
	es.Sold = append(es.Sold, ShopSale{
		ItemID:  model.GetID(),
		Bundles: bundles,
		Total:   total,
		Buyer:   ch.GetName(),
	})
	if len(es.Sold) > math.MaxUint8 {
		es.Sold = es.Sold[len(es.Sold)-math.MaxUint8:]
	}
	view := es.view()
	for _, member := range es.Members() {
		member.Listener.OnMiniRoomItems(member, view)
	}
	es.save(actx, false, ch)

	if !es.soldInform {
		return nil
	}
	message := fmt.Sprintf("고용상점에서 %s %d개가 판매되었습니다.", es.GameWorld.GetResources().GetItemName(model.GetID()), bundles*listed.PerBundle)
	es.GameWorld.Systems().Get[*DispatchSystem]().Visit(es.OwnerID, func(ctx actor.Context, owner *Character) {
		if owner == nil {
			return
		}
		owner.Listener.OnMessage(owner, constant.MsgLightBlueText, message)
	})
	return nil
}

func (es *EntrustedShop) Arrange(actx actor.Context, ch *Character) error {
	if ch.GetID() != es.OwnerID {
		return ErrMiniRoomNotOwner
	}

	es.GameWorld.Systems().Get[*MapSystem]().Visit(es.home, func(ctx actor.Context) {
		if !es.managing {
			es.notifyOwner(func(owner *Character) {
				owner.Listener.OnUnlockAction(owner)
			})
			return
		}

		meso := es.Meso
		es.Meso = 0
		items := es.Items[:0]
		for _, listed := range es.Items {
			if listed.Bundles > 0 {
				items = append(items, listed)
			}
		}
		es.Items = items
		view := es.view()
		es.payOwner(meso, nil, false, func(owner *Character) {
			owner.Listener.OnMiniRoomArranged(owner, view)
		})
	})
	return nil
}

func (es *EntrustedShop) WithdrawMeso(actx actor.Context, ch *Character) error {
	if ch.GetID() != es.OwnerID {
		return ErrMiniRoomNotOwner
	}

	es.GameWorld.Systems().Get[*MapSystem]().Visit(es.home, func(ctx actor.Context) {
		if !es.managing || es.Meso <= 0 {
			es.notifyOwner(func(owner *Character) {
				owner.Listener.OnUnlockAction(owner)
			})
			return
		}

		meso := es.Meso
		es.Meso = 0
		es.payOwner(meso, nil, false, func(owner *Character) {
			owner.Listener.OnMiniRoomMesoWithdrawn(owner)
		})
	})
	return nil
}

func (es *EntrustedShop) Close(actx actor.Context, ch *Character) error {
	if ch.GetID() != es.OwnerID {
		return ErrMiniRoomNotOwner
	}

	es.GameWorld.Systems().Get[*MapSystem]().Visit(es.home, func(ctx actor.Context) {
		if !es.managing || es.Map == nil {
			es.notifyOwner(func(owner *Character) {
				owner.Listener.OnUnlockAction(owner)
			})
			return
		}

		meso := es.Meso
		var items []*ShopItem
		for _, listed := range es.Items {
			if listed.Bundles > 0 {
				items = append(items, listed)
			}
		}
		es.Meso = 0
		es.Items = nil
		es.managing = false
		es.Map.RemoveEntrustedShop(es)
		es.payOwner(meso, items, true, func(owner *Character) {
			owner.MiniRoom = nil
			owner.Listener.OnMiniRoomClosed(owner, pconst.MiniRoomCloseAll)
		})
	})
	return nil
}

func (es *EntrustedShop) close(actx actor.Context) {
	m := es.Map
	if m == nil {
		return
	}
	m.RemoveEntrustedShop(es)
	es.save(actx, true)
}

func NewEntrustedShop(owner *Character, title string, slot int16, itemID uint32) (*EntrustedShop, error) {
	m := owner.GetMap()
	if m == nil || !m.Wz.EntrustedShop || owner.MiniRoom != nil {
		return nil, &MiniRoomEnterError{Code: pconst.MiniRoomEnterCannotOpen}
	}
	if owner.Unsettled.FindByTarget(internal.UnsettledType_UNSETTLED_TYPE_SHOP_OPEN, uint64(owner.AccountID)) != nil {
		return nil, ErrMiniRoomBusy
	}
	permit := owner.Inventory.GetItem(constant.InventoryTypeCash, slot)
	if permit == nil || permit.GetModel().GetID() != itemID || permit.GetModel().GetTraits()&constant.ItemTraitEntrustedShopPermit == 0 {
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

	soldInform := false
	if model, ok := permit.GetModel().(*wz.CashItem); ok {
		soldInform = model.SoldInform
	}
	es := &EntrustedShop{
		kind:       pconst.MiniRoomTypeEntrustedShop,
		AccountID:  owner.AccountID,
		OwnerID:    owner.GetID(),
		OwnerName:  owner.GetName(),
		MapID:      m.TemplateID(),
		ItemID:     itemID,
		Title:      title,
		OpenedAt:   clock.Now(),
		entries:    make(map[uint32]*internal.CharacterSaveEntry),
		soldInform: soldInform,
	}
	es.ObjectCore.self = es
	es.Position = owner.Position
	es.GameWorld = owner.GameWorld
	es.home = m
	return es, nil
}

func (es *EntrustedShop) Place(actx actor.Context, owner *Character) {
	id := owner.Unsettled.Add(internal.UnsettledType_UNSETTLED_TYPE_SHOP_OPEN, uint64(owner.AccountID), nil)
	owner.GameWorld.Systems().Get[*InternalSystem]().OpenShopAsync(actx, owner, id, es.BuildRPC()).Do(func(v *internal.OpenShopReply) error {
		owner.Resume(func(actx actor.Context) {
			owner.Unsettled.Erase(id)
			if v.GetExisting() != nil || v.GetAborted() || v.GetShopId() == 0 {
				owner.Listener.OnMiniRoomEnterFailed(owner, pconst.MiniRoomEnterCannotOpen)
				return
			}

			es.ID = v.GetShopId()
			if owner.GetMap() != es.home || owner.MiniRoom != nil || es.home.checkShopSpot(es.Position) != nil {
				es.save(actx, true)
				owner.Listener.OnMiniRoomEnterFailed(owner, pconst.MiniRoomEnterCannotOpen)
				return
			}
			es.home.AddEntrustedShop(es)
			es.managing = true
			owner.MiniRoom = es
			owner.Listener.OnMiniRoomEntered(owner, 0, es.view(), true)
		})
		return nil
	}).OnError(func(err error) {
		log.Printf("EntrustedShop.Place character=%d: %v", owner.GetID(), err)
		owner.Resume(func(actx actor.Context) {
			owner.Listener.OnMiniRoomEnterFailed(owner, pconst.MiniRoomEnterCannotOpen)
			if request := owner.Unsettled.Find(id); request != nil {
				owner.StoreBank.settleShopOpen(actx, request)
			}
		})
	})
}
