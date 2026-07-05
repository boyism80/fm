package wz

import (
	"encoding/xml"
	"fmt"
	"os"
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

	for _, strField := range info.Strings {
		if idx, err := strconv.Atoi(strField.Name); err == nil {
			meta.Descriptions[idx] = strField.Value
		}
	}

	return meta
}

func parseQuestRequirements(phase *node) []QuestRequirement {
	if phase == nil {
		return nil
	}

	reqs := make([]QuestRequirement, 0, len(phase.Children))
	for _, child := range phase.Children {
		kind := QuestRequirementKind(child.Name)
		req := QuestRequirement{Kind: kind}

		switch kind {
		case QuestReqJob:
			req.Jobs = collectChildIntValues(&child)
		case QuestReqPet:
			req.PetIDs = collectChildUint32IDs(&child)
		case QuestReqItem, QuestReqMob, QuestReqMBCard:
			if kind == QuestReqItem {
				req.Items = parseQuestItemCounts(&child)
			} else if kind == QuestReqMob {
				req.Mobs = parseQuestMobCounts(&child)
			} else {
				req.Items = parseQuestItemCountsWithField(&child, "min")
			}
		case QuestReqQuest:
			req.Quests = parseQuestStateRefs(&child)
		case QuestReqSkill:
			req.Skills = parseQuestSkillRefs(&child)
		case QuestReqTimeStart, QuestReqTimeEnd, QuestReqStartScript, QuestReqEndScript:
			req.StrValue = nodeStringValue(&child)
		case QuestReqFieldEnter:
			req.IntValue = nodeInt(&child, "0", nodeInt(&child, "", 0))
			if req.IntValue == 0 {
				req.IntValue = firstChildIntValue(&child)
			}
		default:
			req.IntValue = firstIntValue(&child)
		}

		reqs = append(reqs, req)
	}

	for _, intField := range phase.Ints {
		kind := QuestRequirementKind(intField.Name)
		reqs = append(reqs, QuestRequirement{
			Kind:     kind,
			IntValue: intField.Value,
		})
	}

	for _, strField := range phase.Strings {
		kind := QuestRequirementKind(strField.Name)
		reqs = append(reqs, QuestRequirement{
			Kind:     kind,
			StrValue: strField.Value,
		})
	}

	return reqs
}

func parseQuestActions(phase *node) []QuestAction {
	if phase == nil {
		return nil
	}

	acts := make([]QuestAction, 0, len(phase.Children))
	for _, child := range phase.Children {
		kind := QuestActionKind(child.Name)
		act := QuestAction{Kind: kind}

		switch kind {
		case QuestActItem:
			act.Items = parseQuestRewardItems(&child)
		case QuestActSkill:
			act.Skills = parseQuestRewardSkills(&child)
		case QuestActQuest:
			act.Quests = parseQuestStateRefs(&child)
		case QuestActSP:
			act.IntValue = nodeInt(child.find("0"), "sp_value", firstIntValue(&child))
		case QuestActInfo, QuestActNPCAct:
			act.StrValue = nodeStringValue(&child)
		default:
			act.IntValue = firstIntValue(&child)
		}

		if kind == QuestActSP || kind == QuestActSkill {
			act.ApplicableJobs = collectQuestActionJobs(&child)
		} else if jobs := child.find("job"); jobs != nil {
			act.ApplicableJobs = collectChildIntValues(jobs)
		}

		acts = append(acts, act)
	}

	for _, intField := range phase.Ints {
		kind := QuestActionKind(intField.Name)
		acts = append(acts, QuestAction{
			Kind:     kind,
			IntValue: intField.Value,
		})
	}

	for _, strField := range phase.Strings {
		kind := QuestActionKind(strField.Name)
		acts = append(acts, QuestAction{
			Kind:     kind,
			StrValue: strField.Value,
		})
	}

	return acts
}

func parseQuestItemCounts(n *node) []QuestItemCount {
	return parseQuestItemCountsWithField(n, "count")
}

