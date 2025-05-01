package actor

import (
	"fmt"
	"net"

	protoactor "github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/handler"
	"github.com/boyism80/fm/msg"
)

func RegisterLoginServerHandlers(ctx protoactor.Context, m *LoginServerActor, h *handler.MessageHandler) {
	handler.RegisterHandler(ctx, m, h, onLoginStart)
	handler.RegisterHandler(ctx, m, h, onLoginClientAccepted)
}

func onLoginStart(ctx protoactor.Context, obj *LoginServerActor, m *msg.StartListening) {
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

func onLoginClientAccepted(ctx protoactor.Context, obj *LoginServerActor, m *msg.ClientConnected) {
	props := protoactor.PropsFromProducer(func() protoactor.Actor {
		return NewLoginClientActor(ctx, m.ServerCtx, m.Conn)
	})
	pid := ctx.Spawn(props)
	ctx.Send(pid, m)
}
