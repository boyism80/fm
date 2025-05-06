package data

import (
	"fmt"
	"log"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
)

type GameData struct {
	Maps     map[uint32]*MapTemplate
	Monsters map[uint32]*MonsterTemplate
	Items    map[uint32]*ItemTemplate
}

// 생성자
func NewGameData() *GameData {

	workerCount := runtime.NumCPU() * 2

	items := map[uint32]*ItemTemplate{}
	err := LoadXmlFiles("D:/git/fm-backup/wz/Character.wz/Weapon",
		workerCount,
		func(path string) (result *ItemTemplate, err error) {
			return loadWeaponFromXML(path)
		},
		func(percent float32, value *ItemTemplate) {
			items[value.Id] = value
			fmt.Printf("\r무기 데이터 로딩 중: (%.1f%%)\n", percent)
		})
	if err != nil {
		log.Fatal(err)
		return nil
	}

	stringResult := map[uint32]*StringTemplate{}
	_ = LoadXmlFiles("D:/git/fm/wz/String.wz", workerCount, func(path string) (result *StringTemplate, err error) {

		// TODO: 파일명으로 분기해서 데이터 파싱 다르게 하기
		m, err := loadStringFromXML(path)
		if err != nil {
			return nil, err
		}
		return m, nil
	}, func(percent float32, value *StringTemplate) {

		// TODO: 파일명으로 분기해서 데이터 적재 다르게 하기
		stringResult[value.Id] = value
		fmt.Printf("\r문자열 데이터 로딩 중: (%.1f%%)\n", percent)
	})

	maps := map[uint32]*MapTemplate{}
	err = LoadXmlFiles("D:/git/fm/wz/Map.wz/Map", workerCount, func(path string) (result *MapTemplate, err error) {
		base := filepath.Base(path)
		idStr := strings.TrimSuffix(base, ".img.xml")
		mapId, err := strconv.Atoi(idStr)
		if err != nil {
			return nil, err
		}

		m, err := loadMapFromXML(path, uint32(mapId))
		if err != nil {
			return nil, err
		}
		return m, nil
	}, func(percent float32, value *MapTemplate) {
		maps[value.Id] = value
		fmt.Printf("\r맵 데이터 로딩 중: (%.1f%%)\n", percent)
	})
	if err != nil {
		log.Fatal(err)
		return nil
	}
	fmt.Println("\n모든 맵 로딩 완료.")

	return &GameData{
		Maps:     maps,
		Monsters: map[uint32]*MonsterTemplate{},
		Items:    map[uint32]*ItemTemplate{},
	}
}
