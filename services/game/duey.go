package game

import (
	"encoding/json"
	"errors"
	"log"
	"math"
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/maple/constant"
	pconst "github.com/boyism80/fm/protocol/constant"
	"github.com/boyism80/fm/protocol/dto"
	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
)

const (
	ParcelArrivalDelay        = 12 * time.Hour
	ParcelKeepDuration        = 30 * 24 * time.Hour
	ParcelMessageBytes        = 100
	QuickParcelCoupon  uint32 = 5330000
	systemParcelSender        = "택배"
)

var (
	ErrDueyClosed         = errors.New("duey is not open")
	ErrDueyBusy           = errors.New("duey request is pending")
	ErrDueyInvalid        = errors.New("invalid duey request")
	ErrDueyNotEnoughMeso  = errors.New("not enough meso for duey")
	ErrDueyNotSendable    = errors.New("item cannot be sent")
	ErrDueyParcelNotFound = errors.New("parcel not found")
	ErrDueyOnlyHeld       = errors.New("one-of-a-kind item already held")
	ErrDueyCannotReceive  = errors.New("parcel cannot be received")
)

type DueyWindow uint8

const (
	DueyWindowClosed DueyWindow = iota
	DueyWindowIdentity
	DueyWindowBox
	DueyWindowQuick
)

type Parcel struct {
	id      uint32
	sender  string
	meso    int32
	quick   bool
	message string
	sentAt  time.Time
	item    Item
}

func (p *Parcel) ArrivedAt() time.Time {
	if p.quick {
		return p.sentAt
	}
	return p.sentAt.Add(ParcelArrivalDelay)
}

func (p *Parcel) ExpiresAt() time.Time {
	return p.ArrivedAt().Add(ParcelKeepDuration)
}

func (p *Parcel) BuildS2C() dto.Parcel {
	parcel := dto.Parcel{
		ID:      p.id,
		Sender:  p.sender,
		Meso:    p.meso,
		Expire:  p.ExpiresAt(),
		Quick:   p.quick,
		Message: p.message,
	}
	if p.item != nil {
		parcel.Item = p.item.BuildS2C()
	}
	return parcel
}

type parcelSendDetail struct {
	Recipient  string          `json:"recipient"`
	ReceiverID uint32          `json:"receiver_id,omitempty"`
	Meso       int32           `json:"meso"`
	Quick      bool            `json:"quick"`
	Item       json.RawMessage `json:"item,omitempty"`
	System     bool            `json:"system,omitempty"`
}

type Duey struct {
	owner       *Character
	parcels     []*Parcel
	expired     []*Parcel
	window      DueyWindow
	fromArrival bool
	pending     bool
}

func (d *Duey) Open(actx actor.Context, fromArrival bool) {
	if d.pending {
		return
	}

	d.pending = true
	d.owner.gameWorld.Systems().Get[*InternalSystem]().LoadParcelsAsync(actx, d.owner).Bind(d.owner).Do(func(v *internal.LoadParcelsReply) error {
		d.pending = false
		reply := v
		d.parcels = d.parcelsOf(reply.GetParcels())
		d.expired = d.parcelsOf(reply.GetExpired())
		d.fromArrival = fromArrival
		if d.owner.gameWorld.DueyIdentityPrompt() {
			d.window = DueyWindowIdentity
			d.owner.listener.OnDueyResult(d.owner, pconst.DueyResultIdentity)
			return nil
		}
		d.show()
		return nil
	}).OnError(func(err error) {
		d.pending = false

		log.Printf("Duey.Open character=%d: %v", d.owner.GetID(), err)
	})
}

func (d *Duey) parcelsOf(pbs []*internal.Parcel) []*Parcel {
	parcels := make([]*Parcel, 0, len(pbs))
	for _, pb := range pbs {
		parcels = append(parcels, NewParcelFromInternalProto(pb, d.owner.gameWorld))
	}
	return parcels
}

func (d *Duey) show() {
	d.window = DueyWindowBox
	d.owner.listener.OnOpenDuey(d.owner, d.fromArrival)
	d.expired = nil
}

func (d *Duey) ConfirmIdentity() error {
	if d.window != DueyWindowIdentity {
		return ErrDueyClosed
	}
	d.show()
	return nil
}

