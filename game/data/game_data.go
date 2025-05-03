package data

import (
	"log"
	"runtime"
)

type GameData struct {
	Maps     map[uint32]*MapTemplate
	Monsters map[uint32]*MonsterTemplate
	Items    map[uint32]*ItemTemplate
}

// 생성자
func NewGameData() *GameData {

	workerCount := runtime.NumCPU() * 2
	maps, err := LoadAllMapsParallel("D:/git/fm/wz/Map.wz/Map", workerCount)
	if err != nil {
		log.Fatal(err)
		return nil
	}

	return &GameData{
		Maps:     maps,
		Monsters: map[uint32]*MonsterTemplate{},
		Items:    map[uint32]*ItemTemplate{},
	}
}
