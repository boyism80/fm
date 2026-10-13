package wz

import "github.com/boyism80/fm/services/game/constant"

type Class struct {
	ID          constant.ClassType
	Parent      *Class
	Family      constant.ClassFamily
	Branch      constant.ClassBranch
	Specialty   uint8
	Advancement uint8
}

func NewClass(id constant.ClassType) *Class {
	c := &Class{
		ID:        id,
		Branch:    constant.ClassBranch(id % 1000 / 100),
		Specialty: uint8(id / 10 % 10),
	}
	switch {
	case id >= constant.ClassLegend:
		c.Family = constant.ClassFamilyLegend
	case id >= constant.ClassNoblesse:
		c.Family = constant.ClassFamilyCygnus
	}
	switch {
	case id%1000 == 0:
	case id%100 == 0:
		c.Advancement = 1
	default:
		c.Advancement = uint8(id%10) + 2
	}

	switch {
	case id%10 != 0:
		c.Parent = NewClass(id - 1)
	case id%1000 == 0 || id < 1000 && id%100 == 0:
	case id%100 == 0 && c.Family == constant.ClassFamilyCygnus:
		c.Parent = NewClass(constant.ClassNoblesse)
	case id%100 == 0:
		c.Parent = NewClass(constant.ClassLegend)
	case id < 1000:
		c.Parent = NewClass(id / 100 * 100)
	default:
		c.Parent = NewClass((id/10 - 1) * 10)
	}
	return c
}

// A quest class code ending in 00 stands for its whole branch.
func (c *Class) InQuestClasses(codes []int) bool {
	if len(codes) == 0 {
		return true
	}
	class := int(c.ID)
	for _, code := range codes {
		if code == class {
			return true
		}
		if code%100 == 0 && class/100 == code/100 {
			return true
		}
	}
	return false
}

func (c *Class) Of(ancestor constant.ClassType) bool {
	for cur := c; cur != nil; cur = cur.Parent {
		if cur.ID == ancestor {
			return true
		}
	}
	return false
}
