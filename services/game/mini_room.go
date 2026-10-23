package game

import (
	"errors"
	"fmt"
	"log"
	"math"
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/maple/constant"
	pconst "github.com/boyism80/fm/protocol/constant"
	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
)

const (
	ShopVisitors       = 3
	ShopSpacingSq      = 15000
	ShopMaxBundleTotal = math.MaxInt16
	shopPortalType     = 2
	shopSaveRetryDelay = 5 * time.Second
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

type MiniRoom interface {
	Leave(ch *Character)
	Chat(ch *Character, message string) error
}

type ShopItem struct {
	item      Item
	bundles   uint16
	perBundle uint16
	price     int32
}

func (r *shopRoom) Title() string { return r.title }

func (r *shopRoom) Items() []*ShopItem { return r.items }

func (si *ShopItem) count() uint16 {
	return si.bundles * si.perBundle
}

type ShopSale struct {
	ItemID  uint32
	Bundles uint16
	Total   int32
	Buyer   string
}

// closedShop is goods waiting in the store bank until the shop save is accepted.
type closedShop struct {
	id    uint64
	meso  int32
	items []*ShopItem
}

type shopRoom struct {
	objectCore
	kind      uint8
	id        uint32
	accountID uint32
	ownerID   uint32
	ownerName string
	mapID     uint32
	itemID    uint32
	title     string
	meso      int32
	items     []*ShopItem
	sold      []ShopSale
	openedAt  time.Time
	home      *Map
	published bool
	owner     *Character
	visitors  [ShopVisitors]*Character
	saving    bool
	dirty     bool
	closing   bool
	entries   map[uint32]*internal.CharacterSaveEntry
	storeBank []closedShop
}

func (r *shopRoom) SlotOf(ch *Character) (uint8, bool) {
	if r.owner == ch {
		return 0, true
	}
	for i, visitor := range r.visitors {
		if visitor == ch {
			return uint8(i + 1), true
		}
	}
	return 0, false
}

func (r *shopRoom) Members() []*Character {
	members := make([]*Character, 0, ShopVisitors+1)
	if r.owner != nil {
		members = append(members, r.owner)
	}
	for _, visitor := range r.visitors {
		if visitor != nil {
			members = append(members, visitor)
		}
	}
	return members
}

func (r *shopRoom) users() uint8 {
	users := uint8(1)
	for _, visitor := range r.visitors {
		if visitor != nil {
			users++
		}
	}
	return users
}

func (r *shopRoom) freeSlot() (int, bool) {
	for i, visitor := range r.visitors {
		if visitor == nil {
			return i, true
		}
	}
	return 0, false
}

func (r *shopRoom) tax(meso int32) int32 {
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

func (r *shopRoom) Chat(ch *Character, message string) error {
	slot, ok := r.SlotOf(ch)
	if !ok {
		return ErrMiniRoomInvalid
	}

	text := fmt.Sprintf("%s : %s", ch.GetName(), message)
	for _, member := range r.Members() {
		member.listener.OnMiniRoomChat(member, slot, text)
	}
	return nil
}

func (r *shopRoom) list(ch *Character, invType constant.InventoryType, slot int16, bundles uint16, perBundle uint16, price int32) error {
	if item := ch.inventory.GetItem(invType, slot); item != nil && r.holds(item.GetModel().GetID()) {
		return ErrMiniRoomOnlyHeld
	}

	listed, err := ch.inventory.takeShopItem(invType, slot, bundles, perBundle, price)
	if err != nil {
		return err
	}
	r.items = append(r.items, listed)
	return nil
}

func (r *shopRoom) holds(itemID uint32) bool {
	for _, listed := range r.items {
		model := listed.item.GetModel()
		if model.GetTraits()&constant.ItemTraitUnique != 0 && model.GetID() == itemID {
			return true
		}
	}
	return false
}

func (inv *Inventory) takeShopItem(invType constant.InventoryType, slot int16, bundles uint16, perBundle uint16, price int32) (*ShopItem, error) {
	if bundles == 0 || perBundle == 0 || price <= 0 {
		return nil, ErrMiniRoomInvalid
	}

	item := inv.GetItem(invType, slot)
	if item == nil {
		return nil, ErrMiniRoomItemNotFound
	}
	if !item.Tradable() {
		return nil, ErrMiniRoomInvalid
	}
	if item.GetModel().GetTraits()&constant.ItemTraitRechargeable != 0 {
		bundles = 1
		perBundle = item.GetCount()
	}
	total := int(bundles) * int(perBundle)
	if total > ShopMaxBundleTotal || total > int(item.GetCount()) {
		return nil, ErrMiniRoomInvalid
	}
	if int64(price)*int64(bundles) > math.MaxInt32 {
		return nil, ErrMiniRoomInvalid
	}

	listed := &ShopItem{
		item:      item.Clone(perBundle),
		bundles:   bundles,
		perBundle: perBundle,
		price:     price,
	}
	inv.RemoveItem(invType, slot, uint16(total))
	return listed, nil
}

func (inv *Inventory) receive(item Item) {
	item.consumeKarma()
	inv.addItemUnchecked(item, true)
}

func (inv *Inventory) collectShopGoods(meso int32, items []*ShopItem) (int32, []*ShopItem) {
	if meso > 0 && (ExchangeSpec{Reward: ExchangeSide{Meso: meso}}).Valid(inv.owner, inv.owner.escrow) == ExchangeOK {
		inv.addMesoUnchecked(meso)
		meso = 0
	}

	var kept []*ShopItem
	for _, listed := range items {
		var spec ExchangeSpec
		spec.Reward.addItem(listed.item, listed.count())
		if inv.hasUniqueItem(listed.item.GetModel()) || spec.Valid(inv.owner, inv.owner.escrow) != ExchangeOK {
			kept = append(kept, listed)
			continue
		}
		inv.addItemUnchecked(listed.item.Clone(listed.count()), true)
	}
	return meso, kept
}

func (r *shopRoom) unlist(ch *Character, index uint16) error {
	if int(index) >= len(r.items) {
		return ErrMiniRoomItemNotFound
	}

	listed := r.items[index]
	if listed.bundles > 0 {
		model := listed.item.GetModel()
		if ch.inventory.hasUniqueItem(model) {
			return ErrMiniRoomOnlyHeld
		}
		var spec ExchangeSpec
		spec.Reward.addItem(listed.item, listed.count())
		if spec.Valid(ch, ch.escrow) != ExchangeOK {
			return ErrInventoryFull
		}
		ch.inventory.addItemUnchecked(listed.item.Clone(listed.count()), true)
	}

	r.items = append(r.items[:index:index], r.items[index+1:]...)
	return nil
}

func (r *shopRoom) sell(ch *Character, index uint16, bundles uint16, held int32) (*ShopItem, int32, int32, error) {
	if int(index) >= len(r.items) || bundles == 0 {
		return nil, 0, 0, &MiniRoomBuyError{Result: pconst.MiniRoomBuyNotEnoughItem}
	}

	listed := r.items[index]
	if bundles > listed.bundles {
		return nil, 0, 0, &MiniRoomBuyError{Result: pconst.MiniRoomBuyNotEnoughItem}
	}
	total := int64(listed.price) * int64(bundles)
	if total > math.MaxInt32 {
		return nil, 0, 0, &MiniRoomBuyError{Result: pconst.MiniRoomBuyUnknown}
	}
	if ch.inventory.meso < int32(total) {
		return nil, 0, 0, &MiniRoomBuyError{Result: pconst.MiniRoomBuyNotEnoughMeso}
	}
	income := int32(total) - r.tax(int32(total))
	if income > math.MaxInt32-held {
		return nil, 0, 0, &MiniRoomBuyError{Result: pconst.MiniRoomBuySellerLimit}
	}
	model := listed.item.GetModel()
	if ch.inventory.hasUniqueItem(model) {
		return nil, 0, 0, &MiniRoomBuyError{Result: pconst.MiniRoomBuyOnlyOne}
	}
	count := bundles * listed.perBundle
	spec := ExchangeSpec{Cost: ExchangeSide{Meso: int32(total)}}
	spec.Reward.addItem(listed.item, count)
	if spec.Valid(ch, ch.escrow) != ExchangeOK {
		return nil, 0, 0, &MiniRoomBuyError{Result: pconst.MiniRoomBuyInventoryFull}
	}

	ch.inventory.removeMesoUnchecked(int32(total))
	ch.inventory.receive(listed.item.Clone(count))
	listed.bundles -= bundles
	return listed, int32(total), income, nil
}

func (r *shopRoom) save(actx actor.Context, closing bool, characters ...*Character) {
	for _, ch := range characters {
		r.entries[ch.GetID()] = ch.BuildRPC(r.gameWorld.GetWorldID())
	}
	if closing {
		r.closing = true
	}
	if r.saving {
		r.dirty = true
		return
	}

	entries := make([]*internal.CharacterSaveEntry, 0, len(r.entries))
	for _, entry := range r.entries {
		entries = append(entries, entry)
	}
	r.entries = make(map[uint32]*internal.CharacterSaveEntry)
	storeBank := r.storeBank
	r.storeBank = nil
	banks := make([]*internal.Shop, 0, len(storeBank))
	for _, shop := range storeBank {
		banks = append(banks, shop.BuildRPC(r))
	}
	r.saving = true
	r.dirty = false
	done := func() {
		r.saving = false
		if r.dirty {
			r.save(actx, r.closing)
		}
	}
	r.gameWorld.Systems().Get[*InternalSystem]().SaveShopAsync(actx, r.BuildRPC(), entries, banks, r.closing).Do(func(*internal.SaveShopReply) error {
		done()
		return nil
	}).OnError(func(err error) {
		log.Printf("shopRoom.save shop=%d: %v", r.id, err)
		for _, entry := range entries {
			id := entry.GetCharacter().GetCharacterId()
			if _, ok := r.entries[id]; !ok {
				r.entries[id] = entry
			}
		}
		r.storeBank = append(storeBank, r.storeBank...)
		time.AfterFunc(shopSaveRetryDelay, func() {
			r.gameWorld.Systems().Get[*MapSystem]().Visit(r.home, func(ctx actor.Context) {
				r.saving = false
				r.save(ctx, r.closing)
			})
		})
	})
}

func (m *Map) VisitMiniRoom(ch *Character, sn uint32, password string) error {
	if t, ok := ch.miniRoom.(*Trade); ok && t.invited == ch && t.sn == sn {
		return t.Visit(ch)
	}
	// The sn of a remotely managed shop is an OID of another map and can match an object here.
	if ch.remoteShops.managed.sn != 0 && ch.remoteShops.managed.sn == sn {
		return ch.remoteShops.visitManaged(sn)
	}
	if es, ok := m.GetObject(constant.ObjectTypeEntrustedShop, sn).(*EntrustedShop); ok {
		if ch.GetID() != es.ownerID {
			return es.Visit(ch)
		}
		if ch.miniRoom != nil {
			return ErrMiniRoomInvalid
		}
		es.Manage(ch.GetID())
		return nil
	}
	if ps, ok := m.GetObject(constant.ObjectTypePersonalShop, sn).(*PersonalShop); ok {
		return ps.Visit(ch)
	}
	if g, ok := m.GetObject(constant.ObjectTypeMiniGame, sn).(*MiniGame); ok {
		return g.Visit(ch, password)
	}
	return &MiniRoomEnterError{Code: pconst.MiniRoomEnterClosed}
}
