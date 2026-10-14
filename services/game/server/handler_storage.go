package server

import (
	"errors"
	"fmt"

	"github.com/boyism80/fm/core"
	pconst "github.com/boyism80/fm/protocol/constant"
	"github.com/boyism80/fm/protocol/request"
	"github.com/boyism80/fm/services/game/client"
	"github.com/boyism80/fm/services/game/constant"
	"github.com/boyism80/fm/services/game/domain"
)

type Storage struct {
	gs *GameServer
}

func (Storage) New(gs *GameServer) *Storage {
	return &Storage{
		gs: gs,
	}
}

func (h *Storage) Handle(ctx *core.ClientContext, req *request.Storage) error {
	client, ok := ctx.Client.(*client.GameClient)
	if !ok {
		return fmt.Errorf("client is not a GameClient")
	}

	ch := client.GetCharacter()
	if ch == nil {
		return nil
	}
	if ch.UsingCashItem() {
		ch.Listener.OnUpdateStats(ch, nil, true)
		return nil
	}

	var err error
	switch req.Mode {
	case pconst.StorageTakeOut:
		err = ch.Account.Storage.TakeOut(constant.InventoryType(req.InventoryType), req.Index)
	case pconst.StorageStore:
		err = ch.Account.Storage.Store(req.Slot, req.ItemID, req.Count)
	case pconst.StorageArrange:
		err = ch.Account.Storage.Arrange()
	case pconst.StorageMeso:
		if req.Meso > 0 {
			err = ch.Account.Storage.WithdrawMeso(req.Meso)
		} else {
			err = ch.Account.Storage.DepositMeso(-req.Meso)
		}
	case pconst.StorageClose:
		ch.Account.Storage.Close()
	}

	switch {
	case errors.Is(err, domain.ErrStorageFull):
		ch.Listener.OnStorageError(ch, pconst.StorageResultFull)
	case errors.Is(err, domain.ErrStorageNotEnoughMeso):
		ch.Listener.OnStorageError(ch, pconst.StorageResultNotEnoughMeso)
	case errors.Is(err, domain.ErrInventoryFull):
		ch.Listener.OnStorageError(ch, pconst.StorageResultInventoryFull)
	case err != nil:
		ch.Listener.OnUnlockAction(ch)
	}
	return nil
}
