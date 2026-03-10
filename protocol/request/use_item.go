package request

import (
	"github.com/boyism80/fm/stream"
)

// UseItem represents the USE_ITEM packet (0x37)
type UseItem struct {
	Tick   uint32
	Slot   uint16
	ItemID uint32
}

func (p *UseItem) Serialize(writer *stream.StreamWriter) error {
	return nil
}

func (p *UseItem) Deserialize(reader *stream.StreamReader) error {
	var err error

	if p.Tick, err = reader.ReadU32(); err != nil {
		return err
	}
	if p.Slot, err = reader.ReadU16(); err != nil {
		return err
	}
	if p.ItemID, err = reader.ReadU32(); err != nil {
		return err
	}

	return nil
}
