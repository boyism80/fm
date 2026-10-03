package response

import "github.com/boyism80/fm/stream"

type ServerChannel struct {
	ChannelID uint16
	Name      string
	Load      uint32
}

type ServerList struct {
	ServerId     uint8
	Channels     []ServerChannel
	WorldName    string
	Flag         byte
	EventMessage string
}

func (s *ServerList) Opcode() uint16 {
	return 0x02
}

func (s *ServerList) Serialize(writer *stream.StreamWriter) error {
	writer.WriteU8(s.ServerId)
	writer.WriteStr16(s.WorldName)
	writer.WriteU8(s.Flag)
	writer.WriteStr16(s.EventMessage)
	writer.WriteU16(100)
	writer.WriteU16(100)
	writer.WriteU8(uint8(len(s.Channels)))

	for _, ch := range s.Channels {
		writer.WriteStr16(ch.Name)
		writer.WriteU32(ch.Load)
		writer.WriteU8(s.ServerId)
		writer.WriteU16(ch.ChannelID)
	}

	writer.WriteU16(1)
	writer.WriteU16(400)
	writer.WriteU16(300)
	writer.WriteStr16("unknown string value")

	return nil
}

func (s *ServerList) Deserialize(reader *stream.StreamReader) {
	s.ServerId = reader.ReadU8()
	s.WorldName = reader.ReadStr16()
	s.Flag = reader.ReadU8()
	s.EventMessage = reader.ReadStr16()
	reader.ReadU16()
	reader.ReadU16()
	n := int(reader.ReadU8())
	if n > 0 {
		s.Channels = make([]ServerChannel, n)
	}
	for i := 0; i < n; i++ {
		s.Channels[i].Name = reader.ReadStr16()
		s.Channels[i].Load = reader.ReadU32()
		reader.ReadU8()
		s.Channels[i].ChannelID = reader.ReadU16()
	}
	reader.ReadU16()
	reader.ReadU16()
	reader.ReadU16()
	reader.ReadStr16()
}
