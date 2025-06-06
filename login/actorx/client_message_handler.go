package actorx

import (
	"log"

	protoactor "github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/common/crypt"
	"github.com/boyism80/fm/common/handler"
	"github.com/boyism80/fm/common/msg"
	common_resp "github.com/boyism80/fm/common/protocol/resp"
	"github.com/boyism80/fm/common/stream"
	"github.com/boyism80/fm/common/types"
	"github.com/boyism80/fm/common/util"
	"github.com/boyism80/fm/login/protocol/resp"
)

func RegisterLoginClientMessageHandlers(ctx protoactor.Context, m *LoginClientActor, h *handler.MessageHandler) {
	handler.RegisterHandler(ctx, m, h, onLoginClientConnected)
	handler.RegisterHandler(ctx, m, h, onLoginClientStopping)
	handler.RegisterHandler(ctx, m, h, onLoginClientInvokeHandler)
	handler.RegisterHandler(ctx, m, h, onLoginClientSendProtocol)
	handler.RegisterHandler(ctx, m, h, onLoginClientPing)
}

func (client *LoginClientActor) ReceivePackets(ctx protoactor.Context, pid *protoactor.PID) {
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

			len, err := crypt.GetPacketLength(header)
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
			reader.DiscardRead()
		}
	}
}

func onLoginClientConnected(ctx protoactor.Context, client *LoginClientActor, m *msg.ClientConnected) {
	log.Println("클라이언트 연결됨")
	// timer := scheduler.NewTimerScheduler(ctx.ActorSystem().Root)
	// client.stopTimer = timer.SendRepeatedly(10*time.Second, 10*time.Second, ctx.Self(), &msg.Ping{})

	// 이건 로그인서버에서만 보내줘야함
	// 게임서버에서 보내주면 안됨
	client.Send(&common_resp.Welcome{
		SendIv: client.sendEncryption.IV(),
		RecvIv: client.recvEncryption.IV()}, types.SEND_POLICY_RAW)

	client.Send(&resp.LoginFailed{
		Reason: resp.LoginFailedReasonNoPopup}, types.SEND_POLICY_ENCRYPT)

	go client.ReceivePackets(ctx, ctx.Self())
}

func onLoginClientStopping(ctx protoactor.Context, client *LoginClientActor, m *protoactor.Stopping) {
	log.Println("클라이언트 접속 종료")
	client.Conn.Close()
	if client.stopTimer != nil {
		client.stopTimer()
	}
}

func onLoginClientInvokeHandler(ctx protoactor.Context, client *LoginClientActor, m *msg.InvokeHandler) {
	client.Invoke(ctx, int(m.Opcode), m.Data)
}

func onLoginClientSendProtocol(ctx protoactor.Context, state *LoginClientActor, m *msg.SendProtocol) {
	state.Send(m.Protocol, m.Policy)
}

func onLoginClientPing(ctx protoactor.Context, client *LoginClientActor, m *msg.Ping) {
	client.Send(&common_resp.Ping{}, types.SEND_POLICY_ENCRYPT)
}
