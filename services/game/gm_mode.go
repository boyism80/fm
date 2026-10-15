package game

import "github.com/boyism80/fm/maple/constant"

type GMMode struct {
	Hidden      bool
	InstantKill bool
	PlayerMode  bool
	TimerLimit  uint32
}

func (ch *Character) ActsAsGM() bool {
	return ch.HasRoleAtLeast(constant.RoleAdmin) && !ch.GM.PlayerMode
}
