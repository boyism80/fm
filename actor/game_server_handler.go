package actor

import (
	"fmt"
	"net"

	protoactor "github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/handler"
	"github.com/boyism80/fm/msg"
)

func RegisterGameServerHandlers(ctx protoactor.Context, m *GameServerActor, h *handler.MessageHandler) {
	handler.RegisterHandler(ctx, m, h, onGameStart)
	handler.RegisterHandler(ctx, m, h, onGameClientAccepted)
}

func onGameStart(ctx protoactor.Context, obj *GameServerActor, m *msg.StartListening) {
	port := fmt.Sprintf(":%d", m.Port)

	// 임시 맵 추가
	props := NewMapActorProps(ctx, m.ServerCtx)
	pid := ctx.Spawn(props)
	m.ServerCtx.Load(200000301, pid)

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
			return NewBotActor(ctx, m.ServerCtx)
		})
		pid := ctx.Spawn(props)
		ctx.Send(pid, &msg.BotConnect{Port: m.Port})

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
