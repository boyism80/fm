package entity

import (
	"github.com/boyism80/fm/types"
)

func (m *Map) TryAcquireMysticReturnPortal(preferredSlot int) (portalID uint8, pos types.Vector2[int16], ok bool) {
	if m == nil || m.Wz == nil {
		return 0, types.Vector2[int16]{}, false
	}
	slots := m.Wz.DoorReturnPortalSlots()
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
	pref := slots[idx]
	if _, used := m.UsedDoorPortalIDs[pref.ID]; !used {
		m.UsedDoorPortalIDs[pref.ID] = struct{}{}
		var v types.Vector2[int16]
		if sp, ok2 := m.Wz.GetSpawnPosition(pref.ID); ok2 {
			v.X, v.Y = sp.X, sp.Y
		} else if pt, ok2 := m.Wz.Portals[pref.ID]; ok2 {
			v.X, v.Y = pt.Position.X, pt.Position.Y
		}
		return pref.ID, v, true
	}
	for _, p := range slots {
		if _, used := m.UsedDoorPortalIDs[p.ID]; used {
			continue
		}
		m.UsedDoorPortalIDs[p.ID] = struct{}{}
		var v types.Vector2[int16]
		if sp, ok2 := m.Wz.GetSpawnPosition(p.ID); ok2 {
			v.X, v.Y = sp.X, sp.Y
		} else if pt, ok2 := m.Wz.Portals[p.ID]; ok2 {
			v.X, v.Y = pt.Position.X, pt.Position.Y
		}
		return p.ID, v, true
	}
	return 0, types.Vector2[int16]{}, false
}

func (m *Map) ReleaseMysticReturnPortal(portalID uint8) {
	if m == nil || m.UsedDoorPortalIDs == nil {
		return
	}
	delete(m.UsedDoorPortalIDs, portalID)
}
