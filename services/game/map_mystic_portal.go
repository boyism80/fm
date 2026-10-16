package game

import (
	"github.com/boyism80/fm/types"
)

func (m *Map) FindMysticReturnPortal(partySlot int) (portalID uint8, pos types.Vector2[int16], ok bool) {
	portalID, ok = m.wz.DoorReturnPortalSpawnID(partySlot)
	if !ok {
		return 0, types.Vector2[int16]{}, false
	}

	if sp, found := m.wz.GetSpawnPosition(portalID); found {
		pos.X, pos.Y = sp.X, sp.Y
		return portalID, pos, true
	}
	if p := m.FindPortal(portalID); p != nil && p.wz != nil {
		pos.X, pos.Y = p.wz.Position.X, p.wz.Position.Y
	}
	return portalID, pos, true
}
