package request

import (
	pconst "github.com/boyism80/fm/protocol/constant"
	"github.com/boyism80/fm/stream"
)

type GivePopulation struct {
	TargetID uint32
	Mode     pconst.PopulationMode
}

func (*GivePopulation) Opcode() byte { return 0x4E }

func (p *GivePopulation) Serialize(w *stream.StreamWriter) error {
	w.WriteU32(p.TargetID)
	w.WriteU8(uint8(p.Mode))
	return nil
}

func (p *GivePopulation) Deserialize(r *stream.StreamReader) {
	p.TargetID = r.ReadU32()
	p.Mode = pconst.PopulationMode(r.ReadU8())
}
