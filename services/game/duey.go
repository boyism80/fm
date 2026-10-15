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
	ID      uint32
	Sender  string
	Meso    int32
	Quick   bool
	Message string
	SentAt  time.Time
	Item    Item
}

func (p *Parcel) ArrivedAt() time.Time {
	if p.Quick {
		return p.SentAt
	}
	return p.SentAt.Add(ParcelArrivalDelay)
}

func (p *Parcel) ExpiresAt() time.Time {
	return p.ArrivedAt().Add(ParcelKeepDuration)
}

func (p *Parcel) BuildS2C() dto.Parcel {
	parcel := dto.Parcel{
		ID:      p.ID,
		Sender:  p.Sender,
		Meso:    p.Meso,
		Expire:  p.ExpiresAt(),
		Quick:   p.Quick,
		Message: p.Message,
	}
	if p.Item != nil {
		parcel.Item = p.Item.BuildS2C()
	}
	return parcel
}

type parcelSendDetail struct {
	Recipient string              `json:"recipient"`
	Meso      int32               `json:"meso"`
	Quick     bool                `json:"quick"`
	Item      *internal.Inventory `json:"item,omitempty"`
}

type Duey struct {
	owner       *Character
	Parcels     []*Parcel
	Expired     []*Parcel
	window      DueyWindow
	fromArrival bool
	pending     bool
}

func (d *Duey) Open(actx actor.Context, fromArrival bool) {
	if d.pending {
		return
	}

	d.pending = true
	d.owner.GameWorld.Systems().Get[*InternalSystem]().LoadParcelsAsync(actx, d.owner).Do(func(v *internal.LoadParcelsReply) error {
		d.owner.Resume(func(actor.Context) {
			d.pending = false
			reply := v
			d.Parcels = d.parcelsOf(reply.GetParcels())
			d.Expired = d.parcelsOf(reply.GetExpired())
			d.fromArrival = fromArrival
			if d.owner.GameWorld.DueyIdentityPrompt() {
				d.window = DueyWindowIdentity
				d.owner.Listener.OnDueyResult(d.owner, pconst.DueyResultIdentity)
				return
			}
			d.show()
		})
		return nil
	}).OnError(func(err error) {
		d.owner.Resume(func(actor.Context) {
			d.pending = false
		})
		log.Printf("Duey.Open character=%d: %v", d.owner.GetID(), err)
	})
}

func (d *Duey) parcelsOf(pbs []*internal.Parcel) []*Parcel {
	parcels := make([]*Parcel, 0, len(pbs))
	for _, pb := range pbs {
		parcels = append(parcels, NewParcelFromInternalProto(pb, d.owner.GameWorld))
	}
	return parcels
}

func (d *Duey) show() {
	d.window = DueyWindowBox
	d.owner.Listener.OnOpenDuey(d.owner, d.fromArrival)
	d.Expired = nil
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
	d.owner.Listener.OnDueyResult(d.owner, pconst.DueyResultQuickOpen)
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
		item = d.owner.Inventory.GetItem(invType, slot)
		if item == nil {
			return ErrDueyInvalid
		}
		if !item.Tradable() {
			return ErrDueyNotSendable
		}
		if count == 0 || item.GetCount() < uint16(count) {
			return ErrDueyInvalid
		}
		if item.GetModel().GetTraits()&constant.ItemTraitRechargeable != 0 {
			count = int16(item.GetCount())
		}
	}
	if item == nil && meso == 0 {
		return ErrDueyInvalid
	}

	if quick && !d.owner.Inventory.HasItem(QuickParcelCoupon) {
		return ErrDueyInvalid
	}

	cost := int64(meso) + int64(d.Fee(meso, quick))
	if cost > math.MaxInt32 || (ExchangeSpec{Cost: ExchangeSide{Meso: int32(cost)}}).Valid(d.owner) != ExchangeOK {
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
	if item != nil {
		parcel.Item = item.Clone(uint16(count)).BuildRPC(0, 0)
		oneOfAKind = item.GetModel().GetTraits()&constant.ItemTraitOnly != 0
		detail.Item = parcel.Item
		d.owner.Inventory.RemoveItem(invType, slot, uint16(count))
	}
	d.owner.Inventory.removeMesoUnchecked(int32(cost))
	if quick {
		d.owner.Inventory.RemoveByItemIDCount(QuickParcelCoupon, 1)
	}
	id := d.owner.Unsettled.Add(internal.UnsettledType_UNSETTLED_TYPE_PARCEL_SEND, 0, detail)

	d.pending = true
	d.owner.GameWorld.Systems().Get[*InternalSystem]().SendParcelAsync(actx, d.owner, id, d.owner.AccountID, recipient, parcel, oneOfAKind).Do(func(v *internal.SendParcelReply) error {
		d.owner.Resume(func(actor.Context) {
			d.pending = false
			result := v.GetResult()
			if result != internal.ParcelResult_PARCEL_RESULT_OK {
				d.refundSend(detail)
			}
			d.owner.Unsettled.Erase(id)
			switch result {
			case internal.ParcelResult_PARCEL_RESULT_OK:
				d.owner.Listener.OnDueyResult(d.owner, pconst.DueyResultSent)
			case internal.ParcelResult_PARCEL_RESULT_RECIPIENT_NOT_FOUND:
				d.owner.Listener.OnDueyResult(d.owner, pconst.DueyResultRecipientNotFound)
			case internal.ParcelResult_PARCEL_RESULT_SAME_ACCOUNT:
				d.owner.Listener.OnDueyResult(d.owner, pconst.DueyResultSameAccount)
			case internal.ParcelResult_PARCEL_RESULT_RECIPIENT_FULL:
				d.owner.Listener.OnDueyResult(d.owner, pconst.DueyResultRecipientFull)
			case internal.ParcelResult_PARCEL_RESULT_ONE_OF_A_KIND:
				d.owner.Listener.OnDueyResult(d.owner, pconst.DueyResultOnlyInBox)
			default:
				d.owner.Listener.OnDueyResult(d.owner, pconst.DueyResultUnknown)
			}
		})
		return nil
	}).OnError(func(err error) {
		d.owner.Resume(func(actor.Context) {
			d.pending = false
			d.owner.Unsettled.Find(id).InFlight = false
			d.owner.Listener.OnDueyResult(d.owner, pconst.DueyResultUnknown)
		})
		log.Printf("Duey.Send character=%d: %v", d.owner.GetID(), err)
	})
	return nil
}

