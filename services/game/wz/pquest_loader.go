package wz

import (
	"encoding/xml"
	"os"
	"path/filepath"
	"strconv"
)

func attachPartyQuestRules(quests map[uint32]*Quest, path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	var root node
	if err := xml.NewDecoder(file).Decode(&root); err != nil {
		return err
	}

	for _, questNode := range root.Children {
		questID, err := strconv.ParseUint(questNode.Name, 10, 32)
		if err != nil {
			continue
		}
		quest := quests[uint32(questID)]
		if quest == nil {
			continue
		}
		rankNode := questNode.find("rank")
		if rankNode == nil {
			continue
		}
		ranks := make(map[string][]PartyQuestRankCheck, len(rankNode.Children))
		for _, rankChild := range rankNode.Children {
			checks := make([]PartyQuestRankCheck, 0)
			for _, modeNode := range rankChild.Children {
				mode := PartyQuestRankMode(modeNode.Name)
				for _, intField := range modeNode.Ints {
					checks = append(checks, PartyQuestRankCheck{
						Mode:     mode,
						Property: intField.Name,
						Value:    intField.Value,
					})
				}
			}
			if len(checks) > 0 {
				ranks[rankChild.Name] = checks
			}
		}
		if len(ranks) > 0 {
			quest.PartyRanks = ranks
		}
	}
	return nil
}

func loadPartyQuestRules(quests map[uint32]*Quest, wzPath string) {
	path := filepath.Join(wzPath, "Quest.wz", "PQuest.img.xml")
	if err := attachPartyQuestRules(quests, path); err != nil {
		return
	}
}
