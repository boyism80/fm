package response

import (
	"github.com/boyism80/fm/stream"
)

type Transfer struct {
	IP          string
	Port        uint16
	CharacterId uint32
}

func (a *Transfer) Opcode() uint16 {
	return 0x04
}

func (a *Transfer) Serialize(writer *stream.StreamWriter) error {
	writer.WriteU16(0)
	writer.WriteIPAddress(a.IP)
	writer.WriteU16(a.Port)
	writer.WriteU32(a.CharacterId)
	writer.WriteU8(0)
	writer.WriteU32(0)
	return nil
}

func (a *Transfer) Deserialize(reader *stream.StreamReader) {
	reader.ReadU16()
	a.IP = reader.ReadIPAddress()
	a.Port = reader.ReadU16()
	a.CharacterId = reader.ReadU32()
	reader.ReadU8()
	reader.ReadU32()
}
