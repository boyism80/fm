package wz

import (
	"encoding/xml"
	"fmt"
	"os"
	"path/filepath"
)

// movePlayerNpcSeats moves each map npc whose Npc.wz info has imitate, an empty seat a player npc of that template
// takes, from NpcSpawns to PlayerNpcSeats, as the original never showed them.
func movePlayerNpcSeats(maps map[uint32]*Map, npcPath string) {
	imitates := map[uint32]bool{}
	for _, m := range maps {
		for spawnID, spawn := range m.NpcSpawns {
			imitate, ok := imitates[spawn.ID]
			if !ok {
				imitate = npcImitates(filepath.Join(npcPath, fmt.Sprintf("%07d.img.xml", spawn.ID)))
				imitates[spawn.ID] = imitate
			}
			if imitate {
				m.PlayerNpcSeats[spawn.ID] = spawn
				delete(m.NpcSpawns, spawnID)
			}
		}
	}
}

func npcImitates(path string) bool {
	file, err := os.Open(path)
	if err != nil {
		return false
	}
	defer file.Close()

	var root node
	if err := xml.NewDecoder(file).Decode(&root); err != nil {
		return false
	}
	return root.find("info").Int("imitate", 0) != 0
}