func parseQuestItemCountsWithField(n *node, countField string) []QuestItemCount {
	if n == nil {
		return nil
	}
	out := make([]QuestItemCount, 0, len(n.Children))
	for _, child := range n.Children {
		itemID := uint32(nodeInt(&child, "id", 0))
		count := nodeInt(&child, countField, 0)
		if itemID == 0 && count == 0 {
			continue
		}
		out = append(out, QuestItemCount{
			ItemID: itemID,
			Count:  count,
		})
	}
	return out
}

func parseQuestMobCounts(n *node) []QuestMobCount {
	if n == nil {
		return nil
	}
	out := make([]QuestMobCount, 0, len(n.Children))
	for _, child := range n.Children {
		mobID := uint32(nodeInt(&child, "id", 0))
		count := nodeInt(&child, "count", 0)
		if mobID == 0 && count == 0 {
			continue
		}
		out = append(out, QuestMobCount{
			MobID: mobID,
			Count: count,
		})
	}
	return out
}

func parseQuestStateRefs(n *node) []QuestStateRef {
	if n == nil {
		return nil
	}
	out := make([]QuestStateRef, 0, len(n.Children))
	for _, child := range n.Children {
		questID := uint32(nodeInt(&child, "id", 0))
		state := nodeInt(&child, "state", 0)
		if questID == 0 {
			continue
		}
		out = append(out, QuestStateRef{
			QuestID: questID,
			State:   state,
		})
	}
	return out
}

func parseQuestSkillRefs(n *node) []QuestSkillRef {
	if n == nil {
		return nil
	}
	out := make([]QuestSkillRef, 0, len(n.Children))
	for _, child := range n.Children {
		skillID := uint32(nodeInt(&child, "id", 0))
		acquire := nodeInt(&child, "acquire", 0)
		if skillID == 0 {
			continue
		}
		out = append(out, QuestSkillRef{
			SkillID: skillID,
			Acquire: acquire,
		})
	}
	return out
}

func parseQuestRewardItems(n *node) []QuestRewardItem {
	if n == nil {
		return nil
	}
	out := make([]QuestRewardItem, 0, len(n.Children))
	for _, child := range n.Children {
		item := QuestRewardItem{
			ItemID:     uint32(nodeInt(&child, "id", 0)),
			Count:      nodeInt(&child, "count", 0),
			Job:        nodeInt(&child, "job", -1),
			JobEx:      nodeInt(&child, "jobEx", -1),
			Gender:     nodeInt(&child, "gender", 2),
			Period:     nodeInt(&child, "period", 0),
			Prop:       -2,
			DateExpire: nodeString(&child, "dateExpire"),
		}
		if child.find("prop") != nil {
			item.Prop = nodeInt(&child, "prop", -1)
		}
		if item.ItemID == 0 && item.Count == 0 {
			continue
		}
		out = append(out, item)
	}
	return out
}

func parseQuestRewardSkills(n *node) []QuestRewardSkill {
	if n == nil {
		return nil
	}
	out := make([]QuestRewardSkill, 0, len(n.Children))
	for _, child := range n.Children {
		skill := QuestRewardSkill{
			SkillID:     uint32(nodeInt(&child, "id", 0)),
			SkillLevel:  nodeInt(&child, "skillLevel", 0),
			MasterLevel: nodeInt(&child, "masterLevel", 0),
		}
		if jobs := child.find("job"); jobs != nil {
			skill.Jobs = collectChildIntValues(jobs)
		}
		if skill.SkillID == 0 {
			continue
		}
		out = append(out, skill)
	}
	return out
}

func collectQuestActionJobs(n *node) []int {
	if n == nil {
		return nil
	}
	var jobs []int
	index := 0
	for {
		entry := n.find(strconv.Itoa(index))
		if entry == nil {
			break
		}
		if jobNode := entry.find("job"); jobNode != nil {
			jobs = append(jobs, collectChildIntValues(jobNode)...)
		}
		index++
	}
	if len(jobs) > 0 {
		return jobs
	}
	if jobNode := n.find("job"); jobNode != nil {
		return collectChildIntValues(jobNode)
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
