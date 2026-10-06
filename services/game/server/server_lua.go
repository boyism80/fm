package server

import (
	"fmt"
	"log"
	"sort"
	"strings"
	"time"

	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/core/async"
	"github.com/boyism80/fm/core/luax"
	pconst "github.com/boyism80/fm/protocol/constant"
	"github.com/boyism80/fm/protocol/response"
	"github.com/boyism80/fm/services/game/constant"
	"github.com/boyism80/fm/services/game/entity"
	"github.com/boyism80/fm/services/game/wz"
	"github.com/boyism80/fm/types"
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

func luaArgToInterface(v lua.LValue) (interface{}, error) {
	if v == nil || v == lua.LNil {
		return nil, nil
	}
	switch v.Type() {
	case lua.LTBool:
		return lua.LVAsBool(v), nil
	case lua.LTNumber:
		return float64(lua.LVAsNumber(v)), nil
	case lua.LTString:
		return string(v.(lua.LString)), nil
	default:
		return nil, fmt.Errorf("unsupported argument type %s", v.Type().String())
	}
}

func (gs *GameServer) registerSkillConstants(luaState *lua.LState) {
	skillTable := luaState.NewTable()
	for key, skillID := range constant.AllSkillConstants() {
		skillTable.RawSetString(key, lua.LNumber(skillID))
	}
	luaState.SetGlobal("Skill", skillTable)
}

func (gs *GameServer) registerBuffFlag(luaState *lua.LState) {
	buffFlagTable := luaState.NewTable()
	for name, bf := range constant.AllBuffFlags() {
		entry := luaState.NewTable()
		entry.RawSetString("mask", lua.LNumber(bf.Mask))
		entry.RawSetString("position", lua.LNumber(bf.Position))
		buffFlagTable.RawSetString(name, entry)
	}
	luaState.SetGlobal("BuffFlag", buffFlagTable)
}

func (gs *GameServer) registerDebuffFlag(luaState *lua.LState) {
	debuffFlagTable := luaState.NewTable()
	for name, df := range constant.AllDebuffFlags() {
		entry := luaState.NewTable()
		entry.RawSetString("mask", lua.LNumber(df.Mask))
		entry.RawSetString("position", lua.LNumber(df.Position))
		entry.RawSetString("debuff", lua.LNumber(df.DiseaseSkillID))
		debuffFlagTable.RawSetString(name, entry)
	}
	luaState.SetGlobal("DebuffFlag", debuffFlagTable)
}

func (gs *GameServer) registerMobBuff(luaState *lua.LState) {
	mobBuffTable := luaState.NewTable()
	for name, st := range constant.AllMobBuffs() {
		entry := luaState.NewTable()
		entry.RawSetString("mask", lua.LNumber(st))
		mobBuffTable.RawSetString(name, entry)
	}
	luaState.SetGlobal("MobBuff", mobBuffTable)
}

func (gs *GameServer) registerBuffType(luaState *lua.LState) {
	buffTypeTable := luaState.NewTable()
	for name, buffType := range constant.AllBuffTypes() {
		buffTypeTable.RawSetString(name, lua.LNumber(buffType))
	}
	luaState.SetGlobal("BuffType", buffTypeTable)
}

func (gs *GameServer) registerWeaponType(luaState *lua.LState) {
	weaponTypeTable := luaState.NewTable()
	for name, wt := range constant.AllWeaponTypes() {
		weaponTypeTable.RawSetString(name, lua.LNumber(wt))
	}
	luaState.SetGlobal("WeaponType", weaponTypeTable)
}

func (gs *GameServer) registerConsumeType(luaState *lua.LState) {
	consumeTypeTable := luaState.NewTable()
	for name, ct := range constant.AllConsumeTypes() {
		consumeTypeTable.RawSetString(name, lua.LNumber(ct))
	}
	luaState.SetGlobal("ConsumeType", consumeTypeTable)
}

func (gs *GameServer) registerEquipmentPartConstants(luaState *lua.LState) {
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
	t.RawSetString("Ring2", lua.LNumber(constant.EquipmentPartsRing2))
	t.RawSetString("Ring3", lua.LNumber(constant.EquipmentPartsRing3))
	t.RawSetString("Ring4", lua.LNumber(constant.EquipmentPartsRing4))
	t.RawSetString("Pendant", lua.LNumber(constant.EquipmentPartsPendant))
	t.RawSetString("Medal", lua.LNumber(constant.EquipmentPartsMedal))
	luaState.SetGlobal("EquipmentPart", t)
}

func (gs *GameServer) registerInventoryTypeConstants(luaState *lua.LState) {
	t := luaState.NewTable()
	t.RawSetString("Equipment", lua.LNumber(constant.InventoryTypeEquipment))
	t.RawSetString("Use", lua.LNumber(constant.InventoryTypeConsume))
	t.RawSetString("Consume", lua.LNumber(constant.InventoryTypeConsume))
	t.RawSetString("Installation", lua.LNumber(constant.InventoryTypeInstallation))
	t.RawSetString("Etc", lua.LNumber(constant.InventoryTypeETC))
	t.RawSetString("Cash", lua.LNumber(constant.InventoryTypeCash))
	luaState.SetGlobal("InventoryType", t)
}

func (gs *GameServer) registerGenderConstants(luaState *lua.LState) {
	t := luaState.NewTable()
	t.RawSetString("Male", lua.LNumber(constant.GenderMale))
	t.RawSetString("Female", lua.LNumber(constant.GenderFemale))
	luaState.SetGlobal("Gender", t)
}

func (gs *GameServer) registerMorphConstants(luaState *lua.LState) {
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

func (gs *GameServer) registerStatConstants(luaState *lua.LState) {
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
	t.RawSetString("Population", lua.LNumber(constant.StatPopulation))
	t.RawSetString("Meso", lua.LNumber(constant.StatMeso))
	luaState.SetGlobal("STAT", t)
}

func (gs *GameServer) registerClassConstants(luaState *lua.LState) {
	t := luaState.NewTable()
	for name, code := range constant.AllClassConstants() {
		t.RawSetString(name, lua.LNumber(code))
	}
	luaState.SetGlobal("Class", t)
}

func (gs *GameServer) registerStanceConstants(luaState *lua.LState) {
	t := luaState.NewTable()
	for name, value := range constant.AllStanceConstants() {
		t.RawSetString(name, lua.LNumber(value))
	}
	luaState.SetGlobal("Stance", t)
}

func (gs *GameServer) registerObjectTypeConstants(luaState *lua.LState) {
	t := luaState.NewTable()
	for name, value := range constant.AllObjectTypeConstants() {
		t.RawSetString(name, lua.LNumber(value))
	}
	luaState.SetGlobal("ObjectType", t)
}

func (gs *GameServer) registerRoleConstants(luaState *lua.LState) {
	t := luaState.NewTable()
	for name, value := range constant.AllCharacterRoles() {
		t.RawSetString(name, lua.LNumber(value))
	}
	luaState.SetGlobal("ROLE", t)
}

func (gs *GameServer) registerMobDieAnimationConstants(luaState *lua.LState) {
	t := luaState.NewTable()
	for name, value := range constant.AllMobDieAnimationTypes() {
		t.RawSetString(name, lua.LNumber(value))
	}
	luaState.SetGlobal("MobDieAnimation", t)
}

func (gs *GameServer) registerMobSpawnTypeConstants(luaState *lua.LState) {
	t := luaState.NewTable()
	for name, value := range constant.AllMobSpawnTypes() {
		t.RawSetString(name, lua.LNumber(value))
	}
	luaState.SetGlobal("MobSpawnType", t)
}

func (gs *GameServer) registerSummonConstants(luaState *lua.LState) {
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

func (gs *GameServer) registerIncomingHitConstants(luaState *lua.LState) {
	t := luaState.NewTable()
	for name, v := range constant.AllIncomingHitConstants() {
		t.RawSetString(name, lua.LNumber(v))
	}
	luaState.SetGlobal("IncomingHit", t)
}

func (gs *GameServer) registerMistTypeConstants(luaState *lua.LState) {
	t := luaState.NewTable()
	for name, value := range constant.AllMistTypes() {
		t.RawSetString(name, lua.LNumber(value))
	}
	luaState.SetGlobal("MistType", t)
}

func (gs *GameServer) registerSkillEffectTypeConstants(luaState *lua.LState) {
	t := luaState.NewTable()
	t.RawSetString("Cast", lua.LNumber(pconst.SkillEffectTypeCast))
	t.RawSetString("Affected", lua.LNumber(pconst.SkillEffectTypeAffected))
	luaState.SetGlobal("SkillEffectType", t)
}

func (gs *GameServer) registerEffectTypeConstants(luaState *lua.LState) {
	t := luaState.NewTable()
	t.RawSetString("LevelUp", lua.LNumber(response.EffectTypeLevelUp))
	t.RawSetString("ClassChange", lua.LNumber(response.EffectTypeClassChange))
	t.RawSetString("QuestCompletion", lua.LNumber(response.EffectTypeQuestCompletion))
	t.RawSetString("RegisterCard", lua.LNumber(response.EffectTypeRegisterCard))
	t.RawSetString("ItemLevelUp", lua.LNumber(response.EffectTypeItemLevelUp))
	luaState.SetGlobal("EffectType", t)
}

func (gs *GameServer) registerExchangeResultConstants(luaState *lua.LState) {
	t := luaState.NewTable()
	t.RawSetString("OK", lua.LNumber(entity.ExchangeOK))
	t.RawSetString("LackCost", lua.LNumber(entity.ExchangeLackCost))
	t.RawSetString("LackCapacity", lua.LNumber(entity.ExchangeLackCapacity))
	luaState.SetGlobal("ExchangeResult", t)
}

func (gs *GameServer) registerServerMessageConstants(luaState *lua.LState) {
	t := luaState.NewTable()
	for name, value := range constant.AllServerMessageTypes() {
		t.RawSetString(name, lua.LNumber(value))
	}
	luaState.SetGlobal("Msg", t)

	scopes := luaState.NewTable()
	for name, value := range constant.AllMessageScopes() {
		scopes.RawSetString(name, lua.LNumber(value))
	}
	luaState.SetGlobal("MessageScope", scopes)
}

func (gs *GameServer) runStartupScript() {
	luaState := luax.NewState()
	defer luaState.Close()

	luax.RegisterFunc(luaState, "register_quest_item", func(L *lua.LState) int {
		gs.resources.QuestItems[uint32(L.CheckInt(1))] = uint32(L.CheckInt(2))
		return 0
	})
	if err := luaState.DoFile(constant.StartupScriptPath); err != nil {
		log.Printf("startup script: %v", err)
	}
}

func skillToLuaWz(luaState *lua.LState, skill *wz.Skill) lua.LValue {
	if skill == nil {
		return lua.LNil
	}
	return luax.NewLuable(luaState, skill)
}

func registerWzMobName(gs *GameServer, L *lua.LState) {
	var zero *wz.Mob
	meta := L.GetTypeMetatable(zero.LuaTypeName()).(*lua.LTable)
	L.SetFuncs(meta, map[string]lua.LGFunction{
		"name": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			m, ok := ud.Value.(*wz.Mob)
			if !ok || m == nil {
				L.ArgError(1, "WzMob expected")
				return 0
			}
			name := ""
			if gs.resources != nil {
				name = gs.resources.GetMobName(m.ID)
			}
			L.Push(lua.LString(name))
			return 1
		},
	})
}

func (gs *GameServer) registerGameLuaState(luaState *lua.LState) {
	luax.RegisterLuaType[*entity.PartyMember](luaState)
	luax.RegisterLuaType[*entity.Party](luaState)
	luax.RegisterLuaType[*entity.StateMachine](luaState)
	luax.RegisterLuaType[*entity.StateMachineGroup](luaState)
	luax.RegisterLuaType[*entity.Portal](luaState)
	luax.RegisterLuaType[*entity.GuildMember](luaState)
	luax.RegisterLuaType[*entity.Guild](luaState)
	luax.RegisterLuaType[*entity.Quest](luaState)
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
	luax.RegisterLuaType[*wz.Map](luaState)
	registerWzMapCreateInstance(gs, luaState)
	luax.RegisterLuaType[*wz.Mob](luaState)
	registerWzMobName(gs, luaState)
	luax.RegisterLuaType[*wz.Skill](luaState)
	luax.RegisterLuaType[*wz.Quest](luaState)
	luax.RegisterLuaType[*wz.NpcSpawn](luaState)
	luax.RegisterLuaType[*entity.SkillEntry](luaState)
	luax.RegisterLuaType[*entity.MobSkill](luaState)
	luax.RegisterLuaDerivedType[*entity.Reactor, *entity.ObjectCore](luaState)
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

	gs.registerBuffFlag(luaState)
	gs.registerDebuffFlag(luaState)
	gs.registerMobBuff(luaState)
	gs.registerBuffType(luaState)
	gs.registerWeaponType(luaState)
	gs.registerConsumeType(luaState)
	gs.registerSkillConstants(luaState)
	gs.registerEquipmentPartConstants(luaState)
	gs.registerInventoryTypeConstants(luaState)
	gs.registerStatConstants(luaState)
	gs.registerGenderConstants(luaState)
	gs.registerMorphConstants(luaState)
	gs.registerClassConstants(luaState)
	gs.registerStanceConstants(luaState)
	gs.registerObjectTypeConstants(luaState)
	gs.registerRoleConstants(luaState)
	gs.registerMobDieAnimationConstants(luaState)
	gs.registerMobSpawnTypeConstants(luaState)
	gs.registerSummonConstants(luaState)
	gs.registerIncomingHitConstants(luaState)
	gs.registerMistTypeConstants(luaState)
	gs.registerSkillEffectTypeConstants(luaState)
	gs.registerEffectTypeConstants(luaState)
	gs.registerExchangeResultConstants(luaState)
	gs.registerServerMessageConstants(luaState)
	registerClockLuaFuncs(gs, luaState)
	registerFaultLuaFuncs(gs, luaState)
	entity.RegisterCarnivalLua(luaState, gs)

	luax.RegisterFunc(luaState, "log", func(L *lua.LState) int {
		parts := make([]string, L.GetTop())
		for i := 1; i <= L.GetTop(); i++ {
			parts[i-1] = L.Get(i).String()
		}
		log.Printf("[lua] %s", strings.Join(parts, " "))
		return 0
	})

	luax.RegisterFunc(luaState, "state_machine", func(L *lua.LState) int {
		name := L.CheckString(1)
		reg := gs.GetStateMachineRegistry()
		if reg == nil {
			L.Push(lua.LNil)
			return 1
		}
		group := reg.Get(name)
		if group == nil {
			L.Push(lua.LNil)
			return 1
		}
		L.Push(luax.NewLuable(L, group))
		return 1
	})

	luax.RegisterFunc(luaState, "state_machines", func(L *lua.LState) int {
		tbl := L.NewTable()
		if gs.stateMachines == nil {
			L.Push(tbl)
			return 1
		}
		for _, group := range gs.stateMachines.Groups() {
			tbl.Append(luax.NewLuable(L, group))
		}
		L.Push(tbl)
		return 1
	})

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
		m := gs.resources.Maps[id]
		if m == nil {
			L.Push(lua.LNil)
			return 1
		}
		L.Push(luax.NewLuable(L, m))
		return 1
	})

	luax.RegisterFunc(luaState, "id2map", func(L *lua.LState) int {
		id := uint32(L.CheckNumber(1))
		if gs.resources == nil {
			L.Push(lua.LNil)
			return 1
		}
		m := gs.resources.Maps[id]
		if m == nil {
			L.Push(lua.LNil)
			return 1
		}
		L.Push(luax.NewLuable(L, m))
		return 1
	})

	luax.RegisterFunc(luaState, "instance_map", func(L *lua.LState) int {
		key := uint32(L.CheckNumber(1))
		ms := gs.GetMapSystem()
		if ms == nil {
			L.Push(lua.LNil)
			return 1
		}
		m := ms.GetInstance(key)
		if m == nil {
			L.Push(lua.LNil)
			return 1
		}
		L.Push(luax.NewLuable(L, m))
		return 1
	})

	luax.RegisterFunc(luaState, "npc_spawns", func(L *lua.LState) int {
		npcID := uint32(L.CheckNumber(1))
		out := L.NewTable()
		if gs.resources == nil {
			L.Push(out)
			return 1
		}
		spawns := gs.resources.FindNpcSpawns(npcID)
		for i, spawn := range spawns {
			entry := L.NewTable()
			entry.RawSetString("map_id", lua.LNumber(spawn.MapID))
			entry.RawSetString("x", lua.LNumber(spawn.X))
			entry.RawSetString("y", lua.LNumber(spawn.Y))
			out.RawSetInt(i+1, entry)
		}
		L.Push(out)
		return 1
	})

	luax.RegisterFunc(luaState, "closest_spawn", func(L *lua.LState) int {
		mapID := uint32(L.CheckNumber(1))
		x := int16(L.CheckNumber(2))
		y := int16(L.CheckNumber(3))
		if gs.resources == nil {
			L.Push(lua.LNil)
			return 1
		}
		mapWz := gs.resources.Maps[mapID]
		if mapWz == nil {
			L.Push(lua.LNil)
			return 1
		}
		spawnID := mapWz.FindClosestPortalSpawnID(types.Point[int16]{X: x, Y: y})
		L.Push(lua.LNumber(spawnID))
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
		L.Push(luax.NewLuable(L, mob))
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
		L.Push(luax.NewLuable(L, mob))
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

	luax.RegisterFunc(luaState, "item2name", func(L *lua.LState) int {
		id := uint32(L.CheckNumber(1))
		if gs.resources == nil {
			L.Push(lua.LNil)
			return 1
		}
		name := gs.resources.GetItemName(id)
		if name == "" {
			L.Push(lua.LNil)
			return 1
		}
		L.Push(lua.LString(name))
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
		L.Push(skillToLuaWz(L, wzSkill))
		return 1
	})

	luax.RegisterFunc(luaState, "id2quest", func(L *lua.LState) int {
		id := uint32(L.CheckNumber(1))
		if gs.resources == nil {
			L.Push(lua.LNil)
			return 1
		}
		wzQuest := gs.resources.GetQuest(id)
		if wzQuest == nil {
			L.Push(lua.LNil)
			return 1
		}
		L.Push(luax.NewLuable(L, wzQuest))
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
			result.RawSetInt(index, skillToLuaWz(L, skillByID[uint32(skillID)]))
			index++
		}

		L.Push(result)
		return 1
	})

	luax.RegisterFunc(luaState, "save", func(L *lua.LState) int {
		cfg, ok := luax.GetConfiguration(L)
		if !ok || cfg.ActorContext == nil {
			L.Push(lua.LBool(false))
			L.Push(lua.LString("save: actor context not found"))
			return 2
		}
		return entity.LuaYieldPromise(L, gs, gs.SaveAllCharactersAsync(cfg.ActorContext), func(_ interface{}, err error) []lua.LValue {
			if err != nil {
				log.Printf("save: %v", err)
				return []lua.LValue{lua.LFalse, lua.LString(err.Error())}
			}
			return []lua.LValue{lua.LTrue, lua.LNil}
		})
	})

	luax.RegisterFunc(luaState, "sleep", func(L *lua.LState) int {
		d := time.Duration(float64(L.CheckNumber(1)) * float64(time.Millisecond))
		promise := async.NewDeferred(nil)
		time.AfterFunc(d, func() {
			promise.SetResult(nil)
		})
		return entity.LuaYieldPromise(L, gs, promise, nil)
	})

	luax.RegisterFunc(luaState, "run_on_map", func(L *lua.LState) int {
		mapID := uint32(L.CheckInt(1))
		scriptPath := L.CheckString(2)
		funcName := L.CheckString(3)
		args := make([]interface{}, 0, L.GetTop()-3)
		for i := 4; i <= L.GetTop(); i++ {
			arg, err := luaArgToInterface(L.Get(i))
			if err != nil {
				L.RaiseError("run_on_map: %v", err)
				return 0
			}
			args = append(args, arg)
		}
		cfg, _ := luax.GetConfiguration(L)
		caller := cfg.ActorPID
		if cfg.ActorContext != nil {
			caller = cfg.ActorContext.Self()
		}
		targetMap := gs.GetMapSystem().Find(gs.stateMachines.FindByActor(caller), mapID)
		return gs.GetMapSystem().RunScript(L, cfg.ActorContext, targetMap, scriptPath, funcName, args)
	})

	luax.RegisterFunc(luaState, "set_packet_log_enabled", func(L *lua.LState) int {
		core.SetPacketLogEnabled(L.CheckBool(1))
		return 0
	})
	luax.RegisterFunc(luaState, "get_packet_log_enabled", func(L *lua.LState) int {
		L.Push(lua.LBool(core.GetPacketLogEnabled()))
		return 1
	})
	luax.RegisterFunc(luaState, "channel_id", func(L *lua.LState) int {
		L.Push(lua.LNumber(gs.config.ChannelId))
		return 1
	})
	luax.RegisterFunc(luaState, "set_megaphone_muted", func(L *lua.LState) int {
		gs.megaphoneMuted.Store(L.CheckBool(1))
		return 0
	})
	luax.RegisterFunc(luaState, "get_megaphone_muted", func(L *lua.LState) int {
		L.Push(lua.LBool(gs.megaphoneMuted.Load()))
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
}
