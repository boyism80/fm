package response

import "github.com/boyism80/fm/stream"

type FaceExpression struct {
	CharacterID uint32
	Expression  uint32
}

func (*FaceExpression) Opcode() uint16 { return 0x8A }
func (p *FaceExpression) Serialize(w *stream.StreamWriter) error {
	w.WriteU32(p.CharacterID)
	w.WriteU32(p.Expression)
	w.Write32(-1)
	w.WriteU8(0)
	return nil
}
func (p *FaceExpression) Deserialize(r *stream.StreamReader) {
	p.CharacterID = r.ReadU32()
	p.Expression = r.ReadU32()
	r.Read32()
	r.ReadU8()
}
