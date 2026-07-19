package actor

import (
	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/core/async"
	"github.com/boyism80/fm/core/luax"
	lua "github.com/yuin/gopher-lua"
)

func (a *MapActor) RunScriptCall(ctx actor.Context, scriptPath string, funcName string, args []interface{}) *async.Promise {
	if a.Map == nil {
		return async.NewPromise(ctx, 0).Then(func(_ interface{}) (interface{}, error) {
			return []lua.LValue{lua.LBool(false), lua.LNil, lua.LString("map not found")}, nil
		})
	}
	return a.Map.RunScript(ctx, scriptPath, funcName, args).Then(func(v interface{}) (interface{}, error) {
		vals := luax.ResultValues(v)
		result := lua.LNil
		if len(vals) > 0 && vals[0] != nil {
			result = vals[0]
		}
		return []lua.LValue{lua.LBool(true), result, lua.LNil}, nil
	})
}
