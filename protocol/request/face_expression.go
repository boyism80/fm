package request

import "github.com/boyism80/fm/stream"

type FaceExpression struct{ Expression uint32 }

func (*FaceExpression) Opcode() byte { return 0x22 }
func (p *FaceExpression) Serialize(w *stream.StreamWriter) error {
	w.WriteU32(p.Expression)
	return nil
}
func (p *FaceExpression) Deserialize(r *stream.StreamReader) { p.Expression = r.ReadU32() }
