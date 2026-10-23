package game

import (
	"encoding/json"
	"log"

	"github.com/asynkron/protoactor-go/actor"
	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
)

const unsettledSettleTimer = "unsettled_settle"

type UnsettledRequest struct {
	id       uint64
	onType   internal.UnsettledType
	target   uint64
	detail   string
	inFlight bool
}

type UnsettledRequests struct {
	owner    *Character
	requests map[uint64]*UnsettledRequest
}

func (u *UnsettledRequests) Add(typ internal.UnsettledType, target uint64, detail any) uint64 {
	id := u.owner.gameWorld.NewUniqueID()
	u.add(id, typ, target, detail)
	return id
}

// One remote request that holds value on two characters is tracked on both with the same id.
func (u *UnsettledRequests) add(id uint64, typ internal.UnsettledType, target uint64, detail any) {
	data, err := json.Marshal(detail)
	if err != nil {
		log.Printf("UnsettledRequests.Add character=%d type=%s: %v", u.owner.GetID(), typ, err)
	}

	u.requests[id] = &UnsettledRequest{id: id, onType: typ, target: target, detail: string(data), inFlight: true}
}

func (u *UnsettledRequests) Lost(id uint64) {
	request := u.Find(id)
	if request == nil {
		return
	}
	request.inFlight = false
}

func (u *UnsettledRequests) Find(id uint64) *UnsettledRequest {
	return u.requests[id]
}

func (u *UnsettledRequests) FindByTarget(typ internal.UnsettledType, target uint64) *UnsettledRequest {
	for _, request := range u.requests {
		if request.onType == typ && request.target == target {
			return request
		}
	}
	return nil
}

func (u *UnsettledRequests) Erase(id uint64) {
	delete(u.requests, id)
}

func (u *UnsettledRequests) Settle(actx actor.Context) {
	for _, request := range u.requests {
		if request.inFlight {
			continue
		}

		request.inFlight = true
		switch request.onType {
		case internal.UnsettledType_UNSETTLED_TYPE_PARCEL_SEND:
			u.owner.duey.settleSend(actx, request)
		case internal.UnsettledType_UNSETTLED_TYPE_PARCEL_CLAIM:
			u.owner.duey.settleClaim(actx, request)
		case internal.UnsettledType_UNSETTLED_TYPE_NOTE_REMOVE:
			u.owner.notes.settleRemove(actx, request)
		case internal.UnsettledType_UNSETTLED_TYPE_WEDDING_GIFT_SEND:
			u.owner.wedding.settleGiftSend(actx, request)
		case internal.UnsettledType_UNSETTLED_TYPE_WEDDING_GIFT_CLAIM:
			u.owner.wedding.settleGiftClaim(actx, request)
		case internal.UnsettledType_UNSETTLED_TYPE_WEDDING_PROPOSAL:
			u.owner.wedding.settleProposal(actx, request)
		case internal.UnsettledType_UNSETTLED_TYPE_WEDDING_RESERVE:
			u.owner.wedding.settleReservation(actx, request)
		case internal.UnsettledType_UNSETTLED_TYPE_WEDDING_WISHLIST:
			u.owner.wedding.settleWishlist(actx, request)
		case internal.UnsettledType_UNSETTLED_TYPE_WEDDING_INVITE:
			u.owner.wedding.settleInvite(actx, request)
		case internal.UnsettledType_UNSETTLED_TYPE_SHOP_OPEN:
			u.owner.storeBank.settleShopOpen(actx, request)
		case internal.UnsettledType_UNSETTLED_TYPE_TRADE_RESTORE:
			u.owner.escrow.settleRestore(actx, request)
		case internal.UnsettledType_UNSETTLED_TYPE_FAMILY_MOVE:
			u.owner.family.settleMove(actx, request)
		case internal.UnsettledType_UNSETTLED_TYPE_GUILD_CREATE:
			u.owner.gameWorld.Systems().Get[*GuildSystem]().SettleCreate(actx, u.owner, request)
		case internal.UnsettledType_UNSETTLED_TYPE_GUILD_CAPACITY:
			u.owner.gameWorld.Systems().Get[*GuildSystem]().SettleCapacity(actx, u.owner, request)
		case internal.UnsettledType_UNSETTLED_TYPE_GUILD_EMBLEM:
			u.owner.gameWorld.Systems().Get[*GuildSystem]().SettleEmblem(actx, u.owner, request)
		case internal.UnsettledType_UNSETTLED_TYPE_GUILD_GP:
			u.owner.gameWorld.Systems().Get[*GuildSystem]().SettleGP(actx, u.owner, request)
		case internal.UnsettledType_UNSETTLED_TYPE_ALLIANCE_CREATE:
			u.owner.gameWorld.Systems().Get[*GuildSystem]().SettleAllianceCreate(actx, u.owner, request)
		case internal.UnsettledType_UNSETTLED_TYPE_ALLIANCE_CAPACITY:
			u.owner.gameWorld.Systems().Get[*AllianceSystem]().SettleCapacity(actx, u.owner, request)
		default:
			log.Printf("UnsettledRequests.Settle character=%d id=%d: unknown type %s", u.owner.GetID(), request.id, request.onType)
		}
	}
}

// Timer callbacks run without an actor context; settle through the character's dispatch.
func (u *UnsettledRequests) ScheduleSettle() {
	u.owner.AddTimer(unsettledSettleTimer, u.owner.gameWorld.SettleInterval(), true, func() {
		if u.owner.gameWorld.SettlePaused() {
			return
		}
		u.owner.gameWorld.Systems().Get[*DispatchSystem]().Visit(u.owner.GetID(), func(ctx actor.Context, current *Character) {
			if current != u.owner {
				return
			}
			u.Settle(ctx)
		})
	})
}
