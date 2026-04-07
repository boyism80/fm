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
	sw.WriteU8(3)
	sw.WriteBoolean(m.White)
	sw.WriteU32(m.Gain)
	sw.Write(make([]byte, 19))
	return nil
}

func (m *GainExp) Deserialize(sr *stream.StreamReader) {
}
