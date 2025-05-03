package data

import (
	"encoding/xml"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
)

type xmlNode struct {
	XMLName  xml.Name
	Name     string    `xml:"name,attr"`
	Value    string    `xml:"value,attr"`
	Children []xmlNode `xml:",any"`
}

func loadMapFromXML(path string, mapId uint32) (*MapTemplate, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var root xmlNode
	if err := xml.NewDecoder(file).Decode(&root); err != nil {
		return nil, err
	}

	var m MapTemplate
	m.Id = mapId
	m.Portals = make(map[uint8]Portal)

	for _, child := range root.Children {
		if child.Name == "info" {
			for _, info := range child.Children {
				switch info.Name {
				case "mapName":
					m.Name = info.Value
				case "version":
					m.Version, _ = strconv.Atoi(info.Value)
				case "cloud":
					m.Cloud, _ = strconv.Atoi(info.Value)
				case "returnMap":
					m.ReturnMapId, _ = strconv.Atoi(info.Value)
				case "forcedReturn":
					m.ForcedReturn, _ = strconv.Atoi(info.Value)
				case "fieldLimit":
					m.FieldLimit, _ = strconv.Atoi(info.Value)
				case "VRTop":
					m.VRTop, _ = strconv.Atoi(info.Value)
				case "VRLeft":
					m.VRLeft, _ = strconv.Atoi(info.Value)
				case "VRBottom":
					m.VRBottom, _ = strconv.Atoi(info.Value)
				case "VRRight":
					m.VRRight, _ = strconv.Atoi(info.Value)
				case "hideMinimap":
					m.HideMinimap = info.Value == "1"
				case "town":
					m.IsTown = info.Value == "1"
				case "mobRate":
					f, err := strconv.ParseFloat(info.Value, 32)
					if err == nil {
						m.MobRate = float32(f)
					}
				case "bgm":
					m.BGM = info.Value
				case "mapMark":
					m.MapMark = info.Value
				case "mapDesc":
					m.MapDesc = info.Value
				case "miniMapOnOff":
					m.MiniMapOnOff = info.Value == "1"
				default:
					break
				}
			}
		}

		if child.Name == "portal" {
			for _, pnode := range child.Children {
				var portal Portal
				for _, field := range pnode.Children {
					switch field.Name {
					case "pn":
						portal.Name = field.Value
					case "pt":
						v, _ := strconv.Atoi(field.Value)
						portal.Type = uint8(v)
					case "tm":
						v, _ := strconv.Atoi(field.Value)
						portal.TargetMapId = int32(v)
					case "tn":
						portal.Target = field.Value
					case "x":
						v, _ := strconv.Atoi(field.Value)
						portal.Position.X = int16(v)
					case "y":
						v, _ := strconv.Atoi(field.Value)
						portal.Position.Y = int16(v)
					case "script":
						if field.Value != "" {
							portal.ScriptName = field.Value
						}
					}
				}
				id, _ := strconv.Atoi(pnode.Name)
				portal.Id = uint8(id)
				m.Portals[portal.Id] = portal
			}
		}
	}

	return &m, nil
}

func LoadAllMapsParallel(root string, workerCount int) (map[uint32]*MapTemplate, error) {
	// Step 1: 모든 .img.xml 파일 수집
	var allFiles []string
	err := filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err == nil && !d.IsDir() && strings.HasSuffix(d.Name(), ".img.xml") {
			allFiles = append(allFiles, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	total := len(allFiles)
	if total == 0 {
		return nil, nil
	}

	// Step 2: 병렬 처리
	jobs := make(chan string, total)
	results := make(chan *MapTemplate, total)
	var wg sync.WaitGroup

	// Step 3: 워커 goroutine 생성
	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for path := range jobs {
				base := filepath.Base(path)
				idStr := strings.TrimSuffix(base, ".img.xml")
				mapId, err := strconv.Atoi(idStr)
				if err != nil {
					continue
				}

				m, err := loadMapFromXML(path, uint32(mapId))
				if err == nil {
					results <- m
				}
			}
		}()
	}

	// Step 4: 작업 분배
	for _, path := range allFiles {
		jobs <- path
	}
	close(jobs)

	// Step 5: 결과 수집 goroutine
	go func() {
		wg.Wait()
		close(results)
	}()

	maps := map[uint32]*MapTemplate{}
	for m := range results {
		maps[m.Id] = m
		progress := len(maps)
		percent := float64(progress) / float64(total) * 100
		fmt.Printf("\r로딩 중: %d / %d (%.1f%%)\n", progress, total, percent)
	}

	fmt.Println("\n모든 맵 로딩 완료.")
	return maps, nil

}
