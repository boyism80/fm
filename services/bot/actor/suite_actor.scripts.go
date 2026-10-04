package actor

import (
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strconv"
	"strings"

	"github.com/boyism80/fm/services/game/wz"
	lua "github.com/yuin/gopher-lua"
)

var (
	sevenDigits  = regexp.MustCompile(`\b(\d{7})\b`)
	questCall    = regexp.MustCompile(`quest\(\s*(\d+)\s*\)`)
	itemChange   = regexp.MustCompile(`(?:\[\s*(\d{7})\s*\]\s*=|(?:mkitem|rmitem)\(\s*(\d{7})\b)`)
	mesoChange   = regexp.MustCompile(`\bmeso\s*=`)
	expChange    = regexp.MustCompile(`\bexp\s*=`)
	mapLiteral   = regexp.MustCompile(`:map\(\s*(\d{9})\b`)
	mapVariable  = regexp.MustCompile(`:map\(\s*[^\d\s)]`)
	mapAnyNumber = regexp.MustCompile(`(?:^|[^#m\d])(\d{9})\b`)
	requireCall  = regexp.MustCompile(`require\(\s*"script/([^"]+)"\s*\)`)
	deferredCall = regexp.MustCompile(`start_solo|start_party|state_machine\(|run_on_map`)
)

type npcSpot struct {
	mapID uint32
	spawn *wz.BaseSpawn
}

func (a *SuiteActor) wzFuncs() map[string]lua.LGFunction {
	return map[string]lua.LGFunction{
		"npc_scripts": func(L *lua.LState) int {
			ids, err := a.scriptIDs("npc")
			if err != nil {
				L.RaiseError("wz.npc_scripts: %v", err)
				return 0
			}
			placed := a.npcSpots()

			found, missing := L.NewTable(), L.NewTable()
			for _, id := range ids {
				at, ok := placed[id]
				if ok == false {
					missing.Append(lua.LNumber(id))
					continue
				}
				found.Append(a.npcEntry(L, id, at))
			}
			L.Push(found)
			L.Push(missing)
			return 2
		},
		"quest_scripts": func(L *lua.LState) int {
			ids, err := a.scriptIDs("quest")
			if err != nil {
				L.RaiseError("wz.quest_scripts: %v", err)
				return 0
			}
			placed := a.npcSpots()

			found, missing := L.NewTable(), L.NewTable()
			for _, id := range ids {
				quest := a.wz.GetQuest(id)
				if quest == nil || (quest.HasStartScript() == false && quest.HasEndScript() == false) {
					missing.Append(lua.LNumber(id))
					continue
				}
				entry := L.NewTable()
				entry.RawSetString("quest", lua.LNumber(id))
				entry.RawSetString("start_script", lua.LBool(quest.HasStartScript()))
				entry.RawSetString("end_script", lua.LBool(quest.HasEndScript()))
				if field := quest.Start.Requirements.FieldEnter; field > 0 {
					start := L.NewTable()
					start.RawSetString("npc", lua.LNumber(quest.Start.Requirements.NPC))
					start.RawSetString("map", lua.LNumber(field))
					entry.RawSetString("start", start)
				} else if at, ok := placed[quest.Start.Requirements.NPC]; ok {
					entry.RawSetString("start", a.npcEntry(L, quest.Start.Requirements.NPC, at))
				}
				if at, ok := placed[quest.Complete.Requirements.NPC]; ok {
					entry.RawSetString("finish", a.npcEntry(L, quest.Complete.Requirements.NPC, at))
				}
				found.Append(entry)
			}
			L.Push(found)
			L.Push(missing)
			return 2
		},
		"script_refs": func(L *lua.LState) int {
			data, err := os.ReadFile(filepath.Join(a.cfg.GameScriptDir, L.CheckString(1), L.CheckString(2)+".lua"))
			if err != nil {
				L.Push(lua.LNil)
				return 1
			}
			maps := mapLiteral.FindAllSubmatch(data, -1)
			if mapVariable.Match(data) {
				maps = append(maps, mapAnyNumber.FindAllSubmatch(data, -1)...)
			}
			var changes [][][]byte
			for _, m := range itemChange.FindAllSubmatch(data, -1) {
				if len(m[1]) == 0 {
					m[1] = m[2]
				}
				changes = append(changes, m)
			}
			var quests []uint32
			for _, m := range questCall.FindAllSubmatch(data, -1) {
				id, err := strconv.ParseUint(string(m[1]), 10, 32)
				if err != nil || a.wz.GetQuest(uint32(id)) == nil || slices.Contains(quests, uint32(id)) {
					continue
				}
				quests = append(quests, uint32(id))
			}
			slices.Sort(quests)

			refs := L.NewTable()
			refs.RawSetString("items", a.numberList(L, a.itemIDs(sevenDigits.FindAllSubmatch(data, -1))))
			refs.RawSetString("quests", a.numberList(L, quests))
			refs.RawSetString("maps", a.numberList(L, a.mapIDs(maps)))
			refs.RawSetString("changes", a.numberList(L, a.itemIDs(changes)))
			refs.RawSetString("meso", lua.LBool(mesoChange.Match(data)))
			refs.RawSetString("exp", lua.LBool(expChange.Match(data)))
			deferred := deferredCall.Match(data)
			for _, m := range requireCall.FindAllSubmatch(data, -1) {
				lib, err := os.ReadFile(filepath.Join(a.cfg.GameScriptDir, string(m[1])+".lua"))
				if err == nil && deferredCall.Match(lib) {
					deferred = true
				}
			}
			refs.RawSetString("deferred", lua.LBool(deferred))
			L.Push(refs)
			return 1
		},
		"portal_scripts": func(L *lua.LState) int {
			names, err := a.scriptNames("portal")
			if err != nil {
				L.RaiseError("wz.portal_scripts: %v", err)
				return 0
			}
			var mapIDs []uint32
			for mapID := range a.wz.Maps {
				mapIDs = append(mapIDs, mapID)
			}
			slices.Sort(mapIDs)

			placed := make(map[string]*lua.LTable)
			for _, mapID := range mapIDs {
				m := a.wz.Maps[mapID]
				var portalIDs []uint8
				for id := range m.Portals {
					portalIDs = append(portalIDs, id)
				}
				slices.Sort(portalIDs)
				for _, id := range portalIDs {
					portal := m.Portals[id]
					if portal.ScriptName == "" {
						continue
					}
					key := strings.ToLower(portal.ScriptName)
					if _, ok := placed[key]; ok {
						continue
					}
					entry := L.NewTable()
					entry.RawSetString("map", lua.LNumber(mapID))
					entry.RawSetString("portal", lua.LString(portal.Name))
					entry.RawSetString("script", lua.LString(portal.ScriptName))
					placed[key] = entry
				}
			}

			found, missing := L.NewTable(), L.NewTable()
			for _, name := range names {
				entry, ok := placed[strings.ToLower(name)]
				if ok == false {
					missing.Append(lua.LString(name))
					continue
				}
				found.Append(entry)
			}
			L.Push(found)
			L.Push(missing)
			return 2
		},
	}
}

