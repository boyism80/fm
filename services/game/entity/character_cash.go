package entity

import (
	"fmt"
	"log"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/core/luax"
	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
	"github.com/boyism80/fm/services/game/constant"
	lua "github.com/yuin/gopher-lua"
)

func (ch *Character) UseCashItem(ctx actor.Context, slot int16, itemID uint32, text string, ear bool, petSN uint64) {
	if ch.session.cashItemInUse.CompareAndSwap(false, true) == false {
		ch.Listener.OnUpdateStats(ch, nil, true)
		return
	}
	finish := func() {
		ch.session.cashItemInUse.Store(false)
		ch.Listener.OnUpdateStats(ch, nil, true)
	}

	if ch.GetHp() <= 0 {
		finish()
		return
	}
	cash := ch.Inventory.Tabs[constant.InventoryTypeCash]
	if cash == nil {
		finish()
		return
	}
	item := cash.Get(uint8(slot))
	if item == nil || item.GetCount() < 1 || item.GetModel().GetID() != itemID {
		finish()
		return
	}
	m := ch.GetMap()
	if m == nil || m.GetLuaRoot() == nil {
		finish()
		return
	}

	scriptPath := fmt.Sprintf("script/item/%d.lua", itemID)
	thread, err := luax.NewThread(m.GetLuaRoot(), scriptPath)
	if err != nil {
		finish()
		return
	}
	luax.SetConfiguration(thread, luax.Configuration{ActorContext: ctx})
	luax.CallAsync(ctx, m.GetLuaRoot(), thread, "on_cash", ch, itemID, text, ear, petSN).Do(func(vals []lua.LValue) error {
		if len(vals) == 0 || vals[0] != lua.LTrue {
			finish()
			return nil
		}
		ch.GameWorld.GetDispatchSystem().Call(ch.GetID(), func(actor.Context) {
			item := cash.Get(uint8(slot))
			if item != nil && item.GetModel().GetID() == itemID {
				ch.Inventory.RemoveItem(constant.InventoryTypeCash, slot, 1)
			}
			finish()
		})
		return nil
	}).OnError(func(err error) {
		log.Printf("cash item script %s: %v", scriptPath, err)
		finish()
	})
}

func (ch *Character) AddCash(actx actor.Context, nxCash int32, maplePoint int32) {
	ch.Listener.AddCashAsync(actx, ch, nxCash, maplePoint).Do(func(v *internal.AddCashReply) error {
		reply := v
		ch.Message(fmt.Sprintf("캐시 잔액: NX %d, 메이플포인트 %d", reply.GetNxCash(), reply.GetMaplePoint()))
		return nil
	}).OnError(func(err error) {
		log.Printf("Character.AddCash character=%d: %v", ch.GetID(), err)
	})
}

func (ch *Character) CreateCashCoupons(actx actor.Context, kind internal.CashCouponKind, value uint32, count uint32) {
	ch.Listener.CreateCashCouponsAsync(actx, ch, kind, value, count).Do(func(v *internal.CreateCashCouponsReply) error {
		for _, code := range v.GetCodes() {
			ch.Message(fmt.Sprintf("쿠폰: %s", code))
		}
		return nil
	}).OnError(func(err error) {
		log.Printf("Character.CreateCashCoupons character=%d: %v", ch.GetID(), err)
	})
}
