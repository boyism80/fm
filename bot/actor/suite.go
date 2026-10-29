package actor

import (
	"path/filepath"

	"github.com/boyism80/fm/common/config"
	"github.com/boyism80/fm/core/luax"
	"github.com/boyism80/fm/maple/wz"
	lua "github.com/yuin/gopher-lua"
)

type Suite struct {
	Name      string
	Serial    bool
	Explicit  bool
	TimeoutMs int
}

// Describe loads the suite file without bots and returns the test_suite name and bot count.
func (s *Suite) Describe(cfg *config.Bot, resources *wz.Resources) (string, int, error) {
	a := NewSuiteActor(cfg, resources, nil, s, 0, "", nil)
	err := a.load()
	if a.L != nil {
		defer a.L.Close()
	}
	if err != nil {
		return s.Name, 1, err
	}
	return a.name, a.botCount(), nil
}

// LoadSuites reads registry.lua under dir and returns the suites in registration order.
func LoadSuites(dir string) ([]*Suite, error) {
	L := luax.NewState()
	defer L.Close()
	var suites []*Suite
	L.SetGlobal("register_test", L.NewFunction(func(L *lua.LState) int {
		s := &Suite{Name: L.CheckString(1)}
		if opts := L.OptTable(2, nil); opts != nil {
			s.Serial = lua.LVAsBool(opts.RawGetString("serial"))
			s.Explicit = lua.LVAsBool(opts.RawGetString("explicit"))
			if n, ok := opts.RawGetString("timeout_ms").(lua.LNumber); ok {
				s.TimeoutMs = int(n)
			}
		}
		suites = append(suites, s)
		return 0
	}))
	if err := L.DoFile(filepath.Join(dir, "registry.lua")); err != nil {
		return nil, err
	}
	return suites, nil
}
