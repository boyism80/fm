package entity

import (
	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/common/types"
	lua "github.com/yuin/gopher-lua"
)

var objectBuiltinFuncs = map[string]lua.LGFunction{
	"pid": func(L *lua.LState) int {
		// 1) 첫 번째 인자를 확인: nil이면 nil 반환
		lv := L.Get(1)
		if lv.Type() == lua.LTNil {
			L.Push(lua.LNil)
			return 1
		}

		// 2) userdata인지 검사
		ud, ok := lv.(*lua.LUserData)
		if !ok {
			L.ArgError(1, "Object expected or nil")
			return 0
		}

		// 3) 실제 Go 객체로 언어 변환
		obj, ok := ud.Value.(*Object)
		if !ok {
			L.ArgError(1, "Object expected")
			return 0
		}

		// 4) pid 필드를 숫자로 반환
		newUd := L.NewUserData()
		newUd.Value = &obj.PID
		L.SetMetatable(ud, L.GetTypeMetatable("pid"))
		L.Push(newUd)
		return 1
	},
}

type Object struct {
	PID      *actor.PID
	Position types.Vector2[int16]
}

func (o *Object) LuaTypeName() string {
	return "Object"
}
func (o *Object) LuaBuiltinFuncs() map[string]lua.LGFunction {
	return objectBuiltinFuncs
}
