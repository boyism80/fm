package entity

import (
	"slices"

	"github.com/boyism80/fm/services/game/constant"
	"github.com/boyism80/fm/types"
)

func (m *Map) updateVisibility(obj Object, before types.Vector2[int16]) {
	if m == nil || obj == nil {
		return
	}
	after := obj.GetPosition()
	switch mover := obj.(type) {
	case *Character:
		beforeNear := m.GetObjectsNear(before, constant.ObjectTypeObject, nil)
		afterNear := m.GetObjectsNear(after, constant.ObjectTypeObject, nil)
		for _, o := range afterNear {
			if o == nil || o == mover || slices.Contains(beforeNear, o) {
				continue
			}
			o.SendSpawnSyncToViewer(mover)
			if peer, ok := o.(*Character); ok {
				mover.SendSpawnSyncToViewer(peer)
				syncPartyMemberHP(mover, peer)
			}
		}
		for _, o := range beforeNear {
			if o == nil || o == mover || slices.Contains(afterNear, o) {
				continue
			}
			o.SendDestroySyncToViewer(mover)
			if peer, ok := o.(*Character); ok {
				mover.SendDestroySyncToViewer(peer)
			}
		}
	default:
		beforeViewers := m.GetObjectsNear(before, constant.ObjectTypeCharacter, nil)
		afterViewers := m.GetObjectsNear(after, constant.ObjectTypeCharacter, nil)
		for _, candidate := range afterViewers {
			viewer, ok := candidate.(*Character)
			if !ok || viewer == nil || slices.Contains(beforeViewers, candidate) {
				continue
			}
			obj.SendSpawnSyncToViewer(viewer)
		}
		for _, candidate := range beforeViewers {
			viewer, ok := candidate.(*Character)
			if !ok || viewer == nil || slices.Contains(afterViewers, candidate) {
				continue
			}
			obj.SendDestroySyncToViewer(viewer)
		}
	}
}

func syncPartyMemberHP(a, b *Character) {
	if a == nil || b == nil {
		return
	}
	pa := a.GetPartyID()
	pb := b.GetPartyID()
	if pa == nil || pb == nil || *pa != *pb {
		return
	}
	if a.Listener != nil {
		a.Listener.OnPartyMemberHPChanged(a, b)
	}
	if b.Listener != nil {
		b.Listener.OnPartyMemberHPChanged(b, a)
	}
}
