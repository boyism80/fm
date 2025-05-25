package entity

import "github.com/boyism80/fm/common/types"

type Map struct {
	Id          uint32
	Name        string
	SpawnPoints map[uint8]types.Vector2[int32]
	Bounds      types.Rect[int32]
	IsTown      bool
	HasClock    bool
	Properties  map[string]string
}
