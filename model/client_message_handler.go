package model

import (
	"log"
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/asynkron/protoactor-go/scheduler"
	"github.com/boyism80/fm/encrypt"
	"github.com/boyism80/fm/handler"
	"github.com/boyism80/fm/msg"
	"github.com/boyism80/fm/packet/resp"
	"github.com/boyism80/fm/stream"
)

func RegisterClientHandlers(m *ClientActor, h *handler.MessageHandler) {
	handler.RegisterHandler(m, h, onClientConnected)
	handler.RegisterHandler(m, h, onClientStopping)
	handler.RegisterHandler(m, h, onPing)
}

func (client *ClientActor) ReceivePackets(ctx actor.Context) {
	buf := make([]byte, 1024)
	conn := client.Conn

	for {
		n, err := conn.Read(buf)
		if err != nil {
			ctx.Stop(ctx.Self())
			break
		}

		client.Buffer = append(client.Buffer, buf[:n]...)
		reader := stream.NewStreamReader(&client.Buffer, stream.LittleEndian)

		for {
			header, err := reader.Read(4)
			if err != nil {
				break
			}

			if !client.recvEncryption.CheckPacketHeader(header) {
				break
			}

			len, err := encrypt.GetPacketLength(header)
			if err != nil {
				break
			}

			data, err := reader.Read(int(len))
			if err != nil {
				break
			}

			data = client.recvEncryption.Decrypt(data)
			opcodeReader := stream.NewStreamReader(&data, stream.LittleEndian)
			opcode, err := opcodeReader.ReadU16()
			if err != nil {
				break
			}

			data = data[2:]
			err = client.Invoke(int(opcode), data)
			if err != nil {
				log.Println(err)
			}
			reader.DiscardRead()
		}
	}
}

func onClientConnected(client *ClientActor, ctx actor.Context, m *msg.ClientConnected) {
	log.Println("클라이언트 연결됨")
	timer := scheduler.NewTimerScheduler(ctx.ActorSystem().Root)
	client.stopTimer = timer.SendRepeatedly(10*time.Second, 10*time.Second, ctx.Self(), &msg.Ping{})

	client.Send(&resp.Welcome{
		SendIv: client.sendEncryption.IV(),
		RecvIv: client.recvEncryption.IV()}, SEND_POLICY_RAW)
	go client.ReceivePackets(ctx)
}

func onClientStopping(client *ClientActor, ctx actor.Context, m *actor.Stopping) {
	log.Println("클라이언트 접속 종료")
	client.Conn.Close()
	if client.stopTimer != nil {
		client.stopTimer()
	}
}

func onPing(client *ClientActor, ctx actor.Context, m *msg.Ping) {
	client.Send(&resp.Ping{}, SEND_POLICY_ENCRYPT)
}
