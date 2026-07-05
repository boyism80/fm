package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/boyism80/fm/services/game/wz"
	"gopkg.in/yaml.v3"
)

func writeYAML(path string, value any) error {
	data, err := yaml.Marshal(value)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o644)
}

func meta(wzPath string) Meta {
	return Meta{
		Generator: "tools/wzlookup",
		WzPath:    wzPath,
	}
}

func exportQuests(res *wz.Resources, wzPath, outDir string) error {
	quests := make(map[uint32]Quest, len(res.Quests))
	for id, q := range res.Quests {
		if q == nil {
			continue
		}
		quests[id] = exportQuest(q)
	}
	return writeYAML(filepath.Join(outDir, "quests.yaml"), QuestsFile{
		Meta:   meta(wzPath),
		Quests: quests,
	})
}

func exportQuest(q *wz.Quest) Quest {
	out := Quest{
		ID:       q.ID,
		Meta:     exportQuestMeta(q.Meta),
		Start:    exportQuestPhase(q.Start),
		Complete: exportQuestPhase(q.Complete),
	}
	return out
}

func exportQuestMeta(m wz.QuestMeta) QuestMeta {
	return QuestMeta{
		Name:            m.Name,
		Parent:          m.Parent,
		Order:           m.Order,
		Descriptions:    m.Descriptions,
		Area:            m.Area,
		AutoStart:       m.AutoStart,
		AutoPreComplete: m.AutoPreComplete,
		AutoComplete:    m.AutoComplete,
		AutoAccept:      m.AutoAccept,
		Blocked:         m.Blocked,
		ViewMedalItem:   m.ViewMedalItem,
		SelectedSkillID: m.SelectedSkillID,
	}
}

func exportQuestPhase(p wz.QuestPhase) QuestPhase {
	reqs := make([]QuestRequirement, 0, len(p.Requirements))
	for _, r := range p.Requirements {
		reqs = append(reqs, exportQuestRequirement(r))
	}
	acts := make([]QuestAction, 0, len(p.Actions))
	for _, a := range p.Actions {
		acts = append(acts, exportQuestAction(a))
	}
	return QuestPhase{
		Requirements: reqs,
		Actions:      acts,
	}
}

func exportQuestRequirement(r wz.QuestRequirement) QuestRequirement {
	items := make([]QuestItemCount, 0, len(r.Items))
	for _, item := range r.Items {
		items = append(items, QuestItemCount{ItemID: item.ItemID, Count: item.Count})
	}
	mobs := make([]QuestMobCount, 0, len(r.Mobs))
	for _, mob := range r.Mobs {
		mobs = append(mobs, QuestMobCount{MobID: mob.MobID, Count: mob.Count})
	}
	quests := make([]QuestStateRef, 0, len(r.Quests))
	for _, q := range r.Quests {
		quests = append(quests, QuestStateRef{QuestID: q.QuestID, State: q.State})
	}
	skills := make([]QuestSkillRef, 0, len(r.Skills))
	for _, s := range r.Skills {
		skills = append(skills, QuestSkillRef{SkillID: s.SkillID, Acquire: s.Acquire})
	}
	return QuestRequirement{
		Kind:        string(r.Kind),
		IntValue:    r.IntValue,
		StrValue:    r.StrValue,
		InfoStrings: r.InfoStrings,
		Jobs:        r.Classes,
		PetIDs:      r.PetIDs,
		Items:       items,
		Mobs:        mobs,
		Quests:      quests,
		Skills:      skills,
	}
}

func exportQuestAction(a wz.QuestAction) QuestAction {
	items := make([]QuestRewardItem, 0, len(a.Items))
	for _, item := range a.Items {
		items = append(items, QuestRewardItem{
			ItemID:     item.ItemID,
			Count:      item.Count,
			Job:        item.Job,
			JobEx:      item.JobEx,
			Gender:     item.Gender,
			Period:     item.Period,
			Prop:       item.Prop,
			DateExpire: item.DateExpire,
		})
	}
	skills := make([]QuestRewardSkill, 0, len(a.Skills))
	for _, s := range a.Skills {
		skills = append(skills, QuestRewardSkill{
			SkillID:     s.SkillID,
			SkillLevel:  s.SkillLevel,
			MasterLevel: s.MasterLevel,
			Jobs:        s.Jobs,
		})
	}
	quests := make([]QuestStateRef, 0, len(a.Quests))
	for _, q := range a.Quests {
		quests = append(quests, QuestStateRef{QuestID: q.QuestID, State: q.State})
	}
	return QuestAction{
		Kind:           string(a.Kind),
		IntValue:       a.IntValue,
		StrValue:       a.StrValue,
		ApplicableJobs: a.ApplicableJobs,
		Items:          items,
		Skills:         skills,
		Quests:         quests,
	}
}

