package handler

import (
	"fmt"

	"github.com/asynkron/protoactor-go/actor"
)

type CommandHandler struct {
	handlers map[string]func(ctx actor.Context, params ...string)
}

func NewCommandHandler() *CommandHandler {
	return &CommandHandler{
		handlers: map[string]func(ctx actor.Context, params ...string){},
	}
}

func (state *CommandHandler) Register(cmd string, handler func(ctx actor.Context, params ...string)) {
	state.handlers[cmd] = handler
}

func RegisterCommandHandler[T any](cmd string, ctx actor.Context, target *T, h *CommandHandler, fn func(actor.Context, *T, ...string)) {
	h.Register(cmd, func(ctx actor.Context, params ...string) {
		fn(ctx, target, params...)
	})
}

func (state *CommandHandler) Handle(ctx actor.Context, params ...string) error {
	if len(params) < 1 {
		return fmt.Errorf("no command provided")
	}

	cmd := params[0]
	handler, exists := state.handlers[cmd]
	if !exists {
		return fmt.Errorf("command not found: %s", cmd)
	}

	handler(ctx, params[1:]...)
	return nil
}
