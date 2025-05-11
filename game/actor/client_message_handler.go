package actor

import (
	"log"

	protoactor "github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/common/encrypt"
	"github.com/boyism80/fm/common/handler"
	common_msg "github.com/boyism80/fm/common/msg"
	"github.com/boyism80/fm/common/stream"
	"github.com/boyism80/fm/common/types"
	"github.com/boyism80/fm/common/util"
	"github.com/boyism80/fm/game/entity"
	"github.com/boyism80/fm/game/msg"
	"github.com/boyism80/fm/game/protocol/resp"

	common_resp "github.com/boyism80/fm/common/protocol/resp"
	login_resp "github.com/boyism80/fm/login/protocol/resp"
)

func RegisterGameClientMessageHandlers(ctx protoactor.Context, m *GameClientActor, h *handler.MessageHandler) {
	handler.RegisterHandler(ctx, m, h, onGameClientConnected)
	handler.RegisterHandler(ctx, m, h, onGameClientStopping)
	handler.RegisterHandler(ctx, m, h, onGameClientInvokeHandler)
	handler.RegisterHandler(ctx, m, h, onGameClientSendProtocol)
	handler.RegisterHandler(ctx, m, h, onGameClientPing)
	handler.RegisterHandler(ctx, m, h, onGameClientWarped)
}

func (client *GameClientActor) ReceivePackets(ctx protoactor.Context, pid *protoactor.PID) {
	buf := make([]byte, 1024)
	conn := client.conn

	for {
		n, err := conn.Read(buf)
		if err != nil {
			ctx.Stop(ctx.Self())
			break
		}

		client.buffer = append(client.buffer, buf[:n]...)
		reader := stream.NewStreamReader(&client.buffer, stream.LittleEndian)

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
			ctx.Send(pid, &common_msg.InvokeHandler{
				Opcode: opcode,
				Data:   data,
			})
			reader.DiscardRead()
		}
	}
}

func onGameClientConnected(ctx protoactor.Context, client *GameClientActor, m *common_msg.ClientConnected) {
	log.Println("클라이언트 연결됨")
	// timer := scheduler.NewTimerScheduler(ctx.ActorSystem().Root)
	// client.stopTimer = timer.SendRepeatedly(10*time.Second, 10*time.Second, ctx.Self(), &msg.Ping{})

	// 게임서버에서 보내주면 안됨
	client.Send(&common_resp.Welcome{
		SendIv: client.sendEncryption.IV(),
		RecvIv: client.recvEncryption.IV()}, types.SEND_POLICY_RAW)

	client.Send(&login_resp.LoginFailed{
		Reason: login_resp.LoginFailedReasonNoPopup}, types.SEND_POLICY_ENCRYPT)
	go client.ReceivePackets(ctx, ctx.Self())
}

func onGameClientStopping(ctx protoactor.Context, client *GameClientActor, m *protoactor.Stopping) {
	ch := client.character
	if ch != nil {
		mapActor := client.context.MapActors[client.character.Map]
		if mapActor != nil {
			ctx.Send(mapActor, &msg.LeaveMap{
				Id: ch.Id,
			})
		}
	}

	log.Println("클라이언트 접속 종료")
	client.conn.Close()
	if client.stopTimer != nil {
		client.stopTimer()
	}
}

func onGameClientInvokeHandler(ctx protoactor.Context, client *GameClientActor, m *common_msg.InvokeHandler) {
	client.Invoke(ctx, int(m.Opcode), m.Data)
}

func onGameClientSendProtocol(ctx protoactor.Context, state *GameClientActor, m *common_msg.SendProtocol) {
	state.Send(m.Protocol, m.Policy)
}

func onGameClientPing(ctx protoactor.Context, client *GameClientActor, m *common_msg.Ping) {
	client.Send(&common_resp.Ping{}, types.SEND_POLICY_ENCRYPT)
}

func onGameClientWarped(ctx protoactor.Context, client *GameClientActor, m *msg.Warped) {
	// 어떤 플레이어가 내가 속한 맵으로 왔다.
	// 나의 캐릭터 정보를 해당 플레이어에게 보내준다.

	ctx.Send(m.Sender, &common_msg.SendProtocol{
		Protocol: &resp.SpawnPlayer{
			Character:       client.character,
			BuffStates:      [4]uint32{},
			Diseases:        [4]uint32{},
			CrushRings:      []*entity.Ring{},
			FriendshipRings: []*entity.Ring{},
			MarriageRings:   []*entity.Ring{},
		},
		Policy: types.SEND_POLICY_ENCRYPT,
	})
}
