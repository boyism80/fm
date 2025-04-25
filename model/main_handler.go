package model

import (
	"fmt"
	"net"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/handler"
	"github.com/boyism80/fm/msg"
)

func RegisterMainHandlers(m *MainActor, h *handler.MessageHandler) {
	handler.RegisterHandler(m, h, onStartListening)
	handler.RegisterHandler(m, h, onAcceptedClient)
}

func onStartListening(obj *MainActor, ctx actor.Context, m *msg.StartListening) {
	port := fmt.Sprintf(":%d", m.Port)
	go func() {
		listener, err := net.Listen("tcp", port)
		if err != nil {
			fmt.Println("Error listening:", err)
			return
		}
		defer listener.Close()
		fmt.Println("Listening on", port)

		// 봇 테스트
		props := actor.PropsFromProducer(func() actor.Actor {
			return NewBotActor()
		})
		pid := ctx.Spawn(props)
		ctx.Send(pid, &msg.BotConnect{Port: m.Port})

		for {
			conn, err := listener.Accept()
			if err != nil {
				fmt.Println("Error accepting:", err)
				continue
			}
			ctx.Send(ctx.Self(), &msg.ClientConnected{Conn: conn})
		}
	}()
}

func onAcceptedClient(obj *MainActor, ctx actor.Context, m *msg.ClientConnected) {
	props := actor.PropsFromProducer(func() actor.Actor {
		return NewClientActor(m.Conn)
	})
	pid := ctx.Spawn(props)
	ctx.Send(pid, m)
}
