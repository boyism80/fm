package server

import (
	"log"
	"sort"
	"time"

	"github.com/boyism80/fm/core/luax"
	g_actor "github.com/boyism80/fm/game/actor"
	"github.com/boyism80/fm/game/constant"
	"github.com/boyism80/fm/game/entity"
	"github.com/boyism80/fm/game/wz"
	lua "github.com/yuin/gopher-lua"
)

func registerSkillConstants(luaState *lua.LState) {
	skillTable := luaState.NewTable()
	for key, skillID := range constant.AllSkillConstants() {
		skillTable.RawSetString(key, lua.LNumber(skillID))
	}
	luaState.SetGlobal("Skill", skillTable)
}

func registerBuffFlagAndMobStatus(luaState *lua.LState) {
	buffFlagTable := luaState.NewTable()
	for name, bf := range constant.AllBuffFlags() {
		entry := luaState.NewTable()
		entry.RawSetString("mask", lua.LNumber(bf.Mask))
		entry.RawSetString("position", lua.LNumber(bf.Position))
		buffFlagTable.RawSetString(name, entry)
	}
	luaState.SetGlobal("BuffFlag", buffFlagTable)
	mobStatusTable := luaState.NewTable()
	for name, st := range constant.AllMobStatuses() {
		entry := luaState.NewTable()
		entry.RawSetString("mask", lua.LNumber(st))
		mobStatusTable.RawSetString(name, entry)
	}
	luaState.SetGlobal("MobStatus", mobStatusTable)
}

func registerWeaponTypeAndConsumeType(luaState *lua.LState) {
	weaponTypeTable := luaState.NewTable()
	for name, wt := range constant.AllWeaponTypes() {
		weaponTypeTable.RawSetString(name, lua.LNumber(wt))
	}
	luaState.SetGlobal("WeaponType", weaponTypeTable)
	consumeTypeTable := luaState.NewTable()
	for name, ct := range constant.AllConsumeTypes() {
		consumeTypeTable.RawSetString(name, lua.LNumber(ct))
	}
	luaState.SetGlobal("ConsumeType", consumeTypeTable)
}

func registerEquipmentPartConstants(luaState *lua.LState) {
	t := luaState.NewTable()
	t.RawSetString("Cap", lua.LNumber(constant.EQUIPMENT_PARTS_CAP))
	t.RawSetString("Face", lua.LNumber(constant.EQUIPMENT_PARTS_FACE))
	t.RawSetString("Eye", lua.LNumber(constant.EQUIPMENT_PARTS_EYE))
	t.RawSetString("Ear", lua.LNumber(constant.EQUIPMENT_PARTS_EAR))
	t.RawSetString("Top", lua.LNumber(constant.EQUIPMENT_PARTS_TOP))
	t.RawSetString("Pants", lua.LNumber(constant.EQUIPMENT_PARTS_PANTS))
	t.RawSetString("Shoes", lua.LNumber(constant.EQUIPMENT_PARTS_SHOES))
	t.RawSetString("Glove", lua.LNumber(constant.EQUIPMENT_PARTS_GLOVE))
	t.RawSetString("Cape", lua.LNumber(constant.EQUIPMENT_PARTS_CAPE))
	t.RawSetString("Shield", lua.LNumber(constant.EQUIPMENT_PARTS_SHIELD))
	t.RawSetString("Weapon", lua.LNumber(constant.EQUIPMENT_PARTS_WEAPON))
	t.RawSetString("Ring", lua.LNumber(constant.EQUIPMENT_PARTS_RING))
	luaState.SetGlobal("EquipmentPart", t)
}

func registerInventoryTypeConstants(luaState *lua.LState) {
	t := luaState.NewTable()
	t.RawSetString("Equipment", lua.LNumber(constant.INVENTORY_TYPE_EQUIPMENT))
	t.RawSetString("Use", lua.LNumber(constant.INVENTORY_TYPE_CONSUME))
	t.RawSetString("Consume", lua.LNumber(constant.INVENTORY_TYPE_CONSUME))
	t.RawSetString("Installation", lua.LNumber(constant.INVENTORY_TYPE_INSTALLATION))
	t.RawSetString("Etc", lua.LNumber(constant.INVENTORY_TYPE_ETC))
	t.RawSetString("Cash", lua.LNumber(constant.INVENTORY_TYPE_CASH))
	luaState.SetGlobal("InventoryType", t)
}

