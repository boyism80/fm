// Package data provides MapleStory game data specifications and types.
// This file contains monster drop table specifications.
package data

// DropSpec defines what items/money a monster can drop when defeated.
type DropSpec struct {
	Mob   uint32  // Monster ID that drops this item
	Item  uint32  // Item ID to drop (0 for money only)
	Money uint32  // Meso amount to drop (0 for item only)
	Prob  float32 // Drop probability (0.0-1.0)
	Min   uint16  // Minimum quantity to drop
	Max   uint16  // Maximum quantity to drop
}
