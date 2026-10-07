package entity

import (
	"errors"
	"log"
	"math"
	"time"

	"github.com/asynkron/protoactor-go/actor"
	pconst "github.com/boyism80/fm/protocol/constant"
	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
	"github.com/boyism80/fm/services/game/constant"
)

const (
	ParcelDeliveryDelay        = 12 * time.Hour
	ParcelKeepDuration         = 30 * 24 * time.Hour
	ParcelMessageBytes         = 100
	QuickDeliveryCoupon uint32 = 5330000
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
	return p.SentAt.Add(ParcelDeliveryDelay)
}

func (p *Parcel) ExpiresAt() time.Time {
	return p.ArrivedAt().Add(ParcelKeepDuration)
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
	d.owner.Listener.LoadParcelsAsync(actx, d.owner).Then(func(v interface{}) (interface{}, error) {
		d.pending = false
		reply := v.(*internal.LoadParcelsReply)
		d.Parcels = d.parcelsOf(reply.GetParcels())
		d.Expired = d.parcelsOf(reply.GetExpired())
		d.fromArrival = fromArrival
		if d.owner.GameWorld.DueyIdentityPrompt() {
			d.window = DueyWindowIdentity
			d.owner.Listener.OnDueyResult(d.owner, pconst.DueyResultIdentity)
			return nil, nil
		}
		d.show()
		return nil, nil
	}).OnError(func(err error) {
		d.pending = false
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
	if quick == false {
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
	if quick == false && message != "" {
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
		if item == nil || count == 0 || item.GetCount() < uint16(count) {
			return ErrDueyInvalid
		}
		model := item.GetModel()
		if model.IsTradeBlock() || model.IsAccountSharable() || model.IsQuest() {
			return ErrDueyNotSendable
		}
		if constant.ItemCategoryOf(model.GetID()) == constant.ItemCategoryPet {
			return ErrDueyNotSendable
		}
		if constant.IsRechargeable(model.GetID()) {
			count = int16(item.GetCount())
		}
	}
	if item == nil && meso == 0 {
		return ErrDueyInvalid
	}

	if quick && d.owner.Inventory.HasItem(QuickDeliveryCoupon) == false {
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
	var sent Item
	oneOfAKind := false
	if item != nil {
		sent = item.Clone(uint16(count))
		parcel.Item = sent.ToProto(0, 0)
		oneOfAKind = sent.GetModel().IsOnly()
		d.owner.Inventory.RemoveItem(invType, slot, uint16(count))
	}
	d.owner.Inventory.removeMesoUnchecked(int32(cost))
	if quick {
		d.owner.Inventory.RemoveByItemIDCount(QuickDeliveryCoupon, 1)
	}

	sender := d.owner.ToProto(d.owner.GameWorld.GetWorldID())
	refund := func() {
		d.owner.Inventory.addMesoUnchecked(int32(cost))
		if sent != nil {
			d.owner.Inventory.addItemUnchecked(sent, true)
		}
		if quick == false {
			return
		}
		if restored, err := NewItem(QuickDeliveryCoupon, 1, d.owner.GameWorld); err == nil {
			d.owner.Inventory.addItemUnchecked(restored, true)
		}
	}

	d.pending = true
	d.owner.Listener.SendParcelAsync(actx, d.owner, recipient, parcel, oneOfAKind, sender).Then(func(v interface{}) (interface{}, error) {
		d.pending = false
		result := v.(*internal.SendParcelReply).GetResult()
		if result != internal.ParcelResult_PARCEL_RESULT_OK {
			refund()
		}
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
		return nil, nil
	}).OnError(func(err error) {
		d.pending = false
		refund()
		d.owner.Listener.OnDueyResult(d.owner, pconst.DueyResultUnknown)
		log.Printf("Duey.Send character=%d: %v", d.owner.GetID(), err)
	})
	return nil
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

	d.pending = true
	d.owner.Listener.ClaimParcelAsync(actx, d.owner, parcelID).Then(func(v interface{}) (interface{}, error) {
		d.pending = false
		reply := v.(*internal.ClaimParcelReply)
		if reply.GetResult() != internal.ParcelResult_PARCEL_RESULT_OK {
			d.owner.Listener.OnDueyResult(d.owner, pconst.DueyResultCannotReceive)
			return nil, nil
		}

		if index, _ := d.find(parcelID); index >= 0 {
			d.Parcels = append(d.Parcels[:index:index], d.Parcels[index+1:]...)
		}
		claimed := NewParcelFromInternalProto(reply.GetParcel(), d.owner.GameWorld)
		if d.checkReceive(claimed) != nil {
			d.returnToBox(actx, reply.GetParcel())
			d.owner.Listener.OnDueyResult(d.owner, pconst.DueyResultInventoryFull)
			return nil, nil
		}
		d.owner.Inventory.addMesoUnchecked(claimed.Meso)
		if claimed.Item != nil {
			d.owner.Inventory.addItemUnchecked(claimed.Item, true)
		}
		d.owner.Listener.OnDueyRemoved(d.owner, parcelID, pconst.DueyRemovedReceived)
		d.owner.GameWorld.SaveAsync(actx, []*internal.CharacterSaveEntry{d.owner.ToProto(d.owner.GameWorld.GetWorldID())}).OnError(func(err error) {
			log.Printf("Duey.Receive save character=%d: %v", d.owner.GetID(), err)
		})
		return nil, nil
	}).OnError(func(err error) {
		d.pending = false
		d.owner.Listener.OnDueyResult(d.owner, pconst.DueyResultUnknown)
		log.Printf("Duey.Receive character=%d: %v", d.owner.GetID(), err)
	})
	return nil
}

func (d *Duey) checkReceive(parcel *Parcel) error {
	if parcel.Meso > math.MaxInt32-d.owner.Inventory.Meso {
		return ErrDueyCannotReceive
	}
	if parcel.Item == nil {
		return nil
	}
	if parcel.Item.GetModel().IsOnly() && d.owner.Inventory.HasItem(parcel.Item.GetModel().GetID()) {
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

func (d *Duey) returnToBox(actx actor.Context, pb *internal.Parcel) {
	pb.Quick = true
	d.owner.Listener.SendParcelAsync(actx, d.owner, d.owner.GetName(), pb, false, nil).OnError(func(err error) {
		log.Printf("Duey.returnToBox character=%d parcel=%d: %v", d.owner.GetID(), pb.GetParcelId(), err)
	})
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
	d.owner.Listener.DeleteParcelAsync(actx, d.owner, parcelID).Then(func(v interface{}) (interface{}, error) {
		d.pending = false
		if index, _ := d.find(parcelID); index >= 0 {
			d.Parcels = append(d.Parcels[:index:index], d.Parcels[index+1:]...)
		}
		if v.(*internal.DeleteParcelReply).GetResult() != internal.ParcelResult_PARCEL_RESULT_OK {
			d.owner.Listener.OnDueyResult(d.owner, pconst.DueyResultInvalid)
			return nil, nil
		}
		d.owner.Listener.OnDueyRemoved(d.owner, parcelID, pconst.DueyRemovedDeleted)
		return nil, nil
	}).OnError(func(err error) {
		d.pending = false
		d.owner.Listener.OnDueyResult(d.owner, pconst.DueyResultUnknown)
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
		parcel.Item = item.ToProto(0, 0)
	}
	if parcel.Item == nil && meso == 0 {
		return ErrDueyInvalid
	}

	d.owner.Listener.SendParcelAsync(actx, d.owner, recipient, parcel, false, nil).Then(func(v interface{}) (interface{}, error) {
		result := v.(*internal.SendParcelReply).GetResult()
		if result != internal.ParcelResult_PARCEL_RESULT_OK {
			log.Printf("Duey.SendFromSystem character=%d recipient=%s: %s", d.owner.GetID(), recipient, result)
		}
		return nil, nil
	}).OnError(func(err error) {
		log.Printf("Duey.SendFromSystem character=%d recipient=%s: %v", d.owner.GetID(), recipient, err)
	})
	return nil
}

func (d *Duey) CheckArrivals(actx actor.Context) {
	d.owner.Listener.CheckParcelArrivalsAsync(actx, d.owner).Then(func(v interface{}) (interface{}, error) {
		reply := v.(*internal.CheckParcelArrivalsReply)
		if reply.GetCount() == 0 {
			return nil, nil
		}
		d.owner.Listener.OnDueyArrival(d.owner, reply.GetSenderName(), reply.GetQuick(), int(reply.GetCount()))
		return nil, nil
	}).OnError(func(err error) {
		log.Printf("Duey.CheckArrivals character=%d: %v", d.owner.GetID(), err)
	})
}
