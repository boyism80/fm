package resp

import (
	"time"

	"github.com/boyism80/fm/common/stream"
	"github.com/boyism80/fm/game/entity"
)

type Warp struct {
	Character  *entity.Character
	Channel    uint32
	SpawnPoint uint8
}

func (p *Warp) Serialize(writer *stream.StreamWriter) error {
	writer.WriteU16(0x55)
	writer.WriteU32(p.Channel)
	writer.WriteU16(2)
	writer.WriteU16(0)
	writer.WriteU32(p.Character.Map)
	writer.WriteU8(p.SpawnPoint)
	writer.WriteU16(p.Character.Hp)
	writer.WriteDateTime(time.Now())
	return nil
}

func (a *Warp) Deserialize(reader *stream.StreamReader) error {
	return nil
}
