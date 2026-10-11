package response

import (
	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
	"github.com/boyism80/fm/stream"
)

type FamilyResult struct {
	Code internal.FamilyResult
}

func (p *FamilyResult) Opcode() uint16 {
	return 0x4B
}

func (p *FamilyResult) Serialize(w *stream.StreamWriter) error {
	w.WriteU32(uint32(p.Code))
	return nil
}

func (p *FamilyResult) Deserialize(reader *stream.StreamReader) {
	p.Code = internal.FamilyResult(reader.ReadU32())
}
