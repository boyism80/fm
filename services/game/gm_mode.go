package game

import "github.com/boyism80/fm/maple/constant"

type GMMode struct {
	hidden      bool
	instantKill bool
	playerMode  bool
	timerLimit  uint32
}

func (ch *Character) ActsAsGM() bool {
	return ch.HasRoleAtLeast(constant.RoleAdmin) && !ch.gm.playerMode
}