func (d *Duey) OpenQuick() error {
	if d.window != DueyWindowClosed {
		return ErrDueyInvalid
	}
	d.window = DueyWindowQuick
	d.owner.listener.OnDueyResult(d.owner, pconst.DueyResultQuickOpen)
	return nil
}

func (d *Duey) Close() {
	d.window = DueyWindowClosed
}

func (d *Duey) Fee(meso int32, quick bool) int32 {
	var permille int64
	switch {
	case meso >= 10000000:
		permille = 40
	case meso >= 5000000:
		permille = 30
	case meso >= 1000000:
		permille = 20
	case meso >= 100000:
		permille = 10
	case meso >= 50000:
		permille = 5
	}
	fee := int64(meso) * permille / 1000
	if !quick {
		fee += 5000
	}
	return int32(fee)
}

func (d *Duey) Send(actx actor.Context, invType constant.InventoryType, slot int16, count int16, meso int32, recipient string, quick bool, message string) error {
	if d.window != DueyWindowBox && d.window != DueyWindowQuick {
		return ErrDueyClosed
	}
	if d.pending {
		return ErrDueyBusy
	}
	if meso < 0 || count < 0 || recipient == "" || d.owner.PersistMapID() == 0 {
		return ErrDueyInvalid
	}
	if !quick && message != "" {
		return ErrDueyInvalid
	}
	messageBytes := 0
	for _, r := range message {
		if r < 0x80 {
			messageBytes++
		} else {
			messageBytes += 2
		}
	}
	if messageBytes > ParcelMessageBytes {
		return ErrDueyInvalid
	}

	var item Item
	if slot != 0 {
		item = d.owner.inventory.GetItem(invType, slot)
		if item == nil {
			return ErrDueyInvalid
		}
		if !item.Tradable() {
			return ErrDueyNotSendable
		}
		if item.GetModel().GetTraits()&constant.ItemTraitRechargeable != 0 {
			count = int16(item.GetCount())
		} else if count == 0 || item.GetCount() < uint16(count) {
			return ErrDueyInvalid
		}
	}
	if item == nil && meso == 0 {
		return ErrDueyInvalid
	}

	if quick && !d.owner.inventory.HasItem(QuickParcelCoupon) {
		return ErrDueyInvalid
	}

	cost := int64(meso) + int64(d.Fee(meso, quick))
	if cost > math.MaxInt32 || (ExchangeSpec{Cost: ExchangeSide{Meso: int32(cost)}}).Valid(d.owner, d.owner.escrow) != ExchangeOK {
		return ErrDueyNotEnoughMeso
	}

	parcel := &internal.Parcel{
		SenderName: d.owner.GetName(),
		Meso:       meso,
		Quick:      quick,
		Message:    message,
	}
	detail := parcelSendDetail{Recipient: recipient, Meso: int32(cost), Quick: quick}
	oneOfAKind := false
	var held Item
	if item != nil {
		held = item.Clone(uint16(count))
		if !d.owner.inventory.RemoveItem(invType, slot, uint16(count)) {
			return ErrDueyInvalid
		}
		parcel.Item = held.BuildRPC(0, 0)
		oneOfAKind = held.GetModel().GetTraits()&constant.ItemTraitUnique != 0
		detail.Item, _ = json.Marshal(parcel.Item)
	}
	paid := false
	if cost > 0 {
		if !d.owner.inventory.RemoveMeso(int32(cost)) {
			if held != nil {
				d.owner.inventory.addItemUnchecked(held, true)
			}
			return ErrDueyNotEnoughMeso
		}
		paid = true
	}
	if quick {
		couponType, slots := d.owner.inventory.FindSlots(QuickParcelCoupon)
		if len(slots) == 0 || !d.owner.inventory.RemoveItem(couponType, slots[0], 1) {
			if held != nil {
				d.owner.inventory.addItemUnchecked(held, true)
			}
			if paid {
				d.owner.inventory.addMesoUnchecked(int32(cost))
			}
			return ErrDueyInvalid
		}
	}
	id := d.owner.unsettled.Add(internal.UnsettledType_UNSETTLED_TYPE_PARCEL_SEND, 0, detail)

	d.pending = true
	d.owner.gameWorld.Systems().Get[*InternalSystem]().SendParcelAsync(actx, d.owner, id, d.owner.accountID, recipient, parcel, oneOfAKind, false).Bind(d.owner).Do(func(v *internal.SendParcelReply) error {
		d.pending = false
		result := v.GetResult()
		if result == internal.ParcelResult_PARCEL_RESULT_OK {
			d.owner.unsettled.Erase(id)
		} else {
			d.owner.Post(func(actx actor.Context) {
				d.returnSend(actx, d.owner.unsettled.Find(id))
			})
		}
		switch result {
		case internal.ParcelResult_PARCEL_RESULT_OK:
			d.owner.listener.OnDueyResult(d.owner, pconst.DueyResultSent)
		case internal.ParcelResult_PARCEL_RESULT_RECIPIENT_NOT_FOUND:
			d.owner.listener.OnDueyResult(d.owner, pconst.DueyResultRecipientNotFound)
		case internal.ParcelResult_PARCEL_RESULT_SAME_ACCOUNT:
			d.owner.listener.OnDueyResult(d.owner, pconst.DueyResultSameAccount)
		case internal.ParcelResult_PARCEL_RESULT_RECIPIENT_FULL:
			d.owner.listener.OnDueyResult(d.owner, pconst.DueyResultRecipientFull)
		case internal.ParcelResult_PARCEL_RESULT_ONE_OF_A_KIND:
			d.owner.listener.OnDueyResult(d.owner, pconst.DueyResultOnlyInBox)
		default:
			d.owner.listener.OnDueyResult(d.owner, pconst.DueyResultUnknown)
		}
		return nil
	}).OnError(func(err error) {
		d.pending = false
		d.owner.unsettled.Find(id).inFlight = false
		d.owner.listener.OnDueyResult(d.owner, pconst.DueyResultUnknown)

		log.Printf("Duey.Send character=%d: %v", d.owner.GetID(), err)
	})
	return nil
}

