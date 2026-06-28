package entity

import (
	"fmt"
	"strings"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/core/async"
	"github.com/boyism80/fm/core/luax"
	lua "github.com/yuin/gopher-lua"
)

func validateScriptPath(path string) error {
	if path == "" {
		return fmt.Errorf("empty script path")
	}
	path = strings.ReplaceAll(path, "\\", "/")
	if !strings.HasPrefix(path, "script/") {
		return fmt.Errorf("script path must start with script/")
	}
	if strings.Contains(path, "..") {
		return fmt.Errorf("invalid script path")
	}
	return nil
}

func (m *Map) RunScript(ctx actor.Context, scriptPath string, funcName string, args []interface{}) *async.Promise {
	promise := async.NewPromise(nil, 0)
	if m == nil {
		promise.SetError(fmt.Errorf("map is nil"))
		return promise
	}
	if err := validateScriptPath(scriptPath); err != nil {
		promise.SetError(err)
		return promise
	}
	if funcName == "" {
		promise.SetError(fmt.Errorf("empty function name"))
		return promise
	}
	root := m.EnsureLuaRoot(ctx)
	if root == nil {
		promise.SetError(fmt.Errorf("lua root not found"))
		return promise
	}
	thread, err := luax.NewThread(root, scriptPath)
	if err != nil {
		promise.SetError(err)
		return promise
	}
	luax.SetConfiguration(thread, luax.Configuration{
		MapActorPID: m.GetActorPID(),
	})
	fn := thread.GetGlobal(funcName)
	if fn.Type() != lua.LTFunction {
		luax.Close(thread)
		promise.SetError(fmt.Errorf("function %q not found", funcName))
		return promise
	}
	callArgs := make([]interface{}, 0, len(args)+1)
	callArgs = append(callArgs, m)
	callArgs = append(callArgs, args...)
	return luax.CallAsync(root, thread, funcName, callArgs...)
}
