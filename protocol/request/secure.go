package request

import (
	"fmt"

	"github.com/boyism80/fm/stream"
)

// Secure represents client packet 0x0C. Structure [1][4][4]: Type, ParamA, ParamB (int32 LE).
type Secure struct {
	Type    byte
	ParamA  int32
	ParamB  int32
	Payload []byte
}

func (p *Secure) Opcode() uint16 {
	return 0x0C
}

func (p *Secure) Serialize(writer *stream.StreamWriter) error {
	return nil
}

func (p *Secure) Deserialize(reader *stream.StreamReader) error {
	n := reader.Remaining()
	if n < 9 {
		if n > 0 {
			var err error
			p.Payload, err = reader.Read(n)
			if err != nil {
				return fmt.Errorf("read payload: %w", err)
			}
		}
		return nil
	}
	t, err := reader.ReadU8()
	if err != nil {
		return err
	}
	p.Type = t
	a, err := reader.Read32()
	if err != nil {
		return err
	}
	p.ParamA = a
	b, err := reader.Read32()
	if err != nil {
		return err
	}
	p.ParamB = b
	remain := reader.Remaining()
	if remain > 0 {
		p.Payload, err = reader.Read(remain)
		if err != nil {
			return fmt.Errorf("read remainder: %w", err)
		}
	}
	return nil
}
