package luax

import (
	"log"
	"math/rand"
	"sync"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/common/handler"
	"github.com/boyism80/fm/game/msg"
	lua "github.com/yuin/gopher-lua"
)

type executor struct {
	handler  *handler.MessageHandler
	workers  []*worker
	lookupMu sync.RWMutex
	lookup   map[*lua.LState]int
	running  bool
}

type worker struct {
	Lua     *lua.LState
	Channel chan any
}

func (state *executor) Receive(context actor.Context) {
	state.handler.Handle(context)
}

func convertToLuaValue(co *lua.LState, args []any) []lua.LValue {
	result := make([]lua.LValue, len(args))
	for i, arg := range args {
		switch v := arg.(type) {
		case nil:
			result[i] = lua.LNil
		case string:
			result[i] = lua.LString(v)
		case int:
			result[i] = lua.LNumber(v)
		case float64:
			result[i] = lua.LNumber(v)
		case bool:
			result[i] = lua.LBool(v)
		case Luable:
			result[i] = NewLuable(co, v)
		}
	}
	return result
}

func (state *executor) run(ctx actor.Context, worker *worker, i int) {

	for state.running {
		task := <-worker.Channel
		switch m := task.(type) {
		case *msg.LuaRun:
			if m.Params == nil {
				m.Params = []any{}
			}
			co, err := NewThread(worker.Lua, m.FileName)
			if err != nil {
				log.Println(err)
				return
			}
			resumeState, err := call(worker.Lua, co, m.FuncName, convertToLuaValue(co, m.Params)...)
			if err != nil {
				log.Println(err)
			}

			switch resumeState {
			case lua.ResumeYield:
				state.lookupMu.Lock()
				state.lookup[co] = i
				state.lookupMu.Unlock()
				if m.PID != nil {
					ctx.Send(m.PID, &msg.LuaYield{
						Lua: co,
					})
				}
			case lua.ResumeError:
				log.Println(err)
			}

		case *msg.LuaResume:
			co := m.Lua
			if m.Params == nil {
				m.Params = []any{}
			}
			resumeState, err := resume(worker.Lua, co, convertToLuaValue(co, m.Params)...)
			if err != nil {
				log.Println(err)
			}

			switch resumeState {
			case lua.ResumeError:
				log.Println(err)
				state.lookupMu.Lock()
				delete(state.lookup, co)
				state.lookupMu.Unlock()
			case lua.ResumeYield:
				if m.PID != nil {
					ctx.Send(m.PID, &msg.LuaYield{
						Lua: co,
					})
				}
			case lua.ResumeOK:
				state.lookupMu.Lock()
				delete(state.lookup, co)
				state.lookupMu.Unlock()
			}
		}
	}
}

func newLuaExecutor(ctx actor.Context, n int) actor.Actor {

	actor := &executor{
		handler:  handler.NewMessageHandler(),
		workers:  []*worker{},
		lookup:   map[*lua.LState]int{},
		lookupMu: sync.RWMutex{},
		running:  true,
	}

	for i := range n {
		root := NewState()
		worker := &worker{
			Lua:     root,
			Channel: make(chan any),
		}
		actor.workers = append(actor.workers, worker)
		go actor.run(ctx, worker, i)
	}

	handler.RegisterHandler(ctx, actor, actor.handler, onLuaStopping)
	handler.RegisterHandler(ctx, actor, actor.handler, onLuaRun)
	handler.RegisterHandler(ctx, actor, actor.handler, onLuaResume)

	return actor
}

func onLuaStopping(ctx actor.Context, state *executor, m *actor.Stopping) {
	state.running = false
	for _, worker := range state.workers {
		close(worker.Channel)
	}
}

func onLuaRun(ctx actor.Context, state *executor, m *msg.LuaRun) {

	i := rand.Intn(len(state.workers))
	worker := state.workers[i]
	worker.Channel <- m
}

func onLuaResume(ctx actor.Context, state *executor, m *msg.LuaResume) {
	state.lookupMu.RLock()
	i, ok := state.lookup[m.Lua]
	state.lookupMu.RUnlock()
	if !ok {
		return
	}
	worker := state.workers[i]
	worker.Channel <- m
}
