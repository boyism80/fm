package actor

import (
	"fmt"
	"net"

	protoactor "github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/handler"
	"github.com/boyism80/fm/model"
	"github.com/boyism80/fm/msg"
)

func RegisterMainHandlers(m *model.Main, h *handler.MessageHandler) {
	handler.RegisterHandler(m, h, handleStartListening)
	handler.RegisterHandler(m, h, handleClientConnected)
}

func handleStartListening(obj *model.Main, ctx protoactor.Context, m *msg.StartListening) {
	port := fmt.Sprintf(":%d", m.Port)
	go func() {
		listener, err := net.Listen("tcp", port)
		if err != nil {
			fmt.Println("Error listening:", err)
			return
		}
		defer listener.Close()
		fmt.Println("Listening on", port)

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

func handleClientConnected(obj *model.Main, ctx protoactor.Context, m *msg.ClientConnected) {
	props := protoactor.PropsFromProducer(func() protoactor.Actor { return &ClientActor{} })
	clientPID := ctx.Spawn(props)
	ctx.Send(clientPID, m)
}
