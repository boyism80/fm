package luax

import (
	"testing"

	lua "github.com/yuin/gopher-lua"
)

func TestDetachAttach(t *testing.T) {
	src := lua.NewState()
	defer src.Close()
	dst := lua.NewState()
	defer dst.Close()

	err := src.DoString(`
		local BASE = 10
		local items = { a = 1, b = { 2, 3 } }
		local function double(n) return n * 2 end
		cb = function(x) return double(BASE + items.a + items.b[2] + x) + where end
	`)
	if err != nil {
		t.Fatal(err)
	}
	detached, err := Detach(src, src.GetGlobal("cb").(*lua.LFunction))
	if err != nil {
		t.Fatal(err)
	}

	dst.SetGlobal("where", lua.LNumber(1000))
	fn, err := detached.Attach(dst)
	if err != nil {
		t.Fatal(err)
	}
	err = dst.CallByParam(lua.P{Fn: fn, NRet: 1, Protect: true}, lua.LNumber(1))
	if err != nil {
		t.Fatal(err)
	}
	if got := dst.Get(-1); got != lua.LNumber(1030) {
		t.Fatalf("got %v, want 1030", got)
	}
}

func TestDetachRejectsCoroutine(t *testing.T) {
	L := lua.NewState()
	defer L.Close()
	err := L.DoString(`
		local co = coroutine.create(function() end)
		cb = function() return co end
	`)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := Detach(L, L.GetGlobal("cb").(*lua.LFunction)); err == nil {
		t.Fatal("expected error")
	}
}
