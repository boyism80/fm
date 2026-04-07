package response

import (
	"github.com/boyism80/fm/stream"
)

type Transfer struct {
	IP          string
	Port        uint16
	CharacterId uint32
}

// Opcode returns the packet opcode for Transfer
func (a *Transfer) Opcode() uint16 {
	return 0x04
}

func (a *Transfer) Serialize(writer *stream.StreamWriter) error {
	_ = writer.WriteU16(0)
	_ = writer.WriteIPAddress(a.IP)
	_ = writer.WriteU16(a.Port)
	_ = writer.WriteU32(a.CharacterId)
	_ = writer.WriteU8(0)
	_ = writer.WriteU32(0)
	return nil
}

func (a *Transfer) Deserialize(reader *stream.StreamReader) {
}
