package resp

import (
	"github.com/boyism80/fm/common/stream"
)

type Authenticate struct {
	AccountId     uint32
	Gender        uint8
	Admin         bool
	AccountName   string
	IsChatBlocked bool
	ChatBlockTime uint64
}

func (a *Authenticate) Serialize(writer *stream.StreamWriter) error {
	writer.WriteU16(0x00)
	writer.WriteU8(0)
	writer.WriteU32(a.AccountId)
	writer.WriteU8(a.Gender)
	writer.WriteBoolean(a.Admin)
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

func (a *Authenticate) Deserialize(reader *stream.StreamReader) error {
	return nil
}
