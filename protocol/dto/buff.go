package dto

import "github.com/boyism80/fm/game/constant"

// BuffEntry is one buff type and value for UpdateBuff and UpdateRemoteBuff.
// GiveBuff (self): always writes Value (short), buffid, bufflength per stat.
type BuffEntry struct {
	Buff  constant.BuffFlag
	Value int32
}
