package actor

import (
	"log"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/common/handler"
)

func RegisterGameClientCommandHandler(ctx actor.Context, client *GameClientActor, h *handler.CommandHandler) {
	handler.RegisterCommandHandler("아이템", ctx, client, h, onCreateItem)
}

func onCreateItem(ctx actor.Context, client *GameClientActor, params ...string) {
	log.Println(params)
}
