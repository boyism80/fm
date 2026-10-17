package game

import (
	"encoding/json"
	"errors"
	"log"
	"strings"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/core/async"
	"github.com/boyism80/fm/maple/constant"
	pconst "github.com/boyism80/fm/protocol/constant"
	"github.com/boyism80/fm/protocol/dto"
	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
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
	id     uint32
	sender string
	item   Item
}

type Wedding struct {
	owner    *Character
	marriage *Marriage
	proposal proposal
	gift     weddingGiftWindow
}

type weddingGiftDetail struct {
	Item json.RawMessage `json:"item"`
}

type weddingProposalDetail struct {
	Box  uint32 `json:"box,omitempty"`
	Ring uint32 `json:"ring"`
}

type weddingReservationDetail struct {
	Ticket uint32 `json:"ticket"`
}

type weddingInviteDetail struct {
	Invitation uint32          `json:"invitation"`
	Guest      string          `json:"guest"`
	Parcel     json.RawMessage `json:"parcel"`
}

type weddingWishlistDetail struct {
	Ticket uint32   `json:"ticket"`
	Wishes []string `json:"wishes"`
}

type weddingGiftWindow struct {
	receiverID uint32
	wishes     []string
	gifts      []*WeddingGift
	pending    bool
}

func (w *Wedding) BuildRingS2C() *dto.MarriageRing {
	if w.marriage == nil || w.marriage.status != MarriageStatusMarried {
		return nil
	}
	for _, equipment := range w.owner.inventory.equipped {
		if equipment == nil {
			continue
		}
		model := equipment.GetModel()
		if model.GetTraits()&constant.ItemTraitWeddingRing != 0 {
			return &dto.MarriageRing{CharacterID: w.owner.GetID(), PartnerID: w.marriage.PartnerID(w.owner.GetID()), ItemID: model.GetID()}
		}
	}
	return nil
}

func (w *Wedding) Propose(name string, boxItemID uint32) {
	if w.marriage != nil {
		if w.marriage.status == MarriageStatusMarried {
			w.owner.listener.OnEngageResult(w.owner, pconst.EngageResultAlreadyMarried)
			return
		}
		w.owner.listener.OnEngageResult(w.owner, pconst.EngageResultAlreadyEngaged)
		return
	}

	m := w.owner.GetMap()
	if m == nil {
		w.owner.listener.OnEngageResult(w.owner, pconst.EngageResultWrongName)
		return
	}
	var target *Character
	for _, obj := range m.GetAllPlayers() {
		if player, ok := obj.(*Character); ok && strings.EqualFold(player.GetName(), name) {
			target = player
			break
		}
	}
	if target == nil {
		w.owner.listener.OnEngageResult(w.owner, pconst.EngageResultWrongName)
		return
	}
	if target.GetGender() == w.owner.GetGender() {
		w.owner.listener.OnEngageResult(w.owner, pconst.EngageResultSameGender)
		return
	}
	if boxItemID < constant.RingBoxFirst || boxItemID > constant.RingBoxLast || !w.owner.inventory.HasItem(boxItemID) {
		w.owner.listener.OnEngageResult(w.owner, pconst.EngageResultBroken)
		return
	}
	if target.wedding.marriage != nil && target.wedding.marriage.status == MarriageStatusMarried {
		w.owner.listener.OnEngageResult(w.owner, pconst.EngageResultPartnerMarried)
		return
	}
	if target.wedding.marriage != nil || target.wedding.proposal.targetID != 0 {
		w.owner.listener.OnEngageResult(w.owner, pconst.EngageResultPartnerEngaged)
		return
	}

	ring := constant.EngagementRingFirst + (boxItemID - constant.RingBoxFirst)
	if !w.owner.inventory.canReceive(ring, 1) {
		w.owner.listener.OnEngageResult(w.owner, pconst.EngageResultInventoryFull)
		return
	}
	if !target.inventory.canReceive(ring, 1) {
		w.owner.listener.OnEngageResult(w.owner, pconst.EngageResultPartnerInventory)
		return
	}

	w.proposal = proposal{targetID: target.GetID(), boxItemID: boxItemID}
	target.listener.OnEngageRequest(target, w.owner.GetName(), w.owner.GetID())
}

func (w *Wedding) CancelProposal() {
	w.proposal = proposal{}
}

