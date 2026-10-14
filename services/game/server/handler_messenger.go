package server

import (
	"context"

	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/core/async"
	"github.com/boyism80/fm/core/rpc"
	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
	"github.com/boyism80/fm/protocol/request"
	"github.com/boyism80/fm/services/game/client"
	"github.com/boyism80/fm/services/game/domain"
	"github.com/boyism80/fm/stream"
)

type Messenger struct {
	gs *GameServer
}

func (Messenger) New(gs *GameServer) *Messenger {
	return &Messenger{
		gs: gs,
	}
}

func (h *Messenger) Handle(ctx *core.ClientContext, req *request.Messenger) error {
	if ctx.ActorContext == nil {
		return nil
	}

	gameClient, ok := ctx.Client.(*client.GameClient)
	if !ok {
		return nil
	}
	ch := gameClient.GetCharacter()
	if ch == nil {
		return nil
	}

	switch req.Action {
	case internal.MessengerAction_MESSENGER_ACTION_ENTER, internal.MessengerAction_MESSENGER_ACTION_LEAVE:
	case internal.MessengerAction_MESSENGER_ACTION_INVITE, internal.MessengerAction_MESSENGER_ACTION_DECLINE:
		if !domain.ClientText(req.Text).Valid(13) {
			return nil
		}
	case internal.MessengerAction_MESSENGER_ACTION_CHAT:
		if !domain.ClientText(req.Text).Valid(200) {
			return nil
		}
	default:
		return nil
	}
	if err := ch.Messenger.Begin(req.Action, req.RoomID); err != nil {
		return nil
	}

	w := stream.NewStreamWriter(stream.LittleEndian)
	ch.BuildS2C().SerializeLook(w)
	internalReq := &internal.MessengerRequest{
		WorldId:     h.gs.config.WorldId,
		CharacterId: ch.GetID(),
		ChannelId:   h.gs.config.ChannelId,
		Action:      req.Action,
		RoomId:      req.RoomID,
		Text:        req.Text,
		Look:        w.Bytes(),
	}
	async.NewTask(ctx.ActorContext, core.InternalRPCPerStepTimeout).ThenRPC(
		func(c context.Context) (*internal.MessengerReply, error) {
			return h.gs.systems.Get[*rpc.InternalConnection]().Client.Messenger(c, internalReq)
		},
		func(reply *internal.MessengerReply) error {
			if ch.LoggedOut() {
				return nil
			}

			if !reply.Ok {
				ch.Messenger.Fail(req.Action, req.Text)
				return nil
			}
			ch.Messenger.Finish(req.Action, req.Text, reply)
			return nil
		},
	).OnError(func(error) {
		ch.Messenger.Fail(req.Action, req.Text)
	})
	return nil
}
