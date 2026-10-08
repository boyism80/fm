package entity

import (
	"errors"
	"log"
	"strings"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/core/async"
	pconst "github.com/boyism80/fm/protocol/constant"
	"github.com/boyism80/fm/protocol/dto"
	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
	"github.com/boyism80/fm/services/game/constant"
)

var (
	ErrMarriageInvalid  = errors.New("invalid marriage request")
	ErrWeddingGiftBusy  = errors.New("wedding gift request is pending")
	ErrWeddingGiftEmpty = errors.New("wedding gift not found")
)

type proposal struct {
	targetID  uint32
	boxItemID uint32
}

type WeddingGift struct {
	ID     uint32
	Sender string
	Item   Item
}

type Wedding struct {
	owner    *Character
	Marriage *Marriage
	proposal proposal
	gift     weddingGiftWindow
}

type weddingGiftWindow struct {
	receiverID uint32
	wishes     []string
	gifts      []*WeddingGift
	pending    bool
}

func (w *Wedding) RingToDTO() *dto.MarriageRing {
	if w.Marriage == nil || w.Marriage.Status != MarriageStatusMarried {
		return nil
	}
	for _, equipment := range w.owner.Inventory.Equipped {
		if equipment == nil {
			continue
		}
		itemID := equipment.GetModel().GetID()
		if constant.IsWeddingRing(itemID) {
			return &dto.MarriageRing{CharacterID: w.owner.GetID(), PartnerID: w.Marriage.PartnerID(w.owner.GetID()), ItemID: itemID}
		}
	}
	return nil
}

func (ch *Character) canReceive(itemID uint32, count uint16) bool {
	return (ExchangeSpec{Reward: ExchangeSide{Items: map[uint32]uint16{itemID: count}}}).Valid(ch) == ExchangeOK
}

func (w *Wedding) Propose(name string, boxItemID uint32) {
	if w.Marriage != nil {
		if w.Marriage.Status == MarriageStatusMarried {
			w.owner.Listener.OnEngageResult(w.owner, pconst.EngageResultAlreadyMarried)
			return
		}
		w.owner.Listener.OnEngageResult(w.owner, pconst.EngageResultAlreadyEngaged)
		return
	}

	var target *Character
	for _, obj := range w.owner.GetMap().GetAllPlayers() {
		if player, ok := obj.(*Character); ok && strings.EqualFold(player.GetName(), name) {
			target = player
			break
		}
	}
	if target == nil {
		w.owner.Listener.OnEngageResult(w.owner, pconst.EngageResultWrongName)
		return
	}
	if target.GetGender() == w.owner.GetGender() {
		w.owner.Listener.OnEngageResult(w.owner, pconst.EngageResultSameGender)
		return
	}
	if boxItemID < constant.RingBoxFirst || boxItemID > constant.RingBoxLast || w.owner.Inventory.HasItem(boxItemID) == false {
		w.owner.Listener.OnEngageResult(w.owner, pconst.EngageResultBroken)
		return
	}
	if target.Wedding.Marriage != nil && target.Wedding.Marriage.Status == MarriageStatusMarried {
		w.owner.Listener.OnEngageResult(w.owner, pconst.EngageResultPartnerMarried)
		return
	}
	if target.Wedding.Marriage != nil || target.Wedding.proposal.targetID != 0 {
		w.owner.Listener.OnEngageResult(w.owner, pconst.EngageResultPartnerEngaged)
		return
	}

	ring := constant.EngagementRingFirst + (boxItemID - constant.RingBoxFirst)
	if w.owner.canReceive(ring, 1) == false {
		w.owner.Listener.OnEngageResult(w.owner, pconst.EngageResultInventoryFull)
		return
	}
	if target.canReceive(ring, 1) == false {
		w.owner.Listener.OnEngageResult(w.owner, pconst.EngageResultPartnerInventory)
		return
	}

	w.proposal = proposal{targetID: target.GetID(), boxItemID: boxItemID}
	target.Listener.OnEngageRequest(target, w.owner.GetName(), w.owner.GetID())
}

func (w *Wedding) CancelProposal() {
	w.proposal = proposal{}
}

