package game

import (
	"encoding/json"
	"log"
	"time"

	"github.com/asynkron/protoactor-go/actor"
	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
)

const (
	unsettledSettleTimer    = "unsettled_settle"
	unsettledSettleInterval = 30 * time.Second
)

type UnsettledRequest struct {
	ID       uint64
	Type     internal.UnsettledType
	Target   uint64
	Detail   string
	InFlight bool
}

type UnsettledRequests struct {
	owner    *Character
	requests map[uint64]*UnsettledRequest
}

func (u *UnsettledRequests) Add(typ internal.UnsettledType, target uint64, detail any) uint64 {
	id := u.owner.GameWorld.NewUniqueID()
	u.add(id, typ, target, detail)
	return id
}

// One remote request that holds value on two characters is tracked on both with the same id.
func (u *UnsettledRequests) add(id uint64, typ internal.UnsettledType, target uint64, detail any) {
	data, err := json.Marshal(detail)
	if err != nil {
		log.Printf("UnsettledRequests.Add character=%d type=%s: %v", u.owner.GetID(), typ, err)
	}

	u.requests[id] = &UnsettledRequest{ID: id, Type: typ, Target: target, Detail: string(data), InFlight: true}
}

func (u *UnsettledRequests) Find(id uint64) *UnsettledRequest {
	return u.requests[id]
}

func (u *UnsettledRequests) FindByTarget(typ internal.UnsettledType, target uint64) *UnsettledRequest {
	for _, request := range u.requests {
		if request.Type == typ && request.Target == target {
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
		if request.InFlight {
			continue
		}

		request.InFlight = true
		switch request.Type {
		case internal.UnsettledType_UNSETTLED_TYPE_PARCEL_SEND:
			u.owner.Duey.settleSend(actx, request)
		case internal.UnsettledType_UNSETTLED_TYPE_PARCEL_CLAIM:
			u.owner.Duey.settleClaim(actx, request)
		case internal.UnsettledType_UNSETTLED_TYPE_NOTE_REMOVE:
			u.owner.Notes.settleRemove(actx, request)
		case internal.UnsettledType_UNSETTLED_TYPE_WEDDING_GIFT_SEND:
			u.owner.Wedding.settleGiftSend(actx, request)
		case internal.UnsettledType_UNSETTLED_TYPE_WEDDING_GIFT_CLAIM:
			u.owner.Wedding.settleGiftClaim(actx, request)
		case internal.UnsettledType_UNSETTLED_TYPE_WEDDING_PROPOSAL:
			u.owner.Wedding.settleProposal(actx, request)
		case internal.UnsettledType_UNSETTLED_TYPE_WEDDING_RESERVE:
			u.owner.Wedding.settleReservation(actx, request)
		case internal.UnsettledType_UNSETTLED_TYPE_WEDDING_WISHLIST:
			u.owner.Wedding.settleWishlist(actx, request)
		case internal.UnsettledType_UNSETTLED_TYPE_WEDDING_INVITE:
			u.owner.Wedding.settleInvite(actx, request)
		case internal.UnsettledType_UNSETTLED_TYPE_SHOP_OPEN:
			u.owner.StoreBank.settleShopOpen(actx, request)
		case internal.UnsettledType_UNSETTLED_TYPE_FAMILY_MOVE:
			u.owner.Family.settleMove(actx, request)
		case internal.UnsettledType_UNSETTLED_TYPE_GUILD_CREATE:
			u.owner.GameWorld.Systems().Get[*GuildSystem]().SettleCreate(actx, u.owner, request)
		case internal.UnsettledType_UNSETTLED_TYPE_GUILD_CAPACITY:
			u.owner.GameWorld.Systems().Get[*GuildSystem]().SettleCapacity(actx, u.owner, request)
		case internal.UnsettledType_UNSETTLED_TYPE_GUILD_EMBLEM:
			u.owner.GameWorld.Systems().Get[*GuildSystem]().SettleEmblem(actx, u.owner, request)
		case internal.UnsettledType_UNSETTLED_TYPE_GUILD_GP:
			u.owner.GameWorld.Systems().Get[*GuildSystem]().SettleGP(actx, u.owner, request)
		case internal.UnsettledType_UNSETTLED_TYPE_ALLIANCE_CREATE:
			u.owner.GameWorld.Systems().Get[*GuildSystem]().SettleAllianceCreate(actx, u.owner, request)
		case internal.UnsettledType_UNSETTLED_TYPE_ALLIANCE_CAPACITY:
			u.owner.GameWorld.Systems().Get[*AllianceSystem]().SettleCapacity(actx, u.owner, request)
		default:
			log.Printf("UnsettledRequests.Settle character=%d id=%d: unknown type %s", u.owner.GetID(), request.ID, request.Type)
		}
	}
}

// Timer callbacks run without an actor context; settle through the character's dispatch.
func (u *UnsettledRequests) ScheduleSettle() {
	u.owner.AddTimer(unsettledSettleTimer, unsettledSettleInterval, true, func() {
		u.owner.GameWorld.Systems().Get[*DispatchSystem]().CallCharacter(u.owner.GetID(), func(ctx actor.Context, current *Character) {
			if current != u.owner {
				return
			}
			u.Settle(ctx)
		})
	})
}
