package actor

import (
	"log"

	"github.com/asynkron/protoactor-go/actor"
	protoactor "github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/encrypt"
	"github.com/boyism80/fm/handler"
	"github.com/boyism80/fm/msg"
	"github.com/boyism80/fm/packet/resp"
	"github.com/boyism80/fm/stream"
	"github.com/boyism80/fm/types"
	"github.com/boyism80/fm/util"
)

func RegisterClientHandlers(ctx protoactor.Context, m *ClientActor, h *handler.MessageHandler) {
	handler.RegisterHandler(ctx, m, h, onClientConnected)
	handler.RegisterHandler(ctx, m, h, onClientStopping)
	handler.RegisterHandler(ctx, m, h, onClientInvokeHandler)
	handler.RegisterHandler(ctx, m, h, onSendProtocol)
	handler.RegisterHandler(ctx, m, h, onPing)

	handler.RegisterHandler(ctx, m, h, onClientWarped)
}

func (client *ClientActor) ReceivePackets(ctx protoactor.Context, pid *protoactor.PID) {
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
			log.Println("[R] " + util.ToHexString(data))
			opcodeReader := stream.NewStreamReader(&data, stream.LittleEndian)
			opcode, err := opcodeReader.ReadU16()
			if err != nil {
				break
			}

			data = data[2:]
			ctx.Send(pid, &msg.InvokeHandler{
				Opcode: opcode,
				Data:   data,
			})
			// err = client.Invoke(int(), data)
			// if err != nil {
			// 	log.Println(err)
			// }
			reader.DiscardRead()
		}
	}
}

func onClientConnected(ctx protoactor.Context, client *ClientActor, m *msg.ClientConnected) {
	log.Println("클라이언트 연결됨")
	// timer := scheduler.NewTimerScheduler(ctx.ActorSystem().Root)
	// client.stopTimer = timer.SendRepeatedly(10*time.Second, 10*time.Second, ctx.Self(), &msg.Ping{})

	client.Send(&resp.Welcome{
		SendIv: client.sendEncryption.IV(),
		RecvIv: client.recvEncryption.IV()}, types.SEND_POLICY_RAW)
	go client.ReceivePackets(ctx, ctx.Self())
}

func onClientStopping(ctx protoactor.Context, client *ClientActor, m *protoactor.Stopping) {
	log.Println("클라이언트 접속 종료")
	client.Conn.Close()
	if client.stopTimer != nil {
		client.stopTimer()
	}
}

func onClientInvokeHandler(ctx protoactor.Context, client *ClientActor, m *msg.InvokeHandler) {
	client.Invoke(ctx, int(m.Opcode), m.Data)
}

func onSendProtocol(ctx actor.Context, state *ClientActor, m *msg.SendProtocol) {
	state.Send(m.Protocol, m.Policy)
}

func onPing(ctx protoactor.Context, client *ClientActor, m *msg.Ping) {
	client.Send(&resp.Ping{}, types.SEND_POLICY_ENCRYPT)
}

func onClientWarped(ctx protoactor.Context, client *ClientActor, m *msg.Warped) {
	// 어떤 플레이어가 내가 속한 맵으로 왔다.
	// 나의 캐릭터 정보를 해당 플레이어에게 보내준다.

	ctx.Send(m.Sender, &msg.SendProtocol{
		Protocol: &resp.SpawnPlayer{
			Character: client.Character,
		},
		Policy: types.SEND_POLICY_ENCRYPT,
	})
}
