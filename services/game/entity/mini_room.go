package entity

import (
	"errors"
	"fmt"
	"log"
	"math"
	"time"

	"github.com/asynkron/protoactor-go/actor"
	pconst "github.com/boyism80/fm/protocol/constant"
	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
	"github.com/boyism80/fm/services/game/constant"
)

const (
	ShopVisitors       = 3
	ShopSpacingSq      = 15000
	ShopMaxBundleTotal = math.MaxInt16
	shopPortalType     = 2
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
	Open(actx actor.Context, ch *Character) error
	AddItem(actx actor.Context, ch *Character, invType constant.InventoryType, slot int16, bundles uint16, perBundle uint16, price int32) error
	Buy(actx actor.Context, ch *Character, index uint16, bundles uint16) error
	RemoveItem(actx actor.Context, ch *Character, index uint16) error
}

type ShopItem struct {
	Item      Item
	Bundles   uint16
	PerBundle uint16
	Price     int32
}

func (si *ShopItem) count() uint16 {
	return si.Bundles * si.PerBundle
}

type ShopSale struct {
	ItemID  uint32
	Bundles uint16
	Total   int32
	Buyer   string
}

type remoteShop struct {
	sn    uint32
	mapID uint32
}

type shopRoom struct {
	ObjectCore
	kind      uint8
	ID        uint32
	AccountID uint32
	OwnerID   uint32
	OwnerName string
	MapID     uint32
	ItemID    uint32
	Title     string
	Meso      int32
	Items     []*ShopItem
	Sold      []ShopSale
	OpenedAt  time.Time
	published bool
	owner     *Character
	Visitors  [ShopVisitors]*Character
	saving    bool
	dirty     bool
	closing   bool
	entries   map[uint32]*internal.CharacterSaveEntry
	storeBank []*internal.Shop
}

