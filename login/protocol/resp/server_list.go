package resp

import (
	"fmt"

	"github.com/boyism80/fm/common/stream"
)

type ServerList struct {
	ServerId     uint8
	ChannelSize  uint8
	WorldName    string
	Flag         byte
	EventMessage string
}

func (s *ServerList) Serialize(writer *stream.StreamWriter) error {
	writer.WriteU16(0x02)
	writer.WriteU8(s.ServerId)
	writer.WriteStr16(s.WorldName)
	writer.WriteU8(s.Flag)
	writer.WriteStr16(s.EventMessage)
	writer.WriteU16(100)
	writer.WriteU16(100)
	writer.WriteU8(s.ChannelSize)

	for i := 1; i <= int(s.ChannelSize); i++ {
		channelName := fmt.Sprintf("%s-%d", s.WorldName, i)
		writer.WriteStr16(channelName)
		writer.WriteU32(1200) // 채널 유저 현황인듯
		writer.WriteU8(s.ServerId)
		writer.WriteU16(uint16(i - 1))
	}

	writer.WriteU16(1)
	writer.WriteU16(400)
	writer.WriteU16(300)
	writer.WriteStr16("unknown string value")

	return nil
}

func (s *ServerList) Deserialize(reader *stream.StreamReader) error {
	return nil
}
