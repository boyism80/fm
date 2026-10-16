package game

import (
	"slices"

	"github.com/boyism80/fm/maple/constant"
	"github.com/boyism80/fm/types"
)

func (m *Map) updateVisibility(obj Object, before types.Vector2[int16]) {
	if m == nil || obj == nil {
		return
	}
	after := obj.GetPosition()
	switch mover := obj.(type) {
	case *Character:
		includeHidden := &SearchOption{IncludeHidden: true}
		beforeNear := m.GetObjectsNear(before, constant.ObjectTypeObject, includeHidden)
		afterNear := m.GetObjectsNear(after, constant.ObjectTypeObject, includeHidden)
		for _, o := range afterNear {
			if o == nil || o == mover || slices.Contains(beforeNear, o) {
				continue
			}
			o.SendSpawnSyncToViewer(mover)
			if peer, ok := o.(*Character); ok {
				mover.SendSpawnSyncToViewer(peer)
				mover.party.syncHP(peer)
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
		includeHidden := &SearchOption{IncludeHidden: true}
		beforeViewers := m.GetObjectsNear(before, constant.ObjectTypeCharacter, includeHidden)
		afterViewers := m.GetObjectsNear(after, constant.ObjectTypeCharacter, includeHidden)
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

func (m *Map) BroadcastPartyHP(ch *Character) {
	for _, peer := range m.partyPeers(ch) {
		peer.listener.OnPartyMemberHP(peer, ch)
	}
}

func (m *Map) SyncPartyHP(ch *Character) {
	for _, peer := range m.partyPeers(ch) {
		peer.listener.OnPartyMemberHP(peer, ch)
		ch.listener.OnPartyMemberHP(ch, peer)
	}
}

func (m *Map) partyPeers(ch *Character) []*Character {
	partyID := ch.party.ID()
	if partyID == nil {
		return nil
	}
	var peers []*Character
	for _, obj := range m.GetObjects(constant.ObjectTypeCharacter) {
		peer, ok := obj.(*Character)
		if !ok || peer == ch || !peer.party.Is(partyID) {
			continue
		}
		peers = append(peers, peer)
	}
	return peers
}
