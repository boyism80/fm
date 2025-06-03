package actor

import (
	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/common/handler"
	common_msg "github.com/boyism80/fm/common/msg"
	"github.com/boyism80/fm/common/types"
	"github.com/boyism80/fm/game/data"
	"github.com/boyism80/fm/game/entity"
	"github.com/boyism80/fm/game/msg"
	"github.com/boyism80/fm/game/protocol/resp"
)

type NpcActor struct {
	entity.Npc
	handler *handler.MessageHandler
}

func NewNpcActorProps(ctx actor.Context, spec *data.NPCSpec, sequence uint32) *actor.Props {
	return actor.PropsFromProducer(func() actor.Actor {
		actor := &NpcActor{
			handler: handler.NewMessageHandler(),
			Npc: entity.Npc{
				ID:   sequence,
				Spec: spec,
			},
		}

		handler.RegisterHandler(ctx, actor, actor.handler, onNpcPlayerWarped)

		return actor
	})
}

func (state *NpcActor) Receive(context actor.Context) {
	state.handler.Handle(context)
}

func onNpcPlayerWarped(ctx actor.Context, state *NpcActor, m *msg.Warped) {
	ctx.Send(m.Sender, &common_msg.SendProtocol{
		Protocol: &resp.SpawnNpc{
			NPC:     &state.Npc,
			Visible: true,
		},
		Policy: types.SEND_POLICY_ENCRYPT,
	})

	ctx.Send(m.Sender, &common_msg.SendProtocol{
		Protocol: &resp.NpcControl{
			NPC:     &state.Npc,
			MiniMap: true,
		},
		Policy: types.SEND_POLICY_ENCRYPT,
	})
}
