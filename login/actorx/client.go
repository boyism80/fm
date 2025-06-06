package actorx

import (
	"log"
	"net"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/asynkron/protoactor-go/scheduler"
	"github.com/boyism80/fm/common/context"
	"github.com/boyism80/fm/common/crypt"
	"github.com/boyism80/fm/common/handler"
	"github.com/boyism80/fm/common/stream"
	"github.com/boyism80/fm/common/types"
	"github.com/boyism80/fm/common/util"
)

type LoginClientActor struct {
	Context        *context.ServerContext
	Conn           net.Conn
	Buffer         []byte
	messageHandler *handler.MessageHandler
	packetHandler  *handler.PacketHandler
	sendEncryption crypt.Encryption
	recvEncryption crypt.Encryption
	stopTimer      scheduler.CancelFunc
}

func NewLoginClientActor(ctx actor.Context, serverCtx *context.ServerContext, conn net.Conn) actor.Actor {

	ivSend := []byte{0x2F, 0xA3, 0x65, 0x43}
	ivRecv := []byte{0x65, 0x56, 0x12, 0xFD}

	actor := &LoginClientActor{
		Context:        serverCtx,
		Conn:           conn,
		Buffer:         []byte{},
		messageHandler: handler.NewMessageHandler(),
		packetHandler:  handler.NewPacketHandler(),
		sendEncryption: crypt.NewEncryption(ivSend, -5),
		recvEncryption: crypt.NewEncryption(ivRecv, 5),
	}

	RegisterLoginClientMessageHandlers(ctx, actor, actor.messageHandler)
	RegisterLoginClientPacketHandler(ctx, actor, actor.packetHandler)
	return actor
}

func (state *LoginClientActor) Receive(context actor.Context) {
	state.messageHandler.Handle(context)
}

func (state *LoginClientActor) Invoke(ctx actor.Context, header int, data []byte) error {
	return state.packetHandler.Handle(ctx, header, data)
}

func (state *LoginClientActor) Send(p types.Packet, policy types.SendPolicy) {
	writer := stream.NewStreamWriter(stream.LittleEndian)
	p.Serialize(writer)
	bytes := writer.Bytes()

	log.Println("[S] " + util.ToHexString(bytes))

	if policy == types.SEND_POLICY_RAW {
		state.Conn.Write(bytes)
		return
	}

	writer = stream.NewStreamWriter(stream.LittleEndian)
	if policy&types.SEND_POLICY_ENCRYPT != 0 {
		header := state.sendEncryption.GetPacketHeader(len(bytes))
		writer.Write(header)
		bytes = state.sendEncryption.Encrypt(bytes)
	}

	writer.Write(bytes)
	bytes = writer.Bytes()
	state.Conn.Write(bytes)
}
