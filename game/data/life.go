// Package data provides MapleStory game data specifications and types.
// This file contains life entity spawn specifications (NPCs and monsters).
package data

import (
	"time"

	"github.com/boyism80/fm/core/types"
)

// FacingDirectionType represents the direction a life entity is facing.
type FacingDirectionType uint8

const (
	FACING_DIRECTION_RIGHT FacingDirectionType = iota // Facing right
	FACING_DIRECTION_LEFT                             // Facing left
)

// Spawn is a marker interface for spawn-related types.
type Spawn interface {
}

// SpawnSpec contains spawn point data for life entities (NPCs/monsters).
type SpawnSpec struct {
	ID                 uint32               // Life entity ID
	Position           types.Vector2[int16] // Spawn position
	RenderX0, RenderX1 int16                // Render boundaries
	CollisionY         int16                // Collision Y coordinate
	Hide               bool                 // Hidden spawn flag
	UseDay, UseNight   bool                 // Day/night visibility
	Foothold           int16                // Foothold ID for positioning
	FacingDirection    FacingDirectionType  // Initial facing direction
	MobTime            time.Duration        // Respawn time for mobs
	Info               uint8                // Additional spawn info
	LimitedName        string               // Limited spawn name
	NoFoothold         bool                 // Ignore foothold positioning
}
