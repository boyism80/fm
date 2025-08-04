// Package data provides MapleStory game data specifications and types.
// This file contains map-related data structures and utilities.
package data

import (
	"math"

	"github.com/boyism80/fm/core/types"
)

// Portal represents a map transition point in MapleStory.
type Portal struct {
	ID          uint8              // Portal identifier
	Name        string             // Portal name for targeting
	TargetMapId int32              // Destination map ID
	Target      string             // Target portal name
	Position    types.Point[int16] // Portal position in map
	ScriptName  string             // Lua script for portal logic
	Type        uint8              // Portal type (normal, script, etc.)
}

// NpcSpawnSpec represents an NPC spawn point on a map.
type NpcSpawnSpec struct {
	*SpawnSpec
}

// MobSpawnSpec represents a monster spawn point on a map.
type MobSpawnSpec struct {
	*SpawnSpec
}

// MapSpec contains all data for a MapleStory map.
type MapSpec struct {
	ID           uint32                               // Map identifier
	Name         string                               // Map display name
	Version      int                                  // Map version
	Cloud        int                                  // Cloud effect type
	ReturnMapId  int                                  // Return map for teleport
	ForcedReturn int                                  // Forced return map
	FieldLimit   int                                  // Field restrictions
	VRTop        int                                  // View rectangle top
	VRLeft       int                                  // View rectangle left
	VRBottom     int                                  // View rectangle bottom
	VRRight      int                                  // View rectangle right
	HideMinimap  bool                                 // Hide minimap flag
	IsTown       bool                                 // Town map flag
	MobRate      float32                              // Monster spawn rate
	BGM          string                               // Background music
	MapMark      string                               // Map mark identifier
	MapDesc      string                               // Map description
	MiniMapOnOff bool                                 // Minimap toggle
	Portals      map[uint8]Portal                     // Map portals
	NpcSpawns    map[uint32]NpcSpawnSpec              // NPC spawn points
	MobSpawns    map[uint32]MobSpawnSpec              // Monster spawn points
	Footholds    *types.QuadTreeNode[int16, Foothold] // Foothold collision data
}

// FootholdPoint calculates the foothold position for a given point.
// Returns the Y coordinate where a character should stand.
func (ms *MapSpec) FootholdPoint(point types.Point[int16]) *types.Point[int16] {
	foothold, ok := ms.Footholds.Find(point)
	if !ok {
		return nil
	}

	top := foothold.Y1
	if foothold.X1 != foothold.X2 && foothold.Y1 != foothold.Y2 {
		s1 := float64(math.Abs(float64(foothold.Y2 - foothold.Y1)))
		s2 := float64(math.Abs(float64(foothold.X2 - foothold.X1)))
		dx := float64(math.Abs(float64(point.X - foothold.X1)))

		alpha := math.Atan(s2 / s1)
		beta := math.Atan(s1 / s2)
		offset := math.Cos(alpha) * (dx / math.Cos(beta))

		if foothold.Y2 < foothold.Y1 {
			top = foothold.Y1 - int16(offset)
		} else {
			top = foothold.Y1 + int16(offset)
		}
	}

	pt := types.Point[int16]{X: point.X, Y: top}
	return &pt
}

// DropPoint calculates a valid drop position for items.
// Returns the foothold position below the initial point.
func (ms *MapSpec) DropPoint(initial types.Point[int16]) (types.Point[int16], bool) {
	highest := types.Point[int16]{X: initial.X, Y: initial.Y - int16(100)}
	if result := ms.FootholdPoint(highest); result != nil {
		return *result, true
	}
	return initial, false
}

// FindPortal searches for a portal by name.
func (spec *MapSpec) FindPortal(name string) (*Portal, bool) {
	for _, portal := range spec.Portals {
		if portal.Name == name {
			return &portal, true
		}
	}

	return nil, false
}
