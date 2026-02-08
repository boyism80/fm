package request

import (
	"github.com/boyism80/fm/stream"
)

// HealOverTime represents the HEAL_OVER_TIME packet (0x48)
type HealOverTime struct {
	Tick   uint32
	HealHP uint16
	HealMP uint16
	PRate  uint8
}

func (p *HealOverTime) Opcode() uint16 {
	return 0x48
}

func (p *HealOverTime) Serialize(writer *stream.StreamWriter) error {
	return nil
}

func (p *HealOverTime) Deserialize(reader *stream.StreamReader) error {
	var err error

	if reader.Remaining() >= 8 {
		if err = reader.Skip(4); err != nil {
			return err
		}
	}

	if p.HealHP, err = reader.ReadU16(); err != nil {
		return err
	}
	if p.HealMP, err = reader.ReadU16(); err != nil {
		return err
	}
	if p.PRate, err = reader.ReadU8(); err != nil {
		return err
	}
	return nil
}
