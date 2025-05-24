package data

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"
)

type node struct {
	Name     string `xml:"name,attr"`
	Value    string `xml:"value,attr"`
	Children []node `xml:",any"`
}

type Resources struct {
	Maps     map[uint32]*MapSpec
	Monsters map[uint32]*MobSpec
	Items    map[uint32]ItemSpec
}

func (node *node) find(name string) *node {

	parts := strings.Split(name, ":")
	current := node
	for _, part := range parts {
		found := false
		for _, v := range current.Children {
			if v.Name == part {
				current = &v
				found = true
				break
			}
		}

		if !found {
			return nil
		}
	}

	return current
}

func loadResourceFiles[T any](root string, workerCount int, action func(path string) (result *T, err error), callback func(percent float32, value *T)) error {
	var allFiles []string
	err := filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err == nil && !d.IsDir() && strings.HasSuffix(d.Name(), ".img.xml") {
			allFiles = append(allFiles, path)
		}
		return nil
	})
	if err != nil {
		return err
	}

	total := len(allFiles)
	if total == 0 {
		return nil
	}

	jobs := make(chan string, total)
	results := make(chan *T, total)
	var wg sync.WaitGroup

	for range workerCount {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for path := range jobs {
				m, err := action(path)

				if err == nil {
					results <- m
				} else {
					log.Println(err)
				}
			}
		}()
	}

	for _, path := range allFiles {
		jobs <- path
	}
	close(jobs)

	go func() {
		wg.Wait()
		close(results)
	}()

	count := 0
	for m := range results {
		count++
		percent := float32(count) / float32(total) * 100
		callback(percent, m)
	}

	return nil
}

func NewResources() *Resources {

	workerCount := runtime.NumCPU() * 2

	items := map[uint32]ItemSpec{}
	err := loadResourceFiles("D:/git/fm/wz/Character.wz",
		workerCount,
		func(path string) (result *EquipmentSpec, err error) {
			return loadWeapons(path)
		},
		func(percent float32, value *EquipmentSpec) {
			items[value.Id] = value
			fmt.Printf("장비 데이터 로딩 중: %.1f%%\n", percent)
		})
	if err != nil {
		log.Fatal(err)
		return nil
	}

	err = loadResourceFiles("D:/git/fm/wz/Item.wz/Consume",
		workerCount,
		func(path string) (result *[]*ConsumeSpec, err error) {
			return loadConsumes(path)
		},
		func(percent float32, value *[]*ConsumeSpec) {

			for _, v := range *value {
				items[v.Id] = v
			}
			fmt.Printf("소비 아이템 데이터 로딩 중: %.1f%%\n", percent)
		})
	if err != nil {
		log.Fatal(err)
		return nil
	}

	err = loadResourceFiles("D:/git/fm/wz/Item.wz/Cash",
		workerCount,
		func(path string) (result *[]*CashItemSpec, err error) {
			return loadCashItems(path)
		},
		func(percent float32, value *[]*CashItemSpec) {

			for _, v := range *value {
				items[v.Id] = v
			}
			fmt.Printf("캐시 아이템 데이터 로딩 중: %.1f%%\n", percent)
		})
	if err != nil {
		log.Fatal(err)
		return nil
	}

	err = loadResourceFiles("D:/git/fm/wz/Item.wz/Install",
		workerCount,
		func(path string) (result *[]*InstallationSpec, err error) {
			return loadInstallations(path)
		},
		func(percent float32, value *[]*InstallationSpec) {

			for _, v := range *value {
				items[v.Id] = v
			}
			fmt.Printf("설치 아이템 데이터 로딩 중: %.1f%%\n", percent)
		})
	if err != nil {
		log.Fatal(err)
		return nil
	}

	err = loadResourceFiles("D:/git/fm/wz/Item.wz/Special",
		workerCount,
		func(path string) (result *[]*SpecialItemSpec, err error) {
			return loadSpecialItems(path)
		},
		func(percent float32, value *[]*SpecialItemSpec) {

			for _, v := range *value {
				items[v.Id] = v
			}
			fmt.Printf("설치 아이템 데이터 로딩 중: %.1f%%\n", percent)
		})
	if err != nil {
		log.Fatal(err)
		return nil
	}

	err = loadResourceFiles("D:/git/fm/wz/Item.wz/Etc",
		workerCount,
		func(path string) (result *[]*GeneralItemSpec, err error) {
			return loadGeneralItems(path)
		},
		func(percent float32, value *[]*GeneralItemSpec) {

			for _, v := range *value {
				items[v.Id] = v
			}
			fmt.Printf("일반 아이템 데이터 로딩 중: %.1f%%\n", percent)
		})
	if err != nil {
		log.Fatal(err)
		return nil
	}

	err = loadResourceFiles("D:/git/fm/wz/Item.wz/Pet",
		workerCount,
		func(path string) (result *PetSpec, err error) {
			return loadPets(path)
		},
		func(percent float32, value *PetSpec) {

			items[value.Id] = value
			fmt.Printf("펫 데이터 로딩 중: %.1f%%\n", percent)
		})
	if err != nil {
		log.Fatal(err)
		return nil
	}

	stringResult := map[uint32]*StringSpec{}
	_ = loadResourceFiles("D:/git/fm/wz/String.wz", workerCount, func(path string) (result *[]*StringSpec, err error) {

		m, err := loadStringResources(path)
		if err != nil {
			return nil, err
		}
		return m, nil
	}, func(percent float32, value *[]*StringSpec) {

		for _, v := range *value {
			stringResult[v.Id] = v
		}

		fmt.Printf("문자열 데이터 로딩 중: %.1f%%\n", percent)
	})

	maps := map[uint32]*MapSpec{}
	err = loadResourceFiles("D:/git/fm/wz/Map.wz/Map", workerCount, func(path string) (result *MapSpec, err error) {
		base := filepath.Base(path)
		idStr := strings.TrimSuffix(base, ".img.xml")
		mapId, err := strconv.Atoi(idStr)
		if err != nil {
			return nil, err
		}

		m, err := loadMaps(path, uint32(mapId))
		if err != nil {
			return nil, err
		}
		return m, nil
	}, func(percent float32, value *MapSpec) {
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
		Monsters: map[uint32]*MobSpec{},
		Items:    items,
	}
}
