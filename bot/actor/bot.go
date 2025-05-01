package actor

import (
	"fmt"
	"log"
	"net"

	protoactor "github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/bot/msg"
	"github.com/boyism80/fm/common/context"
	"github.com/boyism80/fm/common/handler"
	"github.com/boyism80/fm/common/stream"
)

type BotActor struct {
	Conn    net.Conn
	Context context.IServerContext
	handler *handler.MessageHandler
}

func NewBotActor(ctx protoactor.Context, serverCtx context.IServerContext) protoactor.Actor {
	act := &BotActor{
		Context: serverCtx,
		handler: handler.NewMessageHandler(),
	}

	handler.RegisterHandler(ctx, act, act.handler, onBotConnect)
	handler.RegisterHandler(ctx, act, act.handler, onBotSendPacket)
	handler.RegisterHandler(ctx, act, act.handler, onBotClose)
	handler.RegisterHandler(ctx, act, act.handler, onBotStopping)
	return act
}

func (state *BotActor) Receive(context protoactor.Context) {
	state.handler.Handle(context)
}

func onBotStopping(ctx protoactor.Context, bot *BotActor, m *protoactor.Stopping) {
	log.Println("봇 제거중")
	bot.Conn.Close()
}

func onBotConnect(ctx protoactor.Context, bot *BotActor, m *msg.BotConnect) {
	serverAddr := fmt.Sprintf("localhost:%d", m.Port)

	conn, err := net.Dial("tcp", serverAddr)
	if err != nil {
		log.Println("연결 실패:", err)
		ctx.Stop(ctx.Self())
		return
	}
	bot.Conn = conn
	log.Println("서버에 연결됨:", serverAddr)

	ctx.Send(ctx.Self(), &msg.BotClose{})
}

func onBotSendPacket(ctx protoactor.Context, bot *BotActor, m *msg.BotSendPacket) {
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

func onBotClose(ctx protoactor.Context, bot *BotActor, m *msg.BotClose) {
	ctx.Stop(ctx.Self())
}
