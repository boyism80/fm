package entity

import (
	"github.com/boyism80/fm/common/stream"
	"github.com/boyism80/fm/game/data"
)

type Mob struct {
	Life
	Spec     *data.MobSpec
	Foothold int16
}

func (m *Mob) Serialize(writer *stream.StreamWriter) error {
	writer.WriteU32(0)
	return nil
}
