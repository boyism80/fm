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

type weddingGiftWindow struct {
	receiverID uint32
	wishes     []string
	gifts      []*WeddingGift
	pending    bool
}

func (ch *Character) MarriageRingToDTO() *dto.MarriageRing {
	if ch.Marriage == nil || ch.Marriage.Status != MarriageStatusMarried {
		return nil
	}
	for _, equipment := range ch.Inventory.Equipped {
		if equipment == nil {
			continue
		}
		itemID := equipment.GetModel().GetID()
		if constant.IsWeddingRing(itemID) {
			return &dto.MarriageRing{CharacterID: ch.GetID(), PartnerID: ch.Marriage.PartnerID(ch.GetID()), ItemID: itemID}
		}
	}
	return nil
}

func (ch *Character) canReceive(itemID uint32, count uint16) bool {
	return (ExchangeSpec{Reward: ExchangeSide{Items: map[uint32]uint16{itemID: count}}}).Valid(ch) == ExchangeOK
}

func (ch *Character) Propose(name string, boxItemID uint32) {
	if ch.Marriage != nil {
		if ch.Marriage.Status == MarriageStatusMarried {
			ch.Listener.OnEngageResult(ch, pconst.EngageResultAlreadyMarried)
			return
		}
		ch.Listener.OnEngageResult(ch, pconst.EngageResultAlreadyEngaged)
		return
	}

	var target *Character
	for _, obj := range ch.GetMap().GetAllPlayers() {
		if player, ok := obj.(*Character); ok && strings.EqualFold(player.GetName(), name) {
			target = player
			break
		}
	}
	if target == nil {
		ch.Listener.OnEngageResult(ch, pconst.EngageResultWrongName)
		return
	}
	if target.GetGender() == ch.GetGender() {
		ch.Listener.OnEngageResult(ch, pconst.EngageResultSameGender)
		return
	}
	if boxItemID < constant.RingBoxFirst || boxItemID > constant.RingBoxLast || ch.Inventory.HasItem(boxItemID) == false {
		ch.Listener.OnEngageResult(ch, pconst.EngageResultBroken)
		return
	}
	if target.Marriage != nil && target.Marriage.Status == MarriageStatusMarried {
		ch.Listener.OnEngageResult(ch, pconst.EngageResultPartnerMarried)
		return
	}
	if target.Marriage != nil || target.proposal.targetID != 0 {
		ch.Listener.OnEngageResult(ch, pconst.EngageResultPartnerEngaged)
		return
	}

	ring := constant.EngagementRingFirst + (boxItemID - constant.RingBoxFirst)
	if ch.canReceive(ring, 1) == false {
		ch.Listener.OnEngageResult(ch, pconst.EngageResultInventoryFull)
		return
	}
	if target.canReceive(ring, 1) == false {
		ch.Listener.OnEngageResult(ch, pconst.EngageResultPartnerInventory)
		return
	}

	ch.proposal = proposal{targetID: target.GetID(), boxItemID: boxItemID}
	target.Listener.OnEngageRequest(target, ch.GetName(), ch.GetID())
}

func (ch *Character) CancelProposal() {
	ch.proposal = proposal{}
}

