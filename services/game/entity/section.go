package entity

import (
	"github.com/boyism80/fm/services/game/constant"
	"github.com/boyism80/fm/types"
)

type sectionPosition struct {
	x int32
	y int32
}

type section struct {
	position sectionPosition
	objects  map[Object]struct{}
}

type sectionContainer struct {
	sections map[sectionPosition]*section
}

func newSectionContainer() *sectionContainer {
	return &sectionContainer{sections: make(map[sectionPosition]*section)}
}

func (c *sectionContainer) add(obj Object) {
	if c == nil || obj == nil {
		return
	}
	current := c.section(obj.GetPosition())
	if obj.getSection() == current {
		return
	}

	c.remove(obj)
	current.objects[obj] = struct{}{}
	obj.setSection(current)
}

func (c *sectionContainer) remove(obj Object) {
	if c == nil || obj == nil {
		return
	}
	current := obj.getSection()
	if current == nil {
		return
	}

	delete(current.objects, obj)
	obj.setSection(nil)
	if len(current.objects) == 0 {
		delete(c.sections, current.position)
	}
}

func (c *sectionContainer) objectsNear(position types.Vector2[int16], filter constant.ObjectType) []Object {
	objects := c.objectsAround(position)
	result := make([]Object, 0, len(objects))
	for _, obj := range objects {
		if !obj.Is(filter) {
			continue
		}
		dx := int64(obj.GetPosition().X) - int64(position.X)
		dy := int64(obj.GetPosition().Y) - int64(position.Y)
		if dx*dx+dy*dy <= constant.MaxViewRangeSq {
			result = append(result, obj)
		}
	}
	return result
}

func (c *sectionContainer) objectsIn(bounds types.Rect[int32], filter constant.ObjectType) []Object {
	minX := floorSection(bounds.Left)
	maxX := floorSection(bounds.Right)
	minY := floorSection(bounds.Top)
	maxY := floorSection(bounds.Bottom)
	result := make([]Object, 0)
	for y := minY; y <= maxY; y++ {
		for x := minX; x <= maxX; x++ {
			current := c.sections[sectionPosition{x: x, y: y}]
			if current == nil {
				continue
			}
			for obj := range current.objects {
				if !obj.Is(filter) {
					continue
				}
				position := obj.GetPosition()
				if bounds.ContainsPoint(types.Point[int32]{X: int32(position.X), Y: int32(position.Y)}) {
					result = append(result, obj)
				}
			}
		}
	}
	return result
}

func (c *sectionContainer) objectsAround(position types.Vector2[int16]) []Object {
	center := sectionFor(position)
	result := make([]Object, 0)
	for y := center.y - 1; y <= center.y+1; y++ {
		for x := center.x - 1; x <= center.x+1; x++ {
			current := c.sections[sectionPosition{x: x, y: y}]
			if current == nil {
				continue
			}
			for obj := range current.objects {
				result = append(result, obj)
			}
		}
	}
	return result
}

func (c *sectionContainer) section(position types.Vector2[int16]) *section {
	key := sectionFor(position)
	current := c.sections[key]
	if current == nil {
		current = &section{
			position: key,
			objects:  make(map[Object]struct{}),
		}
		c.sections[key] = current
	}
	return current
}

func sectionFor(position types.Vector2[int16]) sectionPosition {
	return sectionPosition{
		x: floorSection(int32(position.X)),
		y: floorSection(int32(position.Y)),
	}
}

func floorSection(value int32) int32 {
	current := value / constant.SectionSize
	if value < 0 && value%constant.SectionSize != 0 {
		current--
	}
	return current
}
