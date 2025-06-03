package lua

import (
	"fmt"

	lua "github.com/yuin/gopher-lua"
)

//------------------------------------------------------------------------------
// Luable 인터페이스
//
//   - LuaTypeName(): 이 타입이 Lua 상에서 사용할 메타테이블 이름 반환
//   - LuaBuiltinFuncs(): 이 타입에 바인딩할 Lua 함수 맵 반환
//------------------------------------------------------------------------------

type Luable interface {
	LuaTypeName() string
	LuaBuiltinFuncs() map[string]lua.LGFunction
}

func Register(L *lua.LState, name string, funcs map[string]lua.LGFunction) {
	mt := L.NewTypeMetatable(name)
	L.SetField(mt, "__index", mt)
	L.SetFuncs(mt, funcs)
}

//------------------------------------------------------------------------------
// RegisterLuaType[T] : 부모가 없는 타입을 등록하는 제네릭 함수
//
//   - T 는 Luable 인터페이스를 만족해야 함
//   - 내부에서 `var zero T` 로 제로값(nil)을 만들어, zero.LuaTypeName() 등 호출
//------------------------------------------------------------------------------

func RegisterLuaType[T Luable](L *lua.LState) {
	var zero T
	typeName := zero.LuaTypeName()

	// 1) 메타테이블 생성 (luaL_newmetatable)
	mt := L.NewTypeMetatable(typeName)

	// 2) mt.__index = mt
	L.SetField(mt, "__index", mt)

	// 3) Go 함수를 메타테이블에 등록
	L.SetFuncs(mt, zero.LuaBuiltinFuncs())
}

//------------------------------------------------------------------------------
// RegisterLuaDerivedType[T, B] : 부모가 있는 타입을 등록하는 제네릭 함수
//
//   - T: 자식 타입, B: 부모 타입 (둘 다 Luable 인터페이스 구현)
//   - 부모가 이미 NewTypeMetatable 으로 등록되어 있어야 함
//------------------------------------------------------------------------------

func RegisterLuaDerivedType[T Luable, B Luable](L *lua.LState) {
	var childZero T
	var parentZero B

	childName := childZero.LuaTypeName()
	parentName := parentZero.LuaTypeName()

	// 1) 자식 메타테이블 생성
	childMt := L.NewTypeMetatable(childName)

	// 2) 부모 메타테이블 가져오기
	parentMt := L.GetTypeMetatable(parentName)
	if parentMt == nil {
		panic(fmt.Sprintf("RegisterLuaDerivedType: 부모 메타테이블 '%s' 가 아직 등록되지 않았습니다.", parentName))
	}

	// 3) childMt 의 메타테이블을 parentMt 로 설정 → 상속 체인 연결
	L.SetMetatable(childMt, parentMt)

	// 4) (선택) __parent 필드에 부모 참조 저장
	L.SetField(childMt, "__parent", parentMt)

	// 5) childMt.__index = childMt  (먼저 자식 메서드를 탐색)
	L.SetField(childMt, "__index", childMt)

	// 6) 자식 함수 등록
	L.SetFuncs(childMt, childZero.LuaBuiltinFuncs())
}