func systemParcel(item *internal.Inventory, meso int32) *internal.Parcel {
	return &internal.Parcel{
		SenderName: systemParcelSender,
		Meso:       meso,
		Quick:      true,
		Item:       item,
	}
}

func (d *Duey) Return(actx actor.Context, request *UnsettledRequest, parcels []*internal.Parcel) {
	if request == nil {
		return
	}
	kept := make([]*internal.Parcel, 0, len(parcels))
	for _, parcel := range parcels {
		if parcel != nil && (parcel.Item != nil || parcel.Meso > 0) {
			kept = append(kept, parcel)
		}
	}
	if len(kept) == 0 {
		d.owner.unsettled.Erase(request.id)
		return
	}

	request.inFlight = true
	d.owner.gameWorld.Systems().Get[*InternalSystem]().ReturnParcelAsync(actx, d.owner, request.id, kept).Bind(d.owner).Do(func(v *internal.ReturnParcelReply) error {
		if v.GetResult() == internal.ParcelResult_PARCEL_RESULT_OK {
			d.owner.unsettled.Erase(request.id)
			return nil
		}
		request.inFlight = false
		log.Printf("Duey.Return character=%d id=%d: %s", d.owner.GetID(), request.id, v.GetResult())
		return nil
	}).OnError(func(err error) {
		request.inFlight = false
		log.Printf("Duey.Return character=%d id=%d: %v", d.owner.GetID(), request.id, err)
	})
}

func (d *Duey) returnSend(actx actor.Context, request *UnsettledRequest) {
	if request == nil {
		return
	}
	var detail parcelSendDetail
	if err := json.Unmarshal([]byte(request.detail), &detail); err != nil {
		request.inFlight = false
		log.Printf("Duey.returnSend character=%d id=%d: %v", d.owner.GetID(), request.id, err)
		return
	}

	var item *internal.Inventory
	if len(detail.Item) > 0 {
		item = &internal.Inventory{}
		if err := json.Unmarshal(detail.Item, item); err != nil {
			request.inFlight = false
			log.Printf("Duey.returnSend character=%d id=%d: %v", d.owner.GetID(), request.id, err)
			return
		}
	}
	parcels := []*internal.Parcel{systemParcel(item, detail.Meso)}
	if detail.Quick {
		coupon, err := NewItem(QuickParcelCoupon, 1, d.owner.gameWorld)
		if err != nil {
			request.inFlight = false
			log.Printf("Duey.returnSend character=%d id=%d: %v", d.owner.GetID(), request.id, err)
			return
		}
		parcels = append(parcels, systemParcel(coupon.BuildRPC(0, 0), 0))
	}
	d.Return(actx, request, parcels)
}

