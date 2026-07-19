package actor

import (
	"github.com/asynkron/protoactor-go/actor"
	pconst "github.com/boyism80/fm/protocol/constant"
	"github.com/boyism80/fm/services/game/entity"
)

type DeliverBuddyListUpdateHandler struct{}

func (DeliverBuddyListUpdateHandler) New() *DeliverBuddyListUpdateHandler {
	return &DeliverBuddyListUpdateHandler{}
}

func (h *DeliverBuddyListUpdateHandler) Handle(ctx actor.Context, a *MapActor, msg *DeliverBuddyListUpdate) {
	if a.Map == nil || msg == nil {
		return
	}
	ch := a.Map.GetPlayer(msg.RecipientCharacterID)
	if ch == nil {
		return
	}
	action := pconst.BuddyListSyncAction(msg.SyncAction)
	bl := ch.BuddyList()
	if action == pconst.BuddyListSyncDelete {
		bl.Clear()
	}
	for _, entry := range msg.Entries {
		if entry.CharacterID == 0 {
			continue
		}
		channel := entry.Channel
		if channel < 0 {
			channel = -1
		}
		bl.Upsert(entity.BuddyListEntry{
			CharacterID: entry.CharacterID,
			Name:        entry.Name,
			Group:       entry.Group,
			Pending:     entry.Pending,
			Channel:     channel,
		})
	}
	ch.Listener.OnBuddyListUpdate(ch, action, msg.Entries)
}
