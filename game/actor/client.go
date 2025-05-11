package actor

import (
	"log"
	"net"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/asynkron/protoactor-go/scheduler"
	"github.com/boyism80/fm/common/context"
	"github.com/boyism80/fm/common/encrypt"
	"github.com/boyism80/fm/common/handler"
	"github.com/boyism80/fm/common/stream"
	"github.com/boyism80/fm/common/types"
	"github.com/boyism80/fm/common/util"
	"github.com/boyism80/fm/game/entity"
)

type GameClientActor struct {
	context        *context.ServerContext
	character      *entity.Character
	conn           net.Conn
	buffer         []byte
	messageHandler *handler.MessageHandler
	packetHandler  *handler.PacketHandler
	commandHandler *handler.CommandHandler
	sendEncryption encrypt.Encryption
	recvEncryption encrypt.Encryption
	stopTimer      scheduler.CancelFunc
}

func NewGameClientActor(ctx actor.Context, serverCtx *context.ServerContext, conn net.Conn) actor.Actor {
	ivSend := []byte{0x2F, 0xA3, 0x65, 0x43}
	ivRecv := []byte{0x65, 0x56, 0x12, 0xFD}

	act := &GameClientActor{
		context:        serverCtx,
		conn:           conn,
		buffer:         []byte{},
		messageHandler: handler.NewMessageHandler(),
		packetHandler:  handler.NewPacketHandler(),
		commandHandler: handler.NewCommandHandler(),
		sendEncryption: encrypt.NewEncryption(ivSend, -5),
		recvEncryption: encrypt.NewEncryption(ivRecv, 5),
	}

	RegisterGameClientMessageHandlers(ctx, act, act.messageHandler)
	RegisterGameClientPacketHandler(ctx, act, act.packetHandler)
	RegisterGameClientCommandHandler(ctx, act, act.commandHandler)
	return act
}

func (state *GameClientActor) Receive(context actor.Context) {
	state.messageHandler.Handle(context)
}

func (state *GameClientActor) Invoke(ctx actor.Context, header int, data []byte) error {
	return state.packetHandler.Handle(ctx, header, data)
}

func (c *GameClientActor) BindCharacter(ctx actor.Context, ch *entity.Character) {
	c.character = ch
	RegisterLifeHandlers(ctx, &c.character.Life, c.messageHandler)
}

func (c *GameClientActor) Name() string {
	if c.character == nil {
		return ""
	}
	return c.character.Name
}

func (actor *GameClientActor) Send(p types.Packet, policy types.SendPolicy) {
	writer := stream.NewStreamWriter(stream.LittleEndian)
	p.Serialize(writer)
	bytes := writer.Bytes()

	log.Println("[S] " + util.ToHexString(bytes))

	if policy == types.SEND_POLICY_RAW {
		actor.conn.Write(bytes)
		return
	}

	writer = stream.NewStreamWriter(stream.LittleEndian)
	if policy&types.SEND_POLICY_ENCRYPT != 0 {
		header := actor.sendEncryption.GetPacketHeader(len(bytes))
		writer.Write(header)
		bytes = actor.sendEncryption.Encrypt(bytes)
	}

	writer.Write(bytes)
	bytes = writer.Bytes()
	actor.conn.Write(bytes)
}