func (r *shopRoom) SlotOf(ch *Character) (uint8, bool) {
	if r.owner == ch {
		return 0, true
	}
	for i, visitor := range r.Visitors {
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
	for _, visitor := range r.Visitors {
		if visitor != nil {
			members = append(members, visitor)
		}
	}
	return members
}

func (r *shopRoom) users() uint8 {
	users := uint8(1)
	for _, visitor := range r.Visitors {
		if visitor != nil {
			users++
		}
	}
	return users
}

func (r *shopRoom) freeSlot() (int, bool) {
	for i, visitor := range r.Visitors {
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
	if ok == false {
		return ErrMiniRoomInvalid
	}

	text := fmt.Sprintf("%s : %s", ch.GetName(), message)
	for _, member := range r.Members() {
		member.Listener.OnMiniRoomChat(member, slot, text)
	}
	return nil
}

func (r *shopRoom) list(ch *Character, invType constant.InventoryType, slot int16, bundles uint16, perBundle uint16, price int32) error {
	if item := ch.Inventory.GetItem(invType, slot); item != nil && r.holds(item.GetModel().GetID()) {
		return ErrMiniRoomOnlyHeld
	}

	listed, err := ch.takeShopItem(invType, slot, bundles, perBundle, price)
	if err != nil {
		return err
	}
	r.Items = append(r.Items, listed)
	return nil
}

func (r *shopRoom) holds(itemID uint32) bool {
	for _, listed := range r.Items {
		model := listed.Item.GetModel()
		if model.IsOnly() && model.GetID() == itemID {
			return true
		}
	}
	return false
}

func (ch *Character) takeShopItem(invType constant.InventoryType, slot int16, bundles uint16, perBundle uint16, price int32) (*ShopItem, error) {
	if bundles == 0 || perBundle == 0 || price <= 0 {
		return nil, ErrMiniRoomInvalid
	}

	item := ch.Inventory.GetItem(invType, slot)
	if item == nil {
		return nil, ErrMiniRoomItemNotFound
	}
	model := item.GetModel()
	if model.IsTradeBlock() || model.IsAccountSharable() || model.IsQuest() {
		return nil, ErrMiniRoomInvalid
	}
	if constant.ItemCategoryOf(model.GetID()) == constant.ItemCategoryPet {
		return nil, ErrMiniRoomInvalid
	}
	if constant.IsRechargeable(model.GetID()) {
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
		Item:      item.Clone(perBundle),
		Bundles:   bundles,
		PerBundle: perBundle,
		Price:     price,
	}
	ch.Inventory.RemoveItem(invType, slot, uint16(total))
	return listed, nil
}

func (ch *Character) collectShopGoods(meso int32, items []*ShopItem) (int32, []*ShopItem) {
	if meso > 0 && (ExchangeSpec{Reward: ExchangeSide{Meso: meso}}).Valid(ch) == ExchangeOK {
		ch.Inventory.addMesoUnchecked(meso)
		meso = 0
	}

	var kept []*ShopItem
	for _, listed := range items {
		model := listed.Item.GetModel()
		spec := ExchangeSpec{Reward: ExchangeSide{Items: map[uint32]uint16{model.GetID(): listed.count()}}}
		if (model.IsOnly() && ch.Inventory.HasItem(model.GetID())) || spec.Valid(ch) != ExchangeOK {
			kept = append(kept, listed)
			continue
		}
		ch.Inventory.addItemUnchecked(listed.Item.Clone(listed.count()), true)
	}
	return meso, kept
}

func (r *shopRoom) unlist(ch *Character, index uint16) error {
	if int(index) >= len(r.Items) {
		return ErrMiniRoomItemNotFound
	}

	listed := r.Items[index]
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

	r.Items = append(r.Items[:index:index], r.Items[index+1:]...)
	return nil
}

func (r *shopRoom) sell(ch *Character, index uint16, bundles uint16, held int32) (*ShopItem, int32, int32, error) {
	if int(index) >= len(r.Items) || bundles == 0 {
		return nil, 0, 0, &MiniRoomBuyError{Result: pconst.MiniRoomBuyNotEnoughItem}
	}

	listed := r.Items[index]
	if bundles > listed.Bundles {
		return nil, 0, 0, &MiniRoomBuyError{Result: pconst.MiniRoomBuyNotEnoughItem}
	}
	total := int64(listed.Price) * int64(bundles)
	if total > math.MaxInt32 {
		return nil, 0, 0, &MiniRoomBuyError{Result: pconst.MiniRoomBuyUnknown}
	}
	if ch.Inventory.Meso < int32(total) {
		return nil, 0, 0, &MiniRoomBuyError{Result: pconst.MiniRoomBuyNotEnoughMeso}
	}
	income := int32(total) - r.tax(int32(total))
	if income > math.MaxInt32-held {
		return nil, 0, 0, &MiniRoomBuyError{Result: pconst.MiniRoomBuySellerLimit}
	}
	model := listed.Item.GetModel()
	if model.IsOnly() && ch.Inventory.HasItem(model.GetID()) {
		return nil, 0, 0, &MiniRoomBuyError{Result: pconst.MiniRoomBuyOnlyOne}
	}
	count := bundles * listed.PerBundle
	spec := ExchangeSpec{
		Cost:   ExchangeSide{Meso: int32(total)},
		Reward: ExchangeSide{Items: map[uint32]uint16{model.GetID(): count}},
	}
	if spec.Valid(ch) != ExchangeOK {
		return nil, 0, 0, &MiniRoomBuyError{Result: pconst.MiniRoomBuyInventoryFull}
	}

	ch.Inventory.removeMesoUnchecked(int32(total))
	ch.Inventory.addItemUnchecked(listed.Item.Clone(count), true)
	listed.Bundles -= bundles
	return listed, int32(total), income, nil
}

func (r *shopRoom) save(actx actor.Context, closing bool, characters ...*Character) {
	for _, ch := range characters {
		r.entries[ch.GetID()] = ch.ToProto(r.GameWorld.GetWorldID())
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
	r.saving = true
	r.dirty = false
	done := func() {
		r.saving = false
		if r.dirty {
			r.save(actx, r.closing)
		}
	}
	r.GameWorld.SaveShopAsync(actx, r.ToProto(), entries, storeBank, r.closing).Do(func(*internal.SaveShopReply) error {
		done()
		return nil
	}).OnError(func(err error) {
		log.Printf("shopRoom.save shop=%d: %v", r.ID, err)
		done()
	})
}