func (w *Wedding) AnswerProposal(actx actor.Context, accepted bool, name string, proposerID uint32) {
	proposer := w.owner.GetMap().GetPlayer(proposerID)
	if proposer == nil || strings.EqualFold(proposer.GetName(), name) == false || proposer.Wedding.proposal.targetID != w.owner.GetID() {
		w.owner.Listener.OnEngageResult(w.owner, pconst.EngageResultProposalCancelled)
		return
	}
	boxItemID := proposer.Wedding.proposal.boxItemID
	proposer.Wedding.proposal = proposal{}
	if w.Marriage != nil || proposer.Wedding.Marriage != nil || proposer.Inventory.HasItem(boxItemID) == false || w.owner.IsAlive() == false || proposer.IsAlive() == false {
		w.owner.Listener.OnEngageResult(w.owner, pconst.EngageResultProposalCancelled)
		return
	}
	if accepted == false {
		proposer.Listener.OnEngageResult(proposer, pconst.EngageResultDeclined)
		return
	}

	ring := constant.EngagementRingFirst + (boxItemID - constant.RingBoxFirst)
	if w.owner.canReceive(ring, 1) == false || proposer.canReceive(ring, 1) == false {
		w.owner.Listener.OnEngageResult(w.owner, pconst.EngageResultPartnerInventory)
		return
	}

	groom, bride := proposer, w.owner
	if w.owner.GetGender() == 0 {
		groom, bride = w.owner, proposer
	}
	w.owner.Listener.CreateMarriageAsync(actx, groom, bride, ring).Do(func(v *internal.MarriageReply) error {
		reply := v
		switch reply.GetResult() {
		case internal.MarriageResult_MARRIAGE_RESULT_OK:
		case internal.MarriageResult_MARRIAGE_RESULT_ALREADY_ENGAGED:
			w.owner.Listener.OnEngageResult(w.owner, pconst.EngageResultAlreadyEngaged)
			return nil
		default:
			w.owner.Listener.OnEngageResult(w.owner, pconst.EngageResultPartnerEngaged)
			return nil
		}

		for _, member := range []*Character{w.owner, proposer} {
			if member.GetMap() != w.owner.GetMap() {
				log.Printf("AnswerProposal character=%d left the map before the engagement was stored", member.GetID())
				continue
			}
			if member == proposer {
				member.Inventory.RemoveByItemIDCount(boxItemID, 1)
			}
			if item, err := NewItem(ring, 1, member.GameWorld); err == nil {
				member.Inventory.addItemUnchecked(item, true)
			}
			member.Wedding.Marriage = NewMarriageFromInternalProto(reply.GetMarriage())
			member.Listener.OnEngageResult(member, pconst.EngageResultEngaged)
		}
		return nil
	}).OnError(func(err error) {
		log.Printf("AnswerProposal character=%d: %v", w.owner.GetID(), err)
	})
}

func (w *Wedding) discardRings() {
	for parts, equipment := range w.owner.Inventory.Equipped {
		if equipment != nil && constant.IsWeddingRing(equipment.GetModel().GetID()) {
			w.owner.Inventory.RemoveEquipped(parts)
		}
	}
	for typ, inven := range w.owner.Inventory.Tabs {
		for slot, item := range inven.Items {
			if item == nil {
				continue
			}
			itemID := item.GetModel().GetID()
			if constant.IsEngagementRing(itemID) || constant.IsWeddingRing(itemID) {
				w.owner.Inventory.RemoveItem(typ, slot, item.GetCount())
			}
		}
	}
}

func (w *Wedding) DropItem(actx actor.Context, itemID uint32) {
	if constant.ItemCategoryOf(itemID) != 421 || w.owner.Inventory.HasItem(itemID) == false {
		return
	}
	if constant.IsEngagementRing(itemID) == false || w.Marriage == nil || w.Marriage.Status != MarriageStatusEngaged {
		w.owner.Inventory.RemoveByItemIDCount(itemID, w.owner.Inventory.GetCountByItemID(itemID))
		return
	}
	if w.Marriage.TicketItemID != 0 {
		w.owner.Listener.OnEngageResult(w.owner, pconst.EngageResultCannotCancel)
		return
	}

	promise, err := w.BreakEngagement(actx)
	if err != nil {
		return
	}
	promise.OnError(func(err error) {
		log.Printf("DropMarriageItem character=%d: %v", w.owner.GetID(), err)
	})
}

func (w *Wedding) BreakEngagement(actx actor.Context) (*async.Promise[bool], error) {
	if w.Marriage == nil || w.Marriage.Status != MarriageStatusEngaged {
		return nil, ErrMarriageInvalid
	}

	return w.owner.Listener.BreakEngagementAsync(actx, w.owner, w.Marriage.ID).Then(func(v *internal.MarriageReply) (bool, error) {
		if v.GetResult() != internal.MarriageResult_MARRIAGE_RESULT_OK {
			w.owner.Listener.OnEngageResult(w.owner, pconst.EngageResultCannotCancel)
			return false, nil
		}
		w.Marriage = nil
		w.discardRings()
		w.owner.Listener.OnEngageResult(w.owner, pconst.EngageResultBroken)
		return true, nil
	}), nil
}

