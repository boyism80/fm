package request

import "github.com/boyism80/fm/stream"

type CancelItemEffect struct {
	SourceID int32
}

func (*CancelItemEffect) Opcode() byte { return 0x38 }

func (c *CancelItemEffect) Serialize(writer *stream.StreamWriter) error {
	return nil
}

func (c *CancelItemEffect) Deserialize(reader *stream.StreamReader) {
	c.SourceID = reader.Read32()
}
