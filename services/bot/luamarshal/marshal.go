package luamarshal

import (
	"fmt"
	"reflect"
	"strings"
	"unicode"

	lua "github.com/yuin/gopher-lua"
)

type Marshal struct {
	fields map[reflect.Type]map[string]int
}

func New() *Marshal {
	return &Marshal{fields: make(map[reflect.Type]map[string]int)}
}

func (m *Marshal) Name(goName string) string {
	runes := []rune(goName)
	var sb strings.Builder
	for i, r := range runes {
		if i > 0 && unicode.IsUpper(r) {
			prev := runes[i-1]
			nextLower := i+1 < len(runes) && unicode.IsLower(runes[i+1])
			if unicode.IsLower(prev) || unicode.IsDigit(prev) || (unicode.IsUpper(prev) && nextLower) {
				sb.WriteByte('_')
			}
		}
		sb.WriteRune(unicode.ToLower(r))
	}
	return sb.String()
}

func (m *Marshal) ToLua(L *lua.LState, v any) lua.LValue {
	return m.toLua(L, reflect.ValueOf(v))
}

func (m *Marshal) toLua(L *lua.LState, v reflect.Value) lua.LValue {
	switch v.Kind() {
	case reflect.Invalid:
		return lua.LNil
	case reflect.Pointer, reflect.Interface:
		if v.IsNil() {
			return lua.LNil
		}
		return m.toLua(L, v.Elem())
	case reflect.Bool:
		return lua.LBool(v.Bool())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return lua.LNumber(v.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return lua.LNumber(v.Uint())
	case reflect.Float32, reflect.Float64:
		return lua.LNumber(v.Float())
	case reflect.String:
		return lua.LString(v.String())
	case reflect.Slice, reflect.Array:
		tbl := L.NewTable()
		for i := 0; i < v.Len(); i++ {
			tbl.RawSetInt(i+1, m.toLua(L, v.Index(i)))
		}
		return tbl
	case reflect.Map:
		tbl := L.NewTable()
		iter := v.MapRange()
		for iter.Next() {
			tbl.RawSet(m.toLua(L, iter.Key()), m.toLua(L, iter.Value()))
		}
		return tbl
	case reflect.Struct:
		tbl := L.NewTable()
		t := v.Type()
		for i := 0; i < t.NumField(); i++ {
			if t.Field(i).IsExported() == false {
				continue
			}
			tbl.RawSetString(m.Name(t.Field(i).Name), m.toLua(L, v.Field(i)))
		}
		return tbl
	default:
		return lua.LNil
	}
}

func (m *Marshal) FromLua(lv lua.LValue, out any) error {
	v := reflect.ValueOf(out)
	if v.Kind() != reflect.Pointer || v.IsNil() {
		return fmt.Errorf("FromLua needs a non-nil pointer, got %T", out)
	}
	return m.fromLua(lv, v.Elem())
}

func (m *Marshal) fromLua(lv lua.LValue, v reflect.Value) error {
	if lv == lua.LNil {
		return nil
	}
	switch v.Kind() {
	case reflect.Pointer:
		if v.IsNil() {
			v.Set(reflect.New(v.Type().Elem()))
		}
		return m.fromLua(lv, v.Elem())
	case reflect.Bool:
		v.SetBool(lua.LVAsBool(lv))
		return nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		n, ok := lv.(lua.LNumber)
		if ok == false {
			return fmt.Errorf("want number, got %s", lv.Type())
		}
		v.SetInt(int64(n))
		return nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		n, ok := lv.(lua.LNumber)
		if ok == false {
			return fmt.Errorf("want number, got %s", lv.Type())
		}
		v.SetUint(uint64(n))
		return nil
	case reflect.Float32, reflect.Float64:
		n, ok := lv.(lua.LNumber)
		if ok == false {
			return fmt.Errorf("want number, got %s", lv.Type())
		}
		v.SetFloat(float64(n))
		return nil
	case reflect.String:
		s, ok := lv.(lua.LString)
		if ok == false {
			return fmt.Errorf("want string, got %s", lv.Type())
		}
		v.SetString(string(s))
		return nil
	case reflect.Slice:
		tbl, ok := lv.(*lua.LTable)
		if ok == false {
			return fmt.Errorf("want table, got %s", lv.Type())
		}
		n := tbl.Len()
		s := reflect.MakeSlice(v.Type(), n, n)
		for i := 0; i < n; i++ {
			if err := m.fromLua(tbl.RawGetInt(i+1), s.Index(i)); err != nil {
				return fmt.Errorf("[%d]: %w", i+1, err)
			}
		}
		v.Set(s)
		return nil
	case reflect.Struct:
		tbl, ok := lv.(*lua.LTable)
		if ok == false {
			return fmt.Errorf("want table, got %s", lv.Type())
		}
		fields := m.fieldIndex(v.Type())
		var err error
		tbl.ForEach(func(key, value lua.LValue) {
			if err != nil {
				return
			}
			name := key.String()
			i, ok := fields[name]
			if ok == false {
				err = fmt.Errorf("%s has no field %q", v.Type().Name(), name)
				return
			}
			if fieldErr := m.fromLua(value, v.Field(i)); fieldErr != nil {
				err = fmt.Errorf("%s: %w", name, fieldErr)
			}
		})
		return err
	default:
		return fmt.Errorf("unsupported kind %s", v.Kind())
	}
}

func (m *Marshal) fieldIndex(t reflect.Type) map[string]int {
	if fields, ok := m.fields[t]; ok {
		return fields
	}
	fields := make(map[string]int)
	for i := 0; i < t.NumField(); i++ {
		if t.Field(i).IsExported() {
			fields[m.Name(t.Field(i).Name)] = i
		}
	}
	m.fields[t] = fields
	return fields
}
