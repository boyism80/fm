// Package data provides MapleStory game data specifications and types.
// This file contains string resource specifications for localized text.
package data

// StringSpec contains localized string data for items, skills, etc.
type StringSpec struct {
	ID   uint32 // Resource identifier
	Name string // Display name
	Desc string // Description text
}
