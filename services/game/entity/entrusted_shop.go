package entity

import (
	"errors"
	"fmt"
	"log"
	"math"
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/core/clock"
	pconst "github.com/boyism80/fm/protocol/constant"
	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
	"github.com/boyism80/fm/protocol/response"
	"github.com/boyism80/fm/services/game/constant"
	"github.com/boyism80/fm/services/game/wz"
	"github.com/boyism80/fm/types"
)

const (
	EntrustedShopMaxItems       = 10
	EntrustedShopVisitors       = 3
	EntrustedShopDuration       = 24 * time.Hour
	EntrustedShopSpacingSq      = 15000
	EntrustedShopMaxBundleTotal = math.MaxInt16
	entrustedShopCloseTimer     = "entrusted_shop_close"
	entrustedShopPortalType     = 2
)

var (
	ErrMiniRoomInvalid      = errors.New("invalid mini room request")
	ErrMiniRoomBusy         = errors.New("mini room request is pending")
	ErrMiniRoomNotOwner     = errors.New("not the mini room owner")
	ErrMiniRoomItemNotFound = errors.New("mini room item not found")
	ErrMiniRoomFull         = errors.New("mini room has no more item slots")
	ErrMiniRoomOnlyHeld     = errors.New("one-of-a-kind item already held")
	ErrMiniRoomMesoOver     = errors.New("meso would overflow")
)

type MiniRoomEnterError struct {
	Code pconst.MiniRoomEnterError
}

func (e *MiniRoomEnterError) Error() string {
	return fmt.Sprintf("mini room enter failed: %d", e.Code)
}

type MiniRoomBuyError struct {
	Result pconst.MiniRoomBuyResult
}

func (e *MiniRoomBuyError) Error() string {
	return fmt.Sprintf("mini room buy failed: %d", e.Result)
}

type EntrustedShopItem struct {
	Item      Item
	Bundles   uint16
	PerBundle uint16
	Price     int32
}

func (hi *EntrustedShopItem) count() uint16 {
	return hi.Bundles * hi.PerBundle
}

type EntrustedShopSale struct {
	ItemID  uint32
	Bundles uint16
	Total   int32
	Buyer   string
}

type EntrustedShop struct {
	ObjectCore
	ID         uint32
	AccountID  uint32
	OwnerID    uint32
	OwnerName  string
	MapID      uint32
	ItemID     uint32
	Title      string
	Meso       int32
	Items      []*EntrustedShopItem
	Sold       []EntrustedShopSale
	OpenedAt   time.Time
	soldInform bool
	published  bool
	owner      *Character
	Visitors   [EntrustedShopVisitors]*Character
	saving     bool
	dirty      bool
	closing    bool
	entries    map[uint32]*internal.CharacterSaveEntry
}

func (es *EntrustedShop) GetObjectType() constant.ObjectType {
	return constant.ObjectTypeEntrustedShop
}

func (es *EntrustedShop) Is(typ constant.ObjectType) bool {
	return es.GetObjectType().Has(typ)
}

func (es *EntrustedShop) balloon() response.EntrustedShopBalloon {
	users := uint8(1)
	for _, visitor := range es.Visitors {
		if visitor != nil {
			users++
		}
	}
	return response.EntrustedShopBalloon{
		SN:     es.OID,
		Title:  es.Title,
		ItemID: es.ItemID,
		Users:  users,
	}
}

