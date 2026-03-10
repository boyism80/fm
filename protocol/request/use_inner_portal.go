package request

import (
	"github.com/boyism80/fm/stream"
)

type UseInnerPortal struct {
	Mode       uint8
	PortalName string
	ToX        int16
	ToY        int16
	FromX      int16
	FromY      int16
}

func (p *UseInnerPortal) Serialize(writer *stream.StreamWriter) error {
	return nil
}

func (p *UseInnerPortal) Deserialize(reader *stream.StreamReader) error {
	var err error
	if p.Mode, err = reader.ReadU8(); err != nil {
		return err
	}
	if p.PortalName, err = reader.ReadStr16(); err != nil {
		return err
	}
	if p.ToX, err = reader.Read16(); err != nil {
		return err
	}
	if p.ToY, err = reader.Read16(); err != nil {
		return err
	}

	if reader.Remaining() >= 4 {
		if p.FromX, err = reader.Read16(); err != nil {
			return err
		}
		if p.FromY, err = reader.Read16(); err != nil {
			return err
		}
	}
	return nil
}
