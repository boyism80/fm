package entity

import (
	"github.com/boyism80/fm/services/game/wz"
)

type Portal struct {
	Map    *Map
	Wz     *wz.Portal
	script string
}

func (p *Portal) Script() string {
	if p == nil {
		return ""
	}
	return p.script
}

func (p *Portal) SetScript(name string) {
	if p == nil {
		return
	}
	p.script = name
}

func (p *Portal) ResetScript() {
	if p == nil || p.Wz == nil {
		return
	}
	p.script = p.Wz.ScriptName
}

func (m *Map) initPortals() {
	if m == nil {
		return
	}
	m.portals = make(map[uint8]*Portal)
	m.portalsByName = make(map[string]*Portal)
	m.doorReturnPortalIDs = nil
	if m.Wz == nil {
		return
	}
	for id, wp := range m.Wz.Portals {
		cp := wp
		cp.ID = id
		p := &Portal{
			Map:    m,
			Wz:     &cp,
			script: cp.ScriptName,
		}
		m.portals[id] = p
		if cp.Name != "" {
			m.portalsByName[cp.Name] = p
		}
	}
	slots := m.Wz.DoorReturnPortalSlots()
	if len(slots) == 0 {
		return
	}
	m.doorReturnPortalIDs = make([]uint8, 0, len(slots))
	for _, s := range slots {
		m.doorReturnPortalIDs = append(m.doorReturnPortalIDs, s.ID)
	}
}

func (m *Map) FindPortal(id uint8) *Portal {
	if m == nil || m.portals == nil {
		return nil
	}
	return m.portals[id]
}

func (m *Map) FindPortalByName(name string) *Portal {
	if m == nil || m.portalsByName == nil || name == "" {
		return nil
	}
	return m.portalsByName[name]
}

func (m *Map) resetPortalScripts() {
	if m == nil {
		return
	}
	for _, p := range m.portals {
		if p == nil {
			continue
		}
		p.ResetScript()
	}
}