func (d *Duey) settleSend(actx actor.Context, request *UnsettledRequest) {
	var detail parcelSendDetail
	if err := json.Unmarshal([]byte(request.detail), &detail); err != nil {
		log.Printf("Duey.settleSend character=%d id=%d: %v", d.owner.GetID(), request.id, err)
		return
	}

	d.owner.gameWorld.Systems().Get[*InternalSystem]().AbortParcelSendAsync(actx, d.owner, request.id, detail.Recipient).Bind(d.owner).Do(func(v *internal.AbortParcelSendReply) error {
		if v.GetSent() {
			d.owner.unsettled.Erase(request.id)
			return nil
		}
		d.returnSend(actx, request)
		return nil
	}).OnError(func(err error) {
		request.inFlight = false

		log.Printf("Duey.settleSend character=%d id=%d: %v", d.owner.GetID(), request.id, err)
	})
}

func (d *Duey) find(parcelID uint32) (int, *Parcel) {
	for i, parcel := range d.parcels {
		if parcel.id == parcelID {
			return i, parcel
		}
	}
	return -1, nil
}

func (d *Duey) Receive(actx actor.Context, parcelID uint32) error {
	if d.window != DueyWindowBox {
		return ErrDueyClosed
	}
	if d.pending {
		return ErrDueyBusy
	}
	_, parcel := d.find(parcelID)
	if parcel == nil {
		return ErrDueyParcelNotFound
	}
	if time.Now().Before(parcel.ArrivedAt()) {
		return ErrDueyCannotReceive
	}
	if err := d.checkReceive(parcel); err != nil {
		return err
	}

	id := d.owner.unsettled.Add(internal.UnsettledType_UNSETTLED_TYPE_PARCEL_CLAIM, uint64(parcelID), nil)
	d.pending = true
	d.owner.gameWorld.Systems().Get[*InternalSystem]().ClaimParcelAsync(actx, d.owner, id, parcelID).Bind(d.owner).Do(func(v *internal.ClaimParcelReply) error {
		d.pending = false
		if v.GetResult() != internal.ParcelResult_PARCEL_RESULT_OK {
			d.owner.unsettled.Erase(id)
			d.owner.listener.OnDueyResult(d.owner, pconst.DueyResultCannotReceive)
			return nil
		}

		if index, _ := d.find(parcelID); index >= 0 {
			d.parcels = append(d.parcels[:index:index], d.parcels[index+1:]...)
		}
		if err := d.receiveClaimed(id, v.GetParcel()); err != nil {
			d.owner.listener.OnDueyResult(d.owner, pconst.DueyResultInventoryFull)
			return nil
		}
		d.owner.listener.OnDueyRemoved(d.owner, parcelID, pconst.DueyRemovedReceived)
		return nil
	}).OnError(func(err error) {
		d.pending = false
		d.owner.unsettled.Find(id).inFlight = false
		d.owner.listener.OnDueyResult(d.owner, pconst.DueyResultUnknown)

		log.Printf("Duey.Receive character=%d: %v", d.owner.GetID(), err)
	})
	return nil
}

// A claimed parcel that does not fit stays unsettled; internal returns it again on the next settle.
func (d *Duey) receiveClaimed(id uint64, pb *internal.Parcel) error {
	claimed := NewParcelFromInternalProto(pb, d.owner.gameWorld)
	if err := d.checkReceive(claimed); err != nil {
		d.owner.unsettled.Find(id).inFlight = false
		return err
	}

	d.owner.inventory.addMesoUnchecked(claimed.meso)
	if claimed.item != nil {
		d.owner.inventory.receive(claimed.item)
	}
	d.owner.unsettled.Erase(id)
	return nil
}

func (d *Duey) settleClaim(actx actor.Context, request *UnsettledRequest) {
	d.owner.gameWorld.Systems().Get[*InternalSystem]().ClaimParcelAsync(actx, d.owner, request.id, uint32(request.target)).Bind(d.owner).Do(func(v *internal.ClaimParcelReply) error {
		if v.GetResult() != internal.ParcelResult_PARCEL_RESULT_OK {
			d.owner.unsettled.Erase(request.id)
			return nil
		}
		d.receiveClaimed(request.id, v.GetParcel())
		return nil
	}).OnError(func(err error) {
		request.inFlight = false

		log.Printf("Duey.settleClaim character=%d id=%d: %v", d.owner.GetID(), request.id, err)
	})
}

