package model

import (
	"fmt"
	"log"
	"net"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/handler"
	"github.com/boyism80/fm/msg"
	"github.com/boyism80/fm/packet"
	"github.com/boyism80/fm/stream"
)

type BotActor struct {
	Conn    net.Conn
	handler *handler.MessageHandler
}

func NewBotActor() actor.Actor {
	act := &BotActor{
		handler: handler.NewMessageHandler(),
	}

	handler.RegisterHandler(act, act.handler, onBotConnect)
	handler.RegisterHandler(act, act.handler, onBotSendPacket)
	handler.RegisterHandler(act, act.handler, onBotClose)
	handler.RegisterHandler(act, act.handler, onBotStopping)
	return act
}

func (state *BotActor) Receive(context actor.Context) {
	state.handler.Handle(context)
}

func onBotStopping(bot *BotActor, ctx actor.Context, m *actor.Stopping) {
	log.Println("봇 제거중")
	bot.Conn.Close()
}

func onBotConnect(bot *BotActor, ctx actor.Context, m *msg.BotConnect) {
	serverAddr := fmt.Sprintf("localhost:%d", m.Port)

	conn, err := net.Dial("tcp", serverAddr)
	if err != nil {
		log.Println("연결 실패:", err)
		ctx.Stop(ctx.Self())
		return
	}
	bot.Conn = conn
	log.Println("서버에 연결됨:", serverAddr)

	ctx.Send(ctx.Self(), &msg.BotSendPacket{Packet: &packet.MovePacket{X: 10, Y: 20}})
	ctx.Send(ctx.Self(), &msg.BotClose{})
}

func onBotSendPacket(bot *BotActor, ctx actor.Context, m *msg.BotSendPacket) {
	wr := stream.NewStreamWriter(stream.LittleEndian)
	m.Packet.Serialize(wr)
	bytes := wr.Bytes()
	wr2 := stream.NewStreamWriter(stream.LittleEndian)
	wr2.Write32(1)
	wr2.Write32(int32(len(bytes)))
	wr2.Write(bytes)

	_, err := bot.Conn.Write(wr2.Bytes())
	if err != nil {
		log.Println("데이터 전송 실패:", err)
		return
	}
}

func onBotClose(bot *BotActor, ctx actor.Context, m *msg.BotClose) {
	ctx.Stop(ctx.Self())
}
