package response

import (
	"github.com/boyism80/fm/stream"
)

type Authenticate struct {
	AccountId     uint32
	Gender        uint8
	Role          uint8
	AccountName   string
	IsChatBlocked bool
	ChatBlockTime uint64
}

func (a *Authenticate) Opcode() uint16 {
	return 0x00
}

func (a *Authenticate) Serialize(writer *stream.StreamWriter) error {
	writer.WriteU8(0)
	writer.WriteU32(a.AccountId)
	writer.WriteU8(a.Gender)
	writer.WriteU8(a.Role)
	writer.WriteU8(0)
	writer.WriteStr16(a.AccountName)
	writer.WriteU32(0)
	writer.WriteU8(0)
	writer.WriteU8(0)
	writer.WriteBoolean(a.IsChatBlocked)
	writer.WriteU64(a.ChatBlockTime)
	writer.WriteStr16("")
	writer.WriteStr16("")
	return nil
}

func (a *Authenticate) Deserialize(reader *stream.StreamReader) {
	reader.ReadU8()
	a.AccountId = reader.ReadU32()
	a.Gender = reader.ReadU8()
	a.Role = reader.ReadU8()
	reader.ReadU8()
	a.AccountName = reader.ReadStr16()
	reader.ReadU32()
	reader.ReadU8()
	reader.ReadU8()
	a.IsChatBlocked = reader.ReadBool()
	a.ChatBlockTime = reader.ReadU64()
	reader.ReadStr16()
	reader.ReadStr16()
}
