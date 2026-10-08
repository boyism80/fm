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
	HiredMerchantMaxItems       = 10
	HiredMerchantVisitors       = 3
	HiredMerchantDuration       = 24 * time.Hour
	HiredMerchantSpacingSq      = 15000
	HiredMerchantMaxBundleTotal = math.MaxInt16
	hiredMerchantCloseTimer     = "hired_merchant_close"
	hiredMerchantPortalType     = 2
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

type HiredMerchantItem struct {
	Item      Item
	Bundles   uint16
	PerBundle uint16
	Price     int32
}

func (hi *HiredMerchantItem) count() uint16 {
	return hi.Bundles * hi.PerBundle
}

type HiredMerchantSale struct {
	ItemID  uint32
	Bundles uint16
	Total   int32
	Buyer   string
}

type HiredMerchant struct {
	ObjectCore
	ID         uint32
	AccountID  uint32
	OwnerID    uint32
	OwnerName  string
	MapID      uint32
	ItemID     uint32
	Title      string
	Meso       int32
	Items      []*HiredMerchantItem
	Sold       []HiredMerchantSale
	OpenedAt   time.Time
	soldInform bool
	published  bool
	owner      *Character
	Visitors   [HiredMerchantVisitors]*Character
	saving     bool
	dirty      bool
	closing    bool
	entries    map[uint32]*internal.CharacterSaveEntry
}

func (hm *HiredMerchant) GetObjectType() constant.ObjectType {
	return constant.ObjectTypeHiredMerchant
}

func (hm *HiredMerchant) Is(typ constant.ObjectType) bool {
	return hm.GetObjectType().Has(typ)
}

func (hm *HiredMerchant) balloon() response.HiredMerchantBalloon {
	users := uint8(1)
	for _, visitor := range hm.Visitors {
		if visitor != nil {
			users++
		}
	}
	return response.HiredMerchantBalloon{
		SN:     hm.OID,
		Title:  hm.Title,
		ItemID: hm.ItemID,
		Users:  users,
	}
}

func (hm *HiredMerchant) SendSpawnSyncToViewer(viewer *Character) {
	if hm.published == false {
		return
	}
	footholdID := uint16(0)
	if foothold, ok := hm.Map.Wz.Footholds.Find(types.Point[int16]{X: hm.Position.X, Y: hm.Position.Y}); ok {
		footholdID = uint16(foothold.ID)
	}
	viewer.Send(&response.SpawnHiredMerchant{
		EmployerID:           hm.OwnerID,
		X:                    hm.Position.X,
		Y:                    hm.Position.Y,
		Foothold:             footholdID,
		OwnerName:            hm.OwnerName,
		HiredMerchantBalloon: hm.balloon(),
	}, types.SEND_POLICY_ENCRYPT)
}

func (hm *HiredMerchant) SendDestroySyncToViewer(viewer *Character) {
	if hm.published == false {
		return
	}
	viewer.Send(&response.DestroyHiredMerchant{EmployerID: hm.OwnerID}, types.SEND_POLICY_ENCRYPT)
}

func (hm *HiredMerchant) updateBalloon() {
	if hm.published == false {
		return
	}
	hm.Broadcast(&response.UpdateHiredMerchant{
		EmployerID:           hm.OwnerID,
		HiredMerchantBalloon: hm.balloon(),
	}, nil)
}

func (hm *HiredMerchant) accepting() bool {
	return hm.published && hm.owner == nil
}

func (hm *HiredMerchant) SlotOf(ch *Character) (uint8, bool) {
	if hm.owner == ch {
		return 0, true
	}
	for i, visitor := range hm.Visitors {
		if visitor == ch {
			return uint8(i + 1), true
		}
	}
	return 0, false
}

func (hm *HiredMerchant) members() []*Character {
	members := make([]*Character, 0, HiredMerchantVisitors+1)
	if hm.owner != nil {
		members = append(members, hm.owner)
	}
	for _, visitor := range hm.Visitors {
		if visitor != nil {
			members = append(members, visitor)
		}
	}
	return members
}

