package game

import (
	"fmt"
	"log"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/core/luax"
	"github.com/boyism80/fm/maple/constant"
	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
	lua "github.com/yuin/gopher-lua"
)

// UsingCashItem reports a cash item or teleport stone whose cost is taken
// after an asynchronous step, so its slot must not move in the meantime.
func (ch *Character) UsingCashItem() bool {
	return ch.session.cashItemInUse.Load()
}

type CashItemInput struct {
	Text         string
	Ear          bool
	PetSN        uint64
	Recipient    string
	Lines        []string
	Attached     bool
	AttachedType uint32
	AttachedSlot uint32
}

func (ch *Character) UseCashItem(ctx actor.Context, slot int16, itemID uint32, input CashItemInput) {
	if !ch.session.cashItemInUse.CompareAndSwap(false, true) {
		ch.listener.OnUpdateStats(ch, nil, true)
		return
	}
	finish := func() {
		ch.session.cashItemInUse.Store(false)
		ch.listener.OnUpdateStats(ch, nil, true)
	}

	if ch.Occupied() {
		finish()
		return
	}
	cash := ch.inventory.tabs[constant.InventoryTypeCash]
	if cash == nil {
		finish()
		return
	}
	item := ch.inventory.GetItem(constant.InventoryTypeCash, slot)
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
		ch.listener.OnMessage(ch, constant.MsgPinkText, "지금은 사용할 수 없는 아이템입니다.")
		finish()
		return
	}
	luax.SetConfiguration(thread, luax.Configuration{ActorContext: ctx})

	extra := thread.NewTable()
	extra.RawSetString("recipient", lua.LString(input.Recipient))
	lines := thread.NewTable()
	for _, line := range input.Lines {
		lines.Append(lua.LString(line))
	}
	extra.RawSetString("lines", lines)
	if input.Attached {
		extra.RawSetString("inventory", lua.LNumber(input.AttachedType))
		extra.RawSetString("slot", lua.LNumber(input.AttachedSlot))
	}
	// The item is taken before the script can wait on another server, so a save during the wait (logout, channel change) never keeps it.
	ch.inventory.RemoveItem(constant.InventoryTypeCash, slot, 1)
	refund := func() {
		ch.gameWorld.Systems().Get[*DispatchSystem]().Visit(ch.GetID(), func(actor.Context, *Character) {
			ch.inventory.restoreItem(constant.InventoryTypeCash, slot, item)
			finish()
		})
	}
	luax.CallAsync(ctx, m.GetLuaRoot(), thread, "on_cash", ch, itemID, input.Text, input.Ear, input.PetSN, extra).Do(func(vals []lua.LValue) error {
		if len(vals) == 0 || vals[0] != lua.LTrue {
			refund()
			return nil
		}
		ch.gameWorld.Systems().Get[*DispatchSystem]().Visit(ch.GetID(), func(actor.Context, *Character) {
			finish()
		})
		return nil
	}).OnError(func(err error) {
		log.Printf("cash item script %s: %v", scriptPath, err)
		refund()
	})
}

func (ch *Character) AddCash(actx actor.Context, nxCash int32, maplePoint int32) {
	ch.gameWorld.Systems().Get[*InternalSystem]().AddCashAsync(actx, ch, nxCash, maplePoint).Do(func(v *internal.AddCashReply) error {
		reply := v
		ch.Message(fmt.Sprintf("캐시 잔액: NX %d, 메이플포인트 %d", reply.GetNxCash(), reply.GetMaplePoint()))
		return nil
	}).OnError(func(err error) {
		log.Printf("Character.AddCash character=%d: %v", ch.GetID(), err)
	})
}

func (ch *Character) CreateCashCoupons(actx actor.Context, kind internal.CashCouponKind, value uint32, count uint32) {
	ch.gameWorld.Systems().Get[*InternalSystem]().CreateCashCouponsAsync(actx, ch, kind, value, count).Do(func(v *internal.CreateCashCouponsReply) error {
		for _, code := range v.GetCodes() {
			ch.Message(fmt.Sprintf("쿠폰: %s", code))
		}
		return nil
	}).OnError(func(err error) {
		log.Printf("Character.CreateCashCoupons character=%d: %v", ch.GetID(), err)
	})
}
