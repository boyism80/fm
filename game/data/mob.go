// Package data provides MapleStory game data specifications and types.
// This file contains monster-related data structures.
package data

// MobSpec contains all statistical data for a MapleStory monster.
type MobSpec struct {
	ID         uint32  // Monster identifier
	BodyAttack int     // Body contact damage
	Level      uint8   // Monster level
	MaxHP      int     // Maximum health points
	MaxMP      int     // Maximum magic points
	Speed      int16   // Movement speed
	PADamage   int     // Physical attack damage
	PDDamage   int     // Physical defense damage
	MADamage   int     // Magic attack damage
	MDDamage   int     // Magic defense damage
	ACC        int     // Accuracy rating
	EVA        int     // Evasion rating
	EXP        uint32  // Experience points given when defeated
	Undead     bool    // Undead monster flag
	Pushed     bool    // Can be pushed by attacks
	FS         float32 // Flying speed for aerial monsters
	SummonType uint8   // Summon type identifier
	MobType    uint8   // Monster category type
}
