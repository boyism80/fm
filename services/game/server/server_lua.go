package server

import (
	"log"
	"sort"
	"time"

	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/core/luax"
	pconst "github.com/boyism80/fm/protocol/constant"
	"github.com/boyism80/fm/protocol/response"
	g_actor "github.com/boyism80/fm/services/game/actor"
	"github.com/boyism80/fm/services/game/constant"
	"github.com/boyism80/fm/services/game/entity"
	"github.com/boyism80/fm/services/game/wz"
	lua "github.com/yuin/gopher-lua"
)

func getClassAdvancementClasses(class uint16) []uint16 {
	if class < 100 {
		return []uint16{class}
	}
	classes := make([]uint16, 0, 4)
	baseClass := (class / 100) * 100
	secondClass := (class / 10) * 10
	thirdClass := secondClass + 1
	classes = append(classes, baseClass)
	if secondClass != baseClass {
		classes = append(classes, secondClass)
	}
	if thirdClass != secondClass && thirdClass <= class {
		classes = append(classes, thirdClass)
	}
	if class != thirdClass && class != secondClass && class != baseClass {
		classes = append(classes, class)
	}
	return classes
}

func registerSkillConstants(luaState *lua.LState) {
	skillTable := luaState.NewTable()
	for key, skillID := range constant.AllSkillConstants() {
		skillTable.RawSetString(key, lua.LNumber(skillID))
	}
	luaState.SetGlobal("Skill", skillTable)
}

func registerBuffFlagAndMobBuff(luaState *lua.LState) {
	buffFlagTable := luaState.NewTable()
	for name, bf := range constant.AllBuffFlags() {
		entry := luaState.NewTable()
		entry.RawSetString("mask", lua.LNumber(bf.Mask))
		entry.RawSetString("position", lua.LNumber(bf.Position))
		buffFlagTable.RawSetString(name, entry)
	}
	luaState.SetGlobal("BuffFlag", buffFlagTable)
	debuffFlagTable := luaState.NewTable()
	for name, df := range constant.AllDebuffFlags() {
		entry := luaState.NewTable()
		entry.RawSetString("mask", lua.LNumber(df.Mask))
		entry.RawSetString("position", lua.LNumber(df.Position))
		entry.RawSetString("debuff", lua.LNumber(df.DiseaseSkillID))
		debuffFlagTable.RawSetString(name, entry)
	}
	luaState.SetGlobal("DebuffFlag", debuffFlagTable)
	mobBuffTable := luaState.NewTable()
	for name, st := range constant.AllMobBuffs() {
		entry := luaState.NewTable()
		entry.RawSetString("mask", lua.LNumber(st))
		mobBuffTable.RawSetString(name, entry)
	}
	luaState.SetGlobal("MobBuff", mobBuffTable)
	buffTypeTable := luaState.NewTable()
	for name, buffType := range constant.AllBuffTypes() {
		buffTypeTable.RawSetString(name, lua.LNumber(buffType))
	}
	luaState.SetGlobal("BuffType", buffTypeTable)
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
	t.RawSetString("Cap", lua.LNumber(constant.EquipmentPartsCap))
	t.RawSetString("Face", lua.LNumber(constant.EquipmentPartsFace))
	t.RawSetString("Eye", lua.LNumber(constant.EquipmentPartsEye))
	t.RawSetString("Ear", lua.LNumber(constant.EquipmentPartsEar))
	t.RawSetString("Top", lua.LNumber(constant.EquipmentPartsTop))
	t.RawSetString("Pants", lua.LNumber(constant.EquipmentPartsPants))
	t.RawSetString("Shoes", lua.LNumber(constant.EquipmentPartsShoes))
	t.RawSetString("Glove", lua.LNumber(constant.EquipmentPartsGlove))
	t.RawSetString("Cape", lua.LNumber(constant.EquipmentPartsCape))
	t.RawSetString("Shield", lua.LNumber(constant.EquipmentPartsShield))
	t.RawSetString("Weapon", lua.LNumber(constant.EquipmentPartsWeapon))
	t.RawSetString("Ring", lua.LNumber(constant.EquipmentPartsRing))
	luaState.SetGlobal("EquipmentPart", t)
}

