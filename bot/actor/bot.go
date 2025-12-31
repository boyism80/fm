package actor

import (
	"fmt"
	"log"
	"net"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/bot/msg"
	"github.com/boyism80/fm/core/context"
	"github.com/boyism80/fm/core/handler"
	"github.com/boyism80/fm/stream"
)

type BotActor struct {
	Conn    net.Conn
	Context *context.ServerContext
	handler *handler.MessageHandler
}

func NewBotActor(ctx actor.Context, serverCtx *context.ServerContext) actor.Actor {
	actor := &BotActor{
		Context: serverCtx,
		handler: handler.NewMessageHandler(),
	}

	handler.RegisterHandler(ctx, actor, actor.handler, onBotConnect)
	handler.RegisterHandler(ctx, actor, actor.handler, onBotSendPacket)
	handler.RegisterHandler(ctx, actor, actor.handler, onBotClose)
	handler.RegisterHandler(ctx, actor, actor.handler, onBotStopping)
	return actor
}

func (state *BotActor) Receive(context actor.Context) {
	state.handler.Handle(context)
}

func onBotStopping(ctx actor.Context, bot *BotActor, m *actor.Stopping) {
	log.Println("Î¥??úÍ±∞Ï§?)
	bot.Conn.Close()
}

func onBotConnect(ctx actor.Context, bot *BotActor, m *msg.BotConnect) {
	serverAddr := fmt.Sprintf("localhost:%d", m.Port)

	conn, err := net.Dial("tcp", serverAddr)
	if err != nil {
		log.Println("?∞Í≤∞ ?§Ìå®:", err)
		ctx.Stop(ctx.Self())
		return
	}
	bot.Conn = conn
	log.Println("?úÎ≤Ñ???∞Í≤∞??", serverAddr)

	ctx.Send(ctx.Self(), &msg.BotClose{})
}

func onBotSendPacket(ctx actor.Context, bot *BotActor, m *msg.BotSendPacket) {
	wr := stream.NewStreamWriter(stream.LittleEndian)
	m.Packet.Serialize(wr)
	bytes := wr.Bytes()
	wr2 := stream.NewStreamWriter(stream.LittleEndian)
	wr2.Write32(1)
	wr2.Write32(int32(len(bytes)))
	wr2.Write(bytes)

	_, err := bot.Conn.Write(wr2.Bytes())
	if err != nil {
		log.Println("?∞Ïù¥???ÑÏÜ° ?§Ìå®:", err)
		return
	}
}

func onBotClose(ctx actor.Context, bot *BotActor, m *msg.BotClose) {
	ctx.Stop(ctx.Self())
}
