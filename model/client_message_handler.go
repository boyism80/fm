package model

import (
	"fmt"

	protoactor "github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/handler"
	"github.com/boyism80/fm/msg"
	"github.com/boyism80/fm/stream"
)

func RegisterClientHandlers(m *ClientActor, h *handler.MessageHandler) {
	handler.RegisterHandler(m, h, onClientConnected)
}

func onClientConnected(client *ClientActor, ctx protoactor.Context, m *msg.ClientConnected) {
	conn := client.Conn
	go func() {
		buf := make([]byte, 1024)

		for {
			n, err := conn.Read(buf)
			if err != nil {
				fmt.Println("Error reading from connection:", err)
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
