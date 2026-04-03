package entity

import (
	"github.com/boyism80/fm/game/constant"
	"github.com/boyism80/fm/game/wz"
	"github.com/boyism80/fm/protocol/response"
	"github.com/boyism80/fm/types"
	lua "github.com/yuin/gopher-lua"
)

type Npc struct {
	ObjectCore
	Wz *wz.NpcSpawn
}

func (n *Npc) GetObjectType() constant.ObjectType {
	return constant.ObjectTypeNpc
}

func (n *Npc) Is(typ constant.ObjectType) bool {
	return n.GetObjectType().Has(typ)
}

func (n *Npc) SendSpawnSyncToViewer(viewer *Character) {
	if n == nil || viewer == nil {
		return
	}
	npcDTO := n.ToDTO()
	viewer.Send(&response.SpawnNpc{
		NPC:     npcDTO,
		Visible: true,
	}, types.SEND_POLICY_ENCRYPT)
	viewer.Send(&response.NpcControl{
		NPC:     npcDTO,
		MiniMap: true,
	}, types.SEND_POLICY_ENCRYPT)
}

func (n *Npc) LuaTypeName() string {
	return "LuaNpc"
}

func (n *Npc) LuaBuiltinFuncs() map[string]lua.LGFunction {
	return map[string]lua.LGFunction{
		"id": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			npc, ok := ud.Value.(*Npc)
			if !ok {
				L.ArgError(1, "Npc expected")
				return 0
			}
			if npc.Wz != nil && npc.Wz.BaseSpawn != nil {
				L.Push(lua.LNumber(npc.Wz.BaseSpawn.ID))
			} else {
				L.Push(lua.LNumber(0))
			}
			return 1
		},
		"wz": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			npc, ok := ud.Value.(*Npc)
			if !ok {
				L.ArgError(1, "Npc expected")
				return 0
			}

			tbl := L.NewTable()
			if npc.Wz != nil && npc.Wz.BaseSpawn != nil {
				tbl.RawSetString("id", lua.LNumber(npc.Wz.BaseSpawn.ID))
				if npc.Wz.BaseSpawn.Position.X != 0 || npc.Wz.BaseSpawn.Position.Y != 0 {
					posTbl := L.NewTable()
					posTbl.RawSetString("x", lua.LNumber(npc.Wz.BaseSpawn.Position.X))
					posTbl.RawSetString("y", lua.LNumber(npc.Wz.BaseSpawn.Position.Y))
					tbl.RawSetString("position", posTbl)
				}
				tbl.RawSetString("collision_y", lua.LNumber(npc.Wz.BaseSpawn.CollisionY))
			}
			L.Push(tbl)
			return 1
		},
	}
}

func (n *Npc) String() string {
	return n.LuaTypeName()
}

func (n *Npc) Type() lua.LValueType {
	return lua.LTUserData
}
