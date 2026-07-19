package entity

import (
	"github.com/boyism80/fm/services/game/constant"
	"github.com/boyism80/fm/types"
)

func (m *Map) syncVisibilityAroundMove(obj Object, before types.Vector2[int16]) {
	if m == nil || obj == nil {
		return
	}
	after := obj.GetPosition()
	switch mover := obj.(type) {
	case *Character:
		beforeNear := m.objectSetNear(before, constant.ObjectTypeObject)
		afterNear := m.GetObjectsNear(after, constant.ObjectTypeObject)
		afterSet := make(map[Object]struct{}, len(afterNear))
		for _, o := range afterNear {
			if o == nil || o == mover {
				continue
			}
			afterSet[o] = struct{}{}
			if _, seen := beforeNear[o]; seen {
				continue
			}
			o.SendSpawnSyncToViewer(mover)
			if peer, ok := o.(*Character); ok {
				mover.SendSpawnSyncToViewer(peer)
				syncPartyMemberHP(mover, peer)
			}
		}
		for o := range beforeNear {
			if o == nil || o == mover {
				continue
			}
			if _, still := afterSet[o]; still {
				continue
			}
			o.SendDestroySyncToViewer(mover)
			if peer, ok := o.(*Character); ok {
				mover.SendDestroySyncToViewer(peer)
			}
		}
	default:
		beforeViewers := m.characterSetNear(before)
		afterViewers := m.GetObjectsNear(after, constant.ObjectTypeCharacter)
		afterSet := make(map[*Character]struct{}, len(afterViewers))
		for _, candidate := range afterViewers {
			viewer, ok := candidate.(*Character)
			if !ok || viewer == nil {
				continue
			}
			afterSet[viewer] = struct{}{}
			if _, seen := beforeViewers[viewer]; seen {
				continue
			}
			obj.SendSpawnSyncToViewer(viewer)
		}
		for viewer := range beforeViewers {
			if _, still := afterSet[viewer]; still {
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

func (m *Map) objectSetNear(position types.Vector2[int16], filter constant.ObjectType) map[Object]struct{} {
	objects := m.GetObjectsNear(position, filter)
	out := make(map[Object]struct{}, len(objects))
	for _, o := range objects {
		if o != nil {
			out[o] = struct{}{}
		}
	}
	return out
}

func (m *Map) characterSetNear(position types.Vector2[int16]) map[*Character]struct{} {
	objects := m.GetObjectsNear(position, constant.ObjectTypeCharacter)
	out := make(map[*Character]struct{}, len(objects))
	for _, candidate := range objects {
		if viewer, ok := candidate.(*Character); ok && viewer != nil {
			out[viewer] = struct{}{}
		}
	}
	return out
}
