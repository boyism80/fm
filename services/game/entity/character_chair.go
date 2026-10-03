package entity

import (
	"errors"

	"github.com/boyism80/fm/services/game/constant"
)

var ErrChairNotOwned = errors.New("chair item not owned")

func (ch *Character) SitOnChair(itemID uint32) error {
	tab := ch.Inventory.Tabs[constant.InventoryTypeInstallation]
	if tab == nil || tab.Find(itemID) == nil {
		return ErrChairNotOwned
	}

	ch.Chair = itemID
	ch.Listener.OnSitOnChair(ch, itemID)
	return nil
}

func (ch *Character) SitOnMapSeat(seatID int16) {
	ch.Chair = 0
	ch.Listener.OnSitOnMapSeat(ch, seatID)
}

func (ch *Character) StandUp() {
	ch.Chair = 0
	ch.Listener.OnStandUp(ch)
}