func registerInventoryTypeConstants(luaState *lua.LState) {
	t := luaState.NewTable()
	t.RawSetString("Equipment", lua.LNumber(constant.InventoryTypeEquipment))
	t.RawSetString("Use", lua.LNumber(constant.InventoryTypeConsume))
	t.RawSetString("Consume", lua.LNumber(constant.InventoryTypeConsume))
	t.RawSetString("Installation", lua.LNumber(constant.InventoryTypeInstallation))
	t.RawSetString("Etc", lua.LNumber(constant.InventoryTypeETC))
	t.RawSetString("Cash", lua.LNumber(constant.InventoryTypeCash))
	luaState.SetGlobal("InventoryType", t)
}

func registerGenderConstants(luaState *lua.LState) {
	t := luaState.NewTable()
	t.RawSetString("Male", lua.LNumber(constant.GenderMale))
	t.RawSetString("Female", lua.LNumber(constant.GenderFemale))
	luaState.SetGlobal("Gender", t)
}

func registerMorphConstants(luaState *lua.LState) {
	t := luaState.NewTable()
	t.RawSetString("Barrel", lua.LNumber(constant.MorphBarrel))
	t.RawSetString("Transform", lua.LNumber(constant.MorphTransform))
	t.RawSetString("TransformFemale", lua.LNumber(constant.MorphTransformFemale))
	t.RawSetString("SuperTransform", lua.LNumber(constant.MorphSuperTransform))
	t.RawSetString("SuperTransformFemale", lua.LNumber(constant.MorphSuperTransformFemale))
	t.RawSetString("Albatross", lua.LNumber(constant.MorphAlbatross))
	t.RawSetString("AlbatrossFemale", lua.LNumber(constant.MorphAlbatrossFemale))
	luaState.SetGlobal("Morph", t)
}

