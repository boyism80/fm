package game

import (
	"sync/atomic"

	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/maple/constant"
	"github.com/boyism80/fm/protocol/response"
)

// characterClient is a game client; the client package depends on this one.
type characterClient interface {
	GetCharacter() *Character
}

type NoticeSystem struct {
	server *core.ServerCore
	muted  atomic.Bool
}

func NewNoticeSystem(server *core.ServerCore) *NoticeSystem {
	return &NoticeSystem{server: server}
}

func (s *NoticeSystem) Start() error { return nil }
func (s *NoticeSystem) Stop()        {}

func (s *NoticeSystem) BroadcastNotice(messageType constant.ServerMessageType, message string, channel int, ear bool) {
	for _, c := range s.server.Clients() {
		gameClient, ok := c.(characterClient)
		if !ok {
			continue
		}
		ch := gameClient.GetCharacter()
		if ch == nil {
			continue
		}
		ch.listener.OnNotice(ch, messageType, message, channel, ear)
	}
}

func (s *NoticeSystem) BroadcastMegaphone(megaphone *response.Megaphone) {
	for _, c := range s.server.Clients() {
		gameClient, ok := c.(characterClient)
		if !ok {
			continue
		}
		ch := gameClient.GetCharacter()
		if ch == nil {
			continue
		}
		ch.listener.OnMegaphone(ch, megaphone)
	}
}