func (w *Wedding) AnswerProposal(actx actor.Context, accepted bool, name string, proposerID uint32) {
	m := w.owner.GetMap()
	if m == nil {
		w.owner.listener.OnEngageResult(w.owner, pconst.EngageResultProposalCancelled)
		return
	}
	proposer := m.GetPlayer(proposerID)
	if proposer == nil || !strings.EqualFold(proposer.GetName(), name) || proposer.wedding.proposal.targetID != w.owner.GetID() {
		w.owner.listener.OnEngageResult(w.owner, pconst.EngageResultProposalCancelled)
		return
	}
	boxItemID := proposer.wedding.proposal.boxItemID
	proposer.wedding.proposal = proposal{}
	if w.marriage != nil || proposer.wedding.marriage != nil || !proposer.inventory.HasItem(boxItemID) || !w.owner.IsAlive() || !proposer.IsAlive() {
		w.owner.listener.OnEngageResult(w.owner, pconst.EngageResultProposalCancelled)
		return
	}
	if !accepted {
		proposer.listener.OnEngageResult(proposer, pconst.EngageResultDeclined)
		return
	}

	ring := constant.EngagementRingFirst + (boxItemID - constant.RingBoxFirst)
	if !w.owner.inventory.canReceive(ring, 1) || !proposer.inventory.canReceive(ring, 1) {
		w.owner.listener.OnEngageResult(w.owner, pconst.EngageResultPartnerInventory)
		return
	}

	groom, bride := proposer, w.owner
	if w.owner.GetGender() == 0 {
		groom, bride = w.owner, proposer
	}

	proposer.inventory.RemoveByItemIDCount(boxItemID, 1)
	id := w.owner.unsettled.Add(internal.UnsettledType_UNSETTLED_TYPE_WEDDING_PROPOSAL, 0, weddingProposalDetail{Ring: ring})
	proposer.unsettled.add(id, internal.UnsettledType_UNSETTLED_TYPE_WEDDING_PROPOSAL, 0, weddingProposalDetail{Box: boxItemID, Ring: ring})

	members := []*Character{w.owner, proposer}
	w.owner.gameWorld.Systems().Get[*InternalSystem]().CreateMarriageAsync(actx, groom, bride, id, ring).Bind(w.owner).Do(func(v *internal.MarriageReply) error {
		engaged := v.GetResult() == internal.MarriageResult_MARRIAGE_RESULT_OK
		failure := pconst.EngageResultPartnerEngaged
		if v.GetResult() == internal.MarriageResult_MARRIAGE_RESULT_ALREADY_ENGAGED {
			failure = pconst.EngageResultAlreadyEngaged
		}
		for i, member := range members {
			member.Post(func(actor.Context) {
				member.wedding.closeProposal(id, engaged)
				if engaged {
					member.wedding.marriage = NewMarriageFromInternalProto(v.GetMarriage())
					member.listener.OnEngageResult(member, pconst.EngageResultEngaged)
				} else if i == 0 {
					member.listener.OnEngageResult(member, failure)
				}
			})
		}
		return nil
	}).OnError(func(err error) {
		for _, member := range members {
			member.Post(func(actor.Context) {
				if request := member.unsettled.Find(id); request != nil {
					request.inFlight = false
				}
			})
		}
		log.Printf("AnswerProposal character=%d: %v", w.owner.GetID(), err)
	})
}

// Both partners hold the proposal under the same id, and internal answers both with the same outcome.
func (w *Wedding) closeProposal(id uint64, engaged bool) {
	request := w.owner.unsettled.Find(id)
	if request == nil {
		return
	}
	w.owner.unsettled.Erase(id)

	var detail weddingProposalDetail
	if err := json.Unmarshal([]byte(request.detail), &detail); err != nil {
		log.Printf("Wedding.closeProposal character=%d id=%d: %v", w.owner.GetID(), id, err)
		return
	}
	itemID := detail.Ring
	if !engaged {
		itemID = detail.Box
	}
	if itemID == 0 {
		return
	}
	if item, err := NewItem(itemID, 1, w.owner.gameWorld); err == nil {
		w.owner.inventory.addItemUnchecked(item, true)
	}
}

func (w *Wedding) settleProposal(actx actor.Context, request *UnsettledRequest) {
	w.owner.gameWorld.Systems().Get[*InternalSystem]().AbortMarriageChangeAsync(actx, w.owner, request.id).Bind(w.owner).Do(func(v *internal.AbortMarriageChangeReply) error {
		w.closeProposal(request.id, v.GetChanged())
		if v.GetChanged() {
			w.Refresh(actx, "")
		}
		return nil
	}).OnError(func(err error) {
		request.inFlight = false

		log.Printf("Wedding.settleProposal character=%d id=%d: %v", w.owner.GetID(), request.id, err)
	})
}

func (w *Wedding) discardRings() {
	for parts, equipment := range w.owner.inventory.equipped {
		if equipment != nil && equipment.GetModel().GetTraits()&constant.ItemTraitWeddingRing != 0 {
			w.owner.inventory.RemoveEquipped(parts)
		}
	}
	for typ, inven := range w.owner.inventory.tabs {
		for slot, item := range inven.items {
			if item == nil {
				continue
			}
			if item.GetModel().GetTraits()&(constant.ItemTraitEngagementRing|constant.ItemTraitWeddingRing) != 0 {
				w.owner.inventory.RemoveItem(typ, slot, item.GetCount())
			}
		}
	}
}