func exportMaps(res *wz.Resources, wzPath, outDir string) error {
	maps := make(map[uint32]Map, len(res.Maps))
	for id, m := range res.Maps {
		if m == nil {
			continue
		}
		street, mapName := mapStringNames(res, id)
		maps[id] = exportMap(m, street, mapName, res)
	}
	return writeYAML(filepath.Join(outDir, "maps.yaml"), MapsFile{
		Meta: meta(wzPath),
		Maps: maps,
	})
}

func mapStringNames(res *wz.Resources, mapID uint32) (street, mapName string) {
	if res.Strings == nil {
		return "", ""
	}
	for _, regionMaps := range res.Strings.MapStrings {
		if data, ok := regionMaps[mapID]; ok {
			return data["streetName"], data["mapName"]
		}
	}
	return "", ""
}

func npcName(res *wz.Resources, npcID uint32) string {
	if res.Strings == nil {
		return ""
	}
	if data, ok := res.Strings.NpcStrings[npcID]; ok {
		return data["name"]
	}
	return ""
}

func mobName(res *wz.Resources, mobID uint32) string {
	if res.Strings == nil {
		return ""
	}
	if data, ok := res.Strings.MobStrings[mobID]; ok {
		return data["name"]
	}
	return ""
}

func exportMap(m *wz.Map, street, mapName string, res *wz.Resources) Map {
	out := Map{
		ID:           m.ID,
		Name:         m.Name,
		StreetName:   street,
		MapName:      mapName,
		ReturnMapID:  m.ReturnMapId,
		ForcedReturn: m.ForcedReturn,
		FieldLimit:   m.FieldLimit,
		IsTown:       m.IsTown,
		BGM:          m.BGM,
		MapMark:      m.MapMark,
		MapDesc:      m.MapDesc,
	}
	if len(m.NpcSpawns) > 0 {
		out.NpcSpawns = make(map[uint32]NpcSpawn, len(m.NpcSpawns))
		for spawnID, spawn := range m.NpcSpawns {
			if spawn.BaseSpawn == nil {
				continue
			}
			out.NpcSpawns[spawnID] = NpcSpawn{
				SpawnID:    spawnID,
				NpcID:      spawn.ID,
				NpcName:    npcName(res, spawn.ID),
				X:          spawn.Position.X,
				Y:          spawn.Position.Y,
				Foothold:   spawn.Foothold,
				Hide:       spawn.Hide,
				MobTimeSec: int(spawn.MobTime / time.Second),
			}
		}
	}
	if len(m.MobSpawns) > 0 {
		out.MobSpawns = make(map[uint32]MobSpawn, len(m.MobSpawns))
		for spawnID, spawn := range m.MobSpawns {
			if spawn.BaseSpawn == nil {
				continue
			}
			out.MobSpawns[spawnID] = MobSpawn{
				SpawnID:    spawnID,
				MobID:      spawn.ID,
				MobName:    mobName(res, spawn.ID),
				X:          spawn.Position.X,
				Y:          spawn.Position.Y,
				Foothold:   spawn.Foothold,
				Hide:       spawn.Hide,
				MobTimeSec: int(spawn.MobTime / time.Second),
			}
		}
	}
	if len(m.ReactorSpawns) > 0 {
		out.ReactorSpawns = make(map[uint32]ReactorSpawn, len(m.ReactorSpawns))
		for spawnID, spawn := range m.ReactorSpawns {
			out.ReactorSpawns[spawnID] = ReactorSpawn{
				SpawnID:    spawnID,
				ReactorID:  spawn.ReactorID,
				X:          spawn.Position.X,
				Y:          spawn.Position.Y,
				RespawnSec: int(spawn.RespawnDelay / time.Second),
				Name:       spawn.Name,
			}
		}
	}
	if len(m.Portals) > 0 {
		out.Portals = make(map[uint8]Portal, len(m.Portals))
		for portalID, portal := range m.Portals {
			out.Portals[portalID] = Portal{
				ID:          portal.ID,
				Name:        portal.Name,
				TargetMapID: portal.TargetMapId,
				Target:      portal.Target,
				X:           portal.Position.X,
				Y:           portal.Position.Y,
				ScriptName:  portal.ScriptName,
				Type:        portal.Type,
			}
		}
	}
	return out
}

