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
	if !ok {
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
		case pconst.MiniRoomTypeTrade:
			err = ch.CreateTrade()
		case pconst.MiniRoomTypeOmok, pconst.MiniRoomTypeMatchCard:
			err = ch.CreateMiniGame(req.Type, req.Title, req.Password, req.Piece)
		default:
			err = entity.ErrMiniRoomInvalid
		}
	case pconst.MiniRoomVisit:
		err = ch.VisitMiniRoom(req.SN, req.Password)
	case pconst.MiniRoomDecline:
		err = ch.DeclineTrade(req.SN, req.Reason)
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

	if req.Mode == pconst.MiniRoomChat {
		return ch.MiniRoom.Chat(ch, req.Message)
	}

	switch room := ch.MiniRoom.(type) {
	case *entity.Trade:
		switch req.Mode {
		case pconst.MiniRoomInvite:
			return room.Invite(ch, req.TargetID)
		case pconst.MiniRoomTradePutItem:
			return room.PutItem(ch, constant.InventoryType(req.InventoryType), req.Slot, req.Count, req.TradeSlot)
		case pconst.MiniRoomTradePutMeso:
			return room.PutMeso(ch, req.Meso)
		case pconst.MiniRoomTradeConfirm:
			return room.Confirm(ch)
		}
	case *entity.EntrustedShop:
		switch req.Mode {
		case pconst.MiniRoomOpen:
			return room.Open(ctx.ActorContext, ch)
		case pconst.MiniRoomAddItem:
			return room.AddItem(ctx.ActorContext, ch, constant.InventoryType(req.InventoryType), req.Slot, req.Bundles, req.PerBundle, req.Price)
		case pconst.MiniRoomBuy:
			return room.Buy(ctx.ActorContext, ch, req.Index, req.Bundles)
		case pconst.MiniRoomRemoveItem:
			return room.RemoveItem(ctx.ActorContext, ch, req.Index)
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
		case pconst.MiniRoomOpen:
			return room.Open(ctx.ActorContext, ch)
		case pconst.MiniRoomPersonalAddItem:
			return room.AddItem(ctx.ActorContext, ch, constant.InventoryType(req.InventoryType), req.Slot, req.Bundles, req.PerBundle, req.Price)
		case pconst.MiniRoomPersonalBuy:
			return room.Buy(ctx.ActorContext, ch, req.Index, req.Bundles)
		case pconst.MiniRoomPersonalRemoveItem:
			return room.RemoveItem(ctx.ActorContext, ch, req.Index)
		case pconst.MiniRoomKick:
			return room.Kick(ch, uint8(req.Slot), req.Name, pconst.MiniRoomLeaveKicked)
		case pconst.MiniRoomKickTimeout:
			return room.Kick(ch, uint8(req.Slot), req.Name, pconst.MiniRoomLeaveStayTimeout)
		case pconst.MiniRoomBlacklist:
			return room.Ban(ch, req.Names)
		}
	case *entity.MiniGame:
		switch req.Mode {
		case pconst.MiniRoomReady:
			return room.Ready(ch, true)
		case pconst.MiniRoomUnready:
			return room.Ready(ch, false)
		case pconst.MiniRoomExpel:
			return room.Expel(ch)
		case pconst.MiniRoomStart:
			return room.Start(ch)
		case pconst.MiniRoomMoveOmok:
			return room.MoveOmok(ch, req.X, req.Y, req.Stone)
		case pconst.MiniRoomSelectCard:
			return room.SelectCard(ch, req.FirstPick, req.Card)
		case pconst.MiniRoomSkip:
			return room.Skip(ch)
		case pconst.MiniRoomGiveUp:
			return room.GiveUp(ch)
		case pconst.MiniRoomRequestTie:
			return room.RequestTie(ch)
		case pconst.MiniRoomAnswerTie:
			return room.AnswerTie(ch, req.Accept)
		case pconst.MiniRoomExitAfterGame:
			return room.ExitAfterGame(ch, true)
		case pconst.MiniRoomCancelExit:
			return room.ExitAfterGame(ch, false)
		}
	}
	return nil
}