func (d *Duey) refundSend(detail parcelSendDetail) {
	d.owner.Inventory.addMesoUnchecked(detail.Meso)
	if detail.Item != nil {
		if item, err := NewItemFromInternalProto(detail.Item, d.owner.GameWorld); err == nil {
			d.owner.Inventory.addItemUnchecked(item, true)
		}
	}
	if !detail.Quick {
		return
	}
	if coupon, err := NewItem(QuickParcelCoupon, 1, d.owner.GameWorld); err == nil {
		d.owner.Inventory.addItemUnchecked(coupon, true)
	}
}

func (d *Duey) settleSend(actx actor.Context, request *UnsettledRequest) {
	var detail parcelSendDetail
	if err := json.Unmarshal([]byte(request.Detail), &detail); err != nil {
		log.Printf("Duey.settleSend character=%d id=%d: %v", d.owner.GetID(), request.ID, err)
		return
	}

	d.owner.GameWorld.Systems().Get[*InternalSystem]().AbortParcelSendAsync(actx, d.owner, request.ID, detail.Recipient).Do(func(v *internal.AbortParcelSendReply) error {
		d.owner.Resume(func(actor.Context) {
			if !v.GetSent() {
				d.refundSend(detail)
			}
			d.owner.Unsettled.Erase(request.ID)
		})
		return nil
	}).OnError(func(err error) {
		d.owner.Resume(func(actor.Context) {
			request.InFlight = false
		})
		log.Printf("Duey.settleSend character=%d id=%d: %v", d.owner.GetID(), request.ID, err)
	})
}

