package data

import (
	"encoding/xml"
	"os"
	"strconv"
)

func loadMapFromXML(path string, mapId uint32) (*MapSpec, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var root node
	if err := xml.NewDecoder(file).Decode(&root); err != nil {
		return nil, err
	}

	spec := MapSpec{
		Id:      mapId,
		Portals: map[uint8]Portal{},
	}

	for _, child := range root.Children {
		if child.Name == "info" {
			for _, info := range child.Children {
				switch info.Name {
				case "mapName":
					spec.Name = info.Value
				case "version":
					spec.Version, _ = strconv.Atoi(info.Value)
				case "cloud":
					spec.Cloud, _ = strconv.Atoi(info.Value)
				case "returnMap":
					spec.ReturnMapId, _ = strconv.Atoi(info.Value)
				case "forcedReturn":
					spec.ForcedReturn, _ = strconv.Atoi(info.Value)
				case "fieldLimit":
					spec.FieldLimit, _ = strconv.Atoi(info.Value)
				case "VRTop":
					spec.VRTop, _ = strconv.Atoi(info.Value)
				case "VRLeft":
					spec.VRLeft, _ = strconv.Atoi(info.Value)
				case "VRBottom":
					spec.VRBottom, _ = strconv.Atoi(info.Value)
				case "VRRight":
					spec.VRRight, _ = strconv.Atoi(info.Value)
				case "hideMinimap":
					spec.HideMinimap = info.Value == "1"
				case "town":
					spec.IsTown = info.Value == "1"
				case "mobRate":
					f, err := strconv.ParseFloat(info.Value, 32)
					if err == nil {
						spec.MobRate = float32(f)
					}
				case "bgm":
					spec.BGM = info.Value
				case "mapMark":
					spec.MapMark = info.Value
				case "mapDesc":
					spec.MapDesc = info.Value
				case "miniMapOnOff":
					spec.MiniMapOnOff = info.Value == "1"
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
				spec.Portals[portal.Id] = portal
			}
		}
	}

	return &spec, nil
}
