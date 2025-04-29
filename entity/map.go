package entity

import (
	"github.com/boyism80/fm/ctrl"
	"github.com/boyism80/fm/types"
)

type Map struct {
	MapObjects map[uint32]ctrl.ObjectController // 이제 타입별 분리 없이 한 곳에 다 모은다
}

func (m *Map) GetMapObjectsInRange(from types.Vec2, rangeSq int, targetType types.ObjectType) []ctrl.ObjectController {
	var ret []ctrl.ObjectController

	for _, obj := range m.MapObjects {
		if obj.Type()&targetType == 0 {
			continue
		}
		if from.DistanceSq(obj.Position()) <= rangeSq {
			ret = append(ret, obj)
		}
	}

	return ret
}

func NewMap() Map {
	return Map{
		MapObjects: map[uint32]ctrl.ObjectController{},
	}
}
