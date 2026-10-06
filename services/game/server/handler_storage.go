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
	if ok == false {
		return fmt.Errorf("client is not a GameClient")
	}

	ch := client.GetCharacter()
	if ch == nil {
		return nil
	}
	if ch.Storage == nil {
		return nil
	}

	var err error
	switch req.Mode {
	case pconst.StorageTakeOut:
		err = ch.Storage.TakeOut(constant.InventoryType(req.InventoryType), req.Index)
	case pconst.StorageStore:
		err = ch.Storage.Store(req.Slot, req.ItemID, req.Count)
	case pconst.StorageArrange:
		err = ch.Storage.Arrange()
	case pconst.StorageMeso:
		if req.Meso > 0 {
			err = ch.Storage.WithdrawMeso(req.Meso)
		} else {
			err = ch.Storage.DepositMeso(-req.Meso)
		}
	case pconst.StorageClose:
		ch.Storage.Close()
	}

	switch {
	case errors.Is(err, entity.ErrStorageFull):
		ch.Listener.OnStorageError(ch, pconst.StorageResultFull)
	case errors.Is(err, entity.ErrStorageNotEnoughMeso):
		ch.Listener.OnStorageError(ch, pconst.StorageResultNotEnoughMeso)
	case errors.Is(err, entity.ErrInventoryFull):
		ch.Listener.OnStorageError(ch, pconst.StorageResultInventoryFull)
	case err != nil:
		ch.Listener.OnUnlockAction(ch)
	}
	return nil
}
