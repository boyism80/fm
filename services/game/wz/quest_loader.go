package wz

import (
	"encoding/xml"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

func loadQuest(path string) (*Quest, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var root node
	if err := xml.NewDecoder(file).Decode(&root); err != nil {
		return nil, err
	}

	id, err := strconv.ParseUint(strings.TrimSuffix(root.Name, ".img"), 10, 32)
	if err != nil {
		return nil, fmt.Errorf("invalid quest id in %s: %w", path, err)
	}

	quest := &Quest{
		ID: uint32(id),
	}

	if info := root.find("QuestInfo"); info != nil {
		quest.Meta = parseQuestMeta(info)
	}

	if check := root.find("Check"); check != nil {
		if start := check.find("0"); start != nil {
			quest.Start.Requirements = parseQuestRequirements(start)
			if quest.Start.Requirements.HasInterval || quest.Start.Requirements.DayByDay {
				quest.Meta.Repeatable = true
			}
		}
		if complete := check.find("1"); complete != nil {
			quest.Complete.Requirements = parseQuestRequirements(complete)
		}
	}

	if act := root.find("Act"); act != nil {
		if start := act.find("0"); start != nil {
			quest.Start.Actions = parseQuestActions(start)
		}
		if complete := act.find("1"); complete != nil {
			quest.Complete.Actions = parseQuestActions(complete)
		}
	}

	return quest, nil
}

func parseQuestMeta(info *node) QuestMeta {
	meta := QuestMeta{
		Descriptions: make(map[int]string),
	}
	meta.Name = nodeString(info, "name")
	meta.Parent = nodeString(info, "parent")
	meta.Order = nodeInt(info, "order", 0)
	meta.Area = nodeInt(info, "area", 0)
	meta.AutoStart = nodeInt(info, "autoStart", 0) > 0
	meta.AutoPreComplete = nodeInt(info, "autoPreComplete", 0) > 0
	meta.AutoComplete = nodeInt(info, "autoComplete", 0) > 0
	meta.AutoAccept = nodeInt(info, "autoAccept", 0) > 0
	meta.Blocked = nodeInt(info, "blocked", 0) > 0
	meta.ViewMedalItem = nodeInt(info, "viewMedalItem", 0)
	meta.SelectedSkillID = nodeInt(info, "selectedSkillID", 0)
	meta.TimeLimit = nodeInt(info, "timeLimit", 0)
	meta.TimeLimit2 = nodeInt(info, "timeLimit2", 0)

	for _, strField := range info.Strings {
		if idx, err := strconv.Atoi(strField.Name); err == nil {
			meta.Descriptions[idx] = strField.Value
		}
	}

	return meta
}

func parseQuestRequirements(phase *node) QuestRequirements {
	var reqs QuestRequirements
	if phase == nil {
		return reqs
	}

	for _, child := range phase.Children {
		applyQuestRequirementChild(&reqs, &child)
	}
	for _, intField := range phase.Ints {
		applyQuestRequirementInt(&reqs, intField.Name, intField.Value)
	}
	for _, strField := range phase.Strings {
		applyQuestRequirementStr(&reqs, strField.Name, strField.Value)
	}
	return reqs
}

func applyQuestRequirementChild(reqs *QuestRequirements, child *node) {
	switch child.Name {
	case "job":
		reqs.Job = collectChildIntValues(child)
	case "pet":
		reqs.Pet = collectChildUint32IDs(child)
	case "item":
		reqs.Item = parseQuestItemCounts(child)
	case "mob":
		reqs.Mob = parseQuestMobCounts(child)
	case "mbcard":
		reqs.MBCard = parseQuestItemCountsWithField(child, "min")
	case "quest":
		reqs.Quest = parseQuestStateRefs(child)
	case "skill":
		reqs.Skill = parseQuestSkillRefs(child)
	case "start":
		reqs.Start = nodeStringValue(child)
	case "end":
		reqs.End = nodeStringValue(child)
	case "startscript":
		reqs.StartScript = nodeStringValue(child)
	case "endscript":
		reqs.EndScript = nodeStringValue(child)
	case "fieldEnter":
		reqs.FieldEnter = nodeInt(child, "0", nodeInt(child, "", 0))
		if reqs.FieldEnter == 0 {
			reqs.FieldEnter = firstChildIntValue(child)
		}
	case "info":
		reqs.Info = collectInfoStrings(child)
	case "npc":
		reqs.NPC = uint32(firstIntValue(child))
	case "lvmin":
		reqs.LevelMin = firstIntValue(child)
	case "lvmax":
		reqs.LevelMax = firstIntValue(child)
	case "level":
		reqs.Level = firstIntValue(child)
	case "pop":
		reqs.Pop = firstIntValue(child)
	case "interval":
		reqs.HasInterval = true
		reqs.Interval = firstIntValue(child)
	case "questComplete":
		reqs.QuestComplete = firstIntValue(child)
	case "pettamenessmin":
		reqs.PetTamenessMin = firstIntValue(child)
	case "mbmin":
		reqs.MBMin = firstIntValue(child)
	case "subJobFlags":
		reqs.SubJobFlags = firstIntValue(child)
	case "dayByDay":
		reqs.DayByDay = true
	case "normalAutoStart":
		reqs.NormalAutoStart = true
	case "partyQuest_S":
		reqs.PartyQuestS = firstIntValue(child)
	case "infoNumber":
		reqs.InfoNumber = firstIntValue(child)
	case "endmeso":
		reqs.EndMeso = firstIntValue(child)
	case "equipAllNeed":
		reqs.EquipAllNeed = firstIntValue(child)
	case "equipSelectNeed":
		reqs.EquipSelectNeed = firstIntValue(child)
	case "premium":
		reqs.Premium = firstIntValue(child) > 0
	case "tamingmoblevelmin":
		reqs.TamingMobLevelMin = firstIntValue(child)
	case "worldmin":
		reqs.WorldMin = nodeStringValue(child)
		if reqs.WorldMin == "" {
			reqs.WorldMin = strconv.Itoa(firstIntValue(child))
		}
	case "worldmax":
		reqs.WorldMax = nodeStringValue(child)
		if reqs.WorldMax == "" {
			reqs.WorldMax = strconv.Itoa(firstIntValue(child))
		}
	case "buff":
		reqs.Buff = nodeStringValue(child)
		if reqs.Buff == "" {
			reqs.Buff = strconv.Itoa(firstIntValue(child))
		}
	case "exceptbuff":
		reqs.ExceptBuff = nodeStringValue(child)
		if reqs.ExceptBuff == "" {
			reqs.ExceptBuff = strconv.Itoa(firstIntValue(child))
		}
	default:
		if str := nodeStringValue(child); str != "" {
			applyQuestRequirementStr(reqs, child.Name, str)
		} else {
			applyQuestRequirementInt(reqs, child.Name, firstIntValue(child))
		}
	}
}

func applyQuestRequirementInt(reqs *QuestRequirements, kind string, value int) {
	switch kind {
	case "npc":
		reqs.NPC = uint32(value)
	case "lvmin":
		reqs.LevelMin = value
	case "lvmax":
		reqs.LevelMax = value
	case "level":
		reqs.Level = value
	case "pop":
		reqs.Pop = value
	case "interval":
		reqs.HasInterval = true
		reqs.Interval = value
	case "fieldEnter":
		reqs.FieldEnter = value
	case "questComplete":
		reqs.QuestComplete = value
	case "pettamenessmin":
		reqs.PetTamenessMin = value
	case "mbmin":
		reqs.MBMin = value
	case "subJobFlags":
		reqs.SubJobFlags = value
	case "dayByDay":
		reqs.DayByDay = true
	case "normalAutoStart":
		reqs.NormalAutoStart = true
	case "partyQuest_S":
		reqs.PartyQuestS = value
	case "infoNumber":
		reqs.InfoNumber = value
	case "endmeso":
		reqs.EndMeso = value
	case "equipAllNeed":
		reqs.EquipAllNeed = value
	case "equipSelectNeed":
		reqs.EquipSelectNeed = value
	case "premium":
		reqs.Premium = value > 0
	case "tamingmoblevelmin":
		reqs.TamingMobLevelMin = value
	case "worldmin":
		reqs.WorldMin = strconv.Itoa(value)
	case "worldmax":
		reqs.WorldMax = strconv.Itoa(value)
	}
}

func applyQuestRequirementStr(reqs *QuestRequirements, kind string, value string) {
	switch kind {
	case "start":
		reqs.Start = value
	case "end":
		reqs.End = value
	case "startscript":
		reqs.StartScript = value
	case "endscript":
		reqs.EndScript = value
	case "worldmin":
		reqs.WorldMin = value
	case "worldmax":
		reqs.WorldMax = value
	case "buff":
		reqs.Buff = value
	case "exceptbuff":
		reqs.ExceptBuff = value
	}
}

func parseQuestActions(phase *node) QuestActions {
	var actions QuestActions
	if phase == nil {
		return actions
	}

	for _, child := range phase.Children {
		kind := child.Name
		switch kind {
		case "item":
			actions.Item = append(actions.Item, parseQuestActionItems(&child)...)
		case "skill":
			actions.Skills = append(actions.Skills, parseQuestActionSkills(&child)...)
			actions.SkillJobs = collectQuestActionClasses(&child)
		case "quest":
			actions.Quests = parseQuestStateRefs(&child)
		case "sp":
			actions.SP = nodeInt(child.find("0"), "sp_value", firstIntValue(&child))
			actions.SPJobs = collectQuestActionClasses(&child)
		case "info":
			actions.Info = nodeStringValue(&child)
		case "npcAct":
			actions.NPCAct = nodeStringValue(&child)
		case "exp":
			actions.Exp = firstIntValue(&child)
		case "money":
			actions.Money = firstIntValue(&child)
		case "pop":
			actions.Pop = firstIntValue(&child)
		case "nextQuest":
			actions.NextQuest = uint32(firstIntValue(&child))
		case "buffItemID":
			actions.BuffItemID = uint32(firstIntValue(&child))
		case "infoNumber":
			actions.InfoNumber = uint32(firstIntValue(&child))
		case "npc":
			actions.NPC = firstIntValue(&child)
		case "pettameness":
			actions.PetTameness = firstIntValue(&child)
		case "petspeed":
			actions.PetSpeed = firstIntValue(&child)
		case "map":
			actions.Map = firstIntValue(&child)
		case "job":
			actions.Job = firstIntValue(&child)
		case "lvmin":
			actions.LvMin = firstIntValue(&child)
		case "lvmax":
			actions.LvMax = firstIntValue(&child)
		case "fieldEnter":
			actions.FieldEnter = firstIntValue(&child)
		case "interval":
			actions.Interval = firstIntValue(&child)
		case "ask":
			actions.Ask = firstIntValue(&child)
		case "stop":
			actions.Stop = firstIntValue(&child)
		case "message":
			actions.Message = nodeStringValue(&child)
		case "start":
			actions.Start = nodeStringValue(&child)
		case "end":
			actions.End = nodeStringValue(&child)
		default:
			if str := nodeStringValue(&child); str != "" {
				if actions.Say == nil {
					actions.Say = make(map[string]string)
				}
				actions.Say[kind] = str
			} else if val := firstIntValue(&child); val != 0 {
				if actions.Say == nil {
					actions.Say = make(map[string]string)
				}
				actions.Say[kind] = strconv.Itoa(val)
			} else {
				if actions.Say == nil {
					actions.Say = make(map[string]string)
				}
				actions.Say[kind] = ""
			}
		}
	}

	for _, intField := range phase.Ints {
		applyQuestActionInt(&actions, intField.Name, intField.Value)
	}

	for _, strField := range phase.Strings {
		applyQuestActionStr(&actions, strField.Name, strField.Value)
	}

	return actions
}

func applyQuestActionInt(actions *QuestActions, kind string, value int) {
	switch kind {
	case "exp":
		actions.Exp = value
	case "money":
		actions.Money = value
	case "pop":
		actions.Pop = value
	case "nextQuest":
		actions.NextQuest = uint32(value)
	case "buffItemID":
		actions.BuffItemID = uint32(value)
	case "infoNumber":
		actions.InfoNumber = uint32(value)
	case "npc":
		actions.NPC = value
	case "sp":
		actions.SP = value
	case "pettameness":
		actions.PetTameness = value
	case "petspeed":
		actions.PetSpeed = value
	case "map":
		actions.Map = value
	case "job":
		actions.Job = value
	case "lvmin":
		actions.LvMin = value
	case "lvmax":
		actions.LvMax = value
	case "fieldEnter":
		actions.FieldEnter = value
	case "interval":
		actions.Interval = value
	case "ask":
		actions.Ask = value
	case "stop":
		actions.Stop = value
	default:
		if actions.Say == nil {
			actions.Say = make(map[string]string)
		}
		actions.Say[kind] = strconv.Itoa(value)
	}
}

func applyQuestActionStr(actions *QuestActions, kind string, value string) {
	switch kind {
	case "info":
		actions.Info = value
	case "npcAct":
		actions.NPCAct = value
	case "message":
		actions.Message = value
	case "start":
		actions.Start = value
	case "end":
		actions.End = value
	default:
		if actions.Say == nil {
			actions.Say = make(map[string]string)
		}
		actions.Say[kind] = value
	}
}

func parseQuestItemCounts(n *node) map[uint32]int {
	return parseQuestItemCountsWithField(n, "count")
}

func parseQuestItemCountsWithField(n *node, countField string) map[uint32]int {
	if n == nil {
		return nil
	}
	out := make(map[uint32]int, len(n.Children))
	for _, child := range n.Children {
		itemID := uint32(nodeInt(&child, "id", 0))
		count := nodeInt(&child, countField, 0)
		if itemID == 0 && count == 0 {
			continue
		}
		out[itemID] = count
	}
	if len(out) == 0 {
		return nil
	}
	return out
}

func parseQuestMobCounts(n *node) map[uint32]int {
	if n == nil {
		return nil
	}
	out := make(map[uint32]int, len(n.Children))
	for _, child := range n.Children {
		mobID := uint32(nodeInt(&child, "id", 0))
		count := nodeInt(&child, "count", 0)
		if mobID == 0 && count == 0 {
			continue
		}
		out[mobID] = count
	}
	if len(out) == 0 {
		return nil
	}
	return out
}

func parseQuestStateRefs(n *node) map[uint32]QuestStatus {
	if n == nil {
		return nil
	}
	out := make(map[uint32]QuestStatus, len(n.Children))
	for _, child := range n.Children {
		questID := uint32(nodeInt(&child, "id", 0))
		state := nodeInt(&child, "state", 0)
		if questID == 0 {
			continue
		}
		out[questID] = QuestStatus(state)
	}
	if len(out) == 0 {
		return nil
	}
	return out
}

func parseQuestSkillRefs(n *node) map[uint32]int {
	if n == nil {
		return nil
	}
	out := make(map[uint32]int, len(n.Children))
	for _, child := range n.Children {
		skillID := uint32(nodeInt(&child, "id", 0))
		acquire := nodeInt(&child, "acquire", 0)
		if skillID == 0 {
			continue
		}
		out[skillID] = acquire
	}
	if len(out) == 0 {
		return nil
	}
	return out
}

func parseQuestActionItems(n *node) []QuestActionItem {
	if n == nil {
		return nil
	}
	out := make([]QuestActionItem, 0, len(n.Children))
	for _, child := range n.Children {
		item := QuestActionItem{
			ItemID:     uint32(nodeInt(&child, "id", 0)),
			Count:      nodeInt(&child, "count", 0),
			Class:      nodeInt(&child, "job", -1),
			ClassEx:    nodeInt(&child, "jobEx", -1),
			Gender:     nodeInt(&child, "gender", 2),
			Period:     nodeInt(&child, "period", 0),
			Prop:       QuestRewardPropAlways,
			DateExpire: nodeString(&child, "dateExpire"),
		}
		if prop, ok := nodeIntOptional(&child, "prop"); ok {
			item.Prop = QuestRewardProp(prop)
		}
		if item.ItemID == 0 && item.Count == 0 {
			continue
		}
		out = append(out, item)
	}
	return out
}

func parseQuestActionSkills(n *node) []QuestActionSkill {
	if n == nil {
		return nil
	}
	out := make([]QuestActionSkill, 0, len(n.Children))
	for _, child := range n.Children {
		skill := QuestActionSkill{
			SkillID:     uint32(nodeInt(&child, "id", 0)),
			SkillLevel:  nodeInt(&child, "skillLevel", 0),
			MasterLevel: nodeInt(&child, "masterLevel", 0),
		}
		if classes := child.find("job"); classes != nil {
			skill.Classes = collectChildIntValues(classes)
		}
		if skill.SkillID == 0 {
			continue
		}
		out = append(out, skill)
	}
	return out
}

func collectQuestActionClasses(n *node) []int {
	if n == nil {
		return nil
	}
	var classes []int
	index := 0
	for {
		entry := n.find(strconv.Itoa(index))
		if entry == nil {
			break
		}
		if classNode := entry.find("job"); classNode != nil {
			classes = append(classes, collectChildIntValues(classNode)...)
		}
		index++
	}
	if len(classes) > 0 {
		return classes
	}
	if classNode := n.find("job"); classNode != nil {
		return collectChildIntValues(classNode)
	}
	return nil
}

func collectChildUint32IDs(n *node) []uint32 {
	if n == nil {
		return nil
	}
	out := make([]uint32, 0, len(n.Children))
	for _, child := range n.Children {
		id := uint32(nodeInt(&child, "id", 0))
		if id > 0 {
			out = append(out, id)
		}
	}
	return out
}

func collectChildIntValues(n *node) []int {
	if n == nil {
		return nil
	}
	out := make([]int, 0, len(n.Children)+len(n.Ints))
	for _, child := range n.Children {
		if val, ok := nodeIntOptional(&child, ""); ok {
			out = append(out, val)
			continue
		}
		if val := nodeInt(&child, "id", 0); val != 0 {
			out = append(out, val)
			continue
		}
		if val := firstIntValue(&child); val != 0 {
			out = append(out, val)
		}
	}
	for _, intField := range n.Ints {
		out = append(out, intField.Value)
	}
	return out
}

func collectInfoStrings(n *node) []string {
	if n == nil {
		return nil
	}
	type indexedString struct {
		index int
		value string
	}
	pairs := make([]indexedString, 0, len(n.Strings)+len(n.Children))
	for _, strField := range n.Strings {
		index, err := strconv.Atoi(strField.Name)
		if err != nil {
			continue
		}
		pairs = append(pairs, indexedString{index: index, value: strField.Value})
	}
	for _, child := range n.Children {
		index, err := strconv.Atoi(child.Name)
		if err != nil {
			continue
		}
		value := nodeStringValue(&child)
		pairs = append(pairs, indexedString{index: index, value: value})
	}
	if len(pairs) == 0 {
		return nil
	}
	sort.Slice(pairs, func(i, j int) bool {
		return pairs[i].index < pairs[j].index
	})
	out := make([]string, len(pairs))
	for i, pair := range pairs {
		out[i] = pair.value
	}
	return out
}

func firstIntValue(n *node) int {
	if n == nil {
		return 0
	}
	if len(n.Ints) > 0 {
		return n.Ints[0].Value
	}
	if n.Value != "" {
		if val, err := strconv.Atoi(n.Value); err == nil {
			return val
		}
	}
	for _, child := range n.Children {
		if val := nodeInt(&child, "", 0); val != 0 {
			return val
		}
		if val := firstIntValue(&child); val != 0 {
			return val
		}
	}
	return 0
}

func firstChildIntValue(n *node) int {
	if n == nil {
		return 0
	}
	for _, child := range n.Children {
		if val := firstIntValue(&child); val != 0 {
			return val
		}
	}
	return 0
}

func nodeStringValue(n *node) string {
	if n == nil {
		return ""
	}
	if s := nodeString(n, ""); s != "" {
		return s
	}
	for _, strField := range n.Strings {
		if strField.Value != "" {
			return strField.Value
		}
	}
	if n.Value != "" {
		return n.Value
	}
	return ""
}