func (w *Wedding) Refresh(actx actor.Context, event string) {
	w.owner.Listener.LoadMarriageAsync(actx, w.owner).Do(func(v *internal.MarriageReply) error {
		w.Marriage = NewMarriageFromInternalProto(v.GetMarriage())
		if w.Marriage == nil {
			w.discardRings()
		}
		switch event {
		case "broken":
			w.owner.Listener.OnEngageResult(w.owner, pconst.EngageResultBroken)
		case "divorced":
			w.owner.Listener.OnSpouseMap(w.owner, 0, 0)
			w.owner.Listener.OnEngageResult(w.owner, pconst.EngageResultDivorced)
		case "reserved":
			w.owner.Listener.OnEngageResult(w.owner, pconst.EngageResultReserved)
		case "married":
			if w.Marriage != nil {
				w.owner.Listener.OnEngageResult(w.owner, pconst.EngageResultMarried)
				w.NotifySpouseMap(actx, false)
			}
		}
		return nil
	}).OnError(func(err error) {
		log.Printf("RefreshMarriage character=%d: %v", w.owner.GetID(), err)
	})
}

func (w *Wedding) Reserve(actx actor.Context, ticketItemID uint32) (*async.Promise[bool], error) {
	if w.Marriage == nil || w.Marriage.Status != MarriageStatusEngaged || w.Marriage.TicketItemID != 0 {
		return nil, ErrMarriageInvalid
	}
	if _, ok := constant.WeddingTickets[ticketItemID]; ok == false || w.owner.Inventory.HasItem(ticketItemID) == false {
		return nil, ErrMarriageInvalid
	}
	partner := w.owner.GetMap().GetPlayer(w.Marriage.PartnerID(w.owner.GetID()))
	if partner == nil {
		return nil, ErrMarriageInvalid
	}

	return w.owner.Listener.ReserveWeddingAsync(actx, w.owner, w.Marriage.ID, ticketItemID).Then(func(v *internal.MarriageReply) (bool, error) {
		reply := v
		if reply.GetResult() != internal.MarriageResult_MARRIAGE_RESULT_OK {
			return false, nil
		}
		w.owner.Inventory.RemoveByItemIDCount(ticketItemID, 1)
		for _, member := range []*Character{w.owner, partner} {
			if member != w.owner && member.GetMap() != w.owner.GetMap() {
				continue
			}
			member.Wedding.Marriage = NewMarriageFromInternalProto(reply.GetMarriage())
			member.Listener.OnWeddingWishlistInput(member)
		}
		return true, nil
	}), nil
}

func (w *Wedding) SubmitWishlist(actx actor.Context, wishes []string) {
	if w.Marriage == nil || w.Marriage.TicketItemID == 0 || w.Marriage.Wished(w.owner.GetID()) {
		return
	}
	ticket := constant.WeddingTickets[w.Marriage.TicketItemID]
	if w.owner.canReceive(ticket.Invitation, ticket.InvitationCount) == false {
		w.owner.Listener.OnEngageResult(w.owner, pconst.EngageResultInventoryFull)
		return
	}
	if len(wishes) > 10 {
		wishes = wishes[:10]
	}

	w.owner.Listener.SetWeddingWishlistAsync(actx, w.owner, w.Marriage.ID, wishes).Do(func(v *internal.MarriageReply) error {
		reply := v
		if reply.GetResult() != internal.MarriageResult_MARRIAGE_RESULT_OK {
			return nil
		}
		w.Marriage = NewMarriageFromInternalProto(reply.GetMarriage())
		if item, err := NewItem(ticket.Invitation, ticket.InvitationCount, w.owner.GameWorld); err == nil {
			w.owner.Inventory.addItemUnchecked(item, true)
		}
		if w.Marriage.Reserved() {
			w.owner.Listener.OnEngageResult(w.owner, pconst.EngageResultReserved)
			return nil
		}
		w.owner.Listener.OnMessage(w.owner, constant.MsgPopup, "위시리스트를 등록했습니다. 상대방이 위시리스트 등록을 끝낼 때까지 잠시 기다려주세요.")
		return nil
	}).OnError(func(err error) {
		log.Printf("SubmitWeddingWishlist character=%d: %v", w.owner.GetID(), err)
	})
}

