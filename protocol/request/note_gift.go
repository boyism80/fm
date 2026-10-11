package request

import "github.com/boyism80/fm/stream"

type NoteGift struct {
}

func (*NoteGift) Opcode() byte { return 0x6E }

func (p *NoteGift) Serialize(writer *stream.StreamWriter) error {
	return nil
}

func (p *NoteGift) Deserialize(reader *stream.StreamReader) {
}
