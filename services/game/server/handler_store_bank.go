package server

import (
	"fmt"

	"github.com/boyism80/fm/core"
	pconst "github.com/boyism80/fm/protocol/constant"
	"github.com/boyism80/fm/protocol/request"
	"github.com/boyism80/fm/services/game/client"
)

type StoreBank struct {
	gs *GameServer
}

func (StoreBank) New(gs *GameServer) *StoreBank {
	return &StoreBank{
		gs: gs,
	}
}

func (h *StoreBank) Handle(ctx *core.ClientContext, req *request.StoreBank) error {
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
	case pconst.StoreBankWithdraw:
		err = ch.StoreBank.Withdraw()
	case pconst.StoreBankConfirm:
		err = ch.StoreBank.Confirm(ctx.ActorContext)
	case pconst.StoreBankClose:
		ch.StoreBank.Close()
	}
	if err != nil {
		ch.Listener.OnUnlockAction(ch)
	}
	return nil
}
