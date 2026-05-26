package wz

import (
	"encoding/xml"
	"os"
	"strconv"

	"github.com/boyism80/fm/types"
)

type MobSkillSlot struct {
	Slot    int
	SkillID uint32
	Level   uint8
	Action  int
}

type MobAttack struct {
	Index        int
	DeadlyAttack bool
	MpBurn       int
	MpCon        int
	DiseaseSkill uint32
	DiseaseLevel uint8
	AttackAfter  int
	PADamage     int
	MADamage     int
	Magic        bool
	RangeR       int
}

type MobSkillLevelData struct {
	SkillWz     *Skill
	SkillID     uint32
	Level       uint8
	HpPercent   int
	MpCon       int
	X           int
	Y           int
	DurationMs  int64
	CooltimeMs  int64
	Prop        float32
	Limit       int16
	SpawnEffect int
	SummonOnce  bool
	Summons     []uint32
	Bounds      types.Rect[int32]
}

func (m *Mob) HasSkill(skillID uint32, level uint8) bool {
	if m == nil {
		return false
	}
	for _, slot := range m.Skills {
		if slot.SkillID == skillID && slot.Level == level {
			return true
		}
	}
	return false
}

func nodeInt(n *node, name string, defaultVal int) int {
	if n == nil {
		return defaultVal
	}
	for _, f := range n.Ints {
		if f.Name == name {
			return f.Value
		}
	}
	return defaultVal
}

func parseMobInfoSkills(skillRoot node) []MobSkillSlot {
	var out []MobSkillSlot
	for _, slotDir := range skillRoot.Children {
		slot, err := strconv.Atoi(slotDir.Name)
		if err != nil {
			continue
		}
		skillID := nodeInt(&slotDir, "skill", 0)
		level := nodeInt(&slotDir, "level", 0)
		if skillID <= 0 || level <= 0 {
			continue
		}
		entry := MobSkillSlot{
			Slot:    slot,
			SkillID: uint32(skillID),
			Level:   uint8(level),
			Action:  -1,
		}
		if action := nodeInt(&slotDir, "action", -1); action >= 0 {
			entry.Action = action
		}
		out = append(out, entry)
	}
	return out
}

func parseMobAttackInfo(attackIndex int, info *node) *MobAttack {
	if info == nil {
		return nil
	}
	attack := &MobAttack{
		Index:        attackIndex,
		DiseaseSkill: uint32(nodeInt(info, "disease", 0)),
		DiseaseLevel: uint8(nodeInt(info, "level", 0)),
		MpCon:        nodeInt(info, "conMP", 0),
		MpBurn:       nodeInt(info, "mpBurn", 0),
		AttackAfter:  nodeInt(info, "attackAfter", 0),
		PADamage:     nodeInt(info, "PADamage", 0),
		MADamage:     nodeInt(info, "PADamage", 0),
		Magic:        nodeInt(info, "magic", 0) > 0,
	}
	for _, ch := range info.Children {
		if ch.Name == "deadlyAttack" {
			attack.DeadlyAttack = true
			break
		}
	}
	rangeNode := info.find("range")
	if rangeNode != nil {
		attack.RangeR = nodeInt(rangeNode, "r", 0)
	}
	return attack
}

func parseMobAttacks(root node) []MobAttack {
	var out []MobAttack
	for _, ch := range root.Children {
		if len(ch.Name) < 7 || ch.Name[:6] != "attack" {
			continue
		}
		suffix := ch.Name[6:]
		index, err := strconv.Atoi(suffix)
		if err != nil || index < 1 {
			continue
		}
		info := ch.find("info")
		attack := parseMobAttackInfo(index, info)
		if attack == nil {
			continue
		}
		out = append(out, *attack)
	}
	return out
}

func parseMobSkillSummons(levelNode *node) []uint32 {
	if levelNode == nil {
		return nil
	}
	var summons []uint32
	for i := 0; ; i++ {
		key := strconv.Itoa(i)
		found := false
		for _, intf := range levelNode.Ints {
			if intf.Name == key {
				summons = append(summons, uint32(intf.Value))
				found = true
				break
			}
		}
		if found {
			continue
		}
		child := levelNode.find(key)
		if child == nil {
			break
		}
		val := nodeInt(child, "id", 0)
		if val == 0 {
			for _, intf := range child.Ints {
				if intf.Name == "id" || intf.Name == "value" {
					val = intf.Value
					break
				}
			}
			if val == 0 && len(child.Ints) > 0 {
				val = child.Ints[0].Value
			}
			if val == 0 && child.Value != "" {
				if parsed, err := strconv.Atoi(child.Value); err == nil {
					val = parsed
				}
			}
		}
		if val <= 0 {
			break
		}
		summons = append(summons, uint32(val))
	}
	return summons
}

func parseMobSkillLevel(skillID uint32, level uint8, levelNode *node) *MobSkillLevelData {
	if levelNode == nil {
		return nil
	}
	data := &MobSkillLevelData{
		SkillID:     skillID,
		Level:       level,
		HpPercent:   nodeInt(levelNode, "hp", 100),
		MpCon:       nodeInt(levelNode, "mpCon", 0),
		X:           nodeInt(levelNode, "x", 1),
		Y:           nodeInt(levelNode, "y", 1),
		DurationMs:  int64(nodeInt(levelNode, "time", 0)) * 1000,
		CooltimeMs:  int64(nodeInt(levelNode, "interval", 0)) * 1000,
		Prop:        float32(nodeInt(levelNode, "prop", 100)) / 100,
		Limit:       int16(nodeInt(levelNode, "limit", 0)),
		SpawnEffect: nodeInt(levelNode, "summonEffect", 0),
		SummonOnce:  nodeInt(levelNode, "summonOnce", 0) > 0,
		Summons:     parseMobSkillSummons(levelNode),
		Bounds:      GetRectFromWzNode(levelNode),
	}
	return data
}

func loadMobSkillData(path string) (map[uint32]map[uint8]*MobSkillLevelData, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var root node
	if err := xml.NewDecoder(file).Decode(&root); err != nil {
		return nil, err
	}

	out := make(map[uint32]map[uint8]*MobSkillLevelData)
	for _, skillDir := range root.Children {
		skillID64, err := strconv.ParseUint(skillDir.Name, 10, 32)
		if err != nil {
			continue
		}
		skillID := uint32(skillID64)
		skillWz := &Skill{ID: skillID}
		levelRoot := skillDir.find("level")
		if levelRoot == nil {
			continue
		}
		for _, levelDir := range levelRoot.Children {
			level64, err := strconv.ParseUint(levelDir.Name, 10, 8)
			if err != nil || level64 == 0 {
				continue
			}
			level := uint8(level64)
			data := parseMobSkillLevel(skillID, level, &levelDir)
			if data == nil {
				continue
			}
			data.SkillWz = skillWz
			levels, ok := out[skillID]
			if !ok {
				levels = make(map[uint8]*MobSkillLevelData)
				out[skillID] = levels
			}
			levels[level] = data
		}
	}
	return out, nil
}
