package resp

import (
	"github.com/boyism80/fm/common/stream"
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
	err := writer.WriteU16(0)
	if err != nil {
		return err
	}
	err = writer.WriteIPAddress(a.IP)
	if err != nil {
		return err
	}
	err = writer.WriteU16(a.Port)
	if err != nil {
		return err
	}
	err = writer.WriteU32(a.CharacterId)
	if err != nil {
		return err
	}
	err = writer.WriteU8(0)
	if err != nil {
		return err
	}
	err = writer.WriteU32(0)
	if err != nil {
		return err
	}
	return nil
}

func (a *Transfer) Deserialize(reader *stream.StreamReader) error {
	return nil
}
