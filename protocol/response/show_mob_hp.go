package response

import "github.com/boyism80/fm/stream"

type ShowMobHp struct {
	OID        uint32
	Percentage uint8
}

func (m *ShowMobHp) Opcode() uint16 {
	return 0xB6
}

func (m *ShowMobHp) Serialize(sw *stream.StreamWriter) error {
	sw.WriteU32(m.OID)
	sw.WriteU8(m.Percentage)
	return nil
}

func (m *ShowMobHp) Deserialize(sr *stream.StreamReader) {
}