func (w *Wedding) DropItem(actx actor.Context, itemID uint32) {
	model := w.owner.itemModel(itemID)
	if model == nil || model.GetTraits()&constant.ItemTraitMarriage == 0 || !w.owner.inventory.HasItem(itemID) {
		return
	}
	if model.GetTraits()&constant.ItemTraitEngagementRing == 0 || w.marriage == nil || w.marriage.status != MarriageStatusEngaged {
		w.owner.inventory.RemoveByItemIDCount(itemID, w.owner.inventory.GetCountByItemID(itemID))
		return
	}
	if w.marriage.ticketItemID != 0 {
		w.owner.listener.OnEngageResult(w.owner, pconst.EngageResultCannotCancel)
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
	if w.marriage == nil || w.marriage.status != MarriageStatusEngaged {
		return nil, ErrMarriageInvalid
	}

	return w.owner.gameWorld.Systems().Get[*InternalSystem]().BreakEngagementAsync(actx, w.owner, w.marriage.id).Bind(w.owner).Then(func(v *internal.MarriageReply) (bool, error) {
		broken := v.GetResult() == internal.MarriageResult_MARRIAGE_RESULT_OK
		if !broken {
			w.owner.listener.OnEngageResult(w.owner, pconst.EngageResultCannotCancel)
			return broken, nil
		}
		w.marriage = nil
		w.discardRings()
		w.owner.listener.OnEngageResult(w.owner, pconst.EngageResultBroken)

		return broken, nil
	}), nil
}

func (w *Wedding) Refresh(actx actor.Context, event string) {
	w.owner.gameWorld.Systems().Get[*InternalSystem]().LoadMarriageAsync(actx, w.owner).Bind(w.owner).Do(func(v *internal.MarriageReply) error {
		w.marriage = NewMarriageFromInternalProto(v.GetMarriage())
		if w.marriage == nil {
			w.discardRings()
		}
		switch event {
		case "broken":
			w.owner.listener.OnEngageResult(w.owner, pconst.EngageResultBroken)
		case "divorced":
			w.owner.listener.OnSpouseMap(w.owner, 0, 0)
			w.owner.listener.OnEngageResult(w.owner, pconst.EngageResultDivorced)
		case "booked":
			if w.marriage != nil && !w.marriage.Wished(w.owner.GetID()) {
				w.owner.listener.OnWeddingWishlistInput(w.owner)
			}
		case "reserved":
			w.owner.listener.OnEngageResult(w.owner, pconst.EngageResultReserved)
		case "married":
			if w.marriage != nil {
				w.owner.listener.OnEngageResult(w.owner, pconst.EngageResultMarried)
				w.NotifySpouseMap(actx, false)
			}
		}
		return nil
	}).OnError(func(err error) {
		log.Printf("RefreshMarriage character=%d: %v", w.owner.GetID(), err)
	})
}

func (w *Wedding) Reserve(actx actor.Context, ticketItemID uint32) (*async.Promise[bool], error) {
	if w.marriage == nil || w.marriage.status != MarriageStatusEngaged || w.marriage.ticketItemID != 0 {
		return nil, ErrMarriageInvalid
	}
	if _, ok := constant.WeddingTickets[ticketItemID]; !ok || !w.owner.inventory.HasItem(ticketItemID) {
		return nil, ErrMarriageInvalid
	}
	m := w.owner.GetMap()
	if m == nil {
		return nil, ErrMarriageInvalid
	}
	if m.GetPlayer(w.marriage.PartnerID(w.owner.GetID())) == nil {
		return nil, ErrMarriageInvalid
	}

	w.owner.inventory.RemoveByItemIDCount(ticketItemID, 1)
	id := w.owner.unsettled.Add(internal.UnsettledType_UNSETTLED_TYPE_WEDDING_RESERVE, uint64(w.marriage.id), weddingReservationDetail{Ticket: ticketItemID})

	return w.owner.gameWorld.Systems().Get[*InternalSystem]().ReserveWeddingAsync(actx, w.owner, id, w.marriage.id, ticketItemID).Bind(w.owner).Then(func(v *internal.MarriageReply) (bool, error) {
		reserved := v.GetResult() == internal.MarriageResult_MARRIAGE_RESULT_OK
		w.closeReservation(id, reserved)
		if !reserved {
			return reserved, nil
		}
		w.marriage = NewMarriageFromInternalProto(v.GetMarriage())
		w.owner.listener.OnWeddingWishlistInput(w.owner)

		return reserved, nil
	}).OnError(func(err error) {
		w.owner.unsettled.Find(id).inFlight = false

		log.Printf("ReserveWedding character=%d: %v", w.owner.GetID(), err)
	}), nil
}

func (w *Wedding) closeReservation(id uint64, reserved bool) {
	request := w.owner.unsettled.Find(id)
	if request == nil {
		return
	}
	w.owner.unsettled.Erase(id)
	if reserved {
		return
	}

	var detail weddingReservationDetail
	if err := json.Unmarshal([]byte(request.detail), &detail); err != nil {
		log.Printf("Wedding.closeReservation character=%d id=%d: %v", w.owner.GetID(), id, err)
		return
	}
	if item, err := NewItem(detail.Ticket, 1, w.owner.gameWorld); err == nil {
		w.owner.inventory.addItemUnchecked(item, true)
	}
}

func (w *Wedding) settleReservation(actx actor.Context, request *UnsettledRequest) {
	w.owner.gameWorld.Systems().Get[*InternalSystem]().AbortMarriageChangeAsync(actx, w.owner, request.id).Bind(w.owner).Do(func(v *internal.AbortMarriageChangeReply) error {
		w.closeReservation(request.id, v.GetChanged())
		if v.GetChanged() {
			w.Refresh(actx, "")
		}
		return nil
	}).OnError(func(err error) {
		request.inFlight = false

		log.Printf("Wedding.settleReservation character=%d id=%d: %v", w.owner.GetID(), request.id, err)
	})
}

func (w *Wedding) SubmitWishlist(actx actor.Context, wishes []string) {
	if w.marriage == nil || w.marriage.ticketItemID == 0 || w.marriage.Wished(w.owner.GetID()) {
		return
	}
	if w.owner.unsettled.FindByTarget(internal.UnsettledType_UNSETTLED_TYPE_WEDDING_WISHLIST, uint64(w.marriage.id)) != nil {
		return
	}
	ticket := constant.WeddingTickets[w.marriage.ticketItemID]
	if !w.owner.inventory.canReceive(ticket.Invitation, ticket.InvitationCount) {
		w.owner.listener.OnEngageResult(w.owner, pconst.EngageResultInventoryFull)
		return
	}
	if len(wishes) > 10 {
		wishes = wishes[:10]
	}

	id := w.owner.unsettled.Add(internal.UnsettledType_UNSETTLED_TYPE_WEDDING_WISHLIST, uint64(w.marriage.id), weddingWishlistDetail{Ticket: ticket.Invitation, Wishes: wishes})
	w.owner.gameWorld.Systems().Get[*InternalSystem]().SetWeddingWishlistAsync(actx, w.owner, id, w.marriage.id, wishes).Bind(w.owner).Do(func(v *internal.MarriageReply) error {
		if !w.completeWishlist(id, ticket.Invitation, v) {
			return nil
		}
		if w.marriage.Reserved() {
			w.owner.listener.OnEngageResult(w.owner, pconst.EngageResultReserved)
			return nil
		}
		w.owner.listener.OnMessage(w.owner, constant.MsgPopup, "위시리스트를 등록했습니다. 상대방이 위시리스트 등록을 끝낼 때까지 잠시 기다려주세요.")
		return nil
	}).OnError(func(err error) {
		if request := w.owner.unsettled.Find(id); request != nil {
			request.inFlight = false
		}

		log.Printf("SubmitWeddingWishlist character=%d: %v", w.owner.GetID(), err)
	})
}

func (w *Wedding) completeWishlist(id uint64, invitation uint32, reply *internal.MarriageReply) bool {
	request := w.owner.unsettled.Find(id)
	if request == nil {
		return false
	}
	if reply.GetResult() != internal.MarriageResult_MARRIAGE_RESULT_OK {
		w.owner.unsettled.Erase(id)
		return false
	}
	ticket := constant.WeddingTickets[reply.GetMarriage().GetTicketItemId()]
	if !w.owner.inventory.canReceive(invitation, ticket.InvitationCount) {
		request.inFlight = false
		w.owner.listener.OnEngageResult(w.owner, pconst.EngageResultInventoryFull)
		return false
	}
	item, err := NewItem(invitation, ticket.InvitationCount, w.owner.gameWorld)
	if err != nil {
		request.inFlight = false
		log.Printf("Wedding.completeWishlist character=%d id=%d: %v", w.owner.GetID(), id, err)
		return false
	}
	w.owner.inventory.addItemUnchecked(item, true)
	w.marriage = NewMarriageFromInternalProto(reply.GetMarriage())
	w.owner.unsettled.Erase(id)
	return true
}

func (w *Wedding) settleWishlist(actx actor.Context, request *UnsettledRequest) {
	var detail weddingWishlistDetail
	if err := json.Unmarshal([]byte(request.detail), &detail); err != nil {
		log.Printf("Wedding.settleWishlist character=%d id=%d: %v", w.owner.GetID(), request.id, err)
		return
	}
	w.owner.gameWorld.Systems().Get[*InternalSystem]().SetWeddingWishlistAsync(actx, w.owner, request.id, uint32(request.target), detail.Wishes).Bind(w.owner).Do(func(v *internal.MarriageReply) error {
		w.completeWishlist(request.id, detail.Ticket, v)
		return nil
	}).OnError(func(err error) {
		request.inFlight = false

		log.Printf("Wedding.settleWishlist character=%d id=%d: %v", w.owner.GetID(), request.id, err)
	})
}

func (w *Wedding) InviteGuest(actx actor.Context, guestName string, marriageID uint32, slot int16) {
	if w.marriage == nil || w.marriage.id != marriageID || w.marriage.ticketItemID == 0 {
		return
	}
	ticket := constant.WeddingTickets[w.marriage.ticketItemID]
	invitation := w.owner.inventory.GetItem(constant.InventoryTypeETC, slot)
	if invitation == nil || invitation.GetModel().GetID() != ticket.Invitation {
		return
	}
	invited, err := NewItem(ticket.InvitedItem, 1, w.owner.gameWorld)
	if err != nil {
		log.Printf("InviteWeddingGuest character=%d: %v", w.owner.GetID(), err)
		return
	}
	invited.(*MiscItem).marriageID = marriageID
	parcel := &internal.Parcel{
		SenderName: w.owner.GetName(),
		Quick:      true,
		Message:    "결혼식에 초대되었습니다. 청첩장을 확인해주세요.",
		Item:       invited.BuildRPC(0, 0),
	}

	w.owner.inventory.RemoveByItemIDCount(ticket.Invitation, 1)
	parcelJSON, _ := json.Marshal(parcel)
	id := w.owner.unsettled.Add(internal.UnsettledType_UNSETTLED_TYPE_WEDDING_INVITE, uint64(marriageID), weddingInviteDetail{Invitation: ticket.Invitation, Guest: guestName, Parcel: parcelJSON})
	w.owner.gameWorld.Systems().Get[*InternalSystem]().InviteWeddingGuestAsync(actx, w.owner, id, marriageID, guestName, parcel).Bind(w.owner).Do(func(v *internal.InviteWeddingGuestReply) error {
		w.closeInvite(id, v.GetResult())
		switch v.GetResult() {
		case internal.MarriageResult_MARRIAGE_RESULT_OK:
			w.owner.listener.OnMessage(w.owner, constant.MsgPopup, "청첩장을 보냈습니다.")
		case internal.MarriageResult_MARRIAGE_RESULT_GUEST_ALREADY_INVITED:
			w.owner.listener.OnMessage(w.owner, constant.MsgPopup, "대상은 이미 결혼식에 초대되었습니다.")
		default:
			w.owner.listener.OnEngageResult(w.owner, pconst.EngageResultWrongName)
		}
		return nil
	}).OnError(func(err error) {
		if request := w.owner.unsettled.Find(id); request != nil {
			request.inFlight = false
		}

		log.Printf("InviteWeddingGuest character=%d: %v", w.owner.GetID(), err)
	})
}

func (w *Wedding) closeInvite(id uint64, result internal.MarriageResult) {
	request := w.owner.unsettled.Find(id)
	if request == nil {
		return
	}
	w.owner.unsettled.Erase(id)
	if result == internal.MarriageResult_MARRIAGE_RESULT_OK {
		return
	}

	var detail weddingInviteDetail
	if err := json.Unmarshal([]byte(request.detail), &detail); err != nil {
		log.Printf("Wedding.closeInvite character=%d id=%d: %v", w.owner.GetID(), id, err)
		return
	}
	if item, err := NewItem(detail.Invitation, 1, w.owner.gameWorld); err == nil {
		w.owner.inventory.addItemUnchecked(item, true)
	}
}

func (w *Wedding) settleInvite(actx actor.Context, request *UnsettledRequest) {
	var detail weddingInviteDetail
	if err := json.Unmarshal([]byte(request.detail), &detail); err != nil {
		log.Printf("Wedding.settleInvite character=%d id=%d: %v", w.owner.GetID(), request.id, err)
		return
	}
	var parcel internal.Parcel
	if err := json.Unmarshal(detail.Parcel, &parcel); err != nil {
		log.Printf("Wedding.settleInvite character=%d id=%d: %v", w.owner.GetID(), request.id, err)
		return
	}
	w.owner.gameWorld.Systems().Get[*InternalSystem]().InviteWeddingGuestAsync(actx, w.owner, request.id, uint32(request.target), detail.Guest, &parcel).Bind(w.owner).Do(func(v *internal.InviteWeddingGuestReply) error {
		w.closeInvite(request.id, v.GetResult())
		return nil
	}).OnError(func(err error) {
		request.inFlight = false

		log.Printf("Wedding.settleInvite character=%d id=%d: %v", w.owner.GetID(), request.id, err)
	})
}

func (w *Wedding) InvitedTo(marriageID uint32) bool {
	for _, ticket := range constant.WeddingTickets {
		invType, slots := w.owner.inventory.FindSlots(ticket.InvitedItem)
		for _, slot := range slots {
			item, ok := w.owner.inventory.GetItem(invType, slot).(*MiscItem)
			if ok && item.marriageID == marriageID {
				return true
			}
		}
	}
	return false
}

func (w *Wedding) OpenInvitation(actx actor.Context, slot int16, itemID uint32) {
	item, ok := w.owner.inventory.GetItem(constant.InventoryTypeETC, slot).(*MiscItem)
	if !ok || item.GetModel().GetID() != itemID || item.marriageID == 0 {
		w.owner.listener.OnEngageResult(w.owner, pconst.EngageResultInvalidInvitation)
		return
	}

	w.owner.gameWorld.Systems().Get[*InternalSystem]().GetMarriageAsync(actx, w.owner, item.marriageID).Bind(w.owner).Do(func(v *internal.MarriageReply) error {
		marriage := NewMarriageFromInternalProto(v.GetMarriage())
		if marriage == nil {
			w.owner.listener.OnEngageResult(w.owner, pconst.EngageResultInvalidInvitation)
			return nil
		}
		weddingType := uint16(0)
		if marriage.ticketItemID != 0 {
			weddingType = uint16(marriage.ticketItemID - constant.WeddingTicketFirst)
		}
		w.owner.listener.OnWeddingInvitation(w.owner, marriage.groomName, marriage.brideName, weddingType)
		return nil
	}).OnError(func(err error) {
		log.Printf("OpenWeddingInvitation character=%d: %v", w.owner.GetID(), err)
	})
}

func (w *Wedding) RequestDivorce(actx actor.Context) (*async.Promise[string], error) {
	if w.marriage == nil || w.marriage.status != MarriageStatusMarried {
		return nil, ErrMarriageInvalid
	}

	return w.owner.gameWorld.Systems().Get[*InternalSystem]().RequestDivorceAsync(actx, w.owner, w.marriage.id).Bind(w.owner).Then(func(v *internal.MarriageReply) (string, error) {
		if v.GetResult() != internal.MarriageResult_MARRIAGE_RESULT_OK {
			return "pending", nil
		}
		marriage := NewMarriageFromInternalProto(v.GetMarriage())
		w.marriage = marriage
		if marriage != nil {
			return "requested", nil
		}
		w.discardRings()
		w.owner.listener.OnSpouseMap(w.owner, 0, 0)
		w.owner.listener.OnEngageResult(w.owner, pconst.EngageResultDivorced)

		return "divorced", nil
	}), nil
}

func (w *Wedding) CancelDivorce(actx actor.Context) (*async.Promise[bool], error) {
	if w.marriage == nil || w.marriage.status != MarriageStatusMarried || w.marriage.divorceRequestedAt.IsZero() {
		return nil, ErrMarriageInvalid
	}

	return w.owner.gameWorld.Systems().Get[*InternalSystem]().CancelDivorceAsync(actx, w.owner, w.marriage.id).Bind(w.owner).Then(func(reply *internal.MarriageReply) (bool, error) {
		if reply.GetResult() != internal.MarriageResult_MARRIAGE_RESULT_OK {
			return false, nil
		}
		w.marriage = NewMarriageFromInternalProto(reply.GetMarriage())

		return true, nil
	}), nil
}

func (w *Wedding) NotifySpouseMap(actx actor.Context, reply bool) {
	if w.marriage == nil || w.marriage.status != MarriageStatusMarried {
		return
	}
	m := w.owner.GetMap()
	if m == nil {
		return
	}
	w.owner.gameWorld.Systems().Get[*InternalSystem]().NotifySpouseMapAsync(actx, w.owner, w.marriage.PartnerID(w.owner.GetID()), m.TemplateID(), reply).Bind(w.owner).OnError(func(err error) {
		log.Printf("NotifySpouseMap character=%d: %v", w.owner.GetID(), err)
	})
}

func (w *Wedding) HideFromSpouse(actx actor.Context) {
	if w.marriage == nil || w.marriage.status != MarriageStatusMarried {
		return
	}
	w.owner.gameWorld.Systems().Get[*InternalSystem]().NotifySpouseMapAsync(actx, w.owner, w.marriage.PartnerID(w.owner.GetID()), 0, false).Bind(w.owner).OnError(func(err error) {
		log.Printf("HideFromSpouse character=%d: %v", w.owner.GetID(), err)
	})
}

func (w *Wedding) OpenGift(receiverID uint32, wishes []string) {
	w.gift = weddingGiftWindow{receiverID: receiverID, wishes: wishes}
	w.owner.listener.OnWeddingGift(w.owner, pconst.WeddingGiftOpenGive, wishes, nil)
}

func (w *Wedding) GiveGift(actx actor.Context, slot int16, itemID uint32, count uint16) error {
	if w.gift.receiverID == 0 {
		return ErrMarriageInvalid
	}
	if w.gift.pending {
		w.owner.listener.OnWeddingGift(w.owner, pconst.WeddingGiftOneAtATime, nil, nil)
		return ErrWeddingGiftBusy
	}
	model := w.owner.itemModel(itemID)
	if model == nil {
		return ErrMarriageInvalid
	}
	invType := model.GetInventoryType()
	item := w.owner.inventory.GetItem(invType, slot)
	if item == nil || item.GetModel().GetID() != itemID || count == 0 || item.GetCount() < count {
		return ErrMarriageInvalid
	}
	if !item.Tradable() {
		w.owner.listener.OnWeddingGift(w.owner, pconst.WeddingGiftGiveRejected, nil, nil)
		return nil
	}
	if model.GetTraits()&constant.ItemTraitRechargeable != 0 {
		count = item.GetCount()
	}

	held := item.Clone(count)
	if !w.owner.inventory.RemoveItem(invType, slot, count) {
		return ErrMarriageInvalid
	}
	pb := &internal.WeddingGift{SenderName: w.owner.GetName(), Item: held.BuildRPC(0, 0)}
	itemJSON, _ := json.Marshal(pb.Item)
	id := w.owner.unsettled.Add(internal.UnsettledType_UNSETTLED_TYPE_WEDDING_GIFT_SEND, uint64(w.gift.receiverID), weddingGiftDetail{Item: itemJSON})

	w.gift.pending = true
	w.owner.gameWorld.Systems().Get[*InternalSystem]().GiveWeddingGiftAsync(actx, w.owner, id, w.gift.receiverID, pb).Bind(w.owner).Do(func(v *internal.GiveWeddingGiftReply) error {
		w.gift.pending = false
		if v.GetResult() == internal.MarriageResult_MARRIAGE_RESULT_OK {
			w.owner.unsettled.Erase(id)
			w.owner.listener.OnWeddingGift(w.owner, pconst.WeddingGiftGiven, w.gift.wishes, map[constant.InventoryType][]Item{invType: {held}})
			return nil
		}
		w.owner.Post(func(actx actor.Context) {
			w.returnGift(actx, w.owner.unsettled.Find(id))
		})
		w.owner.listener.OnWeddingGift(w.owner, pconst.WeddingGiftGiveFailed, nil, nil)
		return nil
	}).OnError(func(err error) {
		w.gift.pending = false
		w.owner.unsettled.Find(id).inFlight = false
		w.owner.listener.OnWeddingGift(w.owner, pconst.WeddingGiftGiveFailed, nil, nil)

		log.Printf("GiveWeddingGift character=%d: %v", w.owner.GetID(), err)
	})
	return nil
}

func (w *Wedding) returnGift(actx actor.Context, request *UnsettledRequest) {
	if request == nil {
		return
	}
	var detail weddingGiftDetail
	if err := json.Unmarshal([]byte(request.detail), &detail); err != nil {
		request.inFlight = false
		log.Printf("Wedding.returnGift character=%d id=%d: %v", w.owner.GetID(), request.id, err)
		return
	}
	var item internal.Inventory
	if err := json.Unmarshal(detail.Item, &item); err != nil {
		request.inFlight = false
		log.Printf("Wedding.returnGift character=%d id=%d: %v", w.owner.GetID(), request.id, err)
		return
	}
	w.owner.duey.Return(actx, request, []*internal.Parcel{systemParcel(&item, 0)})
}

func (w *Wedding) settleGiftSend(actx actor.Context, request *UnsettledRequest) {
	w.owner.gameWorld.Systems().Get[*InternalSystem]().AbortWeddingGiftAsync(actx, w.owner, request.id, uint32(request.target)).Bind(w.owner).Do(func(v *internal.AbortWeddingGiftReply) error {
		if v.GetGiven() {
			w.owner.unsettled.Erase(request.id)
			return nil
		}
		w.returnGift(actx, request)
		return nil
	}).OnError(func(err error) {
		request.inFlight = false

		log.Printf("Wedding.settleGiftSend character=%d id=%d: %v", w.owner.GetID(), request.id, err)
	})
}

func (w *Wedding) giftTab(invType constant.InventoryType) []Item {
	items := make([]Item, 0)
	for _, gift := range w.gift.gifts {
		if gift.item.GetInventoryType() == invType {
			items = append(items, gift.item)
		}
	}
	return items
}

func (w *Wedding) OpenGiftBox(actx actor.Context) *async.Promise[bool] {
	return w.owner.gameWorld.Systems().Get[*InternalSystem]().LoadWeddingGiftsAsync(actx, w.owner).Bind(w.owner).Then(func(v *internal.LoadWeddingGiftsReply) (bool, error) {
		var gifts []*WeddingGift
		for _, pb := range v.GetGifts() {
			item, err := NewItemFromInternalProto(pb.GetItem(), w.owner.gameWorld)
			if err != nil {
				continue
			}
			gifts = append(gifts, &WeddingGift{id: pb.GetGiftId(), sender: pb.GetSenderName(), item: item})
		}

		w.gift = weddingGiftWindow{gifts: gifts}
		if len(gifts) == 0 {
			return false, nil
		}
		tabs := map[constant.InventoryType][]Item{}
		for typ := constant.InventoryTypeEquipment; typ <= constant.InventoryTypeCash; typ++ {
			tabs[typ] = w.giftTab(typ)
		}
		w.owner.listener.OnWeddingGift(w.owner, pconst.WeddingGiftOpenReceive, nil, tabs)

		return true, nil
	})
}

func (w *Wedding) ReceiveGift(actx actor.Context, invType constant.InventoryType, index int) error {
	if w.gift.pending {
		return ErrWeddingGiftBusy
	}
	var gift *WeddingGift
	for _, candidate := range w.gift.gifts {
		if candidate.item.GetInventoryType() != invType {
			continue
		}
		if index == 0 {
			gift = candidate
			break
		}
		index--
	}
	if gift == nil {
		w.owner.listener.OnWeddingGift(w.owner, pconst.WeddingGiftReceiveFail, nil, nil)
		return ErrWeddingGiftEmpty
	}
	if w.owner.inventory.hasOnlyItem(gift.item.GetModel()) {
		w.owner.listener.OnWeddingGift(w.owner, pconst.WeddingGiftOneOfAKind, nil, nil)
		return nil
	}
	if !w.owner.inventory.canReceive(gift.item.GetModel().GetID(), gift.item.GetCount()) {
		w.owner.listener.OnWeddingGift(w.owner, pconst.WeddingGiftReceiveFail, nil, nil)
		w.owner.listener.OnMessage(w.owner, constant.MsgPopup, "인벤토리 공간이 부족합니다.")
		return nil
	}

	id := w.owner.unsettled.Add(internal.UnsettledType_UNSETTLED_TYPE_WEDDING_GIFT_CLAIM, uint64(gift.id), nil)
	w.gift.pending = true
	w.owner.gameWorld.Systems().Get[*InternalSystem]().ClaimWeddingGiftAsync(actx, w.owner, id, gift.id).Bind(w.owner).Do(func(v *internal.ClaimWeddingGiftReply) error {
		w.gift.pending = false
		for i, candidate := range w.gift.gifts {
			if candidate == gift {
				w.gift.gifts = append(w.gift.gifts[:i:i], w.gift.gifts[i+1:]...)
				break
			}
		}
		if v.GetResult() != internal.MarriageResult_MARRIAGE_RESULT_OK {
			w.owner.unsettled.Erase(id)
			w.owner.listener.OnWeddingGift(w.owner, pconst.WeddingGiftReceiveFail, nil, nil)
			return nil
		}
		if !w.receiveClaimed(id, v.GetGift()) {
			w.owner.listener.OnWeddingGift(w.owner, pconst.WeddingGiftReceiveFail, nil, nil)
			return nil
		}
		w.owner.listener.OnWeddingGift(w.owner, pconst.WeddingGiftReceived, nil, map[constant.InventoryType][]Item{invType: w.giftTab(invType)})
		return nil
	}).OnError(func(err error) {
		w.gift.pending = false
		w.owner.unsettled.Find(id).inFlight = false
		w.owner.listener.OnWeddingGift(w.owner, pconst.WeddingGiftReceiveFail, nil, nil)

		log.Printf("ReceiveWeddingGift character=%d: %v", w.owner.GetID(), err)
	})
	return nil
}

// A claimed gift that does not fit stays unsettled; internal returns it again on the next settle.
func (w *Wedding) receiveClaimed(id uint64, pb *internal.WeddingGift) bool {
	item, err := NewItemFromInternalProto(pb.GetItem(), w.owner.gameWorld)
	if err != nil {
		log.Printf("Wedding.receiveClaimed character=%d id=%d: %v", w.owner.GetID(), id, err)
		w.owner.unsettled.Erase(id)
		return false
	}
	model := item.GetModel()
	if w.owner.inventory.hasOnlyItem(model) || !w.owner.inventory.canReceive(model.GetID(), item.GetCount()) {
		w.owner.unsettled.Find(id).inFlight = false
		return false
	}

	w.owner.inventory.receive(item)
	w.owner.unsettled.Erase(id)
	return true
}

func (w *Wedding) settleGiftClaim(actx actor.Context, request *UnsettledRequest) {
	w.owner.gameWorld.Systems().Get[*InternalSystem]().ClaimWeddingGiftAsync(actx, w.owner, request.id, uint32(request.target)).Bind(w.owner).Do(func(v *internal.ClaimWeddingGiftReply) error {
		if v.GetResult() != internal.MarriageResult_MARRIAGE_RESULT_OK {
			w.owner.unsettled.Erase(request.id)
			return nil
		}
		w.receiveClaimed(request.id, v.GetGift())
		return nil
	}).OnError(func(err error) {
		request.inFlight = false

		log.Printf("Wedding.settleGiftClaim character=%d id=%d: %v", w.owner.GetID(), request.id, err)
	})
}

func (w *Wedding) CloseGift() {
	w.gift = weddingGiftWindow{pending: w.gift.pending}
}