func (w *Wedding) InviteGuest(actx actor.Context, guestName string, marriageID uint32, slot int16) {
	if w.Marriage == nil || w.Marriage.ID != marriageID || w.Marriage.TicketItemID == 0 {
		return
	}
	ticket := constant.WeddingTickets[w.Marriage.TicketItemID]
	invitation := w.owner.Inventory.GetItem(constant.InventoryTypeETC, slot)
	if invitation == nil || invitation.GetModel().GetID() != ticket.Invitation {
		return
	}

	w.owner.Listener.InviteWeddingGuestAsync(actx, w.owner, marriageID, guestName).Do(func(v *internal.InviteWeddingGuestReply) error {
		reply := v
		switch reply.GetResult() {
		case internal.MarriageResult_MARRIAGE_RESULT_OK:
		case internal.MarriageResult_MARRIAGE_RESULT_GUEST_ALREADY_INVITED:
			w.owner.Listener.OnMessage(w.owner, constant.MsgPopup, "대상은 이미 결혼식에 초대되었습니다.")
			return nil
		default:
			w.owner.Listener.OnEngageResult(w.owner, pconst.EngageResultWrongName)
			return nil
		}

		w.owner.Inventory.RemoveByItemIDCount(ticket.Invitation, 1)
		invited, err := NewItem(ticket.InvitedItem, 1, w.owner.GameWorld)
		if err != nil {
			return err
		}
		invited.(*MiscItem).MarriageID = marriageID
		parcel := &internal.Parcel{
			SenderName: w.owner.GetName(),
			Quick:      true,
			Message:    "결혼식에 초대되었습니다. 청첩장을 확인해주세요.",
			Item:       invited.ToProto(0, 0),
		}
		w.owner.Listener.SendParcelAsync(actx, w.owner, reply.GetGuestName(), parcel, false, nil).OnError(func(err error) {
			log.Printf("InviteWeddingGuest parcel character=%d guest=%s: %v", w.owner.GetID(), reply.GetGuestName(), err)
		})
		w.owner.Listener.OnMessage(w.owner, constant.MsgPopup, "청첩장을 보냈습니다.")
		return nil
	}).OnError(func(err error) {
		log.Printf("InviteWeddingGuest character=%d: %v", w.owner.GetID(), err)
	})
}

func (w *Wedding) InvitedTo(marriageID uint32) bool {
	for _, ticket := range constant.WeddingTickets {
		invType, slots := w.owner.Inventory.FindSlots(ticket.InvitedItem)
		for _, slot := range slots {
			item, ok := w.owner.Inventory.GetItem(invType, slot).(*MiscItem)
			if ok && item.MarriageID == marriageID {
				return true
			}
		}
	}
	return false
}

func (w *Wedding) OpenInvitation(actx actor.Context, slot int16, itemID uint32) {
	item, ok := w.owner.Inventory.GetItem(constant.InventoryTypeETC, slot).(*MiscItem)
	if ok == false || item.GetModel().GetID() != itemID || item.MarriageID == 0 {
		w.owner.Listener.OnEngageResult(w.owner, pconst.EngageResultInvalidInvitation)
		return
	}

	w.owner.Listener.GetMarriageAsync(actx, w.owner, item.MarriageID).Do(func(v *internal.MarriageReply) error {
		marriage := NewMarriageFromInternalProto(v.GetMarriage())
		if marriage == nil {
			w.owner.Listener.OnEngageResult(w.owner, pconst.EngageResultInvalidInvitation)
			return nil
		}
		weddingType := uint16(0)
		if marriage.TicketItemID != 0 {
			weddingType = uint16(marriage.TicketItemID - constant.WeddingTicketFirst)
		}
		w.owner.Listener.OnWeddingInvitation(w.owner, marriage.GroomName, marriage.BrideName, weddingType)
		return nil
	}).OnError(func(err error) {
		log.Printf("OpenWeddingInvitation character=%d: %v", w.owner.GetID(), err)
	})
}

