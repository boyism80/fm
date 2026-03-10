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

func (u *UseChair) Deserialize(reader *stream.StreamReader) error {
	var err error
	if u.ItemID, err = reader.ReadU32(); err != nil {
		return err
	}
	return nil
}