func (hm *HiredMerchant) Elapsed() time.Duration {
	return clock.Now().Sub(hm.OpenedAt)
}

func (hm *HiredMerchant) tax(meso int32) int32 {
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

func (hm *HiredMerchant) Visit(ch *Character) error {
	if ch.MiniRoom != nil {
		return ErrMiniRoomInvalid
	}

	if ch.GetID() == hm.OwnerID {
		if hm.published == false || hm.owner != nil {
			return &MiniRoomEnterError{Code: pconst.MiniRoomEnterClosed}
		}
		for i, visitor := range hm.Visitors {
			if visitor == nil {
				continue
			}
			hm.Visitors[i] = nil
			visitor.MiniRoom = nil
			visitor.Listener.OnMiniRoomLeft(visitor, uint8(i+1), pconst.MiniRoomLeaveOrganizing)
		}
		hm.owner = ch
		ch.MiniRoom = hm
		hm.updateBalloon()
		ch.Listener.OnMiniRoomEntered(ch, hm, false)
		return nil
	}

	if hm.published == false {
		return &MiniRoomEnterError{Code: pconst.MiniRoomEnterClosed}
	}
	if hm.owner != nil {
		return &MiniRoomEnterError{Code: pconst.MiniRoomEnterOrganizing}
	}
	index := -1
	for i, visitor := range hm.Visitors {
		if visitor == nil {
			index = i
			break
		}
	}
	if index < 0 {
		return &MiniRoomEnterError{Code: pconst.MiniRoomEnterFull}
	}

	slot := uint8(index + 1)
	for _, member := range hm.members() {
		member.Listener.OnMiniRoomVisited(member, slot, ch)
	}
	hm.Visitors[index] = ch
	ch.MiniRoom = hm
	hm.updateBalloon()
	ch.Listener.OnMiniRoomEntered(ch, hm, false)
	return nil
}

func (hm *HiredMerchant) Leave(ch *Character) {
	slot, ok := hm.SlotOf(ch)
	if ok == false {
		return
	}

	ch.MiniRoom = nil
	if slot != 0 {
		hm.Visitors[slot-1] = nil
		for _, member := range hm.members() {
			member.Listener.OnMiniRoomLeft(member, slot, pconst.MiniRoomLeaveExit)
		}
		hm.updateBalloon()
		return
	}

	hm.owner = nil
	if hm.published {
		hm.updateBalloon()
		return
	}
	hm.GameWorld.GetMapSystem().Call(hm.Map, hm.close)
}

func (hm *HiredMerchant) Chat(ch *Character, message string) error {
	slot, ok := hm.SlotOf(ch)
	if ok == false {
		return ErrMiniRoomInvalid
	}

	text := fmt.Sprintf("%s : %s", ch.GetName(), message)
	for _, member := range hm.members() {
		member.Listener.OnMiniRoomChat(member, slot, text)
	}
	return nil
}

func (hm *HiredMerchant) AddItem(actx actor.Context, ch *Character, invType constant.InventoryType, slot int16, bundles uint16, perBundle uint16, price int32) error {
	if hm.owner != ch {
		return ErrMiniRoomNotOwner
	}
	if len(hm.Items) >= HiredMerchantMaxItems {
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
	if total > HiredMerchantMaxBundleTotal || total > int(item.GetCount()) {
		return ErrMiniRoomInvalid
	}
	if int64(price)*int64(bundles) > math.MaxInt32 {
		return ErrMiniRoomInvalid
	}
	if model.IsOnly() {
		for _, listed := range hm.Items {
			if listed.Item.GetModel().GetID() == model.GetID() {
				return ErrMiniRoomOnlyHeld
			}
		}
	}

	listed := &HiredMerchantItem{
		Item:      item.Clone(perBundle),
		Bundles:   bundles,
		PerBundle: perBundle,
		Price:     price,
	}
	ch.Inventory.RemoveItem(invType, slot, uint16(total))
	hm.Items = append(hm.Items, listed)
	ch.Listener.OnMiniRoomItems(ch, hm)
	hm.save(actx, false, ch)
	return nil
}

func (hm *HiredMerchant) RemoveItem(actx actor.Context, ch *Character, index uint16) error {
	if hm.owner != ch {
		return ErrMiniRoomNotOwner
	}
	if int(index) >= len(hm.Items) {
		return ErrMiniRoomItemNotFound
	}

	listed := hm.Items[index]
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

	hm.Items = append(hm.Items[:index:index], hm.Items[index+1:]...)
	ch.Listener.OnMiniRoomItems(ch, hm)
	hm.save(actx, false, ch)
	return nil
}

func (hm *HiredMerchant) Open(ch *Character) error {
	if hm.owner != ch || hm.published {
		return ErrMiniRoomNotOwner
	}
	if len(hm.Items) == 0 {
		return ErrMiniRoomInvalid
	}

	hm.owner = nil
	ch.MiniRoom = nil
	hm.published = true
	hm.BroadcastCall(func(obj Object) {
		viewer, ok := obj.(*Character)
		if ok == false {
			return
		}
		hm.SendSpawnSyncToViewer(viewer)
	}, nil)
	hm.scheduleClose()
	return nil
}

func (hm *HiredMerchant) scheduleClose() {
	remaining := max(HiredMerchantDuration-hm.Elapsed(), 0)
	hm.AddTimer(hiredMerchantCloseTimer, remaining, false, func() {
		hm.GameWorld.GetMapSystem().Call(hm.Map, func(ctx actor.Context) {
			if hm.Map == nil {
				return
			}
			for _, member := range hm.members() {
				slot, _ := hm.SlotOf(member)
				member.MiniRoom = nil
				member.Listener.OnMiniRoomLeft(member, slot, pconst.MiniRoomLeaveTimeUp)
			}
			hm.owner = nil
			hm.Visitors = [HiredMerchantVisitors]*Character{}
			hm.close(ctx)
		})
	})
}

func (hm *HiredMerchant) EndMaintenance(ch *Character) error {
	if hm.owner != ch || hm.published == false {
		return ErrMiniRoomNotOwner
	}

	hm.owner = nil
	ch.MiniRoom = nil
	hm.updateBalloon()
	return nil
}

func (hm *HiredMerchant) Buy(actx actor.Context, ch *Character, index uint16, bundles uint16) error {
	if _, ok := hm.SlotOf(ch); ok == false || hm.accepting() == false {
		return ErrMiniRoomInvalid
	}
	if int(index) >= len(hm.Items) || bundles == 0 {
		return &MiniRoomBuyError{Result: pconst.MiniRoomBuyNotEnoughItem}
	}

	listed := hm.Items[index]
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
	income := int32(total) - hm.tax(int32(total))
	if income > math.MaxInt32-hm.Meso {
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
	hm.Meso += income
	hm.Sold = append(hm.Sold, HiredMerchantSale{
		ItemID:  model.GetID(),
		Bundles: bundles,
		Total:   int32(total),
		Buyer:   ch.GetName(),
	})
	if len(hm.Sold) > math.MaxUint8 {
		hm.Sold = hm.Sold[len(hm.Sold)-math.MaxUint8:]
	}
	for _, member := range hm.members() {
		member.Listener.OnMiniRoomItems(member, hm)
	}
	hm.save(actx, false, ch)

	if hm.soldInform == false {
		return nil
	}
	message := fmt.Sprintf("고용상점에서 %s %d개가 판매되었습니다.", hm.GameWorld.GetResources().GetItemName(model.GetID()), count)
	hm.GameWorld.GetDispatchSystem().CallCharacter(hm.OwnerID, func(ctx actor.Context, owner *Character) {
		owner.Listener.OnMessage(owner, constant.MsgLightBlueText, message)
	})
	return nil
}

func (hm *HiredMerchant) Arrange(actx actor.Context, ch *Character) error {
	if hm.owner != ch {
		return ErrMiniRoomNotOwner
	}

	if (ExchangeSpec{Reward: ExchangeSide{Meso: hm.Meso}}).Valid(ch) == ExchangeOK {
		ch.Inventory.addMesoUnchecked(hm.Meso)
		hm.Meso = 0
	}
	items := hm.Items[:0]
	for _, listed := range hm.Items {
		if listed.Bundles > 0 {
			items = append(items, listed)
		}
	}
	hm.Items = items
	ch.Listener.OnMiniRoomArranged(ch, hm)
	hm.save(actx, false, ch)
	return nil
}

func (hm *HiredMerchant) WithdrawMeso(actx actor.Context, ch *Character) error {
	if hm.owner != ch {
		return ErrMiniRoomNotOwner
	}
	if hm.Meso <= 0 {
		return ErrMiniRoomInvalid
	}
	if (ExchangeSpec{Reward: ExchangeSide{Meso: hm.Meso}}).Valid(ch) != ExchangeOK {
		return ErrMiniRoomMesoOver
	}

	ch.Inventory.addMesoUnchecked(hm.Meso)
	hm.Meso = 0
	ch.Listener.OnMiniRoomMesoWithdrawn(ch)
	hm.save(actx, false, ch)
	return nil
}

func (hm *HiredMerchant) Close(actx actor.Context, ch *Character) error {
	if hm.owner != ch {
		return ErrMiniRoomNotOwner
	}

	result := pconst.MiniRoomCloseAll
	if (ExchangeSpec{Reward: ExchangeSide{Meso: hm.Meso}}).Valid(ch) != ExchangeOK {
		result = pconst.MiniRoomCloseMesoOver
	} else {
		ch.Inventory.addMesoUnchecked(hm.Meso)
		hm.Meso = 0
		result = hm.returnItems(ch)
	}

	hm.owner = nil
	ch.MiniRoom = nil
	ch.Listener.OnMiniRoomClosed(ch, result)
	hm.entries[ch.GetID()] = ch.ToProto(hm.GameWorld.GetWorldID())
	hm.close(actx)
	return nil
}

func (hm *HiredMerchant) returnItems(ch *Character) pconst.MiniRoomCloseResult {
	rewards := map[uint32]uint16{}
	for _, listed := range hm.Items {
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

	for _, listed := range hm.Items {
		if listed.Bundles == 0 {
			continue
		}
		ch.Inventory.addItemUnchecked(listed.Item.Clone(listed.count()), true)
	}
	hm.Items = nil
	return pconst.MiniRoomCloseAll
}

func (hm *HiredMerchant) close(actx actor.Context) {
	m := hm.Map
	if m == nil {
		return
	}
	m.RemoveHiredMerchant(hm)
	hm.save(actx, true)
}

func (hm *HiredMerchant) save(actx actor.Context, closing bool, characters ...*Character) {
	for _, ch := range characters {
		hm.entries[ch.GetID()] = ch.ToProto(hm.GameWorld.GetWorldID())
	}
	if closing {
		hm.closing = true
	}
	if hm.saving {
		hm.dirty = true
		return
	}

	entries := make([]*internal.CharacterSaveEntry, 0, len(hm.entries))
	for _, entry := range hm.entries {
		entries = append(entries, entry)
	}
	hm.entries = make(map[uint32]*internal.CharacterSaveEntry)
	hm.saving = true
	hm.dirty = false
	done := func() {
		hm.saving = false
		if hm.dirty {
			hm.save(actx, hm.closing)
		}
	}
	hm.GameWorld.SaveHiredMerchantAsync(actx, hm.ToProto(), entries, hm.closing).Do(func(*internal.SaveHiredMerchantReply) error {
		done()
		return nil
	}).OnError(func(err error) {
		log.Printf("HiredMerchant.save merchant=%d: %v", hm.ID, err)
		done()
	})
}

func (ch *Character) UseHiredMerchant(actx actor.Context) {
	m := ch.GetMap()
	if m == nil || m.Wz.EntrustedShop == false || ch.MiniRoom != nil || ch.miniRoomPending {
		ch.Listener.OnEntrustedShopCheck(ch, pconst.EntrustedShopCannotOpen, 0, 0)
		return
	}

	ch.miniRoomPending = true
	ch.Listener.FindHiredMerchantAsync(actx, ch).Do(func(v *internal.FindHiredMerchantReply) error {
		ch.miniRoomPending = false
		merchant := v.GetMerchant()
		switch {
		case merchant == nil:
			ch.Listener.OnEntrustedShopCheck(ch, pconst.EntrustedShopTitle, 0, 0)
		case merchant.GetClosedAtUnixMs() != 0:
			ch.Listener.OnEntrustedShopCheck(ch, pconst.EntrustedShopStoreBankFull, 0, 0)
		case merchant.GetCharacterId() == ch.GetID():
			ch.Listener.OnEntrustedShopCheck(ch, pconst.EntrustedShopAlreadyOpen, merchant.GetMapId(), uint8(merchant.GetChannelId()))
		default:
			ch.Listener.OnEntrustedShopCheck(ch, pconst.EntrustedShopAccountBusy, 0, 0)
		}
		return nil
	}).OnError(func(err error) {
		ch.miniRoomPending = false
		ch.Listener.OnEntrustedShopCheck(ch, pconst.EntrustedShopCannotOpen, 0, 0)
		log.Printf("UseHiredMerchant character=%d: %v", ch.GetID(), err)
	})
}

func (ch *Character) CreateHiredMerchant(actx actor.Context, title string, slot int16, itemID uint32) error {
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
	if err := m.checkHiredMerchantSpot(ch.Position); err != nil {
		return err
	}

	soldInform := false
	if model, ok := permit.GetModel().(*wz.CashItem); ok {
		soldInform = model.SoldInform
	}
	hm := &HiredMerchant{
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
	hm.ObjectCore.self = hm
	hm.Position = ch.Position

	ch.miniRoomPending = true
	ch.Listener.OpenHiredMerchantAsync(actx, ch, hm.ToProto()).Do(func(v *internal.OpenHiredMerchantReply) error {
		ch.miniRoomPending = false
		if v.GetExisting() != nil {
			ch.Listener.OnMiniRoomEnterFailed(ch, pconst.MiniRoomEnterCannotOpen)
			return nil
		}

		hm.ID = v.GetMerchantId()
		current := ch.GetMap()
		if current != m || ch.MiniRoom != nil || m.checkHiredMerchantSpot(hm.Position) != nil {
			hm.GameWorld = ch.GameWorld
			hm.save(actx, true)
			ch.Listener.OnMiniRoomEnterFailed(ch, pconst.MiniRoomEnterCannotOpen)
			return nil
		}
		m.AddHiredMerchant(hm)
		hm.owner = ch
		ch.MiniRoom = hm
		ch.Listener.OnMiniRoomEntered(ch, hm, true)
		return nil
	}).OnError(func(err error) {
		ch.miniRoomPending = false
		ch.Listener.OnMiniRoomEnterFailed(ch, pconst.MiniRoomEnterCannotOpen)
		log.Printf("CreateHiredMerchant character=%d: %v", ch.GetID(), err)
	})
	return nil
}

func (ch *Character) VisitMiniRoom(sn uint32) error {
	m := ch.GetMap()
	if m == nil {
		return ErrMiniRoomInvalid
	}
	hm, ok := m.GetObject(constant.ObjectTypeHiredMerchant, sn).(*HiredMerchant)
	if ok == false {
		return &MiniRoomEnterError{Code: pconst.MiniRoomEnterClosed}
	}
	return hm.Visit(ch)
}
