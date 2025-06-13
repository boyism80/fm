package luax

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/asynkron/protoactor-go/scheduler"
	"github.com/boyism80/fm/game/msg"
	lua "github.com/yuin/gopher-lua"
)

var (
	onCreateHooks   []func(*lua.LState)
	onCreateHooksMu sync.Mutex
	compileMu       sync.Mutex
	compiledFuncs   = make(map[string]*lua.LFunction)
	useCache        = os.Getenv("GO_ENV") != "development"
	ExecutorPID     *actor.PID
	RootContext     *actor.RootContext
	timerScheduler  *scheduler.TimerScheduler
)

func init() {
	env, ok := os.LookupEnv("GO_ENV")
	if !ok || env == "" {
		env = "development"
	}
	useCache = env != "development"
}

type Luable interface {
	LuaTypeName() string
	LuaBuiltinFuncs() map[string]lua.LGFunction
}

func NewState() *lua.LState {
	luaState := lua.NewState()
	onCreateHooksMu.Lock()
	for _, hook := range onCreateHooks {
		hook(luaState)
	}
	onCreateHooksMu.Unlock()
	return luaState
}

func preloadScript(root *lua.LState, path string) (*lua.LFunction, error) {
	compileMu.Lock()
	defer compileMu.Unlock()

	if fn, ok := compiledFuncs[path]; ok && useCache {
		return fn, nil
	}

	fn, err := root.LoadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to compile %s: %w", path, err)
	}
	compiledFuncs[path] = fn
	return fn, nil
}

func NewThread(root *lua.LState, path string) (*lua.LState, error) {
	fn, err := preloadScript(root, path)
	if err != nil {
		return nil, err
	}

	co, _ := root.NewThread()
	co.Push(fn)
	if err := co.PCall(0, lua.MultRet, nil); err != nil {
		return nil, fmt.Errorf("script runtime error: %w", err)
	}

	return co, nil
}

func call(root *lua.LState, co *lua.LState, funcName string, args ...lua.LValue) (lua.ResumeState, error) {
	fn := co.GetGlobal(funcName)
	if fn.Type() != lua.LTFunction {
		return lua.ResumeYield, fmt.Errorf("on_start is not a function")
	}

	resumeState, err, _ := root.Resume(co, fn.(*lua.LFunction), args...)
	if err != nil {
		return lua.ResumeYield, err
	}
	return resumeState, nil
}

func Call(ctx actor.Context, sender *actor.PID, fileName string, funcName string, args ...any) {

	ctx.Send(ExecutorPID, &msg.LuaRun{
		PID:      sender,
		FileName: fileName,
		FuncName: funcName,
		Params:   args,
	})
}

func resume(root *lua.LState, co *lua.LState, args ...lua.LValue) (lua.ResumeState, error) {
	resumeState, err, _ := root.Resume(co, nil, args...)
	if err != nil {
		return lua.ResumeYield, err
	}
	return resumeState, nil
}

func Resume(ctx actor.Context, sender *actor.PID, co *lua.LState, args ...any) {
	ctx.Send(ExecutorPID, &msg.LuaResume{
		PID:    sender,
		Lua:    co,
		Params: args,
	})
}

func RegisterOnCreateHook(fn func(*lua.LState)) {
	onCreateHooksMu.Lock()
	defer onCreateHooksMu.Unlock()
	onCreateHooks = append(onCreateHooks, fn)
}

func RegisterFunc(L *lua.LState, name string, fn lua.LGFunction) {
	L.SetGlobal(name, L.NewFunction(fn))
}

func RegisterLuaType[T Luable](L *lua.LState) {
	var zero T
	typeName := zero.LuaTypeName()

	mt := L.NewTypeMetatable(typeName)

	L.SetField(mt, "__index", mt)

	L.SetFuncs(mt, zero.LuaBuiltinFuncs())
}

func RegisterLuaDerivedType[T Luable, B Luable](L *lua.LState) {
	var childZero T
	var parentZero B

	childName := childZero.LuaTypeName()
	parentName := parentZero.LuaTypeName()

	childMt := L.NewTypeMetatable(childName)

	parentMt := L.GetTypeMetatable(parentName)
	if parentMt == nil {
		panic(fmt.Sprintf("RegisterLuaDerivedType: 부모 메타테이블 '%s' 가 아직 등록되지 않았습니다.", parentName))
	}

	L.SetMetatable(childMt, parentMt)
	L.SetField(childMt, "__parent", parentMt)
	L.SetField(childMt, "__index", childMt)
	L.SetFuncs(childMt, childZero.LuaBuiltinFuncs())
}

func NewLuable(L *lua.LState, obj Luable) *lua.LUserData {
	ud := L.NewUserData()
	ud.Value = obj
	L.SetMetatable(ud, L.GetTypeMetatable(obj.LuaTypeName()))
	return ud
}

func Setup(ctx actor.Context) {
	RootContext = ctx.ActorSystem().Root

	props := actor.PropsFromProducer(func() actor.Actor { return newLuaExecutor(ctx, 10) })
	ExecutorPID = ctx.Spawn(props)
	timerScheduler = scheduler.NewTimerScheduler(ctx)
}

func SendAfter(duration time.Duration, message interface{}) {
	timerScheduler.SendOnce(duration, ExecutorPID, message)
}