func (ch *Character) AnswerProposal(actx actor.Context, accepted bool, name string, proposerID uint32) {
	proposer := ch.GetMap().GetPlayer(proposerID)
	if proposer == nil || strings.EqualFold(proposer.GetName(), name) == false || proposer.proposal.targetID != ch.GetID() {
		ch.Listener.OnEngageResult(ch, pconst.EngageResultProposalCancelled)
		return
	}
	boxItemID := proposer.proposal.boxItemID
	proposer.proposal = proposal{}
	if ch.Marriage != nil || proposer.Marriage != nil || proposer.Inventory.HasItem(boxItemID) == false || ch.IsAlive() == false || proposer.IsAlive() == false {
		ch.Listener.OnEngageResult(ch, pconst.EngageResultProposalCancelled)
		return
	}
	if accepted == false {
		proposer.Listener.OnEngageResult(proposer, pconst.EngageResultDeclined)
		return
	}

	ring := constant.EngagementRingFirst + (boxItemID - constant.RingBoxFirst)
	if ch.canReceive(ring, 1) == false || proposer.canReceive(ring, 1) == false {
		ch.Listener.OnEngageResult(ch, pconst.EngageResultPartnerInventory)
		return
	}

	groom, bride := proposer, ch
	if ch.GetGender() == 0 {
		groom, bride = ch, proposer
	}
	ch.Listener.CreateMarriageAsync(actx, groom, bride, ring).Do(func(v *internal.MarriageReply) error {
		reply := v
		switch reply.GetResult() {
		case internal.MarriageResult_MARRIAGE_RESULT_OK:
		case internal.MarriageResult_MARRIAGE_RESULT_ALREADY_ENGAGED:
			ch.Listener.OnEngageResult(ch, pconst.EngageResultAlreadyEngaged)
			return nil
		default:
			ch.Listener.OnEngageResult(ch, pconst.EngageResultPartnerEngaged)
			return nil
		}

		for _, member := range []*Character{ch, proposer} {
			if member.GetMap() != ch.GetMap() {
				log.Printf("AnswerProposal character=%d left the map before the engagement was stored", member.GetID())
				continue
			}
			if member == proposer {
				member.Inventory.RemoveByItemIDCount(boxItemID, 1)
			}
			if item, err := NewItem(ring, 1, member.GameWorld); err == nil {
				member.Inventory.addItemUnchecked(item, true)
			}
			member.Marriage = NewMarriageFromInternalProto(reply.GetMarriage())
			member.Listener.OnEngageResult(member, pconst.EngageResultEngaged)
		}
		return nil
	}).OnError(func(err error) {
		log.Printf("AnswerProposal character=%d: %v", ch.GetID(), err)
	})
}

func (ch *Character) discardMarriageRings() {
	for parts, equipment := range ch.Inventory.Equipped {
		if equipment != nil && constant.IsWeddingRing(equipment.GetModel().GetID()) {
			ch.Inventory.RemoveEquipped(parts)
		}
	}
	for typ, inven := range ch.Inventory.Tabs {
		for slot, item := range inven.Items {
			if item == nil {
				continue
			}
			itemID := item.GetModel().GetID()
			if constant.IsEngagementRing(itemID) || constant.IsWeddingRing(itemID) {
				ch.Inventory.RemoveItem(typ, slot, item.GetCount())
			}
		}
	}
}

func (ch *Character) DropMarriageItem(actx actor.Context, itemID uint32) {
	if constant.ItemCategoryOf(itemID) != 421 || ch.Inventory.HasItem(itemID) == false {
		return
	}
	if constant.IsEngagementRing(itemID) == false || ch.Marriage == nil || ch.Marriage.Status != MarriageStatusEngaged {
		ch.Inventory.RemoveByItemIDCount(itemID, ch.Inventory.GetCountByItemID(itemID))
		return
	}
	if ch.Marriage.TicketItemID != 0 {
		ch.Listener.OnEngageResult(ch, pconst.EngageResultCannotCancel)
		return
	}

	promise, err := ch.BreakEngagement(actx)
	if err != nil {
		return
	}
	promise.OnError(func(err error) {
		log.Printf("DropMarriageItem character=%d: %v", ch.GetID(), err)
	})
}

func (ch *Character) BreakEngagement(actx actor.Context) (*async.Promise[bool], error) {
	if ch.Marriage == nil || ch.Marriage.Status != MarriageStatusEngaged {
		return nil, ErrMarriageInvalid
	}

	return ch.Listener.BreakEngagementAsync(actx, ch, ch.Marriage.ID).Then(func(v *internal.MarriageReply) (bool, error) {
		if v.GetResult() != internal.MarriageResult_MARRIAGE_RESULT_OK {
			ch.Listener.OnEngageResult(ch, pconst.EngageResultCannotCancel)
			return false, nil
		}
		ch.Marriage = nil
		ch.discardMarriageRings()
		ch.Listener.OnEngageResult(ch, pconst.EngageResultBroken)
		return true, nil
	}), nil
}

