package entity

import (
	"errors"

	"github.com/boyism80/fm/services/game/constant"
)

var ErrChairNotOwned = errors.New("chair item not owned")

// SitOnChair sits on a chair item from the installation tab. Chair holds the item ID only.
func (ch *Character) SitOnChair(itemID uint32) error {
	tab := ch.Inventory.Tabs[constant.InventoryTypeInstallation]
	if tab == nil || tab.Find(itemID) == nil {
		return ErrChairNotOwned
	}

	ch.Chair = itemID
	ch.Listener.OnChairChanged(ch, itemID)
	return nil
}

// SitOnMapSeat sits on a seat placed in the map. A map seat is not a chair item, so Chair is cleared.
func (ch *Character) SitOnMapSeat(seatID int16) {
	ch.Chair = 0
	ch.Listener.OnMapSeatChanged(ch, seatID)
}

func (ch *Character) StandUp() {
	ch.Chair = 0
	ch.Listener.OnMapSeatChanged(ch, -1)
	ch.Listener.OnChairChanged(ch, 0)
}
