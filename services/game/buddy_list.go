package game

import (
	"encoding/json"
	"errors"
	"log"
	"sort"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/core/async"
	pconst "github.com/boyism80/fm/protocol/constant"
	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
	"github.com/boyism80/fm/protocol/response"
)

var ErrBuddyNotEnoughMeso = errors.New("not enough meso for buddy list")

type buddyCapacityDetail struct {
	Meso int32
}

type BuddyListEntry struct {
	CharacterID uint32
	Name        string
	Group       string
	Pending     bool
	Channel     int32
}

type BuddyList struct {
	owner    *Character
	capacity uint32
	entries  map[uint32]BuddyListEntry
}

func NewBuddyList(owner *Character) *BuddyList {
	return &BuddyList{
		owner:    owner,
		capacity: uint32(pconst.DefaultBuddyCapacity),
		entries:  make(map[uint32]BuddyListEntry),
	}
}

func (bl *BuddyList) SetCapacity(n uint32) {
	if n > 255 {
		n = 255
	}
	bl.capacity = n
	bl.owner.listener.OnBuddyCapacity(bl.owner, uint8(bl.Capacity()))
}

// Expand pays the fee and raises the capacity on internal, which owns it; the promise is false when internal refuses.
func (bl *BuddyList) Expand(actx actor.Context, fee int32) (*async.Promise[bool], error) {
	if !bl.owner.inventory.RemoveMeso(fee) {
		return nil, ErrBuddyNotEnoughMeso
	}
	detail := buddyCapacityDetail{Meso: fee}
	id := bl.owner.unsettled.Add(internal.UnsettledType_UNSETTLED_TYPE_BUDDY_CAPACITY, 0, detail)

	return bl.owner.gameWorld.Systems().Get[*InternalSystem]().IncreaseBuddyCapacityAsync(actx, bl.owner, id).Bind(bl.owner).Then(func(v *internal.IncreaseBuddyCapacityReply) (bool, error) {
		bl.closeExpansion(id, detail, v)
		return v.GetOk(), nil
	}).OnError(func(err error) {
		bl.owner.unsettled.Find(id).inFlight = false

		log.Printf("BuddyList.Expand character=%d id=%d: %v", bl.owner.GetID(), id, err)
	}), nil
}

func (bl *BuddyList) closeExpansion(id uint64, detail buddyCapacityDetail, reply *internal.IncreaseBuddyCapacityReply) {
	bl.owner.unsettled.Erase(id)
	if !reply.GetOk() {
		bl.owner.inventory.AddMeso(detail.Meso)
		return
	}
	bl.SetCapacity(reply.GetCapacity())
}

// A resend under the same id answers with the capacity the first request reached, so it never raises twice.
func (bl *BuddyList) settleExpansion(actx actor.Context, request *UnsettledRequest) {
	var detail buddyCapacityDetail
	if err := json.Unmarshal([]byte(request.detail), &detail); err != nil {
		log.Printf("BuddyList.settleExpansion character=%d id=%d: %v", bl.owner.GetID(), request.id, err)
		return
	}
	bl.owner.gameWorld.Systems().Get[*InternalSystem]().IncreaseBuddyCapacityAsync(actx, bl.owner, request.id).Bind(bl.owner).Do(func(v *internal.IncreaseBuddyCapacityReply) error {
		bl.closeExpansion(request.id, detail, v)
		return nil
	}).OnError(func(err error) {
		request.inFlight = false

		log.Printf("BuddyList.settleExpansion character=%d id=%d: %v", bl.owner.GetID(), request.id, err)
	})
}

func (bl *BuddyList) Capacity() uint32 {
	if bl.capacity == 0 {
		return uint32(pconst.DefaultBuddyCapacity)
	}
	return bl.capacity
}

func (bl *BuddyList) Clear() {
	bl.entries = make(map[uint32]BuddyListEntry)
}

func (bl *BuddyList) LoadFromProto(entries []*internal.BuddyEntry, capacity uint32) {
	bl.entries = make(map[uint32]BuddyListEntry)
	if capacity > 0 {
		bl.capacity = capacity
	}
	for _, pb := range entries {
		if pb == nil {
			continue
		}
		bl.entries[pb.GetCharacterId()] = BuddyListEntryFromProto(pb)
	}
}

func (bl *BuddyList) Upsert(entry BuddyListEntry) {
	bl.entries[entry.CharacterID] = entry
}

func (bl *BuddyList) Put(entry BuddyListEntry) {
	bl.Upsert(entry)
	bl.owner.listener.OnBuddyListUpdate(bl.owner, pconst.BuddyListSyncUpdate, bl.BuildS2C())
}

func (bl *BuddyList) Remove(characterID uint32) bool {
	if _, ok := bl.entries[characterID]; !ok {
		return false
	}
	delete(bl.entries, characterID)
	return true
}

func (bl *BuddyList) SetChannel(characterID uint32, channel int32) bool {
	entry, ok := bl.entries[characterID]
	if !ok {
		return false
	}
	entry.Channel = channel
	bl.entries[characterID] = entry
	return true
}

func (bl *BuddyList) BuildS2C() []response.BuddyEntry {
	out := make([]response.BuddyEntry, 0, len(bl.entries))
	for _, entry := range bl.entries {
		out = append(out, entry.BuildS2C())
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].CharacterID < out[j].CharacterID
	})
	return out
}

func (bl *BuddyList) SyncLogin() {
	bl.owner.listener.OnBuddyListUpdate(bl.owner, pconst.BuddyListSyncLogin, bl.BuildS2C())
}

func BuddyListEntryFromProto(pb *internal.BuddyEntry) BuddyListEntry {
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

func (e BuddyListEntry) BuildS2C() response.BuddyEntry {
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
