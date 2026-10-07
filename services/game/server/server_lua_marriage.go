package server

import (
	"context"
	"log"

	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/core/async"
	"github.com/boyism80/fm/core/luax"
	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
	"github.com/boyism80/fm/services/game/entity"
	lua "github.com/yuin/gopher-lua"
)

func registerMarriageLuaFuncs(gs *GameServer, luaState *lua.LState) {
	luax.RegisterFunc(luaState, "finish_wedding", func(L *lua.LState) int {
		cfg, ok := luax.GetConfiguration(L)
		if !ok || cfg.ActorContext == nil {
			L.Push(lua.LFalse)
			return 1
		}
		req := &internal.FinishWeddingRequest{
			WorldId:    gs.config.WorldId,
			MarriageId: uint32(L.CheckNumber(1)),
		}
		promise := async.ThenRPC(
			async.NewPromise(cfg.ActorContext, core.InternalRPCPerStepTimeout),
			func(c context.Context) (*internal.MarriageReply, error) {
				return gs.internalClient.FinishWedding(c, req)
			},
			func(reply *internal.MarriageReply) error {
				return nil
			},
		)
		return entity.LuaYieldPromise(L, gs, promise, func(result interface{}, err error) []lua.LValue {
			if err != nil {
				log.Printf("finish_wedding marriage=%d: %v", req.GetMarriageId(), err)
				return []lua.LValue{lua.LFalse}
			}
			return []lua.LValue{lua.LBool(result.(*internal.MarriageReply).GetResult() == internal.MarriageResult_MARRIAGE_RESULT_OK)}
		})
	})
}
