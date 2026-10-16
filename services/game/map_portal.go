package game

import (
	"github.com/boyism80/fm/maple/wz"
)

type Portal struct {
	onMap  *Map
	wz     *wz.Portal
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
	if p == nil || p.wz == nil {
		return
	}
	p.script = p.wz.ScriptName
}

func (m *Map) initPortals() {
	if m == nil {
		return
	}
	m.portals = make(map[uint8]*Portal)
	m.portalsByName = make(map[string]*Portal)
	if m.wz == nil {
		return
	}
	for id, wp := range m.wz.Portals {
		cp := wp
		cp.ID = id
		p := &Portal{
			onMap:  m,
			wz:     &cp,
			script: cp.ScriptName,
		}
		m.portals[id] = p
		if cp.Name != "" {
			m.portalsByName[cp.Name] = p
		}
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