func (d *Duey) find(parcelID uint32) (int, *Parcel) {
	for i, parcel := range d.Parcels {
		if parcel.ID == parcelID {
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

	id := d.owner.Unsettled.Add(internal.UnsettledType_UNSETTLED_TYPE_PARCEL_CLAIM, uint64(parcelID), nil)
	d.pending = true
	d.owner.GameWorld.Systems().Get[*InternalSystem]().ClaimParcelAsync(actx, d.owner, id, parcelID).Do(func(v *internal.ClaimParcelReply) error {
		d.owner.Resume(func(actor.Context) {
			d.pending = false
			if v.GetResult() != internal.ParcelResult_PARCEL_RESULT_OK {
				d.owner.Unsettled.Erase(id)
				d.owner.Listener.OnDueyResult(d.owner, pconst.DueyResultCannotReceive)
				return
			}

			if index, _ := d.find(parcelID); index >= 0 {
				d.Parcels = append(d.Parcels[:index:index], d.Parcels[index+1:]...)
			}
			if err := d.receiveClaimed(id, v.GetParcel()); err != nil {
				d.owner.Listener.OnDueyResult(d.owner, pconst.DueyResultInventoryFull)
				return
			}
			d.owner.Listener.OnDueyRemoved(d.owner, parcelID, pconst.DueyRemovedReceived)
		})
		return nil
	}).OnError(func(err error) {
		d.owner.Resume(func(actor.Context) {
			d.pending = false
			d.owner.Unsettled.Find(id).InFlight = false
			d.owner.Listener.OnDueyResult(d.owner, pconst.DueyResultUnknown)
		})
		log.Printf("Duey.Receive character=%d: %v", d.owner.GetID(), err)
	})
	return nil
}

// A claimed parcel that does not fit stays unsettled; internal returns it again on the next settle.
func (d *Duey) receiveClaimed(id uint64, pb *internal.Parcel) error {
	claimed := NewParcelFromInternalProto(pb, d.owner.GameWorld)
	if err := d.checkReceive(claimed); err != nil {
		d.owner.Unsettled.Find(id).InFlight = false
		return err
	}

	d.owner.Inventory.addMesoUnchecked(claimed.Meso)
	if claimed.Item != nil {
		d.owner.Inventory.receive(claimed.Item)
	}
	d.owner.Unsettled.Erase(id)
	return nil
}

func (d *Duey) settleClaim(actx actor.Context, request *UnsettledRequest) {
	d.owner.GameWorld.Systems().Get[*InternalSystem]().ClaimParcelAsync(actx, d.owner, request.ID, uint32(request.Target)).Do(func(v *internal.ClaimParcelReply) error {
		d.owner.Resume(func(actor.Context) {
			if v.GetResult() != internal.ParcelResult_PARCEL_RESULT_OK {
				d.owner.Unsettled.Erase(request.ID)
				return
			}
			d.receiveClaimed(request.ID, v.GetParcel())
		})
		return nil
	}).OnError(func(err error) {
		d.owner.Resume(func(actor.Context) {
			request.InFlight = false
		})
		log.Printf("Duey.settleClaim character=%d id=%d: %v", d.owner.GetID(), request.ID, err)
	})
}

func (d *Duey) checkReceive(parcel *Parcel) error {
	if parcel.Meso > math.MaxInt32-d.owner.Inventory.Meso {
		return ErrDueyCannotReceive
	}
	if parcel.Item == nil {
		return nil
	}
	if parcel.Item.GetModel().GetTraits()&constant.ItemTraitOnly != 0 && d.owner.Inventory.HasItem(parcel.Item.GetModel().GetID()) {
		return ErrDueyOnlyHeld
	}
	spec := ExchangeSpec{
		Reward: ExchangeSide{
			Items: map[uint32]uint16{parcel.Item.GetModel().GetID(): parcel.Item.GetCount()},
		},
	}
	if spec.Valid(d.owner) != ExchangeOK {
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
	d.owner.GameWorld.Systems().Get[*InternalSystem]().DeleteParcelAsync(actx, d.owner, parcelID).Do(func(v *internal.DeleteParcelReply) error {
		d.owner.Resume(func(actor.Context) {
			d.pending = false
			if index, _ := d.find(parcelID); index >= 0 {
				d.Parcels = append(d.Parcels[:index:index], d.Parcels[index+1:]...)
			}
			if v.GetResult() != internal.ParcelResult_PARCEL_RESULT_OK {
				d.owner.Listener.OnDueyResult(d.owner, pconst.DueyResultInvalid)
				return
			}
			d.owner.Listener.OnDueyRemoved(d.owner, parcelID, pconst.DueyRemovedDeleted)
		})
		return nil
	}).OnError(func(err error) {
		d.owner.Resume(func(actor.Context) {
			d.pending = false
			d.owner.Listener.OnDueyResult(d.owner, pconst.DueyResultUnknown)
		})
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
		item, err := NewItem(itemID, count, d.owner.GameWorld)
		if err != nil || count == 0 {
			return ErrDueyInvalid
		}
		parcel.Item = item.BuildRPC(0, 0)
	}
	if parcel.Item == nil && meso == 0 {
		return ErrDueyInvalid
	}

	d.owner.GameWorld.Systems().Get[*InternalSystem]().SendParcelAsync(actx, d.owner, 0, 0, recipient, parcel, false).Do(func(v *internal.SendParcelReply) error {
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
	d.owner.GameWorld.Systems().Get[*InternalSystem]().CheckParcelArrivalsAsync(actx, d.owner).Do(func(v *internal.CheckParcelArrivalsReply) error {
		reply := v
		if reply.GetCount() == 0 {
			return nil
		}
		d.owner.Listener.OnDueyArrival(d.owner, reply.GetSenderName(), reply.GetQuick(), int(reply.GetCount()))
		return nil
	}).OnError(func(err error) {
		log.Printf("Duey.CheckArrivals character=%d: %v", d.owner.GetID(), err)
	})
}

func (d *Duey) BuildS2C() (parcels, expired []dto.Parcel) {
	parcels = make([]dto.Parcel, 0, len(d.Parcels))
	for _, parcel := range d.Parcels {
		parcels = append(parcels, parcel.BuildS2C())
	}
	expired = make([]dto.Parcel, 0, len(d.Expired))
	for _, parcel := range d.Expired {
		expired = append(expired, parcel.BuildS2C())
	}
	return parcels, expired
}
