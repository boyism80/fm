package model

import (
	"crypto/rand"
	"log"
	"net"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/asynkron/protoactor-go/scheduler"
	"github.com/boyism80/fm/encrypt"
	"github.com/boyism80/fm/handler"
	"github.com/boyism80/fm/packet"
	"github.com/boyism80/fm/stream"
	"github.com/boyism80/fm/util"
)

const (
	SEND_POLICY_RAW     = 0
	SEND_POLICY_ENCRYPT = 1 << iota
)

type ClientActor struct {
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

func (state *ClientActor) Send(p packet.Packet, policy int) {
	writer := stream.NewStreamWriter(stream.LittleEndian)
	p.Serialize(writer)
	bytes := writer.Bytes()

	log.Println("[S] " + util.ToHexString(bytes))

	if policy == SEND_POLICY_RAW {
		state.Conn.Write(bytes)
		return
	}

	writer = stream.NewStreamWriter(stream.LittleEndian)
	if policy&SEND_POLICY_ENCRYPT != 0 {
		header := state.sendEncryption.GetPacketHeader(len(bytes))
		writer.Write(header)
		bytes = state.sendEncryption.Encrypt(bytes)
	}

	writer.Write(bytes)
	bytes = writer.Bytes()
	state.Conn.Write(bytes)
}
