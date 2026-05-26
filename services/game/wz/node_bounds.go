package wz

import "github.com/boyism80/fm/types"

func GetRectFromWzNode(parent *node) types.Rect[int32] {
	if parent == nil {
		return types.Rect[int32]{}
	}

	var lt types.Point[int32]
	ltNode := parent.find("lt")
	if ltNode != nil {
		lt.X = int32(nodeInt(ltNode, "x", 0))
		lt.Y = int32(nodeInt(ltNode, "y", 0))
		for _, v := range ltNode.Vectors {
			if v.Name == "lt" || v.Name == "" {
				lt.X = int32(v.X)
				lt.Y = int32(v.Y)
			}
		}
	}

	var rb types.Point[int32]
	rbNode := parent.find("rb")
	if rbNode != nil {
		rb.X = int32(nodeInt(rbNode, "x", 0))
		rb.Y = int32(nodeInt(rbNode, "y", 0))
		for _, v := range rbNode.Vectors {
			if v.Name == "rb" || v.Name == "" {
				rb.X = int32(v.X)
				rb.Y = int32(v.Y)
			}
		}
	}

	return types.NewRect(lt, rb)
}
