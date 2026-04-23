package request

import (
	"github.com/boyism80/fm/stream"
)

type UseChair struct {
	ItemID uint32
}

func (*UseChair) Opcode() byte { return 0x1A }

func (u *UseChair) Serialize(writer *stream.StreamWriter) error {
	return nil
}

func (u *UseChair) Deserialize(reader *stream.StreamReader) {
	u.ItemID = reader.ReadU32()
}