func (ch *Character) RefreshMarriage(actx actor.Context, event string) {
	ch.Listener.LoadMarriageAsync(actx, ch).Do(func(v *internal.MarriageReply) error {
		ch.Marriage = NewMarriageFromInternalProto(v.GetMarriage())
		if ch.Marriage == nil {
			ch.discardMarriageRings()
		}
		switch event {
		case "broken":
			ch.Listener.OnEngageResult(ch, pconst.EngageResultBroken)
		case "divorced":
			ch.Listener.OnEngageResult(ch, pconst.EngageResultDivorced)
		case "reserved":
			ch.Listener.OnEngageResult(ch, pconst.EngageResultReserved)
		case "married":
			if ch.Marriage != nil {
				ch.Listener.OnEngageResult(ch, pconst.EngageResultMarried)
			}
		}
		return nil
	}).OnError(func(err error) {
		log.Printf("RefreshMarriage character=%d: %v", ch.GetID(), err)
	})
}

func (ch *Character) ReserveWedding(actx actor.Context, ticketItemID uint32) (*async.Promise[bool], error) {
	if ch.Marriage == nil || ch.Marriage.Status != MarriageStatusEngaged || ch.Marriage.TicketItemID != 0 {
		return nil, ErrMarriageInvalid
	}
	if _, ok := constant.WeddingTickets[ticketItemID]; ok == false || ch.Inventory.HasItem(ticketItemID) == false {
		return nil, ErrMarriageInvalid
	}
	partner := ch.GetMap().GetPlayer(ch.Marriage.PartnerID(ch.GetID()))
	if partner == nil {
		return nil, ErrMarriageInvalid
	}

	return ch.Listener.ReserveWeddingAsync(actx, ch, ch.Marriage.ID, ticketItemID).Then(func(v *internal.MarriageReply) (bool, error) {
		reply := v
		if reply.GetResult() != internal.MarriageResult_MARRIAGE_RESULT_OK {
			return false, nil
		}
		ch.Inventory.RemoveByItemIDCount(ticketItemID, 1)
		for _, member := range []*Character{ch, partner} {
			if member != ch && member.GetMap() != ch.GetMap() {
				continue
			}
			member.Marriage = NewMarriageFromInternalProto(reply.GetMarriage())
			member.Listener.OnWeddingWishlistInput(member)
		}
		return true, nil
	}), nil
}

func (ch *Character) SubmitWeddingWishlist(actx actor.Context, wishes []string) {
	if ch.Marriage == nil || ch.Marriage.TicketItemID == 0 || ch.Marriage.Wished(ch.GetID()) {
		return
	}
	ticket := constant.WeddingTickets[ch.Marriage.TicketItemID]
	if ch.canReceive(ticket.Invitation, ticket.InvitationCount) == false {
		ch.Listener.OnEngageResult(ch, pconst.EngageResultInventoryFull)
		return
	}
	if len(wishes) > 10 {
		wishes = wishes[:10]
	}

	ch.Listener.SetWeddingWishlistAsync(actx, ch, ch.Marriage.ID, wishes).Do(func(v *internal.MarriageReply) error {
		reply := v
		if reply.GetResult() != internal.MarriageResult_MARRIAGE_RESULT_OK {
			return nil
		}
		ch.Marriage = NewMarriageFromInternalProto(reply.GetMarriage())
		if item, err := NewItem(ticket.Invitation, ticket.InvitationCount, ch.GameWorld); err == nil {
			ch.Inventory.addItemUnchecked(item, true)
		}
		if ch.Marriage.Reserved() {
			ch.Listener.OnEngageResult(ch, pconst.EngageResultReserved)
			return nil
		}
		ch.Listener.OnMessage(ch, constant.MsgPopup, "위시리스트를 등록했습니다. 상대방이 위시리스트 등록을 끝낼 때까지 잠시 기다려주세요.")
		return nil
	}).OnError(func(err error) {
		log.Printf("SubmitWeddingWishlist character=%d: %v", ch.GetID(), err)
	})
}

