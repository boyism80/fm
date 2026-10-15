package game

import (
	"errors"

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
	RoomID  uint32
	pending bool
	joining uint32
	queued  []MessengerEvent
}

func (m *Messenger) Begin(action internal.MessengerAction, roomID uint32) error {
	if m.pending {
		return ErrMessengerBusy
	}
	if action != internal.MessengerAction_MESSENGER_ACTION_ENTER && action != internal.MessengerAction_MESSENGER_ACTION_DECLINE && m.RoomID == 0 {
		return ErrMessengerNotInRoom
	}

	m.pending = true
	if action == internal.MessengerAction_MESSENGER_ACTION_ENTER {
		m.joining = roomID
	}
	return nil
}

func (m *Messenger) Finish(action internal.MessengerAction, name string, reply *internal.MessengerReply) {
	m.pending = false
	switch action {
	case internal.MessengerAction_MESSENGER_ACTION_ENTER:
		m.RoomID = reply.RoomId
		m.owner.Listener.OnMessengerJoined(m.owner, uint8(reply.Position), reply.Members, m.queued)
		m.joining = 0
		m.queued = nil
	case internal.MessengerAction_MESSENGER_ACTION_LEAVE:
		m.RoomID = 0
	case internal.MessengerAction_MESSENGER_ACTION_INVITE:
		m.owner.Listener.OnMessengerInviteResult(m.owner, name, true)
	}
}

func (m *Messenger) Fail(action internal.MessengerAction, name string) {
	m.pending = false
	m.joining = 0
	m.queued = nil
	if action == internal.MessengerAction_MESSENGER_ACTION_INVITE {
		m.owner.Listener.OnMessengerInviteResult(m.owner, name, false)
	}
}

// Room events can arrive before the join reply; they are held until the
// client has entered the room.
func (m *Messenger) Receive(event MessengerEvent) {
	if event.Action == internal.MessengerAction_MESSENGER_ACTION_INVITE || event.RoomID == m.RoomID {
		m.owner.Listener.OnMessengerEvent(m.owner, event)
		return
	}
	if event.RoomID != m.joining || len(m.queued) >= messengerQueueLimit {
		return
	}

	m.queued = append(m.queued, event)
}
