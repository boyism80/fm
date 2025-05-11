package data

import "github.com/boyism80/fm/common/types"

type Portal struct {
	Id          uint8
	Name        string
	TargetMapId int32
	Target      string
	Position    types.Vec2[int16]
	ScriptName  string
	Type        uint8
}

type MapSpec struct {
	Id           uint32
	Name         string
	Version      int
	Cloud        int
	ReturnMapId  int
	ForcedReturn int
	FieldLimit   int
	VRTop        int
	VRLeft       int
	VRBottom     int
	VRRight      int
	HideMinimap  bool
	IsTown       bool
	MobRate      float32
	BGM          string
	MapMark      string
	MapDesc      string
	MiniMapOnOff bool
	Portals      map[uint8]Portal
}
