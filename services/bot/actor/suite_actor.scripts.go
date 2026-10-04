package actor

import (
	"bytes"
	"maps"
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
	questCall    = regexp.MustCompile(`(?:quest\(\s*(\d+)\s*\)|\bquest\s*=\s*(\d+))`)
	itemChange   = regexp.MustCompile(`(?:\[\s*(\d{7})\s*\]\s*=|(?:mkitem|rmitem)\(\s*(\d{7})\b)`)
	mesoChange   = regexp.MustCompile(`\bmeso\s*=`)
	expChange    = regexp.MustCompile(`\bexp\s*=`)
	mapLiteral   = regexp.MustCompile(`:map\(\s*(\d{9})\b`)
	mapVariable  = regexp.MustCompile(`:map\(\s*[^\d\s)]`)
	mapAnyNumber = regexp.MustCompile(`(?:^|[^#m\d])(\d{9})\b`)
	mapKey       = regexp.MustCompile(`^\s*\[\s*(\d{9})\s*\]\s*=`)
	mapCompare   = regexp.MustCompile(`(?:==|~=|<=|>=|<|>)\s*\d{9}\b|\b\d{9}\s*(?:==|~=|<=|>=|<|>)`)
	currentMap   = regexp.MustCompile(`:map\(\)`)
	levelCheck   = regexp.MustCompile(`:level\(\)`)
	mapSpread    = regexp.MustCompile(`(\d{9})\s*\+\s*math\.random\(\s*\d+\s*,\s*(\d+)\s*\)`)
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
				spots, ok := placed[id]
				if ok == false {
					missing.Append(lua.LNumber(id))
					continue
				}
				entry := a.npcEntry(L, id, spots[0])
				list := L.NewTable()
				for _, at := range spots {
					list.Append(a.npcEntry(L, id, at))
				}
				entry.RawSetString("spots", list)
				found.Append(entry)
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
				} else if spots, ok := placed[quest.Start.Requirements.NPC]; ok {
					entry.RawSetString("start", a.npcEntry(L, quest.Start.Requirements.NPC, spots[0]))
				}
				if spots, ok := placed[quest.Complete.Requirements.NPC]; ok {
					entry.RawSetString("finish", a.npcEntry(L, quest.Complete.Requirements.NPC, spots[0]))
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
			var changes [][][]byte
			for _, m := range itemChange.FindAllSubmatch(data, -1) {
				if len(m[1]) == 0 {
					m[1] = m[2]
				}
				changes = append(changes, m)
			}
			var quests []uint32
			for _, m := range questCall.FindAllSubmatch(data, -1) {
				if len(m[1]) == 0 {
					m[1] = m[2]
				}
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
			refs.RawSetString("maps", a.mapRefs(L, data))
			refs.RawSetString("changes", a.numberList(L, a.itemIDs(changes)))
			refs.RawSetString("meso", lua.LBool(mesoChange.Match(data)))
			refs.RawSetString("exp", lua.LBool(expChange.Match(data)))
			refs.RawSetString("here", lua.LBool(currentMap.Match(data)))
			refs.RawSetString("level", lua.LBool(levelCheck.Match(data)))
			refs.RawSetString("party", lua.LBool(bytes.Contains(data, []byte("party"))))
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
		"has_quest": func(L *lua.LState) int {
			L.Push(lua.LBool(a.wz.GetQuest(uint32(L.CheckInt(1))) != nil))
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

func (a *SuiteActor) npcSpots() map[uint32][]npcSpot {
	placed := make(map[uint32][]npcSpot)
	for mapID, m := range a.wz.Maps {
		keys := slices.Sorted(maps.Keys(m.NpcSpawns))
		var ids []uint32
		for _, key := range keys {
			spawn := m.NpcSpawns[key]
			if spawn.BaseSpawn == nil || spawn.Hide || slices.Contains(ids, spawn.ID) {
				continue
			}
			ids = append(ids, spawn.ID)
			placed[spawn.ID] = append(placed[spawn.ID], npcSpot{mapID: mapID, spawn: spawn.BaseSpawn})
		}
	}
	for _, spots := range placed {
		slices.SortFunc(spots, func(x, y npcSpot) int {
			return int(x.mapID) - int(y.mapID)
		})
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

func (a *SuiteActor) mapRefs(L *lua.LState, data []byte) *lua.LTable {
	pattern := mapLiteral
	if mapVariable.Match(data) {
		pattern = mapAnyNumber
	}
	random := bytes.Contains(data, []byte("math.random"))

	refs := L.NewTable()
	var found []uint32
	for i, line := range bytes.Split(data, []byte("\n")) {
		var from []uint32
		if m := mapKey.FindSubmatch(line); m != nil {
			from = a.mapIDs([][][]byte{m})
			line = line[len(m[0]):]
		}
		ids := a.mapIDs(pattern.FindAllSubmatch(mapCompare.ReplaceAll(line, nil), -1))
		for _, id := range ids {
			if slices.Contains(found, id) {
				continue
			}
			found = append(found, id)
			ref := L.NewTable()
			ref.RawSetString("id", lua.LNumber(id))
			if len(from) > 0 {
				ref.RawSetString("from", lua.LNumber(from[0]))
			}
			if random && len(ids) > 1 {
				ref.RawSetString("group", lua.LNumber(i+1))
			}
			for _, m := range mapSpread.FindAllSubmatch(line, -1) {
				if string(m[1]) == strconv.FormatUint(uint64(id), 10) {
					spread, _ := strconv.Atoi(string(m[2]))
					ref.RawSetString("spread", lua.LNumber(spread))
				}
			}
			refs.Append(ref)
		}
	}
	return refs
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
