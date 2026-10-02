package luax

import (
	"fmt"

	lua "github.com/yuin/gopher-lua"
)

// Upvalues are copied when detached; neither side sees later writes from the other.
type DetachedFunc struct {
	proto    *lua.FunctionProto
	gfn      lua.LGFunction
	upvalues []detachedValue
}

type detachedValue struct {
	value  lua.LValue
	luable Luable
	module string
	script string
	table  *detachedTable
	fn     *DetachedFunc
}

type detachedTable struct {
	keys   []detachedValue
	values []detachedValue
}

type detachScope struct {
	modules map[*lua.LTable]string
	scripts map[*lua.LTable]string
	tables  map[*lua.LTable]*detachedTable
	funcs   map[*lua.LFunction]*DetachedFunc
}

func Detach(L *lua.LState, fn *lua.LFunction) (*DetachedFunc, error) {
	scope := &detachScope{
		modules: make(map[*lua.LTable]string),
		scripts: make(map[*lua.LTable]string),
		tables:  make(map[*lua.LTable]*detachedTable),
		funcs:   make(map[*lua.LFunction]*DetachedFunc),
	}
	if pt, ok := L.GetGlobal("package").(*lua.LTable); ok {
		if loaded, ok := pt.RawGetString("loaded").(*lua.LTable); ok {
			loaded.ForEach(func(k, v lua.LValue) {
				if t, ok := v.(*lua.LTable); ok {
					scope.modules[t] = k.String()
				}
			})
		}
	}
	if mods := modulesTable(L); mods != nil {
		mods.ForEach(func(k, v lua.LValue) {
			if t, ok := v.(*lua.LTable); ok {
				scope.scripts[t] = k.String()
			}
		})
	}
	return scope.function(fn)
}

func (s *detachScope) function(fn *lua.LFunction) (*DetachedFunc, error) {
	if d, ok := s.funcs[fn]; ok {
		return d, nil
	}
	d := &DetachedFunc{}
	s.funcs[fn] = d
	if fn.IsG {
		if len(fn.Upvalues) > 0 {
			return nil, fmt.Errorf("go closure with upvalues cannot be detached")
		}
		d.gfn = fn.GFunction
		return d, nil
	}
	d.proto = fn.Proto
	d.upvalues = make([]detachedValue, len(fn.Upvalues))
	for i, uv := range fn.Upvalues {
		v, err := s.value(uv.Value())
		if err != nil {
			return nil, fmt.Errorf("upvalue %s: %w", fn.Proto.DbgUpvalues[i], err)
		}
		d.upvalues[i] = v
	}
	return d, nil
}

func (s *detachScope) value(v lua.LValue) (detachedValue, error) {
	switch t := v.(type) {
	case *lua.LNilType, lua.LBool, lua.LNumber, lua.LString:
		return detachedValue{value: v}, nil
	case *lua.LUserData:
		luable, ok := t.Value.(Luable)
		if ok == false {
			return detachedValue{}, fmt.Errorf("userdata %T cannot be detached", t.Value)
		}
		return detachedValue{luable: luable}, nil
	case *lua.LFunction:
		fn, err := s.function(t)
		if err != nil {
			return detachedValue{}, err
		}
		return detachedValue{fn: fn}, nil
	case *lua.LTable:
		if name, ok := s.modules[t]; ok {
			return detachedValue{module: name}, nil
		}
		if path, ok := s.scripts[t]; ok {
			return detachedValue{script: path}, nil
		}
		if dt, ok := s.tables[t]; ok {
			return detachedValue{table: dt}, nil
		}
		if t.Metatable != lua.LNil {
			return detachedValue{}, fmt.Errorf("table with metatable cannot be detached")
		}
		dt := &detachedTable{}
		s.tables[t] = dt
		var err error
		t.ForEach(func(k, val lua.LValue) {
			if err != nil {
				return
			}
			dk, kerr := s.value(k)
			if kerr != nil {
				err = kerr
				return
			}
			dv, verr := s.value(val)
			if verr != nil {
				err = verr
				return
			}
			dt.keys = append(dt.keys, dk)
			dt.values = append(dt.values, dv)
		})
		if err != nil {
			return detachedValue{}, err
		}
		return detachedValue{table: dt}, nil
	default:
		return detachedValue{}, fmt.Errorf("%s cannot be detached", v.Type().String())
	}
}

type attachScope struct {
	L      *lua.LState
	tables map[*detachedTable]*lua.LTable
	funcs  map[*DetachedFunc]*lua.LFunction
}

// Run on the goroutine that owns L.
func (d *DetachedFunc) Attach(L *lua.LState) (*lua.LFunction, error) {
	scope := &attachScope{
		L:      L,
		tables: make(map[*detachedTable]*lua.LTable),
		funcs:  make(map[*DetachedFunc]*lua.LFunction),
	}
	return scope.function(d)
}

func (s *attachScope) function(d *DetachedFunc) (*lua.LFunction, error) {
	if fn, ok := s.funcs[d]; ok {
		return fn, nil
	}
	if d.gfn != nil {
		fn := s.L.NewFunction(d.gfn)
		s.funcs[d] = fn
		return fn, nil
	}
	fn := s.L.NewFunctionFromProto(d.proto)
	s.funcs[d] = fn
	for i, dv := range d.upvalues {
		v, err := s.value(dv)
		if err != nil {
			return nil, err
		}
		uv := &lua.Upvalue{}
		uv.SetValue(v)
		fn.Upvalues[i] = uv
	}
	return fn, nil
}

func (s *attachScope) value(dv detachedValue) (lua.LValue, error) {
	switch {
	case dv.luable != nil:
		return NewLuable(s.L, dv.luable), nil
	case dv.fn != nil:
		return s.function(dv.fn)
	case dv.module != "":
		err := s.L.CallByParam(lua.P{Fn: s.L.GetGlobal("require"), NRet: 1, Protect: true}, lua.LString(dv.module))
		if err != nil {
			return nil, err
		}
		v := s.L.Get(-1)
		s.L.Pop(1)
		return v, nil
	case dv.script != "":
		return LoadModule(s.L, dv.script)
	case dv.table != nil:
		if t, ok := s.tables[dv.table]; ok {
			return t, nil
		}
		t := s.L.NewTable()
		s.tables[dv.table] = t
		for i := range dv.table.keys {
			k, err := s.value(dv.table.keys[i])
			if err != nil {
				return nil, err
			}
			v, err := s.value(dv.table.values[i])
			if err != nil {
				return nil, err
			}
			t.RawSet(k, v)
		}
		return t, nil
	case dv.value != nil:
		return dv.value, nil
	default:
		return lua.LNil, nil
	}
}