func (d *Duey) checkReceive(parcel *Parcel) error {
	if parcel.meso > math.MaxInt32-d.owner.inventory.meso {
		return ErrDueyCannotReceive
	}
	if parcel.item == nil {
		return nil
	}
	if d.owner.inventory.hasUniqueItem(parcel.item.GetModel()) {
		return ErrDueyOnlyHeld
	}
	var spec ExchangeSpec
	spec.Reward.include(parcel.item, parcel.item.GetCount())
	if spec.Valid(d.owner, d.owner.escrow) != ExchangeOK {
		return ErrInventoryFull
	}
	return nil
}

func (d *Duey) Delete(actx actor.Context, parcelID uint32) error {
	if d.window != DueyWindowBox {
		return ErrDueyClosed
	}
	if d.pending {
		return ErrDueyBusy
	}
	_, parcel := d.find(parcelID)
	if parcel == nil {
		return ErrDueyParcelNotFound
	}

	d.pending = true
	d.owner.gameWorld.Systems().Get[*InternalSystem]().DeleteParcelAsync(actx, d.owner, parcelID).Bind(d.owner).Do(func(v *internal.DeleteParcelReply) error {
		d.pending = false
		if index, _ := d.find(parcelID); index >= 0 {
			d.parcels = append(d.parcels[:index:index], d.parcels[index+1:]...)
		}
		if v.GetResult() != internal.ParcelResult_PARCEL_RESULT_OK {
			d.owner.listener.OnDueyResult(d.owner, pconst.DueyResultInvalid)
			return nil
		}
		d.owner.listener.OnDueyRemoved(d.owner, parcelID, pconst.DueyRemovedDeleted)
		return nil
	}).OnError(func(err error) {
		d.pending = false
		d.owner.listener.OnDueyResult(d.owner, pconst.DueyResultUnknown)

		log.Printf("Duey.Delete character=%d: %v", d.owner.GetID(), err)
	})
	return nil
}

func (d *Duey) SendFromSystem(actx actor.Context, recipient string, senderName string, itemID uint32, count uint16, meso int32, message string) error {
	if recipient == "" || senderName == "" || meso < 0 {
		return ErrDueyInvalid
	}
	parcel := &internal.Parcel{
		SenderName: senderName,
		Meso:       meso,
		Quick:      true,
		Message:    message,
	}
	if itemID != 0 {
		item, err := NewItem(itemID, count, d.owner.gameWorld)
		if err != nil || count == 0 {
			return ErrDueyInvalid
		}
		parcel.Item = item.BuildRPC(0, 0)
	}
	if parcel.Item == nil && meso == 0 {
		return ErrDueyInvalid
	}

	d.owner.gameWorld.Systems().Get[*InternalSystem]().SendParcelAsync(actx, d.owner, 0, 0, recipient, parcel, false, true).Bind(d.owner).Do(func(v *internal.SendParcelReply) error {
		result := v.GetResult()
		if result != internal.ParcelResult_PARCEL_RESULT_OK {
			log.Printf("Duey.SendFromSystem character=%d recipient=%s: %s", d.owner.GetID(), recipient, result)
		}
		return nil
	}).OnError(func(err error) {
		log.Printf("Duey.SendFromSystem character=%d recipient=%s: %v", d.owner.GetID(), recipient, err)
	})
	return nil
}

func (d *Duey) CheckArrivals(actx actor.Context) {
	d.owner.gameWorld.Systems().Get[*InternalSystem]().CheckParcelArrivalsAsync(actx, d.owner).Bind(d.owner).Do(func(v *internal.CheckParcelArrivalsReply) error {
		reply := v
		if reply.GetCount() == 0 {
			return nil
		}
		d.owner.listener.OnDueyArrival(d.owner, reply.GetSenderName(), reply.GetQuick(), int(reply.GetCount()))
		return nil
	}).OnError(func(err error) {
		log.Printf("Duey.CheckArrivals character=%d: %v", d.owner.GetID(), err)
	})
}

func (d *Duey) BuildS2C() (parcels, expired []dto.Parcel) {
	parcels = make([]dto.Parcel, 0, len(d.parcels))
	for _, parcel := range d.parcels {
		parcels = append(parcels, parcel.BuildS2C())
	}
	expired = make([]dto.Parcel, 0, len(d.expired))
	for _, parcel := range d.expired {
		expired = append(expired, parcel.BuildS2C())
	}
	return parcels, expired
}