func registerStatConstants(luaState *lua.LState) {
	t := luaState.NewTable()
	t.RawSetString("Level", lua.LNumber(constant.STAT_LEVEL))
	t.RawSetString("Exp", lua.LNumber(constant.STAT_EXP))
	t.RawSetString("Class", lua.LNumber(constant.STAT_CLASS))
	t.RawSetString("Str", lua.LNumber(constant.STAT_STR))
	t.RawSetString("Dex", lua.LNumber(constant.STAT_DEX))
	t.RawSetString("Int", lua.LNumber(constant.STAT_INT))
	t.RawSetString("Luk", lua.LNumber(constant.STAT_LUK))
	t.RawSetString("Hp", lua.LNumber(constant.STAT_HP))
	t.RawSetString("MaxHp", lua.LNumber(constant.STAT_MAX_HP))
	t.RawSetString("Mp", lua.LNumber(constant.STAT_MP))
	t.RawSetString("MaxMp", lua.LNumber(constant.STAT_MAX_MP))
	t.RawSetString("AvailableAP", lua.LNumber(constant.STAT_AVAILABLE_AP))
	t.RawSetString("AvailableSP", lua.LNumber(constant.STAT_AVAILABLE_SP))
	t.RawSetString("Fame", lua.LNumber(constant.STAT_FAME))
	t.RawSetString("Meso", lua.LNumber(constant.STAT_MESO))
	luaState.SetGlobal("STAT", t)
}

func registerClassConstants(luaState *lua.LState) {
	t := luaState.NewTable()
	for name, code := range constant.AllClassConstants() {
		t.RawSetString(name, lua.LNumber(code))
	}
	luaState.SetGlobal("Class", t)
}

func registerStanceConstants(luaState *lua.LState) {
	t := luaState.NewTable()
	for name, value := range constant.AllStanceConstants() {
		t.RawSetString(name, lua.LNumber(value))
	}
	luaState.SetGlobal("Stance", t)
}

func registerObjectTypeConstants(luaState *lua.LState) {
	t := luaState.NewTable()
	for name, value := range constant.AllObjectTypeConstants() {
		t.RawSetString(name, lua.LNumber(value))
	}
	luaState.SetGlobal("ObjectType", t)
}

func skillToLuaWzTable(luaState *lua.LState, skill *wz.Skill) *lua.LTable {
	skillTable := luaState.NewTable()
	skillTable.RawSetString("id", lua.LNumber(skill.ID))
	skillTable.RawSetString("skill_id", lua.LNumber(skill.ID))
	skillTable.RawSetString("max_level", lua.LNumber(skill.MaxLevel))
	skillTable.RawSetString("master_level", lua.LNumber(skill.MasterLevel))
	skillTable.RawSetString("invisible", lua.LBool(skill.Invisible))
	return skillTable
}

