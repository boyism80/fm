package entity

import (
	"fmt"
	"strings"

	"github.com/asynkron/protoactor-go/actor"
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

func (m *Map) RunScript(ctx actor.Context, scriptPath string, funcName string, args []interface{}) (lua.LValue, error) {
	if m == nil {
		return nil, fmt.Errorf("map is nil")
	}
	if err := validateScriptPath(scriptPath); err != nil {
		return nil, err
	}
	if funcName == "" {
		return nil, fmt.Errorf("empty function name")
	}
	root := m.EnsureLuaRoot(ctx)
	if root == nil {
		return nil, fmt.Errorf("lua root not found")
	}
	thread, err := luax.NewThread(root, scriptPath)
	if err != nil {
		return nil, err
	}
	luax.SetConfiguration(thread, luax.Configuration{
		MapActorPID: m.GetActorPID(),
	})
	fn := thread.GetGlobal(funcName)
	if fn.Type() != lua.LTFunction {
		return nil, fmt.Errorf("function %q not found", funcName)
	}
	callArgs := make([]interface{}, 0, len(args)+1)
	callArgs = append(callArgs, m)
	callArgs = append(callArgs, args...)
	return luax.Call(thread, funcName, callArgs...)
}
