package request

import (
	"github.com/boyism80/fm/stream"
)

type ChangeKeymap struct {
	Mode    int32
	Changes []KeymapChange

	Data int32
}

type KeymapChange struct {
	Key    int32
	Type   byte
	Action int32
}

func (*ChangeKeymap) Opcode() byte { return 0x71 }

func (c *ChangeKeymap) Serialize(writer *stream.StreamWriter) error {
	return nil
}

func (c *ChangeKeymap) Deserialize(reader *stream.StreamReader) {
	c.Mode = reader.Read32()
	if c.Mode != 0 {
		c.Data = reader.Read32()
		return
	}

	count := reader.Read32()
	for range count {
		c.Changes = append(c.Changes, KeymapChange{
			Key:    reader.Read32(),
			Type:   reader.ReadU8(),
			Action: reader.Read32(),
		})
	}
}