func registerGameLuaState(gs *GameServer, luaState *lua.LState) {
	luax.RegisterLuaType[*entity.Object](luaState)
	luax.RegisterLuaDerivedType[*entity.Drop, *entity.Object](luaState)
	luax.RegisterLuaDerivedType[*entity.Meso, *entity.Drop](luaState)
	luax.RegisterLuaDerivedType[*entity.Life, *entity.Object](luaState)
	luax.RegisterLuaDerivedType[*entity.Character, *entity.Life](luaState)
	luax.RegisterLuaDerivedType[*entity.Mob, *entity.Life](luaState)
	luax.RegisterLuaDerivedType[*entity.Npc, *entity.Object](luaState)
	luax.RegisterLuaType[*entity.Map](luaState)
	luax.RegisterLuaType[*entity.SkillEntry](luaState)
	luax.RegisterLuaType[*entity.ItemCore](luaState)
	luax.RegisterLuaDerivedType[*entity.EquipmentCore, *entity.ItemCore](luaState)
	luax.RegisterLuaDerivedType[*entity.Weapon, *entity.EquipmentCore](luaState)
	luax.RegisterLuaDerivedType[*entity.Shield, *entity.EquipmentCore](luaState)
	luax.RegisterLuaDerivedType[*entity.Cap, *entity.EquipmentCore](luaState)
	luax.RegisterLuaDerivedType[*entity.Face, *entity.EquipmentCore](luaState)
	luax.RegisterLuaDerivedType[*entity.Accessory, *entity.EquipmentCore](luaState)
	luax.RegisterLuaDerivedType[*entity.Top, *entity.EquipmentCore](luaState)
	luax.RegisterLuaDerivedType[*entity.Pants, *entity.EquipmentCore](luaState)
	luax.RegisterLuaDerivedType[*entity.Shoes, *entity.EquipmentCore](luaState)
	luax.RegisterLuaDerivedType[*entity.Glove, *entity.EquipmentCore](luaState)
	luax.RegisterLuaDerivedType[*entity.Cape, *entity.EquipmentCore](luaState)
	luax.RegisterLuaDerivedType[*entity.RingEquip, *entity.EquipmentCore](luaState)
	luax.RegisterLuaType[*wz.ItemWzCore](luaState)
	luax.RegisterLuaDerivedType[*wz.ItemWzEquipment, *wz.ItemWzCore](luaState)
	luax.RegisterLuaDerivedType[*wz.ItemWzWeapon, *wz.ItemWzEquipment](luaState)
	luax.RegisterLuaDerivedType[*wz.ItemWzArmor, *wz.ItemWzEquipment](luaState)
	luax.RegisterLuaDerivedType[*wz.ItemWzConsume, *wz.ItemWzCore](luaState)
	luax.RegisterLuaDerivedType[*wz.ItemWzPet, *wz.ItemWzCore](luaState)
	luax.RegisterLuaDerivedType[*wz.ItemWzGeneralItem, *wz.ItemWzCore](luaState)
	luax.RegisterLuaDerivedType[*wz.ItemWzCashItem, *wz.ItemWzCore](luaState)
	luax.RegisterLuaDerivedType[*wz.ItemWzInstallation, *wz.ItemWzCore](luaState)
	luax.RegisterLuaDerivedType[*entity.Consume, *entity.ItemCore](luaState)
	luax.RegisterLuaDerivedType[*entity.CashItem, *entity.ItemCore](luaState)
	luax.RegisterLuaDerivedType[*entity.GeneralItem, *entity.ItemCore](luaState)
	luax.RegisterLuaDerivedType[*entity.Installation, *entity.ItemCore](luaState)
	luax.RegisterLuaDerivedType[*entity.Pet, *entity.ItemCore](luaState)

	registerBuffFlagAndMobStatus(luaState)
	registerWeaponTypeAndConsumeType(luaState)
	registerSkillConstants(luaState)
	registerEquipmentPartConstants(luaState)
	registerInventoryTypeConstants(luaState)
	registerStatConstants(luaState)
	registerClassConstants(luaState)
	registerStanceConstants(luaState)
	registerObjectTypeConstants(luaState)

	luax.RegisterFunc(luaState, "class_learnable_skill_wzs", func(L *lua.LState) int {
		class := uint16(L.CheckInt(1))
		result := L.NewTable()
		if gs.resources == nil {
			L.Push(result)
			return 1
		}

		classIDs := getClassAdvancementClasses(class)
		skillByID := make(map[uint32]*wz.Skill)
		for _, classID := range classIDs {
			skillIDStart := uint32(classID) * 10000
			skillIDEnd := skillIDStart + 9999
			for skillID := skillIDStart; skillID <= skillIDEnd; skillID++ {
				wzSkill := gs.resources.GetSkill(skillID)
				if wzSkill == nil {
					continue
				}
				skillByID[skillID] = wzSkill
			}
		}

		orderedSkillIDs := make([]int, 0, len(skillByID))
		for skillID := range skillByID {
			orderedSkillIDs = append(orderedSkillIDs, int(skillID))
		}
		sort.Ints(orderedSkillIDs)

		index := 1
		for _, skillID := range orderedSkillIDs {
			skillTable := skillToLuaWzTable(L, skillByID[uint32(skillID)])
			result.RawSetInt(index, skillTable)
			index++
		}

		L.Push(result)
		return 1
	})

	luax.RegisterFunc(luaState, "run_script", func(L *lua.LState) int {
		path := L.CheckString(1)
		fn, err := L.LoadFile(path)
		if err != nil {
			L.RaiseError("run_script: %v", err)
			return 0
		}
		L.Push(fn)
		if err := L.PCall(0, 0, nil); err != nil {
			L.RaiseError("run_script %s: %v", path, err)
			return 0
		}
		return 0
	})

	luax.RegisterFunc(luaState, "sleep", func(L *lua.LState) int {
		duration := L.CheckNumber(1)
		pid := luax.GetThreadPID(L)
		if pid == nil {
			return 0
		}
		root := luax.GetRootLuaState(pid.String())
		if root == nil {
			return 0
		}
		d := time.Duration(float64(duration) * float64(time.Millisecond))
		time.AfterFunc(d, func() {
			gs.GetRootContext().Send(pid, &g_actor.ResumeLua{Root: root, Thread: L})
		})
		return L.Yield(lua.LNil)
	})

	if fn, err := luaState.LoadFile("script/skill/skill.lua"); err != nil {
		log.Printf("Failed to load script/skill/skill.lua: %v", err)
	} else {
		luaState.Push(fn)
		if err := luaState.PCall(0, 0, nil); err != nil {
			log.Printf("Failed to run script/skill/skill.lua: %v", err)
		}
	}
}