func (w *Wedding) RequestDivorce(actx actor.Context) (*async.Promise[string], error) {
	if w.Marriage == nil || w.Marriage.Status != MarriageStatusMarried {
		return nil, ErrMarriageInvalid
	}

	return w.owner.Listener.RequestDivorceAsync(actx, w.owner, w.Marriage.ID).Then(func(v *internal.MarriageReply) (string, error) {
		reply := v
		if reply.GetResult() != internal.MarriageResult_MARRIAGE_RESULT_OK {
			return "pending", nil
		}
		w.Marriage = NewMarriageFromInternalProto(reply.GetMarriage())
		if w.Marriage != nil {
			return "requested", nil
		}
		w.discardRings()
		w.owner.Listener.OnSpouseMap(w.owner, 0, 0)
		w.owner.Listener.OnEngageResult(w.owner, pconst.EngageResultDivorced)
		return "divorced", nil
	}), nil
}

func (w *Wedding) CancelDivorce(actx actor.Context) (*async.Promise[bool], error) {
	if w.Marriage == nil || w.Marriage.Status != MarriageStatusMarried || w.Marriage.DivorceRequestedAt.IsZero() {
		return nil, ErrMarriageInvalid
	}

	return w.owner.Listener.CancelDivorceAsync(actx, w.owner, w.Marriage.ID).Then(func(reply *internal.MarriageReply) (bool, error) {
		if reply.GetResult() != internal.MarriageResult_MARRIAGE_RESULT_OK {
			return false, nil
		}
		w.Marriage = NewMarriageFromInternalProto(reply.GetMarriage())
		return true, nil
	}), nil
}

func (w *Wedding) NotifySpouseMap(actx actor.Context, reply bool) {
	if w.Marriage == nil || w.Marriage.Status != MarriageStatusMarried {
		return
	}
	w.owner.Listener.NotifySpouseMapAsync(actx, w.owner, w.Marriage.PartnerID(w.owner.GetID()), w.owner.GetMap().TemplateID(), reply).OnError(func(err error) {
		log.Printf("NotifySpouseMap character=%d: %v", w.owner.GetID(), err)
	})
}

func (w *Wedding) HideFromSpouse(actx actor.Context) {
	if w.Marriage == nil || w.Marriage.Status != MarriageStatusMarried {
		return
	}
	w.owner.Listener.NotifySpouseMapAsync(actx, w.owner, w.Marriage.PartnerID(w.owner.GetID()), 0, false).OnError(func(err error) {
		log.Printf("HideFromSpouse character=%d: %v", w.owner.GetID(), err)
	})
}

func (w *Wedding) OpenGift(receiverID uint32, wishes []string) {
	w.gift = weddingGiftWindow{receiverID: receiverID, wishes: wishes}
	w.owner.Listener.OnWeddingGift(w.owner, pconst.WeddingGiftOpenGive, wishes, nil)
}

func (w *Wedding) GiveGift(actx actor.Context, slot int16, itemID uint32, count uint16) error {
	if w.gift.receiverID == 0 {
		return ErrMarriageInvalid
	}
	if w.gift.pending {
		w.owner.Listener.OnWeddingGift(w.owner, pconst.WeddingGiftOneAtATime, nil, nil)
		return ErrWeddingGiftBusy
	}
	invType := constant.GetInventoryTypeByItemID(itemID)
	item := w.owner.Inventory.GetItem(invType, slot)
	if item == nil || item.GetModel().GetID() != itemID || count == 0 || item.GetCount() < count {
		return ErrMarriageInvalid
	}
	model := item.GetModel()
	if model.IsTradeBlock() || model.IsQuest() || model.IsAccountSharable() || constant.ItemCategoryOf(itemID) == constant.ItemCategoryPet {
		w.owner.Listener.OnWeddingGift(w.owner, pconst.WeddingGiftGiveRejected, nil, nil)
		return nil
	}
	if constant.IsRechargeable(itemID) {
		count = item.GetCount()
	}

	gift := item.Clone(count)
	w.owner.Inventory.RemoveItem(invType, slot, count)
	w.gift.pending = true
	pb := &internal.WeddingGift{SenderName: w.owner.GetName(), Item: gift.ToProto(0, 0)}
	w.owner.Listener.GiveWeddingGiftAsync(actx, w.owner, w.gift.receiverID, pb, w.owner.ToProto(w.owner.GameWorld.GetWorldID())).Do(func(*internal.GiveWeddingGiftReply) error {
		w.gift.pending = false
		w.owner.Listener.OnWeddingGift(w.owner, pconst.WeddingGiftGiven, w.gift.wishes, map[constant.InventoryType][]Item{invType: {gift}})
		return nil
	}).OnError(func(err error) {
		w.gift.pending = false
		w.owner.Inventory.addItemUnchecked(gift, true)
		w.owner.Listener.OnWeddingGift(w.owner, pconst.WeddingGiftGiveFailed, nil, nil)
		log.Printf("GiveWeddingGift character=%d: %v", w.owner.GetID(), err)
	})
	return nil
}