func exportMobs(res *wz.Resources, wzPath, outDir string) error {
	mobs := make(map[uint32]Mob, len(res.Monsters))
	for id, m := range res.Monsters {
		if m == nil {
			continue
		}
		mobs[id] = exportMob(m, mobName(res, id))
	}
	return writeYAML(filepath.Join(outDir, "mobs.yaml"), MobsFile{
		Meta: meta(wzPath),
		Mobs: mobs,
	})
}

func exportMob(m *wz.Mob, name string) Mob {
	skills := make([]MobSkillSlot, 0, len(m.Skills))
	for _, s := range m.Skills {
		skills = append(skills, MobSkillSlot{
			Slot:    s.Slot,
			SkillID: s.SkillID,
			Level:   s.Level,
			Action:  s.Action,
		})
	}
	attacks := make([]MobAttack, 0, len(m.Attacks))
	for _, a := range m.Attacks {
		attacks = append(attacks, MobAttack{
			Index:        a.Index,
			DeadlyAttack: a.DeadlyAttack,
			MpBurn:       a.MpBurn,
			MpCon:        a.MpCon,
			DiseaseSkill: a.DiseaseSkill,
			DiseaseLevel: a.DiseaseLevel,
			AttackAfter:  a.AttackAfter,
			PADamage:     a.PADamage,
			MADamage:     a.MADamage,
			Magic:        a.Magic,
			RangeR:       a.RangeR,
		})
	}
	var banish *MobBanish
	if m.Banish != nil {
		banish = &MobBanish{
			Message: m.Banish.Message,
			MapID:   m.Banish.MapID,
			Portal:  m.Banish.Portal,
		}
	}
	return Mob{
		ID:                    m.ID,
		Name:                  name,
		BodyAttack:            m.BodyAttack,
		Level:                 m.Level,
		MaxHP:                 m.MaxHP,
		MaxMP:                 m.MaxMP,
		Speed:                 m.Speed,
		PADamage:              m.PADamage,
		PDDamage:              m.PDDamage,
		MADamage:              m.MADamage,
		MDDamage:              m.MDDamage,
		ACC:                   m.ACC,
		EVA:                   m.EVA,
		EXP:                   m.EXP,
		Undead:                m.Undead,
		Pushed:                m.Pushed,
		Boss:                  m.Boss,
		FfaLoot:               m.FfaLoot,
		ExplosiveReward:       m.ExplosiveReward,
		FS:                    m.FS,
		SummonType:            m.SummonType,
		MobType:               m.MobType,
		Link:                  m.Link,
		ElemResist:            m.ElemResist,
		Skills:                skills,
		Attacks:               attacks,
		Banish:                banish,
		Revives:               m.Revives,
		RemoveAfter:           m.RemoveAfter,
		SelfDestructionAction: m.SelfDestructionAction,
		HpTagColor:            m.HpTagColor,
		HpTagBgColor:          m.HpTagBgColor,
	}
}

func exportReactors(res *wz.Resources, wzPath, outDir string) error {
	reactors := make(map[uint32]Reactor, len(res.Reactors))
	for id, r := range res.Reactors {
		if r == nil {
			continue
		}
		reactors[id] = exportReactor(r)
	}
	return writeYAML(filepath.Join(outDir, "reactors.yaml"), ReactorsFile{
		Meta:     meta(wzPath),
		Reactors: reactors,
	})
}

func exportReactor(r *wz.Reactor) Reactor {
	states := make(map[byte]*ReactorEvent, len(r.States))
	for stateID, event := range r.States {
		if event == nil {
			continue
		}
		states[stateID] = &ReactorEvent{
			Type:         int(event.Type),
			NextState:    event.NextState,
			TimeOut:      event.TimeOut,
			ItemID:       event.ItemID,
			ItemQuantity: event.ItemQuantity,
			LTX:          event.LT.X,
			LTY:          event.LT.Y,
			RBX:          event.RB.X,
			RBY:          event.RB.Y,
			TouchFlag:    event.TouchFlag,
			HasClickArea: event.HasClickArea,
		}
	}
	return Reactor{
		ID:              r.ID,
		Link:            r.Info.Link,
		ActivateByTouch: r.Info.ActivateByTouch,
		Action:          r.Action,
		States:          states,
	}
}

