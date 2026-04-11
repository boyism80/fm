package request

import (
	"github.com/boyism80/fm/stream"
)

type UseChair struct {
	ItemID uint32
}

func (u *UseChair) Serialize(writer *stream.StreamWriter) error {
	return nil
}

func (u *UseChair) Deserialize(reader *stream.StreamReader) {
	u.ItemID = reader.ReadU32()
}
