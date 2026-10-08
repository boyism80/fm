package server

import (
	"errors"
	"fmt"

	"github.com/boyism80/fm/core"
	pconst "github.com/boyism80/fm/protocol/constant"
	"github.com/boyism80/fm/protocol/request"
	"github.com/boyism80/fm/services/game/client"
	"github.com/boyism80/fm/services/game/constant"
	"github.com/boyism80/fm/services/game/entity"
)

type MiniRoom struct {
	gs *GameServer
}

func (MiniRoom) New(gs *GameServer) *MiniRoom {
	return &MiniRoom{
		gs: gs,
	}
}

func (h *MiniRoom) Handle(ctx *core.ClientContext, req *request.MiniRoom) error {
	client, ok := ctx.Client.(*client.GameClient)
	if ok == false {
		return fmt.Errorf("client is not a GameClient")
	}

	ch := client.GetCharacter()
	if ch == nil {
		return nil
	}

	var err error
	switch req.Mode {
	case pconst.MiniRoomCreate:
		switch req.Type {
		case pconst.MiniRoomTypeEntrustedShop:
			err = ch.CreateEntrustedShop(ctx.ActorContext, req.Title, req.Slot, req.ItemID)
		case pconst.MiniRoomTypePersonalShop:
			err = ch.CreatePersonalShop(ctx.ActorContext, req.Title, req.Slot, req.ItemID)
		default:
			err = entity.ErrMiniRoomInvalid
		}
	case pconst.MiniRoomVisit:
		err = ch.VisitMiniRoom(req.SN)
	case pconst.MiniRoomExit:
		if ch.MiniRoom != nil {
			ch.MiniRoom.Leave(ch)
		}
	default:
		err = h.handleRoom(ctx, ch, req)
	}

	var enterErr *entity.MiniRoomEnterError
	var buyErr *entity.MiniRoomBuyError
	switch {
	case errors.As(err, &enterErr):
		ch.Listener.OnMiniRoomEnterFailed(ch, enterErr.Code)
	case errors.As(err, &buyErr):
		ch.Listener.OnMiniRoomBuyFailed(ch, buyErr.Result)
	case err != nil:
		ch.Listener.OnUnlockAction(ch)
	}
	return nil
}

func (h *MiniRoom) handleRoom(ctx *core.ClientContext, ch *entity.Character, req *request.MiniRoom) error {
	if ch.MiniRoom == nil {
		return entity.ErrMiniRoomInvalid
	}

	switch req.Mode {
	case pconst.MiniRoomChat:
		return ch.MiniRoom.Chat(ch, req.Message)
	case pconst.MiniRoomOpen:
		return ch.MiniRoom.Open(ctx.ActorContext, ch)
	case pconst.MiniRoomAddItem, pconst.MiniRoomPersonalAddItem:
		return ch.MiniRoom.AddItem(ctx.ActorContext, ch, constant.InventoryType(req.InventoryType), req.Slot, req.Bundles, req.PerBundle, req.Price)
	case pconst.MiniRoomBuy, pconst.MiniRoomPersonalBuy:
		return ch.MiniRoom.Buy(ctx.ActorContext, ch, req.Index, req.Bundles)
	case pconst.MiniRoomRemoveItem, pconst.MiniRoomPersonalRemoveItem:
		return ch.MiniRoom.RemoveItem(ctx.ActorContext, ch, req.Index)
	}

	switch room := ch.MiniRoom.(type) {
	case *entity.EntrustedShop:
		switch req.Mode {
		case pconst.MiniRoomMaintenanceOff:
			return room.EndMaintenance(ch)
		case pconst.MiniRoomArrange:
			return room.Arrange(ctx.ActorContext, ch)
		case pconst.MiniRoomClose:
			return room.Close(ctx.ActorContext, ch)
		case pconst.MiniRoomWithdrawMeso:
			return room.WithdrawMeso(ctx.ActorContext, ch)
		}
	case *entity.PersonalShop:
		switch req.Mode {
		case pconst.MiniRoomKick:
			return room.Kick(ch, uint8(req.Slot), req.Name, pconst.MiniRoomLeaveKicked)
		case pconst.MiniRoomKickTimeout:
			return room.Kick(ch, uint8(req.Slot), req.Name, pconst.MiniRoomLeaveStayTimeout)
		case pconst.MiniRoomBlacklist:
			return room.Ban(ch, req.Names)
		}
	}
	return nil
}
