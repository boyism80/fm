package entity

import (
	"github.com/boyism80/fm/types"
)

func (m *Map) FindMysticReturnPortal(partySlot int) (portalID uint8, pos types.Vector2[int16], ok bool) {
	portalID, ok = m.Wz.DoorReturnPortalSpawnID(partySlot)
	if !ok {
		return 0, types.Vector2[int16]{}, false
	}

	if sp, found := m.Wz.GetSpawnPosition(portalID); found {
		pos.X, pos.Y = sp.X, sp.Y
		return portalID, pos, true
	}
	if p := m.FindPortal(portalID); p != nil && p.Wz != nil {
		pos.X, pos.Y = p.Wz.Position.X, p.Wz.Position.Y
	}
	return portalID, pos, true
}
