package request

import "github.com/boyism80/fm/stream"

type DeleteCharacter struct {
	ID uint32
}

func (*DeleteCharacter) Opcode() byte { return 0x09 }

func (a *DeleteCharacter) Serialize(writer *stream.StreamWriter) error {
	writer.Write(make([]byte, 5))
	writer.WriteU32(a.ID)
	return nil
}

func (a *DeleteCharacter) Deserialize(reader *stream.StreamReader) {
	reader.Read(5)
	a.ID = reader.ReadU32()
}