func (ch *Character) InviteWeddingGuest(actx actor.Context, guestName string, marriageID uint32, slot int16) {
	if ch.Marriage == nil || ch.Marriage.ID != marriageID || ch.Marriage.TicketItemID == 0 {
		return
	}
	ticket := constant.WeddingTickets[ch.Marriage.TicketItemID]
	invitation := ch.Inventory.GetItem(constant.InventoryTypeETC, slot)
	if invitation == nil || invitation.GetModel().GetID() != ticket.Invitation {
		return
	}

	ch.Listener.InviteWeddingGuestAsync(actx, ch, marriageID, guestName).Do(func(v *internal.InviteWeddingGuestReply) error {
		reply := v
		switch reply.GetResult() {
		case internal.MarriageResult_MARRIAGE_RESULT_OK:
		case internal.MarriageResult_MARRIAGE_RESULT_GUEST_ALREADY_INVITED:
			ch.Listener.OnMessage(ch, constant.MsgPopup, "대상은 이미 결혼식에 초대되었습니다.")
			return nil
		default:
			ch.Listener.OnEngageResult(ch, pconst.EngageResultWrongName)
			return nil
		}

		ch.Inventory.RemoveByItemIDCount(ticket.Invitation, 1)
		invited, err := NewItem(ticket.InvitedItem, 1, ch.GameWorld)
		if err != nil {
			return err
		}
		invited.(*MiscItem).MarriageID = marriageID
		parcel := &internal.Parcel{
			SenderName: ch.GetName(),
			Quick:      true,
			Message:    "결혼식에 초대되었습니다. 청첩장을 확인해주세요.",
			Item:       invited.ToProto(0, 0),
		}
		ch.Listener.SendParcelAsync(actx, ch, reply.GetGuestName(), parcel, false, nil).OnError(func(err error) {
			log.Printf("InviteWeddingGuest parcel character=%d guest=%s: %v", ch.GetID(), reply.GetGuestName(), err)
		})
		ch.Listener.OnMessage(ch, constant.MsgPopup, "청첩장을 보냈습니다.")
		return nil
	}).OnError(func(err error) {
		log.Printf("InviteWeddingGuest character=%d: %v", ch.GetID(), err)
	})
}

func (ch *Character) InvitedTo(marriageID uint32) bool {
	for _, ticket := range constant.WeddingTickets {
		invType, slots := ch.Inventory.FindSlots(ticket.InvitedItem)
		for _, slot := range slots {
			item, ok := ch.Inventory.GetItem(invType, slot).(*MiscItem)
			if ok && item.MarriageID == marriageID {
				return true
			}
		}
	}
	return false
}

func (ch *Character) OpenWeddingInvitation(actx actor.Context, slot int16, itemID uint32) {
	item, ok := ch.Inventory.GetItem(constant.InventoryTypeETC, slot).(*MiscItem)
	if ok == false || item.GetModel().GetID() != itemID || item.MarriageID == 0 {
		ch.Listener.OnEngageResult(ch, pconst.EngageResultInvalidInvitation)
		return
	}

	ch.Listener.GetMarriageAsync(actx, ch, item.MarriageID).Do(func(v *internal.MarriageReply) error {
		marriage := NewMarriageFromInternalProto(v.GetMarriage())
		if marriage == nil {
			ch.Listener.OnEngageResult(ch, pconst.EngageResultInvalidInvitation)
			return nil
		}
		weddingType := uint16(0)
		if marriage.TicketItemID != 0 {
			weddingType = uint16(marriage.TicketItemID - constant.WeddingTicketFirst)
		}
		ch.Listener.OnWeddingInvitation(ch, marriage.GroomName, marriage.BrideName, weddingType)
		return nil
	}).OnError(func(err error) {
		log.Printf("OpenWeddingInvitation character=%d: %v", ch.GetID(), err)
	})
}

func (ch *Character) RequestDivorce(actx actor.Context) (*async.Promise[string], error) {
	if ch.Marriage == nil || ch.Marriage.Status != MarriageStatusMarried {
		return nil, ErrMarriageInvalid
	}

	return ch.Listener.RequestDivorceAsync(actx, ch, ch.Marriage.ID).Then(func(v *internal.MarriageReply) (string, error) {
		reply := v
		if reply.GetResult() != internal.MarriageResult_MARRIAGE_RESULT_OK {
			return "pending", nil
		}
		ch.Marriage = NewMarriageFromInternalProto(reply.GetMarriage())
		if ch.Marriage != nil {
			return "requested", nil
		}
		ch.discardMarriageRings()
		ch.Listener.OnEngageResult(ch, pconst.EngageResultDivorced)
		return "divorced", nil
	}), nil
}

