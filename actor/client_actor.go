package actor

import (
	"fmt"

	protoactor "github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/handler"
	"github.com/boyism80/fm/model"
	"github.com/boyism80/fm/msg"
	"github.com/boyism80/fm/packet"
	"github.com/boyism80/fm/stream"
)

type ClientActor struct {
	Client        *model.Client
	handler       *handler.MessageHandler
	packetHandler *handler.PacketHandler
}

func NewClientActor(Client *model.Client) protoactor.Actor {
	act := &ClientActor{
		Client:        Client,
		handler:       handler.NewMessageHandler(),
		packetHandler: handler.NewPacketHandler(),
	}

	RegisterClientHandlers(act.Client, act.handler)

	handler.RegisterPacketHandler(1, act.Client, act.packetHandler, handleMove)
	handler.RegisterPacketHandler(2, act.Client, act.packetHandler, handleAttack)
	return act
}

func RegisterClientHandlers(m *model.Client, h *handler.MessageHandler) {
	handler.RegisterHandler(m, h, handleClientConnected)
}

func (state *ClientActor) Receive(context protoactor.Context) {
	state.handler.Handle(context)
}

func handleClientConnected(client *model.Client, ctx protoactor.Context, m *msg.ClientConnected) {
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

				fmt.Println(header)
				fmt.Println(data)

				reader.DiscardRead()
			}
		}
	}()
}

func handleMove(client *model.Client, request *packet.MovePacket) {

}

func handleAttack(client *model.Client, request *packet.MovePacket) {

}
