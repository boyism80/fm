package actor

import (
	"crypto/rand"
	"log"
	"net"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/asynkron/protoactor-go/scheduler"
	"github.com/boyism80/fm/encrypt"
	"github.com/boyism80/fm/entity"
	"github.com/boyism80/fm/handler"
	"github.com/boyism80/fm/stream"
	"github.com/boyism80/fm/types"
	"github.com/boyism80/fm/util"
)

type ClientActor struct {
	*entity.Character
	Conn           net.Conn
	Buffer         []byte
	messageHandler *handler.MessageHandler
	packetHandler  *handler.PacketHandler
	sendEncryption encrypt.Encryption
	recvEncryption encrypt.Encryption
	stopTimer      scheduler.CancelFunc
}

func generateRandomBytes(size int) []byte {
	b := make([]byte, size)
	_, err := rand.Read(b)
	if err != nil {
		log.Fatal("failed to generate random bytes:", err)
	}
	return b
}

func NewClientActor(conn net.Conn) actor.Actor {
	act := &ClientActor{
		Conn:           conn,
		Buffer:         []byte{},
		messageHandler: handler.NewMessageHandler(),
		packetHandler:  handler.NewPacketHandler(),
		sendEncryption: encrypt.NewEncryption(generateRandomBytes(4), -5),
		recvEncryption: encrypt.NewEncryption(generateRandomBytes(4), 5),
	}

	RegisterCharacterHandlers(act.Character, act.messageHandler)
	RegisterClientHandlers(act, act.messageHandler)
	RegisterPacketHandler(act, act.packetHandler)
	return act
}

func (state *ClientActor) Receive(context actor.Context) {
	state.messageHandler.Handle(context)
}

func (state *ClientActor) Invoke(header int, data []byte) error {
	return state.packetHandler.Handle(header, data)
}

func (c *ClientActor) ID() int64 {
	if c.Character == nil {
		return 0
	}
	return c.Character.Object.ID
}

func (c *ClientActor) Position() types.Vec2 {
	if c.Character == nil {
		return types.Vec2{}
	}
	return c.Character.Object.Position
}

func (c *ClientActor) Type() types.ObjectType {
	return types.ObjectTypePlayer
}

func (c *ClientActor) Name() string {
	if c.Character == nil {
		return ""
	}
	return c.Character.Name
}

func (state *ClientActor) Send(p types.Packet, policy types.SendPolicy) {
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
