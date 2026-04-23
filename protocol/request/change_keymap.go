package request

import (
	"github.com/boyism80/fm/stream"
)

type ChangeKeymap struct {
	Skip       int32
	NumChanges int32
	Changes    []KeymapChange

	Type int32
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
	available := reader.Remaining()
	if available > 8 {
		c.Skip = reader.Read32()
		c.NumChanges = reader.Read32()
		c.Changes = make([]KeymapChange, c.NumChanges)
		for i := int32(0); i < c.NumChanges; i++ {
			c.Changes[i].Key = reader.Read32()
			c.Changes[i].Type = reader.ReadU8()
			c.Changes[i].Action = reader.Read32()
		}
	} else {
		c.Type = reader.Read32()
		c.Data = reader.Read32()
	}
}
