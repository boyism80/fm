package wz

import (
	"encoding/xml"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/boyism80/fm/services/game/constant"
	"github.com/boyism80/fm/types"
)

func loadReactor(path string) (*Reactor, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var root node
	if err := xml.NewDecoder(file).Decode(&root); err != nil {
		return nil, err
	}
	if filepath.Base(path) != root.Name+".xml" {
		return nil, fmt.Errorf("reactor file %s does not match %s", path, root.Name)
	}

	id, err := strconv.Atoi(strings.TrimSuffix(root.Name, ".img"))
	if err != nil {
		return nil, err
	}

	model := &Reactor{
		ID:     uint32(id),
		States: make(map[byte]*ReactorEvent),
	}

	info := root.find("info")
	if info != nil {
		for _, intField := range info.Ints {
			switch intField.Name {
			case "link":
				model.Info.Link = uint32(intField.Value)
			case "activateByTouch":
				model.Info.ActivateByTouch = intField.Value
			}
		}
		for _, strField := range info.Strings {
			switch strField.Name {
			case "link":
				if id, err := strconv.Atoi(strField.Value); err == nil {
					model.Info.Link = uint32(id)
				}
			}
		}
	}

	for _, strField := range root.Strings {
		if strField.Name == "action" {
			model.Action = strField.Value
			break
		}
	}
	if model.Action == "" {
		if actionNode := root.find("action"); actionNode != nil {
			model.Action = actionNode.Value
		}
	}

	for stateIndex := byte(0); ; stateIndex++ {
		stateNode := root.find(strconv.Itoa(int(stateIndex)))
		if stateNode == nil {
			break
		}

		var event *ReactorEvent
		eventNode := stateNode.find("event")
		if eventNode != nil {
			event0 := eventNode.find("0")
			if event0 != nil {
				event = parseReactorEvent(eventNode, event0)
				if event != nil && int(event.Type) >= 999 {
					event = nil
				}
			}
		}
		model.States[stateIndex] = event
	}

	return model, nil
}

func parseReactorEvent(eventNode *node, event0 *node) *ReactorEvent {
	event := &ReactorEvent{
		Type:      constant.ReactorEventType(nodeInt(event0, "type", 0)),
		NextState: byte(nodeInt(event0, "state", 0)),
		TimeOut:   nodeInt(eventNode, "timeOut", -1),
		TouchFlag: nodeInt(event0, "2", 0),
	}

	for _, child := range eventNode.Children {
		if nodeInt(&child, "type", 0) != int(constant.ReactorEventTypeItem) {
			continue
		}
		itemID, ok := nodeIntOptional(&child, "0")
		for _, f := range child.Strings {
			if f.Name == "0" {
				value, err := strconv.Atoi(f.Value)
				itemID, ok = value, err == nil
			}
		}
		if ok == false {
			continue
		}
		quantity, _ := nodeIntOptional(&child, "1")
		if quantity <= 0 {
			quantity = 1
		}
		event.Items = append(event.Items, ReactorItem{ID: itemID, Quantity: quantity})
	}

	if lt, ok := parsePointFromNode(event0, "lt"); ok {
		event.LT = lt
		event.HasLT = true
	}
	if rb, ok := parsePointFromNode(event0, "rb"); ok {
		event.RB = rb
		event.HasRB = true
	}

	if event0.find("clickArea") != nil {
		event.HasClickArea = true
	}

	return event
}

func parsePointFromNode(parent *node, name string) (types.Vector2[int32], bool) {
	if parent == nil {
		return types.Vector2[int32]{}, false
	}

	for _, vecField := range parent.Vectors {
		if vecField.Name == name {
			return types.Vector2[int32]{
				X: int32(vecField.X),
				Y: int32(vecField.Y),
			}, true
		}
	}

	child := parent.find(name)
	if child == nil {
		return types.Vector2[int32]{}, false
	}

	for _, vecField := range child.Vectors {
		return types.Vector2[int32]{
			X: int32(vecField.X),
			Y: int32(vecField.Y),
		}, true
	}

	x, xOk := nodeIntOptional(child, "x")
	y, yOk := nodeIntOptional(child, "y")
	if xOk && yOk {
		return types.Vector2[int32]{
			X: int32(x),
			Y: int32(y),
		}, true
	}

	return types.Vector2[int32]{}, false
}

func nodeIntOptional(n *node, name string) (int, bool) {
	if n == nil {
		return 0, false
	}
	for _, f := range n.Ints {
		if f.Name == name {
			return f.Value, true
		}
	}
	if child := n.find(name); child != nil {
		if val, err := strconv.Atoi(child.Value); err == nil {
			return val, true
		}
	}
	return 0, false
}
