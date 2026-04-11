package response

import (
	"time"

	"github.com/boyism80/fm/protocol/dto"
	"github.com/boyism80/fm/stream"
)

type Warp struct {
	Character *dto.Character
	Channel   uint32
}

func (p *Warp) Serialize(writer *stream.StreamWriter) error {
	writer.WriteU32(p.Channel)
	writer.WriteU16(2)
	writer.WriteU16(0)
	writer.WriteU32(p.Character.Map)
	writer.WriteU8(p.Character.SpawnPoint)
	writer.WriteU16(p.Character.Hp)
	writer.WriteDateTime(time.Now())
	return nil
}

func (a *Warp) Deserialize(reader *stream.StreamReader) {
}

func (w *Warp) Opcode() uint16 {
	return 0x55 // Warp opcode (same as Login)
}