func (ch *Character) CancelDivorce(actx actor.Context) (*async.Promise[bool], error) {
	if ch.Marriage == nil || ch.Marriage.Status != MarriageStatusMarried || ch.Marriage.DivorceRequestedAt.IsZero() {
		return nil, ErrMarriageInvalid
	}

	return ch.Listener.CancelDivorceAsync(actx, ch, ch.Marriage.ID).Then(func(reply *internal.MarriageReply) (bool, error) {
		if reply.GetResult() != internal.MarriageResult_MARRIAGE_RESULT_OK {
			return false, nil
		}
		ch.Marriage = NewMarriageFromInternalProto(reply.GetMarriage())
		return true, nil
	}), nil
}

func (ch *Character) NotifySpouseMap(actx actor.Context, reply bool) {
	if ch.Marriage == nil {
		return
	}
	ch.Listener.NotifySpouseMapAsync(actx, ch, ch.Marriage.PartnerID(ch.GetID()), ch.GetMap().TemplateID(), reply).OnError(func(err error) {
		log.Printf("NotifySpouseMap character=%d: %v", ch.GetID(), err)
	})
}

func (ch *Character) OpenWeddingGift(receiverID uint32, wishes []string) {
	ch.weddingGift = weddingGiftWindow{receiverID: receiverID, wishes: wishes}
	ch.Listener.OnWeddingGift(ch, pconst.WeddingGiftOpenGive, wishes, nil)
}

func (ch *Character) GiveWeddingGift(actx actor.Context, slot int16, itemID uint32, count uint16) error {
	if ch.weddingGift.receiverID == 0 {
		return ErrMarriageInvalid
	}
	if ch.weddingGift.pending {
		ch.Listener.OnWeddingGift(ch, pconst.WeddingGiftOneAtATime, nil, nil)
		return ErrWeddingGiftBusy
	}
	invType := constant.GetInventoryTypeByItemID(itemID)
	item := ch.Inventory.GetItem(invType, slot)
	if item == nil || item.GetModel().GetID() != itemID || count == 0 || item.GetCount() < count {
		return ErrMarriageInvalid
	}
	model := item.GetModel()
	if model.IsTradeBlock() || model.IsQuest() || model.IsAccountSharable() || constant.ItemCategoryOf(itemID) == constant.ItemCategoryPet {
		ch.Listener.OnWeddingGift(ch, pconst.WeddingGiftGiveRejected, nil, nil)
		return nil
	}
	if constant.IsRechargeable(itemID) {
		count = item.GetCount()
	}

	gift := item.Clone(count)
	ch.Inventory.RemoveItem(invType, slot, count)
	ch.weddingGift.pending = true
	pb := &internal.WeddingGift{SenderName: ch.GetName(), Item: gift.ToProto(0, 0)}
	ch.Listener.GiveWeddingGiftAsync(actx, ch, ch.weddingGift.receiverID, pb, ch.ToProto(ch.GameWorld.GetWorldID())).Do(func(*internal.GiveWeddingGiftReply) error {
		ch.weddingGift.pending = false
		ch.Listener.OnWeddingGift(ch, pconst.WeddingGiftGiven, ch.weddingGift.wishes, map[constant.InventoryType][]Item{invType: {gift}})
		return nil
	}).OnError(func(err error) {
		ch.weddingGift.pending = false
		ch.Inventory.addItemUnchecked(gift, true)
		ch.Listener.OnWeddingGift(ch, pconst.WeddingGiftGiveFailed, nil, nil)
		log.Printf("GiveWeddingGift character=%d: %v", ch.GetID(), err)
	})
	return nil
}

