package actor

import (
	"github.com/boyism80/fm/maple/wz"
)

const newCharacterMap = 10000

// Atlas knows which maps a character walks to through plain portals, from a
// return scroll town or from the map it stands on. Any other map is reached
// through an NPC or a portal script, which a suite has to take itself.
type Atlas struct {
	wz   *wz.Resources
	open map[uint32]bool
}

func NewAtlas(resources *wz.Resources) *Atlas {
	starts := []uint32{newCharacterMap}
	for _, item := range resources.Items {
		consume, ok := item.(*wz.Consume)
		if !ok || consume.MoveTo <= 0 || consume.MoveTo == wz.ConsumeMoveToReturnMap {
			continue
		}
		starts = append(starts, uint32(consume.MoveTo))
	}

	at := &Atlas{wz: resources}
	at.open = at.walk(starts...)
	return at
}

func (at *Atlas) walk(from ...uint32) map[uint32]bool {
	seen := make(map[uint32]bool, len(from))
	queue := make([]uint32, 0, len(from))
	for _, id := range from {
		if _, ok := at.wz.Maps[id]; ok && !seen[id] {
			seen[id] = true
			queue = append(queue, id)
		}
	}
	for len(queue) > 0 {
		m := at.wz.Maps[queue[0]]
		queue = queue[1:]
		for _, portal := range m.Portals {
			if portal.ScriptName != "" || portal.TargetMapId < 0 {
				continue
			}
			to := uint32(portal.TargetMapId)
			if _, ok := at.wz.Maps[to]; !ok || seen[to] {
				continue
			}
			seen[to] = true
			queue = append(queue, to)
		}
	}
	return seen
}

func (at *Atlas) Walkable(from, to uint32) bool {
	if at.open[to] {
		return true
	}
	return at.walk(from)[to]
}
