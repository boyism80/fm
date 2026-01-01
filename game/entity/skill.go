package entity

import (
	"time"

	"github.com/boyism80/fm/game/wz"
)

type SkillEntry struct {
	*wz.Skill
	SkillLevel  int
	MasterLevel int
	Expiration  time.Time
}
