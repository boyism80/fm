package server

import (
	"github.com/boyism80/fm/protocol/response"
	"github.com/boyism80/fm/services/game/client"
	"github.com/boyism80/fm/services/game/constant"
)

func (gs *GameServer) BroadcastNotice(messageType constant.ServerMessageType, message string, channel int, ear bool) {
	for _, c := range gs.Clients() {
		gameClient, ok := c.(*client.GameClient)
		if !ok {
			continue
		}
		ch := gameClient.GetCharacter()
		if ch == nil {
			continue
		}
		ch.Listener.OnNotice(ch, messageType, message, channel, ear)
	}
}

func (gs *GameServer) BroadcastMegaphone(megaphone *response.Megaphone) {
	for _, c := range gs.Clients() {
		gameClient, ok := c.(*client.GameClient)
		if !ok {
			continue
		}
		ch := gameClient.GetCharacter()
		if ch == nil {
			continue
		}
		ch.Listener.OnMegaphone(ch, megaphone)
	}
}