func (a *SuiteActor) npcSpots() map[uint32]npcSpot {
	placed := make(map[uint32]npcSpot)
	for mapID, m := range a.wz.Maps {
		for _, spawn := range m.NpcSpawns {
			if spawn.BaseSpawn == nil || spawn.Hide {
				continue
			}
			if cur, ok := placed[spawn.ID]; ok && cur.mapID <= mapID {
				continue
			}
			placed[spawn.ID] = npcSpot{mapID: mapID, spawn: spawn.BaseSpawn}
		}
	}
	return placed
}

func (a *SuiteActor) npcEntry(L *lua.LState, npcID uint32, at npcSpot) *lua.LTable {
	entry := L.NewTable()
	entry.RawSetString("npc", lua.LNumber(npcID))
	entry.RawSetString("map", lua.LNumber(at.mapID))
	entry.RawSetString("x", lua.LNumber(at.spawn.Position.X))
	entry.RawSetString("y", lua.LNumber(at.spawn.Position.Y))
	entry.RawSetString("foothold", lua.LNumber(at.spawn.Foothold))
	return entry
}

func (a *SuiteActor) numberList(L *lua.LState, ids []uint32) *lua.LTable {
	t := L.NewTable()
	for _, id := range ids {
		t.Append(lua.LNumber(id))
	}
	return t
}

func (a *SuiteActor) itemIDs(matches [][][]byte) []uint32 {
	var ids []uint32
	for _, m := range matches {
		id, err := strconv.ParseUint(string(m[1]), 10, 32)
		if err != nil || a.wz.GetItemName(uint32(id)) == "" || slices.Contains(ids, uint32(id)) {
			continue
		}
		ids = append(ids, uint32(id))
	}
	slices.Sort(ids)
	return ids
}

func (a *SuiteActor) mapIDs(matches [][][]byte) []uint32 {
	var ids []uint32
	for _, m := range matches {
		id, err := strconv.ParseUint(string(m[1]), 10, 32)
		if err != nil || slices.Contains(ids, uint32(id)) {
			continue
		}
		if _, ok := a.wz.Maps[uint32(id)]; ok == false {
			continue
		}
		ids = append(ids, uint32(id))
	}
	slices.Sort(ids)
	return ids
}

func (a *SuiteActor) scriptIDs(kind string) ([]uint32, error) {
	names, err := a.scriptNames(kind)
	if err != nil {
		return nil, err
	}
	var ids []uint32
	for _, name := range names {
		id, err := strconv.ParseUint(name, 10, 32)
		if err != nil {
			continue
		}
		ids = append(ids, uint32(id))
	}
	slices.Sort(ids)
	return ids, nil
}

func (a *SuiteActor) scriptNames(kind string) ([]string, error) {
	files, err := os.ReadDir(filepath.Join(a.cfg.GameScriptDir, kind))
	if err != nil {
		return nil, err
	}
	var names []string
	for _, f := range files {
		if f.IsDir() || filepath.Ext(f.Name()) != ".lua" {
			continue
		}
		names = append(names, strings.TrimSuffix(f.Name(), ".lua"))
	}
	slices.Sort(names)
	return names, nil
}
