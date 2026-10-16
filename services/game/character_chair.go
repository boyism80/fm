package game

import (
	"errors"

	"github.com/boyism80/fm/maple/constant"
)

var ErrChairNotOwned = errors.New("chair item not owned")

func (ch *Character) SitOnChair(itemID uint32) error {
	tab := ch.inventory.tabs[constant.InventoryTypeInstallation]
	if tab == nil || tab.Find(itemID) == nil {
		return ErrChairNotOwned
	}

	ch.chair = itemID
	ch.listener.OnSitOnChair(ch, itemID)
	return nil
}

func (ch *Character) SitOnMapSeat(seatID int16) {
	ch.chair = 0
	ch.listener.OnSitOnMapSeat(ch, seatID)
}

func (ch *Character) StandUp() {
	ch.chair = 0
	ch.listener.OnStandUp(ch)
}
