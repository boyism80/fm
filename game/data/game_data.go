package data

import (
	"fmt"
	"log"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
)

type Resources struct {
	Maps     map[uint32]*MapTemplate
	Monsters map[uint32]*MonsterTemplate
	Items    map[uint32]ItemTemplate
}

func NewGameData() *Resources {

	workerCount := runtime.NumCPU() * 2

	items := map[uint32]ItemTemplate{}
	err := LoadXmlFiles("D:/git/fm-backup/wz/Character.wz",
		workerCount,
		func(path string) (result *EquipmentTemplate, err error) {
			return loadWeaponFromXML(path)
		},
		func(percent float32, value *EquipmentTemplate) {
			items[value.Id] = value
			fmt.Printf("장비 데이터 로딩 중: %.1f%%\n", percent)
		})
	if err != nil {
		log.Fatal(err)
		return nil
	}

	err = LoadXmlFiles("D:/git/fm-backup/wz/Item.wz/Consume",
		workerCount,
		func(path string) (result *[]*ConsumeTemplate, err error) {
			return loadConsumeFromXML(path)
		},
		func(percent float32, value *[]*ConsumeTemplate) {

			for _, v := range *value {
				items[v.Id] = v
			}
			fmt.Printf("소비 아이템 데이터 로딩 중: %.1f%%\n", percent)
		})
	if err != nil {
		log.Fatal(err)
		return nil
	}

	err = LoadXmlFiles("D:/git/fm-backup/wz/Item.wz/Cash",
		workerCount,
		func(path string) (result *[]*CashItemTemplate, err error) {
			return loadCashItemFromXML(path)
		},
		func(percent float32, value *[]*CashItemTemplate) {

			for _, v := range *value {
				items[v.Id] = v
			}
			fmt.Printf("캐시 아이템 데이터 로딩 중: %.1f%%\n", percent)
		})
	if err != nil {
		log.Fatal(err)
		return nil
	}

	err = LoadXmlFiles("D:/git/fm-backup/wz/Item.wz/Install",
		workerCount,
		func(path string) (result *[]*InstallationTemplate, err error) {
			return loadInstallationFromXML(path)
		},
		func(percent float32, value *[]*InstallationTemplate) {

			for _, v := range *value {
				items[v.Id] = v
			}
			fmt.Printf("설치 아이템 데이터 로딩 중: %.1f%%\n", percent)
		})
	if err != nil {
		log.Fatal(err)
		return nil
	}

	err = LoadXmlFiles("D:/git/fm-backup/wz/Item.wz/Special",
		workerCount,
		func(path string) (result *[]*SpecialItemTemplate, err error) {
			return loadSpecialItemFromXML(path)
		},
		func(percent float32, value *[]*SpecialItemTemplate) {

			for _, v := range *value {
				items[v.Id] = v
			}
			fmt.Printf("설치 아이템 데이터 로딩 중: %.1f%%\n", percent)
		})
	if err != nil {
		log.Fatal(err)
		return nil
	}

	err = LoadXmlFiles("D:/git/fm-backup/wz/Item.wz/Etc",
		workerCount,
		func(path string) (result *[]*GeneralItemTemplate, err error) {
			return loadGeneralItemFromXML(path)
		},
		func(percent float32, value *[]*GeneralItemTemplate) {

			for _, v := range *value {
				items[v.Id] = v
			}
			fmt.Printf("일반 아이템 데이터 로딩 중: %.1f%%\n", percent)
		})
	if err != nil {
		log.Fatal(err)
		return nil
	}

	err = LoadXmlFiles("D:/git/fm-backup/wz/Item.wz/Pet",
		workerCount,
		func(path string) (result *PetTemplate, err error) {
			return loadPetFromXML(path)
		},
		func(percent float32, value *PetTemplate) {

			items[value.Id] = value
			fmt.Printf("펫 데이터 로딩 중: %.1f%%\n", percent)
		})
	if err != nil {
		log.Fatal(err)
		return nil
	}

	stringResult := map[uint32]*StringTemplate{}
	_ = LoadXmlFiles("D:/git/fm/wz/String.wz", workerCount, func(path string) (result *[]*StringTemplate, err error) {

		m, err := loadStringFromXML(path)
		if err != nil {
			return nil, err
		}
		return m, nil
	}, func(percent float32, value *[]*StringTemplate) {

		for _, v := range *value {
			stringResult[v.Id] = v
		}

		fmt.Printf("문자열 데이터 로딩 중: %.1f%%\n", percent)
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
		fmt.Printf("맵 데이터 로딩 중: %.1f%%\n", percent)
	})
	if err != nil {
		log.Fatal(err)
		return nil
	}
	fmt.Println("\n모든 맵 로딩 완료.")

	return &Resources{
		Maps:     maps,
		Monsters: map[uint32]*MonsterTemplate{},
		Items:    items,
	}
}
