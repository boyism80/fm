package actor

import (
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"

	"github.com/boyism80/fm/services/game/wz"
	lua "github.com/yuin/gopher-lua"
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