func registerStatConstants(luaState *lua.LState) {
	t := luaState.NewTable()
	t.RawSetString("Level", lua.LNumber(constant.StatLevel))
	t.RawSetString("Exp", lua.LNumber(constant.StatEXP))
	t.RawSetString("Class", lua.LNumber(constant.StatClass))
	t.RawSetString("Str", lua.LNumber(constant.StatStr))
	t.RawSetString("Dex", lua.LNumber(constant.StatDex))
	t.RawSetString("Int", lua.LNumber(constant.StatInt))
	t.RawSetString("Luk", lua.LNumber(constant.StatLuk))
	t.RawSetString("Hp", lua.LNumber(constant.StatHP))
	t.RawSetString("MaxHp", lua.LNumber(constant.StatMaxHP))
	t.RawSetString("Mp", lua.LNumber(constant.StatMP))
	t.RawSetString("MaxMp", lua.LNumber(constant.StatMaxMP))
	t.RawSetString("AvailableAP", lua.LNumber(constant.StatAvailableAP))
	t.RawSetString("AvailableSP", lua.LNumber(constant.StatAvailableSP))
	t.RawSetString("Fame", lua.LNumber(constant.StatFame))
	t.RawSetString("Meso", lua.LNumber(constant.StatMeso))
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

func registerRoleConstants(luaState *lua.LState) {
	t := luaState.NewTable()
	for name, value := range constant.AllCharacterRoles() {
		t.RawSetString(name, lua.LNumber(value))
	}
	luaState.SetGlobal("ROLE", t)
}

func registerMobDieAnimationConstants(luaState *lua.LState) {
	t := luaState.NewTable()
	for name, value := range constant.AllMobDieAnimationTypes() {
		t.RawSetString(name, lua.LNumber(value))
	}
	luaState.SetGlobal("MobDieAnimation", t)
}

func registerSummonConstants(luaState *lua.LState) {
	moveTable := luaState.NewTable()
	for name, value := range constant.AllSummonMovementTypes() {
		moveTable.RawSetString(name, lua.LNumber(value))
	}
	luaState.SetGlobal("SummonMovementType", moveTable)

	typeTable := luaState.NewTable()
	for name, value := range constant.AllSummonTypes() {
		typeTable.RawSetString(name, lua.LNumber(value))
	}
	luaState.SetGlobal("SummonType", typeTable)
}

func registerIncomingHitConstants(luaState *lua.LState) {
	t := luaState.NewTable()
	for name, v := range constant.AllIncomingHitConstants() {
		t.RawSetString(name, lua.LNumber(v))
	}
	luaState.SetGlobal("IncomingHit", t)
}

func registerMistTypeConstants(luaState *lua.LState) {
	t := luaState.NewTable()
	for name, value := range constant.AllMistTypes() {
		t.RawSetString(name, lua.LNumber(value))
	}
	luaState.SetGlobal("MistType", t)
}

func registerSkillEffectTypeConstants(luaState *lua.LState) {
	t := luaState.NewTable()
	t.RawSetString("Cast", lua.LNumber(pconst.SkillEffectTypeCast))
	t.RawSetString("Affected", lua.LNumber(pconst.SkillEffectTypeAffected))
	luaState.SetGlobal("SkillEffectType", t)
}

func registerEffectTypeConstants(luaState *lua.LState) {
	t := luaState.NewTable()
	t.RawSetString("LevelUp", lua.LNumber(response.EffectTypeLevelUp))
	t.RawSetString("ClassChange", lua.LNumber(response.EffectTypeClassChange))
	t.RawSetString("QuestCompletion", lua.LNumber(response.EffectTypeQuestCompletion))
	t.RawSetString("RegisterCard", lua.LNumber(response.EffectTypeRegisterCard))
	t.RawSetString("ItemLevelUp", lua.LNumber(response.EffectTypeItemLevelUp))
	luaState.SetGlobal("EffectType", t)
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
	luax.RegisterLuaType[*entity.PartyMember](luaState)
	luax.RegisterLuaType[*entity.Party](luaState)
	luax.RegisterLuaType[*entity.GuildMember](luaState)
	luax.RegisterLuaType[*entity.Guild](luaState)
	luax.RegisterLuaType[*entity.Alliance](luaState)
	luax.RegisterLuaType[*entity.ObjectCore](luaState)
	luax.RegisterLuaDerivedType[*entity.FieldPlacement, *entity.ObjectCore](luaState)
	luax.RegisterLuaDerivedType[*entity.Meso, *entity.FieldPlacement](luaState)
	luax.RegisterLuaDerivedType[*entity.LifeCore, *entity.ObjectCore](luaState)
	luax.RegisterLuaDerivedType[*entity.Character, *entity.LifeCore](luaState)
	luax.RegisterLuaDerivedType[*entity.Mob, *entity.LifeCore](luaState)
	luax.RegisterLuaDerivedType[*entity.Summon, *entity.LifeCore](luaState)
	luax.RegisterLuaDerivedType[*entity.Mist, *entity.ObjectCore](luaState)
	luax.RegisterLuaDerivedType[*entity.Npc, *entity.ObjectCore](luaState)
	luax.RegisterLuaType[*entity.Map](luaState)
	luax.RegisterLuaType[*entity.SkillEntry](luaState)
	luax.RegisterLuaType[*entity.MobSkill](luaState)
	luax.RegisterLuaType[*entity.SkillBuff](luaState)
	luax.RegisterLuaType[*entity.ItemBuff](luaState)
	luax.RegisterLuaType[*entity.MobBuff](luaState)
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
	luax.RegisterLuaType[*wz.ItemCore](luaState)
	luax.RegisterLuaDerivedType[*wz.EquipmentCore, *wz.ItemCore](luaState)
	luax.RegisterLuaDerivedType[*wz.Weapon, *wz.EquipmentCore](luaState)
	luax.RegisterLuaDerivedType[*wz.Armor, *wz.EquipmentCore](luaState)
	luax.RegisterLuaDerivedType[*wz.Consume, *wz.ItemCore](luaState)
	luax.RegisterLuaDerivedType[*wz.Pet, *wz.ItemCore](luaState)
	luax.RegisterLuaDerivedType[*wz.MiscItem, *wz.ItemCore](luaState)
	luax.RegisterLuaDerivedType[*wz.CashItem, *wz.ItemCore](luaState)
	luax.RegisterLuaDerivedType[*wz.Installation, *wz.ItemCore](luaState)
	luax.RegisterLuaDerivedType[*entity.Consume, *entity.ItemCore](luaState)
	luax.RegisterLuaDerivedType[*entity.CashItem, *entity.ItemCore](luaState)
	luax.RegisterLuaDerivedType[*entity.MiscItem, *entity.ItemCore](luaState)
	luax.RegisterLuaDerivedType[*entity.Installation, *entity.ItemCore](luaState)
	luax.RegisterLuaDerivedType[*entity.Pet, *entity.ItemCore](luaState)

	registerBuffFlagAndMobBuff(luaState)
	registerWeaponTypeAndConsumeType(luaState)
	registerSkillConstants(luaState)
	registerEquipmentPartConstants(luaState)
	registerInventoryTypeConstants(luaState)
	registerStatConstants(luaState)
	registerGenderConstants(luaState)
	registerMorphConstants(luaState)
	registerClassConstants(luaState)
	registerStanceConstants(luaState)
	registerObjectTypeConstants(luaState)
	registerRoleConstants(luaState)
	registerMobDieAnimationConstants(luaState)
	registerSummonConstants(luaState)
	registerIncomingHitConstants(luaState)
	registerMistTypeConstants(luaState)
	registerSkillEffectTypeConstants(luaState)
	registerEffectTypeConstants(luaState)

	luax.RegisterFunc(luaState, "name2map", func(L *lua.LState) int {
		name := L.CheckString(1)
		if gs.resources == nil {
			L.Push(lua.LNil)
			return 1
		}
		id, ok := gs.resources.NameToMap(name)
		if !ok {
			L.Push(lua.LNil)
			return 1
		}
		spec := gs.resources.Maps[id]
		if spec == nil {
			L.Push(lua.LNil)
			return 1
		}
		tbl := L.NewTable()
		tbl.RawSetString("id", lua.LNumber(spec.ID))
		tbl.RawSetString("name", lua.LString(spec.Name))
		L.Push(tbl)
		return 1
	})

	luax.RegisterFunc(luaState, "id2map", func(L *lua.LState) int {
		id := uint32(L.CheckNumber(1))
		if gs.resources == nil {
			L.Push(lua.LNil)
			return 1
		}
		spec := gs.resources.Maps[id]
		if spec == nil {
			L.Push(lua.LNil)
			return 1
		}
		tbl := L.NewTable()
		tbl.RawSetString("id", lua.LNumber(spec.ID))
		tbl.RawSetString("name", lua.LString(spec.Name))
		L.Push(tbl)
		return 1
	})

	luax.RegisterFunc(luaState, "name2mob", func(L *lua.LState) int {
		name := L.CheckString(1)
		if gs.resources == nil {
			L.Push(lua.LNil)
			return 1
		}
		id, ok := gs.resources.NameToMob(name)
		if !ok {
			L.Push(lua.LNil)
			return 1
		}
		mob := gs.resources.Monsters[id]
		if mob == nil {
			L.Push(lua.LNil)
			return 1
		}
		tbl := L.NewTable()
		tbl.RawSetString("id", lua.LNumber(mob.ID))
		L.Push(tbl)
		return 1
	})

	luax.RegisterFunc(luaState, "id2mob", func(L *lua.LState) int {
		id := uint32(L.CheckNumber(1))
		if gs.resources == nil {
			L.Push(lua.LNil)
			return 1
		}
		mob := gs.resources.Monsters[id]
		if mob == nil {
			L.Push(lua.LNil)
			return 1
		}
		tbl := L.NewTable()
		tbl.RawSetString("id", lua.LNumber(mob.ID))
		name := gs.resources.GetMobName(mob.ID)
		if name != "" {
			tbl.RawSetString("name", lua.LString(name))
		}
		L.Push(tbl)
		return 1
	})

	luax.RegisterFunc(luaState, "name2npc", func(L *lua.LState) int {
		name := L.CheckString(1)
		if gs.resources == nil {
			L.Push(lua.LNil)
			return 1
		}
		id, ok := gs.resources.NameToNpc(name)
		if !ok {
			L.Push(lua.LNil)
			return 1
		}
		L.Push(lua.LNumber(id))
		return 1
	})

	luax.RegisterFunc(luaState, "name2item", func(L *lua.LState) int {
		name := L.CheckString(1)
		if gs.resources == nil {
			L.Push(lua.LNil)
			return 1
		}
		id, ok := gs.resources.NameToItem(name)
		if !ok {
			L.Push(lua.LNil)
			return 1
		}
		L.Push(lua.LNumber(id))
		return 1
	})

	luax.RegisterFunc(luaState, "name2skill", func(L *lua.LState) int {
		name := L.CheckString(1)
		if gs.resources == nil {
			L.Push(lua.LNil)
			return 1
		}
		id, ok := gs.resources.NameToSkill(name)
		if !ok {
			L.Push(lua.LNil)
			return 1
		}
		wzSkill := gs.resources.GetSkill(id)
		if wzSkill == nil {
			L.Push(lua.LNil)
			return 1
		}
		L.Push(skillToLuaWzTable(L, wzSkill))
		return 1
	})

	luax.RegisterFunc(luaState, "item_wz", func(L *lua.LState) int {
		if gs.resources == nil {
			L.Push(lua.LNil)
			return 1
		}

		argc := L.GetTop()
		if argc < 1 {
			L.ArgError(1, "item_wz(itemIdOrName [, count]) requires at least one argument")
			return 0
		}

		count := uint16(1)
		switch argc {
		case 1:
		case 2:
			if n := L.CheckInt(2); n >= 1 {
				count = uint16(n)
			}
		default:
			L.ArgError(2, "item_wz(itemIdOrName [, count]) expects 1 or 2 arguments")
			return 0
		}

		var itemId uint32
		switch lv := L.Get(1).(type) {
		case lua.LString:
			id, ok := gs.resources.NameToItem(string(lv))
			if !ok {
				L.Push(lua.LNil)
				return 1
			}
			itemId = id
		case lua.LNumber:
			itemId = uint32(lv)
		default:
			L.ArgError(1, "item id (number) or item name (string) expected")
			return 0
		}

		model, ok := gs.resources.Items[itemId]
		if !ok {
			L.Push(lua.LNil)
			return 1
		}

		_ = count
		if luable, ok := model.(luax.Luable); ok && luable != nil {
			L.Push(luax.NewLuable(L, luable))
			return 1
		}
		L.Push(lua.LNil)
		return 1
	})

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

	luax.RegisterFunc(luaState, "save", func(L *lua.LState) int {
		cfg, ok := luax.GetConfiguration(L)
		if !ok || cfg.ActorContext == nil {
			L.Push(lua.LBool(false))
			L.Push(lua.LString("save: actor context not found"))
			return 2
		}
		actorCtx := cfg.ActorContext
		if actorCtx == nil {
			L.Push(lua.LBool(false))
			L.Push(lua.LString("save: actor context not found"))
			return 2
		}
		pid := actorCtx.Self()
		if pid == nil {
			L.Push(lua.LBool(false))
			L.Push(lua.LString("save: actor PID not found"))
			return 2
		}
		mapInstance := gs.getMapByActorPID(pid)
		if mapInstance == nil {
			L.Push(lua.LBool(false))
			L.Push(lua.LString("save: map not found"))
			return 2
		}
		root := mapInstance.GetLuaRoot()
		if root == nil {
			L.Push(lua.LBool(false))
			L.Push(lua.LString("save: root lua state not found"))
			return 2
		}
		cfg, _ = luax.GetConfiguration(L)
		cfg.ActorContext = actorCtx
		cfg.KeepAlive = true
		luax.SetConfiguration(L, cfg)

		p := gs.SaveAllCharactersAsync(actorCtx)
		if p == nil {
			L.Push(lua.LBool(false))
			L.Push(lua.LString("save: promise creation failed"))
			return 2
		}
		var succeeded = true
		var errMsg string
		p.OnError(func(err error) {
			succeeded = false
			if err != nil {
				errMsg = err.Error()
				log.Printf("save: %v", err)
			}
		}).Finally(func() {
			args := []lua.LValue{lua.LBool(succeeded)}
			if errMsg != "" {
				args = append(args, lua.LString(errMsg))
			} else {
				args = append(args, lua.LNil)
			}
			gs.GetRootContext().Send(pid, &g_actor.ResumeLua{Root: root, Thread: L, Args: args})
		}).Run()
		return L.Yield(lua.LNil, lua.LNil)
	})

	luax.RegisterFunc(luaState, "sleep", func(L *lua.LState) int {
		duration := L.CheckNumber(1)
		cfg, ok := luax.GetConfiguration(L)
		if !ok || cfg.ActorContext == nil {
			return 0
		}
		pid := cfg.ActorContext.Self()
		if pid == nil {
			return 0
		}
		mapInstance := gs.getMapByActorPID(pid)
		if mapInstance == nil {
			return 0
		}
		root := mapInstance.GetLuaRoot()
		if root == nil {
			return 0
		}
		d := time.Duration(float64(duration) * float64(time.Millisecond))
		time.AfterFunc(d, func() {
			gs.GetRootContext().Send(pid, &g_actor.ResumeLua{Root: root, Thread: L})
		})
		return L.Yield(lua.LNil)
	})

	luax.RegisterFunc(luaState, "set_packet_log_enabled", func(L *lua.LState) int {
		core.SetPacketLogEnabled(L.CheckBool(1))
		return 0
	})
	luax.RegisterFunc(luaState, "get_packet_log_enabled", func(L *lua.LState) int {
		L.Push(lua.LBool(core.GetPacketLogEnabled()))
		return 1
	})
	luax.RegisterFunc(luaState, "set_packet_log", func(L *lua.LState) int {
		core.SetPacketLogEnabled(L.CheckBool(1))
		return 0
	})
	luax.RegisterFunc(luaState, "get_packet_log", func(L *lua.LState) int {
		L.Push(lua.LBool(core.GetPacketLogEnabled()))
		return 1
	})

	luax.RegisterFunc(luaState, "run_on_script", func(L *lua.LState) int {
		ud := L.CheckUserData(1)
		ch, ok := ud.Value.(*entity.Character)
		if !ok || ch == nil {
			L.RaiseError("run_on_script: Character expected")
			return 0
		}
		cfg, ok := luax.GetConfiguration(L)
		if !ok || cfg.ActorContext == nil {
			L.RaiseError("run_on_script: thread has no actor PID (call from command context)")
			return 0
		}
		actorCtx := cfg.ActorContext
		pid := actorCtx.Self()
		if pid == nil {
			L.RaiseError("run_on_script: thread has no actor PID")
			return 0
		}
		mapInstance := ch.GetMap()
		if mapInstance == nil {
			L.RaiseError("run_on_script: character map not found")
			return 0
		}
		root := mapInstance.GetLuaRoot()
		if root == nil {
			L.RaiseError("run_on_script: root lua state not found")
			return 0
		}
		thread, err := luax.NewThread(root, "script/script.lua")
		if err == nil {
			_, err = luax.Call(thread, "on_script", ch)
		}
		if err != nil {
			L.RaiseError("run_on_script: %v", err)
			return 0
		}
		L.Push(lua.LBool(true))
		return 1
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
