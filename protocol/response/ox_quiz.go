package response

import "github.com/boyism80/fm/stream"

// OXQuiz shows a question of Etc.wz/OXQuiz.img, or its answer when Ask is false.
type OXQuiz struct {
	Ask bool
	Set uint8
	ID  uint16
}

func (p *OXQuiz) Opcode() uint16 {
	return 0x63
}

func (p *OXQuiz) Serialize(sw *stream.StreamWriter) error {
	sw.WriteBoolean(p.Ask)
	sw.WriteU8(p.Set)
	sw.WriteU16(p.ID)
	return nil
}

func (p *OXQuiz) Deserialize(reader *stream.StreamReader) {
	p.Ask = reader.ReadBool()
	p.Set = reader.ReadU8()
	p.ID = reader.ReadU16()
}