func (es *EntrustedShop) SendSpawnSyncToViewer(viewer *Character) {
	if es.published == false {
		return
	}
	footholdID := uint16(0)
	if foothold, ok := es.Map.Wz.Footholds.Find(types.Point[int16]{X: es.Position.X, Y: es.Position.Y}); ok {
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
	if es.published == false {
		return
	}
	viewer.Send(&response.DestroyEntrustedShop{EmployerID: es.OwnerID}, types.SEND_POLICY_ENCRYPT)
}

func (es *EntrustedShop) updateBalloon() {
	if es.published == false {
		return
	}
	es.Broadcast(&response.UpdateEntrustedShop{
		EmployerID:           es.OwnerID,
		EntrustedShopBalloon: es.balloon(),
	}, nil)
}

func (es *EntrustedShop) accepting() bool {
	return es.published && es.owner == nil
}

func (es *EntrustedShop) SlotOf(ch *Character) (uint8, bool) {
	if es.owner == ch {
		return 0, true
	}
	for i, visitor := range es.Visitors {
		if visitor == ch {
			return uint8(i + 1), true
		}
	}
	return 0, false
}

func (es *EntrustedShop) members() []*Character {
	members := make([]*Character, 0, EntrustedShopVisitors+1)
	if es.owner != nil {
		members = append(members, es.owner)
	}
	for _, visitor := range es.Visitors {
		if visitor != nil {
			members = append(members, visitor)
		}
	}
	return members
}

func (es *EntrustedShop) Elapsed() time.Duration {
	return clock.Now().Sub(es.OpenedAt)
}

func (es *EntrustedShop) tax(meso int32) int32 {
	var permille int64
	switch {
	case meso >= 100000000:
		permille = 30
	case meso >= 25000000:
		permille = 25
	case meso >= 10000000:
		permille = 20
	case meso >= 5000000:
		permille = 15
	case meso >= 1000000:
		permille = 9
	case meso >= 100000:
		permille = 4
	}
	return int32(int64(meso) * permille / 1000)
}

func (es *EntrustedShop) Visit(ch *Character) error {
	if ch.MiniRoom != nil {
		return ErrMiniRoomInvalid
	}

	if ch.GetID() == es.OwnerID {
		if es.published == false || es.owner != nil {
			return &MiniRoomEnterError{Code: pconst.MiniRoomEnterClosed}
		}
		for i, visitor := range es.Visitors {
			if visitor == nil {
				continue
			}
			es.Visitors[i] = nil
			visitor.MiniRoom = nil
			visitor.Listener.OnMiniRoomLeft(visitor, uint8(i+1), pconst.MiniRoomLeaveOrganizing)
		}
		es.owner = ch
		ch.MiniRoom = es
		es.updateBalloon()
		ch.Listener.OnMiniRoomEntered(ch, es, false)
		return nil
	}

	if es.published == false {
		return &MiniRoomEnterError{Code: pconst.MiniRoomEnterClosed}
	}
	if es.owner != nil {
		return &MiniRoomEnterError{Code: pconst.MiniRoomEnterOrganizing}
	}
	index := -1
	for i, visitor := range es.Visitors {
		if visitor == nil {
			index = i
			break
		}
	}
	if index < 0 {
		return &MiniRoomEnterError{Code: pconst.MiniRoomEnterFull}
	}

	slot := uint8(index + 1)
	for _, member := range es.members() {
		member.Listener.OnMiniRoomVisited(member, slot, ch)
	}
	es.Visitors[index] = ch
	ch.MiniRoom = es
	es.updateBalloon()
	ch.Listener.OnMiniRoomEntered(ch, es, false)
	return nil
}

func (es *EntrustedShop) Leave(ch *Character) {
	slot, ok := es.SlotOf(ch)
	if ok == false {
		return
	}

	ch.MiniRoom = nil
	if slot != 0 {
		es.Visitors[slot-1] = nil
		for _, member := range es.members() {
			member.Listener.OnMiniRoomLeft(member, slot, pconst.MiniRoomLeaveExit)
		}
		es.updateBalloon()
		return
	}

	es.owner = nil
	if es.published {
		es.updateBalloon()
		return
	}
	es.GameWorld.GetMapSystem().Call(es.Map, es.close)
}

func (es *EntrustedShop) Chat(ch *Character, message string) error {
	slot, ok := es.SlotOf(ch)
	if ok == false {
		return ErrMiniRoomInvalid
	}

	text := fmt.Sprintf("%s : %s", ch.GetName(), message)
	for _, member := range es.members() {
		member.Listener.OnMiniRoomChat(member, slot, text)
	}
	return nil
}

func (es *EntrustedShop) AddItem(actx actor.Context, ch *Character, invType constant.InventoryType, slot int16, bundles uint16, perBundle uint16, price int32) error {
	if es.owner != ch {
		return ErrMiniRoomNotOwner
	}
	if len(es.Items) >= EntrustedShopMaxItems {
		return ErrMiniRoomFull
	}
	if bundles == 0 || perBundle == 0 || price <= 0 {
		return ErrMiniRoomInvalid
	}

	item := ch.Inventory.GetItem(invType, slot)
	if item == nil {
		return ErrMiniRoomItemNotFound
	}
	model := item.GetModel()
	if model.IsTradeBlock() || model.IsAccountSharable() || model.IsQuest() {
		return ErrMiniRoomInvalid
	}
	if constant.ItemCategoryOf(model.GetID()) == constant.ItemCategoryPet {
		return ErrMiniRoomInvalid
	}
	if constant.IsRechargeable(model.GetID()) {
		bundles = 1
		perBundle = item.GetCount()
	}
	total := int(bundles) * int(perBundle)
	if total > EntrustedShopMaxBundleTotal || total > int(item.GetCount()) {
		return ErrMiniRoomInvalid
	}
	if int64(price)*int64(bundles) > math.MaxInt32 {
		return ErrMiniRoomInvalid
	}
	if model.IsOnly() {
		for _, listed := range es.Items {
			if listed.Item.GetModel().GetID() == model.GetID() {
				return ErrMiniRoomOnlyHeld
			}
		}
	}

	listed := &EntrustedShopItem{
		Item:      item.Clone(perBundle),
		Bundles:   bundles,
		PerBundle: perBundle,
		Price:     price,
	}
	ch.Inventory.RemoveItem(invType, slot, uint16(total))
	es.Items = append(es.Items, listed)
	ch.Listener.OnMiniRoomItems(ch, es)
	es.save(actx, false, ch)
	return nil
}

func (es *EntrustedShop) RemoveItem(actx actor.Context, ch *Character, index uint16) error {
	if es.owner != ch {
		return ErrMiniRoomNotOwner
	}
	if int(index) >= len(es.Items) {
		return ErrMiniRoomItemNotFound
	}

	listed := es.Items[index]
	if listed.Bundles > 0 {
		model := listed.Item.GetModel()
		if model.IsOnly() && ch.Inventory.HasItem(model.GetID()) {
			return ErrMiniRoomOnlyHeld
		}
		spec := ExchangeSpec{Reward: ExchangeSide{Items: map[uint32]uint16{model.GetID(): listed.count()}}}
		if spec.Valid(ch) != ExchangeOK {
			return ErrInventoryFull
		}
		ch.Inventory.addItemUnchecked(listed.Item.Clone(listed.count()), true)
	}

	es.Items = append(es.Items[:index:index], es.Items[index+1:]...)
	ch.Listener.OnMiniRoomItems(ch, es)
	es.save(actx, false, ch)
	return nil
}

func (es *EntrustedShop) Open(ch *Character) error {
	if es.owner != ch || es.published {
		return ErrMiniRoomNotOwner
	}
	if len(es.Items) == 0 {
		return ErrMiniRoomInvalid
	}

	es.owner = nil
	ch.MiniRoom = nil
	es.published = true
	es.BroadcastCall(func(obj Object) {
		viewer, ok := obj.(*Character)
		if ok == false {
			return
		}
		es.SendSpawnSyncToViewer(viewer)
	}, nil)
	es.scheduleClose()
	return nil
}

func (es *EntrustedShop) scheduleClose() {
	remaining := max(EntrustedShopDuration-es.Elapsed(), 0)
	es.AddTimer(entrustedShopCloseTimer, remaining, false, func() {
		es.GameWorld.GetMapSystem().Call(es.Map, func(ctx actor.Context) {
			if es.Map == nil {
				return
			}
			for _, member := range es.members() {
				slot, _ := es.SlotOf(member)
				member.MiniRoom = nil
				member.Listener.OnMiniRoomLeft(member, slot, pconst.MiniRoomLeaveTimeUp)
			}
			es.owner = nil
			es.Visitors = [EntrustedShopVisitors]*Character{}
			es.close(ctx)
		})
	})
}

func (es *EntrustedShop) EndMaintenance(ch *Character) error {
	if es.owner != ch || es.published == false {
		return ErrMiniRoomNotOwner
	}

	es.owner = nil
	ch.MiniRoom = nil
	es.updateBalloon()
	return nil
}

func (es *EntrustedShop) Buy(actx actor.Context, ch *Character, index uint16, bundles uint16) error {
	if _, ok := es.SlotOf(ch); ok == false || es.accepting() == false {
		return ErrMiniRoomInvalid
	}
	if int(index) >= len(es.Items) || bundles == 0 {
		return &MiniRoomBuyError{Result: pconst.MiniRoomBuyNotEnoughItem}
	}

	listed := es.Items[index]
	if bundles > listed.Bundles {
		return &MiniRoomBuyError{Result: pconst.MiniRoomBuyNotEnoughItem}
	}
	total := int64(listed.Price) * int64(bundles)
	if total > math.MaxInt32 {
		return &MiniRoomBuyError{Result: pconst.MiniRoomBuyUnknown}
	}
	if ch.Inventory.Meso < int32(total) {
		return &MiniRoomBuyError{Result: pconst.MiniRoomBuyNotEnoughMeso}
	}
	income := int32(total) - es.tax(int32(total))
	if income > math.MaxInt32-es.Meso {
		return &MiniRoomBuyError{Result: pconst.MiniRoomBuySellerLimit}
	}
	model := listed.Item.GetModel()
	if model.IsOnly() && ch.Inventory.HasItem(model.GetID()) {
		return &MiniRoomBuyError{Result: pconst.MiniRoomBuyOnlyOne}
	}
	count := bundles * listed.PerBundle
	spec := ExchangeSpec{
		Cost:   ExchangeSide{Meso: int32(total)},
		Reward: ExchangeSide{Items: map[uint32]uint16{model.GetID(): count}},
	}
	if spec.Valid(ch) != ExchangeOK {
		return &MiniRoomBuyError{Result: pconst.MiniRoomBuyInventoryFull}
	}

	ch.Inventory.removeMesoUnchecked(int32(total))
	ch.Inventory.addItemUnchecked(listed.Item.Clone(count), true)
	listed.Bundles -= bundles
	es.Meso += income
	es.Sold = append(es.Sold, EntrustedShopSale{
		ItemID:  model.GetID(),
		Bundles: bundles,
		Total:   int32(total),
		Buyer:   ch.GetName(),
	})
	if len(es.Sold) > math.MaxUint8 {
		es.Sold = es.Sold[len(es.Sold)-math.MaxUint8:]
	}
	for _, member := range es.members() {
		member.Listener.OnMiniRoomItems(member, es)
	}
	es.save(actx, false, ch)

	if es.soldInform == false {
		return nil
	}
	message := fmt.Sprintf("고용상점에서 %s %d개가 판매되었습니다.", es.GameWorld.GetResources().GetItemName(model.GetID()), count)
	es.GameWorld.GetDispatchSystem().CallCharacter(es.OwnerID, func(ctx actor.Context, owner *Character) {
		owner.Listener.OnMessage(owner, constant.MsgLightBlueText, message)
	})
	return nil
}

func (es *EntrustedShop) Arrange(actx actor.Context, ch *Character) error {
	if es.owner != ch {
		return ErrMiniRoomNotOwner
	}

	if (ExchangeSpec{Reward: ExchangeSide{Meso: es.Meso}}).Valid(ch) == ExchangeOK {
		ch.Inventory.addMesoUnchecked(es.Meso)
		es.Meso = 0
	}
	items := es.Items[:0]
	for _, listed := range es.Items {
		if listed.Bundles > 0 {
			items = append(items, listed)
		}
	}
	es.Items = items
	ch.Listener.OnMiniRoomArranged(ch, es)
	es.save(actx, false, ch)
	return nil
}

func (es *EntrustedShop) WithdrawMeso(actx actor.Context, ch *Character) error {
	if es.owner != ch {
		return ErrMiniRoomNotOwner
	}
	if es.Meso <= 0 {
		return ErrMiniRoomInvalid
	}
	if (ExchangeSpec{Reward: ExchangeSide{Meso: es.Meso}}).Valid(ch) != ExchangeOK {
		return ErrMiniRoomMesoOver
	}

	ch.Inventory.addMesoUnchecked(es.Meso)
	es.Meso = 0
	ch.Listener.OnMiniRoomMesoWithdrawn(ch)
	es.save(actx, false, ch)
	return nil
}

func (es *EntrustedShop) Close(actx actor.Context, ch *Character) error {
	if es.owner != ch {
		return ErrMiniRoomNotOwner
	}

	result := pconst.MiniRoomCloseAll
	if (ExchangeSpec{Reward: ExchangeSide{Meso: es.Meso}}).Valid(ch) != ExchangeOK {
		result = pconst.MiniRoomCloseMesoOver
	} else {
		ch.Inventory.addMesoUnchecked(es.Meso)
		es.Meso = 0
		result = es.returnItems(ch)
	}

	es.owner = nil
	ch.MiniRoom = nil
	ch.Listener.OnMiniRoomClosed(ch, result)
	es.entries[ch.GetID()] = ch.ToProto(es.GameWorld.GetWorldID())
	es.close(actx)
	return nil
}

func (es *EntrustedShop) returnItems(ch *Character) pconst.MiniRoomCloseResult {
	rewards := map[uint32]uint16{}
	for _, listed := range es.Items {
		if listed.Bundles == 0 {
			continue
		}
		model := listed.Item.GetModel()
		if model.IsOnly() && ch.Inventory.HasItem(model.GetID()) {
			return pconst.MiniRoomCloseOnlyOne
		}
		rewards[model.GetID()] += listed.count()
	}
	if (ExchangeSpec{Reward: ExchangeSide{Items: rewards}}).Valid(ch) != ExchangeOK {
		return pconst.MiniRoomCloseInventoryFull
	}

	for _, listed := range es.Items {
		if listed.Bundles == 0 {
			continue
		}
		ch.Inventory.addItemUnchecked(listed.Item.Clone(listed.count()), true)
	}
	es.Items = nil
	return pconst.MiniRoomCloseAll
}

func (es *EntrustedShop) close(actx actor.Context) {
	m := es.Map
	if m == nil {
		return
	}
	m.RemoveEntrustedShop(es)
	es.save(actx, true)
}

func (es *EntrustedShop) save(actx actor.Context, closing bool, characters ...*Character) {
	for _, ch := range characters {
		es.entries[ch.GetID()] = ch.ToProto(es.GameWorld.GetWorldID())
	}
	if closing {
		es.closing = true
	}
	if es.saving {
		es.dirty = true
		return
	}

	entries := make([]*internal.CharacterSaveEntry, 0, len(es.entries))
	for _, entry := range es.entries {
		entries = append(entries, entry)
	}
	es.entries = make(map[uint32]*internal.CharacterSaveEntry)
	es.saving = true
	es.dirty = false
	done := func() {
		es.saving = false
		if es.dirty {
			es.save(actx, es.closing)
		}
	}
	es.GameWorld.SaveEntrustedShopAsync(actx, es.ToProto(), entries, es.closing).Do(func(*internal.SaveEntrustedShopReply) error {
		done()
		return nil
	}).OnError(func(err error) {
		log.Printf("EntrustedShop.save shop=%d: %v", es.ID, err)
		done()
	})
}

func (ch *Character) UseEntrustedShop(actx actor.Context) {
	m := ch.GetMap()
	if m == nil || m.Wz.EntrustedShop == false || ch.MiniRoom != nil || ch.miniRoomPending {
		ch.Listener.OnEntrustedShopCheck(ch, pconst.EntrustedShopCannotOpen, 0, 0)
		return
	}

	ch.miniRoomPending = true
	ch.Listener.FindEntrustedShopAsync(actx, ch).Do(func(v *internal.FindEntrustedShopReply) error {
		ch.miniRoomPending = false
		shop := v.GetShop()
		switch {
		case shop == nil:
			ch.Listener.OnEntrustedShopCheck(ch, pconst.EntrustedShopTitle, 0, 0)
		case shop.GetClosedAtUnixMs() != 0:
			ch.Listener.OnEntrustedShopCheck(ch, pconst.EntrustedShopStoreBankFull, 0, 0)
		case shop.GetCharacterId() == ch.GetID():
			ch.Listener.OnEntrustedShopCheck(ch, pconst.EntrustedShopAlreadyOpen, shop.GetMapId(), uint8(shop.GetChannelId()))
		default:
			ch.Listener.OnEntrustedShopCheck(ch, pconst.EntrustedShopAccountBusy, 0, 0)
		}
		return nil
	}).OnError(func(err error) {
		ch.miniRoomPending = false
		ch.Listener.OnEntrustedShopCheck(ch, pconst.EntrustedShopCannotOpen, 0, 0)
		log.Printf("UseEntrustedShop character=%d: %v", ch.GetID(), err)
	})
}

func (ch *Character) CreateEntrustedShop(actx actor.Context, title string, slot int16, itemID uint32) error {
	m := ch.GetMap()
	if m == nil || m.Wz.EntrustedShop == false || ch.MiniRoom != nil {
		return &MiniRoomEnterError{Code: pconst.MiniRoomEnterCannotOpen}
	}
	if ch.miniRoomPending {
		return ErrMiniRoomBusy
	}
	permit := ch.Inventory.GetItem(constant.InventoryTypeCash, slot)
	if permit == nil || permit.GetModel().GetID() != itemID || itemID/10000 != 503 {
		return ErrMiniRoomInvalid
	}
	if expiration := permit.GetExpiration(); expiration.IsZero() == false && clock.Now().After(expiration) {
		return ErrMiniRoomInvalid
	}
	if title == "" {
		return ErrMiniRoomInvalid
	}
	if err := m.checkEntrustedShopSpot(ch.Position); err != nil {
		return err
	}

	soldInform := false
	if model, ok := permit.GetModel().(*wz.CashItem); ok {
		soldInform = model.SoldInform
	}
	es := &EntrustedShop{
		AccountID:  ch.AccountID,
		OwnerID:    ch.GetID(),
		OwnerName:  ch.GetName(),
		MapID:      m.TemplateID(),
		ItemID:     itemID,
		Title:      title,
		OpenedAt:   clock.Now(),
		soldInform: soldInform,
		entries:    make(map[uint32]*internal.CharacterSaveEntry),
	}
	es.ObjectCore.self = es
	es.Position = ch.Position

	ch.miniRoomPending = true
	ch.Listener.OpenEntrustedShopAsync(actx, ch, es.ToProto()).Do(func(v *internal.OpenEntrustedShopReply) error {
		ch.miniRoomPending = false
		if v.GetExisting() != nil {
			ch.Listener.OnMiniRoomEnterFailed(ch, pconst.MiniRoomEnterCannotOpen)
			return nil
		}

		es.ID = v.GetShopId()
		current := ch.GetMap()
		if current != m || ch.MiniRoom != nil || m.checkEntrustedShopSpot(es.Position) != nil {
			es.GameWorld = ch.GameWorld
			es.save(actx, true)
			ch.Listener.OnMiniRoomEnterFailed(ch, pconst.MiniRoomEnterCannotOpen)
			return nil
		}
		m.AddEntrustedShop(es)
		es.owner = ch
		ch.MiniRoom = es
		ch.Listener.OnMiniRoomEntered(ch, es, true)
		return nil
	}).OnError(func(err error) {
		ch.miniRoomPending = false
		ch.Listener.OnMiniRoomEnterFailed(ch, pconst.MiniRoomEnterCannotOpen)
		log.Printf("CreateEntrustedShop character=%d: %v", ch.GetID(), err)
	})
	return nil
}

func (ch *Character) VisitMiniRoom(sn uint32) error {
	m := ch.GetMap()
	if m == nil {
		return ErrMiniRoomInvalid
	}
	es, ok := m.GetObject(constant.ObjectTypeEntrustedShop, sn).(*EntrustedShop)
	if ok == false {
		return &MiniRoomEnterError{Code: pconst.MiniRoomEnterClosed}
	}
	return es.Visit(ch)
}
