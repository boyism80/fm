package actor

import (
	"fmt"
	"net"

	protoactor "github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/common/handler"
	"github.com/boyism80/fm/common/msg"

	bot "github.com/boyism80/fm/bot/actor"
	bot_msg "github.com/boyism80/fm/bot/msg"
)

type GameServerActor struct {
	handler *handler.MessageHandler
}

func NewGameServerActor() protoactor.Actor {
	actor := &GameServerActor{
		handler: handler.NewMessageHandler(),
	}
	handler.RegisterHandler(nil, actor, actor.handler, onGameStart)
	handler.RegisterHandler(nil, actor, actor.handler, onGameClientAccepted)
	return actor
}

func (state *GameServerActor) Receive(context protoactor.Context) {
	state.handler.Handle(context)
}

func onGameStart(ctx protoactor.Context, obj *GameServerActor, m *msg.StartListening) {
	port := fmt.Sprintf(":%d", m.Port)

	for _, spec := range m.ServerCtx.Resources.Maps {
		props := NewMapActorProps(ctx, m.ServerCtx, spec)
		pid := ctx.Spawn(props)
		m.ServerCtx.MapActors[spec.Id] = pid
	}

	go func() {
		listener, err := net.Listen("tcp", port)
		if err != nil {

			fmt.Println("Error listening:", err)
			return
		}
		defer listener.Close()
		fmt.Println("Listening on", port)

		// 봇 테스트
		props := protoactor.PropsFromProducer(func() protoactor.Actor {
			return bot.NewBotActor(ctx, m.ServerCtx)
		})
		pid := ctx.Spawn(props)
		ctx.Send(pid, &bot_msg.BotConnect{Port: m.Port})

		for {
			conn, err := listener.Accept()
			if err != nil {
				fmt.Println("Error accepting:", err)
				continue
			}
			ctx.Send(ctx.Self(), &msg.ClientConnected{Conn: conn, ServerCtx: m.ServerCtx})
		}
	}()
}

func onGameClientAccepted(ctx protoactor.Context, obj *GameServerActor, m *msg.ClientConnected) {
	props := protoactor.PropsFromProducer(func() protoactor.Actor {
		return NewGameClientActor(ctx, m.ServerCtx, m.Conn)
	})
	pid := ctx.Spawn(props)
	ctx.Send(pid, m)
}
