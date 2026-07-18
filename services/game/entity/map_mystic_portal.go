package entity

import (
	"github.com/boyism80/fm/types"
)

func (m *Map) TryAcquireMysticReturnPortal(preferredSlot int) (portalID uint8, pos types.Vector2[int16], ok bool) {
	if m == nil || m.Wz == nil {
		return 0, types.Vector2[int16]{}, false
	}
	slots := m.doorReturnPortalIDs
	if len(slots) == 0 {
		return 0, types.Vector2[int16]{}, false
	}
	if m.UsedDoorPortalIDs == nil {
		m.UsedDoorPortalIDs = make(map[uint8]struct{})
	}
	idx := preferredSlot
	if idx < 0 {
		idx = 0
	}
	if idx >= len(slots) {
		idx = len(slots) - 1
	}
	prefID := slots[idx]
	if _, used := m.UsedDoorPortalIDs[prefID]; !used {
		m.UsedDoorPortalIDs[prefID] = struct{}{}
		return prefID, m.mysticPortalPosition(prefID), true
	}
	for _, id := range slots {
		if _, used := m.UsedDoorPortalIDs[id]; used {
			continue
		}
		m.UsedDoorPortalIDs[id] = struct{}{}
		return id, m.mysticPortalPosition(id), true
	}
	return 0, types.Vector2[int16]{}, false
}

func (m *Map) mysticPortalPosition(portalID uint8) types.Vector2[int16] {
	var v types.Vector2[int16]
	if sp, ok := m.Wz.GetSpawnPosition(portalID); ok {
		v.X, v.Y = sp.X, sp.Y
		return v
	}
	if p := m.FindPortal(portalID); p != nil && p.Wz != nil {
		v.X, v.Y = p.Wz.Position.X, p.Wz.Position.Y
	}
	return v
}

func (m *Map) ReleaseMysticReturnPortal(portalID uint8) {
	if m == nil || m.UsedDoorPortalIDs == nil {
		return
	}
	delete(m.UsedDoorPortalIDs, portalID)
}
