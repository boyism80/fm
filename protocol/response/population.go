package response

import (
	pconst "github.com/boyism80/fm/protocol/constant"
	"github.com/boyism80/fm/stream"
)

type Population struct {
	Status pconst.PopulationResult
	Mode   pconst.PopulationMode
	Name   string
	Value  int32
}

func (p *Population) Opcode() uint16 {
	return 0x1B
}

func (p *Population) Serialize(w *stream.StreamWriter) error {
	w.WriteU8(uint8(p.Status))
	switch p.Status {
	case pconst.PopulationResultGiven:
		w.WriteStr16(p.Name)
		w.WriteU8(uint8(p.Mode))
		w.Write32(p.Value)
	case pconst.PopulationResultReceived:
		w.WriteStr16(p.Name)
		w.WriteU8(uint8(p.Mode))
	}
	return nil
}

func (p *Population) Deserialize(r *stream.StreamReader) {
	p.Status = pconst.PopulationResult(r.ReadU8())
	switch p.Status {
	case pconst.PopulationResultGiven:
		p.Name = r.ReadStr16()
		p.Mode = pconst.PopulationMode(r.ReadU8())
		p.Value = r.Read32()
	case pconst.PopulationResultReceived:
		p.Name = r.ReadStr16()
		p.Mode = pconst.PopulationMode(r.ReadU8())
	}
}
