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
)

type Whisper struct {
	gs *GameServer
}

func (Whisper) New(gs *GameServer) *Whisper {
	return &Whisper{
		gs: gs,
	}
}

func (h *Whisper) Handle(ctx *core.ClientContext, req *request.Whisper) error {
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

	fail := func() {
		ch.Listener.OnWhisperFailed(ch, req.Target)
	}
	switch req.Mode {
	case internal.WhisperMode_WHISPER_MODE_FIND, internal.WhisperMode_WHISPER_MODE_SEND, internal.WhisperMode_WHISPER_MODE_FIND_BUDDY:
	default:
		fail()
		return nil
	}
	if !domain.ClientText(req.Target).Valid(13) || (req.Mode == internal.WhisperMode_WHISPER_MODE_SEND && !domain.ClientText(req.Message).Valid(200)) {
		fail()
		return nil
	}

	internalReq := &internal.WhisperRequest{
		WorldId:     h.gs.config.WorldId,
		CharacterId: ch.GetID(),
		ChannelId:   h.gs.config.ChannelId,
		Mode:        req.Mode,
		TargetName:  req.Target,
		Message:     req.Message,
	}
	async.NewTask(ctx.ActorContext, core.InternalRPCPerStepTimeout).ThenRPC(
		func(c context.Context) (*internal.WhisperReply, error) {
			return h.gs.systems.Get[*rpc.InternalConnection]().Client.Whisper(c, internalReq)
		},
		func(reply *internal.WhisperReply) error {
			if ch.LoggedOut() {
				return nil
			}

			if !reply.Found {
				fail()
				return nil
			}

			if req.Mode == internal.WhisperMode_WHISPER_MODE_SEND {
				ch.Listener.OnWhisperSent(ch, reply.Name)
			} else {
				ch.Listener.OnWhisperFound(ch, reply, req.Mode == internal.WhisperMode_WHISPER_MODE_FIND_BUDDY)
			}
			return nil
		},
	).OnError(func(error) {
		fail()
	})
	return nil
}
