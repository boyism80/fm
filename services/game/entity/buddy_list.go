package entity

import (
	"sort"

	pconst "github.com/boyism80/fm/protocol/constant"
	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
	"github.com/boyism80/fm/protocol/response"
)

type BuddyListEntry struct {
	CharacterID uint32
	Name        string
	Group       string
	Pending     bool
	Channel     int32
}

type BuddyList struct {
	capacity uint32
	entries  map[uint32]BuddyListEntry
}

func NewBuddyList() *BuddyList {
	return &BuddyList{
		capacity: uint32(pconst.DefaultBuddyCapacity),
		entries:  make(map[uint32]BuddyListEntry),
	}
}

func (bl *BuddyList) Capacity() uint32 {
	if bl == nil {
		return uint32(pconst.DefaultBuddyCapacity)
	}
	if bl.capacity == 0 {
		return uint32(pconst.DefaultBuddyCapacity)
	}
	return bl.capacity
}

func (bl *BuddyList) Clear() {
	if bl == nil {
		return
	}
	bl.entries = make(map[uint32]BuddyListEntry)
}

func (bl *BuddyList) LoadFromProto(entries []*internal.BuddyEntry, capacity uint32) {
	if bl == nil {
		return
	}
	bl.entries = make(map[uint32]BuddyListEntry)
	if capacity > 0 {
		bl.capacity = capacity
	}
	for _, pb := range entries {
		if pb == nil || pb.GetCharacterId() == 0 {
			continue
		}
		bl.entries[pb.GetCharacterId()] = buddyEntryFromProto(pb)
	}
}

func (bl *BuddyList) Upsert(entry BuddyListEntry) {
	if bl == nil || entry.CharacterID == 0 {
		return
	}
	bl.entries[entry.CharacterID] = entry
}

func (bl *BuddyList) Remove(characterID uint32) bool {
	if bl == nil || characterID == 0 {
		return false
	}
	if _, ok := bl.entries[characterID]; !ok {
		return false
	}
	delete(bl.entries, characterID)
	return true
}

func (bl *BuddyList) SetChannel(characterID uint32, channel int32) bool {
	if bl == nil || characterID == 0 {
		return false
	}
	entry, ok := bl.entries[characterID]
	if !ok {
		return false
	}
	entry.Channel = channel
	bl.entries[characterID] = entry
	return true
}

func (bl *BuddyList) SnapshotForClient() []response.BuddyEntry {
	if bl == nil {
		return nil
	}
	out := make([]response.BuddyEntry, 0, len(bl.entries))
	for _, entry := range bl.entries {
		out = append(out, entry.ToResponse())
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].CharacterID < out[j].CharacterID
	})
	return out
}

func buddyEntryFromProto(pb *internal.BuddyEntry) BuddyListEntry {
	ch := pb.GetChannelIndex()
	if ch < 0 {
		ch = -1
	}
	return BuddyListEntry{
		CharacterID: pb.GetCharacterId(),
		Name:        pb.GetName(),
		Group:       pb.GetGroupName(),
		Pending:     pb.GetPending(),
		Channel:     ch,
	}
}

func BuddyListEntryFromProto(pb *internal.BuddyEntry) BuddyListEntry {
	return buddyEntryFromProto(pb)
}

func (e BuddyListEntry) ToResponse() response.BuddyEntry {
	ch := e.Channel
	if ch < 0 {
		ch = -1
	}
	return response.BuddyEntry{
		CharacterID: e.CharacterID,
		Name:        e.Name,
		Pending:     e.Pending,
		Channel:     ch,
		Group:       e.Group,
	}
}
