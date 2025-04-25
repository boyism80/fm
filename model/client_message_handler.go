package model

import (
	"log"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/handler"
	"github.com/boyism80/fm/msg"
	"github.com/boyism80/fm/stream"
)

func RegisterClientHandlers(m *ClientActor, h *handler.MessageHandler) {
	handler.RegisterHandler(m, h, onClientConnected)
	handler.RegisterHandler(m, h, onClientStopping)
}

func onClientConnected(client *ClientActor, ctx actor.Context, m *msg.ClientConnected) {
	conn := client.Conn
	go func() {
		buf := make([]byte, 1024)

		for {
			n, err := conn.Read(buf)
			if err != nil {
				ctx.Stop(ctx.Self())
				conn.Close()
				break
			}

			client.Buffer = append(client.Buffer, buf[:n]...)
			reader := stream.NewStreamReader(client.Buffer, stream.LittleEndian)

			for {
				header, err := reader.Read32()
				if err != nil {
					break
				}

				len, err := reader.ReadU32()
				if err != nil {
					break
				}

				data, err := reader.Read(int(len))
				if err != nil {
					break
				}

				client.Invoke(int(header), data)
				reader.DiscardRead()
			}
		}
	}()
}

func onClientStopping(client *ClientActor, ctx actor.Context, m *actor.Stopping) {
	log.Println("클라이언트 접속 종료")
	client.Conn.Close()
}
