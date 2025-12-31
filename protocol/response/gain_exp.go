package response

import "github.com/boyism80/fm/stream"

type GainExp struct {
	Gain     uint32
	White    bool
	PartyInc int32
}

func (m *GainExp) Opcode() uint16 {
	return 0x1C
}

func (m *GainExp) Serialize(sw *stream.StreamWriter) error {
	if err := sw.WriteU8(3); err != nil {
		return err
	}
	if err := sw.WriteBoolean(m.White); err != nil {
		return err
	}
	if err := sw.WriteU32(m.Gain); err != nil {
		return err
	}
	if err := sw.Write(make([]byte, 19)); err != nil {
		return err
	}
	return nil
}

func (m *GainExp) Deserialize(sr *stream.StreamReader) error {
	return nil
}
