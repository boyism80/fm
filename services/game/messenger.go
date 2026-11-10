package game

import (
	"errors"

	"github.com/asynkron/protoactor-go/actor"
	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
)

var (
	ErrMessengerBusy      = errors.New("messenger request is pending")
	ErrMessengerNotInRoom = errors.New("not in a messenger room")
)

const messengerQueueLimit = 32

type MessengerEvent struct {
	Action   internal.MessengerAction
	RoomID   uint32
	Position uint8
	Name     string
	Text     string
	Channel  uint16
	Look     []byte
}

type Messenger struct {
	owner   *Character
	roomID  uint32
	pending bool
	joining uint32
	queued  []MessengerEvent
	// lookChanged holds a look change made while a join waits, shared once the owner is in the room.
	lookChanged bool
}

// RequestAsync asks internal for one messenger action; one request is in flight at a time, and the room follows the
// reply.
func (m *Messenger) RequestAsync(actx actor.Context, action internal.MessengerAction, roomID uint32, text string) error {
	if m.pending {
		return ErrMessengerBusy
	}
	if action != internal.MessengerAction_MESSENGER_ACTION_ENTER && action != internal.MessengerAction_MESSENGER_ACTION_DECLINE && m.roomID == 0 {
		return ErrMessengerNotInRoom
	}

	m.pending = true
	if action == internal.MessengerAction_MESSENGER_ACTION_ENTER {
		m.joining = roomID
	}
	m.owner.gameWorld.Systems().Get[*InternalSystem]().MessengerAsync(actx, m.owner, action, roomID, text).Do(func(reply *internal.MessengerReply) error {
		if m.owner.LoggedOut() {
			return nil
		}
		if !reply.Ok {
			m.fail(action, text)
			return nil
		}
		m.finish(action, text, reply)
		return nil
	}).OnError(func(error) {
		m.fail(action, text)
	})
	return nil
}

func (m *Messenger) finish(action internal.MessengerAction, name string, reply *internal.MessengerReply) {
	m.pending = false
	switch action {
	case internal.MessengerAction_MESSENGER_ACTION_ENTER:
		m.roomID = reply.RoomId
		m.owner.listener.OnMessengerJoined(m.owner, uint8(reply.Position), reply.Members, m.queued)
		m.joining = 0
		m.queued = nil
		if m.lookChanged {
			m.lookChanged = false
			m.shareLook()
		}
	case internal.MessengerAction_MESSENGER_ACTION_LEAVE:
		m.roomID = 0
	case internal.MessengerAction_MESSENGER_ACTION_INVITE:
		m.owner.listener.OnMessengerInviteResult(m.owner, name, true)
	}
}

func (m *Messenger) fail(action internal.MessengerAction, name string) {
	m.pending = false
	m.joining = 0
	m.queued = nil
	m.lookChanged = false
	if action == internal.MessengerAction_MESSENGER_ACTION_INVITE {
		m.owner.listener.OnMessengerInviteResult(m.owner, name, false)
	}
}

// shareLook shows the owner's new look to its room. A look changed while a join waits is shared once the join ends.
func (m *Messenger) shareLook() {
	if m.roomID == 0 {
		m.lookChanged = m.joining != 0
		return
	}
	m.owner.gameWorld.Systems().Get[*InternalSystem]().MessengerAsync(nil, m.owner, internal.MessengerAction_MESSENGER_ACTION_AVATAR, 0, "").OnError(func(error) {})
}

// Room events can arrive before the join reply; they are held until the
// client has entered the room.
func (m *Messenger) Receive(event MessengerEvent) {
	if event.Action == internal.MessengerAction_MESSENGER_ACTION_INVITE || event.RoomID == m.roomID {
		m.owner.listener.OnMessengerEvent(m.owner, event)
		return
	}
	if event.RoomID != m.joining || len(m.queued) >= messengerQueueLimit {
		return
	}

	m.queued = append(m.queued, event)
}
