package wz

import (
	"encoding/xml"
	"os"
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

	if itemID, ok := nodeIntOptional(event0, "0"); ok {
		event.ItemID = itemID
	}
	if itemQty, ok := nodeIntOptional(event0, "1"); ok {
		event.ItemQuantity = itemQty
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
