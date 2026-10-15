package actor

import (
	"github.com/asynkron/protoactor-go/actor"
	pconst "github.com/boyism80/fm/protocol/constant"
	"github.com/boyism80/fm/services/game"
)

type SyncBuddyListHandler struct{}

func (SyncBuddyListHandler) New() *SyncBuddyListHandler {
	return &SyncBuddyListHandler{}
}

func (h *SyncBuddyListHandler) Handle(ctx actor.Context, a *GameLogicActor, msg *game.SyncBuddyList) {
	if msg == nil {
		return
	}
	m := a.GetCharacter(msg.RecipientCharacterID)
	if m == nil {
		return
	}
	ch := m.GetPlayer(msg.RecipientCharacterID)
	if ch == nil {
		return
	}
	action := pconst.BuddyListSyncAction(msg.SyncAction)
	if action == pconst.BuddyListSyncDelete {
		ch.Buddies.Clear()
	}
	for _, entry := range msg.Entries {
		if entry.CharacterID == 0 {
			continue
		}
		channel := entry.Channel
		if channel < 0 {
			channel = -1
		}
		ch.Buddies.Upsert(game.BuddyListEntry{
			CharacterID: entry.CharacterID,
			Name:        entry.Name,
			Group:       entry.Group,
			Pending:     entry.Pending,
			Channel:     channel,
		})
	}
	ch.Listener.OnBuddyListUpdate(ch, action, msg.Entries)
}