func (w *Wedding) giftTab(invType constant.InventoryType) []Item {
	items := make([]Item, 0)
	for _, gift := range w.gift.gifts {
		if gift.Item.GetInventoryType() == invType {
			items = append(items, gift.Item)
		}
	}
	return items
}

func (w *Wedding) OpenGiftBox(actx actor.Context) *async.Promise[bool] {
	return w.owner.Listener.LoadWeddingGiftsAsync(actx, w.owner).Then(func(v *internal.LoadWeddingGiftsReply) (bool, error) {
		w.gift = weddingGiftWindow{}
		for _, pb := range v.GetGifts() {
			item, err := NewItemFromInternalProto(pb.GetItem(), w.owner.GameWorld)
			if err != nil {
				continue
			}
			w.gift.gifts = append(w.gift.gifts, &WeddingGift{ID: pb.GetGiftId(), Sender: pb.GetSenderName(), Item: item})
		}
		if len(w.gift.gifts) == 0 {
			return false, nil
		}

		tabs := map[constant.InventoryType][]Item{}
		for typ := constant.InventoryTypeEquipment; typ <= constant.InventoryTypeCash; typ++ {
			tabs[typ] = w.giftTab(typ)
		}
		w.owner.Listener.OnWeddingGift(w.owner, pconst.WeddingGiftOpenReceive, nil, tabs)
		return true, nil
	})
}

func (w *Wedding) ReceiveGift(actx actor.Context, invType constant.InventoryType, index int) error {
	if w.gift.pending {
		return ErrWeddingGiftBusy
	}
	var gift *WeddingGift
	for _, candidate := range w.gift.gifts {
		if candidate.Item.GetInventoryType() != invType {
			continue
		}
		if index == 0 {
			gift = candidate
			break
		}
		index--
	}
	if gift == nil {
		w.owner.Listener.OnWeddingGift(w.owner, pconst.WeddingGiftReceiveFail, nil, nil)
		return ErrWeddingGiftEmpty
	}
	if gift.Item.GetModel().IsOnly() && w.owner.Inventory.HasItem(gift.Item.GetModel().GetID()) {
		w.owner.Listener.OnWeddingGift(w.owner, pconst.WeddingGiftOneOfAKind, nil, nil)
		return nil
	}
	if w.owner.canReceive(gift.Item.GetModel().GetID(), gift.Item.GetCount()) == false {
		w.owner.Listener.OnWeddingGift(w.owner, pconst.WeddingGiftReceiveFail, nil, nil)
		w.owner.Listener.OnMessage(w.owner, constant.MsgPopup, "인벤토리 공간이 부족합니다.")
		return nil
	}

	w.gift.pending = true
	w.owner.Listener.ClaimWeddingGiftAsync(actx, w.owner, gift.ID).Do(func(v *internal.ClaimWeddingGiftReply) error {
		w.gift.pending = false
		for i, candidate := range w.gift.gifts {
			if candidate == gift {
				w.gift.gifts = append(w.gift.gifts[:i:i], w.gift.gifts[i+1:]...)
				break
			}
		}
		if v.GetResult() != internal.MarriageResult_MARRIAGE_RESULT_OK {
			w.owner.Listener.OnWeddingGift(w.owner, pconst.WeddingGiftReceiveFail, nil, nil)
			return nil
		}
		w.owner.Inventory.addItemUnchecked(gift.Item, true)
		w.owner.Listener.OnWeddingGift(w.owner, pconst.WeddingGiftReceived, nil, map[constant.InventoryType][]Item{invType: w.giftTab(invType)})
		w.owner.GameWorld.SaveAsync(actx, []*internal.CharacterSaveEntry{w.owner.ToProto(w.owner.GameWorld.GetWorldID())}).OnError(func(err error) {
			log.Printf("ReceiveWeddingGift save character=%d: %v", w.owner.GetID(), err)
		})
		return nil
	}).OnError(func(err error) {
		w.gift.pending = false
		w.owner.Listener.OnWeddingGift(w.owner, pconst.WeddingGiftReceiveFail, nil, nil)
		log.Printf("ReceiveWeddingGift character=%d: %v", w.owner.GetID(), err)
	})
	return nil
}

func (w *Wedding) CloseGift() {
	w.gift = weddingGiftWindow{pending: w.gift.pending}
}