func (ch *Character) weddingGiftTab(invType constant.InventoryType) []Item {
	items := make([]Item, 0)
	for _, gift := range ch.weddingGift.gifts {
		if gift.Item.GetInventoryType() == invType {
			items = append(items, gift.Item)
		}
	}
	return items
}

func (ch *Character) OpenWeddingGiftBox(actx actor.Context) *async.Promise[bool] {
	return ch.Listener.LoadWeddingGiftsAsync(actx, ch).Then(func(v *internal.LoadWeddingGiftsReply) (bool, error) {
		ch.weddingGift = weddingGiftWindow{}
		for _, pb := range v.GetGifts() {
			item, err := NewItemFromInternalProto(pb.GetItem(), ch.GameWorld)
			if err != nil {
				continue
			}
			ch.weddingGift.gifts = append(ch.weddingGift.gifts, &WeddingGift{ID: pb.GetGiftId(), Sender: pb.GetSenderName(), Item: item})
		}
		if len(ch.weddingGift.gifts) == 0 {
			return false, nil
		}

		tabs := map[constant.InventoryType][]Item{}
		for typ := constant.InventoryTypeEquipment; typ <= constant.InventoryTypeCash; typ++ {
			tabs[typ] = ch.weddingGiftTab(typ)
		}
		ch.Listener.OnWeddingGift(ch, pconst.WeddingGiftOpenReceive, nil, tabs)
		return true, nil
	})
}

func (ch *Character) ReceiveWeddingGift(actx actor.Context, invType constant.InventoryType, index int) error {
	if ch.weddingGift.pending {
		return ErrWeddingGiftBusy
	}
	var gift *WeddingGift
	for _, candidate := range ch.weddingGift.gifts {
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
		ch.Listener.OnWeddingGift(ch, pconst.WeddingGiftReceiveFail, nil, nil)
		return ErrWeddingGiftEmpty
	}
	if gift.Item.GetModel().IsOnly() && ch.Inventory.HasItem(gift.Item.GetModel().GetID()) {
		ch.Listener.OnWeddingGift(ch, pconst.WeddingGiftOneOfAKind, nil, nil)
		return nil
	}
	if ch.canReceive(gift.Item.GetModel().GetID(), gift.Item.GetCount()) == false {
		ch.Listener.OnWeddingGift(ch, pconst.WeddingGiftReceiveFail, nil, nil)
		ch.Listener.OnMessage(ch, constant.MsgPopup, "인벤토리 공간이 부족합니다.")
		return nil
	}

	ch.weddingGift.pending = true
	ch.Listener.ClaimWeddingGiftAsync(actx, ch, gift.ID).Do(func(v *internal.ClaimWeddingGiftReply) error {
		ch.weddingGift.pending = false
		for i, candidate := range ch.weddingGift.gifts {
			if candidate == gift {
				ch.weddingGift.gifts = append(ch.weddingGift.gifts[:i:i], ch.weddingGift.gifts[i+1:]...)
				break
			}
		}
		if v.GetResult() != internal.MarriageResult_MARRIAGE_RESULT_OK {
			ch.Listener.OnWeddingGift(ch, pconst.WeddingGiftReceiveFail, nil, nil)
			return nil
		}
		ch.Inventory.addItemUnchecked(gift.Item, true)
		ch.Listener.OnWeddingGift(ch, pconst.WeddingGiftReceived, nil, map[constant.InventoryType][]Item{invType: ch.weddingGiftTab(invType)})
		ch.GameWorld.SaveAsync(actx, []*internal.CharacterSaveEntry{ch.ToProto(ch.GameWorld.GetWorldID())}).OnError(func(err error) {
			log.Printf("ReceiveWeddingGift save character=%d: %v", ch.GetID(), err)
		})
		return nil
	}).OnError(func(err error) {
		ch.weddingGift.pending = false
		ch.Listener.OnWeddingGift(ch, pconst.WeddingGiftReceiveFail, nil, nil)
		log.Printf("ReceiveWeddingGift character=%d: %v", ch.GetID(), err)
	})
	return nil
}

func (ch *Character) CloseWeddingGift() {
	ch.weddingGift = weddingGiftWindow{pending: ch.weddingGift.pending}
}
