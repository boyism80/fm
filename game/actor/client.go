package actor

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
	"github.com/boyism80/fm/game/entity"
)

type GameClientActor struct {
	ctx            *context.ServerContext
	ch             *entity.Character
	conn           net.Conn
	buffer         []byte
	messageHandler *handler.MessageHandler
	packetHandler  *handler.PacketHandler
	commandHandler *handler.CommandHandler
	sendCrypt      crypt.Encryption
	receiveCrypt   crypt.Encryption
	stopTimer      scheduler.CancelFunc
}

func NewGameClientActor(ctx actor.Context, serverCtx *context.ServerContext, conn net.Conn) actor.Actor {
	ivSend := []byte{0x2F, 0xA3, 0x65, 0x43}
	ivRecv := []byte{0x65, 0x56, 0x12, 0xFD}

	actor := &GameClientActor{
		ctx:            serverCtx,
		conn:           conn,
		buffer:         []byte{},
		messageHandler: handler.NewMessageHandler(),
		packetHandler:  handler.NewPacketHandler(),
		commandHandler: handler.NewCommandHandler(),
		sendCrypt:      crypt.NewEncryption(ivSend, -5),
		receiveCrypt:   crypt.NewEncryption(ivRecv, 5),
	}

	RegisterGameClientMessageHandlers(ctx, actor, actor.messageHandler)
	RegisterGameClientPacketHandler(ctx, actor, actor.packetHandler)
	RegisterGameClientCommandHandler(ctx, actor, actor.commandHandler)
	return actor
}

func (state *GameClientActor) Receive(context actor.Context) {
	state.messageHandler.Handle(context)
}

func (state *GameClientActor) Invoke(ctx actor.Context, header int, data []byte) error {
	return state.packetHandler.Handle(ctx, header, data)
}

func (state *GameClientActor) Name() string {
	if state.ch == nil {
		return ""
	}
	return state.ch.Name
}

func (state *GameClientActor) Send(p types.Packet, policy types.SendPolicy) {
	writer := stream.NewStreamWriter(stream.LittleEndian)
	p.Serialize(writer)
	bytes := writer.Bytes()

	log.Println("[S] " + util.ToHexString(bytes))

	if policy == types.SEND_POLICY_RAW {
		state.conn.Write(bytes)
		return
	}

	writer = stream.NewStreamWriter(stream.LittleEndian)
	if policy&types.SEND_POLICY_ENCRYPT != 0 {
		header := state.sendCrypt.GetPacketHeader(len(bytes))
		writer.Write(header)
		bytes = state.sendCrypt.Encrypt(bytes)
	}

	writer.Write(bytes)
	bytes = writer.Bytes()
	state.conn.Write(bytes)
}
