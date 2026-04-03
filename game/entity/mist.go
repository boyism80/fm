package entity

import (
	"time"

	"github.com/boyism80/fm/core/luax"
	"github.com/boyism80/fm/game/constant"
	"github.com/boyism80/fm/game/wz"
	"github.com/boyism80/fm/protocol/response"
	"github.com/boyism80/fm/types"
	lua "github.com/yuin/gopher-lua"
)

type Mist struct {
	ObjectCore
	OwnerID    uint32
	SkillWz    *wz.Skill
	SkillLevel uint8

	PoisonMist           uint8
	MobMist              bool
	MobSkill             bool
	SkillDelay           uint16
	Bounds               types.Rect[int32]
	ExpiresAt            time.Time
	NextPoisonTickAt     time.Time
	PoisonTickMultiplier float64
}

func MistWorldBounds(pos types.Point[int16], ld *wz.SkillLevelData) types.Rect[int32] {
	px := int32(pos.X)
	py := int32(pos.Y)
	if ld == nil {
		return types.Rect[int32]{Left: px, Top: py, Right: px, Bottom: py}
	}
	x1 := px + ld.LT.X
	x2 := px + ld.RB.X
	y1 := py + ld.LT.Y
	y2 := py + ld.RB.Y
	return types.Rect[int32]{
		Left:   min(x1, x2),
		Top:    min(y1, y2),
		Right:  max(x1, x2),
		Bottom: max(y1, y2),
	}
}

func (mist *Mist) GetObjectType() constant.ObjectType {
	return constant.ObjectTypeMist
}

func (mist *Mist) Is(typ constant.ObjectType) bool {
	return mist.GetObjectType().Has(typ)
}

func (mist *Mist) IsExpired(now time.Time) bool {
	return !mist.ExpiresAt.IsZero() && !now.Before(mist.ExpiresAt)
}

func (mist *Mist) SendSpawnSyncToViewer(viewer *Character) {
	if mist == nil || viewer == nil {
		return
	}
	skillID := uint32(0)
	if mist.SkillWz != nil {
		skillID = mist.SkillWz.ID
	}
	viewer.Send(&response.SpawnMist{
		OID:        mist.OID,
		PoisonMist: mist.PoisonMist,
		MobMist:    mist.MobMist,
		OwnerID:    mist.OwnerID,
		SkillID:    skillID,
		SkillLevel: mist.SkillLevel,
		SkillDelay: mist.SkillDelay,
		Bounds:     mist.Bounds,
		MobSkill:   mist.MobSkill,
	}, types.SEND_POLICY_ENCRYPT)
}

func (mist *Mist) LuaTypeName() string {
	return "LuaMist"
}

func (mist *Mist) String() string {
	return mist.LuaTypeName()
}

func (mist *Mist) Type() lua.LValueType {
	return lua.LTUserData
}

func (mist *Mist) LuaBuiltinFuncs() map[string]lua.LGFunction {
	return map[string]lua.LGFunction{
		"owner_id": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			m, ok := ud.Value.(*Mist)
			if !ok {
				L.ArgError(1, "Mist expected")
				return 0
			}
			if L.GetTop() != 1 {
				L.ArgError(2, "owner_id() is read-only")
				return 0
			}
			L.Push(lua.LNumber(m.OwnerID))
			return 1
		},
		"skill_id": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			m, ok := ud.Value.(*Mist)
			if !ok {
				L.ArgError(1, "Mist expected")
				return 0
			}
			if L.GetTop() != 1 {
				L.ArgError(2, "skill_id() is read-only")
				return 0
			}
			skillID := uint32(0)
			if m.SkillWz != nil {
				skillID = m.SkillWz.ID
			}
			L.Push(lua.LNumber(skillID))
			return 1
		},
		"skill_level": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			m, ok := ud.Value.(*Mist)
			if !ok {
				L.ArgError(1, "Mist expected")
				return 0
			}
			if L.GetTop() != 1 {
				L.ArgError(2, "skill_level() is read-only")
				return 0
			}
			L.Push(lua.LNumber(m.SkillLevel))
			return 1
		},
		"skill": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			m, ok := ud.Value.(*Mist)
			if !ok {
				L.ArgError(1, "Mist expected")
				return 0
			}
			if m.SkillWz == nil {
				L.Push(lua.LNil)
				return 1
			}
			entry := &SkillEntry{
				Wz:         m.SkillWz,
				SkillLevel: int(m.SkillLevel),
				Owner:      nil,
			}
			L.Push(luax.NewLuable(L, entry))
			return 1
		},
		"poison_tick_multiplier": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			m, ok := ud.Value.(*Mist)
			if !ok {
				L.ArgError(1, "Mist expected")
				return 0
			}
			if L.GetTop() != 1 {
				L.ArgError(2, "poison_tick_multiplier() is read-only")
				return 0
			}
			L.Push(lua.LNumber(m.PoisonTickMultiplier))
			return 1
		},
	}
}
