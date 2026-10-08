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
		if req.Type != pconst.MiniRoomTypeHiredMerchant {
			ch.Listener.OnUnlockAction(ch)
			return nil
		}
		err = ch.CreateHiredMerchant(ctx.ActorContext, req.Title, req.Slot, req.ItemID)
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
		return ch.MiniRoom.Open(ch)
	case pconst.MiniRoomAddItem:
		return ch.MiniRoom.AddItem(ctx.ActorContext, ch, constant.InventoryType(req.InventoryType), req.Slot, req.Bundles, req.PerBundle, req.Price)
	case pconst.MiniRoomBuy:
		return ch.MiniRoom.Buy(ctx.ActorContext, ch, req.Index, req.Bundles)
	case pconst.MiniRoomRemoveItem:
		return ch.MiniRoom.RemoveItem(ctx.ActorContext, ch, req.Index)
	case pconst.MiniRoomMaintenanceOff:
		return ch.MiniRoom.EndMaintenance(ch)
	case pconst.MiniRoomArrange:
		return ch.MiniRoom.Arrange(ctx.ActorContext, ch)
	case pconst.MiniRoomClose:
		return ch.MiniRoom.Close(ctx.ActorContext, ch)
	case pconst.MiniRoomWithdrawMeso:
		return ch.MiniRoom.WithdrawMeso(ctx.ActorContext, ch)
	}
	return nil
}
