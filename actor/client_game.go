package actor

import (
	"log"
	"net"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/asynkron/protoactor-go/scheduler"
	"github.com/boyism80/fm/context"
	"github.com/boyism80/fm/encrypt"
	"github.com/boyism80/fm/entity"
	"github.com/boyism80/fm/handler"
	"github.com/boyism80/fm/stream"
	"github.com/boyism80/fm/types"
	"github.com/boyism80/fm/util"
)

type GameClientActor struct {
	Context        context.IServerContext
	Character      *entity.Character
	Conn           net.Conn
	Buffer         []byte
	messageHandler *handler.MessageHandler
	packetHandler  *handler.PacketHandler
	sendEncryption encrypt.Encryption
	recvEncryption encrypt.Encryption
	stopTimer      scheduler.CancelFunc
}

func NewGameClientActor(ctx actor.Context, serverCtx context.IServerContext, conn net.Conn) actor.Actor {
	ivSend := []byte{0x2F, 0xA3, 0x65, 0x43}
	ivRecv := []byte{0x65, 0x56, 0x12, 0xFD}

	act := &GameClientActor{
		Context:        serverCtx,
		Conn:           conn,
		Buffer:         []byte{},
		messageHandler: handler.NewMessageHandler(),
		packetHandler:  handler.NewPacketHandler(),
		sendEncryption: encrypt.NewEncryption(ivSend, -5),
		recvEncryption: encrypt.NewEncryption(ivRecv, 5),
	}

	RegisterGameClientMessageHandlers(ctx, act, act.messageHandler)
	RegisterGameClientPacketHandler(ctx, act, act.packetHandler)
	return act
}

func (state *GameClientActor) Receive(context actor.Context) {
	state.messageHandler.Handle(context)
}

func (state *GameClientActor) Invoke(ctx actor.Context, header int, data []byte) error {
	return state.packetHandler.Handle(ctx, header, data)
}

func (c *GameClientActor) BindCharacter(ctx actor.Context, ch *entity.Character) {
	c.Character = ch
	RegisterCharacterHandlers(ctx, c.Character, c.messageHandler)
}

func (c *GameClientActor) ID() int64 {
	if c.Character == nil {
		return 0
	}
	return c.Character.Object.ID
}

func (c *GameClientActor) Position() types.Vec2 {
	if c.Character == nil {
		return types.Vec2{}
	}
	return c.Character.Object.Position
}

func (c *GameClientActor) Type() types.ObjectType {
	return types.ObjectTypePlayer
}

func (c *GameClientActor) Name() string {
	if c.Character == nil {
		return ""
	}
	return c.Character.Name
}

func (state *GameClientActor) Send(p types.Packet, policy types.SendPolicy) {
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
