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

type Duey struct {
	gs *GameServer
}

func (Duey) New(gs *GameServer) *Duey {
	return &Duey{
		gs: gs,
	}
}

func (h *Duey) Handle(ctx *core.ClientContext, req *request.Duey) error {
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
	case pconst.DueyOpenFromArrival:
		ch.Duey.Open(ctx.ActorContext, true)
	case pconst.DueyIdentity:
		err = ch.Duey.ConfirmIdentity()
	case pconst.DueySend:
		err = ch.Duey.Send(ctx.ActorContext, constant.InventoryType(req.InventoryType), req.Slot, req.Count, req.Meso, req.Recipient, req.Quick, req.Message)
	case pconst.DueyReceive:
		err = ch.Duey.Receive(ctx.ActorContext, req.ParcelID)
	case pconst.DueyDelete:
		err = ch.Duey.Delete(ctx.ActorContext, req.ParcelID)
	case pconst.DueyClose:
		ch.Duey.Close()
	}

	switch {
	case errors.Is(err, entity.ErrDueyNotEnoughMeso):
		ch.Listener.OnDueyResult(ch, pconst.DueyResultNotEnoughMeso)
	case errors.Is(err, entity.ErrDueyOnlyHeld):
		ch.Listener.OnDueyResult(ch, pconst.DueyResultOnlyHeld)
	case errors.Is(err, entity.ErrInventoryFull):
		ch.Listener.OnDueyResult(ch, pconst.DueyResultInventoryFull)
	case errors.Is(err, entity.ErrDueyCannotReceive):
		ch.Listener.OnDueyResult(ch, pconst.DueyResultCannotReceive)
	case errors.Is(err, entity.ErrDueyNotSendable), errors.Is(err, entity.ErrDueyInvalid), errors.Is(err, entity.ErrDueyParcelNotFound):
		ch.Listener.OnDueyResult(ch, pconst.DueyResultInvalid)
	case err != nil:
		ch.Listener.OnUnlockAction(ch)
	}
	return nil
}