func exportShops(res *wz.Resources, wzPath, outDir string) error {
	shops := make(map[uint32]Shop, len(res.Shops))
	for npcID, shop := range res.Shops {
		if shop == nil {
			continue
		}
		items := make([]ShopItem, 0, len(shop.Items))
		for _, item := range shop.Items {
			items = append(items, ShopItem{
				ItemID:    item.ItemID,
				Price:     item.Price,
				Period:    item.Period,
				Stock:     item.Stock,
				UnitPrice: item.UnitPrice,
			})
		}
		shops[npcID] = Shop{
			NpcID:   shop.NpcID,
			NpcName: npcName(res, shop.NpcID),
			Items:   items,
		}
	}
	return writeYAML(filepath.Join(outDir, "npc_shops.yaml"), ShopsFile{
		Meta:  meta(wzPath),
		Shops: shops,
	})
}

func exportStrings(res *wz.Resources, wzPath, outDir string) error {
	out := StringsFile{Meta: meta(wzPath)}
	if res.Strings == nil {
		return writeYAML(filepath.Join(outDir, "strings.yaml"), out)
	}
	if len(res.Strings.NpcStrings) > 0 {
		out.Npcs = make(map[uint32]Name, len(res.Strings.NpcStrings))
		for id, data := range res.Strings.NpcStrings {
			name := data["name"]
			if name == "" {
				continue
			}
			out.Npcs[id] = Name{Name: name}
		}
	}
	if len(res.Strings.MobStrings) > 0 {
		out.Mobs = make(map[uint32]Name, len(res.Strings.MobStrings))
		for id, data := range res.Strings.MobStrings {
			name := data["name"]
			if name == "" {
				continue
			}
			out.Mobs[id] = Name{Name: name}
		}
	}
	if len(res.Strings.MapStrings) > 0 {
		out.Maps = make(map[uint32]MapName)
		for _, regionMaps := range res.Strings.MapStrings {
			for id, data := range regionMaps {
				street := data["streetName"]
				mapName := data["mapName"]
				if street == "" && mapName == "" {
					continue
				}
				out.Maps[id] = MapName{
					StreetName: street,
					MapName:    mapName,
				}
			}
		}
	}
	return writeYAML(filepath.Join(outDir, "strings.yaml"), out)
}

func exportDropsFile(res *wz.Resources, wzPath, outDir string) error {
	mobDrops := make(map[uint32][]Drop, len(res.MobDrops))
	for id, drops := range res.MobDrops {
		mobDrops[id] = convertDrops(drops)
	}
	reactorDrops := make(map[uint32][]Drop, len(res.ReactorDrops))
	for id, drops := range res.ReactorDrops {
		reactorDrops[id] = convertDrops(drops)
	}
	return writeYAML(filepath.Join(outDir, "drops.yaml"), DropsFile{
		Meta:         meta(wzPath),
		MobDrops:     mobDrops,
		ReactorDrops: reactorDrops,
	})
}

func convertDrops(drops []wz.Drop) []Drop {
	out := make([]Drop, 0, len(drops))
	for _, d := range drops {
		out = append(out, Drop{
			ItemID:  d.Item,
			Money:   d.Money,
			Prob:    d.Prob,
			Min:     d.Min,
			Max:     d.Max,
			QuestID: d.QuestID,
		})
	}
	return out
}

func exportAll(res *wz.Resources, wzPath, outDir, only string) error {
	type job struct {
		name string
		fn   func(*wz.Resources, string, string) error
	}
	jobs := []job{
		{name: "quests", fn: exportQuests},
		{name: "maps", fn: exportMaps},
		{name: "mobs", fn: exportMobs},
		{name: "reactors", fn: exportReactors},
		{name: "shops", fn: exportShops},
		{name: "strings", fn: exportStrings},
		{name: "drops", fn: exportDropsFile},
	}
	for _, j := range jobs {
		if only != "all" && only != j.name {
			continue
		}
		fmt.Printf("writing %s.yaml...\n", j.name)
		if err := j.fn(res, wzPath, outDir); err != nil {
			return fmt.Errorf("%s: %w", j.name, err)
		}
	}
	return nil
}
